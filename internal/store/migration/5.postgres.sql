alter table cases."case" add column if not exists is_overdue bool default false;

create index concurrently if not exists case_planned_resolve_at_notify_index
    on cases."case" (planned_resolve_at) include (id)
    where not "case".is_overdue and resolved_at isnull;

create or replace function update_case_timings() returns trigger
    language plpgsql
as
$$
BEGIN
    -- Set reacted_at if status_condition is not initial
    IF (NEW.status_condition IS NOT NULL) THEN
        -- Check if the status_condition is initial
        PERFORM initial
        FROM cases.status_condition
        WHERE id = NEW.status_condition AND initial = TRUE;

        -- If no initial status found, set reacted_at
        IF NOT FOUND AND NEW.reacted_at ISNULL THEN
            NEW.reacted_at = timezone('utc', now());
        END IF;
    END IF;

    -- Set resolved_at if status_condition is final
    IF (NEW.status_condition IS NOT NULL) THEN
        -- Check if the status_condition is final
        PERFORM final
        FROM cases.status_condition
        WHERE id = NEW.status_condition AND final = TRUE;

        -- If final status found, set resolved_at
        IF FOUND THEN
            NEW.resolved_at = timezone('utc', now());
        ELSE
            -- If no final status, reset resolved_at to NULL
            NEW.resolved_at = NULL;
        END IF;
    END IF;

    if (TG_OP = 'UPDATE' AND NEW.resolved_at ISNULL AND NEW.is_overdue AND NEW.planned_resolve_at != OLD.planned_resolve_at) THEN
        NEW.is_overdue = false;
    END IF;

    RETURN NEW;
END;
$$;