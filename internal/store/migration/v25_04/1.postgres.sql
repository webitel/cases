alter table cases.case_communication
    add constraint case_communication_cc_communication_id_fk
        foreign key (communication_type) references call_center.cc_communication
            on delete cascade
            deferrable initially deferred;

