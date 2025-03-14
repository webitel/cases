package postgres

import (
	"fmt"
	util2 "github.com/webitel/cases/internal/store/util"
	"github.com/webitel/cases/model/options"
	"time"

	"github.com/jackc/pgtype"
	"github.com/webitel/cases/internal/store/postgres/scanner"

	sq "github.com/Masterminds/squirrel"
	"github.com/lib/pq"
	"github.com/webitel/cases/api/cases"
	dberr "github.com/webitel/cases/internal/errors"
	"github.com/webitel/cases/internal/store"
	"github.com/webitel/cases/util"
)

type ServiceStore struct {
	storage *Store
}

func (s *ServiceStore) Create(rpc options.CreateOptions, add *cases.Service) (*cases.Service, error) {
	// Establish a connection to the database
	db, dbErr := s.storage.Database()
	if dbErr != nil {
		return nil, dberr.NewDBInternalError("postgres.service.create.database_connection_error", dbErr)
	}

	// Build the combined query for inserting Service and related entities
	query, args := s.buildCreateServiceQuery(rpc, add)

	var (
		createdByLookup, updatedByLookup, slaLookup cases.Lookup
		createdAt, updatedAt                        time.Time
		groupLookup                                 cases.ExtendedLookup
		assigneeLookup                              cases.Lookup
	)

	err := db.QueryRow(rpc, query, args...).Scan(
		&add.Id, &add.Name, &add.Description, &add.Code, &add.State,
		&createdAt, &updatedAt,
		&slaLookup.Id, &slaLookup.Name,
		&groupLookup.Id, &groupLookup.Name, &groupLookup.Type,
		&assigneeLookup.Id, &assigneeLookup.Name,
		&createdByLookup.Id, &createdByLookup.Name,
		&updatedByLookup.Id, &updatedByLookup.Name,
		&add.RootId, &add.CatalogId,
	)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.service.create.scan_error", err)
	}

	// Prepare the Service to return
	add.CreatedAt = util.Timestamp(createdAt)
	add.UpdatedAt = util.Timestamp(updatedAt)
	add.CreatedBy = &createdByLookup
	add.UpdatedBy = &updatedByLookup
	add.Group = &groupLookup
	add.Assignee = &assigneeLookup
	add.Sla = &slaLookup

	// Return the created Service
	return add, nil
}

// Delete implements store.ServiceStore.
func (s *ServiceStore) Delete(rpc options.DeleteOptions) error {
	// Establish a connection to the database
	db, dbErr := s.storage.Database()
	if dbErr != nil {
		return dberr.NewDBInternalError("postgres.service.delete.db_connection_error", dbErr)
	}

	// Ensure that there are IDs to delete
	if len(rpc.GetIDs()) == 0 {
		return dberr.NewDBError("postgres.service.delete.no_ids_provided", "No IDs provided for deletion")
	}

	// Build the delete query
	query, args := s.buildDeleteServiceQuery(rpc)

	// Execute the delete query
	res, err := db.Exec(rpc, query, args...)
	if err != nil {
		return dberr.NewDBInternalError("postgres.service.delete.execution_error", err)
	}

	// Check how many rows were affected
	if res.RowsAffected() == 0 {
		return dberr.NewDBNoRowsError("postgres.service.delete.no_rows_deleted")
	}

	return nil
}

