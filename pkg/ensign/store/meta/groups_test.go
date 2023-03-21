package meta_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	api "github.com/rotationalio/ensign/pkg/ensign/api/v1beta1"
	"github.com/rotationalio/ensign/pkg/ensign/store/errors"
	"github.com/rotationalio/ensign/pkg/ensign/store/meta"
	"github.com/rotationalio/ensign/pkg/utils/ulids"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *metaTestSuite) TestSameGroupName() {
	// This is an essential safety test that exercises most of the group store.
	require := s.Require()
	require.False(s.store.ReadOnly())
	defer s.ResetDatabase()

	projectA := ulids.New().Bytes()
	projectB := ulids.New().Bytes()

	// Two different projects should be able to store a group with the same name without
	// conflict with or without group IDs specified by the user.
	testCases := []struct {
		group *api.ConsumerGroup
		msg   string
	}{
		{
			&api.ConsumerGroup{Id: []byte{1, 134, 179, 108, 62, 211, 134, 53, 49, 102, 31, 33, 40, 215, 58, 245}},
			"same 16 byte group id",
		},
		{
			&api.ConsumerGroup{Name: "unicorn"},
			"same name",
		},
		{
			&api.ConsumerGroup{Id: []byte{1, 134, 179, 108, 62, 211, 134, 53, 49, 102, 31, 33, 40, 215, 58, 245, 43, 193, 203, 6}},
			"same 20 byte group id",
		},
	}

	// Expected empty database for this test
	objs, err := s.store.Count(nil)
	require.NoError(err, "could not count db")
	require.Equal(uint64(0), objs, "unexpected empty database to start test")

	for _, tc := range testCases {
		// Create two groups, A and B with the same name and ID as specified in the
		// test case but with different project IDs and data.
		groupA := &api.ConsumerGroup{
			Id:              tc.group.Id,
			ProjectId:       projectA,
			Name:            tc.group.Name,
			Delivery:        api.DeliverySemantic_AT_LEAST_ONCE,
			DeliveryTimeout: durationpb.New(30 * time.Second),
			TopicOffsets:    map[string]uint64{"01GTSMQ3V8ASAPNCFEN378T8RD": 83123},
			Consumers:       [][]byte{{1, 2, 3, 4}},
		}

		groupB := &api.ConsumerGroup{
			Id:              tc.group.Id,
			ProjectId:       projectB,
			Name:            tc.group.Name,
			Delivery:        api.DeliverySemantic_EXACTLY_ONCE,
			DeliveryTimeout: durationpb.New(10 * time.Second),
			TopicOffsets:    map[string]uint64{"01GTSN1WF5BA0XCPT6ES64JVGQ": 62},
			Consumers:       [][]byte{{4, 3, 2, 1}},
		}

		// Should be able to independently create the groups
		created, err := s.store.GetOrCreateGroup(groupA)
		require.NoError(err, "could not create group a: %s", tc.msg)
		require.True(created, "group a was not created: %s", tc.msg)

		created, err = s.store.GetOrCreateGroup(groupB)
		require.NoError(err, "could not create group b: %s", tc.msg)
		require.True(created, "group b was not created: %s", tc.msg)

		// Should be two items in the database
		objs, err := s.store.Count(nil)
		require.NoError(err, "could not count db")
		require.Equal(uint64(2), objs, "unexpected number of objects in database")

		// Should be able to independently update the groups
		groupA.TopicOffsets["01GTSN1139JMK1PS5A524FXWAZ"] = 201
		groupB.TopicOffsets["01GTSN1WF5BA0XCPT6ES64JVGQ"] = 102

		err = s.store.UpdateGroup(groupA)
		require.NoError(err, "could not update group a: %s", tc.msg)

		err = s.store.UpdateGroup(groupB)
		require.NoError(err, "could not update group b: %s", tc.msg)

		// Should be able to retrieve comparable groups
		groupAret := &api.ConsumerGroup{Id: groupA.Id, ProjectId: groupA.ProjectId, Name: groupA.Name}
		groupBret := &api.ConsumerGroup{Id: groupB.Id, ProjectId: groupB.ProjectId, Name: groupB.Name}

		created, err = s.store.GetOrCreateGroup(groupAret)
		require.NoError(err, "could not retrieve group a: %s", tc.msg)
		require.False(created, "group a was created: %s", tc.msg)

		created, err = s.store.GetOrCreateGroup(groupBret)
		require.NoError(err, "could not retrieve group b: %s", tc.msg)
		require.False(created, "group b was created: %s", tc.msg)

		// Compare and contrast retrieved with originals
		require.True(proto.Equal(groupA, groupAret))
		require.True(proto.Equal(groupB, groupBret))
		require.False(proto.Equal(groupB, groupAret))
		require.False(proto.Equal(groupBret, groupAret))
		require.False(proto.Equal(groupA, groupBret))
		require.False(proto.Equal(groupAret, groupBret))

		// Should be able to delete items from database independently
		err = s.store.DeleteGroup(groupA)
		require.NoError(err, "could not delete group a: %s", tc.msg)

		objs, err = s.store.Count(nil)
		require.NoError(err, "could not count db")
		require.Equal(uint64(1), objs, "unexpected number of objects in database")

		err = s.store.DeleteGroup(groupB)
		require.NoError(err, "could not delete group b: %s", tc.msg)
	}
}

