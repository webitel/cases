package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	authmodel "github.com/webitel/cases/auth/model"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx"
	"github.com/lib/pq"

	_go "github.com/webitel/cases/api/cases"
	dberr "github.com/webitel/cases/internal/error"
	"github.com/webitel/cases/internal/store"
	"github.com/webitel/cases/internal/store/scanner"
	"github.com/webitel/cases/model"
	util "github.com/webitel/cases/util"
)

type CaseStore struct {
	storage   store.Store
	mainTable string
}

type CaseScan func(caseItem *_go.Case) any

const (
	caseLeft     = "c"
	relatedAlias = "related"
	linksAlias   = "links"
)

func (c *CaseStore) Create(
	rpc *model.CreateOptions,
	add *_go.Case,
) (*_go.Case, error) {
	// Get the database connection
	d, dbErr := c.storage.Database()
	if dbErr != nil {
		return nil, dberr.NewDBInternalError("postgres.case.create.database_connection_error", dbErr)
	}

	// Begin a transaction
	tx, err := d.Begin(rpc.Context)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.case.create.transaction_error", err)
	}
	defer tx.Rollback(rpc.Context)
	txManager := store.NewTxManager(tx)

	// Scan SLA details
	// Sla_id
	// reaction_at & resolve_at in [milli]seconds
	slaID, slaConditionID, reaction_at, resolve_at, calendarID, err := c.ScanSla(
		rpc,
		txManager,
		add.Service.GetId(),
		add.Priority.GetId(),
	)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.case.create.scan_sla_error", err)
	}

	// Calculate planned times within the transaction
	err = c.calculatePlannedReactionAndResolutionTime(rpc, calendarID, reaction_at, resolve_at, txManager, add)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.case.create.calculate_planned_times_error", err)
	}

	// Build the query
	selectBuilder, plan, err := c.buildCreateCaseSqlizer(rpc, add, slaID, slaConditionID)
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
	if err = txManager.QueryRow(rpc.Context, query, args...).Scan(scanArgs...); err != nil {
		return nil, dberr.NewDBInternalError("postgres.case.create.execution_error", err)
	}

	// Commit the transaction
	if err := tx.Commit(rpc.Context); err != nil {
		return nil, dberr.NewDBInternalError("postgres.case.create.commit_error", err)
	}

	return add, nil
}

