package postgres

import (
	"testing"

	"github.com/Masterminds/squirrel"
	"github.com/webitel/cases/util"
)

func TestApplyFiltersToQuery(t *testing.T) {
	filters := []util.FilterExpr{
		{Field: "case_id", Operator: "=", Value: "123"},
		{Field: "type", Operator: ">=", Value: "2"},
		{Field: "score", Operator: "<", Value: "100"},
		{Field: "name", Operator: "!=", Value: "John"},
	}
	qb := squirrel.Select("*").From("table")
	qb = ApplyFiltersToQuery(qb, "case_id", []util.FilterExpr{filters[0]})
	qb = ApplyFiltersToQuery(qb, "type", []util.FilterExpr{filters[1]})
	qb = ApplyFiltersToQuery(qb, "score", []util.FilterExpr{filters[2]})
	qb = ApplyFiltersToQuery(qb, "name", []util.FilterExpr{filters[3]})
	query, args, err := qb.ToSql()
	if err != nil {
		t.Fatalf("ToSql() error: %v", err)
	}
	if !(contains(query, "case_id = ") && contains(query, "type >= ") && contains(query, "score < ") && contains(query, "name <> ")) {
		t.Errorf("ApplyFiltersToQuery() query = %v, want contains all filter ops including !=", query)
	}
	if len(args) != 4 {
		t.Errorf("ApplyFiltersToQuery() args = %v, want 4", args)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || (len(s) > len(substr) && (s[:len(substr)] == substr || contains(s[1:], substr))))
}
