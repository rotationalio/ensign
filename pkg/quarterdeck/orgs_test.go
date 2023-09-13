package quarterdeck_test

import (
	"context"
	"net/http"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg/quarterdeck/api/v1"
	perms "github.com/rotationalio/ensign/pkg/quarterdeck/permissions"
	"github.com/rotationalio/ensign/pkg/quarterdeck/tokens"
	"github.com/rotationalio/ensign/pkg/utils/responses"
	"github.com/rotationalio/ensign/pkg/utils/ulids"
)

func (s *quarterdeckTestSuite) TestOrganizationDetail() {
	require := s.Require()
	defer s.ResetDatabase()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Retrieving an organization requires authentication
	_, err := s.client.OrganizationDetail(ctx, "invalid")
	s.CheckError(err, http.StatusUnauthorized, "this endpoint requires authentication")

	// Retrieving an organization requires the read:organizations permission
	claims := &tokens.Claims{
		Name:  "Jannel P. Hudson",
		Email: "jannel@example.com",
	}
	ctx = s.AuthContext(ctx, claims)
	orgID := ulids.New()
	_, err = s.client.OrganizationDetail(ctx, orgID.String())
	s.CheckError(err, http.StatusForbidden, "user does not have permission to perform this operation")

	// Valid ID is required in the URL
	claims.Permissions = []string{perms.ReadOrganizations}
	ctx = s.AuthContext(ctx, claims)
	_, err = s.client.OrganizationDetail(ctx, "invalid")
	s.CheckError(err, http.StatusNotFound, "organization not found")

	// Specified organization must match the user's organization
	claims.OrgID = ulids.New().String()
	ctx = s.AuthContext(ctx, claims)
	_, err = s.client.OrganizationDetail(ctx, orgID.String())
	s.CheckError(err, http.StatusNotFound, "organization not found")

	// Organization must exist
	claims.OrgID = orgID.String()
	ctx = s.AuthContext(ctx, claims)
	_, err = s.client.OrganizationDetail(ctx, orgID.String())
	s.CheckError(err, http.StatusNotFound, "organization not found")

	// Successfully retrieving the organization
	claims.OrgID = "01GKHJRF01YXHZ51YMMKV3RCMK"
	ctx = s.AuthContext(ctx, claims)
	_, err = s.client.OrganizationDetail(ctx, claims.OrgID)
	require.NoError(err, "could not retrieve organization details")
}

func (s *quarterdeckTestSuite) TestOrganizationUpdate() {
	require := s.Require()
	defer s.ResetDatabase()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Updating an organization requires authentication
	req := &api.Organization{
		ID: ulids.New(),
	}
	_, err := s.client.OrganizationUpdate(ctx, req)
	s.CheckError(err, http.StatusUnauthorized, "this endpoint requires authentication")

	// Updating an organization requires the edit:organizations permission
	claims := &tokens.Claims{
		Name:  "Jannel P. Hudson",
		Email: "jannel@example.com",
	}
	ctx = s.AuthContext(ctx, claims)
	_, err = s.client.OrganizationUpdate(ctx, req)
	s.CheckError(err, http.StatusForbidden, "user does not have permission to perform this operation")

	// Specified organization must match the user's organization
	claims.Permissions = []string{perms.EditOrganizations}
	claims.OrgID = ulids.New().String()
	ctx = s.AuthContext(ctx, claims)
	_, err = s.client.OrganizationUpdate(ctx, req)
	s.CheckError(err, http.StatusNotFound, responses.ErrOrganizationNotFound)

	// Organization must contain a name
	claims.OrgID = req.ID.String()
	ctx = s.AuthContext(ctx, claims)
	_, err = s.client.OrganizationUpdate(ctx, req)
	s.CheckError(err, http.StatusBadRequest, "missing required field: name")

	// Organization must contain a domain
	req.Name = "My Organization"
	_, err = s.client.OrganizationUpdate(ctx, req)
	s.CheckError(err, http.StatusBadRequest, "missing required field: domain")

	// Organization must exist
	req.Domain = "checkers-io"
	claims.OrgID = req.ID.String()
	ctx = s.AuthContext(ctx, claims)
	_, err = s.client.OrganizationUpdate(ctx, req)
	s.CheckError(err, http.StatusNotFound, responses.ErrOrganizationNotFound)

	// Error should be returned if domain already exists
	req.ID = ulid.MustParse("01GKHJRF01YXHZ51YMMKV3RCMK")
	claims.OrgID = req.ID.String()
	ctx = s.AuthContext(ctx, claims)
	_, err = s.client.OrganizationUpdate(ctx, req)
	s.CheckError(err, http.StatusConflict, responses.ErrDomainAlreadyExists)

	// Successfully updating the organization
	req.Domain = "newdomain.com"
	rep, err := s.client.OrganizationUpdate(ctx, req)
	require.NoError(err, "could not update organization")
	require.Equal("My Organization", rep.Name, "expected new name to be returned")
	require.Equal("newdomain.com", rep.Domain, "expected new domain to be returned")

	// Organization should be updated in the database
	claims.Permissions = []string{perms.ReadOrganizations}
	ctx = s.AuthContext(ctx, claims)
	org, err := s.client.OrganizationDetail(ctx, rep.ID.String())
	require.NoError(err, "could not retrieve organization")
	require.Equal("My Organization", org.Name, "expected name to be updated")
	require.Equal("newdomain.com", org.Domain, "expected domain to be updated")
}

