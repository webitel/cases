package postgres

import (
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/lib/pq"
	errors "github.com/webitel/cases/internal/errors"
	"github.com/webitel/cases/internal/model"
	"github.com/webitel/cases/internal/model/options"
	"github.com/webitel/cases/internal/store"
	util2 "github.com/webitel/cases/internal/store/util"
	"github.com/webitel/cases/util"
)

type ServiceStore struct {
	storage *Store
}

func (s *ServiceStore) Create(rpc options.Creator, add *model.Service) (*model.Service, error) {
	// Establish a connection to the database
	db, err := s.storage.Database()
	if err != nil {
		return nil, fmt.Errorf("failed to get database connection for service creation: %w", err)
	}
	// Build the combined query for inserting Service and related entities
	query, args, err := s.buildCreateServiceQuery(rpc, add)
	if err != nil {
		return nil, fmt.Errorf("failed to build create service query: %w", err)
	}
	var res model.Service
	err = pgxscan.Get(rpc, db, &res, query, args...)
	if err != nil {
		return nil, ParseError(err)
	}
	// Return the created Service
	return &res, nil
}

// Delete implements store.ServiceStore.
func (s *ServiceStore) Delete(rpc options.Deleter) (*model.Service, error) {
	// Establish a connection to the database
	db, err := s.storage.Database()
	if err != nil {
		return nil, fmt.Errorf("failed to get database connection for service deletion: %w", err)
	}

	// Ensure that there are IDs to delete
	if len(rpc.GetIDs()) == 0 {
		return nil, errors.InvalidArgument("no IDs provided for deletion")
	}

	// Build the delete query
	query, args := s.buildDeleteServiceQuery(rpc)

	// Execute the delete query
	res, err := db.Exec(rpc, query, args...)
	if err != nil {
		return nil, ParseError(err)
	}

	// Check how many rows were affected
	if res.RowsAffected() == 0 {
		return nil, errors.NotFound("no rows affected by delete operation")
	}

	return nil, nil
}

// List implements store.ServiceStore.
func (s *ServiceStore) List(rpc options.Searcher) ([]*model.Service, error) {
	// Establish a connection to the database
	db, err := s.storage.Database()
	if err != nil {
		return nil, fmt.Errorf("failed to get database connection for service listing: %w", err)
	}

	// Build SQL query with filtering by root_id
	query, args, err := s.buildSearchServiceQuery(rpc)
	if err != nil {
		return nil, fmt.Errorf("failed to build search service query: %w", err)
	}

	// Execute the query
	var services []*model.Service
	err = pgxscan.Select(rpc, db, &services, query, args...)
	if err != nil {
		return nil, ParseError(err)
	}

	return services, nil
}

// Update implements store.ServiceStore.
func (s *ServiceStore) Update(rpc options.Updator, lookup *model.Service) (*model.Service, error) {
	// Establish a connection to the database
	db, err := s.storage.Database()
	if err != nil {
		return nil, fmt.Errorf("failed to get database connection for service update: %w", err)
	}

	// Build the update query for the Service
	query, args, err := s.buildUpdateServiceQuery(rpc, lookup)
	if err != nil {
		return nil, fmt.Errorf("failed to build update service query: %w", err)
	}

	var res model.Service
	err = pgxscan.Get(rpc, db, &res, query, args...)
	if err != nil {
		return nil, ParseError(err)
	}
	// Return the updated Service
	return &res, nil
}

