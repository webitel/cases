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

type CatalogStore struct {
	storage store.Store
}

// Create implements store.CatalogStore.
func (s *CatalogStore) Create(rpc *model.CreateOptions, add *cases.Catalog) (*cases.Catalog, error) {
	// Establish a connection to the database
	db, dbErr := s.storage.Database()
	if dbErr != nil {
		return nil, model.NewInternalError("postgres.catalog.create.database_connection_error", dbErr.Error())
	}

	// Build the combined query for inserting Catalog, teams, and skills
	query, args := s.buildCreateCatalogQuery(rpc, add)

	var (
		createdByLookup, updatedByLookup cases.Lookup
		createdAt, updatedAt             time.Time
		teamLookups, skillLookups        []byte
	)

	err := db.QueryRow(rpc.Context, query, args...).Scan(
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

func (s *CatalogStore) buildCreateCatalogQuery(rpc *model.CreateOptions, add *cases.Catalog) (string, []interface{}) {
	// Define arguments for the query
	args := []interface{}{
		add.Name,                  // $1: name (cannot be null)
		add.Description,           // $2: description (could be null)
		add.Prefix,                // $3: prefix (could be null)
		add.Code,                  // $4: code (could be null)
		rpc.Time,                  // $5: created_at, updated_at
		rpc.Session.GetUserId(),   // $6: created_by, updated_by
		add.Sla.Id,                // $7: sla_id (could be null)
		add.Status.Id,             // $8: status_id (could be null)
		add.CloseReason.Id,        // $9: close_reason_id (could be null)
		add.State,                 // $10: state (cannot be null)
		rpc.Session.GetDomainId(), // $11: domain ID (dc)
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

// Delete implements store.CatalogStore.
func (s *CatalogStore) Delete(rpc *model.DeleteOptions) error {
	// Establish a connection to the database
	db, dbErr := s.storage.Database()
	if dbErr != nil {
		return model.NewInternalError("postgres.catalog.delete.db_connection_error", dbErr.Error())
	}

	// Ensure that there are IDs to delete
	if len(rpc.IDs) == 0 {
		return model.NewBadRequestError("postgres.catalog.delete.no_ids_provided", "No IDs provided for deletion")
	}

	// Build the delete query
	query, args := s.buildDeleteCatalogQuery(rpc)

	// Execute the delete query
	res, err := db.Exec(rpc.Context, query, args...)
	if err != nil {
		return model.NewInternalError("postgres.catalog.delete.execution_error", err.Error())
	}

	// Check how many rows were affected
	if res.RowsAffected() == 0 {
		return model.NewNotFoundError("postgres.catalog.delete.no_rows_deleted", "No Catalog entries were deleted")
	}

	return nil
}

// Helper method to build the delete query for Catalog
func (s *CatalogStore) buildDeleteCatalogQuery(rpc *model.DeleteOptions) (string, []interface{}) {
	// Build the SQL query using the provided IDs in rpc.IDs
	query := `
		DELETE FROM cases.service_catalog
		WHERE id = ANY($1) AND dc = $2
	`

	// Use the array of IDs and domain ID (dc) for the deletion
	args := []interface{}{
		pq.Array(rpc.IDs),         // $1: array of catalog IDs to delete
		rpc.Session.GetDomainId(), // $2: domain ID to ensure proper scoping
	}

	return store.CompactSQL(query), args
}

// List implements store.CatalogStore.
func (s *CatalogStore) List(
	rpc *model.SearchOptions,
	depth int64,
) (*cases.CatalogList, error) {
	// Establish a connection to the database
	db, dbErr := s.storage.Database()
	if dbErr != nil {
		return nil, model.NewInternalError("postgres.catalog.list.database_connection_error", dbErr.Error())
	}

	// Build SQL query
	query, args, err := s.buildSearchCatalogQuery(rpc, depth)
	if err != nil {
		return nil, model.NewInternalError("postgres.catalog.list.query_build_error", err.Error())
	}

	// Execute the query
	rows, err := db.Query(rpc.Context, query, args...)
	if err != nil {
		return nil, model.NewInternalError("postgres.catalog.list.query_execution_error", err.Error())
	}
	defer rows.Close()

	// Parse the result
	var catalogs []*cases.Catalog
	lCount := 0
	next := false
	fetchAll := rpc.GetSize() == -1

	for rows.Next() {
		// If not fetching all records, check the size limit
		if !fetchAll && lCount >= rpc.GetSize() {
			next = true
			break
		}

		// Initialize catalog and related fields
		catalog := &cases.Catalog{
			Sla:         &cases.Lookup{},
			Status:      &cases.Lookup{},
			CloseReason: &cases.Lookup{},
		}

		var (
			createdBy, updatedBy                      cases.Lookup
			createdAt, updatedAt                      time.Time
			teamLookups, skillLookups, serviceLookups []byte
			rootID                                    int64

			// Services slice to hold the nested services
			services []map[string]interface{}
		)

		// Build scan arguments to include services
		scanArgs, err := s.buildCatalogScanArgs(
			catalog,
			&createdBy, &updatedBy,
			&createdAt, &updatedAt,
			&teamLookups, &skillLookups, &serviceLookups,
			&rootID,
			// ----- Fields to scan -----
			rpc.Fields,
		)
		if err != nil {
			return nil, model.NewInternalError("postgres.catalog.list.scan_args_error", err.Error())
		}

		// Scan the result into the appropriate fields
		err = rows.Scan(scanArgs...)
		if err != nil {
			return nil, model.NewInternalError("postgres.catalog.list.scan_error", err.Error())
		}

		// If rootID is not 0, it's a subservice, so we skip it
		if rootID != 0 {
			continue
		}

		// Unmarshal the JSON arrays into the Lookup slices
		if err := json.Unmarshal(teamLookups, &catalog.Teams); err != nil {
			return nil, model.NewInternalError("postgres.catalog.list.unmarshal_teams_error", err.Error())
		}
		if err := json.Unmarshal(skillLookups, &catalog.Skills); err != nil {
			return nil, model.NewInternalError("postgres.catalog.list.unmarshal_skills_error", err.Error())
		}

		// Handle services unmarshal
		if len(serviceLookups) > 0 {
			if err := json.Unmarshal(serviceLookups, &services); err != nil {
				return nil, model.NewInternalError("postgres.catalog.list.unmarshal_services_error", err.Error())
			}

			// Nest services by root_id
			nestedServices, err := s.nestServicesByRootID(catalog.Id, services)
			if err != nil {
				return nil, model.NewInternalError("postgres.catalog.list.nesting_services_error", err.Error())
			}

			// Add the nested services to the catalog
			catalog.Service = nestedServices
		}

		// Set timestamps and created/updated by fields for the catalog
		catalog.CreatedAt = util.Timestamp(createdAt)
		catalog.UpdatedAt = util.Timestamp(updatedAt)
		catalog.CreatedBy = &createdBy
		catalog.UpdatedBy = &updatedBy

		catalogs = append(catalogs, catalog)
		lCount++
	}

	return &cases.CatalogList{
		Page:  int32(rpc.Page),
		Next:  next,
		Items: catalogs,
	}, nil
}

// buildCatalogScanArgs prepares scan arguments based on rpc.Fields.
// If rpc.Fields contains only "-", all fields will be scanned. Otherwise, fields are selectively scanned.
func (s *CatalogStore) buildCatalogScanArgs(
	catalog *cases.Catalog, // The catalog object to populate
	createdBy, updatedBy *cases.Lookup, // Lookup objects for created_by and updated_by
	createdAt, updatedAt *time.Time, // Temporary variables for created_at and updated_at
	teamLookups, skillLookups, serviceLookups *[]byte, // Byte arrays for teams, skills, and services (as JSON or binary)
	rootId *int64, // Root ID for hierarchy placement
	rpcFields []string, // List of fields to scan based on the request
) ([]interface{}, error) {
	// ------ If rpc.Fields is "-", scan all fields ------
	if rpcFields[0] == "-" {
		return []interface{}{
			// ------ Catalog Basic Information ------
			&catalog.Id,     // Catalog ID
			&catalog.Name,   // Catalog Name
			&catalog.Prefix, // Catalog Prefix

			// ------ SLA Fields ------
			&catalog.Sla.Id,   // SLA ID
			&catalog.Sla.Name, // SLA Name

			// ------ Status Fields ------
			&catalog.Status.Id,   // Status ID
			&catalog.Status.Name, // Status Name

			// ------ Catalog Metadata ------
			&catalog.Code,        // Catalog Code
			&catalog.Description, // Catalog Description

			// ------ Close Reason Fields ------
			&catalog.CloseReason.Id,   // Close Reason ID
			&catalog.CloseReason.Name, // Close Reason Name

			// ------ Catalog State ------
			&catalog.State, // Catalog State (active/inactive)

			// ------ Created By and Updated By Fields ------
			&createdBy.Id,   // Created By User ID
			&createdBy.Name, // Created By User Name
			&updatedBy.Id,   // Updated By User ID
			&updatedBy.Name, // Updated By User Name

			// ------ Timestamps ------
			createdAt, // Created At Timestamp
			updatedAt, // Updated At Timestamp

			// ------ Teams, Skills, and Services Lookups ------
			teamLookups,    // Team Lookups (JSON/binary)
			skillLookups,   // Skill Lookups (JSON/binary)
			serviceLookups, // Service Lookups (JSON/binary)

			// ------ Root ID and Hierarchy Info ------
			rootId, // Root ID for hierarchy
		}, nil
	}

	// ------ If rpc.Fields contains specific fields, scan accordingly ------
	var scanArgs []interface{}

	for _, field := range rpcFields {
		switch field {

		// ------ Catalog Basic Information ------
		case "id":
			scanArgs = append(scanArgs, &catalog.Id) // Catalog ID
		case "name":
			scanArgs = append(scanArgs, &catalog.Name) // Catalog Name
		case "prefix":
			scanArgs = append(scanArgs, &catalog.Prefix) // Catalog Prefix

		// ------ SLA Fields ------
		case "sla":
			scanArgs = append(scanArgs, &catalog.Sla.Id, &catalog.Sla.Name) // SLA ID and Name

		// ------ Status Fields ------
		case "status":
			scanArgs = append(scanArgs, &catalog.Status.Id, &catalog.Status.Name) // Status ID and Name

		// ------ Catalog Metadata ------
		case "code":
			scanArgs = append(scanArgs, &catalog.Code) // Catalog Code
		case "description":
			scanArgs = append(scanArgs, &catalog.Description) // Catalog Description

		// ------ Close Reason Fields ------
		case "close_reason":
			scanArgs = append(scanArgs, &catalog.CloseReason.Id, &catalog.CloseReason.Name) // Close Reason ID and Name

		// ------ Catalog State ------
		case "state":
			scanArgs = append(scanArgs, &catalog.State) // Catalog State (active/inactive)

		// ------ Created By and Updated By Fields ------
		case "created_by":
			scanArgs = append(scanArgs, &createdBy.Id, &createdBy.Name) // Created By User ID and Name
		case "updated_by":
			scanArgs = append(scanArgs, &updatedBy.Id, &updatedBy.Name) // Updated By User ID and Name

		// ------ Timestamps ------
		case "created_at":
			scanArgs = append(scanArgs, createdAt) // Created At Timestamp
		case "updated_at":
			scanArgs = append(scanArgs, updatedAt) // Updated At Timestamp

		// ------ Teams, Skills, and Services Lookups ------
		case "teams":
			scanArgs = append(scanArgs, teamLookups) // Team Lookups (JSON/binary)
		case "skills":
			scanArgs = append(scanArgs, skillLookups) // Skill Lookups (JSON/binary)
		case "services":
			scanArgs = append(scanArgs, serviceLookups) // Service Lookups (JSON/binary)

		// ------ Root ID and Hierarchy Info ------
		case "root_id":
			scanArgs = append(scanArgs, rootId) // Root ID for hierarchy
		}
	}

	return scanArgs, nil
}

func (s *CatalogStore) mapServiceData(serviceData map[string]interface{}) (*cases.Service, error) {
	// Extract necessary fields from the service data map
	serviceSlaID := int64(serviceData["sla_id"].(float64))
	serviceSlaName := serviceData["sla_name"].(string)
	serviceGroupID := int64(serviceData["group_id"].(float64))
	serviceGroupName := serviceData["group_name"].(string)
	serviceAssigneeID := int64(serviceData["assignee_id"].(float64))
	serviceAssigneeName := serviceData["assignee_name"].(string)

	// Extract created_at, updated_at as strings and convert them to int64 timestamps
	createdAtStr := serviceData["created_at"].(string)
	updatedAtStr := serviceData["updated_at"].(string)

	// Convert time strings to timestamps (using time.RFC3339Nano format)
	createdAt, err := util.TimeStringToTimestamp(createdAtStr, time.RFC3339Nano)
	if err != nil {
		return nil, fmt.Errorf("Error parsing created_at: %v", err)
	}
	updatedAt, err := util.TimeStringToTimestamp(updatedAtStr, time.RFC3339Nano)
	if err != nil {
		return nil, fmt.Errorf("Error parsing updated_at: %v", err)
	}

	createdByID := int64(serviceData["created_by_id"].(float64))
	createdByName := serviceData["created_by"].(string)
	updatedByID := int64(serviceData["updated_by_id"].(float64))
	updatedByName := serviceData["updated_by"].(string)

	// Construct the service object
	service := &cases.Service{
		Id:          int64(serviceData["id"].(float64)),
		Name:        serviceData["name"].(string),
		Description: serviceData["description"].(string),
		Code:        serviceData["code"].(string),
		State:       serviceData["state"].(bool),
		RootId:      int64(serviceData["root_id"].(float64)),
		CatalogId:   int64(serviceData["catalog_id"].(float64)),
		Sla: &cases.Lookup{
			Id:   serviceSlaID,
			Name: serviceSlaName,
		},
		Group: &cases.Lookup{
			Id:   serviceGroupID,
			Name: serviceGroupName,
		},
		Assignee: &cases.Lookup{
			Id:   serviceAssigneeID,
			Name: serviceAssigneeName,
		},
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		CreatedBy: &cases.Lookup{
			Id:   createdByID,
			Name: createdByName,
		},
		UpdatedBy: &cases.Lookup{
			Id:   updatedByID,
			Name: updatedByName,
		},
		Service: []*cases.Service{}, // Initialize empty slice for children services
	}

	return service, nil
}

func (s *CatalogStore) nestServicesByRootID(catalogID int64, services []map[string]interface{}) ([]*cases.Service, error) {
	// Step 1: Group services by their root_id
	serviceMap := make(map[int64][]map[string]interface{})
	var topServices []map[string]interface{}

	for _, serviceData := range services {
		rootID := int64(serviceData["root_id"].(float64))
		if rootID == catalogID {
			// Top-level services (directly under the catalog)
			topServices = append(topServices, serviceData)
		} else {
			// Group by root_id (services under other services)
			serviceMap[rootID] = append(serviceMap[rootID], serviceData)
		}
	}

	// Step 2: Recursively build the hierarchy
	return s.buildServiceHierarchy(topServices, serviceMap)
}

func (s *CatalogStore) buildServiceHierarchy(
	serviceDataList []map[string]interface{},
	serviceMap map[int64][]map[string]interface{},
) ([]*cases.Service, error) {
	var services []*cases.Service

	for _, serviceData := range serviceDataList {
		// Map service data into a Service object
		service, err := s.mapServiceData(serviceData)
		if err != nil {
			return nil, err
		}

		// Recursively build children (sub-services)
		if childrenData, exists := serviceMap[service.Id]; exists {
			children, err := s.buildServiceHierarchy(childrenData, serviceMap)
			if err != nil {
				return nil, err
			}
			// Assign children to the service
			service.Service = children
		}

		// Append the service to the list of services
		services = append(services, service)
	}

	return services, nil
}

func (s *CatalogStore) buildSearchCatalogQuery(
	rpc *model.SearchOptions,
	depth int64,
) (string, []interface{}, error) {
	// Default fields: includes all teams, skills, services, etc.
	defaultCatalogFields := []string{
		"catalog.id",
		"catalog.name",                                            // Mandatory
		"COALESCE(catalog.prefix, '') AS prefix",                  // Use COALESCE for prefix to handle null values
		"COALESCE(catalog.sla_id, 0) AS sla_id",                   // Use COALESCE for SLA ID to handle null values
		"COALESCE(sla.name, '') AS sla_name",                      // Use COALESCE for SLA name to handle null values
		"COALESCE(catalog.status_id, 0) AS status_id",             // Use COALESCE for status ID to handle null values
		"COALESCE(status.name, '') AS status_name",                // Use COALESCE for status name to handle null values
		"COALESCE(catalog.code, '') AS code",                      // Optional
		"COALESCE(catalog.description, '') AS description",        // Optional
		"COALESCE(catalog.close_reason_id, 0) AS close_reason_id", // Optional
		"COALESCE(close_reason.name, '') AS close_reason_name",    // Optional
		"catalog.state AS state",
		"COALESCE(catalog.created_by, 0) AS created_by",         // Handle null with default 0 for ID
		"COALESCE(created_by_user.name, '') AS created_by_name", // Handle null with default empty string for name
		"COALESCE(catalog.updated_by, 0) AS updated_by",         // Handle null with default 0 for ID
		"COALESCE(updated_by_user.name, '') AS updated_by_name", // Handle null with default empty string for name
		"catalog.created_at AS created_at",
		"catalog.updated_at AS updated_at",
		"COALESCE(teams_agg.teams, '[]') AS teams",          // Aggregated teams from the CTE
		"COALESCE(skills_agg.skills, '[]') AS skills",       // Aggregated skills from the CTE
		"COALESCE(services_agg.services, '[]') AS services", // Aggregated services from the recursive CTE
		"COALESCE(catalog.root_id, 0) AS root_id",           // Aggregated services from the recursive CTE
	}

	// Initialize flags for recursion
	selectFlags := map[string]bool{
		"services": false,
		"teams":    false,
		"skills":   false,
	}

	// Selected fields handling
	var selectedFields []string

	// If fields are set to "-", use defaultCatalogFields and enable recursion for all entities
	if rpc.Fields[0] == "-" {
		selectedFields = defaultCatalogFields
		selectFlags["services"] = true
		selectFlags["teams"] = true
		selectFlags["skills"] = true
	} else {
		// Handle specific fields from rpc.Fields
		for _, field := range rpc.Fields {
			switch field {
			case "id":
				selectedFields = append(selectedFields, "catalog.id")
			case "name":
				selectedFields = append(selectedFields, "catalog.name")
			case "prefix":
				selectedFields = append(selectedFields, "COALESCE(catalog.prefix, '') AS prefix")
				selectFlags["services"] = true // Enable services recursion if prefix is selected
			case "sla":
				selectedFields = append(selectedFields, "COALESCE(catalog.sla_id, 0) AS sla_id", "COALESCE(sla.name, '') AS sla_name")
			case "status":
				selectedFields = append(selectedFields, "COALESCE(catalog.status_id, 0) AS status_id", "COALESCE(status.name, '') AS status_name")
			case "code":
				selectedFields = append(selectedFields, "COALESCE(catalog.code, '') AS code")
			case "description":
				selectedFields = append(selectedFields, "COALESCE(catalog.description, '') AS description")
			case "close_reason":
				selectedFields = append(selectedFields, "COALESCE(catalog.close_reason_id, 0) AS close_reason_id", "COALESCE(close_reason.name, '') AS close_reason_name")
			case "state":
				selectedFields = append(selectedFields, "catalog.state AS state")
			case "created_by":
				selectedFields = append(selectedFields, "COALESCE(catalog.created_by, 0) AS created_by", "COALESCE(created_by_user.name, '') AS created_by_name")
			case "updated_by":
				selectedFields = append(selectedFields, "COALESCE(catalog.updated_by, 0) AS updated_by", "COALESCE(updated_by_user.name, '') AS updated_by_name")
			case "created_at":
				selectedFields = append(selectedFields, "catalog.created_at AS created_at")
			case "updated_at":
				selectedFields = append(selectedFields, "catalog.updated_at AS updated_at")
			case "teams":
				selectedFields = append(selectedFields, "COALESCE(teams_agg.teams, '[]') AS teams")
				selectFlags["teams"] = true
			case "skills":
				selectedFields = append(selectedFields, "COALESCE(skills_agg.skills, '[]') AS skills")
				selectFlags["skills"] = true
			case "services":
				selectedFields = append(selectedFields, "COALESCE(services_agg.services, '[]') AS services")
				selectFlags["services"] = true // Enable services recursion
			case "root_id":
				selectedFields = append(selectedFields, "COALESCE(catalog.root_id, 0) AS root_id")
			}
		}
	}

	// Build the base query with the selected fields
	queryBuilder := sq.Select(selectedFields...).
		From("cases.service_catalog AS catalog").
		LeftJoin("cases.sla ON sla.id = catalog.sla_id").
		LeftJoin("cases.status ON status.id = catalog.status_id").
		LeftJoin("cases.close_reason ON close_reason.id = catalog.close_reason_id").
		LeftJoin("directory.wbt_user AS created_by_user ON created_by_user.id = catalog.created_by").
		LeftJoin("directory.wbt_user AS updated_by_user ON updated_by_user.id = catalog.updated_by").
		PlaceholderFormat(sq.Dollar)

		// Conditionally add LeftJoin for teams if the field is selected
	if selectFlags["teams"] {
		queryBuilder = queryBuilder.LeftJoin("teams_agg ON teams_agg.catalog_id = catalog.id")
	}

	// Conditionally add LeftJoin for skills if the field is selected
	if selectFlags["skills"] {
		queryBuilder = queryBuilder.LeftJoin("skills_agg ON skills_agg.catalog_id = catalog.id")
	}

	// Conditionally add LeftJoin for services if the field is selected
	if selectFlags["services"] {
		queryBuilder = queryBuilder.LeftJoin("services_agg ON services_agg.catalog_id = catalog.id")
	}

	// Apply filtering by name
	if name, ok := rpc.Filter["name"].(string); ok && len(name) > 0 {
		queryBuilder = queryBuilder.Where(sq.ILike{"catalog.name": "%" + name + "%"})
	}

	// Apply filtering by state
	if state, ok := rpc.Filter["state"]; ok {
		queryBuilder = queryBuilder.Where(sq.Eq{"catalog.state": state})
	}

	// Apply filtering by IDs if provided
	if len(rpc.IDs) > 0 {
		queryBuilder = queryBuilder.Where(sq.Eq{"catalog.id": rpc.IDs})
	}

	// Apply sorting
	for _, sort := range rpc.Sort {
		queryBuilder = queryBuilder.OrderBy(sort)
	}

	// Pagination: Apply limit and offset
	size := rpc.GetSize()
	if size != -1 {
		queryBuilder = queryBuilder.Limit(uint64(size + 1)) // Request one more record to check if there's a next page
	}
	if rpc.Page > 1 {
		queryBuilder = queryBuilder.Offset(uint64((rpc.Page - 1) * size))
	}

	// Apply recursive CTEs based on flags and depth
	if selectFlags["teams"] || selectFlags["skills"] || selectFlags["services"] {
		var prefixQuery string
		if selectFlags["teams"] {
			prefixQuery += `
WITH inserted_teams AS (SELECT catalog_id, team_id
                        FROM cases.team_catalog),
     teams_agg AS (SELECT inserted_teams.catalog_id,
                          json_agg(json_build_object('id', team.id, 'name', team.name)) AS teams
                   FROM inserted_teams
                            LEFT JOIN call_center.cc_team team ON team.id = inserted_teams.team_id
                   GROUP BY inserted_teams.catalog_id),`
		}
		if selectFlags["skills"] {
			prefixQuery += `
 inserted_skills AS (SELECT catalog_id, skill_id
                         FROM cases.skill_catalog),
     skills_agg AS (SELECT inserted_skills.catalog_id,
                           json_agg(json_build_object('id', skill.id, 'name', skill.name)) AS skills
                    FROM inserted_skills
                             LEFT JOIN call_center.cc_skill skill ON skill.id = inserted_skills.skill_id
                    GROUP BY inserted_skills.catalog_id),`
		}
		if selectFlags["services"] {
			prefixQuery += fmt.Sprintf(`services_agg AS (
   WITH RECURSIVE service_hierarchy AS (SELECT service.id,
                                            service.name,
                                            service.description,
                                            service.code,
                                            service.state,
                                            service.sla_id,
                                            service.group_id,
                                            service.assignee_id,
                                            service.root_id,
                                            service.created_at,
                                            service.updated_at,
                                            service.created_by,
                                            service.updated_by,
                                            catalog.id AS catalog_id,
                                            1          AS level
                                     FROM cases.service_catalog service
                                              JOIN cases.service_catalog catalog ON service.root_id = catalog.id
                                     WHERE catalog.root_id IS NULL

                                     UNION ALL

                                     -- Recursively fetch subservices
                                     SELECT subservice.id,
                                            subservice.name,
                                            subservice.description,
                                            subservice.code,
                                            subservice.state,
                                            subservice.sla_id,
                                            subservice.group_id,
                                            subservice.assignee_id,
                                            subservice.root_id,
                                            subservice.created_at,
                                            subservice.updated_at,
                                            subservice.created_by,
                                            subservice.updated_by,
                                            parent.catalog_id,
                                            parent.level + 1 AS level
                                     FROM cases.service_catalog subservice
                                              JOIN service_hierarchy parent ON subservice.root_id = parent.id
                                     WHERE parent.level < CASE WHEN %[1]d > 0 THEN %[1]d ELSE 100 END)
SELECT service_hierarchy.catalog_id,
       json_agg(json_build_object(
               'id', service_hierarchy.id,
               'name', service_hierarchy.name,
               'description', service_hierarchy.description,
               'code', service_hierarchy.code,
               'state', service_hierarchy.state,
               'sla_id', COALESCE(service_hierarchy.sla_id, 0),
               'sla_name', COALESCE(sla.name, ''),
               'group_id', COALESCE(service_hierarchy.group_id, 0),
               'group_name', COALESCE(grp.name, ''),
               'assignee_id', COALESCE(service_hierarchy.assignee_id, 0),
               'assignee_name', COALESCE(assignee.given_name, ''),
               'has_subservices',
               EXISTS (SELECT 1 FROM cases.service_catalog sc WHERE sc.root_id = service_hierarchy.id),
               'root_id', COALESCE(service_hierarchy.root_id, 0),
               'created_at', service_hierarchy.created_at,
               'updated_at', service_hierarchy.updated_at,
               'created_by', COALESCE(created_by_user.name, ''),
               'created_by_id', COALESCE(service_hierarchy.created_by, 0),
               'updated_by', COALESCE(updated_by_user.name, ''),
               'updated_by_id', COALESCE(service_hierarchy.updated_by, 0),
               'catalog_id', service_hierarchy.catalog_id
                )) AS services
FROM service_hierarchy
         LEFT JOIN cases.sla ON sla.id = service_hierarchy.sla_id
         LEFT JOIN contacts.group AS grp ON grp.id = service_hierarchy.group_id
         LEFT JOIN directory.wbt_user AS created_by_user ON created_by_user.id = service_hierarchy.created_by
         LEFT JOIN directory.wbt_user AS updated_by_user ON updated_by_user.id = service_hierarchy.updated_by
         LEFT JOIN contacts.contact AS assignee ON assignee.id = service_hierarchy.assignee_id
GROUP BY service_hierarchy.catalog_id )`, depth)
		}
		queryBuilder = queryBuilder.Prefix(strings.TrimSuffix(prefixQuery, ","))
	}

	// Build the final SQL query and return
	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return "", nil, model.NewInternalError("postgres.catalog.query_build_error", err.Error())
	}

	return store.CompactSQL(query), args, nil
}

// Update implements store.CatalogStore.
func (s *CatalogStore) Update(rpc *model.UpdateOptions, lookup *cases.Catalog) (*cases.Catalog, error) {
	// Establish a connection to the database
	db, dbErr := s.storage.Database()
	if dbErr != nil {
		return nil, model.NewInternalError("postgres.catalog.update.database_connection_error", dbErr.Error())
	}

	// Start a transaction using the TxManager
	tx, err := db.Begin(rpc.Context)
	if err != nil {
		return nil, model.NewInternalError("postgres.catalog.update.transaction_start_error", err.Error())
	}
	txManager := store.NewTxManager(tx)   // Create a new TxManager instance
	defer txManager.Rollback(rpc.Context) // Ensure rollback on error

	// Check if rpc.Fields contains team_ids or skill_ids
	updateTeams := false
	updateSkills := false

	// Check if the fields exist in rpc.Fields
	for _, field := range rpc.Fields {
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
			rpc,
			lookup.Id,
			teamIDs,  // Pass empty slice if no team IDs are provided
			skillIDs, // Pass empty slice if no skill IDs are provided
			rpc.Session.GetUserId(),
			rpc.Time,
			rpc.Session.GetDomainId(),
		)

		// Execute the teams and skills update query and check for affected rows
		var affectedRows int
		err = txManager.QueryRow(rpc.Context, query, args...).Scan(&affectedRows)
		if err != nil {
			return nil, model.NewInternalError("postgres.catalog.update.teams_skills_update_error", err.Error())
		}

		// Optional check if no rows were affected
		if affectedRows == 0 {
			return nil, model.NewInternalError("postgres.catalog.update.no_teams_skills_affected", "No teams or skills were updated")
		}
	}

	// Build the update query for the Catalog
	query, args, err := s.buildUpdateCatalogQuery(rpc, lookup)
	if err != nil {
		return nil, model.NewInternalError("postgres.catalog.update.query_build_error", err.Error())
	}

	var (
		createdByLookup, updatedByLookup cases.Lookup
		createdAt, updatedAt             time.Time
		teamLookups, skillLookups        []byte
	)

	err = txManager.QueryRow(rpc.Context, query, args...).Scan(
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
	if err := txManager.Commit(rpc.Context); err != nil {
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

func (s *CatalogStore) buildUpdateTeamsAndSkillsQuery(
	rpc *model.UpdateOptions,
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

	// Check if "teams" is in rpc.Fields, even if teamIDs is empty
	if util.FieldExists("teams", rpc.Fields) {
		query += `
 updated_teams AS (
    INSERT INTO cases.team_catalog (catalog_id, team_id, created_by, updated_by, updated_at, dc)
        SELECT $1, unnest(NULLIF($5::bigint[], '{}')), $2, $2, $4, $3 -- created_by and updated_by are both set to $2
        ON CONFLICT (catalog_id, team_id)
            DO UPDATE SET updated_at = EXCLUDED.updated_at, updated_by = EXCLUDED.updated_by
        RETURNING catalog_id
    ),
 deleted_teams AS (
     DELETE FROM cases.team_catalog
     WHERE catalog_id = $1
       AND (
         array_length($5, 1) IS NULL -- If array is empty, delete all teams
         OR team_id != ALL ($5) -- If array is not empty, delete teams not in the array
       )
     RETURNING catalog_id
    )`
		args = append(args, pq.Array(teamIDs)) // Append team IDs to args (even if empty)
		cteAdded = true
	} else {
		// Pass an empty array if "teams" is not provided
		args = append(args, pq.Array([]int64{}))
	}

	// Check if "skills" is in rpc.Fields, even if skillIDs are empty
	if util.FieldExists("skills", rpc.Fields) {
		if cteAdded {
			query += `,` // Only add a comma if there is already a CTE defined (for teams)
		}
		query += `
 updated_skills AS (
    INSERT INTO cases.skill_catalog (catalog_id, skill_id, created_by, updated_by, updated_at, dc)
        SELECT $1, unnest(NULLIF($6::bigint[], '{}')), $2, $2, $4, $3 -- created_by and updated_by are both set to $2
        ON CONFLICT (catalog_id, skill_id)
            DO UPDATE SET updated_at = EXCLUDED.updated_at, updated_by = EXCLUDED.updated_by
        RETURNING catalog_id
    ),
 deleted_skills AS (
     DELETE FROM cases.skill_catalog
     WHERE catalog_id = $1
       AND (
         array_length($6, 1) IS NULL -- if array is empty, delete all skills
         OR skill_id != ALL ($6) -- if array is not empty, delete skills not in the array
       )
     RETURNING catalog_id
    )`
		args = append(args, pq.Array(skillIDs)) // Append skill IDs to args (even if empty)
		cteAdded = true
	} else {
		// Pass an empty array if "skills" is not provided
		args = append(args, pq.Array([]int64{}))
	}

	// Construct the final SELECT query after the CTE block
	query += `
SELECT COUNT(*)
FROM (
    ` + func() string {
		var result string
		if util.FieldExists("teams", rpc.Fields) {
			result += `SELECT catalog_id FROM updated_teams UNION ALL SELECT catalog_id FROM deleted_teams`
		}
		if util.FieldExists("skills", rpc.Fields) {
			if util.FieldExists("teams", rpc.Fields) {
				result += ` UNION ALL `
			}
			result += `SELECT catalog_id FROM updated_skills UNION ALL SELECT catalog_id FROM deleted_skills`
		}
		return result
	}() + `
) AS total_affected;`

	// Return the constructed query and arguments
	return store.CompactSQL(query), args
}

func (s *CatalogStore) buildUpdateCatalogQuery(rpc *model.UpdateOptions, lookup *cases.Catalog) (string, []interface{}, error) {
	// Start the update query with Squirrel Update Builder
	updateQueryBuilder := sq.Update("cases.service_catalog").
		PlaceholderFormat(sq.Dollar).
		Set("updated_at", rpc.Time).
		Set("updated_by", rpc.Session.GetUserId()).
		Where(sq.Eq{"id": lookup.Id, "dc": rpc.Session.GetDomainId()})

	// Dynamically set fields based on what the user wants to update
	for _, field := range rpc.Fields {
		switch field {
		case "name":
			updateQueryBuilder = updateQueryBuilder.Set("name", lookup.Name)
		case "description":
			// Use NULLIF to store NULL if description is an empty string
			updateQueryBuilder = updateQueryBuilder.Set("description", sq.Expr("NULLIF(?, '')", lookup.Description))
		case "prefix":
			updateQueryBuilder = updateQueryBuilder.Set("prefix", lookup.Prefix)
		case "code":
			// Use NULLIF to store NULL if code is an empty string
			updateQueryBuilder = updateQueryBuilder.Set("code", sq.Expr("NULLIF(?, '')", lookup.Code))
		case "state":
			updateQueryBuilder = updateQueryBuilder.Set("state", lookup.State)
		case "sla_id":
			updateQueryBuilder = updateQueryBuilder.Set("sla_id", lookup.Sla.Id)
		case "status_id":
			updateQueryBuilder = updateQueryBuilder.Set("status_id", lookup.Status.Id)
		case "close_reason_id":
			// Use NULLIF to store NULL if close_reason_id is an empty string
			updateQueryBuilder = updateQueryBuilder.Set("close_reason_id", sq.Expr("NULLIF(?, 0)", lookup.CloseReason.Id))
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
       COALESCE(close_reason.name, '')                    AS close_reason_name, -- Handle NULL close_reason as empty string
       catalog.created_by,
       COALESCE(created_by_user.name, '')                 AS created_by_name,   -- Handle NULL created_by as empty string
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

func NewCatalogStore(store store.Store) (store.CatalogStore, model.AppError) {
	if store == nil {
		return nil, model.NewInternalError("postgres.new_catalog.check.bad_arguments",
			"error creating Catalog interface to the service table, main store is nil")
	}
	return &CatalogStore{storage: store}, nil
}
