package postgres

import (
	"encoding/json"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/lib/pq"
	"github.com/webitel/cases/api/cases"
	"github.com/webitel/cases/internal/store"
	"github.com/webitel/cases/model"
	"github.com/webitel/cases/util"
)

type CatalogStore struct {
	storage store.Store
}

// Create implements store.CatalogStore.
func (s *CatalogStore) Create(ctx *model.CreateOptions, add *cases.Catalog) (*cases.Catalog, error) {
	// Establish a connection to the database
	db, dbErr := s.storage.Database()
	if dbErr != nil {
		return nil, model.NewInternalError("postgres.catalog.create.database_connection_error", dbErr.Error())
	}

	// Build the combined query for inserting Catalog, teams, and skills
	query, args := s.buildCreateCatalogQuery(ctx, add)

	// Execute the query and scan the result into the Catalog fields
	var createdByLookup, updatedByLookup cases.Lookup
	var createdAt, updatedAt time.Time
	var teamLookups, skillLookups []byte

	err := db.QueryRow(ctx.Context, query, args...).Scan(
		&add.Id, &add.Name, &add.Description, &add.Prefix, &add.Code, &add.State,
		&createdAt, &updatedAt,
		&add.Sla.Id, &add.Sla.Name,
		&add.Status.Id, &add.Status.Name,
		&add.CloseReason.Id, &add.CloseReason.Name,
		&createdByLookup.Id, &createdByLookup.Name,
		&updatedByLookup.Id, &updatedByLookup.Name,
		&teamLookups,  // JSON array for teams
		&skillLookups, // JSON array for skills
	)
	if err != nil {
		return nil, model.NewInternalError("postgres.catalog.create.scan_error", err.Error())
	}

	// Unmarshal the JSON arrays into the Lookup slices
	if err := json.Unmarshal(teamLookups, &add.Teams); err != nil {
		return nil, model.NewInternalError("postgres.catalog.create.unmarshal_teams_error", err.Error())
	}
	if err := json.Unmarshal(skillLookups, &add.Skills); err != nil {
		return nil, model.NewInternalError("postgres.catalog.create.unmarshal_skills_error", err.Error())
	}

	// Prepare the Catalog to return
	add.CreatedAt = util.Timestamp(createdAt)
	add.UpdatedAt = util.Timestamp(updatedAt)
	add.CreatedBy = &createdByLookup
	add.UpdatedBy = &updatedByLookup

	// Return the created Catalog
	return add, nil
}

// Delete implements store.CatalogStore.
func (s *CatalogStore) Delete(ctx *model.DeleteOptions) error {
	// Establish a connection to the database
	db, dbErr := s.storage.Database()
	if dbErr != nil {
		return model.NewInternalError("postgres.catalog.delete.db_connection_error", dbErr.Error())
	}

	// Ensure that there are IDs to delete
	if len(ctx.IDs) == 0 {
		return model.NewBadRequestError("postgres.catalog.delete.no_ids_provided", "No IDs provided for deletion")
	}

	// Build the delete query
	query, args := s.buildDeleteCatalogQuery(ctx)

	// Execute the delete query
	res, err := db.Exec(ctx.Context, query, args...)
	if err != nil {
		return model.NewInternalError("postgres.catalog.delete.execution_error", err.Error())
	}

	// Check how many rows were affected
	if res.RowsAffected() == 0 {
		return model.NewNotFoundError("postgres.catalog.delete.no_rows_deleted", "No Catalog entries were deleted")
	}

	return nil
}

