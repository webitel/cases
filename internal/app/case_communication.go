package app

import (
	"context"
	"github.com/webitel/cases/api/cases"
	authmodel "github.com/webitel/cases/auth/model"
	errors "github.com/webitel/cases/internal/error"
	"github.com/webitel/cases/model"
	"github.com/webitel/cases/util"
	"github.com/webitel/webitel-go-kit/etag"
)

const (
	chatsScopeName  = "chats"
	callsScopeName  = "calls"
	emailsScopeName = "emails"
)

var CaseCommunicationMetadata = model.NewObjectMetadata(
	[]*model.Field{
		{"etag", true},
		{"communication_type", true},
		{"communication_id", true},
	})

type CaseCommunicationService struct {
	app *App
	cases.UnimplementedCaseCommunicationsServer
}

func (c *CaseCommunicationService) LinkCommunication(ctx context.Context, request *cases.LinkCommunicationRequest) (*cases.LinkCommunicationResponse, error) {
	if len(request.Input) == 0 {
		return nil, errors.NewBadRequestError("app.case_communication.link_communication.check_args.payload", "no payload")
	}
	tag, err := etag.EtagOrId(etag.EtagCase, request.CaseEtag)
	if err != nil {
		return nil, errors.NewBadRequestError("app.case_communication.link_communication.invalid_etag", "Invalid case etag")
	}
	createOpts := model.NewCreateOptions(ctx, request, CaseCommunicationMetadata)
	createOpts.ParentID = tag.GetOid()

	// TODO: how to check permissions on each object?
	if session := createOpts.Session; session != nil {
		checkRbacEnabled := func(scope string) bool {
			s := session.GetScope(scope)
			if s == nil {
				return false
			}
			return s.Rbac
		}
		chatsRbac := checkRbacEnabled(chatsScopeName)
		callsRbac := checkRbacEnabled(callsScopeName)
		emailsRbac := checkRbacEnabled(emailsScopeName)
		for _, communication := range request.Input {
			var (
				rbacEnabled  bool
				scopeName    string
				apiCheckFunc func() error
			)
			switch communication.CommunicationType {
			case cases.CaseCommunicationsTypes_COMMUNICATION_CALL:
				rbacEnabled = callsRbac
				scopeName = callsScopeName
				apiCheckFunc = func() error {
					return nil
				}
			case cases.CaseCommunicationsTypes_COMMUNICATION_CHAT:
				rbacEnabled = chatsRbac
				scopeName = chatsScopeName
				apiCheckFunc = func() error {
					return nil
				}
			case cases.CaseCommunicationsTypes_COMMUNICATION_EMAIL:
				rbacEnabled = emailsRbac
				scopeName = emailsScopeName
				apiCheckFunc = func() error {
					return nil
				}
			}
			// check if rbac enabled
			if rbacEnabled {
				// call api to check access to the object
				err := apiCheckFunc()
				if err != nil {
					return nil, errors.NewForbiddenError("app.case_communication.link_communication.permissions.check", err.Error())
				}
			} else if !session.HasObacAccess(scopeName, authmodel.Read) { // if rbac disabled just check access to object
				return nil, errors.NewBadRequestError("app.case_communication.link_communication.invalid_etag", "Invalid case etag")
			}
		}
	}

	res, dbErr := c.app.Store.CaseCommunication().Link(createOpts, request.Input)
	if dbErr != nil {
		return nil, dbErr
	}
	NormalizeResponseCommunications(res, request.GetFields())
	return &cases.LinkCommunicationResponse{Data: res}, nil
}

func (c *CaseCommunicationService) UnlinkCommunication(ctx context.Context, request *cases.UnlinkCommunicationRequest) (*cases.UnlinkCommunicationResponse, error) {
	tag, err := etag.EtagOrId(etag.EtagCaseCommunication, request.Etag)
	if err != nil {
		return nil, errors.NewBadRequestError("app.case_communication.unlink_communication.invalid_etag", "Invalid case etag")
	}
	deleteOpts := model.NewDeleteOptions(ctx)
	deleteOpts.IDs = []int64{tag.GetOid()}

	res, dbErr := c.app.Store.CaseCommunication().Unlink(deleteOpts)
	if dbErr != nil {
		return nil, dbErr
	}
	NormalizeResponseCommunications(res, request.GetFields())
	if len(res) == 0 {
		return nil, errors.NewBadRequestError("app.case_communication.unlink_communication.no_rows_affected", "No rows were affected while deleting")
	}
	return &cases.UnlinkCommunicationResponse{Data: res[0]}, nil
}

func NewCaseCommunicationService(app *App) (*CaseCommunicationService, errors.AppError) {
	return &CaseCommunicationService{app: app}, nil
}

func NormalizeResponseCommunications(res []*cases.CaseCommunication, requestedFields []string) {
	fields := make([]string, len(requestedFields))
	copy(fields, requestedFields)
	if len(fields) == 0 {
		fields = CaseLinkMetadata.GetDefaultFields()
	}
	hasEtag, hasId, hasVer := util.FindEtagFields(fields)
	for _, re := range res {
		if hasEtag {
			re.Etag = etag.EncodeEtag(etag.EtagCaseCommunication, re.Id, re.Ver)
			// hide
			if !hasId {
				re.Id = 0
			}
			if !hasVer {
				re.Ver = 0
			}
		}
	}
}
