package models_test

import (
	"context"

	"github.com/rotationalio/ensign/pkg/quarterdeck/db/models"
)

// If the fixtures change, these will need to be updated.
var (
	nUserFixtures         = int64(3)
	nOrganizationFixtures = int64(3)
)

func (m *modelTestSuite) TestCountUsers() {
	require := m.Require()

	nUsers, err := models.CountUsers(context.Background())
	require.NoError(err, "could not count the number of users")
	require.Equal(nUserFixtures, nUsers, "unexpected number of users returned, have the fixtures changed?")
}

func (m *modelTestSuite) TestCountOrganizations() {
	require := m.Require()

	nOrgs, err := models.CountOrganizations(context.Background())
	require.NoError(err, "could not count the number of organizations")
	require.Equal(nOrganizationFixtures, nOrgs, "unexpected number of organizations returned, have the fixtures changed?")
}
