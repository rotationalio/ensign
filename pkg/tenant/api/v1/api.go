package api

import (
	"context"
	"errors"
	"strings"
)

//===========================================================================
// Service Interface
//===========================================================================

type TenantClient interface {
	Status(context.Context) (*StatusReply, error)

	Register(context.Context, *RegisterRequest) error
	Login(context.Context, *LoginRequest) (*AuthReply, error)
	Refresh(context.Context, *RefreshRequest) (*AuthReply, error)
	Switch(context.Context, *SwitchRequest) (*AuthReply, error)
	VerifyEmail(context.Context, *VerifyRequest) error
	InvitePreview(context.Context, string) (*MemberInvitePreview, error)

	OrganizationList(context.Context, *PageQuery) (*OrganizationPage, error)
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
	MemberRoleUpdate(ctx context.Context, id string, in *UpdateRoleParams) (*Member, error)
	MemberDelete(ctx context.Context, id string) (*MemberDeleteReply, error)

	TenantProjectList(ctx context.Context, id string, in *PageQuery) (*TenantProjectPage, error)
	TenantProjectCreate(ctx context.Context, id string, in *Project) (*Project, error)
	TenantProjectPatch(ctx context.Context, tenantID, projectID string, in *Project) (*Project, error)

	TenantProjectStats(ctx context.Context, id string) ([]*StatValue, error)

	ProjectList(context.Context, *PageQuery) (*ProjectPage, error)
	ProjectCreate(context.Context, *Project) (*Project, error)
	ProjectDetail(ctx context.Context, id string) (*Project, error)
	ProjectUpdate(context.Context, *Project) (*Project, error)
	ProjectPatch(ctx context.Context, id string, in *Project) (*Project, error)
	ProjectDelete(ctx context.Context, id string) error

	ProjectTopicList(ctx context.Context, id string, in *PageQuery) (*ProjectTopicPage, error)
	ProjectTopicCreate(ctx context.Context, id string, in *Topic) (*Topic, error)

	ProjectQuery(ctx context.Context, in *ProjectQueryRequest) (*ProjectQueryResponse, error)

	TopicList(context.Context, *PageQuery) (*TopicPage, error)
	TopicDetail(ctx context.Context, id string) (*Topic, error)
	TopicEvents(ctx context.Context, id string) ([]*EventTypeInfo, error)
	TopicStats(ctx context.Context, id string) ([]*StatValue, error)
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
	Success          bool                  `json:"success"`
	Error            string                `json:"error,omitempty"`
	ValidationErrors FieldValidationErrors `json:"validation_errors,omitempty"`
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
	Email        string `json:"email"`
	Password     string `json:"password"`
	PwCheck      string `json:"pwcheck"`
	AgreeToS     bool   `json:"terms_agreement"`
	AgreePrivacy bool   `json:"privacy_agreement"`
}

