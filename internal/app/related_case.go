package app

import (
	"context"
	"fmt"
	"github.com/webitel/cases/api/cases"
	"github.com/webitel/cases/auth"
	cerror "github.com/webitel/cases/internal/errors"
	deferr "github.com/webitel/cases/internal/errors/defaults"
	"github.com/webitel/cases/model"
	grpcopts "github.com/webitel/cases/model/options/grpc"
	"github.com/webitel/cases/util"
	"github.com/webitel/webitel-go-kit/etag"
	wlogger "github.com/webitel/webitel-go-kit/infra/logger_client"
	watcherkit "github.com/webitel/webitel-go-kit/pkg/watcher"
	"log/slog"
	"strconv"
)

type RelatedCaseService struct {
	app    *App
	logger *wlogger.ObjectedLogger
	cases.UnimplementedRelatedCasesServer
}

var RelatedCaseMetadata = model.NewObjectMetadata("", caseObjScope, []*model.Field{
	{Name: "id", Default: true},
	{Name: "ver", Default: true},
	{Name: "created_at", Default: true},
	{Name: "created_by", Default: true},
	{Name: "updated_at", Default: false},
	{Name: "updated_by", Default: false},
	{Name: "related_case", Default: true},
	{Name: "primary_case", Default: true},
	{Name: "relation", Default: true},
})

func (r *RelatedCaseService) LocateRelatedCase(ctx context.Context, req *cases.LocateRelatedCaseRequest) (*cases.RelatedCase, error) {
	if req.GetEtag() == "" {
		return nil, cerror.NewBadRequestError("app.related_case.locate_related_case.id_required", "ID is required")
	}
	caseTid, err := etag.EtagOrId(etag.EtagCase, req.GetPrimaryCaseEtag())
	if err != nil {
		return nil, cerror.NewBadRequestError("app.related_case.locate_related_case.invalid_primary_id", "Invalid ID")
	}
	searchOpts, err := grpcopts.NewLocateOptions(
		ctx,
		grpcopts.WithFields(req, CaseCommentMetadata,
			util.DeduplicateFields,
			func(in []string) []string {
				return util.EnsureFields(in, "created_at", "id")
			},
		),
		grpcopts.WithIDsAsEtags(etag.EtagRelatedCase, req.GetEtag()),
	)
	if err != nil {
		return nil, NewBadRequestError(err)
	}
	searchOpts.AddFilter("case_id", caseTid.GetOid())
	logAttributes := slog.Group(
		"context",
		slog.Int64("user_id", searchOpts.GetAuthOpts().GetUserId()),
		slog.Int64("domain_id", searchOpts.GetAuthOpts().GetDomainId()),
		slog.Int64("case_id", caseTid.GetOid()),
	)
	accessMode := auth.Read
	if searchOpts.GetAuthOpts().IsRbacCheckRequired(RelatedCaseMetadata.GetParentScopeName(), accessMode) {
		access, err := r.app.Store.Case().CheckRbacAccess(searchOpts, searchOpts.GetAuthOpts(), accessMode, caseTid.GetOid())
		if err != nil {
			slog.ErrorContext(ctx, err.Error(), logAttributes)
			return nil, deferr.ForbiddenError
		}
		if !access {
			slog.ErrorContext(ctx, "user doesn't have required (READ) access to the case", logAttributes)
			return nil, deferr.ForbiddenError
		}
	}

	output, err := r.app.Store.RelatedCase().List(searchOpts)
	if err != nil {
		return nil, err
	}
	if len(output.Data) == 0 {
		return nil, cerror.NewNotFoundError("app.related_case.locate_related_case.not_found", "Related case not found")
	} else if len(output.Data) > 1 {
		return nil, cerror.NewInternalError("app.related_case.locate_related_cases.multiple_found", "Multiple related cases found")
	}

	// Normalize IDs and handle errors
	if err := normalizeIDs(output.Data); err != nil {
		return nil, cerror.NewInternalError("app.related_case.locate_related_cases.normalize_ids_error", err.Error())
	}

	return output.Data[0], nil
}

