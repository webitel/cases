package interceptor

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"regexp"
	"strings"

	"github.com/webitel/cases/auth"
	"go.opentelemetry.io/otel/trace"

	api "github.com/webitel/cases/api/cases"
	"github.com/webitel/cases/auth/user_auth"
	cerror "github.com/webitel/cases/internal/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// Define a header constant for the token
const (
	hdrTokenAccess = "X-Webitel-Access"
	SessionHeader  = "session"
)

// Regular expression to parse gRPC method information
var reg = regexp.MustCompile(`^(.*\.)`)

// AuthUnaryServerInterceptor authenticates and authorizes unary RPCs.
func AuthUnaryServerInterceptor(authManager user_auth.AuthManager) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Retrieve authorization details
		objClass, licenses, action := objClassWithAction(info)

		// Authorize session with the token
		session, err := authManager.AuthorizeFromContext(ctx, objClass, action)
		if err != nil {
			return nil, cerror.NewUnauthorizedError("auth.session.invalid", fmt.Sprintf("Invalid session or expired token: %v", err))
		}

		//  License validation
		if missingLicenses := checkLicenses(session, licenses); len(missingLicenses) > 0 {
			return nil, cerror.NewPermissionForbiddenError("auth.license.missing", fmt.Sprintf("Missing required licenses: %v", missingLicenses))
		}

		// Permission validation
		if ok := validateSessionPermission(session, objClass, action); !ok {
			return nil, cerror.NewPermissionForbiddenError("auth.permission.denied", "Permission denied for the requested action")
		}

		ctx = context.WithValue(ctx, SessionHeader, session)

		// Proceed with handler after successful validation
		resp, err := handler(ctx, req)
		if err != nil {
			return nil, err
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
	switch e := err.(type) {
	case cerror.AppError:
		return status.Error(httpCodeToGrpc(e.GetStatusCode()), e.ToJson())
	case cerror.AuthError:
		return status.Error(httpCodeToGrpc(e.GetStatusCode()), e.ToJson())
	default:
		slog.ErrorContext(ctx, fmt.Sprintf("not app err returned: %s", err.Error()))
		return status.Error(codes.Internal, cerror.NewInternalError("app.interceptor.parse.error", http.StatusText(http.StatusInternalServerError)).ToJson())
	}
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

// tokenFromContext extracts the authorization token from metadata.
func tokenFromContext(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", cerror.NewUnauthorizedError("auth.metadata.missing", "Metadata is empty; authorization token required")
	}
	token := md.Get(hdrTokenAccess)
	if len(token) < 1 || token[0] == "" {
		return "", cerror.NewUnauthorizedError("auth.token.missing", "Authorization token is missing")
	}
	return token[0], nil
}

func objClassWithAction(info *grpc.UnaryServerInfo) (string, []string, auth.AccessMode) {
	serviceName, methodName := splitFullMethodName(info.FullMethod)
	service := api.WebitelAPI[serviceName]
	objClass := service.ObjClass
	licenses := service.AdditionalLicenses
	action := service.WebitelMethods[methodName].Access
	var accessMode auth.AccessMode
	switch action {
	case 0:
		accessMode = auth.Add
	case 1:
		accessMode = auth.Read
	case 2:
		accessMode = auth.Edit
	case 3:
		accessMode = auth.Delete
	}

	return objClass, licenses, accessMode
}

// checkLicenses verifies that the session has all required licenses.
func checkLicenses(session auth.Auther, licenses []string) []string {
	var missing []string
	for _, license := range licenses {
		if !session.CheckLicenseAccess(license) {
			missing = append(missing, license)
		}
	}
	return missing
}

// validateSessionPermission checks if the session has the required permissions.
func validateSessionPermission(session auth.Auther, objClass string, accessMode auth.AccessMode) bool {
	return session.CheckObacAccess(objClass, accessMode)
}

// splitFullMethodName extracts service and method names from the full gRPC method name.
func splitFullMethodName(fullMethod string) (string, string) {
	fullMethod = strings.TrimPrefix(fullMethod, "/")
	if i := strings.Index(fullMethod, "/"); i >= 0 {
		return reg.ReplaceAllString(fullMethod[:i], ""), fullMethod[i+1:]
	}
	return "unknown", "unknown"
}
