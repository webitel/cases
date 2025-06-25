package app

import (
	"context"
	"github.com/webitel/cases/api/cases"
	"github.com/webitel/cases/auth"
	errors "github.com/webitel/cases/internal/errors"
	"github.com/webitel/cases/internal/model"
	grpcopts "github.com/webitel/cases/internal/model/options/grpc"
	"github.com/webitel/cases/util"
	"github.com/webitel/webitel-go-kit/pkg/etag"
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
		return nil, errors.InvalidArgument("Case Etag is required")
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
		return nil, err
	}

	tag, err := etag.EtagOrId(etag.EtagCase, req.CaseEtag)
	if err != nil {
		return nil, errors.InvalidArgument("Invalid Case Etag", errors.WithCause(err))
	}
	searchOpts.AddFilter("case_id", tag.GetOid())

	accessMode := auth.Read
	if searchOpts.GetAuthOpts().IsRbacCheckRequired(CaseFileMetadata.GetParentScopeName(), accessMode) {
		access, err := c.app.Store.Case().CheckRbacAccess(searchOpts, searchOpts.GetAuthOpts(), accessMode, tag.GetOid())
		if err != nil {
			slog.ErrorContext(ctx, err.Error())
			return nil, err
		}
		if !access {
			return nil, errors.Forbidden("user doesn't have required (READ) access to the case")
		}
	}
	files, err := c.app.Store.CaseFile().List(searchOpts)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return nil, err
	}
	return files, nil
}

func (c *CaseFileService) DeleteFile(ctx context.Context, req *cases.DeleteFileRequest) (*cases.File, error) {
	if req.Id == 0 {
		return nil, errors.InvalidArgument("File ID is required")
	}
	deleteOpts, err := grpcopts.NewDeleteOptions(ctx, grpcopts.WithDeleteID(req.GetId()), grpcopts.WithDeleteParentIDAsEtag(etag.EtagCase, req.GetCaseEtag()))
	if err != nil {
		return nil, err
	}

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
			return nil, err
		}
		if !access {
			return nil, errors.Forbidden("user doesn't have required (DELETE) access to the case")
		}
	}

	// Delete the file from the database
	err = c.app.Store.CaseFile().Delete(deleteOpts)
	if err != nil {
		return nil, err

	}

	return &cases.File{}, nil
}

func NewCaseFileService(app *App) (*CaseFileService, error) {
	if app == nil {
		return nil, errors.New("unable to init service, app is nil")
	}
	return &CaseFileService{app: app}, nil
}
