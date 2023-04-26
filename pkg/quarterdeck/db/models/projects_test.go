package models_test

import (
	"context"
	"testing"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg/quarterdeck/db/models"
	"github.com/rotationalio/ensign/pkg/utils/pagination"
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

	require := m.Require()
	for i, tc := range testCases {
		op := &models.OrganizationProject{
			OrgID:     ulid.MustParse(tc.OrgID),
			ProjectID: ulid.MustParse(tc.ProjectID),
		}
		ok, err := op.Exists(context.Background())
		require.NoError(err, "test case %d failed", i)
		tc.Exists(m.T(), ok, "unexpected response from database")
	}
}

func (m *modelTestSuite) TestListProjects() {
	require := m.Require()
	ctx := context.Background()
	orgID := ulid.MustParse("01GKHJRF01YXHZ51YMMKV3RCMK")

	// Test orgID is required
	projects, cursor, err := models.ListProjects(ctx, ulids.Null, nil)
	require.ErrorIs(err, models.ErrMissingOrgID, "orgID is required for list queries")
	require.Nil(cursor)
	require.Nil(projects)

	// Test pageSize is required
	_, _, err = models.ListProjects(ctx, orgID, &pagination.Cursor{})
	require.ErrorIs(err, models.ErrMissingPageSize, "pagination is required for list queries")

	// Test fetch results in a single page with no previous page
	projects, cursor, err = models.ListProjects(ctx, orgID, nil)
	require.NoError(err, "could not list projects")
	require.Nil(cursor, "expected no next page token")
	require.Len(projects, 3, "expected 3 projects returned")

}

func (m *modelTestSuite) TestFetchProject() {
	require := m.Require()
	nullULID := ulids.Null.String()

	testCases := []struct {
		OrgID, ProjectID       string
		KeyCount, RevokedCount int64
		Err                    error
	}{
		{nullULID, nullULID, 0, 0, models.ErrNotFound},
		{"01GKHJRF01YXHZ51YMMKV3RCMK", nullULID, 0, 0, models.ErrNotFound},
		{nullULID, "01GQ7P8DNR9MR64RJR9D64FFNT", 0, 0, models.ErrNotFound},
		{"01GKHJRF01YXHZ51YMMKV3RCMK", "01GQ7P8DNR9MR64RJR9D64FFNT", 2, 0, nil},
		{"01GQFQ14HXF2VC7C1HJECS60XX", "01GQFQCFC9P3S7QZTPYFVBJD7F", 3, 3, nil},
		{"01GKHJRF01YXHZ51YMMKV3RCMK", "01GQFR0KM5S2SSJ8G5E086VQ9K", 9, 3, nil},
		{"01GKHJRF01YXHZ51YMMKV3RCMK", "01GYYRSHRABN0S04ZZ4PXAK6VV", 0, 0, nil},
		{"01GKHJRF01YXHZ51YMMKV3RCMK", "01GQFQCFC9P3S7QZTPYFVBJD7F", 0, 0, models.ErrNotFound},
		{"01GQFQ14HXF2VC7C1HJECS60XX", "01GQ7P8DNR9MR64RJR9D64FFNT", 0, 0, models.ErrNotFound},
		{ulids.New().String(), "01GQ7P8DNR9MR64RJR9D64FFNT", 0, 0, models.ErrNotFound},
		{"01GQFQ14HXF2VC7C1HJECS60XX", ulids.New().String(), 0, 0, models.ErrNotFound},
		{ulids.New().String(), ulids.New().String(), 0, 0, models.ErrNotFound},
	}

	for i, tc := range testCases {
		project, err := models.FetchProject(context.Background(), ulid.MustParse(tc.ProjectID), ulid.MustParse(tc.OrgID))

		if tc.Err != nil {
			require.Error(err, "expected error on test case %d", i)
			require.ErrorIs(err, tc.Err, "expected error to match on test case %d", i)
		} else {
			require.NoError(err, "expected no error on test case %d", i)
			require.Equal(tc.KeyCount, project.APIKeyCount, "expected key count to match on test case %d", i)
			require.Equal(tc.RevokedCount, project.RevokedCount, "expected revoked key count to match on test case %d", i)

			// Test other fetch details
			require.Equal(tc.OrgID, project.OrgID.String())
			require.Equal(tc.ProjectID, project.ProjectID.String())
			require.NotEmpty(project.Created)
			require.NotEmpty(project.Modified)
		}
	}

}

func TestProjectToAPI(t *testing.T) {
	project := &models.Project{
		OrganizationProject: models.OrganizationProject{
			OrgID:     ulid.MustParse("01GYX96VN5FV9PSV6VDBQJV7BP"),
			ProjectID: ulid.MustParse("01GYX97KZV91M2APJ5SYC0GCYC"),
			Base: models.Base{
				Created:  "2023-04-25T17:41:46-05:00",
				Modified: "2023-04-25T17:41:52-05:00",
			},
		},
		APIKeyCount:  23,
		RevokedCount: 98,
	}

	serial := project.ToAPI()

	require.NotEmpty(t, serial.OrgID, "expected org_id to be set on the API response")
	require.Equal(t, project.OrgID, serial.OrgID, "expected org_id to match the model")

	require.NotEmpty(t, serial.ProjectID, "expected project_id to be set on the API response")
	require.Equal(t, project.ProjectID, serial.ProjectID, "expected project_id to match the model")

	require.NotZero(t, serial.APIKeysCount, "expected apikey count to be set on the API response")
	require.Equal(t, int(project.APIKeyCount), serial.APIKeysCount, "expected apikey count to match the model")

	require.NotZero(t, serial.RevokedCount, "expected revoked count to be set on the API response")
	require.Equal(t, int(project.RevokedCount), serial.RevokedCount, "expected revoked count to match the model")

	require.NotZero(t, serial.Created, "expected created to be set on the API response")
	require.Equal(t, project.Created, serial.Created.Format(time.RFC3339), "expected created to match the model")

	require.NotZero(t, serial.Modified, "expected modified to be set on the API response")
	require.Equal(t, project.Modified, serial.Modified.Format(time.RFC3339), "expected modified to match the model")
}
