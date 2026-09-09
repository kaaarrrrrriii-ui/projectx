package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

type databasePinger interface {
	Ping(context.Context) error
}

type ReadinessHandler struct {
	database databasePinger
}

func NewReadinessHandler(database databasePinger) *ReadinessHandler {
	return &ReadinessHandler{database: database}
}

func (h *ReadinessHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /ready", h.getReadiness)
}

func (h *ReadinessHandler) getReadiness(w http.ResponseWriter, request *http.Request) {
	ctx, cancel := context.WithTimeout(request.Context(), 2*time.Second)
	defer cancel()

	status := http.StatusOK
	response := map[string]string{"status": "ready"}
	if err := h.database.Ping(ctx); err != nil {
		status = http.StatusServiceUnavailable
		response["status"] = "unavailable"
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(response)
}
