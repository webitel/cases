package app

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/google/cel-go/cel"
	"github.com/webitel/cases/internal/api_handler/grpc/options"
	"github.com/webitel/cases/internal/api_handler/grpc/options/shared"
	"github.com/webitel/webitel-go-kit/pkg/filters"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/structpb"

	wlogger "github.com/webitel/webitel-go-kit/infra/logger_client"
	"github.com/webitel/webitel-go-kit/pkg/etag"
	watcherkit "github.com/webitel/webitel-go-kit/pkg/watcher"

	"github.com/webitel/cases/api/cases"
	webitelgo "github.com/webitel/cases/api/webitel-go/contacts"
	"github.com/webitel/cases/auth"
	"github.com/webitel/cases/internal/api_handler/grpc"
	"github.com/webitel/cases/internal/api_handler/grpc/utils"
	"github.com/webitel/cases/internal/errors"
	"github.com/webitel/cases/internal/model"
	"github.com/webitel/cases/util"
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
		{Name: "updated_by", Default: true},
		{Name: "updated_at", Default: true},
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
		{Name: "dc", Default: false},
		{Name: "diff", Default: true},
	}, grpc.CaseCommentMetadata, grpc.CaseLinkMetadata, RelatedCaseMetadata)

	resolutionTimeSO = &options.SearchOptions{
		Context: context.Background(),
		Fields:  util.ParseFieldsForEtag(util.RemoveSliceElement(CaseMetadata.GetAllFields(), "related")),
	}
)

type CaseService struct {
	app *App
	cases.UnimplementedCasesServer
	logger        *wlogger.ObjectedLogger
	filtrationEnv *cel.Env
}

func (c *CaseService) SearchCases(ctx context.Context, req *cases.SearchCasesRequest) (*cases.CaseList, error) {
	searchOpts, err := options.NewSearchOptions(
		ctx,
		options.WithSearch(req),
		options.WithPagination(req),
		options.WithFields(
			req,
			CaseMetadata,
			util.DeduplicateFields,
			util.ParseFieldsForEtag,
			util.EnsureIdField,
		),
		options.WithFiltersV1(c.filtrationEnv, req.GetFiltersV1()),
		options.WithFilters(req.GetFilters()),
		options.WithIDsAsEtags(etag.EtagCase, req.GetIds()...),
		options.WithSort(req),
		options.WithQin(req.GetQin()),
	)
	if err != nil {
		return nil, err
	}

	if contactID := req.GetContactId(); contactID != "" {
		searchOpts.Filters = append(searchOpts.Filters, fmt.Sprintf("contact=%s", contactID))
	}

	list, err := c.app.Store.Case().List(searchOpts)
	if err != nil {
		return nil, err
	}
	err = c.NormalizeResponseCases(list, req)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (c *CaseService) LocateCase(ctx context.Context, req *cases.LocateCaseRequest) (*cases.Case, error) {
	searchOpts, err := options.NewLocateOptions(
		ctx,
		options.WithFields(req, CaseMetadata,
			util.DeduplicateFields,
			util.ParseFieldsForEtag,
			util.EnsureIdField,
			util.EnsureCustomField,
		),
		options.WithIDsAsEtags(etag.EtagCase, req.GetEtag()),
	)
	if err != nil {
		return nil, err
	}
	list, err := c.app.Store.Case().List(searchOpts)
	if err != nil {
		return nil, err
	}
	if len(list.Items) == 0 {
		return nil, errors.NotFound("entity not found")
	}
	err = c.NormalizeResponseCases(list, req)
	if err != nil {
		return nil, err
	}
	return list.Items[0], nil
}

// lookupToService converts a Lookup to a Service struct
func lookupToService(lookup *cases.Lookup) *cases.Service {
	if lookup == nil {
		return nil
	}
	return &cases.Service{
		Id:   lookup.Id,
		Name: lookup.Name,
	}
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
			var relatedID int64
			if inputRelated.GetRelatedTo() != "" {
				relatedID, _ = strconv.ParseInt(inputRelated.GetRelatedTo(), 10, 64)
			}
			relatedItems[i] = &cases.RelatedCase{
				Id:           relatedID,
				Etag:         inputRelated.GetEtag(),
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
	//     * Set based on the provided status....
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
		Service:          lookupToService(req.Input.Service),
		Links:            links,
		Related:          related,
		Custom:           req.Input.GetCustom(),
	}

	createOpts, err := options.NewCreateOptions(
		ctx,
		options.WithCreateFields(
			req,
			CaseMetadata.CopyWithAllFieldsSetToDefault(),
			util.DeduplicateFields,
			util.ParseFieldsForEtag,
		),
	)
	if err != nil {
		return nil, err
	}

	// Add override user ID after options are built
	err = options.WithCreateOverrideUserID(req.Input.UserID.GetId())(createOpts)
	if err != nil {
		return nil, err
	}

	logAttributes := slog.Group(
		"context",
		slog.Int64("user_id", createOpts.GetAuthOpts().GetUserId()),
		slog.Int64("domain_id", createOpts.GetAuthOpts().GetDomainId()),
	)

	res, err = c.app.Store.Case().Create(createOpts, res)
	if err != nil {
		return nil, err
	}
	res.Etag, err = etag.EncodeEtag(etag.EtagCase, res.Id, res.Ver)
	if err != nil {
		return nil, err
	}
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
		return nil, err
	}

	ip := createOpts.GetAuthOpts().GetUserIp()
	if ip == "" {
		ip = "unknown"
	}

	message, _ := wlogger.NewMessage(
		createOpts.GetAuthOpts().GetUserId(),
		ip,
		wlogger.CreateAction,
		strconv.Itoa(int(res.GetId())),
		res,
	)

	_, err = c.logger.SendContext(ctx, createOpts.GetAuthOpts().GetDomainId(), message)
	if err != nil {
		return nil, err
	}

	if req.DisableTrigger == false {
		if notifyErr := c.app.watcherManager.Notify(
			model.ScopeCases,
			watcherkit.EventTypeCreate,
			NewCaseWatcherData(
				createOpts.GetAuthOpts(),
				res,
				id,
				roleIds,
			),
		); notifyErr != nil {
			slog.ErrorContext(ctx, fmt.Sprintf("could not notify case creation: %s", notifyErr.Error()), logAttributes)
		}
	}

	return res, nil
}

