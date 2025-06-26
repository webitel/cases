package app

import (
	"log/slog"

	"github.com/webitel/cases/auth"
	"github.com/webitel/cases/internal/api_handler/grpc"
	"github.com/webitel/cases/internal/errors"
	"github.com/webitel/cases/internal/model"
	"github.com/webitel/cases/internal/model/options"
	"google.golang.org/grpc/codes"
)

func (a *App) ListCaseFiles(rpc options.Searcher) ([]*model.CaseFile, error) {
	tag, ok := rpc.GetFilter("case_id").(int64)
	if !ok || tag == 0 {
		return nil, errors.New("case id required", errors.WithCode(codes.InvalidArgument))
	}
	accessMode := auth.Read
	logAttributes := slog.Group(
		"context",
		slog.Int64("user_id", rpc.GetAuthOpts().GetUserId()),
		slog.Int64("domain_id", rpc.GetAuthOpts().GetDomainId()),
	)
	if rpc.GetAuthOpts().IsRbacCheckRequired(grpc.CaseFileMetadata.GetParentScopeName(), accessMode) {
		access, err := a.Store.Case().CheckRbacAccess(rpc, rpc.GetAuthOpts(), accessMode, tag)
		if err != nil {
			slog.ErrorContext(rpc, err.Error(), logAttributes)
			return nil, errors.New("unable access resource", errors.WithCode(codes.PermissionDenied))
		}
		if !access {
			slog.ErrorContext(rpc, "user doesn't have required (READ) access to the case", logAttributes)
			return nil, errors.New("unable access resource", errors.WithCode(codes.PermissionDenied))
		}
	}
	files, err := a.Store.CaseFile().List(rpc)
	if err != nil {
		return nil, err
	}
	return files, nil
}

func (a *App) DeleteCaseFile(rpc options.Deleter) (*model.CaseFile, error) {
	if rpc.GetIDs() == nil || len(rpc.GetIDs()) == 0 {
		return nil, errors.New("file id required", errors.WithCode(codes.InvalidArgument))
	}
	logAttributes := slog.Group(
		"context",
		slog.Int64(
			"user_id",
			rpc.GetAuthOpts().GetUserId(),
		),
		slog.Int64(
			"domain_id",
			rpc.GetAuthOpts().GetDomainId(),
		))
	// Check if the user has permission to delete the file
	accessMode := auth.Delete
	if rpc.GetAuthOpts().IsRbacCheckRequired(grpc.CaseFileMetadata.GetParentScopeName(), accessMode) {
		access, err := a.Store.Case().CheckRbacAccess(
			rpc,
			rpc.GetAuthOpts(),
			accessMode,
			rpc.GetParentID(),
		)
		if err != nil {
			slog.ErrorContext(rpc, err.Error(), logAttributes)
			return nil, errors.New("unable access resource", errors.WithCode(codes.PermissionDenied))
		}
		if !access {
			slog.ErrorContext(rpc, "user doesn't have required (DELETE) access to the case", logAttributes)
			return nil, errors.New("unable access resource", errors.WithCode(codes.PermissionDenied))
		}
	}

	// Delete the file from the database
	file, err := a.Store.CaseFile().Delete(rpc)
	if err != nil {
		slog.ErrorContext(rpc, err.Error(), logAttributes)
		return nil, err
	}
	return file, nil
}