// List implements store.CatalogStore.
func (s *CatalogStore) List(ctx *model.SearchOptions) (*cases.CatalogList, error) {
	// Establish a connection to the database
	db, dbErr := s.storage.Database()
	if dbErr != nil {
		return nil, model.NewInternalError("postgres.catalog.list.database_connection_error", dbErr.Error())
	}

	// Build SQL query
	query, args, err := s.buildSearchCatalogQuery(ctx)
	if err != nil {
		return nil, model.NewInternalError("postgres.catalog.list.query_build_error", err.Error())
	}

	// Execute the query
	rows, err := db.Query(ctx.Context, query, args...)
	if err != nil {
		return nil, model.NewInternalError("postgres.catalog.list.query_execution_error", err.Error())
	}
	defer rows.Close()

	// Parse the result
	var catalogs []*cases.Catalog
	count := 0
	next := false

	for rows.Next() {
		if count >= ctx.Size {
			next = true
			break
		}

		var catalog cases.Catalog
		var teamLookups, skillLookups []byte

		err = rows.Scan(
			&catalog.Id, &catalog.Name, &catalog.CreatedAt,
			&catalog.Sla.Id, &catalog.Sla.Name,
			&catalog.Status.Id, &catalog.Status.Name,
			&catalog.CloseReason.Id, &catalog.CloseReason.Name,
			&catalog.CreatedBy.Id, &catalog.CreatedBy.Name,
			&catalog.UpdatedBy.Id, &catalog.UpdatedBy.Name, &catalog.UpdatedAt,
			&teamLookups, &skillLookups, &catalog.HasServices,
		)
		if err != nil {
			return nil, model.NewInternalError("postgres.catalog.list.scan_error", err.Error())
		}

		// Unmarshal the JSON arrays into the Lookup slices
		if err := json.Unmarshal(teamLookups, &catalog.Teams); err != nil {
			return nil, model.NewInternalError("postgres.catalog.list.unmarshal_teams_error", err.Error())
		}
		if err := json.Unmarshal(skillLookups, &catalog.Skills); err != nil {
			return nil, model.NewInternalError("postgres.catalog.list.unmarshal_skills_error", err.Error())
		}

		catalogs = append(catalogs, &catalog)
		count++
	}

	return &cases.CatalogList{
		Page:  int32(ctx.Page),
		Next:  next,
		Items: catalogs,
	}, nil
}

// Update implements store.CatalogStore.
func (s *CatalogStore) Update(ctx *model.UpdateOptions, lookup *cases.Catalog) (*cases.Catalog, error) {
	// Establish a connection to the database
	db, dbErr := s.storage.Database()
	if dbErr != nil {
		return nil, model.NewInternalError("postgres.catalog.update.database_connection_error", dbErr.Error())
	}

	// Start a transaction using the TxManager
	tx, err := db.Begin(ctx.Context)
	if err != nil {
		return nil, model.NewInternalError("postgres.catalog.update.transaction_start_error", err.Error())
	}
	txManager := store.NewTxManager(tx)   // Create a new TxManager instance
	defer txManager.Rollback(ctx.Context) // Ensure rollback on error

	// Handle teams and skills updates if they exist
	if len(lookup.Teams) > 0 || len(lookup.Skills) > 0 {
		// Extract team and skill IDs from Lookup
		teamIDs := make([]int64, len(lookup.Teams))
		for i, team := range lookup.Teams {
			teamIDs[i] = team.Id
		}

		skillIDs := make([]int64, len(lookup.Skills))
		for i, skill := range lookup.Skills {
			skillIDs[i] = skill.Id
		}

		// Build query to update teams and skills
		query, args := s.buildUpdateTeamsAndSkillsQuery(
			lookup.Id, teamIDs,
			skillIDs, ctx.Session.GetUserId(),
			ctx.Time, ctx.Session.GetDomainId(),
		)

		// Execute the teams and skills update query and check for affected rows
		var affectedRows int
		err = txManager.QueryRow(ctx.Context, query, args...).Scan(&affectedRows)
		if err != nil {
			return nil, model.NewInternalError("postgres.catalog.update.teams_skills_update_error", err.Error())
		}

		// Optional check if no rows were affected
		if affectedRows == 0 {
			return nil, model.NewInternalError("postgres.catalog.update.no_teams_skills_affected", "No teams or skills were updated")
		}
	}

	// Build the update query for the Catalog
	query, args, err := s.buildUpdateCatalogQuery(ctx, lookup)
	if err != nil {
		return nil, model.NewInternalError("postgres.catalog.update.query_build_error", err.Error())
	}

	// Execute the update query for the catalog
	var createdByLookup, updatedByLookup cases.Lookup
	var createdAt, updatedAt time.Time
	var teamLookups, skillLookups []byte

	err = txManager.QueryRow(ctx.Context, query, args...).Scan(
		&lookup.Id, &lookup.Name, &createdAt,
		&lookup.Sla.Id, &lookup.Sla.Name,
		&lookup.Status.Id, &lookup.Status.Name,
		&lookup.CloseReason.Id, &lookup.CloseReason.Name,
		&createdByLookup.Id, &createdByLookup.Name,
		&updatedByLookup.Id, &updatedByLookup.Name, &updatedAt,
		&teamLookups, &skillLookups,
	)
	if err != nil {
		return nil, model.NewInternalError("postgres.catalog.update.execution_error", err.Error())
	}

	// Commit the transaction
	if err := txManager.Commit(ctx.Context); err != nil {
		return nil, model.NewInternalError("postgres.catalog.update.transaction_commit_error", err.Error())
	}

	// Unmarshal the JSON arrays for teams and skills
	if err := json.Unmarshal(teamLookups, &lookup.Teams); err != nil {
		return nil, model.NewInternalError("postgres.catalog.update.unmarshal_teams_error", err.Error())
	}
	if err := json.Unmarshal(skillLookups, &lookup.Skills); err != nil {
		return nil, model.NewInternalError("postgres.catalog.update.unmarshal_skills_error", err.Error())
	}

	// Prepare the updated Catalog to return
	lookup.CreatedAt = util.Timestamp(createdAt)
	lookup.UpdatedAt = util.Timestamp(updatedAt)
	lookup.CreatedBy = &createdByLookup
	lookup.UpdatedBy = &updatedByLookup

	// Return the updated Catalog
	return lookup, nil
}

