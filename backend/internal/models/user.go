package models

// User maps to the users table.
type User struct {
	ID            int64   `json:"id"`
	Username      string  `json:"username"`
	PasswordHash  string  `json:"-"`
	Role          string  `json:"role"`
	ExpertGroupID int64   `json:"expert_group_id"`
	FullName      string  `json:"full_name"`
	AvgRating     float64 `json:"avg_rating"`
	MaxTickets    int     `json:"max_tickets"`
}
