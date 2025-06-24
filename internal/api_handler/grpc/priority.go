package grpc

import (
	"context"
	api "github.com/webitel/cases/api/cases"
	"github.com/webitel/cases/internal/api_handler/grpc/utils"
	"github.com/webitel/cases/internal/errors"
	"github.com/webitel/cases/internal/model"
	"github.com/webitel/cases/internal/model/options"
	grpcopts "github.com/webitel/cases/internal/model/options/grpc"
	"github.com/webitel/cases/util"
	"google.golang.org/grpc/codes"
)

type PriorityHandler interface {
	ListPriorities(options.Searcher, int64, int64) ([]*model.Priority, error)
	LocatePriority(options.Searcher) (*model.Priority, error)
	CreatePriority(options.Creator, *model.Priority) (*model.Priority, error)
	UpdatePriority(options.Updator, *model.Priority) (*model.Priority, error)
	DeletePriority(options.Deleter) (*model.Priority, error)
}

type PriorityService struct {
	app PriorityHandler
	api.UnimplementedPrioritiesServer
}

func NewPriorityService(app PriorityHandler) (*PriorityService, error) {
	if app == nil {
		return nil, errors.New("priority handler is nil")
	}
	return &PriorityService{app: app}, nil
}

var PriorityMetadata = model.NewObjectMetadata(model.ScopeDictionary, "", []*model.Field{
	{Name: "id", Default: true},
	{Name: "created_by", Default: true},
	{Name: "created_at", Default: true},
	{Name: "updated_by", Default: false},
	{Name: "updated_at", Default: false},
	{Name: "name", Default: true},
	{Name: "description", Default: true},
	{Name: "color", Default: true},
})

func (s *PriorityService) CreatePriority(ctx context.Context, req *api.CreatePriorityRequest) (*api.Priority, error) {
	createOpts, err := grpcopts.NewCreateOptions(
		ctx,
		grpcopts.WithCreateFields(req, PriorityMetadata),
	)
	if err != nil {
		return nil, err
	}
	input := &model.Priority{
		Name:        req.Input.Name,
		Description: &req.Input.Description,
		Color:       req.Input.Color,
	}
	m, err := s.app.CreatePriority(createOpts, input)
	if err != nil {
		return nil, err
	}
	return s.Marshal(m)
}

func (s *PriorityService) ListPriorities(ctx context.Context, req *api.ListPriorityRequest) (*api.PriorityList, error) {
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
		return nil, err
	}
	searchOpts.AddFilter("name", req.Q)

	items, err := s.app.ListPriorities(searchOpts, req.NotInSla, req.InSlaCond)
	if err != nil {
		return nil, err
	}
	var res api.PriorityList
	res.Items, err = utils.ConvertToOutputBulk(items, s.Marshal)
	if err != nil {
		return nil, err
	}

	res.Next, res.Items = utils.GetListResult(searchOpts, res.Items)
	res.Page = req.GetPage()

	return &res, nil
}

func (s *PriorityService) UpdatePriority(ctx context.Context, req *api.UpdatePriorityRequest) (*api.Priority, error) {
	updateOpts, err := grpcopts.NewUpdateOptions(
		ctx,
		grpcopts.WithUpdateFields(req, PriorityMetadata),
		grpcopts.WithUpdateMasker(req),
		grpcopts.WithUpdateIDs([]int64{req.GetId()}),
	)
	if err != nil {
		return nil, err
	}
	input := &model.Priority{
		Id:          req.Id,
		Name:        req.Input.Name,
		Description: &req.Input.Description,
		Color:       req.Input.Color,
	}
	m, err := s.app.UpdatePriority(updateOpts, input)
	if err != nil {
		return nil, err
	}
	return s.Marshal(m)
}

func (s *PriorityService) DeletePriority(ctx context.Context, req *api.DeletePriorityRequest) (*api.Priority, error) {
	deleteOpts, err := grpcopts.NewDeleteOptions(ctx, grpcopts.WithDeleteID(req.Id))
	if err != nil {
		return nil, err
	}
	item, err := s.app.DeletePriority(deleteOpts)
	if err != nil {
		return nil, err
	}
	return s.Marshal(item)
}

func (s *PriorityService) LocatePriority(ctx context.Context, req *api.LocatePriorityRequest) (*api.LocatePriorityResponse, error) {
	opts, err := grpcopts.NewLocateOptions(ctx, grpcopts.WithFields(req, PriorityMetadata,
		util.DeduplicateFields,
		util.EnsureIdField,
	), grpcopts.WithID(req.Id))
	if err != nil {
		return nil, err
	}
	items, err := s.app.ListPriorities(opts, 0, 0)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, errors.New("no records found", errors.WithCode(codes.NotFound))
	}
	if len(items) > 1 {
		return nil, errors.New("too many records found", errors.WithCode(codes.InvalidArgument))
	}
	res, err := s.Marshal(items[0])
	if err != nil {
		return nil, err
	}
	return &api.LocatePriorityResponse{Priority: res}, nil
}

func (s *PriorityService) Marshal(model *model.Priority) (*api.Priority, error) {
	if model == nil {
		return nil, nil
	}
	return &api.Priority{
		Id:          model.Id,
		Name:        model.Name,
		Description: utils.Dereference(model.Description),
		Color:       model.Color,
		CreatedAt:   utils.MarshalTime(model.CreatedAt),
		UpdatedAt:   utils.MarshalTime(model.UpdatedAt),
		CreatedBy:   utils.MarshalLookup(model.Author),
		UpdatedBy:   utils.MarshalLookup(model.Editor),
	}, nil
}
