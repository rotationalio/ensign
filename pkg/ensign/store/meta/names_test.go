package meta_test

import (
	"bytes"
	"testing"

	"github.com/oklog/ulid/v2"
	api "github.com/rotationalio/ensign/pkg/ensign/api/v1beta1"
	"github.com/rotationalio/ensign/pkg/ensign/store"
	"github.com/rotationalio/ensign/pkg/ensign/store/errors"
	"github.com/rotationalio/ensign/pkg/ensign/store/meta"
	"github.com/rotationalio/ensign/pkg/utils/ulids"
	"github.com/stretchr/testify/require"
)

func (s *metaTestSuite) TestListTopicNames() {
	require := s.Require()
	require.False(s.store.ReadOnly())

	_, err := s.LoadAllFixtures()
	require.NoError(err, "could not load all fixtures")
	defer s.ResetDatabase()

	testListTopicNames(require, s.store)
}

func (s *readonlyMetaTestSuite) TestListTopicNames() {
	require := s.Require()
	require.True(s.store.ReadOnly())
	testListTopicNames(require, s.store)
}

func testListTopicNames(require *require.Assertions, store store.TopicNamesStore) {
	topics := store.ListTopicNames(ulids.MustParse("01GTSMMC152Q95RD4TNYDFJGHT"))
	defer topics.Release()

	nTopics := 0
	for topics.Next() {
		nTopics++
		topic, err := topics.TopicName()
		require.NoError(err, "could not deserialize topic name")
		require.Equal("01GTSMMC152Q95RD4TNYDFJGHT", topic.ProjectId)
		require.NotEmpty(topic.Name, "missing topic name hash")
		require.NotEmpty(topic.TopicId, "missing topic id")
	}
	require.Equal(5, nTopics)

	err := topics.Error()
	require.NoError(err, "could not list topics from database")
}

func (s *metaTestSuite) TestListTopicNamesPagination() {
	require := s.Require()
	require.False(s.store.ReadOnly())

	_, err := s.LoadAllFixtures()
	require.NoError(err, "could not load all fixtures")
	defer s.ResetDatabase()

	testListTopicNamesPagination(require, s.store)
}

func (s *readonlyMetaTestSuite) TestListTopicNamesPagination() {
	require := s.Require()
	require.True(s.store.ReadOnly())
	testListTopicNamesPagination(require, s.store)
}

func testListTopicNamesPagination(require *require.Assertions, store store.TopicNamesStore) {
	topics := store.ListTopicNames(ulids.MustParse("01GTSMMC152Q95RD4TNYDFJGHT"))
	defer topics.Release()

	pages := 0
	items := 0
	info := &api.PageInfo{PageSize: uint32(2)}

	// Only paginate for a maximum of 10 iterations
	for i := 0; i < 10; i++ {
		page, err := topics.NextPage(info)
		require.NoError(err, "could not fetch page %d", i+1)
		require.LessOrEqual(len(page.TopicNames), int(info.PageSize))

		pages++
		items += len(page.TopicNames)

		if page.NextPageToken == "" {
			break
		}

		info.NextPageToken = page.NextPageToken
	}

	require.NoError(topics.Error(), "could not list topic names from database")
	require.Equal(3, pages)
	require.Equal(5, items)
}

func (s *metaTestSuite) TestTopicExists() {
	require := s.Require()
	require.False(s.store.ReadOnly())

	_, err := s.LoadAllFixtures()
	require.NoError(err, "could not load all fixtures")
	defer s.ResetDatabase()

	testTopicExists(require, s.store)
}

func (s *readonlyMetaTestSuite) TestTopicExists() {
	require := s.Require()
	require.True(s.store.ReadOnly())
	testTopicExists(require, s.store)
}

