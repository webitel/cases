-- cases."case" definition

-- Drop table

-- DROP TABLE cases."case";

CREATE TABLE cases."case" (
	id int8 DEFAULT nextval('cases.case_id'::regclass) NOT NULL,
	dc int8 NOT NULL,
	"name" varchar NOT NULL,
	subject varchar NOT NULL,
	description varchar NULL,
	ver int4 DEFAULT 0 NOT NULL,
	created_by int8 NULL,
	created_at timestamp DEFAULT timezone('utc'::text, now()) NOT NULL,
	updated_by int8 NULL,
	updated_at timestamp DEFAULT timezone('utc'::text, now()) NOT NULL,
	status int8 NOT NULL,
	close_reason_group int8 NOT NULL,
	assignee int8 NULL,
	reporter int8 NULL,
	impacted int8 NULL,
	contact_group int8 NULL,
	priority int8 NOT NULL,
	"source" int8 NOT NULL,
	status_condition int8 NULL,
	close_result varchar NULL,
	close_reason int8 NULL,
	rating int8 NULL,
	rating_comment varchar NULL,
	service int8 NOT NULL,
	sla int8 NOT NULL,
	planned_reaction_at timestamp DEFAULT timezone('utc'::text, now()) NULL,
	planned_resolve_at timestamp DEFAULT timezone('utc'::text, now()) NULL,
	sla_condition_id int8 NULL,
	contact_info text NULL,
	reacted_at timestamp NULL,
	resolved_at timestamp NULL,
	CONSTRAINT case_fk UNIQUE (id, dc),
	CONSTRAINT case_pk PRIMARY KEY (id)
);
CREATE INDEX case_dc ON cases."case" USING btree (dc);

-- Table Triggers

create trigger tg_case_rbac after
insert
    on
    cases."case" for each row execute function directory.tg_obj_default_rbac('cases');
create trigger trigger_update_case_timings before
insert
    or
update
    on
    cases."case" for each row execute function cases.update_case_timings();


-- cases.case_acl definition

-- Drop table

-- DROP TABLE cases.case_acl;

CREATE TABLE cases.case_acl (
	dc int8 NOT NULL,
	grantor int8 NULL,
	subject int8 NOT NULL,
	"object" int8 NOT NULL,
	"access" int2 NULL
);
CREATE INDEX case_acl_grantor_index ON cases.case_acl USING btree (grantor);
CREATE UNIQUE INDEX case_acl_object_subject_pk ON cases.case_acl USING btree (object, subject);
CREATE UNIQUE INDEX case_acl_subject_object_uindex ON cases.case_acl USING btree (subject, object);


-- cases.case_comment definition

-- Drop table

-- DROP TABLE cases.case_comment;

CREATE TABLE cases.case_comment (
	id int8 DEFAULT nextval('cases.case_comment_id'::regclass) NOT NULL,
	ver int4 DEFAULT 0 NOT NULL,
	created_by int8 NULL,
	created_at timestamp DEFAULT timezone('utc'::text, now()) NOT NULL,
	updated_by int8 NULL,
	updated_at timestamp DEFAULT timezone('utc'::text, now()) NOT NULL,
	"comment" varchar NOT NULL,
	case_id int8 NOT NULL,
	dc int8 NOT NULL,
	CONSTRAINT case_comment_id_case_id_unique UNIQUE (id, case_id),
	CONSTRAINT comment_case_fk UNIQUE (id, dc),
	CONSTRAINT comment_case_pk PRIMARY KEY (id)
);
CREATE INDEX comment_case_dc ON cases.case_comment USING btree (dc);
CREATE INDEX idx_link_case_id ON cases.case_comment USING btree (case_id);

-- Table Triggers

create trigger tg_case_comment_rbac after
insert
    on
    cases.case_comment for each row execute function directory.tg_obj_default_rbac('case_comments');


-- cases.case_comment_acl definition

-- Drop table

-- DROP TABLE cases.case_comment_acl;

CREATE TABLE cases.case_comment_acl (
	dc int8 NOT NULL,
	grantor int8 NULL,
	subject int8 NOT NULL,
	"object" int8 NOT NULL,
	"access" int2 NULL
);
CREATE INDEX case_comment_acl_grantor_uindex ON cases.case_comment_acl USING btree (grantor);
CREATE UNIQUE INDEX case_comment_acl_object_subject_pk ON cases.case_comment_acl USING btree (object, subject);
CREATE UNIQUE INDEX case_comment_acl_subject_object_uindex ON cases.case_comment_acl USING btree (subject, object);


