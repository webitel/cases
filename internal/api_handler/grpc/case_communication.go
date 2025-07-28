package grpc

import (
	"context"
	defErr "errors"
	"fmt"

	"google.golang.org/grpc/codes"

	"github.com/webitel/webitel-go-kit/pkg/etag"

	api "github.com/webitel/cases/api/cases"
	"github.com/webitel/cases/internal/api_handler/grpc/utils"
	"github.com/webitel/cases/internal/errors"
	"github.com/webitel/cases/internal/model"
	"github.com/webitel/cases/internal/model/options"
	grpcopts "github.com/webitel/cases/internal/model/options/grpc"
	"github.com/webitel/cases/util"
)

type CaseCommunicationHandler interface {
	ListCommunications(options.Searcher) ([]*model.CaseCommunication, error)
	LinkCommunication(createOpts options.Creator, input []*model.CaseCommunication) ([]*model.CaseCommunication, error)
	UnlinkCommunication(deleteOpts options.Deleter) (int64, error)
}

type CaseCommunicationService struct {
	api.UnimplementedCaseCommunicationsServer

	app CaseCommunicationHandler
}

var CaseCommunicationMetadata = model.NewObjectMetadata("", "cases", []*model.Field{
	{Name: "etag", Default: true},
	{Name: "ver", Default: false},
	{Name: "id", Default: true},
	{Name: "communication_type", Default: true},
	{Name: "communication_id", Default: true},
})

func NewCaseCommunicationService(app CaseCommunicationHandler) (*CaseCommunicationService, error) {
	if app == nil {
		return nil, errors.New("case communication handler is nil", errors.WithCode(codes.InvalidArgument))
	}
	return &CaseCommunicationService{app: app}, nil
}

func (s *CaseCommunicationService) ListCommunications(ctx context.Context, req *api.ListCommunicationsRequest) (*api.ListCommunicationsResponse, error) {
	searchOpts, err := grpcopts.NewSearchOptions(
		ctx,
		grpcopts.WithSearch(req),
		grpcopts.WithPagination(req),
		grpcopts.WithFields(req, CaseCommunicationMetadata,
			util.DeduplicateFields,
			util.ParseFieldsForEtag,
			util.EnsureIdField,
		),
		grpcopts.WithSort(req),
	)
	if err != nil {
		return nil, err
	}
	// Parse etag and add case_id filter
	if req.GetCaseEtag() != "" {
		caseTid, err := etag.EtagOrId(etag.EtagCase, req.GetCaseEtag())
		if err != nil {
			return nil, errors.InvalidArgument("invalid case etag", errors.WithCause(err))
		}

		searchOpts.AddFilter(util.EqualFilter("case_id", caseTid.GetOid()))
	}

	items, err := s.app.ListCommunications(searchOpts)
	if err != nil {
		return nil, err
	}

	var out []*api.CaseCommunication

	out, err = utils.ConvertToOutputBulk(items, s.Marshal)
	if err != nil {
		return nil, err
	}
	// Normalize response
	if err := NormalizeResponseCommunications(out, req.GetFields()); err != nil {
		return nil, err
	}
	return &api.ListCommunicationsResponse{
		Data: out,
		Page: req.GetPage(),
		Next: false, // TODO: handle pagination
	}, nil
}

func (s *CaseCommunicationService) LinkCommunication(ctx context.Context, req *api.LinkCommunicationRequest) (*api.LinkCommunicationResponse, error) {
	input := []*model.CaseCommunication{unmarshalInputCaseCommunication(req.GetInput())}
	// Validate input
	if err := ValidateCaseCommunicationsCreate(input...); err != nil {
		return nil, err
	}
	// Parse case etag and get parent ID
	tag, err := etag.EtagOrId(etag.EtagCase, req.GetCaseEtag())
	if err != nil {
		return nil, errors.InvalidArgument("invalid case etag", errors.WithCause(err))
	}

	createOpts, err := grpcopts.NewCreateOptions(
		ctx,
		grpcopts.WithCreateFields(req, CaseCommunicationMetadata,
			util.DeduplicateFields,
			util.ParseFieldsForEtag,
			util.EnsureIdField,
		),
		grpcopts.WithCreateParentID(tag.GetOid()),
	)
	if err != nil {
		return nil, err
	}

	items, err := s.app.LinkCommunication(createOpts, input)
	if err != nil {
		return nil, err
	}

	var out []*api.CaseCommunication
	out, err = utils.ConvertToOutputBulk(items, s.Marshal)
	if err != nil {
		return nil, err
	}
	// Normalize response
	if err := NormalizeResponseCommunications(out, createOpts.GetFields()); err != nil {
		return nil, err
	}
	return &api.LinkCommunicationResponse{Data: out}, nil
}

