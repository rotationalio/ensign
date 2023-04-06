package models_test

import (
	"context"

	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg/quarterdeck/db/models"
	"github.com/rotationalio/ensign/pkg/utils/pagination"
	"github.com/rotationalio/ensign/pkg/utils/ulids"
	"github.com/stretchr/testify/require"
)

func (m *modelTestSuite) TestGetOrg() {
	require := m.Require()
	ctx := context.Background()

	// Ensure GetOrg returns not found error
	org, err := models.GetOrg(ctx, ulids.New())
	require.ErrorIs(err, models.ErrNotFound)
	require.Nil(org)

	org, err = models.GetOrg(ctx, ulid.MustParse("01GKHJRF01YXHZ51YMMKV3RCMK"))
	require.NoError(err, "could not fetch organization from database")

	// Ensure model is fully populated
	require.Equal("01GKHJRF01YXHZ51YMMKV3RCMK", org.ID.String())
	require.Equal("Testing", org.Name)
	require.Equal("example.com", org.Domain)
	require.Equal(2, org.ProjectCount())
	require.NotEmpty(org.Created, "no created timestamp")
	require.NotEmpty(org.Modified, "no modified timestamp")

	// Test GetOrg by string ID
	org2, err := models.GetOrg(ctx, "01GKHJRF01YXHZ51YMMKV3RCMK")
	require.NoError(err, "could not fetch organization from database")
	require.Equal(org, org2)
}

func (m *modelTestSuite) TestCreateOrg() {
	require := m.Require()
	ctx := context.Background()
	defer m.ResetDB()

	// Ensure name and domain are required on the organization
	org := &models.Organization{}
	require.ErrorIs(org.Create(ctx), models.ErrInvalidOrganization)

	org.Name = "Testing Foundation"
	require.ErrorIs(org.Create(ctx), models.ErrInvalidOrganization)
	org.Domain = "testing"

	err := org.Create(ctx)
	require.NoError(err, "could not create valid organization")

	// Ensure model has been updated
	require.False(ulids.IsZero(org.ID))
	require.NotEmpty(org.Created)
	require.NotEmpty(org.Modified)

	// Fetch the organization from the database
	cmpt, err := models.GetOrg(ctx, org.ID)
	require.NoError(err, "could not fetch org from the database")
	require.Equal(org, cmpt)

	// Should not be able to create a duplicate organization with same domain
	err = org.Create(ctx)
	require.ErrorIs(err, models.ErrDuplicate)
}

func (m *modelTestSuite) TestListOrgs() {
	require := m.Require()
	ctx := context.Background()

	//test passing null userID results in error
	orgs, cursor, err := models.ListOrgs(ctx, ulids.Null, nil)
	require.ErrorIs(err, models.ErrMissingModelID, "userID is required for list queries")
	require.NotNil(err)
	require.Nil(cursor)
	require.Nil(orgs)

	//test passing invalid orgID results in error
	orgs, cursor, err = models.ListOrgs(ctx, 1, nil)
	require.Contains("cannot parse input: unknown type", err.Error())
	require.NotNil(err)
	require.Nil(cursor)
	require.Nil(orgs)

	// test passing in pagination.Cursor without page size results in error
	userID := ulid.MustParse("01GQYYKY0ECGWT5VJRVR32MFHM")
	_, _, err = models.ListUsers(ctx, userID, &pagination.Cursor{})
	require.ErrorIs(err, models.ErrMissingPageSize, "page size is required for list users queries with pagination")

	// Should return the two organizations Zendaya belongs to
	orgs, cursor, err = models.ListOrgs(ctx, userID, nil)
	require.NoError(err, "could not fetch all orgs for Zendaya")
	require.Nil(cursor, "should be no next page so no cursor")
	require.Len(orgs, 2, "expected 2 users for Zendaya")
	org := orgs[0]
	require.NotNil(org.ID)
	require.NotNil(org.Name)
	require.NotNil(org.Domain)
	require.Equal(org.ProjectCount(), 2, "expected 2 projects for organization Testing")

	// test pagination
	pages := 0
	nRows := 0
	cursor = pagination.New("", "", 1)
	for cursor != nil && pages < 100 {
		orgs, nextPage, err := models.ListOrgs(ctx, userID, cursor)
		require.NoError(err, "could not fetch page from server")
		if nextPage != nil {
			require.NotEqual(cursor.StartIndex, nextPage.StartIndex)
			require.NotEqual(cursor.EndIndex, nextPage.EndIndex)
			require.Equal(cursor.PageSize, nextPage.PageSize)
		}

		pages++
		nRows += len(orgs)
		cursor = nextPage
	}

	require.Equal(2, pages, "expected 2 pages")
	require.Equal(2, nRows, "expected 2 results")
}

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
