package model

import "time"

type SLACondition struct {
	*Author
	*Editor
	Id             int64       `db:"id"`
	Name           *string     `db:"name"`
	Priorities     []*Priority `db:"priorities"`
	ReactionTime   *int        `db:"reaction_time"`
	ResolutionTime *int        `db:"resolution_time"`
	SlaId          *int64      `db:"sla_id"` // <-- should be *int64
	CreatedAt      *time.Time  `db:"created_at"`
	UpdatedAt      *time.Time  `db:"updated_at"`
}
