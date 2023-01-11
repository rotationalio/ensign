package db

import "errors"

var (
	ErrNotConnected         = errors.New("not connected to trtl database")
	ErrNotFound             = errors.New("object not found for the specified key")
	ErrMissingID            = errors.New("object requires id for serialization")
	ErrMissingTenantID      = errors.New("object requires tenant id for serialization")
	ErrMissingProjectID     = errors.New("object requires project id for serialization")
	ErrNumberFirstCharacter = errors.New("name cannot begin with a number")
	ErrSpecialCharacters    = errors.New("name cannot include special characters")
	ErrMissingMemberName    = errors.New("member name is required")
	ErrMissingProjectName   = errors.New("project name is required")
	ErrMissingTopicName     = errors.New("topic name is required")
)
