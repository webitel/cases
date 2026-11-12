package postgres

import (
	"context"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/webitel/webitel-go-kit/pkg/etag"

	_go "github.com/webitel/cases/api/cases"
	"github.com/webitel/cases/auth"
	"github.com/webitel/cases/internal/model/options"
	"github.com/webitel/cases/internal/server/interceptor"
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

type placeholderSession struct {
	auth.Auther
	rbac bool
}

func (s placeholderSession) GetDomainId() int64 { return 1 }

func (s placeholderSession) GetUserId() int64 { return 7 }

func (s placeholderSession) GetRoles() []int64 { return []int64{7, 8} }

func (s placeholderSession) IsRbacCheckRequired(string, auth.AccessMode) bool { return s.rbac }

func (s placeholderSession) CheckObacAccess(string, auth.AccessMode) bool { return true }

type placeholderRPC struct {
	context.Context
	session auth.Auther
}

func (r placeholderRPC) Value(key any) any {
	if key == interceptor.SessionHeader {
		return r.session
	}

	return r.Context.Value(key)
}

func (r placeholderRPC) GetAuthOpts() auth.Auther { return r.session }

func (r placeholderRPC) RequestTime() time.Time { return time.Unix(0, 0) }

func (r placeholderRPC) GetFields() []string {
	return []string{"id", "ver", "subject", "comments"}
}

func (r placeholderRPC) GetDerivedSearchOpts() map[string]*options.Searcher { return nil }

func (r placeholderRPC) GetUnknownFields() []string { return nil }

func (r placeholderRPC) GetIDs() []int64 { return nil }

func (r placeholderRPC) GetParentID() int64 { return 0 }

func (r placeholderRPC) GetChildID() int64 { return 0 }

func (r placeholderRPC) GetMask() []string { return []string{"subject"} }

func (r placeholderRPC) GetEtags() []*etag.Tid { return []*etag.Tid{{}} }

func TestCaseWriteQueryPlaceholders(t *testing.T) {
	create := func(store *CaseStore, rpc placeholderRPC) (*Select, error) {
		input := &_go.Case{
			Subject:  "subject",
			Service:  &_go.Service{Id: 1},
			Priority: &_go.Priority{Id: 1},
			Source:   &_go.SourceTypeLookup{Id: 1},
			Status:   &_go.Lookup{Id: 1},
			Reporter: &_go.Lookup{Id: 1},
		}
		base, _, err := store.buildCreateCaseSqlizer(rpc, input, &ServiceRelatedDefs{StatusID: 1, CloseReasonGroupID: 1})

		return base, err
	}
	update := func(store *CaseStore, rpc placeholderRPC) (*Select, error) {
		base, _, err := store.buildUpdateCaseSqlizer(rpc, &_go.Case{Subject: "subject"})

		return base, err
	}

	tests := []struct {
		name  string
		rbac  bool
		build func(store *CaseStore, rpc placeholderRPC) (*Select, error)
	}{
		{name: "create without rbac", build: create},
		{name: "create with rbac", rbac: true, build: create},
		{name: "update without rbac", build: update},
		{name: "update with rbac", rbac: true, build: update},
	}

	placeholder := regexp.MustCompile(`\$(\d+)`)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			session := placeholderSession{rbac: tt.rbac}
			rpc := placeholderRPC{
				Context: context.Background(),
				session: session,
			}

			base, err := tt.build(&CaseStore{mainTable: "cases.case"}, rpc)
			require.NoError(t, err)

			query, args, err := base.ToSql()
			require.NoError(t, err)
			require.NotContains(t, query, "?")

			seen := make(map[int]struct{}, len(args))

			for _, match := range placeholder.FindAllStringSubmatch(query, -1) {
				number, err := strconv.Atoi(match[1])
				require.NoError(t, err)

				seen[number] = struct{}{}
			}

			require.Len(t, seen, len(args))

			for number := 1; number <= len(args); number++ {
				require.Contains(t, seen, number)
			}
		})
	}
}

func TestRenumberPlaceholders(t *testing.T) {
	tests := []struct {
		name   string
		query  string
		offset int
		want   string
	}{
		{name: "no placeholders", query: "SELECT $1", offset: 1, want: "SELECT $1"},
		{name: "continues after offset", query: "a = $1 AND b = ? AND c = ?", offset: 1, want: "a = $1 AND b = $2 AND c = $3"},
		{name: "keeps escaped mark", query: "a ?? 'k' AND b = ?", offset: 4, want: "a ? 'k' AND b = $5"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, renumberPlaceholders(tt.query, tt.offset))
		})
	}
}
