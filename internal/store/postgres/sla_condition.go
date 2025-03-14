package postgres

import (
	"database/sql"
	"encoding/json"
	"fmt"
	util2 "github.com/webitel/cases/internal/store/util"
	"github.com/webitel/cases/model/options"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/lib/pq"
	"github.com/webitel/cases/api/cases"
	dberr "github.com/webitel/cases/internal/errors"
	"github.com/webitel/cases/internal/store"
	"github.com/webitel/cases/internal/store/postgres/transaction"

	"github.com/webitel/cases/util"
)

const (
	slaConditionDefaultSort = "name"
)

type SLAConditionStore struct {
	storage *Store
}

func (s *SLAConditionStore) Create(rpc options.CreateOptions, add *cases.SLACondition, priorities []int64) (*cases.SLACondition, error) {
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

	rows, err := db.Query(rpc, query, args...)
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
	t := rpc.RequestTime()
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
func (s *SLAConditionStore) Delete(rpc options.DeleteOptions) error {
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
	res, err := d.Exec(rpc, query, args...)
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
func (s *SLAConditionStore) List(rpc options.SearchOptions) (*cases.SLAConditionList, error) {
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
	rows, err := d.Query(rpc, query, args...)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.sla_condition.list.execution_error", err)
	}
	defer rows.Close()

	// Prepare the result list
	var slaConditionList []*cases.SLACondition
	lCount := 0
	next := false
	fetchAll := rpc.GetSize() == -1

	for rows.Next() {
		if !fetchAll && lCount >= int(rpc.GetSize()) {
			next = true
			break
		}

		slaCondition := &cases.SLACondition{}
		var (
			createdBy, updatedBy         cases.Lookup
			tempCreatedAt, tempUpdatedAt time.Time
			prioritiesJSON               []byte // JSON field for priorities
		)

		// Build scan arguments dynamically
		scanArgs := s.buildScanArgs(
			rpc.GetFields(), slaCondition, &createdBy, &updatedBy, &tempCreatedAt, &tempUpdatedAt, &prioritiesJSON,
		)

		if err := rows.Scan(scanArgs...); err != nil {
			return nil, dberr.NewDBInternalError("postgres.sla_condition.list.row_scan_error", err)
		}

		// Populate SLACondition fields
		s.populateSLAConditionFields(
			rpc.GetFields(), slaCondition, &createdBy, &updatedBy, tempCreatedAt, tempUpdatedAt,
		)

		// Parse JSON priorities into SLACondition.Priorities
		if len(prioritiesJSON) > 0 {
			if err := json.Unmarshal(prioritiesJSON, &slaCondition.Priorities); err != nil {
				return nil, dberr.NewDBInternalError("postgres.sla_condition.list.priorities_parse_error", err)
			}
		}

		// Add SLACondition to the result list
		slaConditionList = append(slaConditionList, slaCondition)
		lCount++
	}

	if err := rows.Err(); err != nil {
		return nil, dberr.NewDBInternalError("postgres.sla_condition.list.rows_iteration_error", err)
	}

	return &cases.SLAConditionList{
		Page:  int32(rpc.GetPage()),
		Next:  next,
		Items: slaConditionList,
	}, nil
}

// Update implements store.SLAConditionStore.
func (s *SLAConditionStore) Update(rpc options.UpdateOptions, l *cases.SLACondition) (*cases.SLACondition, error) {
	d, dbErr := s.storage.Database()
	if dbErr != nil {
		return nil, dberr.NewDBInternalError("postgres.sla_condition.update.database_connection_error", dbErr)
	}

	// Begin a transaction
	tx, err := d.Begin(rpc)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.sla_condition.update.transaction_begin_error", err)
	}
	defer tx.Rollback(rpc) // Ensure rollback on error

	txManager := transaction.NewTxManager(tx)

	// Update priorities first if there are any IDs

	for _, fields := range rpc.GetMask() {
		if fields == "priorities" {
			priorityQuery, priorityArgs := s.buildUpdatePrioritiesQuery(rpc, l)

			// Execute the query and get the total number of rows affected
			var totalRowsAffected int
			err = txManager.QueryRow(rpc, priorityQuery, priorityArgs...).Scan(&totalRowsAffected)
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
	err = txManager.QueryRow(rpc, query, args...).Scan(
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
	if err := tx.Commit(rpc); err != nil {
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

func (s *SLAConditionStore) buildCreateSLAConditionQuery(rpc options.CreateOptions, sla *cases.SLACondition) (string, []interface{}) {
	// Create arguments for the SQL query
	args := []interface{}{
		sla.Name,                        // $1
		rpc.RequestTime(),               // $2
		rpc.GetAuthOpts().GetUserId(),   // $3
		sla.ReactionTime,                // $4
		sla.ResolutionTime,              // $5
		sla.SlaId,                       // $6
		rpc.GetAuthOpts().GetDomainId(), // $7
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
	// TODO REMOVE %d
	for i, priorityId := range rpc.GetIDs() {
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

	return util2.CompactSQL(query), args
}

// Helper function to build the delete query for SLACondition
func (s *SLAConditionStore) buildDeleteSLAConditionQuery(rpc options.DeleteOptions) (string, []interface{}, error) {
	// Create base query for deletion
	query := deleteSLAConditionQuery

	// Arguments for the query
	args := []interface{}{
		rpc.GetIDs()[0],                 // $1 is the SLA Condition ID to delete
		rpc.GetAuthOpts().GetDomainId(), // $2 is the domain context (dc)
	}

	return query, args, nil
}

// buildSearchSLAConditionQuery constructs the SQL search query for SLAConditions.
func (s *SLAConditionStore) buildSearchSLAConditionQuery(rpc options.SearchOptions) (string, []interface{}, error) {
	convertedIds := util.Int64SliceToStringSlice(rpc.GetIDs())
	ids := util.FieldsFunc(convertedIds, util.InlineFields)
	queryBuilder := sq.Select().
		From("cases.sla_condition AS g").
		Where(sq.Eq{"g.dc": rpc.GetAuthOpts().GetDomainId(), "g.sla_id": rpc.GetFilter("sla_id")}).
		LeftJoin("cases.priority_sla_condition AS ps ON ps.sla_condition_id = g.id").
		PlaceholderFormat(sq.Dollar)

	groupByFields := []string{"g.id"} // Start with the mandatory fields for GROUP BY

	for _, field := range rpc.GetFields() {
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
			queryBuilder = queryBuilder.
				Column(`COALESCE(
					JSON_AGG(
						JSON_BUILD_OBJECT('id', p.id, 'name', p.name)
					) FILTER (WHERE p.id IS NOT NULL),
					'[]'
				) AS priorities`).
				LeftJoin("cases.priority AS p ON ps.priority_id = p.id")
		}
	}
	if len(ids) > 0 {
		queryBuilder = queryBuilder.Where(sq.Eq{"g.id": ids})
	}

	if priorityId := rpc.GetFilter("priority_id"); priorityId != nil {
		// Join cases.priority_sla_condition only if filtering by priority_id
		queryBuilder = queryBuilder.
			Where(sq.Eq{"ps.priority_id": priorityId})
	}

	if name, ok := rpc.GetFilter("name").(string); ok && len(name) > 0 {
		queryBuilder = util2.AddSearchTerm(queryBuilder, name, "g.name")
	}

	// Adjust sort if calendar is present
	sortField := rpc.GetSort()
	// Remove any leading "+" or "-" for comparison
	field := strings.TrimPrefix(strings.TrimPrefix(sortField, "-"), "+")

	if field == "priorities" {
		s := "p.name"
		desc := strings.HasPrefix(sortField, "-")
		// Determine sort direction
		if desc {
			s += " DESC"
		} else {
			s += " ASC"
		}
		queryBuilder = queryBuilder.OrderBy(s)
	} else {
		// -------- Apply sorting ----------
		queryBuilder = util2.ApplyDefaultSorting(rpc, queryBuilder, slaConditionDefaultSort)
	}

	// ---------Apply paging based on Search Opts ( page ; size ) -----------------
	queryBuilder = util2.ApplyPaging(rpc.GetPage(), rpc.GetSize(), queryBuilder)

	// Apply GROUP BY clause
	queryBuilder = queryBuilder.GroupBy(groupByFields...)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return "", nil, dberr.NewDBInternalError("postgres.sla_condition.query_build.sql_generation_error", err)
	}

	return util2.CompactSQL(query), args, nil
}

func (s *SLAConditionStore) buildUpdatePrioritiesQuery(rpc options.UpdateOptions, l *cases.SLACondition) (string, []interface{}) {
	// Prepare arguments for the SQL query
	args := []interface{}{
		l.Id,                            // $1: sla_condition_id
		rpc.GetAuthOpts().GetUserId(),   // $2: created_by and updated_by
		rpc.GetAuthOpts().GetDomainId(), // $3: dc
		pq.Array(rpc.GetIDs()),          // $4: ARRAY of priority IDs
		rpc.RequestTime(),               // $5: timestamp for updated_at
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
func (s *SLAConditionStore) buildUpdateSLAConditionQuery(rpc options.UpdateOptions, l *cases.SLACondition) (string, []interface{}, error) {
	updateBuilder := sq.Update("cases.sla_condition").
		PlaceholderFormat(sq.Dollar). // Set placeholder format to Dollar for PostgreSQL
		Set("updated_at", rpc.RequestTime()).
		Set("updated_by", rpc.GetAuthOpts().GetUserId()).
		Where(sq.Eq{"id": l.Id, "dc": rpc.GetAuthOpts().GetDomainId()})

	// Dynamically add fields to the update builder based on provided fields
	for _, field := range rpc.GetMask() {
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

// buildScanArgs dynamically constructs the scan arguments based on the requested fields.
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
		case "priorities":
			scanArgs = append(scanArgs, prioritiesJSON)
		}
	}

	return scanArgs
}

// populateSLAConditionFields populates SLACondition fields after scanning.
func (s *SLAConditionStore) populateSLAConditionFields(
	fields []string,
	slaCondition *cases.SLACondition,
	createdBy, updatedBy *cases.Lookup,
	createdAt, updatedAt time.Time,
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

// populatePriorities appends a priority to the SLACondition's Priorities slice.
func (s *SLAConditionStore) populatePriorities(
	slaCondition *cases.SLACondition,
	priorityID *sql.NullInt64,
	priorityName *sql.NullString,
) {
	if priorityID.Valid && priorityName.Valid {
		slaCondition.Priorities = append(slaCondition.Priorities, &cases.Lookup{
			Id:   priorityID.Int64,
			Name: priorityName.String,
		})
	}
}

var deleteSLAConditionQuery = util2.CompactSQL(
	`DELETE FROM cases.sla_condition
	 WHERE id = $1 AND dc = $2
	`)

func NewSLAConditionStore(store *Store) (store.SLAConditionStore, error) {
	if store == nil {
		return nil, dberr.NewDBError("postgres.new_sla_condition.check.bad_arguments",
			"error creating SLACondition interface to the status_condition table, main store is nil")
	}
	return &SLAConditionStore{storage: store}, nil
}
