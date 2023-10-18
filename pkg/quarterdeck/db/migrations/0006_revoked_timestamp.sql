-- Add a column for tracking when an API key was revoked.

BEGIN;

ALTER TABLE revoked_api_keys ADD COLUMN revoked TEXT DEFAULT NULL;

COMMIT;