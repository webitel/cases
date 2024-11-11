package app

import (
	"context"
	"strings"
	"time"

	cases "github.com/webitel/cases/api/cases"

	cerror "github.com/webitel/cases/internal/error"
	"github.com/webitel/cases/model"
	"github.com/webitel/cases/util"
	"github.com/webitel/webitel-go-kit/etag"
)

const (
	defaultFieldsCaseComments = "id, comment"
)

type CaseCommentService struct {
	app *App
	cases.UnimplementedCaseCommentsServer
}

func (c *CaseCommentService) LocateComment(
	ctx context.Context,
	req *cases.LocateCommentRequest,
) (*cases.CaseComment, error) {
	// Validate required fields
	if req.Etag == "" {
		return nil, cerror.NewBadRequestError("case_comment_service.locate_comment.etag.required", "Etag is required")
	}

	// Get the session from the context
	session, err := c.app.AuthorizeFromContext(ctx)
	if err != nil {
		return nil, cerror.NewUnauthorizedError("case_comment_service.locate_comment.authorization.failed", err.Error())
	}

	// Convert the etag to an internal identifier (Tid) for filtering by ID
	id, err := etag.EtagOrId(etag.EtagCaseComment, req.Etag)
	if err != nil {
		return nil, cerror.NewBadRequestError("case_comment_service.locate_comment.invalid_etag", "Invalid etag")
	}

	fields := util.FieldsFunc(req.Fields, util.InlineFields)

	if len(fields) == 0 {
		fields = strings.Split(defaultFieldsCaseComments, ", ")
	}

	t := time.Now()

	searchOpts := model.SearchOptions{
		IDs:     []int64{id.GetOid()},
		Session: session,
		Fields:  fields,
		Context: ctx,
		Time:    t,
	}

	// Use ListComments to retrieve the specific comment
	commentList, err := c.app.Store.CaseComment().List(&searchOpts)
	if err != nil {
		return nil, cerror.NewInternalError("case_comment_service.locate_comment.fetch_error", err.Error())
	}

	// Ensure we found exactly one comment
	if len(commentList.Items) == 0 {
		return nil, cerror.NewNotFoundError("case_comment_service.locate_comment.not_found", "Comment not found")
	} else if len(commentList.Items) > 1 {
		return nil, cerror.NewInternalError("case_comment_service.locate_comment.multiple_found", "Multiple comments found")
	}

	// Return the located comment
	return commentList.Items[0], nil
}

func (c *CaseCommentService) UpdateComment(
	ctx context.Context,
	req *cases.UpdateCommentRequest,
) (*cases.CaseComment, error) {
	if req.Input.Etag == "" {
		return nil, cerror.NewBadRequestError("case_comment_service.update_comment.etag.required", "Etag is required")
	}
	// Do NOT allow empty text ---- Comment text is required
	if req.Input.Text == "" {
		return nil, cerror.NewBadRequestError("case_comment_service.update_comment.text.required", "Text is required")
	}

	// Set session, xJsonMask, time, fields, ctx
	updateOpts := model.NewUpdateOptions(ctx, req)

	// Prepare the update model
	comment := &cases.CaseComment{
		Id:   req.Input.Etag,
		Text: req.Input.Text,
	}

	// Execute the update in the store
	updatedComment, err := c.app.Store.CaseComment().Update(updateOpts, comment)
	if err != nil {
		return nil, cerror.NewInternalError("case_comment_service.update_comment.store_update_failed", err.Error())
	}

	return updatedComment, nil
}

