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

// CloseReasonGroupHandler defines the interface for managing close reason groups.
type CloseReasonGroupHandler interface {
	ListCloseReasonGroup(options.Searcher) ([]*model.CloseReasonGroup, error)
	CreateCloseReasonGroup(options.Creator, *model.CloseReasonGroup) (*model.CloseReasonGroup, error)
	UpdateCloseReasonGroup(options.Updator, *model.CloseReasonGroup) (*model.CloseReasonGroup, error)
	DeleteCloseReasonGroup(options.Deleter) (*model.CloseReasonGroup, error)
}

// CloseReasonGroupService implements the gRPC server for close reason groups.
type CloseReasonGroupService struct {
	app CloseReasonGroupHandler
	_go.UnimplementedCloseReasonGroupsServer
	objClassName string
}

// CloseReasonGroupMetadata defines the fields available for close reason group objects.
var CloseReasonGroupMetadata = model.NewObjectMetadata(model.ScopeDictionary, "", []*model.Field{
	{Name: "id", Default: true},
	{Name: "created_by", Default: true},
	{Name: "created_at", Default: true},
	{Name: "updated_by", Default: false},
	{Name: "updated_at", Default: false},
	{Name: "name", Default: true},
	{Name: "description", Default: true},
})

// CreateCloseReasonGroup handles the gRPC request to create a new close reason group.
func (s *CloseReasonGroupService) CreateCloseReasonGroup(
	ctx context.Context,
	req *_go.CreateCloseReasonGroupRequest,
) (*_go.CloseReasonGroup, error) {
	createOpts, err := grpcopts.NewCreateOptions(
		ctx,
		grpcopts.WithCreateFields(req, CloseReasonGroupMetadata),
	)
	if err != nil {
		return nil, err
	}

	input := &model.CloseReasonGroup{
		Name:        &req.Input.Name,
		Description: &req.Input.Description,
	}

	// Create the close reason group in the store
	m, err := s.app.CreateCloseReasonGroup(createOpts, input)
	if err != nil {
		return nil, err
	}

	return s.Marshal(m)
}

// ListCloseReasonGroups handles the gRPC request to list close reason groups with filters and pagination.
func (s *CloseReasonGroupService) ListCloseReasonGroups(
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
		return nil, err
	}
	searchOpts.AddFilter(util.EqualFilter("name", req.Q))

	items, err := s.app.ListCloseReasonGroup(searchOpts)
	if err != nil {
		return nil, err
	}
	var res _go.CloseReasonGroupList
	converted, err := utils.ConvertToOutputBulk(items, s.Marshal)
	if err != nil {
		return nil, err
	}
	res.Next, res.Items = utils.GetListResult(searchOpts, converted)
	res.Page = req.GetPage()

	return &res, nil
}

// UpdateCloseReasonGroup handles the gRPC request to update an existing close reason group.
func (s *CloseReasonGroupService) UpdateCloseReasonGroup(
	ctx context.Context,
	req *_go.UpdateCloseReasonGroupRequest,
) (*_go.CloseReasonGroup, error) {
	updateOpts, err := grpcopts.NewUpdateOptions(
		ctx,
		grpcopts.WithUpdateFields(req, CloseReasonGroupMetadata),
		grpcopts.WithUpdateMasker(req),
		grpcopts.WithUpdateIDs([]int64{req.GetId()}),
	)
	if err != nil {
		return nil, err
	}

	// Update lookup user_session
	input := &model.CloseReasonGroup{
		Id:          int(req.Id),
		Name:        &req.Input.Name,
		Description: &req.Input.Description,
	}

	// Update the lookup in the store
	item, err := s.app.UpdateCloseReasonGroup(updateOpts, input)
	if err != nil {
		return nil, err
	}

	return s.Marshal(item)
}

// DeleteCloseReasonGroup handles the gRPC request to delete a close reason group.
func (s *CloseReasonGroupService) DeleteCloseReasonGroup(
	ctx context.Context,
	req *_go.DeleteCloseReasonGroupRequest,
) (*_go.CloseReasonGroup, error) {
	deleteOpts, err := grpcopts.NewDeleteOptions(ctx, grpcopts.WithDeleteID(req.Id))
	if err != nil {
		return nil, err
	}

	// Delete the lookup in the store
	item, err := s.app.DeleteCloseReasonGroup(deleteOpts)
	if err != nil {
		return nil, err
	}

	return s.Marshal(item)
}

// LocateCloseReasonGroup finds a close reason group by ID and returns it, or an error if not found or ambiguous.
func (s *CloseReasonGroupService) LocateCloseReasonGroup(
	ctx context.Context,
	req *_go.LocateCloseReasonGroupRequest,
) (*_go.LocateCloseReasonGroupResponse, error) {
	opts, err := grpcopts.NewLocateOptions(ctx, grpcopts.WithFields(req, CloseReasonGroupMetadata,
		util.DeduplicateFields,
		util.EnsureIdField,
	), grpcopts.WithID(req.Id))
	if err != nil {
		return nil, err
	}
	// Call the ListCloseReasonGroups method
	items, err := s.app.ListCloseReasonGroup(opts)
	if err != nil {
		return nil, err
	}

	if len(items) == 0 {
		return nil, errors.New("no records found", errors.WithCode(codes.NotFound))
	}
	if len(items) > 1 {
		return nil, errors.New("too many records found", errors.WithCode(codes.InvalidArgument))
	}

	// Return the found close reason group
	res, err := s.Marshal(items[0])
	if err != nil {
		return nil, err
	}
	return &_go.LocateCloseReasonGroupResponse{CloseReasonGroup: res}, nil
}

// NewCloseReasonGroupsService constructs a new CloseReasonGroupService.
func NewCloseReasonGroupsService(app CloseReasonGroupHandler) (*CloseReasonGroupService, error) {
	if app == nil {
		return nil, errors.New("close reason handler is required")
	}

	return &CloseReasonGroupService{app: app, objClassName: "dictionaries"}, nil
}

// Marshal converts a model.CloseReasonGroup to its gRPC representation.
func (s *CloseReasonGroupService) Marshal(model *model.CloseReasonGroup) (*_go.CloseReasonGroup, error) {
	return &_go.CloseReasonGroup{
		Id:          int64(model.Id),
		Name:        utils.Dereference(model.Name),
		Description: utils.Dereference(model.Description),
		CreatedAt:   utils.MarshalTime(model.CreatedAt),
		UpdatedAt:   utils.MarshalTime(model.UpdatedAt),
		CreatedBy:   utils.MarshalLookup(model.Author),
		UpdatedBy:   utils.MarshalLookup(model.Editor),
	}, nil
}
