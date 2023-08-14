package events_test

import (
	"github.com/oklog/ulid/v2"
	api "github.com/rotationalio/ensign/pkg/ensign/api/v1beta1"
	"github.com/rotationalio/ensign/pkg/ensign/rlid"
	"github.com/rotationalio/ensign/pkg/ensign/store/errors"
)

func (s *eventsTestSuite) TestInsert() {
	require := s.Require()
	require.False(s.store.ReadOnly())

	defer s.ResetDatabase()

	// Database should be empty to begin
	count, err := s.store.Count(nil)
	require.NoError(err, "could not count database")
	require.Equal(uint64(0), count, "expected no objects in the database")

	event := &api.EventWrapper{
		Id:      rlid.Make(100).Bytes(),
		TopicId: ulid.Make().Bytes(),
		LocalId: ulid.Make().Bytes(),
	}

	err = s.store.Insert(event)
	require.NoError(err, "could not insert event into the database")

	// Check to make sure the event was inserted into the database
	count, err = s.store.Count(nil)
	require.NoError(err, "could not count database")
	require.Equal(uint64(1), count, "expected an event inserted into the database")

	// Ensure that the localID was set to nil
	require.Nil(event.LocalId, "expected local id to be nil")

	// Cannot insert an empty event
	err = s.store.Insert(&api.EventWrapper{})
	require.Error(err, "expected error trying to insert an empty event")
}

func (s *readonlyEventsTestSuite) TestInsert() {
	require := s.Require()
	require.True(s.store.ReadOnly())

	event := &api.EventWrapper{
		Id:      rlid.Make(100).Bytes(),
		TopicId: ulid.Make().Bytes(),
		LocalId: ulid.Make().Bytes(),
	}

	err := s.store.Insert(event)
	require.ErrorIs(err, errors.ErrReadOnly, "expected readonly error on insert event")
}

func (s *eventsTestSuite) TestList() {
	require := s.Require()
	require.False(s.store.ReadOnly())

	_, err := s.LoadAllFixtures()
	require.NoError(err, "could not load fixtures")
	defer s.ResetDatabase()

	topicID := ulid.MustParse("01GTSN1139JMK1PS5A524FXWAZ")
	iter := s.store.List(topicID)
	defer iter.Release()
}

func (s *readonlyEventsTestSuite) TestList() {
	require := s.Require()
	require.True(s.store.ReadOnly())

	topicID := ulid.MustParse("01GTSN1139JMK1PS5A524FXWAZ")
	iter := s.store.List(topicID)
	defer iter.Release()
}

func (s *eventsTestSuite) TestRetrieve() {
	require := s.Require()
	require.False(s.store.ReadOnly())

	_, err := s.LoadAllFixtures()
	require.NoError(err, "could not load fixtures")
	defer s.ResetDatabase()

	// Should be able to retrieve an event in the database
	event, err := s.store.Retrieve(ulid.MustParse("01GTSN1139JMK1PS5A524FXWAZ"), rlid.MustParse("064yrcthc000000d"))
	require.NoError(err, "could not retrieve event in the database")
	require.NotNil(event, "no event returned from the database")
	require.Equal(uint64(13), event.Offset)
	require.Len(event.Event, 67)

	// Should not be able to retreive an event from a null topic
	_, err = s.store.Retrieve(ulid.ULID{}, rlid.MustParse("064yrcthc000000d"))
	require.ErrorIs(err, errors.ErrInvalidKey)

	// Should not be able to retrieve an event with a null id
	_, err = s.store.Retrieve(ulid.MustParse("01GTSN1139JMK1PS5A524FXWAZ"), rlid.RLID{})
	require.ErrorIs(err, errors.ErrInvalidKey)
}

func (s *readonlyEventsTestSuite) TestRetrieve() {
	require := s.Require()
	require.True(s.store.ReadOnly())

	// Should be able to retrieve an event in the database
	event, err := s.store.Retrieve(ulid.MustParse("01GTSN1139JMK1PS5A524FXWAZ"), rlid.MustParse("064yrcthc000000d"))
	require.NoError(err, "could not retrieve event in the database")
	require.NotNil(event, "no event returned from the database")
	require.Equal(uint64(13), event.Offset)
	require.Len(event.Event, 67)

	// Should not be able to retreive an event from a null topic
	_, err = s.store.Retrieve(ulid.ULID{}, rlid.MustParse("064yrcthc000000d"))
	require.ErrorIs(err, errors.ErrInvalidKey)

	// Should not be able to retrieve an event with a null id
	_, err = s.store.Retrieve(ulid.MustParse("01GTSN1139JMK1PS5A524FXWAZ"), rlid.RLID{})
	require.ErrorIs(err, errors.ErrInvalidKey)
}
