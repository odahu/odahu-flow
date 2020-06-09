BEGIN;

CREATE TABLE IF NOT EXISTS odahu_operator_training
(
    id   VARCHAR(64) PRIMARY KEY,
    spec JSONB,
    status JSONB
);


CREATE TABLE IF NOT EXISTS odahu_operator_packaging
(
    id   VARCHAR(64) PRIMARY KEY,
    spec JSONB,
    status JSONB
);


CREATE TABLE IF NOT EXISTS odahu_operator_deployment
(
    id   VARCHAR(64) PRIMARY KEY,
    spec JSONB,
    status JSONB
);


CREATE TABLE IF NOT EXISTS odahu_operator_toolchain_integration
(
    id   VARCHAR(64) PRIMARY KEY,
    spec JSONB,
    status JSONB
);


CREATE TABLE IF NOT EXISTS odahu_operator_packaging_integration
(
    id   VARCHAR(64) PRIMARY KEY,
    spec JSONB,
    status JSONB
);

COMMIT;