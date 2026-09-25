package postgres

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/webitel/cases/internal/model/options"

	dberr "github.com/webitel/cases/internal/errors"
	"github.com/webitel/cases/internal/model"
	"github.com/webitel/cases/internal/store"
	storeutil "github.com/webitel/cases/internal/store/util"
)

var CaseTimelineFields = []string{
	"call", "chat", "email",
}

type CaseTimelineStore struct {
	storage *Store
}

func (c *CaseTimelineStore) Get(rpc options.Searcher) (*model.CaseTimeline, error) {
	filters := rpc.GetFilter("case_id")
	if len(filters) == 0 || filters[0].Operator != "=" {
		return nil, dberr.NewDBError("postgres.case_timeline.get.check_args.case_id", "case id required and must be '='")
	}

	caseID, err := strconv.ParseInt(filters[0].Value, 10, 64)
	if err != nil || caseID == 0 {
		return nil, dberr.NewDBError("postgres.case_timeline.get.check_args.case_id", "case id required")
	}

	fields := rpc.GetFields()
	if len(fields) == 0 {
		fields = CaseTimelineFields
	}

	// Build the SQL query with JSONB
	query, args, err := c.buildTimelineQuery(caseID, fields, rpc)
	if err != nil {
		return nil, dberr.NewDBError("postgres.case_timeline.get.build_query.error", err.Error())
	}

	// Get database connection
	db, err := c.storage.Database()
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.case_timeline.get.database.error", err)
	}

	var days []*model.DayTimeline
	err = pgxscan.Select(rpc, db, &days, query, args...)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.case_timeline.get.scan.error", err)
	}

	// Process each day to unmarshal items and event data
	for _, day := range days {
		// unmarshal the items JSONB field
		if err := day.UnmarshalItems(); err != nil {
			return nil, dberr.NewDBInternalError("postgres.case_timeline.get.unmarshal_items.error", err)
		}

		// unmarshal each event's data
		for _, event := range day.Items {
			if err := event.UnmarshalEventData(); err != nil {
				return nil, dberr.NewDBInternalError("postgres.case_timeline.get.unmarshal_event_data.error", err)
			}
		}
	}

	// apply pagination
	var result model.CaseTimeline
	result.Days, result.Next = storeutil.ResolvePaging(rpc.GetSize(), days)
	result.Page = int32(rpc.GetPage())

	return &result, nil
}

