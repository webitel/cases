package model

import (
	"time"
)

type Priority struct {
	*Author
	*Editor
	Id                 int64     `json:"id" db:"id"`
	Name                string    `json:"name" db:"name"`
	Description         *string   `json:"description" db:"description"`
	Color               string    `json:"color" db:"color"`
	CreatedAt           *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt           *time.Time `json:"updated_at" db:"updated_at"`
}

type PriorityList struct {
	Page  int
	Next  bool
	Items []*Priority
}
