package app

import (
	"context"
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

	NormalizeCommentsResponse(commentList.Items[0], req)

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

	updateOpts := model.NewUpdateOptions(ctx, req, CaseCommentMetadata)
	updateOpts.Etags = []*etag.Tid{&tag}

	comment := &cases.CaseComment{
		Id:   strconv.Itoa(int(tag.GetOid())),
		Text: req.Input.Text,
		Ver:  tag.GetVer(),
	}

	updatedComment, err := c.app.Store.CaseComment().Update(updateOpts, comment)
	if err != nil {
		return nil, cerror.NewInternalError("app.case_comment.update_comment.store_update_failed", err.Error())
	}

	NormalizeCommentsResponse(updatedComment, req)
	return updatedComment, nil
}

func (c *CaseCommentService) DeleteComment(
	ctx context.Context,
	req *cases.DeleteCommentRequest,
) (*cases.CaseComment, error) {
	if req.Id == "" {
		return nil, cerror.NewBadRequestError("app.case_comment.delete_comment.etag_required", "ID is required")
	}

	deleteOpts := model.NewDeleteOptions(ctx)

	tag, err := etag.EtagOrId(etag.EtagCaseComment, req.Id)
	if err != nil {
		return nil, cerror.NewBadRequestError("app.case_comment.delete_comment.invalid_etag", "Invalid ID")
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
	searchOpts := model.NewSearchOptions(ctx, req, CaseCommentMetadata)
	searchOpts.ParentId = tag.GetOid()
	searchOpts.IDs = ids

	comments, err := c.app.Store.CaseComment().List(searchOpts)
	if err != nil {
		return nil, cerror.NewInternalError("app.case_comment.list_comments.fetch_error", err.Error())
	}

	NormalizeCommentsResponse(comments, req)

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

	createOpts := model.NewCreateOptions(ctx, req, CaseCommentMetadata)

	tag, err := etag.EtagOrId(etag.EtagCaseComment, req.CaseId)
	if err != nil {
		return nil, cerror.NewBadRequestError("app.case_comment.publish_comment.invalid_etag", "Invalid ID")
	}
	createOpts.ParentID = tag.GetOid()

	comment, err := c.app.Store.CaseComment().Publish(createOpts, &cases.CaseComment{Text: req.Input.Text})
	if err != nil {
		return nil, err
	}

	NormalizeCommentsResponse(comment, req)
	return comment, nil
}

func NormalizeCommentsResponse(res interface{}, opts model.Fielder) {
	fields := util.FieldsFunc(opts.GetFields(), util.InlineFields)
	if len(fields) == 0 {
		fields = CaseCommentMetadata.GetDefaultFields()
	}
	hasEtag, hasId, hasVer := util.FindEtagFields(fields)

	processComment := func(comment *cases.CaseComment) {
		if hasEtag {
			id, _ := strconv.Atoi(comment.Id)
			comment.Id = etag.EncodeEtag(etag.EtagCaseComment, int64(id), comment.Ver)
			// if NOT provided in requested fields - hide them in response
			if !hasId {
				comment.Id = ""
			}
			if !hasVer {
				comment.Ver = 0
			}
		}
	}

	switch v := res.(type) {
	case *cases.CaseComment:
		processComment(v)
	case *cases.CaseCommentList:
		for _, comment := range v.Items {
			processComment(comment)
		}
	}
}

func NewCaseCommentService(app *App) (*CaseCommentService, cerror.AppError) {
	if app == nil {
		return nil, cerror.NewInternalError("app.case_comment.new_case_comment_service.app_required", "Unable to initialize service, app is nil")
	}
	return &CaseCommentService{app: app}, nil
}
