package api

import (
	"context"
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

	// Projects Resource
	ProjectList(context.Context, *PageQuery) (*ProjectList, error)
	ProjectCreate(context.Context, *Project) (*Project, error)
	ProjectDetail(context.Context, string) (*Project, error)
	ProjectUpdate(context.Context, *Project) (*Project, error)
	ProjectDelete(context.Context, string) error

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
	APIKey string `json:"api_key"`
}

//===========================================================================
// Project Resource
//===========================================================================

type Project struct {
	ID int `json:"id"`
}

type ProjectList struct{}

//===========================================================================
// API Key Resource
//===========================================================================

type APIKey struct {
	ID int `json:"id"`
}

type APIKeyList struct{}
