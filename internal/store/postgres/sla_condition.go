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

	// Build the SQL update query using a WITH clause
	query, args, err := s.buildUpdateSLAConditionAndPrioritiesQuery(rpc, l)
	if err != nil {
		return nil, model.NewInternalError("postgres.sla_condition.update.query_build_error", err.Error())
	}

	var createdBy, updatedBy cases.Lookup
	var createdAt, updatedAt time.Time

	// Execute the query with the combined update and insert operations
	err = d.QueryRow(rpc.Context, query, args...).Scan(
		&l.Id, &l.Name, &createdAt, &updatedAt,
		&l.ReactionTimeHours, &l.ReactionTimeMinutes,
		&l.ResolutionTimeHours, &l.ResolutionTimeMinutes, &l.SlaId,
		&createdBy.Id, &createdBy.Name, &updatedBy.Id, &updatedBy.Name,
	)
	if err != nil {
		return nil, model.NewInternalError("postgres.sla_condition.update.execution_error", err.Error())
	}

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

	for _, field := range fields {
		switch field {
		case "id", "name", "reaction_time_hours", "reaction_time_minutes",
			"resolution_time_hours", "resolution_time_minutes", "sla_id",
			"created_at", "updated_at":
			queryBuilder = queryBuilder.Column("g." + field)
		case "created_by":
			// cbi = created_by_id,
			// cbn = created_by_name
			queryBuilder = queryBuilder.Column("created_by.id AS cbi, created_by.name AS cbn").
				LeftJoin("directory.wbt_auth AS created_by ON g.created_by = created_by.id")
		case "updated_by":
			// ubi = updated_by_id,
			// ubn = updated_by_name
			queryBuilder = queryBuilder.Column("updated_by.id AS ubi, updated_by.name AS ubn").
				LeftJoin("directory.wbt_auth AS updated_by ON g.updated_by = updated_by.id")
		case "priority":
			// Aggregate priorities as JSON array
			queryBuilder = queryBuilder.
				Column("json_agg(json_build_object('id', p.id, 'name', p.name)) AS priorities").
				LeftJoin("cases.priority_sla_condition AS ps ON ps.sla_condition_id = g.id").
				LeftJoin("cases.priority AS p ON p.id = ps.priority_id").
				GroupBy("g.id")
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

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return "", nil, model.NewInternalError("postgres.sla_condition.query_build.sql_generation_error", err.Error())
	}

	return store.CompactSQL(query), args, nil
}

