-- WTEL-10766: user status Active/Deactivated.
ALTER TABLE cases.case_communication DROP CONSTRAINT IF EXISTS case_communication_wbt_user_id_fk;
ALTER TABLE cases.case_communication
    ADD CONSTRAINT case_communication_wbt_user_id_fk FOREIGN KEY (created_by) REFERENCES directory.wbt_user(id) ON DELETE SET NULL NOT VALID;
ALTER TABLE cases.case_communication VALIDATE CONSTRAINT case_communication_wbt_user_id_fk;

ALTER TABLE cases.close_reason_group DROP CONSTRAINT IF EXISTS close_reason_updated_id_fk;
ALTER TABLE cases.close_reason_group
    ADD CONSTRAINT close_reason_updated_id_fk FOREIGN KEY (updated_by) REFERENCES directory.wbt_user(id) ON DELETE SET NULL DEFERRABLE INITIALLY DEFERRED NOT VALID;
ALTER TABLE cases.close_reason_group VALIDATE CONSTRAINT close_reason_updated_id_fk;

ALTER TABLE cases.priority DROP CONSTRAINT IF EXISTS priority_updated_id_fk;
ALTER TABLE cases.priority
    ADD CONSTRAINT priority_updated_id_fk FOREIGN KEY (updated_by) REFERENCES directory.wbt_user(id) ON DELETE SET NULL DEFERRABLE INITIALLY DEFERRED NOT VALID;
ALTER TABLE cases.priority VALIDATE CONSTRAINT priority_updated_id_fk;

ALTER TABLE cases.service_catalog DROP CONSTRAINT IF EXISTS service_catalog_updated_id_fk;
ALTER TABLE cases.service_catalog
    ADD CONSTRAINT service_catalog_updated_id_fk FOREIGN KEY (updated_by) REFERENCES directory.wbt_user(id) ON DELETE SET NULL DEFERRABLE INITIALLY DEFERRED NOT VALID;
ALTER TABLE cases.service_catalog VALIDATE CONSTRAINT service_catalog_updated_id_fk;

ALTER TABLE cases.sla_condition DROP CONSTRAINT IF EXISTS sla_condition_updated_id_fk;
ALTER TABLE cases.sla_condition
    ADD CONSTRAINT sla_condition_updated_id_fk FOREIGN KEY (updated_by) REFERENCES directory.wbt_user(id) ON DELETE SET NULL DEFERRABLE INITIALLY DEFERRED NOT VALID;
ALTER TABLE cases.sla_condition VALIDATE CONSTRAINT sla_condition_updated_id_fk;

ALTER TABLE cases.sla DROP CONSTRAINT IF EXISTS sla_updated_id_fk;
ALTER TABLE cases.sla
    ADD CONSTRAINT sla_updated_id_fk FOREIGN KEY (updated_by) REFERENCES directory.wbt_user(id) ON DELETE SET NULL DEFERRABLE INITIALLY DEFERRED NOT VALID;
ALTER TABLE cases.sla VALIDATE CONSTRAINT sla_updated_id_fk;

ALTER TABLE cases.source DROP CONSTRAINT IF EXISTS source_updated_id_fk;
ALTER TABLE cases.source
    ADD CONSTRAINT source_updated_id_fk FOREIGN KEY (updated_by) REFERENCES directory.wbt_user(id) ON DELETE SET NULL DEFERRABLE INITIALLY DEFERRED NOT VALID;
ALTER TABLE cases.source VALIDATE CONSTRAINT source_updated_id_fk;

ALTER TABLE cases.status DROP CONSTRAINT IF EXISTS status_updated_id_fk;
ALTER TABLE cases.status
    ADD CONSTRAINT status_updated_id_fk FOREIGN KEY (updated_by) REFERENCES directory.wbt_user(id) ON DELETE SET NULL DEFERRABLE INITIALLY DEFERRED NOT VALID;
ALTER TABLE cases.status VALIDATE CONSTRAINT status_updated_id_fk;
