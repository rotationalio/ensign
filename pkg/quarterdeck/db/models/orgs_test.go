package models_test

import (
	"context"

	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg/quarterdeck/db/models"
	ulids "github.com/rotationalio/ensign/pkg/utils/ulid"
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
	require.NotEmpty(org.Created, "no created timestamp")
	require.NotEmpty(org.Modified, "no modified timestamp")
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
