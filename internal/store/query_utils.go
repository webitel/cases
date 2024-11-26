package store

import (
	"fmt"
	"github.com/Masterminds/squirrel"
)

// Ident returns a string that represents a qualified identifier.
// For example, Ident("cc", "case_link") returns "cc.case_link".
var Ident = func(left, right string) string {
	return fmt.Sprintf("%s.%s", left, right)
}

func FormAsCTE(in squirrel.Sqlizer, alias string) (string, []any, error) {
	query, args, err := in.ToSql()
	if err != nil {
		return "", nil, err
	}
	query = fmt.Sprintf("WITH %s AS (%s)", alias, query)
	return query, args, nil
}
