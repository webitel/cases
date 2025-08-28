package util

import (
	"fmt"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/webitel/cases/internal/model/options"
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
		stack = []filters.Filterer{nodes}
	)

	for len(stack) > 0 {
		// Pop from stack
		current := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		switch data := current.(type) {
		case *filters.FilterNode:
			// Add all child nodes to stack (in reverse order to maintain processing order)
			for i := len(data.Nodes) - 1; i >= 0; i-- {
				stack = append(stack, data.Nodes[i])
			}
		case *filters.Filter:
			var (
				splittedNaming = strings.Split(data.Column, ".")
				value          = data.Value
			)
			if encoder, ok := base.FilterSpecialFieldsEncoder[data.Column]; ok {
				data.Value = encoder(value)
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
				data.Column = Ident(base.TableAlias, column)
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
				data.Column = Ident(fkTable, referencedColumn)

			default:
				return fmt.Errorf("unsupported nest depth, max 1 level of nesting")
			}

		default:
			return fmt.Errorf("unsupported filter type: %T", current)
		}
	}

	return nil
}

func ApplyFilters(base sq.SelectBuilder, filters filters.Filterer) (sq.SelectBuilder, error) {
	parsedFilters, err := ParseFilters(filters)
	if err != nil {
		return base, err
	}
	return base.Where(parsedFilters), nil
}

func ParseFilters(nodes filters.Filterer) (sq.Sqlizer, error) {
	if nodes == nil {
		return sq.Expr("1=1"), nil
	}
	var (
		res sq.Sqlizer
	)
	switch data := nodes.(type) {
	case *filters.FilterNode:
		switch data.Connection {
		case filters.And:
			and := sq.And{}
			for _, bunch := range data.Nodes {
				switch bunchType := bunch.(type) {
				case *filters.FilterNode:
					lowerResult, err := ParseFilters(bunchType)
					if err != nil {
						return nil, err
					}
					and = append(and, lowerResult)
				case *filters.Filter:
					filter, err := applyFilter(bunchType)
					if err != nil {
						return nil, err
					}
					and = append(and, filter)
				}

			}
			res = and
		case filters.Or:
			or := sq.Or{}
			for _, bunch := range data.Nodes {
				switch v := bunch.(type) {
				case *filters.FilterNode:
					lowerResult, err := ParseFilters(v)
					if err != nil {
						return nil, err
					}
					or = append(or, lowerResult)
				case *filters.Filter:
					filter, err := applyFilter(v)
					if err != nil {
						return nil, err
					}
					or = append(or, filter)
				}
			}
			res = or
		default:
			return nil, fmt.Errorf("invalid connection type in filter node: %d", data.Connection)
		}
	case *filters.Filter:
		filter, err := applyFilter(data)
		if err != nil {
			return nil, err
		}
		return filter, nil
	default:
		return nil, fmt.Errorf("unsupported filter type: %T", nodes)
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
