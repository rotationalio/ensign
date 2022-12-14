package models

import (
	"database/sql"

	"github.com/oklog/ulid/v2"
)

// APIKey is a model that represents a row in the api_keys table and provides database
// functionality for interacting with api key data. It should not be used for API
// serialization.
type APIKey struct {
	Base
	ID        ulid.ULID
	KeyID     string
	Secret    string
	Name      string
	ProjectID string
	CreatedBy sql.NullByte
}

// APIKeyPermission is a model representing a many-to-many mapping between api keys and
// permissions. This model is primarily used by the APIKey and Permission models and is
// not intended for direct use generally.
type APIKeyPermission struct {
	Base
	RoleID       ulid.ULID
	PermissionID int64
}
