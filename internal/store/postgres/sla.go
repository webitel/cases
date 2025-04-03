package postgres

import (
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/webitel/cases/api/cases"
	_go "github.com/webitel/cases/api/cases"
	dberr "github.com/webitel/cases/internal/errors"
	"github.com/webitel/cases/internal/store"
	"github.com/webitel/cases/internal/store/postgres/scanner"
	"github.com/webitel/cases/internal/store/postgres/transaction"
	util2 "github.com/webitel/cases/internal/store/util"
	"github.com/webitel/cases/model/options"
	"github.com/webitel/cases/util"
	"regexp"
	"strconv"
	"time"
)

type SLAScan func(sla *cases.SLA) any

const (
	slaLeft        = "s"
	slaDefaultSort = "name"
)

type SLAStore struct {
	storage *Store
}

// Helper function to convert plan to scan arguments.
func convertToSLAScanArgs(plan []SLAScan, sla *cases.SLA) []any {
	var scanArgs []any
	for _, scan := range plan {
		scanArgs = append(scanArgs, scan(sla))
	}
	return scanArgs
}

// Helper function to dynamically build select columns and plan.
func (s *SLAStore) buildSLASelectColumnsAndPlan(
	base sq.SelectBuilder,
	fields []string,
) (sq.SelectBuilder, []SLAScan, error) {
	var plan []SLAScan
	for _, field := range fields {
		switch field {
		case "id":
			base = base.Column(util2.Ident(slaLeft, "id"))
			plan = append(plan, func(sla *cases.SLA) any {
				return &sla.Id
			})
		case "name":
			base = base.Column(util2.Ident(slaLeft, "name"))
			plan = append(plan, func(sla *cases.SLA) any {
				return &sla.Name
			})
		case "description":
			base = base.Column(util2.Ident(slaLeft, "description"))
			plan = append(plan, func(sla *cases.SLA) any {
				return scanner.ScanText(&sla.Description)
			})
		case "valid_from":
			base = base.Column(util2.Ident(slaLeft, "valid_from"))
			plan = append(plan, func(sla *cases.SLA) any {
				return scanner.ScanTimestamp(&sla.ValidFrom)
			})
		case "valid_to":
			base = base.Column(util2.Ident(slaLeft, "valid_to"))
			plan = append(plan, func(sla *cases.SLA) any {
				return scanner.ScanTimestamp(&sla.ValidTo)
			})
		case "calendar":
			base = base.Column(
				fmt.Sprintf(
					"(SELECT ROW(id, name)::text FROM flow.calendar WHERE id = %s.calendar_id) calendar", slaLeft))
			plan = append(plan, func(sla *cases.SLA) any {
				return scanner.ScanRowLookup(&sla.Calendar)
			})
		case "reaction_time":
			base = base.
				Column(util2.Ident(slaLeft, "reaction_time")).
				Column(util2.Ident(slaLeft, "reaction_duration"))

			plan = append(plan,
				func(sla *cases.SLA) any {
					return &sla.ReactionTimeMillis
				},
				func(sla *cases.SLA) any {
					return scanner.ScanTimingsString(&sla.ReactionTime)
				})

		case "resolution_time":
			base = base.
				Column(util2.Ident(slaLeft, "resolution_time")).
				Column(util2.Ident(slaLeft, "resolution_duration"))

			plan = append(plan,
				func(sla *cases.SLA) any {
					return &sla.ResolutionTimeMillis
				},
				func(sla *cases.SLA) any {
					return scanner.ScanTimingsString(&sla.ResolutionTime)
				})

		case "created_at":
			base = base.Column(util2.Ident(slaLeft, "created_at"))
			plan = append(plan, func(sla *cases.SLA) any {
				return scanner.ScanTimestamp(&sla.CreatedAt)
			})
		case "updated_at":
			base = base.Column(util2.Ident(slaLeft, "updated_at"))
			plan = append(plan, func(sla *cases.SLA) any {
				return scanner.ScanTimestamp(&sla.UpdatedAt)
			})
		case "created_by":
			base = base.Column(fmt.Sprintf("(SELECT ROW(id, name)::text FROM directory.wbt_user WHERE id = %s.created_by) created_by", slaLeft))
			plan = append(plan, func(sla *cases.SLA) any {
				return scanner.ScanRowLookup(&sla.CreatedBy)
			})
		case "updated_by":
			base = base.Column(fmt.Sprintf("(SELECT ROW(id, name)::text FROM directory.wbt_user WHERE id = %s.updated_by) updated_by", slaLeft))
			plan = append(plan, func(sla *cases.SLA) any {
				return scanner.ScanRowLookup(&sla.UpdatedBy)
			})
		default:
			return base, nil, dberr.NewDBInternalError("postgres.sla.unknown_field", fmt.Errorf("unknown field: %s", field))
		}
	}
	return base, plan, nil
}