// ScanSla fetches the SLA ID, reaction time, resolution time, calendar ID, and SLA condition ID for the last child service with a non-NULL SLA ID.
func (c *CaseStore) ScanSla(
	rpc *model.CreateOptions,
	txManager *store.TxManager,
	serviceID int64,
	priorityID int64,
) (
	slaID,
	slaConditionID,
	reactionTime,
	resolutionTime,
	calendarID int,
	err error,
) {
	// var slaId, reactionTime, resolutionTime, calendarId, slaConditionId int

	err = txManager.QueryRow(rpc.Context, `
WITH RECURSIVE
    service_hierarchy AS (SELECT id, root_id, sla_id, 1 AS level
                          FROM cases.service_catalog
                          WHERE id = $1 -- Start with the given service ID provided as $1

                          UNION ALL

                          SELECT sc.id, sc.root_id, sc.sla_id, sh.level + 1
                          FROM cases.service_catalog sc
                                   INNER JOIN service_hierarchy sh ON sc.root_id = sh.id
        -- Recursively traverse downward to find all child services, incrementing the level with each step
    ),
    valid_sla_hierarchy AS (SELECT sh.id AS service_id, -- Current service ID
                                   sh.root_id,          -- Parent service ID
                                   sh.sla_id,           -- SLA ID for the current service
                                   sh.level,            -- Depth level in the hierarchy
                                   sla.reaction_time,   -- Reaction time from the SLA
                                   sla.resolution_time, -- Resolution time from the SLA
                                   sla.calendar_id      -- Calendar ID associated with the SLA
                            FROM service_hierarchy sh
                                     LEFT JOIN cases.sla sla ON sh.sla_id = sla.id
                            WHERE sh.sla_id IS NOT NULL -- Keep only services with non-NULL SLA
        -- Here, we extract details of all services with SLAs, preparing them for prioritization
    ),
    deepest_sla
        AS (SELECT DISTINCT ON (sh.level, sh.id) sh.id AS service_id, -- Service ID for the deepest child or nearest valid SLA
                                                 sh.root_id,          -- Parent service ID
                                                 sh.sla_id,           -- SLA ID for the selected service
                                                 sh.level,            -- Depth level in the hierarchy
                                                 sla.reaction_time,   -- Reaction time from SLA
                                                 sla.resolution_time, -- Resolution time from SLA
                                                 sla.calendar_id      -- Calendar ID associated with the SLA
            FROM service_hierarchy sh
                     LEFT JOIN cases.sla sla ON sh.sla_id = sla.id
            ORDER BY sh.level DESC, sh.id
        -- Select the "deepest" child service by level, falling back to the next service upward if necessary
    ),
    priority_condition AS (SELECT sc.id AS sla_condition_id, -- Fetch the SLA condition ID
                                  sc.reaction_time,
                                  sc.resolution_time
                           FROM cases.sla_condition sc
                                    INNER JOIN cases.priority_sla_condition psc ON sc.id = psc.sla_condition_id
                                    INNER JOIN deepest_sla ON sc.sla_id = deepest_sla.sla_id
                           WHERE psc.priority_id = $2 -- Match the given priority ID provided as $2
                           LIMIT 1
        -- Extract reaction and resolution times from SLA conditions if a priority-specific condition exists
    )
SELECT deepest_sla.sla_id,                                                                           -- Final SLA ID
       COALESCE(priority_condition.reaction_time, deepest_sla.reaction_time)     AS reaction_time,
       -- Use priority-specific reaction time if available, otherwise fall back to SLA reaction time
       COALESCE(priority_condition.resolution_time, deepest_sla.resolution_time) AS resolution_time,
       -- Use priority-specific resolution time if available, otherwise fall back to SLA resolution time
       deepest_sla.calendar_id,                                                                      -- Calendar ID associated with the final SLA
       COALESCE(priority_condition.sla_condition_id, 0)                          AS sla_condition_id -- Return SLA condition ID if a priority match is found
FROM deepest_sla
         LEFT JOIN priority_condition ON true;
-- Combine the results to ensure we have reaction and resolution times even if no priority-specific condition exists

	`, serviceID, priorityID).Scan(
		&slaID,
		&reactionTime,
		&resolutionTime,
		&calendarID,
		&slaConditionID,
	)
	if err != nil {
		return 0, 0, 0, 0, 0, dberr.NewDBInternalError("failed to scan SLA: %w", err)
	}

	return slaID, slaConditionID, reactionTime, resolutionTime, calendarID, nil
}

