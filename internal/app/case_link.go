package app

import (
	"context"
	cases "github.com/webitel/cases/api/cases"
	authmodel "github.com/webitel/cases/auth/user_auth"
	cerror "github.com/webitel/cases/internal/error"
	"github.com/webitel/cases/model"
	"github.com/webitel/cases/util"
	"github.com/webitel/webitel-go-kit/etag"
	"log/slog"
)

// In search options extract from context user
// Remove from search options fields functions

type CaseLinkService struct {
	app *App
	cases.UnimplementedCaseLinksServer
}

var CaseLinkMetadata = model.NewObjectMetadata("", caseObjScope, []*model.Field{
	{"etag", true},
	{"id", false},
	{"ver", false},
	{"created_by", true},
	{"created_at", true},
	{"updated_by", false},
	{"updated_at", false},
	{"author", true},
	{"name", true},
	{"url", true},
	{"case_id", false},
})

func (c *CaseLinkService) LocateLink(ctx context.Context, req *cases.LocateLinkRequest) (*cases.CaseLink, error) {
	// Validate required fields
	if req.Etag == "" {
		return nil, cerror.NewBadRequestError("app.case_link.locate.check_args.etag", "Etag is required")
	}

	etg, err := etag.EtagOrId(etag.EtagCaseLink, req.GetEtag())
	if err != nil {
		return nil, cerror.NewBadRequestError("app.case_link.locate.parse_etag.error", err.Error())
	}

	caseEtg, err := etag.EtagOrId(etag.EtagCase, req.GetCaseEtag())
	if err != nil {
		return nil, cerror.NewBadRequestError("app.case_link.locate.parse_case_etag.error", err.Error())
	}

	searchOpts, err := model.NewLocateOptions(ctx, req, CaseLinkMetadata)
	if err != nil {
		slog.Error(err.Error())
		return nil, AppInternalError
	}
	searchOpts.IDs = []int64{etg.GetOid()}
	searchOpts.ParentId = caseEtg.GetOid()
	logAttributes := slog.Group(
		"context",
		slog.Int64("user_id", searchOpts.GetAuthOpts().GetUserId()),
		slog.Int64("domain_id", searchOpts.GetAuthOpts().GetDomainId()),
		slog.Int64("id", searchOpts.IDs[0]),
		slog.Int64("case_id", searchOpts.ParentId),
	)

	if searchOpts.GetAuthOpts().GetObjectScope(CaseLinkMetadata.GetMainScopeName()).IsRbacUsed() {
		access, err := c.app.Store.Case().CheckRbacAccess(searchOpts, searchOpts.GetAuthOpts(), authmodel.Read, searchOpts.ParentId)
		if err != nil {
			slog.Error(err.Error(), logAttributes)
			return nil, AppForbiddenError
		}
		if !access {
			slog.Error("user doesn't have required (READ) access to the case", logAttributes)
			return nil, AppForbiddenError
		}
	}
	links, err := c.app.Store.CaseLink().List(searchOpts)
	if err != nil {
		slog.Error(err.Error(), logAttributes)
		return nil, AppDatabaseError
	}
	if len(links.Items) == 0 {
		return nil, cerror.NewNotFoundError("app.case_link.locate.check_items.error", "not found")
	}
	res := links.Items[0]
	// hide etag if needed
	err = NormalizeResponseLink(res, req)
	if err != nil {
		slog.Error(err.Error(), logAttributes)
		return nil, AppResponseNormalizingError
	}
	return res, nil
}