func (s *metaTestSuite) TestListGroups() {
	require := s.Require()
	require.False(s.store.ReadOnly())

	_, err := s.LoadAllFixtures()
	require.NoError(err, "could not load all fixtures")
	defer s.ResetDatabase()

	groups := s.store.ListGroups(ulids.MustParse("01GTSMMC152Q95RD4TNYDFJGHT"))
	defer groups.Release()

	nGroups := 0
	for groups.Next() {
		nGroups++
		group, err := groups.Group()
		require.NoError(err, "could not deserialize group")
		require.True(strings.HasPrefix(group.Name, "testing.group"))
	}
	require.Equal(10, nGroups)

	err = groups.Error()
	require.NoError(err, "could not list groups from database")
}

func (s *readonlyMetaTestSuite) TestListGroups() {
	require := s.Require()
	require.True(s.store.ReadOnly())

	groups := s.store.ListGroups(ulids.MustParse("01GTSMMC152Q95RD4TNYDFJGHT"))
	defer groups.Release()

	nGroups := 0
	for groups.Next() {
		nGroups++
		group, err := groups.Group()
		require.NoError(err, "could not deserialize group")
		require.True(strings.HasPrefix(group.Name, "testing.group"))
	}
	require.Equal(10, nGroups)

	err := groups.Error()
	require.NoError(err, "could not list groups from database")
}

func (s *metaTestSuite) TestGetOrCreateGroup() {
	require := s.Require()
	require.False(s.store.ReadOnly())

	defer s.ResetDatabase()

	// Database should be empty to begin
	count, err := s.store.Count(nil)
	require.NoError(err, "could not count database")
	require.Equal(uint64(0), count, "expected no objects in the database")

	// Should not be able to create an empty group
	_, err = s.store.GetOrCreateGroup(&api.ConsumerGroup{})
	require.ErrorIs(err, errors.ErrInvalidGroup)

	// Should be able to create a valid group
	group := &api.ConsumerGroup{
		ProjectId: ulids.MustBytes("01GTSRBV1HRZ3PPETSM3YF1N79"),
		Name:      "testing.groups.group1",
	}

	created, err := s.store.GetOrCreateGroup(group)
	require.NoError(err, "could not create valid group")
	require.True(created, "group was not created")

	// DB should have one object in it
	count, err = s.store.Count(nil)
	require.NoError(err, "could not count database")
	require.Equal(uint64(1), count, "expected 1 group in the database")

	// Second call should retrieve the group
	groupb := &api.ConsumerGroup{ProjectId: group.ProjectId, Name: "testing.groups.group1"}
	created, err = s.store.GetOrCreateGroup(groupb)
	require.NoError(err, "could not create valid group")
	require.False(created, "group was not created")

	require.True(proto.Equal(group, groupb))

	// DB should have one object in it
	count, err = s.store.Count(nil)
	require.NoError(err, "could not count database")
	require.Equal(uint64(1), count, "expected 1 group in the database")
}

