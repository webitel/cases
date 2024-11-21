package postgres

import (
	"encoding/json"
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

type SLAConditionStore struct {
	storage store.Store
}

func (s *SLAConditionStore) Create(rpc *model.CreateOptions, add *cases.SLACondition, priorities []int64) (*cases.SLACondition, error) {
	db, dbErr := s.storage.Database()
	if dbErr != nil {
		return nil, dberr.NewDBInternalError("postgres.sla_condition.create.db_connection_error", dbErr)
	}

	// Build the combined SLACondition and Priority insert query
	query, args := s.buildCreateSLAConditionQuery(rpc, add)

	var (
		createdByLookup, updatedByLookup cases.Lookup
		createdAt, updatedAt             time.Time
	)

	prio := []*cases.Lookup{}

	rows, err := db.Query(rpc.Context, query, args...)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.sla_condition.create.execution_error", err)
	}
	defer rows.Close()

	// Iterate over the result set and collect all priorities
	for rows.Next() {
		var lookup cases.Lookup
		if err := rows.Scan(
			&add.Id, &add.Name, &createdAt,
			&add.ReactionTime, &add.ResolutionTime,
			&add.SlaId, &createdByLookup.Id,
			&createdByLookup.Name, &updatedAt, &updatedByLookup.Id,
			&updatedByLookup.Name, &lookup.Id, &lookup.Name,
		); err != nil {
			return nil, dberr.NewDBInternalError("postgres.sla_condition.create.scan_error", err)
		}
		prio = append(prio, &lookup)
	}

	// Check for errors after the iteration is complete
	if err := rows.Err(); err != nil {
		return nil, dberr.NewDBInternalError("postgres.sla_condition.create.iteration_error", err)
	}

	// Prepare the SLACondition object to return
	t := rpc.Time
	return &cases.SLACondition{
		Id:             add.Id,
		Name:           add.Name,
		ReactionTime:   add.ReactionTime,
		ResolutionTime: add.ResolutionTime,
		SlaId:          add.SlaId,
		CreatedAt:      util.Timestamp(t),
		UpdatedAt:      util.Timestamp(t),
		CreatedBy:      &createdByLookup,
		UpdatedBy:      &updatedByLookup,
		Priorities:     prio,
	}, nil
}

// Delete implements store.SLAConditionStore.
func (s *SLAConditionStore) Delete(rpc *model.DeleteOptions) error {
	// Establish a connection to the database
	d, dbErr := s.storage.Database()
	if dbErr != nil {
		return dberr.NewDBInternalError("postgres.sla_condition.delete.database_connection_error", dbErr)
	}

	// Build the delete query for SLACondition
	query, args, err := s.buildDeleteSLAConditionQuery(rpc)
	if err != nil {
		return dberr.NewDBInternalError("postgres.sla_condition.delete.query_build_error", err)
	}

	// Execute the delete query
	res, err := d.Exec(rpc.Context, query, args...)
	if err != nil {
		return dberr.NewDBInternalError("postgres.sla_condition.delete.execution_error", err)
	}

	// Check how many rows were affected by the delete operation
	affected := res.RowsAffected()
	if affected == 0 {
		return dberr.NewDBNoRowsError("postgres.sla_condition.delete.no_rows_affected")
	}

	return nil
}

