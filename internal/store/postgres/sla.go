package postgres

import (
	"fmt"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/lib/pq"
	cases "github.com/webitel/cases/api/cases"
	dberr "github.com/webitel/cases/internal/error"
	"github.com/webitel/cases/internal/store"
	"github.com/webitel/cases/model"
	"github.com/webitel/cases/util"
)

type SLAStore struct {
	storage store.Store
}

// Create implements store.SLAStore.
func (s *SLAStore) Create(rpc *model.CreateOptions, add *cases.SLA) (*cases.SLA, error) {
	d, dbErr := s.storage.Database()
	if dbErr != nil {
		return nil, dberr.NewDBInternalError("postgres.sla.create.database_connection_error", dbErr)
	}

	query, args, err := s.buildCreateSLAQuery(rpc, add)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.sla.create.query_build_error", err)
	}

	var (
		createdByLookup, updatedByLookup, calendar cases.Lookup
		createdAt, updatedAt                       time.Time
		validFrom, validTo                         time.Time
	)

	err = d.QueryRow(rpc.Context, query, args...).Scan(
		&add.Id, &add.Name, &createdAt, &add.Description,
		&validFrom, &validTo, &calendar.Id, &calendar.Name,
		&add.ReactionTime.Hours, &add.ReactionTime.Minutes,
		&add.ResolutionTime.Hours, &add.ResolutionTime.Minutes,
		&createdByLookup.Id, &createdByLookup.Name,
		&updatedAt, &updatedByLookup.Id, &updatedByLookup.Name,
	)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.sla.create.execution_error", err)
	}

	t := rpc.Time
	return &cases.SLA{
		Id:          add.Id,
		Name:        add.Name,
		Description: add.Description,
		ValidFrom:   util.Timestamp(validFrom),
		ValidTo:     util.Timestamp(validTo),
		Calendar:    &calendar,
		ReactionTime: &cases.ReactionTime{
			Hours:   add.ReactionTime.Hours,
			Minutes: add.ReactionTime.Minutes,
		},
		ResolutionTime: &cases.ResolutionTime{
			Hours:   add.ResolutionTime.Hours,
			Minutes: add.ResolutionTime.Minutes,
		},
		CreatedAt: util.Timestamp(t),
		UpdatedAt: util.Timestamp(t),
		CreatedBy: &createdByLookup,
		UpdatedBy: &updatedByLookup,
	}, nil
}

// Delete implements store.SLAStore.
func (s *SLAStore) Delete(rpc *model.DeleteOptions) error {
	d, dbErr := s.storage.Database()
	if dbErr != nil {
		return dberr.NewDBInternalError("postgres.sla.delete.database_connection_error", dbErr)
	}

	query, args, err := s.buildDeleteSLAQuery(rpc)
	if err != nil {
		return dberr.NewDBInternalError("postgres.sla.delete.query_build_error", err)
	}

	res, err := d.Exec(rpc.Context, query, args...)
	if err != nil {
		return dberr.NewDBInternalError("postgres.sla.delete.execution_error", err)
	}

	affected := res.RowsAffected()
	if affected == 0 {
		return dberr.NewDBNoRowsError("postgres.sla.delete.no_rows_affected")
	}

	return nil
}

// List implements store.SLAStore.
func (s *SLAStore) List(rpc *model.SearchOptions) (*cases.SLAList, error) {
	d, dbErr := s.storage.Database()
	if dbErr != nil {
		return nil, dberr.NewDBInternalError("postgres.sla.list.database_connection_error", dbErr)
	}

	query, args, err := s.buildSearchSLAQuery(rpc)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.sla.list.query_build_error", err)
	}

	rows, err := d.Query(rpc.Context, query, args...)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.sla.list.execution_error", err)
	}
	defer rows.Close()

	var slaList []*cases.SLA
	lCount := 0
	next := false
	// Check if we want to fetch all records
	//
	// If the size is -1, we want to fetch all records
	fetchAll := rpc.GetSize() == -1

	for rows.Next() {
		// If not fetching all records, check the size limit
		if !fetchAll && lCount >= int(rpc.GetSize()) {
			next = true
			break
		}

		sla := &cases.SLA{}

		var (
			createdBy, updatedBy, calendar cases.Lookup
			tempCreatedAt, tempUpdatedAt   time.Time
			tempValidFrom, tempValidTo     time.Time
		)

		scanArgs := s.buildScanArgs(
			rpc.Fields, sla,
			&createdBy, &updatedBy, &calendar,
			&tempCreatedAt, &tempUpdatedAt,
			&tempValidFrom, &tempValidTo,
		)
		if err := rows.Scan(scanArgs...); err != nil {
			return nil, dberr.NewDBInternalError("postgres.sla.list.row_scan_error", err)
		}

		s.populateSLAFields(rpc.Fields, sla, &createdBy, &updatedBy, &calendar, tempCreatedAt, tempUpdatedAt, tempValidFrom, tempValidTo)
		slaList = append(slaList, sla)
		lCount++
	}

	return &cases.SLAList{
		Page:  int32(rpc.Page),
		Next:  next,
		Items: slaList,
	}, nil
}