func (c *CaseStore) buildCreateCaseSqlizer(
	rpc *model.CreateOptions,
	caseItem *_go.Case,
	sla int,
	slaCondition int,
) (sq.SelectBuilder, []CaseScan, error) {
	// Parameters for the main case and nested JSON arrays
	params := map[string]interface{}{
		// Case-level parameters
		"date":                rpc.CurrentTime(),
		"user":                rpc.Session.GetUserId(),
		"dc":                  rpc.Session.GetDomainId(),
		"sla":                 sla,
		"sla_condition":       slaCondition,
		"status":              caseItem.Status.GetId(),
		"service":             caseItem.Service.GetId(),
		"rating":              caseItem.Rate.GetRating(),
		"close_result":        caseItem.Close.CloseResult,
		"priority":            caseItem.Priority.GetId(),
		"source":              caseItem.Source.GetId(),
		"close_reason":        caseItem.Close.CloseReason.GetId(),
		"contact_group":       caseItem.Group.GetId(),
		"close_reason_group":  caseItem.CloseReasonGroup.GetId(),
		"subject":             caseItem.Subject,
		"planned_reaction_at": util.LocalTime(caseItem.PlannedReactionAt),
		"planned_resolve_at":  util.LocalTime(caseItem.PlannedResolveAt),
		"reporter":            caseItem.Reporter.GetId(),
		"impacted":            caseItem.Impacted.GetId(),
		"description":         caseItem.Description,
		"assignee":            caseItem.Assignee.GetId(),
		//-------------------------------------------------//
		//------ CASE One-to-Many ( 1 : n ) Attributes ----//
		//-------------------------------------------------//
		// Links and related cases as JSON arrays
		"links":   extractLinksJSON(caseItem.Links),
		"related": extractRelatedJSON(caseItem.Related),
	}

	// Define CTEs for the main case
	statusConditionCTE := `
		status_condition_cte AS (
			SELECT sc.id AS status_condition_id
			FROM cases.status_condition sc
			WHERE sc.status_id = :status AND sc.initial = true
		)`

	prefixCTE := `
		prefix_cte AS (
			SELECT prefix
			FROM cases.service_catalog
			WHERE id = :service
			LIMIT 1
		), id_cte AS (
			SELECT nextval('cases.case_id'::regclass) AS id
		)`

	// Consolidated query for inserting the case, links, and related cases
	query := `
	WITH
		` + statusConditionCTE + `,
		` + prefixCTE + `,
		` + caseLeft + ` AS (
			INSERT INTO cases.case (
				id, name, rating, dc, created_at, created_by, updated_at, updated_by, close_result,
				priority, source, close_reason, status, contact_group, close_reason_group,
				subject, planned_reaction_at, planned_resolve_at, reporter, impacted,
				service, description, assignee, sla, sla_condition_id, status_condition
			) VALUES (
				(SELECT id FROM id_cte),
				CONCAT((SELECT prefix FROM prefix_cte), '_', (SELECT id FROM id_cte)),
				:rating, :dc, :date, :user, :date, :user, :close_result,
				:priority, :source, :close_reason, :status, :contact_group, :close_reason_group,
				:subject, :planned_reaction_at, :planned_resolve_at, :reporter, :impacted,
				:service, :description, :assignee, :sla, :sla_condition,
				(SELECT status_condition_id FROM status_condition_cte)
			)
			RETURNING *
		),
		` + linksAlias + ` AS (
			INSERT INTO cases.case_link (
				name, url, dc, created_by, created_at, updated_by, updated_at, case_id
			)
			SELECT
				item ->> 'name',
				item ->> 'url',
				:dc, :user, :date, :user, :date, (SELECT id FROM ` + caseLeft + `)
			FROM jsonb_array_elements(:links) AS item
		),
		` + relatedAlias + ` AS (
			INSERT INTO cases.related_case (
				parent_case_id, child_case_id, relation_type, dc, created_by, created_at, updated_by, updated_at
			)
			SELECT
				(SELECT id FROM ` + caseLeft + `),
				(item ->> 'id')::bigint,
				(item ->> 'type')::int,
				:dc, :user, :date, :user, :date
			FROM jsonb_array_elements(:related) AS item
		)
	`

	// **Bind named query and parameters**
	// **This binds the named parameters in the query to the provided params map, converting it into a positional query with arguments.**
	// **Example:**
	// **  Query: "INSERT INTO cases.case (name, subject) VALUES (:name, :subject)"**
	// **  Params: map[string]interface{}{"name": "test_name", "subject": "test_subject"}**
	// **  Result: "INSERT INTO cases.case (name, col2) subject ($1, $2)", []interface{}{"test_name", "test_subject"}**
	boundQuery, args, err := store.BindNamed(query, params)
	if err != nil {
		return sq.SelectBuilder{}, nil, dberr.NewDBInternalError("postgres.case.create.bind_named_error", err)
	}

	// Construct SELECT query to return case data
	selectBuilder, plan, err := c.buildCaseSelectColumnsAndPlan(
		rpc.Session,
		sq.Select().PrefixExpr(sq.Expr(boundQuery, args...)),
		rpc.Fields,
	)
	if err != nil {
		return sq.SelectBuilder{}, nil, dberr.NewDBInternalError("postgres.case.create.build_select_query_error", err)
	}

	selectBuilder = selectBuilder.From(caseLeft)

	return selectBuilder, plan, nil
}

// Helper functions to generate JSON arrays for links and related cases
func extractLinksJSON(links *_go.CaseLinkList) []byte {
	if links == nil || len(links.Items) == 0 {
		return []byte("[]")
	}
	var jsonArray []map[string]interface{}
	for _, link := range links.Items {
		jsonArray = append(jsonArray, map[string]interface{}{
			"name": link.Name,
			"url":  link.Url,
		})
	}
	jsonData, _ := json.Marshal(jsonArray)
	return jsonData
}

func extractRelatedJSON(related *_go.RelatedCaseList) []byte {
	if related == nil || len(related.Items) == 0 {
		return []byte("[]")
	}
	var jsonArray []map[string]interface{}
	for _, item := range related.Items {
		jsonArray = append(jsonArray, map[string]interface{}{
			"id":   item.GetId(),
			"type": item.GetRelationType(),
		})
	}
	jsonData, _ := json.Marshal(jsonArray)
	return jsonData
}