// buildTimelineQuery creates SQL query using JSONB for event data
func (c *CaseTimelineStore) buildTimelineQuery(caseID int64, fields []string, rpc options.Searcher) (string, []interface{}, error) {
	if caseID == 0 {
		return "", nil, fmt.Errorf("case id required")
	}

	var cteQueries []string
	var args []interface{}
	argIndex := 1

	// Map to track which event types are included
	includeType := map[string]bool{
		"call":  false,
		"chat":  false,
		"email": false,
	}

	// Build CTEs for each requested field
	for _, field := range fields {
		switch field {
		case "call":
			includeType["call"] = true
			cteQueries = append(cteQueries, fmt.Sprintf(CallsJSONBCTE, argIndex, argIndex+1))
			args = append(args, caseID, store.CommunicationCall)
			argIndex += 2
		case "chat":
			includeType["chat"] = true
			cteQueries = append(cteQueries, fmt.Sprintf(ChatsJSONBCTE, argIndex, argIndex+1))
			args = append(args, caseID, store.CommunicationChat)
			argIndex += 2
		case "email":
			includeType["email"] = true
			cteQueries = append(cteQueries, fmt.Sprintf(EmailsJSONBCTE, argIndex, argIndex+1))
			args = append(args, caseID, store.CommunicationEmail)
			argIndex += 2
		default:
			return "", nil, fmt.Errorf("unknown field: %s", field)
		}
	}

	// If no communication types are requested, return empty result query
	if len(cteQueries) == 0 {
		return "SELECT NULL::timestamp AS day_timestamp, 0::bigint AS chats_count, 0::bigint AS calls_count, 0::bigint AS emails_count, '[]'::jsonb AS items WHERE false", []interface{}{}, nil
	}

	// Build union parts for the main query
	var unionParts []string
	if includeType["call"] {
		unionParts = append(unionParts, `SELECT
			DATE_TRUNC('day', created_at) AS day,
			created_at::timestamp AS created_at,
			'call' AS event_type,
			jsonb_build_object(
				'id', id,
				'closed_at', (EXTRACT(EPOCH FROM closed_at) * 1000)::bigint,
				'duration', duration,
				'total_duration', total_duration,
				'is_inbound', is_inbound,
				'is_missed', is_missed,
				'is_detailed', is_detailed,
				'participants', participants,
				'gateway', gateway,
				'flow_scheme', flow_scheme,
				'queue', queue,
				'files', files,
				'transcripts', transcripts
			) AS event_data
			FROM call_data`)
	}
	if includeType["chat"] {
		unionParts = append(unionParts, `SELECT
			DATE_TRUNC('day', created_at) AS day,
			created_at::timestamp AS created_at,
			'chat' AS event_type,
			jsonb_build_object(
				'id', id,
				'closed_at', (EXTRACT(EPOCH FROM closed_at) * 1000)::bigint,
				'duration', duration,
				'is_inbound', is_inbound,
				'is_missed', is_missed,
				'is_detailed', is_detailed,
				'participants', participants,
				'gateway', gateway,
				'flow_scheme', flow_scheme,
				'queue', queue
			) AS event_data
			FROM chat_data`)
	}
	if includeType["email"] {
		unionParts = append(unionParts, `SELECT
			DATE_TRUNC('day', created_at) AS day,
			created_at::timestamp AS created_at,
			'email' AS event_type,
			jsonb_build_object(
				'id', id,
				'closed_at', (EXTRACT(EPOCH FROM closed_at) * 1000)::bigint,
				'duration', duration,
				'from', "from",
				'to', "to",
				'sender', sender,
				'cc', cc,
				'is_inbound', is_inbound,
				'subject', subject,
				'body', body,
				'html', html,
				'is_detailed', is_detailed,
				'profile', profile,
				'owner', owner,
				'attachments', attachments
			) AS event_data
			FROM email_data`)
	}

	page := rpc.GetPage()
	size := rpc.GetSize()

	// Build the main query
	mainQuery := "SELECT " +
		"    (EXTRACT(EPOCH FROM day) * 1000)::bigint AS day_timestamp," +
		"    COALESCE(SUM(CASE WHEN event_type = 'chat' THEN 1 ELSE 0 END), 0) AS chats_count," +
		"    COALESCE(SUM(CASE WHEN event_type = 'call' THEN 1 ELSE 0 END), 0) AS calls_count," +
		"    COALESCE(SUM(CASE WHEN event_type = 'email' THEN 1 ELSE 0 END), 0) AS emails_count," +
		"    COALESCE(jsonb_agg(" +
		"        jsonb_build_object(" +
		"            'type', event_type," +
		"            'created_at', (EXTRACT(EPOCH FROM created_at) * 1000)::bigint," +
		"            'event_data', event_data" +
		"        ) ORDER BY created_at DESC" +
		"    ), '[]'::jsonb) AS items " +
		"FROM (" + strings.Join(unionParts, " UNION ALL ") + ") combined " +
		"GROUP BY day " +
		"ORDER BY day DESC"

	// apply pagination
	if size > 0 {
		mainQuery += fmt.Sprintf(" LIMIT %d", size+1)
		if page > 1 {
			mainQuery += fmt.Sprintf(" OFFSET %d", (page-1)*size)
		}
	}

	// Combine with CTEs
	finalQuery := "WITH " + strings.Join(cteQueries, ",\n") + "\n" + mainQuery

	return finalQuery, args, nil
}

// endregion

