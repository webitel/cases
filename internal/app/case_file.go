package app

import (
	"context"
	"github.com/webitel/cases/api/cases"
	"github.com/webitel/cases/auth"
	cerror "github.com/webitel/cases/internal/errors"
	deferr "github.com/webitel/cases/internal/errors/defaults"
	"github.com/webitel/cases/model"
	grpcopts "github.com/webitel/cases/model/options/grpc"
	"github.com/webitel/cases/util"
	"github.com/webitel/webitel-go-kit/etag"
	"log/slog"
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
	if err != nil {
		return nil, NewBadRequestError(err)
	}

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
			return nil, deferr.ForbiddenError
		}
		if !access {
			slog.ErrorContext(ctx, "user doesn't have required (READ) access to the case", logAttributes)
			return nil, deferr.ForbiddenError
		}
	}
	files, err := c.app.Store.CaseFile().List(searchOpts)
	if err != nil {
		slog.ErrorContext(ctx, err.Error(), logAttributes)
		return nil, deferr.DatabaseError
	}
	return files, nil
}

func (c *CaseFileService) DeleteFile(ctx context.Context, req *cases.DeleteFileRequest) (*cases.File, error) {
	if req.Id == 0 {
		return nil, cerror.NewBadRequestError("app.case_file.delete_file.file_id_required", "File ID is required")
	}
	deleteOpts, err := grpcopts.NewDeleteOptions(ctx, grpcopts.WithDeleteID(req.GetId()), grpcopts.WithDeleteParentIDAsEtag(etag.EtagCase, req.GetCaseEtag()))
	if err != nil {
		return nil, NewBadRequestError(err)
	}

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
			return nil, deferr.ForbiddenError
		}
		if !access {
			slog.ErrorContext(ctx, "user doesn't have required (DELETE) access to the case", logAttributes)
			return nil, deferr.ForbiddenError
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
			return nil, deferr.DatabaseError
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
