CREATE TABLE IF NOT EXISTS subscriptions(
    id              uuid    PRIMARY KEY,
    tenant_id       uuid    NOT NULL,
    application_id  uuid    NOT NULL,
    name            TEXT    NOT NULL,
    types           TEXT[],
    state           TEXT    NOT NULL,

    create_time timestamp   NOT NULL,
    update_time timestamp   NOT NULL
);