func (c *CaseLinkService) CreateLink(ctx context.Context, req *cases.CreateLinkRequest) (*cases.CaseLink, error) {

	// Validate request
	if req.CaseEtag == "" {
		return nil, cerror.NewBadRequestError("app.case_link.create.case_etag.check_args.case_etag", "Case etag is required")
	} else if req.Input.GetUrl() == "" {
		return nil, cerror.NewBadRequestError("app.case_link.create.case_etag.check_args.url", "Url is required for each link")
	}
	caseTID, err := etag.EtagOrId(etag.EtagCase, req.CaseEtag)
	if err != nil {
		return nil, cerror.NewBadRequestError("app.case_link.create.case_etag.parse.error", err.Error())
	}

	createOpts, err := model.NewCreateOptions(ctx, req, CaseLinkMetadata)
	if err != nil {
		slog.Error(err.Error())
		return nil, AppInternalError
	}
	createOpts.ParentID = caseTID.GetOid()
	logAttributes := slog.Group(
		"context",
		slog.Int64("user_id", createOpts.GetAuthOpts().GetUserId()),
		slog.Int64("domain_id", createOpts.GetAuthOpts().GetDomainId()),
		slog.Int64("case_id", createOpts.ParentID),
	)

	if createOpts.GetAuthOpts().GetObjectScope(CaseLinkMetadata.GetMainScopeName()).IsRbacUsed() {
		access, err := c.app.Store.Case().CheckRbacAccess(createOpts, createOpts.GetAuthOpts(), authmodel.Edit, createOpts.ParentID)
		if err != nil {
			slog.Error(err.Error(), logAttributes)
			return nil, AppForbiddenError
		}
		if !access {
			slog.Error("user doesn't have required (EDIT) access to the case", logAttributes)
			return nil, AppForbiddenError
		}
	}
	res, dbErr := c.app.Store.CaseLink().Create(createOpts, req.Input)
	if dbErr != nil {
		slog.Error(dbErr.Error(), logAttributes)
		return nil, AppDatabaseError
	}

	err = NormalizeResponseLink(res, req)
	if err != nil {
		slog.Error(err.Error(), logAttributes)
		return nil, AppResponseNormalizingError
	}
	return res, nil
}

func (c *CaseLinkService) UpdateLink(ctx context.Context, req *cases.UpdateLinkRequest) (*cases.CaseLink, error) {
	if req.Input == nil {
		return nil, cerror.NewBadRequestError("app.case_link.update.check_args.input", "input required")
	}
	if req.Input.GetEtag() == "" {
		return nil, cerror.NewBadRequestError("app.case_link.update.check_args.id", "case ID required")
	}
	linkTID, err := etag.EtagOrId(etag.EtagCaseLink, req.GetInput().GetEtag())
	if err != nil {
		return nil, cerror.NewBadRequestError("app.case_link.create.link_etag.parse.error", err.Error())
	}
	caseTID, err := etag.EtagOrId(etag.EtagCase, req.GetCaseEtag())
	if err != nil {
		return nil, cerror.NewBadRequestError("app.case_link.create.case_etag.parse.error", err.Error())
	}
	updateOpts, err := model.NewUpdateOptions(ctx, req, CaseLinkMetadata)
	if err != nil {
		slog.Error(err.Error())
		return nil, AppInternalError
	}
	updateOpts.Etags = []*etag.Tid{&linkTID}
	updateOpts.ParentID = caseTID.GetOid()
	logAttributes := slog.Group(
		"context",
		slog.Int64("user_id", updateOpts.GetAuthOpts().GetUserId()),
		slog.Int64("domain_id", updateOpts.GetAuthOpts().GetDomainId()),
		slog.Int64("id", linkTID.GetOid()),
		slog.Int64("case_id", updateOpts.ParentID),
	)

	if updateOpts.GetAuthOpts().GetObjectScope(CaseLinkMetadata.GetMainScopeName()).IsRbacUsed() {
		access, err := c.app.Store.Case().CheckRbacAccess(updateOpts, updateOpts.GetAuthOpts(), authmodel.Edit, updateOpts.ParentID)
		if err != nil {
			slog.Error(err.Error(), logAttributes)
			return nil, AppForbiddenError
		}
		if !access {
			slog.Error("user doesn't have required (EDIT) access to the case", logAttributes)
			return nil, AppForbiddenError
		}
	}
	updated, err := c.app.Store.CaseLink().Update(updateOpts, req.Input)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	err = NormalizeResponseLink(updated, req)
	if err != nil {
		slog.Error(err.Error(), logAttributes)
		return nil, AppResponseNormalizingError
	}

	return updated, nil
}

