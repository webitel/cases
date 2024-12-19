package postgres

import (
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgtype"
	"github.com/webitel/cases/api/cases"
	dberr "github.com/webitel/cases/internal/error"
	"github.com/webitel/cases/internal/store"
	"github.com/webitel/cases/internal/store/scanner"
	"github.com/webitel/cases/model"
)

var CaseTimelineFields = []string{
	"calls", "chats", "emails",
}

type CaseTimelineStore struct {
	storage store.Store
}

func (c *CaseTimelineStore) Get(rpc *model.SearchOptions) (*cases.GetTimelineResponse, error) {
	query, scanPlan, dbErr := buildCaseTimelineSqlizer(rpc)
	if dbErr != nil {
		return nil, dbErr
	}
	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}
	db, dbErr := c.storage.Database()
	if dbErr != nil {
		return nil, dbErr
	}

	rows, err := db.Query(rpc, store.CompactSQL(sql), args...)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.case_timeline.get.exec.error", err)
	}
	result := &cases.GetTimelineResponse{}
	for rows.Next() {
		node := &cases.DayTimeline{}
		var scanValue []any
		for _, f := range scanPlan {
			scanValue = append(scanValue, f(node))
		}
		err = rows.Scan(scanValue...)
		if err != nil {
			return nil, dberr.NewDBInternalError("postgres.case_timeline.get.exec.error", err)
		}
		result.Days = append(result.Days, node)
	}

	return result, nil
}

