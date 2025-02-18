package postgres

import (
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/webitel/cases/api/cases"
	"github.com/webitel/cases/auth"
	dberr "github.com/webitel/cases/internal/errors"
	"github.com/webitel/cases/internal/store"
	"github.com/webitel/cases/internal/store/postgres/scanner"
	"github.com/webitel/cases/model"
)

type CaseCommunicationStore struct {
	storage   store.Store
	mainTable string
}

func (c *CaseCommunicationStore) Link(options *model.CreateOptions, communications []*cases.InputCaseCommunication) ([]*cases.CaseCommunication, error) {
	if len(communications) == 0 {
		return nil, dberr.NewDBError("postgres.case_communication.link.check_args.communications", "empty communications")
	}
	base, plan, dbErr := c.buildCreateCaseCommunicationSqlizer(options, communications)
	if dbErr != nil {
		return nil, dbErr
	}
	db, dbErr := c.storage.Database()
	if dbErr != nil {
		return nil, dbErr
	}
	sql, args, err := base.ToSql()
	if err != nil {
		return nil, dberr.NewDBError("postgres.case_communication.link.convert_to_sql.err", err.Error())
	}
	rows, err := db.Query(options, store.CompactSQL(sql), args...)
	if err != nil {
		return nil, dberr.NewDBError("postgres.case_communication.exec.error", err.Error())
	}
	res, dbErr := c.scanCommunications(rows, plan)
	if dbErr != nil {
		return nil, dbErr
	}
	return res, nil
}

func (c *CaseCommunicationStore) Unlink(options *model.DeleteOptions) (int64, error) {
	base, dbErr := c.buildDeleteCaseCommunicationSqlizer(options)
	if dbErr != nil {
		return 0, dbErr
	}
	db, dbErr := c.storage.Database()
	if dbErr != nil {
		return 0, dbErr
	}
	sql, args, err := base.ToSql()
	if err != nil {
		return 0, dberr.NewDBError("postgres.case_communication.link.convert_to_sql.err", err.Error())
	}
	res, err := db.Exec(options, sql, args...)
	if err != nil {
		return 0, dberr.NewDBError("postgres.case_communication.exec.error", err.Error())
	}
	return res.RowsAffected(), nil
}

func (c *CaseCommunicationStore) List(opts *model.SearchOptions) (*cases.ListCommunicationsResponse, error) {
	base, plan, dbErr := c.buildListCaseCommunicationSqlizer(opts)
	if dbErr != nil {
		return nil, dbErr
	}
	db, dbErr := c.storage.Database()
	if dbErr != nil {
		return nil, dbErr
	}
	sql, args, err := base.ToSql()
	if err != nil {
		return nil, dberr.NewDBError("postgres.case_communication.link.convert_to_sql.err", err.Error())
	}
	rows, err := db.Query(opts, sql, args...)
	if err != nil {
		return nil, dberr.NewDBError("postgres.case_communication.exec.error", err.Error())
	}
	items, dbErr := c.scanCommunications(rows, plan)
	if dbErr != nil {
		return nil, dbErr
	}
	var res cases.ListCommunicationsResponse
	res.Data, res.Next = store.ResolvePaging(opts.GetSize(), items)
	res.Page = int32(opts.GetPage())
	return &res, nil
}

func (c *CaseCommunicationStore) buildListCaseCommunicationSqlizer(options *model.SearchOptions) (query squirrel.Sqlizer, plan []func(caseCommunication *cases.CaseCommunication) any, dbError *dberr.DBError) {
	if options == nil {
		return nil, nil, dberr.NewDBError("postgres.case_communication.build_list_case_communication_sqlizer.check_args.options", "search options required")
	}
	if options.ParentId <= 0 {
		return nil, nil, dberr.NewDBError("postgres.case_communication.build_list_case_communication_sqlizer.check_args.case_id", "case id required")
	}
	alias := "s"
	base := squirrel.Select().From(fmt.Sprintf("%s %s", c.mainTable, alias)).PlaceholderFormat(squirrel.Dollar)
	base = store.ApplyPaging(options.GetPage(), options.GetSize(), base)
	return c.buildSelectColumnsAndPlan(base, alias, options.Fields)
}

