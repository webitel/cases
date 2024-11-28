package app

import (
	"context"
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
		{"etag", true},
		{"created_by", true},
		{"created_at", true},
		{"updated_by", false},
		{"updated_at", false},
		{"author", true},
		{"name", true},
		{"url", true},
		{"case_id", true},
	})

func (c *CaseLinkService) LocateLink(ctx context.Context, req *cases.LocateLinkRequest) (*cases.CaseLink, error) {
	// Validate required fields
	if req.Etag == "" {
		return nil, cerror.NewBadRequestError("app.case_link.locate.check_args.etag", "Etag is required")
	}

	etg, err := etag.EtagOrId(etag.EtagCaseLink, req.Etag)
	if err != nil {
		return nil, cerror.NewBadRequestError("app.case_link.locate.parse_etag.error", err.Error())
	}

	searchOpts := model.NewLocateOptions(ctx, req, CaseLinkMetadata)
	searchOpts.IDs = []int64{etg.GetOid()}

	links, err := c.app.Store.CaseLink().List(searchOpts)
	if err != nil {
		return nil, cerror.NewInternalError("app.case_link.locate.get_list.error", err.Error())
	}
	res := links.Items[0]
	// hide etag if needed
	NormalizeResponseLink(res, req)

	//Return the located comment
	return links.Items[0], nil
}

func (c *CaseLinkService) CreateLink(ctx context.Context, req *cases.CreateLinkRequest) (*cases.CaseLink, error) {

	// Validate request
	if req.CaseEtag == "" {
		return nil, cerror.NewBadRequestError("app.case_link.create.case_etag.check_args.case_etag", "Case etag is required")
	} else if req.Input.GetUrl() == "" {
		return nil, cerror.NewBadRequestError("app.case_link.create.case_etag.check_args.url", "Url is required for each link")
	}
	caseTID, err := etag.EtagOrId(etag.EtagCaseLink, req.CaseEtag)
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
	if req.GetEtag() == "" {
		return nil, cerror.NewBadRequestError("app.case_link.update.check_args.etag", "case etag required")
	}
	linkTID, err := etag.EtagOrId(etag.EtagCaseLink, req.GetEtag())
	if err != nil {
		return nil, cerror.NewBadRequestError("app.case_link.create.case_etag.parse.error", err.Error())
	}
	updateOpts := model.NewUpdateOptions(ctx, req)
	updateOpts.Etags = []*etag.Tid{&linkTID}
	updated, err := c.app.Store.CaseLink().Update(updateOpts, req.Input)
	if err != nil {
		return nil, err
	}
	NormalizeResponseLink(updated, req)

	return updated, nil
}

func (c *CaseLinkService) DeleteLink(ctx context.Context, req *cases.DeleteLinkRequest) (*cases.CaseLink, error) {
	if req.GetEtag() == "" {
		return nil, cerror.NewBadRequestError("app.case_link.update.check_args.etag", "case etag required")
	}
	linkTID, err := etag.EtagOrId(etag.EtagCaseLink, req.GetEtag())
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
	if req.GetCaseEtag() == "" {
		return nil, cerror.NewBadRequestError("app.case_link.list.case_etag.check_args.etag", "case etag is required")
	}

	etg, err := etag.EtagOrId(etag.EtagCase, req.GetCaseEtag())
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

	NormalizeResponseLinks(links, req)

	//Return the located comment
	return links, nil
}

func NewCaseLinkService(app *App) (*CaseLinkService, cerror.AppError) {
	if app == nil {
		return nil, cerror.NewBadRequestError("app.case.new_case_comment_service.check_args.app", "unable to init service, app is nil")
	}
	return &CaseLinkService{app: app}, nil
}

func NormalizeResponseLink(res *cases.CaseLink, opts model.Locator) {
	fields := opts.GetFields()
	if len(opts.GetFields()) == 0 {
		fields = CaseLinkMetadata.GetDefaultFields()
	}
	hasEtag, hasId, hasVer := util.FindEtagFields(fields)
	if hasEtag {
		res.Etag = etag.EncodeEtag(etag.EtagCaseLink, res.Id, res.Ver)
		// hide
		if !hasId {
			res.Id = 0
		}
		if !hasVer {
			res.Ver = 0
		}
	}
}

func NormalizeResponseLinks(res *cases.CaseLinkList, opts model.Locator) {
	fields := opts.GetFields()
	if len(fields) == 0 {
		fields = CaseLinkMetadata.GetDefaultFields()
	}
	hasEtag, hasId, hasVer := util.FindEtagFields(fields)
	for _, re := range res.Items {
		if hasEtag {
			re.Etag = etag.EncodeEtag(etag.EtagCaseLink, re.Id, re.Ver)
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
