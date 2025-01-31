package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v5"
	"github.com/lib/pq"
	_go "github.com/webitel/cases/api/cases"
	"github.com/webitel/cases/auth"
	dberr "github.com/webitel/cases/internal/errors"
	"github.com/webitel/cases/internal/store"
	"github.com/webitel/cases/internal/store/scanner"
	"github.com/webitel/cases/model"
	util "github.com/webitel/cases/util"
)

type CaseStore struct {
	storage   store.Store
	mainTable string
}

const (
	caseLeft               = "c"
	relatedAlias           = "related"
	linksAlias             = "links"
	caseDefaultSort        = "created_at"
	casesObjClassScopeName = "cases"
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
) (sq.SelectBuilder, []func(caseItem *_go.Case) any, error) {
	// Parameters for the main case and nested JSON arrays
	var (
		reporter    *int64
		assignee    *int64
		closeReason *int64
		closeResult *string
	)
	if cl := caseItem.GetClose(); cl != nil {
		if cl.CloseReason != nil && cl.CloseReason.GetId() > 0 {
			closeReason = &cl.CloseReason.Id
		}
		if cl.CloseResult != "" {
			closeResult = &cl.CloseResult
		}
	}
	if caseItem.Reporter != nil && caseItem.Reporter.GetId() > 0 {
		reporter = &caseItem.Reporter.Id
	}
	if caseItem.Assignee != nil && caseItem.Assignee.GetId() > 0 {
		assignee = &caseItem.Assignee.Id
	}
	params := map[string]interface{}{
		// Case-level parameters
		"date":                rpc.CurrentTime(),
		"contact_info":        caseItem.GetContactInfo(),
		"user":                rpc.GetAuthOpts().GetUserId(),
		"dc":                  rpc.GetAuthOpts().GetDomainId(),
		"sla":                 sla,
		"sla_condition":       slaCondition,
		"status":              caseItem.Status.GetId(),
		"service":             caseItem.Service.GetId(),
		"priority":            caseItem.Priority.GetId(),
		"source":              caseItem.Source.GetId(),
		"contact_group":       caseItem.Group.GetId(),
		"close_reason_group":  caseItem.CloseReasonGroup.GetId(),
		"close_result":        closeResult,
		"close_reason":        closeReason,
		"subject":             caseItem.Subject,
		"planned_reaction_at": util.LocalTime(caseItem.PlannedReactionAt),
		"planned_resolve_at":  util.LocalTime(caseItem.PlannedResolveAt),
		"reporter":            reporter,
		"impacted":            caseItem.Impacted.GetId(),
		"description":         caseItem.Description,
		"assignee":            assignee,
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
	    service_cte AS(
		SELECT catalog_id
		FROM cases.service_catalog
			WHERE id = :service
			LIMIT 1
		),
		prefix_cte AS (
			SELECT prefix
			FROM cases.service_catalog
			WHERE id = any(SELECT catalog_id FROM service_cte)
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
				id, name, dc, created_at, created_by, updated_at, updated_by,
				priority, source, status, contact_group, close_reason_group,
				subject, planned_reaction_at, planned_resolve_at, reporter, impacted,
				service, description, assignee, sla, sla_condition_id, status_condition, contact_info,
				close_result, close_reason
			) VALUES (
				(SELECT id FROM id_cte),
				CONCAT((SELECT prefix FROM prefix_cte), '_', (SELECT id FROM id_cte)),
				:dc, :date, :user, :date, :user,
				:priority, :source, :status, :contact_group, :close_reason_group,
				:subject, :planned_reaction_at, :planned_resolve_at, :reporter, :impacted,
				:service, :description, :assignee, :sla, :sla_condition,
				(SELECT status_condition_id FROM status_condition_cte), :contact_info, :close_result, :close_reason
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
				primary_case_id, related_case_id, relation_type, dc, created_by, created_at, updated_by, updated_at
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
		&model.SearchOptions{Context: rpc.Context, Fields: rpc.Fields},
		sq.Select().PrefixExpr(sq.Expr(boundQuery, args...)),
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
	if related == nil || len(related.Data) == 0 {
		return []byte("[]")
	}
	var jsonArray []map[string]interface{}
	for _, item := range related.Data {
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
	case _go.RelationType_RELATION_TYPE_UNSPECIFIED:
		return 0, nil
	case _go.RelationType_DUPLICATES:
		return 1, nil
	case _go.RelationType_IS_DUPLICATED_BY:
		return 2, nil
	case _go.RelationType_BLOCKS:
		return 3, nil
	case _go.RelationType_IS_BLOCKED_BY:
		return 4, nil
	case _go.RelationType_CAUSES:
		return 5, nil
	case _go.RelationType_IS_CAUSED_BY:
		return 6, nil
	case _go.RelationType_IS_CHILD_OF:
		return 7, nil
	case _go.RelationType_IS_PARENT_OF:
		return 8, nil
	case _go.RelationType_RELATES_TO:
		return 9, nil
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
		FROM flow.calendar cl,  UNNEST(accepts::flow.calendar_accept_time[]) x
		WHERE cl.id = $1
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

	var offset time.Duration
	offsetRow := txManager.QueryRow(rpc.Context, `
		SELECT tz.utc_offset
		FROM flow.calendar cl
		    LEFT JOIN flow.calendar_timezones tz ON tz.id = cl.timezone_id
		WHERE cl.id = $1`, calendarID)
	err = offsetRow.Scan(&offset)
	if err != nil {
		return fmt.Errorf("failed to fetch calendar offset details: %w", err)
	}

	// Convert reaction and resolution times from milliseconds to minutes
	reactionMinutes := reactionTime / 60
	resolutionMinutes := resolutionTime / 60

	currentTime := rpc.CurrentTime()
	reactionTimestamp, err := calculateTimestampFromCalendar(currentTime, offset, reactionMinutes, calendar)
	if err != nil {
		return fmt.Errorf("failed to calculate planned reaction time: %w", err)
	}

	//?? TODO
	// resolveTimestamp, err := calculateTimestampFromCalendar(reactionTimestamp, resolutionMinutes, calendar)
	resolveTimestamp, err := calculateTimestampFromCalendar(currentTime, offset, resolutionMinutes, calendar)
	if err != nil {
		return fmt.Errorf("failed to calculate planned resolution time: %w", err)
	}

	caseItem.PlannedReactionAt = reactionTimestamp.UnixMilli()
	caseItem.PlannedResolveAt = resolveTimestamp.UnixMilli()

	return nil
}

func calculateTimestampFromCalendar(
	startTime time.Time,
	calendarOffset time.Duration,
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
	var (
		startingAtMinutes int
		startDayProcessed bool
		addDays           int
	)

	for {
		for _, slot := range calendar {
			slotStartOfTheDayUtc := int((time.Duration(slot.StartTimeOfDay)*time.Minute - calendarOffset).Minutes())
			slotEndOfTheDayUtc := int((time.Duration(slot.EndTimeOfDay)*time.Minute - calendarOffset).Minutes())
			// Match the current day and ensure the slot is in the future
			if slot.Day == currentDay && slotEndOfTheDayUtc > currentTimeInMinutes {
				if startDayProcessed {
					currentTimeInMinutes = slotStartOfTheDayUtc
				}
				if currentTimeInMinutes < slotStartOfTheDayUtc {
					startingAtMinutes = slotStartOfTheDayUtc
				} else {
					startingAtMinutes = currentTimeInMinutes
				}
				availableMinutes := slotEndOfTheDayUtc - startingAtMinutes
				if availableMinutes >= remainingMinutes {
					// if we need to add days, then we should add minutes from the start of the working day
					if addDays != 0 {
						startDate := time.Date(startTime.Year(), startTime.Month(), startTime.Day(), 0, 0, 0, 0, startTime.Location())
						// Add days required
						startDate = startDate.Add(time.Duration(24*addDays) * time.Hour)
						// Add minutes to the start of the working day
						startDate = startDate.Add(time.Duration(slotStartOfTheDayUtc) * time.Minute)
						// Add remaining minutes to the start of the working day
						startDate = startDate.Add(time.Duration(remainingMinutes) * time.Minute)
						return startDate, nil
					}
					// Calculate the exact timestamp
					return startTime.Add(time.Duration(remainingMinutes) * time.Minute), nil
				}
				if slot.Day == currentDay && !startDayProcessed {
					startDayProcessed = true
				}
				remainingMinutes -= availableMinutes
			}

		}

		// If no slots available, move to the next day
		currentDay = (currentDay + 1) % len(calendar) // Wrap around to the start of the week if necessary
		addDays++
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
	var err error
	convertedIds := util.Int64SliceToStringSlice(rpc.IDs)
	ids := util.FieldsFunc(convertedIds, util.InlineFields)
	query := sq.Delete("cases.case").Where("id = ANY(?)", ids).Where("dc = ?", rpc.GetAuthOpts().GetDomainId()).PlaceholderFormat(sq.Dollar)
	query, err = addCaseRbacConditionForDelete(rpc.GetAuthOpts(), auth.Delete, query, "case.id")
	if err != nil {
		return "", nil, err
	}

	return query.ToSql()
}

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
	slct = store.CompactSQL(slct)

	db, dbErr := c.storage.Database()
	if dbErr != nil {
		return nil, dbErr
	}
	rows, err := db.Query(opts.Context, store.CompactSQL(slct), args...)
	if err != nil {
		return nil, dberr.NewDBError("postgres.case.list.exec.error", err.Error())
	}
	var res _go.CaseList
	res.Items, err = c.scanCases(rows, plan)
	if err != nil {
		return nil, err
	}
	res.Items, res.Next = store.ResolvePaging(opts.GetSize(), res.Items)
	res.Page = int64(opts.GetPage())
	return &res, nil
}

func (c *CaseStore) CheckRbacAccess(ctx context.Context, auth auth.Auther, access auth.AccessMode, caseId int64) (bool, error) {
	if auth == nil {
		return false, nil
	}
	if !auth.GetObjectScope(casesObjClassScopeName).IsRbacUsed() {
		return true, nil
	}
	q := sq.Select("1").From("cases.case_acl acl").
		Where("acl.dc = ?", auth.GetDomainId()).
		Where("acl.object = ?", caseId).
		Where("acl.subject = any( ?::int[])", pq.Array(auth.GetRoles())).
		Where("acl.access & ? = ?", int64(access), int64(access)).
		Limit(1)
	db, err := c.storage.Database()
	if err != nil {
		return false, err
	}
	sql, args, defErr := q.ToSql()
	if defErr != nil {
		return false, defErr
	}
	res, defErr := db.Exec(ctx, sql, args...)
	if defErr != nil {
		return false, defErr
	}
	if res.RowsAffected() == 1 {
		return true, nil
	}
	return false, nil
}

func (c *CaseStore) buildListCaseSqlizer(opts *model.SearchOptions) (sq.SelectBuilder, []func(caseItem *_go.Case) any, error) {
	base := sq.Select().From(fmt.Sprintf("%s %s", c.mainTable, caseLeft)).PlaceholderFormat(sq.Dollar)
	base, plan, err := c.buildCaseSelectColumnsAndPlan(opts, base)
	if search := opts.Search; search != "" {
		searchTerm, operator := store.ParseSearchTerm(search)
		searchNumber := store.PrepareSearchNumber(search)
		where := sq.Or{
			sq.Expr(fmt.Sprintf(`%s.reporter = ANY (SELECT contact_id
                        FROM contacts.contact_phone ct_ph
                        WHERE ct_ph.reverse like
						'%%' ||
							overlay(? placing '' from coalesce(
								(select value::int from call_center.system_settings s where s.domain_id = ? and s.name = 'search_number_length'),
								 ?)+1 for ?)
						 || '%%' )`, caseLeft),
				searchNumber, opts.GetAuthOpts().GetDomainId(), len(searchNumber), len(searchNumber)),
			sq.Expr(fmt.Sprintf(`%s.reporter = ANY (SELECT contact_id
                        FROM contacts.contact_email ct_em
                        WHERE ct_em.email %s ?)`, caseLeft, operator),
				searchTerm),
			sq.Expr(fmt.Sprintf(`%s.reporter = ANY (SELECT contact_id
                        FROM contacts.contact_imclient ct_im
                        WHERE ct_im.user_id IN (SELECT id FROM chat.client WHERE name %s ?))`, caseLeft, operator),
				searchTerm),
			sq.Expr(fmt.Sprintf("%s %s ?", store.Ident(caseLeft, "subject"), operator), searchTerm),
			sq.Expr(fmt.Sprintf("%s %s ?", store.Ident(caseLeft, "name"), operator), searchTerm),
			// sq.Expr(fmt.Sprintf("%s = ?", store.Ident(caseLeft, "contact_info")), search),
		}
		base = base.Where(where)

	}
	for _, d := range opts.IDs {
		base = base.Where(fmt.Sprintf("%s = ?", store.Ident(caseLeft, "id")), d)
	}
	for column, value := range opts.Filter {
		switch column {
		case "created_by",
			"updated_by",
			"assignee",         // +
			"reporter",         // +
			"source",           // +
			"priority",         // +
			"impacted",         // +
			"author",           // +
			"close_reason",     // +
			"contact_group",    // +
			"service",          // +
			"status_condition", // +
			"sla":              // +
			if value == "" {
				base = base.Where(fmt.Sprintf("%s ISNULL", store.Ident(caseLeft, column)))
				continue
			}
			switch typedValue := value.(type) {
			case string:
				values := strings.Split(typedValue, ",")
				var valuesInt []int64
				for _, s := range values {
					converted, err := strconv.ParseInt(s, 10, 64)
					if err != nil {
						return base, nil, dberr.NewDBInternalError("postgres.case.build_list_case_sqlizer.convert_to_int_array.error", err)
					}
					valuesInt = append(valuesInt, converted)
				}
				base = base.Where(fmt.Sprintf("%s =  ANY(?::int[])", store.Ident(caseLeft, column)), valuesInt)
			}

		case "rating.from":
			cutted, _ := strings.CutSuffix(column, ".from")
			base = base.Where(fmt.Sprintf("%s > ?::INT", store.Ident(caseLeft, cutted)), value)
		case "rating.to":
			cutted, _ := strings.CutSuffix(column, ".to")
			base = base.Where(fmt.Sprintf("%s < ?::INT", store.Ident(caseLeft, cutted)), value)
		case "sla_condition":
			base = base.Where(fmt.Sprintf("? = ANY(%s)", store.Ident(caseLeft, column)), value)
		case "reacted_at.from", "resolved_at.from", "planned_reaction_at.from", "planned_resolved_at.from":
			cutted, _ := strings.CutSuffix(column, ".from")
			base = base.Where(fmt.Sprintf("extract(epoch from %s)::INT > ?::INT", store.Ident(caseLeft, cutted)), value)
		case "reacted_at.to", "resolved_at.to", "planned_reaction_at.to", "planned_resolved_at.to":
			cutted, _ := strings.CutSuffix(column, ".to")
			base = base.Where(fmt.Sprintf("extract(epoch from %s)::INT < ?::INT", store.Ident(caseLeft, cutted)), value)
		case "attachments":
			var operator string
			if value != "true" {
				operator = "NOT "
			}
			base = base.Where(sq.Expr(fmt.Sprintf(operator+"EXISTS (SELECT id FROM storage.files WHERE uuid = %s::varchar UNION SELECT id FROM cases.case_link WHERE case_link.case_id = %[1]s)", store.Ident(caseLeft, "id"))))

		}
	}
	if err != nil {
		return base, nil, err
	}

	if sess := opts.GetAuthOpts(); sess != nil {
		base = base.Where(store.Ident(caseLeft, "dc = ?"), opts.GetAuthOpts().GetDomainId())
		base, err = addCaseRbacCondition(sess, auth.Read, base, store.Ident(caseLeft, "id"))
	}
	// pagination
	base = store.ApplyPaging(opts.GetPage(), opts.GetSize(), base)
	// sort
	base = store.ApplyDefaultSorting(opts, base, caseDefaultSort)

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

	if err := db.QueryRow(rpc.Context, query, args...).Scan(scanArgs...); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, dberr.NewDBNoRowsError("postgres.case.update.update.scan_ver.not_found")
		}
		return nil, dberr.NewDBInternalError("postgres.case.update.update.execution_error", err)
	}

	return upd, nil
}

func (c *CaseStore) buildUpdateCaseSqlizer(
	rpc *model.UpdateOptions,
	upd *_go.Case,
) (sq.Sqlizer, []func(caseItem *_go.Case) any, error) {
	// Ensure required fields (ID and Version) are included
	rpc.Fields = util.EnsureIdAndVerField(rpc.Fields)
	var err error

	// Initialize the update query
	updateBuilder := sq.Update(c.mainTable).
		PlaceholderFormat(sq.Dollar).
		Set("updated_at", rpc.CurrentTime()).
		Set("updated_by", rpc.GetAuthOpts().GetUserId()).
		Where(sq.Eq{
			"id":  rpc.Etags[0].GetOid(),
			"ver": rpc.Etags[0].GetVer(),
			"dc":  rpc.GetAuthOpts().GetDomainId(),
		})

	updateBuilder, err = addCaseRbacConditionForUpdate(rpc.GetAuthOpts(), auth.Edit, updateBuilder, "case.id")
	if err != nil {
		return nil, nil, err
	}
	// Increment version
	updateBuilder = updateBuilder.Set("ver", sq.Expr("ver + 1"))

	// Handle nested fields using switch-case on req.Mask
	for _, field := range rpc.Mask {
		switch field {
		case "subject":
			updateBuilder = updateBuilder.Set("subject", upd.GetSubject())
		case "description":
			updateBuilder = updateBuilder.Set("description", sq.Expr("NULLIF(?, '')", upd.Description))
		case "priority":
			updateBuilder = updateBuilder.Set("priority", upd.Priority.GetId())
		case "source":
			updateBuilder = updateBuilder.Set("source", upd.Source.GetId())
		case "status":
			updateBuilder = updateBuilder.Set("status", upd.Status.GetId())
		case "status_condition":
			updateBuilder = updateBuilder.Set("status_condition", upd.StatusCondition.GetId())
		case "service":
			updateBuilder = updateBuilder.Set("service", upd.Service.GetId())
			updateBuilder = updateBuilder.Set("sla", sq.Expr("(SELECT sla_id FROM cases.service_catalog WHERE id = ? LIMIT 1)", upd.Service.GetId()))
		case "assignee":
			updateBuilder = updateBuilder.Set("assignee", upd.Assignee.GetId())
		case "reporter":
			updateBuilder = updateBuilder.Set("reporter", upd.Reporter.GetId())
		case "contact_info":
			updateBuilder = updateBuilder.Set("contact_info", upd.GetContactInfo())
		case "impacted":
			updateBuilder = updateBuilder.Set("impacted", upd.Impacted.GetId())
		case "group":
			updateBuilder = updateBuilder.Set("contact_group", upd.Group.GetId())
		case "close.close_reason":
			if upd.Close != nil {
				updateBuilder = updateBuilder.Set("close_reason", upd.Close.CloseReason.GetId())
			}
		case "close.close_result":
			if upd.Close != nil {
				updateBuilder = updateBuilder.Set("close_result", upd.Close.GetCloseResult())
			}
		case "rate.rating":
			if upd.Rate != nil {
				updateBuilder = updateBuilder.Set("rating", upd.Rate.Rating)
			}
		case "rate.rating_comment":
			if upd.Rate != nil {
				updateBuilder = updateBuilder.Set("rating_comment", sq.Expr("NULLIF(?, '')", upd.Rate.RatingComment))
			}
		}
	}

	// Define SELECT query for returning updated fields
	selectBuilder, plan, err := c.buildCaseSelectColumnsAndPlan(
		&model.SearchOptions{
			Size:   -1,
			Fields: rpc.Fields,
			Auth:   rpc.GetAuthOpts(),
			Time:   rpc.Time,
		},
		sq.Select().PrefixExpr(sq.Expr("WITH "+caseLeft+" AS (?)", updateBuilder.Suffix("RETURNING *"))),
	)
	if err != nil {
		return nil, nil, dberr.NewDBError("postgres.case.update.select_query_build_error", err.Error())
	}

	selectBuilder = selectBuilder.From(caseLeft)

	return selectBuilder, plan, nil
}

// session required to get some columns
func (c *CaseStore) buildCaseSelectColumnsAndPlan(opts *model.SearchOptions,
	base sq.SelectBuilder,
) (sq.SelectBuilder, []func(caseItem *_go.Case) any, error) {
	var plan []func(caseItem *_go.Case) any

	for _, field := range opts.Fields {
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
				`(
					SELECT
						ROW(g.id, g.name,
							CASE
								WHEN g.id IN (SELECT id FROM contacts.dynamic_group) THEN 'dynamic'
								ELSE 'static'
							END
						)::text
					FROM contacts.group g
					WHERE g.id = %s.contact_group
				) AS contact_group`, caseLeft))
			plan = append(plan, func(caseItem *_go.Case) any {
				return scanner.ScanRowExtendedLookup(&caseItem.Group)
			})
		case "source":
			base = base.Column(fmt.Sprintf(
				"(SELECT ROW(src.id, src.name, src.type)::text FROM cases.source src WHERE src.id = %s.source) AS source", caseLeft))
			plan = append(plan, func(caseItem *_go.Case) any {
				return scanner.TextDecoder(func(src []byte) error {
					if len(src) == 0 {
						return nil // NULL
					}
					// pointer on pointer on source
					if caseItem.Source == nil {
						caseItem.Source = new(_go.SourceTypeLookup)
					}

					var (
						ok  bool
						str pgtype.Text
						row = []pgtype.TextDecoder{
							scanner.TextDecoder(func(src []byte) error {
								if len(src) == 0 {
									return nil
								}
								err := str.DecodeText(nil, src)
								if err != nil {
									return err
								}
								id, err := strconv.ParseInt(str.String, 10, 64)
								if err != nil {
									return err
								}
								caseItem.Source.Id = id
								return nil
							}),
							scanner.TextDecoder(func(src []byte) error {
								if len(src) == 0 {
									return nil
								}
								err := str.DecodeText(nil, src)
								if err != nil {
									return err
								}
								caseItem.Source.Name = str.String
								ok = ok || (str.String != "" && str.String != "[deleted]") // && str.Status == pgtype.Present
								return nil
							}),
							scanner.TextDecoder(func(src []byte) error {
								if len(src) == 0 {
									return nil
								}
								err := str.DecodeText(nil, src)
								if err != nil {
									return err
								}
								for i, text := range _go.SourceType_name {
									if text == str.String {
										caseItem.Source.Type = _go.SourceType(i)
										return nil
									}
								}
								caseItem.Source.Type = _go.SourceType_TYPE_UNSPECIFIED
								return nil
							}),
						}
						raw = pgtype.NewCompositeTextScanner(nil, src)
					)

					var err error
					for _, col := range row {

						raw.ScanDecoder(col)

						err = raw.Err()
						if err != nil {
							return err
						}
					}

					return nil
				})
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

		case "close":
			base = base.Column(store.Ident(caseLeft, "close_result"))
			plan = append(plan, func(caseItem *_go.Case) any {
				if caseItem.Close == nil {
					caseItem.Close = &_go.CloseInfo{}
				}
				return scanner.ScanText(&caseItem.Close.CloseResult)
			})
			base = base.Column(fmt.Sprintf(
				"(SELECT ROW(cr.id, cr.name)::text FROM cases.close_reason cr WHERE cr.id = %s.close_reason) AS close_reason", caseLeft))
			plan = append(plan, func(caseItem *_go.Case) any {
				if caseItem.Close == nil {
					caseItem.Close = &_go.CloseInfo{}
				}
				return scanner.ScanRowLookup(&caseItem.Close.CloseReason)
			})
		case "rate":
			base = base.Column(store.Ident(caseLeft, "rating"))
			plan = append(plan, func(caseItem *_go.Case) any {
				if caseItem.Rate == nil {
					caseItem.Rate = &_go.RateInfo{}
				}
				return scanner.ScanInt64(&caseItem.Rate.Rating)
			})
			base = base.Column(store.Ident(caseLeft, "rating_comment"))
			plan = append(plan, func(caseItem *_go.Case) any {
				if caseItem.Rate == nil {
					caseItem.Rate = &_go.RateInfo{}
				}
				return scanner.ScanText(&caseItem.Rate.RatingComment)
			})
		case "timing":
			base = base.
				Column(fmt.Sprintf("COALESCE(%s.resolved_at, '1970-01-01 00:00:00') AS resolved_at", caseLeft)).
				Column(fmt.Sprintf("COALESCE(%s.reacted_at, '1970-01-01 00:00:00') AS reacted_at", caseLeft)).
				Column(fmt.Sprintf(
					"COALESCE(CAST(EXTRACT(EPOCH FROM %[1]s.reacted_at - %[1]s.created_at) * 1000 AS bigint), 0) AS difference_in_reaction",
					caseLeft,
				)).
				Column(fmt.Sprintf(
					"COALESCE(CAST(EXTRACT(EPOCH FROM %[1]s.resolved_at - %[1]s.created_at) * 1000 AS bigint), 0) AS difference_in_resolve",
					caseLeft,
				))

			plan = append(plan, func(caseItem *_go.Case) any {
				if caseItem.Timing == nil {
					caseItem.Timing = &_go.TimingInfo{}
				}
				return scanner.ScanTimestamp(&caseItem.Timing.ResolvedAt)
			})
			plan = append(plan, func(caseItem *_go.Case) any {
				return scanner.ScanTimestamp(&caseItem.Timing.ReactedAt)
			})
			plan = append(plan, func(caseItem *_go.Case) any {
				return scanner.ScanInt64(&caseItem.Timing.DifferenceInReaction)
			})
			plan = append(plan, func(caseItem *_go.Case) any {
				return scanner.ScanInt64(&caseItem.Timing.DifferenceInResolve)
			})
		case "sla":
			base = base.Column(fmt.Sprintf(
				"(SELECT ROW(sla.id, sla.name)::text FROM cases.sla sla WHERE sla.id = %s.sla) AS sla", caseLeft))
			plan = append(plan, func(caseItem *_go.Case) any {
				return scanner.ScanRowLookup(&caseItem.Sla)
			})
		case "status_condition":
			base = base.Column(fmt.Sprintf(`
				(SELECT ROW(stc.id, stc.name, stc.initial, stc.final)::text
				 FROM cases.status_condition stc
				 WHERE stc.id = %s.status_condition) AS status_condition`, caseLeft))
			plan = append(plan, func(caseItem *_go.Case) any {
				return scanner.TextDecoder(func(src []byte) error {
					if len(src) == 0 {
						return nil // NULL
					}
					// pointer on pointer on source
					if caseItem.StatusCondition == nil {
						caseItem.StatusCondition = new(_go.StatusCondition)
					}

					var (
						str pgtype.Text
						bl  pgtype.Bool
						row = []pgtype.TextDecoder{
							scanner.TextDecoder(func(src []byte) error {
								if len(src) == 0 {
									return nil
								}
								err := str.DecodeText(nil, src)
								if err != nil {
									return err
								}
								id, err := strconv.ParseInt(str.String, 10, 64)
								if err != nil {
									return err
								}
								caseItem.StatusCondition.Id = id
								return nil
							}),
							scanner.TextDecoder(func(src []byte) error {
								if len(src) == 0 {
									return nil
								}
								err := str.DecodeText(nil, src)
								if err != nil {
									return err
								}
								caseItem.StatusCondition.Name = str.String
								return nil
							}),
							scanner.TextDecoder(func(src []byte) error {
								if len(src) == 0 {
									return nil
								}
								err := bl.Scan(src)
								if err != nil {
									return err
								}
								caseItem.StatusCondition.Initial = bl.Bool
								return nil
							}),
							scanner.TextDecoder(func(src []byte) error {
								if len(src) == 0 {
									return nil
								}
								err := bl.Scan(src)
								if err != nil {
									return err
								}
								caseItem.StatusCondition.Final = bl.Bool
								return nil
							}),
						}
						raw = pgtype.NewCompositeTextScanner(nil, src)
					)

					var err error
					for _, col := range row {

						raw.ScanDecoder(col)

						err = raw.Err()
						if err != nil {
							return err
						}
					}

					return nil
				})
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
				"(SELECT ROW(p.id, p.name, p.color)::text FROM cases.priority p WHERE p.id = %s.priority) AS priority", caseLeft))
			plan = append(plan, func(caseItem *_go.Case) any {
				return scanner.TextDecoder(func(src []byte) error {
					if len(src) == 0 {
						return nil // NULL
					}
					// pointer on pointer on source
					if caseItem.Priority == nil {
						caseItem.Priority = new(_go.Priority)
					}

					var (
						str pgtype.Text
						row = []pgtype.TextDecoder{
							scanner.TextDecoder(func(src []byte) error {
								if len(src) == 0 {
									return nil
								}
								err := str.DecodeText(nil, src)
								if err != nil {
									return err
								}
								id, err := strconv.ParseInt(str.String, 10, 64)
								if err != nil {
									return err
								}
								caseItem.Priority.Id = id
								return nil
							}),
							scanner.TextDecoder(func(src []byte) error {
								if len(src) == 0 {
									return nil
								}
								err := str.DecodeText(nil, src)
								if err != nil {
									return err
								}
								caseItem.Priority.Name = str.String
								return nil
							}),
							scanner.TextDecoder(func(src []byte) error {
								if len(src) == 0 {
									return nil
								}
								err := str.DecodeText(nil, src)
								if err != nil {
									return err
								}
								caseItem.Priority.Color = str.String
								return nil
							}),
						}
						raw = pgtype.NewCompositeTextScanner(nil, src)
					)

					var err error
					for _, col := range row {

						raw.ScanDecoder(col)

						err = raw.Err()
						if err != nil {
							return err
						}
					}

					return nil
				})
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
		case "contact_info":
			base = base.Column(store.Ident(caseLeft, field))
			plan = append(plan, func(caseItem *_go.Case) any {
				return scanner.ScanText(&caseItem.ContactInfo)
			})
		case "impacted":
			base = base.Column(fmt.Sprintf(
				"(SELECT ROW(i.id, i.common_name)::text FROM contacts.contact i WHERE i.id = %s.impacted) AS impacted", caseLeft))
			plan = append(plan, func(caseItem *_go.Case) any {
				return scanner.ScanRowLookup(&caseItem.Impacted)
			})
		case "sla_condition":
			base = base.Column(`
				(SELECT ROW(sc.id, sc.name)::text
				FROM cases.sla_condition sc
				WHERE sc.sla_id = c.sla AND sc.id = ANY(SELECT sla_condition_id FROM cases.priority_sla_condition WHERE priority_id = c.priority LIMIT 1)) AS sla_condition`)
			plan = append(plan, func(caseItem *_go.Case) any {
				return scanner.ScanRowLookup(&caseItem.SlaCondition)
			})
		case "comments":
			derivedOpts := opts.SearchDerivedOptionByField(field)
			if derivedOpts == nil {
				// default
				derivedOpts = &model.SearchOptions{
					Context: opts.Context,
					Fields:  []string{"id", "ver", "text", "created_by", "author", "created_at", "can_edit"},
					Size:    10,
					Page:    1,
					Sort:    []string{"-created_at"},
					Auth:    opts.Auth,
				}
			}

			subquery, scanPlan, filtersApplied, dbErr := buildCommentsSelectAsSubquery(derivedOpts, caseLeft)
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
					res := &_go.CaseCommentList{}
					res.Items, res.Next = store.ResolvePaging(opts.GetSize(), items)
					res.Page = int64(opts.GetPage())
					value.Comments = res
					return nil
				}
				return scanner.GetCompositeTextScanFunction(scanPlan, &items, postProcessing)
			})
		case "links":
			derivedOpts := opts.SearchDerivedOptionByField(field)
			if derivedOpts == nil {
				// default
				derivedOpts = &model.SearchOptions{
					Context: opts.Context,
					Fields:  []string{"id", "ver", "name", "url", "created_by", "author", "created_at"},
					Size:    10,
					Page:    1,
					Sort:    []string{"-created_at"},
				}
			}
			subquery, scanPlan, filtersApplied, dbErr := buildLinkSelectAsSubquery(derivedOpts, caseLeft)
			if dbErr != nil {
				return base, nil, dbErr
			}
			base = AddSubqueryAsColumn(base, subquery, field, filtersApplied > 0)
			plan = append(plan, func(value *_go.Case) any {
				var items []*_go.CaseLink
				postProcessing := func() error {
					if len(items) == 0 {
						return nil
					}
					res := &_go.CaseLinkList{}
					res.Items, res.Next = store.ResolvePaging(opts.GetSize(), items)
					res.Page = int64(opts.GetPage())
					value.Links = res
					return nil
				}
				return scanner.GetCompositeTextScanFunction(scanPlan, &items, postProcessing)
			})
		case "files":
			derivedOpts := opts.SearchDerivedOptionByField(field)
			if derivedOpts == nil {
				derivedOpts = &model.SearchOptions{
					Page: 1,
					Size: 10,
					Sort: []string{"-created_at"},
					Fields: []string{
						"id",
						"size",
						"mime",
						"name",
						"created_at",
					},
				}
			}
			subquery, scanPlan, filtersApplied, dbErr := buildFilesSelectAsSubquery(derivedOpts, caseLeft)
			if dbErr != nil {
				return base, nil, dbErr
			}
			base = AddSubqueryAsColumn(base, subquery, field, filtersApplied > 0)
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
		case "related":
			subquery, err := buildRelatedCasesSubquery(caseLeft)
			if err != nil {
				return base, nil, err
			}

			sqlStr, _, sqlErr := subquery.ToSql()
			if sqlErr != nil {
				return base, nil, sqlErr
			}

			// Add the subquery as a column
			base = base.Column(fmt.Sprintf("(%s) AS related_cases", sqlStr))

			plan = append(plan, func(caseItem *_go.Case) any {
				if caseItem.Related == nil {
					caseItem.Related = &_go.RelatedCaseList{}
				}
				return scanner.ScanJSONToStructList(&caseItem.Related.Data)
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

func buildRelatedCasesSubquery(caseAlias string) (sq.SelectBuilder, error) {
	return sq.Select(`
        JSON_AGG(JSON_BUILD_OBJECT(
            'id', rc.id,
            'related_case', JSON_BUILD_OBJECT(
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
            'relation_type', rc.relation_type -- No casting needed for enum type
        )) AS related_cases
    `).
		From("cases.related_case rc").
		Join("cases.case c_child ON rc.related_case_id = c_child.id").
		LeftJoin("directory.wbt_user u ON rc.created_by = u.id").
		Where(fmt.Sprintf("%s = %s.id", store.Ident("rc", "primary_case_id"), caseAlias)), nil
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
func convertToCaseScanArgs(plan []func(caseItem *_go.Case) any, caseItem *_go.Case) []any {
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

func addCaseRbacCondition(auth auth.Auther, access auth.AccessMode, query sq.SelectBuilder, dependencyColumn string) (sq.SelectBuilder, error) {
	if auth != nil && auth.GetObjectScope(casesObjClassScopeName).IsRbacUsed() {
		subquery := sq.Select("acl.object").From("cases.case_acl acl").
			Where("acl.dc = ?", auth.GetDomainId()).
			Where(fmt.Sprintf("acl.object = %s", dependencyColumn)).
			Where("acl.subject = any( ?::int[])", pq.Array(auth.GetRoles())).
			Where("acl.access & ? = ?", int64(access), int64(access)).
			Limit(1)
		return query.Where("exists(?)", subquery), nil

	}
	return query, nil
}

func addCaseRbacConditionForDelete(auth auth.Auther, access auth.AccessMode, query sq.DeleteBuilder, dependencyColumn string) (sq.DeleteBuilder, error) {
	if auth != nil && auth.GetObjectScope(casesObjClassScopeName).IsRbacUsed() {
		subquery := sq.Select("acl.object").From("cases.case_acl acl").
			Where("acl.dc = ?", auth.GetDomainId()).
			Where(fmt.Sprintf("acl.object = %s", dependencyColumn)).
			Where("acl.subject = any( ?::int[])", pq.Array(auth.GetRoles())).
			Where("acl.access & ? = ?", int64(access), int64(access)).
			Limit(1)
		return query.Where("exists(?)", subquery), nil

	}
	return query, nil
}

func addCaseRbacConditionForUpdate(auth auth.Auther, access auth.AccessMode, query sq.UpdateBuilder, dependencyColumn string) (sq.UpdateBuilder, error) {
	if auth != nil && auth.GetObjectScope(casesObjClassScopeName).IsRbacUsed() {
		subquery := sq.Select("acl.object").From("cases.case_acl acl").
			Where("acl.dc = ?", auth.GetDomainId()).
			Where(fmt.Sprintf("acl.object = %s", dependencyColumn)).
			Where("acl.subject = any( ?::int[])", pq.Array(auth.GetRoles())).
			Where("acl.access & ? = ?", int64(access), int64(access)).
			Limit(1)
		return query.Where("exists(?)", subquery), nil

	}
	return query, nil
}

func addCaseRbacConditionForInsert(auth auth.Auther, access auth.AccessMode, query sq.InsertBuilder, caseId int64, alias string) (sq.InsertBuilder, error) {
	var subquery sq.SelectBuilder
	if auth != nil && auth.GetObjectScope(casesObjClassScopeName).IsRbacUsed() {
		subquery = sq.Select("acl.object").From("cases.case_acl acl").
			Where("acl.dc = ?", auth.GetDomainId()).
			Where("acl.object = ?", caseId).
			Where("acl.subject = any( ?::int[])", pq.Array(auth.GetRoles())).
			Where("acl.access & ? = ?", int64(access), int64(access)).
			Limit(1)
	} else {
		subquery = sq.Select("id").From("cases.case").Where("id = ?")
	}
	sql, args, err := store.FormAsCTE(subquery, alias)
	if err != nil {
		return query, err
	}
	return query.Prefix(sql, args...), nil
}

func (l *CaseStore) scanCases(rows pgx.Rows, plan []func(link *_go.Case) any) ([]*_go.Case, error) {
	var res []*_go.Case

	for rows.Next() {
		link, err := l.scanCase(pgx.Row(rows), plan)
		if err != nil {
			return nil, err
		}
		res = append(res, link)
	}
	return res, nil
}

func (l *CaseStore) scanCase(row pgx.Row, plan []func(link *_go.Case) any) (*_go.Case, error) {
	var link _go.Case
	var scanPlan []any
	for _, scan := range plan {
		scanPlan = append(scanPlan, scan(&link))
	}
	err := row.Scan(scanPlan...)
	if err != nil {
		return nil, err
	}

	return &link, nil
}
