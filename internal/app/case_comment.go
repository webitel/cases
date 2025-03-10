package app

import (
	"context"
	"errors"
	"log/slog"

	"github.com/webitel/cases/api/cases"
	"github.com/webitel/cases/auth"
	cerror "github.com/webitel/cases/internal/errors"
	"github.com/webitel/cases/model"
	"github.com/webitel/cases/util"
	"github.com/webitel/webitel-go-kit/etag"
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

	tag, err := etag.EtagOrId(etag.EtagCaseComment, req.Etag)
	if err != nil {
		return nil, cerror.NewBadRequestError("app.case_comment.locate_comment.invalid_etag", "Invalid etag")
	}

	searchOpts, err := model.NewLocateOptions(ctx, req, CaseCommentMetadata)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return nil, InternalError
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
		slog.ErrorContext(ctx, err.Error(), logAttributes)
		return nil, ResponseNormalizingError
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

	updateOpts, err := model.NewUpdateOptions(ctx, req, CaseCommentMetadata.CopyWithAllFieldsSetToDefault())
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return nil, InternalError
	}
	updateOpts.Etags = []*etag.Tid{&tag}
	logAttributes := slog.Group("context", slog.Int64("user_id", updateOpts.GetAuthOpts().GetUserId()), slog.Int64("domain_id", updateOpts.GetAuthOpts().GetDomainId()), slog.Int64("id", tag.GetOid()))

	comment := &cases.CaseComment{
		Id:   tag.GetOid(),
		Text: req.Input.Text,
		Ver:  tag.GetVer(),
	}

	updatedComment, err := c.app.Store.CaseComment().Update(updateOpts, comment)
	if err != nil {
		slog.ErrorContext(ctx, err.Error(), logAttributes)
		return nil, cerror.NewInternalError("app.case_comment.update_comment.store_update_failed", "database error")
	}

	id := updatedComment.GetId()
	roleIds := updatedComment.GetRoleIds()
	parentId := comment.GetCaseId()

	err = NormalizeCommentsResponse(updatedComment, req)
	if err != nil {
		slog.ErrorContext(ctx, err.Error(), logAttributes)
		return nil, ResponseNormalizingError
	}
	ftsErr := c.SendFtsUpdateEvent(id, updateOpts.GetAuthOpts().GetDomainId(), roleIds, parentId, updatedComment)
	if ftsErr != nil {
		slog.ErrorContext(ctx, ftsErr.Error(), logAttributes)
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

	deleteOpts, err := model.NewDeleteOptions(ctx, CaseCommentMetadata.CopyWithAllFieldsSetToDefault())
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return nil, InternalError
	}
	tag, err := etag.EtagOrId(etag.EtagCaseComment, req.GetEtag())
	if err != nil {
		return nil, cerror.NewBadRequestError("app.case_comment.delete_comment.invalid_etag", "Invalid etag")
	}
	deleteOpts.IDs = []int64{tag.GetOid()}
	logAttributes := slog.Group("context", slog.Int64("user_id", deleteOpts.GetAuthOpts().GetUserId()), slog.Int64("domain_id", deleteOpts.GetAuthOpts().GetDomainId()), slog.Int64("id", tag.GetOid()))

	err = c.app.Store.CaseComment().Delete(deleteOpts)
	if err != nil {
		slog.ErrorContext(ctx, err.Error(), logAttributes)
		return nil, DatabaseError
	}

	ftsErr := c.SendFtsDeleteEvent(tag.GetOid(), deleteOpts.GetAuthOpts().GetDomainId())
	if ftsErr != nil {
		slog.ErrorContext(ctx, ftsErr.Error(), logAttributes)
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
	searchOpts, err := model.NewSearchOptions(ctx, req, CaseCommentMetadata)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return nil, InternalError
	}
	searchOpts.ParentId = tag.GetOid()
	searchOpts.IDs = ids
	logAttributes := slog.Group("context", slog.Int64("user_id", searchOpts.GetAuthOpts().GetUserId()), slog.Int64("domain_id", searchOpts.GetAuthOpts().GetDomainId()), slog.Int64("case_id", tag.GetOid()))

	comments, err := c.app.Store.CaseComment().List(searchOpts)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return nil, DatabaseError
	}

	err = NormalizeCommentsResponse(comments, req)
	if err != nil {
		slog.ErrorContext(ctx, err.Error(), logAttributes)
		return nil, ResponseNormalizingError
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

	createOpts, err := model.NewCreateOptions(ctx, req, CaseCommentMetadata.CopyWithAllFieldsSetToDefault())
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return nil, InternalError
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
		return nil, ForbiddenError
	}
	if createOpts.GetAuthOpts().IsRbacCheckRequired(CaseCommentMetadata.GetParentScopeName(), accessMode) {
		access, err := c.app.Store.Case().CheckRbacAccess(createOpts, createOpts.GetAuthOpts(), accessMode, createOpts.ParentID)
		if err != nil {
			slog.ErrorContext(ctx, err.Error(), logAttributes)
			return nil, ForbiddenError
		}
		if !access {
			slog.ErrorContext(ctx, "user doesn't have required (EDIT) access to the case", logAttributes)
			return nil, ForbiddenError
		}
	}
	comment, err := c.app.Store.CaseComment().Publish(createOpts, &cases.CaseComment{Text: req.Input.Text})
	if err != nil {
		slog.ErrorContext(ctx, err.Error(), logAttributes)
		return nil, DatabaseError
	}

	id := comment.GetId()
	roleId := comment.GetRoleIds()
	parentId := comment.GetCaseId()

	err = NormalizeCommentsResponse(comment, req)
	if err != nil {
		slog.ErrorContext(ctx, err.Error(), logAttributes)
		return nil, ResponseNormalizingError
	}
	ftsErr := c.SendFtsCreateEvent(id, createOpts.GetAuthOpts().GetDomainId(), roleId, parentId, comment)
	if ftsErr != nil {
		slog.ErrorContext(ctx, ftsErr.Error(), logAttributes)
	}
	return comment, nil
}

