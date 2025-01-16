package app

import (
	"context"
	"strings"
	"time"

	cases "github.com/webitel/cases/api/cases"
	authmodel "github.com/webitel/cases/auth/model"
	cerror "github.com/webitel/cases/internal/error"
	"github.com/webitel/cases/model"
	"github.com/webitel/cases/util"
)

type CatalogService struct {
	app *App
	cases.UnimplementedCatalogsServer
}

const (
	defaultCatalogFields = "id, root_id, name, description, prefix, code, state, sla, status, close_reason_group, teams, skills, created_at, created_by, updated_at, updated_by, services"
	defaultSubfields     = "id, name, description, root_id"
)

// CreateCatalog implements cases.CatalogsServer.
func (s *CatalogService) CreateCatalog(ctx context.Context, req *cases.CreateCatalogRequest) (*cases.Catalog, error) {
	// Validate required fields
	if req.Name == "" {
		return nil, cerror.NewBadRequestError("catalog.create_catalog.name.required", "Catalog name is required")
	}
	if req.Prefix == "" {
		return nil, cerror.NewBadRequestError("catalog.create_catalog.prefix.required", "Catalog prefix is required")
	}
	if req.Sla == nil || req.Sla.GetId() == 0 {
		return nil, cerror.NewBadRequestError("catalog.create_catalog.sla.required", "SLA is required")
	}
	if req.Status == nil || req.Status.GetId() == 0 {
		return nil, cerror.NewBadRequestError("catalog.create_catalog.status.required", "Status is required")
	}
	if req.CloseReasonGroup == nil || req.CloseReasonGroup.GetId() == 0 {
		return nil, cerror.NewBadRequestError("catalog.create_catalog.close_reason_group.required", "Close reason group is required")
	}

	session, err := s.app.AuthorizeFromContext(ctx)
	if err != nil {
		return nil, cerror.NewUnauthorizedError("catalog.create_catalog.authorization.failed", err.Error())
	}

	// OBAC check
	accessMode := authmodel.Add
	scope := session.GetScope(model.ScopeDictionary)
	if !session.HasObacAccess(scope.Class, accessMode) {
		return nil, cerror.MakeScopeError(session.GetUserId(), scope.Class, int(accessMode))
	}

	// Define the current user as the creator and updater
	currentU := &cases.Lookup{
		Id:   session.GetUserId(),
		Name: session.GetUserName(),
	}

	// Create a new Catalog model
	catalog := &cases.Catalog{
		Name:             req.Name,
		Description:      req.Description,
		Prefix:           req.Prefix,
		Code:             req.Code,
		State:            req.State,
		Sla:              req.Sla,
		Status:           req.Status,
		CloseReasonGroup: req.CloseReasonGroup,
		CreatedBy:        currentU,
		UpdatedBy:        currentU,
	}

	// Handle multiselect fields: teams and skills
	if len(req.Teams) > 0 {
		catalog.Teams = make([]*cases.Lookup, len(req.Teams))
		copy(catalog.Teams, req.Teams)
	}

	if len(req.Skills) > 0 {
		catalog.Skills = make([]*cases.Lookup, len(req.Skills))
		copy(catalog.Skills, req.Skills)
	}

	t := time.Now()

	// Define create options
	createOpts := model.CreateOptions{
		Session: session,
		Context: ctx,
		Time:    t,
	}

	// Create the Catalog in the store
	r, e := s.app.Store.Catalog().Create(&createOpts, catalog)
	if e != nil {
		return nil, cerror.NewInternalError("catalog.create_catalog.store.create.failed", e.Error())
	}

	return r, nil
}

