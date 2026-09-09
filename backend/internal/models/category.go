package models

// Category maps to the categories table.
type Category struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}
