// File: interceptor/logging_unary_server_interceptor.go

package interceptor

import (
	"context"
	"log/slog"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// LoggingUnaryServerInterceptor logs the details of each request and its duration.
func LoggingUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()

		// Get client IP from metadata
		var ip string
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			ip = getClientIp(md)
		}
		if ip == "" {
			ip = "<unknown>"
		}

		// Log the start of the request
		slog.Info("gRPC request started",
			slog.String("method", info.FullMethod),
			slog.String("client_ip", ip),
			slog.Time("start_time", start),
		)

		// Process the request
		h, err := handler(ctx, req)

		// Log the end of the request with the duration
		duration := time.Since(start)
		if err != nil {
			slog.ErrorContext(ctx, "gRPC request error",
				slog.String("method", info.FullMethod),
				slog.String("client_ip", ip),
				slog.Duration("duration", duration),
				slog.String("error", err.Error()),
			)
		} else {
			slog.Info("gRPC request successful",
				slog.String("method", info.FullMethod),
				slog.String("client_ip", ip),
				slog.Duration("duration", duration),
			)
		}

		return h, err
	}
}

// Helper function to get client IP from metadata
func getClientIp(md metadata.MD) string {
	ip := strings.Join(md.Get("x-real-ip"), ",")
	if ip == "" {
		ip = strings.Join(md.Get("x-forwarded-for"), ",")
	}
	return ip
}