func (s *ServiceStore) buildCreateServiceQuery(rpc options.Creator, add *model.Service) (string, []interface{}, error) {
	if add.Assignee == nil || add.Assignee.GetId() == nil || *add.Assignee.GetId() == 0 {
		return "", nil, errors.InvalidArgument("assignee must be set for service creation")
	}
	if add.Group == nil || add.Group.GetId() == nil || *add.Group.GetId() == 0 {
		return "", nil, errors.InvalidArgument("group must be set for service creation")
	}
	if add.Sla == nil || add.Sla.GetId() == nil || *add.Sla.GetId() == 0 {
		return "", nil, errors.InvalidArgument("SLA must be set for service creation")
	}
	from := "inserted_service"
	insert := sq.Insert("cases.service_catalog").
		Columns(
			"name", "description", "code", "created_at", "created_by", "updated_at",
			"updated_by", "sla_id", "group_id", "assignee_id", "state", "dc", "root_id", "catalog_id",
		).
		Values(
			add.Name,
			add.Description,
			add.Code,
			rpc.RequestTime(),
			rpc.GetAuthOpts().GetUserId(),
			rpc.RequestTime(),
			rpc.GetAuthOpts().GetUserId(),
			add.Sla.Id,
			add.Group.Id,
			add.Assignee.Id,
			add.State,
			rpc.GetAuthOpts().GetDomainId(),
			add.RootId,
			add.CatalogId,
		).
		Suffix(`RETURNING *`).
		PlaceholderFormat(sq.Dollar)
	insertSQL, args, err := util2.FormAsCTE(insert, from)
	if err != nil {
		return "", nil, fmt.Errorf("failed to form CTE for service insert: %w", err)
	}
	// Build the final query with a WITH clause to return the inserted service
	slct := sq.Select().From("inserted_service").PlaceholderFormat(sq.Dollar).Prefix(insertSQL, args...)
	slct, err = s.buildSelectColumns(slct, rpc.GetFields(), from)
	if err != nil {
		return "", nil, fmt.Errorf("failed to build select columns for service creation: %w", err)
	}
	query, args, err := slct.ToSql()
	if err != nil {
		return "", nil, fmt.Errorf("failed to generate SQL for service creation: %w", err)
	}
	return query, args, nil
}

