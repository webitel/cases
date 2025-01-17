package app

import (
	"context"
	authmodel "github.com/webitel/cases/auth/model"
	"github.com/webitel/webitel-go-kit/errors"
	"log/slog"
	"strconv"

	cases "github.com/webitel/cases/api/cases"
	cerror "github.com/webitel/cases/internal/error"
	"github.com/webitel/cases/model"
	"github.com/webitel/cases/util"
	"github.com/webitel/webitel-go-kit/etag"
)

var CaseCommentMetadata = model.NewObjectMetadata(
	"case_comments",
	[]*model.Field{
		{Name: "id", Default: true},
		{Name: "ver", Default: false},
		{Name: "created_at", Default: true},
		{Name: "created_by", Default: true},
		{Name: "updated_at", Default: true},
		{Name: "updated_by", Default: false},
		{Name: "text", Default: true},
		{Name: "edited", Default: true},
		{Name: "can_edit", Default: true},
		{Name: "author", Default: true},
	})

type CaseCommentService struct {
	app *App
	cases.UnimplementedCaseCommentsServer
}

func (c *CaseCommentService) LocateComment(
	ctx context.Context,
	req *cases.LocateCommentRequest,
) (*cases.CaseComment, error) {
	if req.Id == "" {
		return nil, cerror.NewBadRequestError("app.case_comment.locate_comment.etag_required", "ID is required")
	}

	tag, err := etag.EtagOrId(etag.EtagCaseComment, req.Id)
	if err != nil {
		return nil, cerror.NewBadRequestError("app.case_comment.locate_comment.invalid_etag", "Invalid ID")
	}

	searchOpts, err := model.NewLocateOptions(ctx, req, CaseCommentMetadata)
	if err != nil {
		slog.Error(err.Error())
		return nil, AppInternalError
	}
	searchOpts.IDs = []int64{tag.GetOid()}
	logAttributes := slog.Group("context", slog.Int64("user_id", searchOpts.GetAuthOpts().GetUserId()), slog.Int64("domain_id", searchOpts.GetAuthOpts().GetDomainId()))

	commentList, err := c.app.Store.CaseComment().List(searchOpts)
	if err != nil {
		return nil, cerror.NewInternalError("app.case_comment.locate_comment.fetch_error", err.Error())
	}

	if len(commentList.Items) == 0 {
		return nil, cerror.NewNotFoundError("app.case_comment.locate_comment.not_found", "Comment not found")
	} else if len(commentList.Items) > 1 {
		return nil, cerror.NewInternalError("app.case_comment.locate_comment.multiple_found", "Multiple comments found")
	}

	err = NormalizeCommentsResponse(commentList.Items[0], req)
	if err != nil {
		slog.Error(err.Error(), logAttributes)
		return nil, AppResponseNormalizingError
	}

	return commentList.Items[0], nil
}

func (c *CaseCommentService) UpdateComment(
	ctx context.Context,
	req *cases.UpdateCommentRequest,
) (*cases.CaseComment, error) {
	if req.Input.Id == "" {
		return nil, cerror.NewBadRequestError("app.case_comment.update_comment.etag_required", "ID is required")
	}
	if req.Input.Text == "" {
		return nil, cerror.NewBadRequestError("app.case_comment.update_comment.text_required", "Text is required")
	}

	tag, err := etag.EtagOrId(etag.EtagCaseComment, req.Input.Id)
	if err != nil {
		return nil, cerror.NewBadRequestError("app.case_comment.update_comment.invalid_etag", "Invalid ID")
	}

	updateOpts, err := model.NewUpdateOptions(ctx, req, CaseCommentMetadata)
	if err != nil {
		slog.Error(err.Error())
		return nil, AppInternalError
	}
	updateOpts.Etags = []*etag.Tid{&tag}
	logAttributes := slog.Group("context", slog.Int64("user_id", updateOpts.GetAuthOpts().GetUserId()), slog.Int64("domain_id", updateOpts.GetAuthOpts().GetDomainId()), slog.Int64("id", tag.GetOid()))

	comment := &cases.CaseComment{
		Id:   strconv.Itoa(int(tag.GetOid())),
		Text: req.Input.Text,
		Ver:  tag.GetVer(),
	}

	updatedComment, err := c.app.Store.CaseComment().Update(updateOpts, comment)
	if err != nil {
		slog.Error(err.Error(), logAttributes)
		return nil, cerror.NewInternalError("app.case_comment.update_comment.store_update_failed", "database error")
	}

	err = NormalizeCommentsResponse(updatedComment, req)
	if err != nil {
		slog.Error(err.Error(), logAttributes)
		return nil, AppResponseNormalizingError
	}
	return updatedComment, nil
}

