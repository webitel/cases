package grpc

import (
	"context"
	deferror "errors"
	_go "github.com/webitel/cases/api/cases"
	"github.com/webitel/cases/internal/api_handler/grpc/utils"
	"github.com/webitel/cases/internal/errors"
	"github.com/webitel/cases/internal/model"
	"github.com/webitel/cases/internal/model/options"
	grpcopts "github.com/webitel/cases/internal/model/options/grpc"
	"github.com/webitel/cases/util"
	"google.golang.org/grpc/codes"
)

// StatusHandler defines the interface for managing statuses.
type StatusHandler interface {
	ListStatus(options.Searcher) ([]*model.Status, error)
	CreateStatus(options.Creator, *model.Status) (*model.Status, error)
	UpdateStatus(options.Updator, *model.Status) (*model.Status, error)
	DeleteStatus(options.Deleter) (*model.Status, error)
}

// StatusService implements the gRPC server for statuses.
type StatusService struct {
	app StatusHandler
	_go.UnimplementedStatusesServer
	objClassName string
}

// StatusMetadata defines the fields available for status objects.
var StatusMetadata = model.NewObjectMetadata(model.ScopeDictionary, "", []*model.Field{
	{"id", true},
	{"created_by", true},
	{"created_at", true},
	{"updated_by", false},
	{"updated_at", false},
	{"name", true},
	{"description", true},
})

// CreateStatus handles the gRPC request to create a new status.
func (s *StatusService) CreateStatus(ctx context.Context, req *_go.CreateStatusRequest) (*_go.Status, error) {
	createOpts, err := grpcopts.NewCreateOptions(
		ctx,
		grpcopts.WithCreateFields(req, StatusMetadata),
	)

	if err != nil {
		return nil, err
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

// ListStatuses handles the gRPC request to list statuses with filters and pagination.
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
		return nil, err
	}
	searchOpts.AddFilter("name", req.Q)

	items, err := s.app.ListStatus(searchOpts)
	if err != nil {
		return nil, err
	}
	var res _go.StatusList
	res.Items, err = utils.ConvertToOutputBulk(items, s.Marshal)
	if err != nil {
		return nil, err
	}
	res.Next, res.Items = utils.GetListResult(searchOpts, res.Items)
	res.Page = req.GetPage()
	return &res, nil
}

// UpdateStatus handles the gRPC request to update an existing status.
func (s *StatusService) UpdateStatus(ctx context.Context, req *_go.UpdateStatusRequest) (*_go.Status, error) {
	updateOpts, err := grpcopts.NewUpdateOptions(
		ctx,
		grpcopts.WithUpdateFields(req, StatusMetadata),
		grpcopts.WithUpdateMasker(req),
		grpcopts.WithUpdateIDs([]int64{req.GetId()}),
	)
	if err != nil {
		return nil, err
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

// DeleteStatus handles the gRPC request to delete a status.
func (s *StatusService) DeleteStatus(ctx context.Context, req *_go.DeleteStatusRequest) (*_go.Status, error) {
	deleteOpts, err := grpcopts.NewDeleteOptions(ctx, grpcopts.WithDeleteID(req.Id))
	if err != nil {
		return nil, err
	}

	// Delete the lookup in the store
	item, err := s.app.DeleteStatus(deleteOpts)
	if err != nil {
		return nil, err
	}

	return s.Marshal(item)
}

// LocateStatus finds a status by ID and returns it, or an error if not found or ambiguous.
func (s *StatusService) LocateStatus(ctx context.Context, req *_go.LocateStatusRequest) (*_go.LocateStatusResponse, error) {
	opts, err := grpcopts.NewLocateOptions(ctx, grpcopts.WithFields(req, StatusMetadata), grpcopts.WithID(req.GetId()))
	if err != nil {
		return nil, err
	}
	// Call the ListStatuses method
	items, err := s.app.ListStatus(opts)
	if err != nil {
		return nil, err
	}

	// Check if the lookup was found
	if len(items) == 0 {
		return nil, errors.New("status lookup not found", errors.WithCode(codes.NotFound))
	}
	if len(items) > 1 {
		return nil, errors.New("status lookup not found", errors.WithCode(codes.InvalidArgument))
	}
	res, err := s.Marshal(items[0])
	if err != nil {
		return nil, err
	}
	// Return the found status lookup
	return &_go.LocateStatusResponse{Status: res}, nil
}

// NewStatusService constructs a new StatusService.
func NewStatusService(app StatusHandler) (*StatusService, error) {
	if app == nil {
		return nil, deferror.New("status handler is nil")
	}
	return &StatusService{app: app}, nil
}

// Marshal converts a model.Status to its gRPC representation.
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
