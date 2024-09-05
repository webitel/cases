package postgres

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/lib/pq"
	"github.com/webitel/cases/api/cases"
	"github.com/webitel/cases/internal/store"
	"github.com/webitel/cases/model"
	"github.com/webitel/cases/util"
)

type SLAConditionStore struct {
	storage store.Store
}

func (s *SLAConditionStore) Create(rpc *model.CreateOptions, add *cases.SLACondition, priorities []int64) (*cases.SLACondition, error) {
	db, dbErr := s.storage.Database()
	if dbErr != nil {
		return nil, model.NewInternalError("postgres.sla_condition.create.db_connection_error", dbErr.Error())
	}

	// Build the combined SLACondition and Priority insert query
	query, args := s.buildCreateSLAConditionQuery(rpc, add)

	// Execute the combined insert query and get the resulting fields
	var createdByLookup, updatedByLookup cases.Lookup
	var createdAt, updatedAt time.Time

	prio := []*cases.Lookup{}

	rows, err := db.Query(rpc.Context, query, args...)
	if err != nil {
		return nil, model.NewInternalError("postgres.sla_condition.create.execution_error", err.Error())
	}
	defer rows.Close()

	// Iterate over the result set and collect all priorities
	for rows.Next() {
		var lookup cases.Lookup
		if err := rows.Scan(
			&add.Id, &add.Name, &createdAt,
			&add.ReactionTimeHours, &add.ReactionTimeMinutes, &add.ResolutionTimeHours,
			&add.ResolutionTimeMinutes, &add.SlaId, &createdByLookup.Id,
			&createdByLookup.Name, &updatedAt, &updatedByLookup.Id,
			&updatedByLookup.Name, &lookup.Id, &lookup.Name,
		); err != nil {
			return nil, model.NewInternalError("postgres.sla_condition.create.scan_error", err.Error())
		}
		prio = append(prio, &lookup)
	}

	// Check for errors after the iteration is complete
	if err := rows.Err(); err != nil {
		return nil, model.NewInternalError("postgres.sla_condition.create.iteration_error", err.Error())
	}

	// Prepare the SLACondition object to return
	t := rpc.Time
	return &cases.SLACondition{
		Id:                    add.Id,
		Name:                  add.Name,
		ReactionTimeHours:     add.ReactionTimeHours,
		ReactionTimeMinutes:   add.ReactionTimeMinutes,
		ResolutionTimeHours:   add.ResolutionTimeHours,
		ResolutionTimeMinutes: add.ResolutionTimeMinutes,
		SlaId:                 add.SlaId,
		CreatedAt:             util.Timestamp(t),
		UpdatedAt:             util.Timestamp(t),
		CreatedBy:             &createdByLookup,
		UpdatedBy:             &updatedByLookup,
		Priorities:            prio,
	}, nil
}

// Delete implements store.SLAConditionStore.
func (s *SLAConditionStore) Delete(rpc *model.DeleteOptions) error {
	// Establish a connection to the database
	d, dbErr := s.storage.Database()
	if dbErr != nil {
		return model.NewInternalError("postgres.sla_condition.delete.database_connection_error", dbErr.Error())
	}

	// Build the delete query for SLACondition
	query, args, err := s.buildDeleteSLAConditionQuery(rpc)
	if err != nil {
		return model.NewInternalError("postgres.sla_condition.delete.query_build_error", err.Error())
	}

	// Execute the delete query
	res, err := d.Exec(rpc.Context, query, args...)
	if err != nil {
		return model.NewInternalError("postgres.sla_condition.delete.execution_error", err.Error())
	}

	// Check how many rows were affected by the delete operation
	affected := res.RowsAffected()
	if affected == 0 {
		return model.NewNotFoundError("postgres.sla_condition.delete.no_rows_affected", "No rows affected for deletion")
	}

	return nil
}

