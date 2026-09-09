package models

// QA maps to the quoted PostgreSQL table "QA" from the UML.
type QA struct {
	QuestionID int64 `json:"question_id"`
	AnswerID   int64 `json:"answer_id"`
	TicketID   int64 `json:"ticket_id"`
}