// ConvertRelationType validates the cases.RelationType and returns its integer representation.
func ConvertRelationType(relationType _go.RelationType) (int, error) {
	switch relationType {
	case _go.RelationType_BlockedBy:
		return 0, nil
	case _go.RelationType_Blocks:
		return 1, nil
	case _go.RelationType_Duplicates:
		return 2, nil
	case _go.RelationType_DuplicatedBy:
		return 3, nil
	case _go.RelationType_Causes:
		return 4, nil
	case _go.RelationType_CausedBy:
		return 5, nil
	case _go.RelationType_IsChildOf:
		return 6, nil
	case _go.RelationType_IsParentOf:
		return 7, nil
	case _go.RelationType_RelatesTo:
		return 8, nil
	default:
		return -1, fmt.Errorf("invalid relation type: %v", relationType)
	}
}

func (c *CaseStore) calculatePlannedReactionAndResolutionTime(
	rpc *model.CreateOptions,
	calendarID int,
	reactionTime int,
	resolutionTime int,
	txManager *store.TxManager,
	caseItem *_go.Case,
) error {
	rows, err := txManager.Query(rpc.Context, `
		SELECT day, start_time_of_day, end_time_of_day, special, disabled
		FROM flow.calendar, UNNEST(accepts::flow.calendar_accept_time[]) x
		WHERE id = $1
		ORDER BY day, start_time_of_day`, calendarID)
	if err != nil {
		return fmt.Errorf("failed to fetch calendar details: %w", err)
	}
	defer rows.Close()

	var calendar []struct {
		Day            int
		StartTimeOfDay int
		EndTimeOfDay   int
		Special        bool
		Disabled       bool
	}
	for rows.Next() {
		var entry struct {
			Day            int
			StartTimeOfDay int
			EndTimeOfDay   int
			Special        bool
			Disabled       bool
		}
		if err = rows.Scan(
			&entry.Day,
			&entry.StartTimeOfDay,
			&entry.EndTimeOfDay,
			&entry.Special,
			&entry.Disabled,
		); err != nil {
			return fmt.Errorf("failed to scan calendar entry: %w", err)
		}
		if !entry.Disabled {
			calendar = append(calendar, entry)
		}
	}
	if err = rows.Err(); err != nil {
		return fmt.Errorf("error iterating over calendar rows: %w", err)
	}

	// Convert reaction and resolution times from milliseconds to minutes
	reactionMinutes := reactionTime / 60000
	resolutionMinutes := resolutionTime / 60000

	currentTime := rpc.CurrentTime()
	reactionTimestamp, err := calculateTimestampFromCalendar(currentTime, reactionMinutes, calendar)
	if err != nil {
		return fmt.Errorf("failed to calculate planned reaction time: %w", err)
	}

	//?? TODO
	// resolveTimestamp, err := calculateTimestampFromCalendar(reactionTimestamp, resolutionMinutes, calendar)
	resolveTimestamp, err := calculateTimestampFromCalendar(currentTime, resolutionMinutes, calendar)
	if err != nil {
		return fmt.Errorf("failed to calculate planned resolution time: %w", err)
	}

	caseItem.PlannedReactionAt = util.Timestamp(reactionTimestamp)
	caseItem.PlannedResolveAt = util.Timestamp(resolveTimestamp)

	return nil
}

func calculateTimestampFromCalendar(
	startTime time.Time,
	requiredMinutes int,
	calendar []struct {
		Day            int
		StartTimeOfDay int
		EndTimeOfDay   int
		Special        bool
		Disabled       bool
	},
) (time.Time, error) {
	remainingMinutes := requiredMinutes
	currentDay := int(startTime.Weekday())
	currentTimeInMinutes := startTime.Hour()*60 + startTime.Minute()

	for {
		for _, slot := range calendar {
			// Match the current day and ensure the slot is in the future
			if slot.Day == currentDay && slot.StartTimeOfDay >= currentTimeInMinutes {
				availableMinutes := slot.EndTimeOfDay - slot.StartTimeOfDay
				if availableMinutes >= remainingMinutes {
					// Calculate the exact timestamp
					return startTime.Add(time.Duration(remainingMinutes) * time.Minute), nil
				}
				remainingMinutes -= availableMinutes
				currentTimeInMinutes = slot.EndTimeOfDay // Move to the end of the current slot
			}
		}

		// If no slots available, move to the next day
		currentDay = (currentDay + 1) % 7 // Wrap around to the start of the week if necessary
		currentTimeInMinutes = 0          // Reset to start of the day
		startTime = startTime.Add(24 * time.Hour)
	}
}

