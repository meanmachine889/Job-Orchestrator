CREATE TABLE workers (
    id UUID PRIMARY KEY,
    hostname TEXT NOT NULL,
    last_heartbeat TIMESTAMPTZ NOT NULL,
    status TEXT NOT NULL
);