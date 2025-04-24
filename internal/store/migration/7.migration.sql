CREATE OR REPLACE FUNCTION check_sla_deletion() RETURNS trigger AS $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM cases.service_catalog sc
        WHERE sc.sla_id = OLD.id AND sc.root_id IS NULL
    ) THEN
        RAISE EXCEPTION 'Cannot delete SLA with id %, it is referenced by a root service_catalog entry', OLD.id;
    END IF;

    -- Set sla_id = null for all referencing rows with non-null root_id
    UPDATE cases.service_catalog
    SET sla_id = NULL
    WHERE sla_id = OLD.id AND root_id IS NOT NULL;

    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_restrict_or_nullify_sla
    BEFORE DELETE ON cases.sla
    FOR EACH ROW
EXECUTE FUNCTION check_sla_deletion();
