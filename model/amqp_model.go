package model

import (
	"github.com/webitel/cases/api/cases"
)

type CaseAMQPMessage struct {
	Case *cases.Case `json:"case"`
}