-- cases.case_communication definition

-- Drop table

-- DROP TABLE cases.case_communication;

CREATE TABLE cases.case_communication (
	id bigserial NOT NULL,
	ver int4 DEFAULT 0 NOT NULL,
	case_id int8 NOT NULL,
	communication_type int4 NOT NULL, -- type of communication case connected to
	communication_id text NOT NULL, -- Id of communication case connected to
	dc int8 NOT NULL, -- Domain component
	created_by int8 NULL,
	created_at timestamp DEFAULT now() NOT NULL
);
CREATE UNIQUE INDEX unique_event_uindex ON cases.case_communication USING btree (communication_id, communication_type, case_id);
COMMENT ON TABLE cases.case_communication IS 'Table connects case with other possible communications like chats, emails, calls';

-- Column comments

COMMENT ON COLUMN cases.case_communication.communication_type IS 'type of communication case connected to';
COMMENT ON COLUMN cases.case_communication.communication_id IS 'Id of communication case connected to';
COMMENT ON COLUMN cases.case_communication.dc IS 'Domain component';


-- cases.case_link definition

-- Drop table

-- DROP TABLE cases.case_link;

CREATE TABLE cases.case_link (
	id bigserial NOT NULL,
	ver int4 DEFAULT 0 NOT NULL,
	created_by int8 NULL,
	created_at timestamp DEFAULT timezone('utc'::text, now()) NOT NULL,
	updated_by int8 NULL,
	updated_at timestamp DEFAULT timezone('utc'::text, now()) NOT NULL,
	"name" varchar NULL,
	url varchar NOT NULL,
	case_id int8 NOT NULL,
	dc int8 NOT NULL,
	CONSTRAINT case_link_fk UNIQUE (id, dc),
	CONSTRAINT case_link_id_case_id_unique UNIQUE (id, case_id),
	CONSTRAINT link_case_pk PRIMARY KEY (id)
);
CREATE INDEX idx_comment_case_id ON cases.case_link USING btree (case_id);
CREATE INDEX link_case_dc ON cases.case_link USING btree (dc);


-- cases.close_reason definition

-- Drop table

-- DROP TABLE cases.close_reason;

CREATE TABLE cases.close_reason (
	id int8 DEFAULT nextval('cases.reason_id'::regclass) NOT NULL,
	"name" text NOT NULL,
	description text NULL,
	created_at timestamp DEFAULT timezone('utc'::text, now()) NOT NULL,
	updated_at timestamp DEFAULT timezone('utc'::text, now()) NOT NULL,
	created_by int8 NULL,
	updated_by int8 NULL,
	close_reason_id int8 NOT NULL,
	dc int8 NOT NULL,
	CONSTRAINT reason_fk UNIQUE (id, dc),
	CONSTRAINT reason_pk PRIMARY KEY (id)
);
CREATE INDEX reason_id_dc ON cases.close_reason USING btree (dc);
CREATE INDEX reason_source ON cases.close_reason USING btree (close_reason_id);


-- cases.close_reason_group definition

-- Drop table

-- DROP TABLE cases.close_reason_group;

CREATE TABLE cases.close_reason_group (
	id int8 DEFAULT nextval('cases.close_reason_id'::regclass) NOT NULL,
	"name" text NOT NULL,
	description text NULL,
	created_at timestamp DEFAULT timezone('utc'::text, now()) NOT NULL,
	updated_at timestamp DEFAULT timezone('utc'::text, now()) NOT NULL,
	created_by int8 NULL,
	updated_by int8 NULL,
	dc int8 NOT NULL,
	CONSTRAINT close_reason_fk UNIQUE (id, dc),
	CONSTRAINT close_reason_pk PRIMARY KEY (id)
);
CREATE INDEX close_reason_dc ON cases.close_reason_group USING btree (dc);


-- cases.priority definition

-- Drop table

-- DROP TABLE cases.priority;

