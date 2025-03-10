package app

import (
	"context"
	grpcopts "github.com/webitel/cases/model/options/grpc"
	"github.com/webitel/cases/util"
	"log/slog"

	"github.com/webitel/cases/auth"

	"github.com/webitel/cases/api/cases"
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
	//{Name: "url", Default: true},
	{Name: "author", Default: true},
})

func (c *CaseFileService) ListFiles(ctx context.Context, req *cases.ListFilesRequest) (*cases.CaseFileList, error) {
	if req.CaseEtag == "" {
		return nil, cerror.NewBadRequestError("app.case_file.list_files.case_etag_required", "Case Etag is required")
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

	tag, err := etag.EtagOrId(etag.EtagCase, req.CaseEtag)
	if err != nil {
		return nil, cerror.NewBadRequestError("app.case_file.list_files.invalid_case_etag", "Invalid Case Etag")
	}
	searchOpts.AddFilter("case_id", tag.GetOid())
	logAttributes := slog.Group("context", slog.Int64("user_id", searchOpts.GetAuthOpts().GetUserId()), slog.Int64("domain_id", searchOpts.GetAuthOpts().GetDomainId()))
	accessMode := auth.Read
	if searchOpts.GetAuthOpts().IsRbacCheckRequired(CaseFileMetadata.GetParentScopeName(), accessMode) {
		access, err := c.app.Store.Case().CheckRbacAccess(searchOpts, searchOpts.GetAuthOpts(), accessMode, tag.GetOid())
		if err != nil {
			slog.ErrorContext(ctx, err.Error(), logAttributes)
			return nil, ForbiddenError
		}
		if !access {
			slog.ErrorContext(ctx, "user doesn't have required (READ) access to the case", logAttributes)
			return nil, ForbiddenError
		}
	}
	files, err := c.app.Store.CaseFile().List(searchOpts)
	if err != nil {
		slog.ErrorContext(ctx, err.Error(), logAttributes)
		return nil, DatabaseError
	}
	return files, nil
}

func (c *CaseFileService) DeleteFile(ctx context.Context, req *cases.DeleteFileRequest) (*cases.File, error) {
	if req.Id == 0 {
		return nil, cerror.NewBadRequestError("app.case_file.delete_file.file_id_required", "File ID is required")
	}

	deleteOpts, err := model.NewDeleteOptions(ctx, CaseFileMetadata)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return nil, InternalError
	}

	caseTID, err := etag.EtagOrId(etag.EtagCase, req.GetCaseEtag())
	if err != nil {
		return nil, cerror.NewBadRequestError("app.case_file.delete.case_etag.parse.error", err.Error())
	}

	deleteOpts.ID = req.Id
	deleteOpts.ParentID = caseTID.GetOid()

	logAttributes := slog.Group(
		"context",
		slog.Int64(
			"user_id",
			deleteOpts.GetAuthOpts().GetUserId(),
		),
		slog.Int64(
			"domain_id",
			deleteOpts.GetAuthOpts().GetDomainId(),
		))
	// Check if the user has permission to delete the file
	accessMode := auth.Delete
	if deleteOpts.GetAuthOpts().IsRbacCheckRequired(CaseFileMetadata.GetParentScopeName(), accessMode) {
		access, err := c.app.Store.Case().CheckRbacAccess(
			deleteOpts,
			deleteOpts.GetAuthOpts(),
			accessMode,
			deleteOpts.ParentID,
		)
		if err != nil {
			slog.ErrorContext(ctx, err.Error(), logAttributes)
			return nil, ForbiddenError
		}
		if !access {
			slog.ErrorContext(ctx, "user doesn't have required (DELETE) access to the case", logAttributes)
			return nil, ForbiddenError
		}
	}

	// Delete the file from the database
	err = c.app.Store.CaseFile().Delete(deleteOpts)
	if err != nil {
		switch err.(type) {
		case *cerror.DBNoRowsError:
			return nil, cerror.NewBadRequestError(
				"app.case_file.delete.not_found",
				"delete not allowed",
			)
		default:
			slog.ErrorContext(ctx, err.Error(), logAttributes)
			return nil, DatabaseError
		}
	}

	return &cases.File{}, nil
}

func NewCaseFileService(app *App) (*CaseFileService, cerror.AppError) {
	if app == nil {
		return nil, cerror.NewBadRequestError("app.case.new_case_file_service.check_args.app", "unable to init service, app is nil")
	}
	return &CaseFileService{app: app}, nil
}
