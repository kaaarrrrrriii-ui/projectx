package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"example.com/german/backend/internal/service"
)

type stubAdminHandlerService struct {
	adminService
	dashboardCalled bool
}

func (stub *stubAdminHandlerService) Dashboard(context.Context) (service.AdminDashboardResponse, error) {
	stub.dashboardCalled = true
	return service.AdminDashboardResponse{NewCount: 3}, nil
}

func (stub *stubAdminHandlerService) Ticket(context.Context, string) (service.AdminTicketResponse, error) {
	return service.AdminTicketResponse{
		TrackID: "ОТК-ABCD-2345", Status: "assigned", Priority: "urgent",
		ApplicantType: "schoolchild", CreatedAt: time.Date(2026, 9, 10, 8, 0, 0, 0, time.UTC),
	}, nil
}

func TestAdminHandlerRejectsOperator(t *testing.T) {
	t.Parallel()
	auth := &stubAuthService{user: service.AuthUser{ID: 2, Role: "operator"}}
	admin := &stubAdminHandlerService{}
	handler, _ := NewAdminHandler(admin, auth)
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)
	request := httptest.NewRequest(http.MethodGet, "/api/admin/dashboard", nil)
	request.Header.Set("Authorization", "Bearer token")
	response := httptest.NewRecorder()
	mux.ServeHTTP(response, request)
	if response.Code != http.StatusForbidden || admin.dashboardCalled {
		t.Fatalf("status = %d, dashboard called = %v", response.Code, admin.dashboardCalled)
	}
}

func TestAdminHandlerAllowsAdmin(t *testing.T) {
	t.Parallel()
	auth := &stubAuthService{user: service.AuthUser{ID: 1, Role: "admin"}}
	admin := &stubAdminHandlerService{}
	handler, _ := NewAdminHandler(admin, auth)
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)
	request := httptest.NewRequest(http.MethodGet, "/api/admin/dashboard", nil)
	request.Header.Set("Authorization", "Bearer token")
	response := httptest.NewRecorder()
	mux.ServeHTTP(response, request)
	if response.Code != http.StatusOK || !admin.dashboardCalled || !strings.Contains(response.Body.String(), `"new_count":3`) {
		t.Fatalf("status = %d, body = %s", response.Code, response.Body.String())
	}
}

func TestAdminTicketDoesNotExposeConversation(t *testing.T) {
	t.Parallel()
	auth := &stubAuthService{user: service.AuthUser{ID: 1, Role: "admin"}}
	handler, _ := NewAdminHandler(&stubAdminHandlerService{}, auth)
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)
	request := httptest.NewRequest(http.MethodGet, "/api/admin/tickets/%D0%9E%D0%A2%D0%9A-ABCD-2345", nil)
	request.Header.Set("Authorization", "Bearer token")
	response := httptest.NewRecorder()
	mux.ServeHTTP(response, request)
	body := response.Body.String()
	for _, forbidden := range []string{"messages", "attachments", "notes", "description", "clarifications"} {
		if strings.Contains(body, forbidden) {
			t.Fatalf("admin response contains forbidden field %q: %s", forbidden, body)
		}
	}
}
