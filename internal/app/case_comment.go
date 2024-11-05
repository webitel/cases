package app

import (
	"context"

	cases "buf.build/gen/go/webitel/cases/protocolbuffers/go"
	"google.golang.org/protobuf/proto"

	authmodel "github.com/webitel/cases/auth/model"
	casegraph "github.com/webitel/cases/internal/app/graph"
	cerror "github.com/webitel/cases/internal/error"
	"github.com/webitel/cases/model"
	"github.com/webitel/cases/model/graph"
)

type CaseCommentService struct {
	app *App
}

func (c *CaseCommentService) LocateComment(ctx context.Context, req *cases.LocateCommentRequest) (*cases.CaseComment, error) {
	// TODO implement me
	panic("implement me")
}

func (c *CaseCommentService) UpdateComment(ctx context.Context, req *cases.UpdateCommentRequest) (*cases.CaseComment, error) {
	// TODO implement me
	panic("implement me")
}

func (c *CaseCommentService) DeleteComment(ctx context.Context, req *cases.DeleteCommentRequest) (*cases.CaseComment, error) {
	// TODO implement me
	panic("implement me")
}

func (c *CaseCommentService) ListComments(
	ctx context.Context,
	req *cases.ListCommentsRequest,
) (*cases.CaseCommentList, error) {
	// Validate required fields
	if req.CaseEtag == "" {
		return nil, cerror.NewBadRequestError("app.case_comment.list_comments.case_etag.required", "Case etag is required")
	}

	// Parse and validate the `etag` using the provided `EtagOrId` function
	// caseTag, err := etag.EtagOrId(etag.EtagCaseComment, req.CaseEtag)
	// if err != nil {
	// 	return nil, cerror.NewBadRequestError("app.case_comment.list_comments.invalid_case_etag", "Invalid case etag")
	// }

	// Initialize search options based on the request
	searchOpts := model.NewSearchOptions(ctx, req)

	// Setup GraphQL Query for fields parsing and output
	graphQ := struct {
		FieldsParse func(vs []string, decode ...graph.FieldEncoding) (fields graph.FieldsQ, err error)
		Output      func(*cases.CaseCommentList, *graph.Query)
		graph.Query
	}{
		Query: graph.Query{
			Name: "listComments",
		},
		FieldsParse: casegraph.Schema.Case.Comment.Output.ParseFields,
	}

	// Parse requested fields
	graphParsedFields, err := graphQ.FieldsParse(searchOpts.Fields)
	if err != nil {
		return nil, cerror.NewInternalError("app.case_comment.list_comments.fields_parse", err.Error())
	}
	graphQ.Fields = graphParsedFields

	session, err := c.app.AuthorizeFromContext(ctx)
	if err != nil {
		return nil, cerror.NewUnauthorizedError("app.case_comment.list_comments.no_session", "No session found")
	}
	// OBAC check
	accessMode := authmodel.Read
	scope := session.GetScope(model.ScopeDictionary)
	if !session.HasObacAccess(scope.Class, accessMode) {
		return nil, cerror.MakeScopeError(session.GetUserId(), scope.Class, int(accessMode))
	}

	// Execute search operation to retrieve comments from the database
	comments, err := c.app.Store.CommentCase().List(ctx, searchOpts)
	if err != nil {
		return nil, cerror.NewInternalError("app.case_comment.list_comments.fetch_error", err.Error())
	}

	// Prepare response
	resp := &cases.CaseCommentList{}
	graphQ.Output(resp, &graphQ.Query)

	// Step 10: Populate the response with fetched comments
	for _, comment := range comments.Items {
		// Create a new protobuf CaseComment to merge each comment into
		protoComment := &cases.CaseComment{}
		proto.Merge(protoComment, comment)
		resp.Items = append(resp.Items, protoComment)
	}

	return resp, nil
}

func (c *CaseCommentService) MergeComments(ctx context.Context, req *cases.MergeCommentsRequest) (*cases.CaseCommentList, error) {
	// TODO implement me
	panic("implement me")
}

func NewCaseCommentService(app *App) (*CaseCommentService, cerror.AppError) {
	if app == nil {
		return nil, cerror.NewInternalError("app.case.new_case_comment_service.check_args.app", "unable to init service, app is nil")
	}
	return &CaseCommentService{app: app}, nil
}
