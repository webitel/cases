package app

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	cases "github.com/webitel/cases/api/cases"
	"github.com/webitel/cases/model"
	"github.com/webitel/cases/util"
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
	"service",
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
	"related_cases",
	"links",
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

	if req.Input.Subject == "" {
		return nil, cerror.NewBadRequestError("app.case.create_case.subject_required", "Case subject is required")
	}

	if req.Input.Status == 0 {
		return nil, cerror.NewBadRequestError("app.case.create_case.status_required", "Case status is required")
	}

	if req.Input.CloseReason == 0 {
		return nil, cerror.NewBadRequestError("app.case.create_case.close_reason_required", "Case close reason is required")
	}

	if req.Input.Source == 0 {
		return nil, cerror.NewBadRequestError("app.case.create_case.source_required", "Case source is required")
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

	var (
		related *cases.RelatedCaseList
		links   *cases.CaseLinkList
	)

	if len(req.Input.Links) > 0 {
		linkItems := make([]*cases.CaseLink, len(req.Input.Links))
		for i, inputLink := range req.Input.Links {
			linkItems[i] = &cases.CaseLink{
				Url:  inputLink.GetUrl(),
				Name: inputLink.GetName(),
			}
		}
		links = &cases.CaseLinkList{Items: linkItems}
	}

	if len(req.Input.Related) > 0 {
		relatedItems := make([]*cases.RelatedCase, len(req.Input.Related))
		for i, inputRelated := range req.Input.Related {
			relatedItems[i] = &cases.RelatedCase{
				Id:           inputRelated.GetRelatedTo(),
				RelationType: inputRelated.RelationType,
			}
		}
		related = &cases.RelatedCaseList{Items: relatedItems}
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
		Links:    links,
		Related:  related,
		Rate: &cases.RateInfo{
			Rating:        req.Input.Rate.GetRating(),
			RatingComment: req.Input.Rate.GetRatingComment(),
		},
	}

	fields := util.FieldsFunc(req.Fields, util.InlineFields)
	if len(fields) == 0 {
		fields = defaultFieldsCase
	}
	t := time.Now().UTC()

	// Define create options
	createOpts := model.CreateOptions{
		Session: session,
		Context: ctx,
		Time:    t,
		Fields:  fields,
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
