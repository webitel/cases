package app

import (
	"strconv"

	"google.golang.org/grpc/codes"

	"github.com/webitel/cases/auth"
	"github.com/webitel/cases/internal/api_handler/grpc"
	"github.com/webitel/cases/internal/errors"
	"github.com/webitel/cases/internal/model"
	"github.com/webitel/cases/internal/model/options"
)

func (c *App) ListCommunications(searcher options.Searcher) ([]*model.CaseCommunication, error) {
	filters := searcher.GetFilter("case_id")
	if len(filters) == 0 {
		return nil, errors.New("case id is required", errors.WithCode(codes.InvalidArgument))
	}
	caseID, err := strconv.ParseInt(filters[0].Value, 10, 64)
	if err != nil {
		return nil, errors.New("invalid case id", errors.WithCode(codes.InvalidArgument))
	}
	accessMode := auth.Read
	authOpts := searcher.GetAuthOpts()
	if authOpts.GetObjectScope(grpc.CaseCommunicationMetadata.GetParentScopeName()).IsRbacUsed() {
		access, err := c.Store.Case().CheckRbacAccess(searcher, authOpts, accessMode, caseID)
		if err != nil {
			return nil, err
		}

		if !access {
			return nil, errors.Forbidden("user doesn't have required (READ) access to the case")
		}
	}
	res, err := c.Store.CaseCommunication().List(searcher)
	if err != nil {
		return nil, err
	}
	return res, nil // Only return internal model
}

func (c *App) LinkCommunication(createOpts options.Creator, input []*model.CaseCommunication) ([]*model.CaseCommunication, error) {
	accessMode := auth.Edit
	if !createOpts.GetAuthOpts().CheckObacAccess(grpc.CaseCommunicationMetadata.GetParentScopeName(), accessMode) {
		return nil, errors.New("user doesn't have required (EDIT) access to the case", errors.WithCode(codes.PermissionDenied))
	}

	if createOpts.GetAuthOpts().GetObjectScope(grpc.CaseCommunicationMetadata.GetParentScopeName()).IsRbacUsed() {
		access, err := c.Store.Case().CheckRbacAccess(createOpts, createOpts.GetAuthOpts(), accessMode, createOpts.GetParentID())
		if err != nil {
			return nil, errors.New("user doesn't have required (EDIT) access to the case", errors.WithCode(codes.PermissionDenied))
		}

		if !access {
			return nil, errors.New("user doesn't have required (EDIT) access to the case", errors.WithCode(codes.PermissionDenied))
		}
	}
	res, err := c.Store.CaseCommunication().Link(createOpts, input)
	if err != nil {
		return nil, err
	}

	if len(res) == 0 {
		return nil, errors.New("no rows were affected (wrong ids or insufficient rights)", errors.WithCode(codes.InvalidArgument))
	}
	return res, nil // Only return internal model
}

func (c *App) UnlinkCommunication(deleteOpts options.Deleter) (int64, error) {
	caseIDStr := deleteOpts.GetFilter("case_id")
	if len(caseIDStr) == 0 {
		return 0, errors.New("case id is required", errors.WithCode(codes.InvalidArgument))
	}
	caseID, err := strconv.ParseInt(caseIDStr[0].Value, 10, 64)
	if err != nil {
		return 0, errors.New("invalid case id", errors.WithCode(codes.InvalidArgument))
	}
	accessMode := auth.Edit
	if deleteOpts.GetAuthOpts().GetObjectScope(grpc.CaseCommunicationMetadata.GetParentScopeName()).IsRbacUsed() {
		access, err := c.Store.Case().CheckRbacAccess(deleteOpts, deleteOpts.GetAuthOpts(), accessMode, caseID)
		if err != nil {
			return 0, errors.New("user doesn't have required (EDIT) access to the case", errors.WithCode(codes.PermissionDenied))
		}

		if !access {
			return 0, errors.New("user doesn't have required (EDIT) access to the case", errors.WithCode(codes.PermissionDenied))
		}
	}
	affected, err := c.Store.CaseCommunication().Unlink(deleteOpts)
	if err != nil {
		return 0, err
	}
	return affected, nil
}
