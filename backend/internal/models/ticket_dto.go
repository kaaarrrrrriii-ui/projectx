package models

import (
	"time"

	"github.com/google/uuid"
)

type CreateTicketRequest struct {
	ApplicantType ApplicantType `json:"applicant_type"`
	CategoryID    *uuid.UUID    `json:"category_id,omitempty"`
	Text          string        `json:"text"`

	ClarificationAnswers []ClarificationAnswer `json:"clarification_answers,omitempty"`
}

type CreateTicketResponse struct {
	ID        uuid.UUID    `json:"id"`
	TrackCode string       `json:"track_code"`
	Status    TicketStatus `json:"status"`
	CreatedAt time.Time    `json:"created_at"`
}