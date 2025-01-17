package app

import (
	"context"
	"fmt"
	authmodel "github.com/webitel/cases/auth/model"
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

var RelatedCaseMetadata = model.NewObjectMetadata(
	"cases",
	[]*model.Field{
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
	if req.GetId() == "" {
		return nil, cerror.NewBadRequestError("app.related_case.locate_related_case.id_required", "ID is required")
	}

	tag, err := etag.EtagOrId(etag.EtagRelatedCase, req.Id)
	if err != nil {
		return nil, cerror.NewBadRequestError("app.related_case.locate_related_case.invalid_id", "Invalid ID")
	}
	caseTid, err := etag.EtagOrId(etag.EtagCase, req.GetPrimaryCaseId())
	if err != nil {
		return nil, cerror.NewBadRequestError("app.related_case.locate_related_case.invalid_primary_id", "Invalid ID")
	}
	searchOpts, err := model.NewLocateOptions(ctx, req, RelatedCaseMetadata)
	if err != nil {
		slog.Error(err.Error())
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
	if searchOpts.GetAuthOpts().GetObjectScope(RelatedCaseMetadata.GetMainScopeName()).IsRbacUsed() {
		access, err := r.app.Store.Case().CheckRbacAccess(searchOpts, searchOpts.GetAuthOpts(), authmodel.Read, searchOpts.ParentId)
		if err != nil {
			slog.Error(err.Error(), logAttributes)
			return nil, AppForbiddenError
		}
		if !access {
			slog.Error("user doesn't have required (READ) access to the case", logAttributes)
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
	if req.GetPrimaryCaseId() == "" {
		return nil, cerror.NewBadRequestError("app.related_case.create_related_case.primary_case_id_required", "Primary case id required")
	}

	createOpts, err := model.NewCreateOptions(ctx, req, RelatedCaseMetadata)
	if err != nil {
		slog.Error(err.Error())
		return nil, AppInternalError
	}
	primaryCaseTag, err := etag.EtagOrId(etag.EtagCase, req.GetPrimaryCaseId())
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
	if createOpts.GetAuthOpts().GetObjectScope(RelatedCaseMetadata.GetMainScopeName()).IsRbacUsed() {
		primaryAccess, err := r.app.Store.Case().CheckRbacAccess(createOpts, createOpts.GetAuthOpts(), authmodel.Edit, createOpts.ParentID)
		if err != nil {
			slog.Error(err.Error(), logAttributes)
			return nil, AppForbiddenError
		}
		secondaryAccess, err := r.app.Store.Case().CheckRbacAccess(createOpts, createOpts.GetAuthOpts(), authmodel.Read, createOpts.ChildID)
		if err != nil {
			slog.Error(err.Error(), logAttributes)
			return nil, AppForbiddenError
		}
		if !(primaryAccess && secondaryAccess) {
			slog.Error("user doesn't have required (EDIT) access to case", logAttributes)
			return nil, AppForbiddenError
		}
	}

	relatedCase, err := r.app.Store.RelatedCase().Create(createOpts, &req.GetInput().RelationType)
	if err != nil {
		return nil, err
	}

	parsedID, err := strconv.Atoi(relatedCase.Id)
	if err != nil {
		return nil, cerror.NewInternalError(
			"app.related_case.create_related_case.invalid_id",
			"Failed to parse relation ID",
		)
	}

	relatedID, err := strconv.Atoi(relatedCase.RelatedCase.GetId())
	if err != nil {
		return nil, cerror.NewInternalError(
			"app.related_case.create_related_case.invalid_id",
			"Failed to parse related ID",
		)
	}

	primaryID, err := strconv.Atoi(relatedCase.PrimaryCase.GetId())
	if err != nil {
		return nil, cerror.NewInternalError(
			"app.related_case.create_related_case.invalid_id",
			"Failed to parse related ID",
		)
	}

	relatedCase.Id, err = etag.EncodeEtag(etag.EtagRelatedCase, int64(parsedID), relatedCase.Ver)
	if err != nil {
		slog.Error(err.Error(), logAttributes)
		return nil, AppResponseNormalizingError
	}
	relatedCase.RelatedCase.Id, err = etag.EncodeEtag(etag.EtagRelatedCase, int64(relatedID), relatedCase.Ver)
	if err != nil {
		slog.Error(err.Error(), logAttributes)
		return nil, AppResponseNormalizingError
	}
	relatedCase.PrimaryCase.Id, err = etag.EncodeEtag(etag.EtagRelatedCase, int64(primaryID), relatedCase.Ver)
	if err != nil {
		slog.Error(err.Error(), logAttributes)
		return nil, AppResponseNormalizingError
	}
	return relatedCase, nil
}

func (r *RelatedCaseService) UpdateRelatedCase(ctx context.Context, req *cases.UpdateRelatedCaseRequest) (*cases.RelatedCase, error) {
	if req.GetId() == "" {
		return nil, cerror.NewBadRequestError("app.related_case.update_related_case.id_required", "ID required")
	}

	updateOpts, err := model.NewUpdateOptions(ctx, req, RelatedCaseMetadata)
	if err != nil {
		slog.Error(err.Error())
		return nil, AppInternalError
	}
	tag, err := etag.EtagOrId(etag.EtagRelatedCase, req.GetId())
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
	if updateOpts.GetAuthOpts().GetObjectScope(RelatedCaseMetadata.GetMainScopeName()).IsRbacUsed() {
		primaryAccess, err := r.app.Store.Case().CheckRbacAccess(updateOpts, updateOpts.GetAuthOpts(), authmodel.Edit, primaryId)
		if err != nil {
			slog.Error(err.Error(), logAttributes)
			return nil, AppForbiddenError
		}
		secondaryAccess, err := r.app.Store.Case().CheckRbacAccess(updateOpts, updateOpts.GetAuthOpts(), authmodel.Read, relatedId)
		if err != nil {
			slog.Error(err.Error(), logAttributes)
			return nil, AppForbiddenError
		}
		if !(primaryAccess && secondaryAccess) {
			slog.Error("user doesn't have required access to case", logAttributes)
			return nil, AppForbiddenError
		}
	}

	output, err := r.app.Store.RelatedCase().Update(updateOpts, input)
	if err != nil {
		return nil, err
	}

	id, err := strconv.Atoi(output.Id)
	if err != nil {
		// Return the error if ID conversion fails
		return nil, cerror.NewInternalError("failed encoding id, error", err.Error())
	}

	output.Id, err = etag.EncodeEtag(etag.EtagRelatedCase, int64(id), output.Ver)
	if err != nil {
		slog.Error(err.Error(), logAttributes)
		return nil, AppResponseNormalizingError
	}
	return output, nil
}

func (r *RelatedCaseService) DeleteRelatedCase(ctx context.Context, req *cases.DeleteRelatedCaseRequest) (*cases.RelatedCase, error) {
	if req.GetId() == "" {
		return nil, cerror.NewBadRequestError("app.related_case.delete_related_case.id_required", "ID required")
	}

	tag, err := etag.EtagOrId(etag.EtagRelatedCase, req.Id)
	if err != nil {
		return nil, cerror.NewBadRequestError("app.related_case.delete_related_case.invalid_etag", "Invalid etag")
	}
	caseEtag, err := etag.EtagOrId(etag.EtagCase, req.Id)
	if err != nil {
		return nil, cerror.NewBadRequestError("app.related_case.delete_related_case.invalid_etag", "Invalid etag")
	}

	deleteOpts, err := model.NewDeleteOptions(ctx, RelatedCaseMetadata)
	if err != nil {
		slog.Error(err.Error())
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
	if deleteOpts.GetAuthOpts().GetObjectScope(RelatedCaseMetadata.GetMainScopeName()).IsRbacUsed() {
		if deleteOpts.GetAuthOpts().GetObjectScope(RelatedCaseMetadata.GetMainScopeName()).IsRbacUsed() {
			access, err := r.app.Store.Case().CheckRbacAccess(deleteOpts, deleteOpts.GetAuthOpts(), authmodel.Edit, deleteOpts.ParentID)
			if err != nil {
				slog.Error(err.Error(), logAttributes)
				return nil, AppForbiddenError
			}
			if !access {
				slog.Error("user doesn't have required (READ) access to the case", logAttributes)
				return nil, AppForbiddenError
			}
		}
	}

	err = r.app.Store.RelatedCase().Delete(deleteOpts)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (r *RelatedCaseService) ListRelatedCases(ctx context.Context, req *cases.ListRelatedCasesRequest) (*cases.RelatedCaseList, error) {
	if req.GetPrimaryCaseId() == "" {
		return nil, cerror.NewBadRequestError("app.related_case.list_related_case.id_required", "ID required")
	}
	tag, err := etag.EtagOrId(etag.EtagCase, req.PrimaryCaseId)
	if err != nil {
		return nil, cerror.NewBadRequestError("app.related_case.list_related_cases.invalid_etag", "Invalid etag")
	}

	ids, err := util.ParseIds(req.Ids, etag.EtagRelatedCase)
	if err != nil {
		return nil, cerror.NewBadRequestError("app.related_case.list_related_cases.invalid_ids", "Invalid ids format")
	}
	searchOpts, err := model.NewSearchOptions(ctx, req, RelatedCaseMetadata)
	if err != nil {
		slog.Error(err.Error())
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

		// Normalize related case ID
		id, err := strconv.Atoi(relatedCase.Id)
		if err != nil {
			return fmt.Errorf("failed encoding related_case id: %w", err)
		}
		relatedCase.Id, err = etag.EncodeEtag(etag.EtagRelatedCase, int64(id), relatedCase.Ver)
		if err != nil {
			return err
		}

		// Normalize primary case ID
		if relatedCase.PrimaryCase != nil {
			primaryCaseID, err := strconv.Atoi(relatedCase.PrimaryCase.GetId())
			if err != nil {
				return fmt.Errorf("failed encoding primary_case id: %w", err)
			}
			relatedCase.PrimaryCase.Id, err = etag.EncodeEtag(etag.EtagCase, int64(primaryCaseID), relatedCase.PrimaryCase.GetVer())
			if err != nil {
				return err
			}
			// Set PrimaryCase Ver to nil
			relatedCase.PrimaryCase.Ver = 0
		}

		// Normalize related case ID inside related case
		if relatedCase.RelatedCase != nil {
			relatedCaseID, err := strconv.Atoi(relatedCase.RelatedCase.Id)
			if err != nil {
				return fmt.Errorf("failed encoding related_case id: %w", err)
			}
			relatedCase.RelatedCase.Id, err = etag.EncodeEtag(etag.EtagCase, int64(relatedCaseID), relatedCase.RelatedCase.GetVer())
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
