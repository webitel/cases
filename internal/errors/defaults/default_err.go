package errors

import "github.com/webitel/cases/internal/errors"

var (
	DatabaseError = errors.NewInternalError(
		"app.process_api.database.perform_query.error",
		"database error occurred",
	)
	ResponseNormalizingError = errors.NewInternalError(
		"app.process_api.response.normalize.error",
		"error occurred while normalizing response",
	)
	ForbiddenError = errors.NewForbiddenError(
		"app.process_api.response.access.error",
		"unable access resource",
	)
	InternalError = errors.NewInternalError(
		"app.process_api.execution.error",
		"error occurred while processing request",
	)
)