// Helper method to build the delete query for Service
func (s *ServiceStore) buildDeleteServiceQuery(rpc options.Deleter) (string, []interface{}) {
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

func (s *ServiceStore) buildSearchServiceQuery(rpc options.Searcher) (string, []interface{}, error) {
	// Initialize query builder
	queryBuilder := sq.Select().From("cases.service_catalog AS service").
		PlaceholderFormat(sq.Dollar).
		Where("service.root_id IS NOT NULL").
		Where(sq.Eq{"service.dc": rpc.GetAuthOpts().GetDomainId()})

	// Include requested fields in the SELECT clause
	queryBuilder, err := s.buildSelectColumns(queryBuilder, rpc.GetFields(), "service")
	if err != nil {
		return "", nil, fmt.Errorf("failed to build select columns: %w", err)
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
		return "", nil, fmt.Errorf("failed to build SQL query for service search: %w", err)
	}

	return util2.CompactSQL(query), args, nil
}

func applyServiceSorting(queryBuilder sq.SelectBuilder, rpc options.Searcher) sq.SelectBuilder {
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
func (s *ServiceStore) buildUpdateServiceQuery(rpc options.Updator, input *model.Service) (string, []interface{}, error) {
	// Start the update query with Squirrel Update Builder
	updateQueryBuilder := sq.Update("cases.service_catalog").
		PlaceholderFormat(sq.Dollar).
		Set("updated_at", rpc.RequestTime()).
		Set("updated_by", rpc.GetAuthOpts().GetUserId()).
		Where(sq.Eq{"id": input.Id, "dc": rpc.GetAuthOpts().GetDomainId()}).
		Suffix("RETURNING *")

	// Dynamically set fields based on what the user wants to update
	for _, field := range rpc.GetMask() {
		switch field {
		case "name":
			updateQueryBuilder = updateQueryBuilder.Set("name", input.Name)
		case "description":
			// Use NULLIF to store NULL if description is an empty string
			updateQueryBuilder = updateQueryBuilder.Set("description", input.Description)
		case "code":
			// Use NULLIF to store NULL if code is an empty string
			updateQueryBuilder = updateQueryBuilder.Set("code", input.Code)
		case "sla":
			// Use NULLIF to store NULL if sla_id is 0
			updateQueryBuilder = updateQueryBuilder.Set("sla_id", input.Sla.Id)
		case "group":
			if input.Group != nil && input.Group.Id != nil {
				updateQueryBuilder = updateQueryBuilder.Set("group_id", input.Group.Id)
			}

		case "assignee":
			if input.Assignee != nil && input.Assignee.Id != nil {
				updateQueryBuilder = updateQueryBuilder.Set("assignee_id", input.Assignee.Id)
			}
		case "state":
			updateQueryBuilder = updateQueryBuilder.Set("state", input.State)
		case "root_id":
			updateQueryBuilder = updateQueryBuilder.Set("root_id", input.RootId)
		}
	}

	// Convert the update query to SQL
	from := "updated_service"
	updateSQL, updateArgs, err := util2.FormAsCTE(updateQueryBuilder, from)
	if err != nil {
		return "", nil, fmt.Errorf("failed to form CTE for service update: %w", err)
	}

	// Now build the select query with a static SQL using a WITH clause
	selectSQL, err := s.buildSelectColumns(sq.Select().From("cases.service_catalog AS service").
		Where(sq.Eq{"service.dc": rpc.GetAuthOpts().GetDomainId()}).
		PlaceholderFormat(sq.Dollar).Prefix(updateSQL, updateArgs...), rpc.GetFields(), from)
	if err != nil {
		return "", nil, fmt.Errorf("failed to build select columns for service update: %w", err)
	}
	query, args, err := selectSQL.ToSql()
	if err != nil {
		return "", nil, fmt.Errorf("failed to generate SQL for service update: %w", err)
	}
	return query, args, nil
}

func (s *ServiceStore) buildSelectColumns(base sq.SelectBuilder, fields []string, mainTableAlias string) (sq.SelectBuilder, error) {
	if len(fields) == 0 {
		return base, nil
	}
	fields = util.DeduplicateFields(fields)
	for _, field := range fields {
		switch field {
		case "id":
			base = base.Column(util2.Ident(mainTableAlias, "id"))
		case "name":
			base = base.Column(util2.Ident(mainTableAlias, "name"))
		case "description":
			base = base.Column(util2.Ident(mainTableAlias, "description"))
		case "code":
			base = base.Column(util2.Ident(mainTableAlias, "code"))
		case "state":
			base = base.Column(util2.Ident(mainTableAlias, "state"))
		case "sla":
			base = base.Column(`jsonb_build_object(
				'id', sla.id,
				'name', sla.name
			) AS "sla"`)
			base = base.LeftJoin(fmt.Sprintf("cases.sla AS sla ON sla.id = %s AND sla.dc = %s",
				util2.Ident(mainTableAlias, "sla_id"),
				util2.Ident(mainTableAlias, "dc")))
		case "group":
			base = base.Column(`jsonb_build_object(
				'id', grp.id,
				'name', grp.name,
				'type', CASE WHEN grp.id IN (SELECT id FROM contacts.dynamic_group) THEN 'DYNAMIC' ELSE 'STATIC' END
			) AS "group"`)
			base = base.LeftJoin(fmt.Sprintf("contacts.group AS grp ON grp.id = %s AND grp.dc = %s",
				util2.Ident(mainTableAlias, "group_id"),
				util2.Ident(mainTableAlias, "dc")))
		case "assignee":
			base = base.Column(`jsonb_build_object(
				'id', assignee.id,
				'name', assignee.common_name
			) AS "assignee"`)
			base = base.LeftJoin(fmt.Sprintf("contacts.contact AS assignee ON assignee.id = %s AND assignee.dc = %s",
				util2.Ident(mainTableAlias, "assignee_id"),
				util2.Ident(mainTableAlias, "dc")))
		case "created_at":
			base = base.Column(util2.Ident(mainTableAlias, "created_at"))
		case "updated_at":
			base = base.Column(util2.Ident(mainTableAlias, "updated_at"))
		case "created_by":
			base = util2.SetUserColumn(base, mainTableAlias, "crb", "created_by")
		case "updated_by":
			base = util2.SetUserColumn(base, mainTableAlias, "upb", "updated_by")
		case "catalog_id":
			base = base.Column(util2.Ident(mainTableAlias, "catalog_id"))
		case "root_id":
			base = base.Column(util2.Ident(mainTableAlias, "root_id"))
		default:
		}
	}
	return base, nil
}

func NewServiceStore(store *Store) (store.ServiceStore, error) {
	if store == nil {
		return nil, fmt.Errorf("failed to create ServiceStore: main store is nil")
	}
	return &ServiceStore{storage: store}, nil
}
