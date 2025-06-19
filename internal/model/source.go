package model

import "time"

type Source struct {
	*Author
	*Editor
	Id          int        `db:"id"`
	Name        *string    `db:"name"`
	Description *string    `db:"description"`
	Type        *string    `db:"type"`
	CreatedAt   *time.Time `db:"created_at"`
	UpdatedAt   *time.Time `db:"updated_at"`
}
