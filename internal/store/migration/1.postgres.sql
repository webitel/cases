-- cases.appeal definition

-- Drop table

-- DROP TABLE cases.appeal;

CREATE TABLE cases.appeal (
	id int8 DEFAULT nextval('cases.appeal_id'::regclass) NOT NULL,
	"name" text NOT NULL,
	description text NULL,
	created_at timestamp DEFAULT timezone('utc'::text, now()) NOT NULL,
	updated_at timestamp DEFAULT timezone('utc'::text, now()) NOT NULL,
	created_by int8 NULL,
	updated_by int8 NULL,
	"type" text NOT NULL,
	dc int8 NOT NULL,
	CONSTRAINT appeal_pk PRIMARY KEY (id),
	CONSTRAINT apppeal_fk UNIQUE (id, dc)
);
CREATE INDEX appeal_dc ON cases.appeal USING btree (dc);


-- cases.appeal foreign keys

ALTER TABLE cases.appeal ADD CONSTRAINT appeal_created_dc_fk FOREIGN KEY (created_by,dc) REFERENCES directory.wbt_user(id,dc) DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE cases.appeal ADD CONSTRAINT appeal_created_id_fk FOREIGN KEY (created_by) REFERENCES directory.wbt_user(id) ON DELETE SET NULL DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE cases.appeal ADD CONSTRAINT appeal_domain_fk FOREIGN KEY (dc) REFERENCES directory.wbt_domain(dc) ON DELETE CASCADE;
ALTER TABLE cases.appeal ADD CONSTRAINT appeal_updated_dc_fk FOREIGN KEY (updated_by,dc) REFERENCES directory.wbt_user(id,dc) DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE cases.appeal ADD CONSTRAINT appeal_updated_id_fk FOREIGN KEY (updated_by) REFERENCES directory.wbt_user(id) DEFERRABLE INITIALLY DEFERRED;

-- cases.close_reason definition

-- Drop table

-- DROP TABLE cases.close_reason;

CREATE TABLE cases.close_reason (
	id int8 DEFAULT nextval('cases.close_reason_id'::regclass) NOT NULL,
	"name" text NOT NULL,
	description text NOT NULL,
	created_at timestamp DEFAULT timezone('utc'::text, now()) NOT NULL,
	updated_at timestamp DEFAULT timezone('utc'::text, now()) NOT NULL,
	created_by int8 NULL,
	updated_by int8 NULL,
	dc int8 NOT NULL,
	CONSTRAINT close_reason_fk UNIQUE (id, dc),
	CONSTRAINT close_reason_pk PRIMARY KEY (id)
);
CREATE INDEX close_reason_dc ON cases.close_reason USING btree (dc);


-- cases.close_reason foreign keys

ALTER TABLE cases.close_reason ADD CONSTRAINT close_reason_created_dc_fk FOREIGN KEY (created_by,dc) REFERENCES directory.wbt_user(id,dc) DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE cases.close_reason ADD CONSTRAINT close_reason_created_id_fk FOREIGN KEY (created_by) REFERENCES directory.wbt_user(id) ON DELETE SET NULL DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE cases.close_reason ADD CONSTRAINT close_reason_domain_fk FOREIGN KEY (dc) REFERENCES directory.wbt_domain(dc) ON DELETE CASCADE;
ALTER TABLE cases.close_reason ADD CONSTRAINT close_reason_updated_dc_fk FOREIGN KEY (updated_by,dc) REFERENCES directory.wbt_user(id,dc) DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE cases.close_reason ADD CONSTRAINT close_reason_updated_id_fk FOREIGN KEY (updated_by) REFERENCES directory.wbt_user(id) DEFERRABLE INITIALLY DEFERRED;

-- cases.priority definition

-- Drop table

-- DROP TABLE cases.priority;

CREATE TABLE cases.priority (
	id int8 DEFAULT nextval('cases.priority_id'::regclass) NOT NULL,
	"name" text NOT NULL,
	description text NOT NULL,
	created_at timestamp DEFAULT timezone('utc'::text, now()) NOT NULL,
	updated_at timestamp DEFAULT timezone('utc'::text, now()) NOT NULL,
	created_by int8 NULL,
	updated_by int8 NULL,
	dc int8 NOT NULL,
	color varchar(25) NOT NULL,
	CONSTRAINT priority_fk UNIQUE (id, dc),
	CONSTRAINT priority_pk PRIMARY KEY (id)
);
CREATE INDEX priority_dc ON cases.priority USING btree (dc);


