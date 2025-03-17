package app

import (
	"context"
	"github.com/webitel/cases/api/cases"
	"github.com/webitel/cases/api/engine"
	"github.com/webitel/cases/auth"
	cerror "github.com/webitel/cases/internal/errors"
	"github.com/webitel/cases/model"
	grpcopts "github.com/webitel/cases/model/options/grpc"
	"github.com/webitel/cases/util"
	"google.golang.org/grpc/metadata"
	"log/slog"
)

type CatalogService struct {
	app *App
	cases.UnimplementedCatalogsServer
	objClassName string
}

var CatalogMetadata = model.NewObjectMetadata(model.ScopeDictionary, "", []*model.Field{
	{Name: "id", Default: true},
	{Name: "root_id", Default: true},
	{Name: "name", Default: true},
	{Name: "description", Default: true},
	{Name: "prefix", Default: true},
	{Name: "code", Default: true},
	{Name: "state", Default: true},
	{Name: "sla", Default: true},
	{Name: "status", Default: true},
	{Name: "close_reason_group", Default: true},
	{Name: "teams", Default: true},
	{Name: "skills", Default: true},
	{Name: "created_at", Default: true},
	{Name: "created_by", Default: true},
	{Name: "updated_at", Default: false},
	{Name: "updated_by", Default: false},
	{Name: "services", Default: true},
})

const (
	defaultSubfields = "id, name, description, root_id"
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
	createOpts, err := grpcopts.NewCreateOptions(
		ctx,
		grpcopts.WithCreateFields(req, CatalogMetadata),
	)
	if err != nil {
		return nil, NewBadRequestError(err)
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
	r, e := s.app.Store.Catalog().Create(createOpts, catalog)
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
	deleteOpts, err := grpcopts.NewDeleteOptions(ctx, grpcopts.WithDeleteIDs(req.Id))
	if err != nil {
		return nil, NewBadRequestError(err)
	}

	e := s.app.Store.Catalog().Delete(deleteOpts)
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
	opts := []grpcopts.SearchOption{
		grpcopts.WithPagination(req),
		grpcopts.WithFields(req, CatalogMetadata,
			util.DeduplicateFields,
			func(in []string) []string {
				return util.EnsureFields(in, "id", "services")
			},
		),
		grpcopts.WithIDs(req.Id),
		grpcopts.WithSort(req),
	}

	// Conditionally add search if query is not empty
	if req.Query != "" {
		opts = append(opts, grpcopts.WithSearchAsParam(req.Query))
		opts = append(opts, func(options *grpcopts.SearchOptions) error {
			options.Fields = util.EnsureFields(options.Fields, "searched")
			return nil
		})
	}
	searchOptions, err := grpcopts.NewSearchOptions(ctx, opts...)
	if err != nil {
		return nil, NewBadRequestError(err)
	}
	if req.State {
		searchOptions.AddFilter("state", req.State)
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
						searchOptions.AddFilter("team", agent.Team.Id)
					}
				}
				if agent.Skills != nil && len(agent.Skills) != 0 {
					for _, skill := range agent.Skills {
						skills = append(skills, skill.GetId())
					}
					searchOptions.AddFilter("skills", skills)

				}
			}
		} else {
			slog.WarnContext(ctx, err.Error()) // log and skip
			err = nil
		}
	}

	if req.Query != "" {
		searchOptions.AddFilter("name", req.Query)
		req.Fields = append(req.Fields, "searched")
	}

	catalogs, e := s.app.Store.Catalog().List(
		searchOptions,
		req.Depth,
		util.FieldsFunc(req.SubFields, util.InlineFields),
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

	// Build update options
	updateOpts, err := grpcopts.NewUpdateOptions(
		ctx,
		grpcopts.WithUpdateFields(req, CaseMetadata),
		grpcopts.WithUpdateMasker(req),
	)
	if err != nil {
		return nil, NewBadRequestError(err)
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
	r, e := s.app.Store.Catalog().Update(updateOpts, catalog)
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