func (r *RelatedCaseService) CreateRelatedCase(ctx context.Context, req *cases.CreateRelatedCaseRequest) (*cases.RelatedCase, error) {
	if req.GetPrimaryCaseEtag() == "" {
		return nil, cerror.NewBadRequestError("app.related_case.create_related_case.primary_case_id_required", "Primary case id required")
	}

	primaryCaseTag, err := etag.EtagOrId(etag.EtagCase, req.GetPrimaryCaseEtag())
	if err != nil {
		return nil, cerror.NewBadRequestError(
			"app.related_case.created_related_case.invalid_etag",
			"Invalid primary case etag",
		)
	}

	relatedCaseTag, err := etag.EtagOrId(etag.EtagCase, strconv.Itoa(int(req.GetInput().GetRelatedCase().GetId())))
	if err != nil {
		return nil, cerror.NewBadRequestError(
			"app.related_case.created_related_case.invalid_etag",
			"Invalid output etag",
		)
	}

	createOpts, err := grpcopts.NewCreateOptions(
		ctx,
		grpcopts.WithCreateFields(req, RelatedCaseMetadata,
			util.DeduplicateFields,
			util.ParseFieldsForEtag,
			util.EnsureIdField),
		grpcopts.WithCreateParentID(primaryCaseTag.GetOid()),
		grpcopts.WithCreateChildID(relatedCaseTag.GetOid()),
	)
	if err != nil {
		return nil, NewBadRequestError(err)
	}

	logAttributes := slog.Group(
		"context",
		slog.Int64("user_id", createOpts.GetAuthOpts().GetUserId()),
		slog.Int64("domain_id", createOpts.GetAuthOpts().GetDomainId()),
		slog.Int64("parent_id", createOpts.ParentID),
		slog.Int64("child_id", createOpts.ChildID),
	)
	primaryAccessMode := auth.Edit
	if createOpts.GetAuthOpts().IsRbacCheckRequired(RelatedCaseMetadata.GetParentScopeName(), primaryAccessMode) {
		primaryAccess, err := r.app.Store.Case().CheckRbacAccess(createOpts, createOpts.GetAuthOpts(), primaryAccessMode, createOpts.ParentID)
		if err != nil {
			slog.ErrorContext(ctx, err.Error(), logAttributes)
			return nil, deferr.ForbiddenError
		}
		if !primaryAccess {
			slog.ErrorContext(ctx, "user doesn't have required access to the primary case", logAttributes)
			return nil, deferr.ForbiddenError
		}
	}
	secondaryAccessMode := auth.Read
	if createOpts.GetAuthOpts().IsRbacCheckRequired(RelatedCaseMetadata.GetParentScopeName(), secondaryAccessMode) {
		secondaryAccess, err := r.app.Store.Case().CheckRbacAccess(createOpts, createOpts.GetAuthOpts(), secondaryAccessMode, createOpts.ChildID)
		if err != nil {
			slog.ErrorContext(ctx, err.Error(), logAttributes)
			return nil, deferr.ForbiddenError
		}
		if !secondaryAccess {
			slog.ErrorContext(ctx, "user doesn't have required access to the secondary case", logAttributes)
			return nil, deferr.ForbiddenError
		}
	}

	output, err := r.app.Store.RelatedCase().Create(
		createOpts,
		&req.GetInput().RelationType,
		// Used if explicitly set the case creator / updater instead of deriving it from the auth token.
		req.Input.GetUserID().GetId(),
	)
	if err != nil {
		return nil, err
	}

	userIP := createOpts.GetAuthOpts().GetUserIp()
	if userIP == "" {
		userIP = "unknown"
	}

	message, _ := wlogger.NewMessage(
		createOpts.GetAuthOpts().GetUserId(),
		userIP,
		wlogger.UpdateAction,
		strconv.Itoa(int(primaryCaseTag.GetOid())),
		req,
	)

	_, err = r.logger.SendContext(ctx, createOpts.GetAuthOpts().GetDomainId(), message)
	if err != nil {
		return nil, err
	}

	if notifyErr := r.app.watcherManager.Notify(
		model.BrokerScopeRelatedCases,
		watcherkit.EventTypeCreate,
		NewRelatedCaseWatcherData(
			createOpts.GetAuthOpts(),
			output,
			output.GetId(),
			createOpts.GetAuthOpts().GetDomainId(),
		)); notifyErr != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("could not notify related case create: %s, ", notifyErr.Error()), logAttributes)
	}

	output.Etag, err = etag.EncodeEtag(etag.EtagRelatedCase, output.GetId(), output.Ver)
	if err != nil {
		slog.ErrorContext(ctx, err.Error(), logAttributes)
		return nil, deferr.ResponseNormalizingError
	}
	output.RelatedCase.Etag, err = etag.EncodeEtag(etag.EtagCase, output.RelatedCase.GetId(), output.Ver)
	if err != nil {
		slog.ErrorContext(ctx, err.Error(), logAttributes)
		return nil, deferr.ResponseNormalizingError
	}
	output.PrimaryCase.Etag, err = etag.EncodeEtag(etag.EtagCase, output.PrimaryCase.GetId(), output.Ver)
	if err != nil {
		slog.ErrorContext(ctx, err.Error(), logAttributes)
		return nil, deferr.ResponseNormalizingError
	}

	return output, nil
}

