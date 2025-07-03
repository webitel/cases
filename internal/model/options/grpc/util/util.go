package util

import (
	"context"
	"github.com/webitel/cases/auth"
	"github.com/webitel/cases/internal/server/interceptor"
	"strings"
)

func DeduplicateMaskPrefixes(mask []string) []string {
	uniquePrefixes := make(map[string]struct{})
	var trimmedMask []string
	for _, field := range mask {
		prefix := field
		if dotIndex := strings.Index(field, "."); dotIndex > 0 {
			prefix = field[:dotIndex]
		}
		if _, exists := uniquePrefixes[prefix]; !exists {
			uniquePrefixes[prefix] = struct{}{}
			trimmedMask = append(trimmedMask, prefix)
		}
	}
	return trimmedMask
}

const (
	AppServiceName = "cases"
	NamespaceName  = "webitel"
)

func GetAutherOutOfContext(ctx context.Context) auth.Auther {
	return ctx.Value(interceptor.SessionHeader).(auth.Auther)
}
