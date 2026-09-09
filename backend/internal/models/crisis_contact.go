package models

import "time"

// CrisisContact maps to the crisis_contact table.
type CrisisContact struct {
	ID        int64     `json:"id"`
	Contact   string    `json:"contact"`
	TicketID  int64     `json:"ticket_id"`
	CreatedAt time.Time `json:"created_at"`
}
