package postgres

import (
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

type CatalogStore struct {
	storage *Store
}

// Create implements store.CatalogStore.
func (s *CatalogStore) Create(rpc options.CreateOptions, add *cases.Catalog) (*cases.Catalog, error) {
	// Establish a connection to the database
	db, dbErr := s.storage.Database()
	if dbErr != nil {
		return nil, dberr.NewDBInternalError("postgres.catalog.create.database_connection_error", dbErr)
	}

	// Build the combined query for inserting Catalog, teams, and skills
	query, args := s.buildCreateCatalogQuery(rpc, add)

	var (
		createdByLookup, updatedByLookup cases.Lookup
		createdAt, updatedAt             time.Time
		teamLookups, skillLookups        []byte
	)

	err := db.QueryRow(rpc, query, args...).Scan(
		&add.Id, &add.Name, &add.Description, &add.Prefix,
		&add.Code, &add.State,
		&createdAt, &updatedAt,
		&add.Sla.Id, &add.Sla.Name,
		&add.Status.Id, &add.Status.Name,
		&add.CloseReasonGroup.Id, &add.CloseReasonGroup.Name,
		&createdByLookup.Id, &createdByLookup.Name,
		&updatedByLookup.Id, &updatedByLookup.Name,
		&teamLookups,  // JSON array for teams
		&skillLookups, // JSON array for skills
	)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.catalog.create.scan_error", err)
	}

	// Unmarshal the JSON arrays into the Lookup slices
	if err := json.Unmarshal(teamLookups, &add.Teams); err != nil {
		return nil, dberr.NewDBInternalError("postgres.catalog.create.unmarshal_teams_error", err)
	}
	if err := json.Unmarshal(skillLookups, &add.Skills); err != nil {
		return nil, dberr.NewDBInternalError("postgres.catalog.create.unmarshal_skills_error", err)
	}

	// Prepare the Catalog to return
	add.CreatedAt = util.Timestamp(createdAt)
	add.UpdatedAt = util.Timestamp(updatedAt)
	add.CreatedBy = &createdByLookup
	add.UpdatedBy = &updatedByLookup

	// Return the created Catalog
	return add, nil
}

