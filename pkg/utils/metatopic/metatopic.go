package metatopic

import (
	"regexp"
	"strconv"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/vmihailenco/msgpack/v5"
)

// Constants useful for creating ensign Events.
const (
	Mimetype      = "application/msgpack"
	SchemaName    = "metatopic.TopicUpdate"
	SchemaVersion = "1.0.0"
)

// This is the top-level event type that is sent on the metatopic topic. It describes an
// update to a specific topic and provides information about the topic and how it was
// modified, and, if available, who it was modified by.
type TopicUpdate struct {
	UpdateType TopicUpdateType `msgpack:"update_type"`
	OrgID      ulid.ULID       `msgpack:"org_id,omitempty"`
	ProjectID  ulid.ULID       `msgpack:"project_id"`
	TopicID    ulid.ULID       `msgpack:"topic_id"`
	ClientID   string          `msgpack:"client_id"`
	Topic      *Topic          `msgpack:"topic,omitempty"`
}

// A non-protocol buffer representation of the Topic. In the Topic Update struct it
// represents the modified topic (e.g. the current version of the topic).
// TODO: add placements and types to this struct.
type Topic struct {
	ID        []byte    `msgpack:"id"`
	ProjectID []byte    `msgpack:"project_id"`
	Name      string    `msgpack:"name"`
	ReadOnly  bool      `msgpack:"readonly"`
	Offset    uint64    `msgpack:"offset"`
	Shards    uint32    `msgpack:"shards"`
	Created   time.Time `msgpack:"created"`
	Modified  time.Time `msgpack:"modified"`
}

// The type of update made to the topic, e.g. created, modified, deleted, etc.
type TopicUpdateType uint8

const (
	TopicUpdateUnknown TopicUpdateType = iota
	TopicUpdateCreated
	TopicUpdateModified
	TopicUpdateStateChange
	TopicUpdateDeleted
)

var topicUpdateTypeNames = []string{
	"unknown", "created", "modified", "state_change", "deleted",
}

func (t TopicUpdateType) String() string {
	return topicUpdateTypeNames[t]
}

func (t *TopicUpdate) Marshal() ([]byte, error) {
	return msgpack.Marshal(t)
}

func (t *TopicUpdate) Unmarshal(data []byte) error {
	return msgpack.Unmarshal(data, t)
}

var semver = regexp.MustCompile(`^(\d+)\.(\d+)\.(\d+)$`)

// Parse the version components of the schema version to create an event version.
func ParseVersion() (major, minor, patch uint32) {
	if !semver.MatchString(SchemaVersion) {
		panic("cannot parse schema version")
	}

	groups := semver.FindStringSubmatch(SchemaVersion)
	if len(groups) != 4 {
		panic("cannot parse schema version - not enough digits")
	}

	if num, err := strconv.ParseUint(groups[1], 10, 32); err != nil {
		panic("could not parse major schema version component")
	} else {
		major = uint32(num)
	}

	if num, err := strconv.ParseUint(groups[2], 10, 32); err != nil {
		panic("could not parse minor schema version component")
	} else {
		minor = uint32(num)
	}

	if num, err := strconv.ParseUint(groups[3], 10, 32); err != nil {
		panic("could not parse patch schema version component")
	} else {
		patch = uint32(num)
	}

	return major, minor, patch
}