func (s *readonlyMetaTestSuite) TestGetOrCreateGroup() {
	require := s.Require()
	require.True(s.store.ReadOnly())

	// Should not be able to create a new group
	group := &api.ConsumerGroup{
		ProjectId: ulids.New().Bytes(),
		Name:      "newgroup",
	}

	created, err := s.store.GetOrCreateGroup(group)
	require.ErrorIs(err, errors.ErrReadOnly)
	require.False(created)

	// Should be able to get a group that has already been created
	group = &api.ConsumerGroup{
		ProjectId: ulids.MustBytes("01GTSMZNRYXNAZQF5R8NHQ14NM"),
		Id:        ulids.MustBytes("01GVP6XTNT1FM1XWA2Q4Q0VBKQ"),
	}

	created, err = s.store.GetOrCreateGroup(group)
	require.NoError(err)
	require.False(created)

	require.Equal(api.DeliverySemantic_AT_MOST_ONCE, group.Delivery)
}

func (s *metaTestSuite) TestCreateGroup() {
	require := s.Require()
	require.False(s.store.ReadOnly())

	defer s.ResetDatabase()

	// Database should be empty to begin
	count, err := s.store.Count(nil)
	require.NoError(err, "could not count database")
	require.Equal(uint64(0), count, "expected no objects in the database")

	// Should not be able to create an empty group
	err = s.store.CreateGroup(&api.ConsumerGroup{})
	require.ErrorIs(err, errors.ErrInvalidGroup)

	// Should be able to create a valid group
	group := &api.ConsumerGroup{
		ProjectId: ulids.MustBytes("01GTSRBV1HRZ3PPETSM3YF1N79"),
		Name:      "testing.groups.group1",
	}

	err = s.store.CreateGroup(group)
	require.NoError(err, "could not create valid group")

	// DB should have one object in it
	count, err = s.store.Count(nil)
	require.NoError(err, "could not count database")
	require.Equal(uint64(1), count, "expected 1 group in the database")
}

func (s *readonlyMetaTestSuite) TestCreateGroup() {
	require := s.Require()
	require.True(s.store.ReadOnly())

	group := &api.ConsumerGroup{
		ProjectId: ulids.MustBytes("01GTSRBV1HRZ3PPETSM3YF1N79"),
		Name:      "testing.group.test",
	}

	err := s.store.CreateGroup(group)
	require.ErrorIs(err, errors.ErrReadOnly, "expected readonly error on create group")
}

func (s *metaTestSuite) TestRetrieveGroup() {
	require := s.Require()
	require.False(s.store.ReadOnly())

	_, err := s.LoadAllFixtures()
	require.NoError(err, "could not load topic fixtures")
	defer s.ResetDatabase()

	group := &api.ConsumerGroup{
		Id:        ulids.MustBytes("01GVP6XTNT1FM1XWA2Q4Q0VBKQ"),
		ProjectId: ulids.MustBytes("01GTSMZNRYXNAZQF5R8NHQ14NM"),
	}

	err = s.store.RetrieveGroup(group)
	require.NoError(err, "could not retrieve topic")
	require.Equal("feed-monitor", group.Name)
}

func (s *readonlyMetaTestSuite) TestRetrieveGroup() {
	require := s.Require()
	require.True(s.store.ReadOnly())

	group := &api.ConsumerGroup{
		Id:        ulids.MustBytes("01GVP6XTNT1FM1XWA2Q4Q0VBKQ"),
		ProjectId: ulids.MustBytes("01GTSMZNRYXNAZQF5R8NHQ14NM"),
	}

	err := s.store.RetrieveGroup(group)
	require.NoError(err, "could not retrieve topic")
	require.Equal("feed-monitor", group.Name)
}

