BEGIN;

CREATE TABLE IF NOT EXISTS odahu_outbox
(
    id  BIGSERIAL,
    entity_id VARCHAR(64),
    event_type VARCHAR(128),
    event_group VARCHAR(128),
    datetime TIMESTAMPTZ,
    payload JSONB
);

COMMIT;