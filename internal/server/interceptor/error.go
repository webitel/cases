package interceptor

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/webitel/cases/internal/errors"
	outerror "github.com/webitel/webitel-go-kit/pkg/errors"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AuthUnaryServerInterceptor authenticates and authorizes unary RPCs.
func OuterInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		defer func() {
			if panicErr := recover(); panicErr != nil {
				slog.ErrorContext(ctx, "[PANIC RECOVER]", slog.Any("err", panicErr), slog.String("stack", string(debug.Stack())))
				// TODO: Error returning!
			}
		}()
		resp, err := handler(ctx, req)
		if err != nil {
			return nil, logAndReturnGRPCError(ctx, err, info)
		}
		return resp, nil
	}
}

// logAndReturnGRPCError logs the error and converts it to a gRPC error response.
func logAndReturnGRPCError(ctx context.Context, err error, info *grpc.UnaryServerInfo) error {
	if err == nil {
		return nil
	}
	slog.WarnContext(ctx, fmt.Sprintf("method %s, error: %v", info.FullMethod, err.Error()))
	span := trace.SpanFromContext(ctx) // OpenTelemetry tracing
	span.RecordError(err)

	// Determine the correct gRPC error response
	var (
		grpcCode codes.Code
		httpCode int
		id       string
	)
	slog.ErrorContext(ctx, errors.Details(err))
	switch grpcCode = errors.Code(err); grpcCode {
	case codes.Unauthenticated:
		httpCode = http.StatusUnauthorized
		id = "api.process.unauthenticated"
	case codes.PermissionDenied:
		httpCode = http.StatusForbidden
		id = "api.process.unauthorized"
	case codes.NotFound:
		httpCode = http.StatusNotFound
		id = "api.process.not_found"
	case codes.Aborted, codes.InvalidArgument, codes.AlreadyExists:
		httpCode = http.StatusBadRequest
		id = "api.process.bad_args"
	default:
		httpCode = http.StatusInternalServerError
		id = "api.process.internal"

	}
	grpcErr := &outerror.ApplicationError{
		Id:            id,
		DetailedError: err.Error(),
		StatusCode:    httpCode,
		Status:        http.StatusText(httpCode),
	}
	marshaledErr, _ := json.Marshal(grpcErr)
	return status.Error(grpcCode, string(marshaledErr))

}

// httpCodeToGrpc maps HTTP status codes to gRPC error codes.
func httpCodeToGrpc(c int) codes.Code {
	switch c {
	case http.StatusOK:
		return codes.OK
	case http.StatusBadRequest:
		return codes.InvalidArgument
	case http.StatusUnauthorized:
		return codes.Unauthenticated
	case http.StatusForbidden:
		return codes.PermissionDenied
	case http.StatusNotFound:
		return codes.NotFound
	case http.StatusRequestTimeout:
		return codes.DeadlineExceeded
	case http.StatusConflict:
		return codes.Aborted
	case http.StatusGone:
		return codes.NotFound
	case http.StatusTooManyRequests:
		return codes.ResourceExhausted
	case http.StatusInternalServerError:
		return codes.Internal
	case http.StatusNotImplemented:
		return codes.Unimplemented
	case http.StatusServiceUnavailable:
		return codes.Unavailable
	case http.StatusGatewayTimeout:
		return codes.DeadlineExceeded
	default:
		return codes.Unknown
	}
}
