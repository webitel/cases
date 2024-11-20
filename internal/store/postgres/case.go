package postgres

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
	_go "github.com/webitel/cases/api/cases"
	dberr "github.com/webitel/cases/internal/error"
	"github.com/webitel/cases/internal/store"
	"github.com/webitel/cases/internal/store/scanner"
	"github.com/webitel/cases/model"
	util "github.com/webitel/cases/util"
)

type CaseStore struct {
	storage store.Store
}

type CaseScan func(caseItem *_go.Case) any

const (
	caseLeft = "c"
)

// Create implements store.CaseStore.
func (c *CaseStore) Create(rpc *model.CreateOptions, add *_go.Case) (*_go.Case, error) {
	// Get the database connection
	d, dbErr := c.storage.Database()
	if dbErr != nil {
		return nil, dberr.NewDBInternalError("postgres.case.create.database_connection_error", dbErr)
	}

	// Build the query
	selectBuilder, plan, err := c.buildCreateCaseSqlizer(rpc, add)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.case.create.build_query_error", err)
	}

	// Generate the SQL and arguments
	query, args, err := selectBuilder.ToSql()
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.case.create.query_to_sql_error", err)
	}

	query = store.CompactSQL(query)

	// Prepare the scan arguments
	scanArgs := convertToCaseScanArgs(plan, add)

	// Execute the query
	if err := d.QueryRow(rpc.Context, query, args...).Scan(scanArgs...); err != nil {
		return nil, dberr.NewDBInternalError("postgres.case.create.execution_error", err)
	}

	return add, nil
}

func (c *CaseStore) buildCreateCaseSqlizer(
	rpc *model.CreateOptions,
	caseItem *_go.Case,
) (sq.SelectBuilder, []CaseScan, error) {
	rpc.Fields = util.EnsureIdField(rpc.Fields)

	// CTEs for Service Catalog, SLA, and Status Condition
	serviceCatalogCTE := `
		service_catalog_cte AS (
			SELECT sc.sla_id
			FROM cases.service_catalog sc
			WHERE sc.id = $21 -- service_id
		)`

	slaCTE := `
		sla_cte AS (
			SELECT sla.id AS sla_id
			FROM cases.sla sla
			WHERE sla.id = (SELECT sla_id FROM service_catalog_cte)
		)`

	statusConditionCTE := `
		status_condition_cte AS (
			SELECT sc.id AS status_condition_id
			FROM cases.status_condition sc
			WHERE sc.status_id = $13 AND sc.initial = true
		)`

	insertBuilder := sq.Insert("cases.case").
		Columns("rating", "dc", "created_at", "created_by", "updated_at", "updated_by",
			"close_result", "priority", "source", "close_reason",
			"rating_comment", "name", "status", "close_reason_group", "\"group\"",
			"subject", "planned_reaction_at", "planned_resolve_at", "reporter",
			"impacted", "service", "description", "assignee", "sla", "status_condition").
		Values(
			caseItem.Rate.GetRating(),
			rpc.Session.GetDomainId(),
			rpc.CurrentTime(),
			rpc.Session.GetUserId(),
			rpc.CurrentTime(),
			rpc.Session.GetUserId(),
			caseItem.Close.CloseResult,
			caseItem.Priority.GetId(),
			caseItem.Source.GetId(),
			caseItem.Close.CloseReason.GetId(),
			caseItem.Rate.GetRatingComment(),
			caseItem.Name,
			caseItem.Status.GetId(),
			caseItem.CloseReasonGroup.GetId(),
			caseItem.Group.GetId(),
			caseItem.Subject,
			10, // Mocked planned_reaction_at
			10, // Mocked planned_resolve_at
			caseItem.Reporter.GetId(),
			caseItem.Impacted.GetId(),
			caseItem.Service.GetId(),
			sq.Expr("NULLIF(?, '')", caseItem.Description),
			caseItem.Assignee.GetId(),
			sq.Expr("(SELECT sla_id FROM sla_cte)"),
			sq.Expr("(SELECT status_condition_id FROM status_condition_cte)")).
		PlaceholderFormat(sq.Dollar).
		Suffix("RETURNING *")

	insertSQL, insertArgs, err := insertBuilder.ToSql()
	if err != nil {
		return sq.SelectBuilder{}, nil, dberr.NewDBInternalError("postgres.case.create.insert_query_build_error", err)
	}

	// Construct the final CTE query
	cteSQL := fmt.Sprintf(`
		WITH %s,
		%s,
		%s,
		c AS (%s)
	`, serviceCatalogCTE, slaCTE, statusConditionCTE, insertSQL)

	// Build the SELECT query to fetch the returned columns
	selectBuilder, plan, err := c.buildCaseSelectColumnsAndPlan(sq.Select(), rpc.Fields)
	if err != nil {
		return sq.SelectBuilder{}, nil, err
	}

	// Combine the CTE with the SELECT query
	selectBuilder = selectBuilder.
		PrefixExpr(sq.Expr(cteSQL, insertArgs...)).
		From("c")

	return selectBuilder, plan, nil
}

