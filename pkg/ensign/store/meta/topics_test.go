package meta_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	api "github.com/rotationalio/ensign/pkg/ensign/api/v1beta1"
	"github.com/rotationalio/ensign/pkg/ensign/store/errors"
	"github.com/rotationalio/ensign/pkg/ensign/store/meta"
	"github.com/rotationalio/ensign/pkg/utils/ulids"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *metaTestSuite) TestListTopics() {
	require := s.Require()
	require.False(s.store.ReadOnly())

	_, err := s.LoadAllFixtures()
	require.NoError(err, "could not load all fixtures")
	defer s.ResetDatabase()

	topics := s.store.ListTopics(ulids.MustParse("01GTSMMC152Q95RD4TNYDFJGHT"))
	defer topics.Release()

	nTopics := 0
	for topics.Next() {
		nTopics++
		topic, err := topics.Topic()
		require.NoError(err, "could not deserialize topic")
		require.True(strings.HasPrefix(topic.Name, "testing.testapp"))
	}
	require.Equal(5, nTopics)

	err = topics.Error()
	require.NoError(err, "could not list topics from database")
}

func (s *readonlyMetaTestSuite) TestListTopics() {
	require := s.Require()
	require.True(s.store.ReadOnly())

	topics := s.store.ListTopics(ulids.MustParse("01GTSMMC152Q95RD4TNYDFJGHT"))
	defer topics.Release()

	nTopics := 0
	for topics.Next() {
		nTopics++
		topic, err := topics.Topic()
		require.NoError(err, "could not deserialize topic")
		require.True(strings.HasPrefix(topic.Name, "testing.testapp"))
	}
	require.Equal(5, nTopics)

	err := topics.Error()
	require.NoError(err, "could not list topics from database")
}

func (s *metaTestSuite) TestAllowedTopics() {
	require := s.Require()
	require.False(s.store.ReadOnly())

	_, err := s.LoadAllFixtures()
	require.NoError(err, "could not load all fixtures")
	defer s.ResetDatabase()

	topics, err := s.store.AllowedTopics(ulids.MustParse("01GTSMMC152Q95RD4TNYDFJGHT"))
	require.NoError(err, "could not fetch allowed topics")
	require.Len(topics, 5, "unexpected number of topics returned")
}

func (s *readonlyMetaTestSuite) TestAllowedTopics() {
	require := s.Require()
	require.True(s.store.ReadOnly())

	topics, err := s.store.AllowedTopics(ulids.MustParse("01GTSMMC152Q95RD4TNYDFJGHT"))
	require.NoError(err, "could not fetch allowed topics")
	require.Len(topics, 5, "unexpected number of topics returned")
}

func (s *metaTestSuite) TestListTopicsPagination() {
	require := s.Require()
	require.False(s.store.ReadOnly())

	_, err := s.LoadAllFixtures()
	require.NoError(err, "could not load all fixtures")
	defer s.ResetDatabase()

	topics := s.store.ListTopics(ulids.MustParse("01GTSMMC152Q95RD4TNYDFJGHT"))
	defer topics.Release()

	pages := 0
	items := 0
	info := &api.PageInfo{PageSize: uint32(2)}

	// Only paginate for a maximum of 10 iterations
	for i := 0; i < 10; i++ {
		page, err := topics.NextPage(info)
		require.NoError(err, "could not fetch page %d", i+1)
		require.LessOrEqual(len(page.Topics), int(info.PageSize))

		pages++
		items += len(page.Topics)

		if page.NextPageToken == "" {
			break
		}

		info.NextPageToken = page.NextPageToken
	}

	require.NoError(topics.Error(), "could not list topics from database")
	require.Equal(3, pages)
	require.Equal(5, items)
}

