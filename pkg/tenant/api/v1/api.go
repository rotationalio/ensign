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

	TenantProjectList(ctx context.Context, id string, in PageQuery) (*TenantProjectPage, error)
	TenantProjectCreate(ctx context.Context, id string, in *Project) (*Project, error)

	ProjectList(context.Context, PageQuery) (*ProjectPage, error)
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

type TenantProjectPage struct {
	TenantID       string `json:"id"`
	TenantProjects []*Project
	PrevPageToken  string
	NextPageToken  string
}

type Project struct {
	ID   string `json:"id" uri:"id" binding:"required"`
	Name string `json:"name"`
}

type ProjectPage struct {
	Projects      []*Project
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
