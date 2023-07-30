package api_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/rotationalio/ensign/pkg/quarterdeck/api/v1"
	"github.com/rotationalio/ensign/pkg/utils/ulids"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
)

func TestClient(t *testing.T) {
	// Create a Test Server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			require.Equal(t, int64(0), r.ContentLength)
			w.Header().Add("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, "{\"hello\":\"world\"}")
			return
		}

		require.Equal(t, int64(18), r.ContentLength)
		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "{\"error\":\"bad request\"}")
	}))
	defer ts.Close()

	// Create a Client that makes requests to the test server
	client, err := api.New(ts.URL)
	require.NoError(t, err)

	// Ensure that the latest version of the client is returned
	apiv1, ok := client.(*api.APIv1)
	require.True(t, ok)

	// Create a new GET request to a basic path
	req, err := apiv1.NewRequest(context.TODO(), http.MethodGet, "/foo", nil, nil)
	require.NoError(t, err)

	require.Equal(t, "/foo", req.URL.Path)
	require.Equal(t, "", req.URL.RawQuery)
	require.Equal(t, http.MethodGet, req.Method)
	require.Equal(t, "Quarterdeck API Client/v1", req.Header.Get("User-Agent"))
	require.Equal(t, "application/json", req.Header.Get("Accept"))
	require.Equal(t, "application/json; charset=utf-8", req.Header.Get("Content-Type"))

	// Create a new GET request with query params
	params := url.Values{}
	params.Add("q", "searching")
	params.Add("key", "open says me")
	req, err = apiv1.NewRequest(context.TODO(), http.MethodGet, "/foo", nil, &params)
	require.NoError(t, err)
	require.Equal(t, "key=open+says+me&q=searching", req.URL.RawQuery)

	data := make(map[string]string)
	rep, err := apiv1.Do(req, &data, true)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rep.StatusCode)
	require.Contains(t, data, "hello")
	require.Equal(t, "world", data["hello"])

	// Create a new POST request and check error handling
	req, err = apiv1.NewRequest(context.TODO(), http.MethodPost, "/bar", data, nil)
	require.NoError(t, err)
	rep, err = apiv1.Do(req, nil, false)
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, rep.StatusCode)

	req, err = apiv1.NewRequest(context.TODO(), http.MethodPost, "/bar", data, nil)
	require.NoError(t, err)
	_, err = apiv1.Do(req, nil, true)
	require.EqualError(t, err, "[400] bad request")

	// Test supplying an authentication override token in the request context
	ctx := api.ContextWithToken(context.Background(), "newtoken")
	req, err = apiv1.NewRequest(ctx, http.MethodPost, "/bar", data, nil)
	require.NoError(t, err, "could not create request")
	require.Equal(t, "Bearer newtoken", req.Header.Get("Authorization"), "expected the authorization header to be set")

	// Test that default credentials are used if no credentials are supplied in the request context
	defaultCreds := api.Token("default")
	client, err = api.New(ts.URL, api.WithCredentials(defaultCreds))
	require.NoError(t, err, "could not create client")
	apiv1, ok = client.(*api.APIv1)
	require.True(t, ok, "could not cast client to APIv1")
	req, err = apiv1.NewRequest(context.Background(), http.MethodPost, "/bar", data, nil)
	require.NoError(t, err, "could not create request")
	require.Equal(t, "Bearer default", req.Header.Get("Authorization"), "expected the authorization header to be set to default")

	// Test that request credentials override default credentials
	ctx = api.ContextWithToken(context.Background(), "newtoken")
	req, err = apiv1.NewRequest(ctx, http.MethodPost, "/bar", data, nil)
	require.NoError(t, err, "could not create request")
	require.Equal(t, "Bearer newtoken", req.Header.Get("Authorization"), "expected the authorization header to be set to newtoken")
}

//===========================================================================
// Client Methods
//===========================================================================

