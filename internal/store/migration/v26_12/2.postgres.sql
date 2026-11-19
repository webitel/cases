-- WTEL-10766: deleting a user clears Created by / Modified by instead of failing.

ALTER TABLE cases.skill_catalog ALTER COLUMN updated_by DROP NOT NULL;
ALTER TABLE cases.sla_condition ALTER COLUMN created_by DROP NOT NULL;
