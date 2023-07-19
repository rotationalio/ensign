package metatopic

import "errors"

var (
	// TopicUpdate validation errors.
	ErrMissingOrgID      = errors.New("missing organization id")
	ErrMissingProjectID  = errors.New("missing project id")
	ErrMissingTopicID    = errors.New("missing topic id")
	ErrMissingTopic      = errors.New("missing topic details")
	ErrUnknownUpdateType = errors.New("unknown topic update type")

	// Topic validation errors.
	ErrMissingName        = errors.New("missing topic name")
	ErrInvalidEvents      = errors.New("total events must be greater than or equal to 0")
	ErrInvalidStorage     = errors.New("data storage must be greater than or equal to 0")
	ErrMissingPublishers  = errors.New("missing topic publishers")
	ErrMissingSubscribers = errors.New("missing topic subscribers")
	ErrMissingCreated     = errors.New("missing topic created timestamp")
	ErrMissingModified    = errors.New("missing topic modified timestamp")
)
