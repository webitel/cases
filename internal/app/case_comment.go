package app

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/webitel/cases/api/cases"
	"github.com/webitel/cases/auth"
	"github.com/webitel/cases/internal/errors"
	"github.com/webitel/cases/internal/model"
	grpcopts "github.com/webitel/cases/internal/model/options/grpc"
	"github.com/webitel/cases/internal/model/options/grpc/shared"
	"github.com/webitel/cases/util"
	"github.com/webitel/webitel-go-kit/pkg/etag"
	watcherkit "github.com/webitel/webitel-go-kit/pkg/watcher"
	"google.golang.org/grpc/codes"
	"log/slog"
	"time"
)

const caseCommentsObjScope = model.ScopeCaseComments

var CaseCommentMetadata = model.NewObjectMetadata(caseCommentsObjScope, caseObjScope, []*model.Field{
	{Name: "id", Default: false},
	{Name: "etag", Default: true},
	{Name: "ver", Default: false},
	{Name: "created_at", Default: true},
	{Name: "created_by", Default: true},
	{Name: "updated_at", Default: true},
	{Name: "updated_by", Default: true},
	{Name: "text", Default: true},
	{Name: "edited", Default: true},
	{Name: "can_edit", Default: true},
	{Name: "author", Default: true},
	{Name: "role_ids", Default: false},
	{Name: "case_id", Default: false},
})

type CaseCommentService struct {
	app *App
	cases.UnimplementedCaseCommentsServer
}

func (c *CaseCommentService) LocateComment(
	ctx context.Context,
	req *cases.LocateCommentRequest,
) (*cases.CaseComment, error) {
	if req.Etag == "" {
		return nil, errors.InvalidArgument("Etag is required")
	}

	searchOpts, err := grpcopts.NewLocateOptions(
		ctx,
		grpcopts.WithFields(req, CaseCommentMetadata,
			util.DeduplicateFields,
			util.ParseFieldsForEtag,
			func(in []string) []string {
				if util.ContainsField(in, "edited") {
					return util.EnsureFields(in, "updated_at", "created_at")
				}
				return in
			},
		),
		grpcopts.WithIDsAsEtags(etag.EtagCaseComment, req.GetEtag()),
	)
	if err != nil {
		return nil, err
	}
	commentList, err := c.app.Store.CaseComment().List(searchOpts)
	if err != nil {
		return nil, err
	}

	if len(commentList.Items) == 0 {
		return nil, errors.NotFound("Comment not found")
	} else if len(commentList.Items) > 1 {
		return nil, errors.New("too many items found", errors.WithCode(codes.AlreadyExists))
	}

	err = NormalizeCommentsResponse(commentList.Items[0], req)
	if err != nil {
		return nil, err
	}

	return commentList.Items[0], nil
}

func (c *CaseCommentService) UpdateComment(
	ctx context.Context,
	req *cases.UpdateCommentRequest,
) (*cases.CaseComment, error) {
	if req.Input.Etag == "" {
		return nil, errors.InvalidArgument("Etag is required")
	}
	if req.Input.Text == "" {
		return nil, errors.InvalidArgument("Text is required")
	}

	tag, err := etag.EtagOrId(etag.EtagCaseComment, req.Input.Etag)
	if err != nil {
		return nil, errors.InvalidArgument("invalid Etag", errors.WithCause(err))
	}

	opts := []grpcopts.UpdateOption{
		grpcopts.WithUpdateFields(req, CaseCommentMetadata),
		grpcopts.WithUpdateEtag(&tag),
		grpcopts.WithUpdateMasker(req),
	}

	if ts := req.GetUpdatedAt(); ts != 0 {
		opts = append(opts, grpcopts.WithUpdateTime(time.UnixMilli(ts)))
	}

	updateOpts, err := grpcopts.NewUpdateOptions(ctx, opts...)
	if err != nil {
		return nil, err
	}

	input := &cases.CaseComment{
		// Used if explicitly set the case creator / updater instead of deriving it from the auth token.
		UpdatedBy: req.Input.GetUserID(),
		Id:        tag.GetOid(),
		Text:      req.Input.Text,
		Ver:       tag.GetVer(),
	}

	updatedComment, err := c.app.Store.CaseComment().Update(updateOpts, input)
	if err != nil {
		return nil, err
	}

	id := updatedComment.GetId()
	roleIds := updatedComment.GetRoleIds()
	parentId := input.GetCaseId()

	err = NormalizeCommentsResponse(updatedComment, req)
	if err != nil {
		return nil, err
	}

	if notifyErr := c.app.watcherManager.Notify(
		caseCommentsObjScope,
		watcherkit.EventTypeUpdate,
		NewCaseCommentWatcherData(updateOpts.GetAuthOpts(), updatedComment, id, parentId, roleIds),
	); notifyErr != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("could not notify input update: %s, ", notifyErr.Error()))
	}

	return updatedComment, nil
}

