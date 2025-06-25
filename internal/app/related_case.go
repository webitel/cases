package app

import (
	"context"
	"github.com/webitel/cases/api/cases"
	"github.com/webitel/cases/auth"
	"github.com/webitel/cases/internal/errors"
	"github.com/webitel/cases/internal/model"
	grpcopts "github.com/webitel/cases/internal/model/options/grpc"
	"github.com/webitel/cases/util"
	"github.com/webitel/webitel-go-kit/pkg/etag"
	"log/slog"
	"strconv"
)

type RelatedCaseService struct {
	app *App
	cases.UnimplementedRelatedCasesServer
}

var RelatedCaseMetadata = model.NewObjectMetadata("", caseObjScope, []*model.Field{
	{Name: "id", Default: true},
	{Name: "ver", Default: true},
	{Name: "created_at", Default: true},
	{Name: "created_by", Default: true},
	{Name: "updated_at", Default: false},
	{Name: "updated_by", Default: false},
	{Name: "related_case", Default: true},
	{Name: "primary_case", Default: true},
	{Name: "relation", Default: true},
})

func (r *RelatedCaseService) LocateRelatedCase(ctx context.Context, req *cases.LocateRelatedCaseRequest) (*cases.RelatedCase, error) {
	if req.GetEtag() == "" {
		return nil, errors.InvalidArgument("ID is required")
	}
	caseTid, err := etag.EtagOrId(etag.EtagCase, req.GetPrimaryCaseEtag())
	if err != nil {
		return nil, errors.InvalidArgument("Invalid primary case etag", errors.WithCause(err))
	}
	searchOpts, err := grpcopts.NewLocateOptions(
		ctx,
		grpcopts.WithFields(req, CaseCommentMetadata,
			util.DeduplicateFields,
			func(in []string) []string {
				return util.EnsureFields(in, "created_at", "id")
			},
		),
		grpcopts.WithIDsAsEtags(etag.EtagRelatedCase, req.GetEtag()),
	)
	if err != nil {
		return nil, errors.InvalidArgument("Invalid search options", errors.WithCause(err))
	}
	searchOpts.AddFilter("case_id", caseTid.GetOid())
	logAttributes := slog.Group(
		"context",
		slog.Int64("user_id", searchOpts.GetAuthOpts().GetUserId()),
		slog.Int64("domain_id", searchOpts.GetAuthOpts().GetDomainId()),
		slog.Int64("case_id", caseTid.GetOid()),
	)
	accessMode := auth.Read
	if searchOpts.GetAuthOpts().IsRbacCheckRequired(RelatedCaseMetadata.GetParentScopeName(), accessMode) {
		access, err := r.app.Store.Case().CheckRbacAccess(searchOpts, searchOpts.GetAuthOpts(), accessMode, caseTid.GetOid())
		if err != nil {
			slog.ErrorContext(ctx, err.Error(), logAttributes)
			return nil, errors.Forbidden("user doesn't have required (READ) access to the case", errors.WithCause(err))
		}
		if !access {
			slog.ErrorContext(ctx, "user doesn't have required (READ) access to the case", logAttributes)
			return nil, errors.Forbidden("user doesn't have required (READ) access to the case")
		}
	}

	output, err := r.app.Store.RelatedCase().List(searchOpts)
	if err != nil {
		return nil, err
	}
	if len(output.Data) == 0 {
		return nil, errors.NotFound("Related case not found")
	} else if len(output.Data) > 1 {
		return nil, errors.Internal("Multiple related cases found")
	}

	// Normalize IDs and handle errors
	if err := normalizeIDs(output.Data); err != nil {
		return nil, errors.Internal("Failed to normalize related case IDs", errors.WithCause(err))
	}
	return output.Data[0], nil
}

