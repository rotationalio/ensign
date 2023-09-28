package api

import (
	"context"
	"strings"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg/quarterdeck/passwd"
	"github.com/rotationalio/ensign/pkg/utils/ulids"
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
	Refresh(context.Context, *RefreshRequest) (*LoginReply, error)
	Switch(context.Context, *SwitchRequest) (*LoginReply, error)
	VerifyEmail(context.Context, *VerifyRequest) (*LoginReply, error)
	ResendEmail(context.Context, *ResendRequest) error
	ForgotPassword(context.Context, *ForgotPasswordRequest) error

	// Organizations Resource
	OrganizationDetail(context.Context, string) (*Organization, error)
	OrganizationUpdate(context.Context, *Organization) (*Organization, error)
	OrganizationList(context.Context, *OrganizationPageQuery) (*OrganizationList, error)
	WorkspaceLookup(context.Context, *WorkspaceQuery) (*Workspace, error)

	// API Keys Resource
	APIKeyList(context.Context, *APIPageQuery) (*APIKeyList, error)
	APIKeyCreate(context.Context, *APIKey) (*APIKey, error)
	APIKeyDetail(context.Context, string) (*APIKey, error)
	APIKeyUpdate(context.Context, *APIKey) (*APIKey, error)
	APIKeyDelete(context.Context, string) error
	APIKeyPermissions(context.Context) ([]string, error)

	// Project Resource
	ProjectList(context.Context, *PageQuery) (*ProjectList, error)
	ProjectCreate(context.Context, *Project) (*Project, error)
	ProjectAccess(context.Context, *Project) (*LoginReply, error)
	ProjectDetail(context.Context, string) (*Project, error)

	// Users Resource
	UserUpdate(context.Context, *User) (*User, error)
	UserRoleUpdate(context.Context, *UpdateRoleRequest) (*User, error)
	UserList(context.Context, *UserPageQuery) (*UserList, error)
	UserDetail(context.Context, string) (*User, error)
	UserRemove(context.Context, string) (*UserRemoveReply, error)
	UserRemoveConfirm(context.Context, *UserRemoveConfirm) error

	// Invites Resource
	InvitePreview(context.Context, string) (*UserInvitePreview, error)
	InviteCreate(context.Context, *UserInviteRequest) (*UserInviteReply, error)
	InviteAccept(context.Context, *UserInviteToken) (*LoginReply, error)

	// Accounts Resource
	AccountUpdate(context.Context, *User) (*User, error)

	// Client Utility Functions
	WaitForReady(context.Context) error
}

//===========================================================================
// Top Level Requests and Responses
//===========================================================================

