-- Initial schema for the quarterdeck application
BEGIN;

--
-- Table Definitions
-- Primary keys are expected to be 16-byte UUID or ULID data structures
--

CREATE TABLE IF NOT EXISTS organizations (
    id                  BLOB PRIMARY KEY,
    name                TEXT NOT NULL,
    domain              TEXT NOT NULL UNIQUE,
    created             TEXT NOT NULL,
    modified            TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS users (
    id                  BLOB PRIMARY KEY,
    name                TEXT NOT NULL,
    email               TEXT NOT NULL UNIQUE,
    password            TEXT NOT NULL UNIQUE,
    last_login          TEXT,
    created             TEXT NOT NULL,
    modified            TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS organization_users (
    organization_id     BLOB NOT NULL,
    user_id             BLOB NOT NULL,
    created             TEXT NOT NULL,
    modified            TEXT NOT NULL,
    PRIMARY KEY (organization_id, user_id),
    FOREIGN KEY (organization_id) REFERENCES organizations (id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS api_keys (
    id                  BLOB PRIMARY KEY,
    key_id              TEXT NOT NULL UNIQUE,
    secret              TEXT NOT NULL UNIQUE,
    name                TEXT NOT NULL,
    project_id          TEXT NOT NULL,
    created_by          BLOB,
    created             TEXT NOT NULL,
    modified            TEXT NOT NULL,
    FOREIGN KEY (created_by) REFERENCES users (id) ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS roles (
    id                  INTEGER PRIMARY KEY,
    name                TEXT NOT NULL UNIQUE,
    description         TEXT,
    created             TEXT NOT NULL,
    modified            TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS user_roles (
    user_id             BLOB NOT NULL,
    role_id             INTEGER NOT NULL,
    created             TEXT NOT NULL,
    modified            TEXT NOT NULL,
    PRIMARY KEY (user_id, role_id),
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (role_id) REFERENCES roles (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS permissions (
    id                  INTEGER PRIMARY KEY,
    name                TEXT NOT NULL UNIQUE,
    description         TEXT,
    allow_api_keys      BOOL DEFAULT false,
    allow_roles         BOOL DEFAULT true,
    created             TEXT NOT NULL,
    modified            TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS role_permissions (
    role_id             INTEGER NOT NULL,
    permission_id       INTEGER NOT NULL,
    created             TEXT NOT NULL,
    modified            TEXT NOT NULL,
    PRIMARY KEY (role_id, permission_id),
    FOREIGN KEY (role_id) REFERENCES roles (id) ON DELETE CASCADE,
    FOREIGN KEY (permission_id) REFERENCES permissions (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS api_key_permissions (
    api_key_id          BLOB NOT NULL,
    permission_id       INTEGER NOT NULL,
    created             TEXT NOT NULL,
    modified            TEXT NOT NULL,
    PRIMARY KEY (api_key_id, permission_id),
    FOREIGN KEY (api_key_id) REFERENCES api_keys (id) ON DELETE CASCADE,
    FOREIGN KEY (permission_id) REFERENCES permissions (id) ON DELETE CASCADE
);

COMMIT;