func (c *CaseService) UpdateCase(ctx context.Context, req *cases.UpdateCaseRequest) (*cases.UpdateCaseResponse, error) {
	var original *cases.Case

	appErr := c.ValidateUpdateInput(req.Input, req.XJsonMask)
	if appErr != nil {
		return nil, appErr
	}

	tag, err := etag.EtagOrId(etag.EtagCase, req.Input.Etag)
	if err != nil {
		return nil, err
	}

	updateOpts, err := options.NewUpdateOptions(
		ctx,
		options.WithUpdateFields(
			req,
			CaseMetadata.CopyWithAllFieldsSetToDefault(),
			util.DeduplicateFields,
			util.ParseFieldsForEtag,
		),
		options.WithUpdateEtag(&tag),
		options.WithUpdateMasker(req),
	)
	if err != nil {
		return nil, err
	}

	// Add override user ID after options are built
	err = options.WithUpdateOverrideUserID(req.Input.UserID.GetId())(updateOpts)
	if err != nil {
		return nil, err
	}

	logAttributes := slog.Group(
		"context",
		slog.Int64("user_id", updateOpts.GetAuthOpts().GetUserId()),
		slog.Int64("domain_id", updateOpts.GetAuthOpts().GetDomainId()),
		slog.Int64("case_id", tag.GetOid()),
	)

	upd := &cases.Case{
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
		Service:          lookupToService(req.Input.GetService()),
		Custom:           req.Input.GetCustom(),
	}
	if reporter := upd.Reporter; reporter != nil && reporter.GetId() == 0 {
		upd.Reporter = nil
	}

	// If diff is requested, get original case before update
	if util.ContainsField(updateOpts.GetFields(), "diff") {
		locateReq := &cases.LocateCaseRequest{
			Etag: req.Input.Etag,
		}
		var err error
		original, err = c.LocateCase(ctx, locateReq)
		if err != nil {
			return nil, err
		}
	}

	output, err := c.app.Store.Case().Update(updateOpts, upd)
	if err != nil {
		return nil, err
	}
	output.Etag, err = etag.EncodeEtag(etag.EtagCase, output.Id, output.Ver)
	if err != nil {
		return nil, err
	}

	output, err = c.handleDynamicGroup(ctx, output)
	if err != nil {
		return nil, err
	}

	req.Fields = updateOpts.Fields
	err = c.NormalizeResponseCase(output, req)
	if err != nil {
		return nil, err
	}

	ip := updateOpts.GetAuthOpts().GetUserIp()
	if ip == "" {
		ip = "unknown"
	}

	message, _ := wlogger.NewMessage(
		updateOpts.GetAuthOpts().GetUserId(),
		ip,
		wlogger.UpdateAction,
		strconv.Itoa(int(output.GetId())),
		output,
	)

	_, err = c.logger.SendContext(ctx, updateOpts.GetAuthOpts().GetDomainId(), message)
	if err != nil {
		return nil, err
	}

	if req.DisableTrigger == false {
		if notifyErr := c.app.watcherManager.Notify(
			model.ScopeCases,
			watcherkit.EventTypeUpdate,
			NewCaseWatcherData(
				updateOpts.GetAuthOpts(),
				upd,
				output.Id,
				output.GetRoleIds(),
			),
		); notifyErr != nil {
			slog.ErrorContext(
				ctx,
				fmt.Sprintf("could not notify case update: %s", notifyErr.Error()), logAttributes)
		}
	}

	// region diff building

	var changes []*cases.FieldChange
	if util.ContainsField(req.Fields, "diff") && original != nil {
		changes = BuildCaseDiff(original, output)
	}

	return &cases.UpdateCaseResponse{
		Case:    output,
		Changes: changes,
	}, nil
}

