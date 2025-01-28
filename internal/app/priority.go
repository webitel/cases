package app

import (
	"context"
	"log/slog"

	api "github.com/webitel/cases/api/cases"
	"github.com/webitel/cases/util"

	cerror "github.com/webitel/cases/internal/error"
	"github.com/webitel/cases/model"
)

var defaultFieldsPriority = []string{
	"id", "name", "description", "color",
}

var PriorityMetadata = model.NewObjectMetadata(model.ScopeDictionary, "", []*model.Field{
	{"id", true},
	{"created_by", false},
	{"created_at", false},
	{"updated_by", false},
	{"updated_at", false},
	{"name", true},
	{"description", true},
	{"color", true},
})

type PriorityService struct {
	app *App
	api.UnimplementedPrioritiesServer
}

// CreatePriority implements api.PrioritiesServer.
func (p *PriorityService) CreatePriority(ctx context.Context, req *api.CreatePriorityRequest) (*api.Priority, error) {
	if req.Input.Name == "" {
		return nil, cerror.NewBadRequestError("app.priority.create_priority.name_required", "Priority name is required")
	}
	if req.Input.Color == "" {
		return nil, cerror.NewBadRequestError("app.priority.create_priority.color_required", "Color is required")
	}

	lookup := &api.Priority{
		Name:        req.Input.Name,
		Description: req.Input.Description,
		Color:       req.Input.Color,
	}

	fields := util.FieldsFunc(req.Fields, util.InlineFields)
	if len(fields) == 0 {
		fields = defaultFieldsPriority
	}

	createOpts, err := model.NewCreateOptions(ctx, req, PriorityMetadata)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return nil, AppInternalError
	}
	createOpts.Fields = fields

	l, err := p.app.Store.Priority().Create(createOpts, lookup)
	if err != nil {
		return nil, cerror.NewInternalError("app.priority.create_priority.store_create_failed", err.Error())
	}

	return l, nil
}

// ListPriorities implements api.PrioritiesServer.
func (p *PriorityService) ListPriorities(ctx context.Context, req *api.ListPriorityRequest) (*api.PriorityList, error) {
	searchOptions, err := model.NewSearchOptions(ctx, req, PriorityMetadata)
	searchOptions.IDs = req.Id

	searchOptions.Sort = req.Sort
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return nil, AppInternalError
	}
	if req.Q != "" {
		searchOptions.Filter["name"] = req.Q
	}

	prios, err := p.app.Store.Priority().List(searchOptions, req.NotInSla)
	if err != nil {
		return nil, cerror.NewInternalError("app.priority.list_priorities.store_list_failed", err.Error())
	}

	return prios, nil
}

// UpdatePriority implements api.PrioritiesServer.
func (p *PriorityService) UpdatePriority(ctx context.Context, req *api.UpdatePriorityRequest) (*api.Priority, error) {
	if req.Id == 0 {
		return nil, cerror.NewBadRequestError("app.priority.update_priority.id_required", "Priority ID is required")
	}

	mask := []string{}

	for _, f := range req.XJsonMask {
		switch f {
		case "name":
			mask = append(mask, "name")
			if req.Input.Name == "" {
				return nil, cerror.NewBadRequestError("app.priority.update_priority.name_required", "Priority name cannot be empty")
			}
		case "description":
			mask = append(mask, "description")
		case "color":
			mask = append(mask, "color")
			if req.Input.Color == "" {
				return nil, cerror.NewBadRequestError("app.priority.update_priority.color_required", "Color is required")
			}
		}
	}

	fields := util.FieldsFunc(req.Fields, util.InlineFields)
	if len(fields) == 0 {
		fields = defaultFieldsPriority
	}

	updateOpts, err := model.NewUpdateOptions(ctx, req, PriorityMetadata)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return nil, AppInternalError
	}
	updateOpts.Fields = fields
	updateOpts.Mask = mask

	lookup := &api.Priority{
		Id:          req.Id,
		Name:        req.Input.Name,
		Description: req.Input.Description,
		UpdatedBy:   &api.Lookup{Id: updateOpts.GetAuthOpts().GetUserId()},
		Color:       req.Input.Color,
	}

	l, err := p.app.Store.Priority().Update(updateOpts, lookup)
	if err != nil {
		return nil, cerror.NewInternalError("app.priority.update_priority.store_update_failed", err.Error())
	}

	return l, nil
}

// DeletePriority implements api.PrioritiesServer.
func (p *PriorityService) DeletePriority(ctx context.Context, req *api.DeletePriorityRequest) (*api.Priority, error) {
	if req.Id == 0 {
		return nil, cerror.NewBadRequestError("app.priority.delete_priority.id_required", "Priority ID is required")
	}
	deleteOpts, err := model.NewDeleteOptions(ctx, PriorityMetadata)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return nil, AppInternalError
	}
	deleteOpts.IDs = []int64{req.Id}

	err = p.app.Store.Priority().Delete(deleteOpts)
	if err != nil {
		return nil, cerror.NewInternalError("app.priority.delete_priority.store_delete_failed", err.Error())
	}

	return &api.Priority{Id: req.Id}, nil
}

// LocatePriority implements api.PrioritiesServer.
func (p *PriorityService) LocatePriority(ctx context.Context, req *api.LocatePriorityRequest) (*api.LocatePriorityResponse, error) {
	if req.Id == 0 {
		return nil, cerror.NewBadRequestError("app.priority.locate_priority.id_required", "Priority ID is required")
	}

	listReq := &api.ListPriorityRequest{
		Id:     []int64{req.Id},
		Fields: req.Fields,
		Page:   1,
		Size:   1,
	}

	listResp, err := p.ListPriorities(ctx, listReq)
	if err != nil {
		return nil, cerror.NewInternalError("app.priority.locate_priority.list_priorities_failed", err.Error())
	}

	if len(listResp.Items) == 0 {
		return nil, cerror.NewNotFoundError("app.priority.locate_priority.not_found", "Priority not found")
	}

	return &api.LocatePriorityResponse{Priority: listResp.Items[0]}, nil
}

func NewPriorityService(app *App) (*PriorityService, cerror.AppError) {
	if app == nil {
		return nil, cerror.NewInternalError("app.priority.new_priority_service.app_nil", "App is nil")
	}
	return &PriorityService{app: app}, nil
}
