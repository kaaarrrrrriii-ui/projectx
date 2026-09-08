package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type TicketEvent struct {
	ID uuid.UUID `json:"id"`

	TicketID uuid.UUID `json:"ticket_id"`

	Type string `json:"type"`

	Payload json.RawMessage `json:"payload"`

	CreatedAt time.Time `json:"created_at"`
}