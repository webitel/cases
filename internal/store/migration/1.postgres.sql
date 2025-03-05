CREATE SCHEMA cases;


--
-- Name: update_case_timings(); Type: FUNCTION; Schema: cases; Owner: -
--

CREATE FUNCTION cases.update_case_timings() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    -- Set reacted_at if status_condition is not initial
    IF (NEW.status_condition IS NOT NULL) THEN
        -- Check if the status_condition is initial
        PERFORM initial
        FROM cases.status_condition
        WHERE id = NEW.status_condition AND initial = TRUE;

        -- If no initial status found, set reacted_at
        IF NOT FOUND THEN
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

    RETURN NEW;
END;
$$;


--
-- Name: appeal_id; Type: SEQUENCE; Schema: cases; Owner: -
--

CREATE SEQUENCE cases.appeal_id
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: case_id; Type: SEQUENCE; Schema: cases; Owner: -
--

CREATE SEQUENCE cases.case_id
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: case; Type: TABLE; Schema: cases; Owner: -
--

CREATE TABLE cases."case" (
    id bigint DEFAULT nextval('cases.case_id'::regclass) NOT NULL,
    dc bigint NOT NULL,
    name character varying NOT NULL,
    subject character varying NOT NULL,
    description character varying,
    ver integer DEFAULT 0 NOT NULL,
    created_by bigint,
    created_at timestamp without time zone DEFAULT timezone('utc'::text, now()) NOT NULL,
    updated_by bigint,
    updated_at timestamp without time zone DEFAULT timezone('utc'::text, now()) NOT NULL,
    status bigint NOT NULL,
    close_reason_group bigint NOT NULL,
    assignee bigint,
    reporter bigint,
    impacted bigint,
    contact_group bigint,
    priority bigint NOT NULL,
    source bigint NOT NULL,
    status_condition bigint,
    close_result character varying,
    close_reason bigint,
    rating bigint,
    rating_comment character varying,
    service bigint NOT NULL,
    sla bigint NOT NULL,
    planned_reaction_at timestamp without time zone DEFAULT timezone('utc'::text, now()),
    planned_resolve_at timestamp without time zone DEFAULT timezone('utc'::text, now()),
    sla_condition_id bigint,
    contact_info text,
    reacted_at timestamp without time zone,
    resolved_at timestamp without time zone
);


--
-- Name: case_acl; Type: TABLE; Schema: cases; Owner: -
--

CREATE TABLE cases.case_acl (
    dc bigint NOT NULL,
    grantor bigint,
    subject bigint NOT NULL,
    object bigint NOT NULL,
    access smallint
);


--
-- Name: case_comment_id; Type: SEQUENCE; Schema: cases; Owner: -
--

CREATE SEQUENCE cases.case_comment_id
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: case_comment; Type: TABLE; Schema: cases; Owner: -
--

CREATE TABLE cases.case_comment (
    id bigint DEFAULT nextval('cases.case_comment_id'::regclass) NOT NULL,
    ver integer DEFAULT 0 NOT NULL,
    created_by bigint,
    created_at timestamp without time zone DEFAULT timezone('utc'::text, now()) NOT NULL,
    updated_by bigint,
    updated_at timestamp without time zone DEFAULT timezone('utc'::text, now()) NOT NULL,
    comment character varying NOT NULL,
    case_id bigint NOT NULL,
    dc bigint NOT NULL
);


--
-- Name: case_comment_acl; Type: TABLE; Schema: cases; Owner: -
--

CREATE TABLE cases.case_comment_acl (
    dc bigint NOT NULL,
    grantor bigint,
    subject bigint NOT NULL,
    object bigint NOT NULL,
    access smallint
);


--
-- Name: case_communication; Type: TABLE; Schema: cases; Owner: -
--

CREATE TABLE cases.case_communication (
    id bigint NOT NULL,
    ver integer DEFAULT 0 NOT NULL,
    case_id bigint NOT NULL,
    communication_type integer NOT NULL,
    communication_id text NOT NULL,
    dc bigint NOT NULL,
    created_by bigint,
    created_at timestamp without time zone DEFAULT now() NOT NULL
);


--
-- Name: TABLE case_communication; Type: COMMENT; Schema: cases; Owner: -
--

COMMENT ON TABLE cases.case_communication IS 'Table connects case with other possible communications like chats, emails, calls';


--
-- Name: COLUMN case_communication.communication_type; Type: COMMENT; Schema: cases; Owner: -
--

COMMENT ON COLUMN cases.case_communication.communication_type IS 'type of communication case connected to';


--
-- Name: COLUMN case_communication.communication_id; Type: COMMENT; Schema: cases; Owner: -
--

COMMENT ON COLUMN cases.case_communication.communication_id IS 'Id of communication case connected to';


--
-- Name: COLUMN case_communication.dc; Type: COMMENT; Schema: cases; Owner: -
--

COMMENT ON COLUMN cases.case_communication.dc IS 'Domain component';


--
-- Name: case_communication_id_seq; Type: SEQUENCE; Schema: cases; Owner: -
--

CREATE SEQUENCE cases.case_communication_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: case_communication_id_seq; Type: SEQUENCE OWNED BY; Schema: cases; Owner: -
--

ALTER SEQUENCE cases.case_communication_id_seq OWNED BY cases.case_communication.id;


--
-- Name: case_link_id_seq; Type: SEQUENCE; Schema: cases; Owner: -
--

CREATE SEQUENCE cases.case_link_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: case_link; Type: TABLE; Schema: cases; Owner: -
--

CREATE TABLE cases.case_link (
    id bigint DEFAULT nextval('cases.case_link_id_seq'::regclass) NOT NULL,
    ver integer DEFAULT 0 NOT NULL,
    created_by bigint,
    created_at timestamp without time zone DEFAULT timezone('utc'::text, now()) NOT NULL,
    updated_by bigint,
    updated_at timestamp without time zone DEFAULT timezone('utc'::text, now()) NOT NULL,
    name character varying,
    url character varying NOT NULL,
    case_id bigint NOT NULL,
    dc bigint NOT NULL
);


--
-- Name: reason_id; Type: SEQUENCE; Schema: cases; Owner: -
--

CREATE SEQUENCE cases.reason_id
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: close_reason; Type: TABLE; Schema: cases; Owner: -
--

CREATE TABLE cases.close_reason (
    id bigint DEFAULT nextval('cases.reason_id'::regclass) NOT NULL,
    name text NOT NULL,
    description text,
    created_at timestamp without time zone DEFAULT timezone('utc'::text, now()) NOT NULL,
    updated_at timestamp without time zone DEFAULT timezone('utc'::text, now()) NOT NULL,
    created_by bigint,
    updated_by bigint,
    close_reason_id bigint NOT NULL,
    dc bigint NOT NULL
);


--
-- Name: close_reason_id; Type: SEQUENCE; Schema: cases; Owner: -
--

CREATE SEQUENCE cases.close_reason_id
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: close_reason_group; Type: TABLE; Schema: cases; Owner: -
--

CREATE TABLE cases.close_reason_group (
    id bigint DEFAULT nextval('cases.close_reason_id'::regclass) NOT NULL,
    name text NOT NULL,
    description text,
    created_at timestamp without time zone DEFAULT timezone('utc'::text, now()) NOT NULL,
    updated_at timestamp without time zone DEFAULT timezone('utc'::text, now()) NOT NULL,
    created_by bigint,
    updated_by bigint,
    dc bigint NOT NULL
);


--
-- Name: priority_id; Type: SEQUENCE; Schema: cases; Owner: -
--

CREATE SEQUENCE cases.priority_id
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: priority; Type: TABLE; Schema: cases; Owner: -
--