func (s *metaTestSuite) TestUpdateGroup() {
	require := s.Require()
	require.False(s.store.ReadOnly())

	nFixtures, err := s.LoadGroupFixtures()
	require.NoError(err, "could not load topic fixtures")
	defer s.ResetDatabase()

	// Database should have the fixtures states to start
	count, err := s.store.Count(nil)
	require.NoError(err, "could not count database")
	require.Equal(nFixtures, count, "expected topic fixtures in the database")

	// Cannot update a group that doesn't exist
	group := &api.ConsumerGroup{
		Id:        ulids.MustBytes("01GTSRBV1HRZ3PPETSM3YF1N39"),
		ProjectId: ulids.MustBytes("01GTSRBV1HRZ3PPETSM3YF1N79"),
		Name:      "testing.group.test",
		Created:   timestamppb.Now(),
		Modified:  timestamppb.Now(),
	}

	err = s.store.UpdateGroup(group)
	require.ErrorIs(err, errors.ErrNotFound)

	// Can update a group that does exist
	group = &api.ConsumerGroup{
		Id:           ulids.MustBytes("01GVP6XTNT1FM1XWA2Q4Q0VBKQ"),
		ProjectId:    ulids.MustBytes("01GTSMZNRYXNAZQF5R8NHQ14NM"),
		Name:         "testing.group.test",
		Created:      timestamppb.Now(),
		Modified:     timestamppb.Now(),
		TopicOffsets: map[string]uint64{"01GTSN1WF5BA0XCPT6ES64JVGQ": 71},
	}

	err = s.store.UpdateGroup(group)
	require.NoError(err)

	// TODO: ensure group was actually updated

	// Database should have the same number of fixtures states to finish
	count, err = s.store.Count(nil)
	require.NoError(err, "could not count database")
	require.Equal(nFixtures, count, "expected no change in the count of objects")
}

func (s *readonlyMetaTestSuite) TestUpdateGroup() {
	require := s.Require()
	require.True(s.store.ReadOnly())

	group := &api.ConsumerGroup{
		ProjectId: ulids.MustBytes("01GTSRBV1HRZ3PPETSM3YF1N79"),
		Name:      "testing.group.test",
		Created:   timestamppb.Now(),
		Modified:  timestamppb.Now(),
	}

	var key [16]byte
	key, err := group.Key()
	require.NoError(err, "could not create group ID")
	group.Id = key[:]

	err = s.store.UpdateGroup(group)
	require.ErrorIs(err, errors.ErrReadOnly, "expected readonly error on update group")
}

func (s *metaTestSuite) TestDeleteGroup() {
	require := s.Require()
	require.False(s.store.ReadOnly())

	nFixtures, err := s.LoadGroupFixtures()
	require.NoError(err, "could not load topic fixtures")
	defer s.ResetDatabase()

	// Database should have the fixtures states to start
	count, err := s.store.Count(nil)
	require.NoError(err, "could not count database")
	require.Equal(nFixtures, count, "expected topic fixtures in the database")

	group := &api.ConsumerGroup{
		Id:        ulids.MustBytes("01GVP6XTNT1FM1XWA2Q4Q0VBKQ"),
		ProjectId: ulids.MustBytes("01GTSMZNRYXNAZQF5R8NHQ14NM"),
	}

	// Should be able to delete the group
	err = s.store.DeleteGroup(group)
	require.NoError(err, "could not delete group")

	// Group should have been deleted
	count, err = s.store.Count(nil)
	require.NoError(err, "could not count database")
	require.Equal(nFixtures-1, count, "expected one less group in the database")

	// Deleting a second time should have no effect
	err = s.store.DeleteGroup(group)
	require.NoError(err, "could not delete group")

	count, err = s.store.Count(nil)
	require.NoError(err, "could not count database")
	require.Equal(nFixtures-1, count, "expected no change in database count")
}

