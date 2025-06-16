package errors

import "github.com/webitel/cases/internal/errors"

var (
	ConversionError = errors.NewInternalError("app.process_api.conversion.error", "conversion error occurred")
	DatabaseError   = errors.NewInternalError("app.process_api.db.error", "DB error occurred")
)

func NewBadRequestError(err error) errors.AppError {
	return errors.NewBadRequestError("app.process_api.validation.error", err.Error())
}