func testTopicExists(require *require.Assertions, store store.TopicNamesStore) {
	testCases := []struct {
		in  *api.TopicName
		out *api.TopicExistsInfo
		err error
	}{
		{
			&api.TopicName{
				ProjectId: "01GTSMMC152Q95RD4TNYDFJGHT",
				TopicId:   "01GTSMQ3V8ASAPNCFEN378T8RD",
			},
			&api.TopicExistsInfo{
				Query:  `topic="01GTSMQ3V8ASAPNCFEN378T8RD"`,
				Exists: true,
			},
			nil,
		},
		{
			&api.TopicName{
				ProjectId: "01GTSMMC152Q95RD4TNYDFJGHT",
				TopicId:   "01GW2MMS27J0Q9BVE4R156G2MD",
			},
			&api.TopicExistsInfo{
				Query:  `topic="01GW2MMS27J0Q9BVE4R156G2MD"`,
				Exists: false,
			},
			nil,
		},
		{
			&api.TopicName{
				ProjectId: "01GTSMMC152Q95RD4TNYDFJGHT",
				Name:      "testing.testapp.alerts",
			},
			&api.TopicExistsInfo{
				Query:  `name="testing.testapp.alerts"`,
				Exists: true,
			},
			nil,
		},
		{
			&api.TopicName{
				ProjectId: "01GTSMMC152Q95RD4TNYDFJGHT",
				Name:      "testing.testapp.notalerts",
			},
			&api.TopicExistsInfo{
				Query:  `name="testing.testapp.notalerts"`,
				Exists: false,
			},
			nil,
		},
		{
			&api.TopicName{
				ProjectId: "01GTSMMC152Q95RD4TNYDFJGHT",
				TopicId:   "01GTSMQ3V8ASAPNCFEN378T8RD",
				Name:      "testing.testapp.alerts",
			},
			&api.TopicExistsInfo{
				Query:  `name="testing.testapp.alerts" and topic="01GTSMQ3V8ASAPNCFEN378T8RD"`,
				Exists: true,
			},
			nil,
		},
		{
			&api.TopicName{
				ProjectId: "01GTSMMC152Q95RD4TNYDFJGHT",
				TopicId:   "01GW2MMS27J0Q9BVE4R156G2MD",
				Name:      "testing.testapp.alerts",
			},
			&api.TopicExistsInfo{
				Query:  `name="testing.testapp.alerts" and topic="01GW2MMS27J0Q9BVE4R156G2MD"`,
				Exists: false,
			},
			nil,
		},
		{
			&api.TopicName{
				ProjectId: "01GTSMMC152Q95RD4TNYDFJGHT",
				TopicId:   "01GTSMQ3V8ASAPNCFEN378T8RD",
				Name:      "testing.testapp.notalerts",
			},
			&api.TopicExistsInfo{
				Query:  `name="testing.testapp.notalerts" and topic="01GTSMQ3V8ASAPNCFEN378T8RD"`,
				Exists: false,
			},
			nil,
		},
		{
			&api.TopicName{
				ProjectId: "",
				TopicId:   "01GTSMQ3V8ASAPNCFEN378T8RD",
			},
			nil,
			errors.ErrTopicInvalidProjectId,
		},
		{
			&api.TopicName{
				ProjectId: "notanulid",
				TopicId:   "01GTSMQ3V8ASAPNCFEN378T8RD",
			},
			nil,
			errors.ErrTopicInvalidProjectId,
		},
		{
			&api.TopicName{
				ProjectId: "01GTSMMC152Q95RD4TNYDFJGHT",
			},
			nil,
			errors.ErrTopicMissingName,
		},
	}

	for i, tc := range testCases {
		info, err := store.TopicExists(tc.in)

		if tc.err == nil {
			require.NoError(err, "an unexpected error was returned on test case %d", i)
			require.Equal(tc.out.Query, info.Query, "query comparison failed on test case %d", i)
			require.Equal(tc.out.Exists, info.Exists, "exists comparison failed on test case %d", i)
		} else {
			require.ErrorIs(err, tc.err, "error is failed on test case %d", i)
			require.Nil(info, "info was not nil on test case %d", i)
		}
	}
}

func (s *metaTestSuite) TestTopicName() {
	require := s.Require()
	require.False(s.store.ReadOnly())

	_, err := s.LoadAllFixtures()
	require.NoError(err, "could not load all fixtures")
	defer s.ResetDatabase()

	testTopicName(require, s.store)
}

func (s *readonlyMetaTestSuite) TestTopicName() {
	require := s.Require()
	require.True(s.store.ReadOnly())
	testTopicName(require, s.store)
}

