package models

import "time"

// Ticket maps to the tickets table.
type Ticket struct {
	ID          int64          `json:"id"`
	TrackID     string         `json:"track_id"`
	Priority    TicketPriority `json:"-"`
	MordaType   ApplicantType  `json:"-"`
	ClosedAt    *time.Time     `json:"closed_at,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	CategoryID  int64          `json:"category_id"`
	Status      TicketStatus   `json:"-"`
	ReturnCount int            `json:"return_count"`
}

// TicketPriority is persisted in tickets.priority. API responses expose its
// string form so clients do not depend on database integers.
type TicketPriority int

const (
	TicketPriorityStandard TicketPriority = iota + 1
	TicketPriorityUrgent
)

func (priority TicketPriority) String() string {
	switch priority {
	case TicketPriorityStandard:
		return "standard"
	case TicketPriorityUrgent:
		return "urgent"
	default:
		return "unknown"
	}
}

// ApplicantType is persisted in tickets.morda_type.
type ApplicantType int

const (
	ApplicantTypeSchoolchild ApplicantType = iota + 1
	ApplicantTypeParent
	ApplicantTypeTeacher
)

func (applicantType ApplicantType) String() string {
	switch applicantType {
	case ApplicantTypeSchoolchild:
		return "schoolchild"
	case ApplicantTypeParent:
		return "parent"
	case ApplicantTypeTeacher:
		return "teacher"
	default:
		return "unknown"
	}
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
	return status == TicketStatusAnswerReady
}
