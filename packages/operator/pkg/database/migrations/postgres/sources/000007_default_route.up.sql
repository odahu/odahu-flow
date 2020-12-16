BEGIN;
alter table odahu_operator_route
    add is_default boolean default FALSE not null;
COMMIT;