CREATE TABLE cases.priority (
    id bigint DEFAULT nextval('cases.priority_id'::regclass) NOT NULL,
    name text NOT NULL,
    description text,
    created_at timestamp without time zone DEFAULT timezone('utc'::text, now()) NOT NULL,
    updated_at timestamp without time zone DEFAULT timezone('utc'::text, now()) NOT NULL,
    created_by bigint,
    updated_by bigint,
    dc bigint NOT NULL,
    color character varying(25) NOT NULL
);


--
-- Name: priority_sla_condition_id; Type: SEQUENCE; Schema: cases; Owner: -
--

CREATE SEQUENCE cases.priority_sla_condition_id
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: priority_sla_condition; Type: TABLE; Schema: cases; Owner: -
--

CREATE TABLE cases.priority_sla_condition (
    id bigint DEFAULT nextval('cases.priority_sla_condition_id'::regclass) NOT NULL,
    created_at timestamp without time zone DEFAULT timezone('utc'::text, now()) NOT NULL,
    updated_at timestamp without time zone DEFAULT timezone('utc'::text, now()) NOT NULL,
    created_by bigint,
    updated_by bigint,
    sla_condition_id bigint NOT NULL,
    priority_id bigint NOT NULL,
    dc bigint NOT NULL
);


--
-- Name: related_case; Type: TABLE; Schema: cases; Owner: -
--

CREATE TABLE cases.related_case (
    id integer NOT NULL,
    dc bigint NOT NULL,
    ver integer DEFAULT 0 NOT NULL,
    created_by bigint,
    created_at timestamp without time zone DEFAULT timezone('utc'::text, now()) NOT NULL,
    updated_by bigint,
    updated_at timestamp without time zone DEFAULT timezone('utc'::text, now()) NOT NULL,
    relation_type smallint NOT NULL,
    primary_case_id bigint NOT NULL,
    related_case_id bigint NOT NULL
);


--
-- Name: related_case_id_seq; Type: SEQUENCE; Schema: cases; Owner: -
--

CREATE SEQUENCE cases.related_case_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: related_case_id_seq; Type: SEQUENCE OWNED BY; Schema: cases; Owner: -
--

ALTER SEQUENCE cases.related_case_id_seq OWNED BY cases.related_case.id;


--
-- Name: service_catalog; Type: TABLE; Schema: cases; Owner: -
--

CREATE TABLE cases.service_catalog (
    id bigint NOT NULL,
    root_id bigint,
    description text,
    code character varying(50),
    prefix character varying(20),
    state boolean DEFAULT true NOT NULL,
    assignee_id bigint,
    created_at timestamp without time zone DEFAULT timezone('utc'::text, now()) NOT NULL,
    updated_at timestamp without time zone DEFAULT timezone('utc'::text, now()) NOT NULL,
    created_by bigint,
    updated_by bigint,
    dc bigint NOT NULL,
    sla_id bigint,
    status_id bigint,
    close_reason_group_id bigint,
    group_id integer,
    name text NOT NULL,
    catalog_id bigint,
    CONSTRAINT chk_root_id_close_reason_id CHECK ((NOT ((root_id IS NULL) AND (close_reason_group_id IS NULL))))
);


--
-- Name: service_id_seq; Type: SEQUENCE; Schema: cases; Owner: -
--

CREATE SEQUENCE cases.service_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: service_id_seq; Type: SEQUENCE OWNED BY; Schema: cases; Owner: -
--

ALTER SEQUENCE cases.service_id_seq OWNED BY cases.service_catalog.id;


--
-- Name: skill_catalog; Type: TABLE; Schema: cases; Owner: -
--

CREATE TABLE cases.skill_catalog (
    id bigint NOT NULL,
    created_at timestamp without time zone DEFAULT timezone('utc'::text, now()) NOT NULL,
    updated_at timestamp without time zone DEFAULT timezone('utc'::text, now()) NOT NULL,
    created_by bigint,
    updated_by bigint NOT NULL,
    catalog_id bigint NOT NULL,
    skill_id bigint NOT NULL,
    dc bigint NOT NULL
);


--
-- Name: skill_service_id_seq; Type: SEQUENCE; Schema: cases; Owner: -
--

CREATE SEQUENCE cases.skill_service_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: skill_service_id_seq; Type: SEQUENCE OWNED BY; Schema: cases; Owner: -
--

ALTER SEQUENCE cases.skill_service_id_seq OWNED BY cases.skill_catalog.id;


--
-- Name: sla_id; Type: SEQUENCE; Schema: cases; Owner: -
--

CREATE SEQUENCE cases.sla_id
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: sla; Type: TABLE; Schema: cases; Owner: -
--

CREATE TABLE cases.sla (
    id bigint DEFAULT nextval('cases.sla_id'::regclass) NOT NULL,
    name text NOT NULL,
    description text,
    created_at timestamp without time zone DEFAULT timezone('utc'::text, now()) NOT NULL,
    updated_at timestamp without time zone DEFAULT timezone('utc'::text, now()) NOT NULL,
    created_by bigint,
    updated_by bigint,
    dc bigint NOT NULL,
    valid_from timestamp without time zone DEFAULT timezone('utc'::text, now()),
    valid_to timestamp without time zone DEFAULT timezone('utc'::text, now()),
    reaction_time integer DEFAULT 0 NOT NULL,
    resolution_time integer DEFAULT 0 NOT NULL,
    calendar_id bigint NOT NULL
);


--
-- Name: sla_condition_id; Type: SEQUENCE; Schema: cases; Owner: -
--

CREATE SEQUENCE cases.sla_condition_id
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: sla_condition; Type: TABLE; Schema: cases; Owner: -
--

CREATE TABLE cases.sla_condition (
    id bigint DEFAULT nextval('cases.sla_condition_id'::regclass) NOT NULL,
    name text NOT NULL,
    created_at timestamp without time zone DEFAULT timezone('utc'::text, now()) NOT NULL,
    updated_at timestamp without time zone DEFAULT timezone('utc'::text, now()) NOT NULL,
    created_by bigint NOT NULL,
    updated_by bigint,
    dc bigint NOT NULL,
    reaction_time integer DEFAULT 0 NOT NULL,
    resolution_time integer DEFAULT 0 NOT NULL,
    sla_id bigint NOT NULL
);


--
-- Name: source; Type: TABLE; Schema: cases; Owner: -
--

CREATE TABLE cases.source (
    id bigint DEFAULT nextval('cases.appeal_id'::regclass) NOT NULL,
    name text NOT NULL,
    description text,
    created_at timestamp without time zone DEFAULT timezone('utc'::text, now()) NOT NULL,
    updated_at timestamp without time zone DEFAULT timezone('utc'::text, now()) NOT NULL,
    created_by bigint,
    updated_by bigint,
    type text NOT NULL,
    dc bigint NOT NULL
);


--
-- Name: status_id; Type: SEQUENCE; Schema: cases; Owner: -
--

CREATE SEQUENCE cases.status_id
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: status; Type: TABLE; Schema: cases; Owner: -
--

CREATE TABLE cases.status (
    id bigint DEFAULT nextval('cases.status_id'::regclass) NOT NULL,
    name text NOT NULL,
    description text,
    created_at timestamp without time zone DEFAULT timezone('utc'::text, now()) NOT NULL,
    updated_at timestamp without time zone DEFAULT timezone('utc'::text, now()) NOT NULL,
    created_by bigint,
    updated_by bigint,
    dc bigint NOT NULL
);


--
-- Name: status_condition_id; Type: SEQUENCE; Schema: cases; Owner: -
--

CREATE SEQUENCE cases.status_condition_id
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: status_condition; Type: TABLE; Schema: cases; Owner: -
--