func (r *RelatedCaseService) UpdateRelatedCase(ctx context.Context, req *cases.UpdateRelatedCaseRequest) (*cases.RelatedCase, error) {
	if req.GetEtag() == "" {
		return nil, cerror.NewBadRequestError("app.related_case.update_related_case.id_required", "ID required")
	}

	tag, err := etag.EtagOrId(etag.EtagRelatedCase, req.GetEtag())
	if err != nil {
		return nil, cerror.NewBadRequestError(
			"app.related_case.created_related_case.invalid_etag",
			"Invalid ID",
		)
	}
	updateOpts, err := grpcopts.NewUpdateOptions(
		ctx,
		grpcopts.WithUpdateFields(req, RelatedCaseMetadata),
		grpcopts.WithUpdateEtag(&tag),
		grpcopts.WithUpdateMasker(req),
	)
	if err != nil {
		return nil, NewBadRequestError(err)
	}

	primaryCaseTag, err := etag.EtagOrId(etag.EtagCase, strconv.Itoa(int(req.GetInput().GetPrimaryCase().GetId())))
	if err != nil {
		return nil, cerror.NewBadRequestError(
			"app.related_case.created_related_case.invalid_etag",
			"Invalid primary case etag",
		)
	}

	relatedCaseTag, err := etag.EtagOrId(etag.EtagCase, strconv.Itoa(int(req.GetInput().GetRelatedCase().GetId())))
	if err != nil {
		return nil, cerror.NewBadRequestError(
			"app.related_case.created_related_case.invalid_etag",
			"Invalid relatedCase etag",
		)
	}

	if primaryCaseTag.GetOid() == relatedCaseTag.GetOid() {
		return nil, cerror.NewBadRequestError(
			"app.related_case.update_related_case.invalid_ids",
			"A case cannot be related to itself",
		)
	}

	input := &cases.InputRelatedCase{
		PrimaryCase:  req.Input.GetPrimaryCase(),
		RelatedCase:  req.Input.GetRelatedCase(),
		RelationType: req.Input.RelationType,
	}

	primaryId := primaryCaseTag.GetOid()
	relatedId := relatedCaseTag.GetOid()
	logAttributes := slog.Group(
		"context",
		slog.Int64("user_id", updateOpts.GetAuthOpts().GetUserId()),
		slog.Int64("domain_id", updateOpts.GetAuthOpts().GetDomainId()),
		slog.Int64("parent_id", updateOpts.ParentID),
	)
	primaryAccessMode := auth.Edit
	if updateOpts.GetAuthOpts().IsRbacCheckRequired(RelatedCaseMetadata.GetParentScopeName(), primaryAccessMode) {
		primaryAccess, err := r.app.Store.Case().CheckRbacAccess(updateOpts, updateOpts.GetAuthOpts(), primaryAccessMode, primaryId)
		if err != nil {
			slog.ErrorContext(ctx, err.Error(), logAttributes)
			return nil, deferr.ForbiddenError
		}
		if !primaryAccess {
			slog.ErrorContext(ctx, "user doesn't have required access to the primary case", logAttributes)
			return nil, deferr.ForbiddenError
		}
	}
	secondaryAccessMode := auth.Read
	if updateOpts.GetAuthOpts().IsRbacCheckRequired(RelatedCaseMetadata.GetParentScopeName(), secondaryAccessMode) {
		secondaryAccess, err := r.app.Store.Case().CheckRbacAccess(updateOpts, updateOpts.GetAuthOpts(), secondaryAccessMode, relatedId)
		if err != nil {
			slog.ErrorContext(ctx, err.Error(), logAttributes)
			return nil, deferr.ForbiddenError
		}
		if !secondaryAccess {
			slog.ErrorContext(ctx, "user doesn't have required access to the secondary case", logAttributes)
			return nil, deferr.ForbiddenError
		}
	}

	output, err := r.app.Store.RelatedCase().Update(
		updateOpts,
		input,
		// Used if explicitly set the case creator / updater instead of deriving it from the auth token.
		req.Input.GetUserID().GetId(),
	)
	if err != nil {
		return nil, err
	}

	userIP := updateOpts.GetAuthOpts().GetUserIp()
	if userIP == "" {
		userIP = "unknown"
	}

	message, _ := wlogger.NewMessage(
		updateOpts.GetAuthOpts().GetUserId(),
		userIP,
		wlogger.UpdateAction,
		strconv.Itoa(int(primaryCaseTag.GetOid())),
		req,
	)

	_, err = r.logger.SendContext(ctx, updateOpts.GetAuthOpts().GetDomainId(), message)
	if err != nil {
		return nil, err
	}

	if notifyErr := r.app.watcherManager.Notify(
		model.BrokerScopeRelatedCases,
		watcherkit.EventTypeUpdate,
		NewRelatedCaseWatcherData(
			updateOpts.GetAuthOpts(),
			output,
			output.GetId(),
			updateOpts.GetAuthOpts().GetDomainId(),
		)); notifyErr != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("could not notify related case create: %s, ", notifyErr.Error()), logAttributes)
	}

	output.Etag, err = etag.EncodeEtag(etag.EtagRelatedCase, output.GetId(), output.GetVer())
	if err != nil {
		slog.ErrorContext(ctx, err.Error(), logAttributes)
		return nil, deferr.ResponseNormalizingError
	}
	return output, nil
}

