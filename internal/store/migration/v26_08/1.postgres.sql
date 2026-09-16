-- Knowledge base article chosen as the close reason.
alter table cases."case" add column if not exists close_article_id bigint;

-- Drop a leftover of an interrupted concurrent build.
do $$
begin
    if exists (select 1 from pg_index
               where indexrelid = to_regclass('cases.case_close_article_id_index') and not indisvalid) then
        drop index cases.case_close_article_id_index;
    end if;
end
$$;

create index concurrently if not exists case_close_article_id_index
    on cases."case" (close_article_id)
    where close_article_id is not null;
