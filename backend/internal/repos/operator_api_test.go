package repos

import (
	"testing"

	"example.com/german/backend/internal/models"
)

func TestOperatorTransitionRules(t *testing.T) {
	t.Parallel()
	tests := []struct {
		from, to models.TicketStatus
		allowed  bool
	}{
		{models.TicketStatusNew, models.TicketStatusAssigned, true},
		{models.TicketStatusNew, models.TicketStatusRejected, true},
		{models.TicketStatusReturned, models.TicketStatusCompleted, true},
		{models.TicketStatusAssigned, models.TicketStatusCompleted, false},
		{models.TicketStatusNew, models.TicketStatusCompleted, false},
	}
	for _, test := range tests {
		if actual := operatorTransitionAllowed(test.from, test.to); actual != test.allowed {
			t.Fatalf("operatorTransitionAllowed(%s, %s) = %v", test.from.String(), test.to.String(), actual)
		}
	}
}

func TestOperatorTicketDefaultOrder(t *testing.T) {
	t.Parallel()
	if got := operatorTicketOrder(OperatorTicketFilter{Queue: "new"}); got != "crisis_detected DESC, CASE t.priority WHEN 2 THEN 0 WHEN 1 THEN 1 ELSE 2 END, t.created_at ASC, t.id ASC" {
		t.Fatalf("new order = %q", got)
	}
}