-- cases.priority foreign keys

ALTER TABLE cases.priority ADD CONSTRAINT priority_created_dc_fk FOREIGN KEY (created_by,dc) REFERENCES directory.wbt_user(id,dc) DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE cases.priority ADD CONSTRAINT priority_created_id_fk FOREIGN KEY (created_by) REFERENCES directory.wbt_user(id) ON DELETE SET NULL DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE cases.priority ADD CONSTRAINT priority_domain_fk FOREIGN KEY (dc) REFERENCES directory.wbt_domain(dc) ON DELETE CASCADE;
ALTER TABLE cases.priority ADD CONSTRAINT priority_updated_dc_fk FOREIGN KEY (updated_by,dc) REFERENCES directory.wbt_user(id,dc) DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE cases.priority ADD CONSTRAINT priority_updated_id_fk FOREIGN KEY (updated_by) REFERENCES directory.wbt_user(id) DEFERRABLE INITIALLY DEFERRED;

-- cases.priority_sla_condition definition

-- Drop table

-- DROP TABLE cases.priority_sla_condition;

CREATE TABLE cases.priority_sla_condition (
	id int8 DEFAULT nextval('cases.priority_sla_condition_id'::regclass) NOT NULL,
	created_at timestamp DEFAULT timezone('utc'::text, now()) NOT NULL,
	updated_at timestamp DEFAULT timezone('utc'::text, now()) NOT NULL,
	created_by int8 NULL,
	updated_by int8 NULL,
	sla_condition_id int8 NOT NULL,
	priority_id int8 NOT NULL,
	dc int8 NOT NULL,
	CONSTRAINT priority_sla_condition_fk UNIQUE (id, dc),
	CONSTRAINT priority_sla_condition_pk PRIMARY KEY (id)
);
CREATE UNIQUE INDEX idx_sla_condition_priority ON cases.priority_sla_condition USING btree (sla_condition_id, priority_id);
CREATE INDEX priority_sla_condition_dc ON cases.priority_sla_condition USING btree (dc);
CREATE INDEX priority_sla_condition_source ON cases.priority_sla_condition USING btree (sla_condition_id);


-- cases.priority_sla_condition foreign keys

ALTER TABLE cases.priority_sla_condition ADD CONSTRAINT priority_sla_condition_created_by_fk FOREIGN KEY (created_by) REFERENCES directory.wbt_user(id) ON DELETE SET NULL DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE cases.priority_sla_condition ADD CONSTRAINT priority_sla_condition_created_dc_fk FOREIGN KEY (created_by,dc) REFERENCES directory.wbt_user(id,dc) DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE cases.priority_sla_condition ADD CONSTRAINT priority_sla_condition_domain_fk FOREIGN KEY (dc) REFERENCES directory.wbt_domain(dc) ON DELETE CASCADE;
ALTER TABLE cases.priority_sla_condition ADD CONSTRAINT priority_sla_condition_sla_condition_id_fk FOREIGN KEY (sla_condition_id) REFERENCES cases.sla_condition(id) ON DELETE CASCADE;
ALTER TABLE cases.priority_sla_condition ADD CONSTRAINT priority_sla_condition_updated_by_fk FOREIGN KEY (updated_by) REFERENCES directory.wbt_user(id) ON DELETE SET NULL DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE cases.priority_sla_condition ADD CONSTRAINT priority_sla_condition_updated_dc_fk FOREIGN KEY (updated_by,dc) REFERENCES directory.wbt_user(id,dc) DEFERRABLE INITIALLY DEFERRED;

-- cases.reason definition

-- Drop table

-- DROP TABLE cases.reason;

CREATE TABLE cases.reason (
	id int8 DEFAULT nextval('cases.reason_id'::regclass) NOT NULL,
	"name" text NOT NULL,
	description text NOT NULL,
	created_at timestamp DEFAULT timezone('utc'::text, now()) NOT NULL,
	updated_at timestamp DEFAULT timezone('utc'::text, now()) NOT NULL,
	created_by int8 NULL,
	updated_by int8 NULL,
	close_reason_id int8 NOT NULL,
	dc int8 NOT NULL,
	CONSTRAINT reason_fk UNIQUE (id, dc),
	CONSTRAINT reason_pk PRIMARY KEY (id)
);
CREATE INDEX reason_id_dc ON cases.reason USING btree (dc);
CREATE INDEX reason_source ON cases.reason USING btree (close_reason_id);


