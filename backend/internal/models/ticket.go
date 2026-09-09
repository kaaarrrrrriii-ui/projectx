package models

import "time"

// Ticket maps to the tickets table.
type Ticket struct {
	ID          int64      `json:"id"`
	TrackID     string     `json:"track_id"`
	Priority    int        `json:"priority"`
	MordaType   int        `json:"morda_type"`
	ClosedAt    *time.Time `json:"closed_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	CategoryID  int64      `json:"category_id"`
	Status      int        `json:"status"`
	ReturnCount int        `json:"return_count"`
}
