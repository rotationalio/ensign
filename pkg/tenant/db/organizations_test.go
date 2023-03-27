package db_test

import (
	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg/tenant/db"
)

func (s *dbTestSuite) TestVerifyOrg() {
	require := s.Require()

	org := &db.Organization{ID: ulid.ULID{}}
	tenant := &db.Tenant{OrgID: ulid.MustParse("01GMBVR86186E0EKCHQK4ESJB1")}

	// Require error if org ULID is missing.
	ok, err := db.VerifyOrg(org.ID, tenant.OrgID)
	require.ErrorIs(err, db.ErrMissingOrgID, "expected error when org id is missing")
	require.False(ok)

	// Require error if tenant org ID is missing.
	org.ID = ulid.MustParse("01GMTWFK4XZY597Y128KXQ4WHP")
	tenant.OrgID = ulid.ULID{}
	ok, err = db.VerifyOrg(org.ID, tenant.OrgID)
	require.ErrorIs(err, db.ErrMissingID, "expected error when model org id is missing")
	require.False(ok)

	// Require error if org ID and tenant org ID are different.
	tenant.OrgID = ulid.MustParse("01GMBVR86186E0EKCHQK4ESJB1")
	ok, err = db.VerifyOrg(org.ID, tenant.OrgID)
	require.ErrorIs(err, db.ErrOrgNotVerified)
	require.False(ok)

	// Set tenant org ID to org ID and test.
	tenant.OrgID = ulid.MustParse("01GMTWFK4XZY597Y128KXQ4WHP")
	ok, err = db.VerifyOrg(org.ID, tenant.OrgID)
	require.NoError(err, "could not verify org")
	require.True(ok)
}