// List implements store.ServiceStore.
func (s *ServiceStore) List(rpc options.SearchOptions) (*cases.ServiceList, error) {
	// Establish a connection to the database
	db, dbErr := s.storage.Database()
	if dbErr != nil {
		return nil, dberr.NewDBInternalError("postgres.service.list.database_connection_error", dbErr)
	}

	// Build SQL query with filtering by root_id
	query, args, err := s.buildSearchServiceQuery(rpc)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.service.list.query_build_error", err)
	}

	// Execute the query
	rows, err := db.Query(rpc, query, args...)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.service.list.query_execution_error", err)
	}
	defer rows.Close()

	// Parse the result
	var services []*cases.Service
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

		// Initialize service and related lookup objects conditionally
		service := &cases.Service{}
		var createdAt, updatedAt time.Time

		// Build the scan arguments for the current row
		scanArgs, postScanHandler := s.buildServiceScanArgs(service, &createdAt, &updatedAt, rpc.GetFields())

		// Scan the row into the service object
		if err := rows.Scan(scanArgs...); err != nil {
			return nil, dberr.NewDBInternalError("postgres.service.list.scan_error", err)
		}

		// Execute post-processing assignments (handle nullable fields properly)
		postScanHandler()

		// Assign the created and updated timestamp values
		service.CreatedAt = util.Timestamp(createdAt)
		service.UpdatedAt = util.Timestamp(updatedAt)

		// Append the populated service object to the services slice
		services = append(services, service)
		lCount++
	}

	return &cases.ServiceList{
		Page:  int32(rpc.GetPage()),
		Next:  next,
		Items: services,
	}, nil
}

// Update implements store.ServiceStore.
func (s *ServiceStore) Update(rpc options.UpdateOptions, lookup *cases.Service) (*cases.Service, error) {
	// Establish a connection to the database
	db, dbErr := s.storage.Database()
	if dbErr != nil {
		return nil, dberr.NewDBInternalError("postgres.service.update.database_connection_error", dbErr)
	}

	// Build the update query for the Service
	query, args, err := s.buildUpdateServiceQuery(rpc, lookup)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.service.update.query_build_error", err)
	}

	if lookup.Sla == nil {
		lookup.Sla = &cases.Lookup{}
	}

	var (
		createdByLookup, updatedByLookup, assigneeLookup cases.Lookup
		createdAt, updatedAt                             time.Time
		groupLookup                                      cases.ExtendedLookup
	)

	err = db.QueryRow(rpc, query, args...).Scan(
		&lookup.Id, &lookup.Name, &lookup.Description,
		&lookup.Code, &lookup.State, &lookup.Sla.Id,
		&lookup.Sla.Name, scanner.ScanInt64(&groupLookup.Id), scanner.ScanText(&groupLookup.Name), scanner.ScanText(&groupLookup.Type),
		scanner.ScanInt64(&assigneeLookup.Id), scanner.ScanText(&assigneeLookup.Name), &createdByLookup.Id,
		&createdByLookup.Name, &updatedByLookup.Id, &updatedByLookup.Name,
		&createdAt, &updatedAt, &lookup.RootId,
	)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.service.update.execution_error", err)
	}

	// Prepare the updated Service to return
	lookup.CreatedAt = util.Timestamp(createdAt)
	lookup.UpdatedAt = util.Timestamp(updatedAt)
	lookup.CreatedBy = &createdByLookup
	lookup.UpdatedBy = &updatedByLookup
	lookup.Group = &groupLookup
	lookup.Assignee = &assigneeLookup

	// Return the updated Service
	return lookup, nil
}

