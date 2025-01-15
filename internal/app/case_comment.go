package app

import (
	"context"
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
	[]*model.Field{
		{Name: "id", Default: true},
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

	tag, err := etag.EtagOrId(etag.EtagCaseComment, req.Etag)
	if err != nil {
		return nil, cerror.NewBadRequestError("app.case_comment.locate_comment.invalid_etag", "Invalid etag")
	}

	searchOpts := model.NewLocateOptions(ctx, req, CaseCommentMetadata)
	searchOpts.IDs = []int64{tag.GetOid()}

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
		slog.Warn(err.Error(), slog.Int64("user_id", searchOpts.Session.GetUserId()), slog.Int64("domain_id", searchOpts.Session.GetDomainId()))
		return nil, AppResponseNormalizingError
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

	updateOpts := model.NewUpdateOptions(ctx, req, CaseCommentMetadata)
	updateOpts.Etags = []*etag.Tid{&tag}

	comment := &cases.CaseComment{
		Id:   tag.GetOid(),
		Text: req.Input.Text,
		Ver:  tag.GetVer(),
	}

	updatedComment, err := c.app.Store.CaseComment().Update(updateOpts, comment)
	if err != nil {
		slog.Warn(err.Error(), slog.Int64("user_id", updateOpts.Session.GetUserId()), slog.Int64("domain_id", updateOpts.Session.GetDomainId()), slog.Int64("case_id", tag.GetOid()))
		return nil, cerror.NewInternalError("app.case_comment.update_comment.store_update_failed", "database error")
	}

	err = NormalizeCommentsResponse(updatedComment, req)
	if err != nil {
		slog.Warn(err.Error(), slog.Int64("user_id", updateOpts.Session.GetUserId()), slog.Int64("domain_id", updateOpts.Session.GetDomainId()))
		return nil, AppResponseNormalizingError
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

	deleteOpts := model.NewDeleteOptions(ctx)

	tag, err := etag.EtagOrId(etag.EtagCaseComment, req.Etag)
	if err != nil {
		return nil, cerror.NewBadRequestError("app.case_comment.delete_comment.invalid_etag", "Invalid etag")
	}
	deleteOpts.IDs = []int64{tag.GetOid()}

	err = c.app.Store.CaseComment().Delete(deleteOpts)
	if err != nil {
		slog.Warn(err.Error(), slog.Int64("user_id", deleteOpts.Session.GetUserId()), slog.Int64("domain_id", deleteOpts.Session.GetDomainId()), slog.Int64("case_id", tag.GetOid()))
		return nil, AppDatabaseError
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

	tag, err := etag.EtagOrId(etag.EtagCase, req.CaseEtag)
	if err != nil {
		return nil, cerror.NewBadRequestError("app.case_comment.list_comments.invalid_etag", "Invalid etag")
	}

	ids, err := util.ParseIds(req.Ids, etag.EtagCaseComment)
	if err != nil {
		return nil, cerror.NewBadRequestError("app.case_comment.list_comments.invalid_qin", "Invalid Qin format")
	}
	searchOpts := model.NewSearchOptions(ctx, req, CaseCommentMetadata)
	searchOpts.ParentId = tag.GetOid()
	searchOpts.IDs = ids

	comments, err := c.app.Store.CaseComment().List(searchOpts)
	if err != nil {
		slog.Warn(err.Error(), slog.Int64("user_id", searchOpts.Session.GetUserId()), slog.Int64("domain_id", searchOpts.Session.GetDomainId()), slog.Int64("case_id", tag.GetOid()))
		return nil, AppDatabaseError
	}

	err = NormalizeCommentsResponse(comments, req)
	if err != nil {
		slog.Warn(err.Error(), slog.Int64("user_id", searchOpts.Session.GetUserId()), slog.Int64("domain_id", searchOpts.Session.GetDomainId()), slog.Int64("case_id", tag.GetOid()))
		return nil, AppResponseNormalizingError
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

	createOpts := model.NewCreateOptions(ctx, req, CaseCommentMetadata)

	tag, err := etag.EtagOrId(etag.EtagCase, req.CaseEtag)
	if err != nil {
		return nil, cerror.NewBadRequestError("app.case_comment.publish_comment.invalid_etag", "Invalid etag")
	}
	createOpts.ParentID = tag.GetOid()

	comment, err := c.app.Store.CaseComment().Publish(createOpts, &cases.CaseComment{Text: req.Input.Text})
	if err != nil {
		slog.Warn(err.Error(), slog.Int64("user_id", createOpts.Session.GetUserId()), slog.Int64("domain_id", createOpts.Session.GetDomainId()), slog.Int64("case_id", tag.GetOid()))
		return nil, errors.NewInternalError("app.case_comment.publish_comment.database.exec", "database error")
	}

	err = NormalizeCommentsResponse(comment, req)
	if err != nil {
		slog.Warn(err.Error(), slog.Int64("user_id", createOpts.Session.GetUserId()), slog.Int64("domain_id", createOpts.Session.GetDomainId()), slog.Int64("case_id", tag.GetOid()))
		return nil, AppResponseNormalizingError
	}
	return comment, nil
}

func NormalizeCommentsResponse(res interface{}, opts model.Fielder) error {
	processComment := func(comment *cases.CaseComment) error {

		if hasEtag {
			comment.Etag, err = etag.EncodeEtag(etag.EtagCaseComment, comment.Id, comment.Ver)
			if err != nil {
				return err
			}
			// if NOT provided in requested fields - hide them in response
			if !hasId {
				comment.Id = 0
			}
			if !hasVer {
				comment.Ver = 0
			}
		}
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
