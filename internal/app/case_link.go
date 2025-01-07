package app

import (
	"context"
	"strconv"

	cases "github.com/webitel/cases/api/cases"
	cerror "github.com/webitel/cases/internal/error"
	"github.com/webitel/cases/model"
	"github.com/webitel/cases/util"
	"github.com/webitel/webitel-go-kit/etag"
)

// In search options extract from context user
// Remove from search options fields functions

type CaseLinkService struct {
	app *App
	cases.UnimplementedCaseLinksServer
}

var CaseLinkMetadata = model.NewObjectMetadata(
	[]*model.Field{
		{Name: "id", Default: true},
		{Name: "ver", Default: true},
		{Name: "created_by", Default: true},
		{Name: "created_at", Default: true},
		{Name: "updated_by", Default: false},
		{Name: "updated_at", Default: false},
		{Name: "author", Default: true},
		{Name: "name", Default: true},
		{Name: "url", Default: true},
		{Name: "case_id", Default: false},
	})

func (c *CaseLinkService) LocateLink(ctx context.Context, req *cases.LocateLinkRequest) (*cases.CaseLink, error) {
	// Validate required fields
	if req.Id == "" {
		return nil, cerror.NewBadRequestError("app.case_link.locate.check_args.etag", "Etag is required")
	}

	etg, err := etag.EtagOrId(etag.EtagCaseLink, req.Id)
	if err != nil {
		return nil, cerror.NewBadRequestError("app.case_link.locate.parse_etag.error", err.Error())
	}

	searchOpts := model.NewLocateOptions(ctx, req, CaseLinkMetadata)
	searchOpts.IDs = []int64{etg.GetOid()}

	links, err := c.app.Store.CaseLink().List(searchOpts)
	if err != nil {
		return nil, cerror.NewInternalError("app.case_link.locate.get_list.error", err.Error())
	}
	if len(links.Items) == 0 {
		return nil, cerror.NewNotFoundError("app.case_link.locate.check_items.error", "not found")
	}
	res := links.Items[0]
	// hide etag if needed
	NormalizeResponseLink(res, req)

	return res, nil
}

func (c *CaseLinkService) CreateLink(ctx context.Context, req *cases.CreateLinkRequest) (*cases.CaseLink, error) {
	// Validate request
	if req.CaseId == "" {
		return nil, cerror.NewBadRequestError("app.case_link.create.case_etag.check_args.case_etag", "Case ID is required")
	} else if req.Input.GetUrl() == "" {
		return nil, cerror.NewBadRequestError("app.case_link.create.case_etag.check_args.url", "Url is required for each link")
	}
	caseTID, err := etag.EtagOrId(etag.EtagCase, req.CaseId)
	if err != nil {
		return nil, cerror.NewBadRequestError("app.case_link.create.case_etag.parse.error", err.Error())
	}

	createOpts := model.NewCreateOptions(ctx, req, CaseLinkMetadata)
	createOpts.ParentID = caseTID.GetOid()
	res, dbErr := c.app.Store.CaseLink().Create(createOpts, req.Input)
	if dbErr != nil {
		return nil, dbErr
	}

	NormalizeResponseLink(res, req)
	return res, nil
}

func (c *CaseLinkService) UpdateLink(ctx context.Context, req *cases.UpdateLinkRequest) (*cases.CaseLink, error) {
	if req.GetId() == "" {
		return nil, cerror.NewBadRequestError("app.case_link.update.check_args.id", "case ID required")
	}
	linkTID, err := etag.EtagOrId(etag.EtagCaseLink, req.GetId())
	if err != nil {
		return nil, cerror.NewBadRequestError("app.case_link.create.case_etag.parse.error", err.Error())
	}
	updateOpts := model.NewUpdateOptions(ctx, req, CaseLinkMetadata)
	updateOpts.Etags = []*etag.Tid{&linkTID}
	updated, err := c.app.Store.CaseLink().Update(updateOpts, req.Input)
	if err != nil {
		return nil, err
	}
	NormalizeResponseLink(updated, req)

	return updated, nil
}

func (c *CaseLinkService) DeleteLink(ctx context.Context, req *cases.DeleteLinkRequest) (*cases.CaseLink, error) {
	if req.GetId() == "" {
		return nil, cerror.NewBadRequestError("app.case_link.update.check_args.id", "case ID required")
	}
	linkTID, err := etag.EtagOrId(etag.EtagCaseLink, req.GetId())
	if err != nil {
		return nil, cerror.NewBadRequestError("app.case_link.create.case_etag.parse.error", err.Error())
	}
	deleteOpts := model.NewDeleteOptions(ctx)
	deleteOpts.ID = linkTID.GetOid()
	err = c.app.Store.CaseLink().Delete(deleteOpts)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (c *CaseLinkService) ListLinks(ctx context.Context, req *cases.ListLinksRequest) (*cases.CaseLinkList, error) {
	// Validate required fields
	if req.GetCaseId() == "" {
		return nil, cerror.NewBadRequestError("app.case_link.list.case_etag.check_args.id", "case ID is required")
	}

	etg, err := etag.EtagOrId(etag.EtagCase, req.GetCaseId())
	if err != nil {
		return nil, cerror.NewBadRequestError("app.case_link.locate.parse_etag.error", err.Error())
	}

	searchOpts := model.NewSearchOptions(ctx, req, CaseLinkMetadata)
	searchOpts.ParentId = etg.GetOid()
	//
	ids, err := util.ParseIds(req.GetIds(), etag.EtagCaseLink)
	if err != nil {
		return nil, cerror.NewBadRequestError("app.case_link.locate.parse_qin.invalid", err.Error())
	}
	searchOpts.IDs = ids

	links, err := c.app.Store.CaseLink().List(searchOpts)
	if err != nil {
		return nil, cerror.NewInternalError("case_comment_service.locate_comment.fetch_error", err.Error())
	}

	NormalizeResponseLinks(links, req.GetFields())

	// Return the located comment
	return links, nil
}

func NewCaseLinkService(app *App) (*CaseLinkService, cerror.AppError) {
	if app == nil {
		return nil, cerror.NewBadRequestError("app.case.new_case_comment_service.check_args.app", "unable to init service, app is nil")
	}
	return &CaseLinkService{app: app}, nil
}

func NormalizeResponseLink(res *cases.CaseLink, opts model.Fielder) {
	fields := opts.GetFields()
	if len(opts.GetFields()) == 0 {
		fields = CaseLinkMetadata.GetDefaultFields()
	}
	_, hasId, hasVer := util.FindEtagFields(fields)
	if hasId {
		id, _ := strconv.Atoi(res.Id)
		res.Id = etag.EncodeEtag(etag.EtagCaseLink, int64(id), res.Ver)
		// hide
		if !hasId {
			res.Id = ""
		}
		if !hasVer {
			res.Ver = 0
		}
	}
}

func NormalizeResponseLinks(res *cases.CaseLinkList, requestedFields []string) {
	fields := make([]string, len(requestedFields))
	copy(fields, requestedFields)
	if len(fields) == 0 {
		fields = CaseLinkMetadata.GetDefaultFields()
	}
	_, hasId, hasVer := util.FindEtagFields(fields)
	for _, re := range res.Items {
		if hasId {
			id, _ := strconv.Atoi(re.Id)
			re.Id = etag.EncodeEtag(etag.EtagCaseLink, int64(id), re.Ver)
			// hide
			if !hasId {
				re.Id = ""
			}
			if !hasVer {
				re.Ver = 0
			}
		}
	}
}