func (r *RelatedCaseService) CreateRelatedCase(ctx context.Context, req *cases.CreateRelatedCaseRequest) (*cases.RelatedCase, error) {
	if req.GetPrimaryCaseEtag() == "" {
		return nil, errors.InvalidArgument("Primary case id required")
	}

	primaryCaseTag, err := etag.EtagOrId(etag.EtagCase, req.GetPrimaryCaseEtag())
	if err != nil {
		return nil, errors.InvalidArgument("Invalid primary case etag", errors.WithCause(err))
	}

	relatedCaseTag, err := etag.EtagOrId(etag.EtagCase, strconv.Itoa(int(req.GetInput().GetRelatedCase().GetId())))
	if err != nil {
		return nil, errors.InvalidArgument("Invalid related case etag", errors.WithCause(err))
	}

	createOpts, err := grpcopts.NewCreateOptions(
		ctx,
		grpcopts.WithCreateFields(req, RelatedCaseMetadata,
			util.DeduplicateFields,
			util.ParseFieldsForEtag,
			util.EnsureIdField),
		grpcopts.WithCreateParentID(primaryCaseTag.GetOid()),
		grpcopts.WithCreateChildID(relatedCaseTag.GetOid()),
	)
	if err != nil {
		return nil, errors.InvalidArgument("Invalid create options", errors.WithCause(err))
	}

	logAttributes := slog.Group(
		"context",
		slog.Int64("user_id", createOpts.GetAuthOpts().GetUserId()),
		slog.Int64("domain_id", createOpts.GetAuthOpts().GetDomainId()),
		slog.Int64("parent_id", createOpts.ParentID),
		slog.Int64("child_id", createOpts.ChildID),
	)
	primaryAccessMode := auth.Edit
	if createOpts.GetAuthOpts().IsRbacCheckRequired(RelatedCaseMetadata.GetParentScopeName(), primaryAccessMode) {
		primaryAccess, err := r.app.Store.Case().CheckRbacAccess(createOpts, createOpts.GetAuthOpts(), primaryAccessMode, createOpts.ParentID)
		if err != nil {
			slog.ErrorContext(ctx, err.Error(), logAttributes)
			return nil, errors.Forbidden("user doesn't have required (EDIT) access to the primary case", errors.WithCause(err))
		}
		if !primaryAccess {
			slog.ErrorContext(ctx, "user doesn't have required access to the primary case", logAttributes)
			return nil, errors.Forbidden("user doesn't have required (EDIT) access to the primary case")
		}
	}
	secondaryAccessMode := auth.Read
	if createOpts.GetAuthOpts().IsRbacCheckRequired(RelatedCaseMetadata.GetParentScopeName(), secondaryAccessMode) {
		secondaryAccess, err := r.app.Store.Case().CheckRbacAccess(createOpts, createOpts.GetAuthOpts(), secondaryAccessMode, createOpts.ChildID)
		if err != nil {
			slog.ErrorContext(ctx, err.Error(), logAttributes)
			return nil, errors.Forbidden("user doesn't have required (READ) access to the secondary case", errors.WithCause(err))
		}
		if !secondaryAccess {
			slog.ErrorContext(ctx, "user doesn't have required access to the secondary case", logAttributes)
			return nil, errors.Forbidden("user doesn't have required (READ) access to the secondary case")
		}
	}

	relatedCase, err := r.app.Store.RelatedCase().Create(
		createOpts,
		&req.GetInput().RelationType,
		req.Input.GetUserID().GetId(),
	)
	if err != nil {
		return nil, err
	}

	relatedCase.Etag, err = etag.EncodeEtag(etag.EtagRelatedCase, relatedCase.GetId(), relatedCase.Ver)
	if err != nil {
		slog.ErrorContext(ctx, err.Error(), logAttributes)
		return nil, errors.Internal("Failed to encode related case etag", errors.WithCause(err))
	}
	relatedCase.RelatedCase.Etag, err = etag.EncodeEtag(etag.EtagCase, relatedCase.RelatedCase.GetId(), relatedCase.Ver)
	if err != nil {
		slog.ErrorContext(ctx, err.Error(), logAttributes)
		return nil, errors.Internal("Failed to encode related case etag", errors.WithCause(err))
	}
	relatedCase.PrimaryCase.Etag, err = etag.EncodeEtag(etag.EtagCase, relatedCase.PrimaryCase.GetId(), relatedCase.Ver)
	if err != nil {
		slog.ErrorContext(ctx, err.Error(), logAttributes)
		return nil, errors.Internal("Failed to encode primary case etag", errors.WithCause(err))
	}
	return relatedCase, nil
}

