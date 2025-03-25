package errors

import (
	"errors"
	"fmt"
	"github.com/jackc/pgconn"
)

// DBError represents a generic database error.
type DBError struct {
	ID      string
	Message string
}

// NewDBError creates a new DBError with the specified ID and message.
func NewDBError(id, message string) *DBError {
	return &DBError{
		ID:      id,
		Message: message,
	}
}

// Error implements the error interface for DBError.
func (e *DBError) Error() string {
	return fmt.Sprintf("DBError [%s]: %s", e.ID, e.Message)
}

// DBNoRowsError indicates that no rows were found for a query.
type DBNoRowsError struct {
	DBError
}

func NewDBNoRowsError(id string) *DBNoRowsError {
	return &DBNoRowsError{DBError: *NewDBError(id, "entity does not exist")}
}

// DBUniqueViolationError indicates a unique constraint violation.
type DBUniqueViolationError struct {
	DBError
	Column string
	Value  string
}

// DBForeignKeyViolationError indicates a foreign key constraint violation.
type DBForeignKeyViolationError struct {
	DBError
	Column          string
	Value           string
	ForeignKeyTable string
}

// DBCheckViolationError indicates a check constraint violation.
type DBCheckViolationError struct {
	DBError
	Check string
}

func NewDBCheckViolationError(id, check string) *DBCheckViolationError {
	return &DBCheckViolationError{
		DBError: *NewDBError(id, fmt.Sprintf("invalid input: violates check constraint [%s]", check)),
		Check:   check,
	}
}

// DBNotNullViolationError indicates a not-null constraint violation.
type DBNotNullViolationError struct {
	DBError
	Table  string
	Column string
}

// DBEntityConflictError indicates a conflict in entity requests.
type DBEntityConflictError struct {
	DBError
}

// DBConflictError indicates a conflict in the database operation (e.g., version mismatch).
type DBConflictError struct {
	DBError
}

// DBForbiddenError indicates that the user is forbidden from performing an action.
type DBForbiddenError struct {
	DBError
}

// DBInternalError indicates an internal database error.
type DBInternalError struct {
	Reason error
	DBError
}

// DBBadRequestError indicates a bad request due to invalid or missing input.
type DBBadRequestError struct {
	DBError
	MissingParam string
}

// NewDBBadRequestError creates a new DBBadRequestError with the specified ID and missing fields.
func NewDBBadRequestError(id string, param string) *DBBadRequestError {
	message := fmt.Sprintf("missing or invalid required params: %v", param)
	return &DBBadRequestError{
		DBError:      *NewDBError(id, message),
		MissingParam: param,
	}
}

// Error implements the error interface for DBBadRequestError.
func (e *DBBadRequestError) Error() string {
	return fmt.Sprintf("DBBadRequestError [%s]: %s (MissingParam: %v)", e.ID, e.Message, e.MissingParam)
}

// Error implements the error interface for DBInternalError.
func (d *DBInternalError) Error() string {
	if d.Reason != nil {
		return fmt.Sprintf("DBInternalError [%s]: %s (Reason: %s)", d.ID, d.Message, d.Reason.Error())
	}
	return fmt.Sprintf("DBInternalError [%s]: %s", d.ID, d.Message)
}

func NewDBInternalError(id string, reason error) *DBInternalError {
	var detailedMessage string

	// Check if the error is a pgconn.PgError to get additional details
	var pgErr *pgconn.PgError
	if errors.As(reason, &pgErr) {
		// Format a detailed error message from the PgError fields
		detailedMessage = fmt.Sprintf("DB Error: %s - %s. %s", pgErr.Message, pgErr.Detail, pgErr.Hint)
	}

	return &DBInternalError{
		DBError: *NewDBError(id, detailedMessage), // Use the detailed message as the error message
		Reason:  reason,
	}
}

// DBNotFoundError indicates that a specific entity was not found.
type DBNotFoundError struct {
	DBError
}

// NewDBNotFoundError creates a new DBNotFoundError with the specified ID and message.
func NewDBNotFoundError(id, message string) *DBNotFoundError {
	return &DBNotFoundError{
		DBError: *NewDBError(id, message),
	}
}
