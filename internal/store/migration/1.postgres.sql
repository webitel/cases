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

create trigger enforce_single_initial_trigger
    before insert or update
    on cases.status_condition
    for each row
execute procedure public.enforce_single_initial();

create trigger constraints_trigger
    before insert or update or delete
    on cases.status_condition
    for each row
execute procedure cases.conditions_constraints_checker();
