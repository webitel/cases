package grpc

import (
	"context"
	"errors"
	_go "github.com/webitel/cases/api/cases"
	grpcerror "github.com/webitel/cases/internal/api_handler/grpc/errors"
	"github.com/webitel/cases/internal/api_handler/grpc/utils"
	cerror "github.com/webitel/cases/internal/errors"
	"github.com/webitel/cases/internal/model"
	"github.com/webitel/cases/internal/model/options"
	grpcopts "github.com/webitel/cases/internal/model/options/grpc"
	"github.com/webitel/cases/util"
	"log/slog"
)

type StatusConditionHandler interface {
	ListStatusConditions(options.Searcher) ([]*model.StatusCondition, error)
	LocateStatusCondition(options.Searcher) (*model.StatusCondition, error)
	CreateStatusCondition(options.Creator, *model.StatusCondition) (*model.StatusCondition, error)
	UpdateStatusCondition(options.Updator, *model.StatusCondition) (*model.StatusCondition, error)
	DeleteStatusCondition(options.Deleter) (*model.StatusCondition, error)
}

type StatusConditionService struct {
	app StatusConditionHandler
	_go.UnimplementedStatusConditionsServer
	objClassName string
}

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

// CreateStatusCondition implements api.StatusConditionsServer.
func (s *StatusConditionService) CreateStatusCondition(ctx context.Context, req *_go.CreateStatusConditionRequest) (*_go.StatusCondition, error) {
	// Validate required fields
	if req.Input.Name == "" {
		return nil, grpcerror.NewBadRequestError(errors.New("status name is required"))
	}

	// Define create options
	createOpts, err := grpcopts.NewCreateOptions(
		ctx,
		grpcopts.WithCreateFields(req, StatusConditionMetadata),
	)
	if err != nil {
		return nil, grpcerror.NewBadRequestError(err)
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

// ListStatusConditions implements api.StatusConditionsServer.
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
		return nil, grpcerror.NewBadRequestError(err)
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
		slog.ErrorContext(ctx, err.Error())
		return nil, grpcerror.ConversionError
	}
	res.Next, res.Items = utils.GetListResult(searchOptions, res.Items)
	res.Page = req.GetPage()
	return &res, nil
}

// UpdateStatusCondition implements api.StatusConditionsServer.
func (s *StatusConditionService) UpdateStatusCondition(ctx context.Context, req *_go.UpdateStatusConditionRequest) (*_go.StatusCondition, error) {
	// Validate required fields
	if req.Id == 0 {
		return nil, grpcerror.NewBadRequestError(errors.New("status id is required"))
	}

	// Define update options
	updateOpts, err := grpcopts.NewUpdateOptions(
		ctx,
		grpcopts.WithUpdateFields(req, StatusConditionMetadata),
		grpcopts.WithUpdateMasker(req),
	)
	if err != nil {
		return nil, grpcerror.NewBadRequestError(err)
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
		switch err.(type) {
		case *cerror.DBCheckViolationError:
			return nil, cerror.NewBadRequestError(
				"app.status_condition.update.initial_false_not_allowed",
				"update not allowed: there must be at least one initial = TRUE for the given dc and status_id",
			)
		case *cerror.DBInternalError:
			return nil, cerror.NewBadRequestError(
				"app.status_condition.update.error",
				err.Error(),
			)
		}
		return nil, cerror.NewInternalError(
			"app.status_condition.update.error",
			err.Error(),
		)
	}

	return s.Marshal(st)
}

// DeleteStatusCondition implements api.StatusConditionsServer.
func (s *StatusConditionService) DeleteStatusCondition(ctx context.Context, req *_go.DeleteStatusConditionRequest) (*_go.StatusCondition, error) {
	// Validate required fields
	if req.Id == 0 {
		return nil, grpcerror.NewBadRequestError(errors.New("status ID is required"))
	}

	deleteOpts, err := grpcopts.NewDeleteOptions(ctx, grpcopts.WithDeleteID(req.Id), grpcopts.WithDeleteParentID(req.StatusId))
	if err != nil {
		return nil, grpcerror.NewBadRequestError(err)
	}

	// Delete the status in the store
	_, err = s.app.DeleteStatusCondition(deleteOpts)
	if err != nil {
		switch err.(type) {
		case *cerror.DBNoRowsError:
			return nil, cerror.NewBadRequestError(
				"status_condition.delete_status_condition.not_found",
				"delete not allowed",
			)
		}
		return nil, cerror.NewInternalError(
			"status_condition.delete_status_condition.error",
			err.Error(),
		)
	}

	return &(_go.StatusCondition{Id: req.Id}), nil
}

// LocateStatusCondition implements api.StatusConditionsServer.
func (s *StatusConditionService) LocateStatusCondition(ctx context.Context, req *_go.LocateStatusConditionRequest) (*_go.LocateStatusConditionResponse, error) {
	// Validate required fields
	if req.Id == 0 {
		return nil, grpcerror.NewBadRequestError(errors.New("status ID is required"))
	}
	opts, err := grpcopts.NewLocateOptions(ctx, grpcopts.WithID(req.Id), grpcopts.WithFields(req, StatusConditionMetadata, util.EnsureIdField, util.DeduplicateFields))
	if err != nil {
		return nil, grpcerror.NewBadRequestError(err)
	}
	// Call the ListStatusConditions method
	items, err := s.app.ListStatusConditions(opts)
	if err != nil {
		return nil, err
	}

	// Check if the status condition was found
	if len(items) > 1 {
		return nil, grpcerror.NewBadRequestError(errors.New("multiple rows found"))
	}
	if len(items) == 0 {
		return nil, grpcerror.NewBadRequestError(errors.New("not found"))
	}

	// Return the found status condition
	var res _go.LocateStatusConditionResponse
	res.Status, err = s.Marshal(items[0])
	if err != nil {
		return nil, grpcerror.ConversionError
	}
	return &res, nil
}

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

func NewStatusConditionService(app StatusConditionHandler) (*StatusConditionService, error) {
	if app == nil {
		return nil, cerror.NewInternalError("api.config.new_status_condition_service.args_check.app_nil", "internal is nil")
	}
	return &StatusConditionService{app: app, objClassName: model.ScopeDictionary}, nil
}
