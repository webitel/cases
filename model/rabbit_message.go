package model

type RabbitMessage struct {
	UserIp   string `json:"userIp,omitempty"`
	Action   string `json:"action,omitempty"`
	Schema   string
	NewState []byte `json:"newState,omitempty"`
	UserId   int    `json:"userId,omitempty"`
	Date     int64  `json:"date,omitempty"`
	RecordId int64  `json:"recordId,omitempty"`
}
