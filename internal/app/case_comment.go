package app

import (
	"context"
	"encoding/json"
	"fmt"
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

const caseCommentsObjScope = model.ScopeCaseComments

var CaseCommentMetadata = model.NewObjectMetadata(caseCommentsObjScope, caseObjScope, []*model.Field{
	{Name: "id", Default: false},
	{Name: "etag", Default: true},
	{Name: "ver", Default: false},
	{Name: "created_at", Default: true},
	{Name: "created_by", Default: true},
	{Name: "updated_at", Default: true},
	{Name: "updated_by", Default: false},
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
		return nil, cerror.NewBadRequestError("app.case_comment.locate_comment.etag_required", "Etag is required")
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
		return nil, NewBadRequestError(err)
	}
	logAttributes := slog.Group("context", slog.Int64("user_id", searchOpts.GetAuthOpts().GetUserId()), slog.Int64("domain_id", searchOpts.GetAuthOpts().GetDomainId()))

	commentList, err := c.app.Store.CaseComment().List(searchOpts)
	if err != nil {
		return nil, cerror.NewInternalError("app.case_comment.locate_comment.fetch_error", err.Error())
	}

	if len(commentList.Items) == 0 {
		return nil, cerror.NewNotFoundError("app.case_comment.locate_comment.not_found", "Comment not found")
	} else if len(commentList.Items) > 1 {
		return nil, deferr.InternalError
	}

	err = NormalizeCommentsResponse(commentList.Items[0], req)
	if err != nil {
		slog.ErrorContext(ctx, err.Error(), logAttributes)
		return nil, deferr.ResponseNormalizingError
	}

	return commentList.Items[0], nil
}

func (c *CaseCommentService) UpdateComment(
	ctx context.Context,
	req *cases.UpdateCommentRequest,
) (*cases.CaseComment, error) {
	if req.Input.Etag == "" {
		return nil, cerror.NewBadRequestError("app.case_comment.update_comment.etag_required", "Etag is required")
	}
	if req.Input.Text == "" {
		return nil, cerror.NewBadRequestError("app.case_comment.update_comment.text_required", "Text is required")
	}

	tag, err := etag.EtagOrId(etag.EtagCaseComment, req.Input.Etag)
	if err != nil {
		return nil, cerror.NewBadRequestError("app.case_comment.update_comment.invalid_etag", "Invalid etag")
	}

	updateOpts, err := grpcopts.NewUpdateOptions(
		ctx,
		grpcopts.WithUpdateFields(req, CaseCommentMetadata),
		grpcopts.WithUpdateEtag(&tag),
	)
	if err != nil {
		return nil, NewBadRequestError(err)
	}

	logAttributes := slog.Group("context", slog.Int64("user_id", updateOpts.GetAuthOpts().GetUserId()), slog.Int64("domain_id", updateOpts.GetAuthOpts().GetDomainId()), slog.Int64("id", tag.GetOid()))

	input := &cases.CaseComment{
		// Used if explicitly set the case creator / updater instead of deriving it from the auth token.
		UpdatedBy: req.Input.GetUserID(),
		Id:        tag.GetOid(),
		Text:      req.Input.Text,
		Ver:       tag.GetVer(),
	}

	updatedComment, err := c.app.Store.CaseComment().Update(updateOpts, input)
	if err != nil {
		slog.ErrorContext(ctx, err.Error(), logAttributes)
		return nil, cerror.NewInternalError("app.case_comment.update_comment.store_update_failed", "database error")
	}

	id := updatedComment.GetId()
	roleIds := updatedComment.GetRoleIds()
	parentId := input.GetCaseId()

	err = NormalizeCommentsResponse(updatedComment, req)
	if err != nil {
		slog.ErrorContext(ctx, err.Error(), logAttributes)
		return nil, deferr.ResponseNormalizingError
	}

	if notifyErr := c.app.watcherManager.Notify(caseCommentsObjScope, EventTypeUpdate, NewCaseCommentWatcherData(updateOpts.GetAuthOpts(), updatedComment, id, parentId, roleIds)); notifyErr != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("could not notify input update: %s, ", notifyErr.Error()), logAttributes)
	}

	return updatedComment, nil
}

func (c *CaseCommentService) DeleteComment(
	ctx context.Context,
	req *cases.DeleteCommentRequest,
) (*cases.CaseComment, error) {
	if req.Etag == "" {
		return nil, cerror.NewBadRequestError("app.case_comment.delete_comment.etag_required", "Etag is required")
	}
	tag, err := etag.EtagOrId(etag.EtagCaseComment, req.GetEtag())
	if err != nil {
		return nil, cerror.NewBadRequestError("app.case_comment.delete_comment.invalid_etag", "Invalid etag")
	}
	deleteOpts, err := grpcopts.NewDeleteOptions(ctx, grpcopts.WithDeleteID(tag.GetOid()))
	if err != nil {
		return nil, NewBadRequestError(err)
	}
	logAttributes := slog.Group("context", slog.Int64("user_id", deleteOpts.GetAuthOpts().GetUserId()), slog.Int64("domain_id", deleteOpts.GetAuthOpts().GetDomainId()), slog.Int64("id", tag.GetOid()))

	err = c.app.Store.CaseComment().Delete(deleteOpts)
	if err != nil {
		slog.ErrorContext(ctx, err.Error(), logAttributes)
		return nil, deferr.DatabaseError
	}

	if notifyErr := c.app.watcherManager.Notify(caseCommentsObjScope, EventTypeDelete, NewCaseCommentWatcherData(deleteOpts.GetAuthOpts(), nil, tag.GetOid(), 0, nil)); notifyErr != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("could not notify comment delete: %s, ", notifyErr.Error()), logAttributes)
	}
	return nil, nil
}

