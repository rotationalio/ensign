package broker

import "errors"

var (
	ErrBrokerNotRunning = errors.New("operation could not be completed: broker is not running")
	ErrUnknownID        = errors.New("no publisher or subscriber registered with specified id")
)

// Standard error messages for nack codes
const (
	NackMaxEventSizeExceeded = "the maximum size for event data was exceeded, the event could not be processed"
	NackTopicUnknown         = "could not publish event to unknown topic"
	NackTopicArchived        = "cannot publish event to archived/readonly topic"
	NackTopicDeleted         = "cannot publish event to a deleted topic"
	NackPermissionDenied     = "client does not have permission to perform this operation"
	NackConsensusFailure     = "could not commit event, please try again"
	NackShardingFailure      = "unable to assign event key to a shard, please try again"
	NackRedirect             = "redirect to node handling topic or shard, please close and reconnect and try again"
	NackInternal             = "an internal error occurred, please try again shortly"
)
