CREATE TABLE event_store (
    event_id varchar(255) NOT NULL,
    aggregate_id varchar(255) NOT NULL,
    aggregate_type varchar(255) NOT NULL,
    revision integer NOT NULL,
    timestamp bigint NOT NULL,
    event_type varchar(255) NOT NULL,
    event_data blob,
    CONSTRAINT aggregate_revision UNIQUE (aggregate_id, revision)
);

ALTER TABLE event_store ADD PRIMARY KEY (event_id);

CREATE TABLE snapshot (
    id bigint NOT NULL,
    aggregate_id varchar(255) NOT NULL,
    aggregate_type varchar(255) NOT NULL,
    aggregate_version integer NOT NULL,
    revision integer NOT NULL,
    timestamp bigint NOT NULL,
    snapshot_data blob
);

ALTER TABLE snapshot ADD PRIMARY KEY (id);
ALTER TABLE snapshot MODIFY COLUMN id bigint NOT NULL AUTO_INCREMENT;