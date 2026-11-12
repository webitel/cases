-- WTEL-8579: comment type, 1=internal, 2=external.
alter table cases.case_comment
    add column if not exists type smallint default 1 not null;

drop trigger if exists tg_case_comment_rbac on cases.case_comment;
create trigger tg_case_comment_rbac
    after insert
    on cases.case_comment
    for each row
    when (new.type = 2)
execute function directory.tg_obj_default_rbac('case_comments');

drop trigger if exists tg_internal_comment_rbac on cases.case_comment;
create trigger tg_internal_comment_rbac
    after insert
    on cases.case_comment
    for each row
    when (new.type <> 2)
execute function directory.tg_obj_default_rbac('internal_comments');