CREATE TABLE cases.status_condition (
    id bigint DEFAULT nextval('cases.status_condition_id'::regclass) NOT NULL,
    name text NOT NULL,
    description text,
    created_at timestamp without time zone DEFAULT timezone('utc'::text, now()) NOT NULL,
    updated_at timestamp without time zone DEFAULT timezone('utc'::text, now()) NOT NULL,
    created_by bigint,
    updated_by bigint,
    status_id bigint NOT NULL,
    initial boolean NOT NULL,
    final boolean NOT NULL,
    dc bigint NOT NULL
);


--
-- Name: team_catalog; Type: TABLE; Schema: cases; Owner: -
--

CREATE TABLE cases.team_catalog (
    id bigint NOT NULL,
    created_at timestamp without time zone DEFAULT timezone('utc'::text, now()) NOT NULL,
    updated_at timestamp without time zone DEFAULT timezone('utc'::text, now()) NOT NULL,
    created_by bigint,
    updated_by bigint,
    catalog_id bigint NOT NULL,
    team_id bigint NOT NULL,
    dc bigint NOT NULL
);


--
-- Name: team_service_id_seq; Type: SEQUENCE; Schema: cases; Owner: -
--

CREATE SEQUENCE cases.team_service_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: team_service_id_seq; Type: SEQUENCE OWNED BY; Schema: cases; Owner: -
--

ALTER SEQUENCE cases.team_service_id_seq OWNED BY cases.team_catalog.id;


--
-- Name: case_communication id; Type: DEFAULT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.case_communication ALTER COLUMN id SET DEFAULT nextval('cases.case_communication_id_seq'::regclass);


--
-- Name: related_case id; Type: DEFAULT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.related_case ALTER COLUMN id SET DEFAULT nextval('cases.related_case_id_seq'::regclass);


--
-- Name: service_catalog id; Type: DEFAULT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.service_catalog ALTER COLUMN id SET DEFAULT nextval('cases.service_id_seq'::regclass);


--
-- Name: skill_catalog id; Type: DEFAULT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.skill_catalog ALTER COLUMN id SET DEFAULT nextval('cases.skill_service_id_seq'::regclass);


--
-- Name: team_catalog id; Type: DEFAULT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.team_catalog ALTER COLUMN id SET DEFAULT nextval('cases.team_service_id_seq'::regclass);


--
-- Name: case_comment case_comment_id_case_id_unique; Type: CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.case_comment
    ADD CONSTRAINT case_comment_id_case_id_unique UNIQUE (id, case_id);


--
-- Name: case case_fk; Type: CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases."case"
    ADD CONSTRAINT case_fk UNIQUE (id, dc);


--
-- Name: case_link case_link_fk; Type: CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.case_link
    ADD CONSTRAINT case_link_fk UNIQUE (id, dc);


--
-- Name: case_link case_link_id_case_id_unique; Type: CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.case_link
    ADD CONSTRAINT case_link_id_case_id_unique UNIQUE (id, case_id);


--
-- Name: case case_pk; Type: CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases."case"
    ADD CONSTRAINT case_pk PRIMARY KEY (id);


--
-- Name: close_reason_group close_reason_fk; Type: CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.close_reason_group
    ADD CONSTRAINT close_reason_fk UNIQUE (id, dc);


--
-- Name: close_reason_group close_reason_pk; Type: CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.close_reason_group
    ADD CONSTRAINT close_reason_pk PRIMARY KEY (id);


--
-- Name: case_comment comment_case_fk; Type: CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.case_comment
    ADD CONSTRAINT comment_case_fk UNIQUE (id, dc);


--
-- Name: case_comment comment_case_pk; Type: CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.case_comment
    ADD CONSTRAINT comment_case_pk PRIMARY KEY (id);


--
-- Name: case_link link_case_pk; Type: CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.case_link
    ADD CONSTRAINT link_case_pk PRIMARY KEY (id);


--
-- Name: priority priority_fk; Type: CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.priority
    ADD CONSTRAINT priority_fk UNIQUE (id, dc);


--
-- Name: priority priority_pk; Type: CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.priority
    ADD CONSTRAINT priority_pk PRIMARY KEY (id);


--
-- Name: priority_sla_condition priority_sla_condition_fk; Type: CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.priority_sla_condition
    ADD CONSTRAINT priority_sla_condition_fk UNIQUE (id, dc);


--
-- Name: priority_sla_condition priority_sla_condition_pk; Type: CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.priority_sla_condition
    ADD CONSTRAINT priority_sla_condition_pk PRIMARY KEY (id);


--
-- Name: close_reason reason_fk; Type: CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.close_reason
    ADD CONSTRAINT reason_fk UNIQUE (id, dc);


--
-- Name: close_reason reason_pk; Type: CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.close_reason
    ADD CONSTRAINT reason_pk PRIMARY KEY (id);


--
-- Name: related_case related_case_fk; Type: CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.related_case
    ADD CONSTRAINT related_case_fk UNIQUE (id, dc);


--
-- Name: related_case related_case_pk; Type: CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.related_case
    ADD CONSTRAINT related_case_pk PRIMARY KEY (id);


--
-- Name: service_catalog service_fk; Type: CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.service_catalog
    ADD CONSTRAINT service_fk UNIQUE (id, dc);


--
-- Name: service_catalog service_pk; Type: CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.service_catalog
    ADD CONSTRAINT service_pk PRIMARY KEY (id);


--
-- Name: skill_catalog skill_service_fk; Type: CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.skill_catalog
    ADD CONSTRAINT skill_service_fk UNIQUE (id, dc);


--
-- Name: skill_catalog skill_service_pk; Type: CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.skill_catalog
    ADD CONSTRAINT skill_service_pk PRIMARY KEY (id);


--
-- Name: sla_condition sla_condition_fk; Type: CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.sla_condition
    ADD CONSTRAINT sla_condition_fk UNIQUE (id, dc);


--
-- Name: sla_condition sla_condition_pk; Type: CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.sla_condition
    ADD CONSTRAINT sla_condition_pk PRIMARY KEY (id);


--
-- Name: sla_condition sla_condition_sla_name_pk; Type: CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.sla_condition
    ADD CONSTRAINT sla_condition_sla_name_pk UNIQUE (sla_id, name);


--
-- Name: sla sla_fk; Type: CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.sla
    ADD CONSTRAINT sla_fk UNIQUE (id, dc);


--
-- Name: sla sla_pk; Type: CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.sla
    ADD CONSTRAINT sla_pk PRIMARY KEY (id);


--
-- Name: source source_fk; Type: CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.source
    ADD CONSTRAINT source_fk UNIQUE (id, dc);


--
-- Name: source source_pk; Type: CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.source
    ADD CONSTRAINT source_pk PRIMARY KEY (id);


--
-- Name: status_condition status_condition_fk; Type: CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.status_condition
    ADD CONSTRAINT status_condition_fk UNIQUE (id, dc);


--
-- Name: status_condition status_condition_pk; Type: CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.status_condition
    ADD CONSTRAINT status_condition_pk PRIMARY KEY (id);


--
-- Name: status status_fk; Type: CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.status
    ADD CONSTRAINT status_fk UNIQUE (id, dc);


--
-- Name: status status_pk; Type: CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.status
    ADD CONSTRAINT status_pk PRIMARY KEY (id);


--
-- Name: team_catalog team_service_fk; Type: CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.team_catalog
    ADD CONSTRAINT team_service_fk UNIQUE (id, dc);


--
-- Name: team_catalog team_service_pk; Type: CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.team_catalog
    ADD CONSTRAINT team_service_pk PRIMARY KEY (id);


--
-- Name: case_acl_grantor_index; Type: INDEX; Schema: cases; Owner: -
--

