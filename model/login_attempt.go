package model

import "time"

type LoginAttempt struct {
	Date       time.Time
	UserId     *NullInt
	DomainId   *NullInt
	Details    *NullString
	AuthType   string
	UserIp     string
	UserName   string
	UserAgent  string
	DomainName string
	Id         int64
	Success    bool
}

var LoginAttemptFields = struct {
	Id         string
	Success    string
	AuthType   string
	UserIp     string
	Date       string
	User       string
	UserName   string
	UserId     string
	UserAgent  string
	DomainId   string
	DomainName string
	Details    string
}{
	Id:         "id",
	Success:    "success",
	AuthType:   "auth_type",
	UserIp:     "user_ip",
	Date:       "date",
	User:       "user",
	UserId:     "user_id",
	UserName:   "user_name",
	UserAgent:  "user_agent",
	DomainId:   "domain_id",
	DomainName: "domain_name",
	Details:    "details",
}
