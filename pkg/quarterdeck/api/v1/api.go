package api

import (
	"context"
	"strings"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg/quarterdeck/passwd"
	ulids "github.com/rotationalio/ensign/pkg/utils/ulid"
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
	APIKeyList(context.Context, *APIPageQuery) (*APIKeyList, error)
	APIKeyCreate(context.Context, *APIKey) (*APIKey, error)
	APIKeyDetail(context.Context, string) (*APIKey, error)
	APIKeyUpdate(context.Context, *APIKey) (*APIKey, error)
	APIKeyDelete(context.Context, string) error

	// Project Resource
	ProjectCreate(context.Context, *Project) (*Project, error)
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
	ID      string `json:"user_id"`
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
	ID           ulid.ULID `json:"id,omitempty"`            // not allowed on create
	ClientID     string    `json:"client_id"`               // not allowed on create, cannot be updated
	ClientSecret string    `json:"client_secret,omitempty"` // not allowed on created, cannot be updated
	Name         string    `json:"name"`                    // required on create, update
	OrgID        ulid.ULID `json:"org_id"`                  // required on create, cannot be updated
	ProjectID    ulid.ULID `json:"project_id"`              // required on create, cannot be updated
	CreatedBy    ulid.ULID `json:"created_by,omitempty"`    // required on create, cannot be updated
	Source       string    `json:"source,omitempty"`        // not required, but useful
	UserAgent    string    `json:"user_agent,omitempty"`    // not required, but useful
	LastUsed     time.Time `json:"last_used,omitempty"`     // cannot be edited
	Permissions  []string  `json:"permissions,omitempty"`   // required on create, cannot be updated
	Created      time.Time `json:"created,omitempty"`       // cannot be edited
	Modified     time.Time `json:"modified,omitempty"`      // cannot be edited
}

type APIKeyList struct {
	APIKeys       []*APIKey `json:"apikeys"`
	NextPageToken string    `json:"next_page_token,omitempty"`
}

type APIPageQuery struct {
	ProjectID     string `json:"project_id,omitempty" url:"project_id,omitempty" form:"project_id"`
	PageSize      int    `json:"page_size" url:"page_size,omitempty" form:"page_size"`
	NextPageToken string `json:"next_page_token" url:"next_page_token,omitempty" form:"next_page_token"`
}

// ValidateCreate ensures that the APIKey is valid when sent to the Create REST method.
// Validation ensures that the user does not supply data not allowed on create and that
// required fields are present.
func (k *APIKey) ValidateCreate() error {
	switch {
	case !ulids.IsZero(k.ID):
		return RestrictedField("id")
	case k.ClientID != "":
		return RestrictedField("client_id")
	case k.ClientSecret != "":
		return RestrictedField("client_secret")
	case k.Name == "":
		return MissingField("name")
	case !ulids.IsZero(k.OrgID):
		return RestrictedField("org_id")
	case ulids.IsZero(k.ProjectID):
		return MissingField("project_id")
	case !ulids.IsZero(k.CreatedBy):
		return RestrictedField("created_by")
	case k.UserAgent != "":
		return RestrictedField("user_agent")
	case !k.LastUsed.IsZero():
		return RestrictedField("last_used")
	case len(k.Permissions) == 0:
		return MissingField("permissions")
	default:
		return nil
	}
}

// ValidateUpdate ensures that the APIKey is valid when sent to the Update REST method.
// Validation ensures that the user does not supply data not allowed on updated and that
// any required fields are present to update the model.
func (k *APIKey) ValidateUpdate() error {
	switch {
	case ulids.IsZero(k.ID):
		return MissingField("id")
	case k.ClientID != "":
		return RestrictedField("client_id")
	case k.ClientSecret != "":
		return RestrictedField("client_secret")
	case k.Name == "":
		return MissingField("name")
	case !ulids.IsZero(k.OrgID):
		return RestrictedField("org_id")
	case !ulids.IsZero(k.ProjectID):
		return RestrictedField("project_id")
	case !ulids.IsZero(k.CreatedBy):
		return RestrictedField("created_by")
	case k.Source != "":
		return RestrictedField("source")
	case k.UserAgent != "":
		return RestrictedField("user_agent")
	case !k.LastUsed.IsZero():
		return RestrictedField("last_used")
	case len(k.Permissions) != 0:
		return RestrictedField("permissions")
	default:
		return nil
	}
}

//===========================================================================
// Project Resource
//===========================================================================

type Project struct {
	OrgID     ulid.ULID `json:"org_id,omitempty"`   // not allowed on create
	ProjectID ulid.ULID `json:"project_id"`         // required on create
	Created   time.Time `json:"created,omitempty"`  // cannot be edited
	Modified  time.Time `json:"modified,omitempty"` // cannot be edited
}

func (p *Project) Validate() error {
	switch {
	case !ulids.IsZero(p.OrgID):
		return RestrictedField("org_id")
	case ulids.IsZero(p.ProjectID):
		return MissingField("project_id")
	default:
		return nil
	}
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
