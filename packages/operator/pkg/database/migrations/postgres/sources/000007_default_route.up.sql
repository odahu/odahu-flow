BEGIN;
alter table odahu_operator_route
    add "default" boolean default FALSE not null;
COMMIT;