// Update implements store.SLAStore.
func (s *SLAStore) Update(rpc *model.UpdateOptions, l *cases.SLA) (*cases.SLA, error) {
	d, dbErr := s.storage.Database()
	if dbErr != nil {
		return nil, dberr.NewDBInternalError("postgres.sla.update.database_connection_error", dbErr)
	}

	query, args, err := s.buildUpdateSLAQuery(rpc, l)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.sla.update.query_build_error", err)
	}

	var (
		createdBy, updatedBy, calendar cases.Lookup
		createdAt, updatedAt           time.Time
		validFrom, validTo             time.Time
	)

	err = d.QueryRow(rpc.Context, query, args...).Scan(
		&l.Id, &l.Name, &createdAt, &updatedAt, &l.Description,
		&validFrom, &validTo, &calendar.Id, &calendar.Name,
		&l.ReactionTime.Hours, &l.ReactionTime.Minutes,
		&l.ResolutionTime.Hours, &l.ResolutionTime.Minutes,
		&createdBy.Id, &createdBy.Name, &updatedBy.Id, &updatedBy.Name,
	)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.sla.update.execution_error", err)
	}

	// Convert the valid from and valid to timestamps to local time
	l.CreatedAt = util.Timestamp(createdAt)
	l.UpdatedAt = util.Timestamp(updatedAt)
	l.ValidFrom = util.Timestamp(validFrom)
	l.ValidTo = util.Timestamp(validTo)

	l.CreatedBy = &createdBy
	l.UpdatedBy = &updatedBy
	l.Calendar = &calendar

	return l, nil
}

// buildCreateSLAQuery constructs the SQL insert query and returns the query string and arguments.
func (s SLAStore) buildCreateSLAQuery(rpc *model.CreateOptions, sla *cases.SLA) (string, []interface{}, error) {
	// Convert valid_from and valid_to from int64 timestamp to time.Time
	validFrom := util.LocalTime(sla.ValidFrom)
	validTo := util.LocalTime(sla.ValidTo)

	query := createSLAQuery
	args := []interface{}{
		sla.Name,                   // $1 name
		rpc.Session.GetDomainId(),  // $2 dc
		rpc.Time,                   // $3 created_at
		sla.Description,            // $4 description
		rpc.Session.GetUserId(),    // $5 created_by
		validFrom,                  // $6 valid_from
		validTo,                    // $7 valid_to
		sla.Calendar.Id,            // $8 calendar_id
		sla.ReactionTime.Hours,     // $9 reaction_time_hours
		sla.ReactionTime.Minutes,   // $10 reaction_time_minutes
		sla.ResolutionTime.Hours,   // $11 resolution_time_hours
		sla.ResolutionTime.Minutes, // $12 resolution_time_minutes
	}
	return query, args, nil
}

// buildDeleteSLAQuery constructs the SQL delete query and returns the query string and arguments.
func (s SLAStore) buildDeleteSLAQuery(rpc *model.DeleteOptions) (string, []interface{}, error) {
	convertedIds := util.Int64SliceToStringSlice(rpc.IDs)
	ids := util.FieldsFunc(convertedIds, util.InlineFields)

	query := deleteSLAQuery
	args := []interface{}{pq.Array(ids), rpc.Session.GetDomainId()}
	return query, args, nil
}

