package api

const (
	CodeUnknown              = "unknown error"
	CodeMaxEventSizeExceeded = "maximum event size has been exceeded"
	CodeTopicUnknown         = "topic unknown or unhandled"
	CodeTopicArchived        = "cannot publish to an archived topic"
	CodeTopicDeleted         = "topic is currently being deleted"
	CodePermissionDenied     = "not authorized to perform this action"
	CodeConsensusFailure     = "internal replication error"
	CodeShardingFailure      = "wrong node for event sharding policy"
	CodeRedirect             = "wrong ensign node"
	CodeInternal             = "internal error"
	CodeUnprocessed          = "client did not process event"
	CodeTimeout              = "client deadline exceeded"
	CodeUnhandledMimetype    = "unhandled mimetype"
	CodeUnknownType          = "unhandled schema"
	CodeDeliverAgainAny      = "deliver again to any subscriber"
	CodeDeliverAgainNotMe    = "deliver agin to any subscriber but me"
)

func DefaultNackMessage(code Nack_Code) string {
	switch code {
	case Nack_MAX_EVENT_SIZE_EXCEEDED:
		return CodeMaxEventSizeExceeded
	case Nack_TOPIC_UNKNOWN:
		return CodeTopicUnknown
	case Nack_TOPIC_ARCHIVED:
		return CodeTopicArchived
	case Nack_TOPIC_DELETED:
		return CodeTopicDeleted
	case Nack_PERMISSION_DENIED:
		return CodePermissionDenied
	case Nack_CONSENSUS_FAILURE:
		return CodeConsensusFailure
	case Nack_SHARDING_FAILURE:
		return CodeShardingFailure
	case Nack_REDIRECT:
		return CodeRedirect
	case Nack_INTERNAL:
		return CodeInternal
	case Nack_UNPROCESSED:
		return CodeUnprocessed
	case Nack_TIMEOUT:
		return CodeTimeout
	case Nack_UNHANDLED_MIMETYPE:
		return CodeUnhandledMimetype
	case Nack_UNKNOWN_TYPE:
		return CodeUnknownType
	case Nack_DELIVER_AGAIN_ANY:
		return CodeDeliverAgainAny
	case Nack_DELIVER_AGAIN_NOT_ME:
		return CodeDeliverAgainNotMe
	default:
		return CodeUnknown
	}
}
