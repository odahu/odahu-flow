BEGIN;
alter table odahu_operator_training
    add created timestamptz;
alter table odahu_operator_training
    add updated timestamptz;
alter table odahu_operator_packaging
    add created timestamptz;
alter table odahu_operator_packaging
    add updated timestamptz;
alter table odahu_operator_deployment
    add created timestamptz;
alter table odahu_operator_deployment
    add updated timestamptz;
COMMIT;