CREATE TABLE cases.priority (
	id int8 DEFAULT nextval('cases.priority_id'::regclass) NOT NULL,
	"name" text NOT NULL,
	description text NULL,
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


-- cases.related_case definition

-- Drop table

-- DROP TABLE cases.related_case;

CREATE TABLE cases.related_case (
	id serial4 NOT NULL,
	dc int8 NOT NULL,
	ver int4 DEFAULT 0 NOT NULL,
	created_by int8 NULL,
	created_at timestamp DEFAULT timezone('utc'::text, now()) NOT NULL,
	updated_by int8 NULL,
	updated_at timestamp DEFAULT timezone('utc'::text, now()) NOT NULL,
	relation_type int2 NOT NULL,
	primary_case_id int8 NOT NULL,
	related_case_id int8 NOT NULL,
	CONSTRAINT related_case_fk UNIQUE (id, dc),
	CONSTRAINT related_case_pk PRIMARY KEY (id)
);
CREATE INDEX idx_related_case_primary_case ON cases.related_case USING btree (primary_case_id);
CREATE INDEX idx_related_case_related_case ON cases.related_case USING btree (related_case_id);
CREATE INDEX related_case_dc ON cases.related_case USING btree (dc);
CREATE UNIQUE INDEX unique_related_cases_relation ON cases.related_case USING btree (LEAST(primary_case_id, related_case_id), GREATEST(primary_case_id, related_case_id), relation_type);


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
	close_reason_group_id int8 NULL,
	group_id int4 NULL,
	"name" text NOT NULL,
	catalog_id int8 NULL,
	CONSTRAINT chk_root_id_close_reason_id CHECK ((NOT ((root_id IS NULL) AND (close_reason_group_id IS NULL)))),
	CONSTRAINT service_fk UNIQUE (id, dc),
	CONSTRAINT service_pk PRIMARY KEY (id)
);
CREATE INDEX service_dc ON cases.service_catalog USING btree (dc);


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
	reaction_time int4 DEFAULT 0 NOT NULL,
	resolution_time int4 DEFAULT 0 NOT NULL,
	calendar_id int8 NOT NULL,
	CONSTRAINT sla_fk UNIQUE (id, dc),
	CONSTRAINT sla_pk PRIMARY KEY (id)
);
CREATE INDEX sla_dc ON cases.sla USING btree (dc);


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
	reaction_time int4 DEFAULT 0 NOT NULL,
	resolution_time int4 DEFAULT 0 NOT NULL,
	sla_id int8 NOT NULL,
	CONSTRAINT sla_condition_fk UNIQUE (id, dc),
	CONSTRAINT sla_condition_pk PRIMARY KEY (id)
);
CREATE INDEX sla_condition_dc ON cases.sla_condition USING btree (dc);
CREATE INDEX sla_condition_source ON cases.sla_condition USING btree (sla_id);


-- cases."source" definition

-- Drop table

-- DROP TABLE cases."source";

CREATE TABLE cases."source" (
	id int8 DEFAULT nextval('cases.appeal_id'::regclass) NOT NULL,
	"name" text NOT NULL,
	description text NULL,
	created_at timestamp DEFAULT timezone('utc'::text, now()) NOT NULL,
	updated_at timestamp DEFAULT timezone('utc'::text, now()) NOT NULL,
	created_by int8 NULL,
	updated_by int8 NULL,
	"type" text NOT NULL,
	dc int8 NOT NULL,
	CONSTRAINT source_fk UNIQUE (id, dc),
	CONSTRAINT source_pk PRIMARY KEY (id)
);
CREATE INDEX source_dc ON cases.source USING btree (dc);


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


-- cases.status_condition definition

-- Drop table

-- DROP TABLE cases.status_condition;

