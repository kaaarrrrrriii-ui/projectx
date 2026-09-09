package models

// Question maps to the questions table.
type Question struct {
	ID    int64  `json:"id"`
	CatID int64  `json:"cat_id"`
	Text  string `json:"text"`
}