func (s *ServiceStore) buildCreateServiceQuery(rpc options.CreateOptions, add *cases.Service) (string, []interface{}) {
	var assignee, group, sla *int64
	if add.Assignee != nil && add.Assignee.GetId() != 0 {
		assignee = &add.Assignee.Id
	}
	if add.Group != nil && add.Group.GetId() != 0 {
		group = &add.Group.Id
	}
	if add.Sla != nil && add.Sla.GetId() != 0 {
		sla = &add.Sla.Id
	}
	args := []interface{}{
		add.Name,                        // $1: name
		add.Description,                 // $2: description (can be null)
		add.Code,                        // $3: code (can be null)
		rpc.RequestTime(),               // $4: created_at, updated_at
		rpc.GetAuthOpts().GetUserId(),   // $5: created_by, updated_by
		sla,                             // $6: sla_id
		group,                           // $7: group_id
		assignee,                        // $8: assignee_id
		add.State,                       // $9: state
		rpc.GetAuthOpts().GetDomainId(), // $10: domain ID
		add.RootId,                      // $11: root_id (can be null)
		add.CatalogId,                   // $12: catalog_id
	}

	query := `
   WITH inserted_service AS (
    INSERT INTO cases.service_catalog (
                                       name, description, code, created_at, created_by, updated_at,
                                       updated_by, sla_id, group_id, assignee_id, state, dc, root_id, catalog_id
        ) VALUES ($1,
                  COALESCE(NULLIF($2, ''), NULL), -- description (NULL if empty string)
                  COALESCE(NULLIF($3, ''), NULL), -- code (NULL if empty string)
                  $4, $5, $4, $5,
                  COALESCE(NULLIF($6, 0), NULL), -- sla_id (NULL if 0)
                  COALESCE(NULLIF($7, 0), NULL), -- group_id (NULL if 0)
                  COALESCE(NULLIF($8, 0), NULL), -- assignee_id (NULL if 0)
                  $9, $10,
                  COALESCE(NULLIF($11, 0), NULL), -- root_id (NULL if 0)
				  $12
                 )
        RETURNING id, name, description, code, state, sla_id, group_id, assignee_id,
            created_by, updated_by, created_at, updated_at, root_id, catalog_id)
SELECT inserted_service.id,
       inserted_service.name,
       COALESCE(inserted_service.description, '') AS description,     -- Return empty string if null
       COALESCE(inserted_service.code, '')        AS code,            -- Return empty string if null
       inserted_service.state,
       inserted_service.created_at,
       inserted_service.updated_at,
       COALESCE(inserted_service.sla_id, 0)       AS sla_id,          -- Return 0 if null
       COALESCE(sla.name, '')                     AS sla_name,        -- Return empty string if null
       COALESCE(inserted_service.group_id, 0)     AS group_id,        -- Return 0 if null
       COALESCE(grp.name, '')                     AS group_name,      -- Return empty string if null
       CASE WHEN inserted_service.group_id NOTNULL THEN(CASE WHEN inserted_service.group_id IN (SELECT id FROM contacts.dynamic_group) THEN 'DYNAMIC' ELSE 'STATIC' END) ELSE '' END AS group_type,
       COALESCE(inserted_service.assignee_id, 0)  AS assignee_id,     -- Return 0 if null
       COALESCE(assignee.given_name, '')          AS assignee_name,   -- Return empty string if null
       COALESCE(inserted_service.created_by, 0)   AS created_by,      -- Return 0 if null
       COALESCE(created_by_user.name, '')         AS created_by_name, -- Return empty string if null
       COALESCE(inserted_service.updated_by, 0)   AS updated_by,      -- Return 0 if null
       COALESCE(updated_by_user.name, '')         AS updated_by_name, -- Return empty string if null
       COALESCE(inserted_service.root_id, 0)      AS root_id,          -- Return 0 if null
	   inserted_service.catalog_id
FROM inserted_service
         LEFT JOIN cases.sla ON sla.id = inserted_service.sla_id
         LEFT JOIN contacts.group grp ON grp.id = inserted_service.group_id
         LEFT JOIN contacts.contact assignee ON assignee.id = inserted_service.assignee_id
         LEFT JOIN directory.wbt_user created_by_user ON created_by_user.id = inserted_service.created_by
         LEFT JOIN directory.wbt_user updated_by_user ON updated_by_user.id = inserted_service.updated_by;
    `

	return util2.CompactSQL(query), args
}

// Helper method to build the delete query for Service
func (s *ServiceStore) buildDeleteServiceQuery(rpc options.DeleteOptions) (string, []interface{}) {
	query := `
		DELETE FROM cases.service_catalog
		WHERE id = ANY($1) AND dc = $2
	`
	args := []interface{}{
		pq.Array(rpc.GetIDs()),          // $1: array of service IDs to delete
		rpc.GetAuthOpts().GetDomainId(), // $2: domain ID to ensure proper scoping
	}

	return util2.CompactSQL(query), args
}

