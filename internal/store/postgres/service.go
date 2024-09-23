package postgres

import (
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/lib/pq"
	"github.com/webitel/cases/api/cases"
	"github.com/webitel/cases/internal/store"
	"github.com/webitel/cases/model"
	"github.com/webitel/cases/util"
)

type ServiceStore struct {
	storage store.Store
}

func (s *ServiceStore) Create(rpc *model.CreateOptions, add *cases.Service) (*cases.Service, error) {
	// Establish a connection to the database
	db, dbErr := s.storage.Database()
	if dbErr != nil {
		return nil, model.NewInternalError("postgres.service.create.database_connection_error", dbErr.Error())
	}

	// Build the combined query for inserting Service and related entities
	query, args := s.buildCreateServiceQuery(rpc, add)

	// Execute the query and scan the result into the Service fields
	var createdByLookup, updatedByLookup cases.Lookup
	var createdAt, updatedAt time.Time
	var groupLookup, assigneeLookup cases.Lookup

	err := db.QueryRow(rpc.Context, query, args...).Scan(
		&add.Id, &add.Name, &add.Description, &add.Code, &add.State,
		&createdAt, &updatedAt,
		&add.Sla.Id, &add.Sla.Name,
		&groupLookup.Id, &groupLookup.Name,
		&assigneeLookup.Id, &assigneeLookup.Name,
		&createdByLookup.Id, &createdByLookup.Name,
		&updatedByLookup.Id, &updatedByLookup.Name,
		&add.RootId,
	)
	if err != nil {
		return nil, model.NewInternalError("postgres.service.create.scan_error", err.Error())
	}

	// Prepare the Service to return
	add.CreatedAt = util.Timestamp(createdAt)
	add.UpdatedAt = util.Timestamp(updatedAt)
	add.CreatedBy = &createdByLookup
	add.UpdatedBy = &updatedByLookup
	add.Group = &groupLookup
	add.Assignee = &assigneeLookup

	// Return the created Service
	return add, nil
}

// Delete implements store.ServiceStore.
func (s *ServiceStore) Delete(rpc *model.DeleteOptions) error {
	// Establish a connection to the database
	db, dbErr := s.storage.Database()
	if dbErr != nil {
		return model.NewInternalError("postgres.service.delete.db_connection_error", dbErr.Error())
	}

	// Ensure that there are IDs to delete
	if len(rpc.IDs) == 0 {
		return model.NewBadRequestError("postgres.service.delete.no_ids_provided", "No IDs provided for deletion")
	}

	// Build the delete query
	query, args := s.buildDeleteServiceQuery(rpc)

	// Execute the delete query
	res, err := db.Exec(rpc.Context, query, args...)
	if err != nil {
		return model.NewInternalError("postgres.service.delete.execution_error", err.Error())
	}

	// Check how many rows were affected
	if res.RowsAffected() == 0 {
		return model.NewNotFoundError("postgres.service.delete.no_rows_deleted", "No Service entries were deleted")
	}

	return nil
}

// List implements store.ServiceStore.
func (s *ServiceStore) List(rpc *model.SearchOptions) (*cases.ServiceList, error) {
	// Establish a connection to the database
	db, dbErr := s.storage.Database()
	if dbErr != nil {
		return nil, model.NewInternalError("postgres.service.list.database_connection_error", dbErr.Error())
	}

	// Build SQL query with filtering by root_id
	query, args, err := s.buildSearchServiceQuery(rpc)
	if err != nil {
		return nil, model.NewInternalError("postgres.service.list.query_build_error", err.Error())
	}

	// Execute the query
	rows, err := db.Query(rpc.Context, query, args...)
	if err != nil {
		return nil, model.NewInternalError("postgres.service.list.query_execution_error", err.Error())
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
		if !fetchAll && lCount >= rpc.GetSize() {
			next = true
			break
		}

		// Create service and related lookup objects
		service := &cases.Service{
			Sla:      &cases.Lookup{},
			Group:    &cases.Lookup{},
			Assignee: &cases.Lookup{},
		}
		createdBy, updatedBy := &cases.Lookup{}, &cases.Lookup{}
		var createdAt, updatedAt time.Time

		// Build the scan arguments for the current row
		scanArgs := s.buildServiceScanArgs(service, createdBy, updatedBy, &createdAt, &updatedAt, service.Group, service.Assignee)

		// Scan the row into the service object
		err = rows.Scan(scanArgs...)
		if err != nil {
			return nil, model.NewInternalError("postgres.service.list.scan_error", err.Error())
		}

		// Assign the created and updated timestamp values
		service.CreatedAt = util.Timestamp(createdAt)
		service.UpdatedAt = util.Timestamp(updatedAt)

		// Append the populated service object to the services slice
		services = append(services, service)
		lCount++
	}

	return &cases.ServiceList{
		Page:  int32(rpc.Page),
		Next:  next,
		Items: services,
	}, nil
}

