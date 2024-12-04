package store

import (
	"fmt"
	"strings"

	"github.com/Masterminds/squirrel"
)

const (
	ComparisonILike  = "ilike"
	ComparisonRegexp = "~"
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
	query = fmt.Sprintf("%s AS (%s)", alias, query)
	return query, args, nil
}

// ParseSearchTerm delimit searches for the regexp search indicators and if found returns string without indicators and indicator that regexp search found.
// Returns changed copy of the input slice.
func ParseSearchTerm(q string) (s string, operator string) {
	var (
		escapePre = "/"
		escapeSu  = "/"
	)
	if strings.HasPrefix(q, escapePre) && strings.HasSuffix(q, escapeSu) {
		pre, _ := strings.CutPrefix(q, escapePre)
		su, _ := strings.CutSuffix(pre, escapeSu)
		return su, ComparisonRegexp
	} else {
		return "%" + q + "%", ComparisonILike
	}
}

func AddSearchTerm(base squirrel.SelectBuilder, q string, columns ...string) squirrel.SelectBuilder {
	search, operator := ParseSearchTerm(q)
	for _, column := range columns {
		base = base.Where(fmt.Sprintf("%s %s ?", column, operator), search)
	}
	return base
}
