package app

import (
	"testing"
)

func TestNewRegistryError(t *testing.T) {
	id := "reg001"
	message := "Failed to register service"

	err := NewRegistryError(id, message)

	if err.ID != id {
		t.Errorf("expected ID %s, got %s", id, err.ID)
	}
	if err.Message != message {
		t.Errorf("expected Message %s, got %s", message, err.Message)
	}
}

func TestRegistryError_Error(t *testing.T) {
	id := "reg002"
	message := "Service already registered"
	err := NewRegistryError(id, message)

	expectedError := "RegistryError [reg002]: Service already registered"
	if err.Error() != expectedError {
		t.Errorf("expected error string '%s', got '%s'", expectedError, err.Error())
	}
}
