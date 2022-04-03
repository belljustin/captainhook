CREATE TABLE IF NOT EXISTS messages(
    id              uuid    PRIMARY KEY,
    tenant_id       uuid    NOT NULL,
    application_id  uuid    NOT NULL,
    type            TEXT    NOT NULL,
    data            BYTEA   NOT NULL,
    state           TEXT    NOT NULL,

    create_time timestamp NOT NULL,
    update_time timestamp NOT NULL
);