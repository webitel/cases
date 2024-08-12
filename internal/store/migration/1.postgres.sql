create table cases.status
(
    id          bigint    default nextval('cases.status_id'::regclass) not null
        constraint status_pk
            primary key,
    name        text                                                   not null,
    description text,
    created_at  timestamp default timezone('utc'::text, now())         not null,
    updated_at  timestamp default timezone('utc'::text, now())         not null,
    created_by  bigint                                                 not null
        constraint status_created_id_fk
            references directory.wbt_user
            on delete set null
            deferrable initially deferred,
    updated_by  bigint
        constraint status_updated_id_fk
            references directory.wbt_user
            deferrable initially deferred,
    dc          bigint                                                 not null
        constraint status_domain_fk
            references directory.wbt_domain
            on delete cascade,
    constraint status_fk
        unique (id, dc),
    constraint status_created_dc_fk
        foreign key (created_by, dc) references directory.wbt_user ()
            deferrable initially deferred,
    constraint status_updated_dc_fk
        foreign key (updated_by, dc) references directory.wbt_user ()
            deferrable initially deferred
);

alter table cases.status
    owner to opensips;

create index status_dc
    on cases.status (dc);

create trigger tg_cases_status_rbac
    after insert
    on cases.status
    referencing new table inserted
execute procedure directory.tg_stmt_default_rbac('case_statuses');


create table cases.status_condition
(
    id          bigint    default nextval('cases.status_condition_id'::regclass) not null
        constraint status_condition_pk
            primary key,
    name        text                                                             not null,
    description text                                                             not null,
    created_at  timestamp default timezone('utc'::text, now())                   not null,
    updated_at  timestamp default timezone('utc'::text, now())                   not null,
    created_by  bigint
        constraint status_condition_created_by_fk
            references directory.wbt_user
            on delete set null
            deferrable initially deferred,
    updated_by  bigint
        constraint status_condition_updated_by_fk
            references directory.wbt_user
            on delete set null
            deferrable initially deferred,
    status_id   bigint                                                           not null
        constraint status_condition_status_id_fk
            references cases.status,
    initial     boolean                                                          not null,
    final       boolean                                                          not null,
    dc          bigint                                                           not null,
    constraint status_condition_created_dc_fk
        foreign key (created_by, dc) references directory.wbt_user ()
            deferrable initially deferred,
    constraint status_condition_fk
        foreign key (status_id, dc) references cases.status (id, dc)
            on delete cascade
            deferrable initially deferred,
    constraint status_condition_updated_dc_fk
        foreign key (updated_by, dc) references directory.wbt_user ()
            deferrable initially deferred
);

alter table cases.status_condition
    owner to opensips;

create index status_condition_source_index
    on cases.status_condition (status_id);



create table cases.status_acl
(
    dc      bigint not null
        constraint status_acl_domain_fk
            references directory.wbt_domain
            on delete cascade,
    grantor bigint
        constraint status_acl_grantor_id_fk
            references directory.wbt_user
            on delete set null,
    subject bigint not null,
    object  bigint not null,
    access  smallint,
    constraint status_acl_grantor_fk
        foreign key (grantor, dc) references directory.wbt_user ()
            deferrable initially deferred,
    constraint status_acl_object_fk
        foreign key (object, dc) references cases.status (id, dc)
            on delete cascade
            deferrable initially deferred,
    constraint status_acl_subject_fk
        foreign key (subject, dc) references directory.wbt_user ()
            on delete cascade
);

alter table cases.status_acl
    owner to opensips;

create unique index status_acl_object_subject_pk
    on cases.status_acl (object, subject);

create unique index status_acl_subject_object_uindex
    on cases.status_acl (subject, object);

create index status_acl_grantor_index
    on cases.status_acl (grantor);



create table cases.close_reason
(
    id          bigint    default nextval('cases.close_reason_id'::regclass) not null
        constraint close_reason_pk
            primary key,
    name        text                                                         not null,
    description text                                                         not null,
    created_at  timestamp default timezone('utc'::text, now())               not null,
    updated_at  timestamp default timezone('utc'::text, now())               not null,
    created_by  bigint                                                       not null
        constraint close_reason_created_id_fk
            references directory.wbt_user
            on delete set null
            deferrable initially deferred,
    updated_by  bigint                                                       not null
        constraint close_reason_updated_id_fk
            references directory.wbt_user
            deferrable initially deferred,
    dc          bigint                                                       not null
        constraint close_reason_domain_fk
            references directory.wbt_domain
            on delete cascade,
    constraint close_reason_fk
        unique (id, dc),
    constraint close_reason_created_dc_fk
        foreign key (created_by, dc) references directory.wbt_user ()
            deferrable initially deferred,
    constraint close_reason_updated_dc_fk
        foreign key (updated_by, dc) references directory.wbt_user ()
            deferrable initially deferred
);

alter table cases.close_reason
    owner to opensips;

create index close_reason_dc
    on cases.close_reason (dc);

create trigger tg_cases_close_reason_rbac
    after insert
    on cases.close_reason
    referencing new table inserted
execute procedure directory.tg_stmt_default_rbac('close_reasons');

create table cases.reason
(
    id              bigint    default nextval('cases.reason_id'::regclass) not null
        constraint reason_pk
            primary key,
    name            text                                                   not null,
    description     text                                                   not null,
    crated_at       timestamp default timezone('utc'::text, now())         not null,
    updated_at      timestamp default timezone('utc'::text, now())         not null,
    created_by      bigint
        constraint reason_created_by_fk
            references directory.wbt_user
            on delete set null
            deferrable initially deferred,
    updated_by      bigint
        constraint reason_updated_by_fk
            references directory.wbt_user
            on delete set null
            deferrable initially deferred,
    close_reason_id bigint                                                 not null
        constraint reason_close_reason_id_fk
            references cases.close_reason,
    dc              bigint                                                 not null,
    constraint reason_created_dc_fk
        foreign key (created_by, dc) references directory.wbt_user ()
            deferrable initially deferred,
    constraint reason_fk
        foreign key (close_reason_id, dc) references cases.close_reason (id, dc)
            on delete cascade
            deferrable initially deferred,
    constraint reason_updated_dc_fk
        foreign key (updated_by, dc) references directory.wbt_user ()
            deferrable initially deferred
);

alter table cases.reason
    owner to opensips;

create index reason_source_index
    on cases.reason (close_reason_id);


create table cases.close_reason_acl
(
    dc      bigint not null
        constraint close_reason_acl_domain_fk
            references directory.wbt_domain
            on delete cascade,
    grantor bigint
        constraint close_reason_acl_grantor_id_fk
            references directory.wbt_user
            on delete set null,
    subject bigint not null,
    object  bigint not null,
    access  smallint,
    constraint close_reason_acl_grantor_fk
        foreign key (grantor, dc) references directory.wbt_user ()
            deferrable initially deferred,
    constraint close_reason_acl_object_fk
        foreign key (object, dc) references cases.close_reason (id, dc)
            on delete cascade
            deferrable initially deferred,
    constraint close_reason_acl_subject_fk
        foreign key (subject, dc) references directory.wbt_user ()
            on delete cascade
);

alter table cases.close_reason_acl
    owner to opensips;

create unique index close_reason_acl_object_subject_pk
    on cases.close_reason_acl (object, subject);

create unique index close_reason_acl_subject_object_uindex
    on cases.close_reason_acl (subject, object);

create index close_reason_acl_grantor_index
    on cases.close_reason_acl (grantor);
