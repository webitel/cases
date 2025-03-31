package postgres

import (
	"fmt"
	"github.com/webitel/cases/internal/store/postgres/transaction"
	"github.com/webitel/cases/internal/store/util"
	"github.com/webitel/cases/model/options"

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
	storage   *Store
	mainTable string
}

func (c *CaseCommunicationStore) Link(
	options options.CreateOptions,
	communications []*cases.InputCaseCommunication,
) ([]*cases.CaseCommunication, error) {
	if len(communications) == 0 {
		return nil, dberr.NewDBError(
			"postgres.case_communication.link.check_args.communications",
			"empty communications",
		)
	}

	d, dbErr := c.storage.Database()
	if dbErr != nil {
		return nil, dberr.NewDBInternalError("postgres.case_communication.link.database_connection_error", dbErr)
	}

	tx, err := d.Begin(options)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.case_communication.link.transaction_error", err)
	}
	defer func() {
		_ = tx.Rollback(options)
	}()

	txManager := transaction.NewTxManager(tx)

	base, plan, dbErr := c.buildCreateCaseCommunicationSqlizer(txManager, options, communications)
	if dbErr != nil {
		return nil, dbErr
	}

	sql, args, err := base.ToSql()
	if err != nil {
		return nil, dberr.NewDBError("postgres.case_communication.link.convert_to_sql.err", err.Error())
	}

	rows, err := txManager.Query(options, util.CompactSQL(sql), args...)
	if err != nil {
		return nil, dberr.NewDBError("postgres.case_communication.link.exec.error", err.Error())
	}
	defer rows.Close()

	res, dbErr := c.scanCommunications(rows, plan)
	if dbErr != nil {
		return nil, dbErr
	}

	if err := txManager.Commit(options); err != nil {
		return nil, dberr.NewDBInternalError("postgres.case_communication.link.commit_error", err)
	}

	return res, nil
}

