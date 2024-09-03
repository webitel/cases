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
            references cases.status
            on delete cascade,
    initial     boolean                                                          not null,
    final       boolean                                                          not null,
    dc          bigint                                                           not null,
    constraint status_condition_fk
        unique (id, dc),
    constraint status_condition_created_dc_fk
        foreign key (created_by, dc) references directory.wbt_user ()
            deferrable initially deferred,
    constraint status_condition_domain_fk
        foreign key (status_id, dc) references cases.status (id, dc)
            on delete cascade
            deferrable initially deferred,
    constraint status_condition_updated_dc_fk
        foreign key (updated_by, dc) references directory.wbt_user ()
            deferrable initially deferred
);

alter table cases.status_condition
    owner to opensips;

create index status_condition_source
    on cases.status_condition (status_id);

create index status_condition_dc
    on cases.status_condition (dc);

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

create table cases.reason
(
    id              bigint    default nextval('cases.reason_id'::regclass) not null
        constraint reason_pk
            primary key,
    name            text                                                   not null,
    description     text                                                   not null,
    created_at      timestamp default timezone('utc'::text, now())         not null,
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
            references cases.close_reason
            on delete cascade,
    dc              bigint                                                 not null,
    constraint reason_fk
        unique (id, dc),
    constraint reason_created_dc_fk
        foreign key (created_by, dc) references directory.wbt_user ()
            deferrable initially deferred,
    constraint reason_domain_fk
        foreign key (close_reason_id, dc) references cases.close_reason (id, dc)
            on delete cascade
            deferrable initially deferred,
    constraint reason_updated_dc_fk
        foreign key (updated_by, dc) references directory.wbt_user ()
            deferrable initially deferred
);

alter table cases.reason
    owner to opensips;

create index reason_source
    on cases.reason (close_reason_id);

create index reason_id_dc
    on cases.reason (dc);

create table cases.appeal
(
    id          bigint    default nextval('cases.appeal_id'::regclass) not null
        constraint appeal_pk
            primary key,
    name        text                                                   not null,
    description text                                                   not null,
    created_at  timestamp default timezone('utc'::text, now())         not null,
    updated_at  timestamp default timezone('utc'::text, now())         not null,
    created_by  bigint                                                 not null
        constraint appeal_created_id_fk
            references directory.wbt_user
            on delete set null
            deferrable initially deferred,
    updated_by  bigint                                                 not null
        constraint appeal_updated_id_fk
            references directory.wbt_user
            deferrable initially deferred,
    type        text                                                   not null,
    dc          bigint                                                 not null
        constraint appeal_domain_fk
            references directory.wbt_domain
            on delete cascade,
    constraint apppeal_fk
        unique (id, dc),
    constraint appeal_created_dc_fk
        foreign key (created_by, dc) references directory.wbt_user ()
            deferrable initially deferred,
    constraint appeal_updated_dc_fk
        foreign key (updated_by, dc) references directory.wbt_user ()
            deferrable initially deferred
);

alter table cases.appeal
    owner to opensips;

create index appeal_dc
    on cases.appeal (dc);

create table cases.priority
(
    id          bigint    default nextval('cases.priority_id'::regclass) not null
        constraint priority_pk
            primary key,
    name        text                                                     not null,
    description text                                                     not null,
    created_at  timestamp default timezone('utc'::text, now())           not null,
    updated_at  timestamp default timezone('utc'::text, now())           not null,
    created_by  bigint                                                   not null
        constraint priority_created_id_fk
            references directory.wbt_user
            on delete set null
            deferrable initially deferred,
    updated_by  bigint                                                   not null
        constraint priority_updated_id_fk
            references directory.wbt_user
            deferrable initially deferred,
    dc          bigint                                                   not null
        constraint priority_domain_fk
            references directory.wbt_domain
            on delete cascade,
    color       varchar(25)                                              not null,
    constraint priority_fk
        unique (id, dc),
    constraint priority_created_dc_fk
        foreign key (created_by, dc) references directory.wbt_user ()
            deferrable initially deferred,
    constraint priority_updated_dc_fk
        foreign key (updated_by, dc) references directory.wbt_user ()
            deferrable initially deferred
);

alter table cases.priority
    owner to opensips;

create index priority_dc
    on cases.priority (dc);

