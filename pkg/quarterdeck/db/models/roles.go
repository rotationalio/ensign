package models

import (
	"database/sql"
)

// Role is a model that represents a row in the roles table and provides database
// functionality for interacting with role data. It should not be used for API
// serialization.
type Role struct {
	Base
	ID          int64
	Name        string
	Description sql.NullString
}

// Permission is a model that represents a row in the permissions table and provides
// database functionality for interacting with permission data. It should not be used
// for API serialization.
type Permission struct {
	Base
	ID           int64
	Name         string
	Description  sql.NullString
	AllowAPIKeys bool
	AllowRoles   bool
}

// RolePermission is a model representing a many-to-many mapping between roles and
// permissions. This model is primarily used by the Role and Permission models and is
// not intended for direct use generally.
type RolePermission struct {
	Base
	RoleID       int64
	PermissionID int64
}
