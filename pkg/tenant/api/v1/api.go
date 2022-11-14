package api

import "context"

//===========================================================================
// Service Interface
//===========================================================================

type TenantClient interface {
	Status(context.Context) (*StatusReply, error)
	SignUp(context.Context, *ContactInfo) error

	TenantList(context.Context, *PageQuery) (*TenantPage, error)
	TenantCreate(context.Context, *Tenant) error
	TenantDetail(ctx context.Context, id string) (*Tenant, error)
	TenantUpdate(context.Context, *Tenant) (*Tenant, error)
	TenantDelete(ctx context.Context, id string) error

	AppList(context.Context, *AppQuery) (*AppPage, error)
	AppCreate(context.Context, *App) (*App, error)
	AppDetail(ctx context.Context, id string) (*App, error)
	AppDelete(ctx context.Context, id string) error

	TopicList(context.Context, *TopicQuery) (*TopicPage, error)
	TopicCreate(context.Context, *Topic) (*Topic, error)
	TopicDetail(ctx context.Context, id string) (*Topic, error)
	TopicDelete(ctx context.Context, id string) error
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

type Tenant struct {
	ID              string `json:"id" uri:"id" binding:"required"`
	TenantName      string `json:"tenant_name"`
	EnvironmentType string `json:"environment_type"`
}

type PageQuery struct {
	PageSize      uint32 `url:"page_size,omitempty"`
	NextPageToken string `url:"next_page_token,omitempty"`
}

type TenantPage struct {
	Tenants       []*Tenant
	PrevPageToken string
	NextPageToken string
}

type App struct {
	ID      string `json:"id" uri:"id" binding:"required"`
	AppName string `json:"app_name"`
}

type AppQuery struct {
	Query         string
	NextPageToken string
}

type AppPage struct {
	Apps          []*App
	PrevPageToken string
	NextPageToken string
}

type Topic struct {
	ID        string `json:"id" uri:"id" binding:"required"`
	TopicName string `json:"topic_name"`
}

type TopicQuery struct {
	Query         string
	NextPageToken string
}

type TopicPage struct {
	Topics        []*Topic
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
