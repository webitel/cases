package util

import (
	"fmt"
	"github.com/Masterminds/squirrel"
	grpcopts "github.com/webitel/cases/model/options/grpc"
	"strings"
)

const (
	ComparisonILike  = "ilike"
	ComparisonRegexp = "~"
	SortAsc          = "ASC"
	SortDesc         = "DESC"
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

func FormAsCTEs(in map[string]squirrel.Sqlizer) (string, []any, error) {
	var (
		i              int
		resultingQuery string
		resultingArgs  []any
	)
	for alias, sqlizer := range in {
		query, args, _ := sqlizer.ToSql()
		if i == 0 {
			// init
			resultingQuery = fmt.Sprintf("WITH %s AS (%s)", alias, query)
			resultingArgs = args
		} else {
			resultingQuery += fmt.Sprintf("%s AS (%s)", alias, query)
			resultingArgs = append(resultingArgs, args...)
		}

		if len(in)-1 != i {
			resultingQuery += ","
		}
		i++
	}

	return resultingQuery, resultingArgs, nil
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

func ApplyPaging(page int, size int, base squirrel.SelectBuilder) squirrel.SelectBuilder {
	if size > 0 {
		base = base.Limit(uint64(size + 1))
		if page > 1 {
			base = base.Offset(uint64((page - 1) * size))
		}
	}

	return base
}

func ResolvePaging[T any](size int, items []*T) (updatedItems []*T, next bool) {
	updatedItems = items[:]
	if size > 0 {
		if len(updatedItems) > size {
			updatedItems = updatedItems[:size]
			next = true
		}
	}
	return
}

func ApplyDefaultSorting(opts grpcopts.Sorter, base squirrel.SelectBuilder, defaultSort string) squirrel.SelectBuilder {
	if s := opts.GetSort(); len(s) != 0 {

		// Check for + or - prefix
		desc := strings.HasPrefix(s, "-")
		asc := strings.HasPrefix(s, "+")

		// Trim prefix if it exists
		if desc || asc {
			s = strings.TrimPrefix(s, string(s[0]))
		}

		// Determine sort direction
		if desc {
			s += " DESC"
		} else {
			s += " ASC"
		}
		base = base.OrderBy(s)

	} else {
		base = base.OrderBy(fmt.Sprintf(`%s ASC`, defaultSort))
	}

	return base
}

func GetSortingOperator(sort string) (field, direction string) {
	if len(sort) != 0 {

		// Check for + or - prefix
		desc := strings.HasPrefix(sort, "-")
		asc := strings.HasPrefix(sort, "+")

		// Trim prefix if it exists
		if desc || asc {
			sort = strings.TrimPrefix(sort, string(sort[0]))
		}
		var dir string
		// Determine sort direction
		if desc {
			dir = SortDesc
		} else {
			dir = SortAsc
		}
		return sort, dir

	}

	return "", ""
}

// PrepareSearchNumber reverses the first string
func PrepareSearchNumber(number string) string {
	if number != "" {
		r := []rune(number)
		for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
			r[i], r[j] = r[j], r[i]
		}
		return string(r)

	}
	return number
}
