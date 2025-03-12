package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/webitel/cases/api/cases"
	webitelgo "github.com/webitel/cases/api/webitel-go/contacts"
	cerror "github.com/webitel/cases/internal/errors"
	"github.com/webitel/cases/model"
	grpcopts "github.com/webitel/cases/model/options/grpc"
	"github.com/webitel/cases/model/options/grpc/shared"
	"github.com/webitel/cases/util"
	wlogger "github.com/webitel/logger/pkg/client/v2"
	"github.com/webitel/webitel-go-kit/etag"
	"google.golang.org/grpc/metadata"
	"log/slog"
	"strconv"
	"strings"
)

const (
	dynamicGroup = "dynamic"
	caseObjScope = model.ScopeCases
)

var (
	dynamicGroupFields = []string{"id", "name", "type", "conditions", "default_group"}
	groupXJsonMask     = []string{"group", "assignee"}

	CaseMetadata = model.NewObjectMetadata(caseObjScope, "", []*model.Field{
		{Name: "etag", Default: true},
		{Name: "id", Default: true},
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
		{Name: "role_ids", Default: false},
	}, CaseCommentMetadata, CaseLinkMetadata, RelatedCaseMetadata)
)

type CaseService struct {
	app *App
	cases.UnimplementedCasesServer
	logger *wlogger.ObjectedLogger
}

func (c *CaseService) SearchCases(ctx context.Context, req *cases.SearchCasesRequest) (*cases.CaseList, error) {
	searchOpts, err := grpcopts.NewSearchOptions(
		ctx,
		grpcopts.WithSearch(req),
		grpcopts.WithPagination(req),
		grpcopts.WithFilters(req),
		grpcopts.WithFields(req, CaseMetadata,
			util.DeduplicateFields,
			util.ParseFieldsForEtag,
			util.EnsureIdField,
		),
		grpcopts.WithIDsAsEtags(etag.EtagCase, req.GetIds()...),
		grpcopts.WithSort(req),
	)
	if err != nil {
		return nil, NewBadRequestError(err)
	}
	logAttributes := slog.Group("context", slog.Int64("user_id", searchOpts.GetAuthOpts().GetUserId()), slog.Int64("domain_id", searchOpts.GetAuthOpts().GetDomainId()))
	if req.GetContactId() != "" {
		contactId, err := strconv.ParseInt(req.GetContactId(), 10, 64)
		if err != nil {
			slog.ErrorContext(ctx, err.Error(), logAttributes)
			contactId = 0
		}
		searchOpts.AddFilter("contact", contactId)
	}
	list, err := c.app.Store.Case().List(searchOpts)
	if err != nil {
		slog.ErrorContext(ctx, err.Error(), logAttributes)
		return nil, DatabaseError
	}
	err = c.NormalizeResponseCases(list, req)
	if err != nil {
		slog.ErrorContext(ctx, err.Error(), logAttributes)
		return nil, ResponseNormalizingError
	}

	return list, nil
}

