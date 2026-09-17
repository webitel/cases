-- Knowledge base article chosen as the close reason.
alter table cases."case" add column if not exists close_article_id bigint;

create index if not exists case_close_article_id_index
    on cases."case" (close_article_id)
    where close_article_id is not null;

-- Knowledge base articles linked to a case by an operator.
create table if not exists cases.case_article
(
    dc         bigint not null,
    case_id    bigint not null,
    article_id bigint not null,
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
