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

	// List all events in a single topic
	nEvents := 0
	for iter.Next() {
		nEvents++

		// Ensure that the event can be deserialized
		event, err := iter.Event()
		require.NoError(err, "could not extract event from iterator")
		require.Equal(topicID.Bytes(), event.TopicId)
	}

	require.NoError(iter.Error(), "expected no error iterating")
	require.Equal(22, nEvents, "expected the number of events in the fixture")
}

func (s *readonlyEventsTestSuite) TestList() {
	require := s.Require()
	require.True(s.store.ReadOnly())

	topicID := ulid.MustParse("01GTSN1139JMK1PS5A524FXWAZ")
	iter := s.store.List(topicID)
	defer iter.Release()

	// List all events in a single topic
	nEvents := 0
	for iter.Next() {
		nEvents++

		// Ensure that the event can be deserialized
		event, err := iter.Event()
		require.NoError(err, "could not extract event from iterator")
		require.Equal(topicID.Bytes(), event.TopicId)
	}

	require.NoError(iter.Error(), "expected no error iterating")
	require.Equal(22, nEvents, "expected the number of events in the fixture")
}

func (s *eventsTestSuite) TestListSeek() {
	require := s.Require()
	require.False(s.store.ReadOnly())

	_, err := s.LoadAllFixtures()
	require.NoError(err, "could not load fixtures")
	defer s.ResetDatabase()

	topicID := ulid.MustParse("01GTSN1139JMK1PS5A524FXWAZ")
	iter := s.store.List(topicID)
	defer iter.Release()

	ok := iter.Seek(rlid.MustParse("064yrcthc000000d"))
	require.True(ok, "could not seek to event")

	// The iterator should currently be at the seek event without calling Next()
	event, err := iter.Event()
	require.NoError(err, "could not get seek to event")
	require.Equal(rlid.MustParse("064yrcthc000000d").Bytes(), event.Id)

	// List all remaining events in the topic
	nEvents := 0
	for iter.Next() {
		// Ensure that the event can be deserialized
		event, err := iter.Event()
		require.NoError(err, "could not extract event from iterator")
		require.Equal(topicID.Bytes(), event.TopicId)

		nEvents++
	}

	require.NoError(iter.Error(), "expected no error iterating")
	require.Equal(9, nEvents, "expected the number of events in the fixture")
}

func (s *readonlyEventsTestSuite) TestListSeek() {
	require := s.Require()
	require.True(s.store.ReadOnly())

	topicID := ulid.MustParse("01GTSN1139JMK1PS5A524FXWAZ")
	iter := s.store.List(topicID)
	defer iter.Release()

	ok := iter.Seek(rlid.MustParse("064yrcthc000000d"))
	require.True(ok, "could not seek to event")

	// The iterator should currently be at the seek event without calling Next()
	event, err := iter.Event()
	require.NoError(err, "could not get seek to event")
	require.Equal(rlid.MustParse("064yrcthc000000d").Bytes(), event.Id)

	// List all remaining events in the topic
	nEvents := 0
	for iter.Next() {
		// Ensure that the event can be deserialized
		event, err := iter.Event()
		require.NoError(err, "could not extract event from iterator")
		require.Equal(topicID.Bytes(), event.TopicId)

		nEvents++
	}

	require.NoError(iter.Error(), "expected no error iterating")
	require.Equal(9, nEvents, "expected the number of events in the fixture")
}

func (s *eventsTestSuite) TestListError() {
	require := s.Require()
	require.False(s.store.ReadOnly())

	iter := s.store.List(ulid.ULID{})
	defer iter.Release()
	require.ErrorIs(iter.Error(), errors.ErrKeyNull)
	require.False(iter.Seek(rlid.RLID{}))

	_, err := iter.Event()
	require.ErrorIs(err, errors.ErrKeyNull)
}

func (s *readonlyEventsTestSuite) TestListError() {
	require := s.Require()
	require.True(s.store.ReadOnly())

	iter := s.store.List(ulid.ULID{})
	defer iter.Release()
	require.ErrorIs(iter.Error(), errors.ErrKeyNull)
	require.False(iter.Seek(rlid.RLID{}))

	_, err := iter.Event()
	require.ErrorIs(err, errors.ErrKeyNull)
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

func (s *eventsTestSuite) TestDestroy() {
	require := s.Require()
	require.False(s.store.ReadOnly())

	_, err := s.LoadAllFixtures()
	require.NoError(err, "could not load fixtures")
	defer s.ResetDatabase()

	// Database should be empty to begin
	count, err := s.store.Count(nil)
	require.NoError(err, "could not count database")
	require.Equal(uint64(0xee), count, "expected no objects in the database")

	topicID := ulid.MustParse("01GTSN1139JMK1PS5A524FXWAZ")
	err = s.store.Destroy(topicID)
	require.NoError(err, "unable to destroy topic")

	// Check to make sure all objects were destroyed
	count, err = s.store.Count(nil)
	require.NoError(err, "could not count database")
	require.Equal(uint64(0xc2), count, "expected an event inserted into the database")

	// There should be no events in the database
	nEvents := 0
	events := s.store.List(topicID)
	defer events.Release()
	for events.Next() {
		nEvents++
	}
	require.NoError(events.Error(), "could not iterate over events")
	require.Zero(nEvents, "expected no events in the database")

	// There should be no index hashes in the database
	nIndash := 0
	hashes := s.store.LoadIndash(topicID)
	defer hashes.Release()
	for hashes.Next() {
		nIndash++
	}
	require.NoError(hashes.Error(), "could not iterate over hashes")
	require.Zero(nIndash, "expected no index hashes in the database")
}

func (s *readonlyEventsTestSuite) TestDestroy() {
	require := s.Require()
	require.True(s.store.ReadOnly())

	topicID := ulid.MustParse("01GTSN1139JMK1PS5A524FXWAZ")
	err := s.store.Destroy(topicID)
	require.ErrorIs(err, errors.ErrReadOnly, "expected readonly error on destroy topic")
}