func (c *CaseService) LocateCase(ctx context.Context, req *cases.LocateCaseRequest) (*cases.Case, error) {
	searchOpts, err := grpcopts.NewLocateOptions(
		ctx,
		grpcopts.WithFields(req, CaseMetadata),
		grpcopts.WithIDsAsEtags(etag.EtagCase, req.GetEtag()),
	)
	if err != nil {
		return nil, NewBadRequestError(err)
	}
	logAttributes := slog.Group("context", slog.Int64("user_id", searchOpts.GetAuthOpts().GetUserId()), slog.Int64("domain_id", searchOpts.GetAuthOpts().GetDomainId()))
	list, err := c.app.Store.Case().List(searchOpts)
	if err != nil {
		return nil, err
	}
	if len(list.Items) == 0 {
		return nil, cerror.NewBadRequestError("app.case.locate.not_found", "entity not found")
	}
	err = c.NormalizeResponseCases(list, req)
	if err != nil {
		slog.ErrorContext(ctx, err.Error(), logAttributes)
		return nil, ResponseNormalizingError
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
	res := &cases.Case{
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

	createOpts, err := grpcopts.NewCreateOptions(
		ctx,
		grpcopts.WithCreateFields(
			req,
			CaseMetadata.CopyWithAllFieldsSetToDefault(),
			func(fields []string) []string {
				for i, v := range fields {
					if v == "related" {
						fields = append(fields[:i], fields[i+1:]...)
					}
				}
				return fields
			}),
	)
	if err != nil {
		return nil, NewBadRequestError(err)
	}

	logAttributes := slog.Group(
		"context",
		slog.Int64(
			"user_id",
			createOpts.GetAuthOpts().GetUserId()),
		slog.Int64(
			"domain_id",
			createOpts.GetAuthOpts().GetDomainId(),
		))

	res, err = c.app.Store.Case().Create(createOpts, res)
	if err != nil {
		slog.ErrorContext(ctx, err.Error(), logAttributes)
		return nil, DatabaseError
	}
	etag, err := etag.EncodeEtag(etag.EtagCase, res.Id, res.Ver)
	if err != nil {
		return nil, err
	}
	res.Etag = etag

	// save before normalize
	roleIds := res.GetRoleIds()
	id := res.GetId()

	//* Handle dynamic group update if applicable
	res, err = c.handleDynamicGroupUpdate(ctx, res)
	if err != nil {
		return nil, err
	}

	// * CREATE require all fields set to true incase we need to calculate dynamic condition
	req.Fields = createOpts.Fields
	err = c.NormalizeResponseCase(res, req)
	if err != nil {
		slog.ErrorContext(ctx, err.Error(), logAttributes)
		return nil, ResponseNormalizingError
	}

	ftsErr := c.SendFtsCreateEvent(id, createOpts.GetAuthOpts().GetDomainId(), roleIds, res)
	if ftsErr != nil {
		slog.ErrorContext(ctx, ftsErr.Error(), logAttributes)
	}

	logMessage, err := wlogger.NewCreateMessage(
		createOpts.GetAuthOpts().GetUserId(),
		getClientIp(ctx),
		res.Id, res,
	)
	if err != nil {
		return nil, err
	}

	logErr := c.logger.SendContext(ctx, createOpts.GetAuthOpts().GetDomainId(), logMessage)
	if logErr != nil {
		slog.ErrorContext(ctx, logErr.Error(), logAttributes)
	}

	err = c.app.watcher.OnEvent(EventTypeCreate, NewWatcherData(res, createOpts.GetAuthOpts().GetDomainId()))
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("could not notify case creation: %s, ", err.Error()), logAttributes)
	}
	return res, nil
}