func (r *RelatedCaseService) DeleteRelatedCase(ctx context.Context, req *cases.DeleteRelatedCaseRequest) (*cases.RelatedCase, error) {
	if req.GetEtag() == "" {
		return nil, cerror.NewBadRequestError("app.related_case.delete_related_case.id_required", "ID required")
	}
	if req.GetPrimaryCaseEtag() == "" {
		return nil, cerror.NewBadRequestError("app.related_case.delete_related_case.primary_id_required", "Primary case ID required")
	}
	deleteOpts, err := grpcopts.NewDeleteOptions(
		ctx,
		grpcopts.WithDeleteIDsAsEtags(
			etag.EtagRelatedCase,
			req.GetEtag()),
		grpcopts.WithDeleteParentIDAsEtag(
			etag.EtagCase,
			req.GetPrimaryCaseEtag(),
		),
	)

	if err != nil {
		return nil, NewBadRequestError(err)
	}

	primaryCaseTag, err := etag.EtagOrId(etag.EtagCase, req.GetPrimaryCaseEtag())
	if err != nil {
		return nil, cerror.NewBadRequestError(
			"app.related_case.deleted_related_case.invalid_etag",
			"Invalid primary case etag",
		)
	}

	objTag, err := etag.EtagOrId(etag.EtagRelatedCase, req.GetEtag())
	if err != nil {
		return nil, cerror.NewBadRequestError(
			"app.related_case.deleted_related_case.invalid_etag",
			"Invalid relation etag",
		)
	}

	logAttributes := slog.Group(
		"context",
		slog.Int64("user_id", deleteOpts.GetAuthOpts().GetUserId()),
		slog.Int64("domain_id", deleteOpts.GetAuthOpts().GetDomainId()),
		slog.Int64("parent_id", deleteOpts.ParentID),
	)

	accessMode := auth.Edit
	if deleteOpts.GetAuthOpts().IsRbacCheckRequired(RelatedCaseMetadata.GetParentScopeName(), accessMode) {
		access, err := r.app.Store.Case().CheckRbacAccess(deleteOpts, deleteOpts.GetAuthOpts(), accessMode, deleteOpts.GetParentID())
		if err != nil {
			slog.ErrorContext(ctx, err.Error(), logAttributes)
			return nil, deferr.ForbiddenError
		}
		if !access {
			slog.ErrorContext(ctx, "user doesn't have required (EDIT) access to the case", logAttributes)
			return nil, deferr.ForbiddenError
		}

	}

	err = r.app.Store.RelatedCase().Delete(deleteOpts)
	if err != nil {
		slog.ErrorContext(ctx, err.Error(), logAttributes)
		return nil, deferr.DatabaseError
	}

	userIP := deleteOpts.GetAuthOpts().GetUserIp()
	if userIP == "" {
		userIP = "unknown"
	}

	message, _ := wlogger.NewMessage(
		deleteOpts.GetAuthOpts().GetUserId(),
		userIP,
		wlogger.UpdateAction,
		strconv.Itoa(int(primaryCaseTag.GetOid())),
		req,
	)

	_, err = r.logger.SendContext(ctx, deleteOpts.GetAuthOpts().GetDomainId(), message)
	if err != nil {
		return nil, err
	}

	if notifyErr := r.app.watcherManager.Notify(
		model.BrokerScopeRelatedCases,
		watcherkit.EventTypeDelete,
		NewRelatedCaseWatcherData(
			deleteOpts.GetAuthOpts(),
			&cases.RelatedCase{},
			objTag.GetOid(),
			deleteOpts.GetAuthOpts().GetDomainId(),
		)); notifyErr != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("could not notify related case create: %s, ", notifyErr.Error()), logAttributes)
	}

	return nil, nil
}

