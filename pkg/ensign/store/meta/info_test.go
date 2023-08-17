package meta_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/oklog/ulid/v2"
	api "github.com/rotationalio/ensign/pkg/ensign/api/v1beta1"
	"github.com/rotationalio/ensign/pkg/ensign/store"
	"github.com/rotationalio/ensign/pkg/ensign/store/errors"
	"github.com/rotationalio/ensign/pkg/ensign/store/meta"
	"github.com/rotationalio/ensign/pkg/utils/ulids"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *metaTestSuite) TestTopicInfo() {
	require := s.Require()
	require.False(s.store.ReadOnly())

	_, err := s.LoadAllFixtures()
	require.NoError(err, "could not load all fixtures")
	defer s.ResetDatabase()

	testTopicInfo(require, s.store)
}

func (s *readonlyMetaTestSuite) TestTopicInfo() {
	require := s.Require()
	require.True(s.store.ReadOnly())
	testTopicInfo(require, s.store)
}

func testTopicInfo(require *require.Assertions, store store.TopicInfoStore) {
	// Should be able to fetch a topic info that exists in the database
	topicID := ulid.MustParse("01GTSN1139JMK1PS5A524FXWAZ")
	info, err := store.TopicInfo(topicID)
	require.NoError(err, "expected topic with specified ID to exist")
	require.Equal(topicID[:], info.TopicId)
	require.Equal(ulids.MustBytes("01GTSMMC152Q95RD4TNYDFJGHT"), info.ProjectId)

	// Should get an empty topic info if the topic exists but the info does not
	topicID = ulid.MustParse("01GTSN2NQV61P2R4WFYF1NF1JG")
	info, err = store.TopicInfo(topicID)
	require.NoError(err, "expected topic info to be created")
	require.Equal(topicID[:], info.TopicId)
	require.Equal(ulids.MustBytes("01GTSMZNRYXNAZQF5R8NHQ14NM"), info.ProjectId)
	require.Zero(info.Events)
	require.Zero(info.Events)
	require.Zero(info.DataSizeBytes)
	require.Zero(info.Types)
	require.Zero(info.Modified)

	// Should get not found if the topic does not exist in the database
	_, err = store.TopicInfo(ulid.MustParse("01H7V5R4EZ4NATD6DC5RXWJMBG"))
	require.ErrorIs(err, errors.ErrNotFound)
}

func (s *metaTestSuite) TestUpdateTopicInfo() {
	require := s.Require()
	require.False(s.store.ReadOnly())

	info := &api.TopicInfo{
		ProjectId:     ulids.MustBytes("01H7V2HDHM6QH6CZ0KATPSQMF1"),
		TopicId:       ulids.MustBytes("01H7V2HMSR47TQVSFCNTD4D5EE"),
		Events:        uint64(192),
		Duplicates:    uint64(3),
		DataSizeBytes: uint64(5.977e+6),
		Types:         []*api.EventTypeInfo{},
		Modified:      timestamppb.New(time.Date(1999, 9, 9, 9, 9, 9, 0, time.UTC)),
	}

	count, err := s.store.Count(nil)
	require.NoError(err, "could not count db")
	require.Equal(uint64(0), count, "expected nothing in the database")

	err = s.store.UpdateTopicInfo(info)
	require.NoError(err, "expected to be able to update topic info")

	count, err = s.store.Count(nil)
	require.NoError(err, "could not count db")
	require.Equal(uint64(1), count, "expected the info in the database")

	info.Events += 14
	info.Duplicates += 2
	info.DataSizeBytes += 14021

	err = s.store.UpdateTopicInfo(info)
	require.NoError(err, "expected to be able to update topic info")

	count, err = s.store.Count(nil)
	require.NoError(err, "could not count db")
	require.Equal(uint64(1), count, "expected the info in the database")
}

func (s *readonlyMetaTestSuite) TestUpdateTopicInfo() {
	require := s.Require()
	require.True(s.store.ReadOnly())

	info := &api.TopicInfo{
		ProjectId: ulids.MustBytes("01H7V2HDHM6QH6CZ0KATPSQMF1"),
		TopicId:   ulids.MustBytes("01H7V2HMSR47TQVSFCNTD4D5EE"),
		Events:    uint64(192),
	}

	err := s.store.UpdateTopicInfo(info)
	require.ErrorIs(err, errors.ErrReadOnly)
}

func TestValidateTopicInfo(t *testing.T) {
	projectID := ulids.MustBytes("01H7V2HDHM6QH6CZ0KATPSQMF1")
	topicID := ulids.MustBytes("01H7V2HMSR47TQVSFCNTD4D5EE")

	testCases := []struct {
		info   *api.TopicInfo
		target error
	}{
		{nil, errors.ErrTopicInfoInvalidTopicId},
		{&api.TopicInfo{ProjectId: projectID}, errors.ErrTopicInfoMissingTopicId},
		{&api.TopicInfo{TopicId: topicID}, errors.ErrTopicInfoMissingProjectId},
		{&api.TopicInfo{TopicId: topicID, ProjectId: projectID[4:]}, errors.ErrTopicInfoInvalidProjectId},
		{&api.TopicInfo{TopicId: topicID, ProjectId: ulids.Null[:]}, errors.ErrTopicInfoInvalidProjectId},
		{&api.TopicInfo{ProjectId: projectID, TopicId: ulids.Null[:]}, errors.ErrTopicInfoInvalidTopicId},
		{&api.TopicInfo{ProjectId: projectID, TopicId: topicID[7:]}, errors.ErrTopicInfoInvalidTopicId},
		{&api.TopicInfo{ProjectId: projectID, TopicId: topicID}, nil},
	}

	for i, tc := range testCases {
		err := meta.ValidateTopicInfo(tc.info)
		require.ErrorIs(t, err, tc.target, "test %d failed", i)
	}

}

func TestTopicInfoKey(t *testing.T) {
	info := &api.TopicInfo{
		ProjectId: ulids.MustBytes("01H7V2HDHM6QH6CZ0KATPSQMF1"),
		TopicId:   ulids.MustBytes("01H7V2HMSR47TQVSFCNTD4D5EE"),
	}

	key := meta.TopicInfoKey(info)
	require.Len(t, key, 34)
	require.True(t, bytes.HasPrefix(key[:], info.ProjectId))
	require.True(t, bytes.Equal(key[16:18], meta.TopicInfoSegment[:]))
	require.True(t, bytes.HasSuffix(key[:], info.TopicId))
}
