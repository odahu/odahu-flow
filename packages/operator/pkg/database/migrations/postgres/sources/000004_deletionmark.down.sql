BEGIN;
alter table odahu_operator_training drop column "deletionmark";
alter table odahu_operator_packaging drop column "deletionmark";
alter table odahu_operator_deployment drop column "deletionmark";
COMMIT;