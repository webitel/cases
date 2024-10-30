package app

import (
	"testing"
)

func TestNewConfigError(t *testing.T) {
	id := "config.load.error"
	message := "Failed to load configuration"

	err := NewConfigError(id, message)

	if err.ID != id {
		t.Errorf("expected ID %s, got %s", id, err.ID)
	}
	if err.Message != message {
		t.Errorf("expected Message %s, got %s", message, err.Message)
	}
}

func TestConfigError_Error(t *testing.T) {
	id := "config.invalid"
	message := "Invalid configuration value"
	err := NewConfigError(id, message)

	expectedError := "ConfigError [config.invalid]: Invalid configuration value"
	if err.Error() != expectedError {
		t.Errorf("expected error string '%s', got '%s'", expectedError, err.Error())
	}
}