func (s *SLAStore) buildCreateSLAQuery(
	rpc options.CreateOptions,
	sla *cases.SLA,
	reactionTimeMillis int64,
	resolutionTimeMillis int64,
	txManager *transaction.TxManager,
) (sq.SelectBuilder, []SLAScan, error) {
	fields := rpc.GetFields()
	fields = util.EnsureIdField(rpc.GetFields())
	// Build the INSERT query with a RETURNING clause
	insertBuilder := sq.Insert("cases.sla").
		Columns(
			"name", "dc", "created_at",
			"description", "created_by", "updated_at",
			"updated_by", "valid_from", "valid_to",
			"calendar_id", "reaction_time", "resolution_time",
			"reaction_duration", "resolution_duration",
		).
		Values(
			sla.Name,
			rpc.GetAuthOpts().GetDomainId(),
			rpc.RequestTime(),
			sq.Expr("NULLIF(?, '')", sla.Description),
			rpc.GetAuthOpts().GetUserId(),
			rpc.RequestTime(),
			rpc.GetAuthOpts().GetUserId(),
			util.LocalTime(sla.ValidFrom),
			util.LocalTime(sla.ValidTo),
			sla.Calendar.Id,
			reactionTimeMillis,
			resolutionTimeMillis,
			FormatTimings(sla.ReactionTime),
			FormatTimings(sla.ResolutionTime),
		).
		PlaceholderFormat(sq.Dollar).
		Suffix("RETURNING *") // RETURNING all columns for use in the next SELECT

	// Convert the INSERT query into a CTE
	insertSQL, args, err := insertBuilder.ToSql()
	if err != nil {
		return sq.SelectBuilder{}, nil, dberr.NewDBInternalError("postgres.sla.create.query_build_error", err)
	}

	// Use the INSERT query as a CTE (Common Table Expression)
	cte := sq.Expr("WITH s AS ("+insertSQL+")", args...)

	// Dynamically build the SELECT query for the resulting row
	selectBuilder, plan, err := s.buildSLASelectColumnsAndPlan(sq.Select(), fields)
	if err != nil {
		return sq.SelectBuilder{}, nil, err
	}

	// Combine the CTE with the SELECT query
	selectBuilder = selectBuilder.PrefixExpr(cte).From(slaLeft)

	return selectBuilder, plan, nil
}