// Delete implements store.CaseStore.
func (c *CaseStore) Delete(rpc *model.DeleteOptions) error {
	// Establish database connection
	d, err := c.storage.Database()
	if err != nil {
		return dberr.NewDBInternalError("store.case.delete.database_connection_error", err)
	}

	// Build the delete query
	query, args, dbErr := c.buildDeleteCaseQuery(rpc)
	if dbErr != nil {
		return dberr.NewDBInternalError("store.case.delete.query_build_error", dbErr)
	}

	// Execute the query
	res, execErr := d.Exec(rpc.Context, query, args...)
	if execErr != nil {
		return dberr.NewDBInternalError("store.case.delete.exec_error", execErr)
	}

	// Check if any rows were affected
	if res.RowsAffected() == 0 {
		return dberr.NewDBNoRowsError("store.case.delete.not_found")
	}

	return nil
}

func (c CaseStore) buildDeleteCaseQuery(rpc *model.DeleteOptions) (string, []interface{}, error) {
	convertedIds := util.Int64SliceToStringSlice(rpc.IDs)
	ids := util.FieldsFunc(convertedIds, util.InlineFields)

	query := deleteCaseQuery
	args := []interface{}{
		pq.Array(ids),
		rpc.Session.GetDomainId(),
	}
	return query, args, nil
}

var deleteCaseQuery = store.CompactSQL(`
	DELETE FROM cases.case
	WHERE id = ANY($1) AND dc = $2
`)

// List implements store.CaseStore.
func (c *CaseStore) List(opts *model.SearchOptions) (*_go.CaseList, error) {
	if opts == nil {
		return nil, dberr.NewDBError("postgres.case.list.check_args.opts", "search options required")
	}
	query, plan, err := c.buildListCaseSqlizer(opts)
	if err != nil {
		return nil, err
	}
	slct, args, err := query.ToSql()
	if err != nil {
		return nil, dberr.NewDBError("postgres.case.list.to_sql.error", err.Error())
	}
	db, dbErr := c.storage.Database()
	if dbErr != nil {
		return nil, dbErr
	}
	rows, err := db.Query(opts.Context, store.CompactSQL(slct), args...)
	if err != nil {
		return nil, dberr.NewDBError("postgres.case.list.exec.error", err.Error())
	}
	var (
		res _go.CaseList
		i   int
	)
	for ; rows.Next(); i++ {
		if i > int(opts.GetSize())-1 {
			res.Next = true
			res.Page = int64(opts.GetPage())
			break
		}
		var node _go.Case
		scanArgs := convertToCaseScanArgs(plan, &node)
		err = rows.Scan(scanArgs...)
		if err != nil {
			return nil, dberr.NewDBError("postgres.case.list.scan.error", err.Error())
		}
		res.Items = append(res.Items, &node)
	}
	return &res, nil
}

func (c *CaseStore) buildListCaseSqlizer(opts *model.SearchOptions) (sq.SelectBuilder, []CaseScan, error) {
	base := sq.Select().From(fmt.Sprintf("%s %s", c.mainTable, caseLeft)).PlaceholderFormat(sq.Dollar)
	base, plan, err := c.buildCaseSelectColumnsAndPlan(opts.Session, base, opts.Fields)
	if err != nil {
		return base, nil, err
	}

	base = base.Where(store.Ident(caseLeft, "dc = ?"), opts.Session.GetDomainId())
	if opts.Search != "" {
		base = store.AddSearchTerm(base, store.Ident(caseLeft, "name"), store.Ident(caseLeft, "subject"), store.Ident(caseLeft, "contact_info"))
	}
	// pagination
	base = store.ApplyPaging(opts, base)

	// sort
	base = store.ApplyDefaultSorting(opts, base)

	return base, plan, nil
}