func (r *RelatedCaseService) ListRelatedCases(ctx context.Context, req *cases.ListRelatedCasesRequest) (*cases.RelatedCaseList, error) {
	if req.GetPrimaryCaseEtag() == "" {
		return nil, cerror.NewBadRequestError("app.related_case.list_related_case.id_required", "ID required")
	}
	searchOpts, err := grpcopts.NewSearchOptions(
		ctx,
		grpcopts.WithSearch(req),
		grpcopts.WithPagination(req),
		grpcopts.WithFields(req, RelatedCaseMetadata,
			util.DeduplicateFields,
			func(in []string) []string {
				return util.EnsureFields(in, "created_at", "id")
			},
		),
		grpcopts.WithSort(req),
		grpcopts.WithIDsAsEtags(etag.EtagRelatedCase, req.GetIds()...),
	)
	if err != nil {
		return nil, NewBadRequestError(err)
	}
	tag, err := etag.EtagOrId(etag.EtagCase, req.PrimaryCaseEtag)
	if err != nil {
		return nil, cerror.NewBadRequestError("app.related_case.list_related_cases.invalid_etag", "Invalid etag")
	}
	searchOpts.AddFilter("case_id", tag.GetOid())

	output, err := r.app.Store.RelatedCase().List(searchOpts)
	if err != nil {
		return nil, err
	}

	// Normalize IDs and handle errors
	if err := normalizeIDs(output.Data); err != nil {
		return nil, cerror.NewInternalError("app.related_case.list_related_cases.normalize_ids_error", err.Error())
	}
	return output, nil
}

