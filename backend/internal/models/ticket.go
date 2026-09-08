package models

import (
	"time"

	"github.com/google/uuid"
)

type ApplicantType string

const (
	ApplicantTypeStudent ApplicantType = "student"
	ApplicantTypeParent  ApplicantType = "parent"
	ApplicantTypeTeacher ApplicantType = "teacher"
)

type TicketStatus string

const (
	TicketStatusNew TicketStatus = "new"
)

type ClarificationAnswer struct {
	QuestionID string `json:"question_id"`
	Answer     string `json:"answer"`
}

type Ticket struct {
	ID uuid.UUID `json:"id"`

	ApplicantType ApplicantType `json:"applicant_type"`

	CategoryID *uuid.UUID `json:"category_id,omitempty"`

	Text string `json:"text"`

	ClarificationAnswers []ClarificationAnswer `json:"clarification_answers,omitempty"`

	TrackCodeHash []byte `json:"-"`

	Status TicketStatus `json:"status"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}