func (c *CaseStore) Update(
	rpc *model.UpdateOptions,
	upd *_go.Case,
) (*_go.Case, error) {
	// Establish database connection
	db, err := c.storage.Database()
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.case.update.database_connection_error", err)
	}

	// Begin a transaction
	tx, txErr := db.Begin(rpc.Context)
	if txErr != nil {
		return nil, dberr.NewDBInternalError("postgres.cases.case.update.transaction_error", txErr)
	}
	defer tx.Rollback(rpc.Context)
	txManager := store.NewTxManager(tx)

	// Scan the current version of the comment
	ver, verErr := c.ScanVer(rpc.Context, upd.Id, txManager)
	if verErr != nil {
		return nil, verErr
	}

	if upd.Ver != int32(ver) {
		return nil, dberr.NewDBConflictError("postgres.cases.case.update.version_mismatch", "Version mismatch, update failed")
	}

	// Build the SQL query and scan plan
	queryBuilder, plan, sqErr := c.buildUpdateCaseSqlizer(rpc, upd)
	if sqErr != nil {
		return nil, dberr.NewDBInternalError("postgres.case.update.query_build_error", sqErr)
	}

	// Generate the SQL and arguments
	query, args, sqErr := queryBuilder.ToSql()
	if sqErr != nil {
		return nil, dberr.NewDBInternalError("postgres.case.update.query_to_sql_error", sqErr)
	}

	query = store.CompactSQL(query)

	// Prepare scan arguments
	scanArgs := convertToCaseScanArgs(plan, upd)

	// Execute the query
	if sqErr = db.QueryRow(rpc.Context, query, args...).Scan(scanArgs...); sqErr != nil {
		return nil, dberr.NewDBInternalError("postgres.case.update.execution_error", sqErr)
	}

	// Commit the transaction
	if err := tx.Commit(rpc.Context); err != nil {
		return nil, dberr.NewDBInternalError("postgres.cases.case.update.commit_error", err)
	}

	return upd, nil
}

func (c *CaseStore) ScanVer(
	ctx context.Context,
	caseID int64,
	txManager *store.TxManager,
) (int64, error) {
	// Retrieve the current version (`ver`) of the case
	var ver int64
	err := txManager.QueryRow(ctx, "SELECT ver FROM cases.case WHERE id = $1", caseID).Scan(&ver)
	if err != nil {
		if err == pgx.ErrNoRows {
			// Return a specific error if no case with the given ID is found
			return 0, dberr.NewDBNotFoundError("postgres.cases.case.scan_ver.not_found", "Case not found")
		}
		return 0, dberr.NewDBInternalError("postgres.cases.case.scan_ver.query_error", err)
	}
	return ver, nil
}

func (c *CaseStore) buildUpdateCaseSqlizer(
	req *model.UpdateOptions,
	upd *_go.Case,
) (sq.Sqlizer, []CaseScan, error) {
	// Ensure required fields (ID and Version) are included
	req.Fields = util.EnsureIdAndVerField(req.Fields)

	// Initialize the update query
	updateBuilder := sq.Update(c.mainTable).
		PlaceholderFormat(sq.Dollar).
		Set("updated_at", req.CurrentTime()).
		Set("updated_by", req.Session.GetUserId()).
		Where(sq.Eq{"id": upd.Id, "dc": req.Session.GetDomainId()})

	// Increment version
	updateBuilder = updateBuilder.Set("ver", sq.Expr("ver + 1"))

	// Handle nested fields using switch-case on req.Mask
	for _, field := range req.Mask {
		switch field {
		case "subject":
			updateBuilder = updateBuilder.Set("subject", upd.Subject)
		case "description":
			updateBuilder = updateBuilder.Set("description", sq.Expr("NULLIF(?, '')", upd.Description))
		case "priority":
			updateBuilder = updateBuilder.Set("priority", upd.Priority.GetId())
		case "source":
			updateBuilder = updateBuilder.Set("source", upd.Source.GetId())
		case "status":
			updateBuilder = updateBuilder.Set("status", upd.Status.GetId())
		case "close.close_reason":
			if upd.Close != nil {
				updateBuilder = updateBuilder.Set("close_reason", upd.Close.CloseReason.GetId())
			}
		case "close.close_result":
			if upd.Close != nil {
				updateBuilder = updateBuilder.Set("close_result", upd.Close.CloseResult)
			}
		case "assignee":
			updateBuilder = updateBuilder.Set("assignee", upd.Assignee.GetId())
		case "reporter":
			updateBuilder = updateBuilder.Set("reporter", upd.Reporter.GetId())
		case "impacted":
			updateBuilder = updateBuilder.Set("impacted", upd.Impacted.GetId())
		case "contact_group":
			updateBuilder = updateBuilder.Set("contact_group", upd.Group.GetId())
		case "planned_reaction_at":
			updateBuilder = updateBuilder.Set("planned_reaction_at", util.LocalTime(upd.PlannedReactionAt))
		case "planned_resolve_at":
			updateBuilder = updateBuilder.Set("planned_resolve_at", util.LocalTime(upd.PlannedResolveAt))
		case "rate.rating":
			if upd.Rate != nil {
				updateBuilder = updateBuilder.Set("rating", upd.Rate.Rating)
			}
		case "rate.rating_comment":
			if upd.Rate != nil {
				updateBuilder = updateBuilder.Set("rating_comment", sq.Expr("NULLIF(?, '')", upd.Rate.RatingComment))
			}
		default:
			// Optionally handle unknown fields
			return nil, nil, dberr.NewDBError("postgres.case.update.invalid_field", fmt.Sprintf("Unknown field: %s", field))
		}
	}

	// Define SELECT query for returning updated fields
	selectBuilder, plan, err := c.buildCaseSelectColumnsAndPlan(
		req.Session,
		sq.Select().PrefixExpr(sq.Expr("WITH "+caseLeft+" AS (?)", updateBuilder.Suffix("RETURNING *"))),
		req.Fields,
	)
	if err != nil {
		return nil, nil, dberr.NewDBError("postgres.case.update.select_query_build_error", err.Error())
	}

	selectBuilder = selectBuilder.From(caseLeft)

	return selectBuilder, plan, nil
}