func (c *CaseCommentService) DeleteComment(
	ctx context.Context,
	req *cases.DeleteCommentRequest,
) (*cases.CaseComment, error) {
	// Validate required fields
	// Etag is required to delete a comment
	if req.Etag == "" {
		return nil, cerror.NewBadRequestError("case_comment_service.delete_comment.etag.required", "Etag is required")
	}

	// Initialize delete options based on the request
	deleteOpts := model.NewDeleteOptions(ctx)

	//  Convert CaseEtag to an internal identifier (Tid) for processing
	id, err := etag.EtagOrId(etag.EtagCaseComment, req.Etag)
	if err != nil {
		return nil, cerror.NewBadRequestError("case_comment_service.delete_comment.invalid_etag", "Invalid etag")
	}
	deleteOpts.IDs = []int64{id.GetOid()}

	// Call the delete method in the store
	err = c.app.Store.CaseComment().Delete(deleteOpts)
	if err != nil {
		return nil, cerror.NewInternalError("case_comment_service.delete_comment.store_delete_failed", err.Error())
	}
	return nil, nil
}

func (c *CaseCommentService) ListComments(
	ctx context.Context,
	req *cases.ListCommentsRequest,
) (*cases.CaseCommentList, error) {
	// Validate required fields
	if req.CaseEtag == "" {
		return nil, cerror.NewBadRequestError("app.case_comment.list_comments.case_etag.required", "Case etag is required")
	}

	// Get the session from the context
	session, err := c.app.AuthorizeFromContext(ctx)
	if err != nil {
		return nil, cerror.NewUnauthorizedError("app.case_comment.list_comments.authorization.failed", err.Error())
	}

	fields := util.FieldsFunc(req.Fields, util.InlineFields)

	if len(fields) == 0 {
		fields = strings.Split(defaultFieldsCaseComments, ", ")
	}

	// Use default page size and page number if not provided
	page := req.Page
	if page == 0 {
		page = 1
	}

	// Convert the etag to an internal identifier (Tid) for filtering by ID
	id, err := etag.EtagOrId(etag.EtagCaseComment, req.CaseEtag)
	if err != nil {
		return nil, cerror.NewBadRequestError("case_comment_service.locate_comment.invalid_etag", "Invalid etag")
	}

	ids, err := util.ParseQin(req.Qin, etag.EtagCaseComment)
	if err != nil {
		return nil, cerror.NewBadRequestError("app.case_comment.list_comments.invalid_qin", "Invalid Qin format")
	}

	t := time.Now()
	searchOpts := model.SearchOptions{
		IDs:     ids,
		Id:      id.GetOid(),
		Session: session,
		Fields:  fields,
		Context: ctx,
		Sort:    []string{req.Sort},
		Page:    int32(page),
		Size:    int32(req.Size),
		Time:    t,
		Filter:  make(map[string]interface{}),
	}

	// Execute search operation to retrieve comments from the database
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
		return nil, cerror.NewBadRequestError("case_comment_service.merge_comments.case_etag.required", "Case etag is required")
	} else if req.Input.Text == "" {
		return nil, cerror.NewBadRequestError("case_comment_service.merge_comments.text.required", "Text is required for each comment")
	}

	fields := util.FieldsFunc(req.Fields, util.InlineFields)

	if len(fields) == 0 {
		fields = strings.Split(defaultFieldsCaseComments, ", ")
	}

	// Initialize search options based on the request
	createOpts := model.NewCreateOptions(ctx, req)
	// Set the fields to return in the response
	createOpts.Fields = fields

	// Get oid of the Case associated with the comments
	caseID, err := etag.EtagOrId(etag.EtagCaseComment, req.CaseEtag)
	if err != nil {
		return nil, cerror.NewBadRequestError("case_comment_service.locate_comment.invalid_etag", "Invalid etag")
	}
	// Set the Case ID to the comment
	createOpts.ParentID = caseID.GetOid()

	comment, err := c.app.Store.CaseComment().Publish(createOpts, &cases.CaseComment{Text: req.Input.Text})
	if err != nil {
		return nil, cerror.NewInternalError("case_comment_service.merge_comments.merge_error", err.Error())
	}

	return comment, nil
}

func NewCaseCommentService(app *App) (*CaseCommentService, cerror.AppError) {
	if app == nil {
		return nil, cerror.NewInternalError("app.case.new_case_comment_service.check_args.app", "unable to init service, app is nil")
	}
	return &CaseCommentService{app: app}, nil
}
