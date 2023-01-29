package db

import "errors"

var (
	ErrNotConnected       = errors.New("not connected to trtl database")
	ErrNotFound           = errors.New("object not found for the specified key")
	ErrUnavailable        = errors.New("trtl database service is unavailable")
	ErrMissingID          = errors.New("object requires id for serialization")
	ErrMissingOrgID       = errors.New("object requires organization id for serialization")
	ErrMissingTenantID    = errors.New("object requires tenant id for serialization")
	ErrMissingProjectID   = errors.New("object requires project id for serialization")
	ErrValidation         = errors.New("name cannot begin with a number or include a special character")
	ErrMissingMemberName  = errors.New("member name is required")
	ErrMissingProjectName = errors.New("project name is required")
	ErrMissingTopicName   = errors.New("topic name is required")
)
