package model

import (
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
