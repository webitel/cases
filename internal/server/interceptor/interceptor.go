package interceptor

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	api "github.com/webitel/cases/api/cases"

	"github.com/webitel/cases/auth/user_auth"
	autherror "github.com/webitel/cases/internal/error"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
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
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Authorize session with the token
		session, err := authManager.AuthorizeFromContext(ctx)
		if err != nil {
			return nil, autherror.NewUnauthorizedError("auth.session.invalid", fmt.Sprintf("Invalid session or expired token: %v", err))
		}

		// // Retrieve authorization details
		objClass, licenses, action := objClassWithAction(info)

		// // License validation
		if missingLicenses := checkLicenses(session, licenses); len(missingLicenses) > 0 {
			return nil, autherror.NewUnauthorizedError("auth.license.missing", fmt.Sprintf("Missing required licenses: %v", missingLicenses))
		}

		// Permission validation
		if ok, _ := validateSessionPermission(session, objClass, action); !ok {
			return nil, autherror.NewUnauthorizedError("auth.permission.denied", "Permission denied for the requested action")
		}

		ctx = context.WithValue(ctx, SessionHeader, session)

		// Proceed with handler after successful validation
		return handler(ctx, req)
	}
}

// tokenFromContext extracts the authorization token from metadata.
func tokenFromContext(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", autherror.NewUnauthorizedError("auth.metadata.missing", "Metadata is empty; authorization token required")
	}
	token := md.Get(hdrTokenAccess)
	if len(token) < 1 || token[0] == "" {
		return "", autherror.NewUnauthorizedError("auth.token.missing", "Authorization token is missing")
	}
	return token[0], nil
}

func objClassWithAction(info *grpc.UnaryServerInfo) (string, []string, user_auth.AccessMode) {
	serviceName, methodName := splitFullMethodName(info.FullMethod)
	service := api.WebitelAPI[serviceName]
	objClass := service.ObjClass
	licenses := service.AdditionalLicenses
	action := service.WebitelMethods[methodName].Access
	var accessMode user_auth.AccessMode
	switch action {
	case 0:
		accessMode = user_auth.Add
	case 1:
		accessMode = user_auth.Read
	case 2:
		accessMode = user_auth.Edit
	case 3:
		accessMode = user_auth.Delete
	}

	return objClass, licenses, accessMode
}

// checkLicenses verifies that the session has all required licenses.
func checkLicenses(session *user_auth.UserAuthSession, licenses []string) []string {
	var missing []string
	for _, license := range licenses {
		if !session.HasPermission(license) {
			missing = append(missing, license)
		}
	}
	return missing
}

// validateSessionPermission checks if the session has the required permissions.
func validateSessionPermission(session *user_auth.UserAuthSession, objClass string, accessMode user_auth.AccessMode) (bool, bool) {
	scope := session.GetScope(objClass)
	if scope == nil {
		return false, false
	}
	return session.HasObacAccess(scope.Class, accessMode), scope.IsRbacUsed()
}

// splitFullMethodName extracts service and method names from the full gRPC method name.
func splitFullMethodName(fullMethod string) (string, string) {
	fullMethod = strings.TrimPrefix(fullMethod, "/")
	if i := strings.Index(fullMethod, "/"); i >= 0 {
		return reg.ReplaceAllString(fullMethod[:i], ""), fullMethod[i+1:]
	}
	return "unknown", "unknown"
}