-- cases.reason foreign keys

ALTER TABLE cases.reason ADD CONSTRAINT reason_close_reason_id_fk FOREIGN KEY (close_reason_id) REFERENCES cases.close_reason(id) ON DELETE CASCADE;
ALTER TABLE cases.reason ADD CONSTRAINT reason_created_by_fk FOREIGN KEY (created_by) REFERENCES directory.wbt_user(id) ON DELETE SET NULL DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE cases.reason ADD CONSTRAINT reason_created_dc_fk FOREIGN KEY (created_by,dc) REFERENCES directory.wbt_user(id,dc) DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE cases.reason ADD CONSTRAINT reason_domain_fk FOREIGN KEY (close_reason_id,dc) REFERENCES cases.close_reason(id,dc) ON DELETE CASCADE DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE cases.reason ADD CONSTRAINT reason_updated_by_fk FOREIGN KEY (updated_by) REFERENCES directory.wbt_user(id) ON DELETE SET NULL DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE cases.reason ADD CONSTRAINT reason_updated_dc_fk FOREIGN KEY (updated_by,dc) REFERENCES directory.wbt_user(id,dc) DEFERRABLE INITIALLY DEFERRED;

-- cases.service_catalog definition

-- Drop table

-- DROP TABLE cases.service_catalog;

CREATE TABLE cases.service_catalog (
	id int8 DEFAULT nextval('cases.service_id_seq'::regclass) NOT NULL,
	root_id int8 NULL,
	description text NULL,
	code varchar(50) NULL,
	prefix varchar(20) NULL,
	state bool DEFAULT true NOT NULL,
	assignee_id int8 NULL,
	created_at timestamp DEFAULT timezone('utc'::text, now()) NOT NULL,
	updated_at timestamp DEFAULT timezone('utc'::text, now()) NOT NULL,
	created_by int8 NULL,
	updated_by int8 NULL,
	dc int8 NOT NULL,
	sla_id int8 NULL,
	status_id int8 NULL,
	close_reason_id int8 NULL,
	group_id int4 NULL,
	"name" text NOT NULL,
	CONSTRAINT service_fk UNIQUE (id, dc),
	CONSTRAINT service_pk PRIMARY KEY (id)
);
CREATE INDEX service_dc ON cases.service_catalog USING btree (dc);


-- cases.service_catalog foreign keys

ALTER TABLE cases.service_catalog ADD CONSTRAINT service_created_dc_fk FOREIGN KEY (created_by,dc) REFERENCES directory.wbt_user(id,dc) DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE cases.service_catalog ADD CONSTRAINT service_created_id_fk FOREIGN KEY (created_by) REFERENCES directory.wbt_user(id) ON DELETE SET NULL DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE cases.service_catalog ADD CONSTRAINT service_domain_fk FOREIGN KEY (dc) REFERENCES directory.wbt_domain(dc) ON DELETE CASCADE;
ALTER TABLE cases.service_catalog ADD CONSTRAINT service_updated_dc_fk FOREIGN KEY (updated_by,dc) REFERENCES directory.wbt_user(id,dc) DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE cases.service_catalog ADD CONSTRAINT service_updated_id_fk FOREIGN KEY (updated_by) REFERENCES directory.wbt_user(id) DEFERRABLE INITIALLY DEFERRED;

-- cases.skill_catalog definition

-- Drop table

-- DROP TABLE cases.skill_catalog;

