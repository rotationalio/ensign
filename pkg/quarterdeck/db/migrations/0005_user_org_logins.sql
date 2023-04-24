-- Add columns to the organization_users table.

BEGIN;

-- Last login time for the user in the organization.
ALTER TABLE organization_users ADD COLUMN last_login TEXT DEFAULT NULL;

-- Delete confirmation token for the organization user.
ALTER TABLE organization_users ADD COLUMN delete_confirmation_token TEXT DEFAULT NULL;

COMMIT;