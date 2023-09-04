package tenant_test

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"time"

	qd "github.com/rotationalio/ensign/pkg/quarterdeck/api/v1"
	"github.com/rotationalio/ensign/pkg/quarterdeck/mock"
	perms "github.com/rotationalio/ensign/pkg/quarterdeck/permissions"
	"github.com/rotationalio/ensign/pkg/quarterdeck/tokens"
	"github.com/rotationalio/ensign/pkg/tenant/api/v1"
	"github.com/rotationalio/ensign/pkg/tenant/db"
	"github.com/rotationalio/ensign/pkg/utils/responses"
	"github.com/trisacrypto/directory/pkg/trtl/pb/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	accessToken  = "eyJhbGciOiJSUzI1NiIsImtpZCI6IjAxR0U2MkVYWFIwWDA1NjFYRDUzUkRGQlFKIiwidHlwIjoiSldUIn0.eyJpc3MiOiJodHRwOi8vbG9jYWxob3N0OjMwMDEiLCJzdWIiOiIwMUg2UEdGQjRUMzRENFdXRVhRTUFHSk5NSyIsImF1ZCI6WyJodHRwOi8vbG9jYWxob3N0OjMwMDAiXSwiZXhwIjoxNjkzODY2NDY2LCJuYmYiOjE2OTM4NjI4NjYsImlhdCI6MTY5Mzg2Mjg2NiwianRpIjoiMDFoOWgxZ2J6ZHNlMXh2aDU0NnpqZTZ2MGUiLCJuYW1lIjoiS2F0ZSBIb2xsYW5kIiwiZW1haWwiOiJrYXRlQGV4YW1wbGUuY28iLCJwaWN0dXJlIjoiaHR0cHM6Ly93d3cuZ3JhdmF0YXIuY29tL2F2YXRhci84MGViYjNiMGRhZTNmNTUwZGU3MjAyMWJkY2Y0NWQwMCIsIm9yZyI6IjAxSDZQR0ZHNzFOMEFGRVZUSzNOSkI3MVQ5IiwicHJvamVjdCI6IjAxSDZQR0ZUSzJYNTNSR0cyS01TR1IyTTYxIiwicGVybWlzc2lvbnMiOlsicmVhZDpmb28iLCJlZGl0OmZvbyIsImRlbGV0ZTpmb28iXX0.wUxuTUuIOpF-d0lx5oCYCTapoNftjy84U3oK5tn5bQDK3rT-JqrmsXQZ8iXgLePxqenBKzbGmxB4ggXK1SpBukSDiTHC9LTxkfyJFBvSmeZBvqiG_u1T_AtajV-lLBTJaeEdsfsUYDE8V21QmS7vqaPTVm3v4-Uz8ZLSN3U-gsUdIEQPArtKwGep1scje_lMtJbztXGogenQUnXmMdA_wLHJeYZgZVGwJYfnMqXWdtR6R7QvcR9uhcepy1S0y24wtYR_Z3ZKhEJWMgxKnSNidPcWoUmwCbt9T30mePxqs2UxwZ3wsnJLEXF52hLyOsdx8dRG91T3k4JBVz-1tLsZJm8Q-q_vx4llhY-3211FqXa6b9I8GLEcXprDoU5ki3M-DjFhEJ53ICsvuzX-52mFmDNBE4cT-1LY3-kfxiINwXsMbVl6LsE0es1_leXmgfEzkNBkHEfNfmsZhdBZwkMCcq53uKvO2F8zMhZO4Ft8CtaFlBrrIJs1gcDGcFKqBvddlAYBtKXmZ_kpvB-ipxD5uPCZJB57m-m6czjZec21u4n_6yzA6OI8xRnVnw6Ankg7ZvkHs-S5Y4BReharBdk-eabZREVCcD_uOFOJ8z1bYRPYl-4OH4rDvSRuJp0PaWNaqNwLIrD51OPQy8lbwzo9gmsUadwwLNuSGjKHbJYU7n8"
	refreshToken = "eyJhbGciOiJSUzI1NiIsImtpZCI6IjAxR0U2MkVYWFIwWDA1NjFYRDUzUkRGQlFKIiwidHlwIjoiSldUIn0.eyJpc3MiOiJodHRwOi8vbG9jYWxob3N0OjMwMDEiLCJzdWIiOiIwMUg2UEdGQjRUMzRENFdXRVhRTUFHSk5NSyIsImF1ZCI6WyJodHRwOi8vbG9jYWxob3N0OjMwMDAiLCJodHRwOi8vbG9jYWxob3N0OjMwMDEvdjEvcmVmcmVzaCJdLCJleHAiOjE2OTM4NzAwNjYsIm5iZiI6MTY5Mzg2NTU2NiwiaWF0IjoxNjkzODYyODY2LCJqdGkiOiIwMWg5aDFnYnpkc2UxeHZoNTQ2emplNnYwZSIsIm9yZyI6IjAxSDZQR0ZHNzFOMEFGRVZUSzNOSkI3MVQ5IiwicHJvamVjdCI6IjAxSDZQR0ZUSzJYNTNSR0cyS01TR1IyTTYxIn0.n2yIJBmcP2ppf1RumbQUBc4z_ewP1ERtGYzKC4lVoxQyEV9OItHRSEdzW2SKd8duhlOMVxJE46UPxpnW4aYsIzFqOBrgap-Q7q_NP7SiS_GTPQm7Bmy9gKNK2EqEb49fTIGNN77XM6xpP1JbvGh83OZzMr1xYxlhjpw-8d0fM-gPG3tnNiGN8Bqnr68F5BXBc_0hYkJmY5RJ4fWh_ZzD4CaTwINkhZ9leM3eW6S-YKvD5-f67-m7fYzTRmY_t3NRwJDtwfPBNcmPYt-wyOMaGKRTyOGSmi_TbV_4wilOLU_lKf_o2xCkm0PIPt-0x977LGBJweIeQnkQq3E9lqIxZcvpKYit-6O_WFvrZsTDa_DvVwD5RZwhtammFGykn_2TvTp1fA_gnW9yKttu3a9HxUNE609Tt5CBbH0wsHZ83BdLngEmavn9lkqLXRCNlUqfew56N-ENnp1zHOwgEbOoJIqPTFoBN8Cc7-9YotVYR9eJdrKgOvtmTSwYzKtvdSeJv4InBe1wGOH6ljtQWYN2M-B_2wbONBGDrTucSBb_qEEpA0cZrUQi_vKE4ik2pkzLtKL6Vlv8KdEQij4Vh7WnHUmdvFZQjxT0kx0SX5Qv1m1ZtHLQ_Mt6NUgyqatup9jmepnSABFGAlc7ohdCwKFhJqHM6NsMcDoCYVcLs03HfYM"
)

