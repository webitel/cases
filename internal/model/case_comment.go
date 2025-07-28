package model

import "time"

type CaseComment struct {
	*Author
	*Editor
	*Contact
	Id        int64      `json:"id" db:"id"`
	Ver       int32      `json:"ver" db:"ver"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
	Text      string     `json:"text" db:"text"`
	Edited    bool       `json:"edited" db:"edited"`
	CanEdit   bool       `json:"can_edit" db:"can_edit"`
	CaseId    int64      `json:"case_id" db:"case_id"`
	RoleIds   []int64    `json:"role_ids" db:"role_ids"`
}
