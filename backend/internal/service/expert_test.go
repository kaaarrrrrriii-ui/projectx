package service

import (
	"context"
	"testing"
	"time"

	"example.com/german/backend/internal/models"
	"example.com/german/backend/internal/repos"
)

type stubExpertRepository struct {
	filter             repos.ExpertTicketFilter
	openedTrack        string
	openedWorker       int64
	createdRequestType string
	createdReason      string
	createdWorker      int64
	analyticsWorker    int64
	analyticsStart     time.Time
	analyticsEnd       time.Time
}

func (*stubExpertRepository) Profile(context.Context, int64) (repos.ExpertProfileRecord, error) {
	return repos.ExpertProfileRecord{}, nil
}
func (*stubExpertRepository) Dashboard(context.Context, int64) (repos.ExpertDashboardRecord, error) {
	return repos.ExpertDashboardRecord{}, nil
}
func (stub *stubExpertRepository) ListTickets(_ context.Context, filter repos.ExpertTicketFilter) (repos.ExpertTicketPageRecord, error) {
	stub.filter = filter
	return repos.ExpertTicketPageRecord{}, nil
}
func (stub *stubExpertRepository) OpenTicket(_ context.Context, track string, worker int64, _ time.Time) (models.TicketStatus, error) {
	stub.openedTrack, stub.openedWorker = track, worker
	return models.TicketStatusInProgress, nil
}
func (*stubExpertRepository) GetTicket(context.Context, string, int64) (repos.ExpertTicketDetailRecord, error) {
	return repos.ExpertTicketDetailRecord{}, nil
}
func (*stubExpertRepository) AddNote(context.Context, string, int64, string, string, time.Time) (repos.ExpertNoteRecord, error) {
	return repos.ExpertNoteRecord{}, nil
}
func (*stubExpertRepository) AddAnswer(context.Context, string, int64, string, time.Time) (repos.ExpertMessageRecord, error) {
	return repos.ExpertMessageRecord{}, nil
}
func (stub *stubExpertRepository) CreateWorkerRequest(_ context.Context, track string, worker int64, requestType, reason string, createdAt time.Time) (repos.WorkerRequestRecord, error) {
	stub.createdRequestType, stub.createdReason, stub.createdWorker = requestType, reason, worker
	return repos.WorkerRequestRecord{ID: 5, TrackID: track, RequestType: requestType, Reason: reason, Status: "sent", CreatedAt: createdAt, CreatedBy: worker}, nil
}
func (*stubExpertRepository) ListWorkerRequests(context.Context, int64, string, int, int) (repos.WorkerRequestPageRecord, error) {
	return repos.WorkerRequestPageRecord{}, nil
}
func (*stubExpertRepository) CanAccessTicket(context.Context, string, int64) error { return nil }
func (stub *stubExpertRepository) Analytics(_ context.Context, worker int64, start, end time.Time) (repos.AnalyticsRecord, error) {
	stub.analyticsWorker, stub.analyticsStart, stub.analyticsEnd = worker, start, end
	return repos.AnalyticsRecord{}, nil
}
func (*stubExpertRepository) ReportTickets(context.Context, int64, time.Time, time.Time) ([]repos.ReportTicketRecord, error) {
	return nil, nil
}

func newExpertTestService(t *testing.T, repository *stubExpertRepository) *ExpertService {
	t.Helper()
	result, err := NewExpertService(repository)
	if err != nil {
		t.Fatalf("NewExpertService() error = %v", err)
	}
	return result
}

func TestExpertTicketFilterIsScopedToWorker(t *testing.T) {
	t.Parallel()
	repository := &stubExpertRepository{}
	expert := newExpertTestService(t, repository)
	_, err := expert.ListTickets(context.Background(), 7, ExpertTicketListRequest{Queue: "returned", Priority: "urgent", ApplicantType: "parent", Page: "2", Limit: "10"})
	if err != nil {
		t.Fatalf("ListTickets() error = %v", err)
	}
	if repository.filter.WorkerID != 7 || repository.filter.Queue != "returned" || repository.filter.Offset != 10 || repository.filter.Priority != models.TicketPriorityUrgent || repository.filter.ApplicantType != models.ApplicantTypeParent {
		t.Fatalf("expert filter = %+v", repository.filter)
	}
}

func TestExpertOpenNormalizesTrackAndStartsTicket(t *testing.T) {
	t.Parallel()
	repository := &stubExpertRepository{}
	expert := newExpertTestService(t, repository)
	response, err := expert.OpenTicket(context.Background(), " отк-abcd-2345 ", 7)
	if err != nil {
		t.Fatalf("OpenTicket() error = %v", err)
	}
	if repository.openedTrack != "ОТК-ABCD-2345" || repository.openedWorker != 7 || response.Status != "in_progress" {
		t.Fatalf("open result = (%q, %d, %+v)", repository.openedTrack, repository.openedWorker, response)
	}
}

func TestExpertCreatesWorkerRequestWithAgreedTypes(t *testing.T) {
	t.Parallel()
	repository := &stubExpertRepository{}
	expert := newExpertTestService(t, repository)
	response, err := expert.CreateWorkerRequest(context.Background(), "ОТК-ABCD-2345", 7, "add_coworker", " Нужен коллега ")
	if err != nil {
		t.Fatalf("CreateWorkerRequest() error = %v", err)
	}
	if repository.createdWorker != 7 || repository.createdRequestType != "add_coworker" || repository.createdReason != "Нужен коллега" || response.Request.Status != "sent" {
		t.Fatalf("worker request = %+v, response = %+v", repository, response)
	}
	if _, err := expert.CreateWorkerRequest(context.Background(), "ОТК-ABCD-2345", 7, "consultation", "reason"); err != ErrInvalidRequestType {
		t.Fatalf("invalid request type error = %v", err)
	}
}

func TestExpertAnalyticsUsesPersonalWorkerAndYekaterinburgDates(t *testing.T) {
	t.Parallel()
	repository := &stubExpertRepository{}
	expert := newExpertTestService(t, repository)
	if _, err := expert.Analytics(context.Background(), 7, "2026-09-01", "2026-09-10"); err != nil {
		t.Fatalf("Analytics() error = %v", err)
	}
	if repository.analyticsWorker != 7 || repository.analyticsStart.Format(time.RFC3339) != "2026-08-31T19:00:00Z" || repository.analyticsEnd.Format(time.RFC3339) != "2026-09-10T19:00:00Z" {
		t.Fatalf("analytics scope = worker %d, %s - %s", repository.analyticsWorker, repository.analyticsStart, repository.analyticsEnd)
	}
}

func TestRequireExpert(t *testing.T) {
	t.Parallel()
	if err := RequireExpert(AuthUser{Role: "expert"}); err != nil {
		t.Fatalf("RequireExpert(expert) error = %v", err)
	}
	if err := RequireExpert(AuthUser{Role: "operator"}); err != ErrForbidden {
		t.Fatalf("RequireExpert(operator) error = %v", err)
	}
}
