create table if not exists cases.case_communication
(
    id                 bigserial,
    ver                integer default 0 not null,
    created_by bigint,
    created_at timestamp default now() not null,
    case_id            bigint            not null,
    communication_type int              not null,
    communication_id   text              not null,
    dc                 bigint            not null
);

comment on table cases.case_communication is 'Table connects case with other possible communications like chats, emails, calls';

comment on column cases.case_communication.communication_type is 'type of communication case connected to';

comment on column cases.case_communication.communication_id is 'Id of communication case connected to';

comment on column cases.case_communication.dc is 'Domain component';

alter table cases.case_communication
    owner to opensips;

alter table cases.case_communication
    add constraint case_communication_case_id_fk
        foreign key (case_id) references cases."case"
            on delete cascade;
create unique index unique_event_uindex
    on cases.case_communication (communication_id, communication_type, case_id);

alter table cases.case_communication
    add constraint case_communication_wbt_domain_dc_fk
        foreign key (dc) references directory.wbt_domain
            on delete cascade;

alter table cases.case_communication
    add created_by bigint not null;

alter table cases.case_communication
    add constraint case_communication_wbt_user_id_fk
        foreign key (created_by) references directory.wbt_user;

alter table contacts.dynamic_group
    add constraint dynamic_group_default_group_id_fk
        foreign key (default_group_id) references contacts."group"
            on update cascade on delete restrict;

alter table cases.service_catalog
    rename column close_reason_id to close_reason_group_id;


ALTER TABLE cases.status_condition
ADD CONSTRAINT check_initial_final
CHECK (NOT (initial = TRUE AND final = TRUE));