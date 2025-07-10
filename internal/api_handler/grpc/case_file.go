package grpc

import (
	"context"
	"github.com/webitel/cases/api/cases"
	"github.com/webitel/cases/internal/api_handler/grpc/utils"
	"github.com/webitel/cases/internal/errors"
	"github.com/webitel/cases/internal/model"
	"github.com/webitel/cases/internal/model/options"
	grpcopts "github.com/webitel/cases/internal/model/options/grpc"
	"github.com/webitel/cases/util"
	"github.com/webitel/webitel-go-kit/pkg/etag"
	"google.golang.org/grpc/codes"
)

type CaseFileHandler interface {
	ListCaseFiles(options.Searcher) ([]*model.CaseFile, error)
	DeleteCaseFile(options.Deleter) (*model.CaseFile, error)
}

type CaseFileService struct {
	app CaseFileHandler
	cases.UnimplementedCaseFilesServer
}

func NewCaseFileService(app CaseFileHandler) (*CaseFileService, error) {
	if app == nil {
		return nil, errors.New("case file handler is nil")
	}
	return &CaseFileService{app: app}, nil
}

var CaseFileMetadata = model.NewObjectMetadata("", "", []*model.Field{
	{Name: "id", Default: true},
	{Name: "size", Default: true},
	{Name: "mime", Default: true},
	{Name: "name", Default: true},
	{Name: "created_at", Default: true},
	{Name: "created_by", Default: true},
	//{Name: "url", Default: true},
	{Name: "source", Default: true},
})

func (s *CaseFileService) ListFiles(ctx context.Context, req *cases.ListFilesRequest) (*cases.CaseFileList, error) {
	if req.CaseEtag == "" {
		return nil, errors.New("case etag is required", errors.WithCode(codes.InvalidArgument))
	}

	searchOpts, err := grpcopts.NewSearchOptions(
		ctx,
		grpcopts.WithSearch(req),
		grpcopts.WithPagination(req),
		grpcopts.WithSort(req),
		grpcopts.WithFields(req, CaseFileMetadata,
			util.DeduplicateFields,
			util.ParseFieldsForEtag,
			util.EnsureIdField,
		),
	)
	if err != nil {
		return nil, errors.New(err.Error(), errors.WithCode(codes.InvalidArgument))
	}

	tag, err := etag.EtagOrId(etag.EtagCase, req.CaseEtag)
	if err != nil {
		return nil, errors.New("invalid case etag", errors.WithCode(codes.InvalidArgument))
	}
	if tag.GetOid() != 0 {
		searchOpts.AddFilter(util.EqualFilter("case_id", tag.GetOid()))
	}

	files, err := s.app.ListCaseFiles(searchOpts)
	if err != nil {
		return nil, err
	}

	var res cases.CaseFileList
	converted, err := utils.ConvertToOutputBulk(files, s.Marshal)
	if err != nil {
		return nil, err
	}
	res.Next, res.Items = utils.GetListResult(searchOpts, converted)
	res.Page = int64(req.GetPage())

	return &res, nil
}

func (s *CaseFileService) DeleteFile(ctx context.Context, req *cases.DeleteFileRequest) (*cases.File, error) {
	if req.Id == 0 {
		return nil, errors.New("file id is required", errors.WithCode(codes.InvalidArgument))
	}
	deleteOpts, err := grpcopts.NewDeleteOptions(
		ctx,
		grpcopts.WithDeleteID(req.GetId()),
		grpcopts.WithDeleteParentIDAsEtag(etag.EtagCase, req.GetCaseEtag()),
	)
	if err != nil {
		return nil, errors.New(err.Error(), errors.WithCode(codes.InvalidArgument))
	}

	item, err := s.app.DeleteCaseFile(deleteOpts)
	if err != nil {
		return nil, err
	}
	return s.Marshal(item)
}

// Marshal converts a model.CaseFile to cases.File
func (s *CaseFileService) Marshal(m *model.CaseFile) (*cases.File, error) {
	if m == nil {
		return nil, nil
	}
	return &cases.File{
		Id:        int64(m.Id),
		CreatedAt: utils.MarshalTime(m.CreatedAt),
		Size:      m.Size,
		Mime:      m.Mime,
		Name:      m.Name,
		Url:       m.Url,
		CreatedBy: utils.MarshalExtendedLookup(m.Author),
		Source:    m.Source,
	}, nil
}
