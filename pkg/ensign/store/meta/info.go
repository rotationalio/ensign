package meta

import (
	"errors"

	"github.com/oklog/ulid/v2"
	api "github.com/rotationalio/ensign/pkg/ensign/api/v1beta1"
)

func (s *Store) TopicInfo(topicID ulid.ULID) (*api.TopicInfo, error) {
	return nil, errors.New("not implemented yet")
}

func (s *Store) CreateTopicInfo(*api.TopicInfo) error {
	return errors.New("not implemented yet")
}

// To update topic info specify info with deltas that you want to add to the values
// that are currently stored in the info store so that the info can be updated in
// place in a single transaction without concurrency issues.
func (s *Store) UpdateTopicInfo(deltas *api.TopicInfo) error {
	return errors.New("not implemented yet")
}
