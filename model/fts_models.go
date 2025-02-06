package model

type FtsCase struct {
	Description   string  `json:"description,omitempty"`
	CloseResult   string  `json:"close_result,omitempty"`
	RatingComment string  `json:"rating_comment,omitempty"`
	RoleIds       []int64 `json:"_role_ids,omitempty"`
	Subject       string  `json:"subject,omitempty"`
	ContactInfo   string  `json:"contact_info,omitempty"`
	CreatedAt     int64   `json:"created_at,omitempty"`
}

type FtsCaseComment struct {
	ParentId  int64   `json:"parent_id,omitempty"`
	Comment   string  `json:"comment,omitempty"`
	RoleIds   []int64 `json:"_role_ids,omitempty"`
	CreatedAt int64   `json:"created_at,omitempty"`
}