func (s *readonlyMetaTestSuite) TestDeleteGroup() {
	require := s.Require()
	require.True(s.store.ReadOnly())

	group := &api.ConsumerGroup{
		ProjectId: ulids.MustBytes("01GTSRBV1HRZ3PPETSM3YF1N79"),
		Name:      "testing.group.test",
	}

	err := s.store.UpdateGroup(group)
	require.ErrorIs(err, errors.ErrReadOnly, "expected readonly error on update group")
}

func TestGroupKey(t *testing.T) {
	empty16 := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	projectID := ulids.New().Bytes()
	uu := uuid.New()

	// Should be able to create a group key for a valid groups
	testCases := []struct {
		group  *api.ConsumerGroup
		suffix []byte
		msg    string
	}{
		{
			group:  &api.ConsumerGroup{ProjectId: projectID, Name: "testing.groups.group1"},
			suffix: []byte{114, 243, 14, 146, 211, 171, 102, 175, 73, 218, 24, 84, 252, 72, 68, 142},
			msg:    "group with project ID and name",
		},
		{
			group:  &api.ConsumerGroup{ProjectId: projectID, Id: uu[:]},
			suffix: uu[:],
			msg:    "group with project ID and UUID id",
		},
		{
			group:  &api.ConsumerGroup{ProjectId: projectID, Id: []byte{1, 2, 3, 4, 5, 6, 7, 8}},
			suffix: []byte{156, 232, 12, 165, 239, 147, 191, 220, 197, 103, 229, 230, 182, 85, 172, 7},
			msg:    "group with project ID and variable length ID",
		},
		{
			group:  &api.ConsumerGroup{ProjectId: projectID, Id: []byte{1, 2, 3, 4, 5, 6, 7, 8}, Name: "testing.groups.group1"},
			suffix: []byte{156, 232, 12, 165, 239, 147, 191, 220, 197, 103, 229, 230, 182, 85, 172, 7},
			msg:    "group with project ID and variable length ID and name specified",
		},
		{
			group:  &api.ConsumerGroup{ProjectId: projectID, Id: uu[:], Name: "testing.groups.group1"},
			suffix: uu[:],
			msg:    "group with project ID and uuid ID and name specified",
		},
		{
			group:  &api.ConsumerGroup{ProjectId: empty16, Id: empty16},
			suffix: empty16,
			msg:    "group with all zeros for project ID and group ID",
		},
	}

	for _, tc := range testCases {
		// Require test case to be valid
		err := meta.ValidateGroup(tc.group, true)
		require.NoError(t, err, tc.msg)

		key := meta.GroupKey(tc.group)
		require.Len(t, key, 34, "expected the key length to be two ulids long")
		require.True(t, bytes.HasPrefix(key[:], tc.group.ProjectId))
		require.True(t, bytes.Equal(key[16:18], meta.GroupSegment[:]))
		require.True(t, bytes.HasSuffix(key[:], tc.suffix), tc.msg)
	}
}

func TestInvalidGroupKey(t *testing.T) {
	// Should not be able to create a group key for invalid groups.
	testCases := []struct {
		group *api.ConsumerGroup
		msg   string
	}{
		{nil, "should not be able to create group key for nil group"},
		{&api.ConsumerGroup{}, "should not be able to create group key for empty group"},
		{&api.ConsumerGroup{ProjectId: []byte("foo")}, "should not be able to create group with invalid project id"},
		{&api.ConsumerGroup{ProjectId: ulids.New().Bytes()}, "should not be able to create group with only project id"},
		{&api.ConsumerGroup{ProjectId: []byte("foo"), Name: "foo"}, "should not be able to create group with invalid project id and name"},
		{&api.ConsumerGroup{Name: "foo"}, "should not be able to create group without a project id"},
	}

	for _, tc := range testCases {
		// Require test case to be invalid
		err := meta.ValidateGroup(tc.group, true)
		require.ErrorIs(t, err, errors.ErrInvalidGroup, tc.msg)

		require.Panics(t, func() {
			meta.GroupKey(tc.group)
		}, tc.msg)
	}
}

