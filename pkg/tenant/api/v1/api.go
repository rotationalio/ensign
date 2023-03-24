package api

import (
	"context"
	"errors"
)

//===========================================================================
// Service Interface
//===========================================================================

type TenantClient interface {
	Status(context.Context) (*StatusReply, error)

	Register(context.Context, *RegisterRequest) error
	Login(context.Context, *LoginRequest) (*AuthReply, error)
	Refresh(context.Context, *RefreshRequest) (*AuthReply, error)
	VerifyEmail(context.Context, *VerifyRequest) error

	OrganizationDetail(context.Context, string) (*Organization, error)

	TenantList(context.Context, *PageQuery) (*TenantPage, error)
	TenantCreate(context.Context, *Tenant) (*Tenant, error)
	TenantDetail(ctx context.Context, id string) (*Tenant, error)
	TenantUpdate(context.Context, *Tenant) (*Tenant, error)
	TenantDelete(ctx context.Context, id string) error

	TenantStats(ctx context.Context, id string) ([]*StatValue, error)

	MemberList(context.Context, *PageQuery) (*MemberPage, error)
	MemberCreate(context.Context, *Member) (*Member, error)
	MemberDetail(ctx context.Context, id string) (*Member, error)
	MemberUpdate(context.Context, *Member) (*Member, error)
	MemberDelete(ctx context.Context, id string) error

	TenantProjectList(ctx context.Context, id string, in *PageQuery) (*TenantProjectPage, error)
	TenantProjectCreate(ctx context.Context, id string, in *Project) (*Project, error)

	ProjectList(context.Context, *PageQuery) (*ProjectPage, error)
	ProjectCreate(context.Context, *Project) (*Project, error)
	ProjectDetail(ctx context.Context, id string) (*Project, error)
	ProjectUpdate(context.Context, *Project) (*Project, error)
	ProjectDelete(ctx context.Context, id string) error

	ProjectTopicList(ctx context.Context, id string, in *PageQuery) (*ProjectTopicPage, error)
	ProjectTopicCreate(ctx context.Context, id string, in *Topic) (*Topic, error)

	TopicList(context.Context, *PageQuery) (*TopicPage, error)
	TopicDetail(ctx context.Context, id string) (*Topic, error)
	TopicUpdate(context.Context, *Topic) (*Topic, error)
	TopicDelete(ctx context.Context, in *Confirmation) (*Confirmation, error)

	ProjectAPIKeyList(ctx context.Context, id string, in *PageQuery) (*ProjectAPIKeyPage, error)
	ProjectAPIKeyCreate(ctx context.Context, id string, in *APIKey) (*APIKey, error)

	APIKeyCreate(context.Context, *APIKey) (*APIKey, error)
	APIKeyList(context.Context, *PageQuery) (*APIKeyPage, error)
	APIKeyDetail(ctx context.Context, id string) (*APIKey, error)
	APIKeyUpdate(context.Context, *APIKey) (*APIKey, error)
	APIKeyDelete(ctx context.Context, id string) error
	APIKeyPermissions(context.Context) ([]string, error)
}

//===========================================================================
// Top Level Requests and Responses
//===========================================================================

// Confirmation allows APIs to protect users from unintended actions such as deleting
// data by including a confirmation token in the request.
type Confirmation struct {
	ID     string `json:"id,omitempty"`
	Name   string `json:"name,omitempty"`
	Token  string `json:"token,omitempty"`
	Status string `json:"status,omitempty"`
}

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

//===========================================================================
// Tenant Requests and Responses
//===========================================================================