func (s *readonlyMetaTestSuite) TestListTopicsPagination() {
	require := s.Require()
	require.True(s.store.ReadOnly())

	topics := s.store.ListTopics(ulids.MustParse("01GTSMMC152Q95RD4TNYDFJGHT"))
	defer topics.Release()

	pages := 0
	items := 0
	info := &api.PageInfo{PageSize: uint32(2)}

	// Only paginate for a maximum of 10 iterations
	for i := 0; i < 10; i++ {
		page, err := topics.NextPage(info)
		require.NoError(err, "could not fetch page %d", i+1)
		require.LessOrEqual(len(page.Topics), int(info.PageSize))

		pages++
		items += len(page.Topics)

		if page.NextPageToken == "" {
			break
		}

		info.NextPageToken = page.NextPageToken
	}

	require.NoError(topics.Error(), "could not list topics from database")
	require.Equal(3, pages)
	require.Equal(5, items)
}

func (s *metaTestSuite) TestCreateTopic() {
	require := s.Require()
	require.False(s.store.ReadOnly())

	defer s.ResetDatabase()

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

	// Check to make sure the topic and the index entry have been created
	count, err = s.store.Count(nil)
	require.NoError(err, "could not count database")
	require.Equal(uint64(3), count, "expected 3 objects in the database")
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

func (s *metaTestSuite) TestCreateTopicUniqueness() {
	require := s.Require()
	defer s.ResetDatabase()

	// Should not be able to create a topic with the same name in the same project
	// Should be able to create a topic with the same name in a different project
	projectIDa := ulids.New()
	projectIDb := ulids.New()

	topica1 := &api.Topic{
		ProjectId: projectIDa.Bytes(),
		Name:      "test.alerts",
	}

	topica2 := &api.Topic{
		ProjectId: projectIDa.Bytes(),
		Name:      "test.alerts",
	}

	topicb := &api.Topic{
		ProjectId: projectIDb.Bytes(),
		Name:      "test.alerts",
	}

	// Should be able to create topicsa1 and topicb, creating topica2 should fail.
	err := s.store.CreateTopic(topica1)
	require.NoError(err, "could not create topic a1")

	err = s.store.CreateTopic(topica2)
	require.ErrorIs(err, errors.ErrUniqueTopicName)

	err = s.store.CreateTopic(topicb)
	require.NoError(err, "could not create topic b")

	// Check to make sure the topic and the index entry have been created
	count, err := s.store.Count(nil)
	require.NoError(err, "could not count database")
	require.Equal(uint64(6), count, "expected 6 objects in the database")
}

func (s *metaTestSuite) TestRetrieveTopic() {
	require := s.Require()
	require.False(s.store.ReadOnly())

	_, err := s.LoadAllFixtures()
	require.NoError(err, "could not load topic fixtures")
	defer s.ResetDatabase()

	topic, err := s.store.RetrieveTopic(ulids.MustParse("01GTSN1WF5BA0XCPT6ES64JVGQ"))
	require.NoError(err, "could not retrieve topic")
	require.Equal("mock.mockapp.feed", topic.Name)
}

func (s *readonlyMetaTestSuite) TestRetrieveTopic() {
	require := s.Require()
	require.True(s.store.ReadOnly())

	topic, err := s.store.RetrieveTopic(ulids.MustParse("01GTSN1139JMK1PS5A524FXWAZ"))
	require.NoError(err, "could not retrieve topic")
	require.Equal("testing.testapp.shipments", topic.Name)
}

func (s *metaTestSuite) TestUpdateTopic() {
	require := s.Require()
	require.False(s.store.ReadOnly())

	nFixtures, err := s.LoadTopicFixtures()
	require.NoError(err, "could not load topic fixtures")
	defer s.ResetDatabase()

	ts, err := time.Parse(time.RFC3339Nano, "2023-03-05T19:41:59.016422Z")
	require.NoError(err, "could not parse fixture timestamp")

	// Database should have the fixtures states to start
	count, err := s.store.Count(nil)
	require.NoError(err, "could not count database")
	require.Equal(nFixtures, count, "expected topic fixtures in the database")

	// Cannot update a topic that doesn't exist
	notreal := &api.Topic{
		ProjectId: ulids.MustBytes("01GTSMQ3V8ASAPNCFEN378T8RD"),
		Id:        ulids.MustBytes("01GTSMMC152Q95RD4TNYDFJGHT"),
		Name:      "not-real",
		Created:   timestamppb.Now(),
		Modified:  timestamppb.Now(),
	}

	err = s.store.UpdateTopic(notreal)
	require.ErrorIs(err, errors.ErrNotFound)

	// Should not be able to update a topic's name
	topic := &api.Topic{
		Id:        ulids.MustBytes("01GTSMQ3V8ASAPNCFEN378T8RD"),
		ProjectId: ulids.MustBytes("01GTSMMC152Q95RD4TNYDFJGHT"),
		Name:      "testing.testapp.modified_alerts",
		Created:   timestamppb.New(ts),
		Modified:  timestamppb.New(ts),
	}

	err = s.store.UpdateTopic(topic)
	require.ErrorIs(err, errors.ErrTopicNameChanged)

	// Should be able to update a topic that does exist
	topic.Name = "testing.testapp.alerts"
	topic.Offset = 90221
	err = s.store.UpdateTopic(topic)
	require.NoError(err, "could not update topic")

	// Check to ensure the topic has been changed
	compat, err := s.store.RetrieveTopic(ulids.MustParse("01GTSMQ3V8ASAPNCFEN378T8RD"))
	require.NoError(err, "unable to retrieve topic by id")
	require.True(proto.Equal(compat, topic))

	// Database should have the same number of fixtures states to finish
	count, err = s.store.Count(nil)
	require.NoError(err, "could not count database")
	require.Equal(nFixtures, count, "expected no change in the count of objects")
}

func (s *readonlyMetaTestSuite) TestUpdateTopic() {
	require := s.Require()
	require.True(s.store.ReadOnly())

	topic := &api.Topic{
		Id:        ulids.MustBytes("01GTSMQ3V8ASAPNCFEN378T8RD"),
		ProjectId: ulids.MustBytes("01GTSMMC152Q95RD4TNYDFJGHT"),
		Name:      "testing.testapp.alerts",
		Created:   timestamppb.Now(),
		Modified:  timestamppb.Now(),
	}

	err := s.store.UpdateTopic(topic)
	require.ErrorIs(err, errors.ErrReadOnly, "expected readonly error on create topic")
}

func (s *metaTestSuite) TestDeleteTopic() {
	require := s.Require()
	require.False(s.store.ReadOnly())

	nFixtures, err := s.LoadTopicFixtures()
	require.NoError(err, "could not load topic fixtures")
	defer s.ResetDatabase()

	// Database should have the fixtures states to start
	count, err := s.store.Count(nil)
	require.NoError(err, "could not count database")
	require.Equal(nFixtures, count, "expected topic fixtures in the database")

	err = s.store.DeleteTopic(ulids.MustParse("01GTSMSX1M9G2Z45VGG4M12WC0"))
	require.NoError(err, "Could not delete topic")

	// Index and topic should have been deleted
	count, err = s.store.Count(nil)
	require.NoError(err, "could not count database")
	require.Equal(nFixtures-3, count, "expected one less topic fixture and two less indices in the database")

	// Deleting a second time should have no effect
	err = s.store.DeleteTopic(ulids.MustParse("01GTSMSX1M9G2Z45VGG4M12WC0"))
	require.NoError(err, "Could not delete topic")

	count, err = s.store.Count(nil)
	require.NoError(err, "could not count database")
	require.Equal(nFixtures-3, count, "expected no change in database count")
}

func (s *readonlyMetaTestSuite) TestDeleteTopic() {
	require := s.Require()
	require.True(s.store.ReadOnly())

	err := s.store.DeleteTopic(ulids.MustParse("01GTSMQ3V8ASAPNCFEN378T8RD"))
	require.ErrorIs(err, errors.ErrReadOnly, "expected readonly error on create topic")
}

func TestTopicKey(t *testing.T) {
	topic := &api.Topic{
		Id:        ulids.MustBytes("01GTSSDM957VH0GX0RMNKAQM13"),
		ProjectId: ulids.MustBytes("01GTSSCWHMBNCVZFPBPQETXG96"),
	}

	key := meta.TopicKey(topic)
	require.Len(t, key, 34, "expected the key length to be two ulids long")
	require.True(t, bytes.HasPrefix(key[:], topic.ProjectId))
	require.True(t, bytes.Equal(key[16:18], meta.TopicSegment[:]))
	require.True(t, bytes.HasSuffix(key[:], topic.Id))
}

func TestValidateTopic(t *testing.T) {
	testCases := []struct {
		topic   *api.Topic
		partial bool
		err     error
	}{
		{
			nil,
			true,
			errors.ErrTopicInvalidId,
		},
		{
			nil,
			false,
			errors.ErrTopicInvalidId,
		},
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
				Name:      strings.Repeat("ABRACADABRA", 200),
				Created:   timestamppb.Now(),
				Modified:  timestamppb.Now(),
			},
			true,
			errors.ErrTopicNameTooLong,
		},
		{
			&api.Topic{
				Id:        []byte{1, 134, 179, 108, 62, 211, 134, 53, 49, 102, 31, 33, 40, 215, 58, 245},
				ProjectId: []byte{1, 134, 179, 81, 86, 251, 48, 108, 44, 19, 143, 243, 195, 87, 134, 80},
				Name:      strings.Repeat("ABRACADABRA", 200),
				Created:   timestamppb.Now(),
				Modified:  timestamppb.Now(),
			},
			false,
			errors.ErrTopicNameTooLong,
		},
		{
			&api.Topic{
				Id:        []byte{1, 134, 179, 108, 62, 211, 134, 53, 49, 102, 31, 33, 40, 215, 58, 245},
				ProjectId: []byte{1, 134, 179, 81, 86, 251, 48, 108, 44, 19, 143, 243, 195, 87, 134, 80},
				Name:      "1 foo $",
				Created:   timestamppb.Now(),
				Modified:  timestamppb.Now(),
			},
			true,
			errors.ErrInvalidTopicName,
		},
		{
			&api.Topic{
				Id:        []byte{1, 134, 179, 108, 62, 211, 134, 53, 49, 102, 31, 33, 40, 215, 58, 245},
				ProjectId: []byte{1, 134, 179, 81, 86, 251, 48, 108, 44, 19, 143, 243, 195, 87, 134, 80},
				Name:      "1 foo $",
				Created:   timestamppb.Now(),
				Modified:  timestamppb.Now(),
			},
			false,
			errors.ErrInvalidTopicName,
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

func TestValidateTopicName(t *testing.T) {
	valid := []string{
		"topic",
		"snake_case_topic",
		"CamelCaseTopic",
		"dot.separated.topic",
		"dash-separated-topic",
		"topic34",
		"blue_42_blue_42",
		"HutHutHIKE42",
		"dot.dash-underscore_topic",
	}

	invalid := []struct {
		name string
		err  error
	}{
		{"", errors.ErrTopicMissingName},
		{strings.Repeat("ABRACADBRA", 200), errors.ErrTopicNameTooLong},
		{"no spaces", errors.ErrInvalidTopicName},
		{"42noleadnums", errors.ErrInvalidTopicName},
		{"-noleadhypen", errors.ErrInvalidTopicName},
		{".noleaddot", errors.ErrInvalidTopicName},
		{"no$special!chars", errors.ErrInvalidTopicName},
	}

	for _, name := range valid {
		err := meta.ValidateTopicName(name)
		require.NoError(t, err, "expected %q to be a valid topic name", name)
	}

	for _, tc := range invalid {
		err := meta.ValidateTopicName(tc.name)
		require.ErrorIs(t, err, tc.err, "expected %q to be invalid", tc.name)
	}

}
