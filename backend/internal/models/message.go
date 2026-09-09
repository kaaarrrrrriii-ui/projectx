package models

import "time"

// Message maps to the messages table.
type Message struct {
	ID        int64     `json:"id"`
	TicketID  int64     `json:"ticket_id"`
	Text      string    `json:"text"`
	Type      int       `json:"type"`
	CreatedAt time.Time `json:"created_at"`
}
