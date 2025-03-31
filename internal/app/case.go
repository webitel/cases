package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/webitel/cases/auth"

	wlogger "github.com/webitel/logger/pkg/client/v2"

	"github.com/webitel/cases/api/cases"
	webitelgo "github.com/webitel/cases/api/webitel-go/contacts"
	cerror "github.com/webitel/cases/internal/errors"
	deferr "github.com/webitel/cases/internal/errors/defaults"
	"github.com/webitel/cases/model"
	grpcopts "github.com/webitel/cases/model/options/grpc"
	"github.com/webitel/cases/model/options/grpc/shared"
	"github.com/webitel/cases/util"
	"github.com/webitel/webitel-go-kit/etag"
	"google.golang.org/grpc/metadata"
)

const (
	dynamicGroup      = "dynamic"
	caseObjScope      = model.ScopeCases
	defaultLogTimeout = 5 * time.Second
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
		{Name: "close_reason", Default: true},
		{Name: "close_result", Default: true},
		{Name: "rating", Default: true},
		{Name: "rating_comment", Default: true},
		{Name: "sla_condition", Default: true},
		{Name: "service", Default: true},
		{Name: "status_condition", Default: true},
		{Name: "sla", Default: true},
		{Name: "comments", Default: false},
		{Name: "links", Default: false},
		{Name: "files", Default: false},
		{Name: "related", Default: false},
		{Name: "resolved_at", Default: true},
		{Name: "reacted_at", Default: true},
		{Name: "difference_in_reaction", Default: true},
		{Name: "difference_in_resolve", Default: true},
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
	logAttributes := slog.Group(
		"context",
		slog.Int64(
			"user_id",
			searchOpts.GetAuthOpts().GetUserId(),
		),
		slog.Int64(
			"domain_id",
			searchOpts.GetAuthOpts().GetDomainId(),
		),
	)
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
		return nil, deferr.DatabaseError
	}
	err = c.NormalizeResponseCases(list, req)
	if err != nil {
		slog.ErrorContext(ctx, err.Error(), logAttributes)
		return nil, deferr.ResponseNormalizingError
	}

	return list, nil
}