// buildUpdateSLAConditionAndPrioritiesQuery constructs the SQL update query with a MERGE clause for SLACondition and priorities using Squirrel.
func (s *SLAConditionStore) buildUpdateSLAConditionAndPrioritiesQuery(rpc *model.UpdateOptions, l *cases.SLACondition) (string, []interface{}, error) {
	// Initialize base arguments for the SQL query
	args := []interface{}{
		rpc.Time,                  // $1 - updated_at
		rpc.Session.GetUserId(),   // $2 - updated_by
		l.Id,                      // $3 - sla_condition_id
		rpc.Session.GetDomainId(), // $4 - dc id
		pq.Array(rpc.IDs),         // $5 - array of priority_ids
	}

	// Initialize Squirrel update builder for SLACondition table
	updateBuilder := sq.Update("cases.sla_condition").
		Set("updated_at", rpc.Time).
		Set("updated_by", rpc.Session.GetUserId()).
		Where(sq.Eq{"id": l.Id, "dc": rpc.Session.GetDomainId()})

	// Dynamically add fields to the update builder
	if l.Name != "" {
		updateBuilder = updateBuilder.Set("name", l.Name)
		args = append(args, l.Name)
	}
	if l.ReactionTimeHours != 0 {
		updateBuilder = updateBuilder.Set("reaction_time_hours", l.ReactionTimeHours)
		args = append(args, l.ReactionTimeHours)
	}
	if l.ReactionTimeMinutes != 0 {
		updateBuilder = updateBuilder.Set("reaction_time_minutes", l.ReactionTimeMinutes)
		args = append(args, l.ReactionTimeMinutes)
	}
	if l.ResolutionTimeHours != 0 {
		updateBuilder = updateBuilder.Set("resolution_time_hours", l.ResolutionTimeHours)
		args = append(args, l.ResolutionTimeHours)
	}
	if l.ResolutionTimeMinutes != 0 {
		updateBuilder = updateBuilder.Set("resolution_time_minutes", l.ResolutionTimeMinutes)
		args = append(args, l.ResolutionTimeMinutes)
	}
	if l.SlaId != 0 {
		updateBuilder = updateBuilder.Set("sla_id", l.SlaId)
		args = append(args, l.SlaId)
	}

	// Convert the update query to SQL string and arguments
	updateSQL, updateArgs, err := updateBuilder.ToSql()
	if err != nil {
		return "", nil, model.NewInternalError("postgres.sla_condition.update.query_build_error", err.Error())
	}

	// Add any additional arguments generated by the update query
	args = append(args, updateArgs...)

	// Construct the final SQL query using MERGE for upserting priorities and a DELETE operation for missing priorities
	query := fmt.Sprintf(`
WITH updated_sla_condition AS (%s
    RETURNING id, name, created_at, updated_at,
              reaction_time_hours, reaction_time_minutes,
              resolution_time_hours, resolution_time_minutes, sla_id,
              created_by, updated_by),
     merge_priorities AS (MERGE INTO cases.priority_sla_condition AS target
    USING (
        SELECT unnest($5::int[]) AS priority_id
    ) AS source
    ON target.sla_condition_id = $3 AND target.priority_id = source.priority_id
    WHEN MATCHED THEN
        UPDATE SET updated_at = $1,
                   updated_by = $2
    WHEN NOT MATCHED THEN
        INSERT (created_at, updated_at, created_by, updated_by, sla_condition_id, priority_id, dc)
        VALUES ($1, $1, $2, $2, $3, source.priority_id, $4)
    RETURNING sla_condition_id, priority_id),
     deleted_priorities AS (
         DELETE FROM cases.priority_sla_condition
             WHERE sla_condition_id = $3
                 AND priority_id NOT IN (SELECT unnest($5::int[]))
             RETURNING sla_condition_id, priority_id)
SELECT usc.id,
       usc.name,
       usc.created_at,
       usc.updated_at,
       usc.reaction_time_hours,
       usc.reaction_time_minutes,
       usc.resolution_time_hours,
       usc.resolution_time_minutes,
       usc.sla_id,
       usc.created_by                     AS created_by_id,
       COALESCE(c.name::text, c.username) AS created_by_name,
       usc.updated_by                     AS updated_by_id,
       COALESCE(u.name::text, u.username) AS updated_by_name,
       mp.priority_id,
       p.name                             AS priority_name
FROM updated_sla_condition usc
         LEFT JOIN directory.wbt_user u ON u.id = usc.updated_by
         LEFT JOIN directory.wbt_user c ON c.id = usc.created_by
         LEFT JOIN merge_priorities mp ON usc.id = mp.sla_condition_id
         LEFT JOIN cases.priority p ON p.id = mp.priority_id;
`, updateSQL)

	return store.CompactSQL(query), args, nil
}

// // buildUpdateSLAConditionAndPrioritiesQuery constructs the SQL update query with a WITH clause for SLACondition and priorities using Squirrel.
// func (s *SLAConditionStore) buildUpdateSLAConditionAndPrioritiesQuery(rpc *model.UpdateOptions, l *cases.SLACondition) (string, []interface{}, error) {
// 	// Initialize base arguments for the SQL query
// 	args := []interface{}{
// 		rpc.Time,                  // $1 - updated_at
// 		rpc.Session.GetUserId(),   // $2 - updated_by
// 		l.Id,                      // $3 - sla_condition_id
// 		rpc.Session.GetDomainId(), // $4 - dc id
// 		pq.Array(rpc.IDs),         // $5 - array of priority_ids
// 	}

// 	// Initialize Squirrel update builder for SLACondition table
// 	updateBuilder := sq.Update("cases.sla_condition").
// 		Set("updated_at", rpc.Time).
// 		Set("updated_by", rpc.Session.GetUserId()).
// 		Where(sq.Eq{"id": l.Id, "dc": rpc.Session.GetDomainId()})

