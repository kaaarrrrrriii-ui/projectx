package service

import (
	"time"

	"example.com/german/backend/internal/models"
)

type HealthService struct {
	startedAt time.Time
}

func NewHealthService() *HealthService {
	return &HealthService{startedAt: time.Now()}
}

func (s *HealthService) Status() models.HealthResponse {
	return models.HealthResponse{
		Status: "ok",
		Uptime: time.Since(s.startedAt).Round(time.Second).String(),
	}
}
