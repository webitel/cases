package grpc

import (
	"context"

	api "github.com/webitel/cases/api/cases"
	"github.com/webitel/cases/internal/api_handler/grpc/utils"
	"github.com/webitel/cases/internal/errors"
	"github.com/webitel/cases/internal/model"
	"github.com/webitel/cases/internal/model/options"
	grpcopts "github.com/webitel/cases/internal/model/options/grpc"
	"github.com/webitel/cases/internal/model/options/grpc/shared"
	"github.com/webitel/cases/util"
	"github.com/webitel/webitel-go-kit/pkg/etag"
	"google.golang.org/grpc/codes"
)

// CaseCommentHandler defines the interface for managing case comments.
type CaseCommentHandler interface {
	ListCaseComments(options.Searcher) ([]*model.CaseComment, error)
	UpdateCaseComment(options.Updator, *model.CaseComment) (*model.CaseComment, error)
	DeleteCaseComment(options.Deleter) (*model.CaseComment, error)
	PublishCaseComment(options.Creator, *model.CaseComment) (*model.CaseComment, error)
}

// CaseCommentService implements the gRPC server for case comments.
type CaseCommentService struct {
	app CaseCommentHandler
	api.UnimplementedCaseCommentsServer
}

// NewCaseCommentService constructs a new CaseCommentService.
func NewCaseCommentService(app CaseCommentHandler) (*CaseCommentService, error) {
	if app == nil {
		return nil, errors.New("case comment handler is nil")
	}
	return &CaseCommentService{app: app}, nil
}

