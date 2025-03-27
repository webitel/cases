package postgres

import (
	"fmt"
	"github.com/webitel/cases/internal/store/util"
	"github.com/webitel/cases/model/options"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgtype"
	"github.com/webitel/cases/api/cases"
	dberr "github.com/webitel/cases/internal/errors"
	"github.com/webitel/cases/internal/store"
	"github.com/webitel/cases/internal/store/postgres/scanner"
	"github.com/webitel/cases/model"
)

var CaseTimelineFields = []string{
	"calls", "chats", "emails",
}

type CaseTimelineStore struct {
	storage *Store
}

func (c *CaseTimelineStore) Get(rpc options.SearchOptions) (*cases.GetTimelineResponse, error) {
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

	rows, err := db.Query(rpc, util.CompactSQL(sql), args...)
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
			return nil, dberr.NewDBInternalError("postgres.case_timeline.get.scan.error", err)
		}
		result.Days = append(result.Days, node)
	}
	result.Days, result.Next = util.ResolvePaging(rpc.GetSize(), result.Days)
	result.Page = int32(rpc.GetPage())

	return result, nil
}

// region Timeline Build Functions
func buildCaseTimelineSqlizer(rpc options.SearchOptions) (squirrel.Sqlizer, []func(timeline *cases.DayTimeline) any, *dberr.DBError) {
	if rpc == nil {
		return nil, nil, dberr.NewDBError("postgres.case_timeline.build_case_timeline_sqlizer.check_args.rpc", "search options required")
	}
	parentId, ok := rpc.GetFilter("case_id").(int64)
	if !ok || parentId == 0 {
		return nil, nil, dberr.NewDBError("postgres.case_timeline.build_case_timeline_sqlizer.check_args.case_id", "case id required")
	}
	fields := rpc.GetFields()
	if len(fields) == 0 {
		fields = CaseTimelineFields
	}
	var (
		ctes       = make(map[string]squirrel.Sqlizer)
		chatsPlan  []func(timeline **cases.Event) any
		callsPlan  []func(timeline **cases.Event) any
		emailsPlan []func(timeline **cases.Event) any
		dayEvents  []string
		from       = "( SELECT * FROM ("
		plan       []func(timeline *cases.DayTimeline) any
	)
	for i, field := range fields {
		var (
			err       *dberr.DBError
			cte       squirrel.Sqlizer
			eventType string
		)
		switch field {
		case "chat":
			cte, chatsPlan, err = buildTimelineChatsColumn(parentId)
			eventType = cases.CaseTimelineEventType_chat.String()
		case "call":
			cte, callsPlan, err = buildTimelineCallsColumn(parentId)
			eventType = cases.CaseTimelineEventType_call.String()
		case "email":
			cte, emailsPlan, err = buildTimelineEmailsColumn(parentId)
			eventType = cases.CaseTimelineEventType_email.String()
		default:
			return nil, nil, dberr.NewDBError("postgres.case_timeline.build_case_timeline_sqlizer.parse_fields.unknown", "unknown field "+field)
		}
		if err != nil {
			return nil, nil, err
		}
		ctes[field] = cte
		if i != 0 {
			from += " union all "
		}
		from += fmt.Sprintf(`SELECT
					  DATE_TRUNC('day', created_at) "day",
					  created_at,
					  '%s' "event",
					  (%s)::text "data"
					   from %[2]s
						`, eventType, field)
	}
	from += ") AS ec ORDER BY created_at DESC) e"
	cteQuery, args, err := util.FormAsCTEs(ctes)
	if err != nil {
		return nil, nil, dberr.NewDBError("postgres.case_timeline.build_case_timeline_sqlizer.form_cte.error", err.Error())
	}
	query := squirrel.Select("day", "array_agg(e.event)::text \"events\"", "array_agg(e.data)::text \"data\"").
		GroupBy("e.day").
		OrderBy("e.day desc").From(from).
		Prefix(cteQuery, args...).
		PlaceholderFormat(squirrel.Dollar)
	query = util.ApplyPaging(rpc.GetPage(), rpc.GetSize(), query)

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

				size := len(rows.Elements)

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

				var node *cases.Event

				// DECODE
				for r, elem := range rows.Elements {
					// RECORD
					node = &cases.Event{} // NEW
					var (
						eventType string
						scanPlan  []func(**cases.Event) any
					)
					// ALLOC
					switch eventType = dayEvents[r]; eventType {
					case cases.CaseTimelineEventType_email.String():
						actualEvent := &cases.EmailEvent{}
						node.Type = cases.CaseTimelineEventType_email
						node.Event = &cases.Event_Email{Email: actualEvent}
						count := &rec.EmailsCount
						*count++

						// scanning functions set
						scanPlan = emailsPlan
					case cases.CaseTimelineEventType_chat.String():
						actualEvent := &cases.ChatEvent{}
						node.Type = cases.CaseTimelineEventType_chat
						node.Event = &cases.Event_Chat{Chat: actualEvent}
						count := &rec.ChatsCount
						*count++

						// scanning functions set
						scanPlan = chatsPlan
					case cases.CaseTimelineEventType_call.String():
						actualEvent := &cases.CallEvent{}
						node.Type = cases.CaseTimelineEventType_call
						node.Event = &cases.Event_Call{Call: actualEvent}
						count := &rec.CallsCount
						*count++

						// scanning functions set
						scanPlan = callsPlan
					}
					// DECODE
					scan := pgtype.NewCompositeTextScanner(pgtype.NewConnInfo(), []byte(elem))
					for _, bind := range scanPlan {
						df := bind(&node)
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

func buildTimelineChatsColumn(caseId int64) (base squirrel.Sqlizer, plan []func(timeline **cases.Event) any, dbError *dberr.DBError) {
	if caseId == 0 {
		return nil, nil, dberr.NewDBError("postgres.case_timeline.build_timeline_chats_column.check_args.case_id.empty", "case id required")
	}
	base = squirrel.Expr(ChatsCTE, caseId, int32(cases.CaseCommunicationsTypes_COMMUNICATION_CHAT))
	plan = append(plan,
		func(node **cases.Event) any {
			buf := *node
			chat := buf.GetChat()
			return scanner.ScanText(&chat.Id)
		},
		func(node **cases.Event) any {
			buf := *node
			return scanner.ScanTimestamp(&buf.CreatedAt)
		},
		func(node **cases.Event) any {
			buf := *node
			chat := buf.GetChat()
			return scanner.ScanTimestamp(&chat.ClosedAt)
		},
		func(node **cases.Event) any {
			buf := *node
			chat := buf.GetChat()
			return &chat.Duration
		},
		func(node **cases.Event) any {
			buf := *node
			chat := buf.GetChat()
			return scanner.ScanLookupList(&chat.Participants)
		},
		func(node **cases.Event) any {
			buf := *node
			chat := buf.GetChat()
			return scanner.ScanRowExtendedLookup(&chat.Gateway)
		},
		func(node **cases.Event) any {
			buf := *node
			chat := buf.GetChat()
			return scanner.ScanRowLookup(&chat.FlowScheme)
		},
		func(node **cases.Event) any {
			buf := *node
			chat := buf.GetChat()
			return &chat.IsInbound
		},
		func(node **cases.Event) any {
			buf := *node
			chat := buf.GetChat()
			return &chat.IsMissed
		},
		func(node **cases.Event) any {
			buf := *node
			chat := buf.GetChat()
			return scanner.ScanRowLookup(&chat.Queue)
		},
		func(node **cases.Event) any {
			buf := *node
			chat := buf.GetChat()
			return &chat.IsDetailed
		})
	return
}

func buildTimelineCallsColumn(caseId int64) (base squirrel.Sqlizer, plan []func(timeline **cases.Event) any, dbError *dberr.DBError) {
	if caseId == 0 {
		return nil, nil, dberr.NewDBError("postgres.case_timeline.build_timeline_calls_column.check_args.case_id.empty", "case id required")
	}
	base = squirrel.Expr(CallsCTE, caseId, int32(cases.CaseCommunicationsTypes_COMMUNICATION_CALL))

	plan = append(plan,
		func(node **cases.Event) any {
			buf := *node
			call := buf.GetCall()
			return scanner.ScanText(&call.Id)
		},
		func(node **cases.Event) any {
			buf := *node
			return scanner.ScanTimestamp(&buf.CreatedAt)
		},
		func(node **cases.Event) any {
			buf := *node
			call := buf.GetCall()
			return scanner.ScanTimestamp(&call.ClosedAt)
		},
		func(node **cases.Event) any {
			buf := *node
			call := buf.GetCall()
			return &call.Duration
		},
		func(node **cases.Event) any {
			buf := *node
			call := buf.GetCall()
			return &call.TotalDuration
		},
		func(node **cases.Event) any {
			buf := *node
			call := buf.GetCall()
			return scanner.ScanLookupList(&call.Participants)
		},
		func(node **cases.Event) any {
			buf := *node
			call := buf.GetCall()
			return scanner.ScanRowLookup(&call.Gateway)
		},
		func(node **cases.Event) any {
			buf := *node
			call := buf.GetCall()
			return scanner.ScanRowLookup(&call.FlowScheme)
		},
		func(node **cases.Event) any {
			buf := *node
			call := buf.GetCall()
			return scanner.ScanBool(&call.IsInbound)
		},
		func(node **cases.Event) any {
			buf := *node
			call := buf.GetCall()
			return &call.IsMissed
		},
		func(node **cases.Event) any {
			buf := *node
			call := buf.GetCall()
			return scanner.ScanRowLookup(&call.Queue)
		},
		func(node **cases.Event) any {
			buf := *node
			call := buf.GetCall()
			return &call.IsDetailed
		},
		// files
		func(node **cases.Event) any {
			buf := *node
			value := buf.GetCall()
			return scanner.TextDecoder(func(src []byte) error {
				if len(src) == 0 {
					return nil // NULL
				}
				array, inErr := pgtype.ParseUntypedTextArray(string(src))
				if inErr != nil {
					return inErr
				}

				scanPlan := []func(file *cases.CallFile) any{
					// id
					func(file *cases.CallFile) any {
						return scanner.ScanInt64(&file.Id)
					},
					// size
					func(file *cases.CallFile) any {
						return scanner.ScanInt64(&file.Size)
					},
					// mime type
					func(file *cases.CallFile) any {
						return scanner.ScanText(&file.MimeType)
					},
					func(file *cases.CallFile) any {
						return scanner.ScanText(&file.Name)
					},
					func(file *cases.CallFile) any {
						return scanner.ScanInt64(&file.StartAt)
					},
				}

				var err error
				for _, element := range array.Elements {
					var (
						file cases.CallFile
						raw  = pgtype.NewCompositeTextScanner(pgtype.NewConnInfo(), []byte(element))
					)
					for _, bind := range scanPlan {
						raw.ScanValue(bind(&file))
						err = raw.Err()
						if err != nil {
							return err
						}
					}
					value.Files = append(value.Files, &file)
				}
				return nil
			})
		},
		// transcripts
		func(node **cases.Event) any {
			buf := *node
			value := buf.GetCall()
			return scanner.TextDecoder(func(src []byte) error {
				if len(src) == 0 {
					return nil // NULL
				}
				array, inErr := pgtype.ParseUntypedTextArray(string(src))
				if inErr != nil {
					return inErr
				}

				scanPlan := []func(*cases.TranscriptLookup) any{
					// id
					func(transcript *cases.TranscriptLookup) any {
						return scanner.ScanInt64(&transcript.Id)
					},
					// size
					func(transcript *cases.TranscriptLookup) any {
						return scanner.ScanText(&transcript.Locale)
					},
					// file
					func(transcript *cases.TranscriptLookup) any {
						return scanner.ScanRowLookup(&transcript.File)
					},
				}

				var err error
				for _, element := range array.Elements {
					var (
						file cases.TranscriptLookup
						raw  = pgtype.NewCompositeTextScanner(pgtype.NewConnInfo(), []byte(element))
					)
					for _, bind := range scanPlan {
						raw.ScanValue(bind(&file))
						err = raw.Err()
						if err != nil {
							return err
						}
					}
					value.Transcripts = append(value.Transcripts, &file)
				}
				return nil
			})
		},
	)
	return
}

func buildTimelineEmailsColumn(caseId int64) (base squirrel.Sqlizer, plan []func(timeline **cases.Event) any, dbError *dberr.DBError) {
	if caseId == 0 {
		return nil, nil, dberr.NewDBError("postgres.case_timeline.build_timeline_emails_column.check_args.case_id.empty", "case id required")
	}
	base = squirrel.Expr(EmailsCTE, caseId, int32(cases.CaseCommunicationsTypes_COMMUNICATION_EMAIL))

	plan = append(plan,
		func(node **cases.Event) any {
			buf := *node
			email := buf.GetEmail()
			return scanner.ScanText(&email.Id)
		},
		func(node **cases.Event) any {
			buf := *node
			email := buf.GetEmail()
			return &email.From
		},
		func(node **cases.Event) any {
			buf := *node
			email := buf.GetEmail()
			return &email.To
		},
		func(node **cases.Event) any {
			buf := *node
			email := buf.GetEmail()
			return scanner.ScanRowLookup(&email.Profile)
		},
		func(node **cases.Event) any {
			buf := *node
			email := buf.GetEmail()
			return scanner.ScanText(&email.Subject)
		},
		func(node **cases.Event) any {
			buf := *node
			email := buf.GetEmail()
			return &email.Cc
		},
		func(node **cases.Event) any {
			buf := *node
			return scanner.ScanTimestamp(&buf.CreatedAt)
		},
		func(node **cases.Event) any {
			buf := *node
			email := buf.GetEmail()
			return &email.IsInbound
		},
		func(node **cases.Event) any {
			buf := *node
			email := buf.GetEmail()
			return &email.Sender
		},
		func(node **cases.Event) any {
			buf := *node
			email := buf.GetEmail()
			return scanner.ScanText(&email.Body)
		},
		func(node **cases.Event) any {
			buf := *node
			email := buf.GetEmail()
			return scanner.ScanText(&email.Html)
		},
		// attachments
		func(node **cases.Event) any {
			buf := *node
			email := buf.GetEmail()
			return scanner.TextDecoder(func(src []byte) error {
				if len(src) == 0 {
					return nil // NULL
				}
				array, inErr := pgtype.ParseUntypedTextArray(string(src))
				if inErr != nil {
					return inErr
				}

				scanPlan := []func(file *cases.Attachment) any{
					// id
					func(file *cases.Attachment) any {
						return scanner.ScanInt64(&file.Id)
					},
					// mime type
					func(file *cases.Attachment) any {
						return scanner.ScanText(&file.Mime)
					},
					func(file *cases.Attachment) any {
						return scanner.ScanText(&file.Name)
					},
					// size
					func(file *cases.Attachment) any {
						return scanner.ScanInt64(&file.Size)
					},
				}

				var err error
				for _, element := range array.Elements {
					var (
						file cases.Attachment
						raw  = pgtype.NewCompositeTextScanner(pgtype.NewConnInfo(), []byte(element))
					)
					for _, bind := range scanPlan {
						raw.ScanValue(bind(&file))
						err = raw.Err()
						if err != nil {
							return err
						}
					}
					email.Attachments = append(email.Attachments, &file)
				}
				return nil
			})
		},
		func(node **cases.Event) any {
			buf := *node
			email := buf.GetEmail()
			return scanner.ScanRowLookup(&email.Owner)
		})

	return
}

// endregion

func (c *CaseTimelineStore) GetCounter(rpc options.SearchOptions) ([]*model.TimelineCounter, error) {
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
	rows, err := db.Query(rpc, util.CompactSQL(sql), args...)
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

func buildTimelineCounterSqlizer(rpc options.SearchOptions) (query squirrel.Sqlizer, scanPlan []func(response *model.TimelineCounter) any, dbError *dberr.DBError) {
	if rpc == nil {
		return nil, nil, dberr.NewDBError("postgres.case_timeline.build_case_timeline_sqlizer.check_args.rpc", "search options required")
	}
	if len(rpc.GetIDs()) == 0 {
		return nil, nil, dberr.NewDBError("postgres.case_timeline.build_case_timeline_sqlizer.check_args.case_id", "case id empty")
	}
	caseId := rpc.GetIDs()[0]
	if caseId <= 0 {
		return nil, nil, dberr.NewDBError("postgres.case_timeline.build_case_timeline_sqlizer.check_args.case_id", "case id empty")
	}
	fields := rpc.GetFields()
	if len(fields) == 0 {
		fields = CaseTimelineFields
	}
	var (
		ctes = make(map[string]squirrel.Sqlizer)
		from = "("
	)
	for i, field := range fields {
		switch field {
		case cases.CaseTimelineEventType_call.String():
			communicationType := int64(cases.CaseCommunicationsTypes_COMMUNICATION_CALL)
			ctes[field] = squirrel.Expr(CallsCounterCTE, communicationType, caseId, communicationType)
		case cases.CaseTimelineEventType_email.String():
			communicationType := int64(cases.CaseCommunicationsTypes_COMMUNICATION_EMAIL)
			ctes[field] = squirrel.Expr(EmailsCounterCTE, communicationType, caseId, communicationType)
		case cases.CaseTimelineEventType_chat.String():
			communicationType := int64(cases.CaseCommunicationsTypes_COMMUNICATION_CHAT)
			ctes[field] = squirrel.Expr(ChatsCounterCTE, communicationType, caseId, communicationType)
		default:
			return nil, nil, dberr.NewDBError("postgres.case_timeline.build_case_timeline_counter_sqlizer.parse_fields.unknown", "unknown field "+field)
		}
		if i != 0 {
			from += " union all "
		}
		from += fmt.Sprintf("(select * from %s)", field)
	}
	from += ") s"
	cteQuery, args, err := util.FormAsCTEs(ctes)
	if err != nil {
		return nil, nil, dberr.NewDBError("postgres.case_timeline.build_case_timeline_counter_sqlizer.parse_fields.unknown", err.Error())
	}

	query = squirrel.Select("count(*) count",
		"type event",
		"max(closed_at) date_to",
		"min(created_at) date_from").
		From(from).
		Prefix(cteQuery, args...).
		GroupBy("type").
		PlaceholderFormat(squirrel.Dollar)

	// each row represents type of event
	scanPlan = append(scanPlan,
		func(node *model.TimelineCounter) any {
			return scanner.ScanInt64(&node.Count)
		}, func(node *model.TimelineCounter) any {
			return scanner.ScanInt64(&node.EventType)
		}, func(node *model.TimelineCounter) any {
			return scanner.ScanTimestamp(&node.DateTo)
		}, func(node *model.TimelineCounter) any {
			return scanner.ScanTimestamp(&node.DateFrom)
		},
	)

	return
}

// endregion

const (
	CallsCTE = `select c.id::text,
       c.created_at,
       c.hangup_at                  AS                                        closed_at,
       round(case when c.user_id notnull then date_part('epoch'::text, c.hangup_at - c.created_at)::bigint else (select date_part('epoch'::text, hangup_at - created_at)::bigint from call_center.cc_calls_history where parent_id = c.id limit 1) end)::bigint as      duration,
       root.duration                                        total_duration,
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
       files.data                                                                  files,
       transcripts.data                                                          transcripts
from call_center.cc_calls_history c
         left join directory.sip_gateway g on g.id = c.gateway_id
         left join call_center.cc_queue a on a.id = c.queue_id
    -- root reusable columns
         LEFT JOIN LATERAL (SELECT round(date_part('epoch'::text, c.hangup_at - c.created_at)::bigint) duration) root
                   ON true
    -- join flow_scheme
         LEFT JOIN flow.acr_routing_scheme scheme ON scheme.id = c.schema_ids[array_length(c.schema_ids, 1)]
    -- join files
         LEFT jOIN LATERAL (SELECT ARRAY_AGG(ROW (f1.id, f1.size, f1.mime_type, f1.name, f1.created_at)) data
                            FROM storage.files f1
                            WHERE f1.domain_id = c.domain_id
                              AND NOT f1.removed IS TRUE
                              AND f1.uuid = c.id::varchar) files ON true
    -- join transcripts
         LEFT JOIN LATERAL (SELECT ARRAY_AGG(ROW (tr.id, tr.locale, ROW (ff.id, ff.name))) AS data
                            FROM storage.file_transcript tr
                                     LEFT JOIN storage.files ff ON ff.id = tr.file_id
                            WHERE tr.uuid::text = c.id::text
    ) transcripts ON true
where c.id = ANY(SELECT communication_id::uuid FROM cases.case_communication WHERE case_id = ? AND communication_type = ?)
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
       attachments.data                  AS attachments,
       ROW (e."owner_id", u."name") AS "user"
FROM call_center.cc_email e
    -- owner
         LEFT JOIN directory.wbt_user u ON e.owner_id = u.id
    -- email profile
         LEFT JOIN call_center.cc_email_profile p ON e.profile_id = p.id
    -- attachments
         LEFT JOIN LATERAL (SELECT ARRAY_AGG(ROW (f.id, f.mime_type, f.view_name, f.size)) data
                            from storage.files f
                            where f.id = any (e.attachment_ids)

    ) attachments ON true
WHERE e.id = ANY(SELECT communication_id::bigint FROM cases.case_communication WHERE case_id = ? AND communication_type =?)`

	ChatsCTE = `SELECT conv.id::text,
       conv.created_at,
       conv.closed_at,
       round(extract(EPOCH FROM (conv.closed_at - conv.created_at))) as duration,
       participants.data                                                     participants,
       gateway.data                                                          gateway,
       ROW (flow_scheme.id, flow_scheme.name)                        as flow_scheme,
       true                                                             is_inbound,
       false                                                            is_missed,
       null::record                                                     queue,
       true                                                             is_detailed
FROM chat.conversation conv
    -- join participants
         LEFT JOIN LATERAL (SELECT ARRAY_AGG(ROW (usr.id, usr.name)) data
                            FROM chat.channel c
                                     INNER JOIN directory.wbt_auth usr ON user_id = usr.id
                            WHERE conversation_id = conv.id
                              AND internal
                              AND joined_at NOTNULL
                            GROUP BY usr.id) participants ON true
    -- join gateway
         LEFT JOIN LATERAL (SELECT ROW (b.id, b.name, b.provider) data
                            FROM chat.channel
                                     LEFT JOIN chat.bot b ON connection::bigint = b.id
                            WHERE NOT internal
                              AND conversation_id = conv.id
    ) gateway ON true
    -- join flow scheme
         LEFT JOIN flow.acr_routing_scheme flow_scheme ON flow_scheme.id = (conv.props ->> 'flow')::bigint
WHERE conv.id =  ANY(SELECT communication_id::uuid FROM cases.case_communication WHERE case_id = ? AND communication_type = ?)`

	CallsCounterCTE = `SELECT c.id::text,
                      c.created_at,
                      c.hangup_at                                             AS      closed_at,
                      ?::int                                                          type

               FROM call_center.cc_calls_history c
               WHERE c.id = ANY(SELECT communication_id::uuid FROM cases.case_communication WHERE case_id = ? AND communication_type = ?::int)`

	ChatsCounterCTE = `select conv.id::text,
                      conv.created_at,
                      conv.closed_at,
                      ?::int                                                           type
               from chat.conversation conv
               WHERE conv.id = ANY(SELECT communication_id::uuid FROM cases.case_communication WHERE case_id = ? AND communication_type = ?::int)`

	EmailsCounterCTE = `select m.id::text,
                      m.created_at,
                      m.created_at,
                      ?::int                                                           type
               from call_center.cc_email m
               WHERE m.id = ANY(SELECT communication_id::bigint FROM cases.case_communication WHERE case_id = ? AND communication_type = ?::int)`
)

func NewCaseTimelineStore(store *Store) (store.CaseTimelineStore, error) {
	if store == nil {
		return nil, dberr.NewDBError("postgres.case_timeline.new_case_timeline_store.check_args.store",
			"store required")
	}
	return &CaseTimelineStore{storage: store}, nil
}
