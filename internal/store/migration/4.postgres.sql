alter table cases.sla
    alter column reaction_time type bigint using reaction_time::bigint;

alter table cases.sla
    alter column resolution_time type bigint using resolution_time::bigint;

