package postgres

import (
	"fmt"
	"github.com/webitel/cases/internal/store/scanner"
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

const (
	serviceDefaultSort = "name"
)

type ServiceStore struct {
	storage store.Store
}

func (s *ServiceStore) Create(rpc *model.CreateOptions, add *cases.Service) (*cases.Service, error) {
	// Establish a connection to the database
	db, dbErr := s.storage.Database()
	if dbErr != nil {
		return nil, dberr.NewDBInternalError("postgres.service.create.database_connection_error", dbErr)
	}

	// Build the combined query for inserting Service and related entities
	query, args := s.buildCreateServiceQuery(rpc, add)

	var (
		createdByLookup, updatedByLookup cases.Lookup
		createdAt, updatedAt             time.Time
		groupLookup                      cases.ExtendedLookup
		assigneeLookup                   cases.Lookup
	)

	err := db.QueryRow(rpc.Context, query, args...).Scan(
		&add.Id, &add.Name, &add.Description, &add.Code, &add.State,
		&createdAt, &updatedAt,
		&add.Sla.Id, &add.Sla.Name,
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

	// Return the created Service
	return add, nil
}

// Delete implements store.ServiceStore.
func (s *ServiceStore) Delete(rpc *model.DeleteOptions) error {
	// Establish a connection to the database
	db, dbErr := s.storage.Database()
	if dbErr != nil {
		return dberr.NewDBInternalError("postgres.service.delete.db_connection_error", dbErr)
	}

	// Ensure that there are IDs to delete
	if len(rpc.IDs) == 0 {
		return dberr.NewDBError("postgres.service.delete.no_ids_provided", "No IDs provided for deletion")
	}

	// Build the delete query
	query, args := s.buildDeleteServiceQuery(rpc)

	// Execute the delete query
	res, err := db.Exec(rpc.Context, query, args...)
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
func (s *ServiceStore) List(rpc *model.SearchOptions) (*cases.ServiceList, error) {
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
	rows, err := db.Query(rpc.Context, query, args...)
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

		if util.ContainsField(rpc.Fields, "sla") {
			service.Sla = &cases.Lookup{}
		}

		if util.ContainsField(rpc.Fields, "group") {
			service.Group = &cases.ExtendedLookup{}
		}

		if util.ContainsField(rpc.Fields, "assignee") {
			service.Assignee = &cases.Lookup{}
		}

		if util.ContainsField(rpc.Fields, "created_by") {
			service.CreatedBy = &cases.Lookup{}
		}

		if util.ContainsField(rpc.Fields, "updated_by") {
			service.UpdatedBy = &cases.Lookup{}
		}

		var createdAt, updatedAt time.Time

		// Build the scan arguments for the current row
		scanArgs := s.buildServiceScanArgs(
			service,
			&createdAt, &updatedAt,
			rpc.Fields,
		)

		// Scan the row into the service object
		err = rows.Scan(scanArgs...)
		if err != nil {
			return nil, dberr.NewDBInternalError("postgres.service.list.scan_error", err)
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
		return nil, dberr.NewDBInternalError("postgres.service.update.database_connection_error", dbErr)
	}

	// Start a transaction using the TxManager
	tx, err := db.Begin(rpc.Context)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.service.update.transaction_start_error", err)
	}
	txManager := store.NewTxManager(tx)   // Create a new TxManager instance
	defer txManager.Rollback(rpc.Context) // Ensure rollback on error

	// Build the update query for the Service
	query, args, err := s.buildUpdateServiceQuery(rpc, lookup)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.service.update.query_build_error", err)
	}

	if lookup.Sla == nil {
		lookup.Sla = &cases.Lookup{}
	}

	var (
		createdByLookup, updatedByLookup cases.Lookup
		createdAt, updatedAt             time.Time
		groupLookup                      cases.ExtendedLookup
		assigneeLookup                   cases.Lookup
	)

	err = txManager.QueryRow(rpc.Context, query, args...).Scan(
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

	// Commit the transaction
	if err := txManager.Commit(rpc.Context); err != nil {
		return nil, dberr.NewDBInternalError("postgres.service.update.transaction_commit_error", err)
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
	var assignee, group *int64
	if add.Assignee != nil && add.Assignee.GetId() != 0 {
		assignee = &add.Assignee.Id
	}
	if add.Group != nil && add.Group.GetId() != 0 {
		group = &add.Group.Id
	}
	args := []interface{}{
		add.Name,                  // $1: name
		add.Description,           // $2: description (can be null)
		add.Code,                  // $3: code (can be null)
		rpc.Time,                  // $4: created_at, updated_at
		rpc.Session.GetUserId(),   // $5: created_by, updated_by
		add.Sla.Id,                // $6: sla_id
		group,                     // $7: group_id
		assignee,                  // $8: assignee_id
		add.State,                 // $9: state
		rpc.Session.GetDomainId(), // $10: domain ID
		add.RootId,                // $11: root_id (can be null)
		add.CatalogId,             // $12: catalog_id
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

	return store.CompactSQL(query), args
}

// Helper method to build the delete query for Service
func (s *ServiceStore) buildDeleteServiceQuery(rpc *model.DeleteOptions) (string, []interface{}) {
	query := `
		DELETE FROM cases.service_catalog
		WHERE id = ANY($1) AND dc = $2
	`
	args := []interface{}{
		pq.Array(rpc.IDs),         // $1: array of service IDs to delete
		rpc.Session.GetDomainId(), // $2: domain ID to ensure proper scoping
	}

	return store.CompactSQL(query), args
}

func (s *ServiceStore) buildSearchServiceQuery(rpc *model.SearchOptions) (string, []interface{}, error) {
	// Map of fields to their corresponding SQL expressions
	fieldMap := map[string]string{
		"id":          "service.id",
		"name":        "service.name",
		"description": "COALESCE(service.description, '') AS description",
		"root_id":     "service.root_id",
		"code":        "COALESCE(service.code, '') AS code",
		"state":       "service.state",
		"created_at":  "service.created_at",
		"updated_at":  "service.updated_at",
		"catalog_id":  "service.catalog_id",
		"created_by":  "service.created_by",
		"updated_by":  "service.updated_by",
		"sla":         "COALESCE(service.sla_id, 0) AS sla_id, COALESCE(sla.name, '') AS sla_name",
		"group":       "COALESCE(service.group_id, 0) AS group_id, COALESCE(grp.name, '') AS group_name, CASE WHEN service.group_id NOTNULL THEN(CASE WHEN grp.id IN (SELECT id FROM contacts.dynamic_group) THEN 'DYNAMIC' ELSE 'STATIC' END) ELSE NULL END AS group_type",
		"assignee":    "COALESCE(service.assignee_id, 0) AS assignee_id, COALESCE(ass.common_name, '') AS assignee_name",
	}

	// Initialize query builder
	queryBuilder := sq.Select().From("cases.service_catalog AS service").
		PlaceholderFormat(sq.Dollar).
		Where("service.root_id IS NOT NULL")

	// Include requested fields in the SELECT clause
	for _, field := range rpc.Fields {
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
	if rootID, ok := rpc.Filter["root_id"].(int64); ok && rootID > 0 {
		queryBuilder = queryBuilder.Where(sq.Eq{"service.root_id": rootID})
	}

	if name, ok := rpc.Filter["name"].(string); ok && len(name) > 0 {
		substr := util.Substring(name)
		queryBuilder = queryBuilder.Where(sq.ILike{"service.name": "%" + strings.Join(substr, "%") + "%"})
	}

	if state, ok := rpc.Filter["state"]; ok {
		queryBuilder = queryBuilder.Where(sq.Eq{"service.state": state})
	}

	if len(rpc.IDs) > 0 {
		queryBuilder = queryBuilder.Where(sq.Eq{"service.id": rpc.IDs})
	}

	// Apply sorting and pagination
	queryBuilder = queryBuilder.OrderBy("service.name ASC")
	queryBuilder = store.ApplyPaging(rpc.GetPage(), rpc.GetSize(), queryBuilder)

	// Build the query
	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return "", nil, dberr.NewDBInternalError("postgres.service.query_build_error", err)
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
			var val *int64
			if lookup.Group != nil {
				val = &lookup.Group.Id
			}
			// Use NULLIF to store NULL if group_id is 0
			updateQueryBuilder = updateQueryBuilder.Set("group_id", sq.Expr("NULLIF(?, 0)", val))
		case "assignee_id":
			var val *int64
			if lookup.Assignee != nil {
				val = &lookup.Assignee.Id
			}
			// Use NULLIF to store NULL if assignee_id is 0
			updateQueryBuilder = updateQueryBuilder.Set("assignee_id", sq.Expr("NULLIF(?, 0)", val))
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
	// TODO: refactor group type
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
       CASE WHEN service.group_id NOTNULL THEN(CASE WHEN grp.id IN (SELECT id FROM contacts.dynamic_group) THEN 'DYNAMIC' ELSE 'STATIC' END) ELSE NULL END AS group_type,
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

func (s *ServiceStore) buildServiceScanArgs(
	service *cases.Service, // The service object to populate
	createdAt, updatedAt *time.Time, // Temporary variables for created_at and updated_at
	rpcFields []string, // Fields to scan dynamically
) []interface{} {
	scanArgs := []interface{}{}

	// Field map for dynamic scanning
	fieldMap := map[string][]any{
		"id":          {&service.Id},
		"name":        {&service.Name},
		"description": {&service.Description},
		"code":        {&service.Code},
		"state":       {&service.State},
		"created_at":  {createdAt},
		"updated_at":  {updatedAt},
		"root_id":     {&service.RootId},
		"catalog_id":  {&service.CatalogId},
	}

	// Lookup fields that require initialization
	lookupFields := map[string]func(){
		"sla": func() {
			if service.Sla == nil {
				service.Sla = &cases.Lookup{}
			}
			scanArgs = append(scanArgs, &service.Sla.Id, &service.Sla.Name)
		},
		"group": func() {
			if service.Group == nil {
				service.Group = &cases.ExtendedLookup{}
			}
			scanArgs = append(scanArgs, &service.Group.Id, &service.Group.Name, scanner.ScanText(&service.Group.Type))
		},
		"assignee": func() {
			if service.Assignee == nil {
				service.Assignee = &cases.Lookup{}
			}
			scanArgs = append(scanArgs, &service.Assignee.Id, &service.Assignee.Name)
		},
		"created_by": func() {
			if service.CreatedBy == nil {
				service.CreatedBy = &cases.Lookup{}
			}
			scanArgs = append(scanArgs, &service.CreatedBy.Id, &service.CreatedBy.Name)
		},
		"updated_by": func() {
			if service.UpdatedBy == nil {
				service.UpdatedBy = &cases.Lookup{}
			}
			scanArgs = append(scanArgs, &service.UpdatedBy.Id, &service.UpdatedBy.Name)
		},
	}

	// Add scan arguments for regular fields
	for _, field := range rpcFields {
		if args, exists := fieldMap[field]; exists {
			scanArgs = append(scanArgs, args...)
		} else if initFunc, exists := lookupFields[field]; exists {
			initFunc()
		}
	}

	return scanArgs
}

func NewServiceStore(store store.Store) (store.ServiceStore, error) {
	if store == nil {
		return nil, dberr.NewDBError("postgres.new_service.check.bad_arguments",
			"error creating Service interface to the service table, main store is nil")
	}
	return &ServiceStore{storage: store}, nil
}
