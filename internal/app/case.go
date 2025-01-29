package app

import (
	"context"
	"encoding/json"
	"log/slog"
	"strconv"
	"time"

	"github.com/webitel/webitel-go-kit/errors"

	cases "github.com/webitel/cases/api/cases"
	"github.com/webitel/cases/model"
	"github.com/webitel/cases/util"
	"github.com/webitel/webitel-go-kit/etag"

	cerror "github.com/webitel/cases/internal/error"
)

const caseObjScope = "cases"

var CaseMetadata = model.NewObjectMetadata(caseObjScope, "", []*model.Field{
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
	{Name: "close", Default: true},
	{Name: "rate", Default: true},
	{Name: "sla_condition", Default: true},
	{Name: "service", Default: true},
	{Name: "status_condition", Default: true},
	{Name: "sla", Default: true},
	{Name: "comments", Default: false},
	{Name: "links", Default: false},
	{Name: "files", Default: false},
	{Name: "related", Default: false},
	{Name: "timing", Default: true},
	{Name: "contact_info", Default: true},
}, CaseCommentMetadata, CaseLinkMetadata, RelatedCaseMetadata)

type CaseService struct {
	app *App
	cases.UnimplementedCasesServer
}

func (c *CaseService) SearchCases(ctx context.Context, req *cases.SearchCasesRequest) (*cases.CaseList, error) {
	searchOpts, err := model.NewSearchOptions(ctx, req, CaseMetadata)
	if err != nil {
		slog.Error(err.Error())
		return nil, AppInternalError
	}
	logAttributes := slog.Group("context", slog.Int64("user_id", searchOpts.GetAuthOpts().GetUserId()), slog.Int64("domain_id", searchOpts.GetAuthOpts().GetDomainId()))
	ids, err := util.ParseIds(req.GetIds(), etag.EtagCase)
	if err != nil {
		slog.Error(err.Error(), logAttributes)
		return nil, cerror.NewBadRequestError("app.case.search_cases.parse_ids.invalid", err.Error())
	}
	for column, value := range req.GetFilters() {
		if column != "" {
			searchOpts.Filter[column] = value
		}
	}
	searchOpts.IDs = ids
	list, err := c.app.Store.Case().List(searchOpts)
	if err != nil {
		slog.Error(err.Error(), logAttributes)
		return nil, errors.NewInternalError("app.case_communication.search_cases.database.error", "database error")
	}
	err = c.NormalizeResponseCases(list, req, nil)
	if err != nil {
		slog.Error(err.Error(), logAttributes)
		return nil, AppResponseNormalizingError
	}
	return list, nil
}

func (c *CaseService) LocateCase(ctx context.Context, req *cases.LocateCaseRequest) (*cases.Case, error) {
	searchOpts, err := model.NewLocateOptions(ctx, req, CaseMetadata)
	if err != nil {
		slog.Error(err.Error())
		return nil, AppInternalError
	}
	logAttributes := slog.Group("context", slog.Int64("user_id", searchOpts.GetAuthOpts().GetUserId()), slog.Int64("domain_id", searchOpts.GetAuthOpts().GetDomainId()))
	id, err := util.ParseIds([]string{req.GetEtag()}, etag.EtagCase)
	if err != nil {
		slog.Error(err.Error(), logAttributes)
		return nil, cerror.NewBadRequestError("app.case_link.locate.parse_qin.invalid", err.Error())
	}
	searchOpts.IDs = id
	list, err := c.app.Store.Case().List(searchOpts)
	if err != nil {
		return nil, err
	}
	if len(list.Items) == 0 {
		return nil, cerror.NewBadRequestError("app.case_link.locate.not_found", "entity not found")
	}
	err = c.NormalizeResponseCases(list, req, nil)
	if err != nil {
		slog.Error(err.Error(), logAttributes)
		return nil, AppResponseNormalizingError
	}
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
				Etag:         inputRelated.GetRelatedTo(),
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
		Assignee:         req.Input.Assignee,
		Reporter:         req.Input.Reporter,
		Source:           &cases.SourceTypeLookup{Id: req.Input.Source.GetId()},
		Impacted:         req.Input.Impacted,
		Group:            &cases.ExtendedLookup{Id: req.Input.Group.GetId()},
		Status:           req.Input.Status,
		Close:            (*cases.CloseInfo)(req.Input.GetClose()),
		CloseReasonGroup: req.Input.GetCloseReasonGroup(),
		Priority:         &cases.Priority{Id: req.Input.Priority.GetId()},
		Service:          req.Input.Service,
		Links:            links,
		Related:          related,
	}

	createOpts, err := model.NewCreateOptions(ctx, req, CaseMetadata)
	if err != nil {
		slog.Error(err.Error())
		return nil, AppInternalError
	}
	logAttributes := slog.Group("context", slog.Int64("user_id", createOpts.GetAuthOpts().GetUserId()), slog.Int64("domain_id", createOpts.GetAuthOpts().GetDomainId()))
	newCase, err = c.app.Store.Case().Create(createOpts, newCase)
	if err != nil {
		slog.Error(err.Error(), logAttributes)
		return nil, AppDatabaseError
	}

	// Encode etag from the case ID and version
	newCase.Etag, err = etag.EncodeEtag(etag.EtagCase, newCase.Id, newCase.Ver)
	if err != nil {
		slog.Error(err.Error(), logAttributes)
		return nil, AppResponseNormalizingError
	}
	userId := createOpts.GetAuthOpts().GetUserId()

	// Publish an event to RabbitMQ
	event := map[string]interface{}{
		"action":    "CreateCase",
		"user":      userId,
		"case_id":   newCase.Id,
		"case_etag": newCase.Etag,
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

	if newCase.Reporter == nil && util.ContainsField(createOpts.Fields, "reporter") {
		newCase.Reporter = &cases.Lookup{
			Name: AnonymousName,
		}
	}
	return newCase, nil
}

