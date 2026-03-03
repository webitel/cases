ALTER TABLE cases.service_catalog
    ADD COLUMN IF NOT EXISTS default_priority_id bigint;

ALTER TABLE cases.service_catalog
    ADD CONSTRAINT fk_service_catalog_default_priority
        FOREIGN KEY (default_priority_id)
            REFERENCES cases.priority (id)
            ON DELETE SET NULL;