func (c *CaseService) UpdateCase(ctx context.Context, req *cases.UpdateCaseRequest) (*cases.Case, error) {
	// Validate input
	appErr := c.ValidateUpdateInput(req.Input, req.XJsonMask)
	if appErr != nil {
		return nil, appErr
	}

	tag, err := etag.EtagOrId(etag.EtagCase, req.Input.Etag)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return nil, cerror.NewBadRequestError("app.case.update.invalid_etag", "Invalid etag")
	}

	updateOpts, err := grpcopts.NewUpdateOptions(
		ctx,
		grpcopts.WithUpdateFields(
			req,
			CaseMetadata.CopyWithAllFieldsSetToDefault(),
			func(fields []string) []string {
				for i, v := range fields {
					if v == "related" {
						fields = append(fields[:i], fields[i+1:]...)
					}
				}
				return fields
			}),
	)
	if err != nil {
		return nil, NewBadRequestError(err)
	}
	updateOpts.Etags = []*etag.Tid{&tag}

	logAttributes := slog.Group(
		"context",
		slog.Int64(
			"user_id",
			updateOpts.GetAuthOpts().GetUserId(),
		),
		slog.Int64(
			"domain_id",
			updateOpts.GetAuthOpts().GetDomainId(),
		),
		slog.Int64(
			"case_id",
			tag.GetOid(),
		))

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
		Group:            &cases.ExtendedLookup{Id: req.Input.Group.GetId()},
		Priority:         &cases.Priority{Id: req.Input.Priority.GetId()},
		Source:           &cases.SourceTypeLookup{Id: req.Input.Source.GetId()},
		Close:            req.Input.Close,
		Rate:             req.Input.Rate,
		Service:          req.Input.GetService(),
	}

	res, err := c.app.Store.Case().Update(updateOpts, upd)
	if err != nil {
		switch err.(type) {
		case *cerror.DBNoRowsError:
			return nil, cerror.NewBadRequestError("app.case.update.invalid_etag", "Invalid etag")
		}
		slog.ErrorContext(ctx, err.Error())
		return nil, DatabaseError
	}

	etag, err := etag.EncodeEtag(etag.EtagCase, res.Id, res.Ver)
	if err != nil {
		return nil, err
	}
	res.Etag = etag

	// *Handle dynamic group update if applicable
	res, err = c.handleDynamicGroupUpdate(ctx, res)
	if err != nil {
		return nil, err
	}

	roleIds := res.GetRoleIds()
	id := res.GetId()

	// * Update  require all fields set to true incase we need to calculate dynamic condition
	req.Fields = updateOpts.Fields
	err = c.NormalizeResponseCase(res, req)
	if err != nil {
		slog.ErrorContext(ctx, err.Error(), logAttributes)
		return nil, ResponseNormalizingError
	}

	ftsErr := c.SendFtsUpdateEvent(id, updateOpts.GetAuthOpts().GetDomainId(), roleIds, res)
	if ftsErr != nil {
		slog.ErrorContext(ctx, ftsErr.Error(), logAttributes)
	}

	log, err := wlogger.NewCreateMessage(
		updateOpts.GetAuthOpts().GetUserId(),
		getClientIp(ctx),
		res.Id,
		res,
	)
	if err != nil {
		return nil, err
	}
	logErr := c.logger.SendContext(ctx, updateOpts.GetAuthOpts().GetDomainId(), log)
	if logErr != nil {
		slog.ErrorContext(ctx, logErr.Error(), logAttributes)
	}

	err = c.app.watcher.OnEvent(EventTypeUpdate, NewWatcherData(res, updateOpts.GetAuthOpts().GetDomainId()))
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("could not notify case update: %s, ", err.Error()), logAttributes)
	}

	return res, nil
}

// handleDynamicGroupUpdate checks if a dynamic group is needed and updates the case accordingly.
func (c *CaseService) handleDynamicGroupUpdate(
	ctx context.Context,
	input *cases.Case,
) (*cases.Case, error) {
	// *Check if the group is dynamic
	if input.Group != nil && input.Group.Type == dynamicGroup {

		var info metadata.MD
		var ok bool

		info, ok = metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, cerror.NewForbiddenError("internal.grpc.get_context", "Not found")
		}
		newCtx := metadata.NewOutgoingContext(ctx, info)

		id := strconv.Itoa(int(input.Group.GetId()))
		res, err := c.app.webitelgoClient.LocateGroup(
			newCtx,
			&webitelgo.LocateGroupRequest{
				Id:     id,
				Fields: dynamicGroupFields,
			})
		if err != nil {
			return nil, err
		}

		// Resolve dynamic group and update the case
		input, err = c.resolveDynamicContactGroup(ctx, input, res)
		if err != nil {
			return nil, err
		}
	}

	// Return the potentially updated case
	return input, nil
}

