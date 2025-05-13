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
	// Clone
	newSession := *session
	user := *newSession.User
	user.Id = overrideUserID
	newSession.User = &user

	return &newSession
}
