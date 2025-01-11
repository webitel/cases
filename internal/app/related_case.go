package app

import (
	"context"
	"fmt"
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
	searchOpts := model.NewLocateOptions(ctx, req, RelatedCaseMetadata)
	searchOpts.IDs = []int64{tag.GetOid()}

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

	createOpts := model.NewCreateOptions(ctx, req, RelatedCaseMetadata)

	primaryCaseTag, err := etag.EtagOrId(etag.EtagCase, req.GetPrimaryCaseId())
	if err != nil {
		return nil, cerror.NewBadRequestError(
			"app.related_case.created_related_case.invalid_etag",
			"Invalid primary case etag",
		)
	}

	relatedCaseTag, err := etag.EtagOrId(etag.EtagCase, req.GetInput().GetRelatedCaseId())
	if err != nil {
		return nil, cerror.NewBadRequestError(
			"app.related_case.created_related_case.invalid_etag",
			"Invalid relatedCase etag",
		)
	}
	createOpts.ParentID = primaryCaseTag.GetOid()
	createOpts.ChildID = relatedCaseTag.GetOid()

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

	relatedCase.Id = etag.EncodeEtag(etag.EtagRelatedCase, int64(parsedID), relatedCase.Ver)
	relatedCase.RelatedCase.Id = etag.EncodeEtag(etag.EtagRelatedCase, int64(relatedID), relatedCase.Ver)
	relatedCase.PrimaryCase.Id = etag.EncodeEtag(etag.EtagRelatedCase, int64(primaryID), relatedCase.Ver)
	return relatedCase, nil
}

func (r *RelatedCaseService) UpdateRelatedCase(ctx context.Context, req *cases.UpdateRelatedCaseRequest) (*cases.RelatedCase, error) {
	if req.GetId() == "" {
		return nil, cerror.NewBadRequestError("app.related_case.update_related_case.id_required", "ID required")
	}

	updateOpts := model.NewUpdateOptions(ctx, req, RelatedCaseMetadata)
	tag, err := etag.EtagOrId(etag.EtagRelatedCase, req.GetId())
	if err != nil {
		return nil, cerror.NewBadRequestError(
			"app.related_case.created_related_case.invalid_etag",
			"Invalid ID",
		)
	}
	updateOpts.Etags = []*etag.Tid{&tag}

	primaryCaseTag, err := etag.EtagOrId(etag.EtagCase, req.GetInput().GetPrimaryCaseId())
	if err != nil {
		return nil, cerror.NewBadRequestError(
			"app.related_case.created_related_case.invalid_etag",
			"Invalid primary case etag",
		)
	}

	relatedCaseTag, err := etag.EtagOrId(etag.EtagCase, req.GetInput().GetRelatedCaseId())
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
		PrimaryCaseId: strconv.Itoa(int(primaryCaseTag.GetOid())),
		RelatedCaseId: strconv.Itoa(int(relatedCaseTag.GetOid())),
		RelationType:  req.Input.RelationType,
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

	output.Id = etag.EncodeEtag(etag.EtagRelatedCase, int64(id), output.Ver)
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

	deleteOpts := model.NewDeleteOptions(ctx)
	deleteOpts.ID = tag.GetOid()

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
	searchOpts := model.NewSearchOptions(ctx, req, RelatedCaseMetadata)
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
		relatedCase.Id = etag.EncodeEtag(etag.EtagRelatedCase, int64(id), relatedCase.Ver)

		// Normalize primary case ID
		if relatedCase.PrimaryCase != nil {
			primaryCaseID, err := strconv.Atoi(relatedCase.PrimaryCase.GetId())
			if err != nil {
				return fmt.Errorf("failed encoding primary_case id: %w", err)
			}
			relatedCase.PrimaryCase.Id = etag.EncodeEtag(etag.EtagCase, int64(primaryCaseID), relatedCase.PrimaryCase.GetVer())

			// Set PrimaryCase Ver to nil
			relatedCase.PrimaryCase.Ver = 0
		}

		// Normalize related case ID inside related case
		if relatedCase.RelatedCase != nil {
			relatedCaseID, err := strconv.Atoi(relatedCase.RelatedCase.Id)
			if err != nil {
				return fmt.Errorf("failed encoding related_case id: %w", err)
			}
			relatedCase.RelatedCase.Id = etag.EncodeEtag(etag.EtagCase, int64(relatedCaseID), relatedCase.RelatedCase.GetVer())

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
