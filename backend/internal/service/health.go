package service

import "time"

type HealthResponse struct {
	Status string `json:"status"`
	Uptime string `json:"uptime"`
}

type HealthService struct {
	startedAt time.Time
}

func NewHealthService() *HealthService {
	return &HealthService{startedAt: time.Now()}
}

func (s *HealthService) Status() HealthResponse {
	return HealthResponse{
		Status: "ok",
		Uptime: time.Since(s.startedAt).Round(time.Second).String(),
	}
}
