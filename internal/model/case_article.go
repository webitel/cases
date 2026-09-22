package model

import "time"

// Sources of a case article link.
const (
	CaseArticleSourceManual     int32 = 1
	CaseArticleSourceResolution int32 = 2
)

// CaseArticle is a knowledge base article linked to a case.
type CaseArticle struct {
	*Author

	ArticleID int64      `json:"article_id" db:"article_id"`
	Source    int32      `json:"source" db:"source"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	// Set on the article side only.
	CaseID      *int64  `json:"case_id,omitempty" db:"case_id"`
	CaseVer     *int32  `json:"case_ver,omitempty" db:"case_ver"`
	CaseName    *string `json:"case_name,omitempty" db:"case_name"`
	CaseSubject *string `json:"case_subject,omitempty" db:"case_subject"`
}
