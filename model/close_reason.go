package model

import "time"

type CloseReason struct {
	*Author
	*Editor
	Id                 int64     `json:"id" db:"id"`
	Name               string    `json:"name" db:"name"`
	Description        *string    `json:"description" db:"description"`
	CloseReasonGroupId int64     `json:"close_reason_group_id" db:"close_reason_id"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at"`
	Dc                 int64     `json:"dc" db:"dc"`
}

type Author struct {
	Id   *int64  `db:"created_by_id"`   // or updated_by_id
	Name *string `db:"created_by_name"` // or updated_by_name
}

type Editor struct {
	Id   *int64  `db:"updated_by_id"`   // or created_by_id
	Name *string `db:"updated_by_name"` // or created_by_name
}

type CloseReasonSearchOptions struct {
	DomainId           int64
	Page               int
	Size               int
	Fields             []string
	Sort               string
	Ids                []int64
	Q                  string
	CloseReasonGroupId int64
}
type CloseReasonList struct {
	Page  int
	Next  bool
	Items []*CloseReason
}
