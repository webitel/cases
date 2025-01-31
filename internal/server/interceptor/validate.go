// File: interceptor/validate_unary_server_interceptor.go

package interceptor

import (
	"context"
	"errors"

	"github.com/bufbuild/protovalidate-go" // Make sure protovalidate is in your dependencies
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto" // Required for proto.Message type assertion

	cerr "github.com/webitel/cases/internal/errors"
)

// ValidateUnaryServerInterceptor returns a gRPC interceptor for request validation.
func ValidateUnaryServerInterceptor(val *protovalidate.Validator) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Check if the request implements proto.Message
		if v, ok := req.(proto.Message); ok {
			// Perform validation on the message
			if err := val.Validate(v); err != nil {
				var ve *protovalidate.ValidationError
				// Check if the error is a ValidationError
				if errors.As(err, &ve) && len(ve.Violations) > 0 {
					violation := ve.Violations[0]
					return nil, cerr.NewInternalError(
						violation.GetConstraintId(),
						violation.GetMessage(),
					)
				}
				// Return generic validation error if no specific violations found
				return nil, cerr.NewInternalError("unknown", err.Error())
			}
		}
		// Proceed to handler if validation passes
		return handler(ctx, req)
	}
}
