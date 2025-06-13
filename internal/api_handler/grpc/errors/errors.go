package errors

import "github.com/webitel/cases/internal/errors"

func NewBadRequestError(err error) errors.AppError {
	return errors.NewBadRequestError("app.process_api.validation.error", err.Error())
}
