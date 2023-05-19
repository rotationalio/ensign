package api

import (
	"regexp"

	"github.com/oklog/ulid/v2"
	"github.com/twmb/murmur3"
)

const (
	NameHashLength     = 16
	MaxTopicNameLength = 512
)

var topicNameRegex = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9\.\-\_]*$`)

// ULID returns the ULID representation of the topic ID.
func (t *Topic) ULID() (uid ulid.ULID, err error) {
	uid = ulid.ULID{}
	if err = uid.UnmarshalBinary(t.Id); err != nil {
		return uid, err
	}
	return uid, nil
}

// NameHash returns an indexable hash of the topic name using murmur3.
func (t *Topic) NameHash() []byte {
	hash := murmur3.New128()
	hash.Write([]byte(t.Name))
	return hash.Sum(nil)
}

// TopicNameHash returns an indexable hash of a topic name using murmur3.
func TopicNameHash(name string) []byte {
	hash := murmur3.New128()
	hash.Write([]byte(name))
	return hash.Sum(nil)
}

// ValidTopicName returns true if the string is usable as a topic name.
func ValidTopicName(name string) bool {
	if name == "" {
		return false
	}

	if len(name) > MaxTopicNameLength {
		return false
	}

	if !topicNameRegex.MatchString(name) {
		return false
	}

	return true
}
