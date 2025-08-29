package postgres

import sq "github.com/Masterminds/squirrel"

type Select struct {
	TableAlias                 string
	Joins                      map[string]string
	Query                      sq.SelectBuilder
	FilterSpecialFieldsEncoder map[string]func(any) any
}

func NewSelect(rootTableAlias string, builder sq.SelectBuilder, opts ...BuilderOptions) (*Select, error) {
	query := &Select{
		TableAlias:                 rootTableAlias,
		Joins:                      make(map[string]string),
		Query:                      builder,
		FilterSpecialFieldsEncoder: make(map[string]func(any) any),
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

func WithFiltersEncoder(encoder map[string]func(any) any) BuilderOptions {
	return func(s *Select) error {
		s.FilterSpecialFieldsEncoder = encoder
		return nil
	}
}

func (s *Select) ToSql() (string, []any, error) {
	return s.Query.ToSql()
}