CREATE TABLE cases.skill_catalog (
	id int8 DEFAULT nextval('cases.skill_service_id_seq'::regclass) NOT NULL,
	created_at timestamp DEFAULT timezone('utc'::text, now()) NOT NULL,
	updated_at timestamp DEFAULT timezone('utc'::text, now()) NOT NULL,
	created_by int8 NULL,
	updated_by int8 NOT NULL,
	catalog_id int8 NOT NULL,
	skill_id int8 NOT NULL,
	dc int8 NOT NULL,
	CONSTRAINT skill_service_fk UNIQUE (id, dc),
	CONSTRAINT skill_service_pk PRIMARY KEY (id)
);
CREATE UNIQUE INDEX idx_catalog_id_skill_id ON cases.skill_catalog USING btree (catalog_id, skill_id);
CREATE UNIQUE INDEX idx_service_skill ON cases.skill_catalog USING btree (skill_id, catalog_id);
CREATE INDEX skill_service_dc ON cases.skill_catalog USING btree (dc);
CREATE INDEX skill_service_source ON cases.skill_catalog USING btree (catalog_id);


-- cases.skill_catalog foreign keys

ALTER TABLE cases.skill_catalog ADD CONSTRAINT skill_service_created_by_fk FOREIGN KEY (created_by) REFERENCES directory.wbt_user(id) ON DELETE SET NULL DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE cases.skill_catalog ADD CONSTRAINT skill_service_created_dc_fk FOREIGN KEY (created_by,dc) REFERENCES directory.wbt_user(id,dc) DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE cases.skill_catalog ADD CONSTRAINT skill_service_domain_fk FOREIGN KEY (dc) REFERENCES directory.wbt_domain(dc) ON DELETE CASCADE;
ALTER TABLE cases.skill_catalog ADD CONSTRAINT skill_service_service_id_fk FOREIGN KEY (catalog_id) REFERENCES cases.service_catalog(id) ON DELETE CASCADE;
ALTER TABLE cases.skill_catalog ADD CONSTRAINT skill_service_updated_by_fk FOREIGN KEY (updated_by) REFERENCES directory.wbt_user(id) ON DELETE SET NULL DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE cases.skill_catalog ADD CONSTRAINT skill_service_updated_dc_fk FOREIGN KEY (updated_by,dc) REFERENCES directory.wbt_user(id,dc) DEFERRABLE INITIALLY DEFERRED;

-- cases.sla definition

-- Drop table

-- DROP TABLE cases.sla;

CREATE TABLE cases.sla (
	id int8 DEFAULT nextval('cases.sla_id'::regclass) NOT NULL,
	"name" text NOT NULL,
	description text NULL,
	created_at timestamp DEFAULT timezone('utc'::text, now()) NOT NULL,
	updated_at timestamp DEFAULT timezone('utc'::text, now()) NOT NULL,
	created_by int8 NULL,
	updated_by int8 NULL,
	dc int8 NOT NULL,
	valid_from timestamp DEFAULT timezone('utc'::text, now()) NULL,
	valid_to timestamp DEFAULT timezone('utc'::text, now()) NULL,
	reaction_time_hours int4 DEFAULT 0 NOT NULL,
	reaction_time_minutes int4 DEFAULT 0 NULL,
	resolution_time_hours int4 DEFAULT 0 NOT NULL,
	resolution_time_minutes int4 DEFAULT 0 NOT NULL,
	calendar_id int8 NOT NULL,
	CONSTRAINT sla_fk UNIQUE (id, dc),
	CONSTRAINT sla_pk PRIMARY KEY (id)
);
CREATE INDEX sla_dc ON cases.sla USING btree (dc);


-- cases.sla foreign keys

ALTER TABLE cases.sla ADD CONSTRAINT sla_created_dc_fk FOREIGN KEY (created_by,dc) REFERENCES directory.wbt_user(id,dc) DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE cases.sla ADD CONSTRAINT sla_created_id_fk FOREIGN KEY (created_by) REFERENCES directory.wbt_user(id) ON DELETE SET NULL DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE cases.sla ADD CONSTRAINT sla_domain_fk FOREIGN KEY (dc) REFERENCES directory.wbt_domain(dc) ON DELETE CASCADE;
ALTER TABLE cases.sla ADD CONSTRAINT sla_updated_dc_fk FOREIGN KEY (updated_by,dc) REFERENCES directory.wbt_user(id,dc) DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE cases.sla ADD CONSTRAINT sla_updated_id_fk FOREIGN KEY (updated_by) REFERENCES directory.wbt_user(id) DEFERRABLE INITIALLY DEFERRED;

-- cases.sla_condition definition

-- Drop table

-- DROP TABLE cases.sla_condition;