// Validate ensures that all required fields are present without performing complete
// validation checks such as the password strength.
func (r *RegisterRequest) Validate() error {
	if r.Email == "" {
		return errors.New("email is required")
	}

	if r.Password == "" {
		return errors.New("password is required")
	}

	if r.Password != r.PwCheck {
		return errors.New("passwords do not match")
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
	Email       string `json:"email"`
	Password    string `json:"password"`
	OrgID       string `json:"org_id,omitempty"`
	InviteToken string `json:"invite_token,omitempty"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
	OrgID        string `json:"org_id,omitempty"`
}

type SwitchRequest struct {
	OrgID string `json:"org_id"`
}

type VerifyRequest struct {
	Token string `json:"token"`
}

type AuthReply struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	LastLogin    string `json:"last_login,omitempty"`
}

type MemberInvitePreview struct {
	Email       string `json:"email"`
	OrgName     string `json:"org_name"`
	InviterName string `json:"inviter_name"`
	Role        string `json:"role"`
	HasAccount  bool   `json:"has_account"`
}

type PageQuery struct {
	PageSize      uint32 `json:"page_size" url:"page_size,omitempty" form:"page_size"`
	NextPageToken string `json:"next_page_token" url:"next_page_token,omitempty" form:"next_page_token"`
}

type Organization struct {
	ID        string `json:"id" uri:"id"`
	Name      string `json:"name"`
	Owner     string `json:"owner"`
	Domain    string `json:"domain"`
	Projects  int    `json:"projects"`
	LastLogin string `json:"last_login"`
	Created   string `json:"created"`
	Modified  string `json:"modified"`
}

type OrganizationPage struct {
	Organizations []*Organization `json:"organizations"`
	NextPageToken string          `json:"next_page_token,omitempty"`
}

type Tenant struct {
	ID              string `json:"id" uri:"id"`
	Name            string `json:"name"`
	EnvironmentType string `json:"environment_type"`
	Created         string `json:"created,omitempty"`
	Modified        string `json:"modified,omitempty"`
}

type TenantPage struct {
	Tenants       []*Tenant `json:"tenants"`
	NextPageToken string    `json:"next_page_token,omitempty"`
}

// ID must be omitempty so that project owners can be updated on patch.
type Member struct {
	ID                string   `json:"id,omitempty" uri:"id"`
	Email             string   `json:"email"`
	Name              string   `json:"name"`
	Organization      string   `json:"organization"`
	Workspace         string   `json:"workspace"`
	ProfessionSegment string   `json:"profession_segment"`
	DeveloperSegment  []string `json:"developer_segment"`
	Picture           string   `json:"picture"`
	Role              string   `json:"role"`
	Invited           bool     `json:"invited"`
	OnboardingStatus  string   `json:"onboarding_status"`
	Created           string   `json:"created,omitempty"`
	DateAdded         string   `json:"date_added,omitempty"`
	LastActivity      string   `json:"last_activity,omitempty"`
}

// Normalize performs some cleanup on the Member fields to ensure that fields provided
// in the JSON request can be used in comparisons and uniqueness checks.
func (m *Member) Normalize() {
	m.Email = strings.TrimSpace(strings.ToLower(m.Email))
	m.Name = strings.TrimSpace(m.Name)
	m.Organization = strings.TrimSpace(m.Organization)
	m.Workspace = strings.ToLower(strings.TrimSpace(m.Workspace))
	m.ProfessionSegment = strings.TrimSpace(m.ProfessionSegment)
	for i, s := range m.DeveloperSegment {
		m.DeveloperSegment[i] = strings.TrimSpace(s)
	}
	m.Picture = strings.TrimSpace(m.Picture)
	m.Role = strings.TrimSpace(m.Role)
}

type MemberPage struct {
	Members       []*Member `json:"members"`
	NextPageToken string    `json:"next_page_token,omitempty"`
}

type UpdateRoleParams struct {
	Role string `json:"role"`
}

type MemberDeleteReply struct {
	APIKeys []string `json:"api_keys,omitempty"`
	Token   string   `json:"token,omitempty"`
	Deleted bool     `json:"deleted,omitempty"`
}

type TenantProjectPage struct {
	TenantID       string     `json:"id"`
	TenantProjects []*Project `json:"tenant_projects"`
	NextPageToken  string     `json:"next_page_token,omitempty"`
}

// Omitempty should be set on all fields to make sure the project patch endpoints only
// parse fields that were provided in the JSON request.
type Project struct {
	ID           string    `json:"id,omitempty" uri:"id"`
	TenantID     string    `json:"tenant_id,omitempty"`
	Name         string    `json:"name,omitempty"`
	Description  string    `json:"description,omitempty"`
	Owner        Member    `json:"owner,omitempty"`
	Status       string    `json:"status,omitempty"`
	ActiveTopics uint64    `json:"active_topics,omitempty"`
	DataStorage  StatValue `json:"data_storage,omitempty"`
	Created      string    `json:"created,omitempty"`
	Modified     string    `json:"modified,omitempty"`
}

type ProjectPage struct {
	Projects      []*Project `json:"projects"`
	NextPageToken string     `json:"next_page_token,omitempty"`
}

type ProjectTopicPage struct {
	ProjectID     string   `json:"project_id"`
	Topics        []*Topic `json:"topics"`
	NextPageToken string   `json:"next_page_token,omitempty"`
}

type ProjectQueryRequest struct {
	ProjectID  string            `json:"project_id"`
	Query      string            `json:"query"`
	Parameters []*QueryParameter `json:"parameters"`
}

type QueryParameter struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}

type ProjectQueryResponse struct {
	Results     []*QueryResult `json:"results"`
	TotalEvents uint64         `json:"total_events"`
	Error       string         `json:"error,omitempty"`
}

type QueryResult struct {
	Metadata        map[string]string `json:"metadata"`
	Mimetype        string            `json:"mimetype"`
	Version         string            `json:"version"`
	IsBase64Encoded bool              `json:"is_base64_encoded"`
	Data            string            `json:"data"`
	Created         string            `json:"created"`
	Error           string            `json:"error,omitempty"`
}

type Topic struct {
	ID        string `json:"id" uri:"id"`
	ProjectID string `json:"project_id"`
	Name      string `json:"topic_name"`
	Status    string `json:"status"`
	Created   string `json:"created,omitempty"`
	Modified  string `json:"modified,omitempty"`
}

type TopicPage struct {
	Topics        []*Topic `json:"topics"`
	NextPageToken string   `json:"next_page_token,omitempty"`
}

type EventTypeInfo struct {
	Type       string     `json:"type"`
	Version    string     `json:"version"`
	Mimetype   string     `json:"mimetype"`
	Events     *StatValue `json:"events"`
	Duplicates *StatValue `json:"duplicates"`
	Storage    *StatValue `json:"storage"`
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