CREATE INDEX case_acl_grantor_index ON cases.case_acl USING btree (grantor);


--
-- Name: case_acl_object_subject_pk; Type: INDEX; Schema: cases; Owner: -
--

CREATE UNIQUE INDEX case_acl_object_subject_pk ON cases.case_acl USING btree (object, subject);


--
-- Name: case_acl_subject_object_uindex; Type: INDEX; Schema: cases; Owner: -
--

CREATE UNIQUE INDEX case_acl_subject_object_uindex ON cases.case_acl USING btree (subject, object);


--
-- Name: case_comment_acl_grantor_uindex; Type: INDEX; Schema: cases; Owner: -
--

CREATE INDEX case_comment_acl_grantor_uindex ON cases.case_comment_acl USING btree (grantor);


--
-- Name: case_comment_acl_object_subject_pk; Type: INDEX; Schema: cases; Owner: -
--

CREATE UNIQUE INDEX case_comment_acl_object_subject_pk ON cases.case_comment_acl USING btree (object, subject);


--
-- Name: case_comment_acl_subject_object_uindex; Type: INDEX; Schema: cases; Owner: -
--

CREATE UNIQUE INDEX case_comment_acl_subject_object_uindex ON cases.case_comment_acl USING btree (subject, object);


--
-- Name: case_dc; Type: INDEX; Schema: cases; Owner: -
--

CREATE INDEX case_dc ON cases."case" USING btree (dc);


--
-- Name: close_reason_dc; Type: INDEX; Schema: cases; Owner: -
--

CREATE INDEX close_reason_dc ON cases.close_reason_group USING btree (dc);


--
-- Name: comment_case_dc; Type: INDEX; Schema: cases; Owner: -
--

CREATE INDEX comment_case_dc ON cases.case_comment USING btree (dc);


--
-- Name: idx_catalog_id_skill_id; Type: INDEX; Schema: cases; Owner: -
--

CREATE UNIQUE INDEX idx_catalog_id_skill_id ON cases.skill_catalog USING btree (catalog_id, skill_id);


--
-- Name: idx_catalog_id_team_id; Type: INDEX; Schema: cases; Owner: -
--

CREATE UNIQUE INDEX idx_catalog_id_team_id ON cases.team_catalog USING btree (catalog_id, team_id);


--
-- Name: idx_comment_case_id; Type: INDEX; Schema: cases; Owner: -
--

CREATE INDEX idx_comment_case_id ON cases.case_link USING btree (case_id);


--
-- Name: idx_link_case_id; Type: INDEX; Schema: cases; Owner: -
--

CREATE INDEX idx_link_case_id ON cases.case_comment USING btree (case_id);


--
-- Name: idx_related_case_primary_case; Type: INDEX; Schema: cases; Owner: -
--

CREATE INDEX idx_related_case_primary_case ON cases.related_case USING btree (primary_case_id);


--
-- Name: idx_related_case_related_case; Type: INDEX; Schema: cases; Owner: -
--

CREATE INDEX idx_related_case_related_case ON cases.related_case USING btree (related_case_id);


--
-- Name: idx_service_skill; Type: INDEX; Schema: cases; Owner: -
--

CREATE UNIQUE INDEX idx_service_skill ON cases.skill_catalog USING btree (skill_id, catalog_id);


--
-- Name: idx_service_team; Type: INDEX; Schema: cases; Owner: -
--

CREATE INDEX idx_service_team ON cases.team_catalog USING btree (catalog_id, team_id);


--
-- Name: idx_sla_condition_priority; Type: INDEX; Schema: cases; Owner: -
--

CREATE UNIQUE INDEX idx_sla_condition_priority ON cases.priority_sla_condition USING btree (sla_condition_id, priority_id);


--
-- Name: link_case_dc; Type: INDEX; Schema: cases; Owner: -
--

CREATE INDEX link_case_dc ON cases.case_link USING btree (dc);


--
-- Name: priority_dc; Type: INDEX; Schema: cases; Owner: -
--

CREATE INDEX priority_dc ON cases.priority USING btree (dc);


--
-- Name: priority_sla_condition_dc; Type: INDEX; Schema: cases; Owner: -
--

CREATE INDEX priority_sla_condition_dc ON cases.priority_sla_condition USING btree (dc);


--
-- Name: priority_sla_condition_source; Type: INDEX; Schema: cases; Owner: -
--

CREATE INDEX priority_sla_condition_source ON cases.priority_sla_condition USING btree (sla_condition_id);


--
-- Name: reason_id_dc; Type: INDEX; Schema: cases; Owner: -
--

CREATE INDEX reason_id_dc ON cases.close_reason USING btree (dc);


--
-- Name: reason_source; Type: INDEX; Schema: cases; Owner: -
--

CREATE INDEX reason_source ON cases.close_reason USING btree (close_reason_id);


--
-- Name: related_case_dc; Type: INDEX; Schema: cases; Owner: -
--

CREATE INDEX related_case_dc ON cases.related_case USING btree (dc);


--
-- Name: service_dc; Type: INDEX; Schema: cases; Owner: -
--

CREATE INDEX service_dc ON cases.service_catalog USING btree (dc);


--
-- Name: skill_service_dc; Type: INDEX; Schema: cases; Owner: -
--

CREATE INDEX skill_service_dc ON cases.skill_catalog USING btree (dc);


--
-- Name: skill_service_source; Type: INDEX; Schema: cases; Owner: -
--

CREATE INDEX skill_service_source ON cases.skill_catalog USING btree (catalog_id);


--
-- Name: sla_condition_dc; Type: INDEX; Schema: cases; Owner: -
--

CREATE INDEX sla_condition_dc ON cases.sla_condition USING btree (dc);


--
-- Name: sla_condition_source; Type: INDEX; Schema: cases; Owner: -
--

CREATE INDEX sla_condition_source ON cases.sla_condition USING btree (sla_id);


--
-- Name: sla_dc; Type: INDEX; Schema: cases; Owner: -
--

CREATE INDEX sla_dc ON cases.sla USING btree (dc);


--
-- Name: source_dc; Type: INDEX; Schema: cases; Owner: -
--

CREATE INDEX source_dc ON cases.source USING btree (dc);


--
-- Name: status_condition_dc; Type: INDEX; Schema: cases; Owner: -
--

CREATE INDEX status_condition_dc ON cases.status_condition USING btree (dc);


--
-- Name: status_condition_source; Type: INDEX; Schema: cases; Owner: -
--

CREATE INDEX status_condition_source ON cases.status_condition USING btree (status_id);


--
-- Name: status_dc; Type: INDEX; Schema: cases; Owner: -
--

CREATE INDEX status_dc ON cases.status USING btree (dc);


--
-- Name: team_service_dc; Type: INDEX; Schema: cases; Owner: -
--

CREATE INDEX team_service_dc ON cases.team_catalog USING btree (dc);


--
-- Name: team_service_source; Type: INDEX; Schema: cases; Owner: -
--

CREATE INDEX team_service_source ON cases.team_catalog USING btree (catalog_id);


--
-- Name: unique_event_uindex; Type: INDEX; Schema: cases; Owner: -
--

CREATE UNIQUE INDEX unique_event_uindex ON cases.case_communication USING btree (communication_id, communication_type, case_id);


--
-- Name: unique_related_cases_relation; Type: INDEX; Schema: cases; Owner: -
--

CREATE UNIQUE INDEX unique_related_cases_relation ON cases.related_case USING btree (LEAST(primary_case_id, related_case_id), GREATEST(primary_case_id, related_case_id), relation_type);


--
-- Name: case_comment tg_case_comment_rbac; Type: TRIGGER; Schema: cases; Owner: -
--

