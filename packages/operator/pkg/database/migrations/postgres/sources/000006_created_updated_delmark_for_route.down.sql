BEGIN;
alter table odahu_operator_route drop column created;
alter table odahu_operator_route drop column updated;
alter table odahu_operator_route drop column deletionmark;
COMMIT;