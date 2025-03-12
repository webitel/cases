package app

import (
	"context"
	api "github.com/webitel/cases/api/cases"
	cerror "github.com/webitel/cases/internal/errors"
	"github.com/webitel/cases/model"
	grpcopts "github.com/webitel/cases/model/options/grpc"
	"github.com/webitel/cases/util"
)

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

	createOpts, err := grpcopts.NewCreateOptions(
		ctx,
		grpcopts.WithCreateFields(req, PriorityMetadata),
	)
	if err != nil {
		return nil, NewBadRequestError(err)
	}

	l, err := p.app.Store.Priority().Create(createOpts, lookup)
	if err != nil {
		return nil, cerror.NewInternalError("app.priority.create_priority.store_create_failed", err.Error())
	}

	return l, nil
}

// ListPriorities implements api.PrioritiesServer.
func (p *PriorityService) ListPriorities(ctx context.Context, req *api.ListPriorityRequest) (*api.PriorityList, error) {
	searchOpts, err := grpcopts.NewSearchOptions(
		ctx,
		grpcopts.WithSearch(req),
		grpcopts.WithPagination(req),
		grpcopts.WithFields(req, PriorityMetadata,
			util.DeduplicateFields,
			util.EnsureIdField,
		),
		grpcopts.WithSort(req),
		grpcopts.WithIDs(req.GetId()),
	)
	if err != nil {
		return nil, NewBadRequestError(err)
	}
	searchOpts.AddFilter("name", req.Q)

	prios, err := p.app.Store.Priority().List(searchOpts, req.NotInSla, req.InSlaCond)
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

	updateOpts, err := grpcopts.NewUpdateOptions(
		ctx,
		grpcopts.WithUpdateFields(req, PriorityMetadata),
		grpcopts.WithUpdateMasker(req),
	)
	if err != nil {
		return nil, NewBadRequestError(err)
	}

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
	deleteOpts, err := grpcopts.NewDeleteOptions(ctx, grpcopts.WithDeleteID(req.Id))
	if err != nil {
		return nil, NewBadRequestError(err)
	}

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