func (c *CaseService) resolveDynamicContactGroup(
	ctx context.Context,
	inputCase *cases.Case,
	inputGroup *webitelgo.LocateGroupResponse,
) (*cases.Case, error) {
	// Convert the case object to a map for dynamic evaluation
	caseMap, err := caseToMap(inputCase)
	if err != nil {
		fmt.Printf("Error converting case to map: %v\n", err)
		return nil, err
	}

	etag, ok := caseMap["case.etag"].(string)
	if !ok {
		return nil, ForbiddenError
	}

	// Iterate over all dynamic conditions in the group…
	for _, condition := range inputGroup.Group.Conditions {
		// Evaluate the condition against the case map
		if evaluateComplexCondition(caseMap, condition.Expression) {

			groupID, err := strconv.Atoi(condition.Group.GetId())
			if err != nil {
				return nil, ResponseNormalizingError
			}

			var assigneeID int
			if condition.Assignee != nil {
				assigneeID, err = strconv.Atoi(condition.Assignee.GetId())
				if err != nil {
					return nil, ResponseNormalizingError
				}
			} else {
				assigneeID = 0
			}

			// Build request for Case Update API
			req := &cases.UpdateCaseRequest{
				XJsonMask: groupXJsonMask,
				Input: &cases.InputCase{
					Group:    &cases.Lookup{Id: int64(groupID)},
					Assignee: &cases.Lookup{Id: int64(assigneeID)},
					Etag:     etag,
				},
			}

			updCase, err := c.UpdateCase(ctx, req)
			if err != nil {
				return nil, err
			}

			return updCase, nil
		}
	}

	groupID, err := strconv.Atoi(inputGroup.Group.DefaultGroup.GetId())
	if err != nil {
		return nil, fmt.Errorf("app.case.resolveDynamicContactGroup: failed to convert group ID to integer: %w", err)
	}

	req := &cases.UpdateCaseRequest{
		XJsonMask: groupXJsonMask,
		Input: &cases.InputCase{
			Group: &cases.Lookup{Id: int64(groupID)},
			Etag:  etag,
		},
	}
	// Final update if a default group was assigned
	updCase, err := c.UpdateCase(ctx, req)
	if err != nil {
		return nil, err
	}

	return updCase, nil
}

// Converts a Case object to map[string]interface{} with "case." prefixed keys and lowercase values (except case.etag)
func caseToMap(caseObj interface{}) (map[string]interface{}, error) {
	caseJSON, err := json.Marshal(caseObj)
	if err != nil {
		return nil, err
	}

	var caseMap map[string]interface{}
	err = json.Unmarshal(caseJSON, &caseMap)
	if err != nil {
		return nil, err
	}

	// Create a new map with "case." prefixed keys, handling nested fields
	prefixedMap := make(map[string]interface{})
	addPrefixedKeys(prefixedMap, caseMap, "case")

	return prefixedMap, nil
}

// Recursively adds prefixed keys for nested maps (all keys and values converted to lowercase except case.etag)
func addPrefixedKeys(dest map[string]any, source map[string]any, prefix string) {
	for key, value := range source {
		fullKey := strings.ToLower(prefix + "." + key)

		switch v := value.(type) {
		case map[string]any:
			// Recursively process nested maps
			addPrefixedKeys(dest, v, fullKey)
		case string:
			// Keep "case.etag" value unchanged, lowercase everything else
			if fullKey == "case.etag" {
				dest[fullKey] = v
			} else {
				dest[fullKey] = strings.ToLower(v)
			}
		default:
			dest[fullKey] = value // Keep non-string values unchanged
		}
	}
}

// Evaluates complex condition strings with support for AND (&&) and OR (||) operators using bitwise operations
func evaluateComplexCondition(caseMap map[string]any, condition string) bool {
	// Convert condition to lowercase to ensure case-insensitive matching
	condition = strings.ToLower(condition)

	// Split conditions by OR (||)
	orConditions := strings.Split(condition, "||")
	for _, orCondition := range orConditions {
		// Split each OR condition into AND (&&) conditions
		andConditions := strings.Split(orCondition, "&&")
		// Use a bitmask to track whether all AND conditions are met
		var andMask uint = 0
		for i, andCondition := range andConditions {
			andCondition = strings.TrimSpace(andCondition)
			if evaluateSingleCondition(caseMap, andCondition) {
				// Set the corresponding bit in the mask
				andMask |= 1 << i
			}
		}
		// Check if all AND conditions are met (all bits set in the mask)
		if andMask == (1<<len(andConditions) - 1) {
			return true
		}
	}
	return false
}