func (s *CaseCommunicationService) UnlinkCommunication(ctx context.Context, req *api.UnlinkCommunicationRequest) (*api.UnlinkCommunicationResponse, error) {
	// Validate communication etag
	commTag, err := etag.EtagOrId(etag.EtagCaseCommunication, req.GetId())
	if err != nil {
		return nil, errors.InvalidArgument("invalid communication etag", errors.WithCause(err))
	}
	// Validate case etag
	caseTag, err := etag.EtagOrId(etag.EtagCase, req.GetCaseEtag())
	if err != nil {
		return nil, errors.InvalidArgument("invalid case etag", errors.WithCause(err))
	}

	deleteOpts, err := grpcopts.NewDeleteOptions(ctx, grpcopts.WithDeleteID(commTag.GetOid()))
	if err != nil {
		return nil, err
	}

	deleteOpts.IDs = []int64{commTag.GetOid()}

	deleteOpts.AddFilter(util.EqualFilter("case_id", caseTag.GetOid()))

	affected, err := s.app.UnlinkCommunication(deleteOpts)
	if err != nil {
		return nil, err
	}
	return &api.UnlinkCommunicationResponse{Affected: affected}, nil
}

func unmarshalInputCaseCommunication(input *api.InputCaseCommunication) *model.CaseCommunication {
	out := &model.CaseCommunication{
		CommunicationType: utils.UnmarshalLookup(input.GetCommunicationType(), &model.GeneralLookup{}),
		CommunicationId:   input.GetCommunicationId(),
	}
	return out
}

func (s *CaseCommunicationService) Marshal(m *model.CaseCommunication) (*api.CaseCommunication, error) {
	etg, err := etag.EncodeEtag(etag.EtagCaseCommunication, m.Id, m.Ver)
	if err != nil {
		return nil, err
	}
	out := &api.CaseCommunication{
		Id:                m.Id,
		Ver:               m.Ver,
		Etag:              etg,
		CommunicationType: utils.MarshalLookup(m.CommunicationType),
		CommunicationId:   m.CommunicationId,
	}
	return out, nil
}

func ValidateCaseCommunicationsCreate(input ...*model.CaseCommunication) error {
	var errorsSlice []error

	for i, communication := range input {
		if communication == nil {
			errorsSlice = append(errorsSlice, errors.InvalidArgument(fmt.Sprintf("input[%d]: empty entry", i)))

			continue
		}
		if communication.CommunicationId == "" {
			errorsSlice = append(errorsSlice, errors.InvalidArgument(fmt.Sprintf("input[%d]: communication can't be empty", i)))
		}
		if communication.CommunicationType == nil {
			errorsSlice = append(errorsSlice, errors.InvalidArgument(fmt.Sprintf("input[%d]: communication type can't be empty", i)))
		}
	}
	if errorsSlice != nil {
		return errors.InvalidArgument("invalid input", errors.WithCause(defErr.Join(errorsSlice...)))
	}
	return nil
}

func NormalizeResponseCommunications(res []*api.CaseCommunication, requestedFields []string) error {
	if len(requestedFields) == 0 {
		requestedFields = CaseCommunicationMetadata.GetDefaultFields()
	}
	hasEtag, hasID, hasVer := util.FindEtagFields(requestedFields)
	for _, re := range res {
		err := util.NormalizeEtags(etag.EtagCaseCommunication, hasEtag, hasID, hasVer, &re.Etag, &re.Id, &re.Ver)
		if err != nil {
			return err
		}
	}
	return nil
}
