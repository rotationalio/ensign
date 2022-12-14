package api

import (
	"context"
	"strings"

	"github.com/rotationalio/ensign/pkg/quarterdeck/passwd"
)

//===========================================================================
// Service Interface
//===========================================================================

type QuarterdeckClient interface {
	// Unauthenticated endpoints
	Status(context.Context) (*StatusReply, error)
	Register(context.Context, *RegisterRequest) (*RegisterReply, error)
	Login(context.Context, *LoginRequest) (*LoginReply, error)
	Authenticate(context.Context, *APIAuthentication) (*LoginReply, error)

	// Authenticated endpoints
	Refresh(context.Context) (*LoginReply, error)

	// API Keys Resource
	APIKeyList(context.Context, *PageQuery) (*APIKeyList, error)
	APIKeyCreate(context.Context, *APIKey) (*APIKey, error)
	APIKeyDetail(context.Context, string) (*APIKey, error)
	APIKeyUpdate(context.Context, *APIKey) (*APIKey, error)
	APIKeyDelete(context.Context, string) error
}

//===========================================================================
// Top Level Requests and Responses
//===========================================================================

// Reply contains standard fields that are used for generic API responses and errors.
type Reply struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

// Returned on status requests.
type StatusReply struct {
	Status  string `json:"status"`
	Uptime  string `json:"uptime,omitempty"`
	Version string `json:"version,omitempty"`
}

// PageQuery manages paginated list requests.
type PageQuery struct {
	PageSize      int    `json:"page_size" url:"page_size,omitempty" form:"page_size"`
	NextPageToken string `json:"next_page_token" url:"next_page_token,omitempty" form:"next_page_token"`
}

//===========================================================================
// Quarterdeck API Requests and Replies
//===========================================================================

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	PwCheck  string `json:"pwcheck"`
}

// Validate the register request ensuring that the required fields are available and
// that the password is valid - an error is returned if the request is not correct. This
// method also performs some basic data cleanup, trimming whitespace.
func (r *RegisterRequest) Validate() error {
	r.Name = strings.TrimSpace(r.Name)
	r.Email = strings.TrimSpace(r.Email)
	r.Password = strings.TrimSpace(r.Password)
	r.PwCheck = strings.TrimSpace(r.PwCheck)

	if r.Name == "" || r.Email == "" {
		return ErrMissingRegisterField
	}

	if r.Password != r.PwCheck {
		return ErrPasswordMismatch
	}

	if passwd.Strength(r.Password) < passwd.Moderate {
		return ErrPasswordTooWeak
	}
	return nil
}

type RegisterReply struct {
	ID      int    `json:"user_id"`
	Email   string `json:"email"`
	Message string `json:"message"`
	Role    string `json:"role"`
	Created string `json:"created"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginReply struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type APIAuthentication struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

//===========================================================================
// API Key Resource
//===========================================================================

type APIKey struct {
	ID           int      `json:"id,omitempty"`
	ClientID     string   `json:"client_id"`
	ClientSecret string   `json:"client_secret,omitempty"`
	Name         string   `json:"name"`
	ProjectID    string   `json:"project_id"`
	Owner        string   `json:"owner,omitempty"`
	Permissions  []string `json:"permissions,omitempty"`
	Created      string   `json:"created,omitempty"`
	Modified     string   `json:"modified,omitempty"`
}

type APIKeyList struct {
	APIKeys       []*APIKey `json:"apikeys"`
	NextPageToken string    `json:"next_page_token,omitempty"`
}

//===========================================================================
// OpenID Configuration
//===========================================================================

type OpenIDConfiguration struct {
	Issuer                        string   `json:"issuer"`
	AuthorizationEP               string   `json:"authorization_endpoint"`
	TokenEP                       string   `json:"token_endpoint"`
	DeviceAuthorizationEP         string   `json:"device_authorization_endpoint"`
	UserInfoEP                    string   `json:"userinfo_endpoint"`
	MFAChallengeEP                string   `json:"mfa_challenge_endpoint"`
	JWKSURI                       string   `json:"jwks_uri"`
	RegistrationEP                string   `json:"registration_endpoint"`
	RevocationEP                  string   `json:"revocation_endpoint"`
	ScopesSupported               []string `json:"scopes_supported"`
	ResponseTypesSupported        []string `json:"response_types_supported"`
	CodeChallengeMethodsSupported []string `json:"code_challenge_methods_supported"`
	ResponseModesSupported        []string `json:"response_modes_supported"`
	SubjectTypesSupported         []string `json:"subject_types_supported"`
	IDTokenSigningAlgValues       []string `json:"id_token_signing_alg_values_supported"`
	TokenEndpointAuthMethods      []string `json:"token_endpoint_auth_methods_supported"`
	ClaimsSupported               []string `json:"claims_supported"`
	RequestURIParameterSupported  bool     `json:"request_uri_parameter_supported"`
}
