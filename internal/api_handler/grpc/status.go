package grpc

import (
	"context"
	deferror "errors"
	_go "github.com/webitel/cases/api/cases"
	grpcerror "github.com/webitel/cases/internal/api_handler/grpc/errors"
	"github.com/webitel/cases/internal/api_handler/grpc/utils"
	cerror "github.com/webitel/cases/internal/errors"
	"github.com/webitel/cases/internal/model"
	"github.com/webitel/cases/internal/model/options"
	grpcopts "github.com/webitel/cases/internal/model/options/grpc"
	"github.com/webitel/cases/util"
	"log/slog"
	"strings"
)

type StatusHandler interface {
	ListStatus(options.Searcher) ([]*model.Status, error)
	LocateStatus(options.Searcher) (*model.Status, error)
	CreateStatus(options.Creator, *model.Status) (*model.Status, error)
	UpdateStatus(options.Updator, *model.Status) (*model.Status, error)
	DeleteStatus(options.Deleter) (*model.Status, error)
}

type StatusService struct {
	app StatusHandler
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
func (s *StatusService) CreateStatus(ctx context.Context, req *_go.CreateStatusRequest) (*_go.Status, error) {
	// Validate required fields
	if req.Input.Name == "" {
		return nil, cerror.NewBadRequestError("status.create_status.input.name.required", ErrLookupNameReq)
	}

	createOpts, err := grpcopts.NewCreateOptions(
		ctx,
		grpcopts.WithCreateFields(req, StatusMetadata),
	)

	if err != nil {
		return nil, grpcerror.NewBadRequestError(err)
	}

	// Create a new input user_session
	input := &model.Status{
		Name:        &req.Input.Name,
		Description: &req.Input.Description,
	}

	res, err := s.app.CreateStatus(createOpts, input)
	if err != nil {
		return nil, err
	}

	return s.Marshal(res)
}

// ListStatuses implements api.StatusesServer.
func (s *StatusService) ListStatuses(ctx context.Context, req *_go.ListStatusRequest) (*_go.StatusList, error) {
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
		return nil, grpcerror.NewBadRequestError(err)
	}
	searchOpts.AddFilter("name", req.Q)

	items, err := s.app.ListStatus(searchOpts)
	if err != nil {
		return nil, cerror.NewInternalError("status.list_status.store.list.failed", err.Error())
	}
	var res _go.StatusList
	res.Items, err = utils.ConvertToOutputBulk(items, s.Marshal)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return nil, grpcerror.ConversionError
	}
	res.Next, res.Items = utils.GetListResult(searchOpts, res.Items)
	res.Page = req.GetPage()
	return &res, nil
}

// UpdateStatus implements api.StatusesServer.
func (s *StatusService) UpdateStatus(ctx context.Context, req *_go.UpdateStatusRequest) (*_go.Status, error) {
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
		return nil, grpcerror.NewBadRequestError(err)
	}

	// Update input user_session
	input := &model.Status{
		Id:          int(req.Id),
		Name:        &req.Input.Name,
		Description: &req.Input.Description,
	}

	// Update the input in the store
	res, err := s.app.UpdateStatus(updateOpts, input)
	if err != nil {
		return nil, err
	}

	return s.Marshal(res)
}

// DeleteStatus implements api.StatusesServer.
func (s *StatusService) DeleteStatus(ctx context.Context, req *_go.DeleteStatusRequest) (*_go.Status, error) {
	// Validate required fields
	if req.Id == 0 {
		return nil, grpcerror.NewBadRequestError(deferror.New("lookup ID is required"))
	}

	deleteOpts, err := grpcopts.NewDeleteOptions(ctx, grpcopts.WithDeleteID(req.Id))
	if err != nil {
		return nil, grpcerror.NewBadRequestError(err)
	}

	// Delete the lookup in the store
	item, err := s.app.DeleteStatus(deleteOpts)
	if err != nil {
		return nil, cerror.NewInternalError("status.delete_status.store.delete.failed", err.Error())
	}

	return s.Marshal(item)
}

// LocateStatus implements api.StatusesServer.
func (s *StatusService) LocateStatus(ctx context.Context, req *_go.LocateStatusRequest) (*_go.LocateStatusResponse, error) {
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

func NewStatusService(app StatusHandler) (*StatusService, error) {
	if app == nil {
		return nil, deferror.New("status handler is nil")
	}
	return &StatusService{app: app}, nil
}

func (s *StatusService) Marshal(input *model.Status) (*_go.Status, error) {
	if input == nil {
		return nil, nil
	}
	return &_go.Status{
		Id:          int64(input.Id),
		Name:        utils.Dereference(input.Name),
		Description: utils.Dereference(input.Description),
		CreatedAt:   utils.MarshalTime(input.CreatedAt),
		UpdatedAt:   utils.MarshalTime(input.UpdatedAt),
		CreatedBy:   utils.MarshalLookup(input.Author),
		UpdatedBy:   utils.MarshalLookup(input.Editor),
	}, nil
}
