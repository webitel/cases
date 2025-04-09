alter table cases."case" add column if not exists is_overdue bool default false;

create index concurrently if not exists case_planned_resolve_at_notify_index
    on cases."case" (planned_resolve_at) include (id)
    where not "case".is_overdue and close_reason isnull;