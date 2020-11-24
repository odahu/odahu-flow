BEGIN;
alter table odahu.public.odahu_operator_route drop column created;
alter table odahu.public.odahu_operator_route drop column updated;
alter table odahu.public.odahu_operator_route drop column deletionmark;
COMMIT;