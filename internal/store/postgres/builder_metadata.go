package postgres

import (
	"context"
	"strconv"

	sq "github.com/Masterminds/squirrel"
	"github.com/webitel/cases/internal/errors"
	"github.com/webitel/webitel-go-kit/pkg/filters"
)

type Select struct {
	TableAlias          string
	boundArgs           int
	AppliedJoins        map[string]string
	Query               sq.SelectBuilder
	ColumnValueEncoders map[string]func(any) any
	FilterProcessors    map[string]func(expr *filters.FilterExpr) error
	JoinFunc            func(ctx context.Context, rootTableAlias string, column string) (joins []string, tableAlias string, err error)
}

func NewSelect(rootTableAlias string, builder sq.SelectBuilder, opts ...BuilderOptions) (*Select, error) {
	query := &Select{
		TableAlias:          rootTableAlias,
		AppliedJoins:        make(map[string]string),
		Query:               builder,
		ColumnValueEncoders: make(map[string]func(any) any),
	}
	for _, opt := range opts {
		err := opt(query)
		if err != nil {
			return nil, err
		}
	}
	return query, nil
}

type BuilderOptions func(s *Select) error

func WithBoundArgs(count int) BuilderOptions {
	return func(s *Select) error {
		s.boundArgs = count

		return nil
	}
}

func WithColumnValueEncoders(encoders map[string]func(any) any) BuilderOptions {
	return func(s *Select) error {
		s.ColumnValueEncoders = encoders
		return nil
	}
}

func WithFiltersProcessors(processors map[string]func(expr *filters.FilterExpr) error) BuilderOptions {
	return func(s *Select) error {
		if processors == nil {
			return nil
		}
		s.FilterProcessors = processors
		return nil
	}
}

func WithJoinFunc(joinByColumn func(context.Context, string, string) ([]string, string, error)) BuilderOptions {
	return func(s *Select) error {
		if joinByColumn == nil {
			return errors.New("joinByColumn function is required")
		}
		s.JoinFunc = joinByColumn
		return nil
	}
}

func (s *Select) Join(ctx context.Context, column string) (alias string, err error) {
	if s.JoinFunc == nil {
		return "", errors.New("join function is not defined")
	}
	// check if already joined
	if alias, found := s.AppliedJoins[column]; found {
		return alias, nil
	}
	joins, alias, err := s.JoinFunc(ctx, s.TableAlias, column)
	if err != nil {
		return "", err
	}
	for _, join := range joins {
		s.Query = s.Query.JoinClause(join)
	}
	s.AppliedJoins[column] = alias

	return alias, nil
}

func (s *Select) ToSql() (string, []any, error) {
	query, args, err := s.Query.ToSql()
	if err != nil || s.boundArgs == 0 {
		return query, args, err
	}

	return renumberPlaceholders(query, s.boundArgs), args, nil
}

func renumberPlaceholders(query string, offset int) string {
	out := make([]byte, 0, len(query))

	for i := 0; i < len(query); i++ {
		if query[i] != '?' {
			out = append(out, query[i])

			continue
		}

		if i+1 < len(query) && query[i+1] == '?' {
			out = append(out, '?')

			i++

			continue
		}

		offset++

		out = append(out, '$')
		out = strconv.AppendInt(out, int64(offset), 10)
	}

	return string(out)
}