func (c *CaseStore) buildCaseSelectColumnsAndPlan(
	base sq.SelectBuilder,
	fields []string,
) (sq.SelectBuilder, []CaseScan, error) {
	var plan []CaseScan

	for _, field := range fields {
		switch field {
		case "id":
			base = base.Column(store.Ident(caseLeft, "id AS case_id"))
			plan = append(plan, func(caseItem *_go.Case) any {
				return &caseItem.Id
			})
		case "ver":
			base = base.Column(store.Ident(caseLeft, "ver"))
			plan = append(plan, func(caseItem *_go.Case) any {
				return &caseItem.Ver
			})
		case "created_by":
			base = base.Column(fmt.Sprintf(
				"(SELECT ROW(wu.id, wu.name)::text FROM directory.wbt_user wu WHERE wu.id = %s.created_by) AS created_by", caseLeft))
			plan = append(plan, func(caseItem *_go.Case) any {
				return scanner.ScanRowLookup(&caseItem.CreatedBy)
			})
		case "created_at":
			base = base.Column(store.Ident(caseLeft, "created_at"))
			plan = append(plan, func(caseItem *_go.Case) any {
				return scanner.ScanTimestamp(&caseItem.CreatedAt)
			})
		case "updated_by":
			base = base.Column(fmt.Sprintf(
				"(SELECT ROW(wu.id, wu.name)::text FROM directory.wbt_user wu WHERE wu.id = %s.updated_by) AS updated_by", caseLeft))
			plan = append(plan, func(caseItem *_go.Case) any {
				return scanner.ScanRowLookup(&caseItem.UpdatedBy)
			})
		case "updated_at":
			base = base.Column(store.Ident(caseLeft, "updated_at"))
			plan = append(plan, func(caseItem *_go.Case) any {
				return scanner.ScanTimestamp(&caseItem.UpdatedAt)
			})
		case "name":
			base = base.Column(store.Ident(caseLeft, "name"))
			plan = append(plan, func(caseItem *_go.Case) any {
				return &caseItem.Name
			})
		case "subject":
			base = base.Column(store.Ident(caseLeft, "subject"))
			plan = append(plan, func(caseItem *_go.Case) any {
				return &caseItem.Subject
			})
		case "description":
			base = base.Column(store.Ident(caseLeft, "description"))
			plan = append(plan, func(caseItem *_go.Case) any {
				return &caseItem.Description
			})
		case "group":
			base = base.Column(fmt.Sprintf(
				"(SELECT ROW(g.id, g.name)::text FROM contacts.group g WHERE g.id = %s.group) AS contact_group", caseLeft))
			plan = append(plan, func(caseItem *_go.Case) any {
				return scanner.ScanRowLookup(&caseItem.Group)
			})
		case "source":
			base = base.Column(fmt.Sprintf(
				"(SELECT ROW(src.id, src.name)::text FROM cases.source src WHERE src.id = %s.source) AS source", caseLeft))
			plan = append(plan, func(caseItem *_go.Case) any {
				return scanner.ScanRowLookup(&caseItem.Source)
			})
		case "planned_reaction_at":
			base = base.Column(store.Ident(caseLeft, "planned_reaction_at"))
			plan = append(plan, func(caseItem *_go.Case) any {
				// return scanner.ScanTimestamp(&caseItem.PlannedReactionAt)
				return &caseItem.PlannedReactionAt
			})
		case "planned_resolve_at":
			base = base.Column(store.Ident(caseLeft, "planned_resolve_at"))
			plan = append(plan, func(caseItem *_go.Case) any {
				// return scanner.ScanTimestamp(&caseItem.PlannedResolveAt) -- Need to be implemented
				return &caseItem.PlannedResolveAt
			})
		case "close_reason_group":
			base = base.Column(fmt.Sprintf(
				"(SELECT ROW(crg.id, crg.name)::text FROM cases.close_reason_group crg WHERE crg.id = %s.close_reason_group) AS close_reason_group", caseLeft))
			plan = append(plan, func(caseItem *_go.Case) any {
				return scanner.ScanRowLookup(&caseItem.CloseReasonGroup)
			})
		case "author":
			base = base.Column(fmt.Sprintf(`
				(SELECT
					ROW(ca.id, ca.common_name)::text
				FROM directory.wbt_user wu
				LEFT JOIN contacts.contact ca ON wu.contact_id = ca.id
				WHERE wu.id = %s.created_by AND ca.id IS NOT NULL) AS author`, caseLeft))
			plan = append(plan, func(caseItem *_go.Case) any {
				return scanner.ScanRowLookup(&caseItem.Author)
			})

		case "close_result":
			base = base.Column(store.Ident(caseLeft, "close_result"))
			plan = append(plan, func(caseItem *_go.Case) any {
				return &caseItem.Close.CloseResult
			})
		case "close_reason":
			base = base.Column(fmt.Sprintf(
				"(SELECT ROW(cr.id, cr.name)::text FROM cases.close_reason cr WHERE cr.id = %s.close_reason) AS close_reason", caseLeft))
			plan = append(plan, func(caseItem *_go.Case) any {
				return scanner.ScanRowLookup(&caseItem.Close.CloseReason)
			})
		case "rating":
			base = base.Column(store.Ident(caseLeft, "rating"))
			plan = append(plan, func(caseItem *_go.Case) any {
				return &caseItem.Rate.Rating
			})
		case "rating_comment":
			base = base.Column(store.Ident(caseLeft, "rating_comment"))
			plan = append(plan, func(caseItem *_go.Case) any {
				return &caseItem.Rate.RatingComment
			})
		case "sla":
			base = base.Column(fmt.Sprintf(
				"(SELECT ROW(sla.id, sla.name)::text FROM cases.sla sla WHERE sla.id = %s.sla) AS sla", caseLeft))
			plan = append(plan, func(caseItem *_go.Case) any {
				return scanner.ScanRowLookup(&caseItem.Sla)
			})
		case "status_condition":
			base = base.Column(fmt.Sprintf(`
				(SELECT ROW(stc.id, stc.name)::text
				 FROM cases.status_condition stc
				 WHERE stc.id = %s.status_condition) AS status_condition`, caseLeft))
			plan = append(plan, func(caseItem *_go.Case) any {
				return scanner.ScanRowLookup(&caseItem.StatusCondition)
			})
		case "status":
			base = base.Column(fmt.Sprintf(`
				(SELECT
					ROW(st.id, st.name)::text
				FROM cases.status st
				WHERE st.id = %s.status) AS status`, caseLeft))
			plan = append(plan, func(caseItem *_go.Case) any {
				return scanner.ScanRowLookup(&caseItem.Status)
			})
		case "priority":
			base = base.Column(fmt.Sprintf(
				"(SELECT ROW(p.id, p.name)::text FROM cases.priority p WHERE p.id = %s.priority) AS priority", caseLeft))
			plan = append(plan, func(caseItem *_go.Case) any {
				return scanner.ScanRowLookup(&caseItem.Priority)
			})
		case "assignee":
			base = base.Column(fmt.Sprintf(
				"(SELECT ROW(a.id, a.common_name)::text FROM contacts.contact a WHERE a.id = %s.assignee) AS assignee", caseLeft))
			plan = append(plan, func(caseItem *_go.Case) any {
				return scanner.ScanRowLookup(&caseItem.Assignee)
			})

		case "reporter":
			base = base.Column(fmt.Sprintf(
				"(SELECT ROW(r.id, r.common_name)::text FROM contacts.contact r WHERE r.id = %s.reporter) AS reporter", caseLeft))
			plan = append(plan, func(caseItem *_go.Case) any {
				return scanner.ScanRowLookup(&caseItem.Reporter)
			})
		case "impacted":
			base = base.Column(fmt.Sprintf(
				"(SELECT ROW(i.id, i.common_name)::text FROM contacts.contact i WHERE i.id = %s.impacted) AS impacted", caseLeft))
			plan = append(plan, func(caseItem *_go.Case) any {
				return scanner.ScanRowLookup(&caseItem.Impacted)
			})
		case "sla_conditions":
			base = base.Column(`
				(SELECT JSON_AGG(JSON_BUILD_OBJECT(
					'id', sc.id,
					'name', sc.name
				)) FROM cases.sla_condition sc
				WHERE sc.sla_id = (SELECT sla_id FROM sla_cte)) AS sla_conditions`)
			plan = append(plan, func(caseItem *_go.Case) any {
				return scanner.ScanJSONToStructList(&caseItem.SlaCondition)
			})
		case "comments":
			base = base.Column(fmt.Sprintf(`
				(SELECT JSON_AGG(JSON_BUILD_OBJECT(
					'id', cc.id,
					'comment', cc.comment,
					'created_by', JSON_BUILD_OBJECT('id', cc.created_by, 'name', cn.name),
					'updated_by', JSON_BUILD_OBJECT('id', cc.updated_by, 'name', un.name),
					'updated_at', cc.updated_at
				)) FROM cases.case_comment cc
				LEFT JOIN directory.wbt_user cn ON cc.created_by = cn.id
				LEFT JOIN directory.wbt_user un ON cc.updated_by = un.id
				WHERE cc.case_id = %s.id) AS comments`, caseLeft))
			plan = append(plan, func(caseItem *_go.Case) any {
				return scanner.ScanJSONToStructList(&caseItem.Comments.Items)
			})
		case "links":
			base = base.Column(fmt.Sprintf(`
				(SELECT JSON_AGG(JSON_BUILD_OBJECT(
					'id', cl.id,
					'url', cl.url,
					'name', cl.name,
					'created_by', JSON_BUILD_OBJECT('id', cl.created_by, 'name', u.name),
					'created_at', cl.created_at
				)) FROM cases.case_link cl
				LEFT JOIN directory.wbt_user u ON cl.created_by = u.id
				WHERE cl.case_id = %s.id) AS links`, caseLeft))
			plan = append(plan, func(caseItem *_go.Case) any {
				return scanner.ScanJSONToStructList(&caseItem.Links.Items)
			})

		case "related_cases":
			base = base.Column(fmt.Sprintf(`
				(SELECT JSON_AGG(JSON_BUILD_OBJECT(
					'id', rc.related_case_id,
					'name', c.name,
					'subject', c.subject,
					'description', c.description
				)) FROM cases.related_case rc
				JOIN cases.case c ON rc.related_case_id = c.id
				WHERE rc.case_id = %s.id) AS related_cases`, caseLeft))
			plan = append(plan, func(caseItem *_go.Case) any {
				return scanner.ScanJSONToStructList(&caseItem.Related.Items)
			})
		default:
			return sq.SelectBuilder{}, nil, fmt.Errorf("unknown field: %s", field)
		}
	}

	if len(plan) == 0 {
		return sq.SelectBuilder{}, nil, fmt.Errorf("no fields specified for selection")
	}

	return base, plan, nil
}

// Helper function to convert the scan plan to arguments for scanning.
func convertToCaseScanArgs(plan []CaseScan, caseItem *_go.Case) []any {
	var scanArgs []any
	for _, scan := range plan {
		scanArgs = append(scanArgs, scan(caseItem))
	}
	return scanArgs
}

// Delete implements store.CaseStore.
func (c *CaseStore) Delete(req *model.DeleteOptions) (*_go.Case, error) {
	panic("unimplemented")
}

// List implements store.CaseStore.
func (c *CaseStore) List(rpc *model.SearchOptions) (*_go.CaseList, error) {
	panic("unimplemented")
}

// Merge implements store.CaseStore.
func (c *CaseStore) Merge(req *model.CreateOptions) (*_go.CaseList, error) {
	panic("unimplemented")
}

// Update implements store.CaseStore.
func (c *CaseStore) Update(req *model.UpdateOptions) (*_go.Case, error) {
	panic("unimplemented")
}

func NewCaseStore(store store.Store) (store.CaseStore, error) {
	if store == nil {
		return nil, dberr.NewDBError("postgres.new_case.check.bad_arguments",
			"error creating case interface to the case table, main store is nil")
	}
	return &CaseStore{storage: store}, nil
}
