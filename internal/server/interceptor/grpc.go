package interceptor

import (
	"context"
	"google.golang.org/grpc"
)

// AuthUnaryServerInterceptor authenticates and authorizes unary RPCs.
func OuterInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		resp, err := handler(ctx, req)
		if err != nil {
			return nil, logAndReturnGRPCError(ctx, err, info)
		}

		return resp, nil
	}
}
