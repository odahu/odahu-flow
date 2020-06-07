BEGIN;

CREATE TABLE IF NOT EXISTS odahu_operator_training
(
    id   VARCHAR(64) PRIMARY KEY,
    spec JSONB,
    status JSONB
);

-- CREATE TABLE IF NOT EXISTS odahu_operator_training_event
-- (
--     id       SERIAL PRIMARY KEY,
--     text     TEXT,
--     training VARCHAR(64),
--     FOREIGN KEY (training)
--         REFERENCES odahu_operator_training (id)
-- );

CREATE TABLE IF NOT EXISTS odahu_operator_packaging
(
    id   VARCHAR(64) PRIMARY KEY,
    spec JSONB,
    status JSONB
);

-- CREATE TABLE IF NOT EXISTS odahu_operator_packaging_event
-- (
--     id        SERIAL PRIMARY KEY,
--     text      TEXT,
--     packaging VARCHAR(64),
--     FOREIGN KEY (packaging)
--         REFERENCES odahu_operator_packaging (id)
-- );


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