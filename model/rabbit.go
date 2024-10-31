package model

import (
	"time"

	guid "github.com/google/uuid"
	cerror "github.com/webitel/cases/internal/error"
)

type BrokerRecordLogMessage struct {
	Records                        []*LogEntity `json:"records,omitempty"`
	BrokerLogMessageRequiredFields `json:"requiredFields"`
}

type BrokerLogMessageRequiredFields struct {
	UserIp string `json:"userIp,omitempty"`
	Action string `json:"action,omitempty"`
	UserId int    `json:"userId,omitempty"`
	Date   int64  `json:"date,omitempty"`
}

type LogEntity struct {
	NewState BytesJSON `json:"newState,omitempty"`
	Id       int64     `json:"id,omitempty"`
}

type BytesJSON struct {
	Body []byte
}

func (b *BytesJSON) GetBody() []byte {
	return b.Body
}

func (b *BytesJSON) UnmarshalJSON(input []byte) error {
	b.Body = input
	return nil
}

type BrokerLoginMessage struct {
	Login    *Login  `json:"login,omitempty"`
	Status   *Status `json:"status,omitempty"`
	AuthType string  `json:"type,omitempty"`
	Agent    string  `json:"agent,omitempty"`
	From     string  `json:"from,omitempty"`
	Date     int64   `json:"date,omitempty"`
	IsNew    bool    `json:"isNew,omitempty"`
}

func (m *BrokerLoginMessage) ConvertToDatabaseModel() (*LoginAttempt, cerror.AppError) {
	var (
		success       bool
		databaseModel LoginAttempt
		authType      string
	)
	if m.Status != nil {
		success = false
		databaseModel.Details = NewNullString(m.Status.Detail)
	} else {
		success = true
	}
	if user := m.Login.User; user != nil {
		if user.Id != 0 {
			id, err := NewNullInt(user.Id)
			if err != nil {
				return nil, cerror.NewInternalError("app.log.handle_rabbit_login_message.parse_user_id.error", err.Error())
			}
			databaseModel.UserId = id
		}
		databaseModel.UserName = user.Username

	}
	if domain := m.Login.Domain; domain != nil {
		if domain.Id != 0 {
			id, err := NewNullInt(domain.Id)
			if err != nil {
				return nil, cerror.NewInternalError("app.log.handle_rabbit_login_message.parse_domain_id.error", err.Error())
			}
			databaseModel.DomainId = id
		}
		databaseModel.DomainName = domain.Name

	}
	authType = m.AuthType
	if authType == "" {
		authType = "password"
	}
	databaseModel.AuthType = authType
	databaseModel.UserAgent = m.Agent
	databaseModel.Date = time.UnixMilli(m.Date)
	databaseModel.UserIp = m.From
	databaseModel.Success = success

	return &databaseModel, nil
}

type Login struct {
	Id        *guid.UUID        `json:"id,omitempty"`
	Context   map[string]string `json:"context,omitempty"`
	Domain    *Domain           `json:"domain,omitempty"`
	User      *User             `json:"user,omitempty"`
	CreatedAt int64             `json:"created_at,omitempty"`
	ExpiresAt int64             `json:"expires_at,omitempty"`
	MaxAge    int64             `json:"max_age,omitempty"`
}

type User struct {
	Name      string `json:"name,omitempty"`
	Email     string `json:"email,omitempty"`
	Username  string `json:"username,omitempty"`
	Extension string `json:"extension,omitempty"`
	Id        int    `json:"id,omitempty"`
}

type Status struct {
	Id     string `json:"id,omitempty"`
	Status string `json:"status,omitempty"`
	Detail string `json:"detail,omitempty"`
	Code   int    `json:"code,omitempty"`
}

type Domain struct {
	Name string `json:"name,omitempty"`
	Id   int64  `json:"id,omitempty"`
}