func (c *CaseLinkService) DeleteLink(ctx context.Context, req *cases.DeleteLinkRequest) (*cases.CaseLink, error) {
	if req.GetEtag() == "" {
		return nil, cerror.NewBadRequestError("app.case_link.update.check_args.etag", "case etag required")
	}
	linkTID, err := etag.EtagOrId(etag.EtagCaseLink, req.GetEtag())
	if err != nil {
		return nil, cerror.NewBadRequestError("app.case_link.create.link_etag.parse.error", err.Error())
	}
	caseTID, err := etag.EtagOrId(etag.EtagCase, req.GetCaseEtag())
	if err != nil {
		return nil, cerror.NewBadRequestError("app.case_link.create.case_etag.parse.error", err.Error())
	}
	deleteOpts, err := model.NewDeleteOptions(ctx, CaseLinkMetadata)
	if err != nil {
		slog.Error(err.Error())
		return nil, AppInternalError
	}
	deleteOpts.ID = linkTID.GetOid()
	deleteOpts.ParentID = caseTID.GetOid()
	logAttributes := slog.Group(
		"context",
		slog.Int64("user_id", deleteOpts.GetAuthOpts().GetUserId()),
		slog.Int64("domain_id", deleteOpts.GetAuthOpts().GetDomainId()),
		slog.Int64("id", deleteOpts.ID),
		slog.Int64("case_id", deleteOpts.ParentID),
	)
	if deleteOpts.GetAuthOpts().GetObjectScope(CaseLinkMetadata.GetMainScopeName()).IsRbacUsed() {
		access, err := c.app.Store.Case().CheckRbacAccess(deleteOpts, deleteOpts.GetAuthOpts(), authmodel.Edit, deleteOpts.ParentID)
		if err != nil {
			slog.Error(err.Error(), logAttributes)
			return nil, AppForbiddenError
		}
		if !access {
			slog.Error("user doesn't have required (EDIT) access to the case", logAttributes)
			return nil, AppForbiddenError
		}
	}
	err = c.app.Store.CaseLink().Delete(deleteOpts)
	if err != nil {
		slog.Error(err.Error(), logAttributes)
		return nil, AppDatabaseError
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

	searchOpts, err := model.NewSearchOptions(ctx, req, CaseLinkMetadata)
	if err != nil {
		slog.Error(err.Error())
		return nil, AppInternalError
	}
	searchOpts.ParentId = etg.GetOid()
	//
	ids, err := util.ParseIds(req.GetIds(), etag.EtagCaseLink)
	if err != nil {
		return nil, cerror.NewBadRequestError("app.case_link.locate.parse_qin.invalid", err.Error())
	}
	searchOpts.IDs = ids
	logAttributes := slog.Group(
		"context",
		slog.Int64("user_id", searchOpts.GetAuthOpts().GetUserId()),
		slog.Int64("domain_id", searchOpts.GetAuthOpts().GetDomainId()),
		slog.Int64("id", searchOpts.ID),
		slog.Int64("case_id", searchOpts.ParentId),
	)
	if searchOpts.GetAuthOpts().GetObjectScope(CaseLinkMetadata.GetMainScopeName()).IsRbacUsed() {
		access, err := c.app.Store.Case().CheckRbacAccess(searchOpts, searchOpts.GetAuthOpts(), authmodel.Read, searchOpts.ParentId)
		if err != nil {
			slog.Error(err.Error(), logAttributes)
			return nil, AppForbiddenError
		}
		if !access {
			slog.Error("user doesn't have required (READ) access to the case", logAttributes)
			return nil, AppForbiddenError
		}
	}

	links, err := c.app.Store.CaseLink().List(searchOpts)
	if err != nil {
		slog.Error(err.Error(), logAttributes)
		return nil, AppDatabaseError
	}

	err = NormalizeResponseLinks(links, req.GetFields())
	if err != nil {
		slog.Error(err.Error(), logAttributes)
		return nil, AppResponseNormalizingError
	}
	//Return the located comment
	return links, nil
}

func NewCaseLinkService(app *App) (*CaseLinkService, cerror.AppError) {
	if app == nil {
		return nil, cerror.NewBadRequestError("app.case.new_case_comment_service.check_args.app", "unable to init service, app is nil")
	}
	return &CaseLinkService{app: app}, nil
}

func NormalizeResponseLink(res *cases.CaseLink, opts model.Fielder) error {
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
