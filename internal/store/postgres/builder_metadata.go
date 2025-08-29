package postgres

import sq "github.com/Masterminds/squirrel"

type Select struct {
	TableAlias                 string
	Joins                      map[string]string
	Query                      sq.SelectBuilder
	FilterSpecialFieldsEncoder map[string]func(any) any
}

func NewSelectBuilderMetadata(rootTableAlias string, builder sq.SelectBuilder, filtersEncoder map[string]func(v any) any) *Select {
	return &Select{
		TableAlias:                 rootTableAlias,
		Joins:                      make(map[string]string),
		Query:                      builder,
		FilterSpecialFieldsEncoder: filtersEncoder,
	}
}

func (s *Select) ToSql() (string, []any, error) {
	return s.Query.ToSql()
}
