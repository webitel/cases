package error

import (
	"fmt"
)

// ConfigError represents an error that occurs during configuration loading.
type ConfigError struct {
	ID      string
	Message string
}

// NewConfigError creates a new ConfigError with the specified ID and message.
func NewConfigError(id, message string) *ConfigError {
	return &ConfigError{
		ID:      id,
		Message: message,
	}
}

// Error implements the error interface for ConfigError.
func (e *ConfigError) Error() string {
	return fmt.Sprintf("ConfigError [%s]: %s", e.ID, e.Message)
}
