package model

import "time"

const (
	CaseCommentTypeInternal int32 = 1
	CaseCommentTypeExternal int32 = 2

	CaseCommentTypesContextKey = "comment_types"
)

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
	Type      int32      `json:"type" db:"type"`
	PrevType  int32      `json:"-" db:"prev_type"`
}

func CaseCommentScope(commentType int32) string {
	if commentType == CaseCommentTypeExternal {
		return ScopeCaseComments
	}

	return ScopeInternalComments
}

func CaseCommentTypeOrDefault(commentType int32) int32 {
	if commentType == CaseCommentTypeExternal {
		return CaseCommentTypeExternal
	}

	return CaseCommentTypeInternal
}