func (c *CaseService) LocateCase(ctx context.Context, req *cases.LocateCaseRequest) (*cases.Case, error) {
	searchOpts, err := grpcopts.NewLocateOptions(
		ctx,
		grpcopts.WithFields(req, CaseMetadata,
			util.DeduplicateFields,
			util.ParseFieldsForEtag,
			util.EnsureIdField),
		grpcopts.WithIDsAsEtags(etag.EtagCase, req.GetEtag()),
	)
	if err != nil {
		return nil, NewBadRequestError(err)
	}
	logAttributes := slog.Group(
		"context",
		slog.Int64(
			"user_id",
			searchOpts.GetAuthOpts().GetUserId(),
		),
		slog.Int64(
			"domain_id", searchOpts.GetAuthOpts().GetDomainId(),
		),
	)
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
		return nil, deferr.ResponseNormalizingError
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

	var statusCondition *cases.StatusCondition
	if req.Input.StatusCondition != nil {
		statusCondition = &cases.StatusCondition{Id: req.Input.StatusCondition.Id}
	}
	res := &cases.Case{
		// Used if explicitly set the case creator / updater instead of deriving it from the auth token.
		CreatedBy:        req.Input.GetUserID(),
		Subject:          req.Input.Subject,
		Description:      req.Input.Description,
		ContactInfo:      req.Input.ContactInfo,
		Assignee:         req.Input.Assignee,
		Reporter:         req.Input.Reporter,
		Source:           &cases.SourceTypeLookup{Id: req.Input.Source.GetId()},
		Impacted:         req.Input.Impacted,
		Group:            &cases.ExtendedLookup{Id: req.Input.Group.GetId()},
		Status:           req.Input.Status,
		StatusCondition:  statusCondition,
		CloseReason:      req.Input.GetCloseReason(),
		CloseResult:      req.Input.GetCloseResult(),
		CloseReasonGroup: req.Input.GetCloseReasonGroup(),
		Priority:         &cases.Priority{Id: req.Input.Priority.GetId()},
		Rating:           req.Input.Rating,
		RatingComment:    req.Input.RatingComment,
		Service:          req.Input.Service,
		Links:            links,
		Related:          related,
		Custom:           req.Input.GetCustom(),
	}

	createOpts, err := grpcopts.NewCreateOptions(
		ctx,
		grpcopts.WithCreateFields(
			req,
			CaseMetadata.CopyWithAllFieldsSetToDefault(),
			util.DeduplicateFields,
			util.ParseFieldsForEtag,
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
		var DBBadRequestError *cerror.DBBadRequestError
		if errors.As(err, &DBBadRequestError) {
			return nil, cerror.NewBadRequestError(
				"app.case.create_case.param_required",
				err.Error(),
			)
		}
		slog.ErrorContext(ctx, err.Error(), logAttributes)
		return nil, deferr.DatabaseError
	}
	et, err := etag.EncodeEtag(etag.EtagCase, res.Id, res.Ver)
	if err != nil {
		return nil, err
	}
	res.Etag = et

	// save before normalize
	roleIds := res.GetRoleIds()
	id := res.GetId()

	//* Handle dynamic group update if applicable
	res, err = c.handleDynamicGroup(ctx, res)
	if err != nil {
		return nil, err
	}

	// * CREATE require all fields set to true incase we need to calculate dynamic condition
	req.Fields = createOpts.Fields
	err = c.NormalizeResponseCase(res, req)
	if err != nil {
		slog.ErrorContext(ctx, err.Error(), logAttributes)
		return nil, deferr.ResponseNormalizingError
	}

	if notifyErr := c.app.watcherManager.Notify(
		model.ScopeCases,
		EventTypeCreate,
		NewCaseWatcherData(
			createOpts.GetAuthOpts(),
			res,
			id,
			roleIds,
		),
	); notifyErr != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("could not notify case creation: %s, ", notifyErr.Error()), logAttributes)
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
		return nil, cerror.NewBadRequestError("app.case.update.invalid_etag", "Invalid et")
	}

	updateOpts, err := grpcopts.NewUpdateOptions(
		ctx,
		grpcopts.WithUpdateFields(
			req,
			CaseMetadata.CopyWithAllFieldsSetToDefault(),
			util.DeduplicateFields,
			util.ParseFieldsForEtag,
			func(fields []string) []string {
				for i, v := range fields {
					if v == "related" {
						fields = append(fields[:i], fields[i+1:]...)
					}
				}
				return fields
			}),
		grpcopts.WithUpdateEtag(&tag),
		grpcopts.WithUpdateMasker(req),
	)
	if err != nil {
		return nil, NewBadRequestError(err)
	}

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
		// Used if explicitly set the case creator / updater instead of deriving it from the auth token.
		UpdatedBy:        req.Input.GetUserID(),
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
		CloseReason:      req.Input.GetCloseReason(),
		CloseResult:      req.Input.GetCloseResult(),
		Rating:           req.Input.GetRating(),
		RatingComment:    req.Input.GetRatingComment(),
		Service:          req.Input.GetService(),
		Custom:           req.Input.GetCustom(),
	}

	res, err := c.app.Store.Case().Update(updateOpts, upd)
	if err != nil {
		var DBNoRowsError *cerror.DBNoRowsError
		switch {
		case errors.As(err, &DBNoRowsError):
			return nil, cerror.NewBadRequestError("app.case.update.invalid_etag", "Invalid et")
		}
		slog.ErrorContext(ctx, err.Error())
		return nil, deferr.DatabaseError
	}

	et, err := etag.EncodeEtag(etag.EtagCase, res.Id, res.Ver)
	if err != nil {
		return nil, err
	}
	res.Etag = et

	// *Handle dynamic group update if applicable
	res, err = c.handleDynamicGroup(ctx, res)
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
		return nil, deferr.ResponseNormalizingError
	}

	if notifyErr := c.app.watcherManager.Notify(model.ScopeCases, EventTypeUpdate, NewCaseWatcherData(updateOpts.GetAuthOpts(), upd, id, roleIds)); notifyErr != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("could not notify case update: %s, ", notifyErr.Error()), logAttributes)
	}

	return res, nil
}

// handleDynamicGroup checks if a dynamic group is needed and updates the case accordingly.
func (c *CaseService) handleDynamicGroup(
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
		input, err = c.resolveDynamicGroup(ctx, input, res)
		if err != nil {
			return nil, err
		}
	}

	// Return the potentially updated case
	return input, nil
}