type RegisterRequest struct {
	Name         string `json:"name"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	PwCheck      string `json:"pwcheck"`
	Organization string `json:"organization"`
	Domain       string `json:"domain"`
	AgreeToS     bool   `json:"terms_agreement"`
	AgreePrivacy bool   `json:"privacy_agreement"`
}

// Validate ensures that all required fields are present without performing complete
// validation checks such as the password strength.
func (r *RegisterRequest) Validate() error {
	if r.Name == "" {
		return errors.New("name is required")
	}

	if r.Email == "" {
		return errors.New("email is required")
	}

	if r.Password == "" {
		return errors.New("password is required")
	}

	if r.Password != r.PwCheck {
		return errors.New("passwords do not match")
	}

	if r.Organization == "" {
		return errors.New("organization is required")
	}

	if r.Domain == "" {
		return errors.New("domain is required")
	}

	if !r.AgreeToS {
		return errors.New("you must agree to the terms of service")
	}

	if !r.AgreePrivacy {
		return errors.New("you must agree to the privacy policy")
	}
	return nil
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type VerifyRequest struct {
	Token string `json:"token"`
}

type AuthReply struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	LastLogin    string `json:"last_login,omitempty"`
}

type PageQuery struct {
	PageSize      uint32 `url:"page_size,omitempty"`
	NextPageToken string `url:"next_page_token,omitempty"`
}

type Organization struct {
	ID       string `json:"id" uri:"id"`
	Name     string `json:"name"`
	Owner    string `json:"owner"`
	Domain   string `json:"domain"`
	Projects int    `json:"projects"`
	Created  string `json:"created"`
	Modified string `json:"modified"`
}

type Tenant struct {
	OrgID           string `json:"org_id"`
	ID              string `json:"id" uri:"id"`
	Name            string `json:"name"`
	EnvironmentType string `json:"environment_type"`
	Created         string `json:"created,omitempty"`
	Modified        string `json:"modified,omitempty"`
}

type TenantPage struct {
	Tenants       []*Tenant `json:"tenants"`
	PrevPageToken string    `json:"prev_page_token"`
	NextPageToken string    `json:"next_page_token"`
}

type Member struct {
	ID       string `json:"id" uri:"id"`
	Name     string `json:"name"`
	Role     string `json:"role"`
	Created  string `json:"created,omitempty"`
	Modified string `json:"modified,omitempty"`
}

type MemberPage struct {
	Members       []*Member `json:"members"`
	PrevPageToken string    `json:"prev_page_token"`
	NextPageToken string    `json:"next_page_token"`
}

type TenantProjectPage struct {
	TenantID       string     `json:"id"`
	TenantProjects []*Project `json:"tenant_projects"`
	PrevPageToken  string     `json:"prev_page_token"`
	NextPageToken  string     `json:"next_page_token"`
}

type Project struct {
	ID       string `json:"id" uri:"id"`
	TenantID string `json:"tenant_id"`
	Name     string `json:"name"`
	Created  string `json:"created,omitempty"`
	Modified string `json:"modified,omitempty"`
}

type ProjectPage struct {
	Projects      []*Project `json:"projects"`
	PrevPageToken string     `json:"prev_page_token"`
	NextPageToken string     `json:"next_page_token"`
}

type ProjectTopicPage struct {
	ProjectID     string   `json:"project_id"`
	Topics        []*Topic `json:"topics"`
	PrevPageToken string   `json:"prev_page_token"`
	NextPageToken string   `json:"next_page_token"`
}

type Topic struct {
	ID        string `json:"id" uri:"id"`
	ProjectID string `json:"project_id"`
	Name      string `json:"topic_name"`
	State     string `json:"state"`
	Created   string `json:"created,omitempty"`
	Modified  string `json:"modified,omitempty"`
}

type TopicPage struct {
	Topics        []*Topic `json:"topics"`
	PrevPageToken string   `json:"prev_page_token"`
	NextPageToken string   `json:"next_page_token"`
}

type ProjectAPIKeyPage struct {
	ProjectID     string           `json:"project_id"`
	APIKeys       []*APIKeyPreview `json:"api_keys"`
	PrevPageToken string           `json:"prev_page_token"`
	NextPageToken string           `json:"next_page_token"`
}

const (
	PartialPermissions = "Partial"
	FullPermissions    = "Full"
)

type APIKeyPreview struct {
	ID          string `json:"id"`
	ClientID    string `json:"client_id"`
	Name        string `json:"name,omitempty"`
	Permissions string `json:"permissions"`
	Status      string `json:"status"`
	LastUsed    string `json:"last_used,omitempty"`
	Created     string `json:"created"`
	Modified    string `json:"modified"`
}

type APIKey struct {
	ID           string   `json:"id,omitempty"`
	ClientID     string   `json:"client_id"`
	ClientSecret string   `json:"client_secret,omitempty"`
	Name         string   `json:"name"`
	Owner        string   `json:"owner,omitempty"`
	Permissions  []string `json:"permissions,omitempty"`
	Created      string   `json:"created,omitempty"`
	Modified     string   `json:"modified,omitempty"`
}

type APIKeyPage struct {
	APIKeys       []*APIKey `json:"api_keys"`
	PrevPageToken string    `json:"prev_page_token"`
	NextPageToken string    `json:"next_page_token"`
}

// ContactInfo allows users to sign up for email notifications from SendGrid and is
// specifically used to allow users to request Ensign Private Beta access.
type ContactInfo struct {
	FirstName            string `json:"firstName"`
	LastName             string `json:"lastName"`
	Email                string `json:"email"`
	Country              string `json:"country"`
	Title                string `json:"title"`
	Organization         string `json:"organization"`
	CloudServiceProvider string `json:"cloudServiceProvider"`
}

// StatValue contains a value for a named statistic which is meant to support a variety
// of statistics endpoints.
type StatValue struct {
	Name    string  `json:"name"`
	Value   float64 `json:"value"`
	Units   string  `json:"units,omitempty"`
	Percent float64 `json:"percent,omitempty"`
}
