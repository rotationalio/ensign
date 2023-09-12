package tenant_test

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
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

	claims := &tokens.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: "01H6PGFB4T34D4WWEXQMAGJNMK",
		},
		Name:        "Kate Holland",
		Email:       "kate@example.co",
		Picture:     "https://www.gravatar.com/avatar/80ebb3b0dae3f550de72021bdcf45d00",
		OrgID:       "01H6PGFG71N0AFEVTK3NJB71T9",
		ProjectID:   "01H6PGFTK2X53RGG2KMSGR2M61",
		Permissions: []string{"read:foo", "edit:foo", "delete:foo"},
	}

	member := &db.Member{
		OrgID: claims.ParseOrgID(),
		ID:    claims.ParseUserID(),
		Email: claims.Email,
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
	creds := &qd.LoginReply{}
	creds.AccessToken, creds.RefreshToken, err = s.quarterdeck.CreateTokenPair(claims)
	require.NoError(err, "could not create access/refresh token pair")
	s.quarterdeck.OnInvitesAccept(mock.UseStatus(http.StatusOK), mock.UseJSONFixture(creds), mock.RequireAuth())

	// Set the initial claims fixture
	badClaims := &tokens.Claims{
		Name:  "Leopold Wentzel",
		Email: "leopold.wentzel@gmail.com",
	}

	// Endpoint must be authenticated
	req := &api.MemberInviteToken{}
	require.NoError(s.SetClientCSRFProtection(), "could not set csrf protection")
	err = s.client.InviteAccept(ctx, req)
	s.requireError(err, http.StatusUnauthorized, "this endpoint requires authentication", "expected error when user is not authenticated")

	// Should error if the token is missing
	require.NoError(s.SetClientCredentials(badClaims), "could not set client credentials")
	err = s.client.InviteAccept(ctx, req)
	s.requireError(err, http.StatusBadRequest, responses.ErrTryLoginAgain, "expected error when token is missing")

	// Successful invite accept request
	req.Token = "token1234"
	err = s.client.InviteAccept(ctx, req)
	require.NoError(err, "could not accept invite")
	token, err := s.GetClientAccessToken()
	require.NoError(err, "could not get access cookie token")
	require.Equal(creds.AccessToken, token, "expected access token to match")
	token, err = s.GetClientRefreshToken()
	require.NoError(err, "could not get refresh cookie token")
	require.Equal(creds.RefreshToken, token, "expected refresh token to match")

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