func (c *CaseService) resolveDynamicGroup(
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

	et, ok := caseMap["case.et"].(string)
	if !ok {
		return nil, deferr.ForbiddenError
	}

	// Iterate over all dynamic conditions in the group…
	for _, condition := range inputGroup.Group.Conditions {
		// Evaluate the condition against the case map
		if evaluateDynamicCondition(caseMap, condition.Expression) {

			groupID, err := strconv.Atoi(condition.Group.GetId())
			if err != nil {
				return nil, deferr.ResponseNormalizingError
			}

			var assigneeID int
			if condition.Assignee != nil {
				assigneeID, err = strconv.Atoi(condition.Assignee.GetId())
				if err != nil {
					return nil, deferr.ResponseNormalizingError
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
					Etag:     et,
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
		return nil, fmt.Errorf("app.case.resolveDynamicGroup: failed to convert group ID to integer: %w", err)
	}

	req := &cases.UpdateCaseRequest{
		XJsonMask: groupXJsonMask,
		Input: &cases.InputCase{
			Group: &cases.Lookup{Id: int64(groupID)},
			Etag:  et,
		},
	}
	// Final update if a default group was assigned
	updCase, err := c.UpdateCase(ctx, req)
	if err != nil {
		return nil, err
	}

	return updCase, nil
}

// Converts a Case object to map[string]interface{} with "case." prefixed keys and lowercase values (except case.etag).
// Helper method for dynamic contact group resolving.
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

// Recursively adds prefixed keys for nested maps (all keys and values converted to lowercase except case.etag).
// Helper method for dynamic contact group resolving.
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

// Evaluates complex condition strings with support for AND (&&) and OR (||) operators using bitwise operations.
// Helper method for dynamic contact group resolving.
func evaluateDynamicCondition(caseMap map[string]any, condition string) bool {
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

// Evaluates a single condition string, e.g., "case.assignee.name == 'John Wick'".
// Helper method for dynamic contact group resolving.
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

// Resolves a field path (e.g., "case.assignee.name") to its value in the map.
// Helper method for dynamic contact group resolving.
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
	logAttributes := slog.Group(
		"context",
		slog.Int64(
			"user_id",
			deleteOpts.GetAuthOpts().GetUserId(),
		),
		slog.Int64(
			"domain_id",
			deleteOpts.GetAuthOpts().GetDomainId(),
		),
		slog.Int64(
			"case_id",
			tag.GetOid()),
	)

	err = c.app.Store.Case().Delete(deleteOpts)
	if err != nil {
		var DBNoRowsError *cerror.DBNoRowsError
		switch {
		case errors.As(err, &DBNoRowsError):
			return nil, cerror.NewBadRequestError("app.case.delete.invalid_etag", "Invalid etag or insufficient rights")
		}
		slog.ErrorContext(ctx, err.Error(), logAttributes)
		return nil, deferr.DatabaseError
	}
	deleteCase := &cases.Case{
		Id:   tag.GetOid(),
		Ver:  tag.GetVer(),
		Etag: req.Etag,
	}

	if notifyErr := c.app.watcherManager.Notify(
		model.ScopeCases,
		EventTypeDelete,
		NewCaseWatcherData(
			deleteOpts.GetAuthOpts(),
			deleteCase,
			tag.GetOid(),
			nil,
		),
	); notifyErr != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("could not notify case deletion: %s, ", notifyErr.Error()), logAttributes)
	}
	return nil, nil
}

func NewCaseService(app *App) (*CaseService, cerror.AppError) {
	if app == nil {
		return nil, cerror.NewBadRequestError(
			"app.case.new_case_service.check_args.app",
			"unable to init case service, app is nil",
		)
	}

	watcher := NewDefaultWatcher()

	if app.config.LoggerWatcher.Enabled {

		obs, err := NewLoggerObserver(app.wtelLogger, caseObjScope, defaultLogTimeout)
		if err != nil {
			return nil, cerror.NewInternalError("app.case.new_case_service.create_observer.app", err.Error())
		}
		watcher.Attach(EventTypeCreate, obs)
		watcher.Attach(EventTypeUpdate, obs)
		watcher.Attach(EventTypeDelete, obs)
	}

	if app.config.FtsWatcher.Enabled {
		ftsObserver, err := NewFullTextSearchObserver(app.ftsClient, caseObjScope, formCaseFtsModel)
		if err != nil {
			return nil, cerror.NewInternalError("app.case.new_case_service.create_fts_observer.app", err.Error())
		}
		watcher.Attach(EventTypeCreate, ftsObserver)
		watcher.Attach(EventTypeUpdate, ftsObserver)
		watcher.Attach(EventTypeDelete, ftsObserver)
	}

	if app.config.TriggerWatcher.Enabled {
		mq, err := NewTriggerObserver(app.rabbit, app.config.TriggerWatcher, formCaseTriggerModel, slog.With(
			slog.Group("context",
				slog.String("scope", "watcher")),
		))

		if err != nil {
			return nil, cerror.NewInternalError("app.case.new_case_service.create_mq_observer.app", err.Error())
		}
		watcher.Attach(EventTypeCreate, mq)
		watcher.Attach(EventTypeUpdate, mq)
		watcher.Attach(EventTypeDelete, mq)
	}

	app.watcherManager.AddWatcher(caseObjScope, watcher)

	return &CaseService{app: app, logger: app.wtelLogger.GetObjectedLogger(CaseMetadata.GetMainScopeName())}, nil
}

func (c *CaseService) ValidateUpdateInput(
	input *cases.InputCase,
	xJsonMask []string,
) cerror.AppError {
	if input.Etag == "" {
		return cerror.NewBadRequestError("app.case.update_case.etag_required", "Etag is required")
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
		//case "close_reason":
		//	if closeReason := input.GetCloseReason(); closeReason != nil && closeReason.GetId() == 0 {
		//		return cerror.NewBadRequestError("app.case.update_case.close_reason_group_required", "Close Reason group is required")
		//	}
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
			// default:
			// 	if jpath, ok := strings.CutPrefix(field, "custom"); ok {

			// 	}
		}
	}

	return nil
}

