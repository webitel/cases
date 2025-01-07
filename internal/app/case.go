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

var CaseMetadata = model.NewObjectMetadata(
	[]*model.Field{
		{Name: "etag", Default: true},
		{Name: "id", Default: false},
		{Name: "ver", Default: false},
		{Name: "created_by", Default: true},
		{Name: "created_at", Default: true},
		{Name: "updated_by", Default: false},
		{Name: "updated_at", Default: false},
		{Name: "assignee", Default: true},
		{Name: "reporter", Default: true},
		{Name: "name", Default: true},
		{Name: "subject", Default: true},
		{Name: "description", Default: true},
		{Name: "source", Default: true},
		{Name: "priority", Default: true},
		{Name: "impacted", Default: true},
		{Name: "author", Default: true},
		{Name: "planned_reaction_at", Default: true},
		{Name: "planned_resolve_at", Default: true},
		{Name: "status", Default: true},
		{Name: "close_reason_group", Default: true},
		{Name: "group", Default: true},
		{Name: "close_result", Default: false},
		{Name: "close_reason", Default: false},
		{Name: "rating", Default: false},
		{Name: "rating_comment", Default: false},
		{Name: "sla_condition", Default: true},
		{Name: "service", Default: true},
		{Name: "status_condition", Default: true},
		{Name: "sla", Default: true},
		{Name: "comments", Default: false},
		{Name: "links", Default: false},
		{Name: "files", Default: false},
		{Name: "related", Default: false},
		{Name: "contact_info", Default: false},
	})

type CaseService struct {
	app *App
	cases.UnimplementedCasesServer
}

func (c *CaseService) SearchCases(ctx context.Context, req *cases.SearchCasesRequest) (*cases.CaseList, error) {
	searchOpts := model.NewSearchOptions(ctx, req, CaseMetadata)
	ids, err := util.ParseIds(req.GetIds(), etag.EtagCase)
	for column, value := range req.GetFilters() {
		if column != "" {
			searchOpts.Filter[column] = value
		}
	}
	if err != nil {
		return nil, cerror.NewBadRequestError("app.case_link.locate.parse_qin.invalid", err.Error())
	}
	searchOpts.IDs = ids
	list, err := c.app.Store.Case().List(searchOpts)
	if err != nil {
		return nil, err
	}
	c.NormalizeResponseCases(list, req, nil)
	return list, nil
}

func (c *CaseService) LocateCase(ctx context.Context, req *cases.LocateCaseRequest) (*cases.Case, error) {
	searchOpts := model.NewLocateOptions(ctx, req, CaseMetadata)
	id, err := util.ParseIds([]string{req.GetId()}, etag.EtagCase)
	if err != nil {
		return nil, cerror.NewBadRequestError("app.case_link.locate.parse_qin.invalid", err.Error())
	}
	searchOpts.IDs = id
	list, err := c.app.Store.Case().List(searchOpts)
	if err != nil {
		return nil, err
	}
	c.NormalizeResponseCases(list, req, nil)
	return list.Items[0], nil
}

func (c *CaseService) CreateCase(ctx context.Context, req *cases.CreateCaseRequest) (*cases.Case, error) {
	// Validate required fields
	appErr := c.ValidateCreateInput(req.Input)
	if appErr != nil {
		return nil, appErr
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
		related = &cases.RelatedCaseList{Data: relatedItems}
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
		Source:           &cases.SourceTypeLookup{Id: req.Input.Source},
		Impacted:         &cases.Lookup{Id: req.Input.Impacted},
		Group:            &cases.Lookup{Id: req.Input.Group},
		Status:           &cases.Lookup{Id: req.Input.Status},
		CloseReasonGroup: &cases.Lookup{Id: req.Input.CloseReason},
		Priority:         &cases.Lookup{Id: req.Input.Priority},
		Service:          &cases.Lookup{Id: req.Input.Service},
		Links:            links,
		Related:          related,
	}

	createOpts := model.NewCreateOptions(ctx, req, CaseMetadata)

	newCase, err := c.app.Store.Case().Create(createOpts, newCase)
	if err != nil {
		return nil, err
	}
	id, _ := strconv.Atoi(newCase.Id)
	// Encode etag from the case ID and version
	newCase.Id = etag.EncodeEtag(etag.EtagCaseComment, int64(id), newCase.Ver)
	userId := createOpts.Session.GetUserId()

	// Publish an event to RabbitMQ
	event := map[string]interface{}{
		"action":    "CreateCase",
		"user":      userId,
		"case_id":   newCase.Id,
		"case_etag": newCase.Id,
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
		strconv.Itoa(int(userId)),
		time.Now(),
	)
	if err != nil {
		return nil, cerror.NewInternalError("app.case.create_case.event_publish.failed", err.Error())
	}
	c.NormalizeResponseCase(newCase, req)
	return newCase, nil
}

