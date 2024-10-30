package app

import (
	"testing"
)

func TestNewAuthError(t *testing.T) {
	id := "auth001"
	message := "Unauthorized access"

	err := NewAuthError(id, message)

	if err.ID != id {
		t.Errorf("expected ID %s, got %s", id, err.ID)
	}
	if err.Message != message {
		t.Errorf("expected message %s, got %s", message, err.Message)
	}
}

func TestAuthError_Error(t *testing.T) {
	id := "auth002"
	message := "Invalid credentials"
	err := NewAuthError(id, message)

	expected := "AuthError [auth002]: Invalid credentials"
	if err.Error() != expected {
		t.Errorf("expected error string '%s', got '%s'", expected, err.Error())
	}
}

func TestNewUnauthorizedError(t *testing.T) {
	id := "auth003"
	message := "User not authorized"

	err := NewUnauthorizedError(id, message)

	if err.ID != id {
		t.Errorf("expected ID %s, got %s", id, err.ID)
	}
	if err.Message != message {
		t.Errorf("expected message %s, got %s", message, err.Message)
	}
}

func TestUnauthorizedError_Error(t *testing.T) {
	id := "auth004"
	message := "Access denied"
	err := NewUnauthorizedError(id, message)

	expected := "AuthError [auth004]: Access denied"
	if err.Error() != expected {
		t.Errorf("expected error string '%s', got '%s'", expected, err.Error())
	}
}
