package app

import (
	"context"

	"github.com/webitel/cases/api/cases"

	cerror "github.com/webitel/cases/internal/error"
)

type CaseCommentService struct {
	app *App
}

func (c *CaseCommentService) LocateComment(ctx context.Context, request *cases.LocateCommentRequest) (*cases.CaseComment, error) {
	// TODO implement me
	panic("implement me")
}

func (c *CaseCommentService) UpdateComment(ctx context.Context, request *cases.UpdateCommentRequest) (*cases.CaseComment, error) {
	// TODO implement me
	panic("implement me")
}

func (c *CaseCommentService) DeleteComment(ctx context.Context, request *cases.DeleteCommentRequest) (*cases.CaseComment, error) {
	// TODO implement me
	panic("implement me")
}

func (c *CaseCommentService) ListComments(ctx context.Context, request *cases.ListCommentsRequest) (*cases.CaseCommentList, error) {
	// TODO implement me
	panic("implement me")
}

func (c *CaseCommentService) MergeComments(ctx context.Context, request *cases.MergeCommentsRequest) (*cases.CaseCommentList, error) {
	// TODO implement me
	panic("implement me")
}

func (c *CaseCommentService) ResetComments(ctx context.Context, request *cases.ResetCommentsRequest) (*cases.CaseCommentList, error) {
	// TODO implement me
	panic("implement me")
}

func NewCaseCommentService(app *App) (*CaseCommentService, cerror.AppError) {
	if app == nil {
		return nil, cerror.NewBadRequestError("app.case.new_case_comment_service.check_args.app", "unable to init service, app is nil")
	}
	return &CaseCommentService{app: app}, nil
}