func (s *ServiceStore) buildSearchServiceQuery(rpc options.SearchOptions) (string, []interface{}, error) {
	// Map of fields to their corresponding SQL expressions
	fieldMap := map[string]string{
		"id":          "service.id",
		"name":        "service.name",
		"description": "service.description",
		"root_id":     "service.root_id",
		"code":        "service.code",
		"state":       "service.state",
		"created_at":  "service.created_at",
		"updated_at":  "service.updated_at",
		"catalog_id":  "service.catalog_id",
		"created_by":  "service.created_by",
		"updated_by":  "service.updated_by",
		"sla":         "service.sla_id AS sla_id, sla.name AS sla_name",
		"group":       "service.group_id AS group_id, grp.name AS group_name, CASE WHEN service.group_id NOTNULL THEN (CASE WHEN grp.id IN (SELECT id FROM contacts.dynamic_group) THEN 'DYNAMIC' ELSE 'STATIC' END) ELSE NULL END AS group_type",
		"assignee":    "service.assignee_id AS assignee_id, ass.common_name AS assignee_name",
	}

	// Initialize query builder
	queryBuilder := sq.Select().From("cases.service_catalog AS service").
		PlaceholderFormat(sq.Dollar).
		Where("service.root_id IS NOT NULL")

	// Include requested fields in the SELECT clause
	for _, field := range rpc.GetFields() {
		if column, ok := fieldMap[field]; ok {
			queryBuilder = queryBuilder.Column(column)
		}

		// Add necessary JOINs for specific fields
		switch field {
		case "sla":
			queryBuilder = queryBuilder.LeftJoin("cases.sla ON sla.id = service.sla_id")
		case "group":
			queryBuilder = queryBuilder.LeftJoin("contacts.group AS grp ON grp.id = service.group_id")
		case "assignee":
			queryBuilder = queryBuilder.LeftJoin("contacts.contact AS ass ON ass.id = service.assignee_id")
		case "created_by":
			queryBuilder = queryBuilder.LeftJoin("directory.wbt_user AS created_by_user ON created_by_user.id = service.created_by").
				Column("COALESCE(created_by_user.name, '') AS created_by_name")
		case "updated_by":
			queryBuilder = queryBuilder.LeftJoin("directory.wbt_user AS updated_by_user ON updated_by_user.id = service.updated_by").
				Column("COALESCE(updated_by_user.name, '') AS updated_by_name")
		}
	}

	// Apply filters
	if rootID, ok := rpc.GetFilter("root_id").(int64); ok && rootID > 0 {
		queryBuilder = queryBuilder.Where(sq.Eq{"service.root_id": rootID})
	}

	if name, ok := rpc.GetFilter("name").(string); ok && len(name) > 0 {
		queryBuilder = util2.AddSearchTerm(queryBuilder, name, "service.name")
	}

	if state := rpc.GetFilter("state"); state != nil {
		queryBuilder = queryBuilder.Where(sq.Eq{"service.state": state})
	}

	if len(rpc.GetIDs()) > 0 {
		queryBuilder = queryBuilder.Where(sq.Eq{"service.id": rpc.GetIDs()})
	}

	// Apply sorting dynamically
	queryBuilder = applyServiceSorting(queryBuilder, rpc)

	queryBuilder = util2.ApplyPaging(rpc.GetPage(), rpc.GetSize(), queryBuilder)

	// Build the query
	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return "", nil, dberr.NewDBInternalError("postgres.service.query_build_error", err)
	}

	return util2.CompactSQL(query), args, nil
}

func applyServiceSorting(queryBuilder sq.SelectBuilder, rpc options.SearchOptions) sq.SelectBuilder {
	sortableFields := map[string]string{
		"name":        "service.name",
		"code":        "service.code",
		"description": "service.description",
		"state":       "service.state",
		"assignee":    "ass.common_name",
		"group":       "grp.name",
	}

	sortApplied := false

	// Loop through the provided sorting fields
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

		// Apply sorting if the field is valid in the sortableFields map
		if dbField, exists := sortableFields[sortField]; exists {
			queryBuilder = queryBuilder.OrderBy(fmt.Sprintf("%s %s", dbField, sortDirection))
			sortApplied = true
		}
	}

	// Default sorting if no valid sort fields were applied
	if !sortApplied {
		queryBuilder = queryBuilder.OrderBy("service.name ASC")
	}

	return queryBuilder
}

