package grpc

import (
	"context"
	"github.com/webitel/cases/internal/api_handler/grpc/utils"
	"github.com/webitel/cases/internal/errors"
	"github.com/webitel/cases/internal/model/options"
	"google.golang.org/grpc/codes"
	"strings"

	_go "github.com/webitel/cases/api/cases"
	"github.com/webitel/cases/internal/model"
	grpcopts "github.com/webitel/cases/internal/model/options/grpc"
	"github.com/webitel/cases/util"
)

// SourceHandler defines the interface for managing sources.
type SourceHandler interface {
	CreateSource(options.Creator, *model.Source) (*model.Source, error)
	UpdateSource(options.Updator, *model.Source) (*model.Source, error)
	DeleteSource(options.Deleter) (*model.Source, error)
	ListSources(options.Searcher) ([]*model.Source, error)
}

// SourceService implements the gRPC server for sources.
type SourceService struct {
	app SourceHandler
	_go.UnimplementedSourcesServer
	objClassName string
}

// SourceMetadata defines the fields available for source objects.
var SourceMetadata = model.NewObjectMetadata(model.ScopeDictionary, "", []*model.Field{
	{"id", true},
	{"created_by", true},
	{"created_at", true},
	{"updated_by", false},
	{"updated_at", false},
	{"name", true},
	{"description", true},
	{"type", true},
})

// CreateSource handles the gRPC request to create a new source.
func (s *SourceService) CreateSource(
	ctx context.Context,
	req *_go.CreateSourceRequest,
) (*_go.Source, error) {
	if req.Input.Type == _go.SourceType_TYPE_UNSPECIFIED {
		return nil, errors.New("type is required", errors.WithCode(codes.InvalidArgument))
	}
	createOpts, err := grpcopts.NewCreateOptions(
		ctx,
		grpcopts.WithCreateFields(req, SourceMetadata),
	)
	if err != nil {
		return nil, err
	}
	sourceType := req.Input.Type.String()
	input := &model.Source{
		Name:        &req.Input.Name,
		Description: &req.Input.Description,
		Type:        &sourceType,
	}

	// Create the source in the store
	res, err := s.app.CreateSource(createOpts, input)
	if err != nil {
		return nil, err
	}

	return s.Marshal(res)
}

// ListSources handles the gRPC request to list sources with filters and pagination.
func (s *SourceService) ListSources(
	ctx context.Context,
	req *_go.ListSourceRequest,
) (*_go.SourceList, error) {
	searchOpts, err := grpcopts.NewSearchOptions(
		ctx,
		grpcopts.WithSearch(req),
		grpcopts.WithPagination(req),
		grpcopts.WithFields(req, SourceMetadata,
			util.DeduplicateFields,
			util.EnsureIdField,
		),
		grpcopts.WithSort(req),
		grpcopts.WithIDs(req.GetId()),
	)
	if err != nil {
		return nil, err
	}
	if req.Q != "" {
		searchOpts.AddFilter("name", req.Q)
	}
	if len(req.Type) > 0 {
		searchOpts.AddFilter("type", req.Type)
	}

	items, err := s.app.ListSources(searchOpts)
	if err != nil {
		return nil, err
	}

	var res _go.SourceList
	res.Items, err = utils.ConvertToOutputBulk(items, s.Marshal)
	if err != nil {
		return nil, err
	}
	res.Next, res.Items = utils.GetListResult(searchOpts, res.Items)
	res.Page = int32(searchOpts.GetPage())
	return &res, nil
}

// UpdateSource handles the gRPC request to update an existing source.
func (s *SourceService) UpdateSource(
	ctx context.Context,
	req *_go.UpdateSourceRequest,
) (*_go.Source, error) {
	updateOpts, err := grpcopts.NewUpdateOptions(
		ctx,
		grpcopts.WithUpdateFields(req, SourceMetadata),
		grpcopts.WithUpdateMasker(req),
		grpcopts.WithUpdateIDs([]int64{req.GetId()}),
	)
	if err != nil {
		return nil, err
	}
	sourceType := req.Input.Type.String()
	input := &model.Source{
		Id:          int(req.Id),
		Name:        &req.Input.Name,
		Description: &req.Input.Description,
		Type:        &sourceType,
	}

	// Update the source in the store
	res, err := s.app.UpdateSource(updateOpts, input)
	if err != nil {
		return nil, err
	}

	return s.Marshal(res)
}

// DeleteSource handles the gRPC request to delete a source.
func (s *SourceService) DeleteSource(
	ctx context.Context,
	req *_go.DeleteSourceRequest,
) (*_go.Source, error) {
	deleteOpts, err := grpcopts.NewDeleteOptions(ctx, grpcopts.WithDeleteID(req.Id))
	if err != nil {
		return nil, err
	}

	deleteOpts.IDs = []int64{req.Id}

	// Delete the source in the store
	_, err = s.app.DeleteSource(deleteOpts)
	if err != nil {
		return nil, err
	}

	return &(_go.Source{Id: req.Id}), nil
}

// LocateSource finds a source by ID and returns it, or an error if not found or ambiguous.
func (s *SourceService) LocateSource(ctx context.Context, req *_go.LocateSourceRequest) (*_go.LocateSourceResponse, error) {
	opts, err := grpcopts.NewLocateOptions(ctx, grpcopts.WithID(req.Id), grpcopts.WithFields(req, StatusConditionMetadata, util.EnsureIdField, util.DeduplicateFields))
	if err != nil {
		return nil, err
	}
	// Call the ListStatusConditions method
	items, err := s.app.ListSources(opts)
	if err != nil {
		return nil, err
	}

	// Check if the status condition was found
	if len(items) == 0 {
		return nil, errors.New("multiple rows found", errors.WithCode(codes.NotFound))
	}
	if len(items) > 1 {
		return nil, errors.New("not found", errors.WithCode(codes.InvalidArgument))
	}

	// Return the found status condition
	var res _go.LocateSourceResponse
	res.Source, err = s.Marshal(items[0])
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// Marshal converts a model.Source to its gRPC representation.
func (s *SourceService) Marshal(in *model.Source) (*_go.Source, error) {
	return &_go.Source{
		Id:          int64(in.Id),
		Name:        utils.Dereference(in.Name),
		Description: utils.Dereference(in.Description),
		Type:        stringToType(utils.Dereference(in.Type)),
		CreatedAt:   utils.MarshalTime(in.CreatedAt),
		UpdatedAt:   utils.MarshalTime(in.UpdatedAt),
		CreatedBy:   utils.MarshalLookup(in.Author),
		UpdatedBy:   utils.MarshalLookup(in.Author),
	}, nil
}

// NewSourceService constructs a new SourceService.
func NewSourceService(app SourceHandler) (*SourceService, error) {
	if app == nil {
		return nil, errors.New("source handler is nil")
	}
	return &SourceService{app: app, objClassName: model.ScopeDictionary}, nil
}

// stringToType converts a string into the corresponding SourceType enum value.
func stringToType(typeStr string) _go.SourceType {
	switch strings.ToUpper(typeStr) {
	case "CALL":
		return _go.SourceType_CALL
	case "CHAT":
		return _go.SourceType_CHAT
	case "SOCIAL_MEDIA":
		return _go.SourceType_SOCIAL_MEDIA
	case "EMAIL":
		return _go.SourceType_EMAIL
	case "API":
		return _go.SourceType_API
	case "MANUAL":
		return _go.SourceType_MANUAL
	default:
		return _go.SourceType_TYPE_UNSPECIFIED
	}
}