CREATE TRIGGER tg_case_comment_rbac AFTER INSERT ON cases.case_comment FOR EACH ROW EXECUTE FUNCTION directory.tg_obj_default_rbac('case_comments');


--
-- Name: case tg_case_rbac; Type: TRIGGER; Schema: cases; Owner: -
--

CREATE TRIGGER tg_case_rbac AFTER INSERT ON cases."case" FOR EACH ROW EXECUTE FUNCTION directory.tg_obj_default_rbac('cases');


--
-- Name: case trigger_update_case_timings; Type: TRIGGER; Schema: cases; Owner: -
--

CREATE TRIGGER trigger_update_case_timings BEFORE INSERT OR UPDATE ON cases."case" FOR EACH ROW EXECUTE FUNCTION cases.update_case_timings();


--
-- Name: case_acl case_acl_domain_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.case_acl
    ADD CONSTRAINT case_acl_domain_fk FOREIGN KEY (dc) REFERENCES directory.wbt_domain(dc) ON DELETE CASCADE;


--
-- Name: case_acl case_acl_grantor_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.case_acl
    ADD CONSTRAINT case_acl_grantor_fk FOREIGN KEY (grantor, dc) REFERENCES directory.wbt_auth(id, dc) DEFERRABLE INITIALLY DEFERRED;


--
-- Name: case_acl case_acl_grantor_id_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.case_acl
    ADD CONSTRAINT case_acl_grantor_id_fk FOREIGN KEY (grantor) REFERENCES directory.wbt_auth(id) ON DELETE SET NULL;


--
-- Name: case_acl case_acl_object_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.case_acl
    ADD CONSTRAINT case_acl_object_fk FOREIGN KEY (object, dc) REFERENCES cases."case"(id, dc) ON DELETE CASCADE DEFERRABLE INITIALLY DEFERRED;


--
-- Name: case_acl case_acl_subject_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.case_acl
    ADD CONSTRAINT case_acl_subject_fk FOREIGN KEY (subject, dc) REFERENCES directory.wbt_auth(id, dc) ON DELETE CASCADE;


--
-- Name: case case_assignee_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases."case"
    ADD CONSTRAINT case_assignee_fk FOREIGN KEY (assignee) REFERENCES contacts.contact(id) ON DELETE SET NULL DEFERRABLE INITIALLY DEFERRED;


--
-- Name: case case_close_reason_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases."case"
    ADD CONSTRAINT case_close_reason_fk FOREIGN KEY (close_reason) REFERENCES cases.close_reason(id) ON DELETE SET NULL DEFERRABLE INITIALLY DEFERRED;


--
-- Name: case case_close_reason_group_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases."case"
    ADD CONSTRAINT case_close_reason_group_fk FOREIGN KEY (close_reason_group) REFERENCES cases.close_reason_group(id) ON DELETE RESTRICT DEFERRABLE INITIALLY DEFERRED;


--
-- Name: case_comment_acl case_comment_acl_domain_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.case_comment_acl
    ADD CONSTRAINT case_comment_acl_domain_fk FOREIGN KEY (dc) REFERENCES directory.wbt_domain(dc) ON DELETE CASCADE;


--
-- Name: case_comment_acl case_comment_acl_grantor_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.case_comment_acl
    ADD CONSTRAINT case_comment_acl_grantor_fk FOREIGN KEY (grantor, dc) REFERENCES directory.wbt_auth(id, dc) DEFERRABLE INITIALLY DEFERRED;


--
-- Name: case_comment_acl case_comment_acl_grantor_id_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.case_comment_acl
    ADD CONSTRAINT case_comment_acl_grantor_id_fk FOREIGN KEY (grantor) REFERENCES directory.wbt_auth(id) ON DELETE SET NULL;


--
-- Name: case_comment_acl case_comment_acl_object_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.case_comment_acl
    ADD CONSTRAINT case_comment_acl_object_fk FOREIGN KEY (object, dc) REFERENCES cases.case_comment(id, dc) ON DELETE CASCADE DEFERRABLE INITIALLY DEFERRED;


--
-- Name: case_comment_acl case_comment_acl_subject_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.case_comment_acl
    ADD CONSTRAINT case_comment_acl_subject_fk FOREIGN KEY (subject, dc) REFERENCES directory.wbt_auth(id, dc) ON DELETE CASCADE;


--
-- Name: case_communication case_communication_case_id_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.case_communication
    ADD CONSTRAINT case_communication_case_id_fk FOREIGN KEY (case_id) REFERENCES cases."case"(id) ON DELETE CASCADE;


--
-- Name: case_communication case_communication_wbt_domain_dc_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.case_communication
    ADD CONSTRAINT case_communication_wbt_domain_dc_fk FOREIGN KEY (dc) REFERENCES directory.wbt_domain(dc) ON DELETE CASCADE;


--
-- Name: case_communication case_communication_wbt_user_id_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.case_communication
    ADD CONSTRAINT case_communication_wbt_user_id_fk FOREIGN KEY (created_by) REFERENCES directory.wbt_user(id);


--
-- Name: case case_created_dc_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases."case"
    ADD CONSTRAINT case_created_dc_fk FOREIGN KEY (created_by, dc) REFERENCES directory.wbt_user(id, dc) DEFERRABLE INITIALLY DEFERRED;


--
-- Name: case case_created_id_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases."case"
    ADD CONSTRAINT case_created_id_fk FOREIGN KEY (created_by) REFERENCES directory.wbt_user(id) ON DELETE SET NULL DEFERRABLE INITIALLY DEFERRED;


--
-- Name: case case_domain_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases."case"
    ADD CONSTRAINT case_domain_fk FOREIGN KEY (dc) REFERENCES directory.wbt_domain(dc) ON DELETE CASCADE DEFERRABLE INITIALLY DEFERRED;


--
-- Name: case case_group_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases."case"
    ADD CONSTRAINT case_group_fk FOREIGN KEY (contact_group) REFERENCES contacts."group"(id) ON DELETE SET NULL DEFERRABLE INITIALLY DEFERRED;


--
-- Name: case case_impacted_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases."case"
    ADD CONSTRAINT case_impacted_fk FOREIGN KEY (impacted) REFERENCES contacts.contact(id) ON DELETE SET NULL DEFERRABLE INITIALLY DEFERRED;


--
-- Name: case case_priority_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases."case"
    ADD CONSTRAINT case_priority_fk FOREIGN KEY (priority) REFERENCES cases.priority(id) ON DELETE RESTRICT DEFERRABLE INITIALLY DEFERRED;


--
-- Name: case case_reporter_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases."case"
    ADD CONSTRAINT case_reporter_fk FOREIGN KEY (reporter) REFERENCES contacts.contact(id) ON DELETE SET NULL DEFERRABLE INITIALLY DEFERRED;


--
-- Name: case case_service_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases."case"
    ADD CONSTRAINT case_service_fk FOREIGN KEY (service) REFERENCES cases.service_catalog(id) ON DELETE RESTRICT DEFERRABLE INITIALLY DEFERRED;


--
-- Name: case case_sla_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases."case"
    ADD CONSTRAINT case_sla_fk FOREIGN KEY (sla) REFERENCES cases.sla(id) ON DELETE RESTRICT DEFERRABLE INITIALLY DEFERRED;


--
-- Name: case case_source_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases."case"
    ADD CONSTRAINT case_source_fk FOREIGN KEY (source) REFERENCES cases.source(id) ON DELETE RESTRICT DEFERRABLE INITIALLY DEFERRED;


--
-- Name: case case_status_condition_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases."case"
    ADD CONSTRAINT case_status_condition_fk FOREIGN KEY (status_condition) REFERENCES cases.status_condition(id) ON DELETE SET NULL DEFERRABLE INITIALLY DEFERRED;


