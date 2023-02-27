package db

import (
	"errors"
)

var (
	ErrNotConnected = errors.New("not connected to trtl database")
	ErrNotFound     = errors.New("object not found for the specified key")
	ErrUnavailable  = errors.New("trtl database service is unavailable")

	// Missing fields
	ErrMissingID          = errors.New("object requires id for serialization")
	ErrMissingOrgID       = errors.New("object requires organization id for serialization")
	ErrMissingTenantID    = errors.New("object requires tenant id for serialization")
	ErrMissingProjectID   = errors.New("object requires project id for serialization")
	ErrMissingMemberName  = errors.New("member name is required")
	ErrMissingMemberRole  = errors.New("member role is required")
	ErrMissingProjectName = errors.New("project name is required")
	ErrMissingTenantName  = errors.New("tenant name is required")
	ErrMissingEnvType     = errors.New("tenant environment type is required")
	ErrMissingTopicName   = errors.New("topic name is required")

	// Invalid fields
	ErrInvalidMemberName  = errors.New("invalid member name")
	ErrUnknownMemberRole  = errors.New("unknown member role")
	ErrInvalidProjectName = errors.New("invalid project name")
	ErrInvalidTenantName  = errors.New("invalid tenant name")
	ErrInvalidTopicName   = errors.New("invalid topic name")
)
