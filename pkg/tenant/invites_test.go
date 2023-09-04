package tenant_test

import (
	"context"
	"net/http"
	"time"

	qd "github.com/rotationalio/ensign/pkg/quarterdeck/api/v1"
	"github.com/rotationalio/ensign/pkg/quarterdeck/mock"
	"github.com/rotationalio/ensign/pkg/quarterdeck/tokens"
	"github.com/rotationalio/ensign/pkg/tenant/api/v1"
	"github.com/rotationalio/ensign/pkg/utils/responses"
)

const (
	accessToken  = "eyJhbGciOiJSUzI1NiIsImtpZCI6IjAxR1g2NDdTOFBDVkJDUEpIWEdKUjI2UE42IiwidHlwIjoiSldUIn0.eyJpc3MiOiJodHRwOi8vMTI3LjAuMC4xIiwiYXVkIjpbImh0dHA6Ly8xMjcuMC4wLjEiXSwiZXhwIjoxNjgwNjE1MzMwLCJuYmYiOjE2ODA2MTE3MzAsImlhdCI6MTY4MDYxMTczMCwianRpIjoiMDFneDY0N3M4cGN2YmNwamh4Z2pzcG04N3AiLCJuYW1lIjoiSm9obiBEb2UiLCJlbWFpbCI6Impkb2VAZXhhbXBsZS5jb20iLCJvcmciOiIxMjMiLCJwcm9qZWN0IjoiYWJjIiwicGVybWlzc2lvbnMiOlsicmVhZDpkYXRhIiwid3JpdGU6ZGF0YSJdfQ.LLb6c2RdACJmoT3IFgJEwfu2_YJMcKgM2bF3ISF41A37gKTOkBaOe-UuTmjgZ7WEcuQ-cVkht0KI_4zqYYctB_WB9481XoNwff5VgFf3xrPdOYxS00YXQnl09RRqt6Fmca8nvd4mXfdO7uvpyNVuCIqNxBPXdSnRhreSoFB1GtFm42sBPAD7vF-MQUmU0c4PTsbiCfhR1_buH0NYEE1QFp3vYcgoiXOJHh9VStmRscqvLB12AQrcs26G9opdTCCORmvR2W3JLJ_hliHyp-d9lhXmCDFyiGkDEhTAUglqwBjqz5SO1UfAThWJO18PvZl4QPhb724oNT82VPh0DMDwfw"
	refreshToken = "eyJhbGciOiJSUzI1NiIsImtpZCI6IjAxR1g2NDdTOFBDVkJDUEpIWEdKUjI2UE42IiwidHlwIjoiSldUIn0.eyJpc3MiOiJodHRwOi8vMTI3LjAuMC4xIiwiYXVkIjpbImh0dHA6Ly8xMjcuMC4wLjEiLCJodHRwOi8vMTI3LjAuMC4xL3YxL3JlZnJlc2giXSwiZXhwIjoxNjgwNjE4OTMwLCJuYmYiOjE2ODA2MTQ0MzAsImlhdCI6MTY4MDYxMTczMCwianRpIjoiMDFneDY0N3M4cGN2YmNwamh4Z2pzcG04N3AifQ.CLHmtZwSPFCPoMBX06D_C3h3WuEonUbvbfWLvtmrMmIwnTwQ4hxsaRJo_a4qI-emp1HNg-yu_7c3VNwjkti-d0c7CAGApTaf5eRdGJ5HGUkI8RDHbbMFaOK86nAFnzdPJ2JLmGtLzvpF9eFXFllDhRiAB-2t0uKcOdN7cFghdwyWXIVJIJNjngF_WUFklmLKnqORtj_tA6UJ6NJnZln34eMGftAHbuH8x-xUiRePHnro4ydS43CKNOgRP8biMHiRR2broBz0apIt30TeQShaBSbmGx__LYdm7RKPJNVHAn_3h_PwwKQG567-Aqabg6TSmpwhXCk_RfUyQVGv2b997w"
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
	err := s.client.InviteAccept(ctx, req)
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

	// Test that if Quarterdeck returns an error Tenant returns an error
	s.quarterdeck.OnInvitesAccept(mock.UseError(http.StatusNotFound, "invalid token"), mock.RequireAuth())
	err = s.client.InviteAccept(ctx, req)
	s.requireError(err, http.StatusNotFound, "invalid token", "expected error when token is invalid")
}
