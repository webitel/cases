package app

import (
	"context"

	"github.com/webitel/cases/auth/user_auth"
)

// grpcContextKey is a unique type for context keys to avoid collisions.
type grpcContextKey struct{}

// GRPCServerContext holds the custom context data for gRPC requests.
type GRPCServerContext struct {
	SignedInUser *user_auth.User
	RequestId    string
}

// FromContext retrieves GRPCServerContext from the given context, or returns an empty GRPCServerContext if none is found.
func FromContext(ctx context.Context) *GRPCServerContext {
	grpcCtx, ok := ctx.Value(grpcContextKey{}).(*GRPCServerContext)
	if !ok {
		return &GRPCServerContext{}
	}
	return grpcCtx
}

// SetUser adds a SignedInUser to the GRPCServerContext and returns a new context with this data.
func SetUser(ctx context.Context, user *user_auth.User) context.Context {
	grpcCtx := FromContext(ctx)
	if grpcCtx == nil {
		grpcCtx = &GRPCServerContext{}
	}
	grpcCtx.SignedInUser = user
	return context.WithValue(ctx, grpcContextKey{}, grpcCtx)
}

// SetRequestId adds a RequestId to the GRPCServerContext and returns a new context with this data.
func SetRequestId(ctx context.Context, requestId string) context.Context {
	grpcCtx := FromContext(ctx)
	if grpcCtx == nil {
		grpcCtx = &GRPCServerContext{}
	}
	grpcCtx.RequestId = requestId
	return context.WithValue(ctx, grpcContextKey{}, grpcCtx)
}
