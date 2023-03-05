package meta_test

import (
	"bytes"
	"testing"
	"time"

	api "github.com/rotationalio/ensign/pkg/api/v1beta1"
	"github.com/rotationalio/ensign/pkg/ensign/store/errors"
	"github.com/rotationalio/ensign/pkg/ensign/store/meta"
	ulids "github.com/rotationalio/ensign/pkg/utils/ulid"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *metaTestSuite) TestListTopics() {
	require := s.Require()
	require.False(s.store.ReadOnly())
}

func (s *readonlyMetaTestSuite) TestListTopics() {
	require := s.Require()
	require.True(s.store.ReadOnly())
}

func (s *metaTestSuite) TestCreateTopic() {
	require := s.Require()
	require.False(s.store.ReadOnly())

	// Database should be empty to begin
	count, err := s.store.Count(nil)
	require.NoError(err, "could not count database")
	require.Equal(uint64(0), count, "expected no objects in the database")

	// Should not be able to create an empty topic
	err = s.store.CreateTopic(&api.Topic{})
	require.ErrorIs(err, errors.ErrInvalidTopic)

	// Should be able to create a valid topic
	topic := &api.Topic{
		ProjectId: ulids.MustBytes("01GTSRBV1HRZ3PPETSM3YF1N79"),
		Name:      "testing.testapp.test",
	}

	err = s.store.CreateTopic(topic)
	require.NoError(err, "expected to be able to create the valid topic")

	// Check to make sure the topic has been created
	count, err = s.store.Count(nil)
	require.NoError(err, "could not count database")
	require.Equal(uint64(1), count, "expected 1 objects in the database")
}

func (s *readonlyMetaTestSuite) TestCreateTopic() {
	require := s.Require()
	require.True(s.store.ReadOnly())

	topic := &api.Topic{
		ProjectId: ulids.MustBytes("01GTSRBV1HRZ3PPETSM3YF1N79"),
		Name:      "testing.testapp.test",
	}

	err := s.store.CreateTopic(topic)
	require.ErrorIs(err, errors.ErrReadOnly, "expected readonly error on create topic")
}

func (s *metaTestSuite) TestRetrieveTopic() {
	require := s.Require()
	require.False(s.store.ReadOnly())
}

func (s *readonlyMetaTestSuite) TestRetrieveTopic() {
	require := s.Require()
	require.True(s.store.ReadOnly())
}

func (s *metaTestSuite) TestUpdateTopic() {
	require := s.Require()
	require.False(s.store.ReadOnly())
}

func (s *readonlyMetaTestSuite) TestUpdateTopic() {
	require := s.Require()
	require.True(s.store.ReadOnly())
}

func (s *metaTestSuite) TestDeleteTopic() {
	require := s.Require()
	require.False(s.store.ReadOnly())
}

func (s *readonlyMetaTestSuite) TestDeleteTopic() {
	require := s.Require()
	require.True(s.store.ReadOnly())
}

func TestTopicKey(t *testing.T) {
	topic := &api.Topic{
		Id:        ulids.MustBytes("01GTSSDM957VH0GX0RMNKAQM13"),
		ProjectId: ulids.MustBytes("01GTSSCWHMBNCVZFPBPQETXG96"),
	}

	key := meta.TopicKey(topic)
	require.Len(t, key, 32, "expected the key length to be two ulids long")
	require.True(t, bytes.HasPrefix(key, topic.ProjectId))
	require.True(t, bytes.HasSuffix(key, topic.Id))
}

