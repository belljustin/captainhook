CREATE TABLE IF NOT EXISTS applications(
    id          uuid    PRIMARY KEY,
    tenant_id   uuid    NOT NULL,
    name        TEXT    NOT NULL,

    create_time timestamp   NOT NULL,
    update_time timestamp   NOT NULL
);

CREATE UNIQUE INDEX applications_tenant_id_name_unique
    ON applications(tenant_id, name);