package api

// Credentials provides a basic interface for loading an access token from Quarterdeck
// into the Quarterdeck API client. Credentials can be loaded from disk, generated, or
// feched from a passthrough request.
type Credentials interface {
	AccessToken() (string, error)
}

// A Token is just the JWT base64 encoded token string that is obtained from
// Quarterdeck either using the authtest server or from a login with the client.
type Token string

// Token implements the credentials interface and performs limited validation.
func (t Token) AccessToken() (string, error) {
	if string(t) == "" {
		return "", ErrInvalidCredentials
	}
	return string(t), nil
}