// Helper method to build the combined insert query for Catalog and related entities
func (s *CatalogStore) buildCreateCatalogQuery(ctx *model.CreateOptions, add *cases.Catalog) (string, []interface{}) {
	// Define arguments for the query
	args := []interface{}{
		add.Name,                  // $1: name
		add.Description,           // $2: description (can be null)
		add.Prefix,                // $3: prefix (can be null)
		add.Code,                  // $4: code (can be null)
		ctx.Time,                  // $5: created_at, updated_at
		ctx.Session.GetUserId(),   // $6: created_by, updated_by
		add.Sla.Id,                // $7: sla_id
		add.Status.Id,             // $8: status_id
		add.CloseReason.Id,        // $9: close_reason_id
		add.State,                 // $10: state
		ctx.Session.GetDomainId(), // $11: domain ID (dc)
	}

	teamIds := make([]int64, len(add.Teams))
	for i, team := range add.Teams {
		teamIds[i] = team.Id
	}
	args = append(args, pq.Array(teamIds)) // $12: team_ids

	skillIds := make([]int64, len(add.Skills))
	for i, skill := range add.Skills {
		skillIds[i] = skill.Id
	}
	args = append(args, pq.Array(skillIds)) // $13: skill_ids

	// SQL query construction
	query := `
WITH inserted_catalog AS (
    INSERT INTO cases.service_catalog (
        name, description, prefix, code, created_at, created_by, updated_at,
        updated_by, sla_id, status_id, close_reason_id, state, dc
    ) VALUES ($1, $2, $3, $4, $5, $6, $5, $6, $7, $8, $9, $10, $11)
    RETURNING id, name, description, prefix, code, state, sla_id, status_id, close_reason_id,
              created_by, updated_by, created_at, updated_at
),
inserted_teams AS (
    INSERT INTO cases.team_catalog (catalog_id, team_id, created_by, updated_by, created_at, updated_at, dc)
    SELECT inserted_catalog.id, unnest($12::bigint[]), $6, $6, $5, $5, $11
    FROM inserted_catalog
    RETURNING catalog_id, team_id
),
inserted_skills AS (
    INSERT INTO cases.skill_catalog (catalog_id, skill_id, created_by, updated_by, created_at, updated_at, dc)
    SELECT inserted_catalog.id, unnest($13::bigint[]), $6, $6, $5, $5, $11
    FROM inserted_catalog
    RETURNING catalog_id, skill_id
),
teams_agg AS (
    SELECT inserted_teams.catalog_id,
           json_agg(json_build_object('id', team.id, 'name', team.name)) AS teams
    FROM inserted_teams
    LEFT JOIN call_center.cc_team team ON team.id = inserted_teams.team_id
    GROUP BY inserted_teams.catalog_id
),
skills_agg AS (
    SELECT inserted_skills.catalog_id,
           json_agg(json_build_object('id', skill.id, 'name', skill.name)) AS skills
    FROM inserted_skills
    LEFT JOIN call_center.cc_skill skill ON skill.id = inserted_skills.skill_id
    GROUP BY inserted_skills.catalog_id
)
SELECT inserted_catalog.id,
       inserted_catalog.name,
       inserted_catalog.description,
       inserted_catalog.prefix,
       inserted_catalog.code,
       inserted_catalog.state,
       inserted_catalog.created_at,
       inserted_catalog.updated_at,
       inserted_catalog.sla_id,
       sla.name AS sla_name,
       inserted_catalog.status_id,
       status.name AS status_name,
       inserted_catalog.close_reason_id,
       close_reason.name AS close_reason_name,
       inserted_catalog.created_by,
       created_by_user.name  AS created_by_name,
       inserted_catalog.updated_by,
       updated_by_user.name  AS updated_by_name,
       COALESCE(teams_agg.teams, '[]') AS teams,
       COALESCE(skills_agg.skills, '[]') AS skills
FROM inserted_catalog
LEFT JOIN cases.sla ON sla.id = inserted_catalog.sla_id
LEFT JOIN cases.status ON status.id = inserted_catalog.status_id
LEFT JOIN cases.close_reason ON close_reason.id = inserted_catalog.close_reason_id
LEFT JOIN directory.wbt_user created_by_user ON created_by_user.id = inserted_catalog.created_by
LEFT JOIN directory.wbt_user updated_by_user ON updated_by_user.id = inserted_catalog.updated_by
LEFT JOIN teams_agg ON teams_agg.catalog_id = inserted_catalog.id
LEFT JOIN skills_agg ON skills_agg.catalog_id = inserted_catalog.id;
`

	return store.CompactSQL(query), args
}

