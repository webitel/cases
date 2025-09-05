package postgres

import (
	"fmt"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/webitel/cases/internal/model/options"
	util2 "github.com/webitel/cases/internal/store/util"
	"github.com/webitel/cases/util"
	"github.com/webitel/webitel-go-kit/pkg/filters"
)

// ApplyFiltersToQuery applies the filters to the given SelectBuilder query.
func ApplyFiltersToQuery(qb sq.SelectBuilder, column string, filters []util.FilterExpr) sq.SelectBuilder {
	for _, f := range filters {
		switch f.Operator {
		case "=":
			qb = qb.Where(sq.Eq{column: f.Value})
		case "!=":
			qb = qb.Where(sq.NotEq{column: f.Value})
		case ">":
			qb = qb.Where(column+" > ?", f.Value)
		case "<":
			qb = qb.Where(column+" < ?", f.Value)
		case ">=":
			qb = qb.Where(column+" >= ?", f.Value)
		case "<=":
			qb = qb.Where(column+" <= ?", f.Value)
		}
	}
	return qb
}

// NormalizeFilters normalizes the filters by applying the join function to each filter and changing column names that they become valid sql in format: "table.column".
func NormalizeFilters(base *Select, opts options.Searcher, join func(options.Searcher, *Select, string) (string, error)) error {
	if opts.GetFiltersV1() == nil {
		return nil
	}

	// Use a stack to process filters iteratively
	var (
		nodes = opts.GetFiltersV1()
		stack = []*filters.FilterExpr{nodes}
	)

	for len(stack) > 0 {
		// Pop from stack
		current := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if current == nil {
			continue
		}
		if node := current.GetFilterNode(); node != nil {
			// Add all child nodes to stack (in reverse order to maintain processing order)
			for i := len(node.Nodes) - 1; i >= 0; i-- {
				stack = append(stack, node.Nodes[i])
			}
		} else if filter := current.GetFilter(); filter != nil {
			var (
				splittedNaming = strings.Split(filter.Column, ".")
				value          = filter.Value
			)
			if encoder, ok := base.FilterSpecialFieldsEncoder[filter.Column]; ok {
				filter.Value = encoder(value)
			}
			if len(splittedNaming) == 0 || splittedNaming[0] == "" {
				// if column is empty, we cannot apply filter
				return fmt.Errorf("no filter column name")
			}
			switch len(splittedNaming) {
			case 1: // not nested, just column name
				var (
					column = splittedNaming[0]
				)
				filter.Column = util2.Ident(base.TableAlias, column)
			case 2: // nested, table.column
				var (
					fkTable          string
					referencedColumn string
					err              error
					found            bool
				)
				fkColumn := splittedNaming[0]
				referencedColumn = splittedNaming[1]

				if fkTable, found = base.Joins[fkColumn]; !found {
					fkTable, err = join(opts, base, fkColumn)
					if err != nil {
						return err
					}
				}
				filter.Column = util2.Ident(fkTable, referencedColumn)

			default:
				return fmt.Errorf("unsupported nest depth, max 1 level of nesting")
			}
		}
	}

	return nil
}

func ApplyFilters(base sq.SelectBuilder, filters *filters.FilterExpr) (sq.SelectBuilder, error) {
	parsedFilters, err := ParseFilters(filters)
	if err != nil {
		return base, err
	}
	return base.Where(parsedFilters), nil
}

func ParseFilters(expr *filters.FilterExpr) (sq.Sqlizer, error) {
	if expr == nil {
		return sq.Expr("1=1"), nil
	}
	var (
		res        sq.Sqlizer
		parseNodes = func(nodes []*filters.FilterExpr) ([]sq.Sqlizer, error) {
			var sqlizers []sq.Sqlizer
			for _, nestedExpr := range nodes {
				lowerResult, err := ParseFilters(nestedExpr)
				if err != nil {
					return nil, err
				}
				sqlizers = append(sqlizers, lowerResult)

			}
			return sqlizers, nil
		}
	)
	if data := expr.GetFilterNode(); data != nil {
		lowerResult, err := parseNodes(data.Nodes)
		if err != nil {
			return nil, err
		}
		switch data.Connection {
		case filters.And:
			and := append(sq.And{}, lowerResult...)
			res = and
		case filters.Or:
			or := append(sq.Or{}, lowerResult...)
			res = or
		default:
			return nil, fmt.Errorf("invalid connection type in filter node: %d", data.Connection)
		}
	} else if filter := expr.GetFilter(); filter != nil {
		appliedFilter, err := applyFilter(filter)
		if err != nil {
			return nil, err
		}
		return appliedFilter, nil
	}
	return res, nil
}

// Apply filter performs conversion between model.Filter and sq.Sqlizer.
func applyFilter(filter *filters.Filter) (sq.Sqlizer, error) {
	if filter == nil {
		return sq.Expr("1=1"), nil
	}
	var (
		columnName = filter.Column
		value      = filter.Value
	)

	var result sq.Sqlizer
	switch filter.ComparisonType {
	case filters.GreaterThan:
		result = sq.Gt{columnName: value}
	case filters.GreaterThanOrEqual:
		result = sq.GtOrEq{columnName: value}
	case filters.LessThan:
		result = sq.Lt{columnName: value}
	case filters.LessThanOrEqual:
		result = sq.LtOrEq{columnName: value}
	case filters.NotEqual:
		result = sq.NotEq{columnName: value}
	case filters.Like:
		result = sq.Like{columnName: value}
	case filters.ILike:
		result = sq.ILike{columnName: value}
	case filters.Equal:
		result = sq.Eq{columnName: value}
	default:
		return nil, fmt.Errorf("invalid filter type: %d", filter.ComparisonType)
	}
	return result, nil
}

var (
	timeEncoder = func(v any) any {
		switch t := v.(type) {
		case nil:
			return nil
		case time.Time:
			return t
		case *time.Time:
			if t == nil {
				return nil
			}
			return *t
		case int64:
			return time.Unix(t, 0)
		case int:
			return time.Unix(int64(t), 0)
		default:
			return nil
		}
	}
)
