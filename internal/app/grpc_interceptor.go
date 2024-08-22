package app

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/webitel/cases/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var RequestContextName = "grpc_ctx"

func unaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()
	var reqCtx context.Context
	var ip string

	// Extract metadata from incoming context
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		reqCtx = context.WithValue(ctx, RequestContextName, md)
		ip = getClientIp(md)
	} else {
		ip = "<not found>"
		reqCtx = context.WithValue(ctx, RequestContextName, nil)
	}

	// // Start a new span using the centralized tracer
	// span := trace.SpanFromContext(ctx)
	// defer span.End()

	// Log the start of the request for tracing
	slog.Info("cases.grpc_server.request_started",
		slog.String("method", info.FullMethod),
		slog.Time("start_time", start),
	)

	// Handle the request
	h, err := handler(reqCtx, req)

	// Log the result and record any errors in the span
	if err != nil {
		// span.RecordError(err)
		slog.Error("cases.grpc_server.request_error",
			slog.String("ip", ip),
			slog.String("method", info.FullMethod),
			slog.Duration("duration", time.Since(start)),
			slog.String("error", err.Error()))
		var appError model.AppError
		switch {
		case errors.As(err, &appError):
			var e model.AppError
			errors.As(err, &e)
			return h, status.Error(httpCodeToGrpc(e.GetStatusCode()), e.ToJson())
		default:
			return h, err
		}
	} else {
		slog.Info("cases.grpc_server.request_success",
			slog.String("method", info.FullMethod),
			slog.Duration("duration", time.Since(start)))
	}

	// // Set attributes in the span
	// span.SetAttributes(
	// 	attribute.String("ip", ip),
	// 	attribute.String("method", info.FullMethod),
	// 	attribute.String("status", "success"),
	// 	attribute.String("duration", time.Since(start).String()),
	// )

	return h, err
}

func httpCodeToGrpc(c int) codes.Code {
	switch c {
	case http.StatusBadRequest:
		return codes.InvalidArgument
	case http.StatusAccepted:
		return codes.ResourceExhausted
	case http.StatusUnauthorized:
		return codes.Unauthenticated
	case http.StatusForbidden:
		return codes.PermissionDenied
	default:
		return codes.Internal
	}
}

func getClientIp(info metadata.MD) string {
	ip := strings.Join(info.Get("x-real-ip"), ",")
	if ip == "" {
		ip = strings.Join(info.Get("x-forwarded-for"), ",")
	}

	return ip
}
