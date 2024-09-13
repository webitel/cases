package postgres

import (
	"encoding/json"
	"fmt"
	"log"
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
		&add.Id, &add.Name, &add.Description, &add.Prefix,
		&add.Code, &add.State,
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
		if count >= ctx.GetSize() {
			next = true
			break
		}

		// Initialize catalog and related fields
		catalog := &cases.Catalog{
			Sla:         &cases.Lookup{},
			Status:      &cases.Lookup{},
			CloseReason: &cases.Lookup{},
		}
		var createdBy, updatedBy cases.Lookup
		var createdAt, updatedAt time.Time
		var teamLookups, skillLookups []byte

		// Build scan arguments using the helper function
		scanArgs := s.buildCatalogScanArgs(
			catalog, &createdBy,
			&updatedBy, &createdAt,
			&updatedAt, &teamLookups,
			&skillLookups,
		)

		// Debug: print the scan arguments before scanning
		log.Printf("Scan Args: %+v", scanArgs)

		// Scan the result into the appropriate fields
		err = rows.Scan(scanArgs...)
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

		// Set timestamps and created/updated by fields
		catalog.CreatedAt = util.Timestamp(createdAt)
		catalog.UpdatedAt = util.Timestamp(updatedAt)
		catalog.CreatedBy = &createdBy
		catalog.UpdatedBy = &updatedBy

		catalogs = append(catalogs, catalog)
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

	// Check if rpc.Fields contains team_ids or skill_ids
	updateTeams := false
	updateSkills := false

	// Check if the fields exist in rpc.Fields
	for _, field := range ctx.Fields {
		switch field {
		case "teams":
			updateTeams = true
		case "skills":
			updateSkills = true
		}
	}

	// Handle teams and skills updates if rpc.Fields contain team_ids or skill_ids
	if updateTeams || updateSkills {
		// Initialize empty slices for teamIDs and skillIDs
		teamIDs := []int64{}
		skillIDs := []int64{}

		// If the user has provided team updates, extract team IDs
		if updateTeams {
			if len(lookup.Teams) > 0 {
				teamIDs = make([]int64, len(lookup.Teams))
				for i, team := range lookup.Teams {
					teamIDs[i] = team.Id
				}
			} // Else, teamIDs remains as an empty slice
		}

		// If the user has provided skill updates, extract skill IDs
		if updateSkills {
			if len(lookup.Skills) > 0 {
				skillIDs = make([]int64, len(lookup.Skills))
				for i, skill := range lookup.Skills {
					skillIDs[i] = skill.Id
				}
			} // Else, skillIDs remains as an empty slice
		}

		// Build query to update teams and skills
		query, args := s.buildUpdateTeamsAndSkillsQuery(
			ctx,
			lookup.Id,
			teamIDs,  // Pass empty slice if no team IDs are provided
			skillIDs, // Pass empty slice if no skill IDs are provided
			ctx.Session.GetUserId(),
			ctx.Time,
			ctx.Session.GetDomainId(),
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

func (s *CatalogStore) buildCreateCatalogQuery(ctx *model.CreateOptions, add *cases.Catalog) (string, []interface{}) {
	// Define arguments for the query
	args := []interface{}{
		add.Name,                  // $1: name (cannot be null)
		add.Description,           // $2: description (could be null)
		add.Prefix,                // $3: prefix (could be null)
		add.Code,                  // $4: code (could be null)
		ctx.Time,                  // $5: created_at, updated_at
		ctx.Session.GetUserId(),   // $6: created_by, updated_by
		add.Sla.Id,                // $7: sla_id (could be null)
		add.Status.Id,             // $8: status_id (could be null)
		add.CloseReason.Id,        // $9: close_reason_id (could be null)
		add.State,                 // $10: state (cannot be null)
		ctx.Session.GetDomainId(), // $11: domain ID (dc)
	}

	var teamIds []int64
	if len(add.Teams) > 0 {
		teamIds = make([]int64, len(add.Teams))
		for i, team := range add.Teams {
			teamIds[i] = team.Id
		}
	} else {
		teamIds = nil
	}
	args = append(args, pq.Array(teamIds)) // $12: team_ids (could be null)

	var skillIds []int64
	if len(add.Skills) > 0 {
		skillIds = make([]int64, len(add.Skills))
		for i, skill := range add.Skills {
			skillIds[i] = skill.Id
		}
	} else {
		skillIds = nil
	}
	args = append(args, pq.Array(skillIds)) // $13: skill_ids (could be null)

	// SQL query construction
	query := `
WITH inserted_catalog AS (
    INSERT INTO cases.service_catalog (
                                       name, description, prefix, code, created_at, created_by, updated_at,
                                       updated_by, sla_id, status_id, close_reason_id, state, dc
        ) VALUES ($1,
                  COALESCE(NULLIF($2, ''), NULL), -- Description (NULL if empty string)
                  COALESCE(NULLIF($3, ''), NULL), -- Prefix (NULL if empty string)
                  COALESCE(NULLIF($4, ''), NULL), -- Code (NULL if empty string)
                  $5, $6, $5, $6,
                  COALESCE(NULLIF($7, 0), NULL), -- SLA ID (NULL if 0)
                  COALESCE(NULLIF($8, 0), NULL), -- Status ID (NULL if 0)
                  COALESCE(NULLIF($9, 0), NULL), -- Close Reason ID (NULL if 0)
                  $10,
                  $11)
        RETURNING id, name, description, prefix, code, state, sla_id, status_id, close_reason_id,
            created_by, updated_by, created_at, updated_at),
     inserted_teams AS (
         INSERT INTO cases.team_catalog (catalog_id, team_id, created_by, updated_by, created_at, updated_at, dc)
             SELECT inserted_catalog.id, unnest(COALESCE(NULLIF($12::bigint[], '{}'), NULL)), $6, $6, $5, $5, $11
             FROM inserted_catalog
             RETURNING catalog_id, team_id),
     inserted_skills AS (
         INSERT INTO cases.skill_catalog (catalog_id, skill_id, created_by, updated_by, created_at, updated_at, dc)
             SELECT inserted_catalog.id, unnest(COALESCE(NULLIF($13::bigint[], '{}'), NULL)), $6, $6, $5, $5, $11
             FROM inserted_catalog
             RETURNING catalog_id, skill_id),
     teams_agg AS (SELECT inserted_teams.catalog_id,
                          json_agg(json_build_object('id', team.id, 'name', team.name)) AS teams
                   FROM inserted_teams
                            LEFT JOIN call_center.cc_team team ON team.id = inserted_teams.team_id
                   GROUP BY inserted_teams.catalog_id),
     skills_agg AS (SELECT inserted_skills.catalog_id,
                           json_agg(json_build_object('id', skill.id, 'name', skill.name)) AS skills
                    FROM inserted_skills
                             LEFT JOIN call_center.cc_skill skill ON skill.id = inserted_skills.skill_id
                    GROUP BY inserted_skills.catalog_id)
SELECT inserted_catalog.id,
       inserted_catalog.name,
       COALESCE(inserted_catalog.description, '')    AS description,       -- Return empty string if null
       COALESCE(inserted_catalog.prefix, '')         AS prefix,            -- Return empty string if null
       COALESCE(inserted_catalog.code, '')           AS code,              -- Return empty string if null
       inserted_catalog.state,
       inserted_catalog.created_at,
       inserted_catalog.updated_at,
       COALESCE(inserted_catalog.sla_id, 0)          AS sla_id,            -- Return 0 if null
       COALESCE(sla.name, '')                        AS sla_name,          -- Return empty string if null
       COALESCE(inserted_catalog.status_id, 0)       AS status_id,         -- Return 0 if null
       COALESCE(status.name, '')                     AS status_name,       -- Return empty string if null
       COALESCE(inserted_catalog.close_reason_id, 0) AS close_reason_id,   -- Return 0 if null
       COALESCE(close_reason.name, '')               AS close_reason_name, -- Return empty string if null
       COALESCE(inserted_catalog.created_by, 0)      AS created_by,        -- Return 0 if null
       COALESCE(created_by_user.name, '')            AS created_by_name,   -- Return empty string if null
       COALESCE(inserted_catalog.updated_by, 0)      AS updated_by,        -- Return 0 if null
       COALESCE(updated_by_user.name, '')            AS updated_by_name,   -- Return empty string if null
       COALESCE(teams_agg.teams, '[]')               AS teams,             -- Return empty array if null
       COALESCE(skills_agg.skills, '[]')             AS skills             -- Return empty array if null
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

func (s *CatalogStore) buildSearchCatalogQuery(ctx *model.SearchOptions) (string, []interface{}, error) {
	convertedIds := ctx.FieldsUtil.Int64SliceToStringSlice(ctx.IDs)
	ids := ctx.FieldsUtil.FieldsFunc(convertedIds, ctx.FieldsUtil.InlineFields)
	// Initialize query builder with Common Table Expressions (CTEs) for teams and skills aggregation
	queryBuilder := sq.Select(
		"catalog.id",
		"catalog.name",                       // Mandatory
		"catalog.prefix",                     // Mandatory
		"catalog.sla_id",                     // Mandatory
		"sla.name AS sla_name",               // Mandatory
		"catalog.status_id",                  // Mandatory
		"status.name AS status_name",         // Mandatory
		"COALESCE(catalog.code, '') AS code", // Optional
		"COALESCE(catalog.description, '') AS description",        // Optional
		"COALESCE(catalog.close_reason_id, 0) AS close_reason_id", // Optional
		"COALESCE(close_reason.name, '') AS close_reason_name",    // Optional
		"catalog.state AS state",
		"COALESCE(catalog.created_by, 0) AS created_by",
		"COALESCE(created_by_user.name, '') AS created_by_name",
		"COALESCE(catalog.updated_by, 0) AS updated_by",
		"COALESCE(updated_by_user.name, '') AS updated_by_name",
		"catalog.created_at AS created_at",
		"catalog.updated_at AS updated_at",
		"COALESCE(teams_agg.teams, '[]') AS teams",    // Aggregated teams from the CTE
		"COALESCE(skills_agg.skills, '[]') AS skills", // Aggregated skills from the CTE
		// Determine if the catalog has services
		`EXISTS (SELECT 1 FROM cases.service_catalog AS cs WHERE cs.root_id = catalog.id) AS has_services`,
	).
		From("cases.service_catalog AS catalog").
		LeftJoin("cases.sla ON sla.id = catalog.sla_id").
		LeftJoin("cases.status ON status.id = catalog.status_id").
		LeftJoin("cases.close_reason ON close_reason.id = catalog.close_reason_id").
		LeftJoin("directory.wbt_user AS created_by_user ON created_by_user.id = catalog.created_by").
		LeftJoin("directory.wbt_user AS updated_by_user ON updated_by_user.id = catalog.updated_by").
		LeftJoin("teams_agg ON teams_agg.catalog_id = catalog.id").   // Join teams aggregation
		LeftJoin("skills_agg ON skills_agg.catalog_id = catalog.id"). // Join skills aggregation
		Where(sq.Eq{"catalog.root_id": nil}).                         // Fetch only top-level catalogs (where root_id is NULL)
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

	if len(ids) > 0 {
		queryBuilder = queryBuilder.Where(sq.Eq{"catalog.id": ids})
	}

	// Apply sorting
	for _, sort := range ctx.Sort {
		queryBuilder = queryBuilder.OrderBy(sort)
	}

	size := ctx.GetSize()
	page := ctx.Page

	// Apply offset only if page > 1
	if ctx.Page > 1 {
		queryBuilder = queryBuilder.Offset(uint64((page - 1) * size))
	}

	// Apply limit
	if size != -1 {
		queryBuilder = queryBuilder.Limit(uint64(size + 1)) // Request one more record to check if there's a next page
	}

	// Build SQL query with CTEs for teams and skills aggregation
	query, args, err := queryBuilder.Prefix(`
WITH inserted_teams AS (SELECT catalog_id, team_id
                        FROM cases.team_catalog),
     teams_agg AS (SELECT inserted_teams.catalog_id,
                          json_agg(json_build_object('id', team.id, 'name', team.name)) AS teams
                   FROM inserted_teams
                            LEFT JOIN call_center.cc_team team ON team.id = inserted_teams.team_id
                   GROUP BY inserted_teams.catalog_id),
     inserted_skills AS (SELECT catalog_id, skill_id
                         FROM cases.skill_catalog),
     skills_agg AS (SELECT inserted_skills.catalog_id,
                           json_agg(json_build_object('id', skill.id, 'name', skill.name)) AS skills
                    FROM inserted_skills
                             LEFT JOIN call_center.cc_skill skill ON skill.id = inserted_skills.skill_id
                    GROUP BY inserted_skills.catalog_id)
	`).ToSql()
	if err != nil {
		return "", nil, model.NewInternalError("postgres.catalog.query_build_error", err.Error())
	}

	return store.CompactSQL(query), args, nil
}

func (s *CatalogStore) buildUpdateTeamsAndSkillsQuery(ctx *model.UpdateOptions, catalogID int64, teamIDs, skillIDs []int64, updatedBy int64, updatedAt time.Time, domainID int64) (string, []interface{}) {
	args := []interface{}{
		catalogID, // $1: catalog_id
		updatedBy, // $2: updated_by
		domainID,  // $3: dc (domain context)
		updatedAt, // $4: timestamp for updated_at
	}

	// Initialize base query
	query := `WITH `

	// Check if "teams" is in ctx.Fields, even if teamIDs is empty
	if util.FieldExists("teams", ctx.Fields) {
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
		args = append(args, pq.Array(teamIDs)) // Append team IDs to args (even if empty)
	} else {
		// Pass an empty array if "teams" is not provided to avoid null issues
		args = append(args, pq.Array([]int64{}))
	}

	// Check if "skills" is in ctx.Fields, even if skillIDs is empty
	if util.FieldExists("skills", ctx.Fields) {
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
		)`
		args = append(args, pq.Array(skillIDs)) // Append skill IDs to args (even if empty)
	} else {
		// Pass an empty array if "skills" is not provided to avoid null issues
		args = append(args, pq.Array([]int64{}))
	}

	// Finish the query with UNION logic for teams and skills
	query += `
	SELECT COUNT(*)
	FROM (
		` + func() string {
		var result string
		if util.FieldExists("teams", ctx.Fields) {
			result += `SELECT catalog_id FROM updated_teams UNION ALL SELECT catalog_id FROM deleted_teams`
		}
		if util.FieldExists("skills", ctx.Fields) {
			if util.FieldExists("teams", ctx.Fields) {
				result += ` UNION ALL `
			}
			result += `SELECT catalog_id FROM updated_skills UNION ALL SELECT catalog_id FROM deleted_skills`
		}
		return result
	}() + `
	) AS total_affected;
	`

	// Return the constructed query and arguments
	return store.CompactSQL(query), args
}

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
		case "description":
			updateQueryBuilder = updateQueryBuilder.Set("description", lookup.Description)
		case "prefix":
			updateQueryBuilder = updateQueryBuilder.Set("prefix", lookup.Prefix)
		case "code":
			updateQueryBuilder = updateQueryBuilder.Set("code", lookup.Code)
		case "state":
			updateQueryBuilder = updateQueryBuilder.Set("state", lookup.State)
		case "sla_id":
			updateQueryBuilder = updateQueryBuilder.Set("sla_id", lookup.Sla.Id)
		case "status_id":
			updateQueryBuilder = updateQueryBuilder.Set("status_id", lookup.Status.Id)
		case "close_reason_id":
			updateQueryBuilder = updateQueryBuilder.Set("close_reason_id", lookup.CloseReason.Id)
		}
	}

	// Convert the update query to SQL
	updateSQL, args, err := updateQueryBuilder.ToSql()
	if err != nil {
		return "", nil, err
	}

	// Combine the update query with the select query using the WITH clause
	query := fmt.Sprintf(`
		WITH updated_catalog AS (%s
			RETURNING id, name, created_at, updated_at, sla_id, created_by, updated_by, status_id, close_reason_id)
SELECT catalog.id,
       catalog.name,
       catalog.created_at,
       catalog.sla_id,
       sla.name,
       catalog.status_id,
       status.name,
       catalog.close_reason_id,
       close_reason.name,
       catalog.created_by,
       created_by_user.name                               AS created_by_name,
       catalog.updated_by,
       updated_by_user.name                               AS updated_by_name,
       catalog.updated_at,
       COALESCE((SELECT json_agg(json_build_object('id', team.id, 'name', team.name))
                 FROM cases.team_catalog ts
                          LEFT JOIN call_center.cc_team team ON team.id = ts.team_id
                 WHERE ts.catalog_id = catalog.id), '[]') AS teams,
       COALESCE((SELECT json_agg(json_build_object('id', skill.id, 'name', skill.name))
                 FROM cases.skill_catalog ss
                          LEFT JOIN call_center.cc_skill skill ON skill.id = ss.skill_id
                 WHERE ss.catalog_id = catalog.id), '[]') AS skills
FROM updated_catalog AS catalog
         LEFT JOIN cases.sla ON sla.id = catalog.sla_id
         LEFT JOIN cases.status ON status.id = catalog.status_id
         LEFT JOIN cases.close_reason ON close_reason.id = catalog.close_reason_id
         LEFT JOIN directory.wbt_user AS created_by_user ON created_by_user.id = catalog.created_by
         LEFT JOIN directory.wbt_user AS updated_by_user ON updated_by_user.id = catalog.updated_by
GROUP BY catalog.id, catalog.name, catalog.created_at, catalog.sla_id, sla.name, catalog.status_id,
         status.name, catalog.close_reason_id, close_reason.name, catalog.created_by, created_by_user.name,
         catalog.updated_by, updated_by_user.name, catalog.updated_at;
	`, updateSQL)

	// Return the final combined query and arguments
	return store.CompactSQL(query), args, nil
}

