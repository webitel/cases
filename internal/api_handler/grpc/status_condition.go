package grpc

import (
	"context"
	_go "github.com/webitel/cases/api/cases"
	"github.com/webitel/cases/internal/api_handler/grpc/utils"
	"github.com/webitel/cases/internal/errors"
	"github.com/webitel/cases/internal/model"
	"github.com/webitel/cases/internal/model/options"
	grpcopts "github.com/webitel/cases/internal/model/options/grpc"
	"github.com/webitel/cases/util"
	"google.golang.org/grpc/codes"
)

// StatusConditionHandler defines the interface for managing status conditions.
type StatusConditionHandler interface {
	ListStatusConditions(options.Searcher) ([]*model.StatusCondition, error)
	LocateStatusCondition(options.Searcher) (*model.StatusCondition, error)
	CreateStatusCondition(options.Creator, *model.StatusCondition) (*model.StatusCondition, error)
	UpdateStatusCondition(options.Updator, *model.StatusCondition) (*model.StatusCondition, error)
	DeleteStatusCondition(options.Deleter) (*model.StatusCondition, error)
}

// StatusConditionService implements the gRPC server for status conditions.
type StatusConditionService struct {
	app StatusConditionHandler
	_go.UnimplementedStatusConditionsServer
	objClassName string
}

// StatusConditionMetadata defines the fields available for status condition objects.
var StatusConditionMetadata = model.NewObjectMetadata(model.ScopeDictionary, "", []*model.Field{
	{Name: "id", Default: true},
	{Name: "name", Default: true},
	{Name: "description", Default: true},
	{Name: "initial", Default: true},
	{Name: "final", Default: true},
	{Name: "created_by", Default: true},
	{Name: "created_at", Default: true},
	{Name: "updated_by", Default: false},
	{Name: "updated_at", Default: false},
})

// CreateStatusCondition handles the gRPC request to create a new status condition.
func (s *StatusConditionService) CreateStatusCondition(ctx context.Context, req *_go.CreateStatusConditionRequest) (*_go.StatusCondition, error) {
	// Define create options
	createOpts, err := grpcopts.NewCreateOptions(
		ctx,
		grpcopts.WithCreateFields(req, StatusConditionMetadata),
	)
	if err != nil {
		return nil, err
	}

	statusId := int(req.StatusId)
	// Create a new status user_session
	status := &model.StatusCondition{
		Name:        &req.Input.Name,
		Description: &req.Input.Description,
		StatusId:    &statusId,
	}

	// Create the status in the store
	st, e := s.app.CreateStatusCondition(createOpts, status)
	if e != nil {
		return nil, e
	}

	return s.Marshal(st)
}

// ListStatusConditions handles the gRPC request to list status conditions with filters and pagination.
func (s *StatusConditionService) ListStatusConditions(ctx context.Context, req *_go.ListStatusConditionRequest) (*_go.StatusConditionList, error) {
	searchOptions, err := grpcopts.NewSearchOptions(
		ctx,
		grpcopts.WithPagination(req),
		grpcopts.WithFields(req, StatusConditionMetadata,
			util.DeduplicateFields,
			util.EnsureIdField,
		),
		grpcopts.WithIDs(req.Id),
		grpcopts.WithSort(req),
	)
	if err != nil {
		return nil, err
	}
	searchOptions.AddFilter("parent_id", req.StatusId)
	if req.Q != "" {
		searchOptions.AddFilter("name", req.Q)
	}

	statuses, err := s.app.ListStatusConditions(searchOptions)
	if err != nil {
		return nil, err
	}
	var res _go.StatusConditionList
	res.Items, err = utils.ConvertToOutputBulk(statuses, s.Marshal)
	if err != nil {
		return nil, err
	}
	res.Next, res.Items = utils.GetListResult(searchOptions, res.Items)
	res.Page = req.GetPage()
	return &res, nil
}

