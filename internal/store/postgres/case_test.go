package postgres

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseExprToSql(t *testing.T) {
	tests := []struct {
		name     string
		expr     string
		wantSQLs []string
	}{

		{
			name:     "simple equals",
			expr:     "a=1",
			wantSQLs: []string{"c.a = ?"},
		},
		{
			name:     "not equals",
			expr:     "a!=1",
			wantSQLs: []string{"c.a != ?"},
		},
		{
			name:     "and",
			expr:     "a=1&&b=2",
			wantSQLs: []string{"c.a = ?", "c.b = ?"},
		},
		{
			name:     "or",
			expr:     "a=1||b=2",
			wantSQLs: []string{"c.a = ?", "c.b = ?"},
		},
		{
			name:     "complex",
			expr:     "a=1&&b!=2||c=3",
			wantSQLs: []string{"c.a = ?", "c.b != ?", "c.c = ?"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sqlizer := parseExprToSql(tt.expr)
			sql, _, err := sqlizer.ToSql()
			require.NoError(t, err)
			for _, want := range tt.wantSQLs {
				require.True(t, strings.Contains(sql, want), "sql %q does not contain %q", sql, want)
			}
		})
	}
}
