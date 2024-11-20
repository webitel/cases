package app

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	cases "github.com/webitel/cases/api/cases"
	"github.com/webitel/cases/model"
	"github.com/webitel/webitel-go-kit/etag"

	cerror "github.com/webitel/cases/internal/error"
)

/*

API layer

- etag formation
- graphQL types declaration and validation
- authorization interceptor
- rabbitMQ events
- case name forming
-------------------------------------------------
- OpenTelemetry interceptors and init
- storage interfaces
- additional auth layer with context attributes, client name and service registry


Database layer

proto filter parsing
storages (singleton)
calendar storage and calculation module
sql scripts structure
*/

var defaultFieldsCase = []string{
	"id",
	"ver",
	"created_by",
	"created_at",
	"updated_by",
	"updated_at",
	"assignee",
	"reporter",
	"name",
	"subject",
	"description",
	"source",
	"priority",
	"impacted",
	"author",
	"planned_reaction_at",
	"planned_resolve_at",
	"status",
	"close_reason_group",
	"group",
	"close_result",
	"close_reason",
	"rating",
	"rating_comment",
	"sla_conditions",
	"status_condition",
	"sla",
}

type CaseService struct {
	app *App
	cases.UnimplementedCasesServer
}

/*
SearchCases
Authorization
Obac
Rbac
Fields validation with graph
Search options construction with filters
Database layer with search options
Result construction by fields requested
*/
func (c *CaseService) SearchCases(ctx context.Context, req *cases.SearchCasesRequest) (*cases.CaseList, error) {
	// TODO implement me
	panic("implement me")
}

/*
LocateCase
Authorization
Obac
Rbac
Etag parsing
Fields validation with graph
Search options construction with filters
Database layer with search options
Result construction with etag
*/
func (c *CaseService) LocateCase(ctx context.Context, req *cases.LocateCaseRequest) (*cases.Case, error) {
	// TODO implement me
	panic("implement me")
}

