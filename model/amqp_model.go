package model

import (
	"github.com/webitel/cases/api/cases"
)

type CaseAMQPMessage struct {
	DomainId int64
	Case     *cases.Case
}
