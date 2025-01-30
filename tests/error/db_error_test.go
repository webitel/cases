package tests

import (
	"fmt"
	"testing"

	err "github.com/webitel/cases/internal/error"
)

func TestNewDBError(t *testing.T) {
	id := "db.error"
	message := "Database operation failed"

	err := err.NewDBError(id, message)

	if err.ID != id {
		t.Errorf("expected ID %s, got %s", id, err.ID)
	}
	if err.Message != message {
		t.Errorf("expected Message %s, got %s", message, err.Message)
	}
}

func TestDBError_Error(t *testing.T) {
	id := "db.invalid"
	message := "Invalid database query"
	err := err.NewDBError(id, message)

	expectedError := "DBError [db.invalid]: Invalid database query"
	if err.Error() != expectedError {
		t.Errorf("expected error string '%s', got '%s'", expectedError, err.Error())
	}
}

func TestNewDBNoRowsError(t *testing.T) {
	id := "db.no_rows"
	err := err.NewDBNoRowsError(id)

	if err.ID != id {
		t.Errorf("expected ID %s, got %s", id, err.ID)
	}
	expectedError := "DBError [db.no_rows]: entity does not exist"
	if err.Error() != expectedError {
		t.Errorf("expected error string '%s', got '%s'", expectedError, err.Error())
	}
}

func TestNewDBUniqueViolationError(t *testing.T) {
	id := "db.unique_violation"
	column := "username"
	value := "john_doe"
	err := err.NewDBUniqueViolationError(id, column, value)

	if err.ID != id {
		t.Errorf("expected ID %s, got %s", id, err.ID)
	}
	expectedError := "DBError [db.unique_violation]: invalid input: entity [username = john_doe] already exists"
	if err.Error() != expectedError {
		t.Errorf("expected error string '%s', got '%s'", expectedError, err.Error())
	}
}

func TestNewDBForeignKeyViolationError(t *testing.T) {
	id := "db.foreign_key_violation"
	column := "order_id"
	value := "123"
	foreignKey := "users"
	err := err.NewDBForeignKeyViolationError(id, column, value, foreignKey)

	if err.ID != id {
		t.Errorf("expected ID %s, got %s", id, err.ID)
	}
	expectedError := "DBError [db.foreign_key_violation]: invalid input: violates foreign key constraint"
	if err.Error() != expectedError {
		t.Errorf("expected error string '%s', got '%s'", expectedError, err.Error())
	}
}

func TestNewDBCheckViolationError(t *testing.T) {
	id := "db.check_violation"
	check := "amount > 0"
	err := err.NewDBCheckViolationError(id, check)

	if err.ID != id {
		t.Errorf("expected ID %s, got %s", id, err.ID)
	}
	expectedError := "DBError [db.check_violation]: invalid input: violates check constraint [amount > 0]"
	if err.Error() != expectedError {
		t.Errorf("expected error string '%s', got '%s'", expectedError, err.Error())
	}
}

func TestNewDBInternalError(t *testing.T) {
	id := "db.internal_error"
	reason := fmt.Errorf("database connection lost")
	err := err.NewDBInternalError(id, reason)

	if err.ID != id {
		t.Errorf("expected ID %s, got %s", id, err.ID)
	}
	expectedError := "DBError [db.internal_error]: internal server error"
	if err.Error() != expectedError {
		t.Errorf("expected error string '%s', got '%s'", expectedError, err.Error())
	}
}
