package repos

import (
	"testing"

	"example.com/german/backend/internal/models"
)

func TestStatusAfterApplicantMessage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		status models.TicketStatus
		want   models.TicketStatus
		ok     bool
	}{
		{status: models.TicketStatusNew, want: models.TicketStatusNew, ok: true},
		{status: models.TicketStatusAssigned, want: models.TicketStatusAssigned, ok: true},
		{status: models.TicketStatusInProgress, want: models.TicketStatusInProgress, ok: true},
		{status: models.TicketStatusNeedsClarification, want: models.TicketStatusInProgress, ok: true},
		{status: models.TicketStatusAnswerReady, want: models.TicketStatusInProgress, ok: true},
		{status: models.TicketStatusReturned, want: models.TicketStatusReturned, ok: false},
		{status: models.TicketStatusCompleted, want: models.TicketStatusCompleted, ok: false},
		{status: models.TicketStatusRejected, want: models.TicketStatusRejected, ok: false},
		{status: models.TicketStatusClosedWithoutAnswer, want: models.TicketStatusClosedWithoutAnswer, ok: false},
	}

	for _, test := range tests {
		got, ok := statusAfterApplicantMessage(test.status)
		if got != test.want || ok != test.ok {
			t.Errorf("statusAfterApplicantMessage(%s) = (%s, %v), want (%s, %v)", test.status.String(), got.String(), ok, test.want.String(), test.ok)
		}
	}
}
