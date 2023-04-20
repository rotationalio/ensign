-- Add a column to the organization_users to record last login times.

BEGIN;

-- Last login time for the user in the organization.
ALTER TABLE organization_users ADD COLUMN last_login TEXT DEFAULT NULL;

COMMIT;