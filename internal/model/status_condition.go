package model

import "time"

type StatusCondition struct {
	*Author
	*Editor
	Id          int        `db:"id"`
	Name        *string    `db:"name"`
	Description *string    `db:"description"`
	Initial     *bool      `db:"initial"`
	Final       *bool      `db:"final"`
	StatusId    *int       `db:"status_id"`
	CreatedAt   *time.Time `db:"created_at"`
	UpdatedAt   *time.Time `db:"updated_at"`
}
