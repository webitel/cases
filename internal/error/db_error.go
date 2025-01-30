package error

import (
	"fmt"
	"net/http"

	"github.com/jackc/pgconn"
)

// DBError represents a generic database error.
type DBError struct {
	ID      string
	Message string
	Code    int // HTTP status code
}

// NewDBError creates a new DBError with a default HTTP status code (500 Internal Server Error).
func NewDBError(id, message string) *DBError {
	return &DBError{
		ID:      id,
		Message: message,
		Code:    http.StatusInternalServerError,
	}
}

// Error implements the error interface for DBError.
func (e *DBError) Error() string {
	return fmt.Sprintf("DBError [%s]: %s (HTTP %d)", e.ID, e.Message, e.Code)
}

// DBNoRowsError indicates that no rows were found for a query.
type DBNoRowsError struct {
	DBError
}

func NewDBNoRowsError(id string) *DBNoRowsError {
	err := &DBNoRowsError{
		DBError: *NewDBError(id, "entity does not exist"),
	}
	err.Code = http.StatusNotFound
	return err
}

// DBUniqueViolationError indicates a unique constraint violation.
type DBUniqueViolationError struct {
	DBError
	Column string
	Value  string
}

func NewDBUniqueViolationError(id, column, value string) *DBUniqueViolationError {
	err := &DBUniqueViolationError{
		DBError: *NewDBError(id, fmt.Sprintf("invalid input: entity [%s = %s] already exists", column, value)),
		Column:  column,
		Value:   value,
	}
	err.Code = http.StatusConflict // Override default code
	return err
}

// DBForeignKeyViolationError indicates a foreign key constraint violation.
type DBForeignKeyViolationError struct {
	DBError
	Column          string
	Value           string
	ForeignKeyTable string
}

func NewDBForeignKeyViolationError(id, column, value, foreignKey string) *DBForeignKeyViolationError {
	err := &DBForeignKeyViolationError{
		DBError:         *NewDBError(id, "invalid input: violates foreign key constraint"),
		Column:          column,
		Value:           value,
		ForeignKeyTable: foreignKey,
	}
	err.Code = http.StatusBadRequest // Override default code
	return err
}

// DBCheckViolationError indicates a check constraint violation.
type DBCheckViolationError struct {
	DBError
	Check string
}

func NewDBCheckViolationError(id, check string) *DBCheckViolationError {
	err := &DBCheckViolationError{
		DBError: *NewDBError(id, fmt.Sprintf("invalid input: violates check constraint [%s]", check)),
		Check:   check,
	}
	err.Code = http.StatusBadRequest // Override default code
	return err
}

// DBInternalError indicates an internal database error.
type DBInternalError struct {
	Reason error
	DBError
}

func NewDBInternalError(id string, reason error) *DBInternalError {
	var detailedMessage string

	// Check if the error is a pgconn.PgError to get additional details
	if pgErr, ok := reason.(*pgconn.PgError); ok {
		// Format a detailed error message from the PgError fields
		detailedMessage = fmt.Sprintf("DB Error: %s - %s. %s", pgErr.Message, pgErr.Detail, pgErr.Hint)
	} else {
		// If it's not a PgError, use the generic reason's Error() string
		detailedMessage = reason.Error()
	}

	err := &DBInternalError{
		DBError: *NewDBError(id, detailedMessage),
		Reason:  reason,
	}
	err.Code = http.StatusInternalServerError // Override default code
	return err
}
