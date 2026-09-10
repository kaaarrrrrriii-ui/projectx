package repos

import "errors"

var (
	ErrTicketNotFound         = errors.New("ticket not found")
	ErrAttachmentNotFound     = errors.New("attachment not found")
	ErrInvalidTransition      = errors.New("invalid ticket status transition")
	ErrReturnLimitReached     = errors.New("ticket return limit reached")
	ErrAttachmentLimitReached = errors.New("ticket attachment limit reached")
	ErrWorkerNotFound         = errors.New("worker not found")
	ErrWorkerUnavailable      = errors.New("worker is unavailable")
	ErrWorkerAlreadyAssigned  = errors.New("worker is already assigned")
	ErrWorkerNotAssigned      = errors.New("worker is not assigned")
	ErrResponsibleWorker      = errors.New("responsible worker cannot be removed as co-worker")
	ErrResponsibleRequired    = errors.New("responsible worker must be assigned first")
	ErrExpertAccessDenied     = errors.New("expert is not assigned to the ticket")
	ErrWorkerRequestNotFound  = errors.New("worker request not found")
	ErrAdminResourceNotFound  = errors.New("admin resource not found")
	ErrAdminResourceInUse     = errors.New("admin resource is in use")
	ErrUsernameExists         = errors.New("username already exists")
	ErrNameExists             = errors.New("name already exists")
	ErrEmployeeHasTickets     = errors.New("employee has active tickets")
	ErrQuestionInUse          = errors.New("question is in use")
)
