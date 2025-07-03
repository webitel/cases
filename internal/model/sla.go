package model

import "time"

type SLA struct {
	*Author
	*Editor
	*Calendar
	Id             int        `json:"id" db:"id"`
	Name           *string    `json:"name" db:"name"`
	Description    *string    `json:"description" db:"description"`
	ValidFrom      *time.Time `json:"valid_from" db:"valid_from"`
	ValidTo        *time.Time `json:"valid_to" db:"valid_to"`
	ReactionTime   int        `json:"reaction_time" db:"reaction_time"`
	ResolutionTime int        `json:"resolution_time" db:"resolution_time"`
	CreatedAt      *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at" db:"updated_at"`
}