CREATE TABLE cases.sla_condition (
	id int8 DEFAULT nextval('cases.sla_condition_id'::regclass) NOT NULL,
	"name" text NOT NULL,
	created_at timestamp DEFAULT timezone('utc'::text, now()) NOT NULL,
	updated_at timestamp DEFAULT timezone('utc'::text, now()) NOT NULL,
	created_by int8 NOT NULL,
	updated_by int8 NULL,
	dc int8 NOT NULL,
	reaction_time_hours int4 DEFAULT 0 NOT NULL,
	reaction_time_minutes int4 DEFAULT 0 NOT NULL,
	resolution_time_hours int4 DEFAULT 0 NOT NULL,
	resolution_time_minutes int4 DEFAULT 0 NOT NULL,
	sla_id int8 NOT NULL,
	CONSTRAINT sla_condition_fk UNIQUE (id, dc),
	CONSTRAINT sla_condition_pk PRIMARY KEY (id)
);
CREATE INDEX sla_condition_dc ON cases.sla_condition USING btree (dc);
CREATE INDEX sla_condition_source ON cases.sla_condition USING btree (sla_id);


-- cases.sla_condition foreign keys

ALTER TABLE cases.sla_condition ADD CONSTRAINT sla_condition_created_dc_fk FOREIGN KEY (created_by,dc) REFERENCES directory.wbt_user(id,dc) DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE cases.sla_condition ADD CONSTRAINT sla_condition_created_id_fk FOREIGN KEY (created_by) REFERENCES directory.wbt_user(id) ON DELETE SET NULL DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE cases.sla_condition ADD CONSTRAINT sla_condition_domain_fk FOREIGN KEY (dc) REFERENCES directory.wbt_domain(dc) ON DELETE CASCADE;
ALTER TABLE cases.sla_condition ADD CONSTRAINT sla_condition_sla_id_fk FOREIGN KEY (sla_id) REFERENCES cases.sla(id) ON DELETE CASCADE;
ALTER TABLE cases.sla_condition ADD CONSTRAINT sla_condition_updated_dc_fk FOREIGN KEY (updated_by,dc) REFERENCES directory.wbt_user(id,dc) DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE cases.sla_condition ADD CONSTRAINT sla_condition_updated_id_fk FOREIGN KEY (updated_by) REFERENCES directory.wbt_user(id) DEFERRABLE INITIALLY DEFERRED;

-- cases.status definition

-- Drop table

-- DROP TABLE cases.status;

CREATE TABLE cases.status (
	id int8 DEFAULT nextval('cases.status_id'::regclass) NOT NULL,
	"name" text NOT NULL,
	description text NULL,
	created_at timestamp DEFAULT timezone('utc'::text, now()) NOT NULL,
	updated_at timestamp DEFAULT timezone('utc'::text, now()) NOT NULL,
	created_by int8 NULL,
	updated_by int8 NULL,
	dc int8 NOT NULL,
	CONSTRAINT status_fk UNIQUE (id, dc),
	CONSTRAINT status_pk PRIMARY KEY (id)
);
CREATE INDEX status_dc ON cases.status USING btree (dc);


-- cases.status foreign keys

ALTER TABLE cases.status ADD CONSTRAINT status_created_dc_fk FOREIGN KEY (created_by,dc) REFERENCES directory.wbt_user(id,dc) DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE cases.status ADD CONSTRAINT status_created_id_fk FOREIGN KEY (created_by) REFERENCES directory.wbt_user(id) ON DELETE SET NULL DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE cases.status ADD CONSTRAINT status_domain_fk FOREIGN KEY (dc) REFERENCES directory.wbt_domain(dc) ON DELETE CASCADE;
ALTER TABLE cases.status ADD CONSTRAINT status_updated_dc_fk FOREIGN KEY (updated_by,dc) REFERENCES directory.wbt_user(id,dc) DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE cases.status ADD CONSTRAINT status_updated_id_fk FOREIGN KEY (updated_by) REFERENCES directory.wbt_user(id) DEFERRABLE INITIALLY DEFERRED;

-- cases.status_condition definition

-- Drop table

-- DROP TABLE cases.status_condition;

