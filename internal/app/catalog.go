package app

import (
	"context"
	"github.com/webitel/cases/api/engine"
	"github.com/webitel/cases/auth"
	"google.golang.org/grpc/metadata"
	"log/slog"
	"strings"
	"time"

	"github.com/webitel/cases/api/cases"
	cerror "github.com/webitel/cases/internal/errors"
	"github.com/webitel/cases/model"
	"github.com/webitel/cases/util"
)

type CatalogService struct {
	app *App
	cases.UnimplementedCatalogsServer
	objClassName string
}

const (
	defaultCatalogFields = "id, root_id, name, description, prefix, code, state, sla, status, close_reason_group, teams, skills, created_at, created_by, updated_at, updated_by, services"
	defaultSubfields     = "id, name, description, root_id"
)

// CreateCatalog implements cases.CatalogsServer.
func (s *CatalogService) CreateCatalog(ctx context.Context, req *cases.CreateCatalogRequest) (*cases.Catalog, error) {
	// Validate required fields
	if req.Input.Name == "" {
		return nil, cerror.NewBadRequestError("catalog.create_catalog.name.required", "Catalog name is required")
	}
	if req.Input.Prefix == "" {
		return nil, cerror.NewBadRequestError("catalog.create_catalog.prefix.required", "Catalog prefix is required")
	}
	if req.Input.Sla == nil || req.Input.Sla.GetId() == 0 {
		return nil, cerror.NewBadRequestError("catalog.create_catalog.sla.required", "SLA is required")
	}
	if req.Input.Status == nil || req.Input.Status.GetId() == 0 {
		return nil, cerror.NewBadRequestError("catalog.create_catalog.status.required", "Status is required")
	}
	if req.Input.CloseReasonGroup == nil || req.Input.CloseReasonGroup.GetId() == 0 {
		return nil, cerror.NewBadRequestError("catalog.create_catalog.close_reason_group.required", "Close reason group is required")
	}
	// Define create options
	createOpts := model.CreateOptions{
		Auth:    model.GetAutherOutOfContext(ctx),
		Context: ctx,
		Time:    time.Now(),
	}

	// Define the current user as the creator and updater
	currentU := &cases.Lookup{
		Id: createOpts.GetAuthOpts().GetUserId(),
	}

	// Create a new Catalog user_auth
	catalog := &cases.Catalog{
		Name:             req.Input.Name,
		Description:      req.Input.Description,
		Prefix:           req.Input.Prefix,
		Code:             req.Input.Code,
		State:            req.Input.State,
		Sla:              req.Input.Sla,
		Status:           req.Input.Status,
		CloseReasonGroup: req.Input.CloseReasonGroup,
		CreatedBy:        currentU,
		UpdatedBy:        currentU,
	}

	// Handle multiselect fields: teams and skills
	if len(req.Input.Teams) > 0 {
		catalog.Teams = make([]*cases.Lookup, len(req.Input.Teams))
		copy(catalog.Teams, req.Input.Teams)
	}

	if len(req.Input.Skills) > 0 {
		catalog.Skills = make([]*cases.Lookup, len(req.Input.Skills))
		copy(catalog.Skills, req.Input.Skills)
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

	t := time.Now()
	deleteOpts := model.DeleteOptions{
		Auth:    model.GetAutherOutOfContext(ctx),
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
func (s *CatalogService) ListCatalogs(
	ctx context.Context,
	req *cases.ListCatalogRequest,
) (*cases.CatalogList, error) {
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
	searchOptions := &model.SearchOptions{
		IDs:     req.Id,
		Context: ctx,
		Sort:    req.Sort,
		Fields:  req.Fields,
		Page:    int(page),
		Size:    int(req.Size),
		Time:    t,
		Filter:  make(map[string]any),
		Auth:    model.GetAutherOutOfContext(ctx),
	}

	if req.State {
		searchOptions.Filter["state"] = req.State
	}

	if !searchOptions.GetAuthOpts().HasSuperPermission(auth.SuperSelectPermission) { // if user doesn't have super select permission, then apply filters
		var info metadata.MD
		var ok bool

		info, ok = metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, cerror.NewForbiddenError("internal.grpc.get_context", "Not found")
		}
		newCtx := metadata.NewOutgoingContext(ctx, info)
		res, err := s.app.engineAgentClient.SearchAgent(newCtx, &engine.SearchAgentRequest{
			Size:   -1,
			Fields: []string{"id", "team", "skills"},
			UserId: []int64{searchOptions.GetAuthOpts().GetUserId()},
		})
		if err == nil { // passive filter, if we can't receive agent's skills and teams for whatever reason then skip
			if len(res.Items) != 0 {
				var (
					agent  = res.Items[0]
					skills []int64
				)
				if team := agent.Team; team != nil {
					if team.GetId() > 0 {
						searchOptions.Filter["team"] = agent.Team.Id
					}
				}
				if agent.Skills != nil && len(agent.Skills) != 0 {
					for _, skill := range agent.Skills {
						skills = append(skills, skill.GetId())
					}
					searchOptions.Filter["skills"] = skills

				}
			}
		} else {
			slog.WarnContext(ctx, err.Error()) // log and skip
			err = nil
		}
	}

	if req.Query != "" {
		searchOptions.Filter["name"] = req.Query
		req.Fields = append(req.Fields, "searched")
	}

	catalogs, e := s.app.Store.Catalog().List(
		searchOptions,
		req.Depth,
		req.SubFields,
		req.HasSubservices,
	)
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
	if req.Input == nil {
		// TODO: what if input nil? reset all fields?
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

		if strings.HasPrefix(f, "close_reason_group") {
			if !util.ContainsField(fields, "close_reason_group_id") {
				fields = append(fields, "close_reason_group_id")
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

	// Build update options
	updateOpts := model.UpdateOptions{
		Auth:    model.GetAutherOutOfContext(ctx),
		Context: ctx,
		Fields:  fields,
		Time:    time.Now(),
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
		UpdatedBy: &cases.Lookup{
			Id: updateOpts.GetAuthOpts().GetUserId(),
		},
	}
	// Add teams if provided
	if len(req.Input.Teams) > 0 {
		catalog.Teams = req.Input.Teams
	}

	// Add skills if provided
	if len(req.Input.Skills) > 0 {
		catalog.Skills = req.Input.Skills
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
	return &CatalogService{app: app, objClassName: "dictionaries"}, nil
}
