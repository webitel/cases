package app

import (
	"context"
	"strings"
	"time"

	"github.com/webitel/cases/api/cases"
	authmodel "github.com/webitel/cases/auth/model"
	"github.com/webitel/cases/model"
)

type CatalogService struct {
	app *App
}

const (
	defaultCatalogFields = "id, root_id, name, description, prefix, code, state, sla, status, close_reason, teams, skills, created_at, created_by"
)

// CreateCatalog implements cases.CatalogsServer.
func (s *CatalogService) CreateCatalog(ctx context.Context, req *cases.CreateCatalogRequest) (*cases.Catalog, error) {
	// Validate required fields
	if req.Name == "" {
		return nil, model.NewBadRequestError("catalog.create_catalog.name.required", "Catalog name is required")
	}
	if req.Prefix == "" {
		return nil, model.NewBadRequestError("catalog.create_catalog.prefix.required", "Catalog prefix is required")
	}
	if req.SlaId == 0 {
		return nil, model.NewBadRequestError("catalog.create_catalog.sla.required", "SLA is required")
	}
	if req.StatusId == 0 {
		return nil, model.NewBadRequestError("catalog.create_catalog.status.required", "Status is required")
	}

	session, err := s.app.AuthorizeFromContext(ctx)
	if err != nil {
		return nil, model.NewUnauthorizedError("catalog.create_catalog.authorization.failed", err.Error())
	}

	// OBAC check
	accessMode := authmodel.Add
	scope := session.GetScope(model.ScopeDictionary)
	if !session.HasAccess(scope, accessMode) {
		return nil, s.app.MakeScopeError(session, scope, accessMode)
	}

	// Define the current user as the creator and updater
	currentU := &cases.Lookup{
		Id:   session.GetUserId(),
		Name: session.GetUserName(),
	}

	// Create a new Catalog model
	catalog := &cases.Catalog{
		Name:        req.Name,
		Description: req.Description,
		Prefix:      req.Prefix,
		Code:        req.Code,
		Sla:         &cases.Lookup{Id: req.SlaId},
		Status:      &cases.Lookup{Id: req.StatusId},
		CloseReason: &cases.Lookup{Id: req.CloseReasonId},
		CreatedBy:   currentU,
		UpdatedBy:   currentU,
	}

	// Handle multiselect fields: teams and skills
	if len(req.TeamIds) > 0 {
		catalog.Teams = make([]*cases.Lookup, len(req.TeamIds))
		for i, teamId := range req.TeamIds {
			catalog.Teams[i] = &cases.Lookup{Id: teamId}
		}
	}

	if len(req.SkillIds) > 0 {
		catalog.Skills = make([]*cases.Lookup, len(req.SkillIds))
		for i, skillId := range req.SkillIds {
			catalog.Skills[i] = &cases.Lookup{Id: skillId}
		}
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
		return nil, model.NewInternalError("catalog.create_catalog.store.create.failed", e.Error())
	}

	return r, nil
}

// DeleteCatalog implements cases.CatalogsServer.
func (s *CatalogService) DeleteCatalog(ctx context.Context, req *cases.DeleteCatalogRequest) (*cases.CatalogList, error) {
	if len(req.Id) == 0 {
		return nil, model.NewBadRequestError("catalog.delete_catalog.id.required", "Catalog ID is required")
	}

	session, err := s.app.AuthorizeFromContext(ctx)
	if err != nil {
		return nil, model.NewUnauthorizedError("catalog.delete_catalog.authorization.failed", err.Error())
	}

	// OBAC check
	accessMode := authmodel.Delete
	scope := session.GetScope(model.ScopeDictionary)
	if !session.HasAccess(scope, accessMode) {
		return nil, s.app.MakeScopeError(session, scope, accessMode)
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
		return nil, model.NewInternalError("catalog.delete_catalog.store.delete.failed", e.Error())
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
		return nil, model.NewUnauthorizedError("catalog.list_catalogs.authorization.failed", err.Error())
	}

	// OBAC check
	accessMode := authmodel.Read
	scope := session.GetScope(model.ScopeDictionary)
	if !session.HasAccess(scope, accessMode) {
		return nil, s.app.MakeScopeError(session, scope, accessMode)
	}

	page := req.Page
	if page == 0 {
		page = 1
	}

	fields := req.Fields
	if len(fields) == 0 {
		fields = strings.Split(defaultCatalogFields, ", ")
	}

	t := time.Now()
	searchOptions := model.SearchOptions{
		IDs:     req.Id,
		Session: session,
		Context: ctx,
		Sort:    req.Sort,
		Fields:  fields,
		Page:    int(page),
		Size:    int(req.Size),
		Time:    t,
		Filter:  make(map[string]interface{}),
	}

	if req.Q.Query != "" {
		searchOptions.Filter["name"] = req.Q.Query
	}

	catalogs, e := s.app.Store.Catalog().List(&searchOptions, req.Depth, req.Q.Fulfillment.String())
	if e != nil {
		return nil, model.NewInternalError("catalog.list_catalogs.store.list.failed", e.Error())
	}

	return catalogs, nil
}

