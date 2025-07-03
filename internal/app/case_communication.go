package app

import (
	"context"
	defErr "errors"
	"fmt"
	"github.com/webitel/cases/api/cases"
	"github.com/webitel/cases/auth"
	"github.com/webitel/cases/internal/errors"
	"github.com/webitel/cases/internal/model"
	grpcopts "github.com/webitel/cases/internal/model/options/grpc"
	"github.com/webitel/cases/util"
	"github.com/webitel/webitel-go-kit/pkg/etag"
)

var CaseCommunicationMetadata = model.NewObjectMetadata("", caseObjScope, []*model.Field{
	{Name: "etag", Default: true},
	{Name: "ver", Default: false},
	{Name: "id", Default: true},
	{Name: "communication_type", Default: true},
	{Name: "communication_id", Default: true},
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
		return nil, err
	}

	tag, err := etag.EtagOrId(etag.EtagCase, request.GetCaseEtag())
	if err != nil {
		return nil, errors.InvalidArgument(
			"Invalid case etag",
		)
	}
	if tag.GetOid() != 0 {
		searchOpts.AddFilter(util.EqualFilter("case_id=%d", tag.GetOid()))
	}
	if searchOpts.GetAuthOpts().GetObjectScope(CaseCommunicationMetadata.GetParentScopeName()).IsRbacUsed() {
		access, err := c.app.Store.Case().CheckRbacAccess(searchOpts, searchOpts.GetAuthOpts(), auth.Read, tag.GetOid())
		if err != nil {
			return nil, err
		}
		if !access {
			return nil, errors.Forbidden("user doesn't have required (READ) access to the case")
		}
	}

	res, err := c.app.Store.CaseCommunication().List(searchOpts)
	if err != nil {
		return nil, err
	}
	err = NormalizeResponseCommunications(res.Data, request.GetFields())
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *CaseCommunicationService) LinkCommunication(ctx context.Context, request *cases.LinkCommunicationRequest) (*cases.LinkCommunicationResponse, error) {
	err := ValidateCaseCommunicationsCreate(request.Input)
	if err != nil {
		return nil, err
	}
	tag, err := etag.EtagOrId(etag.EtagCase, request.GetCaseEtag())
	if err != nil {
		return nil, errors.InvalidArgument("Invalid case etag", errors.WithCause(err))
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
		return nil, err
	}
	accessMode := auth.Edit
	if !createOpts.GetAuthOpts().CheckObacAccess(CaseCommunicationMetadata.GetParentScopeName(), accessMode) {
		return nil, errors.Forbidden("user doesn't have required (EDIT) access to the case")
	}
	if createOpts.GetAuthOpts().GetObjectScope(CaseCommunicationMetadata.GetParentScopeName()).IsRbacUsed() {
		access, err := c.app.Store.Case().CheckRbacAccess(createOpts, createOpts.GetAuthOpts(), accessMode, createOpts.ParentID)
		if err != nil {
			return nil, errors.Forbidden("user doesn't have required (EDIT) access to the case", errors.WithCause(err))
		}
		if !access {
			return nil, errors.Forbidden("user doesn't have required (EDIT) access to the case")
		}
	}

	res, err := c.app.Store.CaseCommunication().Link(createOpts, []*cases.InputCaseCommunication{request.Input})
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, errors.InvalidArgument("no rows were affected (wrong ids or insufficient rights)")
	}
	err = NormalizeResponseCommunications(res, request.GetFields())
	if err != nil {
		return nil, err
	}
	return &cases.LinkCommunicationResponse{Data: res}, nil
}

func (c *CaseCommunicationService) UnlinkCommunication(ctx context.Context, request *cases.UnlinkCommunicationRequest) (*cases.UnlinkCommunicationResponse, error) {
	tag, err := etag.EtagOrId(etag.EtagCaseCommunication, request.GetId())
	if err != nil {
		return nil, errors.InvalidArgument("Invalid communication etag")
	}
	caseTag, err := etag.EtagOrId(etag.EtagCase, request.GetCaseEtag())
	if err != nil {
		return nil, errors.InvalidArgument("Invalid case etag", errors.WithCause(err))
	}
	deleteOpts, err := grpcopts.NewDeleteOptions(ctx, grpcopts.WithDeleteID(tag.GetOid()))
	if err != nil {
		return nil, err
	}
	deleteOpts.IDs = []int64{tag.GetOid()}
	if deleteOpts.GetAuthOpts().GetObjectScope(CaseCommunicationMetadata.GetParentScopeName()).IsRbacUsed() {
		access, err := c.app.Store.Case().CheckRbacAccess(deleteOpts, deleteOpts.GetAuthOpts(), auth.Edit, caseTag.GetOid())
		if err != nil {
			return nil, err
		}
		if !access {
			return nil, errors.Forbidden("user doesn't have required (EDIT) access to the case")
		}
	}

	affected, err := c.app.Store.CaseCommunication().Unlink(deleteOpts)
	if err != nil {
		return nil, err
	}
	return &cases.UnlinkCommunicationResponse{Affected: affected}, nil
}

func NewCaseCommunicationService(app *App) (*CaseCommunicationService, error) {
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
