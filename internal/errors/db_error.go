package errors

import (
	"fmt"
	"sync"

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

func NewDBUniqueViolationError(id, column, value string) *DBUniqueViolationError {
	return &DBUniqueViolationError{
		DBError: *NewDBError(id, fmt.Sprintf("invalid input: entity [%s = %s] already exists", column, value)),
		Column:  column,
		Value:   value,
	}
}

// DBForeignKeyViolationError indicates a foreign key constraint violation.
type DBForeignKeyViolationError struct {
	DBError
	Column          string
	Value           string
	ForeignKeyTable string
}

func NewDBForeignKeyViolationError(id, column, value, foreignKey string) *DBForeignKeyViolationError {
	return &DBForeignKeyViolationError{
		DBError:         *NewDBError(id, "invalid input: violates foreign key constraint"),
		Column:          column,
		Value:           value,
		ForeignKeyTable: foreignKey,
	}
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

func NewDBNotNullViolationError(id, table, column string) *DBNotNullViolationError {
	return &DBNotNullViolationError{
		DBError: *NewDBError(id, fmt.Sprintf("invalid input: violates not null constraint: column [%s.%s] cannot be null", table, column)),
		Table:   table,
		Column:  column,
	}
}

// DBEntityConflictError indicates a conflict in entity requests.
type DBEntityConflictError struct {
	DBError
}

func NewDBEntityConflictError(id string) *DBEntityConflictError {
	return &DBEntityConflictError{DBError: *NewDBError(id, "found more than one requested entity")}
}

// DBConflictError indicates a conflict in the database operation (e.g., version mismatch).
type DBConflictError struct {
	DBError
}

// NewDBConflictError creates a new DBConflictError with the specified ID and message.
func NewDBConflictError(id, message string) *DBConflictError {
	return &DBConflictError{
		DBError: *NewDBError(id, message),
	}
}

// DBForbiddenError indicates that the user is forbidden from performing an action.
type DBForbiddenError struct {
	DBError
}

// NewDBForbiddenError creates a new DBForbiddenError with the specified ID and message.
func NewDBForbiddenError(id, message string) *DBForbiddenError {
	return &DBForbiddenError{
		DBError: *NewDBError(id, message),
	}
}

// DBInternalError indicates an internal database error.
type DBInternalError struct {
	Reason error
	DBError
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
	if pgErr, ok := reason.(*pgconn.PgError); ok {
		// Format a detailed error message from the PgError fields
		detailedMessage = fmt.Sprintf("DB Error: %s - %s. %s", pgErr.Message, pgErr.Detail, pgErr.Hint)
	} else {
		// If it's not a PgError, use the generic reason's Error() string
		detailedMessage = reason.Error()
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

// Constraint registration for custom check violations.
var (
	checkViolationErrorRegistry = map[string]string{}
	constraintMu                sync.RWMutex
)

// RegisterConstraint registers custom database check constraints with a custom message.
func RegisterConstraint(name, message string) {
	constraintMu.Lock()
	defer constraintMu.Unlock()
	if _, dup := checkViolationErrorRegistry[name]; dup {
		panic("RegisterConstraint called twice for name " + name)
	}
	checkViolationErrorRegistry[name] = message
}
