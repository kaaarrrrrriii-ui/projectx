package models

// Review maps to the reviews table.
type Review struct {
	ID       int64  `json:"id"`
	TicketID int64  `json:"ticket_id"`
	Text     string `json:"text"`
	Rating   int    `json:"rating"`
}
