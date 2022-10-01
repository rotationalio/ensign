-- Populate the database with initial data for roles and permissions.
BEGIN;

INSERT INTO roles (id, name, description, created, modified) VALUES
    (1, 'Owner', 'Manages the organization and billing accounts', datetime('now'), datetime('now')),
    (2, 'Admin', 'Full management access to the organization', datetime('now'), datetime('now')),
    (3, 'Member', 'Can view and manage projects, topics, and keys', datetime('now'), datetime('now')),
    (4, 'Observer', 'Can view projects and topics but cannot manage or edit them', datetime('now'), datetime('now'))
;

INSERT INTO permissions (id, name, description, allow_api_keys, allow_roles, created, modified) VALUES
    (1, 'organization:edit', 'Can make changes to the details of an organization', false, true, datetime('now'), datetime('now')),
    (2, 'collaborators:add', 'Can add users to an organization and change their roles', false, true, datetime('now'), datetime('now')),
    (3, 'collaborators:remove', 'Can remove users from an organization', false, true, datetime('now'), datetime('now')),
    (4, 'collaborators:read', 'Can see the collaborators for an organization', false, true, datetime('now'), datetime('now')),
    (5, 'projects:edit', 'Can create and update projects in the organization', false, true, datetime('now'), datetime('now')),
    (6, 'projects:delete', 'Can delete projects from the organization along with their data', false, true, datetime('now'), datetime('now')),
    (7, 'projects:read', 'Can view projects in the organization', false, true, datetime('now'), datetime('now')),
    (8, 'apikeys:edit', 'Can create and update api keys in a project', false, true, datetime('now'), datetime('now')),
    (9, 'apikeys:delete', 'Can delete api keys from a project', false, true, datetime('now'), datetime('now')),
    (10, 'apikeys:read', 'Can view api keys inside of a project', false, true, datetime('now'), datetime('now')),
    (11, 'topics:create', 'Can create topics in a project', true, true, datetime('now'), datetime('now')),
    (12, 'topics:edit', 'Can modify topic metadata and mark as read only', true, true, datetime('now'), datetime('now')),
    (13, 'topics:destroy', 'Can destroy topics and all their data', true, true, datetime('now'), datetime('now')),
    (14, 'topics:read', 'Can view the topics for a project', true, true, datetime('now'), datetime('now')),
    (15, 'metrics:read', 'Can view the metrics of topics in the specified project', true, true, datetime('now'), datetime('now')),
    (16, 'publisher', 'Is allowed to publish events to the topics in the project', true, false, datetime('now'), datetime('now')),
    (17, 'subscriber', 'Is allowed to subscribe to topics in the project to consume events', true, false, datetime('now'), datetime('now'))
;

INSERT INTO role_permissions (role_id, permission_id, created, modified) VALUES
    (1, 1, datetime('now'), datetime('now')),
    (1, 2, datetime('now'), datetime('now')),
    (1, 3, datetime('now'), datetime('now')),
    (1, 4, datetime('now'), datetime('now')),
    (1, 5, datetime('now'), datetime('now')),
    (1, 6, datetime('now'), datetime('now')),
    (1, 7, datetime('now'), datetime('now')),
    (1, 8, datetime('now'), datetime('now')),
    (1, 9, datetime('now'), datetime('now')),
    (1, 10, datetime('now'), datetime('now')),
    (1, 11, datetime('now'), datetime('now')),
    (1, 12, datetime('now'), datetime('now')),
    (1, 13, datetime('now'), datetime('now')),
    (1, 14, datetime('now'), datetime('now')),
    (1, 15, datetime('now'), datetime('now')),
    (2, 2, datetime('now'), datetime('now')),
    (2, 3, datetime('now'), datetime('now')),
    (2, 4, datetime('now'), datetime('now')),
    (2, 5, datetime('now'), datetime('now')),
    (2, 6, datetime('now'), datetime('now')),
    (2, 7, datetime('now'), datetime('now')),
    (2, 8, datetime('now'), datetime('now')),
    (2, 9, datetime('now'), datetime('now')),
    (2, 10, datetime('now'), datetime('now')),
    (2, 11, datetime('now'), datetime('now')),
    (2, 12, datetime('now'), datetime('now')),
    (2, 13, datetime('now'), datetime('now')),
    (2, 14, datetime('now'), datetime('now')),
    (2, 15, datetime('now'), datetime('now')),
    (3, 4, datetime('now'), datetime('now')),
    (3, 5, datetime('now'), datetime('now')),
    (3, 6, datetime('now'), datetime('now')),
    (3, 7, datetime('now'), datetime('now')),
    (3, 8, datetime('now'), datetime('now')),
    (3, 9, datetime('now'), datetime('now')),
    (3, 10, datetime('now'), datetime('now')),
    (3, 11, datetime('now'), datetime('now')),
    (3, 12, datetime('now'), datetime('now')),
    (3, 14, datetime('now'), datetime('now')),
    (3, 15, datetime('now'), datetime('now')),
    (4, 4, datetime('now'), datetime('now')),
    (4, 7, datetime('now'), datetime('now')),
    (4, 14, datetime('now'), datetime('now')),
    (4, 15, datetime('now'), datetime('now'))
;

COMMIT;