func TestStatus(t *testing.T) {
	t.Run("Ok", func(t *testing.T) {

		fixture := &api.StatusReply{
			Status:  "fine",
			Uptime:  (2 * time.Second).String(),
			Version: "1.0.test",
		}

		// Create a Test Server
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodGet, r.Method)
			require.Equal(t, "/v1/status", r.URL.Path)

			w.Header().Add("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(fixture)
		}))
		defer ts.Close()

		// Create a client to execute tests against the test server
		client, err := api.New(ts.URL)
		require.NoError(t, err)

		out, err := client.Status(context.Background())
		require.NoError(t, err, "could not execute status request")
		require.Equal(t, fixture, out, "expected the fixture to be returned")
	})

	t.Run("Unavailable", func(t *testing.T) {
		fixture := &api.StatusReply{
			Status:  "ack!",
			Uptime:  (9 * time.Second).String(),
			Version: "1.0.panic",
		}

		// Create a Test Server
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodGet, r.Method)
			require.Equal(t, "/v1/status", r.URL.Path)

			w.Header().Add("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusServiceUnavailable)
			json.NewEncoder(w).Encode(fixture)
		}))
		defer ts.Close()

		// Create a client to execute tests against the test server
		client, err := api.New(ts.URL)
		require.NoError(t, err)

		out, err := client.Status(context.Background())
		require.NoError(t, err, "could not execute status request")
		require.Equal(t, fixture, out, "expected the fixture to be returned")
	})
}

func TestRegister(t *testing.T) {
	// Setup the response fixture
	fixture := &api.RegisterReply{
		ID:      ulids.New(),
		OrgID:   ulids.New(),
		Email:   "jb@example.com",
		Message: "Thank you for registering for Ensign!",
		Role:    "Owner",
		Created: time.Now().Format(time.RFC3339Nano),
	}

	// Create a test server
	ts := httptest.NewServer(testhandler(fixture, http.MethodPost, "/v1/register"))
	defer ts.Close()

	// Create a client and execute endpoint request
	client, err := api.New(ts.URL)
	require.NoError(t, err, "could not create api client")

	req := &api.RegisterRequest{
		Name:         "Jane Bartholomew",
		Email:        "jb@example.com",
		Password:     "supers3cr4etsquir!!",
		PwCheck:      "supers3cr4etsquir!!",
		Organization: "Square",
		Domain:       "square",
		AgreeToS:     true,
		AgreePrivacy: true,
	}

	rep, err := client.Register(context.TODO(), req)
	require.NoError(t, err, "could not execute api request")
	require.Equal(t, fixture, rep, "unexpected response returned")
}

func TestLogin(t *testing.T) {
	// Setup the response fixture
	fixture := &api.LoginReply{}

	// Create a test server
	ts := httptest.NewServer(testhandler(fixture, http.MethodPost, "/v1/login"))
	defer ts.Close()

	// Create a client and execute endpoint request
	client, err := api.New(ts.URL)
	require.NoError(t, err, "could not create api client")

	req := &api.LoginRequest{}

	rep, err := client.Login(context.TODO(), req)
	require.NoError(t, err, "could not execute api request")
	require.Equal(t, fixture, rep, "unexpected response returned")
}

func TestAuthenticate(t *testing.T) {
	// Setup the response fixture
	fixture := &api.LoginReply{}

	// Create a test server
	ts := httptest.NewServer(testhandler(fixture, http.MethodPost, "/v1/authenticate"))
	defer ts.Close()

	// Create a client and execute endpoint request
	client, err := api.New(ts.URL)
	require.NoError(t, err, "could not create api client")

	req := &api.APIAuthentication{}

	rep, err := client.Authenticate(context.TODO(), req)
	require.NoError(t, err, "could not execute api request")
	require.Equal(t, fixture, rep, "unexpected response returned")
}

func TestRefresh(t *testing.T) {
	// Setup the response fixture
	fixture := &api.LoginReply{}

	// Create a test server
	ts := httptest.NewServer(testhandler(fixture, http.MethodPost, "/v1/refresh"))
	defer ts.Close()

	// Create a client and execute endpoint request
	client, err := api.New(ts.URL)
	require.NoError(t, err, "could not create api client")

	req := &api.RefreshRequest{}
	rep, err := client.Refresh(context.TODO(), req)
	require.NoError(t, err, "could not execute api request")
	require.Equal(t, fixture, rep, "unexpected response returned")
}

func TestSwitch(t *testing.T) {
	// Setup the response fixture
	fixture := &api.LoginReply{}

	// Create a test server
	ts := httptest.NewServer(testhandler(fixture, http.MethodPost, "/v1/switch"))
	defer ts.Close()

	// Create a client and execute endpoint request
	client, err := api.New(ts.URL)
	require.NoError(t, err, "could not create api client")

	req := &api.SwitchRequest{}
	rep, err := client.Switch(context.TODO(), req)
	require.NoError(t, err, "could not execute api request")
	require.Equal(t, fixture, rep, "unexpected response returned")
}

