package grpc

import (
	"context"

	"github.com/webitel/cases/api/cases"
	"github.com/webitel/cases/internal/api_handler/grpc/utils"
	"github.com/webitel/cases/internal/errors"
	"github.com/webitel/cases/internal/model"
	"github.com/webitel/cases/internal/model/options"
	grpcopts "github.com/webitel/cases/internal/model/options/grpc"
	"github.com/webitel/cases/internal/model/options/grpc/shared"
	"github.com/webitel/cases/util"
	"github.com/webitel/webitel-go-kit/pkg/etag"
)

type CaseLinkHandler interface {
	ListCaseLinks(options.Searcher) ([]*model.CaseLink, error)
	CreateCaseLink(options.Creator, *model.CaseLink) (*model.CaseLink, error)
	UpdateCaseLink(options.Updator, *model.CaseLink) (*model.CaseLink, error)
	DeleteCaseLink(options.Deleter) (*model.CaseLink, error)
}

type CaseLinkService struct {
	app CaseLinkHandler
	cases.UnimplementedCaseLinksServer
}

func NewCaseLinkService(handler CaseLinkHandler) *CaseLinkService {
	return &CaseLinkService{app: handler}
}

var CaseLinkMetadata = model.NewObjectMetadata("", model.ScopeCases, []*model.Field{
	{Name: "etag", Default: true},
	{Name: "id", Default: false},
	{Name: "ver", Default: false},
	{Name: "created_by", Default: true},
	{Name: "created_at", Default: true},
	{Name: "updated_by", Default: false},
	{Name: "updated_at", Default: false},
	{Name: "author", Default: true},
	{Name: "name", Default: true},
	{Name: "url", Default: true},
	{Name: "case_id", Default: false},
})

func (s *CaseLinkService) ListLinks(ctx context.Context, req *cases.ListLinksRequest) (*cases.CaseLinkList, error) {
	if req.GetCaseEtag() == "" {
		return nil, errors.InvalidArgument("case etag is required")
	}
	searchOpts, err := grpcopts.NewSearchOptions(
		ctx,
		grpcopts.WithSearch(req),
		grpcopts.WithPagination(req),
		grpcopts.WithFields(req, CaseLinkMetadata,
			util.DeduplicateFields,
			util.ParseFieldsForEtag,
			util.EnsureIdField,
		),
		grpcopts.WithIDsAsEtags(etag.EtagCaseLink, req.GetIds()...),
		grpcopts.WithSort(req),
	)
	if err != nil {
		return nil, err
	}
	etg, err := etag.EtagOrId(etag.EtagCase, req.GetCaseEtag())
	if err != nil {
		return nil, errors.InvalidArgument("invalid etag", errors.WithCause(err))
	}
	searchOpts.AddFilter("case_id", etg.GetOid())

	links, err := s.app.ListCaseLinks(searchOpts)
	if err != nil {
		return nil, err
	}

	var res cases.CaseLinkList
	converted, err := utils.ConvertToOutputBulk(links, s.Marshal)
	if err != nil {
		return nil, err
	}
	res.Next, res.Items = utils.GetListResult(searchOpts, converted)
	res.Page = int64(req.GetPage())

	// Normalize response
	if err := NormalizeResponseLinks(&res, req.GetFields()); err != nil {
		return nil, err
	}

	return &res, nil
}