/*
CreateCase
Authorization
Obac
Fields validation with graph
Database layer with create options
Calendar's logic
Result construction
Rabbit event publishing
*/
func (c *CaseService) CreateCase(ctx context.Context, req *cases.CreateCaseRequest) (*cases.Case, error) {
	// Authorize the user
	session, err := c.app.AuthorizeFromContext(ctx)
	if err != nil {
		return nil, cerror.NewUnauthorizedError("app.case.create_case.authorization_failed", err.Error())
	}

	// Validate required fields
	if req.Input.Name == "" {
		return nil, cerror.NewBadRequestError("app.case.create_case.name_required", "Case name is required")
	}

	if req.Input.Subject == "" {
		return nil, cerror.NewBadRequestError("app.case.create_case.subject_required", "Case subject is required")
	}

	if req.Input.Reporter == 0 {
		return nil, cerror.NewBadRequestError("app.case.create_case.reporter_required", "Reporter is required")
	}

	if req.Input.Impacted == 0 {
		return nil, cerror.NewBadRequestError("app.case.create_case.impacted_required", "Impacted contact is required")
	}

	// Validate additional optional fields if needed
	if req.Input.Priority == 0 {
		return nil, cerror.NewBadRequestError("app.case.create_case.invalid_priority", "Invalid priority specified")
	}

	if req.Input.Service == 0 {
		return nil, cerror.NewBadRequestError("app.case.create_case.invalid_service", "Invalid service specified")
	}

	// // Convert InputCaseComment to CaseCommentList
	// var comments *cases.CaseCommentList
	// if req.Input.Comments != nil && len(req.Input.Comments) > 0 {
	// 	comments = &cases.CaseCommentList{
	// 		Items: make([]*cases.CaseComment, len(req.Input.Comments)),
	// 	}
	// 	for i, inputComment := range req.Input.Comments {
	// 		tag, err := etag.EtagOrId(etag.EtagCase, inputComment.Etag)
	// 		if err != nil {
	// 			return nil, cerror.NewBadRequestError("app.case.create_case.invalid_etag", "Invalid etag")
	// 		}
	// 		comments.Items[i] = &cases.CaseComment{
	// 			Id: strconv.Itoa(int(tag.GetOid())),
	// 		}
	// 	}
	// }

	// // Convert InputCaseLink to CaseLinkList
	// var links *cases.CaseLinkList
	// if req.Input.Links != nil && len(req.Input.Links) > 0 {
	// 	linkItems := make([]*cases.CaseLink, len(req.Input.Links))
	// 	for i, inputLink := range req.Input.Links {
	// 		tag, err := etag.EtagOrId(etag.EtagCaseLink, inputLink.Etag)
	// 		if err != nil {
	// 			return nil, cerror.NewBadRequestError("app.case.create_case.invalid_etag", "Invalid etag for link")
	// 		}
	// 		linkItems[i] = &cases.CaseLink{
	// 			Id: tag.GetOid(),
	// 		}
	// 	}
	// 	links = &cases.CaseLinkList{Items: linkItems}
	// }

	// // Convert InputRelatedCase to RelatedCaseList
	// var related *cases.RelatedCaseList
	// if req.Input.Related != nil && len(req.Input.Related) > 0 {
	// 	relatedItems := make([]*cases.RelatedCase, len(req.Input.Related))
	// 	for i, inputRelated := range req.Input.Related {
	// 		tag, err := etag.EtagOrId(etag.EtagRelatedCase, inputRelated.Etag)
	// 		if err != nil {
	// 			return nil, cerror.NewBadRequestError("app.case.create_case.invalid_etag", "Invalid etag for related case")
	// 		}
	// 		relatedItems[i] = &cases.RelatedCase{
	// 			Id: tag.GetOid(),
	// 		}
	// 	}
	// 	related = &cases.RelatedCaseList{Items: relatedItems}
	// }

	var (
		commentItems []*cases.CaseComment
		linkItems    []*cases.CaseLink
		relatedItems []*cases.RelatedCase
	)

	if len(req.Input.Comments) > 0 {
		commentItems, err = transformWithEtag(req.Input.Comments, etag.EtagCase, "app.case.create_case.invalid_comment_etag", func(etag int64) *cases.CaseComment {
			return &cases.CaseComment{Id: strconv.Itoa(int(etag))}
		})
		if err != nil {
			return nil, err
		}
	}

	if len(req.Input.Comments) > 0 {
		linkItems, err = transformWithEtag(req.Input.Links, etag.EtagCaseLink, "app.case.create_case.invalid_link_etag", func(etag int64) *cases.CaseLink {
			return &cases.CaseLink{Id: etag}
		})
		if err != nil {
			return nil, err
		}
	}

	if len(req.Input.Comments) > 0 {
		relatedItems, err = transformWithEtag(req.Input.Related, etag.EtagRelatedCase, "app.case.create_case.invalid_related_case_etag", func(etag int64) *cases.RelatedCase {
			return &cases.RelatedCase{Id: etag}
		})
		if err != nil {
			return nil, err
		}
	}

	// -----------------------------------------------------------------------------
	// Special Fields Overview (Computed or Derived Dynamically)
	// -----------------------------------------------------------------------------
	// - planned_reaction_at:
	//     * Automatically calculated based on the SLA and calendar conditions.
	//     * Determines the expected reaction time for the case.
	// - planned_resolve_at:
	//     * Computed dynamically using SLA rules and calendar settings.
	//     * Represents the anticipated resolution time for the case.
	// - source:
	//     * Derived from the case's service.
	// - status_condition:
	//     * Set based on the provided status or inferred through system logic.
	// - timing:
	//     * Calculated dynamically by the SLA engine during case lifecycle.
	//     * Represents SLA-driven timing metrics for reaction and resolution.
	// - SLA:
	//     * Pulled from the associated service's configuration.
	//     * Defines the service-level agreements applicable to this case.
	// - SLA Conditions:
	//     * Derived from the associated service's SLA configuration.
	//     * Specifies conditions and thresholds for SLA adherence.
	// -----------------------------------------------------------------------------
	newCase := &cases.Case{
		Name:             req.Input.Name,
		Subject:          req.Input.Subject,
		Description:      req.Input.Description,
		ContactInfo:      req.Input.ContactInfo,
		Assignee:         &cases.Lookup{Id: req.Input.Assignee},
		Reporter:         &cases.Lookup{Id: req.Input.Reporter},
		Source:           &cases.Lookup{Id: req.Input.Source},
		Impacted:         &cases.Lookup{Id: req.Input.Impacted},
		Group:            &cases.Lookup{Id: req.Input.Group},
		Status:           &cases.Lookup{Id: req.Input.Status},
		CloseReasonGroup: &cases.Lookup{Id: req.Input.CloseReason},
		Close: &cases.CloseInfo{
			CloseResult: req.Input.Close.CloseResult,
			CloseReason: &cases.Lookup{Id: req.Input.Close.CloseReason},
		},
		Priority: &cases.Lookup{Id: req.Input.Priority},
		Service:  &cases.Lookup{Id: req.Input.Service},
		Comments: &cases.CaseCommentList{Items: commentItems},
		Links:    &cases.CaseLinkList{Items: linkItems},
		Related:  &cases.RelatedCaseList{Items: relatedItems},
		Rate: &cases.RateInfo{
			Rating:        req.Input.Rate.GetRating(),
			RatingComment: req.Input.Rate.GetRatingComment(),
		},
	}
	t := time.Now().UTC()

	// Define create options
	createOpts := model.CreateOptions{
		Session: session,
		Context: ctx,
		Time:    t,
		Fields:  defaultFieldsCase,
	}

	newCase, err = c.app.Store.Case().Create(&createOpts, newCase)
	if err != nil {
		return nil, cerror.NewInternalError("app.case_comment.publish_comment.publish_error", err.Error())
	}

	// Encode etag from the case ID and version
	etag := etag.EncodeEtag(etag.EtagCaseComment, int64(newCase.Id), newCase.Ver)
	newCase.Etag = etag

	// Publish an event to RabbitMQ
	event := map[string]interface{}{
		"action":    "CreateCase",
		"user":      session.GetUserId(),
		"case_id":   newCase.Id,
		"case_etag": etag,
		"case_ver":  newCase.Ver,
		"case_name": newCase.Name,
	}

	eventData, err := json.Marshal(event)
	if err != nil {
		return nil, cerror.NewInternalError("app.case.create_case.event_marshal.failed", err.Error())
	}

	err = c.app.rabbit.Publish(
		model.APP_SERVICE_NAME,
		"create_case_key",
		eventData,
		strconv.Itoa(int(session.GetUserId())),
		time.Now(),
	)
	if err != nil {
		return nil, cerror.NewInternalError("app.case.create_case.event_publish.failed", err.Error())
	}

	return newCase, nil
}