// List implements store.SLAConditionStore.
func (s *SLAConditionStore) List(rpc *model.SearchOptions) (*cases.SLAConditionList, error) {
	// Establish a connection to the database
	d, dbErr := s.storage.Database()
	if dbErr != nil {
		return nil, model.NewInternalError("postgres.sla_condition.list.database_connection_error", dbErr.Error())
	}

	// Build the search query for SLACondition
	query, args, err := s.buildSearchSLAConditionQuery(rpc)
	if err != nil {
		return nil, model.NewInternalError("postgres.sla_condition.list.query_build_error", err.Error())
	}

	// Execute the search query
	rows, err := d.Query(rpc.Context, query, args...)
	if err != nil {
		return nil, model.NewInternalError("postgres.sla_condition.list.execution_error", err.Error())
	}
	defer rows.Close()

	var slaConditionList []*cases.SLACondition
	lCount := 0
	next := false

	// Iterate over query results
	for rows.Next() {
		if lCount >= rpc.GetSize() {
			next = true
			break
		}

		slaCondition := &cases.SLACondition{}
		var createdBy, updatedBy cases.Lookup
		var tempCreatedAt, tempUpdatedAt time.Time

		// Priorities will be scanned as JSON
		var prioritiesJSON []byte

		// Build scan arguments dynamically based on the requested fields
		scanArgs := s.buildScanArgs(
			rpc.Fields, slaCondition, &createdBy,
			&updatedBy, &tempCreatedAt, &tempUpdatedAt,
			&prioritiesJSON,
		)

		if err := rows.Scan(scanArgs...); err != nil {
			return nil, model.NewInternalError("postgres.sla_condition.list.row_scan_error", err.Error())
		}

		// Check if prioritiesJSON is not empty or NULL before unmarshalling
		if len(prioritiesJSON) > 0 {
			if err := json.Unmarshal(prioritiesJSON, &slaCondition.Priorities); err != nil {
				return nil, model.NewInternalError("postgres.sla_condition.list.json_unmarshal_error", err.Error())
			}
		} else {
			// Handle NULL or empty JSON by initializing to an empty slice
			slaCondition.Priorities = []*cases.Lookup{}
		}

		// Populate SLACondition fields
		s.populateSLAConditionFields(
			rpc.Fields, slaCondition, &createdBy,
			&updatedBy, tempCreatedAt, tempUpdatedAt,
		)
		slaConditionList = append(slaConditionList, slaCondition)
		lCount++
	}

	return &cases.SLAConditionList{
		Page:  int32(rpc.Page),
		Next:  next,
		Items: slaConditionList,
	}, nil
}

