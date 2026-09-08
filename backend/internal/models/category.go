package models

import "github.com/google/uuid"

type Category struct {
	ID     uuid.UUID `json:"id"`
	Name   string    `json:"name"`
	Slug   string    `json:"slug"`
	Active bool      `json:"active"`
}