// 	// Dynamically add fields to the update builder
// 	if l.Name != "" {
// 		updateBuilder = updateBuilder.Set("name", l.Name)
// 		args = append(args, l.Name) // $6 - name
// 	}
// 	if l.ReactionTimeHours != 0 {
// 		updateBuilder = updateBuilder.Set("reaction_time_hours", l.ReactionTimeHours)
// 		args = append(args, l.ReactionTimeHours) // $7 - reaction_time_hours
// 	}
// 	if l.ReactionTimeMinutes != 0 {
// 		updateBuilder = updateBuilder.Set("reaction_time_minutes", l.ReactionTimeMinutes)
// 		args = append(args, l.ReactionTimeMinutes) // $8 - reaction_time_minutes
// 	}
// 	if l.ResolutionTimeHours != 0 {
// 		updateBuilder = updateBuilder.Set("resolution_time_hours", l.ResolutionTimeHours)
// 		args = append(args, l.ResolutionTimeHours) // $9 - resolution_time_hours
// 	}
// 	if l.ResolutionTimeMinutes != 0 {
// 		updateBuilder = updateBuilder.Set("resolution_time_minutes", l.ResolutionTimeMinutes)
// 		args = append(args, l.ResolutionTimeMinutes) // $10 - resolution_time_minutes
// 	}
// 	if l.SlaId != 0 {
// 		updateBuilder = updateBuilder.Set("sla_id", l.SlaId)
// 		args = append(args, l.SlaId) // $11 - sla_id
// 	}

// 	// Convert the update query to SQL string and arguments
// 	updateSQL, updateArgs, err := updateBuilder.ToSql()
// 	if err != nil {
// 		return "", nil, model.NewInternalError("postgres.sla_condition.update.query_build_error", err.Error())
// 	}

// 	// Add any additional arguments generated by the update query
// 	args = append(args, updateArgs...)

// 	// Construct the final SQL query using the generated update SQL and priorities insert statement
// 	query := fmt.Sprintf(`
// WITH updated_sla_condition AS (%s
//     RETURNING id, name, created_at, updated_at,
//               reaction_time_hours, reaction_time_minutes,
//               resolution_time_hours, resolution_time_minutes, sla_id,
//               created_by, updated_by),
//      deleted_priorities AS (
//          DELETE FROM cases.priority_sla_condition
//              WHERE sla_condition_id = $3
//                  AND NOT (priority_id = ANY ($5))),
//      inserted_priorities AS (
//          INSERT INTO cases.priority_sla_condition (created_at, updated_at, created_by, updated_by, sla_condition_id,
//                                                    priority_id, dc)
//              SELECT $1, $1, $2, $2, $3, unnest($5::int[]), $4
//              ON CONFLICT (sla_condition_id, priority_id) DO NOTHING
//              RETURNING sla_condition_id, priority_id)
// SELECT usc.id,
//        usc.name,
//        usc.created_at,
//        usc.updated_at,
//        usc.reaction_time_hours,
//        usc.reaction_time_minutes,
//        usc.resolution_time_hours,
//        usc.resolution_time_minutes,
//        usc.sla_id,
//        usc.created_by                     AS created_by_id,
//        COALESCE(c.name::text, c.username) AS created_by_name,
//        usc.updated_by                     AS updated_by_id,
//        COALESCE(u.name::text, u.username) AS updated_by_name,
//        ip.priority_id,
//        p.name                             AS priority_name
// FROM updated_sla_condition usc
//          LEFT JOIN directory.wbt_user u ON u.id = usc.updated_by
//          LEFT JOIN directory.wbt_user c ON c.id = usc.created_by
//          LEFT JOIN inserted_priorities ip ON usc.id = ip.sla_condition_id
//          LEFT JOIN cases.priority p ON p.id = ip.priority_id;
// `, updateSQL)

// 	return store.CompactSQL(query), args, nil
// }

// Helper function to build scan arguments based on fields
func (s *SLAConditionStore) buildScanArgs(
	fields []string,
	slaCondition *cases.SLACondition,
	createdBy,
	updatedBy *cases.Lookup,
	createdAt,
	updatedAt *time.Time,
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