// region Timeline Build Functions
func buildCaseTimelineSqlizer(rpc *model.SearchOptions) (query squirrel.Sqlizer, plan []func(timeline *cases.DayTimeline) any, dbError *dberr.DBError) {
	if rpc == nil {
		return nil, nil, dberr.NewDBError("postgres.case_timeline.build_case_timeline_sqlizer.check_args.rpc", "search options required")
	}
	if len(rpc.IDs) == 0 {
		return nil, nil, dberr.NewDBError("postgres.case_timeline.build_case_timeline_sqlizer.check_args.case_id", "case id required")
	}
	caseId := rpc.IDs[0]
	if caseId <= 0 {
		return nil, nil, dberr.NewDBError("postgres.case_timeline.build_case_timeline_sqlizer.check_args.case_id", "case id empty")
	}
	fields := rpc.Fields[:]
	if len(fields) == 0 {
		fields = CaseTimelineFields
	}
	var (
		ctes       = make(map[string]squirrel.Sqlizer)
		chatsPlan  []func(timeline *cases.Event) any
		callsPlan  []func(timeline *cases.Event) any
		emailsPlan []func(timeline *cases.Event) any
		dayEvents  []string
		from       = "( SELECT * FROM ("
	)
	for i, field := range fields {
		var (
			subplan []func(*cases.Event) any
			err     *dberr.DBError
			cte     squirrel.Sqlizer
		)
		switch field {
		case "chats":
			cte, subplan, err = buildTimelineChatsColumn(caseId)
			if err != nil {
				return nil, nil, err
			}
		case "calls":
			cte, subplan, err = buildTimelineCallsColumn(caseId)
			if err != nil {
				return nil, nil, err
			}
		case "emails":
			cte, subplan, err = buildTimelineEmailsColumn(caseId)
			if err != nil {
				return nil, nil, err
			}
		default:
			return nil, nil, dberr.NewDBError("postgres.case_timeline.build_case_timeline_sqlizer.parse_fields.unknown", "unknown field "+field)
		}
		ctes[field] = cte
		callsPlan = subplan
		if i != 0 {
			from += " union all "
		}
		from += fmt.Sprintf(`SELECT
					  DATE_TRUNC('day', created_at) "day",
					  created_at,
					  'chat' "event",
					  (%s) "data"
					   from %[1]s
						`, field)
	}
	from += ") AS ec ORDER BY created_at DESC) e"
	cteQuery, args, err := store.FormAsCTEs(ctes)
	if err != nil {
		return nil, nil, dberr.NewDBError("postgres.case_timeline.build_case_timeline_sqlizer.form_cte.error", err.Error())
	}
	query = squirrel.Select("day", "array_agg(e.event) \"events\"", "array_agg(e.data) \"data\"").
		GroupBy("e.day").
		OrderBy("e.day desc").From(from).
		Prefix(cteQuery, args...).
		PlaceholderFormat(squirrel.Dollar)

	plan = []func(timeline *cases.DayTimeline) any{
		func(rec *cases.DayTimeline) any {
			return scanner.ScanTimestamp(&rec.DayTimestamp)
		},
		// day event signatures
		func(rec *cases.DayTimeline) any {
			return scanner.ScanFunc(func(src interface{}) (err error) {
				if src == nil {
					return
				}
				var text string
				switch src := src.(type) {
				case string:
					text = src
				case []byte:
					text = string(src)
				default:
					return dberr.NewDBError(
						"postgres.case_timeline.build_case_timeline_sqlizer.scan_event_signatures.error",
						"postgres: unknown type to convert into []string",
					)
				}
				if src == "{}" {
					return // nil; NULL
				}

				var rows *pgtype.UntypedTextArray
				rows, err = pgtype.ParseUntypedTextArray(text)
				if err != nil || len(rows.Elements) == 0 {
					return err
				}

				// context:
				// [row] -- fields requested
				var (
					size = len(rows.Elements)
				)

				if size == 0 {
					return nil
				}
				// init dayEvents for every row
				dayEvents = make([]string, 0, size)

				// DECODE
				for _, elem := range rows.Elements {
					// RECORD
					if elem == "" {
						return dberr.NewDBError("postgres.case_timeline.build_case_timeline_sqlizer.scan_event_signatures.element_error", "empty event type")
					}
					dayEvents = append(dayEvents, elem)

				}
				return nil
			})
		},
		// actually day events
		func(rec *cases.DayTimeline) any {
			return scanner.ScanFunc(func(src any) (err error) {
				if src == nil {
					return
				}
				var text string
				switch src := src.(type) {
				case string:
					text = src
				case []byte:
					text = string(src)
				default:
					return dberr.NewDBError(
						"postgres.case_timeline.build_case_timeline_sqlizer.scan_events.error",
						"unknown type input",
					)
				}
				if src == "{}" {
					return // nil; NULL
				}

				var rows *pgtype.UntypedTextArray
				rows, err = pgtype.ParseUntypedTextArray(text)
				if err != nil || len(rows.Elements) == 0 {
					return err
				}

				var (
					node *cases.Event
				)

				// DECODE
				for r, elem := range rows.Elements {
					// RECORD
					node = &cases.Event{} // NEW
					var (
						eventType string
						scanPlan  []func(*cases.Event) any
					)
					// ALLOC
					switch eventType = dayEvents[r]; eventType {
					case cases.CaseTimelineEventType_email.String():
						actualEvent := &cases.EmailEvent{}
						node.Type = cases.CaseTimelineEventType_email
						node.Event = &cases.Event_Email{Email: actualEvent}
						rec.EmailsCount++

						// scanning functions set
						scanPlan = emailsPlan
					case cases.CaseTimelineEventType_chat.String():
						actualEvent := &cases.ChatEvent{}
						node.Type = cases.CaseTimelineEventType_chat
						node.Event = &cases.Event_Chat{Chat: actualEvent}
						rec.ChatsCount++

						// scanning functions set
						scanPlan = chatsPlan
					case cases.CaseTimelineEventType_call.String():
						actualEvent := &cases.CallEvent{}
						node.Type = cases.CaseTimelineEventType_call
						node.Event = &cases.Event_Call{Call: actualEvent}
						rec.CallsCount++

						// scanning functions set
						scanPlan = callsPlan
					}
					// DECODE
					scan := pgtype.NewCompositeTextScanner(pgtype.NewConnInfo(), []byte(elem))
					for _, bind := range scanPlan {
						df := bind(node)
						if df == nil {
							// omit; pseudo calc
							continue
						}
						scan.ScanValue(df)
						err = scan.Err()
						if err != nil {
							return err
						}
					}
					rec.Items = append(rec.Items, node)
				}
				return nil
			})
		},
	}
	return query, plan, nil

}
func buildTimelineChatsColumn(caseId int64) (base squirrel.Sqlizer, plan []func(timeline *cases.Event) any, dbError *dberr.DBError) {
	if caseId == 0 {
		return nil, nil, dberr.NewDBError("postgres.case_timeline.build_timeline_chats_column.check_args.case_id.empty", "case id required")
	}
	base = squirrel.Expr(ChatsCTE, caseId)
	plan = append(plan,
		func(node *cases.Event) any {
			chat := node.GetChat()
			return scanner.ScanText(&chat.Id)
		},
		func(node *cases.Event) any {
			return scanner.ScanTimestamp(&node.CreatedAt)
		},
		func(node *cases.Event) any {
			chat := node.GetChat()
			return scanner.ScanTimestamp(&chat.ClosedAt)
		},
		func(node *cases.Event) any {
			chat := node.GetChat()
			return &chat.Duration
		},
		func(node *cases.Event) any {
			chat := node.GetChat()
			return scanner.ScanLookupList(&chat.Participants)
		},
		func(node *cases.Event) any {
			chat := node.GetChat()
			return scanner.ScanRowExtendedLookup(&chat.Gateway)
		},
		func(node *cases.Event) any {
			return scanner.ScanTimestamp(&node.CreatedAt)
		},
		func(node *cases.Event) any {
			chat := node.GetChat()
			return scanner.ScanRowLookup(&chat.FlowScheme)
		},
		func(node *cases.Event) any {
			// TODO: type
			return nil
		},
		func(node *cases.Event) any {
			chat := node.GetChat()
			return &chat.IsInbound
		},
		func(node *cases.Event) any {
			chat := node.GetChat()
			return &chat.IsMissed
		},
		func(node *cases.Event) any {
			chat := node.GetChat()
			return scanner.ScanRowLookup(&chat.Queue)
		},
		func(node *cases.Event) any {
			chat := node.GetChat()
			return &chat.IsDetailed
		})
	return
}
func buildTimelineCallsColumn(caseId int64) (base squirrel.Sqlizer, plan []func(timeline *cases.Event) any, dbError *dberr.DBError) {
	if caseId == 0 {
		return nil, nil, dberr.NewDBError("postgres.case_timeline.build_timeline_calls_column.check_args.case_id.empty", "case id required")
	}
	base = squirrel.Expr(CallsCTE, caseId)

	plan = append(plan,
		func(node *cases.Event) any {
			call := node.GetCall()
			return scanner.ScanText(&call.Id)
		},
		func(node *cases.Event) any {
			return scanner.ScanTimestamp(&node.CreatedAt)
		},
		func(node *cases.Event) any {
			call := node.GetCall()
			return scanner.ScanTimestamp(&call.ClosedAt)
		},
		func(node *cases.Event) any {
			call := node.GetCall()
			return &call.Duration
		},
		func(node *cases.Event) any {
			call := node.GetCall()
			return &call.TotalDuration
		},
		func(node *cases.Event) any {
			call := node.GetCall()
			return scanner.ScanLookupList(&call.Participants)
		},
		func(node *cases.Event) any {
			call := node.GetCall()
			return scanner.ScanRowLookup(&call.Gateway)
		},
		func(node *cases.Event) any {
			call := node.GetCall()
			return scanner.ScanRowLookup(&call.FlowScheme)
		},
		// type
		func(node *cases.Event) any {
			return nil
		},
		func(node *cases.Event) any {
			call := node.GetCall()
			return &call.IsInbound
		},
		func(node *cases.Event) any {
			call := node.GetCall()
			return &call.IsMissed
		},
		func(node *cases.Event) any {
			call := node.GetCall()
			return scanner.ScanRowLookup(&call.Queue)
		},
		func(node *cases.Event) any {
			call := node.GetCall()
			return &call.IsDetailed
		},
		// files
		func(node *cases.Event) any {
			// TODO
			return nil
		},
		// transcripts
		func(node *cases.Event) any {
			// TODO
			return nil
		},
	)
	return
}
func buildTimelineEmailsColumn(caseId int64) (base squirrel.Sqlizer, plan []func(timeline *cases.Event) any, dbError *dberr.DBError) {
	if caseId == 0 {
		return nil, nil, dberr.NewDBError("postgres.case_timeline.build_timeline_emails_column.check_args.case_id.empty", "case id required")
	}
	base = squirrel.Expr(EmailsCTE, caseId)

	plan = append(plan,
		func(node *cases.Event) any {
			email := node.GetEmail()
			return scanner.ScanText(&email.Id)
		},
		func(node *cases.Event) any {
			email := node.GetEmail()
			return &email.From
		},
		func(node *cases.Event) any {
			email := node.GetEmail()
			return &email.To
		},
		func(node *cases.Event) any {
			email := node.GetEmail()
			return scanner.ScanRowLookup(&email.Profile)
		},
		func(node *cases.Event) any {
			email := node.GetEmail()
			return scanner.ScanText(&email.Subject)
		},
		func(node *cases.Event) any {
			email := node.GetEmail()
			return &email.Cc
		},
		func(node *cases.Event) any {
			return scanner.ScanTimestamp(&node.CreatedAt)
		},
		func(node *cases.Event) any {
			email := node.GetEmail()
			return &email.IsInbound
		},
		func(node *cases.Event) any {
			email := node.GetEmail()
			return &email.Sender
		},
		func(node *cases.Event) any {
			email := node.GetEmail()
			return scanner.ScanText(&email.Body)
		},
		func(node *cases.Event) any {
			email := node.GetEmail()
			return scanner.ScanText(&email.Html)
		},
		// attachments
		func(node *cases.Event) any {
			// TODO
			return nil
		},
		func(node *cases.Event) any {
			email := node.GetEmail()
			return scanner.ScanRowLookup(&email.Owner)
		})

	return

}

