package error

import (
	"testing"

	"google.golang.org/grpc/codes"
)

func TestNewGRPCError(t *testing.T) {
	code := codes.NotFound
	message := "Resource not found"

	err := NewGRPCError(code, message)

	if err.Code != code {
		t.Errorf("expected Code %s, got %s", code.String(), err.Code.String())
	}
	if err.Message != message {
		t.Errorf("expected Message %s, got %s", message, err.Message)
	}
}

func TestGRPCError_Error(t *testing.T) {
	code := codes.PermissionDenied
	message := "Permission denied"
	err := NewGRPCError(code, message)

	expectedError := "GRPCError [Code: PermissionDenied]: Permission denied"
	if err.Error() != expectedError {
		t.Errorf("expected error string '%s', got '%s'", expectedError, err.Error())
	}
}

func TestNewInvalidArgumentError(t *testing.T) {
	message := "Invalid argument provided"
	err := NewInvalidArgumentError(message)

	if grpcErr, ok := err.(*GRPCError); ok {
		if grpcErr.Code != codes.InvalidArgument {
			t.Errorf("expected Code %s, got %s", codes.InvalidArgument.String(), grpcErr.Code.String())
		}
		if grpcErr.Message != message {
			t.Errorf("expected Message %s, got %s", message, grpcErr.Message)
		}
	} else {
		t.Errorf("expected error type to be *GRPCError")
	}
}
