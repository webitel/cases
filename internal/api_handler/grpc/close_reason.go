package grpc

import (
	"context"
	"github.com/webitel/cases/internal/errors"
	"google.golang.org/grpc/codes"

	"github.com/webitel/cases/api/cases"
	"github.com/webitel/cases/internal/api_handler/grpc/utils"
	"github.com/webitel/cases/internal/model"
	"github.com/webitel/cases/internal/model/options"
	grpcopts "github.com/webitel/cases/internal/model/options/grpc"
	"github.com/webitel/cases/util"
)

type CloseReasonHandler interface {
	ListCloseReasons(options.Searcher, int64) ([]*model.CloseReason, error)
	CreateCloseReason(options.Creator, *model.CloseReason) (*model.CloseReason, error)
	UpdateCloseReason(options.Updator, *model.CloseReason) (*model.CloseReason, error)
	DeleteCloseReason(options.Deleter) (*model.CloseReason, error)
}

type CloseReasonService struct {
	app CloseReasonHandler
	cases.UnimplementedCloseReasonsServer
}

func NewCloseReasonService(app CloseReasonHandler) (*CloseReasonService, error) {
	if app == nil {
		return nil, deferror.New("close reason handler is nil")
	}
	return &CloseReasonService{app: app}, nil
}

var CloseReasonMetadata = model.NewObjectMetadata(model.ScopeDictionary, "", []*model.Field{
	{Name: "id", Default: true},
	{Name: "created_by", Default: true},
	{Name: "created_at", Default: true},
	{Name: "updated_by", Default: false},
	{Name: "updated_at", Default: false},
	{Name: "name", Default: true},
	{Name: "description", Default: true},
	{Name: "close_reason_id", Default: false},
})

func (s *CloseReasonService) CreateCloseReason(
	ctx context.Context,
	req *cases.CreateCloseReasonRequest,
) (*cases.CloseReason, error) {
	createOpts, err := grpcopts.NewCreateOptions(
		ctx,
		grpcopts.WithCreateFields(req, CloseReasonMetadata),
	)
	if err != nil {
		return nil, err
	}

	input := &model.CloseReason{
		Name:               req.Input.Name,
		Description:        &req.Input.Description,
		CloseReasonGroupId: req.CloseReasonGroupId,
	}

	m, err := s.app.CreateCloseReason(createOpts, input)
	if err != nil {
		return nil, err
	}
	return s.Marshal(m)
}

func (s *CloseReasonService) ListCloseReasons(
	ctx context.Context,
	req *cases.ListCloseReasonRequest,
) (*cases.CloseReasonList, error) {
	searcher, err := grpcopts.NewSearchOptions(
		ctx,
		grpcopts.WithSearch(req),
		grpcopts.WithPagination(req),
		grpcopts.WithFields(req, CloseReasonMetadata,
			util.DeduplicateFields,
			util.EnsureIdField,
		),
		grpcopts.WithSort(req),
		grpcopts.WithIDs(req.GetId()),
	)
	if err != nil {
		return nil, err
	}
	searcher.AddFilter("name", req.Q)
	searcher.AddFilter("parent_id", req.CloseReasonGroupId)

	items, err := s.app.ListCloseReasons(searcher, req.GetCloseReasonGroupId())
	if err != nil {
		return nil, err
	}

	var res cases.CloseReasonList
	converted, err := utils.ConvertToOutputBulk(items, s.Marshal)
	if err != nil {
		return nil, err
	}
	res.Next, res.Items = utils.GetListResult(searcher, converted)
	res.Page = req.GetPage()

	return &res, nil
}

func (s *CloseReasonService) UpdateCloseReason(
	ctx context.Context,
	req *cases.UpdateCloseReasonRequest,
) (*cases.CloseReason, error) {
	updator, err := grpcopts.NewUpdateOptions(
		ctx,
		grpcopts.WithUpdateFields(req, CloseReasonMetadata),
		grpcopts.WithUpdateMasker(req),
	)
	if err != nil {
		return nil, err
	}

	input := &model.CloseReason{
		Id:                 int64(req.Id),
		Name:               req.Input.Name,
		Description:        &req.Input.Description,
		CloseReasonGroupId: req.CloseReasonGroupId,
	}

	updated, err := s.app.UpdateCloseReason(updator, input)
	if err != nil {
		return nil, err
	}
	return s.Marshal(updated)
}

func (s *CloseReasonService) DeleteCloseReason(
	ctx context.Context,
	req *cases.DeleteCloseReasonRequest,
) (*cases.CloseReason, error) {
	deleteOpts, err := grpcopts.NewDeleteOptions(
		ctx,
		grpcopts.WithDeleteID(req.Id),
	)
	if err != nil {
		return nil, err
	}

	item, err := s.app.DeleteCloseReason(deleteOpts)
	if err != nil {
		return nil, err
	}
	return s.Marshal(item)
}

func (s *CloseReasonService) LocateCloseReason(
	ctx context.Context,
	req *cases.LocateCloseReasonRequest,
) (*cases.LocateCloseReasonResponse, error) {
	opts, err := grpcopts.NewLocateOptions(ctx, grpcopts.WithFields(req, CloseReasonMetadata,
		util.DeduplicateFields,
		util.EnsureIdField,
	), grpcopts.WithID(req.Id))
	if err != nil {
		return nil, err
	}

	items, err := s.app.ListCloseReasons(opts, req.GetCloseReasonGroupId())
	if err != nil {
		return nil, err
	}

	if len(items) == 0 {
		return nil, errors.New("no items found", errors.WithCode(codes.NotFound))
	}
	if len(items) > 1 {
		return nil, errors.New("too many items found", errors.WithCode(codes.InvalidArgument))
	}

	res, err := s.Marshal(items[0])
	if err != nil {
		return nil, err
	}
	return &cases.LocateCloseReasonResponse{CloseReason: res}, nil
}

func (s *CloseReasonService) Marshal(model *model.CloseReason) (*cases.CloseReason, error) {
	if model == nil {
		return nil, nil
	}
	return &cases.CloseReason{
		Id:                 model.Id,
		Name:               model.Name,
		Description:        utils.Dereference(model.Description),
		CloseReasonGroupId: model.CloseReasonGroupId,
		CreatedAt:          utils.MarshalTime(model.CreatedAt),
		UpdatedAt:          utils.MarshalTime(model.UpdatedAt),
		CreatedBy:          utils.MarshalLookup(model.Author),
		UpdatedBy:          utils.MarshalLookup(model.Editor),
	}, nil
}
