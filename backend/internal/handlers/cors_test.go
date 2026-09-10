package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWithCORSAllowsConfiguredFrontend(t *testing.T) {
	t.Parallel()
	handler := WithCORS(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}), "http://localhost:3000")
	request := httptest.NewRequest(http.MethodOptions, "/api/tickets", nil)
	request.Header.Set("Origin", "http://localhost:3000")
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusNoContent)
	}
	if response.Header().Get("Access-Control-Allow-Origin") != "http://localhost:3000" {
		t.Fatalf("Access-Control-Allow-Origin = %q", response.Header().Get("Access-Control-Allow-Origin"))
	}
}

func TestWithCORSRejectsUnknownPreflightOrigin(t *testing.T) {
	t.Parallel()
	handler := WithCORS(http.NotFoundHandler(), "http://localhost:3000")
	request := httptest.NewRequest(http.MethodOptions, "/api/tickets", nil)
	request.Header.Set("Origin", "https://example.org")
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusForbidden)
	}
}