// DeleteCatalog implements cases.CatalogsServer.
func (s *CatalogService) DeleteCatalog(ctx context.Context, req *cases.DeleteCatalogRequest) (*cases.CatalogList, error) {
	if len(req.Id) == 0 {
		return nil, cerror.NewBadRequestError("catalog.delete_catalog.id.required", "Catalog ID is required")
	}

	session, err := s.app.AuthorizeFromContext(ctx)
	if err != nil {
		return nil, cerror.NewUnauthorizedError("catalog.delete_catalog.authorization.failed", err.Error())
	}

	// OBAC check
	accessMode := authmodel.Delete
	scope := session.GetScope(model.ScopeDictionary)
	if !session.HasObacAccess(scope.Class, accessMode) {
		return nil, cerror.MakeScopeError(session.GetUserId(), scope.Class, int(accessMode))
	}

	t := time.Now()
	deleteOpts := model.DeleteOptions{
		Session: session,
		Context: ctx,
		IDs:     req.Id,
		Time:    t,
	}

	e := s.app.Store.Catalog().Delete(&deleteOpts)
	if e != nil {
		return nil, cerror.NewInternalError("catalog.delete_catalog.store.delete.failed", e.Error())
	}

	deletedCatalogs := make([]*cases.Catalog, len(req.Id))
	for i, id := range req.Id {
		deletedCatalogs[i] = &cases.Catalog{Id: id}
	}

	return &cases.CatalogList{
		Items: deletedCatalogs,
	}, nil
}

// ListCatalogs implements cases.CatalogsServer.
func (s *CatalogService) ListCatalogs(ctx context.Context, req *cases.ListCatalogRequest) (*cases.CatalogList, error) {
	session, err := s.app.AuthorizeFromContext(ctx)
	if err != nil {
		return nil, cerror.NewUnauthorizedError("catalog.list_catalogs.authorization.failed", err.Error())
	}

	// OBAC check
	accessMode := authmodel.Read
	scope := session.GetScope(model.ScopeDictionary)
	if !session.HasObacAccess(scope.Class, accessMode) {
		return nil, cerror.MakeScopeError(session.GetUserId(), scope.Class, int(accessMode))
	}

	page := req.Page
	if page == 0 {
		page = 1
	}

	if len(req.Fields) == 0 {
		req.Fields = strings.Split(defaultCatalogFields, ", ")
	} else {
		req.Fields = util.FieldsFunc(req.Fields, util.InlineFields)
	}

	if !util.ContainsField(req.Fields, "services") {
		req.Fields = append(req.Fields, "services")
	}

	if req.Query != "" {
		req.Fields = append(req.Fields, "searched")
	}

	if len(req.SubFields) > 0 {
		req.SubFields = util.FieldsFunc(req.SubFields, util.InlineFields)
	} else if len(req.SubFields) == 0 {
		req.SubFields = strings.Split(defaultSubfields, ", ")
	}

	t := time.Now()
	searchOptions := model.SearchOptions{
		IDs:     req.Id, // TODO check placholders in DB layer
		Session: session,
		Context: ctx,
		Sort:    req.Sort,
		Fields:  req.Fields,
		Page:    int(page),
		Size:    int(req.Size),
		Time:    t,
		Filter:  make(map[string]interface{}),
	}

	if req.Query != "" {
		searchOptions.Filter["name"] = req.Query
		req.Fields = append(req.Fields, "searched")
	}

	catalogs, e := s.app.Store.Catalog().List(&searchOptions, req.Depth, req.SubFields)
	if e != nil {
		return nil, cerror.NewInternalError("catalog.list_catalogs.store.list.failed", e.Error())
	}

	return catalogs, nil
}

// LocateCatalog implements cases.CatalogsServer.
func (s *CatalogService) LocateCatalog(ctx context.Context, req *cases.LocateCatalogRequest) (*cases.LocateCatalogResponse, error) {
	if req.Id == 0 {
		return nil, cerror.NewBadRequestError("catalog.locate_catalog.id.required", "Catalog ID is required")
	}

	listReq := &cases.ListCatalogRequest{
		Id:        []int64{req.Id},
		Page:      1,
		Size:      1,
		Fields:    req.Fields,
		SubFields: req.SubFields,
	}

	listResp, err := s.ListCatalogs(ctx, listReq)
	if err != nil {
		return nil, cerror.NewInternalError("catalog.locate_catalog.list_catalogs.error", err.Error())
	}

	if len(listResp.Items) == 0 {
		return nil, cerror.NewNotFoundError("catalog.locate_catalog.not_found", "Catalog not found")
	}

	return &cases.LocateCatalogResponse{Catalog: listResp.Items[0]}, nil
}

