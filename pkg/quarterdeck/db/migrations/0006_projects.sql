-- Add columns to the organization_projects table.

BEGIN;

-- Owner of the project.
ALTER TABLE organization_projects ADD COLUMN owner_id BLOB NOT NULL DEFAULT '';

UPDATE organization_projects
SET owner_id = (
    SELECT user_id FROM organization_users
    WHERE organization_id = organization_projects.organization_id AND role_id = 1
    ORDER BY created LIMIT 1)
WHERE owner_id = '';

COMMIT;