func TestVerifyEmail(t *testing.T) {
	// Setup the response fixture
	fixture := &api.Reply{
		Success: true,
	}

	// Create a test server
	ts := httptest.NewServer(testhandler(fixture, http.MethodPost, "/v1/verify"))
	defer ts.Close()

	// Create a client and execute endpoint request
	client, err := api.New(ts.URL)
	require.NoError(t, err, "could not create api client")

	req := &api.VerifyRequest{
		Token: "1234567890",
	}
	err = client.VerifyEmail(context.TODO(), req)
	require.NoError(t, err, "could not execute api request")
}

//===========================================================================
// Organization Resource
//===========================================================================

func TestOrganizationDetail(t *testing.T) {
	// Setup the response fixture
	fixture := &api.Organization{
		ID:     ulids.New(),
		Name:   "Events R Us",
		Domain: "events.io",
	}

	// Create a test server
	ts := httptest.NewServer(testhandler(fixture, http.MethodGet, fmt.Sprintf("/v1/organizations/%s", fixture.ID.String())))
	defer ts.Close()

	// Create a client and execute endpoint request
	client, err := api.New(ts.URL)
	require.NoError(t, err, "could not create api client")

	rep, err := client.OrganizationDetail(context.TODO(), fixture.ID.String())
	require.NoError(t, err, "could not execute api request")
	require.Equal(t, fixture, rep, "unexpected response returned")
}

func TestOrganizationList(t *testing.T) {
	// Setup the response fixture
	fixture := &api.OrganizationList{}

	// Create a test server
	ts := httptest.NewServer(testhandler(fixture, http.MethodGet, "/v1/organizations"))
	defer ts.Close()

	// Create a client and execute endpoint request
	client, err := api.New(ts.URL)
	require.NoError(t, err, "could not create api client")

	req := &api.OrganizationPageQuery{}
	rep, err := client.OrganizationList(context.TODO(), req)
	require.NoError(t, err, "could not execute api request")
	require.Equal(t, fixture, rep, "unexpected response returned")
}

//===========================================================================
// API Keys Resource
//===========================================================================

func TestAPIKeyList(t *testing.T) {
	// Setup the response fixture
	fixture := &api.APIKeyList{}

	// Create a test server
	ts := httptest.NewServer(testhandler(fixture, http.MethodGet, "/v1/apikeys"))
	defer ts.Close()

	// Create a client and execute endpoint request
	client, err := api.New(ts.URL)
	require.NoError(t, err, "could not create api client")

	req := &api.APIPageQuery{}
	rep, err := client.APIKeyList(context.TODO(), req)
	require.NoError(t, err, "could not execute api request")
	require.Equal(t, fixture, rep, "unexpected response returned")
}

func TestAPIKeyCreate(t *testing.T) {
	// Setup the response fixture
	fixture := &api.APIKey{}

	// Create a test server
	ts := httptest.NewServer(testhandler(fixture, http.MethodPost, "/v1/apikeys"))
	defer ts.Close()

	// Create a client and execute endpoint request
	client, err := api.New(ts.URL)
	require.NoError(t, err, "could not create api client")

	req := &api.APIKey{}
	rep, err := client.APIKeyCreate(context.TODO(), req)
	require.NoError(t, err, "could not execute api request")
	require.Equal(t, fixture, rep, "unexpected response returned")
}

func TestAPIKeyDetail(t *testing.T) {
	// Setup the response fixture
	fixture := &api.APIKey{}

	// Create a test server
	ts := httptest.NewServer(testhandler(fixture, http.MethodGet, "/v1/apikeys/foo"))
	defer ts.Close()

	// Create a client and execute endpoint request
	client, err := api.New(ts.URL)
	require.NoError(t, err, "could not create api client")

	rep, err := client.APIKeyDetail(context.TODO(), "foo")
	require.NoError(t, err, "could not execute api request")
	require.Equal(t, fixture, rep, "unexpected response returned")
}

func TestAPIKeyUpdate(t *testing.T) {
	// Setup the response fixture
	kid := ulids.New()
	fixture := &api.APIKey{
		ID: kid,
	}

	// Create a test server
	ts := httptest.NewServer(testhandler(fixture, http.MethodPut, fmt.Sprintf("/v1/apikeys/%s", kid.String())))
	defer ts.Close()

	// Create a client and execute endpoint request
	client, err := api.New(ts.URL)
	require.NoError(t, err, "could not create api client")

	req := &api.APIKey{
		ID: kid,
	}

	rep, err := client.APIKeyUpdate(context.TODO(), req)
	require.NoError(t, err, "could not execute api request")
	require.Equal(t, fixture, rep, "unexpected response returned")
}