// buildCatalogScanArgs prepares scan arguments for populating a Catalog object.
func (s *CatalogStore) buildCatalogScanArgs(
	catalog *cases.Catalog, // The catalog object to populate
	createdBy, updatedBy *cases.Lookup, // Lookup objects for created_by and updated_by
	createdAt, updatedAt *time.Time, // Temporary variables for created_at and updated_at
	teamLookups, skillLookups *[]byte, // Byte arrays for teams and skills (as JSON or binary)
) []interface{} {
	return []interface{}{
		// Catalog ID
		&catalog.Id, // Catalog ID

		// Catalog metadata
		&catalog.Name,   // Catalog name
		&catalog.Prefix, // Catalog prefix

		// SLA fields
		&catalog.Sla.Id,   // SLA ID
		&catalog.Sla.Name, // SLA name

		// Status fields
		&catalog.Status.Id,   // Status ID
		&catalog.Status.Name, // Status name

		// Code and description
		&catalog.Code,        // Catalog code (optional, can be null)
		&catalog.Description, // Catalog description (optional, can be null)

		// Close reason fields
		&catalog.CloseReason.Id,   // Close reason ID
		&catalog.CloseReason.Name, // Close reason name

		// State field
		&catalog.State, // Catalog state (active/inactive)

		// Created by and updated by fields (Lookup)
		&createdBy.Id,   // Created by user ID
		&createdBy.Name, // Created by user name
		&updatedBy.Id,   // Updated by user ID
		&updatedBy.Name, // Updated by user name

		// Timestamps
		createdAt, // Created at timestamp
		updatedAt, // Updated at timestamp

		// Teams and skills lookups (usually in JSON or binary format)
		teamLookups,  // Team lookups (can be null if no teams)
		skillLookups, // Skill lookups (can be null if no skills)

		// Additional fields
		&catalog.HasServices, // Whether the catalog has related services
	}
}

func NewCatalogStore(store store.Store) (store.CatalogStore, model.AppError) {
	if store == nil {
		return nil, model.NewInternalError("postgres.new_catalog.check.bad_arguments",
			"error creating Catalog interface to the service table, main store is nil")
	}
	return &CatalogStore{storage: store}, nil
}
