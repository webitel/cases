package postgres

import (
	"testing"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/stretchr/testify/require"

	"github.com/webitel/cases/auth"
	"github.com/webitel/cases/internal/model"
	"github.com/webitel/cases/internal/model/options"
)

type commentAccessSession struct {
	auth.Auther

	obac map[string]bool
	rbac map[string]bool
}

func (s *commentAccessSession) CheckObacAccess(scope string, _ auth.AccessMode) bool {
	return s.obac[scope]
}

func (s *commentAccessSession) IsRbacCheckRequired(scope string, _ auth.AccessMode) bool {
	return s.rbac[scope]
}

func (s *commentAccessSession) GetDomainId() int64 { return 1 }

func (s *commentAccessSession) GetRoles() []int64 { return []int64{10} }

func TestCaseCommentAccessCondition(t *testing.T) {
	const (
		external = model.ScopeCaseComments
		internal = model.ScopeInternalComments
		aclExpr  = "EXISTS(SELECT acl.object FROM cases.case_comment_acl acl WHERE acl.dc = ? AND acl.object = cc.id AND acl.subject = any( ?::int[]) AND acl.access & ? = ? LIMIT 1)"
	)

	tests := []struct {
		name     string
		obac     map[string]bool
		rbac     map[string]bool
		wantSQL  string
		wantArgs []any
	}{
		{
			name:     "both types",
			obac:     map[string]bool{external: true, internal: true},
			wantSQL:  "((cc.type = ?) OR (cc.type = ?))",
			wantArgs: []any{model.CaseCommentTypeExternal, model.CaseCommentTypeInternal},
		},
		{
			name:     "external only",
			obac:     map[string]bool{external: true},
			wantSQL:  "((cc.type = ?))",
			wantArgs: []any{model.CaseCommentTypeExternal},
		},
		{
			name:     "internal only",
			obac:     map[string]bool{internal: true},
			wantSQL:  "((cc.type = ?))",
			wantArgs: []any{model.CaseCommentTypeInternal},
		},
		{
			name:    "no access",
			obac:    map[string]bool{},
			wantSQL: "FALSE",
		},
		{
			name:    "rbac on internal only",
			obac:    map[string]bool{external: true, internal: true},
			rbac:    map[string]bool{internal: true},
			wantSQL: "((cc.type = ?) OR (cc.type = ? AND " + aclExpr + "))",
		},
		{
			name:    "rbac without obac grant is skipped",
			obac:    map[string]bool{internal: true},
			rbac:    map[string]bool{external: true},
			wantSQL: "((cc.type = ?))",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			session := &commentAccessSession{obac: tc.obac, rbac: tc.rbac}

			cond := caseCommentAccessCondition(session, auth.Read, "cc.type", "cc.id")
			require.NotNil(t, cond)

			sql, args, err := cond.ToSql()
			require.NoError(t, err)
			require.Equal(t, tc.wantSQL, sql)

			if tc.wantArgs != nil {
				require.Equal(t, tc.wantArgs, args)
			}
		})
	}
}

func TestCaseCommentAccessConditionWithoutSession(t *testing.T) {
	require.Nil(t, caseCommentAccessCondition(nil, auth.Read, "cc.type", "cc.id"))
}

type commentFieldsOpts struct {
	options.Creator

	fields  []string
	session auth.Auther
}

func (o *commentFieldsOpts) GetFields() []string { return o.fields }

func (o *commentFieldsOpts) GetAuthOpts() auth.Auther { return o.session }

func (o *commentFieldsOpts) RequestTime() time.Time { return time.Time{} }

type commentUpdateOpts struct {
	options.Updator

	*commentFieldsOpts
}

func (o *commentUpdateOpts) GetFields() []string { return o.fields }

func (o *commentUpdateOpts) GetAuthOpts() auth.Auther { return o.session }

func (o *commentUpdateOpts) RequestTime() time.Time { return time.Time{} }

type commentDeleteOpts struct {
	options.Deleter

	*commentFieldsOpts
}

func (o *commentDeleteOpts) GetFields() []string { return o.fields }

func (o *commentDeleteOpts) GetAuthOpts() auth.Auther { return o.session }

func (o *commentDeleteOpts) GetIDs() []int64 { return []int64{1} }

func (o *commentDeleteOpts) RequestTime() time.Time { return time.Time{} }

func (s *commentAccessSession) GetUserId() int64 { return 1 }

func TestCaseCommentWriteQueriesAlwaysSelectType(t *testing.T) {
	session := &commentAccessSession{obac: map[string]bool{model.ScopeCaseComments: true, model.ScopeInternalComments: true}}
	base := &commentFieldsOpts{fields: []string{"text"}, session: session}
	store := &CaseCommentStore{}

	publish, err := store.buildPublishCaseCommentQuery(base, &model.CaseComment{Text: "x"})
	require.NoError(t, err)

	update, err := store.buildUpdateCaseCommentQuery(&commentUpdateOpts{commentFieldsOpts: base}, &model.CaseComment{Text: "x"})
	require.NoError(t, err)

	remove, err := store.buildDeleteCaseCommentQuery(&commentDeleteOpts{commentFieldsOpts: base})
	require.NoError(t, err)

	for name, builder := range map[string]sq.SelectBuilder{"publish": publish, "update": update, "delete": remove} {
		t.Run(name, func(t *testing.T) {
			sql, _, err := builder.ToSql()
			require.NoError(t, err)
			require.Contains(t, sql, "cc.type")
		})
	}
}