func BuildCaseDiff(original, updated *cases.Case) []*cases.FieldChange {
	var changes []*cases.FieldChange

	compare := func(field string, oldVal, newVal any) {
		if !reflect.DeepEqual(oldVal, newVal) {
			changes = append(changes, &cases.FieldChange{
				Field:    field,
				OldValue: toProtoValue(oldVal),
				NewValue: toProtoValue(newVal),
			})
		}
	}

	// Compare scalar fields
	compare("etag", original.Etag, updated.Etag)
	compare("subject", original.Subject, updated.Subject)
	compare("description", original.Description, updated.Description)
	compare("contact_info", original.ContactInfo, updated.ContactInfo)
	compare("close_result", original.CloseResult, updated.CloseResult)
	compare("rating", original.Rating, updated.Rating)
	compare("rating_comment", original.RatingComment, updated.RatingComment)
	compare("assignee", original.Assignee, updated.Assignee)
	compare("reporter", original.Reporter, updated.Reporter)
	compare("impacted", original.Impacted, updated.Impacted)
	compare("group", original.Group, updated.Group)
	compare("status", original.Status, updated.Status)
	compare("priority", original.Priority, updated.Priority)
	compare("source", original.Source, updated.Source)
	compare("service", original.Service, updated.Service)
	compare("close_reason", original.CloseReason, updated.CloseReason)
	compare("userID", original.UpdatedBy, updated.UpdatedBy)
	compare("status_condition", original.StatusCondition, updated.StatusCondition)

	//-------- Custom -------- //
	compare("custom", original.Custom, updated.Custom)

	return changes
}