func (c *CaseService) UpdateCase(ctx context.Context, req *cases.UpdateCaseRequest) (*cases.Case, error) {
	// Validate input
	appErr := c.ValidateUpdateInput(req.Input, req.XJsonMask)
	if appErr != nil {
		return nil, appErr
	}

	tag, err := etag.EtagOrId(etag.EtagCase, req.Input.Etag)
	if err != nil {
		slog.Error(err.Error())
		return nil, cerror.NewBadRequestError("app.case.update.invalid_etag", "Invalid etag")
	}

	updateOpts, err := model.NewUpdateOptions(ctx, req, CaseMetadata)
	if err != nil {
		slog.Error(err.Error())
		return nil, AppInternalError
	}
	updateOpts.Etags = []*etag.Tid{&tag}
	logAttributes := slog.Group("context", slog.Int64("user_id", updateOpts.GetAuthOpts().GetUserId()), slog.Int64("domain_id", updateOpts.GetAuthOpts().GetDomainId()), slog.Int64("case_id", tag.GetOid()))

	upd := &cases.Case{
		Id:               tag.GetOid(),
		Ver:              tag.GetVer(),
		Subject:          req.Input.GetSubject(),
		Description:      req.Input.GetDescription(),
		ContactInfo:      req.Input.GetContactInfo(),
		Status:           req.Input.GetStatus(),
		StatusCondition:  req.Input.GetStatusCondition(),
		CloseReasonGroup: req.Input.GetCloseReason(),
		Assignee:         req.Input.GetAssignee(),
		Reporter:         req.Input.GetReporter(),
		Impacted:         req.Input.GetImpacted(),
		Group:            &cases.ExtendedLookup{Id: req.Input.GetGroup().GetId()},
		Priority:         &cases.Priority{Id: req.Input.Priority.GetId()},
		Source:           &cases.SourceTypeLookup{Id: req.Input.Source.GetId()},
		Close: &cases.CloseInfo{
			CloseResult: req.Input.Close.GetCloseResult(),
			CloseReason: req.Input.Close.GetCloseReason(),
		},
		Rate: &cases.RateInfo{
			Rating:        req.Input.Rate.GetRating(),
			RatingComment: req.Input.Rate.GetRatingComment(),
		},
		Service: req.Input.GetService(),
	}

	updatedCase, err := c.app.Store.Case().Update(updateOpts, upd)
	if err != nil {
		switch err.(type) {
		case *cerror.DBNoRowsError:
			return nil, cerror.NewBadRequestError("app.case.update.invalid_etag", "Invalid etag")
		}
		slog.Error(err.Error())
		return nil, AppDatabaseError
	}

	err = c.NormalizeResponseCase(updatedCase, req)
	if err != nil {
		slog.Error(err.Error(), logAttributes)
		return nil, AppResponseNormalizingError
	}
	return updatedCase, nil
}

