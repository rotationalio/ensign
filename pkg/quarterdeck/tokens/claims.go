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
