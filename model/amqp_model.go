package model

import (
	"github.com/webitel/cases/api/cases"
)

type CaseAMQPMessage struct {
	DomainId int64       `json:"domain_id"`
	Case     *cases.Case `json:"case"`
}