func normalizeIDs(relatedCases []*cases.RelatedCase) error {
	for _, relatedCase := range relatedCases {
		if relatedCase == nil {
			continue
		}
		var err error
		// Normalize related case ID
		relatedCase.Etag, err = etag.EncodeEtag(etag.EtagRelatedCase, relatedCase.GetId(), relatedCase.Ver)
		if err != nil {
			return err
		}

		// Normalize primary case ID
		if relatedCase.PrimaryCase != nil {

			relatedCase.PrimaryCase.Etag, err = etag.EncodeEtag(etag.EtagCase, relatedCase.PrimaryCase.GetId(), relatedCase.PrimaryCase.GetVer())
			if err != nil {
				return err
			}
			// Set PrimaryCase Ver to nil
			relatedCase.PrimaryCase.Ver = 0
		}

		// Normalize related case ID inside related case
		if relatedCase.RelatedCase != nil {
			relatedCase.RelatedCase.Etag, err = etag.EncodeEtag(etag.EtagCase, relatedCase.RelatedCase.Id, relatedCase.RelatedCase.GetVer())
			if err != nil {
				return err
			}
			// Set RelatedCase Ver to nil
			relatedCase.RelatedCase.Ver = 0
		}
	}

	return nil
}

func NewRelatedCaseService(app *App) (*RelatedCaseService, error) {
	if app == nil {
		return nil, cerror.NewBadRequestError(
			"app.case.new_related_case_service.check_args.app",
			"unable to init service, app is nil",
		)
	}
	logger, err := app.wtelLogger.GetObjectedLogger("cases")
	if err != nil {
		return nil, err
	}

	service := &RelatedCaseService{
		app:    app,
		logger: logger,
	}

	watcher := watcherkit.NewDefaultWatcher()

	if app.config.TriggerWatcher.Enabled {
		mq, err := NewTriggerObserver(app.rabbit, app.config.TriggerWatcher, formRelatedCaseTriggerModel, slog.With(
			slog.Group("context",
				slog.String("scope", "watcher")),
		))

		if err != nil {
			return nil, cerror.NewInternalError("app.case.new_related_case_service.create_mq_observer.app", err.Error())
		}
		watcher.Attach(watcherkit.EventTypeCreate, mq)
		watcher.Attach(watcherkit.EventTypeUpdate, mq)
		watcher.Attach(watcherkit.EventTypeDelete, mq)
		watcher.Attach(watcherkit.EventTypeResolutionTime, mq)

		app.caseResolutionTimer.Start()
	}

	app.watcherManager.AddWatcher(model.BrokerScopeRelatedCases, watcher)

	return service, nil
}

func formRelatedCaseTriggerModel(item *cases.RelatedCase) (*model.RelatedCaseAMQPMessage, error) {
	m := &model.RelatedCaseAMQPMessage{
		RelatedCase: item,
	}
	return m, nil
}

type RelatedCaseWatcherData struct {
	relCase *cases.RelatedCase
	Args    map[string]any
}

func (wd *RelatedCaseWatcherData) GetArgs() map[string]any {
	return wd.Args
}

func NewRelatedCaseWatcherData(session auth.Auther, relCase *cases.RelatedCase, relCaseID int64, dc int64) *RelatedCaseWatcherData {
	return &RelatedCaseWatcherData{
		relCase: relCase,
		Args: map[string]any{
			"session":   session,
			"obj":       relCase,
			"id":        relCaseID,
			"domain_id": dc,
		},
	}
}
