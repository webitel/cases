package util

import (
	"fmt"
	"strings"
)

func EqualFilter(field string, value any) string {
	return fmt.Sprintf("%s=%v", field, value)
}

type FilterExpr struct {
	Field    string
	Operator string
	Value    string
}

func ParseFilters(filters []string) []FilterExpr {
	var result []FilterExpr
	operators := []string{"!=", ">=", "<=", "=", ">", "<"}
	for _, filter := range filters {
		filter = strings.TrimSpace(filter)
		for _, op := range operators {
			if idx := strings.Index(filter, op); idx > 0 {
				field := strings.TrimSpace(filter[:idx])
				value := strings.TrimSpace(filter[idx+len(op):])
				result = append(result, FilterExpr{
					Field:    field,
					Operator: op,
					Value:    value,
				})
				break
			}
		}
	}
	return result
}

func GetFilter(filters []string, field string) []FilterExpr {
	parsed := ParseFilters(filters)
	var result []FilterExpr
	for _, f := range parsed {
		if f.Field == field && f.Value != "" {
			result = append(result, f)
		}
	}
	return result
}