func TestAPIKeyDelete(t *testing.T) {
	// Setup the response fixture
	fixture := &api.Reply{Success: true}

	// Create a test server
	ts := httptest.NewServer(testhandler(fixture, http.MethodDelete, "/v1/apikeys/foo"))
	defer ts.Close()

	// Create a client and execute endpoint request
	client, err := api.New(ts.URL)
	require.NoError(t, err, "could not create api client")

	err = client.APIKeyDelete(context.TODO(), "foo")
	require.NoError(t, err, "could not execute api request")
}

func TestAPIKeyPermissions(t *testing.T) {
	// Setup the response fixture
	fixture := []string{"topics:read", "topics:destroy"}

	// Create a test server
	ts := httptest.NewServer(testhandler(fixture, http.MethodGet, "/v1/apikeys/permissions"))
	defer ts.Close()

	// Create a client and execute endpoint request
	client, err := api.New(ts.URL)
	require.NoError(t, err, "could not create api client")

	rep, err := client.APIKeyPermissions(context.TODO())
	require.NoError(t, err, "could not execute api request")
	require.Equal(t, fixture, rep, "unexpected response returned")
}

//===========================================================================
// Project Resource
//===========================================================================

func TestProjectList(t *testing.T) {
	// Setup the response fixture
	dts, _ := time.Parse(time.RFC3339Nano, time.Now().Format(time.RFC3339Nano))

	fixture := &api.ProjectList{
		Projects: []*api.Project{
			{
				OrgID:     ulids.New(),
				ProjectID: ulids.New(),
				Created:   dts,
				Modified:  dts,
			},
		},
	}

	// Create a test server
	ts := httptest.NewServer(testhandler(fixture, http.MethodGet, "/v1/projects"))
	defer ts.Close()

	// Create a client and execute endpoint request
	client, err := api.New(ts.URL)
	require.NoError(t, err, "could not create api client")

	rep, err := client.ProjectList(context.Background(), &api.PageQuery{})
	require.NoError(t, err, "could not execute api request")
	require.Equal(t, fixture, rep, "unexpected response returned")
}

func TestProjectCreate(t *testing.T) {
	// Setup the response fixture
	dts, _ := time.Parse(time.RFC3339Nano, time.Now().Format(time.RFC3339Nano))

	fixture := &api.Project{
		OrgID:     ulids.New(),
		ProjectID: ulids.New(),
		Created:   dts,
		Modified:  dts,
	}

	// Create a test server
	ts := httptest.NewServer(testhandler(fixture, http.MethodPost, "/v1/projects"))
	defer ts.Close()

	// Create a client and execute endpoint request
	client, err := api.New(ts.URL)
	require.NoError(t, err, "could not create api client")

	req := &api.Project{
		ProjectID: fixture.ProjectID,
	}
	rep, err := client.ProjectCreate(context.Background(), req)
	require.NoError(t, err, "could not execute api request")
	require.Equal(t, fixture, rep, "unexpected response returned")
}

func TestProjectAccess(t *testing.T) {
	// Setup the response fixture
	dts, _ := time.Parse(time.RFC3339Nano, time.Now().Format(time.RFC3339Nano))

	fixture := &api.Project{
		OrgID:     ulids.New(),
		ProjectID: ulids.New(),
		Created:   dts,
		Modified:  dts,
	}

	// Create a test server
	ts := httptest.NewServer(testhandler(fixture, http.MethodGet, "/v1/projects/abcd1234"))
	defer ts.Close()

	// Create a client and execute endpoint request
	client, err := api.New(ts.URL)
	require.NoError(t, err, "could not create api client")

	rep, err := client.ProjectDetail(context.Background(), "abcd1234")
	require.NoError(t, err, "could not execute api request")
	require.Equal(t, fixture, rep, "unexpected response returned")
}

func TestProjectDetail(t *testing.T) {
	// Setup the response fixture
	fixture := &api.LoginReply{
		AccessToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
	}

	// Create a test server
	ts := httptest.NewServer(testhandler(fixture, http.MethodPost, "/v1/projects/access"))
	defer ts.Close()

	// Create a client and execute endpoint request
	client, err := api.New(ts.URL)
	require.NoError(t, err, "could not create api client")

	req := &api.Project{
		ProjectID: ulids.New(),
	}
	rep, err := client.ProjectAccess(context.Background(), req)
	require.NoError(t, err, "could not execute api request")
	require.Equal(t, fixture, rep, "unexpected response returned")
}

