package postgres

import (
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"

	"github.com/webitel/cases/internal/errors"
	"github.com/webitel/cases/internal/model"
	"github.com/webitel/cases/internal/model/options"
	"github.com/webitel/cases/internal/store"
	"github.com/webitel/cases/internal/store/postgres/transaction"
	storeutil "github.com/webitel/cases/internal/store/util"
)

type CaseCommunicationStore struct {
	storage   *Store
	mainTable string
}

func (c *CaseCommunicationStore) Link(
	options options.Creator,
	communications []*model.CaseCommunication,
) ([]*model.CaseCommunication, error) {
	if len(communications) == 0 {
		return nil, errors.InvalidArgument("empty communications")
	}
	db, err := c.storage.Database()
	if err != nil {
		return nil, errors.Internal("postgres.case_communication.link.database_connection_error", errors.WithCause(err))
	}
	tx, err := db.Begin(options)
	if err != nil {
		return nil, errors.Internal("postgres.case_communication.link.transaction_error", errors.WithCause(err))
	}

	defer func() { _ = tx.Rollback(options) }()

	// Use internal model.CaseCommunication directly
	base, err := c.buildCreateCaseCommunicationSqlizer(tx, options, communications)
	if err != nil {
		return nil, err
	}
	query, args, err := base.ToSql()
	if err != nil {
		return nil, err
	}
	var result []*model.CaseCommunication
	err = pgxscan.Select(options, tx, &result, query, args...)
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(options); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *CaseCommunicationStore) Unlink(options options.Deleter) (int64, error) {
	base, dbErr := c.buildDeleteCaseCommunicationSqlizer(options)
	if dbErr != nil {
		return 0, dbErr
	}
	db, err := c.storage.Database()
	if err != nil {
		return 0, err
	}
	sql, args, err := base.ToSql()
	if err != nil {
		return 0, errors.Internal("postgres.case_communication.link.convert_to_sql.err", errors.WithCause(err))
	}
	res, err := db.Exec(options, sql, args...)
	if err != nil {
		return 0, errors.Internal("postgres.case_communication.exec.error", errors.WithCause(err))
	}
	return res.RowsAffected(), nil
}