// GetCounter retrieves timeline counter data for a case
func (c *CaseTimelineStore) GetCounter(rpc options.Searcher) ([]*model.TimelineCounter, error) {
	filters := rpc.GetFilter("case_id")
	if len(filters) == 0 || filters[0].Operator != "=" {
		return nil, dberr.NewDBError(
			"postgres.case_timeline.get_counter.check_args.case_id",
			"case id required and must use '=' operator")
	}

	caseID, err := strconv.ParseInt(filters[0].Value, 10, 64)
	if err != nil || caseID == 0 {
		return nil, dberr.NewDBError(
			"postgres.case_timeline.get_counter.check_args.case_id",
			"invalid case id")
	}

	db, err := c.storage.Database()
	if err != nil {
		return nil, dberr.NewDBInternalError(
			"postgres.case_timeline.get_counter.database_connection_error",
			err)
	}

	fields := rpc.GetFields()
	if len(fields) == 0 {
		fields = CaseTimelineFields
	}

	var unionParts []string
	var args []interface{}
	argIndex := 1

	// Build union parts
	for _, field := range fields {
		switch field {
		case "call":
			unionParts = append(unionParts, fmt.Sprintf(CallCounterQuery, argIndex, argIndex+1))
			args = append(args, caseID, store.CommunicationCall)
			argIndex += 2

		case "chat":
			unionParts = append(unionParts, fmt.Sprintf(ChatCounterQuery, argIndex, argIndex+1))
			args = append(args, caseID, store.CommunicationChat)
			argIndex += 2

		case "email":
			unionParts = append(unionParts, fmt.Sprintf(EmailCounterQuery, argIndex, argIndex+1))
			args = append(args, caseID, store.CommunicationEmail)
			argIndex += 2
		}
	}

	// If no communication types are requested, return empty result
	if len(unionParts) == 0 {
		return []*model.TimelineCounter{}, nil
	}

	// Build the query
	finalQuery := strings.Join(unionParts, " UNION ALL ")

	// Execute query using pgxscan
	var counters []*model.TimelineCounter
	err = pgxscan.Select(rpc, db, &counters, finalQuery, args...)
	if err != nil {
		return nil, dberr.NewDBInternalError(
			"postgres.case_timeline.get_counter.scan_error",
			err)
	}

	return counters, nil
}

// endregion

// GetItemInfo retrieves saved variables + postprocessing results for a single
// timeline communication (call | chat | email) that belongs to the case.
func (c *CaseTimelineStore) GetItemInfo(rpc options.Searcher, caseID int64, itemType model.CaseTimelineEventType, itemID string) (*model.CaseTimelineItemInfo, error) {
	var raw string
	switch itemType {
	case model.TimelineEventTypeCall:
		raw = CallItemInfoQuery
	case model.TimelineEventTypeChat:
		raw = ChatItemInfoQuery
	case model.TimelineEventTypeEmail:
		raw = EmailItemInfoQuery
	default:
		return nil, dberr.NewDBBadRequestError(
			"postgres.case_timeline.get_item_info.check_args.type",
			fmt.Sprintf("unknown timeline event type %q", itemType))
	}

	channel, err := communicationChannel(itemType)
	if err != nil {
		return nil, err
	}

	db, err := c.storage.Database()
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.case_timeline.get_item_info.database.error", err)
	}

	var row caseTimelineItemInfoRow
	err = pgxscan.Get(rpc, db, &row, raw, itemID, rpc.GetAuthOpts().GetDomainId(), caseID, channel)
	if err != nil {
		if dberr.Is(err, pgx.ErrNoRows) {
			return &model.CaseTimelineItemInfo{}, nil
		}
		return nil, ParseError(err)
	}

	info := &model.CaseTimelineItemInfo{}
	var exclude map[string]bool
	switch itemType {
	case model.TimelineEventTypeEmail:
		exclude = emailSystemVariableKeys
	case model.TimelineEventTypeChat:
		exclude = chatSystemVariableKeys
	}
	if info.Variables, err = decodeCaseTimelineVariables(row.Variables, exclude); err != nil {
		return nil, dberr.NewDBInternalError("postgres.case_timeline.get_item_info.decode_variables.error", err)
	}
	if info.Postprocessing, err = decodeCaseTimelinePostprocessing(row.Postprocessing); err != nil {
		return nil, dberr.NewDBInternalError("postgres.case_timeline.get_item_info.decode_postprocessing.error", err)
	}
	return info, nil
}

func communicationChannel(itemType model.CaseTimelineEventType) (string, error) {
	switch itemType {
	case model.TimelineEventTypeCall:
		return store.CommunicationCall, nil
	case model.TimelineEventTypeChat:
		return store.CommunicationChat, nil
	case model.TimelineEventTypeEmail:
		return store.CommunicationEmail, nil
	default:
		return "", dberr.NewDBBadRequestError(
			"postgres.case_timeline.get_item_info.check_args.type",
			fmt.Sprintf("unknown timeline event type %q", itemType))
	}
}

type caseTimelineItemInfoRow struct {
	Variables      []byte `db:"variables"`
	Postprocessing []byte `db:"postprocessing"`
}

var emailSystemVariableKeys = map[string]bool{
	"message_id":  true,
	"reply_to":    true,
	"from":        true,
	"cc":          true,
	"sender":      true,
	"in_reply_to": true,
	"body":        true,
	"body_html":   true,
	"subject":     true,
	"id":          true,
	"attachments": true,
}