CREATE TABLE cases.status_condition (
	id int8 DEFAULT nextval('cases.status_condition_id'::regclass) NOT NULL,
	"name" text NOT NULL,
	description text NULL,
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


-- cases."case" foreign keys

ALTER TABLE cases."case" ADD CONSTRAINT case_assignee_fk FOREIGN KEY (assignee) REFERENCES contacts.contact(id) ON DELETE SET NULL DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE cases."case" ADD CONSTRAINT case_close_reason_fk FOREIGN KEY (close_reason) REFERENCES cases.close_reason(id) ON DELETE SET NULL DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE cases."case" ADD CONSTRAINT case_close_reason_group_fk FOREIGN KEY (close_reason_group) REFERENCES cases.close_reason_group(id) ON DELETE RESTRICT DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE cases."case" ADD CONSTRAINT case_created_dc_fk FOREIGN KEY (created_by,dc) REFERENCES directory.wbt_user(id,dc) DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE cases."case" ADD CONSTRAINT case_created_id_fk FOREIGN KEY (created_by) REFERENCES directory.wbt_user(id) ON DELETE SET NULL DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE cases."case" ADD CONSTRAINT case_domain_fk FOREIGN KEY (dc) REFERENCES directory.wbt_domain(dc) ON DELETE CASCADE DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE cases."case" ADD CONSTRAINT case_group_fk FOREIGN KEY (contact_group) REFERENCES contacts."group"(id) ON DELETE SET NULL DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE cases."case" ADD CONSTRAINT case_impacted_fk FOREIGN KEY (impacted) REFERENCES contacts.contact(id) ON DELETE SET NULL DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE cases."case" ADD CONSTRAINT case_priority_fk FOREIGN KEY (priority) REFERENCES cases.priority(id) ON DELETE RESTRICT DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE cases."case" ADD CONSTRAINT case_reporter_fk FOREIGN KEY (reporter) REFERENCES contacts.contact(id) ON DELETE SET NULL DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE cases."case" ADD CONSTRAINT case_service_fk FOREIGN KEY (service) REFERENCES cases.service_catalog(id) ON DELETE RESTRICT DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE cases."case" ADD CONSTRAINT case_sla_fk FOREIGN KEY (sla) REFERENCES cases.sla(id) ON DELETE RESTRICT DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE cases."case" ADD CONSTRAINT case_source_fk FOREIGN KEY ("source") REFERENCES cases."source"(id) ON DELETE RESTRICT DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE cases."case" ADD CONSTRAINT case_status_condition_fk FOREIGN KEY (status_condition) REFERENCES cases.status_condition(id) ON DELETE SET NULL DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE cases."case" ADD CONSTRAINT case_status_fk FOREIGN KEY (status) REFERENCES cases.status(id) ON DELETE RESTRICT DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE cases."case" ADD CONSTRAINT case_updated_dc_fk FOREIGN KEY (updated_by,dc) REFERENCES directory.wbt_user(id,dc) DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE cases."case" ADD CONSTRAINT case_updated_id_fk FOREIGN KEY (updated_by) REFERENCES directory.wbt_user(id) ON DELETE SET NULL DEFERRABLE INITIALLY DEFERRED;


-- cases.case_acl foreign keys

ALTER TABLE cases.case_acl ADD CONSTRAINT case_acl_domain_fk FOREIGN KEY (dc) REFERENCES directory.wbt_domain(dc) ON DELETE CASCADE;
ALTER TABLE cases.case_acl ADD CONSTRAINT case_acl_grantor_fk FOREIGN KEY (grantor,dc) REFERENCES directory.wbt_auth(id,dc) DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE cases.case_acl ADD CONSTRAINT case_acl_grantor_id_fk FOREIGN KEY (grantor) REFERENCES directory.wbt_auth(id) ON DELETE SET NULL;
ALTER TABLE cases.case_acl ADD CONSTRAINT case_acl_object_fk FOREIGN KEY ("object",dc) REFERENCES cases."case"(id,dc) ON DELETE CASCADE DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE cases.case_acl ADD CONSTRAINT case_acl_subject_fk FOREIGN KEY (subject,dc) REFERENCES directory.wbt_auth(id,dc) ON DELETE CASCADE;


-- cases.case_comment foreign keys

ALTER TABLE cases.case_comment ADD CONSTRAINT comment_case_case_fk FOREIGN KEY (case_id) REFERENCES cases."case"(id) ON DELETE CASCADE DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE cases.case_comment ADD CONSTRAINT comment_case_created_dc_fk FOREIGN KEY (created_by,dc) REFERENCES directory.wbt_user(id,dc) DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE cases.case_comment ADD CONSTRAINT comment_case_created_id_fk FOREIGN KEY (created_by) REFERENCES directory.wbt_user(id) ON DELETE SET NULL DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE cases.case_comment ADD CONSTRAINT comment_case_updated_dc_fk FOREIGN KEY (updated_by,dc) REFERENCES directory.wbt_user(id,dc) DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE cases.case_comment ADD CONSTRAINT comment_case_updated_id_fk FOREIGN KEY (updated_by) REFERENCES directory.wbt_user(id) ON DELETE SET NULL DEFERRABLE INITIALLY DEFERRED;


-- cases.case_comment_acl foreign keys

ALTER TABLE cases.case_comment_acl ADD CONSTRAINT case_comment_acl_domain_fk FOREIGN KEY (dc) REFERENCES directory.wbt_domain(dc) ON DELETE CASCADE;
ALTER TABLE cases.case_comment_acl ADD CONSTRAINT case_comment_acl_grantor_fk FOREIGN KEY (grantor,dc) REFERENCES directory.wbt_auth(id,dc) DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE cases.case_comment_acl ADD CONSTRAINT case_comment_acl_grantor_id_fk FOREIGN KEY (grantor) REFERENCES directory.wbt_auth(id) ON DELETE SET NULL;
ALTER TABLE cases.case_comment_acl ADD CONSTRAINT case_comment_acl_object_fk FOREIGN KEY ("object",dc) REFERENCES cases.case_comment(id,dc) ON DELETE CASCADE DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE cases.case_comment_acl ADD CONSTRAINT case_comment_acl_subject_fk FOREIGN KEY (subject,dc) REFERENCES directory.wbt_auth(id,dc) ON DELETE CASCADE;


-- cases.case_communication foreign keys

ALTER TABLE cases.case_communication ADD CONSTRAINT case_communication_case_id_fk FOREIGN KEY (case_id) REFERENCES cases."case"(id) ON DELETE CASCADE;
ALTER TABLE cases.case_communication ADD CONSTRAINT case_communication_wbt_domain_dc_fk FOREIGN KEY (dc) REFERENCES directory.wbt_domain(dc) ON DELETE CASCADE;
ALTER TABLE cases.case_communication ADD CONSTRAINT case_communication_wbt_user_id_fk FOREIGN KEY (created_by) REFERENCES directory.wbt_user(id);


-- cases.case_link foreign keys

ALTER TABLE cases.case_link ADD CONSTRAINT link_case_case_fk FOREIGN KEY (case_id) REFERENCES cases."case"(id) ON DELETE CASCADE DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE cases.case_link ADD CONSTRAINT link_case_created_dc_fk FOREIGN KEY (created_by,dc) REFERENCES directory.wbt_user(id,dc) DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE cases.case_link ADD CONSTRAINT link_case_created_id_fk FOREIGN KEY (created_by) REFERENCES directory.wbt_user(id) ON DELETE SET NULL DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE cases.case_link ADD CONSTRAINT link_case_updated_dc_fk FOREIGN KEY (updated_by,dc) REFERENCES directory.wbt_user(id,dc) DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE cases.case_link ADD CONSTRAINT link_case_updated_id_fk FOREIGN KEY (updated_by) REFERENCES directory.wbt_user(id) ON DELETE SET NULL DEFERRABLE INITIALLY DEFERRED;


-- cases.close_reason foreign keys

ALTER TABLE cases.close_reason ADD CONSTRAINT reason_close_reason_id_fk FOREIGN KEY (close_reason_id) REFERENCES cases.close_reason_group(id) ON DELETE CASCADE;
ALTER TABLE cases.close_reason ADD CONSTRAINT reason_created_by_fk FOREIGN KEY (created_by) REFERENCES directory.wbt_user(id) ON DELETE SET NULL DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE cases.close_reason ADD CONSTRAINT reason_created_dc_fk FOREIGN KEY (created_by,dc) REFERENCES directory.wbt_user(id,dc) DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE cases.close_reason ADD CONSTRAINT reason_domain_fk FOREIGN KEY (close_reason_id,dc) REFERENCES cases.close_reason_group(id,dc) ON DELETE CASCADE DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE cases.close_reason ADD CONSTRAINT reason_updated_by_fk FOREIGN KEY (updated_by) REFERENCES directory.wbt_user(id) ON DELETE SET NULL DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE cases.close_reason ADD CONSTRAINT reason_updated_dc_fk FOREIGN KEY (updated_by,dc) REFERENCES directory.wbt_user(id,dc) DEFERRABLE INITIALLY DEFERRED;


-- cases.close_reason_group foreign keys

ALTER TABLE cases.close_reason_group ADD CONSTRAINT close_reason_created_dc_fk FOREIGN KEY (created_by,dc) REFERENCES directory.wbt_user(id,dc) DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE cases.close_reason_group ADD CONSTRAINT close_reason_created_id_fk FOREIGN KEY (created_by) REFERENCES directory.wbt_user(id) ON DELETE SET NULL DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE cases.close_reason_group ADD CONSTRAINT close_reason_domain_fk FOREIGN KEY (dc) REFERENCES directory.wbt_domain(dc) ON DELETE CASCADE;
ALTER TABLE cases.close_reason_group ADD CONSTRAINT close_reason_updated_dc_fk FOREIGN KEY (updated_by,dc) REFERENCES directory.wbt_user(id,dc) DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE cases.close_reason_group ADD CONSTRAINT close_reason_updated_id_fk FOREIGN KEY (updated_by) REFERENCES directory.wbt_user(id) DEFERRABLE INITIALLY DEFERRED;


-- cases.priority foreign keys

ALTER TABLE cases.priority ADD CONSTRAINT priority_created_dc_fk FOREIGN KEY (created_by,dc) REFERENCES directory.wbt_user(id,dc) DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE cases.priority ADD CONSTRAINT priority_created_id_fk FOREIGN KEY (created_by) REFERENCES directory.wbt_user(id) ON DELETE SET NULL DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE cases.priority ADD CONSTRAINT priority_domain_fk FOREIGN KEY (dc) REFERENCES directory.wbt_domain(dc) ON DELETE CASCADE;
ALTER TABLE cases.priority ADD CONSTRAINT priority_updated_dc_fk FOREIGN KEY (updated_by,dc) REFERENCES directory.wbt_user(id,dc) DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE cases.priority ADD CONSTRAINT priority_updated_id_fk FOREIGN KEY (updated_by) REFERENCES directory.wbt_user(id) DEFERRABLE INITIALLY DEFERRED;


-- cases.priority_sla_condition foreign keys

ALTER TABLE cases.priority_sla_condition ADD CONSTRAINT priority_sla_condition_created_by_fk FOREIGN KEY (created_by) REFERENCES directory.wbt_user(id) ON DELETE SET NULL DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE cases.priority_sla_condition ADD CONSTRAINT priority_sla_condition_created_dc_fk FOREIGN KEY (created_by,dc) REFERENCES directory.wbt_user(id,dc) DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE cases.priority_sla_condition ADD CONSTRAINT priority_sla_condition_domain_fk FOREIGN KEY (dc) REFERENCES directory.wbt_domain(dc) ON DELETE CASCADE;
ALTER TABLE cases.priority_sla_condition ADD CONSTRAINT priority_sla_condition_priority_fk FOREIGN KEY (priority_id,dc) REFERENCES cases.priority(id,dc) ON DELETE RESTRICT DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE cases.priority_sla_condition ADD CONSTRAINT priority_sla_condition_sla_condition_id_fk FOREIGN KEY (sla_condition_id) REFERENCES cases.sla_condition(id) ON DELETE CASCADE;
ALTER TABLE cases.priority_sla_condition ADD CONSTRAINT priority_sla_condition_updated_by_fk FOREIGN KEY (updated_by) REFERENCES directory.wbt_user(id) ON DELETE SET NULL DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE cases.priority_sla_condition ADD CONSTRAINT priority_sla_condition_updated_dc_fk FOREIGN KEY (updated_by,dc) REFERENCES directory.wbt_user(id,dc) DEFERRABLE INITIALLY DEFERRED;


-- cases.related_case foreign keys

ALTER TABLE cases.related_case ADD CONSTRAINT related_case_created_dc_fk FOREIGN KEY (created_by,dc) REFERENCES directory.wbt_user(id,dc) DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE cases.related_case ADD CONSTRAINT related_case_created_id_fk FOREIGN KEY (created_by) REFERENCES directory.wbt_user(id) ON DELETE SET NULL DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE cases.related_case ADD CONSTRAINT related_case_primary_case_fk FOREIGN KEY (primary_case_id) REFERENCES cases."case"(id) ON DELETE CASCADE DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE cases.related_case ADD CONSTRAINT related_case_related_case_fk FOREIGN KEY (related_case_id) REFERENCES cases."case"(id) ON DELETE CASCADE DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE cases.related_case ADD CONSTRAINT related_case_updated_dc_fk FOREIGN KEY (updated_by,dc) REFERENCES directory.wbt_user(id,dc) DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE cases.related_case ADD CONSTRAINT related_case_updated_id_fk FOREIGN KEY (updated_by) REFERENCES directory.wbt_user(id) ON DELETE SET NULL DEFERRABLE INITIALLY DEFERRED;


-- cases.service_catalog foreign keys

ALTER TABLE cases.service_catalog ADD CONSTRAINT service_catalog_created_dc_fk FOREIGN KEY (created_by,dc) REFERENCES directory.wbt_user(id,dc) DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE cases.service_catalog ADD CONSTRAINT service_catalog_domain_fk FOREIGN KEY (dc) REFERENCES directory.wbt_domain(dc) ON DELETE CASCADE;
ALTER TABLE cases.service_catalog ADD CONSTRAINT service_catalog_root_fk FOREIGN KEY (root_id) REFERENCES cases.service_catalog(id) ON DELETE CASCADE;
ALTER TABLE cases.service_catalog ADD CONSTRAINT service_catalog_updated_id_fk FOREIGN KEY (updated_by) REFERENCES directory.wbt_user(id) DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE cases.service_catalog ADD CONSTRAINT service_created_id_fk FOREIGN KEY (created_by) REFERENCES directory.wbt_user(id) ON DELETE SET NULL DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE cases.service_catalog ADD CONSTRAINT service_updated_dc_fk FOREIGN KEY (updated_by,dc) REFERENCES directory.wbt_user(id,dc) DEFERRABLE INITIALLY DEFERRED;


-- cases.skill_catalog foreign keys

ALTER TABLE cases.skill_catalog ADD CONSTRAINT skill_service_created_by_fk FOREIGN KEY (created_by) REFERENCES directory.wbt_user(id) ON DELETE SET NULL DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE cases.skill_catalog ADD CONSTRAINT skill_service_created_dc_fk FOREIGN KEY (created_by,dc) REFERENCES directory.wbt_user(id,dc) DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE cases.skill_catalog ADD CONSTRAINT skill_service_domain_fk FOREIGN KEY (dc) REFERENCES directory.wbt_domain(dc) ON DELETE CASCADE;
ALTER TABLE cases.skill_catalog ADD CONSTRAINT skill_service_service_id_fk FOREIGN KEY (catalog_id) REFERENCES cases.service_catalog(id) ON DELETE CASCADE;
ALTER TABLE cases.skill_catalog ADD CONSTRAINT skill_service_updated_by_fk FOREIGN KEY (updated_by) REFERENCES directory.wbt_user(id) ON DELETE SET NULL DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE cases.skill_catalog ADD CONSTRAINT skill_service_updated_dc_fk FOREIGN KEY (updated_by,dc) REFERENCES directory.wbt_user(id,dc) DEFERRABLE INITIALLY DEFERRED;


-- cases.sla foreign keys

ALTER TABLE cases.sla ADD CONSTRAINT sla_created_dc_fk FOREIGN KEY (created_by,dc) REFERENCES directory.wbt_user(id,dc) DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE cases.sla ADD CONSTRAINT sla_created_id_fk FOREIGN KEY (created_by) REFERENCES directory.wbt_user(id) ON DELETE SET NULL DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE cases.sla ADD CONSTRAINT sla_domain_fk FOREIGN KEY (dc) REFERENCES directory.wbt_domain(dc) ON DELETE CASCADE;
ALTER TABLE cases.sla ADD CONSTRAINT sla_updated_dc_fk FOREIGN KEY (updated_by,dc) REFERENCES directory.wbt_user(id,dc) DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE cases.sla ADD CONSTRAINT sla_updated_id_fk FOREIGN KEY (updated_by) REFERENCES directory.wbt_user(id) DEFERRABLE INITIALLY DEFERRED;


-- cases.sla_condition foreign keys

ALTER TABLE cases.sla_condition ADD CONSTRAINT sla_condition_created_dc_fk FOREIGN KEY (created_by,dc) REFERENCES <?>() DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE cases.sla_condition ADD CONSTRAINT sla_condition_created_id_fk FOREIGN KEY () REFERENCES <?>() ON DELETE SET NULL DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE cases.sla_condition ADD CONSTRAINT sla_condition_domain_fk FOREIGN KEY () REFERENCES <?>() ON DELETE CASCADE;
ALTER TABLE cases.sla_condition ADD CONSTRAINT sla_condition_sla_id_fk FOREIGN KEY () REFERENCES <?>() ON DELETE CASCADE;
ALTER TABLE cases.sla_condition ADD CONSTRAINT sla_condition_updated_dc_fk FOREIGN KEY () REFERENCES <?>() DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE cases.sla_condition ADD CONSTRAINT sla_condition_updated_id_fk FOREIGN KEY () REFERENCES <?>() DEFERRABLE INITIALLY DEFERRED;


-- cases."source" foreign keys

ALTER TABLE cases."source" ADD CONSTRAINT source_created_dc_fk FOREIGN KEY (created_by,dc) REFERENCES directory.wbt_user(id,dc) DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE cases."source" ADD CONSTRAINT source_created_id_fk FOREIGN KEY (created_by) REFERENCES directory.wbt_user(id) ON DELETE SET NULL DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE cases."source" ADD CONSTRAINT source_domain_fk FOREIGN KEY (dc) REFERENCES directory.wbt_domain(dc) ON DELETE CASCADE;
ALTER TABLE cases."source" ADD CONSTRAINT source_updated_dc_fk FOREIGN KEY (updated_by,dc) REFERENCES directory.wbt_user(id,dc) DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE cases."source" ADD CONSTRAINT source_updated_id_fk FOREIGN KEY (updated_by) REFERENCES directory.wbt_user(id) DEFERRABLE INITIALLY DEFERRED;


-- cases.status foreign keys

ALTER TABLE cases.status ADD CONSTRAINT status_created_id_fk FOREIGN KEY (created_by) REFERENCES directory.wbt_user(id) ON DELETE SET NULL DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE cases.status ADD CONSTRAINT status_domain_fk FOREIGN KEY (dc) REFERENCES directory.wbt_domain(dc) ON DELETE CASCADE;
ALTER TABLE cases.status ADD CONSTRAINT status_updated_dc_fk FOREIGN KEY (updated_by,dc) REFERENCES directory.wbt_user(id,dc) DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE cases.status ADD CONSTRAINT status_updated_id_fk FOREIGN KEY (updated_by) REFERENCES directory.wbt_user(id) DEFERRABLE INITIALLY DEFERRED;


-- cases.status_condition foreign keys

ALTER TABLE cases.status_condition ADD CONSTRAINT status_condition_created_by_fk FOREIGN KEY (created_by) REFERENCES directory.wbt_user(id) ON DELETE SET NULL DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE cases.status_condition ADD CONSTRAINT status_condition_created_dc_fk FOREIGN KEY (created_by,dc) REFERENCES directory.wbt_user(id,dc) DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE cases.status_condition ADD CONSTRAINT status_condition_status_id_fk FOREIGN KEY (status_id) REFERENCES cases.status(id) ON DELETE CASCADE;
ALTER TABLE cases.status_condition ADD CONSTRAINT status_condition_updated_by_fk FOREIGN KEY (updated_by) REFERENCES directory.wbt_user(id) ON DELETE SET NULL DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE cases.status_condition ADD CONSTRAINT status_condition_updated_dc_fk FOREIGN KEY (updated_by,dc) REFERENCES directory.wbt_user(id,dc) DEFERRABLE INITIALLY DEFERRED;


-- cases.team_catalog foreign keys

ALTER TABLE cases.team_catalog ADD CONSTRAINT team_service_created_by_fk FOREIGN KEY (created_by) REFERENCES directory.wbt_user(id) ON DELETE SET NULL DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE cases.team_catalog ADD CONSTRAINT team_service_created_dc_fk FOREIGN KEY (created_by,dc) REFERENCES directory.wbt_user(id,dc) DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE cases.team_catalog ADD CONSTRAINT team_service_domain_fk FOREIGN KEY (dc) REFERENCES directory.wbt_domain(dc) ON DELETE CASCADE;
ALTER TABLE cases.team_catalog ADD CONSTRAINT team_service_service_id_fk FOREIGN KEY (catalog_id) REFERENCES cases.service_catalog(id) ON DELETE CASCADE;
ALTER TABLE cases.team_catalog ADD CONSTRAINT team_service_updated_by_fk FOREIGN KEY (updated_by) REFERENCES directory.wbt_user(id) ON DELETE SET NULL DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE cases.team_catalog ADD CONSTRAINT team_service_updated_dc_fk FOREIGN KEY (updated_by,dc) REFERENCES directory.wbt_user(id,dc) DEFERRABLE INITIALLY DEFERRED;