func (c *CaseCommunicationStore) List(options options.Searcher) ([]*model.CaseCommunication, error) {
	if options == nil {
		return nil, errors.InvalidArgument("search options required")
	}
	base, err := c.buildListCaseCommunicationSqlizer(options)
	if err != nil {
		return nil, err
	}
	db, err := c.storage.Database()
	if err != nil {
		return nil, err
	}
	sql, args, err := base.ToSql()
	if err != nil {
		return nil, err
	}
	var items []*model.CaseCommunication
	err = pgxscan.Select(options, db, &items, sql, args...)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (c *CaseCommunicationStore) buildListCaseCommunicationSqlizer(
	options options.Searcher,
) (squirrel.SelectBuilder, error) {
	if options == nil {
		return squirrel.SelectBuilder{}, errors.InvalidArgument("search options required")
	}
	caseIDFilters := options.GetFilter("case_id")
	if len(caseIDFilters) == 0 {
		return squirrel.SelectBuilder{}, errors.InvalidArgument("case id required")
	}
	alias := "s"
	base := squirrel.Select().
		From(fmt.Sprintf("%s %s", c.mainTable, alias)).
		Where(fmt.Sprintf("%s = ?", storeutil.Ident(alias, "dc")), options.GetAuthOpts().GetDomainId()).
		PlaceholderFormat(squirrel.Dollar)
	// Apply all case_id filters (with all supported operators)
	base = ApplyFiltersToQuery(base, storeutil.Ident(alias, "case_id"), caseIDFilters)
	base = storeutil.ApplyPaging(options.GetPage(), options.GetSize(), base)

	// Dynamic columns selection
	fields := options.GetFields()
	base, err := c.buildSelectColumns(base, alias, fields)
	if err != nil {
		return squirrel.SelectBuilder{}, err
	}
	return base, nil
}

func (c *CaseCommunicationStore) buildCreateCaseCommunicationSqlizer(
	tx transaction.Transaction,
	options options.Creator,
	input []*model.CaseCommunication,
) (
	query squirrel.Sqlizer,
	err error,
) {
	if options == nil {
		return nil, errors.InvalidArgument("create options required")
	}
	if options.GetParentID() <= 0 {
		return nil, errors.InvalidArgument("case id required")
	}

	insert := squirrel.Insert(c.mainTable).
		Columns("created_by", "created_at", "dc", "communication_type", "communication_id", "case_id").
		Suffix("RETURNING *")

	var (
		caseID     = options.GetParentID()
		dc, userID *int64
		// roles      []int64
		// callsRbac, caseRbac bool
	)

	if session := options.GetAuthOpts(); session != nil {
		d := session.GetDomainId()
		dc = &d
		u := session.GetUserId()
		userID = &u
		// roles = session.GetRoles()
		// callsRbac = session.IsRbacCheckRequired(model.ScopeCalls, auth.Read)
		// caseRbac = session.IsRbacCheckRequired(model.ScopeCases, auth.Edit)
	}

	//	if caseRbac {
	//		caseSubquery := squirrel.Expr(`(
	//			SELECT object FROM cases.case_acl acl
	//			WHERE acl.dc = ? AND acl.object = ? AND acl.subject = ANY(?::int[]) AND acl.access & ? = ?)`,
	//		dc, caseId, roles, auth.Edit, auth.Edit,
	//		)
	//	}
	//
	// else {
	caseSubquery := squirrel.Expr(`?`, caseID)
	// }

	for _, communication := range input {
		var (
			channel    string
			commTypeID int64
			subquery   squirrel.Sqlizer
		)
		if communication.CommunicationType == nil || communication.CommunicationType.Id == nil || *communication.CommunicationType.Id == 0 {
			// can't determine communication type
			// skip?
			continue
		}
		commTypeID = int64(*communication.CommunicationType.Id)
		err := tx.QueryRow(
			options,
			`SELECT channel FROM call_center.cc_communication WHERE id = $1 AND domain_id = $2`,
			commTypeID,
			options.GetAuthOpts().GetDomainId(),
		).Scan(&channel)
		if err != nil {
			return nil, errors.Internal("postgres.case_communication.resolve_channel", errors.WithCause(err))
		}

		switch channel {
		case store.CommunicationCall:
			//Fixme rbac with active calls
			//if callsRbac {
			//	subquery = squirrel.Expr(`(
			//		SELECT c.id::text FROM call_center.cc_calls_history c
			//		WHERE c.id = ?::uuid AND (
			//			c.user_id = ANY(call_center.cc_calls_rbac_users(?::int8, ?::int8) || ?::int[])
			//			OR c.queue_id = ANY(call_center.cc_calls_rbac_queues(?::int8, ?::int8, ?::int[]))
			//			OR (c.user_ids NOTNULL AND c.user_ids::int[] && call_center.rbac_users_from_group('calls', ?::int8, ?::int2, ?::int[]))
			//			OR c.grantee_id = ANY(?::int[])
			//		)
			//	)`,
			//		communication.CommunicationId,
			//		dc, userId, roles,
			//		dc, userId, roles,
			//		dc, auth.Read, roles,
			//		roles,
			//	)
			//}
			//else {
			subquery = squirrel.Expr(
				`(
		SELECT COALESCE(
			(SELECT id FROM call_center.cc_calls_history WHERE id = ?),
			(SELECT id FROM call_center.cc_calls WHERE id = ?)
		)
	)`,
				communication.CommunicationId,
				communication.CommunicationId,
			)
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
			return nil, errors.InvalidArgument(fmt.Sprintf("unknown communication channel: %s", channel))
		}

		insert = insert.Values(userID, options.RequestTime(), dc, commTypeID, subquery, caseSubquery)
	}

	insertAlias := "i"
	insertCte, args, err := storeutil.FormAsCTE(insert, insertAlias)
	if err != nil {
		return nil, errors.Internal(
			"postgres.case_communication.build_create_case_communication_sqlizer.form_cte.error",
			errors.WithCause(err),
		)
	}

	base := squirrel.Select().From(insertAlias).
		Prefix(insertCte, args...).
		PlaceholderFormat(squirrel.Dollar)

	base, err = c.buildSelectColumns(base, insertAlias, options.GetFields())
	if err != nil {
		return nil, err
	}
	return base, nil
}

func (c *CaseCommunicationStore) buildSelectColumns(
	base squirrel.SelectBuilder,
	left string,
	fields []string,
) (squirrel.SelectBuilder, error) {
	if len(fields) == 0 {
		fields = CaseCommunicationFields
	}

	var communicationAlias string
	joinCommunication := func() {
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

	for _, field := range fields {
		switch field {
		case "id":
			base = base.Column(storeutil.Ident(left, "id"))
		case "ver":
			base = base.Column(storeutil.Ident(left, "ver"))
		case "communication_type":
			joinCommunication()
			base = base.Column(
				fmt.Sprintf(`jsonb_build_object('id', %s.id, 'name', %s.channel) AS communication_type`,
					communicationAlias,
					communicationAlias,
				),
			)
		case "communication_id":
			base = base.Column(storeutil.Ident(left, "communication_id"))
		default:
			// ignore unknown fields for now
		}
	}
	return base, nil
}

func (c *CaseCommunicationStore) buildDeleteCaseCommunicationSqlizer(
	options options.Deleter,
) (
	query squirrel.Sqlizer,
	err error,
) {
	if options == nil {
		return nil, errors.InvalidArgument("delete options required")
	}
	if len(options.GetIDs()) == 0 {
		return nil, errors.InvalidArgument("ids required to delete")
	}
	ids := options.GetIDs()

	del := squirrel.Delete(c.mainTable).
		Where(squirrel.Eq{"id": ids}).
		PlaceholderFormat(squirrel.Dollar)

	return del, nil
}

var _ store.CaseCommunicationStore = &CaseCommunicationStore{}

var CaseCommunicationFields = []string{"id", "ver", "communication_type", "communication_id"}

func NewCaseCommunicationStore(store *Store) (store.CaseCommunicationStore, error) {
	if store == nil {
		return nil, errors.Internal(
			"error creating case communication store, main store is nil")
	}
	return &CaseCommunicationStore{storage: store, mainTable: "cases.case_communication"}, nil
}
