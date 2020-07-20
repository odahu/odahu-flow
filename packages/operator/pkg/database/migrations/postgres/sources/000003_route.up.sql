BEGIN;

CREATE TABLE IF NOT EXISTS odahu_operator_route
(
    id   VARCHAR(64) PRIMARY KEY,
    spec JSONB,
    status JSONB
);

COMMIT;