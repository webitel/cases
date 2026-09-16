package model

import (
	"time"

	"github.com/webitel/cases/api/cases"
)

type CaseAMQPMessage struct {
	Case *cases.Case `json:"case"`
}

type CaseLinkAMQPMessage struct {
	CaseLink *cases.CaseLink `json:"case_link"`
}

type CaseCommentAMQPMessage struct {
	CaseComment *cases.CaseComment `json:"case_comment"`
}

type CaseFileAMQPMessage struct {
	CaseFile *cases.File `json:"case_file"`
}

type RelatedCaseAMQPMessage struct {
	RelatedCase *cases.RelatedCase `json:"related_case"`
}

type CaseResolutionAMQPMessage struct {
	Type       string    `json:"type"`
	Schema     int       `json:"schema"`
	OccurredAt time.Time `json:"occurred_at"`
	DomainID   int64     `json:"domain_id"`
	CaseID     int64     `json:"case_id"`
	ArticleID  int64     `json:"article_id"`
}