func (s *quarterdeckTestSuite) TestOrganizationList() {
	require := s.Require()
	defer s.ResetDatabase()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Listing organizations requires authentication
	req := &api.OrganizationPageQuery{}
	_, err := s.client.OrganizationList(ctx, req)
	s.CheckError(err, http.StatusUnauthorized, "this endpoint requires authentication")

	// Listing organizations requires the read:organizations permission
	claims := &tokens.Claims{
		Name:  "Zendaya Longeye",
		Email: "zendaya@testing.io",
	}
	ctx = s.AuthContext(ctx, claims)
	_, err = s.client.OrganizationList(ctx, req)
	s.CheckError(err, http.StatusForbidden, "user does not have permission to perform this operation")

	// Create valid claims for accessing the API
	claims.Subject = "01GQYYKY0ECGWT5VJRVR32MFHM"
	claims.OrgID = "01GKHJRF01YXHZ51YMMKV3RCMK"
	claims.Permissions = []string{perms.ReadOrganizations}
	ctx = s.AuthContext(ctx, claims)

	// Should be able to list all users for the specified organization
	page, err := s.client.OrganizationList(ctx, req)
	require.NoError(err, "could not fetch users")
	require.Len(page.Organizations, 2, "expected 2 results back from the fixtures")
	require.Empty(page.NextPageToken, "expected no next page token in response")

	// Should be able to paginate the request for the specified organization
	req.PageSize = 1
	page, err = s.client.OrganizationList(ctx, req)
	require.NoError(err, "could not fetch paginated users")
	require.Len(page.Organizations, 1, "expected 1 result back from the fixtures")
	require.NotEmpty(page.NextPageToken, "expected next page token in response")

	// Test fetching the next page with the next page token
	req.NextPageToken = page.NextPageToken
	page2, err := s.client.OrganizationList(ctx, req)
	require.NoError(err, "could not fetch paginated api keys")
	require.Len(page2.Organizations, 1, "expected 1 result back from the fixtures")
	require.Empty(page2.NextPageToken, "expected no next page token in response")
	require.NotEqual(page.Organizations[0].Name, page2.Organizations[0].Name, "expected a new page of results")
	require.Equal(page2.Organizations[0].Projects, 10, "expected 10 projects for the Checkers organization")
	require.NotEmpty(page2.Organizations[0].LastLogin, "expected a last login time for Zendaya in the Checkers organization")

	// maximum number of requests is 2, break when pagination is complete.
	req.NextPageToken = ""
	nPages, nResults := 0, 0
	for {
		page, err = s.client.OrganizationList(ctx, req)
		require.NoError(err, "could not fetch page of results")

		nPages++
		nResults += len(page.Organizations)

		if page.NextPageToken != "" {
			req.NextPageToken = page.NextPageToken
		} else {
			break
		}
	}

	require.Equal(nPages, 2, "expected 2 pages")
	require.Equal(nResults, 2, "expected 2 results")
}

