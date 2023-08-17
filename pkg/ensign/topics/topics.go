/*
Package topics provides some helpers for managing topics in memory.
*/
package topics

import (
	"bytes"
	"encoding/base64"
	"errors"

	"github.com/oklog/ulid/v2"
	api "github.com/rotationalio/ensign/pkg/ensign/api/v1beta1"
	"github.com/rotationalio/ensign/pkg/utils/ulids"
)

var (
	ErrAlreadyExists  = errors.New("topic name or ID already exists in group")
	ErrEmptyReference = errors.New("topic name and ID cannot be empty")
)

// A topics group is a set of related topics (e.g. topics that belong to the same
// project) that allow easy lookups between topic names and IDs for easy referencing.
// TODO: should name groups also handle name hashes?
type NameGroup struct {
	names map[string]ulid.ULID
	ids   map[ulid.ULID]string
}

// Add a topic reference to the names group consisting of the topic name and ID.
func (g *NameGroup) Add(name string, id ulid.ULID) error {
	// Will panic if name is nil and ids is not or vice versa.
	if g.names == nil && g.ids == nil {
		g.names = make(map[string]ulid.ULID)
		g.ids = make(map[ulid.ULID]string)
	}

	if name == "" || ulids.IsZero(id) {
		return ErrEmptyReference
	}

	if _, ok := g.names[name]; ok {
		return ErrAlreadyExists
	}

	if _, ok := g.ids[id]; ok {
		return ErrAlreadyExists
	}

	g.names[name] = id
	g.ids[id] = name
	return nil
}

// Add a topic reference from the topic api struct.
func (g *NameGroup) AddTopic(topic *api.Topic) (err error) {
	var topicID ulid.ULID
	if topicID, err = topic.ParseTopicID(); err != nil {
		return err
	}

	return g.Add(topic.Name, topicID)
}

// Contains checks if the string is contained by the name group. It first checks to see
// if the string is a valid topic name, and if so it checks the names hash; then it
// checks if the string is a parseable ulid, and if so it checks the ID field. Finally,
// it checks if the string is a base64 encoded topic hash, and checks the name hashes.
func (g *NameGroup) Contains(s string) bool {
	// Is this a valid topic name?
	if api.ValidTopicName(s) {
		if g.ContainsTopicName(s) {
			return true
		}
	}

	// Is this a ULID string?
	if topicID, err := ulid.Parse(s); err == nil {
		if g.ContainsTopicID(topicID) {
			return true
		}
	}

	// Is this a base64 encoded topic hash?
	if hash, err := base64.RawURLEncoding.DecodeString(s); err == nil {
		if g.ContainsTopicHash(hash) {
			return true
		}
	}

	return false
}

// Lookup checks if the string is contained in the name group with similar semantics to
// Contains. If found, it will return the name and topicID, otherwise empty values.
func (g *NameGroup) Lookup(s string) (name string, topicID ulid.ULID, ok bool) {
	if api.ValidTopicName(s) {
		if topicID, ok = g.names[s]; ok {
			return s, topicID, ok
		}
	}

	if topicID, err := ulid.Parse(s); err == nil {
		if name, ok := g.ids[topicID]; ok {
			return name, topicID, ok
		}
	}

	if hash, err := base64.RawURLEncoding.DecodeString(s); err == nil {
		if name, topicID, ok = g.LookupTopicHash(hash); ok {
			return name, topicID, ok
		}
	}

	return "", ulid.ULID{}, false
}

// Check if the name group contains the specified topicID
func (g *NameGroup) ContainsTopicID(topicID ulid.ULID) bool {
	_, ok := g.ids[topicID]
	return ok
}

// Lookup the name of a topic by its ID
func (g *NameGroup) LookupTopicID(topicID ulid.ULID) (string, bool) {
	name, ok := g.ids[topicID]
	return name, ok
}

// Check if the name group contains the specified topic name
func (g *NameGroup) ContainsTopicName(name string) bool {
	_, ok := g.names[name]
	return ok
}