// chatSystemVariableKeys mirrors chat_manager's systemChannelVars
// (internal/event_router/event_router.go)
var chatSystemVariableKeys = map[string]bool{
	"chat":             true,
	"cid":              true,
	"externalChatID":   true,
	"flow":             true,
	"old_flow":         true,
	"xfer":             true,
	"chat_transferred": true,
	"chatplan_name":    true,
	"from":             true,
	"user":             true,
	"needs_processing": true,
}

func decodeCaseTimelineVariables(raw []byte, exclude map[string]bool) ([]*model.CaseTimelineVariable, error) {
	if len(raw) == 0 {
		return nil, nil
	}
	var kv map[string]json.RawMessage
	if err := json.Unmarshal(raw, &kv); err != nil {
		return nil, err
	}
	if len(kv) == 0 {
		return nil, nil
	}
	keys := make([]string, 0, len(kv))
	for key := range kv {
		if exclude[key] {
			continue
		}
		keys = append(keys, key)
	}
	sort.Strings(keys)

	vars := make([]*model.CaseTimelineVariable, 0, len(keys))
	for _, key := range keys {
		vars = append(vars, &model.CaseTimelineVariable{
			Key:   key,
			Value: jsonValueToString(kv[key]),
		})
	}
	return vars, nil
}

func jsonValueToString(raw json.RawMessage) string {
	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		return s
	}
	return string(raw)
}

type caseTimelinePostprocessingRow struct {
	Agent       *model.GeneralLookup `json:"agent"`
	Form        json.RawMessage      `json:"form"`
	ReportingAt int64                `json:"reporting_at"`
}

func decodeCaseTimelinePostprocessing(raw []byte) ([]*model.CaseTimelinePostprocessingResult, error) {
	if len(raw) == 0 || string(raw) == "null" {
		return nil, nil
	}
	var rows []caseTimelinePostprocessingRow
	if err := json.Unmarshal(raw, &rows); err != nil {
		return nil, err
	}

	results := make([]*model.CaseTimelinePostprocessingResult, 0, len(rows))
	for _, r := range rows {
		item := &model.CaseTimelinePostprocessingResult{
			Agent:       r.Agent,
			ReportingAt: r.ReportingAt,
		}
		if len(r.Form) > 0 && string(r.Form) != "null" {
			var native any
			if err := json.Unmarshal(r.Form, &native); err != nil {
				return nil, err
			}
			item.Form = native
		}
		results = append(results, item)
	}
	return results, nil
}

