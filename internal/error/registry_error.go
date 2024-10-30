package app

import (
	"fmt"
)

// RegistryError represents an error that occurs during service registration or deregistration.
type RegistryError struct {
	ID      string
	Message string
}

// NewRegistryError creates a new RegistryError with the specified ID and message.
func NewRegistryError(id, message string) *RegistryError {
	return &RegistryError{
		ID:      id,
		Message: message,
	}
}

// Error implements the error interface for RegistryError.
func (e *RegistryError) Error() string {
	return fmt.Sprintf("RegistryError [%s]: %s", e.ID, e.Message)
}
