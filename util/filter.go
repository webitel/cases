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

// PartitionFilter splits filters into entries matching `field`and the remaining raw strings.
func PartitionFilter(filters []string, field string) (matched []FilterExpr, rest []string) {
	operators := []string{"!=", ">=", "<=", "=", ">", "<"}
	for _, raw := range filters {
		trimmed := strings.TrimSpace(raw)
		consumed := false
		for _, op := range operators {
			idx := strings.Index(trimmed, op)
			if idx <= 0 {
				continue
			}
			lhs := strings.TrimSpace(trimmed[:idx])
			value := strings.TrimSpace(trimmed[idx+len(op):])
			if lhs == field && value != "" {
				matched = append(matched, FilterExpr{Field: lhs, Operator: op, Value: value})
				consumed = true
			}
			break
		}
		if !consumed {
			rest = append(rest, raw)
		}
	}
	return matched, rest
}