// endregion

func (c *CaseTimelineStore) GetCounter(rpc *model.SearchOptions) ([]*model.TimelineCounter, error) {
	query, plan, dbErr := buildTimelineCounterSqlizer(rpc)
	if dbErr != nil {
		return nil, dbErr
	}
	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}
	db, dbErr := c.storage.Database()
	if dbErr != nil {
		return nil, dbErr
	}
	rows, err := db.Query(rpc, store.CompactSQL(sql), args...)
	if err != nil {
		return nil, err
	}
	var res []*model.TimelineCounter
	for rows.Next() {
		node := &model.TimelineCounter{}
		var planValues []any
		for _, bind := range plan {
			planValues = append(planValues, bind(node))
		}
		err = rows.Scan(planValues...)
		if err != nil {
			return nil, err
		}
		res = append(res, node)
	}
	return res, nil
}

// region Timeline Counter Build Functions

func buildTimelineCounterSqlizer(rpc *model.SearchOptions) (query squirrel.Sqlizer, scanPlan []func(response *model.TimelineCounter) any, dbError *dberr.DBError) {
	if rpc == nil {
		return nil, nil, dberr.NewDBError("postgres.case_timeline.build_case_timeline_sqlizer.check_args.rpc", "search options required")
	}
	if len(rpc.IDs) == 0 {
		return nil, nil, dberr.NewDBError("postgres.case_timeline.build_case_timeline_sqlizer.check_args.case_id", "case id required")
	}
	caseId := rpc.IDs[0]
	if caseId <= 0 {
		return nil, nil, dberr.NewDBError("postgres.case_timeline.build_case_timeline_sqlizer.check_args.case_id", "case id empty")
	}
	fields := rpc.Fields[:]
	if len(fields) == 0 {
		fields = CaseTimelineFields
	}
	var (
		ctes = make(map[string]squirrel.Sqlizer)
		from = "("
	)
	for i, field := range fields {
		switch field {
		case "calls":
			ctes[field] = squirrel.Expr(CallsCounterCTE, caseId)
		case "emails":
			ctes[field] = squirrel.Expr(EmailsCounterCTE, caseId)
		case "chats":
			ctes[field] = squirrel.Expr(ChatsCounterCTE, caseId)
		default:
			return nil, nil, dberr.NewDBError("postgres.case_timeline.build_case_timeline_counter_sqlizer.parse_fields.unknown", "unknown field "+field)
		}
		if i != 0 {
			from += " union all "
		}
		from += fmt.Sprintf("(select * from %s)", field)
	}
	from += ") s"
	cteQuery, args, err := store.FormAsCTEs(ctes)
	if err != nil {
		return nil, nil, dberr.NewDBError("postgres.case_timeline.build_case_timeline_counter_sqlizer.parse_fields.unknown", err.Error())
	}

	query = squirrel.Select("count(*) count",
		"type event",
		"max(closed_at) date_to",
		"min(created_at) date_from)").
		From(from).
		Prefix(cteQuery, args...)

	// each row represents type of event
	scanPlan = append(scanPlan,
		func(node *model.TimelineCounter) any {
			return &node.Count
		}, func(node *model.TimelineCounter) any {
			return scanner.ScanText(&node.EventType)
		}, func(node *model.TimelineCounter) any {
			return scanner.ScanTimestamp(&node.DateTo)
		}, func(node *model.TimelineCounter) any {
			return scanner.ScanTimestamp(&node.DateFrom)
		},
	)

	return
}

