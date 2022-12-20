package api

import (
	"context"
)

//===========================================================================
// Service Interface
//===========================================================================

type TenantClient interface {
	Status(context.Context) (*StatusReply, error)
	SignUp(context.Context, *ContactInfo) error

	TenantList(context.Context, *PageQuery) (*TenantPage, error)
	TenantCreate(context.Context, *Tenant) (*Tenant, error)
	TenantDetail(ctx context.Context, id string) (*Tenant, error)
	TenantUpdate(context.Context, *Tenant) (*Tenant, error)
	TenantDelete(ctx context.Context, id string) error

	TenantMemberList(ctx context.Context, id string, in *PageQuery) (*TenantMemberPage, error)
	TenantMemberCreate(ctx context.Context, id string, in *Member) (*Member, error)

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
	TopicCreate(context.Context, *Topic) (*Topic, error)
	TopicDetail(ctx context.Context, id string) (*Topic, error)
	TopicUpdate(context.Context, *Topic) (*Topic, error)
	TopicDelete(ctx context.Context, id string) error

	ProjectAPIKeyList(ctx context.Context, id string, in *PageQuery) (*ProjectAPIKeyPage, error)
	ProjectAPIKeyCreate(ctx context.Context, id string, in *APIKey) (*APIKey, error)

	APIKeyList(context.Context, *PageQuery) (*APIKeyPage, error)
	APIKeyCreate(context.Context, *APIKey) (*APIKey, error)
	APIKeyDetail(ctx context.Context, id string) (*APIKey, error)
	APIKeyUpdate(context.Context, *APIKey) (*APIKey, error)
	APIKeyDelete(ctx context.Context, id string) error
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

//===========================================================================
// Tenant Requests and Responses
//===========================================================================

type PageQuery struct {
	PageSize      uint32 `url:"page_size,omitempty"`
	NextPageToken string `url:"next_page_token,omitempty"`
}

type Tenant struct {
	ID              string `json:"id" uri:"id"`
	Name            string `json:"name"`
	EnvironmentType string `json:"environment_type"`
}

type TenantPage struct {
	Tenants       []*Tenant
	PrevPageToken string
	NextPageToken string
}

type TenantMemberPage struct {
	TenantID      string `json:"tenant_id"`
	TenantMembers []*Member
	PrevPageToken string
	NextPageToken string
}

type Member struct {
	ID   string `json:"id" uri:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
}

type MemberPage struct {
	Members       []*Member
	PrevPageToken string
	NextPageToken string
}

type TenantProjectPage struct {
	TenantID       string `json:"id"`
	TenantProjects []*Project
	PrevPageToken  string
	NextPageToken  string
}

type Project struct {
	ID   string `json:"id" uri:"id"`
	Name string `json:"name"`
}

type ProjectPage struct {
	Projects      []*Project
	PrevPageToken string
	NextPageToken string
}

type ProjectTopicPage struct {
	ProjectID     string `json:"project_id"`
	TenantTopics  []*Topic
	PrevPageToken string
	NextPageToken string
}

type Topic struct {
	ID   string `json:"id" uri:"id"`
	Name string `json:"topic_name"`
}

type TopicPage struct {
	Topics        []*Topic
	PrevPageToken string
	NextPageToken string
}

type ProjectAPIKeyPage struct {
	ProjectID     string `json:"project_id"`
	APIKeys       []*APIKey
	PrevPageToken string
	NextPageToken string
}

type APIKey struct {
	ID           int      `json:"id,omitempty"`
	ClientID     string   `json:"client_id"`
	ClientSecret string   `json:"client_secret,omitempty"`
	Name         string   `json:"name"`
	Owner        string   `json:"owner,omitempty"`
	Permissions  []string `json:"permissions,omitempty"`
	Created      string   `json:"created,omitempty"`
	Modified     string   `json:"modified,omitempty"`
}

type APIKeyPage struct {
	APIKeys       []*APIKey
	PrevPageToken string
	NextPageToken string
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
