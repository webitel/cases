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
	"github.com/webitel/cases/internal/store/postgres/scanner"
	"github.com/webitel/cases/internal/store/postgres/transaction"
	"github.com/webitel/cases/model"
	util "github.com/webitel/cases/util"
)

type CaseStore struct {
	storage   store.Store
	mainTable string
}

const (
	caseLeft                  = "c"
	caseDefaultSort           = "created_at"
	casesObjClassScopeName    = model.ScopeCases
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
	caseRelatedAlias          = "related"
	caseLinksAlias            = "links"
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
	txManager := transaction.NewTxManager(tx)

	// Scan SLA details
	// Sla_id
	// reaction_at & resolve_at in [milli]seconds
	slaID, slaConditionID, reactionAt, resolveAt, calendarID, err := c.ScanSla(
		rpc,
		txManager,
		add.Service.GetId(),
		add.Priority.GetId(),
	)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.case.create.scan_sla_error", err)
	}

	// Calculate planned times within the transaction
	err = c.calculatePlannedReactionAndResolutionTime(rpc, calendarID, reactionAt, resolveAt, txManager, add)
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
	txManager *transaction.TxManager,
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
    service_hierarchy AS (
        -- Start with the service with id = $1
        SELECT id, root_id, sla_id, 1 AS level
        FROM cases.service_catalog
        WHERE id = $1  -- Start from the service with id = $1

        UNION ALL

        -- Recursively find parent services based on root_id
        SELECT sc.id, sc.root_id, COALESCE(sc.sla_id, sh.sla_id) AS sla_id, sh.level + 1
        FROM cases.service_catalog sc
        INNER JOIN service_hierarchy sh ON sc.id = sh.root_id  -- Join on root_id to get parents
    ),
    deepest_service AS (
        -- Select the SLA ID and its associated details from the deepest service
        SELECT sla_id, MAX(level) AS max_level
        FROM service_hierarchy
        WHERE sla_id IS NOT NULL  -- Filter out rows where sla_id is null
        GROUP BY sla_id
        ORDER BY max_level ASC
        LIMIT 1  -- Only return the SLA ID of the deepest service
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
-- Final SELECT for deepest service and its priority condition details with COALESCE
SELECT ds.sla_id,
       COALESCE(pc.reaction_time, sla.reaction_time) AS reaction_time,
       COALESCE(pc.resolution_time, sla.resolution_time) AS resolution_time,
       sla.calendar_id,
       pc.sla_condition_id
FROM deepest_service ds
LEFT JOIN priority_condition pc ON true
LEFT JOIN cases.sla sla ON ds.sla_id = sla.id;

	`, serviceID, priorityID).Scan(
		scanner.ScanInt(&slaID),
		scanner.ScanInt(&reactionTime),
		scanner.ScanInt(&resolutionTime),
		scanner.ScanInt(&calendarID),
		scanner.ScanInt(&slaConditionID),
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
		assignee, closeReason, reporter *int64
		closeResult, description        *string
	)

	if id := caseItem.GetClose().GetCloseReason().GetId(); id > 0 {
		closeReason = &id
	}

	if result := caseItem.GetClose().GetCloseResult(); result != "" {
		closeResult = &result
	}

	if id := caseItem.Reporter.GetId(); id > 0 {
		reporter = &id
	}

	if id := caseItem.Assignee.GetId(); id > 0 {
		assignee = &id
	}

	if desc := caseItem.Description; desc != "" {
		description = &desc
	}
	params := map[string]any{
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
		"description":         description,
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
		&model.SearchOptions{
			Context: rpc.Context,
			Fields:  rpc.Fields,
		},
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

func (c *CaseStore) calculatePlannedReactionAndResolutionTime(
	rpc *model.CreateOptions,
	calendarID int,
	reactionTime int,
	resolutionTime int,
	txManager *transaction.TxManager,
	caseItem *_go.Case,
) error {
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
	err = txManager.QueryRow(rpc.Context, `
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

	// Get the current time
	currentTime := rpc.CurrentTime()

	// Calculate planned reaction and resolution timestamps
	reactionTimestamp, err := calculateTimestampFromCalendar(currentTime, offset, reactionMinutes, mergedSlots)
	if err != nil {
		return fmt.Errorf("failed to calculate planned reaction time: %w", err)
	}

	resolveTimestamp, err := calculateTimestampFromCalendar(currentTime, offset, resolutionMinutes, mergedSlots)
	if err != nil {
		return fmt.Errorf("failed to calculate planned resolution time: %w", err)
	}

	caseItem.PlannedReactionAt = reactionTimestamp.UnixMilli()
	caseItem.PlannedResolveAt = resolveTimestamp.UnixMilli()

	return nil
}