// Helper method to build the combined update and select query for Service using Squirrel
func (s *ServiceStore) buildUpdateServiceQuery(rpc options.UpdateOptions, input *cases.Service) (string, []interface{}, error) {
	// Start the update query with Squirrel Update Builder
	updateQueryBuilder := sq.Update("cases.service_catalog").
		PlaceholderFormat(sq.Dollar).
		Set("updated_at", rpc.RequestTime()).
		Set("updated_by", rpc.GetAuthOpts().GetUserId()).
		Where(sq.Eq{"id": input.Id, "dc": rpc.GetAuthOpts().GetDomainId()})

	// Dynamically set fields based on what the user wants to update
	for _, field := range rpc.GetMask() {
		switch field {
		case "name":
			updateQueryBuilder = updateQueryBuilder.Set("name", input.Name)
		case "description":
			// Use NULLIF to store NULL if description is an empty string
			updateQueryBuilder = updateQueryBuilder.Set("description", sq.Expr("NULLIF(?, '')", input.Description))
		case "code":
			// Use NULLIF to store NULL if code is an empty string
			updateQueryBuilder = updateQueryBuilder.Set("code", sq.Expr("NULLIF(?, '')", input.Code))
		case "sla":
			// Use NULLIF to store NULL if sla_id is 0
			updateQueryBuilder = updateQueryBuilder.Set("sla_id", sq.Expr("NULLIF(?, 0)", input.Sla.Id))
		case "group":
			var val *int64
			if input.Group != nil {
				val = &input.Group.Id
			}
			// Use NULLIF to store NULL if group_id is 0
			updateQueryBuilder = updateQueryBuilder.Set("group_id", sq.Expr("NULLIF(?, 0)", val))
		case "assignee":
			var val *int64
			if input.Assignee != nil {
				val = &input.Assignee.Id
			}
			// Use NULLIF to store NULL if assignee_id is 0
			updateQueryBuilder = updateQueryBuilder.Set("assignee_id", sq.Expr("NULLIF(?, 0)", val))
		case "state":
			updateQueryBuilder = updateQueryBuilder.Set("state", input.State)
		case "root_id":
			updateQueryBuilder = updateQueryBuilder.Set("root_id", input.RootId)
		}
	}

	// Convert the update query to SQL
	updateSQL, args, err := updateQueryBuilder.ToSql()
	if err != nil {
		return "", nil, err
	}

	// Now build the select query with a static SQL using a WITH clause
	// TODO: refactor group type
	query := fmt.Sprintf(`
WITH updated_service AS (%s
			RETURNING id, name, description, code, state, sla_id, group_id, assignee_id, created_by, updated_by, created_at, updated_at, root_id)
SELECT service.id,
       service.name,
       COALESCE(service.description, '') AS description,  -- Use COALESCE to return an empty string if description is NULL
       COALESCE(service.code, '')        AS code,         -- Use COALESCE to return an empty string if code is NULL
       service.state,
       COALESCE(service.sla_id, 0)       AS sla_id,
       COALESCE(sla.name, '')            AS sla_name,     -- Handle NULL SLA as empty string
       service.group_id,
       COALESCE(grp.name, '')            AS group_name,   -- Handle NULL group as empty string
       CASE WHEN service.group_id NOTNULL THEN(CASE WHEN grp.id IN (SELECT id FROM contacts.dynamic_group) THEN 'DYNAMIC' ELSE 'STATIC' END) ELSE NULL END AS group_type,
       service.assignee_id,
       COALESCE(assignee.given_name, '') AS assignee_name, -- Handle NULL assignee as empty string
       service.created_by,
       COALESCE(created_by_user.name, created_by_user.username) AS created_by_name,  -- Handle NULL created_by as empty string
       service.updated_by,
       COALESCE(updated_by_user.name, updated_by_user.username) AS updated_by_name,
       service.created_at,
       service.updated_at,
       service.root_id
FROM updated_service AS service
         LEFT JOIN cases.sla ON sla.id = service.sla_id
         LEFT JOIN contacts.group AS grp ON grp.id = service.group_id
         LEFT JOIN contacts.contact AS assignee ON assignee.id = service.assignee_id
         LEFT JOIN directory.wbt_user AS created_by_user ON created_by_user.id = service.created_by
         LEFT JOIN directory.wbt_user AS updated_by_user ON updated_by_user.id = service.updated_by;
	`, updateSQL)

	// Return the final combined query and arguments
	return util2.CompactSQL(query), args, nil
}