func (s *quarterdeckTestSuite) TestWorkspaceLookup() {
	// NOTE: no need to reset database as this is a read-only test
	require := s.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	authctx := s.AuthContext(ctx, &tokens.Claims{OrgID: "01GQFQ14HXF2VC7C1HJECS60XX", Name: "Jannel P. Hudson", Email: "jannel@example.com", Permissions: []string{perms.ReadOrganizations}})

	s.Run("RequireAuth", func() {
		_, err := s.client.WorkspaceLookup(ctx, &api.WorkspaceQuery{Domain: "rotational-labs"})
		s.CheckError(err, http.StatusUnauthorized, "this endpoint requires authentication")
	})

	s.Run("RequirePerms", func() {
		claims := &tokens.Claims{
			Name:        "Jannel P. Hudson",
			Email:       "jannel@example.com",
			Permissions: []string{"read:foo"},
		}

		ctx = s.AuthContext(ctx, claims)
		_, err := s.client.WorkspaceLookup(ctx, &api.WorkspaceQuery{Domain: "rotational-labs"})
		s.CheckError(err, http.StatusForbidden, "user does not have permission to perform this operation")
	})

	s.Run("NoDomain", func() {
		_, err := s.client.WorkspaceLookup(authctx, &api.WorkspaceQuery{})
		s.CheckError(err, http.StatusBadRequest, responses.ErrBadWorkspaceLookup)
	})

	s.Run("WorkspaceNotFound", func() {
		rep, err := s.client.WorkspaceLookup(authctx, &api.WorkspaceQuery{Domain: "walnut-tosser"})
		s.CheckError(err, http.StatusNotFound, responses.ErrWorkspaceNotFound)
		require.Nil(rep, "expected response to be nil")
	})

	s.Run("WorkspaceFoundFull", func() {
		rep, err := s.client.WorkspaceLookup(authctx, &api.WorkspaceQuery{Domain: "checkers-io"})
		require.NoError(err, "could not execute request")
		require.NotNil(rep, "expected a response back")

		// Expect a full response because Jannel belongs to the example organization
		require.Equal(ulid.MustParse("01GQFQ14HXF2VC7C1HJECS60XX"), rep.OrgID)
		require.Equal("Checkers", rep.Name)
		require.Equal("checkers-io", rep.Domain)
		require.False(rep.IsAvailable)
	})

	s.Run("WorkspaceFoundPartial", func() {
		rep, err := s.client.WorkspaceLookup(authctx, &api.WorkspaceQuery{Domain: "ghost-co"})
		require.NoError(err, "could not execute request")
		require.NotNil(rep, "expected a response back")

		// Expect a partial response because Jannel does not belong to the ghost organization
		require.Empty(rep.OrgID)
		require.Empty(rep.Name)
		require.Equal("ghost-co", rep.Domain)
		require.False(rep.IsAvailable)
	})

	s.Run("WorkspaceAvailable", func() {
		rep, err := s.client.WorkspaceLookup(authctx, &api.WorkspaceQuery{Domain: "walnut-tosser", CheckAvailable: true})
		require.NoError(err, "expected a 200 response with check available")
		require.True(rep.IsAvailable)
		require.Equal(rep.Domain, "walnut-tosser")
		require.Empty(rep.Name)
		require.Empty(rep.OrgID)
	})

	s.Run("WorkspaceNotAvailablePartial", func() {
		rep, err := s.client.WorkspaceLookup(authctx, &api.WorkspaceQuery{Domain: "example-com", CheckAvailable: true})
		require.NoError(err, "expected a 200 response with check available")

		require.False(rep.IsAvailable)
		require.Equal(rep.Domain, "example-com")
		require.Empty(rep.Name)
		require.Empty(rep.OrgID)
	})

	s.Run("WorkspaceNotAvailableFull", func() {
		rep, err := s.client.WorkspaceLookup(authctx, &api.WorkspaceQuery{Domain: "checkers-io", CheckAvailable: true})
		require.NoError(err, "expected a 200 response with check available")

		require.False(rep.IsAvailable)
		require.Equal(rep.Domain, "checkers-io")
		require.Equal("Checkers", rep.Name)
		require.Equal(ulid.MustParse("01GQFQ14HXF2VC7C1HJECS60XX"), rep.OrgID)
	})
}