--
-- Name: case case_status_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases."case"
    ADD CONSTRAINT case_status_fk FOREIGN KEY (status) REFERENCES cases.status(id) ON DELETE RESTRICT DEFERRABLE INITIALLY DEFERRED;


--
-- Name: case case_updated_dc_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases."case"
    ADD CONSTRAINT case_updated_dc_fk FOREIGN KEY (updated_by, dc) REFERENCES directory.wbt_user(id, dc) DEFERRABLE INITIALLY DEFERRED;


--
-- Name: case case_updated_id_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases."case"
    ADD CONSTRAINT case_updated_id_fk FOREIGN KEY (updated_by) REFERENCES directory.wbt_user(id) ON DELETE SET NULL DEFERRABLE INITIALLY DEFERRED;


--
-- Name: close_reason_group close_reason_created_dc_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.close_reason_group
    ADD CONSTRAINT close_reason_created_dc_fk FOREIGN KEY (created_by, dc) REFERENCES directory.wbt_user(id, dc) DEFERRABLE INITIALLY DEFERRED;


--
-- Name: close_reason_group close_reason_created_id_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.close_reason_group
    ADD CONSTRAINT close_reason_created_id_fk FOREIGN KEY (created_by) REFERENCES directory.wbt_user(id) ON DELETE SET NULL DEFERRABLE INITIALLY DEFERRED;


--
-- Name: close_reason_group close_reason_domain_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.close_reason_group
    ADD CONSTRAINT close_reason_domain_fk FOREIGN KEY (dc) REFERENCES directory.wbt_domain(dc) ON DELETE CASCADE;


--
-- Name: close_reason_group close_reason_updated_dc_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.close_reason_group
    ADD CONSTRAINT close_reason_updated_dc_fk FOREIGN KEY (updated_by, dc) REFERENCES directory.wbt_user(id, dc) DEFERRABLE INITIALLY DEFERRED;


--
-- Name: close_reason_group close_reason_updated_id_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.close_reason_group
    ADD CONSTRAINT close_reason_updated_id_fk FOREIGN KEY (updated_by) REFERENCES directory.wbt_user(id) DEFERRABLE INITIALLY DEFERRED;


--
-- Name: case_comment comment_case_case_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.case_comment
    ADD CONSTRAINT comment_case_case_fk FOREIGN KEY (case_id) REFERENCES cases."case"(id) ON DELETE CASCADE DEFERRABLE INITIALLY DEFERRED;


--
-- Name: case_comment comment_case_created_dc_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.case_comment
    ADD CONSTRAINT comment_case_created_dc_fk FOREIGN KEY (created_by, dc) REFERENCES directory.wbt_user(id, dc) DEFERRABLE INITIALLY DEFERRED;


--
-- Name: case_comment comment_case_created_id_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.case_comment
    ADD CONSTRAINT comment_case_created_id_fk FOREIGN KEY (created_by) REFERENCES directory.wbt_user(id) ON DELETE SET NULL DEFERRABLE INITIALLY DEFERRED;


--
-- Name: case_comment comment_case_updated_dc_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.case_comment
    ADD CONSTRAINT comment_case_updated_dc_fk FOREIGN KEY (updated_by, dc) REFERENCES directory.wbt_user(id, dc) DEFERRABLE INITIALLY DEFERRED;


--
-- Name: case_comment comment_case_updated_id_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.case_comment
    ADD CONSTRAINT comment_case_updated_id_fk FOREIGN KEY (updated_by) REFERENCES directory.wbt_user(id) ON DELETE SET NULL DEFERRABLE INITIALLY DEFERRED;


--
-- Name: case_link link_case_case_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.case_link
    ADD CONSTRAINT link_case_case_fk FOREIGN KEY (case_id) REFERENCES cases."case"(id) ON DELETE CASCADE DEFERRABLE INITIALLY DEFERRED;


--
-- Name: case_link link_case_created_dc_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.case_link
    ADD CONSTRAINT link_case_created_dc_fk FOREIGN KEY (created_by, dc) REFERENCES directory.wbt_user(id, dc) DEFERRABLE INITIALLY DEFERRED;


--
-- Name: case_link link_case_created_id_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.case_link
    ADD CONSTRAINT link_case_created_id_fk FOREIGN KEY (created_by) REFERENCES directory.wbt_user(id) ON DELETE SET NULL DEFERRABLE INITIALLY DEFERRED;


--
-- Name: case_link link_case_updated_dc_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.case_link
    ADD CONSTRAINT link_case_updated_dc_fk FOREIGN KEY (updated_by, dc) REFERENCES directory.wbt_user(id, dc) DEFERRABLE INITIALLY DEFERRED;


--
-- Name: case_link link_case_updated_id_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.case_link
    ADD CONSTRAINT link_case_updated_id_fk FOREIGN KEY (updated_by) REFERENCES directory.wbt_user(id) ON DELETE SET NULL DEFERRABLE INITIALLY DEFERRED;


--
-- Name: priority priority_created_dc_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.priority
    ADD CONSTRAINT priority_created_dc_fk FOREIGN KEY (created_by, dc) REFERENCES directory.wbt_user(id, dc) DEFERRABLE INITIALLY DEFERRED;


--
-- Name: priority priority_created_id_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.priority
    ADD CONSTRAINT priority_created_id_fk FOREIGN KEY (created_by) REFERENCES directory.wbt_user(id) ON DELETE SET NULL DEFERRABLE INITIALLY DEFERRED;


--
-- Name: priority priority_domain_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.priority
    ADD CONSTRAINT priority_domain_fk FOREIGN KEY (dc) REFERENCES directory.wbt_domain(dc) ON DELETE CASCADE;


--
-- Name: priority_sla_condition priority_sla_condition_created_by_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.priority_sla_condition
    ADD CONSTRAINT priority_sla_condition_created_by_fk FOREIGN KEY (created_by) REFERENCES directory.wbt_user(id) ON DELETE SET NULL DEFERRABLE INITIALLY DEFERRED;


--
-- Name: priority_sla_condition priority_sla_condition_created_dc_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.priority_sla_condition
    ADD CONSTRAINT priority_sla_condition_created_dc_fk FOREIGN KEY (created_by, dc) REFERENCES directory.wbt_user(id, dc) DEFERRABLE INITIALLY DEFERRED;


--
-- Name: priority_sla_condition priority_sla_condition_domain_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.priority_sla_condition
    ADD CONSTRAINT priority_sla_condition_domain_fk FOREIGN KEY (dc) REFERENCES directory.wbt_domain(dc) ON DELETE CASCADE;


--
-- Name: priority_sla_condition priority_sla_condition_priority_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.priority_sla_condition
    ADD CONSTRAINT priority_sla_condition_priority_fk FOREIGN KEY (priority_id, dc) REFERENCES cases.priority(id, dc) ON DELETE RESTRICT DEFERRABLE INITIALLY DEFERRED;


--
-- Name: priority_sla_condition priority_sla_condition_sla_condition_id_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.priority_sla_condition
    ADD CONSTRAINT priority_sla_condition_sla_condition_id_fk FOREIGN KEY (sla_condition_id) REFERENCES cases.sla_condition(id) ON DELETE CASCADE;


--
-- Name: priority_sla_condition priority_sla_condition_updated_by_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.priority_sla_condition
    ADD CONSTRAINT priority_sla_condition_updated_by_fk FOREIGN KEY (updated_by) REFERENCES directory.wbt_user(id) ON DELETE SET NULL DEFERRABLE INITIALLY DEFERRED;