func NormalizeCommentsResponse(res interface{}, opts model.Fielder) error {
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

func (c *CaseCommentService) SendFtsCreateEvent(id int64, domainId int64, roleIds []int64, caseId int64, comment *cases.CaseComment) error {
	if domainId == 0 {
		return errors.New("domain id required")
	}
	if id == 0 {
		return errors.New("id required")
	}
	m, err := c.formFtsModel(roleIds, caseId, comment)
	if err != nil {
		return err
	}
	return c.app.ftsClient.Create(domainId, model.ScopeCaseComments, id, m)
}

func (c *CaseCommentService) SendFtsUpdateEvent(id int64, domainId int64, roleIds []int64, caseId int64, comment *cases.CaseComment) error {
	if domainId == 0 {
		return errors.New("domain id required")
	}
	if id == 0 {
		return errors.New("id required")
	}
	m, err := c.formFtsModel(roleIds, caseId, comment)
	if err != nil {
		return err
	}
	return c.app.ftsClient.Update(domainId, model.ScopeCaseComments, id, m)
}

func (c *CaseCommentService) SendFtsDeleteEvent(id int64, domainId int64) error {
	if domainId == 0 {
		return errors.New("domain id required")
	}
	if id == 0 {
		return errors.New("id required")
	}
	return c.app.ftsClient.Delete(domainId, model.ScopeCaseComments, id)
}

func (c *CaseCommentService) formFtsModel(roleIds []int64, caseId int64, comment *cases.CaseComment) (*model.FtsCaseComment, error) {
	if caseId == 0 {
		return nil, errors.New("case id required")
	}
	return &model.FtsCaseComment{
		ParentId:  caseId,
		Comment:   comment.GetText(),
		RoleIds:   roleIds,
		CreatedAt: comment.GetCreatedAt(),
	}, nil
}

func NewCaseCommentService(app *App) (*CaseCommentService, cerror.AppError) {
	if app == nil {
		return nil, cerror.NewInternalError("app.case_comment.new_case_comment_service.app_required", "Unable to initialize service, app is nil")
	}
	return &CaseCommentService{app: app}, nil
}