func toProtoValue(v any) *structpb.Value {
	val, err := structpb.NewValue(v)
	if err != nil {
		// Fallback to string representation if not convertible
		strVal := fmt.Sprintf("%v", v)
		val, _ = structpb.NewValue(strVal)
	}

	return val
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
			return nil, errors.InvalidArgument("Not found context")
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
		return nil, err
	}

	et, ok := caseMap["case.etag"].(string)
	if !ok {
		return nil, errors.NotFound("not found")
	}

	// Iterate over all dynamic conditions in the group…
	for _, condition := range inputGroup.Group.Conditions {
		// Evaluate the condition against the case map
		if evaluateDynamicCondition(caseMap, condition.Expression) {

			groupID, err := strconv.Atoi(condition.Group.GetId())
			if err != nil {
				return nil, err
			}

			var assigneeID int
			if condition.Assignee != nil {
				assigneeID, err = strconv.Atoi(condition.Assignee.GetId())
				if err != nil {
					return nil, err
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

			return updCase.Case, nil
		}
	}

	groupID, err := strconv.Atoi(inputGroup.Group.DefaultGroup.GetId())
	if err != nil {
		return nil, err
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

	return updCase.Case, nil
}

// Converts a Case object to map[string]interface{} with "case." prefixed keys and lowercase values (except case.etag).
// Helper method for dynamic contact group resolving.
func caseToMap(caseObj any) (map[string]any, error) {
	caseJSON, err := json.Marshal(caseObj)
	if err != nil {
		return nil, err
	}

	var caseMap map[string]any
	err = json.Unmarshal(caseJSON, &caseMap)
	if err != nil {
		return nil, err
	}

	// Create a new map with "case." prefixed keys, handling nested fields
	prefixedMap := make(map[string]any)
	addPrefixedKeys(prefixedMap, caseMap, "case")

	return prefixedMap, nil
}

// Recursively adds prefixed keys for nested maps (all keys and values converted to lowercase except case.etag).
// Helper method for dynamic contact group resolving.
func addPrefixedKeys(dest, source map[string]any, prefix string) {
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

// Evaluates complex condition strings with support for And (&&) and Or (||) operators using bitwise operations.
// Helper method for dynamic contact group resolving.
func evaluateDynamicCondition(caseMap map[string]any, condition string) bool {
	// Convert condition to lowercase to ensure case-insensitive matching
	condition = strings.ToLower(condition)

	// Split conditions by Or (||)
	orConditions := strings.Split(condition, "||")
	for _, orCondition := range orConditions {
		// Split each Or condition into And (&&) conditions
		andConditions := strings.Split(orCondition, "&&")
		// Use a bitmask to track whether all And conditions are met
		var andMask uint = 0
		for i, andCondition := range andConditions {
			andCondition = strings.TrimSpace(andCondition)
			if evaluateSingleCondition(caseMap, andCondition) {
				// Set the corresponding bit in the mask
				andMask |= 1 << i
			}
		}
		// Check if all And conditions are met (all bits set in the mask)
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
func resolveFieldPath(data map[string]any, path string) any {
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
		return nil, errors.InvalidArgument("Etag is required")
	}
	tag, err := etag.EtagOrId(etag.EtagCase, req.Etag)
	if err != nil {
		return nil, errors.InvalidArgument("Invalid case etag")
	}
	deleteOpts, err := options.NewDeleteOptions(ctx, options.WithDeleteID(tag.GetOid()))
	if err != nil {
		return nil, err
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
		return nil, err
	}
	deleteCase := &cases.Case{
		Id:   tag.GetOid(),
		Ver:  tag.GetVer(),
		Etag: req.Etag,
	}

	message, _ := wlogger.NewMessage(
		deleteOpts.GetAuthOpts().GetUserId(),
		deleteOpts.GetAuthOpts().GetUserIp(),
		wlogger.DeleteAction,
		strconv.Itoa(int(tag.GetOid())),
		deleteCase,
	)

	_, err = c.logger.SendContext(ctx, deleteOpts.GetAuthOpts().GetDomainId(), message)
	if err != nil {
		return nil, err
	}

	if notifyErr := c.app.watcherManager.Notify(
		model.ScopeCases,
		watcherkit.EventTypeDelete,
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

func NewCaseService(app *App) (*CaseService, error) {
	if app == nil {
		return nil, errors.InvalidArgument(
			"unable to init case service, app is nil",
		)
	}
	objectedLogger, err := app.wtelLogger.GetObjectedLogger(CaseMetadata.GetMainScopeName())
	if err != nil {
		return nil, err
	}
	service := &CaseService{
		app:    app,
		logger: objectedLogger,
	}
	// Create a new CEL environment for case filtering
	filtrationEnv, err := cel.NewEnv(filters.ProtoToCELVariables(&cases.Case{})...)
	if err != nil {
		return nil, err
	}
	service.filtrationEnv = filtrationEnv

	watcher := watcherkit.NewDefaultWatcher()

	//if app.config.LoggerWatcher.Enabled {
	//
	//	obs, err := NewLoggerObserver(app.wtelLogger, caseObjScope, defaultLogTimeout)
	//	if err != nil {
	//		return nil, cerror.NewInternalError("app.case.new_case_service.create_observer.app", err.Error())
	//	}
	//	watcher.Attach(watcherkit.EventTypeCreate, obs)
	//	watcher.Attach(watcherkit.EventTypeUpdate, obs)
	//	watcher.Attach(watcherkit.EventTypeDelete, obs)
	//}

	if app.config.FtsWatcher.Enabled {
		ftsObserver, err := NewFullTextSearchObserver(app.ftsClient, caseObjScope, formCaseFtsModel)
		if err != nil {
			return nil, err
		}

		watcher.Attach(watcherkit.EventTypeCreate, ftsObserver)
		watcher.Attach(watcherkit.EventTypeUpdate, ftsObserver)
		watcher.Attach(watcherkit.EventTypeDelete, ftsObserver)
	}

	if app.config.TriggerWatcher.Enabled {
		mq, err := NewTriggerObserver(app.rabbitPublisher, app.config.TriggerWatcher, formCaseTriggerModel, slog.With(
			slog.Group("context",
				slog.String("scope", "watcher")),
		))
		if err != nil {
			return nil, err
		}

		watcher.Attach(watcherkit.EventTypeCreate, mq)
		watcher.Attach(watcherkit.EventTypeUpdate, mq)
		watcher.Attach(watcherkit.EventTypeDelete, mq)
		watcher.Attach(watcherkit.EventTypeResolutionTime, mq)
		app.caseResolutionTimer = NewTimerTask[*App](time.Duration(
			app.config.TriggerWatcher.ResolutionCheckInterval)*time.Second,
			service.scheduleResolutionTime,
			app,
		)

		app.caseResolutionTimer.Start()
	}

	app.watcherManager.AddWatcher(caseObjScope, watcher)

	return service, nil
}

func (c *CaseService) ValidateUpdateInput(
	input *cases.InputCase,
	xJsonMask []string,
) error {
	if input == nil {
		return errors.InvalidArgument("Input is required")
	}
	if input.Etag == "" {
		return errors.InvalidArgument("Etag is required")
	}

	// Iterate over xJsonMask and validate corresponding fields
	// Validating fields passed for updating
	for _, field := range xJsonMask {
		switch field {
		case "subject":
			if input.Subject == "" {
				return errors.InvalidArgument("Subject is required")
			}
		case "status":
			if input.Status.GetId() == 0 {
				return errors.InvalidArgument("Status is required")
			}
		case "priority":
			if input.Priority.GetId() == 0 {
				return errors.InvalidArgument("Priority is required")
			}
		case "source":
			if input.Source.GetId() == 0 {
				return errors.InvalidArgument("Source is required")
			}
		case "service":
			if input.Service.GetId() == 0 {
				return errors.InvalidArgument("Service is required")
			}
		}
	}

	return nil
}

// region UTILITY

func (c *CaseService) ValidateCreateInput(input *cases.InputCreateCase) error {
	if input == nil {
		return errors.InvalidArgument("Input is required")
	}
	if input.Subject == "" {
		return errors.InvalidArgument("Case subject is required")
	}
	if input.Source.GetId() == 0 {
		return errors.InvalidArgument("Case source is required")
	}
	if input.Service.GetId() == 0 {
		return errors.InvalidArgument("Case service is required")
	}
	if input.Reporter.GetId() == 0 {
		return errors.InvalidArgument("Case reporter is required")
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
					etag.EtagCase,
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
			if com.Id == 0 {
				// FIXME temp sanitize: skip comments with empty ID
				continue
			}
			err = util.NormalizeEtags(etag.EtagCaseComment, true, false, false, &com.Etag, &com.Id, &com.Ver)
			if err != nil {
				return err
			}
		}
	}

	if re.Links != nil {
		for _, link := range re.Links.Items {
			if link.Id == 0 {
				// FIXME temp sanitize: skip links with empty ID
				continue
			}
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
			if related.Id == 0 {
				// FIXME temp sanitize: skip related cases with empty ID
				continue
			}
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

func formCaseLinkTriggerModel(item *model.CaseLink) (*model.CaseLinkAMQPMessage, error) {
	// Convert model.CaseLink to cases.CaseLink for AMQP message
	protoLink := &cases.CaseLink{
		Id:        item.Id,
		Ver:       item.Ver,
		Name:      *item.Name,
		Url:       item.Url,
		CreatedBy: utils.MarshalLookup(item.Author),
		UpdatedBy: utils.MarshalLookup(item.Editor),
		Author:    utils.MarshalLookup(item.Contact),
		CreatedAt: utils.MarshalTime(item.CreatedAt),
		UpdatedAt: utils.MarshalTime(item.UpdatedAt),
	}
	m := &model.CaseLinkAMQPMessage{
		CaseLink: protoLink,
	}

	return m, nil
}

func formCaseCommentTriggerModel(item *model.CaseComment) (*model.CaseCommentAMQPMessage, error) {
	protoComment := &cases.CaseComment{
		Id:        item.Id,
		Ver:       item.Ver,
		Text:      item.Text,
		CreatedBy: utils.MarshalLookup(item.Author),
		Author:    utils.MarshalLookup(item.Contact),
		CreatedAt: utils.MarshalTime(item.CreatedAt),
		UpdatedAt: utils.MarshalTime(item.UpdatedAt),
		CanEdit:   item.CanEdit,
		CaseId:    item.CaseId,
		UpdatedBy: utils.MarshalLookup(item.Editor),
		Edited:    item.Edited,
	}
	m := &model.CaseCommentAMQPMessage{
		CaseComment: protoComment,
	}

	return m, nil
}

func formCaseFiletriggerModel(item *cases.File) (*model.CaseFileAMQPMessage, error) {
	m := &model.CaseFileAMQPMessage{
		CaseFile: item,
	}

	return m, nil
}

type CaseWatcherData struct {
	case_ *cases.Case
	Args  map[string]any
}

func NewCaseWatcherData(session auth.Auther, case_ *cases.Case, caseId int64, roleIds []int64) *CaseWatcherData {
	return &CaseWatcherData{case_: case_, Args: map[string]any{
		"session":   session,
		"obj":       case_,
		"id":        caseId,
		"role_ids":  roleIds,
		"domain_id": case_.Dc,
	}}
}

func (wd *CaseWatcherData) GetArgs() map[string]any {
	return wd.Args
}

func (c *CaseService) scheduleResolutionTime(app *App) {
	var css []*cases.Case
	var err error
	var retry bool

	if css, retry, err = app.Store.Case().SetOverdueCases(resolutionTimeSO); err != nil {
		slog.Error(errors.Details(err))
		return
	}

	for _, cs := range css {
		err = c.NormalizeResponseCase(cs, resolutionTimeSO)
		if err != nil {
			slog.Error(errors.Details(err))
			continue
		}

		if notifyErr := app.watcherManager.Notify(
			model.ScopeCases,
			watcherkit.EventTypeResolutionTime,
			NewCaseWatcherData(
				nil,
				cs,
				cs.Id,
				nil,
			),
		); notifyErr != nil {
			slog.Error(fmt.Sprintf("could not notify case resolution time: %s", notifyErr.Error()))
		}
	}

	if retry {
		c.scheduleResolutionTime(app)
	}
}