create table cases.sla
(
    id                      bigint    default nextval('cases.sla_id'::regclass) not null
        constraint sla_pk
            primary key,
    name                    text                                                not null,
    description             text,
    created_at              timestamp default timezone('utc'::text, now())      not null,
    updated_at              timestamp default timezone('utc'::text, now())      not null,
    created_by              bigint                                              not null
        constraint sla_created_id_fk
            references directory.wbt_user
            on delete set null
            deferrable initially deferred,
    updated_by              bigint                                              not null
        constraint sla_updated_id_fk
            references directory.wbt_user
            deferrable initially deferred,
    dc                      bigint                                              not null
        constraint sla_domain_fk
            references directory.wbt_domain
            on delete cascade,
    valid_from              timestamp default timezone('utc'::text, now()),
    valid_to                timestamp default timezone('utc'::text, now()),
    reaction_time_hours     integer   default 0                                 not null,
    reaction_time_minutes   integer   default 0,
    resolution_time_hours   integer   default 0                                 not null,
    resolution_time_minutes integer   default 0                                 not null,
    calendar_id             bigint                                              not null,
    constraint sla_fk
        unique (id, dc),
    constraint sla_created_dc_fk
        foreign key (created_by, dc) references directory.wbt_user ()
            deferrable initially deferred,
    constraint sla_updated_dc_fk
        foreign key (updated_by, dc) references directory.wbt_user ()
            deferrable initially deferred
);

alter table cases.sla
    owner to opensips;

create index sla_dc
    on cases.sla (dc);

create table cases.sla_condition
(
    id                      bigint    default nextval('cases.sla_condition_id'::regclass) not null
        constraint sla_condition_pk
            primary key,
    name                    text                                                          not null,
    created_at              timestamp default timezone('utc'::text, now())                not null,
    updated_at              timestamp default timezone('utc'::text, now())                not null,
    created_by              bigint                                                        not null
        constraint sla_condition_created_id_fk
            references directory.wbt_user
            on delete set null
            deferrable initially deferred,
    updated_by              bigint                                                        not null
        constraint sla_condition_updated_id_fk
            references directory.wbt_user
            deferrable initially deferred,
    dc                      bigint                                                        not null
        constraint sla_condition_domain_fk
            references directory.wbt_domain
            on delete cascade,
    reaction_time_hours     integer   default 0                                           not null,
    reaction_time_minutes   integer   default 0                                           not null,
    resolution_time_hours   integer   default 0                                           not null,
    resolution_time_minutes integer   default 0                                           not null,
    sla_id                  bigint                                                        not null
        constraint sla_condition_sla_id_fk
            references cases.sla
            on delete cascade,
    constraint sla_condition_fk
        unique (id, dc),
    constraint sla_condition_created_dc_fk
        foreign key (created_by, dc) references directory.wbt_user ()
            deferrable initially deferred,
    constraint sla_condition_updated_dc_fk
        foreign key (updated_by, dc) references directory.wbt_user ()
            deferrable initially deferred
);

alter table cases.sla_condition
    owner to opensips;

create index sla_condition_dc
    on cases.sla_condition (dc);

create index sla_condition_source
    on cases.sla_condition (sla_id);

create table cases.priority_sla_condition
(
    id               bigint    default nextval('cases.priority_sla_condition_id'::regclass) not null
        constraint priority_sla_condition_pk
            primary key,
    created_at       timestamp default timezone('utc'::text, now())                         not null,
    updated_at       timestamp default timezone('utc'::text, now())                         not null,
    created_by       bigint                                                                 not null
        constraint priority_sla_condition_created_by_fk
            references directory.wbt_user
            on delete set null
            deferrable initially deferred,
    updated_by       bigint                                                                 not null
        constraint priority_sla_condition_updated_by_fk
            references directory.wbt_user
            on delete set null
            deferrable initially deferred,
    sla_condition_id bigint                                                                 not null
        constraint priority_sla_condition_sla_condition_id_fk
            references cases.sla_condition
            on delete cascade,
    priority_id      bigint                                                                 not null,
    dc               bigint                                                                 not null
        constraint priority_sla_condition_domain_fk
            references directory.wbt_domain
            on delete cascade,
    constraint priority_sla_condition_fk
        unique (id, dc),
    constraint priority_sla_condition_created_dc_fk
        foreign key (created_by, dc) references directory.wbt_user ()
            deferrable initially deferred,
    constraint priority_sla_condition_updated_dc_fk
        foreign key (updated_by, dc) references directory.wbt_user ()
            deferrable initially deferred
);

alter table cases.priority_sla_condition
    owner to opensips;

create index priority_sla_condition_dc
    on cases.priority_sla_condition (dc);

create index priority_sla_condition_source
    on cases.priority_sla_condition (sla_condition_id);