// endregion

var (
	CallsCTE = `select c.id::text,
       c.created_at,
       c.hangup_at                  AS                                        closed_at,
       CASE
           WHEN c.user_id NOTNULL THEN root.duration
           ELSE parent.duration END AS                                        duration,
       CASE
           WHEN c.user_id NOTNULL THEN parent.duration
           ELSE root.duration END   AS                                        total_duration,
       (with recursive a as (select *
                             from call_center.cc_calls_history
                             where id in (with recursive a as (select d.id::uuid, d.user_id
                                                               from call_center.cc_calls_history d
                                                               where d.id::uuid = c.id
                                                                 and d.domain_id = 1
                                                               union all
                                                               select d.id::uuid, d.user_id
                                                               from call_center.cc_calls_history d,
                                                                    a
                                                               where (d.parent_id::uuid = a.id::uuid or
                                                                      (d.transfer_from::uuid = a.id::uuid)))
                                          select distinct id
                                          from a))
        SELECT ARRAY_AGG(p.participant)
        FROM (SELECT ROW (usr.id, coalesce(usr.name, usr.username)) participant
              from a
                       inner join directory.wbt_user usr on a.user_id = usr.id
              order by a.created_at) as p)                                    participants,
       null::record                 as                                        gateway,
       ROW (scheme.id, scheme.name) as                                        flow_scheme,
       'call'                                                                 type,
       (c.direction = 'inbound')                                              is_inbound,
       (case when c.bridged_id isnull then c.queue_id notnull else false end) is_missed,
       ROW (c.queue_id, a.name)                                               queue,
       exists(with recursive a as (select*
                                   from call_center.cc_calls_history
                                   where id in (with recursive a as (select d.id::uuid, d.user_id
                                                                     from call_center.cc_calls_history d
                                                                     where d.id::uuid = c.id
                                                                       and d.domain_id = 1
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
                             a.blind_transfer::varchar) notnull)              is_detailed,
       files                                                                  files,
       transcripts                                                            transcripts
from call_center.cc_calls_history c
         left join directory.sip_gateway g on g.id = c.gateway_id
         left join call_center.cc_queue a on a.id = c.queue_id
    -- parent resuable columns
         LEFT JOIN LATERAL (SELECT round(date_part('epoch'::TEXT, hangup_at - created_at)::BIGINT) duration
                            FROM call_center.cc_calls_history
                            WHERE parent_id = c.id
                            LIMIT 1) parent ON true
    -- root reusable columns
         LEFT JOIN LATERAL (SELECT round(date_part('epoch'::text, c.hangup_at - c.created_at)::bigint) duration) root
                   ON true
    -- join flow_scheme
         LEFT JOIN flow.acr_routing_scheme scheme ON scheme.id = c.schema_ids[array_length(c.schema_ids, 1)]
    -- join files
         LEFT jOIN LATERAL (SELECT ARRAY_AGG(ROW (f1.id, f1.size, f1.mime_type, f1.name, f1.created_at))
                            FROM storage.files f1
                            WHERE f1.domain_id = c.domain_id
                              AND NOT f1.removed IS TRUE
                              AND f1.uuid = c.id::varchar) files ON true
    -- join transcripts
         LEFT JOIN LATERAL (SELECT ARRAY_AGG(ROW (tr.id, tr.locale, tr.file_id, ROW (ff.id, ff.name))) AS data
                            FROM storage.file_transcript tr
                                     LEFT JOIN storage.files ff ON ff.id = tr.file_id
                            WHERE tr.uuid::text = c.id::text
    ) transcripts ON true
where c.case_id = ?
  and c.transfer_from isnull`

	EmailsCTE = `SELECT e.id::text,
       e."from",
       e."to",
       ROW (e.profile_id, p.name)   AS profile,
       e.subject,
       e.cc,
       e.created_at,
       (e.direction = 'inbound')    AS is_inbound,
       e.sender,
       e.body,
       e.html,
       attachments                  AS attachments,
       ROW (e."owner_id", u."name") AS "user"
FROM call_center.cc_email e
    -- owner
         LEFT JOIN directory.wbt_user u ON e.owner_id = u.id
    -- email profile
         LEFT JOIN call_center.cc_email_profile p ON e.profile_id = p.id
    -- attachments
         LEFT JOIN LATERAL (SELECT ARRAY_AGG(ROW (f.id, f.mime_type, f.view_name, f.size))
                            from storage.files f
                            where f.id = any (e.attachment_ids)

    ) attachments ON true
WHERE e.case_ids && array [?]::int8[]`

	ChatsCTE = `SELECT conv.id::text,
       conv.created_at,
       conv.closed_at,
       round(extract(EPOCH FROM (conv.closed_at - conv.created_at))) as duration,
       participants                                                     participants,
       gateway                                                          gateway,
       ROW (flow_scheme.id, flow_scheme.name)                        as flow_scheme,
       'chat'                                                           type,
       true                                                             is_inbound,
       false                                                            is_missed,
       null::record                                                     queue,
       true                                                             is_detailed
FROM chat.conversation conv
    -- join participants
         LEFT JOIN LATERAL (SELECT ARRAY_AGG(ROW (usr.id, usr.name))
                            FROM chat.channel c
                                     INNER JOIN directory.wbt_auth usr ON user_id = usr.id
                            WHERE conversation_id = conv.id
                              AND internal
                              AND joined_at NOTNULL
                            GROUP BY usr.id) participants ON true
    -- join gateway
         LEFT JOIN LATERAL (SELECT ROW (b.id, b.name, b.provider) gateway
                            FROM chat.channel
                                     LEFT JOIN chat.bot b ON connection::bigint = b.id
                            WHERE NOT internal
                              AND conversation_id = conv.id
    ) gateway ON true
    -- join flow scheme
         LEFT JOIN flow.acr_routing_scheme flow_scheme ON flow_scheme.id = (conv.props ->> 'flow')::bigint
WHERE conv.id = ANY()`

	CallsCounterCTE = `select c.id::::text,
                      c.created_at,
                      c.hangup_at                                             as      closed_at,
                      'call'                                                          type

               from call_center.cc_calls_history c
               where  c.contact_id = :ContactId`

	ChatsCounterCTE = `select conv.id::::text,
                      conv.created_at,
                      conv.closed_at,
                      'chat'                                                           type
               from chat.conversation conv
               where conv.id = any (
                   select conversation_id
                   from chat.channel
                   where user_id =
                         any (select user_id from contacts.contact_imclient where contact_id = :ContactId) and not channel.internal)`

	EmailsCounterCTE = `select m.id::::text,
                      m.created_at,
                      m.created_at,
                      'email'                                                           type
               from call_center.cc_email m
               where m.contact_ids && array[:ContactId]::::int8[]`
)

func NewCaseTimelineStore(store store.Store) (store.CaseTimelineStore, error) {
	if store == nil {
		return nil, dberr.NewDBError("postgres.case_timeline.new_case_timeline_store.check_args.store",
			"store required")
	}
	return &CaseTimelineStore{storage: store}, nil
}
