BEGIN;
ALTER TABLE IF EXISTS odahu_operator_training_integration RENAME TO odahu_operator_toolchain_integration;
COMMIT;