func (r *RelatedCaseService) UpdateRelatedCase(ctx context.Context, req *cases.UpdateRelatedCaseRequest) (*cases.RelatedCase, error) {
	if req.GetEtag() == "" {
		return nil, errors.InvalidArgument("ID required")
	}

	tag, err := etag.EtagOrId(etag.EtagRelatedCase, req.GetEtag())
	if err != nil {
		return nil, errors.InvalidArgument("Invalid related case etag", errors.WithCause(err))
	}
	updateOpts, err := grpcopts.NewUpdateOptions(
		ctx,
		grpcopts.WithUpdateFields(req, RelatedCaseMetadata),
		grpcopts.WithUpdateEtag(&tag),
		grpcopts.WithUpdateMasker(req),
	)
	if err != nil {
		return nil, errors.InvalidArgument("Invalid update options", errors.WithCause(err))
	}

	primaryCaseTag, err := etag.EtagOrId(etag.EtagCase, strconv.Itoa(int(req.GetInput().GetPrimaryCase().GetId())))
	if err != nil {
		return nil, errors.InvalidArgument("Invalid primary case etag", errors.WithCause(err))
	}

	relatedCaseTag, err := etag.EtagOrId(etag.EtagCase, strconv.Itoa(int(req.GetInput().GetRelatedCase().GetId())))
	if err != nil {
		return nil, errors.InvalidArgument("Invalid related case etag", errors.WithCause(err))
	}

	if primaryCaseTag.GetOid() == relatedCaseTag.GetOid() {
		return nil, errors.InvalidArgument("A case cannot be related to itself")
	}

	input := &cases.InputRelatedCase{
		PrimaryCase:  req.Input.GetPrimaryCase(),
		RelatedCase:  req.Input.GetRelatedCase(),
		RelationType: req.Input.RelationType,
	}

	primaryId := primaryCaseTag.GetOid()
	relatedId := relatedCaseTag.GetOid()
	logAttributes := slog.Group(
		"context",
		slog.Int64("user_id", updateOpts.GetAuthOpts().GetUserId()),
		slog.Int64("domain_id", updateOpts.GetAuthOpts().GetDomainId()),
		slog.Int64("parent_id", updateOpts.ParentID),
	)
	primaryAccessMode := auth.Edit
	if updateOpts.GetAuthOpts().IsRbacCheckRequired(RelatedCaseMetadata.GetParentScopeName(), primaryAccessMode) {
		primaryAccess, err := r.app.Store.Case().CheckRbacAccess(updateOpts, updateOpts.GetAuthOpts(), primaryAccessMode, primaryId)
		if err != nil {
			slog.ErrorContext(ctx, err.Error(), logAttributes)
			return nil, errors.Forbidden("user doesn't have required (EDIT) access to the primary case", errors.WithCause(err))
		}
		if !primaryAccess {
			slog.ErrorContext(ctx, "user doesn't have required access to the primary case", logAttributes)
			return nil, errors.Forbidden("user doesn't have required (EDIT) access to the primary case")
		}
	}
	secondaryAccessMode := auth.Read
	if updateOpts.GetAuthOpts().IsRbacCheckRequired(RelatedCaseMetadata.GetParentScopeName(), secondaryAccessMode) {
		secondaryAccess, err := r.app.Store.Case().CheckRbacAccess(updateOpts, updateOpts.GetAuthOpts(), secondaryAccessMode, relatedId)
		if err != nil {
			slog.ErrorContext(ctx, err.Error(), logAttributes)
			return nil, errors.Forbidden("user doesn't have required (READ) access to the secondary case", errors.WithCause(err))
		}
		if !secondaryAccess {
			slog.ErrorContext(ctx, "user doesn't have required access to the secondary case", logAttributes)
			return nil, errors.Forbidden("user doesn't have required (READ) access to the secondary case")
		}
	}

	output, err := r.app.Store.RelatedCase().Update(
		updateOpts,
		input,
		req.Input.GetUserID().GetId(),
	)
	if err != nil {
		return nil, err
	}

	output.Etag, err = etag.EncodeEtag(etag.EtagRelatedCase, output.GetId(), output.GetVer())
	if err != nil {
		slog.ErrorContext(ctx, err.Error(), logAttributes)
		return nil, errors.Internal("Failed to encode related case etag", errors.WithCause(err))
	}
	return output, nil
}

