package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"example.com/german/backend/internal/service"
)

type stubAuthService struct {
	loginResponse service.LoginResponse
	loginErr      error
	user          service.AuthUser
	authErr       error
}

func (stub *stubAuthService) Login(context.Context, string, string) (service.LoginResponse, error) {
	return stub.loginResponse, stub.loginErr
}

func (stub *stubAuthService) Authenticate(context.Context, string) (service.AuthUser, error) {
	return stub.user, stub.authErr
}

type stubOperatorService struct {
	assignmentCalled bool
	workerID         int64
	actorID          int64
}

func (*stubOperatorService) EligibleWorkers(context.Context, string, string, string) (service.EligibleWorkersResponse, error) {
	return service.EligibleWorkersResponse{}, nil
}

func (stub *stubOperatorService) SetResponsibleWorker(_ context.Context, _ string, workerID, actorID int64) (service.AssignmentResponse, error) {
	stub.assignmentCalled = true
	stub.workerID = workerID
	stub.actorID = actorID
	return service.AssignmentResponse{WorkerID: workerID, Status: "assigned"}, nil
}

func (*stubOperatorService) AddWorker(context.Context, string, int64, int64) (service.AssignmentResponse, error) {
	return service.AssignmentResponse{}, nil
}

func (*stubOperatorService) RemoveWorker(context.Context, string, int64, int64) (service.AssignmentResponse, error) {
	return service.AssignmentResponse{}, nil
}

func (*stubOperatorService) Reject(context.Context, string, int64, string, string) (service.RejectionResponse, error) {
	return service.RejectionResponse{}, nil
}

func (*stubOperatorService) ListTickets(context.Context, service.TicketListRequest) (service.OperatorTicketPageResponse, error) {
	return service.OperatorTicketPageResponse{}, nil
}

func (*stubOperatorService) Analytics(context.Context, string, string) (service.AnalyticsResponse, error) {
	return service.AnalyticsResponse{}, nil
}

func (*stubOperatorService) Report(context.Context, string, string, string) (service.GeneratedReport, error) {
	return service.GeneratedReport{}, nil
}

func TestAuthLoginHandler(t *testing.T) {
	t.Parallel()
	auth := unexpectedStubAuth(&stubAuthService{loginResponse: service.LoginResponse{Token: "signed-token"}})
	handler, err := NewAuthHandler(auth)
	if err != nil {
		t.Fatalf("NewAuthHandler() error = %v", err)
	}
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)
	request := httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(`{"username":"operator","password":"secret"}`))
	response := httptest.NewRecorder()
	mux.ServeHTTP(response, request)
	if response.Code != http.StatusOK || !strings.Contains(response.Body.String(), `"token":"signed-token"`) {
		t.Fatalf("login status = %d, body = %s", response.Code, response.Body.String())
	}
}

func TestOperatorHandlerRequiresAuthentication(t *testing.T) {
	t.Parallel()
	auth := &stubAuthService{authErr: service.ErrInvalidToken}
	operator := &stubOperatorService{}
	handler, _ := NewOperatorHandler(operator, auth)
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)
	request := httptest.NewRequest(http.MethodPut, "/api/operator/tickets/%D0%9E%D0%A2%D0%9A-ABCD-2345/responsible-worker", strings.NewReader(`{"worker_id":7}`))
	response := httptest.NewRecorder()
	mux.ServeHTTP(response, request)
	if response.Code != http.StatusUnauthorized || operator.assignmentCalled {
		t.Fatalf("status = %d, assignment called = %v", response.Code, operator.assignmentCalled)
	}
}

func TestOperatorHandlerAssignsAuthenticatedWorker(t *testing.T) {
	t.Parallel()
	auth := &stubAuthService{user: service.AuthUser{ID: 3, Role: "operator"}}
	operator := &stubOperatorService{}
	handler, _ := NewOperatorHandler(operator, auth)
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)
	request := httptest.NewRequest(http.MethodPut, "/api/operator/tickets/%D0%9E%D0%A2%D0%9A-ABCD-2345/responsible-worker", strings.NewReader(`{"worker_id":7}`))
	request.Header.Set("Authorization", "Bearer token")
	response := httptest.NewRecorder()
	mux.ServeHTTP(response, request)
	if response.Code != http.StatusOK || !operator.assignmentCalled || operator.workerID != 7 || operator.actorID != 3 {
		t.Fatalf("status = %d, operator = %+v, body = %s", response.Code, operator, response.Body.String())
	}
}

func unexpectedStubAuth(stub *stubAuthService) *stubAuthService {
	stub.authErr = errors.New("Authenticate must not be called")
	return stub
}