func (c *CaseCommunicationStore) scanCommunications(rows pgx.Rows, plan []func(*cases.CaseCommunication) any) ([]*cases.CaseCommunication, *dberr.DBError) {
	var res []*cases.CaseCommunication
	for rows.Next() {
		node := &cases.CaseCommunication{}
		var scanValues []any
		for _, f := range plan {
			scanValues = append(scanValues, f(node))
		}
		err := rows.Scan(scanValues...)
		if err != nil {
			return nil, dberr.NewDBError("postgres.case_communication.scan.error", err.Error())
		}
		res = append(res, node)
	}
	return res, nil
}

func (c *CaseCommunicationStore) buildCreateCaseCommunicationSqlizer(options *model.CreateOptions, communications []*cases.InputCaseCommunication) (query squirrel.Sqlizer, plan []func(caseCommunication *cases.CaseCommunication) any, dbError *dberr.DBError) {
	if options == nil {
		return nil, nil, dberr.NewDBError("postgres.case_communication.build_create_case_communication_sqlizer.check_args.options", "create options required")
	}
	if options.ParentID <= 0 {
		return nil, nil, dberr.NewDBError("postgres.case_communication.build_create_case_communication_sqlizer.check_args.case_id", "case id required")
	}
	insert := squirrel.Insert(c.mainTable).Columns("created_by", "created_at", "dc", "communication_type", "communication_id", "case_id").Suffix("ON CONFLICT DO NOTHING RETURNING *")

	var (
		caseId              = options.ParentID
		dc                  *int64
		userId              *int64
		roles               []int64
		callsRbac, caseRbac bool
	)
	if session := options.GetAuthOpts(); session != nil {
		d := session.GetDomainId()
		dc = &d
		u := session.GetUserId()
		userId = &u
		roles = session.GetRoles()
		callsRbac = session.IsRbacCheckRequired("calls", auth.Read)
		caseRbac = session.IsRbacCheckRequired("cases", auth.Edit)
	}
	var caseSubquery squirrel.Sqlizer
	if caseRbac {
		caseSubquery = squirrel.Expr(`(SELECT object FROM cases.case_acl acl WHERE acl.dc = ? AND acl.object = ? AND acl.subject = any(?::int[]) AND acl.access & ? = ?)`,
			dc, caseId, roles, auth.Edit, auth.Edit)
	} else {
		caseSubquery = squirrel.Expr(`?`, caseId)
	}

	for _, communication := range communications {
		switch communication.CommunicationType {
		case cases.CaseCommunicationsTypes_COMMUNICATION_CALL:
			var callsSubquery squirrel.Sqlizer
			if callsRbac {
				callsSubquery = squirrel.Expr(`(SELECT c.id::text
															 FROM call_center.cc_calls_history c
															 WHERE id = ?::uuid
															   AND (
																 (c.user_id = ANY (call_center.cc_calls_rbac_users(?::int8, ?::int8) || ?::int[])
																	 OR c.queue_id = ANY (call_center.cc_calls_rbac_queues(?::int8, ?::int8, ?::int[]))
																	 OR (c.user_ids NOTNULL AND c.user_ids::int[] && call_center.rbac_users_from_group('calls', ?::int8, ?::int2, ?::int[]))
																	 OR (c.grantee_id = ANY (?::int[]))
																  )
																 ))`,
					communication.CommunicationId,
					dc, userId, roles,
					dc, userId, roles,
					dc, auth.Read, roles,
					roles)
			} else {
				callsSubquery = squirrel.Expr(`(SELECT id FROM call_center.cc_calls_history WHERE id = ?)`, communication.CommunicationId)
			}
			insert = insert.Values(userId, options.CurrentTime(), dc, int64(communication.CommunicationType), callsSubquery, caseSubquery)
		case cases.CaseCommunicationsTypes_COMMUNICATION_CHAT:
			insert = insert.Values(userId, options.CurrentTime(), dc, int64(communication.CommunicationType), squirrel.Expr(`(SELECT id FROM chat.conversation WHERE id = ?)`, communication.CommunicationId), caseSubquery)
		case cases.CaseCommunicationsTypes_COMMUNICATION_EMAIL:
			insert = insert.Values(userId, options.CurrentTime(), dc, int64(communication.CommunicationType), squirrel.Expr(`(SELECT id FROM call_center.cc_email WHERE id = ?)`, communication.CommunicationId), caseSubquery)
		default:
			return nil, nil, dberr.NewDBError("postgres.case_communication.build.create_case_communication_sqlizer.switch_types.unknown", "unsupported communication type")
		}
	}
	insertAlias := "i"
	insertCte, args, err := store.FormAsCTE(insert, insertAlias)
	if err != nil {
		return nil, nil, dberr.NewDBError("postgres.case_communication.build.create_case_communication_sqlizer.form_cte.error", err.Error())
	}
	base := squirrel.Select().From(insertAlias).Prefix(insertCte, args...).PlaceholderFormat(squirrel.Dollar)
	return c.buildSelectColumnsAndPlan(base, insertAlias, options.Fields)
}

