-- Migrations for quarterdeck data storage.
-- These migrations target an embedded sqlite3 database that is replicated across all
-- quarterdeck nodes. The migrations table allows a booting node to determine which
-- version its schema is at so that it can quickly make changes to its data store when
-- the node starts or during runtime.
BEGIN;

-- The migrations table stores the migrations applied to arrive at the current schema
-- of the database. The quarterdeck application checks this table for the version the
-- db is at and applies any later migrations as needed.
CREATE TABLE IF NOT EXISTS migrations (
    id      INTEGER PRIMARY KEY,
    name    TEXT NOT NULL,
    version TEXT NOT NULL,
    created TEXT NOT NULL
);

COMMIT;