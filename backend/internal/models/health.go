package models

type HealthResponse struct {
	Status string `json:"status"`
	Uptime string `json:"uptime"`
}
