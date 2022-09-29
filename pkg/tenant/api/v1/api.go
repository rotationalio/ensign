package api

import "context"

//===========================================================================
// Service Interface
//===========================================================================

type TenantClient interface {
	Status(context.Context) (*StatusReply, error)
	SignUp(context.Context, *ContactInfo) error
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
// API Request and Reply Structs
//===========================================================================

// ContactInfo allows users to sign up for email notifications from SendGrid and is
// specifically used to allow users to request Ensign Private Beta access.
type ContactInfo struct {
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Email        string `json:"email"`
	Country      string `json:"country"`
	Title        string `json:"title"`
	Organization string `json:"organization"`
}
