package db_test

import (
	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg/tenant/db"
)

func (s *dbTestSuite) TestVerifyOrg() {
	require := s.Require()

	orgID := ulid.MustParse("01GMTWFK4XZY597Y128KXQ4WHP")
	tenantOrgID := ulid.MustParse("01GMBVR86186E0EKCHQK4ESJB1")

	ok, err := db.VerifyOrg(orgID, tenantOrgID)
	require.NoError(err, "could not verify org id")
	require.Equal(ok, false)

	tenantOrgID2 := ulid.MustParse("01GMTWFK4XZY597Y128KXQ4WHP")
	ok, err = db.VerifyOrg(orgID, tenantOrgID2)
	require.NoError(err, "could not verify org id")
	require.Equal(ok, true)
}
