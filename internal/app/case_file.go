package app

import (
	"context"

	cases "github.com/webitel/cases/api/cases"
	cerror "github.com/webitel/cases/internal/error"
	"github.com/webitel/cases/model"
	"github.com/webitel/cases/util"
	"github.com/webitel/webitel-go-kit/etag"
)

type CaseFileService struct {
	app *App
	cases.UnimplementedCaseFilesServer
}

var CaseFileMetadata = model.NewObjectMetadata(
	[]*model.Field{
		{"size", true},
		{"mime", true},
		{"name", true},
		{"created_at", true},
	})

var defaultFieldsCaseFile = []string{"size", "mime", "name", "created_at"}

func (c *CaseFileService) ListFiles(ctx context.Context, req *cases.ListFilesRequest) (*cases.CaseFileList, error) {
	if req.CaseEtag == "" {
		return nil, cerror.NewBadRequestError("app.case_file.list_files.case_etag_required", "Case Etag is required")
	}

	fields := util.FieldsFunc(req.Fields, util.InlineFields)
	if len(fields) == 0 {
		fields = defaultFieldsCaseFile
	}

	page := req.Page
	if page == 0 {
		page = 1
	}

	tag, err := etag.EtagOrId(etag.EtagCaseComment, req.CaseEtag)
	if err != nil {
		return nil, cerror.NewBadRequestError("app.case_file.list_files.invalid_case_etag", "Invalid Case Etag")
	}

	ids, err := util.ParseIds(req.Ids, etag.EtagCaseComment)
	if err != nil {
		return nil, cerror.NewBadRequestError("app.case_file.list_files.invalid_qin", "Invalid Qin")
	}
	// Build search options
	searchOpts := model.NewSearchOptions(ctx, req, CaseFileMetadata)
	searchOpts.IDs = ids
	searchOpts.ParentId = tag.GetOid()

	files, err := c.app.Store.CaseFile().List(searchOpts)
	if err != nil {
		return nil, cerror.NewInternalError("app.case_file.list_files.fetch_error", err.Error())
	}
	if len(files.Items) == 0 {
		return nil, cerror.NewNotFoundError("app.case_file.list_files.not_found", "Files not found")
	}
	return files, nil
}

func NewCaseFileService(app *App) (*CaseFileService, cerror.AppError) {
	if app == nil {
		return nil, cerror.NewBadRequestError("app.case.new_case_file_service.check_args.app", "unable to init service, app is nil")
	}
	return &CaseFileService{app: app}, nil
}