// Update implements store.SLAConditionStore.
func (s *SLAConditionStore) Update(rpc *model.UpdateOptions, l *cases.SLACondition) (*cases.SLACondition, error) {
	d, dbErr := s.storage.Database()
	if dbErr != nil {
		return nil, model.NewInternalError("postgres.sla_condition.update.database_connection_error", dbErr.Error())
	}

	// Begin a transaction
	tx, err := d.Begin(rpc.Context)
	if err != nil {
		return nil, model.NewInternalError("postgres.sla_condition.update.transaction_begin_error", err.Error())
	}
	defer tx.Rollback(rpc.Context) // Ensure rollback on error

	txManager := store.NewTxManager(tx)

	// Update priorities first if there are any IDs
	if len(rpc.IDs) > 0 {
		// Example usage of Exec and checking rows affected
		priorityQuery, priorityArgs := s.buildUpdatePrioritiesQuery(rpc, l)

		// Execute the query to update priorities
		commandTag, execErr := txManager.Exec(rpc.Context, priorityQuery, priorityArgs...)
		if execErr != nil {
			return nil, model.NewInternalError("postgres.sla_condition.update.priorities_execution_error", execErr.Error())
		}

		// Check rows affected using commandTag
		rowsAffected := commandTag.RowsAffected()
		if rowsAffected == 0 {
			return nil, model.NewInternalError("postgres.sla_condition.update.no_priorities_affected", "No priorities were updated or deleted.")
		}

	}

	// Build and execute the update query for sla_condition and return priorities JSON in one query
	query, args, err := s.buildUpdateSLAConditionQuery(rpc, l)
	if err != nil {
		return nil, model.NewInternalError("postgres.sla_condition.update.query_build_error", err.Error())
	}

	var createdBy, updatedBy cases.Lookup
	var createdAt, updatedAt time.Time
	var prioritiesJSON []byte // For JSON aggregated priorities

	// Execute the update query for sla_condition and fetch priorities JSON
	err = txManager.QueryRow(rpc.Context, query, args...).Scan(
		&l.Id, &l.Name, &createdAt, &updatedAt,
		&l.ReactionTimeHours, &l.ReactionTimeMinutes,
		&l.ResolutionTimeHours, &l.ResolutionTimeMinutes, &l.SlaId,
		&createdBy.Id, &createdBy.Name, &updatedBy.Id, &updatedBy.Name, // Corrected to include user names
		&prioritiesJSON, // Fetch JSON aggregated priorities
	)
	if err != nil {
		return nil, model.NewInternalError("postgres.sla_condition.update.execution_error", err.Error())
	}

	// Commit the transaction
	if err := tx.Commit(rpc.Context); err != nil {
		return nil, model.NewInternalError("postgres.sla_condition.update.transaction_commit_error", err.Error())
	}

	// Process JSON aggregated priorities if not empty
	if len(prioritiesJSON) > 0 {
		if err := json.Unmarshal(prioritiesJSON, &l.Priorities); err != nil {
			return nil, model.NewInternalError("postgres.sla_condition.update.json_unmarshal_error", err.Error())
		}
	} else {
		// Initialize to an empty slice if no priorities
		l.Priorities = []*cases.Lookup{}
	}

	// Set timestamps and user information
	l.CreatedAt = util.Timestamp(createdAt)
	l.UpdatedAt = util.Timestamp(updatedAt)
	l.CreatedBy = &createdBy
	l.UpdatedBy = &updatedBy

	return l, nil
}

func (s *SLAConditionStore) buildCreateSLAConditionQuery(rpc *model.CreateOptions, sla *cases.SLACondition) (string, []interface{}) {
	// Create arguments for the SQL query
	args := []interface{}{
		sla.Name,                  // $1
		rpc.Time,                  // $2
		rpc.Session.GetUserId(),   // $3
		sla.ReactionTimeHours,     // $4
		sla.ReactionTimeMinutes,   // $5
		sla.ResolutionTimeHours,   // $6
		sla.ResolutionTimeMinutes, // $7
		sla.SlaId,                 // $8
		rpc.Session.GetDomainId(), // $9
	}

	// SQL query construction
	query := `
WITH inserted_sla AS (
    INSERT INTO cases.sla_condition (
                                     name, created_at, created_by, updated_at,
                                     updated_by, reaction_time_hours,
                                     reaction_time_minutes, resolution_time_hours,
                                     resolution_time_minutes, sla_id, dc
        )
        VALUES ($1, $2, $3, $2, $3, $4, $5, $6, $7, $8, $9)
        RETURNING id, name, created_at, reaction_time_hours,
            reaction_time_minutes, resolution_time_hours,
            resolution_time_minutes, sla_id, created_by AS created_by_id,
            updated_by AS updated_by_id, updated_at),
     inserted_priorities AS (
         INSERT INTO cases.priority_sla_condition (
                                                   created_at, updated_at, created_by, updated_by,
                                                   sla_condition_id, priority_id, dc
             )
             SELECT $2, $2, $3, $3, inserted_sla.id, p.priority_id, $9
             FROM inserted_sla,
                  (SELECT unnest(ARRAY [`

	// Add placeholders for each priorityId to build the unnest array dynamically
	for i, priorityId := range rpc.Ids {
		if i > 0 {
			query += ", "
		}
		query += fmt.Sprintf("%d", priorityId)
	}

	query += `]) AS priority_id, $9 AS dc) p
        RETURNING sla_condition_id, priority_id
    )
SELECT inserted_sla.id,
       inserted_sla.name,
       inserted_sla.created_at,
       inserted_sla.reaction_time_hours,
       inserted_sla.reaction_time_minutes,
       inserted_sla.resolution_time_hours,
       inserted_sla.resolution_time_minutes,
       inserted_sla.sla_id,
       inserted_sla.created_by_id,
       COALESCE(c.name, c.username) AS created_by_name,
       inserted_sla.updated_at,
       inserted_sla.updated_by_id,
       COALESCE(u.name, u.username) AS updated_by_name,
       inserted_priorities.priority_id,
       p.name
FROM inserted_sla
         LEFT JOIN directory.wbt_user u ON u.id = inserted_sla.updated_by_id
         LEFT JOIN directory.wbt_user c ON c.id = inserted_sla.created_by_id
         LEFT JOIN inserted_priorities ON inserted_sla.id = inserted_priorities.sla_condition_id
         LEFT JOIN cases.priority p ON p.id = inserted_priorities.priority_id;`

	return store.CompactSQL(query), args
}

