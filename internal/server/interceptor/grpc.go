package interceptor

import (
	"context"
	"google.golang.org/grpc"
	"log/slog"
	"runtime/debug"
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
