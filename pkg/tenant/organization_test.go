package tenant_test

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/oklog/ulid/v2"
	qd "github.com/rotationalio/ensign/pkg/quarterdeck/api/v1"
	"github.com/rotationalio/ensign/pkg/quarterdeck/mock"
	perms "github.com/rotationalio/ensign/pkg/quarterdeck/permissions"
	"github.com/rotationalio/ensign/pkg/quarterdeck/tokens"
	"github.com/rotationalio/ensign/pkg/tenant/api/v1"
	"github.com/rotationalio/ensign/pkg/tenant/db"
	"github.com/trisacrypto/directory/pkg/trtl/pb/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *tenantTestSuite) TestOrganizationList() {
	require := s.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	page := &qd.OrganizationList{
		Organizations: []*qd.Organization{
			{
				ID:        ulid.MustParse("01GKHJRF01YXHZ51YMMKV3RCMK"),
				Name:      "Rotational Labs",
				Domain:    "rotational.io",
				Created:   time.Now(),
				LastLogin: time.Now().Add(time.Hour),
			},
			{
				ID:        ulid.MustParse("02GKHJRF01YXHZ51YMMKV3RABC"),
				Name:      "McDowell's",
				Domain:    "mcdowells.com",
				Created:   time.Now(),
				LastLogin: time.Now().Add(time.Hour),
			},
		},
	}

	expected := []*api.Organization{
		{
			ID:        "01GKHJRF01YXHZ51YMMKV3RCMK",
			Name:      "Rotational Labs",
			Domain:    "rotational.io",
			Created:   page.Organizations[0].Created.Format(time.RFC3339Nano),
			LastLogin: page.Organizations[0].LastLogin.Format(time.RFC3339Nano),
		},
		{
			ID:        "02GKHJRF01YXHZ51YMMKV3RABC",
			Name:      "McDowell's",
			Domain:    "mcdowells.com",
			Created:   page.Organizations[1].Created.Format(time.RFC3339Nano),
			LastLogin: page.Organizations[1].LastLogin.Format(time.RFC3339Nano),
		},
	}

	// Initial Quarterdeck mock should return 200 OK with the organization
	s.quarterdeck.OnOrganizations("", mock.UseStatus(http.StatusOK), mock.UseJSONFixture(page), mock.RequireAuth())

	// Setup the initial claims fixture
	claims := &tokens.Claims{
		Name:        "Jannel P. Hudson",
		Email:       "jannel@example.com",
		Permissions: []string{"read:nothing"},
	}

	// Endpoint must be authenticated
	_, err := s.client.OrganizationList(ctx, &api.PageQuery{})
	s.requireError(err, http.StatusUnauthorized, "this endpoint requires authentication")

	// User must have the correct permissions
	require.NoError(s.SetClientCredentials(claims), "could not set client credentials")
	_, err = s.client.OrganizationList(ctx, &api.PageQuery{})
	s.requireError(err, http.StatusUnauthorized, "user does not have permission to perform this operation")

	// Test returning a page of organizations.
	claims.Permissions = []string{perms.ReadOrganizations}
	require.NoError(s.SetClientCredentials(claims), "could not set client credentials")
	rep, err := s.client.OrganizationList(ctx, &api.PageQuery{PageSize: 1})
	require.NoError(err, "could not list organizations")
	require.Equal(page.NextPageToken, rep.NextPageToken, "expected next page token to match")
	require.Equal(len(page.Organizations), len(rep.Organizations), "expected organizations count to match")
	require.Equal(expected, rep.Organizations, "expected organizations data to match")

	// Should return an error if Quarterdeck returns an error.
	s.quarterdeck.OnOrganizations("", mock.UseError(http.StatusInternalServerError, "could not list organizations"), mock.RequireAuth())
	_, err = s.client.OrganizationList(ctx, &api.PageQuery{})
	s.requireError(err, http.StatusInternalServerError, "could not list organizations", "expected error when Quarterdeck returns an error")
}

func (s *tenantTestSuite) TestOrganizationDetail() {
	require := s.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	orgID := "01GKHJRF01YXHZ51YMMKV3RCMK"
	org := &qd.Organization{
		ID:       ulid.MustParse(orgID),
		Name:     "Rotational Labs",
		Domain:   "rotational.io",
		Projects: 1,
		Created:  time.Now(),
		Modified: time.Now().Add(time.Hour),
	}

	members := []*db.Member{
		{
			OrgID: ulid.MustParse(orgID),
			ID:    ulid.MustParse("01GKHJRF01YXHZ51YMMKV3RCMK"),
			Name:  "Jannel P. Hudson",
			Role:  perms.RoleOwner,
		},
		{
			OrgID: ulid.MustParse(orgID),
			ID:    ulid.MustParse("02GKHJRF01YXHZ51YMMKV3RABC"),
			Name:  "John Doe",
			Role:  perms.RoleMember,
		},
	}

	// Setup the trtl mock to list the member fixtures
	trtl := db.GetMock()
	defer trtl.Reset()

	trtl.OnCursor = func(in *pb.CursorRequest, stream pb.Trtl_CursorServer) error {
		if !bytes.Equal(in.Prefix, org.ID[:]) || in.Namespace != db.MembersNamespace {
			return status.Error(codes.FailedPrecondition, "unexpected cursor request")
		}

		// Send back some data and terminate
		for i, member := range members {
			data, err := member.MarshalValue()
			require.NoError(err, "could not marshal data")
			stream.Send(&pb.KVPair{
				Key:       []byte(fmt.Sprintf("key %d", i)),
				Value:     data,
				Namespace: in.Namespace,
			})
		}
		return nil
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
	s.requireError(err, http.StatusNotFound, "organization not found")

	// User can only access their own organization
	_, err = s.client.OrganizationDetail(ctx, orgID)
	s.requireError(err, http.StatusNotFound, "organization not found")

	// Successfully retrieving organization details
	claims.OrgID = orgID
	require.NoError(s.SetClientCredentials(claims), "could not set client credentials")
	expected := &api.Organization{
		ID:       orgID,
		Name:     org.Name,
		Owner:    members[0].Name,
		Domain:   org.Domain,
		Projects: 1,
		Created:  org.Created.Format(time.RFC3339Nano),
	}
	reply, err := s.client.OrganizationDetail(ctx, orgID)
	require.NoError(err, "could not retrieve organization details")
	require.Equal(expected, reply, "organization details did not match")

	// Test that the method returns an error if Quarterdeck returns an error
	s.quarterdeck.OnOrganizations(orgID, mock.UseError(http.StatusNotFound, "organization not found"))
	_, err = s.client.OrganizationDetail(ctx, orgID)
	s.requireError(err, http.StatusNotFound, "organization not found")
}
