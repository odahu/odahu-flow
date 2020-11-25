BEGIN;
alter table odahu_operator_route drop column default;
COMMIT;