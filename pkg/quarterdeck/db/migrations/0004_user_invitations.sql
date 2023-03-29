-- Add an invitations table to support user invitations.
BEGIN;

CREATE TABLE IF NOT EXISTS user_invitations (
    organization_id     BLOB NOT NULL,
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

CREATE UNIQUE INDEX unique_user_invitation_token ON user_invitations (token);

COMMIT;