// Reply contains standard fields that are used for generic API responses and errors.
type Reply struct {
	Success    bool   `json:"success"`
	Error      string `json:"error,omitempty"`
	Unverified bool   `json:"unverified,omitempty"`
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
	ProjectID    string `json:"project_id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	PwCheck      string `json:"pwcheck"`
	Organization string `json:"organization"`
	Domain       string `json:"domain"`
	AgreeToS     bool   `json:"terms_agreement"`
	AgreePrivacy bool   `json:"privacy_agreement"`
}

// Validate the register request ensuring that the required fields are available and
// that the password is valid - an error is returned if the request is not correct. This
// method also performs some basic data cleanup, trimming whitespace.
func (r *RegisterRequest) Validate() error {
	r.Name = strings.TrimSpace(r.Name)
	r.Email = strings.TrimSpace(r.Email)
	r.Password = strings.TrimSpace(r.Password)
	r.PwCheck = strings.TrimSpace(r.PwCheck)
	r.Organization = strings.TrimSpace(r.Organization)
	r.Domain = strings.ToLower(strings.TrimSpace(r.Domain))

	// Required for all requests
	switch {
	case r.Email == "":
		return MissingField("email")
	case r.Password == "":
		return MissingField("password")
	case r.Password != r.PwCheck:
		return ErrPasswordMismatch
	case passwd.Strength(r.Password) < passwd.Moderate:
		return ErrPasswordTooWeak
	case !r.AgreeToS:
		return MissingField("terms_agreement")
	case !r.AgreePrivacy:
		return MissingField("privacy_agreement")
	}

	return nil
}

type RegisterReply struct {
	ID        ulid.ULID `json:"user_id"`
	OrgID     ulid.ULID `json:"org_id"`
	Email     string    `json:"email"`
	OrgName   string    `json:"org_name"`
	OrgDomain string    `json:"org_domain"`
	Message   string    `json:"message"`
	Role      string    `json:"role"`
	Created   string    `json:"created"`
}

type LoginRequest struct {
	Email       string    `json:"email"`
	Password    string    `json:"password"`
	OrgID       ulid.ULID `json:"org_id,omitempty"`
	InviteToken string    `json:"invite_token,omitempty"`
}

type LoginReply struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
	LastLogin    string `json:"last_login,omitempty"`
}

type APIAuthentication struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

type RefreshRequest struct {
	RefreshToken string    `json:"refresh_token"`
	OrgID        ulid.ULID `json:"org_id,omitempty"`
}

type SwitchRequest struct {
	OrgID ulid.ULID `json:"org_id"`
}

type VerifyRequest struct {
	Token string    `json:"token"`
	OrgID ulid.ULID `json:"org_id,omitempty"`
}

type ResendRequest struct {
	Email string    `json:"email"`
	OrgID ulid.ULID `json:"org_id,omitempty"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

//===========================================================================
// Organization Resource
//===========================================================================

type Organization struct {
	ID        ulid.ULID `json:"id"`
	Name      string    `json:"name"`
	Domain    string    `json:"domain"`
	Projects  int       `json:"projects"`
	LastLogin time.Time `json:"last_login,omitempty"`
	Created   time.Time `json:"created,omitempty"`
	Modified  time.Time `json:"modified,omitempty"`
}

func (o *Organization) ValidateUpdate() error {
	switch {
	case ulids.IsZero(o.ID):
		return MissingField("id")
	case o.Name == "":
		return MissingField("name")
	case o.Domain == "":
		return MissingField("domain")
	default:
		return nil
	}
}

type OrganizationList struct {
	Organizations []*Organization `json:"organizations"`
	NextPageToken string          `json:"next_page_token,omitempty"`
}

type OrganizationPageQuery struct {
	PageSize      int    `json:"page_size" url:"page_size,omitempty" form:"page_size"`
	NextPageToken string `json:"next_page_token" url:"next_page_token,omitempty" form:"next_page_token"`
}

type WorkspaceQuery struct {
	Domain         string `json:"domain" url:"domain,omitempty" form:"domain"`
	CheckAvailable bool   `json:"check_available" url:"check_available,omitempty" form:"check_available"`
}

type Workspace struct {
	OrgID       ulid.ULID `json:"org_id,omitempty"`
	Name        string    `json:"name,omitempty"`
	Domain      string    `json:"domain"`
	IsAvailable bool      `json:"is_available"`
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

type APIKeyPreview struct {
	ID        ulid.ULID `json:"id"`
	ClientID  string    `json:"client_id"`
	Name      string    `json:"name,omitempty"`
	ProjectID ulid.ULID `json:"project_id"`
	Partial   bool      `json:"partial"`
	Status    string    `json:"status"`
	LastUsed  time.Time `json:"last_used,omitempty"`
	Created   time.Time `json:"created"`
	Modified  time.Time `json:"modified"`
}

type APIKeyList struct {
	APIKeys       []*APIKeyPreview `json:"apikeys"`
	NextPageToken string           `json:"next_page_token,omitempty"`
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

type ProjectList struct {
	Projects      []*Project `json:"projects"`
	NextPageToken string     `json:"next_page_token,omitempty"`
}

type Project struct {
	OrgID        ulid.ULID `json:"org_id,omitempty"`        // not allowed on create
	ProjectID    ulid.ULID `json:"project_id"`              // required on create and access
	APIKeysCount int       `json:"apikeys_count,omitempty"` // cannot be edited
	RevokedCount int       `json:"revoked_count,omitempty"` // cannot be edited
	Created      time.Time `json:"created,omitempty"`       // cannot be edited
	Modified     time.Time `json:"modified,omitempty"`      // cannot be edited
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

// ===========================================================================
// Users Resource
// ===========================================================================

type User struct {
	UserID    ulid.ULID `json:"user_id"`
	OrgID     ulid.ULID `json:"org_id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	LastLogin time.Time `json:"last_login"`
}

type UpdateRoleRequest struct {
	ID   ulid.ULID `json:"id"`
	Role string    `json:"role"`
}

type UserList struct {
	Users         []*User `json:"users"`
	NextPageToken string  `json:"next_page_token,omitempty"`
}

type UserPageQuery struct {
	PageSize      int    `json:"page_size" url:"page_size,omitempty" form:"page_size"`
	NextPageToken string `json:"next_page_token" url:"next_page_token,omitempty" form:"next_page_token"`
}

type UserRemoveConfirm struct {
	ID    ulid.ULID `json:"id"`
	Token string    `json:"token"`
}

type UserRemoveReply struct {
	APIKeys []string `json:"api_keys,omitempty"`
	Token   string   `json:"token,omitempty"`
	Deleted bool     `json:"deleted"`
}

// TODO: validate Email
func (u *User) ValidateUpdate() error {
	switch {
	case ulids.IsZero(u.UserID):
		return MissingField("user_id")
	case u.Name == "":
		return MissingField("name")
	default:
		return nil
	}
}

// ===========================================================================
// Invites Resource
// ===========================================================================

// UserInviteToken contains a token that is used to accept an invite.
type UserInviteToken struct {
	Token string `json:"token"`
}

// UserInvitePreview contains user-facing information about an invite but not any
// internal details such as IDs.
type UserInvitePreview struct {
	Email       string `json:"email"`
	OrgName     string `json:"org_name"`
	InviterName string `json:"inviter_name"`
	Role        string `json:"role"`
	UserExists  bool   `json:"user_exists"`
}

// NOTE: Users can only invite someone to the organization they are currently logged
// into.
type UserInviteRequest struct {
	Email string `json:"email"`
	Role  string `json:"role"`
}

// UserInviteReply contains detailed information that corresponds to a newly issued
// invite token.
type UserInviteReply struct {
	UserID       ulid.ULID `json:"user_id"`
	OrgID        ulid.ULID `json:"org_id"`
	Email        string    `json:"email"`
	Role         string    `json:"role"`
	Name         string    `json:"name"`
	Organization string    `json:"organization"`
	Workspace    string    `json:"workspace"`
	ExpiresAt    string    `json:"expires_at"`
	CreatedBy    ulid.ULID `json:"created_by"`
	Created      string    `json:"created"`
}
