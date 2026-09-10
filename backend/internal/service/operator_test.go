package service

import (
	"archive/zip"
	"bytes"
	"context"
	"database/sql"
	"encoding/csv"
	"io"
	"strings"
	"testing"
	"time"

	"example.com/german/backend/internal/models"
	"example.com/german/backend/internal/repos"
)

type stubOperatorRepository struct {
	assignment      repos.AssignmentRecord
	assignedTrackID string
	assignedWorker  int64
	assignedActor   int64
	analytics       repos.AnalyticsRecord
	analyticsStart  time.Time
	analyticsEnd    time.Time
	reportRecords   []repos.ReportTicketRecord
	workerRequest   repos.WorkerRequestRecord
	completedID     int64
	completedBy     int64
	addedWorker     int64
}

func (repository *stubOperatorRepository) EligibleWorkers(context.Context, string, string, int64) (repos.EligibleWorkersRecord, error) {
	return repos.EligibleWorkersRecord{}, nil
}

func (repository *stubOperatorRepository) SetResponsibleWorker(_ context.Context, trackID string, workerID, actorID int64, _ time.Time) (repos.AssignmentRecord, error) {
	repository.assignedTrackID = trackID
	repository.assignedWorker = workerID
	repository.assignedActor = actorID
	return repository.assignment, nil
}

func (repository *stubOperatorRepository) AddWorker(_ context.Context, _ string, workerID, _ int64, _ time.Time) (repos.AssignmentRecord, error) {
	repository.addedWorker = workerID
	return repository.assignment, nil
}

func (repository *stubOperatorRepository) RemoveWorker(context.Context, string, int64, int64, time.Time) (repos.AssignmentRecord, error) {
	return repository.assignment, nil
}

func (repository *stubOperatorRepository) Reject(context.Context, string, int64, string, string, time.Time) (repos.RejectionRecord, error) {
	return repos.RejectionRecord{}, nil
}

func (repository *stubOperatorRepository) ListTickets(context.Context, repos.OperatorTicketFilter) (repos.OperatorTicketPageRecord, error) {
	return repos.OperatorTicketPageRecord{}, nil
}

func (repository *stubOperatorRepository) Analytics(_ context.Context, start, end time.Time) (repos.AnalyticsRecord, error) {
	repository.analyticsStart = start
	repository.analyticsEnd = end
	return repository.analytics, nil
}

func (repository *stubOperatorRepository) ReportTickets(context.Context, time.Time, time.Time) ([]repos.ReportTicketRecord, error) {
	return repository.reportRecords, nil
}

func (repository *stubOperatorRepository) ListOperatorWorkerRequests(context.Context, string, int, int) (repos.WorkerRequestPageRecord, error) {
	return repos.WorkerRequestPageRecord{Items: []repos.WorkerRequestRecord{repository.workerRequest}, Total: 1}, nil
}

func (repository *stubOperatorRepository) CompleteWorkerRequest(_ context.Context, requestID, workerID, actorID int64, _ time.Time) (repos.WorkerRequestRecord, error) {
	repository.completedID, repository.completedBy, repository.addedWorker = requestID, actorID, workerID
	result := repository.workerRequest
	result.Status = "completed"
	return result, nil
}

func newOperatorTestService(t *testing.T, repository *stubOperatorRepository) *OperatorService {
	t.Helper()
	result, err := NewOperatorService(repository)
	if err != nil {
		t.Fatalf("NewOperatorService() error = %v", err)
	}
	return result
}

func TestOperatorServiceAssignsResponsibleWorker(t *testing.T) {
	t.Parallel()
	repository := &stubOperatorRepository{assignment: repos.AssignmentRecord{WorkerID: 7, Status: models.TicketStatusAssigned}}
	operator := newOperatorTestService(t, repository)
	response, err := operator.SetResponsibleWorker(context.Background(), " отк-abcd-2345 ", 7, 3)
	if err != nil {
		t.Fatalf("SetResponsibleWorker() error = %v", err)
	}
	if repository.assignedTrackID != "ОТК-ABCD-2345" || repository.assignedWorker != 7 || repository.assignedActor != 3 {
		t.Fatalf("assignment arguments = %q, %d, %d", repository.assignedTrackID, repository.assignedWorker, repository.assignedActor)
	}
	if response.Status != "assigned" || response.WorkerID != 7 {
		t.Fatalf("assignment response = %+v", response)
	}
}