// buildSearchSLAQuery constructs the SQL search query and returns the query builder.
func (s SLAStore) buildSearchSLAQuery(rpc *model.SearchOptions) (string, []interface{}, error) {
	convertedIds := util.Int64SliceToStringSlice(rpc.IDs)
	ids := util.FieldsFunc(convertedIds, util.InlineFields)

	queryBuilder := sq.Select().
		From("cases.sla AS g").
		Where(sq.Eq{"g.dc": rpc.Session.GetDomainId()}).
		PlaceholderFormat(sq.Dollar)

	fields := util.FieldsFunc(rpc.Fields, util.InlineFields)
	rpc.Fields = append(fields, "id")

	for _, field := range rpc.Fields {
		switch field {
		case "id", "name", "reaction_time_hours",
			"reaction_time_minutes", "resolution_time_hours",
			"resolution_time_minutes", "created_at", "updated_at":
			queryBuilder = queryBuilder.Column("g." + field)
		case "description":
			// Use COALESCE to handle null values for description
			queryBuilder = queryBuilder.Column("COALESCE(g.description, '') AS description")
		case "valid_from":
			// Use COALESCE to handle null values for valid_from
			queryBuilder = queryBuilder.Column("COALESCE(g.valid_from, '') AS valid_from")
		case "valid_to":
			// Use COALESCE to handle null values for valid_to
			queryBuilder = queryBuilder.Column("COALESCE(g.valid_to, '') AS valid_to")
		case "calendar":
			// Include calendar_id and calendar_name
			queryBuilder = queryBuilder.
				Column("g.calendar_id").
				Column("COALESCE(cal.name, '') AS calendar_name").
				LeftJoin("flow.calendar AS cal ON cal.id = g.calendar_id")
		case "created_by":
			// Handle nulls using COALESCE for created_by
			queryBuilder = queryBuilder.
				Column("COALESCE(created_by.id, 0) AS cbi").
				Column("COALESCE(created_by.name, '') AS cbn").
				LeftJoin("directory.wbt_auth AS created_by ON g.created_by = created_by.id")
		case "updated_by":
			// Handle nulls using COALESCE for updated_by
			queryBuilder = queryBuilder.
				Column("COALESCE(updated_by.id, 0) AS ubi").
				Column("COALESCE(updated_by.name, '') AS ubn").
				LeftJoin("directory.wbt_auth AS updated_by ON g.updated_by = updated_by.id")
		}
	}

	if len(ids) > 0 {
		queryBuilder = queryBuilder.Where(sq.Eq{"g.id": ids})
	}

	if name, ok := rpc.Filter["name"].(string); ok && len(name) > 0 {
		substrs := util.Substring(name)
		combinedLike := strings.Join(substrs, "%")
		queryBuilder = queryBuilder.Where(sq.ILike{"g.name": combinedLike})
	}

	parsedFields := util.FieldsFunc(rpc.Sort, util.InlineFields)

	var sortFields []string

	for _, sortField := range parsedFields {
		desc := false
		if strings.HasPrefix(sortField, "!") {
			desc = true
			sortField = strings.TrimPrefix(sortField, "!")
		}

		var column string
		switch sortField {
		case "name", "description", "valid_from", "valid_to":
			column = "g." + sortField
		default:
			continue
		}

		if desc {
			column += " DESC"
		} else {
			column += " ASC"
		}

		sortFields = append(sortFields, column)
	}

	// Apply sorting
	queryBuilder = queryBuilder.OrderBy(sortFields...)

	size := rpc.GetSize()
	page := rpc.GetPage()

	// Apply offset only if page > 1
	if rpc.Page > 1 {
		queryBuilder = queryBuilder.Offset(uint64((page - 1) * size))
	}

	// Apply limit
	if size != -1 {
		queryBuilder = queryBuilder.Limit(uint64(size + 1)) // Request one more record to check if there's a next page
	}

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return "", nil, dberr.NewDBInternalError("postgres.sla.query_build.sql_generation_error", err)
	}

	return store.CompactSQL(query), args, nil
}

