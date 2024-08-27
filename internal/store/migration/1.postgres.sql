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



create table cases.close_reason (
  id bigint primary key not null default nextval('close_reason_id'::regclass),
  name text not null,
  description text not null,
  created_at timestamp without time zone not null default timezone('utc'::text, now()),
  updated_at timestamp without time zone not null default timezone('utc'::text, now()),
  created_by bigint not null,
  updated_by bigint not null,
  dc bigint not null,
  foreign key (created_by, dc) references directory.wbt_user (id, dc)
  match simple on update no action on delete no action,
  foreign key (created_by) references directory.wbt_user (id)
  match simple on update no action on delete set null,
  foreign key (dc) references directory.wbt_domain (dc)
  match simple on update no action on delete cascade,
  foreign key (updated_by, dc) references directory.wbt_user (id, dc)
  match simple on update no action on delete no action,
  foreign key (updated_by) references directory.wbt_user (id)
  match simple on update no action on delete no action
);
create index close_reason_dc on close_reason using btree (dc);
create unique index close_reason_fk on close_reason using btree (id, dc);



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



create table cases.reason (
  id bigint primary key not null default nextval('reason_id'::regclass),
  name text not null,
  description text not null,
  created_at timestamp without time zone not null default timezone('utc'::text, now()),
  updated_at timestamp without time zone not null default timezone('utc'::text, now()),
  created_by bigint,
  updated_by bigint,
  close_reason_id bigint not null,
  dc bigint not null,
  foreign key (close_reason_id) references cases.close_reason (id)
  match simple on update no action on delete no action,
  foreign key (created_by) references directory.wbt_user (id)
  match simple on update no action on delete set null,
  foreign key (created_by, dc) references directory.wbt_user (id, dc)
  match simple on update no action on delete no action,
  foreign key (close_reason_id, dc) references cases.close_reason (id, dc)
  match simple on update no action on delete cascade,
  foreign key (updated_by) references directory.wbt_user (id)
  match simple on update no action on delete set null,
  foreign key (updated_by, dc) references directory.wbt_user (id, dc)
  match simple on update no action on delete no action
);
create index reason_source_index on reason using btree (close_reason_id);



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
