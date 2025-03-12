package app

import (
	"context"
	_go "github.com/webitel/cases/api/cases"
	cerror "github.com/webitel/cases/internal/errors"
	"github.com/webitel/cases/model"
	grpcopts "github.com/webitel/cases/model/options/grpc"
	"github.com/webitel/cases/util"
	"strings"
)

type StatusService struct {
	app *App
	_go.UnimplementedStatusesServer
	objClassName string
}

var StatusMetadata = model.NewObjectMetadata(model.ScopeDictionary, "", []*model.Field{
	{"id", true},
	{"created_by", true},
	{"created_at", true},
	{"updated_by", false},
	{"updated_at", false},
	{"name", true},
	{"description", true},
})

const (
	ErrLookupNameReq    = "Lookup name is required"
	statusDefaultFields = "id, name, description, created_by"
)

// CreateStatus implements api.StatusesServer.
func (s StatusService) CreateStatus(ctx context.Context, req *_go.CreateStatusRequest) (*_go.Status, error) {
	// Validate required fields
	if req.Input.Name == "" {
		return nil, cerror.NewBadRequestError("status.create_status.input.name.required", ErrLookupNameReq)
	}

	createOpts, err := grpcopts.NewCreateOptions(
		ctx,
		grpcopts.WithCreateFields(req, StatusMetadata),
	)

	if err != nil {
		return nil, NewBadRequestError(err)
	}

	// Create a new input user_auth
	input := &_go.Status{
		Name:        req.Input.Name,
		Description: req.Input.Description,
	}

	res, err := s.app.Store.Status().Create(createOpts, input)
	if err != nil {
		return nil, cerror.NewInternalError("status.create_status.store.create.failed", err.Error())
	}

	return res, nil
}

// ListStatuses implements api.StatusesServer.
func (s StatusService) ListStatuses(ctx context.Context, req *_go.ListStatusRequest) (*_go.StatusList, error) {
	searchOpts, err := grpcopts.NewSearchOptions(
		ctx,
		grpcopts.WithSearch(req),
		grpcopts.WithPagination(req),
		grpcopts.WithFields(req, StatusMetadata,
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

	res, err := s.app.Store.Status().List(searchOpts)
	if err != nil {
		return nil, cerror.NewInternalError("status.list_status.store.list.failed", err.Error())
	}

	return res, nil
}

// UpdateStatus implements api.StatusesServer.
func (s StatusService) UpdateStatus(ctx context.Context, req *_go.UpdateStatusRequest) (*_go.Status, error) {
	// Validate required fields
	if req.Id == 0 {
		return nil, cerror.NewBadRequestError("status.update_status.input.id.required", "Lookup ID is required")
	}

	updateOpts, err := grpcopts.NewUpdateOptions(
		ctx,
		grpcopts.WithUpdateFields(req, StatusMetadata),
		grpcopts.WithUpdateMasker(req),
	)
	if err != nil {
		return nil, NewBadRequestError(err)
	}

	// Update input user_auth
	input := &_go.Status{
		Id:          req.Id,
		Name:        req.Input.Name,
		Description: req.Input.Description,
	}

	// Update the input in the store
	res, err := s.app.Store.Status().Update(updateOpts, input)
	if err != nil {
		return nil, cerror.NewInternalError("status.update_status.store.update.failed", err.Error())
	}

	return res, nil
}

// DeleteStatus implements api.StatusesServer.
func (s StatusService) DeleteStatus(ctx context.Context, req *_go.DeleteStatusRequest) (*_go.Status, error) {
	// Validate required fields
	if req.Id == 0 {
		return nil, cerror.NewBadRequestError("status.delete_status.lookup.id.required", "Lookup ID is required")
	}

	deleteOpts, err := grpcopts.NewDeleteOptions(ctx, grpcopts.WithDeleteID(req.Id))
	if err != nil {
		return nil, NewBadRequestError(err)
	}

	// Delete the lookup in the store
	err = s.app.Store.Status().Delete(deleteOpts)
	if err != nil {
		return nil, cerror.NewInternalError("status.delete_status.store.delete.failed", err.Error())
	}

	return &(_go.Status{Id: req.Id}), nil
}

// LocateStatus implements api.StatusesServer.
func (s StatusService) LocateStatus(ctx context.Context, req *_go.LocateStatusRequest) (*_go.LocateStatusResponse, error) {
	// Validate required fields
	if req.Id == 0 {
		return nil, cerror.NewBadRequestError("status.locate_status.lookup.id.required", "Lookup ID is required")
	}

	fields := util.FieldsFunc(req.Fields, util.InlineFields)
	if len(fields) == 0 {
		fields = strings.Split(statusDefaultFields, ",")
	}

	// Prepare a list request with necessary parameters
	listReq := &_go.ListStatusRequest{
		Id:     []int64{req.Id},
		Fields: fields,
		Page:   1,
		Size:   1,
	}

	// Call the ListStatuses method
	res, err := s.ListStatuses(ctx, listReq)
	if err != nil {
		return nil, cerror.NewInternalError("status.locate_status.list_status.error", err.Error())
	}

	// Check if the lookup was found
	if len(res.Items) == 0 {
		return nil, cerror.NewNotFoundError("status.locate_status.not_found", "Status lookup not found")
	}

	// Return the found status lookup
	return &_go.LocateStatusResponse{Status: res.Items[0]}, nil
}

func NewStatusService(app *App) (*StatusService, cerror.AppError) {
	if app == nil {
		return nil, cerror.NewInternalError("api.config.new_status.args_check.app_nil", "internal is nil")
	}
	return &StatusService{app: app}, nil
}
