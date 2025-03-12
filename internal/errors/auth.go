package errors

import (
	"encoding/json"
	"fmt"
	"net/http"

	goi18n "github.com/nicksnyder/go-i18n/i18n"
)

type AuthError interface {
	SetTranslationParams(map[string]any) AuthError
	GetTranslationParams() map[string]any
	SetStatusCode(int) AuthError
	GetStatusCode() int
	SetDetailedError(string)
	GetDetailedError() string
	SetRequestId(string)
	GetRequestId() string
	GetId() string

	Error() string
	Translate(goi18n.TranslateFunc)
	SystemMessage(goi18n.TranslateFunc) string
	ToJson() string
	String() string
}

type AuthorizationError struct {
	params        map[string]interface{}
	Id            string `json:"id"`
	Where         string `json:"where,omitempty"`
	Status        string `json:"status"`
	DetailedError string `json:"detail"`
	RequestId     string `json:"request_id,omitempty"`
	StatusCode    int    `json:"code,omitempty"`
}

func (err *AuthorizationError) SetTranslationParams(params map[string]any) AuthError {
	err.params = params
	return err
}

func (err *AuthorizationError) GetTranslationParams() map[string]any {
	return err.params
}

func (err *AuthorizationError) SetStatusCode(code int) AuthError {
	err.StatusCode = code
	err.Status = http.StatusText(err.StatusCode)
	return err
}

func (err *AuthorizationError) GetStatusCode() int {
	return err.StatusCode
}

func (err *AuthorizationError) Error() string {
	return fmt.Sprintf("AuthError [%s]: %s, %s", err.Id, err.Status, err.DetailedError)
}

func (err *AuthorizationError) SetDetailedError(details string) {
	err.DetailedError = details
}

func (err *AuthorizationError) GetDetailedError() string {
	return err.DetailedError
}

func (err *AuthorizationError) Translate(T goi18n.TranslateFunc) {
	if T == nil && err.DetailedError == "" {
		err.DetailedError = err.Id
		return
	}

	var errText string

	if err.params == nil {
		errText = T(err.Id)
	} else {
		errText = T(err.Id, err.params)
	}

	if errText != err.Id {
		err.DetailedError = errText
	}
}

func (err *AuthorizationError) SystemMessage(T goi18n.TranslateFunc) string {
	if err.params == nil {
		return T(err.Id)
	} else {
		return T(err.Id, err.params)
	}
}

func (err *AuthorizationError) SetRequestId(id string) {
	err.RequestId = id
}

func (err *AuthorizationError) GetRequestId() string {
	return err.RequestId
}

func (err *AuthorizationError) GetId() string {
	return err.Id
}

func (err *AuthorizationError) ToJson() string {
	b, _ := json.Marshal(err)
	return string(b)
}

func (err *AuthorizationError) String() string {
	if err.Id == err.Status && err.DetailedError != "" {
		return err.DetailedError
	}
	return err.Status
}

func NewUnauthorizedError(id, details string) AuthError {
	return newAuthError(id, details).SetStatusCode(http.StatusUnauthorized)
}

func NewPermissionForbiddenError(id, details string) AuthError {
	return newAuthError(id, details).SetStatusCode(http.StatusForbidden)
}

func newAuthError(id string, details string) AuthError {
	return &AuthorizationError{Id: id, Status: id, DetailedError: details}
}
