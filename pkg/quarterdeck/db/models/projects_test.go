package models_test

import (
	"context"

	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg/quarterdeck/db/models"
	"github.com/rotationalio/ensign/pkg/utils/ulids"
	"github.com/stretchr/testify/require"
)

func (m *modelTestSuite) TestCreateOrganizationProject() {
	require := m.Require()
	defer m.ResetDB()

	// Test Happy Path
	orgID := ulid.MustParse("01GKHJRF01YXHZ51YMMKV3RCMK")
	projectID := ulids.New()
	op := &models.OrganizationProject{OrgID: orgID, ProjectID: projectID}

	err := op.Save(context.Background())
	require.NoError(err, "could not create an organization project mapping")

	// Test Error Paths
	testCases := []struct {
		OrgID, ProjectID ulid.ULID
		Error            error
		Msg              string
	}{
		{
			OrgID:     ulids.Null,
			ProjectID: projectID,
			Error:     models.ErrMissingOrgID,
			Msg:       "orgID is required",
		},
		{
			OrgID:     orgID,
			ProjectID: ulids.Null,
			Error:     models.ErrMissingProjectID,
			Msg:       "projectID is required",
		},
		{
			OrgID:     orgID,
			ProjectID: projectID,
			Error:     models.ErrDuplicate,
			Msg:       "cannot create a duplicate organization project mapping",
		},
		{
			OrgID:     ulid.MustParse("01GQFQ14HXF2VC7C1HJECS60XX"),
			ProjectID: projectID,
			Error:     models.ErrDuplicate,
			Msg:       "cannot duplicate project ID in database for another organization",
		},
		{
			OrgID:     ulids.New(),
			ProjectID: ulids.New(),
			Error:     models.ErrMissingRelation,
			Msg:       "cannot create a mapping for an organization that doesn't exist",
		},
	}

	for _, tc := range testCases {
		op = &models.OrganizationProject{OrgID: tc.OrgID, ProjectID: tc.ProjectID}
		err = op.Save(context.Background())
		require.Error(err, "expected save to fail for test %q", tc.Msg)
		require.ErrorIs(err, tc.Error, tc.Msg)
	}
}

func (m *modelTestSuite) TestOrganizationProjectExists() {
	nullULID := ulids.Null.String()
	testCases := []struct {
		OrgID, ProjectID string
		Exists           require.BoolAssertionFunc
	}{
		{nullULID, nullULID, require.False},
		{"01GKHJRF01YXHZ51YMMKV3RCMK", nullULID, require.False},
		{nullULID, "01GQ7P8DNR9MR64RJR9D64FFNT", require.False},
		{"01GKHJRF01YXHZ51YMMKV3RCMK", "01GQ7P8DNR9MR64RJR9D64FFNT", require.True},
		{"01GQFQ14HXF2VC7C1HJECS60XX", "01GQFQCFC9P3S7QZTPYFVBJD7F", require.True},
		{"01GKHJRF01YXHZ51YMMKV3RCMK", "01GQFQCFC9P3S7QZTPYFVBJD7F", require.False},
		{"01GQFQ14HXF2VC7C1HJECS60XX", "01GQ7P8DNR9MR64RJR9D64FFNT", require.False},
		{ulids.New().String(), "01GQ7P8DNR9MR64RJR9D64FFNT", require.False},
		{"01GQFQ14HXF2VC7C1HJECS60XX", ulids.New().String(), require.False},
		{ulids.New().String(), ulids.New().String(), require.False},
	}

	for i, tc := range testCases {
		op := &models.OrganizationProject{
			OrgID:     ulid.MustParse(tc.OrgID),
			ProjectID: ulid.MustParse(tc.ProjectID),
		}
		ok, err := op.Exists(context.Background())
		require.NoError(m.T(), err, "test case %d failed", i)
		tc.Exists(m.T(), ok, "unexpected response from database")
	}
}