// List implements store.SLAConditionStore.
func (s *SLAConditionStore) List(rpc *model.SearchOptions) (*cases.SLAConditionList, error) {
	// Establish a connection to the database
	d, dbErr := s.storage.Database()
	if dbErr != nil {
		return nil, dberr.NewDBInternalError("postgres.sla_condition.list.database_connection_error", dbErr)
	}

	// Build the search query for SLACondition
	query, args, err := s.buildSearchSLAConditionQuery(rpc)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.sla_condition.list.query_build_error", err)
	}

	// Execute the search query
	rows, err := d.Query(rpc.Context, query, args...)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.sla_condition.list.execution_error", err)
	}
	defer rows.Close()

	var slaConditionList []*cases.SLACondition
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

		slaCondition := &cases.SLACondition{}

		var (
			createdBy, updatedBy         cases.Lookup
			tempCreatedAt, tempUpdatedAt time.Time
			prioritiesJSON               []byte
		)
		// Build scan arguments dynamically based on the requested fields
		scanArgs := s.buildScanArgs(
			rpc.Fields, slaCondition, &createdBy,
			&updatedBy, &tempCreatedAt, &tempUpdatedAt,
			&prioritiesJSON,
		)

		if err := rows.Scan(scanArgs...); err != nil {
			return nil, dberr.NewDBInternalError("postgres.sla_condition.list.row_scan_error", err)
		}

		// Check if prioritiesJSON is not empty or NULL before unmarshalling
		if len(prioritiesJSON) > 0 {
			if err := json.Unmarshal(prioritiesJSON, &slaCondition.Priorities); err != nil {
				return nil, dberr.NewDBInternalError("postgres.sla_condition.list.json_unmarshal_error", err)
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
		return nil, dberr.NewDBInternalError("postgres.sla_condition.update.database_connection_error", dbErr)
	}

	// Begin a transaction
	tx, err := d.Begin(rpc.Context)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.sla_condition.update.transaction_begin_error", err)
	}
	defer tx.Rollback(rpc.Context) // Ensure rollback on error

	txManager := store.NewTxManager(tx)

	// Update priorities first if there are any IDs

	for _, fields := range rpc.Fields {
		if fields == "priorities" {
			priorityQuery, priorityArgs := s.buildUpdatePrioritiesQuery(rpc, l)

			// Execute the query and get the total number of rows affected
			var totalRowsAffected int
			err = txManager.QueryRow(rpc.Context, priorityQuery, priorityArgs...).Scan(&totalRowsAffected)
			if err != nil {
				return nil, dberr.NewDBInternalError("postgres.sla_condition.update.priorities_execution_error", err)
			}

			// Check if any rows were affected
			if totalRowsAffected == 0 {
				return nil, dberr.NewDBNoRowsError("postgres.sla_condition.update.no_priorities_affected")
			}
		}
	}

	// Build and execute the update query for sla_condition and return priorities JSON in one query
	query, args, err := s.buildUpdateSLAConditionQuery(rpc, l)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.sla_condition.update.query_build_error", err)
	}

	var (
		createdBy, updatedBy cases.Lookup
		createdAt, updatedAt time.Time
		prioritiesJSON       []byte
	)

	// Execute the update query for sla_condition and fetch priorities JSON
	err = txManager.QueryRow(rpc.Context, query, args...).Scan(
		&l.Id, &l.Name, &createdAt, &updatedAt,
		&l.ReactionTime, &l.ResolutionTime,
		&l.SlaId,
		&createdBy.Id, &createdBy.Name, &updatedBy.Id, &updatedBy.Name, // Corrected to include user names
		&prioritiesJSON, // Fetch JSON aggregated priorities
	)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.sla_condition.update.execution_error", err)
	}

	// Commit the transaction
	if err := tx.Commit(rpc.Context); err != nil {
		return nil, dberr.NewDBInternalError("postgres.sla_condition.update.transaction_commit_error", err)
	}

	// Process JSON aggregated priorities if not empty
	if len(prioritiesJSON) > 0 {
		if err := json.Unmarshal(prioritiesJSON, &l.Priorities); err != nil {
			return nil, dberr.NewDBInternalError("postgres.sla_condition.update.json_unmarshal_error", err)
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
		sla.ReactionTime,          // $4
		sla.ResolutionTime,        // $5
		sla.SlaId,                 // $6
		rpc.Session.GetDomainId(), // $7
	}

	// SQL query construction
	query := `
WITH inserted_sla AS (
    INSERT INTO cases.sla_condition (
                                     name, created_at, created_by, updated_at,
                                     updated_by, reaction_time, resolution_time,
									 sla_id, dc
        )
        VALUES ($1, $2, $3, $2, $3, $4, $5, $6, $7)
        RETURNING id, name, created_at,
            reaction_time, resolution_time,
			sla_id, created_by AS created_by_id,
            updated_by AS updated_by_id, updated_at),
     inserted_priorities AS (
         INSERT INTO cases.priority_sla_condition (
                                                   created_at, updated_at, created_by, updated_by,
                                                   sla_condition_id, priority_id, dc
             )
             SELECT $2, $2, $3, $3, inserted_sla.id, p.priority_id, $7
             FROM inserted_sla,
                  (SELECT unnest(ARRAY [`

	// Add placeholders for each priorityId to build the unnest array dynamically
	for i, priorityId := range rpc.Ids {
		if i > 0 {
			query += ", "
		}
		query += fmt.Sprintf("%d", priorityId)
	}

	query += `]) AS priority_id, $7 AS dc) p
        RETURNING sla_condition_id, priority_id
    )
SELECT inserted_sla.id,
       inserted_sla.name,
       inserted_sla.created_at,
       inserted_sla.reaction_time,
       inserted_sla.resolution_time,
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
	convertedIds := util.Int64SliceToStringSlice(rpc.IDs)
	ids := util.FieldsFunc(convertedIds, util.InlineFields)

	queryBuilder := sq.Select().
		From("cases.sla_condition AS g").
		Where(sq.Eq{"g.dc": rpc.Session.GetDomainId(), "g.sla_id": rpc.Id}).
		PlaceholderFormat(sq.Dollar)

	fields := util.FieldsFunc(rpc.Fields, util.InlineFields)

	//-----ALWAYS INCLUDE ID FIELD-----
	rpc.Fields = append(fields, "id") // FIXME ------ NEED to pass only if absent

	groupByFields := []string{"g.id"} // Start with the mandatory fields for GROUP BY

	for _, field := range rpc.Fields {
		switch field {
		case "id", "name", "reaction_time", "resolution_time",
			"sla_id", "created_at", "updated_at":
			queryBuilder = queryBuilder.Column("g." + field)
			groupByFields = append(groupByFields, "g."+field) // Add non-aggregated fields to GROUP BY
		case "created_by":
			// cbi = created_by_id
			// cbn = created_by_name
			// Use COALESCE to handle null values
			queryBuilder = queryBuilder.
				Column("COALESCE(created_by.id, 0) AS cbi").    // Handle NULL as 0 for created_by_id
				Column("COALESCE(created_by.name, '') AS cbn"). // Handle NULL as '' for created_by_name
				LeftJoin("directory.wbt_auth AS created_by ON g.created_by = created_by.id")
			groupByFields = append(groupByFields, "created_by.id", "created_by.name")
		case "updated_by":
			// ubi = updated_by_id
			// ubn = updated_by_name
			// Use COALESCE to handle null values
			queryBuilder = queryBuilder.
				Column("COALESCE(updated_by.id, 0) AS ubi").    // Handle NULL as 0 for updated_by_id
				Column("COALESCE(updated_by.name, '') AS ubn"). // Handle NULL as '' for updated_by_name
				LeftJoin("directory.wbt_auth AS updated_by ON g.updated_by = updated_by.id")
			groupByFields = append(groupByFields, "updated_by.id", "updated_by.name")
		case "priorities":
			// Aggregate priorities as JSON array
			queryBuilder = queryBuilder.
				Column("json_agg(json_build_object('id', p.id, 'name', p.name)) AS priorities").
				LeftJoin("cases.priority_sla_condition AS ps ON ps.sla_condition_id = g.id").
				LeftJoin("cases.priority AS p ON p.id = ps.priority_id")
			// No need to add aggregated columns to GROUP BY
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
		case "name", "reaction_time", "resolution_time", "created_at", "updated_at":
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

	// Apply GROUP BY clause
	queryBuilder = queryBuilder.GroupBy(groupByFields...)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return "", nil, dberr.NewDBInternalError("postgres.sla_condition.query_build.sql_generation_error", err)
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
		rpc.Time,                  // $5: timestamp for updated_at
	}

	// query that updates or inserts priorities and deletes non-selected ones
	query := `
WITH updated_priorities AS (
    INSERT INTO cases.priority_sla_condition (created_at, updated_at, created_by, updated_by, sla_condition_id,
                                              priority_id, dc)
        SELECT $5, $5, $2, $2, $1, unnest($4::bigint[]), $3
        ON CONFLICT (sla_condition_id, priority_id)
            DO UPDATE SET updated_at = EXCLUDED.updated_at, updated_by = EXCLUDED.updated_by
        RETURNING sla_condition_id, priority_id),
     deleted_priorities AS (
         DELETE FROM cases.priority_sla_condition
             WHERE sla_condition_id = $1
                 AND priority_id != ALL ($4::bigint[])
             RETURNING sla_condition_id, priority_id)
SELECT COUNT(*)
FROM (SELECT sla_condition_id
      FROM updated_priorities
      UNION ALL
      SELECT sla_condition_id
      FROM deleted_priorities) AS total_affected;
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
		case "reaction_time":
			updateBuilder = updateBuilder.Set("reaction_time", l.ReactionTime)
		case "resolution_time":
			updateBuilder = updateBuilder.Set("resolution_time", l.ResolutionTime)
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
        RETURNING id, name, created_at, updated_at, reaction_time,resolution_time, sla_id, created_by, updated_by)