func (c *CaseCommentService) ListComments(
	ctx context.Context,
	req *cases.ListCommentsRequest,
) (*cases.CaseCommentList, error) {
	if req.CaseEtag == "" {
		return nil, cerror.NewBadRequestError("app.case_comment.list_comments.case_etag_required", "Case etag is required")
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
		return nil, NewBadRequestError(err)
	}
	tag, err := etag.EtagOrId(etag.EtagCase, req.CaseEtag)
	if err != nil {
		return nil, cerror.NewBadRequestError("app.case_comment.list_comments.invalid_etag", "Invalid etag")
	}
	searchOpts.AddFilter("case_id", tag.GetOid())
	logAttributes := slog.Group("context", slog.Int64("user_id", searchOpts.GetAuthOpts().GetUserId()), slog.Int64("domain_id", searchOpts.GetAuthOpts().GetDomainId()), slog.Int64("case_id", tag.GetOid()))

	comments, err := c.app.Store.CaseComment().List(searchOpts)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return nil, deferr.DatabaseError
	}

	err = NormalizeCommentsResponse(comments, req)
	if err != nil {
		slog.ErrorContext(ctx, err.Error(), logAttributes)
		return nil, deferr.ResponseNormalizingError
	}

	return comments, nil
}

func (c *CaseCommentService) PublishComment(
	ctx context.Context,
	req *cases.PublishCommentRequest,
) (*cases.CaseComment, error) {
	if req.CaseEtag == "" {
		return nil, cerror.NewBadRequestError("app.case_comment.publish_comment.case_etag_required", "Case etag is required")
	} else if req.Input.Text == "" {
		return nil, cerror.NewBadRequestError("app.case_comment.publish_comment.text_required", "Text is required")
	}

	createOpts, err := grpcopts.NewCreateOptions(
		ctx,
		grpcopts.WithCreateFields(req, CaseCommentMetadata,
			util.DeduplicateFields,
			util.ParseFieldsForEtag,
			util.EnsureIdField),
	)
	if err != nil {
		return nil, NewBadRequestError(err)
	}

	tag, err := etag.EtagOrId(etag.EtagCase, req.CaseEtag)
	if err != nil {
		return nil, cerror.NewBadRequestError("app.case_comment.publish_comment.invalid_etag", "Invalid etag")
	}
	createOpts.ParentID = tag.GetOid()
	logAttributes := slog.Group("context", slog.Int64("user_id", createOpts.GetAuthOpts().GetUserId()), slog.Int64("domain_id", createOpts.GetAuthOpts().GetDomainId()), slog.Int64("case_id", tag.GetOid()))

	accessMode := auth.Read
	if !createOpts.GetAuthOpts().CheckObacAccess(CaseCommentMetadata.GetParentScopeName(), accessMode) {
		slog.ErrorContext(ctx, "user doesn't have required (EDIT) access to the case", logAttributes)
		return nil, deferr.ForbiddenError
	}
	if createOpts.GetAuthOpts().IsRbacCheckRequired(CaseCommentMetadata.GetParentScopeName(), accessMode) {
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
	comment, err := c.app.Store.CaseComment().Publish(
		createOpts,
		&cases.CaseComment{
			// Used if explicitly set the case creator / updater instead of deriving it from the auth token.
			CreatedBy: req.Input.GetUserID(),
			Text:      req.Input.Text,
		},
	)
	if err != nil {
		slog.ErrorContext(ctx, err.Error(), logAttributes)
		return nil, deferr.DatabaseError
	}

	id := comment.GetId()
	roleId := comment.GetRoleIds()
	parentId := comment.GetCaseId()

	err = NormalizeCommentsResponse(comment, req)
	if err != nil {
		slog.ErrorContext(ctx, err.Error(), logAttributes)
		return nil, deferr.ResponseNormalizingError
	}
	if notifyErr := c.app.watcherManager.Notify(caseCommentsObjScope, EventTypeCreate, NewCaseCommentWatcherData(createOpts.GetAuthOpts(), comment, id, parentId, roleId)); notifyErr != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("could not notify comment create: %s, ", notifyErr.Error()), logAttributes)
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

func NewCaseCommentService(app *App) (*CaseCommentService, cerror.AppError) {
	if app == nil {
		return nil, cerror.NewInternalError("app.case_comment.new_case_comment_service.app_required", "Unable to initialize service, app is nil")
	}
	watcher := NewDefaultWatcher()
	if app.config.FtsWatcher.Enabled {
		ftsObserver, err := NewFullTextSearchObserver(app.ftsClient, caseCommentsObjScope, formCommentsFtsModel)
		if err != nil {
			return nil, cerror.NewInternalError("app.case.new_case_comment_service.create_observer.app", err.Error())
		}
		watcher.Attach(EventTypeCreate, ftsObserver)
		watcher.Attach(EventTypeUpdate, ftsObserver)
		watcher.Attach(EventTypeDelete, ftsObserver)
	}
	app.watcherManager.AddWatcher(caseCommentsObjScope, watcher)
	return &CaseCommentService{app: app}, nil
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
