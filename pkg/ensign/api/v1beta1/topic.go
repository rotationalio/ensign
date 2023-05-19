package api

import "github.com/twmb/murmur3"

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