func (s *tenantTestSuite) TestInvitePreview() {
	require := s.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Initial Quarterdeck mock returns a valid preview
	preview := &qd.UserInvitePreview{
		Email:       "leopold.wentzel@gmail.com",
		OrgName:     "Events R Us",
		InviterName: "Geoffrey",
		Role:        "Member",
		UserExists:  true,
	}
	s.quarterdeck.OnInvitesPreview("token1234", mock.UseStatus(http.StatusOK), mock.UseJSONFixture(preview))

	// Test successful preview request
	rep, err := s.client.InvitePreview(ctx, "token1234")
	require.NoError(err, "could not get preview invite")
	require.Equal(preview.Email, rep.Email, "expected email to match")
	require.Equal(preview.OrgName, rep.OrgName, "expected org name to match")
	require.Equal(preview.InviterName, rep.InviterName, "expected inviter name to match")
	require.Equal(preview.Role, rep.Role, "expected role to match")
	require.True(rep.HasAccount, "expected user to exist")

	// Test invalid invitation response is correctly forwarded by Tenant
	s.quarterdeck.OnInvitesPreview("token1234", mock.UseError(http.StatusBadRequest, "invalid invitation"))
	_, err = s.client.InvitePreview(ctx, "token1234")
	s.requireError(err, http.StatusBadRequest, "invalid invitation", "expected error when token is invalid")
}

