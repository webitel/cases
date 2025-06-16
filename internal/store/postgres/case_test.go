package postgres

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseComplexFilter(t *testing.T) {
	tests := []struct {
		name     string
		expr     string
		wantSQLs []string
		hasArgs  bool // whether this filter should have arguments
	}{
		{
			name:     "simple equals",
			expr:     "created_by=3",
			wantSQLs: []string{"c.created_by = ANY(?::int[])"},
			hasArgs:  true,
		},
		{
			name:     "not equals",
			expr:     "created_by!=3",
			wantSQLs: []string{"c.created_by != ?"},
			hasArgs:  true,
		},
		{
			name:     "and condition",
			expr:     "created_by=3&&author=4",
			wantSQLs: []string{"c.created_by = ANY(?::int[])", "auth.id = ANY(?::int[])"},
			hasArgs:  true,
		},
		{
			name:     "or condition",
			expr:     "created_by=3||author=4",
			wantSQLs: []string{"c.created_by = ANY(?::int[])", "auth.id = ANY(?::int[])"},
			hasArgs:  true,
		},
		{
			name:     "complex condition",
			expr:     "created_by=3&&author!=4||status=5",
			wantSQLs: []string{"c.created_by = ANY(?::int[])", "auth.id != ?", "c.status = ANY(?::int[])"},
			hasArgs:  true,
		},
		{
			name:     "multiple values",
			expr:     "created_by=3,4,5",
			wantSQLs: []string{"c.created_by = ANY(?::int[])"},
			hasArgs:  true,
		},
		{
			name:     "null value",
			expr:     "created_by=null",
			wantSQLs: []string{"c.created_by ISNULL"},
			hasArgs:  false,
		},
		{
			name:     "not null value",
			expr:     "created_by!=null",
			wantSQLs: []string{"c.created_by NOTNULL"},
			hasArgs:  false,
		},
		{
			name:     "multiple conditions with null",
			expr:     "created_by=3&&author=null||status!=null",
			wantSQLs: []string{"c.created_by = ANY(?::int[])", "auth.id ISNULL", "c.status NOTNULL"},
			hasArgs:  true,
		},
		{
			name:     "attachments true",
			expr:     "attachments=true",
			wantSQLs: []string{"EXISTS (SELECT id FROM storage.files WHERE uuid = c.id::varchar UNION SELECT id FROM cases.case_link WHERE case_link.case_id = c.id)"},
			hasArgs:  false,
		},
		{
			name:     "attachments not true",
			expr:     "attachments!=true",
			wantSQLs: []string{"NOT EXISTS (SELECT id FROM storage.files WHERE uuid = c.id::varchar UNION SELECT id FROM cases.case_link WHERE case_link.case_id = c.id)"},
			hasArgs:  false,
		},
		{
			name:     "empty filter",
			expr:     "",
			wantSQLs: []string{"1=0"},
			hasArgs:  false,
		},
		{
			name:     "invalid filter",
			expr:     "invalid",
			wantSQLs: []string{"1=0"},
			hasArgs:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &CaseStore{}
			var err error
			sqlizer := store.parseComplexFilter(tt.expr, nil, &customFilterContext{}, &err)

			// Print test case details
			t.Logf("Test case: %s", tt.name)
			t.Logf("Input expression: %q", tt.expr)
			t.Logf("Expected SQL fragments: %v", tt.wantSQLs)

			// Check for errors
			if err != nil {
				t.Logf("Error: %v", err)
				if tt.expr != "invalid" && tt.expr != "" {
					require.NoError(t, err, "unexpected error for valid filter")
				}
				return
			}

			// Get SQL and arguments
			sql, args, err := sqlizer.ToSql()
			require.NoError(t, err, "failed to generate SQL")

			t.Logf("Generated SQL: %q", sql)
			t.Logf("Arguments: %v", args)

			// Check SQL fragments
			for _, want := range tt.wantSQLs {
				require.True(t, strings.Contains(sql, want),
					"SQL %q does not contain expected fragment %q", sql, want)
			}

			// Check arguments
			if tt.hasArgs {
				require.NotEmpty(t, args, "expected arguments for filter %q", tt.expr)
				t.Logf("Number of arguments: %d", len(args))
			} else {
				require.Empty(t, args, "unexpected arguments for filter %q", tt.expr)
			}

			// Additional validation for specific cases
			if strings.Contains(tt.expr, "null") {
				if strings.Contains(tt.expr, "!=") {
					require.Contains(t, sql, "NOTNULL", "expected NOTNULL for not null condition")
				} else {
					require.Contains(t, sql, "ISNULL", "expected ISNULL for null condition")
				}
			}

			if strings.Contains(tt.expr, "attachments") {
				if (strings.Contains(tt.expr, "=true") && !strings.Contains(tt.expr, "!=")) ||
					(strings.Contains(tt.expr, "!=false")) {
					require.Contains(t, sql, "EXISTS", "expected EXISTS for attachments=true or attachments!=false")
				} else if (strings.Contains(tt.expr, "=false") && !strings.Contains(tt.expr, "!=")) ||
					(strings.Contains(tt.expr, "!=true")) {
					require.Contains(t, sql, "NOT EXISTS", "expected NOT EXISTS for attachments=false or attachments!=true")
				}
			}
		})
	}
}
