BEGIN;
alter table odahu.public.odahu_operator_route
    add created timestamptz;
alter table odahu.public.odahu_operator_route
    add updated timestamptz;
alter table odahu.public.odahu_operator_route
    add "deletionmark" boolean default FALSE not null;
COMMIT;