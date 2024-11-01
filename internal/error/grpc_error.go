// File: grpc_error/grpc_error.go
package error

import (
	"fmt"

	"google.golang.org/grpc/codes"
)

// GRPCError represents a generic gRPC error with a gRPC code.
type GRPCError struct {
	Message string
	Code    codes.Code
}

// NewGRPCError creates a new GRPCError with the specified code and message.
func NewGRPCError(code codes.Code, message string) *GRPCError {
	return &GRPCError{
		Code:    code,
		Message: message,
	}
}

// Error implements the error interface for GRPCError.
func (e *GRPCError) Error() string {
	return fmt.Sprintf("GRPCError [Code: %s]: %s", e.Code.String(), e.Message)
}

// NewInvalidArgumentError creates a new invalid argument error for gRPC responses.
func NewInvalidArgumentError(message string) error {
	return &GRPCError{
		Code:    codes.InvalidArgument,
		Message: message,
	}
}

// NewValidationError creates a new validation error with details about the field and violation.
func NewValidationError(field, constraint, message string) error {
	detailedMessage := fmt.Sprintf("Validation failed on field '%s' [%s]: %s", field, constraint, message)
	return &GRPCError{
		Code:    codes.InvalidArgument,
		Message: detailedMessage,
	}
}