func TestOperatorCompletesAddCoworkerRequest(t *testing.T) {
	t.Parallel()
	repository := &stubOperatorRepository{
		workerRequest: repos.WorkerRequestRecord{ID: 12, TicketID: 4, TrackID: "ОТК-ABCD-2345", RequestType: "add_coworker", Status: "sent"},
	}
	operator := newOperatorTestService(t, repository)
	response, err := operator.CompleteWorkerRequest(context.Background(), 12, 9, 3)
	if err != nil {
		t.Fatalf("CompleteWorkerRequest() error = %v", err)
	}
	if repository.addedWorker != 9 || repository.completedID != 12 || repository.completedBy != 3 {
		t.Fatalf("completion calls = worker %d, request %d, actor %d", repository.addedWorker, repository.completedID, repository.completedBy)
	}
	if response.Status != "completed" || response.WorkerID != 9 {
		t.Fatalf("completion response = %+v", response)
	}
}

func TestParseTicketFilter(t *testing.T) {
	t.Parallel()
	filter, page, err := parseTicketFilter(TicketListRequest{
		Queue: "assigned", Priority: "urgent", ApplicantType: "teacher", Page: "2", Limit: "30",
	})
	if err != nil {
		t.Fatalf("parseTicketFilter() error = %v", err)
	}
	if page != 2 || filter.Limit != 30 || filter.Offset != 30 || filter.Priority != models.TicketPriorityUrgent || filter.ApplicantType != models.ApplicantTypeTeacher {
		t.Fatalf("filter = %+v, page = %d", filter, page)
	}
}

func TestOperatorAnalyticsUsesYekaterinburgDatesAndShares(t *testing.T) {
	t.Parallel()
	average := 120.0
	repository := &stubOperatorRepository{analytics: repos.AnalyticsRecord{
		Total: 4, UrgentCount: 1, ReturnedCount: 2, DurationSecondsToAccept: &average,
		Categories: []repos.NamedCountRecord{{ID: 2, Name: "Кибербуллинг", Count: 1}},
	}}
	operator := newOperatorTestService(t, repository)
	response, err := operator.Analytics(context.Background(), "2026-09-10", "2026-09-10")
	if err != nil {
		t.Fatalf("Analytics() error = %v", err)
	}
	wantStart := time.Date(2026, 9, 9, 19, 0, 0, 0, time.UTC)
	if !repository.analyticsStart.Equal(wantStart) || repository.analyticsEnd.Sub(repository.analyticsStart) != 24*time.Hour {
		t.Fatalf("analytics range = %v .. %v", repository.analyticsStart, repository.analyticsEnd)
	}
	if response.UrgentSharePercent != 25 || response.ReturnSharePercent != 50 || len(response.Categories) != 1 || response.Categories[0].Percent != 25 {
		t.Fatalf("analytics response = %+v", response)
	}
}

func TestReportsAreAnonymizedAndReadable(t *testing.T) {
	t.Parallel()
	records := []repos.ReportTicketRecord{{
		CreatedAt: time.Date(2026, 9, 10, 8, 0, 0, 0, time.UTC),
		ClosedAt:  sql.NullTime{Time: time.Date(2026, 9, 10, 9, 0, 0, 0, time.UTC), Valid: true},
		Category:  "Кибербуллинг", ApplicantType: models.ApplicantTypeSchoolchild,
		Status: models.TicketStatusCompleted, Priority: models.TicketPriorityUrgent,
	}}
	csvData, err := buildCSVReport(records)
	if err != nil {
		t.Fatalf("buildCSVReport() error = %v", err)
	}
	reader := csv.NewReader(bytes.NewReader(bytes.TrimPrefix(csvData, []byte{0xef, 0xbb, 0xbf})))
	rows, err := reader.ReadAll()
	if err != nil || len(rows) != 2 || rows[1][2] != "Кибербуллинг" {
		t.Fatalf("CSV rows = %#v, error = %v", rows, err)
	}
	if strings.Contains(string(csvData), "track_id") {
		t.Fatal("CSV report exposes track_id")
	}

	xlsxData, err := buildXLSXReport(records)
	if err != nil {
		t.Fatalf("buildXLSXReport() error = %v", err)
	}
	archive, err := zip.NewReader(bytes.NewReader(xlsxData), int64(len(xlsxData)))
	if err != nil {
		t.Fatalf("open XLSX archive: %v", err)
	}
	foundSheet := false
	for _, file := range archive.File {
		if file.Name != "xl/worksheets/sheet1.xml" {
			continue
		}
		foundSheet = true
		content, err := file.Open()
		if err != nil {
			t.Fatalf("open XLSX sheet: %v", err)
		}
		data, err := io.ReadAll(content)
		_ = content.Close()
		if err != nil || !bytes.Contains(data, []byte("Кибербуллинг")) || bytes.Contains(data, []byte("track_id")) {
			t.Fatalf("XLSX sheet is invalid or not anonymized: %v", err)
		}
	}
	if !foundSheet {
		t.Fatal("XLSX worksheet is missing")
	}
}