// session required to get some columns
func (c *CaseStore) buildCaseSelectColumnsAndPlan(session *authmodel.Session,
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
		case "contact_group":
			base = base.Column(fmt.Sprintf(
				"(SELECT ROW(g.id, g.name)::text FROM contacts.group g WHERE g.id = %s.contact_group) AS contact_group", caseLeft))
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
				return scanner.ScanTimestamp(&caseItem.PlannedReactionAt)
			})
		case "planned_resolve_at":
			base = base.Column(store.Ident(caseLeft, "planned_resolve_at"))
			plan = append(plan, func(caseItem *_go.Case) any {
				return scanner.ScanTimestamp(&caseItem.PlannedResolveAt)
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
		case "service":
			base = base.Column(fmt.Sprintf(
				"(SELECT ROW(s.id, s.name)::text FROM cases.service_catalog s WHERE s.id = %s.service) AS service", caseLeft))
			plan = append(plan, func(caseItem *_go.Case) any {
				return scanner.ScanRowLookup(&caseItem.Service)
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
				WHERE sc.sla_id = c.sla) AS sla_conditions`)
			plan = append(plan, func(caseItem *_go.Case) any {
				return scanner.ScanJSONToStructList(&caseItem.SlaCondition)
			})
		case "comments":
			opts := &model.SearchOptions{Session: session, Page: 1, Size: 10, Fields: []string{
				"id",
				"ver",
				"text",
				"created_by",
				"updated_by",
				"updated_at",
			}}
			subquery, scanPlan, filtersApplied, dbErr := buildCommentsSelectAsSubquery(opts, caseLeft)
			if dbErr != nil {
				return base, nil, dbErr
			}
			base = AddSubqueryAsColumn(base, subquery, "comments", filtersApplied > 0)
			plan = append(plan, func(value *_go.Case) any {
				var items []*_go.CaseComment
				postProcessing := func() error {
					if len(items) == 0 {
						return nil
					}
					res := &_go.CaseCommentList{Items: items}
					if len(items) > int(opts.GetSize()) {
						res.Items = res.Items[:len(res.Items)-1]
						res.Next = true
					}
					res.Page = int64(opts.GetPage())
					value.Comments = res
					return nil
				}
				return scanner.GetCompositeTextScanFunction(scanPlan, &items, postProcessing)
			})
		case "links":
			opts := &model.SearchOptions{Page: 1, Size: 10, Fields: []string{
				"id",
				"ver",
				"url",
				"name",
				"author",
				"created_by",
			}}
			subquery, scanPlan, filtersApplied, dbErr := buildLinkSelectAsSubquery(opts, caseLeft)
			if dbErr != nil {
				return base, nil, dbErr
			}
			base = AddSubqueryAsColumn(base, subquery, "links", filtersApplied > 0)
			plan = append(plan, func(value *_go.Case) any {
				var items []*_go.CaseLink
				postProcessing := func() error {
					if len(items) == 0 {
						return nil
					}
					res := &_go.CaseLinkList{Items: items}
					if len(items) > int(opts.GetSize()) {
						res.Items = res.Items[:len(res.Items)-1]
						res.Next = true
					}
					res.Page = int64(opts.GetPage())
					value.Links = res
					return nil
				}
				return scanner.GetCompositeTextScanFunction(scanPlan, &items, postProcessing)
			})
		case "files":
			opts := &model.SearchOptions{Page: 1, Size: 10, Fields: []string{
				"id",
				"size",
				"mime",
				"name",
				"created_at",
			}}
			subquery, scanPlan, dbErr := buildFilesSelectAsSubquery(opts, caseLeft)
			if dbErr != nil {
				return base, nil, dbErr
			}
			var filtersApplied bool
			if len(opts.Filter) > 0 {
				filtersApplied = true
			}
			base = AddSubqueryAsColumn(base, subquery, "files", filtersApplied)
			plan = append(plan, func(value *_go.Case) any {
				var items []*_go.File
				postProcessing := func() error {
					if len(items) == 0 {
						return nil
					}
					res := &_go.CaseFileList{Items: items}
					if len(items) > int(opts.GetSize()) {
						res.Items = res.Items[:len(res.Items)-1]
						res.Next = true
					}
					res.Page = int64(opts.GetPage())
					value.Files = res
					return nil
				}
				return scanner.GetCompositeTextScanFunction(scanPlan, &items, postProcessing)
			})
		case "related_cases":
			base = base.Column(fmt.Sprintf(`
				(SELECT JSON_AGG(JSON_BUILD_OBJECT(
					'id', rc.id, -- ID of the related_case record
					'child', JSON_BUILD_OBJECT( -- Child case details
						'id', c_child.id,
						'name', c_child.name,
						'subject', c_child.subject,
						'description', c_child.description
					),
					'created_at', CAST(EXTRACT(EPOCH FROM rc.created_at) * 1000 AS BIGINT),
					'created_by', JSON_BUILD_OBJECT(
					   'name', u.name,
					   'id', u.id
					),
					'relation_type', rc.relation_type  -- Output numeric enum value directly
				))
				FROM %s rc
                JOIN cases.case c_child
                ON rc.child_case_id = c_child.id -- Fetch details for the child case
				LEFT JOIN directory.wbt_user u ON rc.created_by = u.id
                WHERE rc.parent_case_id = %s.id) AS related_cases`, relatedAlias, caseLeft))
			// parent_case_id -- newly created case
			// child_case_id -- attached case id
			plan = append(plan, func(caseItem *_go.Case) any {
				if caseItem.Related == nil {
					caseItem.Related = &_go.RelatedCaseList{}
				}
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

func AddSubqueryAsColumn(mainQuery sq.SelectBuilder, subquery sq.SelectBuilder, subAlias string, filtersApplied bool) sq.SelectBuilder {
	if filtersApplied {
		subquery = subquery.Prefix("LATERAL (SELECT ARRAY(SELECT (sub) FROM (").Suffix(fmt.Sprintf(") sub) %s) %[1]s ON array_length(%[1]s.%[1]s, 1) > 0", subAlias))
		query, args, _ := subquery.ToSql()
		mainQuery = mainQuery.Join(query, args...)
	} else {
		subquery = subquery.Prefix("LATERAL (SELECT ARRAY(SELECT (sub) FROM (").Suffix(fmt.Sprintf(") sub) %s) %[1]s ON true", subAlias))
		query, args, _ := subquery.ToSql()
		mainQuery = mainQuery.LeftJoin(query, args...)
	}
	mainQuery = mainQuery.Column(subAlias + "::text")

	return mainQuery
}

// Helper function to convert the scan plan to arguments for scanning.
func convertToCaseScanArgs(plan []CaseScan, caseItem *_go.Case) []any {
	var scanArgs []any
	for _, scan := range plan {
		scanArgs = append(scanArgs, scan(caseItem))
	}
	return scanArgs
}

func NewCaseStore(store store.Store) (store.CaseStore, error) {
	if store == nil {
		return nil, dberr.NewDBError("postgres.new_case.check.bad_arguments",
			"error creating case interface to the case table, main store is nil")
	}
	return &CaseStore{storage: store, mainTable: "cases.case"}, nil
}
