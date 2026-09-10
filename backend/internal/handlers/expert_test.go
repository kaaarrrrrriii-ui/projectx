package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"example.com/german/backend/internal/service"
)

type stubExpertService struct {
	openedWorker int64
	openedTrack  string
}

func (*stubExpertService) Profile(context.Context, int64) (service.ExpertProfileResponse, error) {
	return service.ExpertProfileResponse{}, nil
}
func (*stubExpertService) Dashboard(context.Context, int64) (service.ExpertDashboardResponse, error) {
	return service.ExpertDashboardResponse{}, nil
}
func (*stubExpertService) ListTickets(context.Context, int64, service.ExpertTicketListRequest) (service.ExpertTicketPageResponse, error) {
	return service.ExpertTicketPageResponse{}, nil
}
func (stub *stubExpertService) OpenTicket(_ context.Context, track string, worker int64) (service.ExpertOpenResponse, error) {
	stub.openedTrack, stub.openedWorker = track, worker
	return service.ExpertOpenResponse{Status: "in_progress"}, nil
}
func (*stubExpertService) Ticket(context.Context, string, int64) (service.ExpertTicketResponse, error) {
	return service.ExpertTicketResponse{}, nil
}
func (*stubExpertService) AddNote(context.Context, string, service.AuthUser, string) (service.CreateExpertNoteResponse, error) {
	return service.CreateExpertNoteResponse{}, nil
}
func (*stubExpertService) AddAnswer(context.Context, string, int64, string) (service.CreateExpertAnswerResponse, error) {
	return service.CreateExpertAnswerResponse{}, nil
}
func (*stubExpertService) CreateWorkerRequest(context.Context, string, int64, string, string) (service.CreateWorkerRequestResponse, error) {
	return service.CreateWorkerRequestResponse{}, nil
}
func (*stubExpertService) ListWorkerRequests(context.Context, int64, string, string, string) (service.WorkerRequestPageResponse, error) {
	return service.WorkerRequestPageResponse{}, nil
}
func (*stubExpertService) CanAccessTicket(context.Context, string, int64) error { return nil }
func (*stubExpertService) Analytics(context.Context, int64, string, string) (service.AnalyticsResponse, error) {
	return service.AnalyticsResponse{}, nil
}
func (*stubExpertService) Report(context.Context, int64, string, string, string) (service.GeneratedReport, error) {
	return service.GeneratedReport{}, nil
}

type stubExpertAttachmentService struct{}

func (*stubExpertAttachmentService) Open(context.Context, string, int64) (service.AttachmentFile, error) {
	return service.AttachmentFile{}, nil
}

func TestExpertHandlerRejectsOperatorRole(t *testing.T) {
	t.Parallel()
	auth := &stubAuthService{user: service.AuthUser{ID: 3, Role: "operator"}}
	expert := &stubExpertService{}
	handler, _ := NewExpertHandler(expert, auth, &stubExpertAttachmentService{})
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)
	request := httptest.NewRequest(http.MethodGet, "/api/expert/dashboard", nil)
	request.Header.Set("Authorization", "Bearer token")
	response := httptest.NewRecorder()
	mux.ServeHTTP(response, request)
	if response.Code != http.StatusForbidden {
		t.Fatalf("status = %d, body = %s", response.Code, response.Body.String())
	}
}

func TestExpertHandlerUsesAuthenticatedExpertForAutomaticOpen(t *testing.T) {
	t.Parallel()
	auth := &stubAuthService{user: service.AuthUser{ID: 7, Role: "expert"}}
	expert := &stubExpertService{}
	handler, _ := NewExpertHandler(expert, auth, &stubExpertAttachmentService{})
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)
	request := httptest.NewRequest(http.MethodPost, "/api/expert/tickets/%D0%9E%D0%A2%D0%9A-ABCD-2345/open", nil)
	request.Header.Set("Authorization", "Bearer token")
	response := httptest.NewRecorder()
	mux.ServeHTTP(response, request)
	if response.Code != http.StatusOK || expert.openedWorker != 7 || expert.openedTrack != "ОТК-ABCD-2345" {
		t.Fatalf("status = %d, expert = %+v, body = %s", response.Code, expert, response.Body.String())
	}
}
