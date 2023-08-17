package api

import (
	"regexp"

	"github.com/oklog/ulid/v2"
	mimetype "github.com/rotationalio/ensign/pkg/ensign/mimetype/v1beta1"
	"github.com/rotationalio/ensign/pkg/ensign/rlid"
	"github.com/twmb/murmur3"
)

const (
	NameHashLength     = 16
	MaxTopicNameLength = 512
)

var topicNameRegex = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9\.\-\_]*$`)

// ParseTopicID returns the ULID representation of the topic ID.
func (t *Topic) ParseTopicID() (uid ulid.ULID, err error) {
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

func (i *TopicInfo) ParseTopicID() (topicID ulid.ULID, err error) {
	topicID = ulid.ULID{}
	if err = topicID.UnmarshalBinary(i.TopicId); err != nil {
		return topicID, err
	}
	return topicID, nil
}

func (i *TopicInfo) ParseProjectID() (projectID ulid.ULID, err error) {
	projectID = ulid.ULID{}
	if err = projectID.UnmarshalBinary(i.ProjectId); err != nil {
		return projectID, err
	}
	return projectID, nil
}

func (i *TopicInfo) ParseEventOffsetID() (eventID rlid.RLID, err error) {
	eventID = rlid.RLID{}
	if err = eventID.UnmarshalBinary(i.EventOffsetId); err != nil {
		return eventID, err
	}
	return eventID, nil
}

// Finds the event type info for the specified type in the type list. If it does not
// exist, the event type info is created an appended to the type list.
func (i *TopicInfo) FindEventTypeInfo(etype *Type, mime mimetype.MIME) *EventTypeInfo {
	// Look for existing event type info for the specified type
	for _, einfo := range i.Types {
		if einfo.Type.Equals(etype) && einfo.Mimetype == mime {
			return einfo
		}
	}

	// Create event type info for the specified type
	einfo := &EventTypeInfo{Type: etype, Mimetype: mime}
	i.Types = append(i.Types, einfo)
	return einfo
}