func (s SLAStore) buildUpdateSLAQuery(rpc *model.UpdateOptions, l *cases.SLA) (string, []interface{}, error) {
	// Initialize Squirrel builder
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	// Create a Squirrel update builder
	updateBuilder := psql.Update("cases.sla").
		Set("updated_at", rpc.Time).
		Set("updated_by", rpc.Session.GetUserId())

	// Convert the valid from and valid to timestamps to local time
	validFrom := util.LocalTime(l.ValidFrom)
	validTo := util.LocalTime(l.ValidTo)

	// Add the fields to the update query if they are provided
	for _, field := range rpc.Fields {
		switch field {
		case "name":
			if l.Name != "" {
				updateBuilder = updateBuilder.Set("name", l.Name)
			}
		case "description":
			// Use NULLIF to set NULL if description is an empty string
			updateBuilder = updateBuilder.Set("description", sq.Expr("NULLIF(?, '')", l.Description))
		case "valid_from":
			// Use NULLIF to set NULL if valid_from is 0
			updateBuilder = updateBuilder.Set("valid_from", sq.Expr("NULLIF(?, 0)", validFrom))
		case "valid_to":
			// Use NULLIF to set NULL if valid_to is 0
			updateBuilder = updateBuilder.Set("valid_to", sq.Expr("NULLIF(?, 0)", validTo))
		case "calendar_id":
			if l.Calendar.Id != 0 {
				updateBuilder = updateBuilder.Set("calendar_id", l.Calendar.Id)
			}
		case "reaction_time_hours":
			if l.ReactionTime.Hours != 0 {
				updateBuilder = updateBuilder.Set("reaction_time_hours", l.ReactionTime.Hours)
			}
		case "reaction_time_minutes":
			if l.ReactionTime.Minutes != 0 {
				updateBuilder = updateBuilder.Set("reaction_time_minutes", l.ReactionTime.Minutes)
			}
		case "resolution_time_hours":
			if l.ResolutionTime.Hours != 0 {
				updateBuilder = updateBuilder.Set("resolution_time_hours", l.ResolutionTime.Hours)
			}
		case "resolution_time_minutes":
			if l.ResolutionTime.Minutes != 0 {
				updateBuilder = updateBuilder.Set("resolution_time_minutes", l.ResolutionTime.Minutes)
			}
		}
	}

	// Add the WHERE clause for id and dc
	updateBuilder = updateBuilder.Where(sq.Eq{"id": l.Id, "dc": rpc.Session.GetDomainId()})

	// Build the SQL string and the arguments slice
	sql, args, err := updateBuilder.ToSql()
	if err != nil {
		return "", nil, err
	}

	// Construct the final SQL query with joins for created_by and updated_by
	query := fmt.Sprintf(`
WITH upd AS (%s
    RETURNING id, name, created_at, updated_at,
              description, valid_from, valid_to, calendar_id,
              reaction_time_hours, reaction_time_minutes,
              resolution_time_hours, resolution_time_minutes,
              created_by, updated_by)
SELECT
    upd.id,
    upd.name,
    upd.created_at,
    upd.updated_at,
    COALESCE(upd.description, '')              AS description,
    COALESCE(upd.valid_from, NULL)               AS valid_from,
    COALESCE(upd.valid_to, NULL)                 AS valid_to,
    upd.calendar_id,
    COALESCE(cal.name, '')                     AS calendar_name,
    upd.reaction_time_hours,
    upd.reaction_time_minutes,
    upd.resolution_time_hours,
    upd.resolution_time_minutes,
    upd.created_by                             AS created_by_id,
    COALESCE(c.name::text, c.username, '')     AS created_by_name,
    upd.updated_by                             AS updated_by_id,
    COALESCE(u.name::text, u.username, '')     AS updated_by_name
FROM upd
         LEFT JOIN directory.wbt_user u ON u.id = upd.updated_by
         LEFT JOIN directory.wbt_user c ON c.id = upd.created_by
         LEFT JOIN flow.calendar cal ON cal.id = upd.calendar_id;

    `, sql)

	return store.CompactSQL(query), args, nil
}

