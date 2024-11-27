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
	// - status_condition:
	//     * Set based on the provided status.
	// - timing:
	//     * Calculated dynamically by the SLA engine during case lifecycle.
	//     * Represents SLA-driven timing metrics for reaction and resolution.
	// - SLA:
	//     * Pulled from the associated service's configuration using a recursive query.
	//     * The process begins by traversing the service hierarchy, starting from the given service ID,
	//       and recursively finding the "deepest" child service (the lowest level in the hierarchy) that has a non-NULL SLA.
	//     * The SLA is selected from this lowest-level service.
	//     * SLA conditions are further checked for specific priorities.
	//         - If a condition matches the priority, the reaction and resolution times are derived from the SLA condition.
	//         - If no condition matches the priority, the default reaction and resolution times are taken directly from the SLA.
	// - SLA Conditions:
	//     * Derived from the associated service's SLA configuration.
	//     * SLA conditions define specific rules and thresholds for SLA adherence based on the priority of the case.
	//     * During SLA resolution:
	//         - If a condition matches the given priority, the corresponding SLA condition is selected and applied.
	// -----------------------------------------------------------------------------
	newCase := &cases.Case{
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
	if req.Input.Etag == "" {
		return nil, cerror.NewBadRequestError("app.case.update_case.etag_required", "Etag is required")
	}
	if req.Input.Subject == "" {
		return nil, cerror.NewBadRequestError("app.case.update_case.subject_required", "Subject is required")
	}
	if req.Input.Status.GetId() == 0 {
		return nil, cerror.NewBadRequestError("app.case.update_case.status_required", "Status is required")
	}
	if req.Input.CloseReason.GetId() == 0 {
		return nil, cerror.NewBadRequestError("app.case.update_case.close_reason_group_required", "Close Reason group is required")
	}
	if req.Input.Priority.GetId() == 0 {
		return nil, cerror.NewBadRequestError("app.case.update_case.priority_required", "Priority is required")
	}
	if req.Input.Source.GetId() == 0 {
		return nil, cerror.NewBadRequestError("app.case.update_case.source_required", "Source is required")
	}
	if req.Input.Service.GetId() == 0 {
		return nil, cerror.NewBadRequestError("app.case.update_case.service_required", "Service is required")
	}

	fields := util.FieldsFunc(req.Fields, util.InlineFields)
	if len(fields) == 0 {
		fields = defaultFieldsCase
	}

	tag, err := etag.EtagOrId(etag.EtagCaseComment, req.Input.Etag)
	if err != nil {
		return nil, cerror.NewBadRequestError("app.case_comment.update_comment.invalid_etag", "Invalid etag")
	}

	updateOpts := model.NewUpdateOptions(ctx, req)
	updateOpts.IDs = []int64{tag.GetOid()}
	updateOpts.Fields = fields

	upd := &cases.Case{
		Ver:              tag.GetVer(),
		Subject:          req.Input.Subject,
		Description:      req.Input.Description,
		Status:           &cases.Lookup{Id: req.Input.Status.GetId()},
		CloseReasonGroup: &cases.Lookup{Id: req.Input.CloseReason.GetId()},
		Assignee:         &cases.Lookup{Id: req.Input.Assignee.GetId()},
		Reporter:         &cases.Lookup{Id: req.Input.Reporter.GetId()},
		Impacted:         &cases.Lookup{Id: req.Input.Impacted.GetId()},
		Group:            &cases.Lookup{Id: req.Input.Group.GetId()},
		Priority:         &cases.Lookup{Id: req.Input.Priority.GetId()},
		Source:           &cases.Lookup{Id: req.Input.Source.GetId()},
		Close: &cases.CloseInfo{
			CloseResult: req.Input.Close.CloseResult,
			CloseReason: req.Input.GetCloseReason(),
		},
		Rate: &cases.RateInfo{
			Rating:        req.Input.Rate.Rating,
			RatingComment: req.Input.Rate.RatingComment,
		},
		Service: &cases.Lookup{Id: req.Input.Service.GetId()},
	}

	updatedCase, err := c.app.Store.Case().Update(updateOpts, upd)
	if err != nil {
		return nil, cerror.NewInternalError("app.case.update_case.store_update_failed", err.Error())
	}

	// Encode etag from the comment ID and version
	e := etag.EncodeEtag(etag.EtagCaseComment, int64(updatedCase.Id), updatedCase.Ver)
	updatedCase.Etag = e

	return updatedCase, nil
}

func (c *CaseService) DeleteCase(
	ctx context.Context,
	req *cases.DeleteCaseRequest,
) (*cases.Case, error) {
	if req.Etag == "" {
		return nil, cerror.NewBadRequestError("app.case.delete_case.etag_required", "Etag is required")
	}

	deleteOpts := model.NewDeleteOptions(ctx)

	tag, err := etag.EtagOrId(etag.EtagCaseComment, req.Etag)
	if err != nil {
		return nil, cerror.NewBadRequestError("app.case.delete_case.invalid_etag", "Invalid etag")
	}
	deleteOpts.IDs = []int64{tag.GetOid()}

	err = c.app.Store.Case().Delete(deleteOpts)
	if err != nil {
		return nil, cerror.NewInternalError("app.case.delete_case.store_delete_failed", err.Error())
	}
	return nil, nil
}

func NewCaseService(app *App) (*CaseService, cerror.AppError) {
	if app == nil {
		return nil, cerror.NewBadRequestError("app.case.new_case_service.check_args.app", "unable to init case service, app is nil")
	}
	return &CaseService{app: app}, nil
}
