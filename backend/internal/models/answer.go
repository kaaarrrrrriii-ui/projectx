package models

// Answer maps to the answers table.
type Answer struct {
	ID         int64  `json:"id"`
	QuestionID int64  `json:"question_id"`
	Text       string `json:"text"`
}