// UpdateStatusCondition handles the gRPC request to update an existing status condition.
func (s *StatusConditionService) UpdateStatusCondition(ctx context.Context, req *_go.UpdateStatusConditionRequest) (*_go.StatusCondition, error) {
	// Define update options
	updateOpts, err := grpcopts.NewUpdateOptions(
		ctx,
		grpcopts.WithUpdateFields(req, StatusConditionMetadata),
		grpcopts.WithUpdateMasker(req),
	)
	if err != nil {
		return nil, err
	}
	statusId := int(req.StatusId)
	// Update input user_session
	input := &model.StatusCondition{
		Id:          int(req.Id),
		StatusId:    &statusId,
		Name:        &req.Input.Name,
		Description: &req.Input.Description,
	}

	if req.Input.Initial != nil {
		input.Initial = &req.Input.Initial.Value
	}
	if req.Input.Final != nil {
		input.Final = &req.Input.Final.Value
	}

	// Update the input in the store
	st, err := s.app.UpdateStatusCondition(updateOpts, input)
	if err != nil {
		return nil, err
	}

	return s.Marshal(st)
}

// DeleteStatusCondition handles the gRPC request to delete a status condition.
func (s *StatusConditionService) DeleteStatusCondition(ctx context.Context, req *_go.DeleteStatusConditionRequest) (*_go.StatusCondition, error) {
	deleteOpts, err := grpcopts.NewDeleteOptions(ctx, grpcopts.WithDeleteID(req.Id), grpcopts.WithDeleteParentID(req.StatusId))
	if err != nil {
		return nil, err
	}

	// Delete the status in the store
	_, err = s.app.DeleteStatusCondition(deleteOpts)
	if err != nil {
		return nil, err
	}

	return &(_go.StatusCondition{Id: req.Id}), nil
}

// LocateStatusCondition finds a status condition by ID and returns it, or an error if not found or ambiguous.
func (s *StatusConditionService) LocateStatusCondition(ctx context.Context, req *_go.LocateStatusConditionRequest) (*_go.LocateStatusConditionResponse, error) {
	opts, err := grpcopts.NewLocateOptions(ctx, grpcopts.WithID(req.Id), grpcopts.WithFields(req, StatusConditionMetadata, util.EnsureIdField, util.DeduplicateFields))
	if err != nil {
		return nil, err
	}
	// Call the ListStatusConditions method
	items, err := s.app.ListStatusConditions(opts)
	if err != nil {
		return nil, err
	}

	// Check if the status condition was found
	if len(items) > 1 {
		return nil, errors.New("too many records found", errors.WithCode(codes.InvalidArgument))
	}
	if len(items) == 0 {
		return nil, errors.New("not found", errors.WithCode(codes.NotFound))
	}

	// Return the found status condition
	var res _go.LocateStatusConditionResponse
	res.Status, err = s.Marshal(items[0])
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// Marshal converts a model.StatusCondition to its gRPC representation.
func (s *StatusConditionService) Marshal(model *model.StatusCondition) (*_go.StatusCondition, error) {
	return &_go.StatusCondition{
		Id:          int64(model.Id),
		Name:        utils.Dereference(model.Name),
		Description: utils.Dereference(model.Description),
		Initial:     utils.Dereference(model.Initial),
		Final:       utils.Dereference(model.Final),
		StatusId:    int64(utils.Dereference(model.StatusId)),
		CreatedAt:   utils.MarshalTime(model.CreatedAt),
		UpdatedAt:   utils.MarshalTime(model.UpdatedAt),
		CreatedBy:   utils.MarshalLookup(model.Author),
		UpdatedBy:   utils.MarshalLookup(model.Editor),
	}, nil
}

// NewStatusConditionService constructs a new StatusConditionService.
func NewStatusConditionService(app StatusConditionHandler) (*StatusConditionService, error) {
	if app == nil {
		return nil, errors.New("status condition handler is nil")
	}
	return &StatusConditionService{app: app, objClassName: model.ScopeDictionary}, nil
}
