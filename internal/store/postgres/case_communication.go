package postgres

import (
	"fmt"

	"github.com/webitel/cases/auth"
	"github.com/webitel/cases/internal/store/postgres/transaction"
	storeUtil "github.com/webitel/cases/internal/store/util"
	"github.com/webitel/cases/util"

	"github.com/webitel/cases/model/options"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/webitel/cases/api/cases"
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
	options options.Creator,
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

	// Defer rollback, but ignore error if the transaction is already committed/rolled back
	defer func() {
		_ = tx.Rollback(options) // Will be a no-op if already committed
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

	rows, err := txManager.Query(options, storeUtil.CompactSQL(sql), args...)
	if err != nil {
		return nil, dberr.NewDBError("postgres.case_communication.link.exec.error", err.Error())
	}
	defer rows.Close()

	res, dbErr := c.scanCommunications(rows, plan)
	if dbErr != nil {
		return nil, dbErr
	}

	// If Commit succeeds, rollback is now a no-op
	if err := txManager.Commit(options); err != nil {
		return nil, dberr.NewDBInternalError("postgres.case_communication.link.commit_error", err)
	}

	return res, nil
}

func (c *CaseCommunicationStore) Unlink(options options.Deleter) (int64, error) {
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

func (c *CaseCommunicationStore) List(opts options.Searcher) (*cases.ListCommunicationsResponse, error) {
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
	res.Data, res.Next = storeUtil.ResolvePaging(opts.GetSize(), items)
	res.Page = int32(opts.GetPage())
	return &res, nil
}

func (c *CaseCommunicationStore) buildListCaseCommunicationSqlizer(
	options options.Searcher,
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
	caseIDFilters := options.GetFilter("case_id")
	if len(caseIDFilters) == 0 {
		return nil, nil, dberr.NewDBError(
			"postgres.case_communication.build_list_case_communication_sqlizer.check_args.case_id",
			"case id required",
		)
	}
	alias := "s"
	base := squirrel.Select().
		From(fmt.Sprintf("%s %s", c.mainTable, alias)).
		Where(fmt.Sprintf("%s = ?", storeUtil.Ident(alias, "dc")), options.GetAuthOpts().GetDomainId()).
		PlaceholderFormat(squirrel.Dollar)
	// Apply all case_id filters (with all supported operators)
	base = util.ApplyFiltersToQuery(base, storeUtil.Ident(alias, "case_id"), caseIDFilters)
	base = storeUtil.ApplyPaging(options.GetPage(), options.GetSize(), base)

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
	options options.Creator,
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
		var (
			channel    string
			commTypeId int64
			subquery   squirrel.Sqlizer
		)
		if commType := communication.CommunicationType; commType == nil || commType.GetId() == 0 {
			// can't determine communication type
			// skip?
			continue
		} else {
			commTypeId = commType.GetId()
		}
		err := tx.QueryRow(
			options,
			`SELECT channel FROM call_center.cc_communication WHERE id = $1 AND domain_id = $2`,
			commTypeId,
			options.GetAuthOpts().GetDomainId(),
		).Scan(&channel)
		if err != nil {
			return nil, nil, dberr.NewDBError("postgres.case_communication.resolve_channel", err.Error())
		}

		var ()

		switch channel {
		case store.CommunicationCall:
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
		case store.CommunicationChat:
			subquery = squirrel.Expr(
				`(SELECT id FROM chat.conversation WHERE id = ?)`,
				communication.CommunicationId,
			)
		case store.CommunicationEmail:
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

		insert = insert.Values(userId, options.RequestTime(), dc, commTypeId, subquery, caseSubquery)
	}

	insertAlias := "i"
	insertCte, args, err := storeUtil.FormAsCTE(insert, insertAlias)
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
					"call_center.cc_communication %s ON %s.id = %s.communication_type",
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
			base = base.Column(storeUtil.Ident(left, "id"))
			plan = append(plan, func(comm *cases.CaseCommunication) any {
				return &comm.Id
			})
		case "ver":
			base = base.Column(storeUtil.Ident(left, "ver"))
			plan = append(plan, func(comm *cases.CaseCommunication) any {
				return &comm.Ver
			})
		case "communication_type":
			joinCommunication()
			base = base.Column(
				fmt.Sprintf("ROW(%s.id, %s.channel)::text AS communication_type",
					communicationAlias,
					communicationAlias,
				),
			)
			plan = append(plan, func(comm *cases.CaseCommunication) any {
				return scanner.ScanRowLookup(&comm.CommunicationType)
			})
		case "communication_id":
			base = base.Column(storeUtil.Ident(left, "communication_id"))
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
	options options.Deleter,
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
