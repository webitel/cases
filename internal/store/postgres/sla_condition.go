package postgres

import (
	"database/sql"
	"fmt"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/webitel/cases/internal/model"
	"github.com/webitel/cases/internal/model/options"
	storeutil "github.com/webitel/cases/internal/store/util"
	"strings"
	"time"

	storeUtil "github.com/webitel/cases/internal/store/util"

	sq "github.com/Masterminds/squirrel"
	"github.com/lib/pq"
	"github.com/webitel/cases/api/cases"
	"github.com/webitel/cases/internal/errors"
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

func (s *SLAConditionStore) Create(rpc options.Creator, add *model.SLACondition) (*model.SLACondition, error) {
	db, dbErr := s.storage.Database()
	if dbErr != nil {
		return nil, errors.Internal(fmt.Sprintf("postgres.sla_condition.create.db_connection_error: %v", dbErr))
	}

	// Build the combined SLACondition and Priority insert query
	sqlizer, err := s.buildCreateSLAConditionQuery(rpc, add)
	if err != nil {
		return nil, err
	}
	query, args, err := sqlizer.ToSql()
	if err != nil {
		return nil, errors.Internal(fmt.Sprintf("postgres.sla_condition.create.query_build_error: %v", err))
	}
	var res model.SLACondition
	err = pgxscan.Get(rpc, db, &res, query, args...)
	if err != nil {
		return nil, ParseError(err)
	}
	return &res, nil
}

// Delete implements store.SLAConditionStore.
func (s *SLAConditionStore) Delete(rpc options.Deleter) (*model.SLACondition, error) {
	// Establish a connection to the database
	d, dbErr := s.storage.Database()
	if dbErr != nil {
		return nil, errors.Internal(fmt.Sprintf("postgres.sla_condition.delete.database_connection_error: %v", dbErr))
	}

	// Build the delete query for SLACondition
	query, args, err := s.buildDeleteSLAConditionQuery(rpc)
	if err != nil {
		return nil, err
	}

	var res model.SLACondition
	err = pgxscan.Get(rpc, d, &res, query, args...)
	if err != nil {
		return nil, ParseError(err)
	}
	return nil, nil
}

// List implements store.SLAConditionStore.
func (s *SLAConditionStore) List(rpc options.Searcher) ([]*model.SLACondition, error) {
	// Establish a connection to the database
	d, dbErr := s.storage.Database()
	if dbErr != nil {
		return nil, errors.Internal(fmt.Sprintf("postgres.sla_condition.list.database_connection_error: %v", dbErr))
	}

	// Build the search query for SLACondition
	query, args, err := s.buildSearchSLAConditionQuery(rpc)
	if err != nil {
		return nil, err
	}

	var res []*model.SLACondition
	err = pgxscan.Select(rpc, d, &res, query, args...)
	if err != nil {
		return nil, ParseError(err)
	}

	return res, nil
}

// Update implements store.SLAConditionStore.
func (s *SLAConditionStore) Update(rpc options.Updator, l *model.SLACondition) (*model.SLACondition, error) {
	d, dbErr := s.storage.Database()
	if dbErr != nil {
		return nil, errors.Internal(fmt.Sprintf("postgres.sla_condition.update.database_connection_error: %v", dbErr))
	}

	// Begin a transaction
	tx, err := d.Begin(rpc)
	if err != nil {
		return nil, errors.Internal(fmt.Sprintf("postgres.sla_condition.update.transaction_begin_error: %v", err))
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
				return nil, ParseError(err)
			}

			// Check if any rows were affected
			if totalRowsAffected == 0 {
				return nil, errors.NotFound("no rows affected by priorities update")
			}
		}
	}

	// Build and execute the update query for sla_condition and return priorities JSON in one query
	query, args, err := s.buildUpdateSLAConditionQuery(rpc, l)
	if err != nil {
		return nil, err
	}

	var res []*model.SLACondition
	err = pgxscan.Select(rpc, txManager, &res, query, args...)
	if err != nil {
		return nil, ParseError(err)
	}

	// Commit the transaction
	if err := tx.Commit(rpc); err != nil {
		return nil, ParseError(err)
	}
	return l, nil
}

func (s *SLAConditionStore) buildCreateSLAConditionQuery(rpc options.Creator, sla *model.SLACondition) (sq.Sqlizer, error) {
	var ctes []*storeutil.CTE
	conditionCTEName := "inserted_condition"

	// select of the inserted condition
	base, err := buildSLAConditionColumns(sq.Select().From(conditionCTEName).PlaceholderFormat(sq.Dollar), rpc.GetFields(), conditionCTEName)
	if err != nil {
		return nil, err
	}

	// insert condition query
	conditionInsert := sq.Insert("cases.sla_condition").Columns("name", "created_at", "created_by", "updated_at",
		"updated_by", "reaction_time", "resolution_time",
		"sla_id", "dc").Values(sla.Name,
		rpc.RequestTime(),
		rpc.GetAuthOpts().GetUserId(),
		rpc.RequestTime(),
		rpc.GetAuthOpts().GetUserId(),
		sla.ReactionTime,
		sla.ResolutionTime,
		sla.SlaId,
		rpc.GetAuthOpts().GetDomainId(),
	).Suffix("RETURNING *")
	ctes = append(ctes, storeutil.NewCTE(conditionCTEName, conditionInsert))

	if len(sla.Priorities) != 0 { // if priorities not empty insert priorities
		priorityCTEName := "inserted_priorities"
		conditionPriorityInsert := sq.Insert("cases.priority_sla_condition").Columns("created_at", "updated_at", "created_by", "updated_by",
			"sla_condition_id", "priority_id", "dc",
		).Suffix("RETURNING *")
		for _, priority := range sla.Priorities {
			if priority.Id == 0 {
				continue
			}
			conditionPriorityInsert = conditionPriorityInsert.Values(
				sq.Expr(fmt.Sprintf("SELECT created_at, updated_at, created_by, updated_by, id, ?, dc FROM %s", priorityCTEName), priority.Id),
			)
		}
		ctes = append(ctes, storeutil.NewCTE(priorityCTEName, conditionPriorityInsert))
	}
	query, args, err := storeutil.FormAsCTEs(ctes)
	if err != nil {
		return nil, err
	}
	base = base.Prefix(query, args...)
	return base, nil
}