// Helper function to build the delete query for SLACondition
func (s *SLAConditionStore) buildDeleteSLAConditionQuery(rpc *model.DeleteOptions) (string, []interface{}, error) {
	// Create base query for deletion
	query := deleteSLAConditionQuery

	// Arguments for the query
	args := []interface{}{
		rpc.IDs[0],                // $1 is the SLA Condition ID to delete
		rpc.Session.GetDomainId(), // $2 is the domain context (dc)
	}

	return query, args, nil
}

// buildSearchSLAConditionQuery constructs the SQL search query for SLAConditions.
func (s *SLAConditionStore) buildSearchSLAConditionQuery(rpc *model.SearchOptions) (string, []interface{}, error) {
	convertedIds := rpc.FieldsUtil.Int64SliceToStringSlice(rpc.IDs)
	ids := rpc.FieldsUtil.FieldsFunc(convertedIds, rpc.FieldsUtil.InlineFields)

	queryBuilder := sq.Select().
		From("cases.sla_condition AS g").
		Where(sq.Eq{"g.dc": rpc.Session.GetDomainId(), "g.sla_id": rpc.Id}).
		PlaceholderFormat(sq.Dollar)

	fields := rpc.FieldsUtil.FieldsFunc(rpc.Fields, rpc.FieldsUtil.InlineFields)

	groupByFields := []string{"g.id"} // Start with the default group by the field

	for _, field := range fields {
		switch field {
		case "id", "name", "reaction_time_hours", "reaction_time_minutes",
			"resolution_time_hours", "resolution_time_minutes", "sla_id",
			"created_at", "updated_at":
			queryBuilder = queryBuilder.Column("g." + field)
			groupByFields = append(groupByFields, "g."+field) // Add to group by fields
		case "created_by":
			queryBuilder = queryBuilder.Column("created_by.id AS cbi, created_by.name AS cbn").
				LeftJoin("directory.wbt_auth AS created_by ON g.created_by = created_by.id")
			groupByFields = append(groupByFields, "created_by.id", "created_by.name") // Add to group by fields
		case "updated_by":
			queryBuilder = queryBuilder.Column("updated_by.id AS ubi, updated_by.name AS ubn").
				LeftJoin("directory.wbt_auth AS updated_by ON g.updated_by = updated_by.id")
			groupByFields = append(groupByFields, "updated_by.id", "updated_by.name") // Add to group by fields
		case "priority":
			// Aggregate priorities as JSON array
			queryBuilder = queryBuilder.
				Column("json_agg(json_build_object('id', p.id, 'name', p.name)) AS priorities").
				LeftJoin("cases.priority_sla_condition AS ps ON ps.sla_condition_id = g.id").
				LeftJoin("cases.priority AS p ON p.id = ps.priority_id").
				GroupBy("g.id") // Group by SLACondition ID
		}
	}
	if len(ids) > 0 {
		queryBuilder = queryBuilder.Where(sq.Eq{"g.id": ids})
	}

	if name, ok := rpc.Filter["name"].(string); ok && len(name) > 0 {
		substrs := rpc.Match.Substring(name)
		combinedLike := strings.Join(substrs, "%")
		queryBuilder = queryBuilder.Where(sq.ILike{"g.name": "%" + combinedLike + "%"})
	}

	parsedFields := rpc.FieldsUtil.FieldsFunc(rpc.Sort, rpc.FieldsUtil.InlineFields)

	var sortFields []string

	for _, sortField := range parsedFields {
		desc := false
		if strings.HasPrefix(sortField, "!") {
			desc = true
			sortField = strings.TrimPrefix(sortField, "!")
		}

		var column string
		switch sortField {
		case "name", "reaction_time_hours", "reaction_time_minutes",
			"resolution_time_hours", "resolution_time_minutes", "created_at", "updated_at":
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
	// Add the GROUP BY clause with all non-aggregated fields
	queryBuilder = queryBuilder.GroupBy(groupByFields...)
	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return "", nil, model.NewInternalError("postgres.sla_condition.query_build.sql_generation_error", err.Error())
	}

	return store.CompactSQL(query), args, nil
}

func (s *SLAConditionStore) buildUpdatePrioritiesQuery(rpc *model.UpdateOptions, l *cases.SLACondition) (string, []interface{}) {
	// Prepare arguments for the SQL query
	args := []interface{}{
		l.Id,                      // $1: sla_condition_id
		rpc.Session.GetUserId(),   // $2: created_by and updated_by
		rpc.Session.GetDomainId(), // $3: dc
		pq.Array(rpc.IDs),         // $4: ARRAY of priority IDs
	}

	query := `
 WITH updated_priorities AS (
    -- Insert new priorities or update existing ones
    INSERT INTO cases.priority_sla_condition (created_at, updated_at, created_by, updated_by, sla_condition_id,
                                              priority_id, dc)
        SELECT CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, $2, $2, $1, unnest($4::bigint[]), $3
        ON CONFLICT (sla_condition_id, priority_id)
            DO UPDATE SET updated_at = EXCLUDED.updated_at, updated_by = EXCLUDED.updated_by
        RETURNING sla_condition_id, priority_id),
     deleted_priorities AS (
         -- Delete priorities that are not in the selected list
         DELETE FROM cases.priority_sla_condition
             WHERE sla_condition_id = $1
                 AND priority_id != ALL ($4::bigint[])
             RETURNING sla_condition_id, priority_id)
SELECT 1; -- Dummy select to complete the query

    `

	return query, args
}

// Function to build the update query for sla_condition and return priorities JSON
func (s *SLAConditionStore) buildUpdateSLAConditionQuery(rpc *model.UpdateOptions, l *cases.SLACondition) (string, []interface{}, error) {
	updateBuilder := sq.Update("cases.sla_condition").
		PlaceholderFormat(sq.Dollar). // Set placeholder format to Dollar for PostgreSQL
		Set("updated_at", rpc.Time).
		Set("updated_by", rpc.Session.GetUserId()).
		Where(sq.Eq{"id": l.Id, "dc": rpc.Session.GetDomainId()})

	// Dynamically add fields to the update builder based on provided fields
	for _, field := range rpc.Fields {
		switch field {
		case "name":
			updateBuilder = updateBuilder.Set("name", l.Name)
		case "reaction_time_hours":
			updateBuilder = updateBuilder.Set("reaction_time_hours", l.ReactionTimeHours)
		case "reaction_time_minutes":
			updateBuilder = updateBuilder.Set("reaction_time_minutes", l.ReactionTimeMinutes)
		case "resolution_time_hours":
			updateBuilder = updateBuilder.Set("resolution_time_hours", l.ResolutionTimeHours)
		case "resolution_time_minutes":
			updateBuilder = updateBuilder.Set("resolution_time_minutes", l.ResolutionTimeMinutes)
		case "sla_id":
			updateBuilder = updateBuilder.Set("sla_id", l.SlaId)
		}
	}

	// Convert the update query to SQL string and arguments
	updateSQL, args, err := updateBuilder.ToSql()
	if err != nil {
		return "", nil, err
	}

	query := fmt.Sprintf(`
   WITH upd_condition AS (%s
        RETURNING id, name, created_at, updated_at, reaction_time_hours, reaction_time_minutes,
                  resolution_time_hours, resolution_time_minutes, sla_id, created_by, updated_by)
SELECT usc.id,
       usc.name,
       usc.created_at,
       usc.updated_at,
       usc.reaction_time_hours,
       usc.reaction_time_minutes,
       usc.resolution_time_hours,
       usc.resolution_time_minutes,
       usc.sla_id,
       usc.created_by,
       c.name                                                  AS created_by_name,
       usc.updated_by,
       u.name                                                  AS updated_by_name,
       json_agg(json_build_object('id', p.id, 'name', p.name)) AS priorities_json
FROM upd_condition usc
         LEFT JOIN directory.wbt_user c ON c.id = usc.created_by
         LEFT JOIN directory.wbt_user u ON u.id = usc.updated_by
         LEFT JOIN cases.priority_sla_condition psc ON usc.id = psc.sla_condition_id
         LEFT JOIN cases.priority p ON p.id = psc.priority_id
GROUP BY usc.id, usc.name, usc.created_at, usc.updated_at,
         usc.reaction_time_hours, usc.reaction_time_minutes,
         usc.resolution_time_hours, usc.resolution_time_minutes, usc.sla_id,
         usc.created_by, usc.updated_by, c.name, u.name
    `, updateSQL)

	return query, args, nil
}

// Helper function to build scan arguments based on fields
func (s *SLAConditionStore) buildScanArgs(
	fields []string,
	slaCondition *cases.SLACondition,
	createdBy, updatedBy *cases.Lookup,
	createdAt, updatedAt *time.Time,
	prioritiesJSON *[]byte,
) []interface{} {
	scanArgs := []interface{}{&slaCondition.Id}

	for _, field := range fields {
		switch field {
		case "name":
			scanArgs = append(scanArgs, &slaCondition.Name)
		case "reaction_time_hours":
			scanArgs = append(scanArgs, &slaCondition.ReactionTimeHours)
		case "reaction_time_minutes":
			scanArgs = append(scanArgs, &slaCondition.ReactionTimeMinutes)
		case "resolution_time_hours":
			scanArgs = append(scanArgs, &slaCondition.ResolutionTimeHours)
		case "resolution_time_minutes":
			scanArgs = append(scanArgs, &slaCondition.ResolutionTimeMinutes)
		case "sla_id":
			scanArgs = append(scanArgs, &slaCondition.SlaId)
		case "created_at":
			scanArgs = append(scanArgs, createdAt)
		case "updated_at":
			scanArgs = append(scanArgs, updatedAt)
		case "created_by":
			scanArgs = append(scanArgs, &createdBy.Id, &createdBy.Name)
		case "updated_by":
			scanArgs = append(scanArgs, &updatedBy.Id, &updatedBy.Name)
		case "priority":
			scanArgs = append(scanArgs, prioritiesJSON)
		}
	}

	return scanArgs
}

// Helper function to populate SLACondition fields after scanning
func (s *SLAConditionStore) populateSLAConditionFields(
	fields []string,
	slaCondition *cases.SLACondition,
	createdBy,
	updatedBy *cases.Lookup,
	createdAt,
	updatedAt time.Time,
) {
	for _, field := range fields {
		switch field {
		case "created_by":
			slaCondition.CreatedBy = createdBy
		case "updated_by":
			slaCondition.UpdatedBy = updatedBy
		case "created_at":
			slaCondition.CreatedAt = util.Timestamp(createdAt)
		case "updated_at":
			slaCondition.UpdatedAt = util.Timestamp(updatedAt)
		}
	}
}

var deleteSLAConditionQuery = store.CompactSQL(
	`DELETE FROM cases.sla_condition
	 WHERE id = $1 AND dc = $2
	`)

func NewSLAConditionStore(store store.Store) (store.SLAConditionStore, model.AppError) {
	if store == nil {
		return nil, model.NewInternalError("postgres.new_sla_condition.check.bad_arguments",
			"error creating SLACondition interface to the status_condition table, main store is nil")
	}
	return &SLAConditionStore{storage: store}, nil
}