func (c *CaseCommentService) DeleteComment(
	ctx context.Context,
	req *cases.DeleteCommentRequest,
) (*cases.CaseComment, error) {
	if req.Id == "" {
		return nil, cerror.NewBadRequestError("app.case_comment.delete_comment.etag_required", "ID is required")
	}

	deleteOpts, err := model.NewDeleteOptions(ctx, CaseCommentMetadata)
	if err != nil {
		slog.Error(err.Error())
		return nil, AppInternalError
	}
	tag, err := etag.EtagOrId(etag.EtagCaseComment, req.Id)
	if err != nil {
		return nil, cerror.NewBadRequestError("app.case_comment.delete_comment.invalid_etag", "Invalid ID")
	}
	deleteOpts.IDs = []int64{tag.GetOid()}
	logAttributes := slog.Group("context", slog.Int64("user_id", deleteOpts.GetAuthOpts().GetUserId()), slog.Int64("domain_id", deleteOpts.GetAuthOpts().GetDomainId()), slog.Int64("id", tag.GetOid()))

	err = c.app.Store.CaseComment().Delete(deleteOpts)
	if err != nil {
		slog.Error(err.Error(), logAttributes)
		return nil, AppDatabaseError
	}
	return nil, nil
}

func (c *CaseCommentService) ListComments(
	ctx context.Context,
	req *cases.ListCommentsRequest,
) (*cases.CaseCommentList, error) {
	if req.CaseId == "" {
		return nil, cerror.NewBadRequestError("app.case_comment.list_comments.case_etag_required", "Case ID is required")
	}

	tag, err := etag.EtagOrId(etag.EtagCase, req.CaseId)
	if err != nil {
		return nil, cerror.NewBadRequestError("app.case_comment.list_comments.invalid_etag", "Invalid ID")
	}

	ids, err := util.ParseIds(req.Ids, etag.EtagCaseComment)
	if err != nil {
		return nil, cerror.NewBadRequestError("app.case_comment.list_comments.invalid_qin", "Invalid Qin format")
	}
	searchOpts, err := model.NewSearchOptions(ctx, req, CaseCommentMetadata)
	if err != nil {
		slog.Error(err.Error())
		return nil, AppInternalError
	}
	searchOpts.ParentId = tag.GetOid()
	searchOpts.IDs = ids
	logAttributes := slog.Group("context", slog.Int64("user_id", searchOpts.GetAuthOpts().GetUserId()), slog.Int64("domain_id", searchOpts.GetAuthOpts().GetDomainId()), slog.Int64("case_id", tag.GetOid()))

	comments, err := c.app.Store.CaseComment().List(searchOpts)
	if err != nil {
		slog.Error(err.Error())
		return nil, AppDatabaseError
	}

	err = NormalizeCommentsResponse(comments, req)
	if err != nil {
		slog.Error(err.Error(), logAttributes)
		return nil, AppResponseNormalizingError
	}

	return comments, nil
}

func (c *CaseCommentService) PublishComment(
	ctx context.Context,
	req *cases.PublishCommentRequest,
) (*cases.CaseComment, error) {
	if req.CaseId == "" {
		return nil, cerror.NewBadRequestError("app.case_comment.publish_comment.case_etag_required", "Case ID is required")
	} else if req.Input.Text == "" {
		return nil, cerror.NewBadRequestError("app.case_comment.publish_comment.text_required", "Text is required")
	}

	createOpts, err := model.NewCreateOptions(ctx, req, CaseCommentMetadata)
	if err != nil {
		slog.Error(err.Error())
		return nil, AppInternalError
	}

	tag, err := etag.EtagOrId(etag.EtagCase, req.CaseId)
	if err != nil {
		return nil, cerror.NewBadRequestError("app.case_comment.publish_comment.invalid_etag", "Invalid ID")
	}
	createOpts.ParentID = tag.GetOid()
	logAttributes := slog.Group("context", slog.Int64("user_id", createOpts.GetAuthOpts().GetUserId()), slog.Int64("domain_id", createOpts.GetAuthOpts().GetDomainId()), slog.Int64("case_id", tag.GetOid()))
	if createOpts.GetAuthOpts().GetObjectScope(CaseCommunicationMetadata.GetMainScopeName()).IsRbacUsed() {
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
	comment, err := c.app.Store.CaseComment().Publish(createOpts, &cases.CaseComment{Text: req.Input.Text})
	if err != nil {
		slog.Error(err.Error(), logAttributes)
		return nil, errors.NewInternalError("app.case_comment.publish_comment.database.exec", "database error")
	}

	err = NormalizeCommentsResponse(comment, req)
	if err != nil {
		slog.Error(err.Error(), logAttributes)
		return nil, AppResponseNormalizingError
	}
	return comment, nil
}

func NormalizeCommentsResponse(res interface{}, opts model.Fielder) error {
	processComment := func(comment *cases.CaseComment) error {

		id, err := strconv.Atoi(comment.Id)
		if err != nil {
			return err
		}
		comment.Id, err = etag.EncodeEtag(etag.EtagCaseComment, int64(id), comment.Ver)
		if err != nil {
			return err
		}
		comment.Ver = 0
		return nil
	}

	switch v := res.(type) {
	case *cases.CaseComment:
		err := processComment(v)
		if err != nil {
			return err
		}
	case *cases.CaseCommentList:
		for _, comment := range v.Items {
			err := processComment(comment)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func NewCaseCommentService(app *App) (*CaseCommentService, cerror.AppError) {
	if app == nil {
		return nil, cerror.NewInternalError("app.case_comment.new_case_comment_service.app_required", "Unable to initialize service, app is nil")
	}
	return &CaseCommentService{app: app}, nil
}
