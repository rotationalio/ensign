package tokens

import jwt "github.com/golang-jwt/jwt/v4"

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