--
-- Name: priority_sla_condition priority_sla_condition_updated_dc_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.priority_sla_condition
    ADD CONSTRAINT priority_sla_condition_updated_dc_fk FOREIGN KEY (updated_by, dc) REFERENCES directory.wbt_user(id, dc) DEFERRABLE INITIALLY DEFERRED;


--
-- Name: priority priority_updated_dc_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.priority
    ADD CONSTRAINT priority_updated_dc_fk FOREIGN KEY (updated_by, dc) REFERENCES directory.wbt_user(id, dc) DEFERRABLE INITIALLY DEFERRED;


--
-- Name: priority priority_updated_id_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.priority
    ADD CONSTRAINT priority_updated_id_fk FOREIGN KEY (updated_by) REFERENCES directory.wbt_user(id) DEFERRABLE INITIALLY DEFERRED;


--
-- Name: close_reason reason_close_reason_id_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.close_reason
    ADD CONSTRAINT reason_close_reason_id_fk FOREIGN KEY (close_reason_id) REFERENCES cases.close_reason_group(id) ON DELETE CASCADE;


--
-- Name: close_reason reason_created_by_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.close_reason
    ADD CONSTRAINT reason_created_by_fk FOREIGN KEY (created_by) REFERENCES directory.wbt_user(id) ON DELETE SET NULL DEFERRABLE INITIALLY DEFERRED;


--
-- Name: close_reason reason_created_dc_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.close_reason
    ADD CONSTRAINT reason_created_dc_fk FOREIGN KEY (created_by, dc) REFERENCES directory.wbt_user(id, dc) DEFERRABLE INITIALLY DEFERRED;


--
-- Name: close_reason reason_domain_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.close_reason
    ADD CONSTRAINT reason_domain_fk FOREIGN KEY (close_reason_id, dc) REFERENCES cases.close_reason_group(id, dc) ON DELETE CASCADE DEFERRABLE INITIALLY DEFERRED;


--
-- Name: close_reason reason_updated_by_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.close_reason
    ADD CONSTRAINT reason_updated_by_fk FOREIGN KEY (updated_by) REFERENCES directory.wbt_user(id) ON DELETE SET NULL DEFERRABLE INITIALLY DEFERRED;


--
-- Name: close_reason reason_updated_dc_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.close_reason
    ADD CONSTRAINT reason_updated_dc_fk FOREIGN KEY (updated_by, dc) REFERENCES directory.wbt_user(id, dc) DEFERRABLE INITIALLY DEFERRED;


--
-- Name: related_case related_case_created_dc_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.related_case
    ADD CONSTRAINT related_case_created_dc_fk FOREIGN KEY (created_by, dc) REFERENCES directory.wbt_user(id, dc) DEFERRABLE INITIALLY DEFERRED;


--
-- Name: related_case related_case_created_id_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.related_case
    ADD CONSTRAINT related_case_created_id_fk FOREIGN KEY (created_by) REFERENCES directory.wbt_user(id) ON DELETE SET NULL DEFERRABLE INITIALLY DEFERRED;


--
-- Name: related_case related_case_primary_case_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.related_case
    ADD CONSTRAINT related_case_primary_case_fk FOREIGN KEY (primary_case_id) REFERENCES cases."case"(id) ON DELETE CASCADE DEFERRABLE INITIALLY DEFERRED;


--
-- Name: related_case related_case_related_case_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.related_case
    ADD CONSTRAINT related_case_related_case_fk FOREIGN KEY (related_case_id) REFERENCES cases."case"(id) ON DELETE CASCADE DEFERRABLE INITIALLY DEFERRED;


--
-- Name: related_case related_case_updated_dc_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.related_case
    ADD CONSTRAINT related_case_updated_dc_fk FOREIGN KEY (updated_by, dc) REFERENCES directory.wbt_user(id, dc) DEFERRABLE INITIALLY DEFERRED;


--
-- Name: related_case related_case_updated_id_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.related_case
    ADD CONSTRAINT related_case_updated_id_fk FOREIGN KEY (updated_by) REFERENCES directory.wbt_user(id) ON DELETE SET NULL DEFERRABLE INITIALLY DEFERRED;


--
-- Name: service_catalog service_catalog_created_dc_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.service_catalog
    ADD CONSTRAINT service_catalog_created_dc_fk FOREIGN KEY (created_by, dc) REFERENCES directory.wbt_user(id, dc) DEFERRABLE INITIALLY DEFERRED;


--
-- Name: service_catalog service_catalog_domain_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.service_catalog
    ADD CONSTRAINT service_catalog_domain_fk FOREIGN KEY (dc) REFERENCES directory.wbt_domain(dc) ON DELETE CASCADE;


--
-- Name: service_catalog service_catalog_root_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.service_catalog
    ADD CONSTRAINT service_catalog_root_fk FOREIGN KEY (root_id) REFERENCES cases.service_catalog(id) ON DELETE CASCADE;


--
-- Name: service_catalog service_catalog_updated_id_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.service_catalog
    ADD CONSTRAINT service_catalog_updated_id_fk FOREIGN KEY (updated_by) REFERENCES directory.wbt_user(id) DEFERRABLE INITIALLY DEFERRED;


--
-- Name: service_catalog service_created_id_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.service_catalog
    ADD CONSTRAINT service_created_id_fk FOREIGN KEY (created_by) REFERENCES directory.wbt_user(id) ON DELETE SET NULL DEFERRABLE INITIALLY DEFERRED;


--
-- Name: service_catalog service_updated_dc_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.service_catalog
    ADD CONSTRAINT service_updated_dc_fk FOREIGN KEY (updated_by, dc) REFERENCES directory.wbt_user(id, dc) DEFERRABLE INITIALLY DEFERRED;


--
-- Name: skill_catalog skill_service_created_by_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.skill_catalog
    ADD CONSTRAINT skill_service_created_by_fk FOREIGN KEY (created_by) REFERENCES directory.wbt_user(id) ON DELETE SET NULL DEFERRABLE INITIALLY DEFERRED;


--
-- Name: skill_catalog skill_service_created_dc_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.skill_catalog
    ADD CONSTRAINT skill_service_created_dc_fk FOREIGN KEY (created_by, dc) REFERENCES directory.wbt_user(id, dc) DEFERRABLE INITIALLY DEFERRED;


--
-- Name: skill_catalog skill_service_domain_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.skill_catalog
    ADD CONSTRAINT skill_service_domain_fk FOREIGN KEY (dc) REFERENCES directory.wbt_domain(dc) ON DELETE CASCADE;


--
-- Name: skill_catalog skill_service_service_id_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.skill_catalog
    ADD CONSTRAINT skill_service_service_id_fk FOREIGN KEY (catalog_id) REFERENCES cases.service_catalog(id) ON DELETE CASCADE;


--
-- Name: skill_catalog skill_service_updated_by_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.skill_catalog
    ADD CONSTRAINT skill_service_updated_by_fk FOREIGN KEY (updated_by) REFERENCES directory.wbt_user(id) ON DELETE SET NULL DEFERRABLE INITIALLY DEFERRED;


--
-- Name: skill_catalog skill_service_updated_dc_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.skill_catalog
    ADD CONSTRAINT skill_service_updated_dc_fk FOREIGN KEY (updated_by, dc) REFERENCES directory.wbt_user(id, dc) DEFERRABLE INITIALLY DEFERRED;


--
-- Name: sla_condition sla_condition_created_dc_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.sla_condition
    ADD CONSTRAINT sla_condition_created_dc_fk FOREIGN KEY (created_by, dc) REFERENCES directory.wbt_user(id, dc) DEFERRABLE INITIALLY DEFERRED;


--
-- Name: sla_condition sla_condition_created_id_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.sla_condition
    ADD CONSTRAINT sla_condition_created_id_fk FOREIGN KEY (created_by) REFERENCES directory.wbt_user(id) ON DELETE SET NULL DEFERRABLE INITIALLY DEFERRED;