func (c *CaseCommunicationStore) buildSelectColumnsAndPlan(base squirrel.SelectBuilder, left string, fields []string) (query squirrel.SelectBuilder, plan []func(caseCommunication *cases.CaseCommunication) any, dbError *dberr.DBError) {
	if len(fields) == 0 {
		fields = CaseCommunicationFields
	}
	for _, field := range fields {
		switch field {
		case "id":
			base = base.Column(store.Ident(left, "id"))
			plan = append(plan, func(comm *cases.CaseCommunication) any {
				return &comm.Id
			})
		case "ver":
			base = base.Column(store.Ident(left, "ver"))
			plan = append(plan, func(comm *cases.CaseCommunication) any {
				return &comm.Ver
			})
		case "communication_type":
			base = base.Column(store.Ident(left, "communication_type"))
			plan = append(plan, func(comm *cases.CaseCommunication) any {
				return &comm.CommunicationType
			})
		case "communication_id":
			base = base.Column(store.Ident(left, "communication_id"))
			plan = append(plan, func(comm *cases.CaseCommunication) any {
				return scanner.ScanText(&comm.CommunicationId)
			})
		default:
			return base, nil, dberr.NewDBError("postgres.case_communication.build_select_columns_and_plan.cycle_fields.unknown", fmt.Sprintf("%s field is unknown", field))
		}
	}
	return base, plan, nil
}

func (c *CaseCommunicationStore) buildDeleteCaseCommunicationSqlizer(options *model.DeleteOptions) (query squirrel.Sqlizer, dbError *dberr.DBError) {
	if options == nil {
		return nil, dberr.NewDBError("postgres.case_communication.build_delete_case_communication_sqlizer.check_args.options", "delete options required")
	}
	if len(options.IDs) == 0 {
		return nil, dberr.NewDBError("postgres.case_communication.build_delete_case_communication_sqlizer.check_args.ids", "ids required to delete")
	}
	del := squirrel.Delete(c.mainTable).Where("id = ANY(?)", options.IDs)
	return del, nil
}

var s store.CaseCommunicationStore = &CaseCommunicationStore{}

var CaseCommunicationFields = []string{"id", "ver", "communication_type", "communication_id"}

func NewCaseCommunicationStore(store store.Store) (store.CaseCommunicationStore, error) {
	if store == nil {
		return nil, dberr.NewDBError("postgres.new_case_communication.check.bad_arguments",
			"error creating case communication store, main store is nil")
	}
	return &CaseCommunicationStore{storage: store, mainTable: "cases.case_communication"}, nil
}
