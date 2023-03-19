package meta_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/google/uuid"
	api "github.com/rotationalio/ensign/pkg/ensign/api/v1beta1"
	"github.com/rotationalio/ensign/pkg/ensign/store/errors"
	"github.com/rotationalio/ensign/pkg/ensign/store/meta"
	"github.com/rotationalio/ensign/pkg/utils/ulids"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *metaTestSuite) TestSameGroupName() {
	// Two different projects should be able to store a group with the same name without
	// conflict with or without group IDs specified by the user.
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
	s.T().Skip("not implemented yet")
}

func (s *readonlyMetaTestSuite) TestGetOrCreateGroup() {
	s.T().Skip("not implemented yet")
}

func (s *metaTestSuite) TestCreateGroup() {
	s.T().Skip("not implemented yet")
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

func (s *metaTestSuite) TestUpdateGroup() {
	s.T().Skip("not implemented yet")
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
	s.T().Skip("not implemented yet")
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

		s, _ := tc.group.Key()
		fmt.Println(s)

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
