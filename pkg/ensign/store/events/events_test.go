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
