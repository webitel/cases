package postgres

import (
	"testing"

	"reflect"

	"github.com/Masterminds/squirrel"
	"github.com/webitel/cases/util"
	"github.com/webitel/webitel-go-kit/pkg/filters"
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

func TestApplyFilterBool(t *testing.T) {
	tests := []struct {
		name       string
		comparison filters.Comparison
		value      any
		wantSQL    string
		wantArgs   []any
	}{
		{"eq false includes NULL", filters.Equal, false, "xc.flag IS NOT TRUE", nil},
		{"ne true includes NULL", filters.NotEqual, true, "xc.flag IS NOT TRUE", nil},
		{"eq true", filters.Equal, true, "xc.flag = ?", []any{true}},
		{"ne false", filters.NotEqual, false, "xc.flag <> ?", []any{false}},
		{"eq int unaffected", filters.Equal, int64(0), "xc.flag = ?", []any{int64(0)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sqlizer, err := applyFilter(&filters.Filter{Column: "xc.flag", Value: tt.value, ComparisonType: tt.comparison})
			if err != nil {
				t.Fatalf("applyFilter() error: %v", err)
			}
			query, args, err := sqlizer.ToSql()
			if err != nil {
				t.Fatalf("ToSql() error: %v", err)
			}
			if query != tt.wantSQL {
				t.Errorf("applyFilter() query = %q, want %q", query, tt.wantSQL)
			}
			if len(args) != len(tt.wantArgs) || (len(args) > 0 && !reflect.DeepEqual(args, tt.wantArgs)) {
				t.Errorf("applyFilter() args = %v, want %v", args, tt.wantArgs)
			}
		})
	}
}