// UpdateCatalog implements cases.CatalogsServer.
func (s *CatalogService) UpdateCatalog(ctx context.Context, req *cases.UpdateCatalogRequest) (*cases.Catalog, error) {
	if req.Id == 0 {
		return nil, cerror.NewBadRequestError("catalog.update_catalog.id.required", "Catalog ID is required")
	}

	session, err := s.app.AuthorizeFromContext(ctx)
	if err != nil {
		return nil, cerror.NewUnauthorizedError("catalog.update_catalog.authorization.failed", err.Error())
	}

	// OBAC check
	accessMode := authmodel.Edit
	scope := session.GetScope(model.ScopeDictionary)
	if !session.HasObacAccess(scope.Class, accessMode) {
		return nil, cerror.MakeScopeError(session.GetUserId(), scope.Class, int(accessMode))
	}

	u := &cases.Lookup{
		Id:   session.GetUserId(),
		Name: session.GetUserName(),
	}

	// Build catalog from the request input
	catalog := &cases.Catalog{
		Id:               req.Id,
		Name:             req.Input.Name,
		Description:      req.Input.Description,
		Prefix:           req.Input.Prefix,
		Code:             req.Input.Code,
		State:            req.Input.State,
		Sla:              req.Input.Sla,
		Status:           req.Input.Status,
		CloseReasonGroup: req.Input.CloseReasonGroup,
		UpdatedBy:        u,
	}

	// Add teams if provided
	if len(req.Input.Teams) > 0 {
		catalog.Teams = make([]*cases.Lookup, len(req.Input.Teams))
		copy(catalog.Teams, req.Input.Teams)
	}

	// Add skills if provided
	if len(req.Input.Skills) > 0 {
		catalog.Skills = make([]*cases.Lookup, len(req.Input.Skills))
		copy(catalog.Skills, req.Input.Skills)
	}

	// List of fields to update
	fields := []string{"id"}

	// Validate required fields and build the list of fields for update
	for _, f := range req.XJsonMask {
		// Handle prefixed fields
		if strings.HasPrefix(f, "skills") {
			if !util.ContainsField(fields, "skills") {
				fields = append(fields, "skills")
			}
			continue
		}

		if strings.HasPrefix(f, "teams") {
			if !util.ContainsField(fields, "teams") {
				fields = append(fields, "teams")
			}
			continue
		}

		if strings.HasPrefix(f, "close_reason") {
			if !util.ContainsField(fields, "close_reason_id") {
				fields = append(fields, "close_reason_id")
			}
			continue
		}

		if strings.HasPrefix(f, "status") {
			if req.Input.Status.GetId() == 0 {
				return nil, cerror.NewBadRequestError("catalog.update_catalog.status.required", "Catalog status is required and cannot be empty")
			}
			if !util.ContainsField(fields, "status_id") {
				fields = append(fields, "status_id")
			}
			continue
		}

		if strings.HasPrefix(f, "sla") {
			if req.Input.Sla.GetId() == 0 {
				return nil, cerror.NewBadRequestError("catalog.update_catalog.sla.required", "Catalog SLA is required and cannot be empty")
			}
			if !util.ContainsField(fields, "sla_id") {
				fields = append(fields, "sla_id")
			}
			continue
		}

		// Handle exact matches
		switch f {
		case "name":
			if req.Input.Name == "" {
				return nil, cerror.NewBadRequestError("catalog.update_catalog.name.required", "Catalog name is required and cannot be empty")
			}
			fields = append(fields, "name")
		case "prefix":
			if req.Input.Prefix == "" {
				return nil, cerror.NewBadRequestError("catalog.update_catalog.prefix.required", "Catalog prefix is required and cannot be empty")
			}
			fields = append(fields, "prefix")
		case "description":
			fields = append(fields, "description")
		case "code":
			fields = append(fields, "code")
		case "state":
			fields = append(fields, "state")
		}
	}

	// Capture current time
	t := time.Now()

	// Build update options
	updateOpts := model.UpdateOptions{
		Session: session,
		Context: ctx,
		Fields:  fields,
		Time:    t,
	}

	// Update the catalog
	r, e := s.app.Store.Catalog().Update(&updateOpts, catalog)
	if e != nil {
		return nil, cerror.NewInternalError("catalog.update_catalog.store.update.failed", e.Error())
	}

	return r, nil
}

// NewCatalogService creates a new CatalogService.
func NewCatalogService(app *App) (*CatalogService, cerror.AppError) {
	if app == nil {
		return nil, cerror.NewInternalError("api.config.new_catalog.args_check.app_nil", "internal is nil")
	}
	return &CatalogService{app: app}, nil
}
