package postgres

import (
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/webitel/cases/api/cases"
	authmodel "github.com/webitel/cases/auth/model"
	dberr "github.com/webitel/cases/internal/error"
	"github.com/webitel/cases/internal/store"
	"github.com/webitel/cases/internal/store/scanner"
	"github.com/webitel/cases/model"
)

type CaseCommunicationStore struct {
	storage   store.Store
	mainTable string
}

func (c *CaseCommunicationStore) Link(options *model.CreateOptions, communications []*cases.InputCaseCommunication) ([]*cases.CaseCommunication, error) {
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
	if len(res) < len(communications) && len(communications) == 1 {
		return nil, dberr.NewDBError("postgres.case_communication.result_processing.error", "wrong or duplicated communication_id or insufficient permissions")
	}
	return res, nil
}

func (c *CaseCommunicationStore) Unlink(options *model.DeleteOptions) ([]*cases.CaseCommunication, error) {
	base, plan, dbErr := c.buildDeleteCaseCommunicationSqlizer(options)
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
	rows, err := db.Query(options, sql, args...)
	if err != nil {
		return nil, dberr.NewDBError("postgres.case_communication.exec.error", err.Error())
	}
	res, dbErr := c.scanCommunications(rows, plan)
	if dbErr != nil {
		return nil, dbErr
	}
	if len(res) == 0 {
		return nil, dberr.NewDBError("postgres.case_communication.final_check.no_rows", "no rows were affected")
	}
	return res, nil
}

func (c *CaseCommunicationStore) List(opts *model.SearchOptions) ([]*cases.CaseCommunication, error) {
	//TODO implement me
	panic("implement me")
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
	insert := squirrel.Insert(c.mainTable).Columns("created_by", "created_at", "dc", "communication_type", "communication_id", "case_id").Suffix("RETURNING *")
	session := options.Session
	dc := session.GetDomainId()
	userId := session.GetUserId()
	roles := session.GetAclRoles()
	callsRbac := session.GetScope("calls").IsRbacUsed()

	for _, communication := range communications {
		dbErr := ValidateCaseCommunicationCreate(communication)
		if dbErr != nil {
			return nil, nil, dbErr
		}

		switch communication.CommunicationType {
		case cases.CaseCommunicationsTypes_COMMUNICATION_CALL:
			if callsRbac {
				insert = insert.Values(session.GetUserId(), options.CurrentTime(), dc, int32(communication.CommunicationType), squirrel.Expr(`(SELECT c.id::text
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
					dc, authmodel.Read, roles,
					roles), options.ParentID)
				continue
			}
		default:
		}

		insert = insert.Values(session.GetUserId(), options.CurrentTime(), session.GetDomainId(), int32(communication.CommunicationType), communication.CommunicationId, options.ParentID)
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

func (c *CaseCommunicationStore) buildDeleteCaseCommunicationSqlizer(options *model.DeleteOptions) (query squirrel.Sqlizer, plan []func(caseCommunication *cases.CaseCommunication) any, dbError *dberr.DBError) {
	if options == nil {
		return nil, nil, dberr.NewDBError("postgres.case_communication.build_delete_case_communication_sqlizer.check_args.options", "delete options required")
	}
	if len(options.IDs) == 0 {
		return nil, nil, dberr.NewDBError("postgres.case_communication.build_delete_case_communication_sqlizer.check_args.ids", "ids required to delete")
	}

	delCte := squirrel.Delete(c.mainTable).Where("id = ANY(?)", options.IDs).Suffix("RETURNING *")
	delAlias := "d"
	insertCte, args, err := store.FormAsCTE(delCte, delAlias)
	if err != nil {
		return nil, nil, dberr.NewDBError("postgres.case_communication.build.create_case_communication_sqlizer.form_cte.error", err.Error())
	}
	base := squirrel.Select().From(delAlias).Prefix(insertCte, args...).PlaceholderFormat(squirrel.Dollar)
	return c.buildSelectColumnsAndPlan(base, delAlias, CaseCommunicationFields)

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

func ValidateCaseCommunicationCreate(input *cases.InputCaseCommunication) *dberr.DBError {
	if input.CommunicationId == "" {
		return dberr.NewDBError("postgres.case_communication.validate_case_communication_create.validate.communication_id", "communication can't be empty")
	}
	if input.CommunicationType <= 0 {
		return dberr.NewDBError("postgres.case_communication.validate_case_communication_create.validate.type", "communication type can't be empty")
	}
	var typeFound bool
	for _, i := range cases.CaseCommunicationsTypes_value {
		if i == int32(input.CommunicationType) {
			typeFound = true
		}
	}
	if !typeFound {
		return dberr.NewDBError("postgres.case_communication.validate_case_communication_create.validate.type", "communication type not allowed")
	}
	return nil
}
