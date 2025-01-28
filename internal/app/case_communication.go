package app

import (
	"context"
	defErr "errors"
	"fmt"
	"github.com/webitel/cases/api/cases"
	"github.com/webitel/cases/auth"
	errors "github.com/webitel/cases/internal/error"
	"github.com/webitel/cases/model"
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

func (c *CaseCommunicationService) ListCommunications(ctx context.Context, request *cases.ListCommunicationsRequest) (*cases.ListCommunicationsResponse, error) {
	tag, err := etag.EtagOrId(etag.EtagCase, request.GetCaseEtag())
	if err != nil {
		return nil, errors.NewBadRequestError("app.case_communication.list_communication.invalid_etag", "Invalid case etag")
	}
	searchOpts, err := model.NewSearchOptions(ctx, request, CaseCommunicationMetadata)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return nil, AppInternalError
	}
	searchOpts.ParentId = tag.GetOid()
	logAttributes := slog.Group("context", slog.Int64("case_id", tag.GetOid()), slog.Int64("user_id", searchOpts.GetAuthOpts().GetUserId()), slog.Int64("domain_id", searchOpts.GetAuthOpts().GetDomainId()))

	if searchOpts.GetAuthOpts().GetObjectScope(CaseCommunicationMetadata.GetParentScopeName()).IsRbacUsed() {
		access, err := c.app.Store.Case().CheckRbacAccess(searchOpts, searchOpts.GetAuthOpts(), auth.Read, searchOpts.ParentId)
		if err != nil {
			slog.ErrorContext(ctx, err.Error(), logAttributes)
			return nil, AppForbiddenError
		}
		if !access {
			slog.ErrorContext(ctx, "user doesn't have required (READ) access to the case", logAttributes)
			return nil, AppForbiddenError
		}
	}

	res, dbErr := c.app.Store.CaseCommunication().List(searchOpts)
	if dbErr != nil {
		slog.ErrorContext(ctx, dbErr.Error(), slog.Int64("id", tag.GetOid()))
		return nil, AppDatabaseError
	}
	err = NormalizeResponseCommunications(res.Data, request.GetFields())
	if err != nil {
		slog.ErrorContext(ctx, err.Error(), logAttributes)
		return nil, AppResponseNormalizingError
	}
	return res, nil
}

