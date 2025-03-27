package app

import (
	"context"
	"github.com/webitel/cases/api/cases"
	"github.com/webitel/cases/auth"
	cerror "github.com/webitel/cases/internal/errors"
	deferr "github.com/webitel/cases/internal/errors/defaults"
	"github.com/webitel/cases/model"
	grpcopts "github.com/webitel/cases/model/options/grpc"
	"github.com/webitel/cases/model/options/grpc/shared"
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

	caseEtg, err := etag.EtagOrId(etag.EtagCase, req.GetCaseEtag())
	if err != nil {
		return nil, cerror.NewBadRequestError("app.case_link.locate.parse_case_etag.error", err.Error())
	}

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
		return nil, NewBadRequestError(err)
	}
	searchOpts.AddFilter("case_id", caseEtg.GetOid())
	logAttributes := slog.Group(
		"context",
		slog.Int64("user_id", searchOpts.GetAuthOpts().GetUserId()),
		slog.Int64("domain_id", searchOpts.GetAuthOpts().GetDomainId()),
		slog.Int64("id", searchOpts.IDs[0]),
		slog.Int64("case_id", caseEtg.GetOid()),
	)
	accessMode := auth.Read
	if searchOpts.GetAuthOpts().IsRbacCheckRequired(CaseLinkMetadata.GetParentScopeName(), accessMode) {
		access, err := c.app.Store.Case().CheckRbacAccess(searchOpts, searchOpts.GetAuthOpts(), accessMode, caseEtg.GetOid())
		if err != nil {
			slog.ErrorContext(ctx, err.Error(), logAttributes)
			return nil, deferr.ForbiddenError
		}
		if !access {
			slog.ErrorContext(ctx, "user doesn't have required (READ) access to the case", logAttributes)
			return nil, deferr.ForbiddenError
		}
	}
	links, err := c.app.Store.CaseLink().List(searchOpts)
	if err != nil {
		slog.ErrorContext(ctx, err.Error(), logAttributes)
		return nil, deferr.DatabaseError
	}
	if len(links.Items) == 0 {
		return nil, cerror.NewNotFoundError("app.case_link.locate.check_items.error", "not found")
	}
	res := links.Items[0]
	// hide etag if needed
	err = NormalizeResponseLink(res, req)
	if err != nil {
		slog.ErrorContext(ctx, err.Error(), logAttributes)
		return nil, deferr.ResponseNormalizingError
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
	caseTid, err := etag.EtagOrId(etag.EtagCase, req.CaseEtag)
	if err != nil {
		return nil, cerror.NewBadRequestError("app.case_link.create.case_etag.parse.error", err.Error())
	}

	createOpts, err := grpcopts.NewCreateOptions(
		ctx,
		grpcopts.WithCreateFields(req, CaseLinkMetadata,
			util.DeduplicateFields,
			util.ParseFieldsForEtag,
			util.EnsureIdField),
		grpcopts.WithCreateParentID(caseTid.GetOid()),
	)
	if err != nil {
		return nil, NewBadRequestError(err)
	}

	logAttributes := slog.Group(
		"context",
		slog.Int64("user_id", createOpts.GetAuthOpts().GetUserId()),
		slog.Int64("domain_id", createOpts.GetAuthOpts().GetDomainId()),
		slog.Int64("case_id", createOpts.ParentID),
	)
	accessMode := auth.Edit
	if createOpts.GetAuthOpts().IsRbacCheckRequired(CaseLinkMetadata.GetParentScopeName(), accessMode) {
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
	res, dbErr := c.app.Store.CaseLink().Create(createOpts, req.Input)
	if dbErr != nil {
		slog.ErrorContext(ctx, dbErr.Error(), logAttributes)
		return nil, deferr.DatabaseError
	}

	err = NormalizeResponseLink(res, req)
	if err != nil {
		slog.ErrorContext(ctx, err.Error(), logAttributes)
		return nil, deferr.ResponseNormalizingError
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
	linkTid, err := etag.EtagOrId(etag.EtagCaseLink, req.GetInput().GetEtag())
	if err != nil {
		return nil, cerror.NewBadRequestError("app.case_link.create.link_etag.parse.error", err.Error())
	}
	caseTid, err := etag.EtagOrId(etag.EtagCase, req.GetCaseEtag())
	if err != nil {
		return nil, cerror.NewBadRequestError("app.case_link.create.case_etag.parse.error", err.Error())
	}
	updateOpts, err := grpcopts.NewUpdateOptions(
		ctx,
		grpcopts.WithUpdateFields(req, CaseLinkMetadata),
		grpcopts.WithUpdateParentID(caseTid.GetOid()),
		grpcopts.WithUpdateEtag(&linkTid),
	)
	if err != nil {
		return nil, NewBadRequestError(err)
	}

	logAttributes := slog.Group(
		"context",
		slog.Int64("user_id", updateOpts.GetAuthOpts().GetUserId()),
		slog.Int64("domain_id", updateOpts.GetAuthOpts().GetDomainId()),
		slog.Int64("id", linkTid.GetOid()),
		slog.Int64("case_id", updateOpts.ParentID),
	)
	accessMode := auth.Edit
	if updateOpts.GetAuthOpts().IsRbacCheckRequired(CaseLinkMetadata.GetParentScopeName(), accessMode) {
		access, err := c.app.Store.Case().CheckRbacAccess(updateOpts, updateOpts.GetAuthOpts(), auth.Edit, updateOpts.ParentID)
		if err != nil {
			slog.ErrorContext(ctx, err.Error(), logAttributes)
			return nil, deferr.ForbiddenError
		}
		if !access {
			slog.ErrorContext(ctx, "user doesn't have required (EDIT) access to the case", logAttributes)
			return nil, deferr.ForbiddenError
		}
	}
	updated, err := c.app.Store.CaseLink().Update(updateOpts, req.Input)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return nil, err
	}
	err = NormalizeResponseLink(updated, req)
	if err != nil {
		slog.ErrorContext(ctx, err.Error(), logAttributes)
		return nil, deferr.ResponseNormalizingError
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
	deleteOpts, err := grpcopts.NewDeleteOptions(ctx, grpcopts.WithDeleteID(linkTID.GetOid()), grpcopts.WithDeleteParentIDAsEtag(etag.EtagCase, req.GetCaseEtag()))
	if err != nil {
		return nil, NewBadRequestError(err)
	}
	logAttributes := slog.Group(
		"context",
		slog.Int64("user_id", deleteOpts.GetAuthOpts().GetUserId()),
		slog.Int64("domain_id", deleteOpts.GetAuthOpts().GetDomainId()),
		slog.Int64("id", linkTID.GetOid()),
		slog.Int64("case_id", deleteOpts.ParentID),
	)
	accessMode := auth.Edit
	if deleteOpts.GetAuthOpts().IsRbacCheckRequired(CaseLinkMetadata.GetParentScopeName(), accessMode) {
		access, err := c.app.Store.Case().CheckRbacAccess(deleteOpts, deleteOpts.GetAuthOpts(), auth.Edit, deleteOpts.ParentID)
		if err != nil {
			slog.ErrorContext(ctx, err.Error(), logAttributes)
			return nil, deferr.ForbiddenError
		}
		if !access {
			slog.ErrorContext(ctx, "user doesn't have required (EDIT) access to the case", logAttributes)
			return nil, deferr.ForbiddenError
		}
	}
	err = c.app.Store.CaseLink().Delete(deleteOpts)
	if err != nil {
		slog.ErrorContext(ctx, err.Error(), logAttributes)
		return nil, deferr.DatabaseError
	}

	return nil, nil
}

func (c *CaseLinkService) ListLinks(ctx context.Context, req *cases.ListLinksRequest) (*cases.CaseLinkList, error) {
	// Validate required fields
	if req.GetCaseEtag() == "" {
		return nil, cerror.NewBadRequestError("app.case_link.list.case_etag.check_args.etag", "case etag is required")
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
		return nil, NewBadRequestError(err)
	}
	etg, err := etag.EtagOrId(etag.EtagCase, req.GetCaseEtag())
	if err != nil {
		return nil, cerror.NewBadRequestError("app.case_link.locate.parse_etag.error", err.Error())
	}
	searchOpts.AddFilter("case_id", etg.GetOid())
	logAttributes := slog.Group(
		"context",
		slog.Int64("user_id", searchOpts.GetAuthOpts().GetUserId()),
		slog.Int64("domain_id", searchOpts.GetAuthOpts().GetDomainId()),
		slog.Int64("case_id", etg.GetOid()),
	)
	accessMode := auth.Read
	if searchOpts.GetAuthOpts().IsRbacCheckRequired(CaseLinkMetadata.GetParentScopeName(), accessMode) {
		access, err := c.app.Store.Case().CheckRbacAccess(searchOpts, searchOpts.GetAuthOpts(), accessMode, etg.GetOid())
		if err != nil {
			slog.ErrorContext(ctx, err.Error(), logAttributes)
			return nil, deferr.ForbiddenError
		}
		if !access {
			slog.ErrorContext(ctx, "user doesn't have required (READ) access to the case", logAttributes)
			return nil, deferr.ForbiddenError
		}
	}

	links, err := c.app.Store.CaseLink().List(searchOpts)
	if err != nil {
		slog.ErrorContext(ctx, err.Error(), logAttributes)
		return nil, deferr.DatabaseError
	}

	err = NormalizeResponseLinks(links, req.GetFields())
	if err != nil {
		slog.ErrorContext(ctx, err.Error(), logAttributes)
		return nil, deferr.ResponseNormalizingError
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
