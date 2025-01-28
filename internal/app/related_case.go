package app

import (
	"context"
	"github.com/webitel/cases/auth"
	"log/slog"
	"strconv"

	cases "github.com/webitel/cases/api/cases"
	cerror "github.com/webitel/cases/internal/error"
	"github.com/webitel/cases/model"
	"github.com/webitel/cases/util"
	"github.com/webitel/webitel-go-kit/etag"
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
		return nil, cerror.NewBadRequestError("app.related_case.locate_related_case.id_required", "ID is required")
	}

	tag, err := etag.EtagOrId(etag.EtagRelatedCase, req.GetEtag())
	if err != nil {
		return nil, cerror.NewBadRequestError("app.related_case.locate_related_case.invalid_id", "Invalid ID")
	}
	caseTid, err := etag.EtagOrId(etag.EtagCase, req.GetPrimaryCaseEtag())
	if err != nil {
		return nil, cerror.NewBadRequestError("app.related_case.locate_related_case.invalid_primary_id", "Invalid ID")
	}
	searchOpts, err := model.NewLocateOptions(ctx, req, RelatedCaseMetadata)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return nil, AppInternalError
	}
	searchOpts.IDs = []int64{tag.GetOid()}
	searchOpts.ParentId = caseTid.GetOid()
	logAttributes := slog.Group(
		"context",
		slog.Int64("user_id", searchOpts.GetAuthOpts().GetUserId()),
		slog.Int64("domain_id", searchOpts.GetAuthOpts().GetDomainId()),
		slog.Int64("case_id", searchOpts.ParentId),
	)
	accessMode := auth.Read
	if searchOpts.GetAuthOpts().IsRbacCheckRequired(RelatedCaseMetadata.GetParentScopeName(), accessMode) {
		access, err := r.app.Store.Case().CheckRbacAccess(searchOpts, searchOpts.GetAuthOpts(), accessMode, searchOpts.ParentId)
		if err != nil {
			slog.ErrorContext(ctx, err.Error(), logAttributes)
			return nil, AppForbiddenError
		}
		if !access {
			slog.ErrorContext(ctx, "user doesn't have required (READ) access to the case", logAttributes)
			return nil, AppForbiddenError
		}
	}

	output, err := r.app.Store.RelatedCase().List(searchOpts)
	if err != nil {
		return nil, err
	}
	if len(output.Data) == 0 {
		return nil, cerror.NewNotFoundError("app.related_case.locate_related_case.not_found", "Related case not found")
	} else if len(output.Data) > 1 {
		return nil, cerror.NewInternalError("app.related_case.locate_related_cases.multiple_found", "Multiple related cases found")
	}

	// Normalize IDs and handle errors
	if err := normalizeIDs(output.Data); err != nil {
		return nil, cerror.NewInternalError("app.related_case.locate_related_cases.normalize_ids_error", err.Error())
	}

	return output.Data[0], nil
}