func (c *CaseCommentService) DeleteComment(
	ctx context.Context,
	req *cases.DeleteCommentRequest,
) (*cases.CaseComment, error) {
	if req.Etag == "" {
		return nil, errors.InvalidArgument("etag is required")
	}
	tag, err := etag.EtagOrId(etag.EtagCaseComment, req.GetEtag())
	if err != nil {
		return nil, errors.InvalidArgument("invalid Etag", errors.WithCause(err))
	}
	deleteOpts, err := grpcopts.NewDeleteOptions(ctx, grpcopts.WithDeleteID(tag.GetOid()))
	if err != nil {
		return nil, err
	}
	err = c.app.Store.CaseComment().Delete(deleteOpts)
	if err != nil {
		return nil, err
	}

	if notifyErr := c.app.watcherManager.Notify(
		caseCommentsObjScope,
		watcherkit.EventTypeDelete,
		NewCaseCommentWatcherData(deleteOpts.GetAuthOpts(), nil, tag.GetOid(), 0, nil),
	); notifyErr != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("could not notify comment delete: %s, ", notifyErr.Error()))
	}
	return nil, nil
}

func (c *CaseCommentService) ListComments(
	ctx context.Context,
	req *cases.ListCommentsRequest,
) (*cases.CaseCommentList, error) {
	if req.CaseEtag == "" {
		return nil, errors.InvalidArgument("case etag is required")
	}
	searchOpts, err := grpcopts.NewSearchOptions(
		ctx,
		grpcopts.WithSearch(req),
		grpcopts.WithPagination(req),
		grpcopts.WithFields(req, CaseCommentMetadata,
			util.DeduplicateFields,
			util.ParseFieldsForEtag,
			func(in []string) []string {
				if util.ContainsField(in, "edited") {
					return util.EnsureFields(in, "updated_at", "created_at")
				}
				return in
			},
		),
		grpcopts.WithIDsAsEtags(etag.EtagCaseComment, req.GetIds()...),
		grpcopts.WithSort(req),
	)
	if err != nil {
		return nil, err
	}
	tag, err := etag.EtagOrId(etag.EtagCase, req.CaseEtag)
	if err != nil {
		return nil, errors.InvalidArgument("Invalid etag", errors.WithCause(err))
	}
	if tag.GetOid() != 0 {
		searchOpts.AddFilter(fmt.Sprintf("case_id=%d", tag.GetOid()))
	}
	comments, err := c.app.Store.CaseComment().List(searchOpts)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return nil, err
	}

	err = NormalizeCommentsResponse(comments, req)
	if err != nil {
		return nil, err
	}

	return comments, nil
}

func (c *CaseCommentService) PublishComment(
	ctx context.Context,
	req *cases.PublishCommentRequest,
) (*cases.CaseComment, error) {
	if req.CaseEtag == "" {
		return nil, errors.InvalidArgument("case etag is required")
	} else if req.Input.Text == "" {
		return nil, errors.InvalidArgument("text is required")
	}

	opts := []grpcopts.CreateOption{
		grpcopts.WithCreateFields(req, CaseCommentMetadata,
			util.DeduplicateFields,
			util.ParseFieldsForEtag,
			util.EnsureIdField,
		),
	}

	if ts := req.GetCreatedAt(); ts != 0 {
		opts = append(opts, grpcopts.WithCreateTime(time.UnixMilli(ts)))
	}

	createOpts, err := grpcopts.NewCreateOptions(ctx, opts...)
	if err != nil {
		return nil, err
	}

	tag, err := etag.EtagOrId(etag.EtagCase, req.CaseEtag)
	if err != nil {
		return nil, errors.InvalidArgument("Invalid etag")
	}
	createOpts.ParentID = tag.GetOid()

	accessMode := auth.Read
	if !createOpts.GetAuthOpts().CheckObacAccess(CaseCommentMetadata.GetParentScopeName(), accessMode) {
		return nil, errors.New("user doesn't have required (EDIT) access to the case", errors.WithCode(codes.PermissionDenied))
	}
	if createOpts.GetAuthOpts().IsRbacCheckRequired(CaseCommentMetadata.GetParentScopeName(), accessMode) {
		access, err := c.app.Store.Case().CheckRbacAccess(createOpts, createOpts.GetAuthOpts(), accessMode, createOpts.ParentID)
		if err != nil {
			return nil, err
		}
		if !access {
			return nil, errors.New("user doesn't have required (EDIT) access to the case", errors.WithCode(codes.PermissionDenied))
		}
	}
	comment, err := c.app.Store.CaseComment().Publish(
		createOpts,
		&cases.CaseComment{
			// Used if explicitly set the case creator / updater instead of deriving it from the auth token.
			CreatedBy: req.Input.GetUserID(),
			Text:      req.Input.Text,
		},
	)
	if err != nil {
		return nil, err
	}

	id := comment.GetId()
	roleId := comment.GetRoleIds()
	parentId := comment.GetCaseId()

	err = NormalizeCommentsResponse(comment, req)
	if err != nil {
		return nil, err
	}
	if notifyErr := c.app.watcherManager.Notify(
		caseCommentsObjScope,
		watcherkit.EventTypeCreate,
		NewCaseCommentWatcherData(createOpts.GetAuthOpts(), comment, id, parentId, roleId),
	); notifyErr != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("could not notify comment create: %s, ", notifyErr.Error()))
	}

	return comment, nil
}