// CaseCommentMetadata defines the fields available for case comment objects.
var CaseCommentMetadata = model.NewObjectMetadata(model.ScopeCaseComments, model.ScopeCases, []*model.Field{
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

// LocateComment handles the gRPC request to locate a comment by its etag.
func (s *CaseCommentService) LocateComment(ctx context.Context, req *api.LocateCommentRequest) (*api.CaseComment, error) {
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

	comments, err := s.app.ListCaseComments(searchOpts)
	if err != nil {
		return nil, err
	}

	if len(comments) == 0 {
		return nil, errors.NotFound("Comment not found")
	} else if len(comments) > 1 {
		return nil, errors.New("too many items found", errors.WithCode(codes.AlreadyExists))
	}

	comment := comments[0]

	result, err := s.Marshal(comment)
	if err != nil {
		return nil, err
	}

	err = s.NormalizeResponse(result, req)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// UpdateComment handles the gRPC request to update a comment.
func (s *CaseCommentService) UpdateComment(ctx context.Context, req *api.UpdateCommentRequest) (*api.CaseComment, error) {
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

	updateOpts, err := grpcopts.NewUpdateOptions(
		ctx,
		grpcopts.WithUpdateFields(req, CaseCommentMetadata.CopyWithAllFieldsSetToDefault()),
		grpcopts.WithUpdateEtag(&tag),
		grpcopts.WithUpdateMasker(req),
	)
	if err != nil {
		return nil, err
	}

	input := &model.CaseComment{
		Id:   tag.GetOid(),
		Ver:  tag.GetVer(),
		Text: req.Input.Text,
	}

	// Set user ID if provided
	if req.Input.GetUserID() != nil {
		userId := int(req.Input.GetUserID().Id)
		userName := req.Input.GetUserID().Name
		input.Editor = &model.Editor{
			Id:   &userId,
			Name: &userName,
		}
	}

	updatedComment, err := s.app.UpdateCaseComment(updateOpts, input)
	if err != nil {
		return nil, err
	}

	result, err := s.Marshal(updatedComment)
	if err != nil {
		return nil, err
	}

	err = s.NormalizeResponse(result, req)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// DeleteComment handles the gRPC request to delete a comment.
func (s *CaseCommentService) DeleteComment(ctx context.Context, req *api.DeleteCommentRequest) (*api.CaseComment, error) {
	if req.Etag == "" {
		return nil, errors.InvalidArgument("etag is required")
	}

	tag, err := etag.EtagOrId(etag.EtagCaseComment, req.GetEtag())
	if err != nil {
		return nil, errors.InvalidArgument("invalid Etag", errors.WithCause(err))
	}

	deleteOpts, err := grpcopts.NewDeleteOptions(ctx, grpcopts.WithDeleteID(tag.GetOid()), grpcopts.WithDeleteFields(req, CaseCommentMetadata.CopyWithAllFieldsSetToDefault(), util.ParseFieldsForEtag))
	if err != nil {
		return nil, err
	}

	deletedComment, err := s.app.DeleteCaseComment(deleteOpts)
	if err != nil {
		return nil, err
	}

	result, err := s.Marshal(deletedComment)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// ListComments handles the gRPC request to list comments for a case.
func (s *CaseCommentService) ListComments(ctx context.Context, req *api.ListCommentsRequest) (*api.CaseCommentList, error) {
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
		searchOpts.AddFilter(util.EqualFilter("case_id", tag.GetOid()))
	}

	comments, err := s.app.ListCaseComments(searchOpts)
	if err != nil {
		return nil, err
	}

	var res api.CaseCommentList
	res.Items, err = utils.ConvertToOutputBulk(comments, s.Marshal)
	if err != nil {
		return nil, err
	}

	res.Next, res.Items = utils.GetListResult(searchOpts, res.Items)
	res.Page = int64(req.GetPage())

	err = s.NormalizeResponse(&res, req)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

// PublishComment handles the gRPC request to publish a comment to a case.
func (s *CaseCommentService) PublishComment(ctx context.Context, req *api.PublishCommentRequest) (*api.CaseComment, error) {
	if req.CaseEtag == "" {
		return nil, errors.InvalidArgument("case etag is required")
	} else if req.Input.Text == "" {
		return nil, errors.InvalidArgument("text is required")
	}

	createOpts, err := grpcopts.NewCreateOptions(
		ctx,
		grpcopts.WithCreateFields(req, CaseCommentMetadata.CopyWithAllFieldsSetToDefault(),
			util.DeduplicateFields,
			util.ParseFieldsForEtag,
			util.EnsureIdField),
	)
	if err != nil {
		return nil, err
	}

	tag, err := etag.EtagOrId(etag.EtagCase, req.CaseEtag)
	if err != nil {
		return nil, errors.InvalidArgument("Invalid etag")
	}
	createOpts.ParentID = tag.GetOid()

	input := &model.CaseComment{
		Text:   req.Input.Text,
		CaseId: createOpts.ParentID,
	}

	// Set user ID if provided
	if req.Input.GetUserID() != nil {
		userId := int(req.Input.GetUserID().Id)
		userName := req.Input.GetUserID().Name
		input.Author = &model.Author{
			Id:   &userId,
			Name: &userName,
		}
	}

	comment, err := s.app.PublishCaseComment(createOpts, input)
	if err != nil {
		return nil, err
	}

	result, err := s.Marshal(comment)
	if err != nil {
		return nil, err
	}

	err = s.NormalizeResponse(result, req)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Marshal converts a model.CaseComment to its gRPC representation.
func (s *CaseCommentService) Marshal(model *model.CaseComment) (*api.CaseComment, error) {
	if model == nil {
		return nil, nil
	}

	comment := &api.CaseComment{
		Id:      model.Id,
		Ver:     model.Ver,
		Etag:    model.Etag,
		Text:    model.Text,
		Edited:  model.Edited,
		CanEdit: model.CanEdit,
		CaseId:  model.CaseId,
	}

	if model.CreatedAt != nil {
		comment.CreatedAt = model.CreatedAt.Unix() * 1000 // Convert to milliseconds
	}
	if model.UpdatedAt != nil {
		comment.UpdatedAt = model.UpdatedAt.Unix() * 1000 // Convert to milliseconds
	}

	if model.Author != nil {
		comment.CreatedBy = utils.MarshalLookup(model.Author)
	}
	if model.Editor != nil {
		comment.UpdatedBy = utils.MarshalLookup(model.Editor)
	}
	if model.Contact != nil {
		comment.Author = utils.MarshalLookup(model.Contact)
	}

	return comment, nil
}

// NormalizeResponse normalizes the response based on requested fields and etag handling.
func (s *CaseCommentService) NormalizeResponse(res interface{}, opts shared.Fielder) error {
	requestedFields := opts.GetFields()
	if len(requestedFields) == 0 {
		requestedFields = CaseCommentMetadata.GetDefaultFields()
	}
	hasEtag, hasId, hasVer := util.FindEtagFields(requestedFields)
	var err error

	processComment := func(comment *api.CaseComment) error {
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
	case *api.CaseComment:
		err = processComment(v)
		if err != nil {
			return err
		}
	case *api.CaseCommentList:
		for _, comment := range v.Items {
			err = processComment(comment)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
