package service

import (
	"archive/zip"
	"bytes"
	"context"
	"database/sql"
	"encoding/csv"
	"errors"
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
	onlyAvailable   bool
	listFilter      repos.OperatorTicketFilter
	dashboard       repos.OperatorDashboardRecord
	detail          repos.OperatorTicketDetailRecord
	update          repos.OperatorTicketUpdate
	closedMessage   string
}

func (repository *stubOperatorRepository) EligibleWorkers(_ context.Context, _ string, _ string, _ int64, onlyAvailable bool) (repos.EligibleWorkersRecord, error) {
	repository.onlyAvailable = onlyAvailable
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

func (repository *stubOperatorRepository) ListTickets(_ context.Context, filter repos.OperatorTicketFilter) (repos.OperatorTicketPageRecord, error) {
	repository.listFilter = filter
	return repos.OperatorTicketPageRecord{}, nil
}

func (repository *stubOperatorRepository) Dashboard(context.Context, time.Time, time.Time) (repos.OperatorDashboardRecord, error) {
	return repository.dashboard, nil
}

func (repository *stubOperatorRepository) GetOperatorTicket(context.Context, string) (repos.OperatorTicketDetailRecord, error) {
	return repository.detail, nil
}

func (repository *stubOperatorRepository) UpdateOperatorTicket(_ context.Context, _ string, _ int64, update repos.OperatorTicketUpdate, _ time.Time) error {
	repository.update = update
	return nil
}

func (repository *stubOperatorRepository) CloseOperatorTicket(_ context.Context, _ string, _ int64, message string, closedAt time.Time) (repos.OperatorCloseRecord, error) {
	repository.closedMessage = message
	return repos.OperatorCloseRecord{Status: models.TicketStatusCompleted, ClosedAt: closedAt, Message: message}, nil
}

func (*stubOperatorRepository) CanAccessOperatorAttachment(context.Context, string, int64) error {
	return nil
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

func TestParseTicketFilterSupportsNewSortAndLowPriority(t *testing.T) {
	t.Parallel()
	filter, _, err := parseTicketFilter(TicketListRequest{Queue: "new", Priority: "low", Sort: "waiting_desc"})
	if err != nil || filter.Priority != models.TicketPriorityLow || filter.Sort != "waiting_desc" {
		t.Fatalf("filter = %+v, error = %v", filter, err)
	}
	if _, _, err := parseTicketFilter(TicketListRequest{Queue: "new", Sort: "drop table"}); !errors.Is(err, ErrInvalidFilter) {
		t.Fatalf("invalid sort error = %v", err)
	}
}

func TestEligibleWorkersDefaultsToAvailableOnly(t *testing.T) {
	t.Parallel()
	repository := &stubOperatorRepository{}
	operator := newOperatorTestService(t, repository)
	if _, err := operator.EligibleWorkers(context.Background(), "ОТК-ABCD-2345", "", "", ""); err != nil || !repository.onlyAvailable {
		t.Fatalf("only_available = %v, error = %v", repository.onlyAvailable, err)
	}
	if _, err := operator.EligibleWorkers(context.Background(), "ОТК-ABCD-2345", "", "", "false"); err != nil || repository.onlyAvailable {
		t.Fatalf("only_available=false = %v, error = %v", repository.onlyAvailable, err)
	}
}

func TestOperatorDashboardAndClose(t *testing.T) {
	t.Parallel()
	repository := &stubOperatorRepository{dashboard: repos.OperatorDashboardRecord{NewCount: 3, CrisisCount: 1}}
	operator := newOperatorTestService(t, repository)
	dashboard, err := operator.Dashboard(context.Background())
	if err != nil || dashboard.NewCount != 3 || dashboard.CrisisCount != 1 {
		t.Fatalf("dashboard = %+v, error = %v", dashboard, err)
	}
	closed, err := operator.CloseTicket(context.Background(), "ОТК-ABCD-2345", 4, "  Ответ оператора  ")
	if err != nil || closed.Status != "completed" || repository.closedMessage != "Ответ оператора" {
		t.Fatalf("close = %+v, message = %q, error = %v", closed, repository.closedMessage, err)
	}
}

func TestOperatorUpdateTicketMapsAtomicPatch(t *testing.T) {
	t.Parallel()
	repository := &stubOperatorRepository{detail: repos.OperatorTicketDetailRecord{
		TrackID: "ОТК-ABCD-2345", Status: models.TicketStatusNew,
		Priority: models.TicketPriorityLow, CategoryID: 8, Category: "Категория",
		ApplicantType: models.ApplicantTypeParent,
	}}
	operator := newOperatorTestService(t, repository)
	categoryID, priority, status := int64(8), "low", "rejected"
	response, err := operator.UpdateTicket(context.Background(), "ОТК-ABCD-2345", 3, OperatorTicketUpdateRequest{
		CategoryID: &categoryID, Priority: &priority, Status: &status, Reason: "вне компетенции",
	})
	if err != nil || repository.update.Priority == nil || *repository.update.Priority != models.TicketPriorityLow || repository.update.Status == nil || *repository.update.Status != models.TicketStatusRejected {
		t.Fatalf("update = %+v, response = %+v, error = %v", repository.update, response, err)
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