func (s *SLAStore) Create(rpc options.CreateOptions, input *cases.SLA) (*cases.SLA, error) {
	db, dbErr := s.storage.Database()
	if dbErr != nil {
		return nil, dberr.NewDBInternalError("postgres.sla.create.database_connection_error", dbErr)
	}

	tx, err := db.Begin(rpc)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.sla.create.begin_tx_error", err)
	}
	txManager := transaction.NewTxManager(tx)

	defer func() {
		if err != nil {
			_ = txManager.Rollback(rpc)
		} else {
			_ = txManager.Commit(rpc)
		}
	}()

	var (
		resolutionTimeMillis int64
		reactionTimeMillis   int64
	)

	reactionTimeMillis, err = s.CalculateCalendarMillis(
		rpc,
		txManager,
		int(input.Calendar.GetId()),
		input.ReactionTime,
	)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.sla.create.reaction_time_calc_error", err)
	}

	resolutionTimeMillis, err = s.CalculateCalendarMillis(
		rpc,
		txManager,
		int(input.Calendar.GetId()),
		input.ResolutionTime,
	)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.sla.create.resolution_time_calc_error", err)
	}

	selectBuilder, plan, err := s.buildCreateSLAQuery(
		rpc,
		input,
		reactionTimeMillis,
		resolutionTimeMillis,
		txManager,
	)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.sla.create.build_query_error", err)
	}

	query, args, err := selectBuilder.ToSql()
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.sla.create.query_to_sql_error", err)
	}

	scanArgs := convertToSLAScanArgs(plan, input)
	if err := txManager.QueryRow(rpc, query, args...).Scan(scanArgs...); err != nil {
		return nil, dberr.NewDBInternalError("postgres.sla.create.execution_error", err)
	}

	return input, nil
}

func FormatTimings(t *cases.Timings) string {
	if t == nil {
		return ""
	}
	str := ""
	if t.Dd > 0 {
		str += fmt.Sprintf("%dd", t.Dd)
	}
	if t.Hh > 0 {
		str += fmt.Sprintf("%dh", t.Hh)
	}
	if t.Mm > 0 {
		str += fmt.Sprintf("%dm", t.Mm)
	}
	return str
}

var re = regexp.MustCompile(`(?:(\d+)d)?(?:(\d+)h)?(?:(\d+)m)?`)

func ParseTimings(s string) *cases.Timings {
	if s == "" {
		return &cases.Timings{}
	}
	matches := re.FindStringSubmatch(s)
	t := &cases.Timings{}
	if len(matches) == 4 {
		if matches[1] != "" {
			t.Dd, _ = strconv.ParseInt(matches[1], 10, 64)
		}
		if matches[2] != "" {
			t.Hh, _ = strconv.ParseInt(matches[2], 10, 64)
		}
		if matches[3] != "" {
			t.Mm, _ = strconv.ParseInt(matches[3], 10, 64)
		}
	}
	return t
}

func (s *SLAStore) getCalendarConfig(
	rpc TimingOpts,
	txManager *transaction.TxManager,
	calendarID int,
) ([]MergedSlot, time.Duration, error) {
	// Fetch calendar and exceptions
	calendar, err := fetchCalendarSlots(rpc, txManager, calendarID)
	if err != nil {
		return nil, 0, err
	}
	exceptions, err := fetchExceptionSlots(rpc, txManager, calendarID)
	if err != nil {
		return nil, 0, err
	}
	merged := mergeCalendarAndExceptions(calendar, exceptions)

	// Fetch timezone offset
	var offset time.Duration
	err = txManager.QueryRow(rpc, `
		SELECT tz.utc_offset
		FROM flow.calendar cl
		    LEFT JOIN flow.calendar_timezones tz ON tz.id = cl.timezone_id
		WHERE cl.id = $1`, calendarID).Scan(&offset)
	if err != nil {
		return nil, 0, err
	}

	return merged, offset, nil
}