--
-- Name: sla_condition sla_condition_domain_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.sla_condition
    ADD CONSTRAINT sla_condition_domain_fk FOREIGN KEY (dc) REFERENCES directory.wbt_domain(dc) ON DELETE CASCADE;


--
-- Name: sla_condition sla_condition_sla_id_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.sla_condition
    ADD CONSTRAINT sla_condition_sla_id_fk FOREIGN KEY (sla_id) REFERENCES cases.sla(id) ON DELETE CASCADE;


--
-- Name: sla_condition sla_condition_updated_dc_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.sla_condition
    ADD CONSTRAINT sla_condition_updated_dc_fk FOREIGN KEY (updated_by, dc) REFERENCES directory.wbt_user(id, dc) DEFERRABLE INITIALLY DEFERRED;


--
-- Name: sla_condition sla_condition_updated_id_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.sla_condition
    ADD CONSTRAINT sla_condition_updated_id_fk FOREIGN KEY (updated_by) REFERENCES directory.wbt_user(id) DEFERRABLE INITIALLY DEFERRED;


--
-- Name: sla sla_created_dc_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.sla
    ADD CONSTRAINT sla_created_dc_fk FOREIGN KEY (created_by, dc) REFERENCES directory.wbt_user(id, dc) DEFERRABLE INITIALLY DEFERRED;


--
-- Name: sla sla_created_id_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.sla
    ADD CONSTRAINT sla_created_id_fk FOREIGN KEY (created_by) REFERENCES directory.wbt_user(id) ON DELETE SET NULL DEFERRABLE INITIALLY DEFERRED;


--
-- Name: sla sla_domain_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.sla
    ADD CONSTRAINT sla_domain_fk FOREIGN KEY (dc) REFERENCES directory.wbt_domain(dc) ON DELETE CASCADE;


--
-- Name: sla sla_updated_dc_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.sla
    ADD CONSTRAINT sla_updated_dc_fk FOREIGN KEY (updated_by, dc) REFERENCES directory.wbt_user(id, dc) DEFERRABLE INITIALLY DEFERRED;


--
-- Name: sla sla_updated_id_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.sla
    ADD CONSTRAINT sla_updated_id_fk FOREIGN KEY (updated_by) REFERENCES directory.wbt_user(id) DEFERRABLE INITIALLY DEFERRED;


--
-- Name: source source_created_dc_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.source
    ADD CONSTRAINT source_created_dc_fk FOREIGN KEY (created_by, dc) REFERENCES directory.wbt_user(id, dc) DEFERRABLE INITIALLY DEFERRED;


--
-- Name: source source_created_id_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.source
    ADD CONSTRAINT source_created_id_fk FOREIGN KEY (created_by) REFERENCES directory.wbt_user(id) ON DELETE SET NULL DEFERRABLE INITIALLY DEFERRED;


--
-- Name: source source_domain_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.source
    ADD CONSTRAINT source_domain_fk FOREIGN KEY (dc) REFERENCES directory.wbt_domain(dc) ON DELETE CASCADE;


--
-- Name: source source_updated_dc_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.source
    ADD CONSTRAINT source_updated_dc_fk FOREIGN KEY (updated_by, dc) REFERENCES directory.wbt_user(id, dc) DEFERRABLE INITIALLY DEFERRED;


--
-- Name: source source_updated_id_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.source
    ADD CONSTRAINT source_updated_id_fk FOREIGN KEY (updated_by) REFERENCES directory.wbt_user(id) DEFERRABLE INITIALLY DEFERRED;


--
-- Name: status_condition status_condition_created_by_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.status_condition
    ADD CONSTRAINT status_condition_created_by_fk FOREIGN KEY (created_by) REFERENCES directory.wbt_user(id) ON DELETE SET NULL DEFERRABLE INITIALLY DEFERRED;


--
-- Name: status_condition status_condition_created_dc_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.status_condition
    ADD CONSTRAINT status_condition_created_dc_fk FOREIGN KEY (created_by, dc) REFERENCES directory.wbt_user(id, dc) DEFERRABLE INITIALLY DEFERRED;


--
-- Name: status_condition status_condition_status_id_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.status_condition
    ADD CONSTRAINT status_condition_status_id_fk FOREIGN KEY (status_id) REFERENCES cases.status(id) ON DELETE CASCADE;


--
-- Name: status_condition status_condition_updated_by_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.status_condition
    ADD CONSTRAINT status_condition_updated_by_fk FOREIGN KEY (updated_by) REFERENCES directory.wbt_user(id) ON DELETE SET NULL DEFERRABLE INITIALLY DEFERRED;


--
-- Name: status_condition status_condition_updated_dc_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.status_condition
    ADD CONSTRAINT status_condition_updated_dc_fk FOREIGN KEY (updated_by, dc) REFERENCES directory.wbt_user(id, dc) DEFERRABLE INITIALLY DEFERRED;


--
-- Name: status status_created_id_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.status
    ADD CONSTRAINT status_created_id_fk FOREIGN KEY (created_by) REFERENCES directory.wbt_user(id) ON DELETE SET NULL DEFERRABLE INITIALLY DEFERRED;


--
-- Name: status status_domain_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.status
    ADD CONSTRAINT status_domain_fk FOREIGN KEY (dc) REFERENCES directory.wbt_domain(dc) ON DELETE CASCADE;


--
-- Name: status status_updated_dc_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.status
    ADD CONSTRAINT status_updated_dc_fk FOREIGN KEY (updated_by, dc) REFERENCES directory.wbt_user(id, dc) DEFERRABLE INITIALLY DEFERRED;


--
-- Name: status status_updated_id_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.status
    ADD CONSTRAINT status_updated_id_fk FOREIGN KEY (updated_by) REFERENCES directory.wbt_user(id) DEFERRABLE INITIALLY DEFERRED;


--
-- Name: team_catalog team_service_created_by_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.team_catalog
    ADD CONSTRAINT team_service_created_by_fk FOREIGN KEY (created_by) REFERENCES directory.wbt_user(id) ON DELETE SET NULL DEFERRABLE INITIALLY DEFERRED;


--
-- Name: team_catalog team_service_created_dc_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.team_catalog
    ADD CONSTRAINT team_service_created_dc_fk FOREIGN KEY (created_by, dc) REFERENCES directory.wbt_user(id, dc) DEFERRABLE INITIALLY DEFERRED;


--
-- Name: team_catalog team_service_domain_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.team_catalog
    ADD CONSTRAINT team_service_domain_fk FOREIGN KEY (dc) REFERENCES directory.wbt_domain(dc) ON DELETE CASCADE;


--
-- Name: team_catalog team_service_service_id_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.team_catalog
    ADD CONSTRAINT team_service_service_id_fk FOREIGN KEY (catalog_id) REFERENCES cases.service_catalog(id) ON DELETE CASCADE;


--
-- Name: team_catalog team_service_updated_by_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.team_catalog
    ADD CONSTRAINT team_service_updated_by_fk FOREIGN KEY (updated_by) REFERENCES directory.wbt_user(id) ON DELETE SET NULL DEFERRABLE INITIALLY DEFERRED;


--
-- Name: team_catalog team_service_updated_dc_fk; Type: FK CONSTRAINT; Schema: cases; Owner: -
--

ALTER TABLE ONLY cases.team_catalog
    ADD CONSTRAINT team_service_updated_dc_fk FOREIGN KEY (updated_by, dc) REFERENCES directory.wbt_user(id, dc) DEFERRABLE INITIALLY DEFERRED;


--
-- PostgreSQL database dump complete
--