func TestValidateGroup(t *testing.T) {
	testCases := []struct {
		group   *api.ConsumerGroup
		partial bool
		err     error
	}{
		{
			nil,
			true,
			errors.ErrGroupMissingKeyField,
		},
		{
			nil,
			false,
			errors.ErrGroupMissingKeyField,
		},
		{
			&api.ConsumerGroup{},
			true,
			errors.ErrGroupMissingKeyField,
		},
		{
			&api.ConsumerGroup{},
			false,
			errors.ErrGroupMissingKeyField,
		},
		{
			&api.ConsumerGroup{
				Id: []byte{1, 2, 3, 4, 5, 6, 7, 8},
			},
			false,
			errors.ErrGroupMissingProjectId,
		},
		{
			&api.ConsumerGroup{
				Name: "testing.groups.group1",
			},
			false,
			errors.ErrGroupMissingProjectId,
		},
		{
			&api.ConsumerGroup{
				ProjectId: []byte{1, 2, 3, 4},
				Name:      "testing.groups.group1",
			},
			false,
			errors.ErrGroupInvalidProjectId,
		},
		{
			&api.ConsumerGroup{
				ProjectId: []byte{1, 134, 179, 81, 86, 251, 48, 108, 44, 19, 143, 243, 195, 87, 134, 80},
				Name:      "testing.groups.group1",
			},
			true,
			nil,
		},
		{
			&api.ConsumerGroup{
				ProjectId: []byte{1, 134, 179, 81, 86, 251, 48, 108, 44, 19, 143, 243, 195, 87, 134, 80},
				Id:        []byte{1, 134, 179, 108, 62, 211, 134, 53, 49, 102, 31, 33, 40, 215, 58, 245},
			},
			true,
			nil,
		},
		{
			&api.ConsumerGroup{
				ProjectId: []byte{1, 134, 179, 81, 86, 251, 48, 108, 44, 19, 143, 243, 195, 87, 134, 80},
				Id:        []byte{1, 134, 179, 108, 62, 211, 134, 53, 49, 102, 31, 33, 40, 215, 58, 245},
				Name:      "testing.groups.group1",
			},
			true,
			nil,
		},
		{
			&api.ConsumerGroup{
				ProjectId: []byte{1, 134, 179, 81, 86, 251, 48, 108, 44, 19, 143, 243, 195, 87, 134, 80},
				Id:        []byte{1, 2, 3, 4, 5, 6, 7, 8},
				Name:      "testing.groups.group1",
			},
			true,
			nil,
		},
		{
			&api.ConsumerGroup{
				ProjectId: []byte{1, 134, 179, 81, 86, 251, 48, 108, 44, 19, 143, 243, 195, 87, 134, 80},
				Name:      "testing.groups.group1",
				Created:   timestamppb.Now(),
				Modified:  timestamppb.Now(),
			},
			false,
			errors.ErrGroupMissingId,
		},
		{
			&api.ConsumerGroup{
				ProjectId: []byte{1, 134, 179, 81, 86, 251, 48, 108, 44, 19, 143, 243, 195, 87, 134, 80},
				Id:        []byte{1, 134, 179, 108, 62, 211, 134, 53, 49, 102, 31, 33, 40, 215, 58, 245},
				Name:      "testing.groups.group1",
				Modified:  timestamppb.Now(),
			},
			false,
			errors.ErrGroupInvalidCreated,
		},
		{
			&api.ConsumerGroup{
				ProjectId: []byte{1, 134, 179, 81, 86, 251, 48, 108, 44, 19, 143, 243, 195, 87, 134, 80},
				Id:        []byte{1, 134, 179, 108, 62, 211, 134, 53, 49, 102, 31, 33, 40, 215, 58, 245},
				Name:      "testing.groups.group1",
				Created:   timestamppb.Now(),
			},
			false,
			errors.ErrGroupInvalidModified,
		},
	}

	for i, tc := range testCases {
		err := meta.ValidateGroup(tc.group, tc.partial)
		if tc.err == nil {
			require.NoError(t, err, "failed testcase %d -- expected no error", i)
		} else {
			require.ErrorIs(t, err, tc.err, "failed testcase %d -- expected matching error", i)
		}
	}
}
