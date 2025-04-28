create or replace function cases.update_case_timings() returns trigger
    language plpgsql
as
$$
DECLARE
    is_initial BOOLEAN := FALSE;
    is_final BOOLEAN := FALSE;
BEGIN
    IF (NEW.status_condition IS NOT NULL) THEN
        -- Fetch both initial and final flags for the given status_condition
        SELECT initial, final
        INTO is_initial, is_final
        FROM cases.status_condition
        WHERE id = NEW.status_condition;

        -- Set reacted_at if status is not initial and reacted_at hasn't been set
        IF NOT is_initial AND NEW.reacted_at IS NULL THEN
            NEW.reacted_at = timezone('utc', now());
        ELSIF is_initial AND is_final AND NEW.reacted_at IS NULL THEN
            -- Special case: if status is both initial and final, still set reacted_at
            NEW.reacted_at = timezone('utc', now());
        END IF;

        -- Set resolved_at if the status is final
        IF is_final THEN
            NEW.resolved_at = timezone('utc', now());
        ELSE
            -- If it's not a final status, reset resolved_at to NULL
            NEW.resolved_at = NULL;
        END IF;
    END IF;

    IF (TG_OP = 'UPDATE' AND NEW.resolved_at ISNULL AND NEW.is_overdue AND NEW.planned_resolve_at != OLD.planned_resolve_at) THEN
        NEW.is_overdue = false;
    END IF;

    RETURN NEW;
END;
$$;

