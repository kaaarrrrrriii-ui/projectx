package service

import (
	"strings"

	"example.com/german/backend/internal/models"
)

func ticketPriorityFromString(value string) models.TicketPriority {
	switch strings.TrimSpace(value) {
	case "standard":
		return models.TicketPriorityStandard
	case "urgent":
		return models.TicketPriorityUrgent
	case "low":
		return models.TicketPriorityLow
	default:
		return 0
	}
}
