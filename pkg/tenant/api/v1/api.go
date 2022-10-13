package api

import "context"

//===========================================================================
// Service Interface
//===========================================================================

type TenantClient interface {
	Status(context.Context) (*StatusReply, error)
	SignUp(context.Context, *ContactInfo) error

	UserList(context.Context, *UserQuery) (*UserPage, error)
	UserCreate(context.Context, *User) (*User, error)
	UserDetail(ctx context.Context, id string) (*User, error)
	UserUpdate(context.Context, *User) (*User, error)
	UserDelete(ctx context.Context, id string) error

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

type User struct {
	ID       string `json:"id" uri:"id" binding:"required"`
	UserName string `json:"user_name"`
}

type UserQuery struct {
	Query         string
	NextPageToken string
}

type UserPage struct {
	Users         []*User
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
	CloudServiceProvider string `json:"cloud_service_provider"`
}
