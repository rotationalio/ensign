package tokens

import (
	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/oklog/ulid/v2"
	ulids "github.com/rotationalio/ensign/pkg/utils/ulid"
)

// Claims implements custom claims for the Quarterdeck application.
type Claims struct {
	jwt.RegisteredClaims
	Name        string   `json:"name,omitempty"`
	Email       string   `json:"email,omitempty"`
	OrgID       string   `json:"org,omitempty"`
	ProjectID   string   `json:"project,omitempty"`
	Permissions []string `json:"permissions,omitempty"`
}

// HasPermission checks if the claims contain the specified permission.
func (c Claims) HasPermission(requiredPermission string) bool {
	for _, permission := range c.Permissions {
		if permission == requiredPermission {
			return true
		}
	}
	return false
}

// HasAllPermissions checks if all specified permissions are in the claims.
func (c Claims) HasAllPermissions(requiredPermissions ...string) bool {
	for _, requiredPermission := range requiredPermissions {
		if !c.HasPermission(requiredPermission) {
			return false
		}
	}
	return true
}

// ParseOrgID returns the ULID of the organization ID in the claims. If the OrgID is not
// valid then an empty ULID is returned without an error to reduce error checking in
// handlers. If the caller needs to know if the ULID is invalid they should parse it
// themselves; otherwise the Null ULID will prevent most accesses from succeeding.
func (c Claims) ParseOrgID() ulid.ULID {
	orgID, err := ulid.Parse(c.OrgID)
	if err != nil {
		return ulids.Null
	}
	return orgID
}

// ParseUserID returns the ULID of the user from the Subject of the claims. If the
// UserID is not valid then an empty ULID is returned without an error to reduce error
// checking in the handlers. If the caller needs to know if the ULID is invalid, they
// should parse it themsleves or perform an IsZero check on the result.
func (c Claims) ParseUserID() ulid.ULID {
	userID, err := ulid.Parse(c.Subject)
	if err != nil {
		return ulids.Null
	}
	return userID
}
