package app

import (
	"context"
	"strconv"
	"time"

	cases "github.com/webitel/cases/api/cases"
	cerror "github.com/webitel/cases/internal/error"
	"github.com/webitel/cases/model"
	"github.com/webitel/cases/util"
	"github.com/webitel/webitel-go-kit/etag"
)

var defaultFieldsCaseComments = []string{"id", "comment"}

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

	fields := util.FieldsFunc(req.Fields, util.InlineFields)
	if len(fields) == 0 {
		fields = defaultFieldsCaseComments
	}

	searchOpts := model.NewLocateOptions(ctx, req)
	searchOpts.IDs = []int64{tag.GetOid()}
	searchOpts.Fields = fields

	commentList, err := c.app.Store.CaseComment().List(searchOpts)
	if err != nil {
		return nil, cerror.NewInternalError("app.case_comment.locate_comment.fetch_error", err.Error())
	}

	if len(commentList.Items) == 0 {
		return nil, cerror.NewNotFoundError("app.case_comment.locate_comment.not_found", "Comment not found")
	} else if len(commentList.Items) > 1 {
		return nil, cerror.NewInternalError("app.case_comment.locate_comment.multiple_found", "Multiple comments found")
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

	fields := util.FieldsFunc(req.Fields, util.InlineFields)
	if len(fields) == 0 {
		fields = defaultFieldsCaseComments
	}

	tag, err := etag.EtagOrId(etag.EtagCaseComment, req.Input.Etag)
	if err != nil {
		return nil, cerror.NewBadRequestError("app.case_comment.update_comment.invalid_etag", "Invalid etag")
	}

	updateOpts := model.NewUpdateOptions(ctx, req)
	updateOpts.IDs = []int64{tag.GetOid()}
	updateOpts.Fields = fields

	comment := &cases.CaseComment{
		Id:   strconv.Itoa(int(tag.GetOid())),
		Text: req.Input.Text,
		Ver:  tag.GetVer(),
	}

	updatedComment, err := c.app.Store.CaseComment().Update(updateOpts, comment)
	if err != nil {
		return nil, cerror.NewInternalError("app.case_comment.update_comment.store_update_failed", err.Error())
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
		return nil, cerror.NewInternalError("app.case_comment.delete_comment.store_delete_failed", err.Error())
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

	session, err := c.app.AuthorizeFromContext(ctx)
	if err != nil {
		return nil, cerror.NewUnauthorizedError("app.case_comment.list_comments.authorization_failed", err.Error())
	}

	fields := util.FieldsFunc(req.Fields, util.InlineFields)
	if len(fields) == 0 {
		fields = defaultFieldsCaseComments
	}

	page := req.Page
	if page == 0 {
		page = 1
	}

	tag, err := etag.EtagOrId(etag.EtagCaseComment, req.CaseEtag)
	if err != nil {
		return nil, cerror.NewBadRequestError("app.case_comment.list_comments.invalid_etag", "Invalid etag")
	}

	ids, err := util.ParseQin(req.Qin, etag.EtagCaseComment)
	if err != nil {
		return nil, cerror.NewBadRequestError("app.case_comment.list_comments.invalid_qin", "Invalid Qin format")
	}

	t := time.Now()
	searchOpts := model.SearchOptions{
		IDs:     ids,
		Id:      tag.GetOid(),
		Session: session,
		Fields:  fields,
		Context: ctx,
		Sort:    []string{req.Sort},
		Page:    int32(page),
		Size:    int32(req.Size),
		Time:    t,
		Filter:  make(map[string]interface{}),
	}

	comments, err := c.app.Store.CaseComment().List(&searchOpts)
	if err != nil {
		return nil, cerror.NewInternalError("app.case_comment.list_comments.fetch_error", err.Error())
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

	fields := util.FieldsFunc(req.Fields, util.InlineFields)
	if len(fields) == 0 {
		fields = defaultFieldsCaseComments
	}

	createOpts := model.NewCreateOptions(ctx, req)
	createOpts.Fields = fields

	tag, err := etag.EtagOrId(etag.EtagCaseComment, req.CaseEtag)
	if err != nil {
		return nil, cerror.NewBadRequestError("app.case_comment.publish_comment.invalid_etag", "Invalid etag")
	}
	createOpts.ParentID = tag.GetOid()

	comment, err := c.app.Store.CaseComment().Publish(createOpts, &cases.CaseComment{Text: req.Input.Text})
	if err != nil {
		return nil, cerror.NewInternalError("app.case_comment.publish_comment.publish_error", err.Error())
	}

	return comment, nil
}

func NewCaseCommentService(app *App) (*CaseCommentService, cerror.AppError) {
	if app == nil {
		return nil, cerror.NewInternalError("app.case_comment.new_case_comment_service.app_required", "Unable to initialize service, app is nil")
	}
	return &CaseCommentService{app: app}, nil
}
