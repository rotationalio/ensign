package events_test

import (
	"bytes"
	"encoding/base64"

	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg/ensign/rlid"
	"github.com/rotationalio/ensign/pkg/ensign/store/errors"
)

func (s *eventsTestSuite) TestIndash() {
	require := s.Require()
	require.False(s.store.ReadOnly())

	defer s.ResetDatabase()

	// Database should be empty to begin
	count, err := s.store.Count(nil)
	require.NoError(err, "could not count database")
	require.Equal(uint64(0), count, "expected no objects in the database")

	// Create an Indash
	topicID := ulid.MustParse("01GTSN1139JMK1PS5A524FXWAZ")
	hash, _ := base64.RawStdEncoding.DecodeString("skEWbaiEWvNu9CZKZNJuLg")
	eventID := rlid.MustParse("064yrcthc000000d")

	err = s.store.Indash(topicID, hash, eventID)
	require.NoError(err, "could not indash into the database")

	// Check to make sure the hash was inserted into the database
	count, err = s.store.Count(nil)
	require.NoError(err, "could not count database")
	require.Equal(uint64(1), count, "expected an indash in the database")
}

func (s *readonlyEventsTestSuite) TestIndash() {
	require := s.Require()
	require.True(s.store.ReadOnly())

	topicID := ulid.MustParse("01GTSN1139JMK1PS5A524FXWAZ")
	hash, _ := base64.RawStdEncoding.DecodeString("skEWbaiEWvNu9CZKZNJuLg")
	eventID := rlid.MustParse("064yrcthc000000d")

	err := s.store.Indash(topicID, hash, eventID)
	require.ErrorIs(err, errors.ErrReadOnly, "expected readonly error on indash")
}

func (s *eventsTestSuite) TestLoadIndash() {
	require := s.Require()
	require.False(s.store.ReadOnly())

	_, err := s.LoadAllFixtures()
	require.NoError(err, "could not load fixtures")
	defer s.ResetDatabase()

	topicID := ulid.MustParse("01GTSN1139JMK1PS5A524FXWAZ")
	iter := s.store.LoadIndash(topicID)
	defer iter.Release()

	hashes := make(map[string]struct{})
	for iter.Next() {
		// Make sure each hash is unique
		hash, err := iter.Hash()
		require.NoError(err, "could not retrieve hash from object")
		hashes[base64.RawStdEncoding.EncodeToString(hash)] = struct{}{}

		// Ensure hash is a suffix of the key
		key := iter.Key()
		require.True(bytes.HasSuffix(key, hash))

		// Ensure topic ID is a prefix of the key
		require.True(bytes.HasPrefix(key, topicID[:]))

		// Ensure the key is composed correctly
		require.Len(key, len(topicID[:])+2+len(hash))
	}

	require.NoError(iter.Error(), "expected no error iterating")
	require.Len(hashes, 22, "expected the same number of hashes as events in the fixture")
}

func (s *readonlyEventsTestSuite) TestLoadIndash() {
	require := s.Require()
	require.True(s.store.ReadOnly())

	topicID := ulid.MustParse("01GTSN1139JMK1PS5A524FXWAZ")
	iter := s.store.LoadIndash(topicID)
	defer iter.Release()

	hashes := make(map[string]struct{})
	for iter.Next() {
		// Make sure each hash is unique
		hash, err := iter.Hash()
		require.NoError(err, "could not retrieve hash from object")
		hashes[base64.RawStdEncoding.EncodeToString(hash)] = struct{}{}

		// Ensure hash is a suffix of the key
		key := iter.Key()
		require.True(bytes.HasSuffix(key, hash))

		// Ensure topic ID is a prefix of the key
		require.True(bytes.HasPrefix(key, topicID[:]))

		// Ensure the key is composed correctly
		require.Len(key, len(topicID[:])+2+len(hash))
	}

	require.NoError(iter.Error(), "expected no error iterating")
	require.Len(hashes, 22, "expected the same number of hashes as events in the fixture")
}
