package models

import "time"

// TicketWorker maps to the tickets_workers table.
type TicketWorker struct {
	WorkerID      int64     `json:"worker_id"`
	TicketID      int64     `json:"ticket_id"`
	IsResponsible bool      `json:"is_responsible"`
	CreatedAt     time.Time `json:"created_at"`
	Actual        bool      `json:"actual"`
}
