package app

import (
	"fmt"
)

// AuthError represents a generic authentication error.
type AuthError struct {
	ID      string
	Message string
}

// NewAuthError creates a new AuthError with the specified ID and message.
func NewAuthError(id, message string) *AuthError {
	return &AuthError{
		ID:      id,
		Message: message,
	}
}

// Error implements the error interface for AuthError.
func (e *AuthError) Error() string {
	return fmt.Sprintf("AuthError [%s]: %s", e.ID, e.Message)
}

// UnauthorizedError indicates an unauthorized access attempt.
type UnauthorizedError struct {
	AuthError
}

func NewUnauthorizedError(id string, message string) *UnauthorizedError {
	return &UnauthorizedError{AuthError: *NewAuthError(id, message)}
}