func (c *CaseService) DeleteCase(ctx context.Context, req *cases.DeleteCaseRequest) (*cases.Case, error) {
	if req.Etag == "" {
		return nil, cerror.NewBadRequestError("app.case.delete_case.etag_required", "Etag is required")
	}

	deleteOpts, err := model.NewDeleteOptions(ctx, CaseMetadata)
	if err != nil {
		slog.Error(err.Error())
		return nil, AppInternalError
	}

	tag, err := etag.EtagOrId(etag.EtagCase, req.Etag)
	if err != nil {
		return nil, cerror.NewBadRequestError("app.case.delete_case.invalid_etag", "Invalid etag")
	}
	deleteOpts.IDs = []int64{tag.GetOid()}
	logAttributes := slog.Group("context", slog.Int64("user_id", deleteOpts.GetAuthOpts().GetUserId()), slog.Int64("domain_id", deleteOpts.GetAuthOpts().GetDomainId()), slog.Int64("case_id", tag.GetOid()))

	err = c.app.Store.Case().Delete(deleteOpts)
	if err != nil {
		switch err.(type) {
		case *cerror.DBNoRowsError:
			return nil, cerror.NewBadRequestError("app.case.delete.invalid_etag", "Invalid etag")
		}
		slog.Error(err.Error(), logAttributes)
		return nil, AppDatabaseError
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
	if input.Etag == "" {
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
			if input.Close.GetCloseReason().GetId() == 0 {
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

	if input.Status.GetId() == 0 {
		return cerror.NewBadRequestError("app.case.create_case.status_required", "Case status is required")
	}

	if input.GetCloseReasonGroup().GetId() == 0 {
		return cerror.NewBadRequestError("app.case.create_case.close_reason_required", "Case close reason is required")
	}

	if input.Source.GetId() == 0 {
		return cerror.NewBadRequestError("app.case.create_case.source_required", "Case source is required")
	}

	if input.Impacted.GetId() == 0 {
		return cerror.NewBadRequestError("app.case.create_case.impacted_required", "Impacted contact is required")
	}

	// Validate additional optional fields if needed
	if input.Priority.GetId() == 0 {
		return cerror.NewBadRequestError("app.case.create_case.invalid_priority", "Invalid priority specified")
	}

	if input.Service.GetId() == 0 {
		return cerror.NewBadRequestError("app.case.create_case.invalid_service", "Invalid service specified")
	}
	return nil
}

// NormalizeResponseCases validates and normalizes the response cases.CaseList to the front-end side.
func (c *CaseService) NormalizeResponseCases(res *cases.CaseList, mainOpts model.Fielder, subOpts map[string]model.Fielder) error {
	var err error
	fields := mainOpts.GetFields()
	if len(fields) == 0 {
		fields = CaseMetadata.GetDefaultFields()
	}
	hasEtag, hasId, hasVer := util.FindEtagFields(fields)

	fields = util.FieldsFunc(fields, util.InlineFields)

	for _, item := range res.Items {
		err = util.NormalizeEtags(etag.EtagCase, hasEtag, hasId, hasVer, &item.Etag, &item.Id, &item.Ver)
		if err != nil {
			return err
		}
		if item.Reporter == nil && util.ContainsField(fields, "reporter") {
			item.Reporter = &cases.Lookup{
				Name: AnonymousName,
			}
		}

	}

	for _, item := range res.Items {
		if item.Comments != nil {
			for _, com := range item.Comments.Items {
				err = util.NormalizeEtags(etag.EtagCaseComment, true, false, false, &com.Etag, &com.Id, &com.Ver)
				if err != nil {
					return err
				}
			}
		}
		if item.Links != nil {
			for _, link := range item.Links.Items {
				err = util.NormalizeEtags(etag.EtagCaseLink, true, false, false, &link.Etag, &link.Id, &link.Ver)
				if err != nil {
					return err
				}
			}
		}
		if item.Related != nil {
			for _, related := range item.Related.Data {
				err = util.NormalizeEtags(etag.EtagRelatedCase, true, false, false, &related.Etag, &related.Id, &related.Ver)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// NormalizeResponseCase validates and normalizes the response cases.Case to the front-end side.
func (c *CaseService) NormalizeResponseCase(re *cases.Case, opts model.Fielder) error {
	fields := opts.GetFields()
	if len(fields) == 0 {
		fields = CaseMetadata.GetDefaultFields()
	}
	util.NormalizeEtag(fields, &re.Etag, &re.Id, &re.Ver)

	if re.Reporter == nil && util.ContainsField(fields, "reporter") {
		re.Reporter = &cases.Lookup{
			Name: AnonymousName,
		}
	}
	return nil
}

// endregion
