alter table cases.sla
    add reaction_time bigint default 0 not null;

alter table cases.sla
    add resolution_time bigint default 0 not null;



