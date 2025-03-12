package app

import (
	"context"
	_go "github.com/webitel/cases/api/cases"
	cerror "github.com/webitel/cases/internal/errors"
	"github.com/webitel/cases/model"
	grpcopts "github.com/webitel/cases/model/options/grpc"
	"github.com/webitel/cases/util"
)

type CloseReasonGroupService struct {
	app *App
	_go.UnimplementedCloseReasonGroupsServer
	objClassName string
}

var CloseReasonGroupMetadata = model.NewObjectMetadata(model.ScopeDictionary, "", []*model.Field{
	{"id", true},
	{"created_by", true},
	{"created_at", true},
	{"updated_by", false},
	{"updated_at", false},
	{"name", true},
	{"description", true},
})

func (s CloseReasonGroupService) CreateCloseReasonGroup(
	ctx context.Context,
	req *_go.CreateCloseReasonGroupRequest,
) (*_go.CloseReasonGroup, error) {
	// Validate required fields
	if req.Input.Name == "" {
		return nil, cerror.NewBadRequestError("close_reason_group_service.create_close_reason_group.name.required", "Lookup name is required")
	}

	createOpts, err := grpcopts.NewCreateOptions(
		ctx,
		grpcopts.WithCreateFields(req, CloseReasonGroupMetadata),
	)
	if err != nil {
		return nil, NewBadRequestError(err)
	}

	input := &_go.CloseReasonGroup{
		Name:        req.Input.Name,
		Description: req.Input.Description,
	}

	// Create the close reason group in the store
	res, err := s.app.Store.CloseReasonGroup().Create(createOpts, input)
	if err != nil {
		return nil, cerror.NewInternalError("close_reason_group_service.create_close_reason_group.store.create.failed", err.Error())
	}

	return res, nil
}

func (s CloseReasonGroupService) ListCloseReasonGroups(
	ctx context.Context,
	req *_go.ListCloseReasonGroupsRequest,
) (*_go.CloseReasonGroupList, error) {
	searchOpts, err := grpcopts.NewSearchOptions(
		ctx,
		grpcopts.WithSearch(req),
		grpcopts.WithPagination(req),
		grpcopts.WithFields(req, CloseReasonGroupMetadata,
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

	res, err := s.app.Store.CloseReasonGroup().List(searchOpts)
	if err != nil {
		return nil, cerror.NewInternalError("close_reason_group_service.list_close_reason_groups.store.list.failed", err.Error())
	}

	return res, nil
}

func (s CloseReasonGroupService) UpdateCloseReasonGroup(
	ctx context.Context,
	req *_go.UpdateCloseReasonGroupRequest,
) (*_go.CloseReasonGroup, error) {
	// Validate required fields
	if req.Id == 0 {
		return nil, cerror.NewBadRequestError("close_reason_group_service.update_close_reason_group.id.required", "Lookup ID is required")
	}

	updateOpts, err := grpcopts.NewUpdateOptions(
		ctx,
		grpcopts.WithUpdateFields(req, CloseReasonGroupMetadata),
		grpcopts.WithUpdateMasker(req),
	)
	if err != nil {
		return nil, NewBadRequestError(err)
	}

	// Update lookup user_auth
	input := &_go.CloseReasonGroup{
		Id:          req.Id,
		Name:        req.Input.Name,
		Description: req.Input.Description,
	}

	// Update the lookup in the store
	res, err := s.app.Store.CloseReasonGroup().Update(updateOpts, input)
	if err != nil {
		return nil, cerror.NewInternalError("close_reason_group_service.update_close_reason_group.store.update.failed", err.Error())
	}

	return res, nil
}

func (s CloseReasonGroupService) DeleteCloseReasonGroup(
	ctx context.Context,
	req *_go.DeleteCloseReasonGroupRequest,
) (*_go.CloseReasonGroup, error) {
	// Validate required fields
	if req.Id == 0 {
		return nil, cerror.NewBadRequestError("close_reason_group_service.delete_close_reason_group.id.required", "Lookup ID is required")
	}

	deleteOpts, err := grpcopts.NewDeleteOptions(ctx, grpcopts.WithDeleteID(req.Id))
	if err != nil {
		return nil, NewBadRequestError(err)
	}

	// Delete the lookup in the store
	err = s.app.Store.CloseReasonGroup().Delete(deleteOpts)
	if err != nil {
		return nil, cerror.NewInternalError("close_reason_group_service.delete_close_reason_group.store.delete.failed", err.Error())
	}

	return &(_go.CloseReasonGroup{Id: req.Id}), nil
}

func (s CloseReasonGroupService) LocateCloseReasonGroup(
	ctx context.Context,
	req *_go.LocateCloseReasonGroupRequest,
) (*_go.LocateCloseReasonGroupResponse, error) {
	// Validate required fields
	if req.Id == 0 {
		return nil, cerror.NewBadRequestError("close_reason_group_service.locate_close_reason_group.id.required", "Lookup ID is required")
	}

	fields := util.FieldsFunc(req.Fields, util.InlineFields)

	// Prepare a list request with necessary parameters
	listReq := &_go.ListCloseReasonGroupsRequest{
		Id:     []int64{req.Id},
		Fields: fields,
		Page:   1,
		Size:   1,
	}

	// Call the ListCloseReasonGroups method
	res, err := s.ListCloseReasonGroups(ctx, listReq)
	if err != nil {
		return nil, cerror.NewInternalError("close_reason_group_service.locate_close_reason_group.list_close_reason_groups.error", err.Error())
	}

	// Check if the close reason group was found
	if len(res.Items) == 0 {
		return nil, cerror.NewNotFoundError("close_reason_group_service.locate_close_reason_group.not_found", "Close reason group not found")
	}

	// Return the found close reason group
	return &_go.LocateCloseReasonGroupResponse{CloseReasonGroup: res.Items[0]}, nil
}

func NewCloseReasonGroupsService(app *App) (*CloseReasonGroupService, cerror.AppError) {
	if app == nil {
		return nil, cerror.NewInternalError("api.config.new_close_reason_group_service.args_check.app_nil", "internal is nil")
	}

	return &CloseReasonGroupService{app: app, objClassName: "dictionaries"}, nil
}