func (s *SLAStore) CalculateCalendarMillis(
	rpc TimingOpts,
	txManager *transaction.TxManager,
	calendarID int,
	t *_go.Timings,
) (int64, error) {
	if t == nil {
		return 0, nil
	}

	merged, offset, err := s.getCalendarConfig(rpc, txManager, calendarID)
	if err != nil {
		return 0, err
	}

	now := rpc.RequestTime()
	current := now
	remainingDays := int(t.Dd)
	totalWorkingMinutes := 0

	for remainingDays > 0 {
		dayMatched := false

		for _, slot := range merged {
			if slot.Disabled {
				continue
			}
			if !slot.Date.IsZero() && !isSameDate(current, slot.Date) {
				continue
			}
			if slot.Date.IsZero() && int(current.Weekday()) != slot.Day {
				continue
			}

			// Compute slot range
			slotStart := time.Date(current.Year(), current.Month(), current.Day(), 0, 0, 0, 0, current.Location()).
				Add(time.Minute * time.Duration(slot.StartTimeOfDay-int(offset.Minutes())))
			slotEnd := time.Date(current.Year(), current.Month(), current.Day(), 0, 0, 0, 0, current.Location()).
				Add(time.Minute * time.Duration(slot.EndTimeOfDay-int(offset.Minutes())))

			// Skip expired slots on the first day
			if current.Year() == now.Year() && current.YearDay() == now.YearDay() && slotEnd.Before(now) {
				continue
			}

			// If this is the first day and now is within slot → don't trim start
			// because Dd means "full calendar working day"
			start := slotStart
			end := slotEnd

			if end.After(start) {
				diff := int(end.Sub(start).Minutes())
				totalWorkingMinutes += diff
				dayMatched = true
			}
		}

		if dayMatched {
			remainingDays--
		}

		current = current.AddDate(0, 0, 1)
	}

	// Step 2: Add Hh + Mm directly as working minutes
	added := int(t.Hh*60 + t.Mm)
	totalWorkingMinutes += added

	return int64(totalWorkingMinutes) * 60_000, nil
}

func (s *SLAStore) buildUpdateSLAQuery(
	rpc options.UpdateOptions,
	sla *cases.SLA,
	reactionTimeMillis int64,
	resolutionTimeMillis int64,
	txManager *transaction.TxManager,
) (sq.SelectBuilder, []SLAScan, error) {
	fields := rpc.GetFields()
	fields = util.EnsureIdField(rpc.GetFields())
	// Start the UPDATE query
	updateBuilder := sq.Update("cases.sla").
		PlaceholderFormat(sq.Dollar). // Use PostgreSQL-compatible placeholders
		Set("updated_at", rpc.RequestTime()).
		Set("updated_by", rpc.GetAuthOpts().GetUserId()).
		Where(sq.Eq{"id": sla.Id}).
		Where(sq.Eq{"dc": rpc.GetAuthOpts().GetDomainId()})

	// Dynamically add fields to the SET clause
	for _, field := range rpc.GetMask() {
		switch field {
		case "name":
			if sla.Name != "" {
				updateBuilder = updateBuilder.Set("name", sla.Name)
			}
		case "description":
			updateBuilder = updateBuilder.Set("description", sq.Expr("NULLIF(?, '')", sla.Description))
		case "valid_from":
			updateBuilder = updateBuilder.Set("valid_from", util.LocalTime(sla.ValidFrom))
		case "valid_to":
			updateBuilder = updateBuilder.Set("valid_to", util.LocalTime(sla.ValidTo))
		case "calendar_id":
			if sla.Calendar.Id != 0 {
				updateBuilder = updateBuilder.Set("calendar_id", sla.Calendar.Id)
			}
		case "reaction_time":
			updateBuilder = updateBuilder.
				Set("reaction_time", reactionTimeMillis).
				Set("reaction_duration", FormatTimings(sla.ReactionTime))
		case "resolution_time":
			updateBuilder = updateBuilder.
				Set("resolution_time", resolutionTimeMillis).
				Set("resolution_duration", FormatTimings(sla.ResolutionTime))
		}
	}

	// Generate the CTE for the update operation
	updateSQL, args, err := updateBuilder.Suffix("RETURNING *").ToSql()
	if err != nil {
		return sq.SelectBuilder{}, nil, dberr.NewDBInternalError("postgres.sla.update.query_build_error", err)
	}

	// Use the UPDATE query as a CTE
	cte := sq.Expr("WITH s AS ("+updateSQL+")", args...)

	// Build select clause and scan plan dynamically using buildSLASelectColumnsAndPlan
	selectBuilder, plan, err := s.buildSLASelectColumnsAndPlan(sq.Select(), fields)
	if err != nil {
		return sq.SelectBuilder{}, nil, err
	}

	// Combine the CTE with the SELECT query
	selectBuilder = selectBuilder.PrefixExpr(cte).From("s")

	return selectBuilder, plan, nil
}