func (r *RelatedCaseService) DeleteRelatedCase(ctx context.Context, req *cases.DeleteRelatedCaseRequest) (*cases.RelatedCase, error) {
	if req.GetEtag() == "" {
		return nil, errors.InvalidArgument("ID required")
	}
	if req.GetPrimaryCaseEtag() == "" {
		return nil, errors.InvalidArgument("Primary case ID required")
	}
	deleteOpts, err := grpcopts.NewDeleteOptions(ctx, grpcopts.WithDeleteIDsAsEtags(etag.EtagRelatedCase, req.GetEtag()), grpcopts.WithDeleteParentIDAsEtag(etag.EtagCase, req.GetPrimaryCaseEtag()))
	if err != nil {
		return nil, errors.InvalidArgument("Invalid delete options", errors.WithCause(err))
	}
	logAttributes := slog.Group(
		"context",
		slog.Int64("user_id", deleteOpts.GetAuthOpts().GetUserId()),
		slog.Int64("domain_id", deleteOpts.GetAuthOpts().GetDomainId()),
		slog.Int64("parent_id", deleteOpts.ParentID),
	)

	accessMode := auth.Edit
	if deleteOpts.GetAuthOpts().IsRbacCheckRequired(RelatedCaseMetadata.GetParentScopeName(), accessMode) {
		access, err := r.app.Store.Case().CheckRbacAccess(deleteOpts, deleteOpts.GetAuthOpts(), accessMode, deleteOpts.GetParentID())
		if err != nil {
			slog.ErrorContext(ctx, err.Error(), logAttributes)
			return nil, errors.Forbidden("user doesn't have required (EDIT) access to the case", errors.WithCause(err))
		}
		if !access {
			slog.ErrorContext(ctx, "user doesn't have required (EDIT) access to the case", logAttributes)
			return nil, errors.Forbidden("user doesn't have required (EDIT) access to the case")
		}

	}

	err = r.app.Store.RelatedCase().Delete(deleteOpts)
	if err != nil {
		slog.ErrorContext(ctx, err.Error(), logAttributes)
		return nil, err
	}
	return nil, nil
}

func (r *RelatedCaseService) ListRelatedCases(ctx context.Context, req *cases.ListRelatedCasesRequest) (*cases.RelatedCaseList, error) {
	if req.GetPrimaryCaseEtag() == "" {
		return nil, errors.InvalidArgument("Primary case ID required")
	}
	searchOpts, err := grpcopts.NewSearchOptions(
		ctx,
		grpcopts.WithSearch(req),
		grpcopts.WithPagination(req),
		grpcopts.WithFields(req, RelatedCaseMetadata,
			util.DeduplicateFields,
			func(in []string) []string {
				return util.EnsureFields(in, "created_at", "id")
			},
		),
		grpcopts.WithSort(req),
		grpcopts.WithIDsAsEtags(etag.EtagRelatedCase, req.GetIds()...),
	)
	if err != nil {
		return nil, errors.InvalidArgument("Invalid search options", errors.WithCause(err))
	}
	tag, err := etag.EtagOrId(etag.EtagCase, req.PrimaryCaseEtag)
	if err != nil {
		return nil, errors.InvalidArgument("Invalid primary case etag", errors.WithCause(err))
	}
	searchOpts.AddFilter("case_id", tag.GetOid())

	output, err := r.app.Store.RelatedCase().List(searchOpts)
	if err != nil {
		return nil, err
	}

	// Normalize IDs and handle errors
	if err := normalizeIDs(output.Data); err != nil {
		return nil, errors.Internal("Failed to normalize related case IDs", errors.WithCause(err))
	}
	return output, nil
}

func normalizeIDs(relatedCases []*cases.RelatedCase) error {
	for _, relatedCase := range relatedCases {
		if relatedCase == nil {
			continue
		}
		var err error
		// Normalize related case ID
		relatedCase.Etag, err = etag.EncodeEtag(etag.EtagRelatedCase, relatedCase.GetId(), relatedCase.Ver)
		if err != nil {
			return err
		}

		// Normalize primary case ID
		if relatedCase.PrimaryCase != nil {

			relatedCase.PrimaryCase.Etag, err = etag.EncodeEtag(etag.EtagCase, relatedCase.PrimaryCase.GetId(), relatedCase.PrimaryCase.GetVer())
			if err != nil {
				return err
			}
			// Set PrimaryCase Ver to nil
			relatedCase.PrimaryCase.Ver = 0
		}

		// Normalize related case ID inside related case
		if relatedCase.RelatedCase != nil {
			relatedCase.RelatedCase.Etag, err = etag.EncodeEtag(etag.EtagCase, relatedCase.RelatedCase.Id, relatedCase.RelatedCase.GetVer())
			if err != nil {
				return err
			}
			// Set RelatedCase Ver to nil
			relatedCase.RelatedCase.Ver = 0
		}
	}

	return nil
}

func NewRelatedCaseService(app *App) (*RelatedCaseService, error) {
	if app == nil {
		return nil, errors.InvalidArgument("unable to init service, app is nil")
	}
	return &RelatedCaseService{app: app}, nil
}