// Helper function to build the delete query for SLACondition
func (s *SLAConditionStore) buildDeleteSLAConditionQuery(rpc options.Deleter) (string, []interface{}, error) {
	// Create base query for deletion
	query := deleteSLAConditionQuery

	// Arguments for the query
	args := []interface{}{
		rpc.GetIDs()[0],                 // $1 is the SLA Condition ID to delete
		rpc.GetAuthOpts().GetDomainId(), // $2 is the domain context (dc)
	}

	return query, args, nil
}

func buildSLAConditionColumns(queryBuilder sq.SelectBuilder, fields []string, tableAlias string) (sq.SelectBuilder, error) {
	for _, field := range fields {
		switch field {
		case "id", "name", "reaction_time", "resolution_time",
			"sla_id", "created_at", "updated_at":
			queryBuilder = queryBuilder.Column(storeutil.Ident(tableAlias, field))
		case "created_by":
			queryBuilder = storeutil.SetUserColumn(queryBuilder, tableAlias, "created_by", field)
		case "updated_by":
			queryBuilder = storeutil.SetUserColumn(queryBuilder, tableAlias, "updated_by", field)
		case "priorities":
			queryBuilder = queryBuilder.
				Column(`priorities.priorities`).
				LeftJoin(fmt.Sprintf(`LATERAL (
             SELECT JSON_AGG(JSON_BUILD_OBJECT('id', p.id, 'name', p.name)) priorities
                             FROM cases.priority p
                                      INNER JOIN cases.priority_sla_condition ps ON p.id = ps.priority_id AND ps.sla_condition_id = %s
             ) priorities ON true`, storeutil.Ident(tableAlias, "id")))
		}
	}
	return queryBuilder, nil
}

// buildSearchSLAConditionQuery constructs the SQL search query for SLAConditions.
func (s *SLAConditionStore) buildSearchSLAConditionQuery(rpc options.Searcher) (string, []interface{}, error) {
	queryBuilder := sq.Select().
		From("cases.sla_condition AS g").
		Where(sq.Eq{"g.dc": rpc.GetAuthOpts().GetDomainId()}).
		LeftJoin("cases.priority_sla_condition AS ps ON ps.sla_condition_id = g.id").
		PlaceholderFormat(sq.Dollar)
	if slaIdFilters := rpc.GetFilter("sla_id"); len(slaIdFilters) > 0 {
		queryBuilder = storeutil.ApplyFiltersToQuery(queryBuilder, "g.sla_id", slaIdFilters)
	}
	queryBuilder, err := buildSLAConditionColumns(queryBuilder, rpc.GetFields(), "g")
	if err != nil {
		return "", nil, err
	}
	if len(rpc.GetIDs()) > 0 {
		queryBuilder = queryBuilder.Where(sq.Eq{"g.id": rpc.GetIDs()})
	}

	// Move new filter logic for priority_id and name here, after the switch-case
	if priorityIdFilters := rpc.GetFilter("priority_id"); len(priorityIdFilters) > 0 {
		queryBuilder = storeutil.ApplyFiltersToQuery(queryBuilder, "ps.priority_id", priorityIdFilters)
	}

	nameFilters := rpc.GetFilter("name")
	if len(nameFilters) > 0 {
		f := nameFilters[0]
		if (f.Operator == "=" || f.Operator == "") && len(f.Value) > 0 {
			queryBuilder = storeUtil.AddSearchTerm(queryBuilder, f.Value, "g.name")
		}
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
		queryBuilder = storeutil.ApplyDefaultSorting(rpc, queryBuilder, slaConditionDefaultSort)
	}

	// ---------Apply paging based on Search Opts ( page ; size ) -----------------
	queryBuilder = storeutil.ApplyPaging(rpc.GetPage(), rpc.GetSize(), queryBuilder)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return "", nil, err
	}

	return storeutil.CompactSQL(query), args, nil
}

func (s *SLAConditionStore) buildUpdatePrioritiesQuery(rpc options.Updator, l *model.SLACondition) (string, []interface{}) {
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
func (s *SLAConditionStore) buildUpdateSLAConditionQuery(rpc options.Updator, l *model.SLACondition) (string, []interface{}, error) {
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
       COALESCE(u.name, '')                                     AS updated_by_name,
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

var deleteSLAConditionQuery = storeutil.CompactSQL(
	`DELETE FROM cases.sla_condition
	 WHERE id = $1 AND dc = $2
	`)

func NewSLAConditionStore(store *Store) (store.SLAConditionStore, error) {
	if store == nil {
		return nil, errors.New(
			"error creating SLACondition interface to the status_condition table, main store is nil")
	}
	return &SLAConditionStore{storage: store}, nil
}
