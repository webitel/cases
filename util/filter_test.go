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

func TestPartitionFilter(t *testing.T) {
	cases := []struct {
		name        string
		input       []string
		field       string
		wantMatched []FilterExpr
		wantRest    []string
	}{
		{
			name:        "single match strips entry",
			input:       []string{"fts=hello", "status=1"},
			field:       "fts",
			wantMatched: []FilterExpr{{Field: "fts", Operator: "=", Value: "hello"}},
			wantRest:    []string{"status=1"},
		},
		{
			name:        "no match keeps everything",
			input:       []string{"status=1", "type=bug"},
			field:       "fts",
			wantMatched: nil,
			wantRest:    []string{"status=1", "type=bug"},
		},
		{
			name:        "empty value not matched and not stripped",
			input:       []string{"fts=", "status=1"},
			field:       "fts",
			wantMatched: nil,
			wantRest:    []string{"fts=", "status=1"},
		},
		{
			name:  "multiple matches all stripped",
			input: []string{"fts=foo", "fts!=bar", "status=1"},
			field: "fts",
			wantMatched: []FilterExpr{
				{Field: "fts", Operator: "=", Value: "foo"},
				{Field: "fts", Operator: "!=", Value: "bar"},
			},
			wantRest: []string{"status=1"},
		},
		{
			name:        "spaces handled",
			input:       []string{" fts = query text ", "status=1"},
			field:       "fts",
			wantMatched: []FilterExpr{{Field: "fts", Operator: "=", Value: "query text"}},
			wantRest:    []string{"status=1"},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			matched, rest := PartitionFilter(c.input, c.field)
			if !reflect.DeepEqual(matched, c.wantMatched) {
				t.Errorf("matched = %v, want %v", matched, c.wantMatched)
			}
			if !reflect.DeepEqual(rest, c.wantRest) {
				t.Errorf("rest = %v, want %v", rest, c.wantRest)
			}
		})
	}
}
