package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"maps"
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
	"github.com/webitel/cases/internal/store/postgres/scanner"
	"github.com/webitel/cases/internal/store/postgres/transaction"
	dbutil "github.com/webitel/cases/internal/store/util"
	storeutils "github.com/webitel/cases/internal/store/util"
	"github.com/webitel/cases/model"
	"github.com/webitel/cases/model/options"
	"github.com/webitel/cases/model/options/defaults"
	common "github.com/webitel/cases/model/options/grpc"
	"github.com/webitel/cases/util"

	customrel "github.com/webitel/custom/reflect"
	custompgx "github.com/webitel/custom/store/postgres"
)

type CaseStore struct {
	storage   *Store
	mainTable string
}

const (
	caseLeft                  = "c"
	caseDefaultSort           = "-created_at"
	caseCreatedByAlias        = "cb"
	caseUpdatedByAlias        = "ub"
	caseSourceAlias           = "src"
	caseCloseReasonGroupAlias = "crg"
	caseAuthorAlias           = "auth"
	caseCloseReasonAlias      = "cr"
	caseSlaAlias              = "sl"
	caseStatusAlias           = "st"
	casePriorityAlias         = "pr"
	caseServiceAlias          = "svc"
	caseAssigneeAlias         = "ass" // :))
	caseReporterAlias         = "rp"
	caseImpactedAlias         = "im"
	caseGroupAlias            = "grp"
	caseSlaConditionAlias     = "cond"
	caseRelatedAlias          = "related"
	caseLinksAlias            = "links"
)

func (c *CaseStore) Create(
	rpc options.CreateOptions,
	add *_go.Case,
) (*_go.Case, error) {
	// Get the database connection
	d, dbErr := c.storage.Database()
	if dbErr != nil {
		return nil, dberr.NewDBInternalError("postgres.case.create.database_connection_error", dbErr)
	}

	// Begin a transaction
	tx, err := d.Begin(rpc)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.case.create.transaction_error", err)
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		err := tx.Rollback(ctx)
		if err != nil {
			log.Println("postgres.case.create.transaction_error", err)
		}
	}(tx, rpc)
	txManager := transaction.NewTxManager(tx)

	// scan service related defaults
	serviceDefs, err := c.ScanServiceDefs(
		rpc,
		txManager,
		add.Service.GetId(),
		add.Priority.GetId(),
	)
	if err != nil {
		return nil, err
	}

	if serviceDefs.StatusID == 0 {
		return nil, dberr.NewDBBadRequestError("postgres.case.create.missing.params", "StatusID")
	}

	if serviceDefs.CloseReasonGroupID == 0 {
		return nil, dberr.NewDBBadRequestError("postgres.case.create.missing.params", "CloseReasonGroupID")
	}

	// Calculate planned times within the transaction
	err = c.calculateTimings(
		nil,
		rpc,
		serviceDefs.CalendarID,
		serviceDefs.ReactionTime,
		serviceDefs.ResolutionTime,
		txManager,
		add,
	)

	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.case.create.calculate_planned_times_error", err)
	}

	// Build the query
	selectBuilder, plan, err := c.buildCreateCaseSqlizer(
		rpc,
		add,
		serviceDefs,
	)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.case.create.build_query_error", err)
	}

	// Generate the SQL and arguments
	query, args, err := selectBuilder.ToSql()
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.case.create.query_to_sql_error", err)
	}

	query = storeutils.CompactSQL(query)

	// Prepare the scan arguments
	scanArgs := convertToCaseScanArgs(plan, add)

	// Execute the query
	if err = txManager.QueryRow(rpc, query, args...).Scan(scanArgs...); err != nil {
		return nil, dberr.NewDBInternalError("postgres.case.create.execution_error", err)
	}

	// Commit the transaction
	if err := tx.Commit(rpc); err != nil {
		return nil, dberr.NewDBInternalError("postgres.case.create.commit_error", err)
	}

	for _, field := range rpc.GetFields() {
		if field == "role_ids" {
			roles, defErr := c.GetRolesById(rpc, add.GetId(), auth.Read)
			if defErr != nil {
				return nil, defErr
			}
			add.RoleIds = roles
			break
		}
	}

	return add, nil
}

type ServiceRelatedDefs struct {
	SLAID              int
	SLAConditionID     int
	ReactionTime       int
	ResolutionTime     int
	CalendarID         int
	StatusID           int
	CloseReasonGroupID int
}

// ScanServiceDefs fetches the SLA ID, reaction time, resolution time, calendar ID, and SLA condition ID for the last child service with a non-NULL SLA ID.
func (c *CaseStore) ScanServiceDefs(
	ctx context.Context,
	txManager *transaction.TxManager,
	serviceID int64,
	priorityID int64,
) (*ServiceRelatedDefs, error) {
	var res ServiceRelatedDefs

	err := txManager.QueryRow(ctx, `
WITH RECURSIVE
    service_hierarchy AS (
        SELECT id,
               root_id,
               sla_id,
               status_id,
               close_reason_group_id,
               ARRAY[id] AS path
        FROM cases.service_catalog
        WHERE id = $1

        UNION ALL

        SELECT sc.id,
               sc.root_id,
               COALESCE(sc.sla_id, sh.sla_id),
               COALESCE(sc.status_id, sh.status_id),
               COALESCE(sc.close_reason_group_id, sh.close_reason_group_id),
               sh.path || sc.id
        FROM cases.service_catalog sc
        INNER JOIN service_hierarchy sh ON sc.id = sh.root_id
    ),
    deepest_service AS (
        SELECT id,
               sla_id,
               status_id,
               close_reason_group_id,
               path
        FROM service_hierarchy
        WHERE sla_id IS NOT NULL
          AND status_id IS NOT NULL
          AND close_reason_group_id IS NOT NULL
        ORDER BY array_length(path, 1) ASC
        LIMIT 1
    ),
    priority_condition AS (
        SELECT sc.id AS sla_condition_id,
               sc.reaction_time,
               sc.resolution_time
        FROM cases.sla_condition sc
        INNER JOIN cases.priority_sla_condition psc ON sc.id = psc.sla_condition_id
        INNER JOIN cases.sla sla ON sc.sla_id = sla.id
        INNER JOIN deepest_service ds ON sla.id = ds.sla_id
        WHERE psc.priority_id = $2
        LIMIT 1
    )
SELECT ds.sla_id,
       COALESCE(pc.reaction_time, sla.reaction_time),
       COALESCE(pc.resolution_time, sla.resolution_time),
       sla.calendar_id,
       pc.sla_condition_id,
       ds.status_id,
       ds.close_reason_group_id
FROM deepest_service ds
LEFT JOIN priority_condition pc ON true
LEFT JOIN cases.sla sla ON ds.sla_id = sla.id;
`, serviceID, priorityID).Scan(
		scanner.ScanInt(&res.SLAID),
		scanner.ScanInt(&res.ReactionTime),
		scanner.ScanInt(&res.ResolutionTime),
		scanner.ScanInt(&res.CalendarID),
		scanner.ScanInt(&res.SLAConditionID),
		scanner.ScanInt(&res.StatusID),
		scanner.ScanInt(&res.CloseReasonGroupID),
	)
	if err != nil {
		return nil, dberr.NewDBInternalError("failed to scan SLA: %w", err)
	}

	return &res, nil
}

