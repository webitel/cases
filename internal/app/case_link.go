package app

import (
	"context"
	"fmt"
	"github.com/webitel/cases/api/cases"
	"github.com/webitel/cases/auth"
	auth_util "github.com/webitel/cases/auth/util"
	cerror "github.com/webitel/cases/internal/errors"
	deferr "github.com/webitel/cases/internal/errors/defaults"
	"github.com/webitel/cases/model"
	grpcopts "github.com/webitel/cases/model/options/grpc"
	"github.com/webitel/cases/model/options/grpc/shared"
	"github.com/webitel/cases/util"
	"github.com/webitel/webitel-go-kit/etag"
	wlogger "github.com/webitel/webitel-go-kit/infra/logger_client"
	watcherkit "github.com/webitel/webitel-go-kit/pkg/watcher"
	"log/slog"
	"strconv"
)

// In search options extract from context user
// Remove from search options fields functions

type CaseLinkService struct {
	app    *App
	logger *wlogger.ObjectedLogger
	cases.UnimplementedCaseLinksServer
}

var CaseLinkMetadata = model.NewObjectMetadata("", caseObjScope, []*model.Field{
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
	searchOpts.AddFilter(fmt.Sprintf("case_id=%d", caseEtg.GetOid()))
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
	authOpts := createOpts.GetAuthOpts()
	if overrideID := req.Input.UserID.GetId(); overrideID != 0 {
		authOpts = auth_util.CloneWithUserID(authOpts, overrideID)
	}

	message, err := wlogger.NewMessage(
		createOpts.GetAuthOpts().GetUserId(),
		createOpts.GetAuthOpts().GetUserIp(),
		wlogger.UpdateAction,
		strconv.FormatInt(res.GetId(), 10),
		req,
	)
	_, err = c.logger.SendContext(ctx, createOpts.GetAuthOpts().GetDomainId(), message)
	if err != nil {
		return nil, err
	}

	if notifyErr := c.app.watcherManager.Notify(
		model.BrokerScopeCaseLinks,
		watcherkit.EventTypeCreate,
		NewLinkWatcherData(
			authOpts,
			res,
			res.GetId(),
			authOpts.GetDomainId(),
		)); notifyErr != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("could not notify link create: %s, ", notifyErr.Error()), logAttributes)
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
		grpcopts.WithUpdateMasker(req),
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

	authOpts := updateOpts.GetAuthOpts()
	if overrideID := req.Input.UserID.GetId(); overrideID != 0 {
		authOpts = auth_util.CloneWithUserID(authOpts, overrideID)
	}

	message, err := wlogger.NewMessage(
		updateOpts.GetAuthOpts().GetUserId(),
		updateOpts.GetAuthOpts().GetUserIp(),
		wlogger.UpdateAction,
		strconv.FormatInt(linkTid.GetOid(), 10),
		req,
	)
	_, err = c.logger.SendContext(ctx, updateOpts.GetAuthOpts().GetDomainId(), message)
	if err != nil {
		return nil, err
	}

	if notifyErr := c.app.watcherManager.Notify(
		model.BrokerScopeCaseLinks,
		watcherkit.EventTypeUpdate,
		NewLinkWatcherData(
			authOpts,
			updated,
			updated.GetId(),
			authOpts.GetDomainId(),
		)); notifyErr != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("could not notify link update: %s, ", notifyErr.Error()), logAttributes)
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

	message, err := wlogger.NewMessage(
		deleteOpts.GetAuthOpts().GetUserId(),
		deleteOpts.GetAuthOpts().GetUserIp(),
		wlogger.UpdateAction,
		strconv.FormatInt(linkTID.GetOid(), 10),
		req,
	)
	_, err = c.logger.SendContext(ctx, deleteOpts.GetAuthOpts().GetDomainId(), message)
	if err != nil {
		return nil, err
	}

	if notifyErr := c.app.watcherManager.Notify(
		model.BrokerScopeCaseLinks,
		watcherkit.EventTypeDelete,
		NewLinkWatcherData(
			deleteOpts.GetAuthOpts(),
			&cases.CaseLink{},
			linkTID.GetOid(),
			deleteOpts.GetAuthOpts().GetDomainId(),
		)); notifyErr != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("could not notify link delete: %s, ", notifyErr.Error()), logAttributes)
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
	searchOpts.AddFilter(fmt.Sprintf("case_id=%d", etg.GetOid()))
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

func NewCaseLinkService(app *App) (*CaseLinkService, error) {
	if app == nil {
		return nil, cerror.NewBadRequestError(
			"app.case.new_case_comment_service.check_args.app",
			"unable to init service, app is nil",
		)
	}
	logger, err := app.wtelLogger.GetObjectedLogger("cases")
	if err != nil {
		return nil, err
	}
	service := &CaseLinkService{
		app:    app,
		logger: logger,
	}
	watcher := watcherkit.NewDefaultWatcher()

	if app.config.TriggerWatcher.Enabled {
		mq, err := NewTriggerObserver(app.rabbit, app.config.TriggerWatcher, formCaseLinkTriggerModel, slog.With(
			slog.Group("context",
				slog.String("scope", "watcher")),
		))

		if err != nil {
			return nil, cerror.NewInternalError("app.case.new_case_link_service.create_mq_observer.app", err.Error())
		}
		watcher.Attach(watcherkit.EventTypeCreate, mq)
		watcher.Attach(watcherkit.EventTypeUpdate, mq)
		watcher.Attach(watcherkit.EventTypeDelete, mq)
		watcher.Attach(watcherkit.EventTypeResolutionTime, mq)

		app.caseResolutionTimer.Start()
	}

	app.watcherManager.AddWatcher(model.BrokerScopeCaseLinks, watcher)

	return service, nil
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

type CaseLinkWatcherData struct {
	link *cases.CaseLink
	Args map[string]any
}

func (wd *CaseLinkWatcherData) GetArgs() map[string]any {
	return wd.Args
}

func NewLinkWatcherData(session auth.Auther, link *cases.CaseLink, linkId int64, dc int64) *CaseLinkWatcherData {
	return &CaseLinkWatcherData{
		link: link,
		Args: map[string]any{
			"session":   session,
			"obj":       link,
			"id":        linkId,
			"domain_id": dc,
		},
	}
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