//===========================================================================
// Users Resource
//===========================================================================

func TestUserDetail(t *testing.T) {
	// Setup the response fixture
	fixture := &api.User{}

	// Create a test server
	ts := httptest.NewServer(testhandler(fixture, http.MethodGet, "/v1/users/foo"))
	defer ts.Close()

	// Create a client and execute endpoint request
	client, err := api.New(ts.URL)
	require.NoError(t, err, "could not create api client")

	rep, err := client.UserDetail(context.TODO(), "foo")
	require.NoError(t, err, "could not execute api request")
	require.Equal(t, fixture, rep, "unexpected response returned")
}

func TestUserUpdate(t *testing.T) {
	// Setup the response fixture
	userID := ulids.New()
	fixture := &api.User{
		UserID: userID,
	}

	// Create a test server
	ts := httptest.NewServer(testhandler(fixture, http.MethodPut, fmt.Sprintf("/v1/users/%s", userID.String())))
	defer ts.Close()

	// Create a client and execute endpoint request
	client, err := api.New(ts.URL)
	require.NoError(t, err, "could not create api client")

	req := &api.User{
		UserID: userID,
		Name:   "Joan Miller",
	}

	rep, err := client.UserUpdate(context.TODO(), req)
	require.NoError(t, err, "could not execute api request")
	require.Equal(t, fixture, rep, "unexpected response returned")
}

func TestUserRoleUpdate(t *testing.T) {
	// Setup the response fixture
	userID := ulids.New()
	fixture := &api.User{
		UserID: userID,
		Name:   "Joan Miller",
		Role:   "Admin",
	}

	// Create a test server
	ts := httptest.NewServer(testhandler(fixture, http.MethodPost, fmt.Sprintf("/v1/users/%s", userID.String())))
	defer ts.Close()

	// Create a client and execute endpoint request
	client, err := api.New(ts.URL)
	require.NoError(t, err, "could not create api client")

	req := &api.UpdateRoleRequest{
		ID:   userID,
		Role: "Admin",
	}

	rep, err := client.UserRoleUpdate(context.TODO(), req)
	require.NoError(t, err, "could not execute api request")
	require.Equal(t, fixture, rep, "unexpected response returned")
}

func TestUserList(t *testing.T) {
	// Setup the response fixture
	fixture := &api.UserList{}

	// Create a test server
	ts := httptest.NewServer(testhandler(fixture, http.MethodGet, "/v1/users"))
	defer ts.Close()

	// Create a client and execute endpoint request
	client, err := api.New(ts.URL)
	require.NoError(t, err, "could not create api client")

	req := &api.UserPageQuery{}
	rep, err := client.UserList(context.TODO(), req)
	require.NoError(t, err, "could not execute api request")
	require.Equal(t, fixture, rep, "unexpected response returned")
}

func TestUserRemove(t *testing.T) {
	// Setup the response fixture
	fixture := &api.UserRemoveReply{
		APIKeys: []string{"bar", "baz"},
		Token:   "foo",
	}

	// Create a test server
	ts := httptest.NewServer(testhandler(fixture, http.MethodDelete, "/v1/users/foo"))
	defer ts.Close()

	// Create a client and execute endpoint request
	client, err := api.New(ts.URL)
	require.NoError(t, err, "could not create api client")

	out, err := client.UserRemove(context.TODO(), "foo")
	require.NoError(t, err, "could not execute api request")
	require.Equal(t, fixture, out, "unexpected response returned")
}

func TestUserRemoveConfirm(t *testing.T) {
	id := ulids.New()

	// Create a test server
	ts := httptest.NewServer(testhandler(nil, http.MethodDelete, "/v1/users/"+id.String()+"/confirm"))
	defer ts.Close()

	// Create a client and execute endpoint request
	client, err := api.New(ts.URL)
	require.NoError(t, err, "could not create api client")

	req := &api.UserRemoveConfirm{
		ID:    id,
		Token: "foo",
	}
	err = client.UserRemoveConfirm(context.TODO(), req)
	require.NoError(t, err, "could not execute api request")
}

func TestInvitePreview(t *testing.T) {
	// Setup the response fixture
	fixture := &api.UserInvitePreview{
		OrgName:     "Acme Inc.",
		InviterName: "John Doe",
		Role:        "Member",
		UserExists:  true,
	}

	// Create a test server
	ts := httptest.NewServer(testhandler(fixture, http.MethodGet, "/v1/invites/foo"))
	defer ts.Close()

	// Create a client and execute endpoint request
	client, err := api.New(ts.URL)
	require.NoError(t, err, "could not create api client")

	rep, err := client.InvitePreview(context.TODO(), "foo")
	require.NoError(t, err, "could not execute api request")
	require.Equal(t, fixture, rep, "unexpected response returned")
}