SELECT usc.id,
       usc.name,
       usc.created_at,
       usc.updated_at,
       usc.reaction_time,
       usc.resolution_time,
       usc.sla_id,
       usc.created_by,
       COALESCE(c.name, '')                                    AS created_by_name,
       usc.updated_by,
       u.name                                                  AS updated_by_name,
       json_agg(json_build_object('id', p.id, 'name', p.name)) AS priorities_json
FROM upd_condition usc
         LEFT JOIN directory.wbt_user c ON c.id = usc.created_by
         LEFT JOIN directory.wbt_user u ON u.id = usc.updated_by
         LEFT JOIN cases.priority_sla_condition psc ON usc.id = psc.sla_condition_id
         LEFT JOIN cases.priority p ON p.id = psc.priority_id
GROUP BY usc.id, usc.name, usc.created_at, usc.updated_at,
         usc.reaction_time, usc.resolution_time,
		 usc.sla_id, usc.created_by, usc.updated_by, c.name, u.name
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
	var scanArgs []interface{}

	for _, field := range fields {
		switch field {
		case "id":
			scanArgs = append(scanArgs, &slaCondition.Id)
		case "name":
			scanArgs = append(scanArgs, &slaCondition.Name)
		case "reaction_time":
			scanArgs = append(scanArgs, &slaCondition.ReactionTime)
		case "resolution_time":
			scanArgs = append(scanArgs, &slaCondition.ResolutionTime)
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

func NewSLAConditionStore(store store.Store) (store.SLAConditionStore, error) {
	if store == nil {
		return nil, dberr.NewDBError("postgres.new_sla_condition.check.bad_arguments",
			"error creating SLACondition interface to the status_condition table, main store is nil")
	}
	return &SLAConditionStore{storage: store}, nil
}
