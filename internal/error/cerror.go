package error

import (
	"encoding/json"
	"fmt"
	"net/http"

	goi18n "github.com/nicksnyder/go-i18n/i18n"
)

var translateFunc goi18n.TranslateFunc = nil

func AppErrorInit(t goi18n.TranslateFunc) {
	translateFunc = t
}

type AppError interface {
	SetTranslationParams(map[string]any) AppError
	GetTranslationParams() map[string]any
	SetStatusCode(int) AppError
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

type ApplicationError struct {
	params        map[string]interface{}
	Id            string `json:"id"`
	Where         string `json:"where,omitempty"`
	Status        string `json:"status"`
	DetailedError string `json:"detail"`
	RequestId     string `json:"request_id,omitempty"`
	StatusCode    int    `json:"code,omitempty"`
}

func (err *ApplicationError) SetTranslationParams(params map[string]any) AppError {
	err.params = params
	return err
}

func (err *ApplicationError) GetTranslationParams() map[string]any {
	return err.params
}

func (err *ApplicationError) SetStatusCode(code int) AppError {
	err.StatusCode = code
	err.Status = http.StatusText(err.StatusCode)
	return err
}

func (err *ApplicationError) GetStatusCode() int {
	return err.StatusCode
}

func (err *ApplicationError) Error() string {
	var where string
	if err.Where != "" {
		where = err.Where + ": "
	}
	return fmt.Sprintf("%s%s, %s", where, err.Status, err.DetailedError)
}

func (err *ApplicationError) SetDetailedError(details string) {
	err.DetailedError = details
}

func (err *ApplicationError) GetDetailedError() string {
	return err.DetailedError
}

func (err *ApplicationError) Translate(T goi18n.TranslateFunc) {
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

func (err *ApplicationError) SystemMessage(T goi18n.TranslateFunc) string {
	if err.params == nil {
		return T(err.Id)
	} else {
		return T(err.Id, err.params)
	}
}

func (err *ApplicationError) SetRequestId(id string) {
	err.RequestId = id
}

func (err *ApplicationError) GetRequestId() string {
	return err.RequestId
}

func (err *ApplicationError) GetId() string {
	return err.Id
}

func (err *ApplicationError) ToJson() string {
	b, _ := json.Marshal(err)
	return string(b)
}

func (err *ApplicationError) String() string {
	if err.Id == err.Status && err.DetailedError != "" {
		return err.DetailedError
	}
	return err.Status
}

// Error constructors
func NewInternalError(id string, details string) AppError {
	return newAppError(id, details).SetStatusCode(http.StatusInternalServerError)
}

func NewNotFoundError(id string, details string) AppError {
	return newAppError(id, details).SetStatusCode(http.StatusNotFound)
}

func NewBadRequestError(id string, details string) AppError {
	return newAppError(id, details).SetStatusCode(http.StatusBadRequest)
}

func NewForbiddenError(id string, details string) AppError {
	return newAppError(id, details).SetStatusCode(http.StatusForbidden)
}

func NewCustomCodeError(id string, details string, code int) AppError {
	if code > 511 || code < 100 {
		code = http.StatusInternalServerError
	}
	return newAppError(id, details).SetStatusCode(code)
}

func newAppError(id string, details string) AppError {
	return &ApplicationError{Id: id, Status: id, DetailedError: details}
}

func AppErrorFromJson(js string) *ApplicationError {
	var err ApplicationError
	json.Unmarshal([]byte(js), &err)
	if err.Id == "" {
		return nil
	}
	return &err
}

func MakePermissionError(userID int64) AppError {
	return NewForbiddenError(
		"internal.permissions.check_access.denied",
		fmt.Sprintf("userId=%d, access denied", userID),
	)
}

// MakeScopeError returns an AppError when a user lacks the required scope access.
func MakeScopeError(userID int64, scopeName string, access int) AppError {
	return NewForbiddenError(
		"internal.scope.check_access.denied",
		fmt.Sprintf("access denied for user %d on scope '%s' with access level %d", userID, scopeName, access),
	)
}
