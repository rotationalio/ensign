-- Add columns to the users table to support email verification.
BEGIN;

-- Whether or not the user has a verified email address.
ALTER TABLE users ADD COLUMN email_verified BOOL DEFAULT 0;

-- Expiration time for the token.
ALTER TABLE users ADD COLUMN email_verification_expires TEXT DEFAULT NULL;

-- Token provided by the user to verify their email address.
ALTER TABLE users ADD COLUMN email_verification_token TEXT DEFAULT NULL;

-- Unique index on the token to allow for fast lookups.
CREATE UNIQUE INDEX unique_email_verification_token ON users (email_verification_token);

-- Secret key used to sign the token.
ALTER TABLE users ADD COLUMN email_verification_secret BLOB DEFAULT NULL;

COMMIT;