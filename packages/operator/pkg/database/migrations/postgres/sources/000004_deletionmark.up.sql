BEGIN;
alter table odahu_operator_training
    add "deletionmark" boolean default FALSE not null;
alter table odahu_operator_packaging
    add "deletionmark" boolean default FALSE not null;
alter table odahu_operator_deployment
    add "deletionmark" boolean default FALSE not null;
COMMIT;