package app

import (
	"context"
	"github.com/webitel/cases/auth"
	"log/slog"

	cases "github.com/webitel/cases/api/cases"
	cerror "github.com/webitel/cases/internal/errors"
	"github.com/webitel/cases/model"
	"github.com/webitel/webitel-go-kit/etag"
)

type CaseFileService struct {
	app *App
	cases.UnimplementedCaseFilesServer
}

var CaseFileMetadata = model.NewObjectMetadata("", caseObjScope, []*model.Field{
	{Name: "id", Default: true},
	{Name: "size", Default: true},
	{Name: "mime", Default: true},
	{Name: "name", Default: true},
	{Name: "created_at", Default: true},
	{Name: "created_by", Default: true},
	{Name: "url", Default: true},
	{Name: "author", Default: true},
})

func (c *CaseFileService) ListFiles(ctx context.Context, req *cases.ListFilesRequest) (*cases.CaseFileList, error) {
	if req.CaseEtag == "" {
		return nil, cerror.NewBadRequestError("app.case_file.list_files.case_etag_required", "Case Etag is required")
	}
	tag, err := etag.EtagOrId(etag.EtagCase, req.CaseEtag)
	if err != nil {
		return nil, cerror.NewBadRequestError("app.case_file.list_files.invalid_case_etag", "Invalid Case Etag")
	}
	// Build search options
	searchOpts, err := model.NewSearchOptions(ctx, req, CaseFileMetadata)
	if err != nil {
		slog.Error(err.Error())
		return nil, AppInternalError
	}
	searchOpts.ParentId = tag.GetOid()
	logAttributes := slog.Group("context", slog.Int64("user_id", searchOpts.GetAuthOpts().GetUserId()), slog.Int64("domain_id", searchOpts.GetAuthOpts().GetDomainId()))
	accessMode := auth.Read
	if searchOpts.GetAuthOpts().IsRbacCheckRequired(CaseFileMetadata.GetParentScopeName(), accessMode) {
		access, err := c.app.Store.Case().CheckRbacAccess(searchOpts, searchOpts.GetAuthOpts(), accessMode, searchOpts.ParentId)
		if err != nil {
			slog.Error(err.Error(), logAttributes)
			return nil, AppForbiddenError
		}
		if !access {
			slog.Error("user doesn't have required (READ) access to the case", logAttributes)
			return nil, AppForbiddenError
		}
	}
	files, err := c.app.Store.CaseFile().List(searchOpts)
	if err != nil {
		slog.Error(err.Error(), logAttributes)
		return nil, AppDatabaseError
	}
	return files, nil
}

func NewCaseFileService(app *App) (*CaseFileService, cerror.AppError) {
	if app == nil {
		return nil, cerror.NewBadRequestError("app.case.new_case_file_service.check_args.app", "unable to init service, app is nil")
	}
	return &CaseFileService{app: app}, nil
}