// region UTILITY

func (c *CaseService) ValidateCreateInput(input *cases.InputCreateCase) cerror.AppError {
	if input.Subject == "" {
		return cerror.NewBadRequestError("app.case.create_case.subject_required", "Case subject is required")
	}
	if input.Source.GetId() == 0 {
		return cerror.NewBadRequestError("app.case.create_case.source_required", "Case source is required")
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
				err = util.NormalizeEtags(
					etag.EtagRelatedCase,
					true,
					false,
					false,
					&related.Etag,
					&related.Id,
					&related.Ver,
				)
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
			err = util.NormalizeEtags(
				etag.EtagCaseLink,
				true,
				false,
				false,
				&link.Etag,
				&link.Id,
				&link.Ver,
			)
			if err != nil {
				return err
			}
		}
	}
	if re.Related != nil {
		for _, related := range re.Related.Data {
			err = util.NormalizeEtags(
				etag.EtagRelatedCase,
				true,
				false,
				false,
				&related.Etag,
				&related.Id,
				&related.Ver,
			)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func formCaseFtsModel(item *cases.Case, params map[string]any) (*model.FtsCase, error) {
	roles, ok := params["role_ids"].([]int64)
	if !ok {
		return nil, fmt.Errorf("role ids required for FTS model")
	}
	m := &model.FtsCase{
		Description:   item.GetDescription(),
		RoleIds:       roles,
		Subject:       item.GetSubject(),
		ContactInfo:   item.GetContactInfo(),
		CreatedAt:     item.GetCreatedAt(),
		RatingComment: item.GetRatingComment(),
		CloseResult:   item.GetCloseResult(),
	}
	return m, nil
}

func formCaseTriggerModel(item *cases.Case) (*model.CaseAMQPMessage, error) {
	m := &model.CaseAMQPMessage{
		Case: item,
	}
	return m, nil
}

type CaseWatcherData struct {
	case_      *cases.Case
	CaseString string `json:"case"`
	DomainId   int64  `json:"domain_id"`
	Args       map[string]any
}

func NewCaseWatcherData(session auth.Auther, case_ *cases.Case, caseId int64, roleIds []int64) *CaseWatcherData {
	return &CaseWatcherData{case_: case_, Args: map[string]any{"session": session, "obj": case_, "id": caseId, "role_ids": roleIds}}
}

func (wd *CaseWatcherData) GetArgs() map[string]any {
	return wd.Args
}
