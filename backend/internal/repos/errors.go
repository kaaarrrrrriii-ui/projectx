package repos

import "errors"

var (
	ErrTicketNotFound         = errors.New("ticket not found")
	ErrAttachmentNotFound     = errors.New("attachment not found")
	ErrInvalidTransition      = errors.New("invalid ticket status transition")
	ErrReturnLimitReached     = errors.New("ticket return limit reached")
	ErrAttachmentLimitReached = errors.New("ticket attachment limit reached")
)