func (s *SLAStore) Update(rpc options.UpdateOptions, input *cases.SLA) (*cases.SLA, error) {
	db, dbErr := s.storage.Database()
	if dbErr != nil {
		return nil, dberr.NewDBInternalError("postgres.sla.input.database_connection_error", dbErr)
	}

	tx, err := db.Begin(rpc)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.sla.input.begin_tx_error", err)
	}
	txManager := transaction.NewTxManager(tx)

	defer func() {
		if err != nil {
			_ = txManager.Rollback(rpc)
		} else {
			_ = txManager.Commit(rpc)
		}
	}()

	var (
		reactionTimeMillis   int64
		resolutionTimeMillis int64
	)

	if util.ContainsField(rpc.GetMask(), "reaction_time") {
		reactionTimeMillis, err = s.CalculateCalendarMillis(
			rpc,
			txManager,
			int(input.Calendar.GetId()),
			input.ResolutionTime,
		)
		if err != nil {
			return nil, dberr.NewDBInternalError("postgres.sla.input.reaction_time_calc_error", err)
		}
	}

	if util.ContainsField(rpc.GetMask(), "resolution_time") {
		resolutionTimeMillis, err = s.CalculateCalendarMillis(
			rpc,
			txManager,
			int(input.Calendar.GetId()),
			input.ReactionTime,
		)
		if err != nil {
			return nil, dberr.NewDBInternalError("postgres.sla.input.resolution_time_calc_error", err)
		}
	}

	selectBuilder, plan, err := s.buildUpdateSLAQuery(
		rpc,
		input,
		reactionTimeMillis,
		resolutionTimeMillis,
		txManager,
	)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.sla.input.build_query_error", err)
	}

	query, args, err := selectBuilder.ToSql()
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.sla.input.query_to_sql_error", err)
	}

	scanArgs := convertToSLAScanArgs(plan, input)
	if err := txManager.QueryRow(rpc, query, args...).Scan(scanArgs...); err != nil {
		return nil, dberr.NewDBInternalError("postgres.sla.input.execution_error", err)
	}

	return input, nil
}

func (s *SLAStore) buildListSLAQuery(
	rpc options.SearchOptions,
	txManager *transaction.TxManager,
) (sq.SelectBuilder, []SLAScan, error) {

	queryBuilder := sq.Select().
		From("cases.sla AS s").
		Where(sq.Eq{"s.dc": rpc.GetAuthOpts().GetDomainId()}).
		PlaceholderFormat(sq.Dollar)

	// Add ID filter if provided
	if len(rpc.GetIDs()) > 0 {
		queryBuilder = queryBuilder.Where(sq.Eq{"s.id": rpc.GetIDs()})
	}

	// Add name filter if provided
	if name, ok := rpc.GetFilter("name").(string); ok && len(name) > 0 {
		queryBuilder = util2.AddSearchTerm(queryBuilder, name, "s.name")
	}

	// -------- Apply sorting ----------
	queryBuilder = util2.ApplyDefaultSorting(rpc, queryBuilder, slaDefaultSort)

	// ---------Apply paging based on Search Opts ( page ; size ) -----------------
	queryBuilder = util2.ApplyPaging(rpc.GetPage(), rpc.GetSize(), queryBuilder)

	// Add select columns and scan plan for requested fields
	queryBuilder, plan, err := s.buildSLASelectColumnsAndPlan(queryBuilder, rpc.GetFields())
	if err != nil {
		return sq.SelectBuilder{}, nil, dberr.NewDBInternalError("postgres.sla.search.query_build_error", err)
	}

	return queryBuilder, plan, nil
}

