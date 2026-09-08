package handlers

import (
	"encoding/json"
	"net/http"

	"example.com/german/backend/internal/service"
)

type HealthHandler struct {
	service *service.HealthService
}

func NewHealthHandler(healthService *service.HealthService) *HealthHandler {
	return &HealthHandler{service: healthService}
}

func (h *HealthHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /health", h.getHealth)
}

func (h *HealthHandler) getHealth(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if err := json.NewEncoder(w).Encode(h.service.Status()); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
