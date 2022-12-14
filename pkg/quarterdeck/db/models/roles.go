package models

import (
	"database/sql"

	"github.com/oklog/ulid/v2"
)

// Role is a model that represents a row in the roles table and provides database
// functionality for interacting with role data. It should not be used for API
// serialization.
type Role struct {
	Base
	ID          ulid.ULID
	Name        string
	Description sql.NullString
}

// Permission is a model that represents a row in the permissions table and provides
// database functionality for interacting with permission data. It should not be used
// for API serialization.
type Permission struct {
	Base
	ID           ulid.ULID
	Name         string
	Description  sql.NullString
	AllowAPIKeys bool
	AllowRoles   bool
}

// RollPermission is a model representing a many-to-many mapping between roles and
// permissions. This model is primarily used by the Role and Permission models and is
// not intended for direct use generally.
type RollPermission struct {
	Base
	RoleID       ulid.ULID
	PermissionID ulid.ULID
}