// Lookup the ID of a topic by its name
func (g *NameGroup) LookupTopicName(name string) (ulid.ULID, bool) {
	topicID, ok := g.names[name]
	return topicID, ok
}

// Check if the name group contains the specified topic hash.
func (g *NameGroup) ContainsTopicHash(hash []byte) bool {
	if len(hash) != api.NameHashLength {
		return false
	}

	for name := range g.names {
		nameHash := api.TopicNameHash(name)
		if bytes.Equal(nameHash, hash) {
			return true
		}
	}
	return false
}

// Lookup the name and topicID for the specified topic name hash.
func (g *NameGroup) LookupTopicHash(hash []byte) (string, ulid.ULID, bool) {
	if len(hash) != api.NameHashLength {
		return "", ulid.ULID{}, false
	}

	for name := range g.names {
		nameHash := api.TopicNameHash(name)
		if bytes.Equal(nameHash, hash) {
			return name, g.names[name], true
		}
	}

	return "", ulid.ULID{}, false
}

// Filter the topics name group by topic names, IDs, or base64 encoded hashes. This is
// the primary way to filter a topic name group from the user. The returned name group
// is a subset of the topics that are both in the original name group and specified by
// the list of topics. The original NameGroup is not modified.
func (g *NameGroup) Filter(topics ...string) *NameGroup {
	filtered := &NameGroup{
		names: make(map[string]ulid.ULID),
		ids:   make(map[ulid.ULID]string),
	}

	for _, topic := range topics {
		if name, topicID, ok := g.Lookup(topic); ok {
			filtered.Add(name, topicID)
		}
	}

	return filtered
}

// Filter the topics name group by the specified topicIDs, returning the subset of
// topics that are both in the original named group and specified by the list of IDs.
// E.g. if the topicID is in the original name group it is kept, otherwise it is
// omitted. A new NameGroup is returned, the original is not modified.
func (g *NameGroup) FilterTopicID(topicIDs ...ulid.ULID) *NameGroup {
	filtered := &NameGroup{
		names: make(map[string]ulid.ULID),
		ids:   make(map[ulid.ULID]string),
	}

	for _, topicID := range topicIDs {
		if name, ok := g.ids[topicID]; ok {
			filtered.Add(name, topicID)
		}
	}

	return filtered
}

// Filter the topics name group by the specified topic names, returning the subset of
// topics that are both in the original named group and specified by the list of names.
// E.g. if the name is in the original name group it is kept, otherwise it is omitted.
// A new NameGroup is returned, the original is not modified.
func (g *NameGroup) FilterTopicName(names ...string) *NameGroup {
	filtered := &NameGroup{
		names: make(map[string]ulid.ULID),
		ids:   make(map[ulid.ULID]string),
	}

	for _, name := range names {
		if topicID, ok := g.names[name]; ok {
			filtered.Add(name, topicID)
		}
	}

	return filtered
}

// TopicMap returns a map of topic name to topic ID bytes, which is used in StreamReady
// messages from the server, and also to easily perform lookups in leveldb indices.
func (g *NameGroup) TopicMap() map[string][]byte {
	topics := make(map[string][]byte)
	for id, name := range g.ids {
		topics[name] = id.Bytes()
	}
	return topics
}

// TopicIDs returns a slice of all of the topic ULIDs in the map.
func (g *NameGroup) TopicIDs() []ulid.ULID {
	topics := make([]ulid.ULID, 0, g.Length())
	for topicID := range g.ids {
		topics = append(topics, topicID)
	}
	return topics
}

// Returns the number of items in the name group.
func (g *NameGroup) Length() int {
	if len(g.ids) != len(g.names) {
		panic("name group has been corrupted")
	}
	return len(g.ids)
}

// Returns a base64 encoded string of the topic name hash.
func NameHash(name string) string {
	return base64.RawURLEncoding.EncodeToString(api.TopicNameHash(name))
}