func (c *CaseService) UpdateCase(ctx context.Context, req *cases.UpdateCaseRequest) (*cases.Case, error) {
	// Validate input
	appErr := c.ValidateUpdateInput(req.Input, req.XJsonMask)
	if appErr != nil {
		return nil, appErr
	}

	tag, err := etag.EtagOrId(etag.EtagCase, req.Input.Id)
	if err != nil {
		return nil, cerror.NewBadRequestError("app.case_comment.update_comment.invalid_etag", "Invalid etag")
	}

	updateOpts := model.NewUpdateOptions(ctx, req, CaseMetadata)

	upd := &cases.Case{
		Id:               strconv.Itoa(int(tag.GetOid())),
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
		Source:           &cases.SourceTypeLookup{Id: req.Input.Source.GetId()},
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

	c.NormalizeResponseCase(updatedCase, req)
	return updatedCase, nil
}

func (c *CaseService) DeleteCase(ctx context.Context, req *cases.DeleteCaseRequest) (*cases.Case, error) {
	if req.Id == "" {
		return nil, cerror.NewBadRequestError("app.case.delete_case.etag_required", "Etag is required")
	}

	deleteOpts := model.NewDeleteOptions(ctx)

	tag, err := etag.EtagOrId(etag.EtagCaseComment, req.Id)
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

func (c *CaseService) ValidateUpdateInput(
	input *cases.InputCase,
	xJsonMask []string,
) cerror.AppError {
	if input.Id == "" {
		return cerror.NewBadRequestError("app.case.update_case.etag_required", "Etag is required")
	}

	// Ensure nested structures are initialized
	if input.Rate == nil {
		input.Rate = &cases.RateInfo{}
	}
	if input.Close == nil {
		input.Close = &cases.CloseInfo{}
	}

	// Iterate over xJsonMask and validate corresponding fields
	// Validating fields passed for updating
	for _, field := range xJsonMask {
		switch field {
		case "subject":
			if input.Subject == "" {
				return cerror.NewBadRequestError("app.case.update_case.subject_required", "Subject is required")
			}
		case "status":
			if input.Status.GetId() == 0 {
				return cerror.NewBadRequestError("app.case.update_case.status_required", "Status is required")
			}
		case "close.close_reason":
			if input.CloseReason.GetId() == 0 {
				return cerror.NewBadRequestError("app.case.update_case.close_reason_group_required", "Close Reason group is required")
			}
		case "priority":
			if input.Priority.GetId() == 0 {
				return cerror.NewBadRequestError("app.case.update_case.priority_required", "Priority is required")
			}
		case "source":
			if input.Source.GetId() == 0 {
				return cerror.NewBadRequestError("app.case.update_case.source_required", "Source is required")
			}
		case "service":
			if input.Service.GetId() == 0 {
				return cerror.NewBadRequestError("app.case.update_case.service_required", "Service is required")
			}
		}
	}

	return nil
}

// region UTILITY

func (c *CaseService) ValidateCreateInput(input *cases.InputCreateCase) cerror.AppError {
	if input.Subject == "" {
		return cerror.NewBadRequestError("app.case.create_case.subject_required", "Case subject is required")
	}

	if input.Status == 0 {
		return cerror.NewBadRequestError("app.case.create_case.status_required", "Case status is required")
	}

	if input.CloseReason == 0 {
		return cerror.NewBadRequestError("app.case.create_case.close_reason_required", "Case close reason is required")
	}

	if input.Source == 0 {
		return cerror.NewBadRequestError("app.case.create_case.source_required", "Case source is required")
	}

	if input.Impacted == 0 {
		return cerror.NewBadRequestError("app.case.create_case.impacted_required", "Impacted contact is required")
	}

	// Validate additional optional fields if needed
	if input.Priority == 0 {
		return cerror.NewBadRequestError("app.case.create_case.invalid_priority", "Invalid priority specified")
	}

	if input.Service == 0 {
		return cerror.NewBadRequestError("app.case.create_case.invalid_service", "Invalid service specified")
	}
	return nil
}

// NormalizeResponseCases validates and normalizes the response cases.CaseList to the front-end side.
func (c *CaseService) NormalizeResponseCases(res *cases.CaseList, mainOpts model.Fielder, subOpts map[string]model.Fielder) {
	fields := mainOpts.GetFields()
	if len(fields) == 0 {
		fields = CaseMetadata.GetDefaultFields()
	}

	for _, item := range res.Items {
		id, _ := strconv.Atoi(item.Id)
		item.Id = etag.EncodeEtag(etag.EtagCase, int64(id), item.Ver)
		if item.Reporter == nil && util.ContainsField(fields, "reporter") {
			item.Reporter = &cases.Lookup{
				Name: AnonymousName,
			}
		}

	}
	for _, field := range fields {
		switch field {
		case "comments":
			//for _, item := range res.Items {
			//	for _, comment := range item.Comments.Items {
			//		util.NormalizeEtags(etag.EtagCaseComment, true, true, true, &comment.Id, &comment.Id, &comment.Ver)
			//	}
			//}
		case "links":
			for _, item := range res.Items {
				if item.Links != nil {
					for _, link := range item.Links.Items {
						id, _ := strconv.Atoi(link.Id)
						item.Id = etag.EncodeEtag(etag.EtagCaseLink, int64(id), link.Ver)
					}
				}
			}
		case "related_cases":
			for _, item := range res.Items {
				if item.Related != nil {
					for _, related := range item.Related.Data {
						id, _ := strconv.Atoi(related.Id)
						related.Id = etag.EncodeEtag(etag.EtagRelatedCase, int64(id), related.Ver)
					}
				}
			}
		}
	}
}

// NormalizeResponseCase validates and normalizes the response cases.Case to the front-end side.
func (c *CaseService) NormalizeResponseCase(re *cases.Case, opts model.Fielder) {
	fields := opts.GetFields()
	if len(fields) == 0 {
		fields = CaseMetadata.GetDefaultFields()
	}

	id, _ := strconv.Atoi(re.Id)
	re.Id = etag.EncodeEtag(etag.EtagCase, int64(id), re.Ver)

	if re.Reporter == nil && util.ContainsField(fields, "reporter") {
		re.Reporter = &cases.Lookup{
			Name: AnonymousName,
		}
	}
}

// endregion