func (s *CaseLinkService) CreateLink(ctx context.Context, req *cases.CreateLinkRequest) (*cases.CaseLink, error) {
	if req.CaseEtag == "" {
		return nil, errors.InvalidArgument("case etag is required")
	}
	if req.Input == nil || req.Input.Url == "" {
		return nil, errors.InvalidArgument("url is required")
	}
	caseTid, err := etag.EtagOrId(etag.EtagCase, req.CaseEtag)
	if err != nil {
		return nil, errors.InvalidArgument("invalid etag", errors.WithCause(err))
	}
	createOpts, err := grpcopts.NewCreateOptions(
		ctx,
		grpcopts.WithCreateFields(req, CaseLinkMetadata,
			util.DeduplicateFields,
			util.ParseFieldsForEtag,
			util.EnsureIdField,
		),
		grpcopts.WithCreateParentID(caseTid.GetOid()),
	)
	if err != nil {
		return nil, err
	}

	input := &model.CaseLink{
		Name: &req.Input.Name,
		Url:  req.Input.Url,
	}

	m, err := s.app.CreateCaseLink(createOpts, input)
	if err != nil {
		return nil, err
	}
	out, err := s.Marshal(m)
	if err != nil {
		return nil, err
	}
	if err := NormalizeResponseLink(out, req); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *CaseLinkService) UpdateLink(ctx context.Context, req *cases.UpdateLinkRequest) (*cases.CaseLink, error) {
	// Validate input
	if req.Input == nil {
		return nil, errors.InvalidArgument("input required")
	}
	if req.Input.Etag == "" {
		return nil, errors.InvalidArgument("link etag is required")
	}

	// Decode etags to numeric IDs
	linkTid, err := etag.EtagOrId(etag.EtagCaseLink, req.Input.Etag)
	if err != nil {
		return nil, errors.InvalidArgument("invalid link etag", errors.WithCause(err))
	}
	caseTid, err := etag.EtagOrId(etag.EtagCase, req.CaseEtag)
	if err != nil {
		return nil, errors.InvalidArgument("invalid case etag", errors.WithCause(err))
	}

	// Build update options with both IDs
	updateOpts, err := grpcopts.NewUpdateOptions(
		ctx,
		grpcopts.WithUpdateFields(req, CaseLinkMetadata),
		grpcopts.WithUpdateParentID(caseTid.GetOid()),
		grpcopts.WithUpdateEtag(&linkTid),
		grpcopts.WithUpdateMasker(req),
	)
	if err != nil {
		return nil, err
	}

	// Prepare input model
	input := &model.CaseLink{
		Name: &req.Input.Name,
		Url:  req.Input.Url,
	}

	// Call business logic
	m, err := s.app.UpdateCaseLink(updateOpts, input)
	if err != nil {
		return nil, err
	}

	// Marshal and normalize response
	out, err := s.Marshal(m)
	if err != nil {
		return nil, err
	}
	if err := NormalizeResponseLink(out, req); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *CaseLinkService) DeleteLink(ctx context.Context, req *cases.DeleteLinkRequest) (*cases.CaseLink, error) {
	if req.Etag == "" {
		return nil, errors.InvalidArgument("etag is required")
	}
	linkTID, err := etag.EtagOrId(etag.EtagCaseLink, req.GetEtag())
	if err != nil {
		return nil, err
	}
	deleteOpts, err := grpcopts.NewDeleteOptions(ctx, grpcopts.WithDeleteID(linkTID.GetOid()), grpcopts.WithDeleteParentIDAsEtag(etag.EtagCase, req.GetCaseEtag()))
	if err != nil {
		return nil, err
	}

	m, err := s.app.DeleteCaseLink(deleteOpts)
	if err != nil {
		return nil, err
	}
	out, err := s.Marshal(m)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (s *CaseLinkService) LocateLink(
	ctx context.Context,
	req *cases.LocateLinkRequest,
) (*cases.CaseLink, error) {
	// Build search options using the request (similar to ListLinks)
	searchOpts, err := grpcopts.NewLocateOptions(
		ctx,
		grpcopts.WithFields(req, CaseLinkMetadata,
			util.DeduplicateFields,
			util.EnsureIdField,
			util.ParseFieldsForEtag,
		),
		grpcopts.WithIDsAsEtags(etag.EtagCaseLink, req.GetEtag()),
	)
	if err != nil {
		return nil, err
	}

	// Add parent case filter if provided
	if req.GetCaseEtag() != "" {
		caseEtg, err := etag.EtagOrId(etag.EtagCase, req.GetCaseEtag())
		if err != nil {
			return nil, errors.InvalidArgument("invalid case etag", errors.WithCause(err))
		}
		searchOpts.AddFilter("case_id", caseEtg.GetOid())
	}

	// Call business logic
	links, err := s.app.ListCaseLinks(searchOpts)
	if err != nil {
		return nil, err
	}
	if len(links) == 0 {
		return nil, errors.NotFound("not found")
	}
	if len(links) > 1 {
		return nil, errors.InvalidArgument("too many items found")
	}

	out, err := s.Marshal(links[0])
	if err != nil {
		return nil, err
	}
	if err := NormalizeResponseLink(out, req); err != nil {
		return nil, err
	}
	return out, nil
}

// Marshal converts a model.CaseLink to cases.CaseLink
func (s *CaseLinkService) Marshal(m *model.CaseLink) (*cases.CaseLink, error) {
	if m == nil {
		return nil, nil
	}
	return &cases.CaseLink{
		Id:        m.Id,
		Ver:       m.Ver,
		Etag:      m.Etag,
		CreatedBy: utils.MarshalLookup(m.Author),
		CreatedAt: utils.MarshalTime(m.CreatedAt),
		UpdatedBy: utils.MarshalLookup(m.Editor),
		UpdatedAt: utils.MarshalTime(m.UpdatedAt),
		Author:    utils.MarshalLookup(m.Contact),
		Name:      utils.Dereference(m.Name),
		Url:       m.Url,
	}, nil
}

func NormalizeResponseLink(res *cases.CaseLink, opts shared.Fielder) error {
	var err error
	hasEtag, hasId, hasVer := util.FindEtagFields(opts.GetFields())
	if hasEtag {
		res.Etag, err = etag.EncodeEtag(etag.EtagCaseLink, res.GetId(), res.GetVer())
		if err != nil {
			return err
		}

		// hide
		if !hasId {
			res.Id = 0
		}
		if !hasVer {
			res.Ver = 0
		}
	}
	return nil
}

func NormalizeResponseLinks(res *cases.CaseLinkList, requestedFields []string) error {

	if len(requestedFields) == 0 {
		requestedFields = CaseLinkMetadata.GetDefaultFields()
	}
	var err error
	hasEtag, hasId, hasVer := util.FindEtagFields(requestedFields)
	for _, re := range res.Items {
		if hasEtag {
			re.Etag, err = etag.EncodeEtag(etag.EtagCaseLink, re.Id, re.Ver)
			if err != nil {
				return err
			}
			// hide
			if !hasId {
				re.Id = 0
			}
			if !hasVer {
				re.Ver = 0
			}
		}
	}
	return nil
}
