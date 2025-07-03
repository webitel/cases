package model

import (
	"time"
)

type Service struct {
	*Author
	*Editor
	Id          int                    `json:"id,omitempty" db:"id"`
	Name        *string                `json:"name,omitempty" db:"name"`
	RootId      *int                   `json:"root_id,omitempty" db:"root_id"`
	Description *string                `json:"description,omitempty" db:"description"`
	Code        *string                `json:"code,omitempty" db:"code"`
	State       *bool                  `json:"state,omitempty" db:"state"`
	Sla         *GeneralLookup         `json:"sla,omitempty" db:"sla"`
	Group       *GeneralExtendedLookup `json:"group,omitempty" db:"group"`
	Assignee    *GeneralLookup         `json:"assignee,omitempty" db:"assignee"`
	CreatedAt   *time.Time             `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt   *time.Time             `json:"updated_at,omitempty" db:"updated_at"`
	CatalogId   *int                   `json:"catalog_id,omitempty" db:"catalog_id"`
	Services    []*Service             `json:"services,omitempty" db:"services"`
	Searched    *bool                  `json:"searched,omitempty" db:"searched"`
}