func (r *RelatedCaseService) CreateRelatedCase(ctx context.Context, req *cases.CreateRelatedCaseRequest) (*cases.RelatedCase, error) {
	if req.GetPrimaryCaseEtag() == "" {
		return nil, cerror.NewBadRequestError("app.related_case.create_related_case.primary_case_id_required", "Primary case id required")
	}

	createOpts, err := model.NewCreateOptions(ctx, req, RelatedCaseMetadata)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return nil, AppInternalError
	}
	primaryCaseTag, err := etag.EtagOrId(etag.EtagCase, req.GetPrimaryCaseEtag())
	if err != nil {
		return nil, cerror.NewBadRequestError(
			"app.related_case.created_related_case.invalid_etag",
			"Invalid primary case etag",
		)
	}

	relatedCaseTag, err := etag.EtagOrId(etag.EtagCase, strconv.Itoa(int(req.GetInput().GetRelatedCase().GetId())))
	if err != nil {
		return nil, cerror.NewBadRequestError(
			"app.related_case.created_related_case.invalid_etag",
			"Invalid relatedCase etag",
		)
	}
	createOpts.ParentID = primaryCaseTag.GetOid()
	createOpts.ChildID = relatedCaseTag.GetOid()

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
			return nil, AppForbiddenError
		}
		if !primaryAccess {
			slog.ErrorContext(ctx, "user doesn't have required access to the primary case", logAttributes)
			return nil, AppForbiddenError
		}
	}
	secondaryAccessMode := auth.Read
	if createOpts.GetAuthOpts().IsRbacCheckRequired(RelatedCaseMetadata.GetParentScopeName(), secondaryAccessMode) {
		secondaryAccess, err := r.app.Store.Case().CheckRbacAccess(createOpts, createOpts.GetAuthOpts(), secondaryAccessMode, createOpts.ChildID)
		if err != nil {
			slog.ErrorContext(ctx, err.Error(), logAttributes)
			return nil, AppForbiddenError
		}
		if !secondaryAccess {
			slog.ErrorContext(ctx, "user doesn't have required access to the secondary case", logAttributes)
			return nil, AppForbiddenError
		}
	}

	relatedCase, err := r.app.Store.RelatedCase().Create(createOpts, &req.GetInput().RelationType)
	if err != nil {
		return nil, err
	}

	relatedCase.Etag, err = etag.EncodeEtag(etag.EtagRelatedCase, relatedCase.GetId(), relatedCase.Ver)
	if err != nil {
		slog.ErrorContext(ctx, err.Error(), logAttributes)
		return nil, AppResponseNormalizingError
	}
	relatedCase.RelatedCase.Etag, err = etag.EncodeEtag(etag.EtagCase, relatedCase.RelatedCase.GetId(), relatedCase.Ver)
	if err != nil {
		slog.ErrorContext(ctx, err.Error(), logAttributes)
		return nil, AppResponseNormalizingError
	}
	relatedCase.PrimaryCase.Etag, err = etag.EncodeEtag(etag.EtagCase, relatedCase.PrimaryCase.GetId(), relatedCase.Ver)
	if err != nil {
		slog.ErrorContext(ctx, err.Error(), logAttributes)
		return nil, AppResponseNormalizingError
	}
	return relatedCase, nil
}

func (r *RelatedCaseService) UpdateRelatedCase(ctx context.Context, req *cases.UpdateRelatedCaseRequest) (*cases.RelatedCase, error) {
	if req.GetEtag() == "" {
		return nil, cerror.NewBadRequestError("app.related_case.update_related_case.id_required", "ID required")
	}

	updateOpts, err := model.NewUpdateOptions(ctx, req, RelatedCaseMetadata)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return nil, AppInternalError
	}
	tag, err := etag.EtagOrId(etag.EtagRelatedCase, req.GetEtag())
	if err != nil {
		return nil, cerror.NewBadRequestError(
			"app.related_case.created_related_case.invalid_etag",
			"Invalid ID",
		)
	}
	updateOpts.Etags = []*etag.Tid{&tag}

	primaryCaseTag, err := etag.EtagOrId(etag.EtagCase, strconv.Itoa(int(req.GetInput().GetPrimaryCase().GetId())))
	if err != nil {
		return nil, cerror.NewBadRequestError(
			"app.related_case.created_related_case.invalid_etag",
			"Invalid primary case etag",
		)
	}

	relatedCaseTag, err := etag.EtagOrId(etag.EtagCase, strconv.Itoa(int(req.GetInput().GetRelatedCase().GetId())))
	if err != nil {
		return nil, cerror.NewBadRequestError(
			"app.related_case.created_related_case.invalid_etag",
			"Invalid relatedCase etag",
		)
	}

	if primaryCaseTag.GetOid() == relatedCaseTag.GetOid() {
		return nil, cerror.NewBadRequestError(
			"app.related_case.update_related_case.invalid_ids",
			"A case cannot be related to itself",
		)
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
			return nil, AppForbiddenError
		}
		if !primaryAccess {
			slog.ErrorContext(ctx, "user doesn't have required access to the primary case", logAttributes)
			return nil, AppForbiddenError
		}
	}
	secondaryAccessMode := auth.Read
	if updateOpts.GetAuthOpts().IsRbacCheckRequired(RelatedCaseMetadata.GetParentScopeName(), secondaryAccessMode) {
		secondaryAccess, err := r.app.Store.Case().CheckRbacAccess(updateOpts, updateOpts.GetAuthOpts(), secondaryAccessMode, relatedId)
		if err != nil {
			slog.ErrorContext(ctx, err.Error(), logAttributes)
			return nil, AppForbiddenError
		}
		if !secondaryAccess {
			slog.ErrorContext(ctx, "user doesn't have required access to the secondary case", logAttributes)
			return nil, AppForbiddenError
		}
	}

	output, err := r.app.Store.RelatedCase().Update(updateOpts, input)
	if err != nil {
		return nil, err
	}

	output.Etag, err = etag.EncodeEtag(etag.EtagRelatedCase, output.GetId(), output.GetVer())
	if err != nil {
		slog.ErrorContext(ctx, err.Error(), logAttributes)
		return nil, AppResponseNormalizingError
	}
	return output, nil
}

