package auth_util

import (
	"github.com/webitel/cases/auth"
	"github.com/webitel/cases/auth/user_auth"
)

func CloneWithUserID(src auth.Auther, overrideUserID int64) auth.Auther {
	session, ok := src.(*user_auth.UserAuthSession)
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
