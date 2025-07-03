package model

import "time"

type Status struct {
	*Author
	*Editor
	Id          int
	Name        *string
	Description *string
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
}