// Evaluates a single condition string, e.g., "case.assignee.name == 'John Wick'"
func evaluateSingleCondition(caseMap map[string]any, condition string) bool {
	// Convert condition to lowercase
	condition = strings.ToLower(condition)

	// Parse condition into field, operator, and value
	var field, operator, value string
	if strings.Contains(condition, "==") {
		parts := strings.Split(condition, "==")
		field, operator, value = strings.TrimSpace(parts[0]), "==", strings.TrimSpace(parts[1])
	} else if strings.Contains(condition, "!=") {
		parts := strings.Split(condition, "!=")
		field, operator, value = strings.TrimSpace(parts[0]), "!=", strings.TrimSpace(parts[1])
	} else {
		fmt.Printf("Unsupported operator in condition: %s\n", condition)
		return false
	}

	// Remove quotes around the value if present and convert to lowercase
	value = strings.Trim(value, `"'`)

	// Resolve the field path in the caseMap (keys and values are already lowercase)
	fieldValue := resolveFieldPath(caseMap, field)

	// Compare the field value with the condition value
	switch operator {
	case "==":
		return fmt.Sprintf("%v", fieldValue) == value
	case "!=":
		return fmt.Sprintf("%v", fieldValue) != value
	default:
		fmt.Printf("Unsupported operator: %s\n", operator)
		return false
	}
}

// Resolves a field path (e.g., "case.assignee.name") to its value in the map
func resolveFieldPath(data map[string]interface{}, path string) interface{} {
	// Convert path to lowercase to match map keys
	path = strings.ToLower(path)

	// Direct lookup
	if value, exists := data[path]; exists {
		return value
	}

	return nil
}

func (c *CaseService) DeleteCase(ctx context.Context, req *cases.DeleteCaseRequest) (*cases.Case, error) {
	if req.Etag == "" {
		return nil, cerror.NewBadRequestError("app.case.delete_case.etag_required", "Etag is required")
	}
	tag, err := etag.EtagOrId(etag.EtagCase, req.Etag)
	if err != nil {
		return nil, cerror.NewBadRequestError("app.case.delete.invalid_etag", "Invalid case etag")
	}
	deleteOpts, err := grpcopts.NewDeleteOptions(ctx, grpcopts.WithDeleteID(tag.GetOid()))
	if err != nil {
		return nil, NewBadRequestError(err)
	}
	logAttributes := slog.Group("context", slog.Int64("user_id", deleteOpts.GetAuthOpts().GetUserId()), slog.Int64("domain_id", deleteOpts.GetAuthOpts().GetDomainId()), slog.Int64("case_id", tag.GetOid()))

	err = c.app.Store.Case().Delete(deleteOpts)
	if err != nil {
		switch err.(type) {
		case *cerror.DBNoRowsError:
			return nil, cerror.NewBadRequestError("app.case.delete.invalid_etag", "Invalid etag or insufficient rights")
		}
		slog.ErrorContext(ctx, err.Error(), logAttributes)
		return nil, DatabaseError
	}
	log, err := wlogger.NewDeleteMessage(deleteOpts.GetAuthOpts().GetUserId(), getClientIp(ctx), tag.GetOid())
	if err != nil {
		return nil, err
	}
	logErr := c.logger.SendContext(ctx, deleteOpts.GetAuthOpts().GetDomainId(), log)
	if logErr != nil {
		slog.ErrorContext(ctx, logErr.Error(), logAttributes)
	}
	deleteCase := &cases.Case{
		Id:   tag.GetOid(),
		Ver:  tag.GetVer(),
		Etag: req.Etag,
	}
	ftsErr := c.SendFtsDeleteEvent(tag.GetOid(), deleteOpts.GetAuthOpts().GetDomainId())
	if ftsErr != nil {
		slog.ErrorContext(ctx, ftsErr.Error(), logAttributes)
	}
	err = c.app.watcher.OnEvent(EventTypeDelete, NewWatcherData(deleteCase, deleteOpts.GetAuthOpts().GetDomainId()))
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("could not notify case deletion: %s, ", err.Error()), logAttributes)
	}
	return nil, nil
}

