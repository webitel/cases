package auth_util

import (
	"github.com/webitel/cases/auth"
	"github.com/webitel/cases/auth/session/user_session"
)

func CloneWithUserID(src auth.Auther, overrideUserID int64) auth.Auther {
	session, ok := src.(*user_session.UserAuthSession)
	if !ok {
		return src
	}
	// Clone session
	newSession := *session
	// Clone user fully
	if newSession.User != nil {
		newUser := *newSession.User
		newUser.Id = overrideUserID
		newSession.User = &newUser
	}
	return &newSession
}