func (s *SLAStore) List(rpc options.SearchOptions) (*cases.SLAList, error) {
	db, dbErr := s.storage.Database()
	if dbErr != nil {
		return nil, dberr.NewDBInternalError("postgres.sla.list.database_connection_error", dbErr)
	}

	tx, err := db.Begin(rpc)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.sla.list.begin_tx_error", err)
	}
	txManager := transaction.NewTxManager(tx)
	defer func() {
		_ = txManager.Rollback(rpc)
	}()

	selectBuilder, plan, err := s.buildListSLAQuery(rpc, txManager)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.sla.list.build_query_error", err)
	}

	query, args, err := selectBuilder.ToSql()
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.sla.list.query_build_error", err)
	}
	query = util2.CompactSQL(query)

	rows, err := txManager.Query(rpc, query, args...)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.sla.list.execution_error", err)
	}
	defer rows.Close()

	var slas []*cases.SLA
	lCount := 0
	next := false
	fetchAll := rpc.GetSize() == -1

	for rows.Next() {
		sla := &cases.SLA{}
		scanArgs := convertToSLAScanArgs(plan, sla)
		if err := rows.Scan(scanArgs...); err != nil {
			return nil, dberr.NewDBInternalError("postgres.sla.list.row_scan_error", err)
		}
		slas = append(slas, sla)
		lCount++
		if !fetchAll && lCount >= rpc.GetSize() {
			next = true
			break
		}
	}
	rows.Close()

	_ = txManager.Commit(rpc)

	return &cases.SLAList{
		Page:  int32(rpc.GetPage()),
		Next:  next,
		Items: slas,
	}, nil
}

func (s *SLAStore) buildDeleteSLAQuery(
	rpc options.DeleteOptions,
) (sq.DeleteBuilder, error) {
	// Ensure IDs are provided
	if len(rpc.GetIDs()) == 0 {
		return sq.DeleteBuilder{}, dberr.NewDBInternalError("postgres.sla.delete.missing_ids", fmt.Errorf("no IDs provided for deletion"))
	}

	// Build the delete query
	deleteBuilder := sq.Delete("cases.sla").
		Where(sq.Eq{"id": rpc.GetIDs()}).
		Where(sq.Eq{"dc": rpc.GetAuthOpts().GetDomainId()}).
		PlaceholderFormat(sq.Dollar)

	return deleteBuilder, nil
}

func (s *SLAStore) Delete(rpc options.DeleteOptions) error {
	d, dbErr := s.storage.Database()
	if dbErr != nil {
		return dberr.NewDBInternalError("postgres.sla.delete.database_connection_error", dbErr)
	}

	deleteBuilder, err := s.buildDeleteSLAQuery(rpc)
	if err != nil {
		return dberr.NewDBInternalError("postgres.sla.delete.query_build_error", err)
	}

	query, args, err := deleteBuilder.ToSql()
	if err != nil {
		return dberr.NewDBInternalError("postgres.sla.delete.query_to_sql_error", err)
	}

	res, execErr := d.Exec(rpc, query, args...)
	if execErr != nil {
		return dberr.NewDBInternalError("postgres.sla.delete.execution_error", execErr)
	}

	if res.RowsAffected() == 0 {
		return dberr.NewDBNoRowsError("postgres.sla.delete.no_rows_affected")
	}

	return nil
}

func NewSLAStore(store *Store) (store.SLAStore, error) {
	if store == nil {
		return nil, dberr.NewDBError("postgres.new_sla.check.bad_arguments",
			"error creating SLA interface, main store is nil")
	}
	return &SLAStore{storage: store}, nil
}
