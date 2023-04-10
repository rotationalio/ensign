-- Add an invitations table to support user invitations.
BEGIN;

CREATE TABLE IF NOT EXISTS user_invitations (
    user_id             BLOB NOT NULL,
    organization_id     BLOB NOT NULL,
    role                TEXT NOT NULL,
    email               TEXT NOT NULL,
    expires             TEXT NOT NULL,
    token               TEXT NOT NULL UNIQUE,
    secret              BLOB NOT NULL,
    created_by          BLOB NOT NULL,
    created             TEXT NOT NULL,
    modified            TEXT NOT NULL,
    FOREIGN KEY (organization_id) REFERENCES organizations (id) ON DELETE CASCADE,
    FOREIGN KEY (created_by) REFERENCES users (id) ON DELETE CASCADE
);

COMMIT;