// Helper method to build the delete query for Catalog
func (s *CatalogStore) buildDeleteCatalogQuery(ctx *model.DeleteOptions) (string, []interface{}) {
	// Build the SQL query using the provided IDs in ctx.IDs
	query := `
		DELETE FROM cases.service_catalog
		WHERE id = ANY($1) AND dc = $2
	`

	// Use the array of IDs and domain ID (dc) for the deletion
	args := []interface{}{
		pq.Array(ctx.IDs),         // $1: array of catalog IDs to delete
		ctx.Session.GetDomainId(), // $2: domain ID to ensure proper scoping
	}

	return store.CompactSQL(query), args
}

// Helper method to build the search query for Catalog
func (s *CatalogStore) buildSearchCatalogQuery(ctx *model.SearchOptions) (string, []interface{}, error) {
	// Initialize query builder
	queryBuilder := sq.Select(
		"catalog.id",
		"catalog.name",
		"catalog.description",
		"catalog.prefix",
		"catalog.code",
		"catalog.state",
		"catalog.sla_id",
		"sla.name",
		"catalog.status_id",
		"status.name",
		"catalog.close_reason_id",
		"close_reason.name",
		"grp.name",
		"catalog.created_by",
		"created_by_user.name AS created_by_name",
		"catalog.updated_by",
		"updated_by_user.name AS updated_by_name",
		"catalog.created_at",
		"catalog.updated_at",
		`COALESCE(json_agg(json_build_object('id', team.id, 'name', team.name)) FILTER (WHERE team.id IS NOT NULL), '[]') AS teams`,
		`COALESCE(json_agg(json_build_object('id', skill.id, 'name', skill.name)) FILTER (WHERE skill.id IS NOT NULL), '[]') AS skills`,
		// Determine if the catalog has services
		`EXISTS (SELECT 1 FROM cases.catalog cs WHERE cs.parent_id = catalog.id) AS has_services`,
	).
		From("cases.service_catalog AS catalog").
		LeftJoin("cases.sla ON sla.id = catalog.sla_id").
		LeftJoin("cases.status ON status.id = catalog.status_id").
		LeftJoin("cases.close_reason ON close_reason.id = catalog.close_reason_id").
		LeftJoin("directory.wbt_user AS created_by_user ON created_by_user.id = catalog.created_by").
		LeftJoin("directory.wbt_user AS updated_by_user ON updated_by_user.id = catalog.updated_by").
		LeftJoin("call_center.cc_teams_catalog AS tc ON tc.catalog_id = catalog.id").
		LeftJoin("call_center.cc_teams AS team ON team.id = tc.team_id").
		LeftJoin("call_center.cc_skills_catalog AS sc ON sc.catalog_id = catalog.id").
		LeftJoin("call_center.cc_skills AS skill ON skill.id = sc.skill_id").
		GroupBy(
			"catalog.id", "sla.name", "status.name", "close_reason.name", "grp.name", "created_by_user.name", "updated_by_user.name",
		).
		Where(sq.Eq{"catalog.parent_id": nil}). // Fetch only top-level catalogs (where parent_id is NULL)
		PlaceholderFormat(sq.Dollar)

	// Apply filtering by name
	if name, ok := ctx.Filter["name"].(string); ok && len(name) > 0 {
		substr := ctx.Match.Substring(name)
		queryBuilder = queryBuilder.Where(sq.ILike{"catalog.name": substr})
	}

	// Apply filtering by state
	if state, ok := ctx.Filter["state"]; ok {
		queryBuilder = queryBuilder.Where(sq.Eq{"catalog.state": state})
	}

	// Apply sorting
	for _, sort := range ctx.Sort {
		queryBuilder = queryBuilder.OrderBy(sort)
	}

	// Apply pagination
	if ctx.Page > 0 && ctx.Size > 0 {
		queryBuilder = queryBuilder.Limit(uint64(ctx.Size + 1)).Offset(uint64((ctx.Page - 1) * ctx.Size))
	}

	// Build SQL query
	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return "", nil, model.NewInternalError("postgres.catalog.query_build_error", err.Error())
	}

	return store.CompactSQL(query), args, nil
}