func (c *CaseCommunicationStore) Unlink(options options.DeleteOptions) (int64, error) {
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

func (c *CaseCommunicationStore) List(opts options.SearchOptions) (*cases.ListCommunicationsResponse, error) {
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
	res.Data, res.Next = util.ResolvePaging(opts.GetSize(), items)
	res.Page = int32(opts.GetPage())
	return &res, nil
}

func (c *CaseCommunicationStore) buildListCaseCommunicationSqlizer(
	options options.SearchOptions,
) (
	query squirrel.Sqlizer,
	plan []func(caseCommunication *cases.CaseCommunication) any,
	dbError *dberr.DBError,
) {
	if options == nil {
		return nil, nil, dberr.NewDBError(
			"postgres.case_communication.build_list_case_communication_sqlizer.check_args.options",
			"search options required",
		)
	}
	parentId, ok := options.GetFilter("case_id").(int64)
	if !ok || parentId == 0 {
		return nil, nil, dberr.NewDBError(
			"postgres.case_communication.build_list_case_communication_sqlizer.check_args.case_id",
			"case id required",
		)
	}
	alias := "s"
	base := squirrel.Select().
		From(fmt.Sprintf("%s %s", c.mainTable, alias)).
		Where(fmt.Sprintf("%s = ?", util.Ident(alias, "case_id")), parentId).
		PlaceholderFormat(squirrel.Dollar)
	base = util.ApplyPaging(options.GetPage(), options.GetSize(), base)
	return c.buildSelectColumnsAndPlan(base, alias, options.GetFields())
}

func (c *CaseCommunicationStore) scanCommunications(
	rows pgx.Rows,
	plan []func(*cases.CaseCommunication) any,
) (
	[]*cases.CaseCommunication,
	*dberr.DBError,
) {
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

func (c *CaseCommunicationStore) buildCreateCaseCommunicationSqlizer(
	tx transaction.Transaction,
	options options.CreateOptions,
	input []*cases.InputCaseCommunication,
) (
	query squirrel.Sqlizer,
	plan []func(caseCommunication *cases.CaseCommunication) any,
	dbError *dberr.DBError,
) {
	if options == nil {
		return nil, nil, dberr.NewDBError(
			"postgres.case_communication.build_create_case_communication_sqlizer.check_args.options",
			"create options required",
		)
	}
	if options.GetParentID() <= 0 {
		return nil, nil, dberr.NewDBError(
			"postgres.case_communication.build_create_case_communication_sqlizer.check_args.case_id",
			"case id required",
		)
	}

	insert := squirrel.Insert(c.mainTable).
		Columns("created_by", "created_at", "dc", "communication_type", "communication_id", "case_id").
		Suffix("RETURNING *")

	var (
		caseId              = options.GetParentID()
		dc, userId          *int64
		roles               []int64
		callsRbac, caseRbac bool
	)

	if session := options.GetAuthOpts(); session != nil {
		d := session.GetDomainId()
		dc = &d
		u := session.GetUserId()
		userId = &u
		roles = session.GetRoles()
		callsRbac = session.IsRbacCheckRequired(model.ScopeCalls, auth.Read)
		caseRbac = session.IsRbacCheckRequired(model.ScopeCases, auth.Edit)
	}

	var caseSubquery squirrel.Sqlizer
	if caseRbac {
		caseSubquery = squirrel.Expr(`(
			SELECT object FROM cases.case_acl acl
			WHERE acl.dc = ? AND acl.object = ? AND acl.subject = ANY(?::int[]) AND acl.access & ? = ?)`,
			dc, caseId, roles, auth.Edit, auth.Edit,
		)
	} else {
		caseSubquery = squirrel.Expr(`?`, caseId)
	}

	for _, communication := range input {
		var channel string
		err := tx.QueryRow(
			options,
			`SELECT channel FROM call_center.cc_communication WHERE id = $1 AND dc = $2`,
			communication.CommunicationId,
			options.GetAuthOpts().GetDomainId(),
		).Scan(&channel)
		if err != nil {
			return nil, nil, dberr.NewDBError("postgres.case_communication.resolve_channel", err.Error())
		}

		var (
			commType int64
			subquery squirrel.Sqlizer
		)

		switch channel {
		case "Phone":
			commType = int64(cases.CaseTimelineEventType_call)
			if callsRbac {
				subquery = squirrel.Expr(`(
					SELECT c.id::text FROM call_center.cc_calls_history c
					WHERE c.id = ?::uuid AND (
						c.user_id = ANY(call_center.cc_calls_rbac_users(?::int8, ?::int8) || ?::int[])
						OR c.queue_id = ANY(call_center.cc_calls_rbac_queues(?::int8, ?::int8, ?::int[]))
						OR (c.user_ids NOTNULL AND c.user_ids::int[] && call_center.rbac_users_from_group('calls', ?::int8, ?::int2, ?::int[]))
						OR c.grantee_id = ANY(?::int[])
					)
				)`,
					communication.CommunicationId,
					dc, userId, roles,
					dc, userId, roles,
					dc, auth.Read, roles,
					roles,
				)
			} else {
				subquery = squirrel.Expr(
					`(SELECT id FROM call_center.cc_calls_history WHERE id = ?)`,
					communication.CommunicationId,
				)
			}
		case "Messaging":
			commType = int64(cases.CaseTimelineEventType_chat)
			subquery = squirrel.Expr(
				`(SELECT id FROM chat.conversation WHERE id = ?)`,
				communication.CommunicationId,
			)
		case "Email":
			commType = int64(cases.CaseTimelineEventType_email)
			subquery = squirrel.Expr(
				`(SELECT id FROM call_center.cc_email WHERE id = ?)`,
				communication.CommunicationId,
			)
		default:
			return nil, nil, dberr.NewDBError(
				"postgres.case_communication.unknown_channel",
				fmt.Sprintf("unknown communication channel: %s", channel),
			)
		}

		insert = insert.Values(userId, options.RequestTime(), dc, commType, subquery, caseSubquery)
	}

	insertAlias := "i"
	insertCte, args, err := util.FormAsCTE(insert, insertAlias)
	if err != nil {
		return nil, nil, dberr.NewDBError(
			"postgres.case_communication.build_create_case_communication_sqlizer.form_cte.error",
			err.Error(),
		)
	}

	base := squirrel.Select().From(insertAlias).
		Prefix(insertCte, args...).
		PlaceholderFormat(squirrel.Dollar)

	return c.buildSelectColumnsAndPlan(base, insertAlias, options.GetFields())
}

func (c *CaseCommunicationStore) buildSelectColumnsAndPlan(
	base squirrel.SelectBuilder,
	left string,
	fields []string,
) (query squirrel.SelectBuilder, plan []func(comm *cases.CaseCommunication) any, dbError *dberr.DBError) {
	if len(fields) == 0 {
		fields = CaseCommunicationFields
	}

	var (
		communicationAlias string
		joinCommunication  = func() {
			if communicationAlias != "" {
				return
			}
			communicationAlias = "comm"
			base = base.LeftJoin(
				fmt.Sprintf(
					"call_center.cc_communication %s ON %s.id = %s.communication_id",
					communicationAlias,
					communicationAlias,
					left,
				),
			)
		}
	)

	for _, field := range fields {
		switch field {
		case "id":
			base = base.Column(util.Ident(left, "id"))
			plan = append(plan, func(comm *cases.CaseCommunication) any {
				return &comm.Id
			})
		case "ver":
			base = base.Column(util.Ident(left, "ver"))
			plan = append(plan, func(comm *cases.CaseCommunication) any {
				return &comm.Ver
			})
		case "communication_type":
			joinCommunication()
			base = base.Column(
				fmt.Sprintf("ROW(%s.communication_type, %s.channel)::text AS communication_type",
					communicationAlias,
					communicationAlias,
				),
			)
			plan = append(plan, func(comm *cases.CaseCommunication) any {
				return scanner.ScanRowLookup(&comm.CommunicationType)
			})
		case "communication_id":
			base = base.Column(util.Ident(left, "communication_id"))
			plan = append(plan, func(comm *cases.CaseCommunication) any {
				return scanner.ScanText(&comm.CommunicationId)
			})
		default:
			return base, nil, dberr.NewDBError(
				"postgres.case_communication.build_select_columns_and_plan.cycle_fields.unknown",
				fmt.Sprintf("%s field is unknown", field),
			)
		}
	}

	return base, plan, nil
}

func (c *CaseCommunicationStore) buildDeleteCaseCommunicationSqlizer(
	options options.DeleteOptions,
) (
	query squirrel.Sqlizer,
	dbError *dberr.DBError,
) {
	if options == nil {
		return nil, dberr.NewDBError(
			"postgres.case_communication.build_delete_case_communication_sqlizer.check_args.options",
			"delete options required",
		)
	}
	if len(options.GetIDs()) == 0 {
		return nil, dberr.NewDBError(
			"postgres.case_communication.build_delete_case_communication_sqlizer.check_args.ids",
			"ids required to delete",
		)
	}
	del := squirrel.Delete(c.mainTable).Where("id = ANY(?)", options.GetIDs())
	return del, nil
}

var s store.CaseCommunicationStore = &CaseCommunicationStore{}

var CaseCommunicationFields = []string{"id", "ver", "communication_type", "communication_id"}

func NewCaseCommunicationStore(store *Store) (store.CaseCommunicationStore, error) {
	if store == nil {
		return nil, dberr.NewDBError("postgres.new_case_communication.check.bad_arguments",
			"error creating case communication store, main store is nil")
	}
	return &CaseCommunicationStore{storage: store, mainTable: "cases.case_communication"}, nil
}