const (
	// CallItemInfoQuery returns saved variables + postprocessing results for a
	// single call, scoped to a case via cases.case_communication.
	CallItemInfoQuery = `
SELECT
	coalesce(c.payload, '{}'::jsonb) AS variables,
	(
		WITH RECURSIVE legs AS (
			SELECT d.id, d.attempt_id, d.attempt_ids
			FROM call_center.cc_calls_history d
			WHERE d.id = c.id
			  AND d.domain_id = $2
			UNION ALL
			SELECT d.id, d.attempt_id, d.attempt_ids
			FROM call_center.cc_calls_history d, legs
			WHERE d.parent_id = legs.id
			   OR d.transfer_from = legs.id
		),
		attempts AS (
			SELECT attempt_id AS id FROM legs WHERE attempt_id NOTNULL
			UNION
			SELECT unnest(attempt_ids) FROM legs WHERE attempt_ids NOTNULL
		)
		SELECT jsonb_agg(jsonb_build_object(
			'agent', call_center.cc_get_lookup(u.id, coalesce(u.name, u.username)),
			'form', p.form_fields,
			'reporting_at', call_center.cc_view_timestamp(p.reporting_at)
		))
		FROM call_center.cc_member_attempt_history p
			LEFT JOIN call_center.cc_agent a ON a.id = p.agent_id
			LEFT JOIN directory.wbt_user u ON u.id = a.user_id
		WHERE p.domain_id = $2
		  AND p.id = ANY(SELECT id FROM attempts)
		  AND p.form_fields NOTNULL
	) AS postprocessing
FROM call_center.cc_calls_history c
WHERE c.id = $1::uuid
  AND c.domain_id = $2
  AND c.id::text = ANY(
	SELECT communication_id FROM cases.case_communication casecom
	LEFT JOIN call_center.cc_communication com ON com.id = casecom.communication_type
	WHERE casecom.case_id = $3 AND com.channel = $4
  )`

	// ChatItemInfoQuery returns saved variables + postprocessing results for a
	// single chat conversation, scoped to a case via cases.case_communication.
	ChatItemInfoQuery = `
SELECT
	coalesce(conv.props, '{}'::jsonb) AS variables,
	(
		SELECT jsonb_agg(jsonb_build_object(
			'agent', call_center.cc_get_lookup(u.id, coalesce(u.name, u.username)),
			'form', p.form_fields,
			'reporting_at', call_center.cc_view_timestamp(p.reporting_at)
		))
		FROM call_center.cc_member_attempt_history p
			LEFT JOIN call_center.cc_agent a ON a.id = p.agent_id
			LEFT JOIN directory.wbt_user u ON u.id = a.user_id
		WHERE p.domain_id = $2
		  AND p.member_call_id = conv.id::varchar
		  AND p.form_fields NOTNULL
	) AS postprocessing
FROM chat.conversation conv
WHERE conv.id = $1::uuid
  AND conv.domain_id = $2
  AND conv.id::text = ANY(
	SELECT communication_id FROM cases.case_communication casecom
	LEFT JOIN call_center.cc_communication com ON com.id = casecom.communication_type
	WHERE casecom.case_id = $3 AND com.channel = $4
  )`

	// EmailItemInfoQuery returns saved variables + postprocessing results for a
	// single email, scoped to a case via cases.case_communication.
	EmailItemInfoQuery = `
SELECT
	coalesce((
		SELECT m.variables
		FROM call_center.cc_member m
		WHERE m.domain_id = $2
		  AND e.message_id IS NOT NULL
		  AND m.variables ->> 'message_id' = e.message_id
		ORDER BY m.id DESC
		LIMIT 1
	), '{}'::jsonb) AS variables,
	(
		SELECT jsonb_agg(jsonb_build_object(
			'agent', call_center.cc_get_lookup(u.id, coalesce(u.name, u.username)),
			'form', p.form_fields,
			'reporting_at', call_center.cc_view_timestamp(p.reporting_at)
		))
		FROM call_center.cc_member_attempt_history p
			JOIN call_center.cc_member m ON m.id = p.member_id
			LEFT JOIN call_center.cc_agent a ON a.id = p.agent_id
			LEFT JOIN directory.wbt_user u ON u.id = a.user_id
		WHERE p.domain_id = $2
		  AND e.message_id IS NOT NULL
		  AND m.variables ->> 'message_id' = e.message_id
		  AND p.form_fields NOTNULL
	) AS postprocessing
FROM call_center.cc_email e
WHERE e.id = $1::int8
  AND e.id::text = ANY(
	SELECT communication_id FROM cases.case_communication casecom
	LEFT JOIN call_center.cc_communication com ON com.id = casecom.communication_type
	WHERE casecom.case_id = $3 AND com.channel = $4
  )`
)