func (s *CatalogStore) buildUpdateTeamsAndSkillsQuery(catalogID int64, teamIDs, skillIDs []int64, updatedBy int64, updatedAt time.Time, domainID int64) (string, []interface{}) {
	args := []interface{}{
		catalogID, // $1: catalog_id
		updatedBy, // $2: updated_by
		domainID,  // $3: dc (domain context)
		updatedAt, // $4: timestamp for updated_at
	}

	// Initialize base query
	query := `
	WITH `

	// Check if teams are provided and build the query for teams
	if len(teamIDs) > 0 {
		query += `
	updated_teams AS (
			INSERT INTO cases.team_catalog (catalog_id, team_id, updated_by, updated_at, dc)
			SELECT $1, unnest($5::bigint[]), $2, $4, $3
			ON CONFLICT (catalog_id, team_id)
			DO UPDATE SET updated_at = EXCLUDED.updated_at, updated_by = EXCLUDED.updated_by
			RETURNING catalog_id
		),
		deleted_teams AS (
			DELETE FROM cases.team_catalog
			WHERE catalog_id = $1
			AND team_id != ALL ($5::bigint[])
			RETURNING catalog_id
		),`
		args = append(args, pq.Array(teamIDs)) // Append team IDs to the args
	}

	// Check if skills are provided and build the query for skills
	if len(skillIDs) > 0 {
		query += `
		updated_skills AS (
			INSERT INTO cases.skill_catalog (catalog_id, skill_id, updated_by, updated_at, dc)
			SELECT $1, unnest($6::bigint[]), $2, $4, $3
			ON CONFLICT (catalog_id, skill_id)
			DO UPDATE SET updated_at = EXCLUDED.updated_at, updated_by = EXCLUDED.updated_by
			RETURNING catalog_id
		),
		deleted_skills AS (
			DELETE FROM cases.skill_catalog
			WHERE catalog_id = $1
			AND skill_id != ALL ($6::bigint[])
			RETURNING catalog_id
		),`
		args = append(args, pq.Array(skillIDs)) // Append skill IDs to the args
	}

	// Finish the query
	query += `
	SELECT COUNT(*)
	FROM (` +
		// If teams were provided, use the teams' CTE in the result union
		func() string {
			if len(teamIDs) > 0 {
				return `SELECT catalog_id FROM updated_teams
					UNION ALL
					SELECT catalog_id FROM deleted_teams`
			}
			return ""
		}() +
		// If skills were provided, use the skills' CTE in the result union
		func() string {
			if len(skillIDs) > 0 {
				if len(teamIDs) > 0 {
					return ` UNION ALL ` // Add the union operator only if teams also exist
				}
				return ""
			}
			return ""
		}() +
		func() string {
			if len(skillIDs) > 0 {
				return `SELECT catalog_id FROM updated_skills
					UNION ALL
					SELECT catalog_id FROM deleted_skills`
			}
			return ""
		}() + `
	) AS total_affected;
	`

	// Return the constructed query and arguments
	return store.CompactSQL(query), args
}