func (s *tenantTestSuite) TestInviteAccept() {
	require := s.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Setup the trtl mock to return a member record on update
	trtl := db.GetMock()
	defer trtl.Reset()

	newClaims, err := tokens.ParseUnverifiedTokenClaims(accessToken)
	require.NoError(err, "could not parse claims from access token fixture")
	orgID := newClaims.ParseOrgID()
	require.NotEmpty(orgID, "expected org id to be set in access token fixture")
	userID := newClaims.ParseUserID()
	require.NotEmpty(userID, "expected user id to be set in access token fixture")

	member := &db.Member{
		OrgID: orgID,
		ID:    userID,
		Email: "kate@example.co",
		Role:  perms.RoleMember,
	}

	data, err := member.MarshalValue()
	require.NoError(err, "could not marshal member data")

	key, err := member.Key()
	require.NoError(err, "could not create member key from fixture")

	// trtl should return the member data on get
	trtl.OnGet = func(ctx context.Context, gr *pb.GetRequest) (*pb.GetReply, error) {
		switch gr.Namespace {
		case db.MembersNamespace:
			switch {
			case bytes.Equal(gr.Key, key):
				return &pb.GetReply{
					Value: data,
				}, nil
			default:
				return nil, status.Errorf(codes.NotFound, "member not found")
			}
		default:
			return nil, errors.New("unexpected namespace")
		}
	}

	// trtl should return OK on put
	trtl.OnPut = func(ctx context.Context, pr *pb.PutRequest) (_ *pb.PutReply, err error) {
		switch pr.Namespace {
		case db.MembersNamespace:
			switch {
			case bytes.Equal(pr.Key, key):
				// Ensure that the timestamps are updated
				var newMember db.Member
				if err = newMember.UnmarshalValue(pr.Value); err != nil {
					return nil, err
				}

				if newMember.JoinedAt.IsZero() {
					return nil, errors.New("joined at timestamp not set for invited member")
				}

				if newMember.LastActivity.IsZero() {
					return nil, errors.New("last activity timestamp not set for invited member")
				}
				return &pb.PutReply{}, nil
			default:
				return nil, status.Errorf(codes.NotFound, "member not found")
			}
		default:
			return nil, errors.New("unexpected namespace")
		}
	}

	// Initial Quarterdeck mock returns auth credentials
	creds := &qd.LoginReply{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	s.quarterdeck.OnInvitesAccept(mock.UseStatus(http.StatusOK), mock.UseJSONFixture(creds), mock.RequireAuth())

	// Set the initial claims fixture
	claims := &tokens.Claims{
		Name:  "Leopold Wentzel",
		Email: "leopold.wentzel@gmail.com",
	}

	// Endpoint must be authenticated
	req := &api.MemberInviteToken{}
	require.NoError(s.SetClientCSRFProtection(), "could not set csrf protection")
	err = s.client.InviteAccept(ctx, req)
	s.requireError(err, http.StatusUnauthorized, "this endpoint requires authentication", "expected error when user is not authenticated")

	// Should error if the token is missing
	require.NoError(s.SetClientCredentials(claims), "could not set client credentials")
	err = s.client.InviteAccept(ctx, req)
	s.requireError(err, http.StatusBadRequest, responses.ErrTryLoginAgain, "expected error when token is missing")

	// Successful invite accept request
	req.Token = "token1234"
	err = s.client.InviteAccept(ctx, req)
	require.NoError(err, "could not accept invite")
	token, err := s.GetClientAccessToken()
	require.NoError(err, "could not get access cookie token")
	require.Equal(accessToken, token, "expected access token to match")
	token, err = s.GetClientRefreshToken()
	require.NoError(err, "could not get refresh cookie token")
	require.Equal(refreshToken, token, "expected refresh token to match")

	// Test that an errors is returned if trtl returns an error
	trtl.OnGet = func(ctx context.Context, gr *pb.GetRequest) (*pb.GetReply, error) {
		return nil, status.Errorf(codes.NotFound, "member not found")
	}
	err = s.client.InviteAccept(ctx, req)
	s.requireError(err, http.StatusInternalServerError, responses.ErrSomethingWentWrong, "expected error when member is not found in trtl")

	// Test that if Quarterdeck returns an error Tenant returns an error
	s.quarterdeck.OnInvitesAccept(mock.UseError(http.StatusNotFound, "invalid token"), mock.RequireAuth())
	err = s.client.InviteAccept(ctx, req)
	s.requireError(err, http.StatusNotFound, "invalid token", "expected error when token is invalid")
}