const (
	CallCounterQuery = `
SELECT
	'Phone' AS event_type,
	COUNT(*)::bigint AS count,
	COALESCE((EXTRACT(EPOCH FROM MIN(c.created_at)) * 1000)::bigint, 0) AS date_from,
	COALESCE((EXTRACT(EPOCH FROM MAX(c.hangup_at)) * 1000)::bigint, 0) AS date_to
FROM call_center.cc_calls_history c
WHERE c.id = ANY(SELECT communication_id::uuid
	FROM cases.case_communication casecom
	LEFT JOIN call_center.cc_communication com ON com.id = casecom.communication_type
	WHERE case_id = $%d AND com.channel = $%d)
AND c.transfer_from IS NULL`

	ChatCounterQuery = `
SELECT
	'Messaging' AS event_type,
	COUNT(*)::bigint AS count,
	COALESCE((EXTRACT(EPOCH FROM MIN(conv.created_at)) * 1000)::bigint, 0) AS date_from,
	COALESCE((EXTRACT(EPOCH FROM MAX(conv.closed_at)) * 1000)::bigint, 0) AS date_to
FROM chat.conversation conv
WHERE conv.id = ANY(SELECT communication_id::uuid
	FROM cases.case_communication casecom
	LEFT JOIN call_center.cc_communication com ON com.id = casecom.communication_type
	WHERE case_id = $%d AND com.channel = $%d)`

	EmailCounterQuery = `
SELECT
	'Email' AS event_type,
	COUNT(*)::bigint AS count,
	COALESCE((EXTRACT(EPOCH FROM MIN(e.created_at)) * 1000)::bigint, 0) AS date_from,
	COALESCE((EXTRACT(EPOCH FROM MAX(e.created_at)) * 1000)::bigint, 0) AS date_to
FROM call_center.cc_email e
WHERE e.id = ANY(SELECT communication_id::bigint
	FROM cases.case_communication casecom
	LEFT JOIN call_center.cc_communication com ON com.id = casecom.communication_type
	WHERE case_id = $%d AND com.channel = $%d)`

	// JSONB CTE Queries
	CallsJSONBCTE = `
call_data AS (
	SELECT
		c.id::text,
		c.created_at,
		c.hangup_at AS closed_at,
		round(case when c.user_id notnull then date_part('epoch'::text, c.hangup_at - c.created_at)::bigint else (select date_part('epoch'::text, hangup_at - created_at)::bigint from call_center.cc_calls_history where parent_id = c.id limit 1) end)::bigint as duration,
		root.duration AS total_duration,
		(c.direction = 'inbound') AS is_inbound,
		(case when c.bridged_id isnull then c.queue_id notnull else false end) AS is_missed,
		exists(with recursive a as (select*
                from call_center.cc_calls_history
                where id in (with recursive a as (select d.id::uuid, d.user_id
                                                  from call_center.cc_calls_history d
                                                  where d.id::uuid = c.id
                                                  union all
                                                  select d.id::uuid, d.user_id
                                                  from call_center.cc_calls_history d,
                                                       a
                                                  where (d.parent_id::uuid =
                                                         a.id::uuid or
                                                         (d.transfer_from::uuid = a.id::uuid)))
                             select distinct id
                             from a))
           select id::uuid ids
           from a
           where a.parent_id != c.id
              or coalesce(a.transfer_to::varchar, a.transfer_from::varchar,
                          a.blind_transfer::varchar) notnull) AS is_detailed,
		participants.data AS participants,
		gateway.data AS gateway,
		flow_scheme.data AS flow_scheme,
		queue.data AS queue,
		files.data AS files,
		transcripts.data AS transcripts
	FROM call_center.cc_calls_history c
	LEFT JOIN LATERAL (SELECT round(date_part('epoch'::text, c.hangup_at - c.created_at)::bigint) duration) root ON true
	LEFT JOIN LATERAL (
		SELECT jsonb_agg(jsonb_build_object('id', users.id, 'name', users.name)) AS data
		FROM (with recursive a as (select *
                                from call_center.cc_calls_history
                                where id in (with recursive a as (select d.id::uuid, d.user_id
                                                                  from call_center.cc_calls_history d
                                                                  where d.id::uuid = c.id
                                                                  union all
                                                                  select d.id::uuid, d.user_id
                                                                  from call_center.cc_calls_history d,
                                                                       a
                                                                  where (d.parent_id::uuid = a.id::uuid or
                                                                         (d.transfer_from::uuid = a.id::uuid)))
                             select distinct id
                             from a))
              SELECT usr.id, coalesce(usr.name, usr.username) as name
              from a
                       inner join directory.wbt_user usr on a.user_id = usr.id
              order by a.created_at) users
	) participants ON true
	LEFT JOIN LATERAL (
		SELECT jsonb_build_object('id', null, 'name', null) AS data
	) gateway ON true
	LEFT JOIN LATERAL (
		SELECT jsonb_build_object('id', scheme.id, 'name', scheme.name) AS data
		FROM flow.acr_routing_scheme scheme
		WHERE scheme.id = c.schema_ids[array_length(c.schema_ids, 1)]
	) flow_scheme ON true
	LEFT JOIN LATERAL (
		SELECT jsonb_build_object('id', c.queue_id, 'name', a.name) AS data
		FROM call_center.cc_queue a
		WHERE a.id = c.queue_id
	) queue ON true
	LEFT JOIN LATERAL (
		SELECT jsonb_agg(jsonb_build_object('id', f1.id, 'size', f1.size, 'mime_type', f1.mime_type, 'name', f1.name, 'start_at', f1.created_at * 1000, 'channel', f1.channel)) AS data
		FROM storage.files f1
		WHERE f1.domain_id = c.domain_id
		  AND NOT f1.removed IS TRUE
		  AND f1.uuid = c.id::varchar
	) files ON true
	LEFT JOIN LATERAL (
		SELECT jsonb_agg(jsonb_build_object('id', tr.id, 'locale', tr.locale, 'file', jsonb_build_object('id', ff.id, 'name', ff.name))) AS data
		FROM storage.file_transcript tr
		LEFT JOIN storage.files ff ON ff.id = tr.file_id
		WHERE tr.uuid::text = c.id::text
	) transcripts ON true
	WHERE c.id = ANY(SELECT communication_id::uuid FROM cases.case_communication casecom LEFT JOIN call_center.cc_communication com ON com.id = casecom.communication_type WHERE case_id = $%d AND com.channel = $%d)
	  AND c.transfer_from isnull
)`

	ChatsJSONBCTE = `
chat_data AS (
	SELECT
		conv.id::text,
		conv.created_at,
		conv.closed_at AS closed_at,
		round(extract(EPOCH FROM (conv.closed_at - conv.created_at)))::bigint AS duration,
		true AS is_inbound,
		false AS is_missed,
		true AS is_detailed,
		participants.data AS participants,
		gateway.data AS gateway,
		flow_scheme.data AS flow_scheme,
		queue.data AS queue
	FROM chat.conversation conv
	LEFT JOIN LATERAL (
		SELECT jsonb_agg(jsonb_build_object('id', usr.id, 'name', usr.name)) AS data
		FROM chat.channel c
		INNER JOIN directory.wbt_auth usr ON user_id = usr.id
		WHERE conversation_id = conv.id
		  AND internal
		  AND joined_at NOTNULL
		GROUP BY usr.id
	) participants ON true
	LEFT JOIN LATERAL (
		SELECT jsonb_build_object('id', b.id, 'name', b.name, 'provider', b.provider) AS data
		FROM chat.channel
		LEFT JOIN chat.bot b ON connection::bigint = b.id
		WHERE NOT internal
		  AND conversation_id = conv.id
	) gateway ON true
	LEFT JOIN LATERAL (
		SELECT jsonb_build_object('id', flow_scheme.id, 'name', flow_scheme.name) AS data
		FROM flow.acr_routing_scheme flow_scheme
		WHERE flow_scheme.id = (conv.props ->> 'flow')::bigint
	) flow_scheme ON true
	LEFT JOIN LATERAL (
		SELECT null::jsonb AS data
	) queue ON true
	WHERE conv.id = ANY(SELECT communication_id::uuid FROM cases.case_communication casecom LEFT JOIN call_center.cc_communication com ON com.id = casecom.communication_type WHERE case_id = $%d AND com.channel = $%d)
)`

	EmailsJSONBCTE = `
email_data AS (
	SELECT
		e.id::text,
		e.created_at,
		e.created_at AS closed_at,
		0::bigint AS duration,
		e."from",
		e."to",
		e.sender,
		e.cc,
		(e.direction = 'inbound') AS is_inbound,
		e.subject,
		e.body,
		e.html,
		true AS is_detailed,
		profile.data AS profile,
		owner.data AS owner,
		attachments.data AS attachments
	FROM call_center.cc_email e
	LEFT JOIN LATERAL (
		SELECT jsonb_build_object('id', e.profile_id, 'name', p.name) AS data
		FROM call_center.cc_email_profile p
		WHERE e.profile_id = p.id
	) profile ON true
	LEFT JOIN LATERAL (
		SELECT jsonb_build_object('id', e."owner_id", 'name', u."name") AS data
		FROM directory.wbt_user u
		WHERE e.owner_id = u.id
	) owner ON true
	LEFT JOIN LATERAL (
		SELECT jsonb_agg(jsonb_build_object('id', f.id, 'mime', f.mime_type, 'name', f.view_name, 'size', f.size)) AS data
		FROM storage.files f
		WHERE f.id = any (e.attachment_ids)
	) attachments ON true
	WHERE e.id = ANY(SELECT communication_id::bigint FROM cases.case_communication casecom LEFT JOIN call_center.cc_communication com ON com.id = casecom.communication_type WHERE case_id = $%d AND com.channel = $%d)
)`
)

func NewCaseTimelineStore(store *Store) (store.CaseTimelineStore, error) {
	if store == nil {
		return nil, dberr.NewDBError("postgres.case_timeline.new_case_timeline_store.check_args.store",
			"store required")
	}
	return &CaseTimelineStore{storage: store}, nil
}
