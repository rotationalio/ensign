package db

import (
	"errors"
	"fmt"
)

var (
	ErrNotConnected       = errors.New("not connected to trtl database")
	ErrNotFound           = errors.New("object not found for the specified key")
	ErrUnavailable        = errors.New("trtl database service is unavailable")
	ErrMissingID          = errors.New("object requires id for serialization")
	ErrMissingOrgID       = errors.New("object requires organization id for serialization")
	ErrMissingTenantID    = errors.New("object requires tenant id for serialization")
	ErrMissingProjectID   = errors.New("object requires project id for serialization")
	ErrMissingMemberName  = errors.New("member name is required")
	ErrMissingProjectName = errors.New("project name is required")
	ErrMissingTenantName  = errors.New("tenant name is required")
	ErrMissingEnvType     = errors.New("tenant environment type is required")
	ErrMissingTopicName   = errors.New("topic name is required")
)

func ValidatonError(model string) error {
	return fmt.Errorf("%s name cannot begin with a number or include a special character", model)
}
