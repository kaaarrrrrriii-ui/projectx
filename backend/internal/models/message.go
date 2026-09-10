package models

import "time"

// Message maps to the messages table.
type Message struct {
	ID        int64       `json:"id"`
	TicketID  int64       `json:"ticket_id"`
	Text      string      `json:"text"`
	Type      MessageType `json:"-"`
	CreatedAt time.Time   `json:"created_at"`
}

// MessageType is the integer representation persisted in messages.type.
type MessageType int

const (
	MessageTypeApplicant MessageType = iota + 1
	MessageTypeSpecialist
	MessageTypeInternalNote
	MessageTypeSystem
	MessageTypeReturnReason
	MessageTypeOperator
)

func (messageType MessageType) String() string {
	switch messageType {
	case MessageTypeApplicant:
		return "applicant"
	case MessageTypeSpecialist:
		return "specialist"
	case MessageTypeInternalNote:
		return "internal_note"
	case MessageTypeSystem:
		return "system"
	case MessageTypeReturnReason:
		return "return_reason"
	case MessageTypeOperator:
		return "operator"
	default:
		return "unknown"
	}
}