// Generic function to transform a list of inputs with etag validation
func transformWithEtag[T any, R any](
	inputList []*T, // List of input items to process
	etagType etag.EtagType, // Type of etag to validate
	errorIdentifier string, // Error identifier for invalid etag
	transform func(int64) R, // Transformation function to generate output item
) ([]R, error) {
	if len(inputList) == 0 {
		return nil, nil
	}

	// Initialize the output slice
	outputList := make([]R, len(inputList))

	for i, input := range inputList {
		// Ensure the input implements the GetEtag method
		etagSource, ok := any(input).(interface{ GetEtag() string })
		if !ok {
			return nil, fmt.Errorf("input does not implement GetEtag")
		}

		// Extract and validate the etag
		tag, err := etag.EtagOrId(etagType, etagSource.GetEtag())
		if err != nil {
			return nil, cerror.NewBadRequestError(errorIdentifier, "Invalid etag")
		}

		// Transform the validated etag to the desired output type
		outputList[i] = transform(tag.GetOid())
	}

	return outputList, nil
}

func (c *CaseService) UpdateCase(ctx context.Context, req *cases.UpdateCaseRequest) (*cases.Case, error) {
	// TODO implement me
	panic("implement me")
}

func (c *CaseService) DeleteCase(ctx context.Context, req *cases.DeleteCaseRequest) (*cases.Case, error) {
	// TODO implement me
	panic("implement me")
}

func NewCaseService(app *App) (*CaseService, cerror.AppError) {
	if app == nil {
		return nil, cerror.NewBadRequestError("app.case.new_case_service.check_args.app", "unable to init case service, app is nil")
	}
	return &CaseService{app: app}, nil
}