func testTopicName(require *require.Assertions, store store.TopicNamesStore) {
	testCases := []struct {
		topicID  ulid.ULID
		expected string
		err      error
	}{
		{ulids.Null, "", errors.ErrNotFound},
		{ulid.MustParse("01H7D9XQ6FDNKSN0B6M070E0TV"), "", errors.ErrNotFound},
		{ulid.MustParse("01GTSMQ3V8ASAPNCFEN378T8RD"), "testing.testapp.alerts", nil},
		{ulid.MustParse("01GTSN1WF5BA0XCPT6ES64JVGQ"), "mock.mockapp.feed", nil},
	}

	for i, tc := range testCases {
		actual, err := store.TopicName(tc.topicID)
		if tc.err != nil {
			require.Error(err, "expected error for test case %d", i)
			require.ErrorIs(err, tc.err, "expected error for test case %d", i)
		} else {
			require.NoError(err, "expected no error for test case %d", i)
			require.Equal(tc.expected, actual, "expected topic name match for test case %d", i)
		}
	}
}

func (s *metaTestSuite) TestLookupTopicName() {
	require := s.Require()
	require.False(s.store.ReadOnly())

	_, err := s.LoadAllFixtures()
	require.NoError(err, "could not load all fixtures")
	defer s.ResetDatabase()

	testLookupTopicName(require, s.store)
}

func (s *readonlyMetaTestSuite) TestLookupTopicName() {
	require := s.Require()
	require.True(s.store.ReadOnly())
	testLookupTopicName(require, s.store)
}

func testLookupTopicName(require *require.Assertions, store store.TopicNamesStore) {
	testCases := []struct {
		name      string
		projectID ulid.ULID
		expected  ulid.ULID
		err       error
	}{
		{"", ulid.MustParse("01GTSMMC152Q95RD4TNYDFJGHT"), ulids.Null, errors.ErrNotFound},
		{"testing.testapp.receipts", ulids.Null, ulids.Null, errors.ErrNotFound},
		{"testing.testapp.receipts", ulid.MustParse("01H7D9XQ6FDNKSN0B6M070E0TV"), ulids.Null, errors.ErrNotFound},
		{"testing.testapp.receipts", ulid.MustParse("01GTSMZNRYXNAZQF5R8NHQ14NM"), ulids.Null, errors.ErrNotFound},
		{"banana-fruit-wacky", ulid.MustParse("01GTSMMC152Q95RD4TNYDFJGHT"), ulids.Null, errors.ErrNotFound},
		{"testing.testapp.receipts", ulid.MustParse("01GTSMMC152Q95RD4TNYDFJGHT"), ulid.MustParse("01GV6KYPW33RW5D800ERR3NP8S"), nil},
		{"mock.mockapp.post", ulid.MustParse("01GTSMZNRYXNAZQF5R8NHQ14NM"), ulid.MustParse("01GTSN2NQV61P2R4WFYF1NF1JG"), nil},
	}

	for i, tc := range testCases {
		actual, err := store.LookupTopicName(tc.name, tc.projectID)
		if tc.err != nil {
			require.Error(err, "expected error for test case %d", i)
			require.ErrorIs(err, tc.err, "expected error for test case %d", i)
		} else {
			require.NoError(err, "expected no error for test case %d", i)
			require.Equal(tc.expected, actual, "expected topic name match for test case %d", i)
		}
	}
}

func TestTopicNameKey(t *testing.T) {
	topic := &api.Topic{
		ProjectId: ulids.MustBytes("01GTSSDM957VH0GX0RMNKAQM13"),
		Name:      "testing.testapp.foo",
	}

	key := meta.TopicNameKey(topic)

	require.Len(t, key, 34, "expected the key length to be two ulids long")
	require.True(t, bytes.HasPrefix(key[:], topic.ProjectId))
	require.True(t, bytes.Equal(key[16:18], meta.TopicNamesSegment[:]))
	require.True(t, bytes.HasSuffix(key[:], []byte{23, 140, 224, 69, 25, 219, 130, 55, 226, 167, 181, 227, 210, 179, 70, 34}))
}