var (
	createSLAQuery = store.CompactSQL(`
WITH ins AS (
    INSERT INTO cases.sla (
                           name, dc, created_at, description,
                           created_by, updated_at, updated_by,
                           valid_from, valid_to, calendar_id,
                           reaction_time_hours, reaction_time_minutes,
                           resolution_time_hours, resolution_time_minutes
        )
        VALUES ($1, $2, $3, NULLIF($4, ''), $5, $3, $5, $6, $7, $8, $9, $10, $11, $12)
        RETURNING *)
SELECT ins.id,
       ins.name,
       ins.created_at,
       COALESCE(ins.description, '')      AS description,
       COALESCE(ins.valid_from, NULL)       AS valid_from,
       COALESCE(ins.valid_to, NULL)         AS valid_to,
       ins.calendar_id,
       COALESCE(cal.name, '')             AS calendar_name,
       ins.reaction_time_hours,
       ins.reaction_time_minutes,
       ins.resolution_time_hours,
       ins.resolution_time_minutes,
       ins.created_by                     AS created_by_id,
       COALESCE(c.name::text, c.username) AS created_by_name,
       ins.updated_at,
       ins.updated_by                     AS updated_by_id,
       COALESCE(u.name::text, u.username) AS updated_by_name
FROM ins
         LEFT JOIN directory.wbt_user u ON u.id = ins.updated_by
         LEFT JOIN directory.wbt_user c ON c.id = ins.created_by
         LEFT JOIN flow.calendar cal ON cal.id = ins.calendar_id;
	`)

	deleteSLAQuery = store.CompactSQL(`
                   DELETE FROM cases.sla
                   WHERE id = ANY($1) AND dc = $2`)
)

func (s *SLAStore) buildScanArgs(fields []string,
	sla *cases.SLA,
	createdBy, updatedBy, calendar *cases.Lookup,
	tempCreatedAt, tempUpdatedAt *time.Time,
	tempValidFrom, tempValidTo *time.Time,
) []interface{} {
	var scanArgs []interface{}

	for _, field := range fields {
		switch field {
		case "id":
			scanArgs = append(scanArgs, &sla.Id)
		case "name":
			scanArgs = append(scanArgs, &sla.Name)
		case "description":
			scanArgs = append(scanArgs, &sla.Description)
		case "valid_from":
			scanArgs = append(scanArgs, tempValidFrom)
		case "valid_to":
			scanArgs = append(scanArgs, tempValidTo)
		case "calendar":
			scanArgs = append(scanArgs, &calendar.Id, &calendar.Name)
		case "reaction_time_hours":
			scanArgs = append(scanArgs, &sla.ReactionTime.Hours)
		case "reaction_time_minutes":
			scanArgs = append(scanArgs, &sla.ReactionTime.Minutes)
		case "resolution_time_hours":
			scanArgs = append(scanArgs, &sla.ResolutionTime.Hours)
		case "resolution_time_minutes":
			scanArgs = append(scanArgs, &sla.ResolutionTime.Minutes)
		case "created_at":
			scanArgs = append(scanArgs, tempCreatedAt)
		case "updated_at":
			scanArgs = append(scanArgs, tempUpdatedAt)
		case "created_by":
			scanArgs = append(scanArgs, &createdBy.Id, &createdBy.Name)
		case "updated_by":
			scanArgs = append(scanArgs, &updatedBy.Id, &updatedBy.Name)
		}
	}

	return scanArgs
}

func (s *SLAStore) populateSLAFields(
	fields []string,
	sla *cases.SLA,
	createdBy, updatedBy, calendar *cases.Lookup,
	tempCreatedAt, tempUpdatedAt time.Time,
	tempValidFrom, tempValidTo time.Time,
) {
	for _, field := range fields {
		switch field {
		case "created_at":
			sla.CreatedAt = util.Timestamp(tempCreatedAt)
		case "updated_at":
			sla.UpdatedAt = util.Timestamp(tempUpdatedAt)
		case "valid_from":
			sla.ValidFrom = util.Timestamp(tempValidFrom)
		case "valid_to":
			sla.ValidTo = util.Timestamp(tempValidTo)
		case "created_by":
			sla.CreatedBy = createdBy
		case "updated_by":
			sla.UpdatedBy = updatedBy
		case "calendar":
			sla.Calendar = calendar
		}
	}
}

func NewSLAStore(store store.Store) (store.SLAStore, error) {
	if store == nil {
		return nil, dberr.NewDBError("postgres.new_sla.check.bad_arguments",
			"error creating SLA interface to the status_condition table, main store is nil")
	}
	return &SLAStore{storage: store}, nil
}
