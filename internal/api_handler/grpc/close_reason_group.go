package grpc

import (
	"context"
	deferror "errors"
	_go "github.com/webitel/cases/api/cases"
	grpcerror "github.com/webitel/cases/internal/api_handler/grpc/errors"
	"github.com/webitel/cases/internal/api_handler/grpc/utils"
	"github.com/webitel/cases/internal/model"
	"github.com/webitel/cases/internal/model/options"
	grpcopts "github.com/webitel/cases/internal/model/options/grpc"
	"github.com/webitel/cases/util"
)

type CloseReasonGroupHandler interface {
	ListCloseReasonGroup(options.Searcher) ([]*model.CloseReasonGroup, error)
	LocateCloseReasonGroup(options.Searcher) (*model.CloseReasonGroup, error)
	CreateCloseReasonGroup(options.Creator, *model.CloseReasonGroup) (*model.CloseReasonGroup, error)
	UpdateCloseReasonGroup(options.Updator, *model.CloseReasonGroup) (*model.CloseReasonGroup, error)
	DeleteCloseReasonGroup(options.Deleter) (*model.CloseReasonGroup, error)
}

type CloseReasonGroupService struct {
	app CloseReasonGroupHandler
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

func (s *CloseReasonGroupService) CreateCloseReasonGroup(
	ctx context.Context,
	req *_go.CreateCloseReasonGroupRequest,
) (*_go.CloseReasonGroup, error) {
	// Validate required fields
	if req.Input.Name == "" {
		return nil, grpcerror.NewBadRequestError(deferror.New("lookup name is required"))
	}

	createOpts, err := grpcopts.NewCreateOptions(
		ctx,
		grpcopts.WithCreateFields(req, CloseReasonGroupMetadata),
	)
	if err != nil {
		return nil, grpcerror.NewBadRequestError(err)
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
		return nil, grpcerror.NewBadRequestError(err)
	}
	searchOpts.AddFilter("name", req.Q)

	items, err := s.app.ListCloseReasonGroup(searchOpts)
	if err != nil {
		return nil, err
	}
	var res _go.CloseReasonGroupList
	converted, err := utils.ConvertToOutputBulk(items, s.Marshal)
	if err != nil {
		return nil, err
	}
	res.Next, res.Items = utils.GetListResult(req, converted)
	res.Page = req.GetPage()

	return &res, nil
}

func (s *CloseReasonGroupService) UpdateCloseReasonGroup(
	ctx context.Context,
	req *_go.UpdateCloseReasonGroupRequest,
) (*_go.CloseReasonGroup, error) {
	// Validate required fields
	if req.Id == 0 {
		return nil, grpcerror.NewBadRequestError(deferror.New("lookup ID is required"))
	}

	updateOpts, err := grpcopts.NewUpdateOptions(
		ctx,
		grpcopts.WithUpdateFields(req, CloseReasonGroupMetadata),
		grpcopts.WithUpdateMasker(req),
	)
	if err != nil {
		return nil, grpcerror.NewBadRequestError(err)
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

func (s *CloseReasonGroupService) DeleteCloseReasonGroup(
	ctx context.Context,
	req *_go.DeleteCloseReasonGroupRequest,
) (*_go.CloseReasonGroup, error) {
	// Validate required fields
	if req.Id == 0 {
		return nil, grpcerror.NewBadRequestError(deferror.New("lookup ID is required"))
	}

	deleteOpts, err := grpcopts.NewDeleteOptions(ctx, grpcopts.WithDeleteID(req.Id))
	if err != nil {
		return nil, grpcerror.NewBadRequestError(err)
	}

	// Delete the lookup in the store
	item, err := s.app.DeleteCloseReasonGroup(deleteOpts)
	if err != nil {
		return nil, err
	}

	return s.Marshal(item)
}

func (s *CloseReasonGroupService) LocateCloseReasonGroup(
	ctx context.Context,
	req *_go.LocateCloseReasonGroupRequest,
) (*_go.LocateCloseReasonGroupResponse, error) {
	// Validate required fields
	if req.Id == 0 {
		return nil, grpcerror.NewBadRequestError(deferror.New("Lookup ID is required"))
	}
	opts, err := grpcopts.NewLocateOptions(ctx, grpcopts.WithFields(req, CloseReasonGroupMetadata,
		util.DeduplicateFields,
		util.EnsureIdField,
	), grpcopts.WithID(req.Id))
	if err != nil {
		return nil, grpcerror.NewBadRequestError(err)
	}
	// Call the ListCloseReasonGroups method
	item, err := s.app.LocateCloseReasonGroup(opts)
	if err != nil {
		return nil, err
	}

	// Return the found close reason group
	res, err := s.Marshal(item)
	if err != nil {

	}
	return &_go.LocateCloseReasonGroupResponse{CloseReasonGroup: res}, nil
}

func NewCloseReasonGroupsService(app CloseReasonGroupHandler) (*CloseReasonGroupService, error) {
	if app == nil {
		return nil, deferror.New("close reason handler is required")
	}

	return &CloseReasonGroupService{app: app, objClassName: "dictionaries"}, nil
}

func (s *CloseReasonGroupService) Marshal(model *model.CloseReasonGroup) (*_go.CloseReasonGroup, error) {
	return &_go.CloseReasonGroup{
		Id:          int64(model.Id),
		Name:        utils.Dereference(model.Name),
		Description: utils.Dereference(model.Description),
		CreatedAt:   model.CreatedAt.UnixMilli(),
		UpdatedAt:   model.UpdatedAt.UnixMilli(),
		CreatedBy:   utils.MarshalLookup(model.Author),
		UpdatedBy:   utils.MarshalLookup(model.Editor),
	}, nil
}

func (s *CloseReasonGroupService) Unmarshal(model *model.CloseReasonGroup) (*_go.CloseReasonGroup, error) {
	return &_go.CloseReasonGroup{
		Id:          int64(model.Id),
		Name:        utils.Dereference(model.Name),
		Description: utils.Dereference(model.Description),
		CreatedAt:   model.CreatedAt.UnixMilli(),
		UpdatedAt:   model.UpdatedAt.UnixMilli(),
		CreatedBy:   utils.MarshalLookup(model.Author),
		UpdatedBy:   utils.MarshalLookup(model.Editor),
	}, nil
}
