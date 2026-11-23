alter table cases.status_condition
    add column if not exists archived boolean default false not null;

alter table cases.status_condition
    drop constraint if exists status_condition_archived_final_check;

alter table cases.status_condition
    add constraint status_condition_archived_final_check
        check (not archived or final);