// LocateCatalog implements cases.CatalogsServer.
func (s *CatalogService) LocateCatalog(ctx context.Context, req *cases.LocateCatalogRequest) (*cases.LocateCatalogResponse, error) {
	if req.Id == 0 {
		return nil, model.NewBadRequestError("catalog.locate_catalog.id.required", "Catalog ID is required")
	}

	listReq := &cases.ListCatalogRequest{
		Id:   []int64{req.Id},
		Page: 1,
		Size: 1,
	}

	listResp, err := s.ListCatalogs(ctx, listReq)
	if err != nil {
		return nil, model.NewInternalError("catalog.locate_catalog.list_catalogs.error", err.Error())
	}

	if len(listResp.Items) == 0 {
		return nil, model.NewNotFoundError("catalog.locate_catalog.not_found", "Catalog not found")
	}

	return &cases.LocateCatalogResponse{Catalog: listResp.Items[0]}, nil
}

// UpdateCatalog implements cases.CatalogsServer.
func (s *CatalogService) UpdateCatalog(ctx context.Context, req *cases.UpdateCatalogRequest) (*cases.Catalog, error) {
	if req.Id == 0 {
		return nil, model.NewBadRequestError("catalog.update_catalog.id.required", "Catalog ID is required")
	}

	session, err := s.app.AuthorizeFromContext(ctx)
	if err != nil {
		return nil, model.NewUnauthorizedError("catalog.update_catalog.authorization.failed", err.Error())
	}

	// OBAC check
	accessMode := authmodel.Edit
	scope := session.GetScope(model.ScopeDictionary)
	if !session.HasAccess(scope, accessMode) {
		return nil, s.app.MakeScopeError(session, scope, accessMode)
	}

	u := &cases.Lookup{
		Id:   session.GetUserId(),
		Name: session.GetUserName(),
	}

	// Build catalog from the request input
	catalog := &cases.Catalog{
		Id:          req.Id,
		Name:        req.Input.Name,
		Description: req.Input.Description,
		Prefix:      req.Input.Prefix,
		Code:        req.Input.Code,
		Sla:         &cases.Lookup{Id: req.Input.SlaId},
		Status:      &cases.Lookup{Id: req.Input.StatusId},
		CloseReason: &cases.Lookup{Id: req.Input.CloseReasonId},
		UpdatedBy:   u,
	}

	// Add teams if provided
	if len(req.Input.TeamIds) > 0 {
		catalog.Teams = make([]*cases.Lookup, len(req.Input.TeamIds))
		for i, teamId := range req.Input.TeamIds {
			catalog.Teams[i] = &cases.Lookup{Id: teamId}
		}
	}

	// Add skills if provided
	if len(req.Input.SkillIds) > 0 {
		catalog.Skills = make([]*cases.Lookup, len(req.Input.SkillIds))
		for i, skillId := range req.Input.SkillIds {
			catalog.Skills[i] = &cases.Lookup{Id: skillId}
		}
	}

	// List of fields to update
	fields := []string{"id"}

	// Validate required fields and build the list of fields for update
	for _, f := range req.XJsonMask {
		switch f {
		case "name":
			if req.Input.Name == "" {
				return nil, model.NewBadRequestError("catalog.update_catalog.name.required", "Catalog name is required and cannot be empty")
			}
			fields = append(fields, "name")
		case "prefix":
			if req.Input.Prefix == "" {
				return nil, model.NewBadRequestError("catalog.update_catalog.prefix.required", "Catalog prefix is required and cannot be empty")
			}
			fields = append(fields, "prefix")
		case "sla_id":
			if req.Input.SlaId == 0 {
				return nil, model.NewBadRequestError("catalog.update_catalog.sla.required", "Catalog SLA is required and cannot be empty")
			}
			fields = append(fields, "sla_id")
		case "status_id":
			if req.Input.StatusId == 0 {
				return nil, model.NewBadRequestError("catalog.update_catalog.status.required", "Catalog status is required and cannot be empty")
			}
			fields = append(fields, "status_id")
		case "description":
			fields = append(fields, "description")
		case "code":
			fields = append(fields, "code")
		case "close_reason_id":
			fields = append(fields, "close_reason_id")
		case "team_ids":
			fields = append(fields, "teams")
		case "skill_ids":
			fields = append(fields, "skills")
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
		return nil, model.NewInternalError("catalog.update_catalog.store.update.failed", e.Error())
	}

	return r, nil
}

// NewCatalogService creates a new CatalogService.
func NewCatalogService(app *App) (*CatalogService, model.AppError) {
	if app == nil {
		return nil, model.NewInternalError("api.config.new_catalog.args_check.app_nil", "internal is nil")
	}
	return &CatalogService{app: app}, nil
}
