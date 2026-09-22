-- Knowledge base articles linked to a case: manual links and the close article.
create table if not exists cases.case_article
(
    dc         bigint not null,
    case_id    bigint not null,
    article_id bigint not null,
    source     smallint default 1 not null, -- 1=manual, 2=resolution
    created_by bigint,
    created_at timestamp without time zone default timezone('utc'::text, now()) not null,
    constraint case_article_pk primary key (case_id, article_id),
    constraint case_article_case_fk foreign key (case_id)
        references cases."case" (id) on delete cascade deferrable initially deferred,
    constraint case_article_created_id_fk foreign key (created_by)
        references directory.wbt_user (id) on delete set null deferrable initially deferred,
    constraint case_article_created_dc_fk foreign key (created_by, dc)
        references directory.wbt_user (id, dc) deferrable initially deferred
);

create index if not exists case_article_article_id_index
    on cases.case_article (article_id);

-- One close article per case.
create unique index if not exists case_article_resolution_uindex
    on cases.case_article (case_id)
    where source = 2;
