package auth_util

import (
	"github.com/webitel/cases/auth"
	"github.com/webitel/cases/auth/user_auth"
)

func CloneWithUserID(original auth.Auther, overrideUserID int64) auth.Auther {
	session, ok := original.(*user_auth.UserAuthSession)
	if !ok {
		return original
	}
	// Clone
	newSession := *session
	user := *newSession.User
	user.Id = overrideUserID
	newSession.User = &user

	return &newSession
}
