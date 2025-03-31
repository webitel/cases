package app

import (
	"context"
	defErr "errors"
	"fmt"
	"github.com/webitel/cases/api/cases"
	"github.com/webitel/cases/auth"
	"github.com/webitel/cases/internal/errors"
	deferr "github.com/webitel/cases/internal/errors/defaults"
	"github.com/webitel/cases/model"
	grpcopts "github.com/webitel/cases/model/options/grpc"
	"github.com/webitel/cases/util"
	"github.com/webitel/webitel-go-kit/etag"
	"log/slog"
)

var CaseCommunicationMetadata = model.NewObjectMetadata("", caseObjScope, []*model.Field{
	{Name: "etag", Default: true},
	{Name: "ver", Default: false},
	{"id", true},
	{"communication_type", true},
	{"communication_id", true},
})

type CaseCommunicationService struct {
	app *App
	cases.UnimplementedCaseCommunicationsServer
}

func (c *CaseCommunicationService) ListCommunications(
	ctx context.Context,
	request *cases.ListCommunicationsRequest,
) (*cases.ListCommunicationsResponse, error) {
	searchOpts, err := grpcopts.NewSearchOptions(
		ctx,
		grpcopts.WithSearch(request),
		grpcopts.WithPagination(request),
		grpcopts.WithFields(request, CaseCommunicationMetadata,
			util.DeduplicateFields,
			util.ParseFieldsForEtag,
		),
		grpcopts.WithSort(request),
	)
	if err != nil {
		return nil, NewBadRequestError(err)
	}

	tag, err := etag.EtagOrId(etag.EtagCase, request.GetCaseEtag())
	if err != nil {
		return nil, errors.NewBadRequestError(
			"app.case_communication.list_communication.invalid_etag",
			"Invalid case etag",
		)
	}
	searchOpts.AddFilter("case_id", tag.GetOid())
	logAttributes := slog.Group("context", slog.Int64("case_id", tag.GetOid()), slog.Int64("user_id", searchOpts.GetAuthOpts().GetUserId()), slog.Int64("domain_id", searchOpts.GetAuthOpts().GetDomainId()))

	if searchOpts.GetAuthOpts().GetObjectScope(CaseCommunicationMetadata.GetParentScopeName()).IsRbacUsed() {
		access, err := c.app.Store.Case().CheckRbacAccess(searchOpts, searchOpts.GetAuthOpts(), auth.Read, tag.GetOid())
		if err != nil {
			slog.ErrorContext(ctx, err.Error(), logAttributes)
			return nil, deferr.ForbiddenError
		}
		if !access {
			slog.ErrorContext(ctx, "user doesn't have required (READ) access to the case", logAttributes)
			return nil, deferr.ForbiddenError
		}
	}

	res, dbErr := c.app.Store.CaseCommunication().List(searchOpts)
	if dbErr != nil {
		slog.ErrorContext(ctx, dbErr.Error(), slog.Int64("id", tag.GetOid()))
		return nil, deferr.DatabaseError
	}
	err = NormalizeResponseCommunications(res.Data, request.GetFields())
	if err != nil {
		slog.ErrorContext(ctx, err.Error(), logAttributes)
		return nil, deferr.ResponseNormalizingError
	}
	return res, nil
}

func (c *CaseCommunicationService) LinkCommunication(ctx context.Context, request *cases.LinkCommunicationRequest) (*cases.LinkCommunicationResponse, error) {
	err := ValidateCaseCommunicationsCreate(request.Input)
	if err != nil {
		return nil, errors.NewBadRequestError("app.case_communication.link_communication.validate_payload.error", err.Error())
	}
	tag, err := etag.EtagOrId(etag.EtagCase, request.GetCaseEtag())
	if err != nil {
		return nil, errors.NewBadRequestError("app.case_communication.link_communication.invalid_etag", "Invalid case etag")
	}
	createOpts, err := grpcopts.NewCreateOptions(
		ctx,
		grpcopts.WithCreateFields(request, CaseCommunicationMetadata,
			util.DeduplicateFields,
			util.ParseFieldsForEtag,
			util.EnsureIdField,
		),
		grpcopts.WithCreateParentID(tag.GetOid()),
	)
	if err != nil {
		return nil, NewBadRequestError(err)
	}
	logAttributes := slog.Group("context", slog.Int64("user_id", createOpts.GetAuthOpts().GetUserId()), slog.Int64("domain_id", createOpts.GetAuthOpts().GetDomainId()), slog.Int64("case_id", createOpts.ParentID))
	accessMode := auth.Edit
	if !createOpts.GetAuthOpts().CheckObacAccess(CaseCommunicationMetadata.GetParentScopeName(), accessMode) {
		slog.ErrorContext(ctx, "user doesn't have required (EDIT) access to the case", logAttributes)
		return nil, deferr.ForbiddenError
	}
	if createOpts.GetAuthOpts().GetObjectScope(CaseCommunicationMetadata.GetParentScopeName()).IsRbacUsed() {
		access, err := c.app.Store.Case().CheckRbacAccess(createOpts, createOpts.GetAuthOpts(), accessMode, createOpts.ParentID)
		if err != nil {
			slog.ErrorContext(ctx, err.Error(), logAttributes)
			return nil, deferr.ForbiddenError
		}
		if !access {
			slog.ErrorContext(ctx, "user doesn't have required (EDIT) access to the case", logAttributes)
			return nil, deferr.ForbiddenError
		}
	}

	res, dbErr := c.app.Store.CaseCommunication().Link(createOpts, []*cases.InputCaseCommunication{request.Input})
	if dbErr != nil {
		slog.ErrorContext(ctx, dbErr.Error(), logAttributes)
		return nil, deferr.InternalError
	}
	if len(res) == 0 {
		return nil, errors.NewBadRequestError("app.case_communication.link_communication.result.no_response", "no rows were affected (wrong ids or insufficient rights)")
	}
	err = NormalizeResponseCommunications(res, request.GetFields())
	if err != nil {
		slog.ErrorContext(ctx, err.Error(), logAttributes)
		return nil, deferr.ResponseNormalizingError
	}
	return &cases.LinkCommunicationResponse{Data: res}, nil
}

