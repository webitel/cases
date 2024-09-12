package postgres

import (
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

func (s *ServiceStore) Create(ctx *model.CreateOptions, add *cases.Service) (*cases.Service, error) {
	// Establish a connection to the database
	db, dbErr := s.storage.Database()
	if dbErr != nil {
		return nil, model.NewInternalError("postgres.service.create.database_connection_error", dbErr.Error())
	}

	// Build the combined query for inserting Service and related entities
	query, args := s.buildCreateServiceQuery(ctx, add)

	// Execute the query and scan the result into the Service fields
	var createdByLookup, updatedByLookup cases.Lookup
	var createdAt, updatedAt time.Time
	var groupLookup, assigneeLookup cases.Lookup

	err := db.QueryRow(ctx.Context, query, args...).Scan(
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
func (s *ServiceStore) Delete(ctx *model.DeleteOptions) error {
	// Establish a connection to the database
	db, dbErr := s.storage.Database()
	if dbErr != nil {
		return model.NewInternalError("postgres.service.delete.db_connection_error", dbErr.Error())
	}

	// Ensure that there are IDs to delete
	if len(ctx.IDs) == 0 {
		return model.NewBadRequestError("postgres.service.delete.no_ids_provided", "No IDs provided for deletion")
	}

	// Build the delete query
	query, args := s.buildDeleteServiceQuery(ctx)

	// Execute the delete query
	res, err := db.Exec(ctx.Context, query, args...)
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
func (s *ServiceStore) List(ctx *model.SearchOptions) (*cases.ServiceList, error) {
	// Establish a connection to the database
	db, dbErr := s.storage.Database()
	if dbErr != nil {
		return nil, model.NewInternalError("postgres.service.list.database_connection_error", dbErr.Error())
	}

	// Build SQL query with filtering by root_id
	query, args, err := s.buildSearchServiceQuery(ctx)
	if err != nil {
		return nil, model.NewInternalError("postgres.service.list.query_build_error", err.Error())
	}

	// Execute the query
	rows, err := db.Query(ctx.Context, query, args...)
	if err != nil {
		return nil, model.NewInternalError("postgres.service.list.query_execution_error", err.Error())
	}
	defer rows.Close()

	// Parse the result
	var services []*cases.Service
	count := 0
	next := false

	for rows.Next() {
		if count >= ctx.Size {
			next = true
			break
		}

		var service cases.Service
		var groupLookup, assigneeLookup cases.Lookup

		err = rows.Scan(
			&service.Id, &service.Name, &service.Description, &service.Code, &service.State,
			&service.CreatedAt, &service.UpdatedAt,
			&service.Sla.Id, &service.Sla.Name,
			&groupLookup.Id, &groupLookup.Name,
			&assigneeLookup.Id, &assigneeLookup.Name,
			&service.CreatedBy.Id, &service.CreatedBy.Name,
			&service.UpdatedBy.Id, &service.UpdatedBy.Name,
			&service.HasServices, &service.RootId,
		)
		if err != nil {
			return nil, model.NewInternalError("postgres.service.list.scan_error", err.Error())
		}

		service.Group = &groupLookup
		service.Assignee = &assigneeLookup

		services = append(services, &service)
		count++
	}

	return &cases.ServiceList{
		Page:  int32(ctx.Page),
		Next:  next,
		Items: services,
	}, nil
}

// Update implements store.ServiceStore.
func (s *ServiceStore) Update(ctx *model.UpdateOptions, lookup *cases.Service) (*cases.Service, error) {
	// Establish a connection to the database
	db, dbErr := s.storage.Database()
	if dbErr != nil {
		return nil, model.NewInternalError("postgres.service.update.database_connection_error", dbErr.Error())
	}

	// Start a transaction using the TxManager
	tx, err := db.Begin(ctx.Context)
	if err != nil {
		return nil, model.NewInternalError("postgres.service.update.transaction_start_error", err.Error())
	}
	txManager := store.NewTxManager(tx)   // Create a new TxManager instance
	defer txManager.Rollback(ctx.Context) // Ensure rollback on error

	// Build the update query for the Service
	query, args, err := s.buildUpdateServiceQuery(ctx, lookup)
	if err != nil {
		return nil, model.NewInternalError("postgres.service.update.query_build_error", err.Error())
	}

	// Execute the update query for the service
	var createdByLookup, updatedByLookup cases.Lookup
	var createdAt, updatedAt time.Time
	var groupLookup, assigneeLookup cases.Lookup

	err = txManager.QueryRow(ctx.Context, query, args...).Scan(
		&lookup.Id, &lookup.Name, &createdAt,
		&lookup.Sla.Id, &lookup.Sla.Name,
		&groupLookup.Id, &groupLookup.Name,
		&assigneeLookup.Id, &assigneeLookup.Name,
		&createdByLookup.Id, &createdByLookup.Name,
		&updatedByLookup.Id, &updatedByLookup.Name, &updatedAt, &lookup.RootId,
	)
	if err != nil {
		return nil, model.NewInternalError("postgres.service.update.execution_error", err.Error())
	}

	// Commit the transaction
	if err := txManager.Commit(ctx.Context); err != nil {
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

// Helper method to build the combined insert query for Service and related entities
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
		add.RootId,                // $11: root_id
	}

	query := `
    WITH inserted_service AS (
        INSERT INTO cases.service_catalog (
            name, description, code, created_at, created_by, updated_at,
            updated_by, sla_id, group_id, assignee_id, state, dc, root_id
        ) VALUES ($1, $2, $3, $4, $5, $4, $5, $6, $7, $8, $9, $10, $11)
        RETURNING id, name, description, code, state, sla_id, group_id, assignee_id,
                  created_by, updated_by, created_at, updated_at, root_id
    )
    SELECT inserted_service.id,
           inserted_service.name,
           inserted_service.description,
           inserted_service.code,
           inserted_service.state,
           inserted_service.created_at,
           inserted_service.updated_at,
           inserted_service.sla_id,
           sla.name AS sla_name,
           inserted_service.group_id,
           grp.name AS group_name,
           inserted_service.assignee_id,
           assignee.given_name AS assignee_name,
           inserted_service.created_by,
           created_by_user.name AS created_by_name,
           inserted_service.updated_by,
           updated_by_user.name AS updated_by_name,
		   inserted_service.root_id
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
func (s *ServiceStore) buildDeleteServiceQuery(ctx *model.DeleteOptions) (string, []interface{}) {
	query := `
		DELETE FROM cases.service
		WHERE id = ANY($1) AND dc = $2
	`
	args := []interface{}{
		pq.Array(ctx.IDs),         // $1: array of service IDs to delete
		ctx.Session.GetDomainId(), // $2: domain ID to ensure proper scoping
	}

	return store.CompactSQL(query), args
}

// Helper method to build the search query for Service
func (s *ServiceStore) buildSearchServiceQuery(ctx *model.SearchOptions) (string, []interface{}, error) {
	// Initialize query builder
	queryBuilder := sq.Select(
		"service.id",
		"service.name",
		"service.description",
		"service.code",
		"service.state",
		"service.sla_id",
		"sla.name",
		"service.group_id",
		"grp.name",
		"service.assignee_id",
		"assignee.name",
		"service.created_by",
		"created_by_user.name AS created_by_name",
		"service.updated_by",
		"updated_by_user.name AS updated_by_name",
		"service.created_at",
		"service.updated_at",
		// Determine if the service has subservices
		`EXISTS (SELECT 1 FROM cases.service cs WHERE cs.root_id = service.id) AS has_subservices`,
		"service.root_id",
	).
		From("cases.service_catalog AS service").
		LeftJoin("cases.sla ON sla.id = service.sla_id").
		LeftJoin("contacts.group AS grp ON grp.id = service.group_id").
		LeftJoin("contacts.contact AS assignee ON assignee.id = service.assignee_id").
		LeftJoin("directory.wbt_user AS created_by_user ON created_by_user.id = service.created_by").
		LeftJoin("directory.wbt_user AS updated_by_user ON updated_by_user.id = service.updated_by").
		GroupBy(
			"service.id", "sla.name", "grp.name", "assignee.name", "created_by_user.name", "updated_by_user.name", "service.root_id",
		).
		PlaceholderFormat(sq.Dollar)

	// Apply filtering by catalog_id (using root_id)
	if rootID, ok := ctx.Filter["root_id"].(int64); ok && rootID > 0 {
		queryBuilder = queryBuilder.Where(sq.Eq{"service.root_id": rootID})
	}

	// Apply filtering by name
	if name, ok := ctx.Filter["name"].(string); ok && len(name) > 0 {
		substr := ctx.Match.Substring(name)
		queryBuilder = queryBuilder.Where(sq.ILike{"service.name": substr})
	}

	// Apply filtering by state
	if state, ok := ctx.Filter["state"]; ok {
		queryBuilder = queryBuilder.Where(sq.Eq{"service.state": state})
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
		return "", nil, model.NewInternalError("postgres.service.query_build_error", err.Error())
	}

	return store.CompactSQL(query), args, nil
}

// Helper method to build the combined update and select query for Service using Squirrel
func (s *ServiceStore) buildUpdateServiceQuery(ctx *model.UpdateOptions, lookup *cases.Service) (string, []interface{}, error) {
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
		case "group_id":
			updateQueryBuilder = updateQueryBuilder.Set("group_id", lookup.Group.Id)
		case "assignee_id":
			updateQueryBuilder = updateQueryBuilder.Set("assignee_id", lookup.Assignee.Id)
		}
	}

	// Convert the update query to SQL
	updateQuery, args, err := updateQueryBuilder.ToSql()
	if err != nil {
		return "", nil, err
	}

	// Now build the select query to return the updated service
	selectQueryBuilder := sq.Select(
		"service.id",
		"service.name",
		"service.description",
		"service.code",
		"service.state",
		"service.sla_id",
		"sla.name",
		"service.group_id",
		"grp.name",
		"service.assignee_id",
		"assignee.name",
		"service.created_by",
		"created_by_user.name AS created_by_name",
		"service.updated_by",
		"updated_by_user.name AS updated_by_name",
		"service.created_at",
		"service.updated_at",
		"service.root_id",
	).
		From("cases.service_catalog AS service").
		LeftJoin("cases.sla ON sla.id = service.sla_id").
		LeftJoin("contacts.group AS grp ON grp.id = service.group_id").
		LeftJoin("contacts.contact AS assignee ON assignee.id = service.assignee_id").
		LeftJoin("directory.wbt_user AS created_by_user ON created_by_user.id = service.created_by").
		LeftJoin("directory.wbt_user AS updated_by_user ON updated_by_user.id = service.updated_by").
		Where(sq.Eq{"service.id": lookup.Id, "service.dc": ctx.Session.GetDomainId()}).
		GroupBy(
			"service.id",
			"sla.name",
			"grp.name",
			"assignee.name",
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

func NewServiceStore(store store.Store) (store.ServiceStore, model.AppError) {
	if store == nil {
		return nil, model.NewInternalError("postgres.new_service.check.bad_arguments",
			"error creating Service interface to the service table, main store is nil")
	}
	return &ServiceStore{storage: store}, nil
}
