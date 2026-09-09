package models

// CategoryExpertGroup maps to the cats_expert_groups table.
type CategoryExpertGroup struct {
	GroupID int64 `json:"group_id"`
	CatID   int64 `json:"cat_id"`
}
