package model

import "time"

type SLACondition struct {
	*Author
	*Editor
	Id             int         `db:"id"`
	Name           *string     `db:"name"`
	Priorities     []*Priority `db:"priorities"`
	ReactionTime   *int        `db:"reaction_time"`
	ResolutionTime *int        `db:"resolution_time"`
	SlaId          *int        `db:"sla_id"`
	CreatedAt      *time.Time  `db:"created_at"`
	UpdatedAt      *time.Time  `db:"updated_at"`
}
