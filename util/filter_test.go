package util

import (
	"reflect"
	"testing"
)

func TestParseFilters(t *testing.T) {
	cases := []struct {
		name    string
		input   []string
		expects []FilterExpr
	}{
		{
			name:    "simple eq",
			input:   []string{"case_id=123"},
			expects: []FilterExpr{{Field: "case_id", Operator: "=", Value: "123"}},
		},
		{
			name:  "multiple ops",
			input: []string{"name!=John", "type>=2", "score<100"},
			expects: []FilterExpr{
				{Field: "name", Operator: "!=", Value: "John"},
				{Field: "type", Operator: ">=", Value: "2"},
				{Field: "score", Operator: "<", Value: "100"},
			},
		},
		{
			name:  "spaces",
			input: []string{" type = bug ", "score > 1 "},
			expects: []FilterExpr{
				{Field: "type", Operator: "=", Value: "bug"},
				{Field: "score", Operator: ">", Value: "1"},
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := ParseFilters(c.input)
			if !reflect.DeepEqual(got, c.expects) {
				t.Errorf("ParseFilters() = %v, want %v", got, c.expects)
			}
		})
	}
}

func TestGetFilter(t *testing.T) {
	filters := []string{"case_id=123", "case_id!=456", "type>=2", "score<100"}
	res := GetFilter(filters, "case_id")
	want := []FilterExpr{
		{Field: "case_id", Operator: "=", Value: "123"},
		{Field: "case_id", Operator: "!=", Value: "456"},
	}
	if !reflect.DeepEqual(res, want) {
		t.Errorf("GetFilter() = %v, want %v", res, want)
	}
}
