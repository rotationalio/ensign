package tokens

import jwt "github.com/golang-jwt/jwt/v4"

// Claims implements custom claims for the Quarterdeck application.
type Claims struct {
	jwt.RegisteredClaims
	Domain  string `json:"hd,omitempty"`
	Email   string `json:"email,omitempty"`
	Name    string `json:"name,omitempty"`
	Picture string `json:"picture,omitempty"`
}
