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
	storeutil "github.com/webitel/cases/internal/store/util"
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
			add.Sla.GetId(),
			add.Group.GetId(),
			add.Assignee.GetId(),
			add.State,
			rpc.GetAuthOpts().GetDomainId(),
			add.RootId,
			add.CatalogId,
		).
		Suffix(`RETURNING *`).
		PlaceholderFormat(sq.Dollar)
	insertSQL, args, err := storeutil.FormAsCTE(insert, from)
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

	return storeutil.CompactSQL(query), args
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
	if rootIDFilters := rpc.GetFilter("root_id"); len(rootIDFilters) > 0 {
		queryBuilder = storeutil.ApplyFiltersToQuery(queryBuilder, "service.root_id", rootIDFilters)
	}

	// Updated name filter logic for consistency
	nameFilters := rpc.GetFilter("name")
	if len(nameFilters) > 0 {
		f := nameFilters[0]
		if (f.Operator == "=" || f.Operator == "") && len(f.Value) > 0 {
			queryBuilder = storeutil.AddSearchTerm(queryBuilder, f.Value, "service.name")
		}
	}

	if stateFilters := rpc.GetFilter("state"); len(stateFilters) > 0 {
		queryBuilder = storeutil.ApplyFiltersToQuery(queryBuilder, "service.state", stateFilters)
	}

	if len(rpc.GetIDs()) > 0 {
		queryBuilder = queryBuilder.Where(sq.Eq{"service.id": rpc.GetIDs()})
	}

	// Apply sorting dynamically
	queryBuilder = applyServiceSorting(queryBuilder, rpc)

	queryBuilder = storeutil.ApplyPaging(rpc.GetPage(), rpc.GetSize(), queryBuilder)

	// Build the query
	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return "", nil, fmt.Errorf("failed to build SQL query for service search: %w", err)
	}

	return storeutil.CompactSQL(query), args, nil
}

func applyServiceSorting(queryBuilder sq.SelectBuilder, rpc options.Searcher) sq.SelectBuilder {
	sortableFields := map[string]string{
		"name":        "service.name",
		"code":        "service.code",
		"description": "service.description",
		"state":       "service.state",
		"assignee":    "assignee.common_name",
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
			updateQueryBuilder = updateQueryBuilder.Set("description", input.Description)
		case "code":
			updateQueryBuilder = updateQueryBuilder.Set("code", input.Code)
		case "sla":
			updateQueryBuilder = updateQueryBuilder.Set("sla_id", input.Sla.Id)
		case "group":
			if input.Group == nil || input.Group.Id == nil {
				updateQueryBuilder = updateQueryBuilder.Set("group_id", nil)
			} else {
				updateQueryBuilder = updateQueryBuilder.Set("group_id", sq.Expr("NULLIF(?, 0)", input.Group.Id))
			}
		case "assignee":
			if input.Assignee == nil || input.Assignee.Id == nil {
				updateQueryBuilder = updateQueryBuilder.Set("assignee_id", nil)
			} else {
				updateQueryBuilder = updateQueryBuilder.Set("assignee_id", sq.Expr("NULLIF(?, 0)", input.Assignee.Id))
			}
		case "state":
			updateQueryBuilder = updateQueryBuilder.Set("state", input.State)
		case "root_id":
			updateQueryBuilder = updateQueryBuilder.Set("root_id", input.RootId)
		}
	}

	// Convert the update query to SQL
	from := "updated_service"
	updateSQL, updateArgs, err := storeutil.FormAsCTE(updateQueryBuilder, from)
	if err != nil {
		return "", nil, fmt.Errorf("failed to form CTE for service update: %w", err)
	}
	selectSql := sq.Select().From(from).Prefix(updateSQL, updateArgs...).PlaceholderFormat(sq.Dollar)
	// Now build the select query with a static SQL using a WITH clause
	selectSQL, err := s.buildSelectColumns(selectSql, rpc.GetFields(), from)
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
			base = base.Column(storeutil.Ident(mainTableAlias, "id"))
		case "name":
			base = base.Column(storeutil.Ident(mainTableAlias, "name"))
		case "description":
			base = base.Column(storeutil.Ident(mainTableAlias, "description"))
		case "code":
			base = base.Column(storeutil.Ident(mainTableAlias, "code"))
		case "state":
			base = base.Column(storeutil.Ident(mainTableAlias, "state"))
		case "sla":
			base = base.Column(`jsonb_build_object(
				'id', sla.id,
				'name', sla.name
			) AS "sla"`)
			base = base.LeftJoin(fmt.Sprintf("cases.sla AS sla ON sla.id = %s AND sla.dc = %s",
				storeutil.Ident(mainTableAlias, "sla_id"),
				storeutil.Ident(mainTableAlias, "dc")))
		case "group":
			base = base.Column(`jsonb_build_object(
				'id', grp.id,
				'name', grp.name,
				'type', CASE WHEN grp.id IS NULL THEN NULL WHEN grp.id IN (SELECT id FROM contacts.dynamic_group) THEN 'DYNAMIC' ELSE 'STATIC' END
			) AS "group"`)
			base = base.LeftJoin(fmt.Sprintf("contacts.group AS grp ON grp.id = %s AND grp.dc = %s",
				storeutil.Ident(mainTableAlias, "group_id"),
				storeutil.Ident(mainTableAlias, "dc")))
		case "assignee":
			base = base.Column(`jsonb_build_object(
				'id', assignee.id,
				'name', assignee.common_name
			) AS "assignee"`)
			base = base.LeftJoin(fmt.Sprintf("contacts.contact AS assignee ON assignee.id = %s AND assignee.dc = %s",
				storeutil.Ident(mainTableAlias, "assignee_id"),
				storeutil.Ident(mainTableAlias, "dc")))
		case "created_at":
			base = base.Column(storeutil.Ident(mainTableAlias, "created_at"))
		case "updated_at":
			base = base.Column(storeutil.Ident(mainTableAlias, "updated_at"))
		case "created_by":
			base = storeutil.SetUserColumn(base, mainTableAlias, "crb", "created_by")
		case "updated_by":
			base = storeutil.SetUserColumn(base, mainTableAlias, "upb", "updated_by")
		case "catalog_id":
			base = base.Column(storeutil.Ident(mainTableAlias, "catalog_id"))
		case "root_id":
			base = base.Column(storeutil.Ident(mainTableAlias, "root_id"))
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