func TestValidateTopic(t *testing.T) {
	testCases := []struct {
		topic   *api.Topic
		partial bool
		err     error
	}{
		{
			&api.Topic{
				Id:       []byte{1, 134, 179, 81, 86, 251, 48, 108, 44, 19, 143, 243, 195, 87, 134, 80},
				Name:     "testing.testapp.test",
				Created:  timestamppb.Now(),
				Modified: timestamppb.Now(),
			},
			true,
			errors.ErrTopicMissingProjectId,
		},
		{
			&api.Topic{
				Id:       []byte{1, 134, 179, 81, 86, 251, 48, 108, 44, 19, 143, 243, 195, 87, 134, 80},
				Name:     "testing.testapp.test",
				Created:  timestamppb.Now(),
				Modified: timestamppb.Now(),
			},
			false,
			errors.ErrTopicMissingProjectId,
		},

		{
			&api.Topic{
				Id:        []byte{1, 134, 179, 81, 86, 251, 48, 108, 44, 19, 143, 243, 195, 87, 134, 80},
				ProjectId: []byte{42},
				Name:      "testing.testapp.test",
				Created:   timestamppb.Now(),
				Modified:  timestamppb.Now(),
			},
			true,
			errors.ErrTopicInvalidProjectId,
		},
		{
			&api.Topic{
				Id:        []byte{1, 134, 179, 81, 86, 251, 48, 108, 44, 19, 143, 243, 195, 87, 134, 80},
				ProjectId: []byte{42},
				Name:      "testing.testapp.test",
				Created:   timestamppb.Now(),
				Modified:  timestamppb.Now(),
			},
			false,
			errors.ErrTopicInvalidProjectId,
		},
		{
			&api.Topic{
				Id:        []byte{1, 134, 179, 108, 62, 211, 134, 53, 49, 102, 31, 33, 40, 215, 58, 245},
				ProjectId: []byte{1, 134, 179, 81, 86, 251, 48, 108, 44, 19, 143, 243, 195, 87, 134, 80},
				Name:      "",
				Created:   timestamppb.Now(),
				Modified:  timestamppb.Now(),
			},
			true,
			errors.ErrTopicMissingName,
		},
		{
			&api.Topic{
				Id:        []byte{1, 134, 179, 108, 62, 211, 134, 53, 49, 102, 31, 33, 40, 215, 58, 245},
				ProjectId: []byte{1, 134, 179, 81, 86, 251, 48, 108, 44, 19, 143, 243, 195, 87, 134, 80},
				Name:      "",
				Created:   timestamppb.Now(),
				Modified:  timestamppb.Now(),
			},
			false,
			errors.ErrTopicMissingName,
		},
		{
			&api.Topic{
				Id:        []byte{},
				ProjectId: []byte{1, 134, 179, 81, 86, 251, 48, 108, 44, 19, 143, 243, 195, 87, 134, 80},
				Name:      "testing.testapp.test",
				Created:   timestamppb.Now(),
				Modified:  timestamppb.Now(),
			},
			false,
			errors.ErrTopicMissingId,
		},
		{
			&api.Topic{
				Id:        []byte{},
				ProjectId: []byte{1, 134, 179, 81, 86, 251, 48, 108, 44, 19, 143, 243, 195, 87, 134, 80},
				Name:      "testing.testapp.test",
				Created:   timestamppb.Now(),
				Modified:  timestamppb.Now(),
			},
			true,
			nil,
		},
		{
			&api.Topic{
				Id:        []byte{42},
				ProjectId: []byte{1, 134, 179, 81, 86, 251, 48, 108, 44, 19, 143, 243, 195, 87, 134, 80},
				Name:      "testing.testapp.test",
				Created:   timestamppb.Now(),
				Modified:  timestamppb.Now(),
			},
			false,
			errors.ErrTopicInvalidId,
		},
		{
			&api.Topic{
				Id:        []byte{42},
				ProjectId: []byte{1, 134, 179, 81, 86, 251, 48, 108, 44, 19, 143, 243, 195, 87, 134, 80},
				Name:      "testing.testapp.test",
				Created:   timestamppb.Now(),
				Modified:  timestamppb.Now(),
			},
			true,
			errors.ErrTopicInvalidId,
		},
		{
			&api.Topic{
				Id:        []byte{1, 134, 179, 108, 62, 211, 134, 53, 49, 102, 31, 33, 40, 215, 58, 245},
				ProjectId: []byte{1, 134, 179, 81, 86, 251, 48, 108, 44, 19, 143, 243, 195, 87, 134, 80},
				Name:      "testing.testapp.test",
				Created:   timestamppb.New(time.Time{}),
				Modified:  timestamppb.Now(),
			},
			false,
			errors.ErrTopicInvalidCreated,
		},
		{
			&api.Topic{
				Id:        []byte{1, 134, 179, 108, 62, 211, 134, 53, 49, 102, 31, 33, 40, 215, 58, 245},
				ProjectId: []byte{1, 134, 179, 81, 86, 251, 48, 108, 44, 19, 143, 243, 195, 87, 134, 80},
				Name:      "testing.testapp.test",
				Created:   timestamppb.New(time.Time{}),
				Modified:  timestamppb.Now(),
			},
			true,
			nil,
		},
		{
			&api.Topic{
				Id:        []byte{1, 134, 179, 108, 62, 211, 134, 53, 49, 102, 31, 33, 40, 215, 58, 245},
				ProjectId: []byte{1, 134, 179, 81, 86, 251, 48, 108, 44, 19, 143, 243, 195, 87, 134, 80},
				Name:      "testing.testapp.test",
				Created:   timestamppb.Now(),
				Modified:  nil,
			},
			false,
			errors.ErrTopicInvalidModified,
		},
		{
			&api.Topic{
				Id:        []byte{1, 134, 179, 108, 62, 211, 134, 53, 49, 102, 31, 33, 40, 215, 58, 245},
				ProjectId: []byte{1, 134, 179, 81, 86, 251, 48, 108, 44, 19, 143, 243, 195, 87, 134, 80},
				Name:      "testing.testapp.test",
				Created:   timestamppb.Now(),
				Modified:  nil,
			},
			true,
			nil,
		},
	}

	for i, tc := range testCases {
		err := meta.ValidateTopic(tc.topic, tc.partial)
		if tc.err == nil {
			require.NoError(t, err, "failed testcase %d -- expected no error", i)
		} else {
			require.ErrorIs(t, err, tc.err, "failed testcase %d -- expected matching error", i)
		}
	}
}
