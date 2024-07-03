-- Add a column for tracking the user account type.

BEGIN;

ALTER TABLE users ADD COLUMN account_type TEXT DEFAULT 'sandbox';

COMMIT;