func (s *CatalogStore) buildCreateCatalogQuery(rpc options.CreateOptions, add *cases.Catalog) (string, []interface{}) {
	// Define arguments for the query
	args := []any{
		add.Name,                        // $1: name (cannot be null)
		add.Description,                 // $2: description (could be null)
		add.Prefix,                      // $3: prefix (could be null)
		add.Code,                        // $4: code (could be null)
		rpc.RequestTime(),               // $5: created_at, updated_at
		rpc.GetAuthOpts().GetUserId(),   // $6: created_by, updated_by
		add.Sla.Id,                      // $7: sla_id (could be null)
		add.Status.Id,                   // $8: status_id (could be null)
		add.CloseReasonGroup.Id,         // $9: close_reason_id (could be null)
		add.State,                       // $10: state (cannot be null)
		rpc.GetAuthOpts().GetDomainId(), // $11: domain ID (dc)
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
                                       updated_by, sla_id, status_id, close_reason_group_id, state, dc
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
        RETURNING id, name, description, prefix, code, state, sla_id, status_id, close_reason_group_id,
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
       COALESCE(inserted_catalog.close_reason_group_id, 0) AS close_reason_group_id,   -- Return 0 if null
       COALESCE(close_reason_group.name, '')               AS close_reason_group_name, -- Return empty string if null
       COALESCE(inserted_catalog.created_by, 0)      AS created_by,        -- Return 0 if null
       COALESCE(created_by_user.name, '')            AS created_by_name,   -- Return empty string if null
       COALESCE(inserted_catalog.updated_by, 0)      AS updated_by,        -- Return 0 if null
       COALESCE(updated_by_user.name, '')            AS updated_by_name,   -- Return empty string if null
       COALESCE(teams_agg.teams, '[]')               AS teams,             -- Return empty array if null
       COALESCE(skills_agg.skills, '[]')             AS skills             -- Return empty array if null
FROM inserted_catalog
         LEFT JOIN cases.sla ON sla.id = inserted_catalog.sla_id
         LEFT JOIN cases.status ON status.id = inserted_catalog.status_id
         LEFT JOIN cases.close_reason_group ON close_reason_group.id = inserted_catalog.close_reason_group_id
         LEFT JOIN directory.wbt_user created_by_user ON created_by_user.id = inserted_catalog.created_by
         LEFT JOIN directory.wbt_user updated_by_user ON updated_by_user.id = inserted_catalog.updated_by
         LEFT JOIN teams_agg ON teams_agg.catalog_id = inserted_catalog.id
         LEFT JOIN skills_agg ON skills_agg.catalog_id = inserted_catalog.id;
`

	return util2.CompactSQL(query), args
}

// Delete implements store.CatalogStore.
func (s *CatalogStore) Delete(rpc options.DeleteOptions) error {
	// Establish a connection to the database
	db, dbErr := s.storage.Database()
	if dbErr != nil {
		return dberr.NewDBInternalError("postgres.catalog.delete.db_connection_error", dbErr)
	}

	// Ensure that there are IDs to delete
	if len(rpc.GetIDs()) == 0 {
		return dberr.NewDBNoRowsError("postgres.catalog.delete.no_ids_provided")
	}

	// Build the delete query
	query, args := s.buildDeleteCatalogQuery(rpc)

	// Execute the delete query
	res, err := db.Exec(rpc, query, args...)
	if err != nil {
		return dberr.NewDBInternalError("postgres.catalog.delete.execution_error", err)
	}

	// Check how many rows were affected
	if res.RowsAffected() == 0 {
		return dberr.NewDBNoRowsError("postgres.catalog.delete.no_rows_deleted")
	}

	return nil
}

// Helper method to build the delete query for Catalog
func (s *CatalogStore) buildDeleteCatalogQuery(rpc options.DeleteOptions) (string, []any) {
	// Build the SQL query using the provided IDs in rpc.IDs
	query := `
		DELETE FROM cases.service_catalog
		WHERE id = ANY($1) AND dc = $2
	`

	// Use the array of IDs and domain ID (dc) for the deletion
	args := []any{
		pq.Array(rpc.GetIDs()),          // $1: array of catalog IDs to delete
		rpc.GetAuthOpts().GetDomainId(), // $2: domain ID to ensure proper scoping
	}

	return util2.CompactSQL(query), args
}

// List implements store.CatalogStore.
func (s *CatalogStore) List(
	rpc options.SearchOptions,
	depth int64,
	subfields []string,
	hasSubservices bool,
) (*cases.CatalogList, error) {
	// 1. Connect to DB
	db, dbErr := s.storage.Database()
	if dbErr != nil {
		return nil, dberr.NewDBInternalError("postgres.catalog.list.database_connection_error", dbErr)
	}

	// 3. Build SQL query
	query, args, err := s.buildSearchCatalogQuery(rpc, depth, subfields, hasSubservices)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.catalog.list.query_build_error", err)
	}

	// 4. Execute the query
	rows, err := db.Query(rpc, query, args...)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.catalog.list.query_execution_error", err)
	}
	defer rows.Close()

	// 5. Prepare containers
	var (
		catalogs []*cases.Catalog
		lCount   int
		next     bool
		fetchAll = rpc.GetSize() == -1
	)

	// 6. Single-pass read
	for rows.Next() {
		// If not fetching all, check the size limit
		if !fetchAll && lCount >= rpc.GetSize() {
			next = true
			break
		}

		// Prepare the top-level Catalog object
		catalog := &cases.Catalog{}

		// Initialize lookups if requested
		if util.ContainsField(rpc.GetFields(), "sla") {
			catalog.Sla = &cases.Lookup{}
		}
		if util.ContainsField(rpc.GetFields(), "status") {
			catalog.Status = &cases.Lookup{}
		}
		if util.ContainsField(rpc.GetFields(), "close_reason_group") {
			catalog.CloseReasonGroup = &cases.Lookup{}
		}
		if util.ContainsField(rpc.GetFields(), "teams") {
			catalog.Teams = []*cases.Lookup{}
		}
		if util.ContainsField(rpc.GetFields(), "skills") {
			catalog.Skills = []*cases.Lookup{}
		}
		if util.ContainsField(rpc.GetFields(), "services") {
			catalog.Service = []*cases.Service{}
		}
		if util.ContainsField(rpc.GetFields(), "created_by") {
			catalog.CreatedBy = &cases.Lookup{}
		}
		if util.ContainsField(rpc.GetFields(), "updated_by") {
			catalog.UpdatedBy = &cases.Lookup{}
		}

		var (
			createdAt, updatedAt              time.Time
			rootID                            int64
			teamData, skillData, servicesData string
		)

		// Build placeholders
		scanArgs, err := s.buildCatalogScanArgs(catalog, &createdAt, &updatedAt, &rootID, rpc.GetFields(), &teamData, &skillData, &servicesData)
		if err != nil {
			return nil, dberr.NewDBInternalError("postgres.catalog.list.scan_args_error", err)
		}

		// Single row scan
		if err := rows.Scan(scanArgs...); err != nil {
			return nil, dberr.NewDBInternalError("postgres.catalog.list.scan_error", err)
		}

		catalog.CreatedAt = util.Timestamp(createdAt)
		catalog.UpdatedAt = util.Timestamp(updatedAt)

		if util.ContainsField(rpc.GetFields(), "teams") && teamData != "" {
			var parsedTeams []*cases.Lookup
			if err := json.Unmarshal([]byte(teamData), &parsedTeams); err != nil {
				return nil, dberr.NewDBInternalError("postgres.catalog.list.teams_parse_error", err)
			}
			catalog.Teams = parsedTeams
		} else {
			catalog.Teams = []*cases.Lookup{}
		}

		if util.ContainsField(rpc.GetFields(), "skills") && skillData != "" {
			var parsedSkills []*cases.Lookup
			if err := json.Unmarshal([]byte(skillData), &parsedSkills); err != nil {
				return nil, dberr.NewDBInternalError("postgres.catalog.list.skills_parse_error", err)
			}
			catalog.Skills = parsedSkills
		} else {
			catalog.Skills = []*cases.Lookup{}
		}

		// Handle servicesData: Trim "{}" to ""
		if util.ContainsField(rpc.GetFields(), "services") && servicesData == "{}" {
			servicesData = ""
		}

		if util.ContainsField(rpc.GetFields(), "services") && servicesData != "" {
			// Parse JSON into a generic map
			var rawServices []map[string]interface{}
			if err := json.Unmarshal([]byte(servicesData), &rawServices); err != nil {
				return nil, dberr.NewDBInternalError("postgres.catalog.list.services_parse_error", err)
			}

			var parsedServices []*cases.Service
			for _, raw := range rawServices {
				// Skip empty objects
				if len(raw) == 0 {
					continue
				}

				service := &cases.Service{}

				// Directly assign simple fields
				if id, ok := raw["id"].(float64); ok {
					service.Id = int64(id)
				}
				if code, ok := raw["code"].(string); ok {
					service.Code = code
				}
				if name, ok := raw["name"].(string); ok {
					service.Name = name
				}
				if rootID, ok := raw["root_id"].(float64); ok {
					service.RootId = int64(rootID)
				}
				if description, ok := raw["description"].(string); ok {
					service.Description = description
				}
				// Parse created_at and updated_at from Unix milliseconds
				if createdAt, ok := raw["created_at"].(float64); ok {
					service.CreatedAt = int64(createdAt)
				}
				if updatedAt, ok := raw["updated_at"].(float64); ok {
					service.UpdatedAt = int64(updatedAt)
				}
				if catalogId, ok := raw["catalog_id"].(float64); ok {
					service.CatalogId = int64(catalogId)
				}

				// Map SLA to Lookup
				if slaID, ok := raw["sla_id"].(float64); ok {
					if slaName, ok := raw["sla_name"].(string); ok {
						service.Sla = &cases.Lookup{
							Id:   int64(slaID),
							Name: slaName,
						}
					}
				}

				// Map Group to Lookup
				if groupID, ok := raw["group_id"].(float64); ok {
					if groupName, ok := raw["group_name"].(string); ok {
						if groupType, ok := raw["group_type"].(string); ok {
							if groupID == 0 {
								service.Group = nil
							} else {
								service.Group = &cases.ExtendedLookup{
									Id:   int64(groupID),
									Name: groupName,
									Type: strings.ToUpper(groupType),
								}
							}
						}
					}
				}

				// Map Assignee to Lookup
				if assigneeID, ok := raw["assignee_id"].(float64); ok {
					if assigneeName, ok := raw["assignee_name"].(string); ok {
						service.Assignee = &cases.Lookup{
							Id:   int64(assigneeID),
							Name: assigneeName,
						}
					}
				}

				// Map CreatedBy to Lookup
				if createdByID, ok := raw["created_by"].(float64); ok {
					if createdByName, ok := raw["created_by_name"].(string); ok {
						service.CreatedBy = &cases.Lookup{
							Id:   int64(createdByID),
							Name: createdByName,
						}
					}
				}

				// Map UpdatedBy to Lookup
				if updatedByID, ok := raw["updated_by"].(float64); ok {
					if updatedByName, ok := raw["updated_by_name"].(string); ok {
						service.UpdatedBy = &cases.Lookup{
							Id:   int64(updatedByID),
							Name: updatedByName,
						}
					}
				}

				// Map searched field
				if searched, ok := raw["searched"].(bool); ok {
					service.Searched = searched
				}

				// Map state field
				if state, ok := raw["state"].(bool); ok {
					service.State = state
				}

				parsedServices = append(parsedServices, service)
			}

			catalog.Service = parsedServices
		}

		if rootID == 0 {
			catalogs = append(catalogs, catalog)
			lCount++
		}

	}

	// 12. Check for row errors
	if rows.Err() != nil {
		return nil, dberr.NewDBInternalError("postgres.catalog.list.query_rows_error", rows.Err())
	}

	// 13. Build nested hierarchy if “services” is requested
	if util.ContainsField(rpc.GetFields(), "services") {
		// For each top-level catalog, nest subservices
		for _, cat := range catalogs {
			nested, err := s.nestServicesByRootID(cat.Id, cat.Service)
			if err != nil {
				return nil, dberr.NewDBInternalError("postgres.catalog.list.nesting_services_error", err)
			}
			cat.Service = nested
		}
	}

	// 14. Return the final result
	return &cases.CatalogList{
		Page:  int32(rpc.GetPage()),
		Next:  next,
		Items: catalogs,
	}, nil
}

func (s *CatalogStore) buildCatalogScanArgs(
	catalog *cases.Catalog,
	createdAt, updatedAt *time.Time,
	rootId *int64,
	rpcFields []string,
	teamData, skillData, serviceData *string) (scanArgs []interface{}, err error) {
	for _, field := range rpcFields {
		switch field {

		// ---------- Catalog Fields ----------
		case "id":
			scanArgs = append(scanArgs, &catalog.Id)

		case "name":
			scanArgs = append(scanArgs, &catalog.Name)

		case "description":
			scanArgs = append(scanArgs, &catalog.Description)

		case "root_id":
			scanArgs = append(scanArgs, rootId)

		case "prefix":
			scanArgs = append(scanArgs, &catalog.Prefix)

		case "code":
			scanArgs = append(scanArgs, &catalog.Code)

		case "state":
			scanArgs = append(scanArgs, &catalog.State)

		case "sla":
			scanArgs = append(scanArgs, &catalog.Sla.Id, &catalog.Sla.Name)

		case "status":
			scanArgs = append(scanArgs, &catalog.Status.Id, &catalog.Status.Name)

		case "close_reason_group":
			scanArgs = append(scanArgs, &catalog.CloseReasonGroup.Id, &catalog.CloseReasonGroup.Name)

		case "created_by":
			scanArgs = append(scanArgs, &catalog.CreatedBy.Id, &catalog.CreatedBy.Name)

		case "updated_by":
			scanArgs = append(scanArgs, &catalog.UpdatedBy.Id, &catalog.UpdatedBy.Name)

		case "created_at":
			scanArgs = append(scanArgs, createdAt)

		case "updated_at":
			scanArgs = append(scanArgs, updatedAt)

		case "teams":
			scanArgs = append(scanArgs, teamData)

		case "skills":
			scanArgs = append(scanArgs, skillData)

			// ---------- Searched Field ----------
		case "searched":
			scanArgs = append(scanArgs, &catalog.Searched)

			// ---------- Service Fields ----------
		case "services":
			scanArgs = append(scanArgs, serviceData)
		}
	}

	return scanArgs, nil
}

func (s *CatalogStore) nestServicesByRootID(
	rootCatalogID int64,
	services []*cases.Service,
) ([]*cases.Service, error) {
	// Map services by their RootId
	serviceMap := make(map[int64][]*cases.Service)
	for _, svc := range services {
		// Group services by their RootId
		serviceMap[svc.RootId] = append(serviceMap[svc.RootId], svc)
	}

	// Start building the hierarchy from the rootCatalogID
	hierarchy := s.buildServiceHierarchy(rootCatalogID, serviceMap)
	return hierarchy, nil
}

func (s *CatalogStore) buildServiceHierarchy(
	rootID int64,
	serviceMap map[int64][]*cases.Service,
) []*cases.Service {
	// Retrieve all children of the current rootID
	children := serviceMap[rootID]
	for _, child := range children {
		// Recursively attach sub-services to the current child
		child.Service = s.buildServiceHierarchy(child.Id, serviceMap)
	}
	return children
}

func (s *CatalogStore) buildSearchCatalogQuery(
	rpc options.SearchOptions,
	depth int64,
	subfields []string,
	hasSubservices bool,
) (string, []interface{}, error) {
	// Catalog-level field map (removed "services": "" entry)
	fieldMap := map[string]string{
		"id":                 "catalog.id",
		"name":               "catalog.name",
		"prefix":             "COALESCE(catalog.prefix, '') AS prefix",
		"sla":                "COALESCE(catalog.sla_id, 0) AS sla_id, COALESCE(sla.name, '') AS sla_name",
		"group":              "COALESCE(catalog.group_id, 0) AS group_id, COALESCE(group_lookup.name, '') AS group_name",
		"assignee":           "COALESCE(catalog.assignee_id, 0) AS assignee_id, COALESCE(assignee_user.name, '') AS assignee_name",
		"status":             "COALESCE(catalog.status_id, 0) AS status_id, COALESCE(status.name, '') AS status_name",
		"code":               "COALESCE(catalog.code, '') AS code",
		"description":        "COALESCE(catalog.description, '') AS description",
		"close_reason_group": "COALESCE(catalog.close_reason_group_id, 0) AS close_reason_group_id, COALESCE(close_reason_group.name, '') AS close_reason_name",
		"state":              "catalog.state AS state",
		"created_by":         "COALESCE(catalog.created_by, 0) AS created_by, COALESCE(created_by_user.name, '') AS created_by_name",
		"updated_by":         "COALESCE(catalog.updated_by, 0) AS updated_by, COALESCE(updated_by_user.name, '') AS updated_by_name",
		"created_at":         "catalog.created_at AS created_at",
		"updated_at":         "catalog.updated_at AS updated_at",
		"teams":              "COALESCE(JSONB_AGG(DISTINCT JSONB_BUILD_OBJECT('id', teams.team_id, 'name', teams.team_name)) FILTER (WHERE teams.team_id IS NOT NULL), '[]') AS team_data",
		"skills":             "COALESCE(JSONB_AGG(DISTINCT JSONB_BUILD_OBJECT('id', skills.skill_id, 'name', skills.skill_name)) FILTER (WHERE skills.skill_id IS NOT NULL), '[]') AS skill_data",
		"root_id":            "COALESCE(catalog.root_id, 0) AS root_id",
	}

	// Flags for applying CTEs or conditional joins
	selectFlags := map[string]bool{
		"teams":    false,
		"skills":   false,
		"search":   false,
		"services": false,
	}

	if name, ok := rpc.GetFilter("name").(string); ok && len(name) > 0 {
		selectFlags["search"] = true
	}

	var selectedFields []string

	// 1) Build the catalog columns from rpc.Fields
	for _, field := range rpc.GetFields() {
		if mappedField, ok := fieldMap[field]; ok {
			cols := strings.Split(mappedField, ", ")
			selectedFields = append(selectedFields, cols...)
			switch field {
			case "prefix":
				// If prefix => maybe you want something else?
			case "teams":
				selectFlags["teams"] = true
			case "skills":
				selectFlags["skills"] = true
			}
		}
	}

	// 2) Check if the user requested "services" in rpc.Fields.
	// If so, we build columns from serviceFields
	// Check if the user requested "services" in rpc.Fields.
	if util.ContainsField(rpc.GetFields(), "services") {
		selectFlags["services"] = true

		// Build dynamic JSONB_AGG for services
		serviceAgg := buildServiceJSONBAgg(subfields, selectFlags["search"])
		selectedFields = append(selectedFields, serviceAgg)
	}

	if selectFlags["search"] {
		// Add the "searched" field dynamically to the field map
		fieldMap["searched"] = "CASE WHEN catalog.id IN (SELECT catalog_id FROM search_catalog) THEN true ELSE false END::boolean AS searched"
		selectedFields = append(selectedFields, fieldMap["searched"])
	}

	// Named parameters
	params := map[string]interface{}{
		"dc":     rpc.GetAuthOpts().GetDomainId(),
		"limit":  rpc.GetSize() + 1,
		"offset": (rpc.GetPage() - 1) * rpc.GetSize(),
	}

	// Build the base query
	queryBuilder := sq.Select(strings.Join(selectedFields, ", ")).
		From("limited_catalogs AS catalog").
		Where("catalog.dc = :dc").
		PlaceholderFormat(sq.Dollar)

	// Add the standard joins for catalog fields
	queryBuilder = queryBuilder.
		LeftJoin("cases.sla AS sla ON sla.id = catalog.sla_id").
		LeftJoin("contacts.group AS group_lookup ON group_lookup.id = catalog.group_id").
		LeftJoin("directory.wbt_user AS assignee_user ON assignee_user.id = catalog.assignee_id").
		LeftJoin("cases.status ON status.id = catalog.status_id").
		LeftJoin("cases.close_reason_group ON close_reason_group.id = catalog.close_reason_group_id").
		LeftJoin("directory.wbt_user AS created_by_user ON created_by_user.id = catalog.created_by").
		LeftJoin("directory.wbt_user AS updated_by_user ON updated_by_user.id = catalog.updated_by")

	// 4) Add conditional WHERE clause for search
	if selectFlags["search"] {
		queryBuilder = queryBuilder.Where(
			sq.Or{
				sq.Expr("catalog.id IN (SELECT catalog_id FROM search_catalog)"),
				sq.Expr("catalog.id IN (SELECT id FROM service_hierarchy)"),
			},
		)
	}

	// 5) Conditional joins for teams, skills
	if selectFlags["teams"] {
		queryBuilder = queryBuilder.LeftJoin("teams ON teams.catalog_id = catalog.id")
	}
	if selectFlags["skills"] {
		queryBuilder = queryBuilder.LeftJoin("skills ON skills.catalog_id = catalog.id")
	}

	// 6) If "services" is requested, join services_hierarchy
	if selectFlags["services"] {
		queryBuilder = queryBuilder.LeftJoin("service_hierarchy ON service_hierarchy.catalog_id = catalog.id")
	}

	//FIXME make services json building in separate cte -> then make
	// SELECT ... services_cte.services
	// AND JSONB_ARRAY_LENGTH(services_cte.services) > 0

	// 7) State + ID filters
	if state, ok := rpc.GetFilter("state").(bool); ok {
		params["state"] = state
		queryBuilder = queryBuilder.Where("catalog.state = :state")

		// Add EXISTS condition if hasSubservices is true
		if hasSubservices && state {
			queryBuilder = queryBuilder.Where(`EXISTS (
            SELECT 1
            FROM cases.service_catalog sc
            WHERE sc.root_id = catalog.id
              AND sc.state = :state
        )`)
		}
	}

	teamFilter, teamFilterFound := rpc.GetFilter("team").(int64)
	skillsFilter, skillFilterFound := rpc.GetFilter("skills").([]int64)
	if teamFilterFound || skillFilterFound {
		or := sq.Or{}
		if teamFilter > 0 {
			or = append(or, sq.Expr("catalog.id IN (SELECT DISTINCT catalog_id FROM cases.team_catalog WHERE team_id = :team)"))
			params["team"] = teamFilter
		}
		if len(skillsFilter) > 0 {
			or = append(or, sq.Expr("catalog.id IN (SELECT DISTINCT catalog_id FROM cases.skill_catalog WHERE skill_id = ANY(:skills))", skillsFilter))
			params["skills"] = skillsFilter
		}
		or = append(or, sq.Expr("(NOT EXISTS(SELECT catalog_id FROM cases.team_catalog WHERE catalog_id = catalog.id) AND NOT EXISTS(SELECT catalog_id FROM cases.skill_catalog WHERE catalog_id = catalog.id))"))
		queryBuilder = queryBuilder.Where(or)
	}

	// Add condition for rpc.IDs if provided
	if len(rpc.GetIDs()) > 0 {
		params["id"] = rpc.GetIDs()[0]
		queryBuilder = queryBuilder.Where("catalog.id = :id")
	}

	// Search CTE
	var searchQ string
	var searchCondition string
	var hasSubservicesFilter string

	// Define hasSubservicesFilter conditionally
	if hasSubservices {
		hasSubservicesFilter = `
			AND EXISTS (
				SELECT 1
				FROM cases.service_catalog lc
				WHERE lc.root_id = catalog.id
			)
		`
	}

	if selectFlags["search"] {
		if name, ok := rpc.GetFilter("name").(string); ok {
			params["name"] = "%" + strings.Join(util.Substring(name), "%") + "%"
		}

		searchQ = `
		search_catalog AS (
			SELECT
				catalog.id AS catalog_id,
				catalog.catalog_id AS service_catalog_id,
				CASE
					WHEN catalog.catalog_id IS NULL THEN catalog.id
					ELSE catalog.catalog_id
				END AS target_catalog_id,
				catalog.id AS searched_id
			FROM cases.service_catalog catalog
			WHERE catalog.name ILIKE :name
		),`
		searchCondition = "AND id IN (SELECT target_catalog_id FROM search_catalog)"
	}

	queryBuilder = util2.ApplyPaging(rpc.GetPage(), rpc.GetSize(), queryBuilder)

	// Prefix query
	prefixQuery := fmt.Sprintf(`
	limited_catalogs AS (
		SELECT *
		FROM cases.service_catalog catalog
		WHERE root_id IS NULL
		%s -- Conditionally include the search condition
		%s -- Conditionally include hasSubservicesFilter
	),
`, searchCondition, hasSubservicesFilter)

	// Add the prefix query with or without search_catalog based on search condition
	if selectFlags["search"] {
		queryBuilder = queryBuilder.Prefix(fmt.Sprintf(`WITH %s%s`, searchQ, prefixQuery))
	} else {
		queryBuilder = queryBuilder.Prefix(fmt.Sprintf(`WITH %s`, prefixQuery))
	}

	// Sorting logic
	queryBuilder = applySorting(queryBuilder, rpc)

	// 10) If we need CTE(s)
	if selectFlags["teams"] || selectFlags["skills"] || selectFlags["services"] || selectFlags["search"] {
		prefixQuery := buildCTEs(selectFlags, depth, subfields, rpc.GetFilters())
		queryBuilder = queryBuilder.Prefix(prefixQuery)
	}
	// Add GROUP BY for catalog fields
	groupByFields := buildCatalogGroupByFields(rpc.GetFields())

	if len(groupByFields) != 0 {
		queryBuilder = queryBuilder.GroupBy(strings.Join(groupByFields, ", "))
	}

	// 11) Build final query
	sqlQuery, _, err := queryBuilder.ToSql()
	if err != nil {
		return "", nil, dberr.NewDBInternalError("postgres.catalog.query_build_error", err)
	}

	q, args, err := util2.BindNamed(sqlQuery, params)
	if err != nil {
		return "", nil, fmt.Errorf("failed to bind named parameters: %w", err)
	}

	return util2.CompactSQL(q), args, nil
}

func buildServiceJSONBAgg(subfields []string, searched bool) string {
	var jsonFields strings.Builder

	// Start building JSONB_AGG with JSONB_BUILD_OBJECT
	jsonFields.WriteString(`
COALESCE(
    JSONB_AGG(
        DISTINCT JSONB_BUILD_OBJECT(
`)

	// Conditionally append fields based on subfields
	if util.ContainsField(subfields, "id") {
		jsonFields.WriteString("'id', service_hierarchy.id,\n")
	}
	if util.ContainsField(subfields, "state") {
		jsonFields.WriteString("'state', service_hierarchy.state,\n")
	}
	if util.ContainsField(subfields, "name") {
		jsonFields.WriteString("'name', service_hierarchy.name,\n")
	}
	if util.ContainsField(subfields, "code") {
		jsonFields.WriteString("'code', service_hierarchy.code,\n")
	}
	if util.ContainsField(subfields, "description") {
		jsonFields.WriteString("'description', service_hierarchy.description,\n")
	}
	if util.ContainsField(subfields, "sla") {
		jsonFields.WriteString("'sla_id', COALESCE(service_hierarchy.sla_id, 0),\n")
		jsonFields.WriteString("'sla_name', COALESCE(service_hierarchy.sla_name, ''),\n")
	}
	if util.ContainsField(subfields, "group") {
		jsonFields.WriteString("'group_id', COALESCE(service_hierarchy.group_id, 0),\n")
		jsonFields.WriteString("'group_name', COALESCE(service_hierarchy.group_name, ''),\n")
		jsonFields.WriteString("'group_type', COALESCE(service_hierarchy.group_type, ''),\n")
	}
	if util.ContainsField(subfields, "assignee") {
		jsonFields.WriteString("'assignee_id', COALESCE(service_hierarchy.assignee_id, 0),\n")
		jsonFields.WriteString("'assignee_name', COALESCE(service_hierarchy.assignee_name, ''),\n")
	}
	if util.ContainsField(subfields, "created_by") {
		jsonFields.WriteString("'created_by', COALESCE(service_hierarchy.created_by, 0),\n")
		jsonFields.WriteString("'created_by_name', COALESCE(service_hierarchy.created_by_name, ''),\n")
	}
	if util.ContainsField(subfields, "updated_by") {
		jsonFields.WriteString("'updated_by', COALESCE(service_hierarchy.updated_by, 0),\n")
		jsonFields.WriteString("'updated_by_name', COALESCE(service_hierarchy.updated_by_name, ''),\n")
	}
	if util.ContainsField(subfields, "created_at") {
		// Convert created_at to milliseconds
		jsonFields.WriteString("'created_at', FLOOR(EXTRACT(EPOCH FROM service_hierarchy.created_at) * 1000),\n")
	}
	if util.ContainsField(subfields, "updated_at") {
		// Convert updated_at to milliseconds
		jsonFields.WriteString("'updated_at', FLOOR(EXTRACT(EPOCH FROM service_hierarchy.updated_at) * 1000),\n")
	}
	if util.ContainsField(subfields, "root_id") {
		jsonFields.WriteString("'root_id', service_hierarchy.root_id,\n")
	}
	if util.ContainsField(subfields, "catalog_id") {
		jsonFields.WriteString("'catalog_id', service_hierarchy.catalog_id,\n")
	}
	// Add the searched field if search is active
	if searched {
		jsonFields.WriteString("'searched', service_hierarchy.searched,\n")
	}

	// Remove the trailing comma and close JSONB_BUILD_OBJECT
	jsonAgg := strings.TrimSuffix(jsonFields.String(), ",\n")
	jsonAgg += `
        )
    ) FILTER (WHERE service_hierarchy.root_id IS NOT NULL),
    '{}'::jsonb
) AS services
`
	return jsonAgg
}

// applySorting applies dynamic sorting based on rpc.Sort and defaults to sorting by name in ascending order.
func applySorting(queryBuilder sq.SelectBuilder, rpc options.SearchOptions) sq.SelectBuilder {
	sortableFields := map[string]string{
		"name":               "catalog.name",
		"prefix":             "catalog.prefix",
		"sla":                "sla.name",
		"code":               "catalog.code",
		"status":             "status.name",
		"close_reason_group": "close_reason_group.name",
		"description":        "catalog.description",
		"state":              "catalog.state",
	}

	sortApplied := false

	sortField := rpc.GetSort()
	sortDirection := "ASC"
	if len(sortField) > 0 {
		switch sortField[0] {
		case '-':
			sortDirection = "DESC"
			sortField = sortField[1:]
		case '+':
			sortField = sortField[1:]
		}

		if dbField, exists := sortableFields[sortField]; exists {
			queryBuilder = queryBuilder.OrderBy(fmt.Sprintf("%s %s", dbField, sortDirection))
			sortApplied = true
		}
	}

	if !sortApplied {
		queryBuilder = queryBuilder.OrderBy("catalog.name ASC")
	}

	return queryBuilder
}

// buildCatalogGroupByFields constructs the GROUP BY fields for catalog-level columns conditionally.
func buildCatalogGroupByFields(requestedFields []string) []string {
	fieldMap := map[string]string{
		"id":                 "catalog.id",
		"name":               "catalog.name",
		"prefix":             "catalog.prefix",
		"sla":                "catalog.sla_id, sla.name",
		"group":              "catalog.group_id, group_lookup.name",
		"assignee":           "catalog.assignee_id, assignee_user.name",
		"status":             "catalog.status_id, status.name",
		"code":               "catalog.code",
		"description":        "catalog.description",
		"close_reason_group": "catalog.close_reason_group_id, close_reason_group.name",
		"state":              "catalog.state",
		"created_by":         "catalog.created_by, created_by_user.name",
		"updated_by":         "catalog.updated_by, updated_by_user.name",
		"created_at":         "catalog.created_at",
		"updated_at":         "catalog.updated_at",
		"root_id":            "catalog.root_id",
	}

	var groupByFields []string

	for _, field := range requestedFields {
		if groupField, ok := fieldMap[field]; ok {
			cols := strings.Split(groupField, ", ")
			groupByFields = append(groupByFields, cols...)
		}
	}

	return groupByFields
}

func buildCTEs(
	selectFlags map[string]bool,
	depth int64,
	subfields []string,
	filter map[string]any,
) string {
	var prefixQuery strings.Builder
	// prefixQuery.WriteString("WITH ")

	// Teams CTE
	if selectFlags["teams"] {
		if selectFlags["search"] {
			prefixQuery.WriteString(`
		teams AS (
			SELECT
				catalog_id,
				COALESCE(team_id, 0) AS team_id,
				COALESCE(team.name, '') AS team_name
			FROM cases.team_catalog
			JOIN call_center.cc_team team ON team.id = team_catalog.team_id
			WHERE catalog_id IN (SELECT target_catalog_id FROM search_catalog)
		),`)
		} else {
			prefixQuery.WriteString(`
		teams AS (
			SELECT
				catalog_id,
				COALESCE(team_id, 0) AS team_id,
				COALESCE(team.name, '') AS team_name
			FROM cases.team_catalog
			JOIN call_center.cc_team team ON team.id = team_catalog.team_id
		),`)
		}
	}

	// Skills CTE
	if selectFlags["skills"] {
		if selectFlags["search"] {
			prefixQuery.WriteString(`
		skills AS (
			SELECT
				catalog_id,
				COALESCE(skill_id, 0) AS skill_id,
				COALESCE(skill.name, '') AS skill_name
			FROM cases.skill_catalog
			JOIN call_center.cc_skill skill ON skill.id = skill_catalog.skill_id
			WHERE catalog_id IN (SELECT target_catalog_id FROM search_catalog)
		),`)
		} else {
			prefixQuery.WriteString(`
		skills AS (
			SELECT
				catalog_id,
				COALESCE(skill_id, 0) AS skill_id,
				COALESCE(skill.name, '') AS skill_name
			FROM cases.skill_catalog
			JOIN call_center.cc_skill skill ON skill.id = skill_catalog.skill_id
		),`)
		}
	}

	// Services Hierarchy CTE
	if selectFlags["services"] {
		anchorSQL := buildAnchorServiceSelect(subfields, selectFlags)
		subserviceSQL := buildSubserviceSelect(subfields, depth, selectFlags, filter)

		prefixQuery.WriteString(fmt.Sprintf(`service_hierarchy AS (
			WITH RECURSIVE recursive_hierarchy AS (
				%s
				UNION ALL
				%s
			)
			SELECT *
			FROM recursive_hierarchy
		),`, anchorSQL, subserviceSQL))
	}

	return strings.TrimSuffix(prefixQuery.String(), ",")
}

func buildAnchorServiceSelect(
	serviceFields []string,
	selectFlags map[string]bool,
) string {
	var sb strings.Builder

	sb.WriteString(`
SELECT catalog.id,
       catalog.name,
       catalog.description,
       catalog.root_id,
       catalog.id AS catalog_id,
       1 AS level
`)

	//----------------------------------------------------------------------
	// Conditionally insert code/state columns in SELECT for anchor
	//----------------------------------------------------------------------
	if util.ContainsField(serviceFields, "code") {
		sb.WriteString(`,
       COALESCE(catalog.code, '') AS code
`)
	}
	if util.ContainsField(serviceFields, "state") {
		sb.WriteString(`,
       catalog.state
`)
	}

	if util.ContainsField(serviceFields, "created_at") {
		sb.WriteString(`,
       catalog.created_at
`)
	}

	// If user wants "updated_at"
	if util.ContainsField(serviceFields, "updated_at") {
		sb.WriteString(`,
       catalog.updated_at
`)
	}

	// If user wants, say, "sla" in services, join SLA columns
	if util.ContainsField(serviceFields, "sla") {
		sb.WriteString(`,
       COALESCE(catalog.sla_id, 0) AS sla_id,
       COALESCE(service_sla.name, '') AS sla_name
`)
	}

	// If user wants "group"
	if util.ContainsField(serviceFields, "group") {
		sb.WriteString(`,
       COALESCE(catalog.group_id, 0) AS group_id,
       COALESCE(service_group.name, '') AS group_name,
       CASE
           WHEN catalog.group_id IN (SELECT id FROM contacts.dynamic_group) THEN 'dynamic'
           ELSE 'static'
       END AS group_type
`)
	}

	// If user wants "assignee"
	if util.ContainsField(serviceFields, "assignee") {
		sb.WriteString(`,
       COALESCE(catalog.assignee_id, 0) AS assignee_id,
       COALESCE(service_assignee.common_name, '') AS assignee_name
`)
	}

	// If user wants "created_by"
	if util.ContainsField(serviceFields, "created_by") {
		sb.WriteString(`,
       COALESCE(catalog.created_by, 0) AS created_by,
       COALESCE(service_created_by_user.name, '') AS created_by_name
`)
	}

	// If user wants "updated_by"
	if util.ContainsField(serviceFields, "updated_by") {
		sb.WriteString(`,
       COALESCE(catalog.updated_by, 0) AS updated_by,
       COALESCE(service_updated_by_user.name, '') AS updated_by_name
`)
	}

	sb.WriteString(`,
	false AS searched
`)

	// etc. for code, state, etc. if you want them in anchor
	sb.WriteString(`
FROM limited_catalogs AS catalog
`)

	// Now add conditional left joins in the anchor part
	if util.ContainsField(serviceFields, "sla") {
		sb.WriteString(`
LEFT JOIN cases.sla AS service_sla
       ON service_sla.id = catalog.sla_id
`)
	}
	if util.ContainsField(serviceFields, "group") {
		sb.WriteString(`
LEFT JOIN contacts.group AS service_group
       ON service_group.id = catalog.group_id
`)
	}
	if util.ContainsField(serviceFields, "assignee") {
		sb.WriteString(`
LEFT JOIN contacts.contact AS service_assignee
       ON service_assignee.id = catalog.assignee_id
`)
	}
	if util.ContainsField(serviceFields, "created_by") {
		sb.WriteString(`
LEFT JOIN directory.wbt_user AS service_created_by_user
       ON service_created_by_user.id = catalog.created_by
`)
	}
	if util.ContainsField(serviceFields, "updated_by") {
		sb.WriteString(`
LEFT JOIN directory.wbt_user AS service_updated_by_user
       ON service_updated_by_user.id = catalog.updated_by
`)
	}

	if selectFlags["search"] {
		sb.WriteString(`WHERE catalog.id IN (SELECT target_catalog_id FROM search_catalog)`)
	}

	return sb.String()
}

func buildSubserviceSelect(
	serviceFields []string,
	depth int64,
	selectFlags map[string]bool,
	filter map[string]any,
) string {
	var sb strings.Builder

	sb.WriteString(`
SELECT subservice.id,
       subservice.name,
       subservice.description,
       subservice.root_id,
	   parent.catalog_id,
       parent.level + 1 AS level
`)

	//----------------------------------------------------------------------
	// Conditionally insert code/state columns in SELECT for subservices
	//----------------------------------------------------------------------
	if util.ContainsField(serviceFields, "code") {
		sb.WriteString(`,
       COALESCE(subservice.code, '') AS code
`)
	}
	if util.ContainsField(serviceFields, "state") {
		sb.WriteString(`,
       subservice.state
`)
	}

	if util.ContainsField(serviceFields, "created_at") {
		sb.WriteString(`,
       subservice.created_at
`)
	}
	if util.ContainsField(serviceFields, "updated_at") {
		sb.WriteString(`,
       subservice.updated_at
`)
	}

	if util.ContainsField(serviceFields, "sla") {
		sb.WriteString(`,
       COALESCE(subservice.sla_id, 0) AS sla_id,
       COALESCE(service_sla2.name, '') AS sla_name
`)
	}
	// If user wants "group"
	if util.ContainsField(serviceFields, "group") {
		sb.WriteString(`,
       COALESCE(subservice.group_id, 0) AS group_id,
       COALESCE(service_group2.name, '') AS group_name,
       CASE
           WHEN subservice.group_id IN (SELECT id FROM contacts.dynamic_group) THEN 'dynamic'
           ELSE 'static'
       END AS group_type
`)
	}

	if util.ContainsField(serviceFields, "assignee") {
		sb.WriteString(`,
       COALESCE(subservice.assignee_id, 0) AS assignee_id,
       COALESCE(service_assignee2.common_name, '') AS assignee_name
`)
	}
	if util.ContainsField(serviceFields, "created_by") {
		sb.WriteString(`,
       COALESCE(subservice.created_by, 0) AS created_by,
       COALESCE(service_created_by_user2.name, '') AS created_by_name
`)
	}
	if util.ContainsField(serviceFields, "updated_by") {
		sb.WriteString(`,
       COALESCE(subservice.updated_by, 0) AS updated_by,
       COALESCE(service_updated_by_user2.name, '') AS updated_by_name
`)
	}

	if selectFlags["search"] {
		// etc. for code, state, etc. if you want them in anchor
		sb.WriteString(`
       ,CASE
    WHEN subservice.id IN (SELECT searched_id FROM search_catalog)
        THEN true
    ELSE false
END AS searched
	`)
	} else {
		// Add "searched" column for subservices
		sb.WriteString(`,
	false AS searched
`)
	}

	sb.WriteString(`
FROM cases.service_catalog AS subservice
`)

	// subservice-level left joins
	if util.ContainsField(serviceFields, "sla") {
		sb.WriteString(`
LEFT JOIN cases.sla AS service_sla2
       ON service_sla2.id = subservice.sla_id
`)
	}
	if util.ContainsField(serviceFields, "group") {
		sb.WriteString(`
LEFT JOIN contacts.group AS service_group2
       ON service_group2.id = subservice.group_id
`)
	}
	if util.ContainsField(serviceFields, "assignee") {
		sb.WriteString(`
LEFT JOIN contacts.contact AS service_assignee2
       ON service_assignee2.id = subservice.assignee_id
`)
	}
	if util.ContainsField(serviceFields, "created_by") {
		sb.WriteString(`
LEFT JOIN directory.wbt_user AS service_created_by_user2
       ON service_created_by_user2.id = subservice.created_by
`)
	}
	if util.ContainsField(serviceFields, "updated_by") {
		sb.WriteString(`
LEFT JOIN directory.wbt_user AS service_updated_by_user2
       ON service_updated_by_user2.id = subservice.updated_by
`)
	}

	// Ensure depth is valid and set defaults
	if depth == 0 {
		depth = 100
	} else if depth > 100 {
		depth = 100
	}

	sb.WriteString(`
JOIN recursive_hierarchy AS parent ON subservice.root_id = parent.id
WHERE parent.level < CASE WHEN `)

	// Automatically increase by +1 to include the requested level
	sb.WriteString(fmt.Sprintf(`%[1]d > 0 THEN %[1]d ELSE 3 END`, depth+1))

	// Add state filter for subservices
	if _, ok := filter["state"].(bool); ok {
		sb.WriteString(` AND subservice.state = :state`)
	}

	sb.WriteString(`
`)

	return sb.String()
}

// Update implements store.CatalogStore.
func (s *CatalogStore) Update(rpc options.UpdateOptions, lookup *cases.Catalog) (*cases.Catalog, error) {
	// Establish a connection to the database
	db, dbErr := s.storage.Database()
	if dbErr != nil {
		return nil, dberr.NewDBInternalError("postgres.catalog.update.database_connection_error", dbErr)
	}

	// Start a transaction using the TxManager
	tx, err := db.Begin(rpc)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.catalog.update.transaction_start_error", err)
	}
	txManager := transaction.NewTxManager(tx) // Create a new TxManager instance
	defer txManager.Rollback(rpc)             // Ensure rollback on error

	// Check if rpc.Fields contains team_ids or skill_ids
	updateTeams := false
	updateSkills := false

	// Check if the fields exist in rpc.Fields
	for _, field := range rpc.GetMask() {
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
		var teamIDs []int64
		var skillIDs []int64

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
			rpc,
			lookup.Id,
			teamIDs,  // Pass empty slice if no team IDs are provided
			skillIDs, // Pass empty slice if no skill IDs are provided
			rpc.GetAuthOpts().GetUserId(),
			rpc.RequestTime(),
			rpc.GetAuthOpts().GetDomainId(),
		)

		// Execute the teams and skills update query and check for affected rows
		var affectedRows int
		err = txManager.QueryRow(rpc, query, args...).Scan(&affectedRows)
		if err != nil {
			return nil, dberr.NewDBInternalError("postgres.catalog.update.teams_skills_update_error", err)
		}

		// Optional check if no rows were affected
		if affectedRows == 0 {
			return nil, dberr.NewDBNoRowsError("postgres.catalog.update.no_teams_skills_affected")
		}
	}

	// Build the update query for the Catalog
	query, args, err := s.buildUpdateCatalogQuery(rpc, lookup)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.catalog.update.query_build_error", err)
	}

	var (
		createdByLookup, updatedByLookup cases.Lookup
		createdAt, updatedAt             time.Time
		teamLookups, skillLookups        []byte
	)

	err = txManager.QueryRow(rpc, query, args...).Scan(
		&lookup.Id, &lookup.Name, &createdAt,
		&lookup.Sla.Id, &lookup.Sla.Name,
		&lookup.Status.Id, &lookup.Status.Name,
		&lookup.CloseReasonGroup.Id, &lookup.CloseReasonGroup.Name,
		&createdByLookup.Id, &createdByLookup.Name,
		&updatedByLookup.Id, &updatedByLookup.Name, &updatedAt,
		&lookup.State,
		&teamLookups, &skillLookups,
	)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.catalog.update.execution_error", err)
	}

	// Commit the transaction
	if err := txManager.Commit(rpc); err != nil {
		return nil, dberr.NewDBInternalError("postgres.catalog.update.transaction_commit_error", err)
	}

	// Unmarshal the JSON arrays for teams and skills
	if err := json.Unmarshal(teamLookups, &lookup.Teams); err != nil {
		return nil, dberr.NewDBInternalError("postgres.catalog.update.unmarshal_teams_error", err)
	}
	if err := json.Unmarshal(skillLookups, &lookup.Skills); err != nil {
		return nil, dberr.NewDBInternalError("postgres.catalog.update.unmarshal_skills_error", err)
	}

	// Prepare the updated Catalog to return
	lookup.CreatedAt = util.Timestamp(createdAt)
	lookup.UpdatedAt = util.Timestamp(updatedAt)
	lookup.CreatedBy = &createdByLookup
	lookup.UpdatedBy = &updatedByLookup

	// Return the updated Catalog
	return lookup, nil
}

func (s *CatalogStore) buildUpdateTeamsAndSkillsQuery(
	rpc options.UpdateOptions,
	catalogID int64,
	teamIDs,
	skillIDs []int64,
	updatedBy int64,
	updatedAt time.Time,
	domainID int64,
) (string, []interface{}) {
	args := []interface{}{
		catalogID, // $1: catalog_id
		updatedBy, // $2: updated_by (will also be used for created_by)
		domainID,  // $3: dc (domain context)
		updatedAt, // $4: timestamp for updated_at
	}

	// Initialize base query
	query := `WITH`

	// Flag to manage if we've added any CTEs
	cteAdded := false
	placeholderIndex := 5 // Start placeholder index after the initial args

	// Check if "teams" is in rpc.Fields, even if teamIDs is empty
	if util.FieldExists("teams", rpc.GetMask()) {
		query += `
 updated_teams AS (
    INSERT INTO cases.team_catalog (catalog_id, team_id, created_by, updated_by, updated_at, dc)
        SELECT $1, unnest(NULLIF($` + fmt.Sprintf("%d", placeholderIndex) + `::bigint[], '{}')), $2, $2, $4, $3
        ON CONFLICT (catalog_id, team_id)
            DO UPDATE SET updated_at = EXCLUDED.updated_at, updated_by = EXCLUDED.updated_by
        RETURNING catalog_id
    ),
 deleted_teams AS (
     DELETE FROM cases.team_catalog
     WHERE catalog_id = $1
       AND (
         array_length($` + fmt.Sprintf("%d", placeholderIndex) + `, 1) IS NULL
         OR team_id != ALL ($` + fmt.Sprintf("%d", placeholderIndex) + `)
       )
     RETURNING catalog_id
    )`
		args = append(args, pq.Array(teamIDs)) // Append team IDs to args (even if empty)
		placeholderIndex++                     // Increment placeholder index
		cteAdded = true
	}

	// Check if "skills" is in rpc.Fields, even if skillIDs are empty
	if util.FieldExists("skills", rpc.GetMask()) {
		if cteAdded {
			query += `,` // Only add a comma if there is already a CTE defined (for teams)
		}
		query += `
 updated_skills AS (
    INSERT INTO cases.skill_catalog (catalog_id, skill_id, created_by, updated_by, updated_at, dc)
        SELECT $1, unnest(NULLIF($` + fmt.Sprintf("%d", placeholderIndex) + `::bigint[], '{}')), $2, $2, $4, $3
        ON CONFLICT (catalog_id, skill_id)
            DO UPDATE SET updated_at = EXCLUDED.updated_at, updated_by = EXCLUDED.updated_by
        RETURNING catalog_id
    ),
 deleted_skills AS (
     DELETE FROM cases.skill_catalog
     WHERE catalog_id = $1
       AND (
         array_length($` + fmt.Sprintf("%d", placeholderIndex) + `, 1) IS NULL
         OR skill_id != ALL ($` + fmt.Sprintf("%d", placeholderIndex) + `)
       )
     RETURNING catalog_id
    )`
		args = append(args, pq.Array(skillIDs)) // Append skill IDs to args (even if empty)
		placeholderIndex++                      // Increment placeholder index
		cteAdded = true
	}

	// Construct the final SELECT query after the CTE block
	query += `
SELECT COUNT(*)
FROM (
    ` + func() string {
		var result string
		if util.FieldExists("teams", rpc.GetMask()) {
			result += `SELECT catalog_id FROM updated_teams UNION ALL SELECT catalog_id FROM deleted_teams`
		}
		if util.FieldExists("skills", rpc.GetMask()) {
			if util.FieldExists("teams", rpc.GetMask()) {
				result += ` UNION ALL `
			}
			result += `SELECT catalog_id FROM updated_skills UNION ALL SELECT catalog_id FROM deleted_skills`
		}
		return result
	}() + `
) AS total_affected;`

	// Return the constructed query and arguments
	return util2.CompactSQL(query), args
}

func (s *CatalogStore) buildUpdateCatalogQuery(
	rpc options.UpdateOptions,
	lookup *cases.Catalog,
) (string, []interface{}, error) {
	// Start the WITH clause to check if root_id is NULL
	// Checking whether root is NULL or not
	//
	// If ROOT is NOT NULL ---- user try to update service
	// Status / Prefix / Close reason could not be set for service
	//TODO REMOVE %d
	checkRoot := fmt.Sprintf(`
WITH root_check AS (
    SELECT root_id
    FROM cases.service_catalog
    WHERE id = %d
)
`, lookup.Id)

	// Start the update query with Squirrel Update Builder
	updateQueryBuilder := sq.Update("cases.service_catalog").
		PlaceholderFormat(sq.Dollar).
		Set("updated_at", rpc.RequestTime()).
		Set("updated_by", rpc.GetAuthOpts().GetUserId()).
		Where(sq.Eq{"id": lookup.Id, "dc": rpc.GetAuthOpts().GetDomainId()})

	// Dynamically set fields based on user update preferences
	for _, field := range rpc.GetMask() {
		switch field {
		case "name":
			updateQueryBuilder = updateQueryBuilder.Set("name", lookup.Name)
		case "description":
			updateQueryBuilder = updateQueryBuilder.Set("description",
				sq.Expr("NULLIF(?, '')",
					lookup.Description,
				))
		case "prefix":
			updateQueryBuilder = updateQueryBuilder.Set("prefix",
				sq.Expr("(CASE WHEN (SELECT root_id FROM root_check) IS NULL THEN ? ELSE prefix END)",
					lookup.Prefix,
				))
		case "code":
			updateQueryBuilder = updateQueryBuilder.Set("code", sq.Expr("NULLIF(?, '')", lookup.Code))
		case "state":
			updateQueryBuilder = updateQueryBuilder.Set("state", lookup.State)
		case "sla":
			updateQueryBuilder = updateQueryBuilder.Set("sla_id", lookup.Sla.Id)
		case "status":
			updateQueryBuilder = updateQueryBuilder.Set("status_id",
				sq.Expr("(CASE WHEN (SELECT root_id FROM root_check) IS NULL THEN ? ELSE status_id END)",
					lookup.Status.Id,
				))
		case "close_reason_group":
			updateQueryBuilder = updateQueryBuilder.Set("close_reason_group_id",
				sq.Expr("(CASE WHEN (SELECT root_id FROM root_check) IS NULL THEN NULLIF(?, 0) ELSE close_reason_group_id END)",
					lookup.CloseReasonGroup.Id,
				))
		}
	}

	// Convert the update query to SQL
	updateSQL, args, err := updateQueryBuilder.ToSql()
	if err != nil {
		return "", nil, err
	}

	// Combine the WITH clause and update query
	query := fmt.Sprintf(`
%s, updated_catalog AS (
    %s
    RETURNING *
)
SELECT catalog.id,
       catalog.name,
       catalog.created_at,
       catalog.sla_id,
       sla.name,
       COALESCE(catalog.status_id, 0) AS status_id,
       COALESCE(status.name, '') AS status_name,
       COALESCE(catalog.close_reason_group_id, 0) AS close_reason_group_id,
       COALESCE(close_reason_group.name, '')                    AS close_reason_group_name,
       catalog.created_by,
       COALESCE(created_by_user.name, '')                 AS created_by_name,
       catalog.updated_by,
       COALESCE(updated_by_user.name, '')                               AS updated_by_name,
       catalog.updated_at,
	   catalog.state,
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
         LEFT JOIN cases.close_reason_group ON close_reason_group.id = catalog.close_reason_group_id
         LEFT JOIN directory.wbt_user AS created_by_user ON created_by_user.id = catalog.created_by
         LEFT JOIN directory.wbt_user AS updated_by_user ON updated_by_user.id = catalog.updated_by
GROUP BY catalog.id, catalog.name, catalog.created_at, catalog.sla_id, sla.name, catalog.status_id,
         status.name, catalog.close_reason_group_id, close_reason_group.name, catalog.created_by, created_by_user.name,
         catalog.updated_by, updated_by_user.name, catalog.updated_at, catalog.state;
	`, checkRoot, updateSQL)

	// Return the final combined query and arguments
	return util2.CompactSQL(query), args, nil
}

func NewCatalogStore(store *Store) (store.CatalogStore, error) {
	if store == nil {
		return nil, dberr.NewDBError("postgres.new_catalog.check.bad_arguments",
			"error creating Catalog interface to the service table, main store is nil")
	}
	return &CatalogStore{storage: store}, nil
}