func (c *CaseCommunicationService) UnlinkCommunication(ctx context.Context, request *cases.UnlinkCommunicationRequest) (*cases.UnlinkCommunicationResponse, error) {
	tag, err := etag.EtagOrId(etag.EtagCaseCommunication, request.GetId())
	if err != nil {
		return nil, errors.NewBadRequestError("app.case_communication.unlink_communication.invalid_etag", "Invalid communication etag")
	}
	caseTag, err := etag.EtagOrId(etag.EtagCase, request.GetCaseEtag())
	if err != nil {
		return nil, errors.NewBadRequestError("app.case_communication.unlink_communication.invalid_etag", "Invalid case etag")
	}
	deleteOpts, err := grpcopts.NewDeleteOptions(ctx, grpcopts.WithDeleteID(tag.GetOid()))
	if err != nil {
		return nil, NewBadRequestError(err)
	}
	deleteOpts.IDs = []int64{tag.GetOid()}
	logAttributes := slog.Group("context", slog.Int64("user_id", deleteOpts.GetAuthOpts().GetUserId()), slog.Int64("domain_id", deleteOpts.GetAuthOpts().GetDomainId()))

	if deleteOpts.GetAuthOpts().GetObjectScope(CaseCommunicationMetadata.GetParentScopeName()).IsRbacUsed() {
		access, err := c.app.Store.Case().CheckRbacAccess(deleteOpts, deleteOpts.GetAuthOpts(), auth.Edit, caseTag.GetOid())
		if err != nil {
			slog.ErrorContext(ctx, err.Error(), logAttributes)
			return nil, deferr.ForbiddenError
		}
		if !access {
			slog.ErrorContext(ctx, "user doesn't have required (EDIT) access to the case", logAttributes)
			return nil, deferr.ForbiddenError
		}
	}

	affected, dbErr := c.app.Store.CaseCommunication().Unlink(deleteOpts)
	if dbErr != nil {
		slog.ErrorContext(ctx, dbErr.Error(), slog.Int64("id", tag.GetOid()))
		return nil, deferr.DatabaseError
	}
	return &cases.UnlinkCommunicationResponse{Affected: affected}, nil
}

func NewCaseCommunicationService(app *App) (*CaseCommunicationService, errors.AppError) {
	return &CaseCommunicationService{app: app}, nil
}

func NormalizeResponseCommunications(res []*cases.CaseCommunication, requestedFields []string) error {
	if len(requestedFields) == 0 {
		requestedFields = CaseCommunicationMetadata.GetDefaultFields()
	}
	hasEtag, hasId, hasVer := util.FindEtagFields(requestedFields)
	for _, re := range res {
		err := util.NormalizeEtags(etag.EtagCaseCommunication, hasEtag, hasId, hasVer, &re.Etag, &re.Id, &re.Ver)
		if err != nil {
			return err
		}

	}
	return nil
}

func ValidateCaseCommunicationsCreate(input ...*cases.InputCaseCommunication) error {

	var errorsSlice []error
	for i, communication := range input {
		if communication == nil {
			errorsSlice = append(errorsSlice, fmt.Errorf("input[%d]: empty entry", i))
			continue
		}
		if communication.CommunicationId == "" {
			errorsSlice = append(errorsSlice, fmt.Errorf("input[%d]: communication can't be empty", i))
		}
		if communication.CommunicationType.GetId() == 0 {
			errorsSlice = append(errorsSlice, fmt.Errorf("input[%d]: communication id can't be empty", i))
		}
	}
	if errorsSlice != nil {
		return defErr.Join(errorsSlice...)
	}
	return nil
}
