package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

type stubDatabasePinger struct {
	err error
}

func (p stubDatabasePinger) Ping(context.Context) error {
	return p.err
}

func TestReadinessHandler(t *testing.T) {
	tests := []struct {
		name       string
		pingError  error
		wantStatus int
		wantBody   string
	}{
		{name: "database available", wantStatus: http.StatusOK, wantBody: `{"status":"ready"}` + "\n"},
		{name: "database unavailable", pingError: errors.New("database unavailable"), wantStatus: http.StatusServiceUnavailable, wantBody: `{"status":"unavailable"}` + "\n"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			handler := NewReadinessHandler(stubDatabasePinger{err: test.pingError})
			request := httptest.NewRequest(http.MethodGet, "/ready", nil)
			response := httptest.NewRecorder()

			handler.getReadiness(response, request)

			if response.Code != test.wantStatus {
				t.Fatalf("status = %d, want %d", response.Code, test.wantStatus)
			}
			if response.Body.String() != test.wantBody {
				t.Fatalf("body = %q, want %q", response.Body.String(), test.wantBody)
			}
		})
	}
}
