package tenant_test

import (
	"context"
	"net/http"
	"time"

	"github.com/oklog/ulid/v2"
	qd "github.com/rotationalio/ensign/pkg/quarterdeck/api/v1"
	"github.com/rotationalio/ensign/pkg/quarterdeck/mock"
	perms "github.com/rotationalio/ensign/pkg/quarterdeck/permissions"
	"github.com/rotationalio/ensign/pkg/quarterdeck/tokens"
	"github.com/rotationalio/ensign/pkg/tenant/api/v1"
)

func (s *tenantTestSuite) TestOrganizationDetail() {
	require := s.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	orgID := "01GKHJRF01YXHZ51YMMKV3RCMK"
	org := &qd.Organization{
		ID:       ulid.MustParse(orgID),
		Name:     "Rotational Labs",
		Domain:   "rotational.io",
		Created:  time.Now(),
		Modified: time.Now().Add(time.Hour),
	}

	// Initial Quarterdeck mock should return 200 OK with the organization
	s.quarterdeck.OnOrganizations(orgID, mock.UseStatus(http.StatusOK), mock.UseJSONFixture(org))

	// Setup the initial claims fixture
	claims := &tokens.Claims{
		Name:        "Jannel P. Hudson",
		Email:       "jannel@example.com",
		OrgID:       "02ABCDEF01YXHZ51YMMKV3RCMK",
		Permissions: []string{"read:nothing"},
	}

	// Endpoint must be authenticated
	_, err := s.client.OrganizationDetail(ctx, "invalid")
	s.requireError(err, http.StatusUnauthorized, "this endpoint requires authentication")

	// User must have the correct permissions
	require.NoError(s.SetClientCredentials(claims), "could not set client credentials")
	_, err = s.client.OrganizationDetail(ctx, "invalid")
	s.requireError(err, http.StatusUnauthorized, "user does not have permission to perform this operation")

	// Organization must be parseable
	claims.Permissions = []string{perms.ReadOrganizations}
	require.NoError(s.SetClientCredentials(claims), "could not set client credentials")
	_, err = s.client.OrganizationDetail(ctx, "invalid")
	s.requireError(err, http.StatusBadRequest, "could not parse organization ID")

	// User can only access their own organization
	_, err = s.client.OrganizationDetail(ctx, orgID)
	s.requireError(err, http.StatusNotFound, "organization not found")

	// Successfully retrieving organization details
	claims.OrgID = orgID
	require.NoError(s.SetClientCredentials(claims), "could not set client credentials")
	expected := &api.Organization{
		ID:       orgID,
		Name:     org.Name,
		Domain:   org.Domain,
		Created:  org.Created.Format(time.RFC3339Nano),
		Modified: org.Modified.Format(time.RFC3339Nano),
	}
	reply, err := s.client.OrganizationDetail(ctx, orgID)
	require.NoError(err, "could not retrieve organization details")
	require.Equal(expected, reply, "organization details did not match")

	// Test that the method returns an error if Quarterdeck returns an error
	s.quarterdeck.OnOrganizations(orgID, mock.UseStatus(http.StatusUnauthorized))
	_, err = s.client.OrganizationDetail(ctx, orgID)
	s.requireError(err, http.StatusUnauthorized, "could not detail organization")
}