func NormalizeCommentsResponse(res interface{}, opts shared.Fielder) error {
	requestedFields := opts.GetFields()
	if len(requestedFields) == 0 {
		requestedFields = CaseCommentMetadata.GetDefaultFields()
	}
	hasEtag, hasId, hasVer := util.FindEtagFields(requestedFields)
	var err error
	processComment := func(comment *cases.CaseComment) error {
		comment.RoleIds = nil
		comment.CaseId = 0
		if hasEtag {
			comment.Etag, err = etag.EncodeEtag(etag.EtagCaseComment, comment.Id, comment.Ver)
			if err != nil {
				return err
			}
			if !hasId {
				comment.Id = 0
			}
			if !hasVer {
				comment.Ver = 0
			}
		}
		return nil
	}

	switch v := res.(type) {
	case *cases.CaseComment:
		err = processComment(v)
		if err != nil {
			return err
		}
	case *cases.CaseCommentList:
		for _, comment := range v.Items {
			err = processComment(comment)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func formCommentsFtsModel(comment *cases.CaseComment, params map[string]any) (*model.FtsCaseComment, error) {
	roles, ok := params["role_ids"].([]int64)
	if !ok {
		return nil, fmt.Errorf("role ids required for FTS model")
	}
	caseId, ok := params["case_id"].(int64)
	if !ok {
		return nil, fmt.Errorf("case id required for FTS model")
	}

	return &model.FtsCaseComment{
		ParentId:  caseId,
		Comment:   comment.GetText(),
		RoleIds:   roles,
		CreatedAt: comment.GetCreatedAt(),
	}, nil
}

func NewCaseCommentService(app *App) (*CaseCommentService, error) {
	if app == nil {
		return nil, errors.New("Unable to initialize case comment service, app is nil")
	}
	watcher := watcherkit.NewDefaultWatcher()

	service := &CaseCommentService{
		app: app,
	}

	if app.config.LoggerWatcher.Enabled {

		obs, err := NewLoggerObserver(app.wtelLogger, caseCommentsObjScope, defaultLogTimeout)
		if err != nil {
			return nil, err
		}
		watcher.Attach(watcherkit.EventTypeCreate, obs)
		watcher.Attach(watcherkit.EventTypeUpdate, obs)
		watcher.Attach(watcherkit.EventTypeDelete, obs)
	}

	if app.config.FtsWatcher.Enabled {
		ftsObserver, err := NewFullTextSearchObserver(app.ftsClient, caseCommentsObjScope, formCommentsFtsModel)
		if err != nil {
			return nil, err
		}
		watcher.Attach(watcherkit.EventTypeCreate, ftsObserver)
		watcher.Attach(watcherkit.EventTypeUpdate, ftsObserver)
		watcher.Attach(watcherkit.EventTypeDelete, ftsObserver)
	}

	if app.config.TriggerWatcher.Enabled {
		mq, err := NewTriggerObserver(app.rabbitPublisher, app.config.TriggerWatcher, formCaseCommentTriggerModel, slog.With(
			slog.Group("context",
				slog.String("scope", "watcher")),
		))

		if err != nil {
			return nil, err
		}
		watcher.Attach(watcherkit.EventTypeCreate, mq)
		watcher.Attach(watcherkit.EventTypeUpdate, mq)
		watcher.Attach(watcherkit.EventTypeDelete, mq)
		watcher.Attach(watcherkit.EventTypeResolutionTime, mq)

		app.caseResolutionTimer.Start()
	}

	app.watcherManager.AddWatcher(caseCommentsObjScope, watcher)

	return service, nil
}

type CaseCommentWatcherData struct {
	comment *cases.CaseComment
	Args    map[string]any
}

func NewCaseCommentWatcherData(session auth.Auther, comment *cases.CaseComment, id, caseId int64, roleIds []int64) *CaseCommentWatcherData {
	return &CaseCommentWatcherData{comment: comment, Args: map[string]any{"session": session, "obj": comment, "case_id": caseId, "role_ids": roleIds, "id": id}}
}

func (wd *CaseCommentWatcherData) Marshal() ([]byte, error) {
	return json.Marshal(wd.comment)
}

func (wd *CaseCommentWatcherData) GetArgs() map[string]any {
	return wd.Args
}