// Update implements store.ServiceStore.
func (s *ServiceStore) Update(rpc *model.UpdateOptions, lookup *cases.Service) (*cases.Service, error) {
	// Establish a connection to the database
	db, dbErr := s.storage.Database()
	if dbErr != nil {
		return nil, model.NewInternalError("postgres.service.update.database_connection_error", dbErr.Error())
	}

	// Start a transaction using the TxManager
	tx, err := db.Begin(rpc.Context)
	if err != nil {
		return nil, model.NewInternalError("postgres.service.update.transaction_start_error", err.Error())
	}
	txManager := store.NewTxManager(tx)   // Create a new TxManager instance
	defer txManager.Rollback(rpc.Context) // Ensure rollback on error

	// Build the update query for the Service
	query, args, err := s.buildUpdateServiceQuery(rpc, lookup)
	if err != nil {
		return nil, model.NewInternalError("postgres.service.update.query_build_error", err.Error())
	}

	// Execute the update query for the service
	var createdByLookup, updatedByLookup cases.Lookup
	var createdAt, updatedAt time.Time
	var groupLookup, assigneeLookup cases.Lookup

	err = txManager.QueryRow(rpc.Context, query, args...).Scan(
		&lookup.Id, &lookup.Name, &lookup.Description,
		&lookup.Code, &lookup.State, &lookup.Sla.Id,
		&lookup.Sla.Name, &groupLookup.Id, &groupLookup.Name,
		&assigneeLookup.Id, &assigneeLookup.Name, &createdByLookup.Id,
		&createdByLookup.Name, &updatedByLookup.Id, &updatedByLookup.Name,
		&createdAt, &updatedAt, &lookup.RootId,
	)
	if err != nil {
		return nil, model.NewInternalError("postgres.service.update.execution_error", err.Error())
	}

	// Commit the transaction
	if err := txManager.Commit(rpc.Context); err != nil {
		return nil, model.NewInternalError("postgres.service.update.transaction_commit_error", err.Error())
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

func (s *ServiceStore) buildCreateServiceQuery(rpc *model.CreateOptions, add *cases.Service) (string, []interface{}) {
	args := []interface{}{
		add.Name,                  // $1: name
		add.Description,           // $2: description (can be null)
		add.Code,                  // $3: code (can be null)
		rpc.Time,                  // $4: created_at, updated_at
		rpc.Session.GetUserId(),   // $5: created_by, updated_by
		add.Sla.Id,                // $6: sla_id
		add.Group.Id,              // $7: group_id
		add.Assignee.Id,           // $8: assignee_id
		add.State,                 // $9: state
		rpc.Session.GetDomainId(), // $10: domain ID
		add.RootId,                // $11: root_id (can be null)
	}

	query := `
   WITH inserted_service AS (
    INSERT INTO cases.service_catalog (
                                       name, description, code, created_at, created_by, updated_at,
                                       updated_by, sla_id, group_id, assignee_id, state, dc, root_id
        ) VALUES ($1,
                  COALESCE(NULLIF($2, ''), NULL), -- description (NULL if empty string)
                  COALESCE(NULLIF($3, ''), NULL), -- code (NULL if empty string)
                  $4, $5, $4, $5,
                  COALESCE(NULLIF($6, 0), NULL), -- sla_id (NULL if 0)
                  COALESCE(NULLIF($7, 0), NULL), -- group_id (NULL if 0)
                  COALESCE(NULLIF($8, 0), NULL), -- assignee_id (NULL if 0)
                  $9, $10,
                  COALESCE(NULLIF($11, 0), NULL) -- root_id (NULL if 0)
                 )
        RETURNING id, name, description, code, state, sla_id, group_id, assignee_id,
            created_by, updated_by, created_at, updated_at, root_id)
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
       COALESCE(inserted_service.assignee_id, 0)  AS assignee_id,     -- Return 0 if null
       COALESCE(assignee.given_name, '')          AS assignee_name,   -- Return empty string if null
       COALESCE(inserted_service.created_by, 0)   AS created_by,      -- Return 0 if null
       COALESCE(created_by_user.name, '')         AS created_by_name, -- Return empty string if null
       COALESCE(inserted_service.updated_by, 0)   AS updated_by,      -- Return 0 if null
       COALESCE(updated_by_user.name, '')         AS updated_by_name, -- Return empty string if null
       COALESCE(inserted_service.root_id, 0)      AS root_id          -- Return 0 if null
FROM inserted_service
         LEFT JOIN cases.sla ON sla.id = inserted_service.sla_id
         LEFT JOIN contacts.group grp ON grp.id = inserted_service.group_id
         LEFT JOIN contacts.contact assignee ON assignee.id = inserted_service.assignee_id
         LEFT JOIN directory.wbt_user created_by_user ON created_by_user.id = inserted_service.created_by
         LEFT JOIN directory.wbt_user updated_by_user ON updated_by_user.id = inserted_service.updated_by;
    `

	return store.CompactSQL(query), args
}

// Helper method to build the delete query for Service
func (s *ServiceStore) buildDeleteServiceQuery(rpc *model.DeleteOptions) (string, []interface{}) {
	query := `
		DELETE FROM cases.service
		WHERE id = ANY($1) AND dc = $2
	`
	args := []interface{}{
		pq.Array(rpc.IDs),         // $1: array of service IDs to delete
		rpc.Session.GetDomainId(), // $2: domain ID to ensure proper scoping
	}

	return store.CompactSQL(query), args
}

// Helper method to build the search query for Service
func (s *ServiceStore) buildSearchServiceQuery(rpc *model.SearchOptions) (string, []interface{}, error) {
	convertedIds := rpc.FieldsUtil.Int64SliceToStringSlice(rpc.IDs)
	ids := rpc.FieldsUtil.FieldsFunc(convertedIds, rpc.FieldsUtil.InlineFields)

	// Initialize query builder with COALESCE for optional fields
	queryBuilder := sq.Select(
		"service.id",
		"service.name", // Name can't be null
		"COALESCE(service.description, '') AS description",      // Default to empty string if NULL
		"COALESCE(service.code, '') AS code",                    // Default to empty string if NULL
		"service.state",                                         // State can't be null
		"COALESCE(service.sla_id, 0) AS sla_id",                 // Default to 0 if NULL
		"COALESCE(sla.name, '') AS sla_name",                    // Default to empty string if NULL
		"COALESCE(service.group_id, 0) AS group_id",             // Default to 0 if NULL
		"COALESCE(grp.name, '') AS group_name",                  // Default to empty string if NULL
		"COALESCE(service.assignee_id, 0) AS assignee_id",       // Default to 0 if NULL
		"COALESCE(assignee.given_name, '') AS assignee_name",    // Default to empty string if NULL
		"service.created_by",                                    // created_by can't be null
		"COALESCE(created_by_user.name, '') AS created_by_name", // Default to empty string if NULL
		"service.updated_by",                                    // updated_by can't be null
		"COALESCE(updated_by_user.name, '') AS updated_by_name", // Default to empty string if NULL
		"service.created_at",                                    // created_at can't be null
		"service.updated_at",                                    // updated_at can't be null
		// Determine if the service has subservices
		`EXISTS (SELECT 1 FROM cases.service_catalog cs WHERE cs.root_id = service.id) AS has_subservices`,
		"service.root_id AS root_id",
	).
		From("cases.service_catalog AS service").
		LeftJoin("cases.sla ON sla.id = service.sla_id").
		LeftJoin("contacts.group AS grp ON grp.id = service.group_id").
		LeftJoin("contacts.contact AS assignee ON assignee.id = service.assignee_id").
		LeftJoin("directory.wbt_user AS created_by_user ON created_by_user.id = service.created_by").
		LeftJoin("directory.wbt_user AS updated_by_user ON updated_by_user.id = service.updated_by").
		GroupBy(
			"service.id", "sla.name", "grp.name", "assignee.given_name", "created_by_user.name", "updated_by_user.name", "service.root_id",
		).
		Where("service.root_id IS NOT NULL").
		PlaceholderFormat(sq.Dollar)

	// Apply filtering by root_id (using root_id from context)
	if rootID, ok := rpc.Filter["root_id"].(int64); ok && rootID > 0 {
		queryBuilder = queryBuilder.Where(sq.Eq{"service.root_id": rootID})
	}

	// Apply filtering by name (using case-insensitive matching)
	if name, ok := rpc.Filter["name"].(string); ok && len(name) > 0 {
		substr := rpc.Match.Substring(name)
		queryBuilder = queryBuilder.Where(sq.ILike{"service.name": substr})
	}

	// Apply filtering by state
	if state, ok := rpc.Filter["state"]; ok {
		queryBuilder = queryBuilder.Where(sq.Eq{"service.state": state})
	}

	// Apply filtering by IDs if provided
	if len(ids) > 0 {
		queryBuilder = queryBuilder.Where(sq.Eq{"service.id": ids})
	}

	// Apply sorting based on context
	for _, sort := range rpc.Sort {
		queryBuilder = queryBuilder.OrderBy(sort)
	}

	// Get size and page for pagination
	size := rpc.GetSize()
	page := rpc.GetPage()

	// Apply offset only if page > 1
	if rpc.Page > 1 {
		queryBuilder = queryBuilder.Offset(uint64((page - 1) * size))
	}

	// Apply limit for pagination
	if size != -1 {
		queryBuilder = queryBuilder.Limit(uint64(size + 1)) // Request one more record to check if there's a next page
	}

	// Build SQL query
	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return "", nil, model.NewInternalError("postgres.service.query_build_error", err.Error())
	}

	return store.CompactSQL(query), args, nil
}

// Helper method to build the combined update and select query for Service using Squirrel
func (s *ServiceStore) buildUpdateServiceQuery(rpc *model.UpdateOptions, lookup *cases.Service) (string, []interface{}, error) {
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
		case "code":
			// Use NULLIF to store NULL if code is an empty string
			updateQueryBuilder = updateQueryBuilder.Set("code", sq.Expr("NULLIF(?, '')", lookup.Code))
		case "sla_id":
			// Use NULLIF to store NULL if sla_id is 0
			updateQueryBuilder = updateQueryBuilder.Set("sla_id", sq.Expr("NULLIF(?, 0)", lookup.Sla.Id))
		case "group_id":
			// Use NULLIF to store NULL if group_id is 0
			updateQueryBuilder = updateQueryBuilder.Set("group_id", sq.Expr("NULLIF(?, 0)", lookup.Group.Id))
		case "assignee_id":
			// Use NULLIF to store NULL if assignee_id is 0
			updateQueryBuilder = updateQueryBuilder.Set("assignee_id", sq.Expr("NULLIF(?, 0)", lookup.Assignee.Id))
		case "state":
			updateQueryBuilder = updateQueryBuilder.Set("state", lookup.State)
		case "root_id":
			updateQueryBuilder = updateQueryBuilder.Set("root_id", lookup.RootId)
		}
	}

	// Convert the update query to SQL
	updateSQL, args, err := updateQueryBuilder.ToSql()
	if err != nil {
		return "", nil, err
	}

	// Now build the select query with a static SQL using a WITH clause
	query := fmt.Sprintf(`
WITH updated_service AS (%s
			RETURNING id, name, description, code, state, sla_id, group_id, assignee_id, created_by, updated_by, created_at, updated_at, root_id)
SELECT service.id,
       service.name,
       COALESCE(service.description, '') AS description,  -- Use COALESCE to return an empty string if description is NULL
       COALESCE(service.code, '')        AS code,         -- Use COALESCE to return an empty string if code is NULL
       service.state,
       service.sla_id,
       COALESCE(sla.name, '')            AS sla_name,     -- Handle NULL SLA as empty string
       service.group_id,
       COALESCE(grp.name, '')            AS group_name,   -- Handle NULL group as empty string
       service.assignee_id,
       COALESCE(assignee.given_name, '') AS assignee_name, -- Handle NULL assignee as empty string
       service.created_by,
       COALESCE(created_by_user.name, '') AS created_by_name,  -- Handle NULL created_by as empty string
       service.updated_by,
       updated_by_user.name              AS updated_by_name,
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
	return store.CompactSQL(query), args, nil
}

// buildServiceScanArgs prepares scan arguments for populating a Service object.
func (s *ServiceStore) buildServiceScanArgs(
	service *cases.Service, // The service object to populate
	createdBy, updatedBy *cases.Lookup, // Lookup objects for created_by and updated_by
	createdAt, updatedAt *time.Time, // Temporary variables for created_at and updated_at
	groupLookup, assigneeLookup *cases.Lookup, // Lookup objects for group and assignee
) []interface{} {
	return []interface{}{
		// Service fields
		&service.Id,          // Service ID
		&service.Name,        // Service name
		&service.Description, // Service description
		&service.Code,        // Service code
		&service.State,       // Service state

		// SLA fields
		&service.Sla.Id,   // SLA ID
		&service.Sla.Name, // SLA name

		// Group fields
		&groupLookup.Id,   // Group ID
		&groupLookup.Name, // Group name

		// Assignee fields
		&assigneeLookup.Id,   // Assignee ID
		&assigneeLookup.Name, // Assignee name

		// Created and updated by fields
		&createdBy.Id,   // Created by user ID
		&createdBy.Name, // Created by user name
		&updatedBy.Id,   // Updated by user ID
		&updatedBy.Name, // Updated by user name

		// Timestamps
		createdAt, // Created at timestamp
		updatedAt, // Updated at timestamp

		// Has services and root ID fields
		&service.HasServices, // Whether service has related services
		&service.RootId,      // Root service ID
	}
}

func NewServiceStore(store store.Store) (store.ServiceStore, model.AppError) {
	if store == nil {
		return nil, model.NewInternalError("postgres.new_service.check.bad_arguments",
			"error creating Service interface to the service table, main store is nil")
	}
	return &ServiceStore{storage: store}, nil
}