// Helper method to build the combined update and select query for Catalog using Squirrel
func (s *CatalogStore) buildUpdateCatalogQuery(ctx *model.UpdateOptions, lookup *cases.Catalog) (string, []interface{}, error) {
	// Start the update query with Squirrel Update Builder
	updateQueryBuilder := sq.Update("cases.service_catalog").
		PlaceholderFormat(sq.Dollar).
		Set("updated_at", ctx.Time).
		Set("updated_by", ctx.Session.GetUserId()).
		Where(sq.Eq{"id": lookup.Id, "dc": ctx.Session.GetDomainId()})

	// Dynamically set fields based on what the user wants to update
	for _, field := range ctx.Fields {
		switch field {
		case "name":
			updateQueryBuilder = updateQueryBuilder.Set("name", lookup.Name)
		case "sla_id":
			updateQueryBuilder = updateQueryBuilder.Set("sla_id", lookup.Sla.Id)
		case "status_id":
			updateQueryBuilder = updateQueryBuilder.Set("status_id", lookup.Status.Id)
		case "close_reason_id":
			updateQueryBuilder = updateQueryBuilder.Set("close_reason_id", lookup.CloseReason.Id)
		}
	}

	// Convert the update query to SQL
	updateQuery, args, err := updateQueryBuilder.ToSql()
	if err != nil {
		return "", nil, err
	}

	// Now build the select query to return the updated catalog
	selectQueryBuilder := sq.Select(
		"catalog.id",
		"catalog.name",
		"catalog.created_at",
		"catalog.sla_id",
		"sla.name",
		"catalog.status_id",
		"status.name",
		"catalog.close_reason_id",
		"close_reason.name",
		"grp.name",
		"catalog.created_by",
		"created_by_user.name AS created_by_name",
		"catalog.updated_by",
		"updated_by_user.name AS updated_by_name",
		"catalog.updated_at",
		// Teams JSON aggregation
		`COALESCE(json_agg(json_build_object('id', team.id, 'name', team.name)) FILTER (WHERE team.id IS NOT NULL), '[]') AS teams`,
		// Skills JSON aggregation
		`COALESCE(json_agg(json_build_object('id', skill.id, 'name', skill.name)) FILTER (WHERE skill.id IS NOT NULL), '[]') AS skills`,
	).
		From("cases.service_catalog AS catalog").
		LeftJoin("cases.sla ON sla.id = catalog.sla_id").
		LeftJoin("cases.status ON status.id = catalog.status_id").
		LeftJoin("cases.close_reason ON close_reason.id = catalog.close_reason_id").
		LeftJoin("directory.wbt_user AS created_by_user ON created_by_user.id = catalog.created_by").
		LeftJoin("directory.wbt_user AS updated_by_user ON updated_by_user.id = catalog.updated_by").
		LeftJoin("call_center.cc_teams_catalog AS tc ON tc.catalog_id = catalog.id").
		LeftJoin("call_center.cc_teams AS team ON team.id = tc.team_id").
		LeftJoin("call_center.cc_skills_catalog AS sc ON sc.catalog_id = catalog.id").
		LeftJoin("call_center.cc_skills AS skill ON skill.id = sc.skill_id").
		Where(sq.Eq{"catalog.id": lookup.Id, "catalog.dc": ctx.Session.GetDomainId()}).
		GroupBy(
			"catalog.id",
			"sla.name",
			"status.name",
			"close_reason.name",
			"grp.name",
			"created_by_user.name",
			"updated_by_user.name",
		)

	// Convert the select query to SQL
	selectQuery, selectArgs, err := selectQueryBuilder.ToSql()
	if err != nil {
		return "", nil, err
	}

	// Combine update and select query
	query := updateQuery + "; " + selectQuery
	combinedArgs := append(args, selectArgs...)

	return store.CompactSQL(query), combinedArgs, nil
}

func NewCatalogStore(store store.Store) (store.CatalogStore, model.AppError) {
	if store == nil {
		return nil, model.NewInternalError("postgres.new_catalog.check.bad_arguments",
			"error creating Catalog interface to the service table, main store is nil")
	}
	return &CatalogStore{storage: store}, nil
}
