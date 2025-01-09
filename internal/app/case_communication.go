package app

import (
	"context"
	defErr "errors"
	"fmt"
	"github.com/webitel/cases/api/cases"
	errors "github.com/webitel/cases/internal/error"
	"github.com/webitel/cases/model"
	"github.com/webitel/cases/util"
	"github.com/webitel/webitel-go-kit/etag"
	"log/slog"
	"strconv"
)

var CaseCommunicationMetadata = model.NewObjectMetadata(
	[]*model.Field{
		{"id", true},
		{"communication_type", true},
		{"communication_id", true},
	})

type CaseCommunicationService struct {
	app *App
	cases.UnimplementedCaseCommunicationsServer
}

func (c *CaseCommunicationService) ListCommunications(ctx context.Context, request *cases.ListCommunicationsRequest) (*cases.ListCommunicationsResponse, error) {
	// TODO: RBAC check by request.CaseEtag
	tag, err := etag.EtagOrId(etag.EtagCase, request.CaseId)
	if err != nil {
		return nil, errors.NewBadRequestError("app.case_communication.list_communication.invalid_etag", "Invalid case etag")
	}
	searchOpts := model.NewSearchOptions(ctx, request, CaseCommunicationMetadata)
	searchOpts.ParentId = tag.GetOid()

	res, dbErr := c.app.Store.CaseCommunication().List(searchOpts)
	if dbErr != nil {
		slog.Warn(dbErr.Error(), slog.Int64("id", tag.GetOid()))
		return nil, AppDatabaseError
	}
	err = NormalizeResponseCommunications(res.Data, request.GetFields())
	if err != nil {
		slog.Warn(err.Error(), slog.Int64("user_id", searchOpts.Session.GetUserId()), slog.Int64("domain_id", searchOpts.Session.GetDomainId()), slog.Int64("case_id", tag.GetOid()))
		return nil, AppResponseNormalizingError
	}
	return res, nil
}

func (c *CaseCommunicationService) LinkCommunication(ctx context.Context, request *cases.LinkCommunicationRequest) (*cases.LinkCommunicationResponse, error) {
	// TODO: RBAC check by request.CaseEtag
	if len(request.Input) == 0 {
		return nil, errors.NewBadRequestError("app.case_communication.link_communication.check_args.payload", "no payload")
	}
	err := ValidateCaseCommunicationsCreate(request.Input...)
	if err != nil {
		return nil, errors.NewBadRequestError("app.case_communication.link_communication.validate_payload.error", err.Error())
	}
	tag, err := etag.EtagOrId(etag.EtagCase, request.GetCaseId())
	if err != nil {
		return nil, errors.NewBadRequestError("app.case_communication.link_communication.invalid_etag", "Invalid case etag")
	}
	createOpts := model.NewCreateOptions(ctx, request, CaseCommunicationMetadata)
	createOpts.ParentID = tag.GetOid()

	res, dbErr := c.app.Store.CaseCommunication().Link(createOpts, request.Input)
	if dbErr != nil {
		slog.Warn(dbErr.Error(), slog.Int64("id", tag.GetOid()))
		return nil, AppDatabaseError
	}
	if len(res) == 0 {
		return nil, errors.NewBadRequestError("app.case_communication.link_communication.result.no_response", "no rows were affected (wrong ids or insufficient rights)")
	}
	err = NormalizeResponseCommunications(res, request.GetFields())
	if err != nil {
		slog.Warn(err.Error(), slog.Int64("user_id", createOpts.Session.GetUserId()), slog.Int64("domain_id", createOpts.Session.GetDomainId()), slog.Int64("case_id", tag.GetOid()))
		return nil, AppResponseNormalizingError
	}
	return &cases.LinkCommunicationResponse{Data: res}, nil
}

func (c *CaseCommunicationService) UnlinkCommunication(ctx context.Context, request *cases.UnlinkCommunicationRequest) (*cases.UnlinkCommunicationResponse, error) {
	// TODO: RBAC check by request.CaseEtag
	tag, err := etag.EtagOrId(etag.EtagCaseCommunication, request.GetId())
	if err != nil {
		return nil, errors.NewBadRequestError("app.case_communication.unlink_communication.invalid_etag", "Invalid case etag")
	}
	deleteOpts := model.NewDeleteOptions(ctx)
	deleteOpts.IDs = []int64{tag.GetOid()}

	affected, dbErr := c.app.Store.CaseCommunication().Unlink(deleteOpts)
	if dbErr != nil {
		slog.Warn(dbErr.Error(), slog.Int64("id", tag.GetOid()))
		return nil, AppDatabaseError
	}
	return &cases.UnlinkCommunicationResponse{Affected: affected}, nil
}

func NewCaseCommunicationService(app *App) (*CaseCommunicationService, errors.AppError) {
	return &CaseCommunicationService{app: app}, nil
}

func NormalizeResponseCommunications(res []*cases.CaseCommunication, requestedFields []string) error {
	if len(requestedFields) == 0 {
		requestedFields = CaseLinkMetadata.GetDefaultFields()
	}
	_, hasId, hasVer := util.FindEtagFields(requestedFields)
	for _, re := range res {
		if hasId {
			id, err := strconv.ParseInt(re.Id, 10, 64)
			if err != nil {
				return err
			}
			re.Id, err = etag.EncodeEtag(etag.EtagCaseCommunication, id, re.Ver)
			if err != nil {
				return err
			}
			if !hasVer {
				re.Ver = 0
			}
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