// buildServiceScanArgs builds scan arguments dynamically and returns a post-processing function.
func (s *ServiceStore) buildServiceScanArgs(
	service *cases.Service, // The service object to populate
	createdAt, updatedAt *time.Time, // Temporary variables for timestamps
	rpcFields []string, // Fields to scan dynamically
) ([]interface{}, func()) {
	scanArgs := []interface{}{}

	// Temporary variables for nullable fields
	var slaID, assigneeID, groupID, createdByID, updatedByID pgtype.Int8
	var slaName, assigneeName, groupName, groupType, createdByName, updatedByName pgtype.Text

	// Field map for dynamic scanning
	fieldMap := map[string][]any{
		"id":          {&service.Id},
		"name":        {&service.Name},
		"description": {scanner.ScanText(&service.Description)},
		"code":        {scanner.ScanText(&service.Code)},
		"state":       {&service.State},
		"created_at":  {createdAt},
		"updated_at":  {updatedAt},
		"root_id":     {&service.RootId},
		"catalog_id":  {&service.CatalogId},
	}

	// Lookup fields that require initialization
	lookupFields := map[string]func(){
		"sla":        func() { scanArgs = append(scanArgs, &slaID, &slaName) },
		"group":      func() { scanArgs = append(scanArgs, &groupID, &groupName, &groupType) },
		"assignee":   func() { scanArgs = append(scanArgs, &assigneeID, &assigneeName) },
		"created_by": func() { scanArgs = append(scanArgs, &createdByID, &createdByName) },
		"updated_by": func() { scanArgs = append(scanArgs, &updatedByID, &updatedByName) },
	}

	// Add scan arguments for regular fields
	for _, field := range rpcFields {
		if args, exists := fieldMap[field]; exists {
			scanArgs = append(scanArgs, args...)
		} else if initFunc, exists := lookupFields[field]; exists {
			initFunc()
		}
	}

	// Function to assign scanned values after scanning
	postScanHandler := func() {
		// Set Sla to nil if ID is 0 or absent
		if slaID.Status == pgtype.Present && slaID.Int != 0 {
			service.Sla = &cases.Lookup{Id: slaID.Int, Name: safePgText(slaName)}
		} else {
			service.Sla = nil
		}

		// Set Group to nil if ID is 0 or absent
		if groupID.Status == pgtype.Present && groupID.Int != 0 {
			service.Group = &cases.ExtendedLookup{Id: groupID.Int, Name: safePgText(groupName), Type: safePgText(groupType)}
		} else {
			service.Group = nil
		}

		// Set Assignee to nil if ID is 0 or absent
		if assigneeID.Status == pgtype.Present && assigneeID.Int != 0 {
			service.Assignee = &cases.Lookup{Id: assigneeID.Int, Name: safePgText(assigneeName)}
		} else {
			service.Assignee = nil
		}

		// Set CreatedBy to nil if ID is 0 or absent
		if createdByID.Status == pgtype.Present && createdByID.Int != 0 {
			service.CreatedBy = &cases.Lookup{Id: createdByID.Int, Name: safePgText(createdByName)}
		} else {
			service.CreatedBy = nil
		}

		// Set UpdatedBy to nil if ID is 0 or absent
		if updatedByID.Status == pgtype.Present && updatedByID.Int != 0 {
			service.UpdatedBy = &cases.Lookup{Id: updatedByID.Int, Name: safePgText(updatedByName)}
		} else {
			service.UpdatedBy = nil
		}
	}

	return scanArgs, postScanHandler
}

// Helper function to safely extract pgtype.Text values
func safePgText(text pgtype.Text) string {
	if text.Status == pgtype.Present {
		return text.String
	}
	return ""
}

func NewServiceStore(store *Store) (store.ServiceStore, error) {
	if store == nil {
		return nil, dberr.NewDBError("postgres.new_service.check.bad_arguments",
			"error creating Service interface to the service table, main store is nil")
	}
	return &ServiceStore{storage: store}, nil
}