func TestInviteCreate(t *testing.T) {
	// Setup the response fixture
	fixture := &api.UserInviteReply{
		UserID:    ulids.New(),
		OrgID:     ulids.New(),
		Email:     "leopold.wentzel@gmail.com",
		Role:      "admin",
		CreatedBy: ulids.New(),
	}

	// Create a test server
	ts := httptest.NewServer(testhandler(fixture, http.MethodPost, "/v1/invites"))
	defer ts.Close()

	// Create a client and execute endpoint request
	client, err := api.New(ts.URL)
	require.NoError(t, err, "could not create api client")

	req := &api.UserInviteRequest{
		Email: "leopold.wentzel@gmail.com",
		Role:  "admin",
	}
	reply, err := client.InviteCreate(context.TODO(), req)
	require.NoError(t, err, "could not execute api request")
	require.Equal(t, fixture, reply, "unexpected response returned")
}

//===========================================================================
// Accounts Resource
//===========================================================================

func TestAccountUpdate(t *testing.T) {
	// Setup the response fixture
	userID := ulids.New()
	fixture := &api.User{
		UserID: userID,
	}

	// Create a test server
	ts := httptest.NewServer(testhandler(fixture, http.MethodPut, fmt.Sprintf("/v1/accounts/%s", userID.String())))
	defer ts.Close()

	// Create a client and execute endpoint request
	client, err := api.New(ts.URL)
	require.NoError(t, err, "could not create api client")

	req := &api.User{
		UserID: userID,
		Name:   "Joan Miller",
	}

	rep, err := client.AccountUpdate(context.TODO(), req)
	require.NoError(t, err, "could not execute api request")
	require.Equal(t, fixture, rep, "unexpected response returned")
}

func TestWaitForReady(t *testing.T) {
	// This is a long running test, skip if in short mode
	if testing.Short() {
		t.Skip("skipping long running test in short mode")
	}

	// Backoff Interval should be as follows:
	// Request #   Retry Interval (sec)     Randomized Interval (sec)
	//   1          0.5                     [0.25,   0.75]
	//   2          0.75                    [0.375,  1.125]
	//   3          1.125                   [0.562,  1.687]
	//   4          1.687                   [0.8435, 2.53]
	//   5          2.53                    [1.265,  3.795]
	//   6          3.795                   [1.897,  5.692]
	//   7          5.692                   [2.846,  8.538]
	//   8          8.538                   [4.269, 12.807]
	//   9         12.807                   [6.403, 19.210]
	//  10         19.210                   backoff.Stop

	fixture := &api.StatusReply{
		Version: "1.0.test",
	}

	// Create a Test Server
	tries := 0
	started := time.Now()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/status" {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		tries++
		var status int
		if tries < 4 {
			status = http.StatusServiceUnavailable
			fixture.Status = "maintenance"
		} else {
			status = http.StatusOK
			fixture.Status = "fine"
		}

		fixture.Uptime = time.Since(started).String()

		log.Info().Int("status", status).Int("tries", tries).Msg("responding to status request")
		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(fixture)
	}))
	defer ts.Close()

	// Create a client to execute tests against the test server
	client, err := api.New(ts.URL)
	require.NoError(t, err)

	// We expect it takes 5 tries before a good response is returned that means that
	// the minimum delay according to the above table is 1.187 seconds
	err = client.WaitForReady(context.Background())
	require.NoError(t, err)
	require.GreaterOrEqual(t, time.Since(started), 1187*time.Millisecond)

	// Should not have any wait since the test server will respond true
	started = time.Now()
	err = client.WaitForReady(context.Background())
	require.NoError(t, err)
	require.LessOrEqual(t, time.Since(started), 250*time.Millisecond)

	// Test timeout
	tries = 0
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	err = client.WaitForReady(ctx)
	require.ErrorIs(t, err, context.DeadlineExceeded)
}

//===========================================================================
// Helper Methods
//===========================================================================

func testhandler(fixture interface{}, expectedMethod, expectedPath string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json; charset=utf-8")

		if r.Method != expectedMethod {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(api.ErrorResponse("unexpected http method"))
			return
		}

		if r.URL.Path != expectedPath {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(api.ErrorResponse("unexpected endpoint path"))
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(fixture)
	})
}