func (c *CaseStore) buildCreateCaseSqlizer(
	rpc options.CreateOptions,
	input *_go.Case,
	serviceDefs *ServiceRelatedDefs,
) (
	sq.SelectBuilder,
	[]func(caseItem *_go.Case) any,
	error,
) {

	// Extract optional fields via helper utils
	assignee := dbutil.IDPtr(input.GetAssignee())
	closeReason := dbutil.IDPtr(input.GetCloseReason())
	reporter := dbutil.IDPtr(input.GetReporter())
	impacted := dbutil.IDPtr(input.GetImpacted())
	group := dbutil.IDPtr(input.Group)
	description := dbutil.StringPtr(input.Description)
	closeResult := dbutil.StringPtr(input.GetCloseResult())

	// Set fallback defaults for status and close reason group
	defStatusID := input.Status.GetId()
	if defStatusID == 0 {
		defStatusID = int64(serviceDefs.StatusID)
	}

	defCloseReasonGroupID := input.CloseReasonGroup.GetId()
	if defCloseReasonGroupID == 0 {
		defCloseReasonGroupID = int64(serviceDefs.CloseReasonGroupID)
	}

	// Default user from token
	userID := rpc.GetAuthOpts().GetUserId()

	// Override if input.CreatedBy is explicitly provided
	if createdBy := input.GetCreatedBy(); createdBy != nil && createdBy.Id != 0 {
		userID = createdBy.Id
	}

	params := map[string]any{
		// Case-level parameters
		"date":                rpc.RequestTime(),
		"contact_info":        input.GetContactInfo(),
		"user":                userID,
		"dc":                  rpc.GetAuthOpts().GetDomainId(),
		"sla":                 serviceDefs.SLAID,
		"sla_condition":       serviceDefs.SLAConditionID,
		"status":              defStatusID,
		"status_condition":    input.StatusCondition.GetId(),
		"service":             input.Service.GetId(),
		"priority":            input.Priority.GetId(),
		"source":              input.Source.GetId(),
		"contact_group":       group,
		"close_reason_group":  defCloseReasonGroupID,
		"close_result":        closeResult,
		"close_reason":        closeReason,
		"rating":              input.Rating,
		"rating_comment":      input.RatingComment,
		"subject":             input.Subject,
		"planned_reaction_at": util.LocalTime(input.PlannedReactionAt),
		"planned_resolve_at":  util.LocalTime(input.PlannedResolveAt),
		"reporter":            reporter,
		"impacted":            impacted,
		"description":         description,
		"assignee":            assignee,
		//-------------------------------------------------//
		//------ CASE One-to-Many ( 1 : n ) Attributes ----//
		//-------------------------------------------------//
		// Links and related cases as JSON arrays
		"links":   extractLinksJSON(input.Links),
		"related": extractRelatedJSON(input.Related),
	}

	priorityCTE := `
	priority_cte AS (
		SELECT COALESCE(NULLIF(:priority, 0), id) AS priority_id
		FROM cases.priority
		ORDER BY id
		LIMIT 1
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

	statusConditionCTE := ""
	useStatusConditionRef := ":status_condition" // default

	if input.GetStatusCondition().GetId() == 0 {
		statusConditionCTE = `
		status_condition_cte AS (
			SELECT sc.id AS status_condition_id
			FROM cases.status_condition sc
			WHERE sc.status_id = :status AND sc.initial = true
		),`
		useStatusConditionRef = "(SELECT status_condition_id FROM status_condition_cte)"
	}

	// Consolidated query for inserting the case, links, and related cases
	query := `
	WITH
		` + prefixCTE + `,
        ` + priorityCTE + `,
        ` + statusConditionCTE + `
		` + caseLeft + ` AS (
			INSERT INTO cases.case (
				id, name, dc, created_at, created_by, updated_at, updated_by,
				priority, source, status, contact_group, close_reason_group,
				subject, planned_reaction_at, planned_resolve_at, reporter, impacted,
				service, description, assignee, sla, sla_condition_id, status_condition, contact_info,
				close_result, close_reason, rating, rating_comment
			) VALUES (
				(SELECT id FROM id_cte),
				CONCAT((SELECT prefix FROM prefix_cte), '_', (SELECT id FROM id_cte)),
				:dc, :date, :user, :date, :user,
				(SELECT priority_id FROM priority_cte), :source, :status, :contact_group, :close_reason_group,
				:subject, :planned_reaction_at, :planned_resolve_at, :reporter, :impacted,
				:service, :description, :assignee, :sla, :sla_condition,
				` + useStatusConditionRef + `, :contact_info, :close_result, :close_reason, 
                NULLIF(:rating, 0), NULLIF(:rating_comment, '')
			)
			RETURNING *
		),
		` + caseLinksAlias + ` AS (
			INSERT INTO cases.case_link (
				name, url, dc, created_by, created_at, updated_by, updated_at, case_id
			)
			SELECT
				item ->> 'name',
				item ->> 'url',
				:dc, :user, :date, :user, :date, (SELECT id FROM ` + caseLeft + `)
			FROM jsonb_array_elements(:links) AS item
		),
		` + caseRelatedAlias + ` AS (
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

	// region: [custom] fields ..
	var custom *customCtx
	if data := input.GetCustom(); len(data.GetFields()) > 0 {
		custom = c.custom(rpc)
		if custom == nil || custom.refer == nil {
			// No [custom] extensions/cases -BUT- case.Custom data specified !
			err := fmt.Errorf("custom: no specification")
			return sq.SelectBuilder{}, nil, err
		}
		// PREPARE Statement !..
		oid := sq.Expr("(SELECT id FROM " + caseLeft + ")")
		insertQ, args, err := custom.refer.Update(
			oid, data, false, // [!]partial
		)
		if err != nil {
			return sq.SelectBuilder{}, nil, dberr.NewDBInternalError(
				"postgres.custom.prepare.error", err,
			)
		}
		if insertQ != nil {
			insertQ, _, err := insertQ.ToSql()
			if err != nil {
				return sq.SelectBuilder{}, nil, dberr.NewDBInternalError(
					"postgres.custom.prepare.error", err,
				)
			}
			cte := "x" // alias
			query += ", " + cte + " AS (" + insertQ + ")"
			for name, value := range args {
				params[name] = value
			}
			custom.table = cte
			// Return INSERT[ed] data record ! normalized
			custom.fields = make([]string, 0, len(data.Fields))
			maps.Keys(data.Fields)(func(name string) bool {
				// [NOTE]: MAY be unknown field name !
				custom.fields = append(custom.fields, name)
				return true
			})
		} // else { // No INSERT to perform ! }
	}
	// sanitize: no source output !
	input.Custom = nil
	// endregion: [custom] fields ..

	// **Bind named query and parameters**
	// **This binds the named parameters in the query to the provided params map, converting it into a positional query with arguments.**
	// **Example:**
	// **  Query: "INSERT INTO cases.case (name, subject) VALUES (:name, :subject)"**
	// **  Params: map[string]interface{}{"name": "test_name", "subject": "test_subject"}**
	// **  Result: "INSERT INTO cases.case (name, col2) subject ($1, $2)", []interface{}{"test_name", "test_subject"}**
	boundQuery, args, err := storeutils.BindNamed(query, params)
	if err != nil {
		return sq.SelectBuilder{}, nil, dberr.NewDBInternalError("postgres.case.create.bind_named_error", err)
	}

	// Construct SELECT query to return case data
	selectBuilder, plan, err := c.buildCaseSelectColumnsAndPlan(
		withSearchOptions(rpc, func(search *common.SearchOptions) (_ error) {
			if custom != nil {
				// Output query custom fields ..
				search.UnknownFields = append(
					search.UnknownFields, custom.fields..., // customFieldName,
				)
				// Chain prepared query context
				search.Filters[customCtxState] = custom
			}
			return
		}),
		sq.Select().PrefixExpr(
			sq.Expr(
				boundQuery,
				args...,
			),
		),
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
	var jsonArray []map[string]any
	for _, link := range links.Items {
		jsonArray = append(jsonArray, map[string]any{
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
	var jsonArray []map[string]any
	for _, item := range related.Data {
		jsonArray = append(jsonArray, map[string]any{
			"id":   item.GetId(),
			"type": item.GetRelationType(),
		})
	}
	jsonData, _ := json.Marshal(jsonArray)
	return jsonData
}

type CalendarSlot struct {
	Day            int
	StartTimeOfDay int
	EndTimeOfDay   int
	Disabled       bool
}

type ExceptionSlot struct {
	Date           time.Time
	StartTimeOfDay int
	EndTimeOfDay   int
	Disabled       bool
	Repeat         bool
	Working        bool
}

type MergedSlot struct {
	Day            int       // Weekday (0-6, Sunday-Saturday)
	Date           time.Time // Specific date (can be empty if not an exception)
	StartTimeOfDay int       // Start time of the slot (in minutes from midnight)
	EndTimeOfDay   int       // End time of the slot (in minutes from midnight)
	Disabled       bool      // Is the slot disabled
}

type TimingOpts interface {
	RequestTime() time.Time
	GetAuthOpts() auth.Auther
	context.Context
}

func (c *CaseStore) calculateTimings(
	caseID *int64,
	rpc TimingOpts,
	calendarID int,
	reactionTime int,
	resolutionTime int,
	txManager *transaction.TxManager,
	caseItem *_go.Case,
) error {
	// Determine the pivot time
	var pivotTime time.Time
	if caseID == nil {
		pivotTime = rpc.RequestTime()
	} else {
		err := txManager.QueryRow(rpc, `
			SELECT created_at FROM cases.case WHERE id = $1`, *caseID).Scan(&pivotTime)
		if err != nil {
			return fmt.Errorf("failed to fetch created_at for caseID %d: %w", *caseID, err)
		}
	}

	// Fetch standard calendar working hours
	calendar, err := fetchCalendarSlots(rpc, txManager, calendarID)
	if err != nil {
		return err
	}

	// Fetch exceptions (overrides for specific days)
	exceptions, err := fetchExceptionSlots(rpc, txManager, calendarID)
	if err != nil {
		return err
	}

	// Merge calendar and exceptions into a single slice
	mergedSlots := mergeCalendarAndExceptions(calendar, exceptions)

	// Fetch timezone offset
	var offset time.Duration
	err = txManager.QueryRow(rpc, `
		SELECT tz.utc_offset
		FROM flow.calendar cl
		    LEFT JOIN flow.calendar_timezones tz ON tz.id = cl.timezone_id
		WHERE cl.id = $1`, calendarID).Scan(&offset)
	if err != nil {
		return fmt.Errorf("failed to fetch calendar offset: %w", err)
	}

	// Convert reaction and resolution times from seconds to minutes
	reactionMinutes := reactionTime / 60
	resolutionMinutes := resolutionTime / 60

	// Calculate planned reaction and resolution timestamps
	reactionTimestamp, err := calculateTimestampFromCalendar(pivotTime, offset, reactionMinutes, mergedSlots)
	if err != nil {
		return fmt.Errorf("failed to calculate planned reaction time: %w", err)
	}

	resolveTimestamp, err := calculateTimestampFromCalendar(pivotTime, offset, resolutionMinutes, mergedSlots)
	if err != nil {
		return fmt.Errorf("failed to calculate planned resolution time: %w", err)
	}

	caseItem.PlannedReactionAt = reactionTimestamp.UnixMilli()
	caseItem.PlannedResolveAt = resolveTimestamp.UnixMilli()

	return nil
}

// fetchCalendarSlots retrieves working hours for a calendar
func fetchCalendarSlots(rpc TimingOpts, txManager *transaction.TxManager, calendarID int) ([]CalendarSlot, error) {
	rows, err := txManager.Query(rpc, `
		SELECT day, start_time_of_day, end_time_of_day, disabled
		FROM flow.calendar cl,
		     UNNEST(accepts::flow.calendar_accept_time[]) x
		WHERE cl.id = $1
		    AND (NOT x.special OR x.special IS NULL)
		ORDER BY day, start_time_of_day`, calendarID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch calendar details: %w", err)
	}
	defer rows.Close()

	var calendar []CalendarSlot
	for rows.Next() {
		var entry CalendarSlot
		if err := rows.Scan(&entry.Day, &entry.StartTimeOfDay, &entry.EndTimeOfDay, &entry.Disabled); err != nil {
			return nil, fmt.Errorf("failed to scan calendar entry: %w", err)
		}

		// Adjust day value to make Sunday 0, Monday 1, Tuesday 2, ..., Saturday 6
		// If the DB starts with Monday as 0, we need to adjust it by adding 1
		// Example: Monday (0) becomes 1, Tuesday (1) becomes 2, ..., Sunday (6) becomes 0
		entry.Day = (entry.Day + 1) % 7 // This ensures Sunday is 0, Monday is 1, ..., Saturday is 6

		// Add adjusted entry to the calendar slice
		calendar = append(calendar, entry)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over calendar rows: %w", err)
	}
	return calendar, nil
}

// fetchExceptionSlots retrieves exceptions for specific days (overrides)
func fetchExceptionSlots(rpc TimingOpts, txManager *transaction.TxManager, calendarID int) ([]ExceptionSlot, error) {
	rows, err := txManager.Query(rpc, `
		SELECT
			to_timestamp(x.date / 1000) AS date,
			x.work_start AS start_time_of_day,
			x.work_stop AS end_time_of_day,
			x.disabled,
			x.repeat,
			x.working
		FROM flow.calendar cl, UNNEST(cl.excepts::flow.calendar_except_date[]) x
		WHERE cl.id = $1
		ORDER BY x.date, x.work_start`, calendarID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch calendar exceptions: %w", err)
	}
	defer rows.Close()

	var exceptions []ExceptionSlot
	for rows.Next() {
		var entry ExceptionSlot
		if err := rows.Scan(&entry.Date, &entry.StartTimeOfDay, &entry.EndTimeOfDay, &entry.Disabled, &entry.Repeat, &entry.Working); err != nil {
			return nil, fmt.Errorf("failed to scan exception entry: %w", err)
		}
		exceptions = append(exceptions, entry)
	}
	return exceptions, nil
}

// mergeCalendarAndExceptions merges calendar and exceptions into a single slice
func mergeCalendarAndExceptions(calendar []CalendarSlot, exceptions []ExceptionSlot) []MergedSlot {
	mergedSlots := make([]MergedSlot, 0)

	// Convert calendar slots to merged slots
	for _, cal := range calendar {
		// Adjust weekday to start from Sunday as 0, Monday as 1, etc.
		adjustedDay := cal.Day % 7 // Adjust to make sure it's in [0, 6] range
		mergedSlots = append(mergedSlots, MergedSlot{
			Day:            adjustedDay,
			Date:           time.Time{}, // Calendar slots don't have a specific date
			StartTimeOfDay: cal.StartTimeOfDay,
			EndTimeOfDay:   cal.EndTimeOfDay,
			Disabled:       cal.Disabled,
		})
	}

	// Override with exceptions
	for _, exception := range exceptions {
		// If working is false, set disabled to true
		if !exception.Working {
			exception.Disabled = true
		}

		// If the exception is set to repeat annually (not weekly)
		if exception.Repeat {
			// Set the exception to repeat every year on the same date
			mergedSlots = append(mergedSlots, MergedSlot{
				Day:            -1,             // Special indicator for a specific date
				Date:           exception.Date, // Use the specific exception date
				StartTimeOfDay: exception.StartTimeOfDay,
				EndTimeOfDay:   exception.EndTimeOfDay,
				Disabled:       exception.Disabled,
			})
		} else {
			// Specific date exception (non-repeating)
			mergedSlots = append(mergedSlots, MergedSlot{
				Day:            -1, // Special indicator for a specific date
				Date:           exception.Date,
				StartTimeOfDay: exception.StartTimeOfDay,
				EndTimeOfDay:   exception.EndTimeOfDay,
				Disabled:       exception.Disabled,
			})
		}
	}

	return mergedSlots
}

func calculateTimestampFromCalendar(
	startTime time.Time,
	calendarOffset time.Duration,
	requiredMinutes int,
	mergedSlots []MergedSlot,
) (time.Time, error) {
	remainingMinutes := requiredMinutes
	currentTimeInMinutes := startTime.Hour()*60 + startTime.Minute() // UTC
	addDays := 0

	// Process each day while there are remaining minutes
	for remainingMinutes > 0 {
		// Calculate current day date
		currentDayDate := startTime.AddDate(0, 0, addDays)

		// Check if today is a disabled exception and skip if true
		// This ensures that we skip the day only once if both calendar and exception are disabled
		skipDay := false
		for _, slot := range mergedSlots {
			if slot.Disabled && !slot.Date.IsZero() && isSameDate(slot.Date, currentDayDate) {
				// If today is marked as disabled in exception, skip this day
				skipDay = true
				break
			}
		}

		// Skip the whole day if it's an exception or calendar day
		if skipDay {
			addDays++
			currentTimeInMinutes = 0
			continue
		}

		// Check for date-specific slots first (exceptions)
		dateSpecificSlotFound := false
		for _, slot := range mergedSlots {
			if slot.Disabled {
				continue
			}

			// If this is a date-specific slot, check if it matches the current date
			if !slot.Date.IsZero() && isSameDate(slot.Date, currentDayDate) {
				dateSpecificSlotFound = true

				// Convert slot times to UTC minutes considering the calendar offset
				slotStartUtc := slot.StartTimeOfDay - int(calendarOffset.Minutes())
				slotEndUtc := slot.EndTimeOfDay - int(calendarOffset.Minutes())

				// Ensure we start counting from the correct time (taking currentTimeInMinutes into account)
				startingAt := max(currentTimeInMinutes, slotStartUtc)
				if slotEndUtc <= startingAt {
					continue
				}

				// Calculate the available minutes in the interval
				availableMinutes := slotEndUtc - startingAt

				// If enough minutes are available, finalize the time
				if availableMinutes >= remainingMinutes {
					finalTime := currentDayDate
					finalTime = time.Date(
						finalTime.Year(),
						finalTime.Month(),
						finalTime.Day(),
						0, 0, 0, 0,
						finalTime.Location(),
					)
					finalTime = finalTime.Add(time.Duration(startingAt+remainingMinutes) * time.Minute)
					return finalTime, nil
				}

				// Deduct available minutes and move to the next interval
				remainingMinutes -= availableMinutes
				currentTimeInMinutes = slotEndUtc
			}
		}

		// If no date-specific slot was found, fall back to regular day-of-week slots
		if !dateSpecificSlotFound {
			for _, slot := range mergedSlots {
				if slot.Disabled {
					continue
				}

				// Skip date-specific slots (already processed above)
				if !slot.Date.IsZero() {
					continue
				}

				// For calendar slots, ensure we match the correct weekday
				if int(currentDayDate.Weekday()) != slot.Day {
					continue
				}

				// Convert slot times to UTC minutes considering the calendar offset
				slotStartUtc := slot.StartTimeOfDay - int(calendarOffset.Minutes())
				slotEndUtc := slot.EndTimeOfDay - int(calendarOffset.Minutes())

				// Ensure we start counting from the correct time (taking currentTimeInMinutes into account)
				startingAt := max(currentTimeInMinutes, slotStartUtc)
				if slotEndUtc <= startingAt {
					continue
				}

				// Calculate the available minutes in the interval
				availableMinutes := slotEndUtc - startingAt

				// If enough minutes are available, finalize the time
				if availableMinutes >= remainingMinutes {
					finalTime := currentDayDate
					finalTime = time.Date(
						finalTime.Year(),
						finalTime.Month(),
						finalTime.Day(),
						0, 0, 0, 0,
						finalTime.Location(),
					)
					finalTime = finalTime.Add(time.Duration(startingAt+remainingMinutes) * time.Minute)
					return finalTime, nil
				}

				// Deduct available minutes and move to the next interval
				remainingMinutes -= availableMinutes
				currentTimeInMinutes = slotEndUtc
			}
		}

		// Move to the next day if we haven't allocated all the required minutes
		addDays++

		// Reset the start time for the next day to 00:00
		currentTimeInMinutes = 0
	}

	return time.Time{}, errors.New("unable to allocate required minutes")
}

// Helper function to check if two dates are the same (ignoring time)
func isSameDate(date1, date2 time.Time) bool {
	return date1.Year() == date2.Year() &&
		date1.Month() == date2.Month() &&
		date1.Day() == date2.Day()
}

// Delete implements store.CaseStore.
func (c *CaseStore) Delete(rpc options.DeleteOptions) error {
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
	res, execErr := d.Exec(rpc, query, args...)
	if execErr != nil {
		return dberr.NewDBInternalError("store.case.delete.exec_error", execErr)
	}

	// Check if any rows were affected
	if res.RowsAffected() == 0 {
		return dberr.NewDBNoRowsError("store.case.delete.not_found")
	}

	return nil
}

func (c CaseStore) buildDeleteCaseQuery(rpc options.DeleteOptions) (string, []interface{}, error) {
	var err error
	convertedIds := util.Int64SliceToStringSlice(rpc.GetIDs())
	ids := util.FieldsFunc(convertedIds, util.InlineFields)
	query := sq.Delete("cases.case").Where("id = ANY(?)", ids).Where("dc = ?", rpc.GetAuthOpts().GetDomainId()).PlaceholderFormat(sq.Dollar)
	query, err = addCaseRbacConditionForDelete(rpc.GetAuthOpts(), auth.Delete, query, "id")
	if err != nil {
		return "", nil, err
	}

	return query.ToSql()
}

// List implements store.CaseStore.
func (c *CaseStore) List(opts options.SearchOptions) (*_go.CaseList, error) {
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
	slct = storeutils.CompactSQL(slct)
	db, dbErr := c.storage.Database()
	if dbErr != nil {
		return nil, dbErr
	}
	rows, err := db.Query(opts, storeutils.CompactSQL(slct), args...)
	if err != nil {
		return nil, dberr.NewDBError("postgres.case.list.exec.error", err.Error())
	}
	var res _go.CaseList
	res.Items, err = c.scanCases(rows, plan)
	if err != nil {
		return nil, err
	}
	res.Items, res.Next = storeutils.ResolvePaging(opts.GetSize(), res.Items)
	res.Page = int64(opts.GetPage())
	return &res, nil
}

func (c *CaseStore) CheckRbacAccess(ctx context.Context, auth auth.Auther, access auth.AccessMode, caseId int64) (bool, error) {
	if auth == nil {
		return false, nil
	}
	if !auth.IsRbacCheckRequired(model.ScopeCases, access) {
		return true, nil
	}
	q := sq.Select("1").From("cases.case_acl acl").
		Where("acl.dc = ?", auth.GetDomainId()).
		Where("acl.object = ?", caseId).
		Where("acl.subject = any( ?::int[])", pq.Array(auth.GetRoles())).
		Where("acl.access & ? = ?", int64(access), int64(access)).
		Limit(1).PlaceholderFormat(sq.Dollar)
	db, err := c.storage.Database()
	if err != nil {
		return false, err
	}
	query, args, defErr := q.ToSql()
	if defErr != nil {
		return false, defErr
	}
	res, defErr := db.Exec(ctx, query, args...)
	if defErr != nil {
		return false, defErr
	}
	if res.RowsAffected() == 1 {
		return true, nil
	}
	return false, nil
}

func (c *CaseStore) buildListCaseSqlizer(opts options.SearchOptions) (sq.SelectBuilder, []func(caseItem *_go.Case) any, error) {
	base := sq.Select().From(fmt.Sprintf("%s %s", c.mainTable, caseLeft)).PlaceholderFormat(sq.Dollar)
	base, plan, err := c.buildCaseSelectColumnsAndPlan(opts, base)
	if search := opts.GetSearch(); search != "" {
		searchTerm, operator := storeutils.ParseSearchTerm(search)
		searchNumber := storeutils.PrepareSearchNumber(search)
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
			sq.Expr(fmt.Sprintf("%s %s ?", storeutils.Ident(caseLeft, "subject"), operator), searchTerm),
			sq.Expr(fmt.Sprintf("%s %s ?", storeutils.Ident(caseLeft, "name"), operator), searchTerm),
			sq.Expr(fmt.Sprintf("%s %s ?", storeutils.Ident(caseLeft, "contact_info"), operator), searchTerm),
		}
		base = base.Where(where)

	}
	if len(opts.GetIDs()) != 0 {
		base = base.Where(fmt.Sprintf("%s = ANY(?)", storeutils.Ident(caseLeft, "id")), opts.GetIDs())
	}
	for column, value := range opts.GetFilters() {
		switch column {
		case "created_by",
			"updated_by",
			"assignee",         // +
			"reporter",         // +
			"source",           // +
			"priority",         // +
			"status",           // +
			"impacted",         // +
			"close_reason",     // +
			"service",          // +
			"status_condition", // +
			"sla_condition",
			"group",
			"sla": // +
			dbColumn := column
			switch column {
			case "group":
				dbColumn = "contact_group"
			case "sla_condition":
				dbColumn = "sla_condition_id"
			}
			switch typedValue := value.(type) {
			case string:
				values := strings.Split(typedValue, ",")
				var (
					valuesInt []int64
					isNull    bool
					expr      sq.Or
				)
				for _, s := range values {
					if s == "" {
						continue
					}
					if s == "null" {
						isNull = true
						continue
					}
					converted, err := strconv.ParseInt(s, 10, 64)
					if err != nil {
						return base, nil, dberr.NewDBInternalError("postgres.case.build_list_case_sqlizer.convert_to_int_array.error", err)
					}
					valuesInt = append(valuesInt, converted)
				}
				col := storeutils.Ident(caseLeft, dbColumn)
				expr = append(expr, sq.Expr(fmt.Sprintf("%s = ANY(?::int[])", col), valuesInt))
				if isNull {
					expr = append(expr, sq.Expr(fmt.Sprintf("%s ISNULL", col)))
				}
				base = base.Where(expr)
			}
		case "status_condition.final":
			var final bool
			switch typedValue := value.(type) {
			case string:
				if typedValue == "true" {
					final = true
				}
			}
			base = base.Where(
				fmt.Sprintf("EXISTS(SELECT id FROM cases.status_condition WHERE id = %s AND final = ?)",
					storeutils.Ident(caseLeft, "status_condition"),
				),
				final,
			)
		case "author":
			switch typedValue := value.(type) {
			case string:
				values := strings.Split(typedValue, ",")
				var (
					valuesInt []int64
					isNull    bool
					expr      sq.Or
				)
				for _, s := range values {
					if s == "" {
						continue
					}
					if s == "null" {
						isNull = true
						continue
					}
					converted, err := strconv.ParseInt(s, 10, 64)
					if err != nil {
						return base, nil, dberr.NewDBInternalError("postgres.case.build_list_case_sqlizer.convert_to_int_array.error", err)
					}
					valuesInt = append(valuesInt, converted)
				}
				col := storeutils.Ident(caseAuthorAlias, "id")
				expr = append(expr, sq.Expr(fmt.Sprintf("%s = ANY(?::int[])", col), valuesInt))
				if isNull {
					expr = append(expr, sq.Expr(fmt.Sprintf("%s ISNULL", col)))
				}
				base = base.Where(expr)
			}
		case "communication_id":
			switch typedValue := value.(type) {
			case string:
				values := strings.Split(typedValue, ",")
				var (
					communicationUUIDs []string
					isNull             bool
					expr               sq.Or
				)
				for _, s := range values {
					if s == "" {
						continue
					}
					if s == "null" {
						isNull = true
						continue
					}
					communicationUUIDs = append(communicationUUIDs, s)
				}

				if len(communicationUUIDs) > 0 {
					expr = append(expr, sq.Expr(
						fmt.Sprintf(`EXISTS (
					SELECT 1 FROM cases.case_communication cc
					WHERE cc.case_id = %s AND cc.communication_id = ANY(?::text[])
				)`, storeutils.Ident(caseLeft, "id")),
						communicationUUIDs,
					))
				}

				if isNull {
					expr = append(expr, sq.Expr(
						fmt.Sprintf(`NOT EXISTS (
					SELECT 1 FROM cases.case_communication cc
					WHERE cc.case_id = %s
				)`, storeutils.Ident(caseLeft, "id")),
					))
				}

				if len(expr) > 0 {
					base = base.Where(expr)
				}
			}
		case "rating.from":
			cutted, _ := strings.CutSuffix(column, ".from")
			base = base.Where(fmt.Sprintf("%s >= ?::INT", storeutils.Ident(caseLeft, cutted)), value)
		case "rating.to":
			cutted, _ := strings.CutSuffix(column, ".to")
			base = base.Where(fmt.Sprintf("%s <= ?::INT", storeutils.Ident(caseLeft, cutted)), value)
		case "reacted_at.from", "resolved_at.from", "planned_reaction_at.from", "planned_resolve_at.from", "created_at.from":
			cutted, _ := strings.CutSuffix(column, ".from")
			base = base.Where(fmt.Sprintf("extract(epoch from %s)*1000::BIGINT > ?::BIGINT", storeutils.Ident(caseLeft, cutted)), value)
		case "reacted_at.to", "resolved_at.to", "planned_reaction_at.to", "planned_resolve_at.to", "created_at.to":
			cutted, _ := strings.CutSuffix(column, ".to")
			base = base.Where(fmt.Sprintf("extract(epoch from %s)*1000::BIGINT < ?::BIGINT", storeutils.Ident(caseLeft, cutted)), value)
		case "attachments":
			var operator string
			if value != "true" {
				operator = "NOT "
			}
			base = base.Where(sq.Expr(fmt.Sprintf(operator+"EXISTS (SELECT id FROM storage.files WHERE uuid = %s::varchar UNION SELECT id FROM cases.case_link WHERE case_link.case_id = %[1]s)", storeutils.Ident(caseLeft, "id"))))
		case "contact":
			base = base.Where(sq.Or{
				sq.Expr(fmt.Sprintf("%s.reporter = ?", caseLeft), value),
				sq.Expr(fmt.Sprintf("%s.assignee = ?", caseLeft), value),
			})
		}
	}
	if err != nil {
		return base, nil, err
	}

	if sess := opts.GetAuthOpts(); sess != nil {
		base = base.Where(storeutils.Ident(caseLeft, "dc = ?"), opts.GetAuthOpts().GetDomainId())
		base, err = addCaseRbacCondition(sess, auth.Read, base, storeutils.Ident(caseLeft, "id"))
	}
	// pagination
	base = storeutils.ApplyPaging(opts.GetPage(), opts.GetSize(), base)
	// sort
	sort := opts.GetSort()
	if sort == "" {
		sort = caseDefaultSort
	}
	field, direction := storeutils.GetSortingOperator(sort)
	var tableAlias string
	if !util.ContainsStringIgnoreCase(opts.GetFields(), field) { // not joined yet
		base, tableAlias, err = c.joinRequiredTable(base, field)
		if err != nil {
			return base, nil, err
		}
	} else { // get alias
		switch field {
		case "created_by":
			tableAlias = caseCreatedByAlias
		case "updated_by":
			tableAlias = caseUpdatedByAlias
		case "source":
			tableAlias = caseSourceAlias
		case "close_reason_group":
			tableAlias = caseCloseReasonGroupAlias
		case "sla":
			tableAlias = caseSlaAlias
		case "status":
			tableAlias = caseStatusAlias
		case "priority":
			tableAlias = casePriorityAlias
		case "service":
			tableAlias = caseServiceAlias
		case "author":
			tableAlias = caseAuthorAlias
		case "assignee":
			tableAlias = caseAssigneeAlias
		case "reporter":
			tableAlias = caseReporterAlias
		case "impacted":
			tableAlias = caseImpactedAlias
		case "group":
			tableAlias = caseGroupAlias
		case "close_reason":
			tableAlias = caseCloseReasonAlias
		case "sla_condition":
			tableAlias = caseSlaConditionAlias
		}
	}
	if tableAlias == "" {
		tableAlias = caseLeft
	}
	switch field {
	case "id", "ver", "created_at", "updated_at", "name", "subject", "description", "planned_reaction_at", "planned_resolve_at", "reacted_at", "resolved_at", "contact_info", "close_result", "rating_comment", "rating":
		base = base.OrderBy(fmt.Sprintf("%s %s", storeutils.Ident(tableAlias, field), direction))
	case "created_by", "updated_by", "source", "close_reason_group", "close_reason", "sla", "status_condition", "status", "priority", "service", "group":
		base = base.OrderBy(fmt.Sprintf("%s %s", storeutils.Ident(tableAlias, "name"), direction))
	case "author", "assignee", "reporter", "impacted":
		base = base.OrderBy(fmt.Sprintf("%s %s", storeutils.Ident(tableAlias, "common_name"), direction))
	case "sla_condition":
		base = base.OrderBy(fmt.Sprintf("%s %s", storeutils.Ident(tableAlias, "name"), direction))
	}

	return base, plan, nil
}

// region UPDATE
func (c *CaseStore) Update(
	rpc options.UpdateOptions,
	upd *_go.Case,
) (*_go.Case, error) {
	// Establish database connection
	db, err := c.storage.Database()
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.case.update.database_connection_error", err)
	}

	// Begin a transaction
	tx, txErr := db.Begin(rpc)
	if txErr != nil {
		return nil, dberr.NewDBInternalError("postgres.case.create.transaction_error", txErr)
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		err := tx.Rollback(ctx)
		if err != nil {
			log.Printf("postgres.case.update.rollback_error: %v\n", err)
		}
	}(tx, rpc)
	txManager := transaction.NewTxManager(tx)

	// * if user change Service -- SLA ; SLA Condition ; Planned Reaction / Resolve at ; Calendar could be changed
	if util.ContainsField(rpc.GetMask(), "service") {
		serviceDefs, err := c.ScanServiceDefs(
			rpc,
			txManager,
			upd.Service.GetId(),
			upd.Priority.GetId(),
		)
		if err != nil {
			return nil, dberr.NewDBInternalError("postgres.case.update.scan_sla_error", err)
		}

		oid := rpc.GetEtags()[0].GetOid()

		// Calculate planned times within the transaction
		err = c.calculateTimings(
			&oid,
			rpc,
			serviceDefs.CalendarID,
			serviceDefs.ReactionTime,
			serviceDefs.ResolutionTime,
			txManager,
			upd,
		)
		if err != nil {
			return nil, dberr.NewDBInternalError("postgres.case.update.calculate_planned_times_error", err)
		}

		// * assign new values ( SLA ; SLA Condition ; Planned Reaction / Resolve at ) to update (input) object
		if upd.Sla == nil {
			upd.Sla = &_go.Lookup{}
		}
		upd.Sla.Id = int64(serviceDefs.SLAID)
		if upd.SlaCondition == nil {
			upd.SlaCondition = &_go.Lookup{}
		}
		upd.SlaCondition.Id = int64(serviceDefs.SLAConditionID)
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

	query = storeutils.CompactSQL(query)

	// Prepare scan arguments
	scanArgs := convertToCaseScanArgs(plan, upd)

	if err := txManager.QueryRow(rpc, query, args...).Scan(scanArgs...); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, dberr.NewDBNoRowsError("postgres.case.update.update.scan_ver.not_found")
		}
		return nil, dberr.NewDBInternalError("postgres.case.update.update.execution_error", err)
	}

	// Commit the transaction
	if err := tx.Commit(rpc); err != nil {
		return nil, dberr.NewDBInternalError("postgres.case.update.commit_error", err)
	}
	for _, field := range rpc.GetFields() {
		if field == "role_ids" {
			roles, defErr := c.GetRolesById(rpc, upd.GetId(), auth.Read)
			if defErr != nil {
				return nil, defErr
			}
			upd.RoleIds = roles
			break
		}
	}

	return upd, nil
}

func (c *CaseStore) buildUpdateCaseSqlizer(
	rpc options.UpdateOptions,
	input *_go.Case,
) (sq.Sqlizer, []func(caseItem *_go.Case) any, error) {
	// Ensure required fields (ID and Version) are included
	fields := rpc.GetFields()
	fields = util.EnsureIdAndVerField(fields)
	var err error

	userID := rpc.GetAuthOpts().GetUserId()
	if util.ContainsField(rpc.GetMask(), "userID") {
		if updatedBy := input.GetUpdatedBy(); updatedBy != nil && updatedBy.Id != 0 {
			userID = updatedBy.Id
		}
	}

	// Initialize the update query
	updateBuilder := sq.Update(c.mainTable).
		PlaceholderFormat(sq.Dollar).
		Set("updated_at", rpc.RequestTime()).
		Set("updated_by", userID).
		Where(sq.Eq{
			"id":  rpc.GetEtags()[0].GetOid(),
			"ver": rpc.GetEtags()[0].GetVer(),
			"dc":  rpc.GetAuthOpts().GetDomainId(),
		})

	updateBuilder, err = addCaseRbacConditionForUpdate(rpc.GetAuthOpts(), auth.Edit, updateBuilder, "id")
	if err != nil {
		return nil, nil, err
	}
	// Increment version
	updateBuilder = updateBuilder.Set("ver", sq.Expr("ver + 1"))

	// region: [custom] fields ..
	var custom struct {
		customCtx
		update sq.Sqlizer
		params custompgx.Parameters
		output []common.SearchOption
	}
	// endregion: [custom] fields ..

	// Handle nested fields using switch-case on req.Mask
	for _, field := range rpc.GetMask() {
		switch field {
		case "subject":
			updateBuilder = updateBuilder.Set("subject", input.GetSubject())
		case "description":
			updateBuilder = updateBuilder.Set("description", sq.Expr("NULLIF(?, '')", input.Description))
		case "priority":
			updateBuilder = updateBuilder.Set("priority", input.Priority.GetId())
		case "source":
			updateBuilder = updateBuilder.Set("source", input.Source.GetId())
		case "status":
			updateBuilder = updateBuilder.Set("status", input.Status.GetId())
		case "status_condition":
			updateBuilder = updateBuilder.Set("status_condition", input.StatusCondition.GetId())
		case "service":
			prefixCTE := `
			WITH service_cte AS (
				SELECT catalog_id
				FROM cases.service_catalog
				WHERE id = ?
				LIMIT 1
			),
			prefix_cte AS (
				SELECT prefix
				FROM cases.service_catalog
				WHERE id = ANY(SELECT catalog_id FROM service_cte)
				LIMIT 1
			)
			SELECT prefix FROM prefix_cte`

			updateBuilder = updateBuilder.Set("service", input.Service.GetId())

			// Update SLA, SLA condition, and planned times
			updateBuilder = updateBuilder.Set("sla", input.Sla.GetId())
			updateBuilder = updateBuilder.Set("sla_condition_id", input.SlaCondition.GetId())
			updateBuilder = updateBuilder.Set("planned_resolve_at", util.LocalTime(input.GetPlannedResolveAt()))
			updateBuilder = updateBuilder.Set("planned_reaction_at", util.LocalTime(input.GetPlannedReactionAt()))

			caseIDString := strconv.FormatInt(rpc.GetEtags()[0].GetOid(), 10)

			updateBuilder = updateBuilder.Set("name",
				sq.Expr("CONCAT(("+prefixCTE+"), '_', CAST(? AS TEXT))",
					input.Service.GetId(), caseIDString))

		case "assignee":
			if input.Assignee.GetId() == 0 {
				updateBuilder = updateBuilder.Set("assignee", nil)
			} else {
				updateBuilder = updateBuilder.Set("assignee", input.Assignee.GetId())
			}
		case "reporter":
			updateBuilder = updateBuilder.Set("reporter", input.Reporter.GetId())
		case "contact_info":
			updateBuilder = updateBuilder.Set("contact_info", input.GetContactInfo())
		case "impacted":
			var impacted *int64
			if imp := input.GetImpacted().GetId(); imp != 0 {
				impacted = &imp
			}
			updateBuilder = updateBuilder.Set("impacted", impacted)
		case "group":
			if input.Group.GetId() == 0 {
				updateBuilder = updateBuilder.Set("contact_group", nil)
			} else {
				updateBuilder = updateBuilder.Set("contact_group", input.Group.GetId())
			}
		case "close_reason":
			if input.GetCloseReason() != nil {
				var closeReason *int64
				if reas := input.CloseReason.GetId(); reas > 0 {
					closeReason = &reas
				}
				updateBuilder = updateBuilder.Set("close_reason", closeReason)
			} else {
				updateBuilder = updateBuilder.Set("close_reason", nil)
			}
		case "close_result":
			var closeResult *string
			if res := input.GetCloseResult(); res != "" {
				closeResult = &input.CloseResult
			}
			updateBuilder = updateBuilder.Set("close_result", closeResult)
		case "rating":
			updateBuilder = updateBuilder.Set("rating", input.Rating)
		case "rating_comment":
			updateBuilder = updateBuilder.Set("rating_comment", sq.Expr("NULLIF(?, '')", input.RatingComment))
		// region: [custom] fields ..
		case "custom": // customFieldName
			{
				// [NOTE]: PATCH {"custom":null} !
				// get has [custom] extension defined !?
				if e := c.custom(rpc); e != nil {
					custom.customCtx = *e // shallowcopy
				}
				// record changes for update ..
				data := input.GetCustom()
				// sanitize: no source for output !
				input.Custom = nil
				// extension querier available ?
				if custom.refer == nil {
					// NO [custom] extension descriptor !
					if data != nil {
						err = fmt.Errorf("custom: no specification")
						return nil, nil, err
					}
					// no extension & data provided
					continue // ok ; next field ..
				}
				// PREPARE Statement !..
				// oid := rpc.GetEtags()[0].GetOid()
				oid := sq.Expr("(SELECT id FROM " + caseLeft + ")")
				const partial = true // [FIXME]: !
				updateQ, params, re := custom.refer.Update(
					oid, data, partial,
				)
				if err = re; err != nil {
					// failed to prepare UPDATE statement
					return nil, nil, err
				}
				if updateQ == nil {
					// No UPDATE to perform !
					continue // ok ; next field ..
				}
				custom.update = updateQ
				custom.params = params

				// tblname := strings.Split(custom.refer.Table(), ".")
				// ctename := tblname[len(tblname)-1]
				custom.table = "x" // ctename

				custom.fields = make([]string, 0, len(data.Fields))
				maps.Keys(data.Fields)(func(name string) bool {
					// [NOTE]: MAY be unknown field name !
					custom.fields = append(custom.fields, name)
					return true
				})
			}
			// endregion: [custom] fields ..
		}
	}

	WITH := sq.Select().PrefixExpr(
		sq.Expr("WITH "+caseLeft+" AS (?)",
			updateBuilder.Suffix("RETURNING *"),
		),
	) //.PlaceholderFormat(sq.Dollar)

	if custom.update != nil {
		// [RE]Bind (inject) :named paramenters !
		_, args, _ := WITH.Column("_").ToSql()
		query, _, _ := custom.update.ToSql()
		query, args, re := custompgx.BindNamedOffset(
			query, custom.params, len(args), // offset,
		)
		if err = re; err != nil {
			return nil, nil, err
		}
		//
		custom.update = sq.Expr(
			query, args...,
		)
		// WITH custom (..UPDATE..)
		WITH = WITH.PrefixExpr(sq.Expr(
			", "+custom.table+" AS (?)",
			custom.update,
		))
		// Return UPDATE[d] field(s) ...
		custom.output = append(custom.output,
			func(search *common.SearchOptions) (_ error) {
				search.UnknownFields = append(
					search.UnknownFields, custom.fields..., // customFieldName,
				)
				search.Filters[customCtxState] = &custom.customCtx
				return
			},
		)
	}

	// Define SELECT query for returning updated fields
	selectBuilder, plan, err := c.buildCaseSelectColumnsAndPlan(
		withSearchOptions(rpc, custom.output...), WITH,
	)
	if err != nil {
		return nil, nil, dberr.NewDBError("postgres.case.update.select_query_build_error", err.Error())
	}

	selectBuilder = selectBuilder.From(caseLeft)

	return selectBuilder, plan, nil
}

func (c *CaseStore) joinRequiredTable(base sq.SelectBuilder, field string) (q sq.SelectBuilder, joinedTableAlias string, err error) {
	var (
		tableAlias string
		joinTable  = func(neededAlias string, table string, connection string) {
			base = base.LeftJoin(fmt.Sprintf("%s %s ON %[2]s.id = %s", table, neededAlias, connection))
		}
	)

	switch field {
	case "created_by":
		tableAlias = caseCreatedByAlias
		joinTable(tableAlias, "directory.wbt_user", storeutils.Ident(caseLeft, "created_by"))
	case "updated_by":
		tableAlias = caseUpdatedByAlias
		joinTable(tableAlias, "directory.wbt_user", storeutils.Ident(caseLeft, "updated_by"))
	case "source":
		tableAlias = caseSourceAlias
		joinTable(tableAlias, "cases.source", storeutils.Ident(caseLeft, "source"))
	case "close_reason_group":
		tableAlias = caseCloseReasonGroupAlias
		joinTable(tableAlias, "cases.close_reason_group", storeutils.Ident(caseLeft, "close_reason_group"))
	case "author":
		createdByAlias := "cb_au"
		tableAlias = caseAuthorAlias
		joinTable(createdByAlias, "directory.wbt_user", storeutils.Ident(caseLeft, "created_by"))
		joinTable(tableAlias, "contacts.contact", storeutils.Ident(createdByAlias, "contact_id"))
	case "close_reason":
		tableAlias = caseCloseReasonAlias
		joinTable(tableAlias, "cases.close_reason", storeutils.Ident(caseLeft, "close_reason"))
	case "sla":
		tableAlias = caseSlaAlias
		joinTable(tableAlias, "cases.sla", storeutils.Ident(caseLeft, "sla"))
	case "status":
		tableAlias = caseStatusAlias
		joinTable(tableAlias, "cases.status", storeutils.Ident(caseLeft, "status"))
	case "priority":
		tableAlias = casePriorityAlias
		joinTable(tableAlias, "cases.priority", storeutils.Ident(caseLeft, "priority"))
	case "service":
		tableAlias = caseServiceAlias
		joinTable(tableAlias, "cases.service_catalog", storeutils.Ident(caseLeft, "service"))
	case "assignee":
		tableAlias = caseAssigneeAlias
		joinTable(tableAlias, "contacts.contact", storeutils.Ident(caseLeft, "assignee"))
	case "reporter":
		tableAlias = caseReporterAlias
		joinTable(tableAlias, "contacts.contact", storeutils.Ident(caseLeft, "reporter"))
	case "impacted":
		tableAlias = caseImpactedAlias
		joinTable(tableAlias, "contacts.contact", storeutils.Ident(caseLeft, "impacted"))
	case "group":
		tableAlias = caseGroupAlias
		joinTable(tableAlias, "contacts.group", storeutils.Ident(caseLeft, "contact_group"))
	case "sla_condition":
		tableAlias = caseSlaConditionAlias
		joinTable(tableAlias, "cases.sla_condition", storeutils.Ident(caseLeft, "sla_condition_id"))
	}
	return base, tableAlias, nil
}

// session required to get some columns
func (c *CaseStore) buildCaseSelectColumnsAndPlan(
	req options.SearchOptions, base sq.SelectBuilder,
) (
	sq.SelectBuilder, []func(caseItem *_go.Case) any, error,
) {
	var (
		plan       []func(caseItem *_go.Case) any
		tableAlias string
		err        error

		fields = req.GetFields()
		auther = req.GetAuthOpts()
	)

	for _, field := range fields {
		base, tableAlias, err = c.joinRequiredTable(base, field)
		if err != nil {
			return base, nil, err
		}
		// no table was joined
		if tableAlias == "" {
			tableAlias = caseLeft
		}
		switch field {
		case "id":
			base = base.Column(storeutils.Ident(tableAlias, "id AS case_id"))
			plan = append(plan, func(caseItem *_go.Case) any {
				return &caseItem.Id
			})
		case "ver":
			base = base.Column(storeutils.Ident(tableAlias, "ver"))
			plan = append(plan, func(caseItem *_go.Case) any {
				return &caseItem.Ver
			})
		case "created_by":
			base = base.Column(fmt.Sprintf(
				"ROW(%s.id, %[1]s.name)::text AS created_by", tableAlias))
			plan = append(plan, func(caseItem *_go.Case) any {
				return scanner.ScanRowLookup(&caseItem.CreatedBy)
			})
		case "created_at":
			base = base.Column(storeutils.Ident(tableAlias, "created_at"))
			plan = append(plan, func(caseItem *_go.Case) any {
				return scanner.ScanTimestamp(&caseItem.CreatedAt)
			})
		case "updated_by":
			base = base.Column(fmt.Sprintf(
				"ROW(%s.id, %[1]s.name)::text AS updated_by", tableAlias))
			plan = append(plan, func(caseItem *_go.Case) any {
				return scanner.ScanRowLookup(&caseItem.UpdatedBy)
			})
		case "updated_at":
			base = base.Column(storeutils.Ident(tableAlias, "updated_at"))
			plan = append(plan, func(caseItem *_go.Case) any {
				return scanner.ScanTimestamp(&caseItem.UpdatedAt)
			})
		case "name":
			base = base.Column(storeutils.Ident(tableAlias, "name"))
			plan = append(plan, func(caseItem *_go.Case) any {
				return &caseItem.Name
			})
		case "subject":
			base = base.Column(storeutils.Ident(tableAlias, "subject"))
			plan = append(plan, func(caseItem *_go.Case) any {
				return &caseItem.Subject
			})
		case "description":
			base = base.Column(storeutils.Ident(tableAlias, "description"))
			plan = append(plan, func(caseItem *_go.Case) any {
				return scanner.ScanText(&caseItem.Description)
			})
		case "group":
			base = base.Column(fmt.Sprintf(
				`ROW(%s.id, %[1]s.name,
							CASE
								WHEN EXISTS(SELECT id FROM contacts.dynamic_group WHERE id = %[1]s.id) THEN 'dynamic'
								ELSE 'static'
							END
						)::text  AS contact_group`, tableAlias))
			plan = append(plan, func(caseItem *_go.Case) any {
				return scanner.ScanRowExtendedLookup(&caseItem.Group)
			})
		case "role_ids":
			// skip
		case "source":
			base = base.Column(fmt.Sprintf(
				"ROW(%s.source, %[2]s.name, %[2]s.type)::text AS source", caseLeft, tableAlias))
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
			base = base.Column(storeutils.Ident(caseLeft, "planned_reaction_at"))
			plan = append(plan, func(caseItem *_go.Case) any {
				return scanner.ScanTimestamp(&caseItem.PlannedReactionAt)
			})
		case "planned_resolve_at":
			base = base.Column(storeutils.Ident(caseLeft, "planned_resolve_at"))
			plan = append(plan, func(caseItem *_go.Case) any {
				return scanner.ScanTimestamp(&caseItem.PlannedResolveAt)
			})
		case "close_reason_group":
			base = base.Column(fmt.Sprintf(
				"ROW(%s.id, %[1]s.name)::text  AS close_reason_group", tableAlias))
			plan = append(plan, func(caseItem *_go.Case) any {
				return scanner.ScanRowLookup(&caseItem.CloseReasonGroup)
			})
		case "author":
			base = base.Column(fmt.Sprintf(`ROW(%s.id, %[1]s.common_name)::text AS author`, tableAlias))
			plan = append(plan, func(caseItem *_go.Case) any {
				return scanner.ScanRowLookup(&caseItem.Author)
			})
		case "close_result":
			base = base.Column(storeutils.Ident(caseLeft, "close_result"))
			plan = append(plan, func(caseItem *_go.Case) any {
				return scanner.ScanText(&caseItem.CloseResult)
			})
		case "close_reason":
			base = base.Column(fmt.Sprintf(
				"ROW(%s.id, %[1]s.name)::text AS close_reason", tableAlias))
			plan = append(plan, func(caseItem *_go.Case) any {
				return scanner.ScanRowLookup(&caseItem.CloseReason)
			})
		case "rating":
			base = base.Column(storeutils.Ident(caseLeft, "rating"))
			plan = append(plan, func(caseItem *_go.Case) any {
				return scanner.ScanInt64(&caseItem.Rating)
			})
		case "rating_comment":
			base = base.Column(storeutils.Ident(caseLeft, "rating_comment"))
			plan = append(plan, func(caseItem *_go.Case) any {
				return scanner.ScanText(&caseItem.RatingComment)
			})
		case "resolved_at":
			base = base.
				Column(storeutils.Ident(caseLeft, "resolved_at"))
			plan = append(plan, func(caseItem *_go.Case) any {
				return scanner.ScanTimestamp(&caseItem.ResolvedAt)
			})
		case "reacted_at":
			base = base.
				Column(storeutils.Ident(caseLeft, "reacted_at"))
			plan = append(plan, func(caseItem *_go.Case) any {
				return scanner.ScanTimestamp(&caseItem.ReactedAt)
			})
		case "difference_in_reaction":
			base = base.
				Column(fmt.Sprintf(
					"COALESCE(CAST(EXTRACT(EPOCH FROM %s.reacted_at - %[1]s.created_at) * 1000 AS bigint), 0) AS difference_in_reaction",
					caseLeft,
				))
			plan = append(plan, func(caseItem *_go.Case) any {
				return scanner.ScanTimestamp(&caseItem.DifferenceInReaction)
			})
		case "difference_in_resolve":
			base = base.
				Column(fmt.Sprintf(
					"COALESCE(CAST(EXTRACT(EPOCH FROM %s.resolved_at - %[1]s.created_at) * 1000 AS bigint), 0) AS difference_in_resolve",
					caseLeft,
				))
			plan = append(plan, func(caseItem *_go.Case) any {
				return scanner.ScanTimestamp(&caseItem.DifferenceInResolve)
			})
		case "sla":
			base = base.Column(fmt.Sprintf(
				"ROW(%s.id, %[1]s.name)::text AS sla", tableAlias))
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
			base = base.Column(fmt.Sprintf(`ROW(%s.id, %[1]s.name)::text AS status`, tableAlias))
			plan = append(plan, func(caseItem *_go.Case) any {
				return scanner.ScanRowLookup(&caseItem.Status)
			})
		case "priority":
			base = base.Column(fmt.Sprintf("ROW(%s.id, %[1]s.name, %[1]s.color)::text AS priority", tableAlias))
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
			base = base.Column(fmt.Sprintf("ROW(%s.id, %[1]s.name)::text AS service", tableAlias))
			plan = append(plan, func(caseItem *_go.Case) any {
				return scanner.ScanRowLookup(&caseItem.Service)
			})
		case "assignee":
			base = base.Column(fmt.Sprintf(
				"ROW(%s.id, %[1]s.common_name)::text AS assignee", tableAlias))
			plan = append(plan, func(caseItem *_go.Case) any {
				return scanner.ScanRowLookup(&caseItem.Assignee)
			})
		case "reporter":
			base = base.Column(fmt.Sprintf(
				"ROW(%s.id, %[1]s.common_name)::text AS reporter", tableAlias))
			plan = append(plan, func(caseItem *_go.Case) any {
				return scanner.ScanRowLookup(&caseItem.Reporter)
			})
		case "contact_info":
			base = base.Column(storeutils.Ident(caseLeft, field))
			plan = append(plan, func(caseItem *_go.Case) any {
				return scanner.ScanText(&caseItem.ContactInfo)
			})
		case "impacted":
			base = base.Column(fmt.Sprintf(
				"ROW(%s.id, %[1]s.common_name)::text AS impacted", tableAlias))
			plan = append(plan, func(caseItem *_go.Case) any {
				return scanner.ScanRowLookup(&caseItem.Impacted)
			})
		case "sla_condition":
			base = base.Column(fmt.Sprintf("ROW(%s.id, %[1]s.name)::text AS sla_condition", tableAlias))
			plan = append(plan, func(caseItem *_go.Case) any {
				return scanner.ScanRowLookup(&caseItem.SlaCondition)
			})
		case "comments":
			commentFields := []string{"id", "ver", "text", "created_by", "author", "created_at", "can_edit"}
			subquery, scanPlan, dbErr := buildCommentsSelectAsSubquery(auther, commentFields, caseLeft)
			if dbErr != nil {
				return base, nil, dbErr
			}
			base = AddSubqueryAsColumn(base, subquery, "comments", false)
			plan = append(plan, func(value *_go.Case) any {
				var items []*_go.CaseComment
				postProcessing := func() error {
					if len(items) == 0 {
						return nil
					}
					res := &_go.CaseCommentList{}
					res.Items, res.Next = storeutils.ResolvePaging(defaults.DefaultSearchSize, items)
					res.Page = 1
					value.Comments = res
					return nil
				}
				return scanner.GetCompositeTextScanFunction(scanPlan, &items, postProcessing)
			})
		case "links":
			linksFields := []string{"id", "ver", "name", "url", "created_by", "author", "created_at"}
			subquery, scanPlan, dbErr := buildLinkSelectAsSubquery(linksFields, caseLeft)
			if dbErr != nil {
				return base, nil, dbErr
			}
			base = AddSubqueryAsColumn(base, subquery, field, false)
			plan = append(plan, func(value *_go.Case) any {
				var items []*_go.CaseLink
				postProcessing := func() error {
					if len(items) == 0 {
						return nil
					}
					res := &_go.CaseLinkList{}
					res.Items, res.Next = storeutils.ResolvePaging(defaults.DefaultSearchSize, items)
					res.Page = 1
					value.Links = res
					return nil
				}
				return scanner.GetCompositeTextScanFunction(scanPlan, &items, postProcessing)
			})
		case "files":
			filesFields := []string{
				"id",
				"size",
				"mime",
				"name",
				"created_at",
			}
			subquery, scanPlan, filtersApplied, dbErr := buildFilesSelectAsSubquery(filesFields, caseLeft)
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
					res.Items, res.Next = storeutils.ResolvePaging(defaults.DefaultSearchSize, items)
					res.Page = 1
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

	if unknown := req.GetUnknownFields(); len(unknown) > 0 {
		// custom [extensions/cases] configuration ?!
		custom := c.custom(req)
		// found & available ?
		if custom != nil && custom.refer != nil {
			// [TODO]: grab known fields for single request
			var (
				// known {nested..} fields query
				nested = make([]string, 0, len(unknown))
				field  customrel.FieldDescriptor
				ok     bool // e.g.: ?fields=custom{*}
			)
			for _, name := range unknown {
				switch name {
				case customFieldName, "*", "+":
					{
						// common field name ; ALL {nested..}
						ok = true
						continue
					}
				}
				field = custom.typof.Fields().ByName(name)
				if field == nil {
					// case.custom{%name}; no such field !
					continue
				}
				nested = append(nested, field.Name())
			}
			// ?fields=custom{..} requested ?
			if ok || len(nested) > 0 {
				// custom.table AS alias
				base, custom.table = custom.refer.Join(
					base, tableAlias, custom.table, "",
				)
				var scan func(custompgx.RecordExtendable) sql.Scanner
				base, scan, err = custom.refer.Columns(
					base, custom.table, nested...,
				)
				if err != nil {
					// failed to build query columns
					return sq.SelectBuilder{}, nil, err
				}
				plan = append(plan, func(row *_go.Case) any {
					return scan(custompgx.ProtoExtendable(row))
				})
			}
		}
		// continue
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
		Where(fmt.Sprintf("%s = %s.id", storeutils.Ident("rc", "primary_case_id"), caseAlias)), nil
}

func AddSubqueryAsColumn(mainQuery sq.SelectBuilder, subquery sq.SelectBuilder, subAlias string, filtersApplied bool) sq.SelectBuilder {
	if filtersApplied {
		subquery = subquery.Prefix("LATERAL (SELECT ARRAY(SELECT (subq) FROM (").Suffix(fmt.Sprintf(") subq) %s) %[1]s ON array_length(%[1]s.%[1]s, 1) > 0", subAlias))
		query, args, _ := subquery.ToSql()
		mainQuery = mainQuery.Join(query, args...)
	} else {
		subquery = subquery.Prefix("LATERAL (SELECT ARRAY(SELECT (subq) FROM (").Suffix(fmt.Sprintf(") subq) %s) %[1]s ON true", subAlias))
		query, args, _ := subquery.ToSql()
		mainQuery = mainQuery.LeftJoin(query, args...)
	}
	mainQuery = mainQuery.Column(subAlias + "::text")

	return mainQuery
}

func (c *CaseStore) GetRolesById(
	ctx context.Context,
	caseId int64,
	access auth.AccessMode,
) ([]int64, error) {

	db, err := c.storage.Database()
	if err != nil {
		return nil, err
	}
	//// Establish database connection
	//query := "(SELECT ARRAY_AGG(DISTINCT subject) rbac_r FROM cases.case_acl WHERE object = ? AND access & ? = ?)"
	query := sq.Select("ARRAY_AGG(DISTINCT subject)").From("cases.case_acl").Where("object = ?", caseId).Where("access & ? = ?", uint8(access), uint8(access)).PlaceholderFormat(sq.Dollar)
	q, args, _ := query.ToSql()
	row := db.QueryRow(ctx, q, args...)

	var res []int64
	defErr := row.Scan(&res)
	if defErr != nil {
		return nil, defErr
	}

	return res, nil
}

// Helper function to convert the scan plan to arguments for scanning.
func convertToCaseScanArgs(plan []func(caseItem *_go.Case) any, caseItem *_go.Case) []any {
	var scanArgs []any
	for _, scan := range plan {
		scanArgs = append(scanArgs, scan(caseItem))
	}
	return scanArgs
}

func NewCaseStore(store *Store) (store.CaseStore, error) {
	if store == nil {
		return nil, dberr.NewDBError("postgres.new_case.check.bad_arguments",
			"error creating case interface to the case table, main store is nil")
	}
	return &CaseStore{storage: store, mainTable: "cases.case"}, nil
}

func addCaseRbacCondition(auth auth.Auther, access auth.AccessMode, query sq.SelectBuilder, dependencyColumn string) (sq.SelectBuilder, error) {
	if auth != nil && auth.IsRbacCheckRequired(model.ScopeCases, access) {
		return query.Where(sq.Expr(fmt.Sprintf("EXISTS(SELECT acl.object FROM cases.case_acl acl WHERE acl.dc = ? AND acl.object = %s AND acl.subject = any( ?::int[]) AND acl.access & ? = ? LIMIT 1)", dependencyColumn),
			auth.GetDomainId(), pq.Array(auth.GetRoles()), int64(access), int64(access))), nil

	}
	return query, nil
}

func addCaseRbacConditionForDelete(auth auth.Auther, access auth.AccessMode, query sq.DeleteBuilder, dependencyColumn string) (sq.DeleteBuilder, error) {
	if auth != nil && auth.IsRbacCheckRequired(model.ScopeCases, access) {
		return query.Where(sq.Expr(fmt.Sprintf("EXISTS(SELECT acl.object FROM cases.case_acl acl WHERE acl.dc = ? AND acl.object = %s AND acl.subject = any( ?::int[]) AND acl.access & ? = ? LIMIT 1)", dependencyColumn),
			auth.GetDomainId(), pq.Array(auth.GetRoles()), int64(access), int64(access))), nil

	}
	return query, nil
}

func addCaseRbacConditionForUpdate(auth auth.Auther, access auth.AccessMode, query sq.UpdateBuilder, dependencyColumn string) (sq.UpdateBuilder, error) {
	if auth != nil && auth.IsRbacCheckRequired(model.ScopeCases, access) {
		return query.Where(sq.Expr(fmt.Sprintf("EXISTS(SELECT acl.object FROM cases.case_acl acl WHERE acl.dc = ? AND acl.object = %s AND acl.subject = any( ?::int[]) AND acl.access & ? = ? LIMIT 1)", dependencyColumn),
			auth.GetDomainId(), pq.Array(auth.GetRoles()), int64(access), int64(access))), nil

	}
	return query, nil
}
func (c *CaseStore) scanCases(rows pgx.Rows, plan []func(link *_go.Case) any) ([]*_go.Case, error) {
	var res []*_go.Case

	for rows.Next() {
		link, err := c.scanCase(pgx.Row(rows), plan)
		if err != nil {
			return nil, err
		}
		res = append(res, link)
	}
	return res, nil
}

func (c *CaseStore) scanCase(row pgx.Row, plan []func(link *_go.Case) any) (*_go.Case, error) {
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