func (r *RelatedCaseService) DeleteRelatedCase(ctx context.Context, req *cases.DeleteRelatedCaseRequest) (*cases.RelatedCase, error) {
	if req.GetEtag() == "" {
		return nil, cerror.NewBadRequestError("app.related_case.delete_related_case.id_required", "ID required")
	}

	tag, err := etag.EtagOrId(etag.EtagRelatedCase, req.GetEtag())
	if err != nil {
		return nil, cerror.NewBadRequestError("app.related_case.delete_related_case.invalid_etag", "Invalid etag")
	}
	caseEtag, err := etag.EtagOrId(etag.EtagCase, req.GetPrimaryCaseEtag())
	if err != nil {
		return nil, cerror.NewBadRequestError("app.related_case.delete_related_case.invalid_etag", "Invalid etag")
	}

	deleteOpts, err := model.NewDeleteOptions(ctx, RelatedCaseMetadata)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return nil, AppInternalError
	}
	deleteOpts.ID = tag.GetOid()
	deleteOpts.ParentID = caseEtag.GetOid()
	logAttributes := slog.Group(
		"context",
		slog.Int64("user_id", deleteOpts.GetAuthOpts().GetUserId()),
		slog.Int64("domain_id", deleteOpts.GetAuthOpts().GetDomainId()),
		slog.Int64("parent_id", deleteOpts.ParentID),
	)
	accessMode := auth.Edit
	if deleteOpts.GetAuthOpts().IsRbacCheckRequired(RelatedCaseMetadata.GetParentScopeName(), accessMode) {
		access, err := r.app.Store.Case().CheckRbacAccess(deleteOpts, deleteOpts.GetAuthOpts(), accessMode, deleteOpts.ParentID)
		if err != nil {
			slog.ErrorContext(ctx, err.Error(), logAttributes)
			return nil, AppForbiddenError
		}
		if !access {
			slog.ErrorContext(ctx, "user doesn't have required (READ) access to the case", logAttributes)
			return nil, AppForbiddenError
		}

	}

	err = r.app.Store.RelatedCase().Delete(deleteOpts)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (r *RelatedCaseService) ListRelatedCases(ctx context.Context, req *cases.ListRelatedCasesRequest) (*cases.RelatedCaseList, error) {
	if req.GetPrimaryCaseEtag() == "" {
		return nil, cerror.NewBadRequestError("app.related_case.list_related_case.id_required", "ID required")
	}
	tag, err := etag.EtagOrId(etag.EtagCase, req.PrimaryCaseEtag)
	if err != nil {
		return nil, cerror.NewBadRequestError("app.related_case.list_related_cases.invalid_etag", "Invalid etag")
	}

	ids, err := util.ParseIds(req.Ids, etag.EtagRelatedCase)
	if err != nil {
		return nil, cerror.NewBadRequestError("app.related_case.list_related_cases.invalid_ids", "Invalid ids format")
	}
	searchOpts, err := model.NewSearchOptions(ctx, req, RelatedCaseMetadata)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return nil, AppInternalError
	}
	searchOpts.ParentId = tag.GetOid()
	searchOpts.IDs = ids

	output, err := r.app.Store.RelatedCase().List(searchOpts)
	if err != nil {
		return nil, err
	}

	// Normalize IDs and handle errors
	if err := normalizeIDs(output.Data); err != nil {
		return nil, cerror.NewInternalError("app.related_case.list_related_cases.normalize_ids_error", err.Error())
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

func NewRelatedCaseService(app *App) (*RelatedCaseService, cerror.AppError) {
	if app == nil {
		return nil, cerror.NewBadRequestError("app.case.new_related_case_service.check_args.app", "unable to init service, app is nil")
	}
	return &RelatedCaseService{app: app}, nil
}