func NewCaseService(app *App) (*CaseService, cerror.AppError) {
	if app == nil {
		return nil, cerror.NewBadRequestError("app.case.new_case_service.check_args.app", "unable to init case service, app is nil")
	}
	return &CaseService{app: app, logger: app.wtelLogger.GetObjectedLogger(CaseMetadata.GetMainScopeName())}, nil
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
func (c *CaseService) NormalizeResponseCases(res *cases.CaseList, mainOpts shared.Fielder) error {
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
		// always hide
		item.RoleIds = nil
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
func (c *CaseService) NormalizeResponseCase(re *cases.Case, opts shared.Fielder) error {
	fields := opts.GetFields()
	if len(fields) == 0 {
		fields = CaseMetadata.GetDefaultFields()
	}
	err := util.NormalizeEtag(etag.EtagCase, fields, &re.Etag, &re.Id, &re.Ver)
	if err != nil {
		return err
	}
	re.RoleIds = nil

	if re.Reporter == nil && util.ContainsField(fields, "reporter") {
		re.Reporter = &cases.Lookup{
			Name: AnonymousName,
		}
	}
	if re.Comments != nil {
		for _, com := range re.Comments.Items {
			err = util.NormalizeEtags(etag.EtagCaseComment, true, false, false, &com.Etag, &com.Id, &com.Ver)
			if err != nil {
				return err
			}
		}
	}
	if re.Links != nil {
		for _, link := range re.Links.Items {
			err = util.NormalizeEtags(etag.EtagCaseLink, true, false, false, &link.Etag, &link.Id, &link.Ver)
			if err != nil {
				return err
			}
		}
	}
	if re.Related != nil {
		for _, related := range re.Related.Data {
			err = util.NormalizeEtags(etag.EtagRelatedCase, true, false, false, &related.Etag, &related.Id, &related.Ver)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *CaseService) formFtsModel(roleIds []int64, item *cases.Case) (*model.FtsCase, error) {
	m := &model.FtsCase{
		Description: item.GetDescription(),
		RoleIds:     roleIds,
		Subject:     item.GetSubject(),
		ContactInfo: item.GetContactInfo(),
		CreatedAt:   item.GetCreatedAt(),
	}

	if rate := item.GetRate(); rate != nil {
		m.RatingComment = rate.GetRatingComment()
	}

	if cl := item.GetClose(); cl != nil {
		m.CloseResult = cl.GetCloseResult()
	}
	return m, nil
}

func (c *CaseService) SendFtsCreateEvent(id int64, domainId int64, roleIds []int64, item *cases.Case) error {
	if domainId == 0 {
		return errors.New("domain id required")
	}
	if id == 0 {
		return errors.New("id required")
	}
	m, err := c.formFtsModel(roleIds, item)
	if err != nil {
		return err
	}
	m.RoleIds = roleIds
	return c.app.ftsClient.Create(domainId, model.ScopeCases, id, m)
}

func (c *CaseService) SendFtsUpdateEvent(id int64, domainId int64, roleIds []int64, item *cases.Case) error {
	if domainId == 0 {
		return errors.New("domain id required")
	}
	if id == 0 {
		return errors.New("id required")
	}
	m, err := c.formFtsModel(roleIds, item)
	if err != nil {
		return err
	}
	return c.app.ftsClient.Update(domainId, model.ScopeCases, id, m)
}

func (c *CaseService) SendFtsDeleteEvent(id int64, domainId int64) error {
	if domainId == 0 {
		return errors.New("domain id required")
	}
	if id == 0 {
		return errors.New("id required")
	}
	return c.app.ftsClient.Delete(domainId, model.ScopeCases, id)
}