// fetchCalendarSlots retrieves working hours for a calendar
func fetchCalendarSlots(rpc *model.CreateOptions, txManager *transaction.TxManager, calendarID int) ([]CalendarSlot, error) {
	rows, err := txManager.Query(rpc.Context, `
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
func fetchExceptionSlots(rpc *model.CreateOptions, txManager *transaction.TxManager, calendarID int) ([]ExceptionSlot, error) {
	rows, err := txManager.Query(rpc.Context, `
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
		adjustedDay := (cal.Day % 7) // Adjust to make sure it's in [0, 6] range
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
	opts.Fields = util.EnsureFields(opts.Fields, "id")
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
			"status",           // +
			"impacted",         // +
			"close_reason",     // +
			"contact_group",    // +
			"service",          // +
			"status_condition", // +
			"sla":              // +
			if value == "null" {
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
		case "author":
			if value == "null" {
				base = base.Where(fmt.Sprintf("%s IS NULL", store.Ident(caseAuthorAlias, "id")))
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
				base = base.Where(fmt.Sprintf("%s = ANY(?::int[])", store.Ident(caseAuthorAlias, "id")), valuesInt)
			}
		// Filter for the date range (created_at)
		case "created_at.from", "created_at.to":
			// Check if both from and to are provided, create the range filter
			fromValue, hasFrom := opts.Filter["created_at.from"]
			toValue, hasTo := opts.Filter["created_at.to"]
			if hasFrom && hasTo {
				// Apply range filtering using both `from` and `to` values
				base = base.Where(fmt.Sprintf("%s >= ?::timestamp AND c.created_at <= ?::timestamp", store.Ident(caseLeft, "created_at")), fromValue, toValue)
			} else if hasFrom {
				// Only "from" filter is provided
				base = base.Where(fmt.Sprintf("%s >= ?::timestamp", store.Ident(caseLeft, "created_at")), fromValue)
			} else if hasTo {
				// Only "to" filter is provided
				base = base.Where(fmt.Sprintf("%s <= ?::timestamp", store.Ident(caseLeft, "created_at")), toValue)
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
	field, direction := store.GetSortingOperator(opts)
	if field == "" {
		field = caseDefaultSort
		direction = "DESC"
	}
	var tableAlias string
	if !util.ContainsStringIgnoreCase(opts.Fields, field) { // not joined yet
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
		}
	}
	if tableAlias == "" {
		tableAlias = caseLeft
	}
	switch field {
	case "id", "ver", "created_at", "updated_at", "name", "subject", "description", "planned_reaction_at", "planned_resolve_at", "contact_info":
		base = base.OrderBy(fmt.Sprintf("%s %s", store.Ident(tableAlias, field), direction))
	case "created_by", "updated_by", "source", "close_reason_group", "sla", "status_condition", "status", "priority", "service":
		base = base.OrderBy(fmt.Sprintf("%s %s", store.Ident(tableAlias, "name"), direction))
	case "author", "assignee", "reporter", "impacted":
		base = base.OrderBy(fmt.Sprintf("%s %s", store.Ident(tableAlias, "common_name"), direction))
	}

	return base, plan, nil
}

// region UPDATE
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
		return nil, dberr.NewDBInternalError("postgres.case.create.transaction_error", txErr)
	}
	defer tx.Rollback(rpc.Context)
	txManager := transaction.NewTxManager(tx)

	// * if user change Service -- SLA ; SLA Condition ; Planned Reaction / Resolve at ; Calendar could be changed
	if util.ContainsField(rpc.Mask, "service") {
		slaID, slaConditionID, reaction_at, resolve_at, calendarID, err := c.ScanSla(
			&model.CreateOptions{Context: rpc.Context},
			txManager,
			upd.Service.GetId(),
			upd.Priority.GetId(),
		)
		if err != nil {
			return nil, dberr.NewDBInternalError("postgres.case.update.scan_sla_error", err)
		}

		// Calculate planned times within the transaction
		err = c.calculatePlannedReactionAndResolutionTime(
			&model.CreateOptions{
				Context: rpc.Context,
				Time:    rpc.CurrentTime(),
			},
			calendarID,
			reaction_at,
			resolve_at,
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
		upd.Sla.Id = int64(slaID)
		if upd.SlaCondition == nil {
			upd.SlaCondition = &_go.Lookup{}
		}
		upd.SlaCondition.Id = int64(slaConditionID)
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

	if err := txManager.QueryRow(rpc.Context, query, args...).Scan(scanArgs...); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, dberr.NewDBNoRowsError("postgres.case.update.update.scan_ver.not_found")
		}
		return nil, dberr.NewDBInternalError("postgres.case.update.update.execution_error", err)
	}

	// Commit the transaction
	if err := tx.Commit(rpc.Context); err != nil {
		return nil, dberr.NewDBInternalError("postgres.case.update.commit_error", err)
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
			updateBuilder = updateBuilder.Set("sla", upd.Sla.GetId())
			updateBuilder = updateBuilder.Set("sla_condition_id", upd.SlaCondition.GetId())
			updateBuilder = updateBuilder.Set("planned_resolve_at", util.LocalTime(upd.GetPlannedResolveAt()))
			updateBuilder = updateBuilder.Set("planned_reaction_at", util.LocalTime(upd.GetPlannedReactionAt()))
		case "assignee":
			if upd.Assignee.GetId() == 0 {
				updateBuilder = updateBuilder.Set("assignee", nil)
			} else {
				updateBuilder = updateBuilder.Set("assignee", upd.Assignee.GetId())
			}
		case "reporter":
			updateBuilder = updateBuilder.Set("reporter", upd.Reporter.GetId())
		case "contact_info":
			updateBuilder = updateBuilder.Set("contact_info", upd.GetContactInfo())
		case "impacted":
			updateBuilder = updateBuilder.Set("impacted", upd.Impacted.GetId())
		case "group":
			if upd.Group.GetId() == 0 {
				updateBuilder = updateBuilder.Set("contact_group", nil)
			} else {
				updateBuilder = updateBuilder.Set("contact_group", upd.Group.GetId())
			}
		case "close":
			if upd.Close != nil && upd.Close.CloseReason != nil {
				var closeReason *int64
				if reas := upd.Close.CloseReason.GetId(); reas > 0 {
					closeReason = &reas
				}
				updateBuilder = updateBuilder.Set("close_reason", closeReason)
				updateBuilder = updateBuilder.Set("close_result", upd.Close.GetCloseResult())
			}
		case "rate":
			if upd.Rate != nil {
				updateBuilder = updateBuilder.Set("rating", upd.Rate.Rating)
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
		joinTable(tableAlias, "directory.wbt_user", store.Ident(caseLeft, "created_by"))
	case "updated_by":
		tableAlias = caseUpdatedByAlias
		joinTable(tableAlias, "directory.wbt_user", store.Ident(caseLeft, "updated_by"))
	case "source":
		tableAlias = caseSourceAlias
		joinTable(tableAlias, "cases.source", store.Ident(caseLeft, "source"))
	case "close_reason_group":
		tableAlias = caseCloseReasonGroupAlias
		joinTable(tableAlias, "cases.close_reason_group", store.Ident(caseLeft, "close_reason_group"))
	case "author":
		createdByAlias := "cb_au"
		tableAlias = caseAuthorAlias
		joinTable(createdByAlias, "directory.wbt_user", store.Ident(caseLeft, "created_by"))
		joinTable(tableAlias, "contacts.contact", store.Ident(createdByAlias, "contact_id"))
	case "close":
		tableAlias = caseCloseReasonAlias
		joinTable(tableAlias, "cases.close_reason", store.Ident(caseLeft, "close_reason"))
	case "sla":
		tableAlias = caseSlaAlias
		joinTable(tableAlias, "cases.sla", store.Ident(caseLeft, "sla"))
	case "status":
		tableAlias = caseStatusAlias
		joinTable(tableAlias, "cases.status", store.Ident(caseLeft, "status"))
	case "priority":
		tableAlias = casePriorityAlias
		joinTable(tableAlias, "cases.priority", store.Ident(caseLeft, "priority"))
	case "service":
		tableAlias = caseServiceAlias
		joinTable(tableAlias, "cases.service_catalog", store.Ident(caseLeft, "service"))
	case "assignee":
		tableAlias = caseAssigneeAlias
		joinTable(tableAlias, "contacts.contact", store.Ident(caseLeft, "assignee"))
	case "reporter":
		tableAlias = caseReporterAlias
		joinTable(tableAlias, "contacts.contact", store.Ident(caseLeft, "reporter"))
	case "impacted":
		tableAlias = caseImpactedAlias
		joinTable(tableAlias, "contacts.contact", store.Ident(caseLeft, "impacted"))
	case "group":
		tableAlias = caseGroupAlias
		joinTable(tableAlias, "contacts.group", store.Ident(caseLeft, "contact_group"))
	}
	return base, tableAlias, nil
}

// session required to get some columns
func (c *CaseStore) buildCaseSelectColumnsAndPlan(opts *model.SearchOptions,
	base sq.SelectBuilder,
) (sq.SelectBuilder, []func(caseItem *_go.Case) any, error) {
	var (
		plan       []func(caseItem *_go.Case) any
		tableAlias string
		err        error
	)

	for _, field := range opts.Fields {
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
			base = base.Column(store.Ident(tableAlias, "id AS case_id"))
			plan = append(plan, func(caseItem *_go.Case) any {
				return &caseItem.Id
			})
		case "ver":
			base = base.Column(store.Ident(tableAlias, "ver"))
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
			base = base.Column(store.Ident(tableAlias, "created_at"))
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
			base = base.Column(store.Ident(tableAlias, "updated_at"))
			plan = append(plan, func(caseItem *_go.Case) any {
				return scanner.ScanTimestamp(&caseItem.UpdatedAt)
			})
		case "name":
			base = base.Column(store.Ident(tableAlias, "name"))
			plan = append(plan, func(caseItem *_go.Case) any {
				return &caseItem.Name
			})
		case "subject":
			base = base.Column(store.Ident(tableAlias, "subject"))
			plan = append(plan, func(caseItem *_go.Case) any {
				return &caseItem.Subject
			})
		case "description":
			base = base.Column(store.Ident(tableAlias, "description"))
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
				"ROW(%s.id, %[1]s.name)::text  AS close_reason_group", tableAlias))
			plan = append(plan, func(caseItem *_go.Case) any {
				return scanner.ScanRowLookup(&caseItem.CloseReasonGroup)
			})
		case "author":
			base = base.Column(fmt.Sprintf(`ROW(%s.id, %[1]s.common_name)::text AS author`, tableAlias))
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
				"ROW(%s.id, %[1]s.name)::text AS close_reason", tableAlias))
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
					"COALESCE(CAST(EXTRACT(EPOCH FROM %s.reacted_at - %[1]s.created_at) * 1000 AS bigint), 0) AS difference_in_reaction",
					caseLeft,
				)).
				Column(fmt.Sprintf(
					"COALESCE(CAST(EXTRACT(EPOCH FROM %s.resolved_at - %[1]s.created_at) * 1000 AS bigint), 0) AS difference_in_resolve",
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
			base = base.Column(store.Ident(caseLeft, field))
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
					Sort:    "-created_at",
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
					Sort:    "-created_at",
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
					Sort: "-created_at",
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
