package models

import "time"

// Ticket maps to the tickets table.
type Ticket struct {
	ID          int64        `json:"id"`
	TrackID     string       `json:"track_id"`
	Priority    int          `json:"priority"`
	MordaType   int          `json:"morda_type"`
	ClosedAt    *time.Time   `json:"closed_at,omitempty"`
	CreatedAt   time.Time    `json:"created_at"`
	CategoryID  int64        `json:"category_id"`
	Status      TicketStatus `json:"-"`
	ReturnCount int          `json:"return_count"`
}

// TicketStatus is the integer representation persisted in tickets.status.
// Public API responses use String so database values do not leak to clients.
type TicketStatus int

const (
	TicketStatusNew TicketStatus = iota + 1
	TicketStatusAssigned
	TicketStatusInProgress
	TicketStatusNeedsClarification
	TicketStatusAnswerReady
	TicketStatusReturned
	TicketStatusCompleted
	TicketStatusRejected
	TicketStatusClosedWithoutAnswer
)

func (status TicketStatus) String() string {
	switch status {
	case TicketStatusNew:
		return "new"
	case TicketStatusAssigned:
		return "assigned"
	case TicketStatusInProgress:
		return "in_progress"
	case TicketStatusNeedsClarification:
		return "needs_clarification"
	case TicketStatusAnswerReady:
		return "answer_ready"
	case TicketStatusReturned:
		return "returned"
	case TicketStatusCompleted:
		return "completed"
	case TicketStatusRejected:
		return "rejected"
	case TicketStatusClosedWithoutAnswer:
		return "closed_without_answer"
	default:
		return "unknown"
	}
}

func (status TicketStatus) IsApplicantCompletable() bool {
	switch status {
	case TicketStatusNew,
		TicketStatusAssigned,
		TicketStatusInProgress,
		TicketStatusNeedsClarification,
		TicketStatusAnswerReady,
		TicketStatusReturned:
		return true
	default:
		return false
	}
}