func (c *CaseCommunicationService) LinkCommunication(ctx context.Context, request *cases.LinkCommunicationRequest) (*cases.LinkCommunicationResponse, error) {
	if len(request.Input) == 0 {
		return nil, errors.NewBadRequestError("app.case_communication.link_communication.check_args.payload", "no payload")
	}
	err := ValidateCaseCommunicationsCreate(request.Input...)
	if err != nil {
		return nil, errors.NewBadRequestError("app.case_communication.link_communication.validate_payload.error", err.Error())
	}
	tag, err := etag.EtagOrId(etag.EtagCase, request.GetCaseEtag())
	if err != nil {
		return nil, errors.NewBadRequestError("app.case_communication.link_communication.invalid_etag", "Invalid case etag")
	}
	createOpts, err := model.NewCreateOptions(ctx, request, CaseCommunicationMetadata)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return nil, AppInternalError
	}
	createOpts.ParentID = tag.GetOid()
	logAttributes := slog.Group("context", slog.Int64("user_id", createOpts.GetAuthOpts().GetUserId()), slog.Int64("domain_id", createOpts.GetAuthOpts().GetDomainId()), slog.Int64("case_id", createOpts.ParentID))
	accessMode := auth.Edit
	if !createOpts.GetAuthOpts().CheckObacAccess(CaseCommunicationMetadata.GetParentScopeName(), accessMode) {
		slog.ErrorContext(ctx, "user doesn't have required (EDIT) access to the case", logAttributes)
		return nil, AppForbiddenError
	}
	if createOpts.GetAuthOpts().GetObjectScope(CaseCommunicationMetadata.GetParentScopeName()).IsRbacUsed() {
		access, err := c.app.Store.Case().CheckRbacAccess(createOpts, createOpts.GetAuthOpts(), accessMode, createOpts.ParentID)
		if err != nil {
			slog.ErrorContext(ctx, err.Error(), logAttributes)
			return nil, AppForbiddenError
		}
		if !access {
			slog.ErrorContext(ctx, "user doesn't have required (EDIT) access to the case", logAttributes)
			return nil, AppForbiddenError
		}
	}

	res, dbErr := c.app.Store.CaseCommunication().Link(createOpts, request.Input)
	if dbErr != nil {
		slog.ErrorContext(ctx, dbErr.Error(), logAttributes)
		return nil, AppInternalError
	}
	if len(res) == 0 {
		return nil, errors.NewBadRequestError("app.case_communication.link_communication.result.no_response", "no rows were affected (wrong ids or insufficient rights)")
	}
	err = NormalizeResponseCommunications(res, request.GetFields())
	if err != nil {
		slog.ErrorContext(ctx, err.Error(), logAttributes)
		return nil, AppResponseNormalizingError
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
	deleteOpts, err := model.NewDeleteOptions(ctx, CaseCommunicationMetadata)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return nil, AppInternalError
	}
	deleteOpts.IDs = []int64{tag.GetOid()}
	logAttributes := slog.Group("context", slog.Int64("user_id", deleteOpts.GetAuthOpts().GetUserId()), slog.Int64("domain_id", deleteOpts.GetAuthOpts().GetDomainId()))

	if deleteOpts.GetAuthOpts().GetObjectScope(CaseCommunicationMetadata.GetParentScopeName()).IsRbacUsed() {
		access, err := c.app.Store.Case().CheckRbacAccess(deleteOpts, deleteOpts.GetAuthOpts(), auth.Edit, caseTag.GetOid())
		if err != nil {
			slog.ErrorContext(ctx, err.Error(), logAttributes)
			return nil, AppForbiddenError
		}
		if !access {
			slog.ErrorContext(ctx, "user doesn't have required (EDIT) access to the case", logAttributes)
			return nil, AppForbiddenError
		}
	}

	affected, dbErr := c.app.Store.CaseCommunication().Unlink(deleteOpts)
	if dbErr != nil {
		slog.ErrorContext(ctx, dbErr.Error(), slog.Int64("id", tag.GetOid()))
		return nil, AppDatabaseError
	}
	return &cases.UnlinkCommunicationResponse{Affected: affected}, nil
}

func NewCaseCommunicationService(app *App) (*CaseCommunicationService, errors.AppError) {
	return &CaseCommunicationService{app: app}, nil
}

func NormalizeResponseCommunications(res []*cases.CaseCommunication, requestedFields []string) error {
	if len(requestedFields) == 0 {
		requestedFields = CaseCommentMetadata.GetDefaultFields()
	}
	hasEtag, hasId, hasVer := util.FindEtagFields(requestedFields)
	for _, re := range res {
		err := util.NormalizeEtags(etag.EtagCase, hasEtag, hasId, hasVer, &re.Etag, &re.Id, &re.Ver)
		if err != nil {
			return err
		}

	}
	return nil
}

func ValidateCaseCommunicationsCreate(input ...*cases.InputCaseCommunication) error {
	errText := "validation errors: "
	for i, communication := range input {
		errText := fmt.Sprintf("([%v]: ", i)
		if communication.CommunicationId == "" {
			errText += "communication can't be empty;"
		}
		if communication.CommunicationType <= 0 {
			errText += "communication type can't be empty;"
		}
		var typeFound bool
		for _, i := range cases.CaseCommunicationsTypes_value {
			if i == int32(communication.CommunicationType) {
				typeFound = true
				break
			}
		}
		if !typeFound {
			errText += "communication type not allowed;"
		}
		errText += ") "
	}
	return defErr.New(errText)
}