CREATE TABLE cases.status_condition (
	id int8 DEFAULT nextval('cases.status_condition_id'::regclass) NOT NULL,
	"name" text NOT NULL,
	description text NOT NULL,
	created_at timestamp DEFAULT timezone('utc'::text, now()) NOT NULL,
	updated_at timestamp DEFAULT timezone('utc'::text, now()) NOT NULL,
	created_by int8 NULL,
	updated_by int8 NULL,
	status_id int8 NOT NULL,
	initial bool NOT NULL,
	"final" bool NOT NULL,
	dc int8 NOT NULL,
	CONSTRAINT status_condition_fk UNIQUE (id, dc),
	CONSTRAINT status_condition_pk PRIMARY KEY (id)
);
CREATE INDEX status_condition_dc ON cases.status_condition USING btree (dc);
CREATE INDEX status_condition_source ON cases.status_condition USING btree (status_id);


-- cases.status_condition foreign keys

ALTER TABLE cases.status_condition ADD CONSTRAINT status_condition_created_by_fk FOREIGN KEY (created_by) REFERENCES directory.wbt_user(id) ON DELETE SET NULL DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE cases.status_condition ADD CONSTRAINT status_condition_created_dc_fk FOREIGN KEY (created_by,dc) REFERENCES directory.wbt_user(id,dc) DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE cases.status_condition ADD CONSTRAINT status_condition_domain_fk FOREIGN KEY (status_id,dc) REFERENCES cases.status(id,dc) ON DELETE CASCADE DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE cases.status_condition ADD CONSTRAINT status_condition_status_id_fk FOREIGN KEY (status_id) REFERENCES cases.status(id) ON DELETE CASCADE;
ALTER TABLE cases.status_condition ADD CONSTRAINT status_condition_updated_by_fk FOREIGN KEY (updated_by) REFERENCES directory.wbt_user(id) ON DELETE SET NULL DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE cases.status_condition ADD CONSTRAINT status_condition_updated_dc_fk FOREIGN KEY (updated_by,dc) REFERENCES directory.wbt_user(id,dc) DEFERRABLE INITIALLY DEFERRED;

-- cases.team_catalog definition

-- Drop table

-- DROP TABLE cases.team_catalog;

CREATE TABLE cases.team_catalog (
	id int8 DEFAULT nextval('cases.team_service_id_seq'::regclass) NOT NULL,
	created_at timestamp DEFAULT timezone('utc'::text, now()) NOT NULL,
	updated_at timestamp DEFAULT timezone('utc'::text, now()) NOT NULL,
	created_by int8 NULL,
	updated_by int8 NULL,
	catalog_id int8 NOT NULL,
	team_id int8 NOT NULL,
	dc int8 NOT NULL,
	CONSTRAINT team_service_fk UNIQUE (id, dc),
	CONSTRAINT team_service_pk PRIMARY KEY (id)
);
CREATE UNIQUE INDEX idx_catalog_id_team_id ON cases.team_catalog USING btree (catalog_id, team_id);
CREATE INDEX idx_service_team ON cases.team_catalog USING btree (catalog_id, team_id);
CREATE INDEX team_service_dc ON cases.team_catalog USING btree (dc);
CREATE INDEX team_service_source ON cases.team_catalog USING btree (catalog_id);


-- cases.team_catalog foreign keys

ALTER TABLE cases.team_catalog ADD CONSTRAINT team_service_created_by_fk FOREIGN KEY (created_by) REFERENCES directory.wbt_user(id) ON DELETE SET NULL DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE cases.team_catalog ADD CONSTRAINT team_service_created_dc_fk FOREIGN KEY (created_by,dc) REFERENCES directory.wbt_user(id,dc) DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE cases.team_catalog ADD CONSTRAINT team_service_domain_fk FOREIGN KEY (dc) REFERENCES directory.wbt_domain(dc) ON DELETE CASCADE;
ALTER TABLE cases.team_catalog ADD CONSTRAINT team_service_service_id_fk FOREIGN KEY (catalog_id) REFERENCES cases.service_catalog(id) ON DELETE CASCADE;
ALTER TABLE cases.team_catalog ADD CONSTRAINT team_service_updated_by_fk FOREIGN KEY (updated_by) REFERENCES directory.wbt_user(id) ON DELETE SET NULL DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE cases.team_catalog ADD CONSTRAINT team_service_updated_dc_fk FOREIGN KEY (updated_by,dc) REFERENCES directory.wbt_user(id,dc) DEFERRABLE INITIALLY DEFERRED;