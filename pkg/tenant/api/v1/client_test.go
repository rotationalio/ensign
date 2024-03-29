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

	"github.com/rotationalio/ensign/pkg/quarterdeck/middleware"
	"github.com/rotationalio/ensign/pkg/quarterdeck/permissions"
	"github.com/rotationalio/ensign/pkg/tenant/api/v1"
	"github.com/stretchr/testify/require"
)

func TestClient(t *testing.T) {
	// Creates a Test Server
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

	// Creates a Client that makes requests to the test server
	client, err := api.New(ts.URL)
	require.NoError(t, err)

	// Ensures that the latest version of the client is returned
	apiv1, ok := client.(*api.APIv1)
	require.True(t, ok)

	// Creates a new GET request to a basic path
	req, err := apiv1.NewRequest(context.TODO(), http.MethodGet, "/foo", nil, nil)
	require.NoError(t, err)

	require.Equal(t, "/foo", req.URL.Path)
	require.Equal(t, "", req.URL.RawQuery)
	require.Equal(t, http.MethodGet, req.Method)
	require.Equal(t, "Tenant API Client/v1", req.Header.Get("User-Agent"))
	require.Equal(t, "application/json", req.Header.Get("Accept"))
	require.Equal(t, "application/json; charset=utf-8", req.Header.Get("Content-Type"))

	// Creates a new GET request with query params
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

	// Creates a new POST request and checks error handling
	req, err = apiv1.NewRequest(context.TODO(), http.MethodPost, "/bar", data, nil)
	require.NoError(t, err)
	rep, err = apiv1.Do(req, nil, false)
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, rep.StatusCode)
}

func TestStatus(t *testing.T) {
	t.Run("Ok", func(t *testing.T) {
		fixture := &api.StatusReply{
			Status:  "fine",
			Uptime:  (2 * time.Second).String(),
			Version: "1.0.test",
		}

		// Creates a Test Server
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodGet, r.Method)
			require.Equal(t, "/v1/status", r.URL.Path)

			w.Header().Add("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(fixture)

		}))
		defer ts.Close()

		// Creates a client to execute tests against the test server
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

		// Creates a Test Server
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodGet, r.Method)
			require.Equal(t, "/v1/status", r.URL.Path)

			w.Header().Add("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(fixture)

		}))
		defer ts.Close()

		// Creates a client to execute tests against the test server
		client, err := api.New(ts.URL)
		require.NoError(t, err)

		out, err := client.Status(context.Background())
		require.NoError(t, err, "could not execute status request")
		require.Equal(t, fixture, out, "expected the fixture to be returned")
	})
}

func TestRegister(t *testing.T) {
	// Create a test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/v1/register", r.URL.Path)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	// Create a client to execute tests against the test server
	client, err := api.New(ts.URL)
	require.NoError(t, err, "could not create client")

	// Create a new register request
	req := &api.RegisterRequest{
		Email:    "leopold.wentzel@gmail.com",
		Password: "hunter2",
		PwCheck:  "hunter2",
	}
	err = client.Register(context.Background(), req)
	require.NoError(t, err, "could not execute register request")
}

func TestLogin(t *testing.T) {
	fixture := &api.AuthReply{
		AccessToken:  "access",
		RefreshToken: "refresh",
	}

	// Create a test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/v1/login", r.URL.Path)

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(fixture)
	}))
	defer ts.Close()

	// Create a client to execute tests against the test server
	client, err := api.New(ts.URL)
	require.NoError(t, err, "could not create client")

	// Create a new login request
	req := &api.LoginRequest{
		Email:    "leopold.wentzel@gmail.com",
		Password: "hunter2",
	}
	rep, err := client.Login(context.Background(), req)
	require.NoError(t, err, "could not execute login request")
	require.Equal(t, fixture, rep, "expected the fixture to be returned")
}

func TestRefresh(t *testing.T) {
	fixture := &api.AuthReply{
		AccessToken:  "access",
		RefreshToken: "refresh",
	}

	// Create a test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/v1/refresh", r.URL.Path)

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(fixture)
	}))
	defer ts.Close()

	// Create a client to execute tests against the test server
	client, err := api.New(ts.URL)
	require.NoError(t, err, "could not create client")

	// Create a new refresh request
	req := &api.RefreshRequest{
		RefreshToken: "refresh",
	}
	rep, err := client.Refresh(context.Background(), req)
	require.NoError(t, err, "could not execute refresh request")
	require.Equal(t, fixture, rep, "expected the fixture to be returned")
}

func TestSwitch(t *testing.T) {
	fixture := &api.AuthReply{
		AccessToken:  "access",
		RefreshToken: "refresh",
	}

	// Create a test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/v1/switch", r.URL.Path)

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(fixture)
	}))
	defer ts.Close()

	// Create a client to execute tests against the test server
	client, err := api.New(ts.URL)
	require.NoError(t, err, "could not create client")

	// Create a new switch request
	req := &api.SwitchRequest{
		OrgID: "001",
	}
	rep, err := client.Switch(context.Background(), req)
	require.NoError(t, err, "could not execute switch request")
	require.Equal(t, fixture, rep, "expected the fixture to be returned")
}

func TestVerifyEmail(t *testing.T) {
	fixture := &api.AuthReply{
		AccessToken:  "access",
		RefreshToken: "refresh",
	}

	// Create a test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/v1/verify", r.URL.Path)

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(fixture)
	}))
	defer ts.Close()

	// Create a client to execute tests against the test server
	client, err := api.New(ts.URL)
	require.NoError(t, err, "could not create client")

	// Create a new verify request
	req := &api.VerifyRequest{
		Token: "token",
	}
	rep, err := client.VerifyEmail(context.Background(), req)
	require.NoError(t, err, "could not execute verify request")
	require.Equal(t, fixture, rep, "expected the fixture to be returned")
}

func TestResendEmail(t *testing.T) {
	fixture := &api.Reply{}

	// Create a test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/v1/resend", r.URL.Path)

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(fixture)
	}))
	defer ts.Close()

	// Create a client to execute tests against the test server
	client, err := api.New(ts.URL)
	require.NoError(t, err, "could not create client")

	// Create a new resend request
	req := &api.ResendRequest{
		Email: "leopold.wentzel@gmail.com",
		OrgID: "001",
	}
	err = client.ResendEmail(context.Background(), req)
	require.NoError(t, err, "could not execute resend request")
}

func TestForgotPassword(t *testing.T) {
	// Create a test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/v1/forgot-password", r.URL.Path)

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	// Create a client to execute tests against the test server
	client, err := api.New(ts.URL)
	require.NoError(t, err)

	// Create a new forgot password request
	req := &api.ForgotPasswordRequest{
		Email: "leopold.wentzel@gmail.com",
	}
	err = client.ForgotPassword(context.Background(), req)
	require.NoError(t, err, "could not execute forgot password request")
}

func TestResetPassword(t *testing.T) {
	// Create a test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/v1/reset-password", r.URL.Path)

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	// Create a client to execute tests against the test server
	client, err := api.New(ts.URL)
	require.NoError(t, err)

	// Create a new reset password request
	req := &api.ResetPasswordRequest{
		Token:    "token",
		Password: "password",
		PwCheck:  "password",
	}
	err = client.ResetPassword(context.Background(), req)
	require.NoError(t, err, "could not execute reset password request")
}

func TestInvitePreview(t *testing.T) {
	fixture := &api.MemberInvitePreview{
		Email:       "leopold.wentzel@checkers.io",
		OrgName:     "Checkers",
		InviterName: "Alice Smith",
		Role:        "Member",
		HasAccount:  true,
	}

	// Create a test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/v1/invites/1234", r.URL.Path)

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(fixture)
	}))
	defer ts.Close()

	// Create a client to execute tests against the test server
	client, err := api.New(ts.URL)
	require.NoError(t, err, "could not create client")

	// Do the request
	out, err := client.InvitePreview(context.Background(), "1234")
	require.NoError(t, err, "could not execute invite preview request")
	require.Equal(t, fixture, out, "expected the fixture to be returned")
}

func TestInviteAccept(t *testing.T) {
	fixture := &api.AuthReply{
		AccessToken:  "access",
		RefreshToken: "refresh",
	}

	// Create a test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/v1/invites/accept", r.URL.Path)

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(fixture)
	}))
	defer ts.Close()

	// Create a client to execute tests against the test server
	client, err := api.New(ts.URL)
	require.NoError(t, err, "could not create client")

	req := &api.MemberInviteToken{
		Token: "token",
	}
	rep, err := client.InviteAccept(context.Background(), req)
	require.NoError(t, err, "could not execute invite accept request")
	require.Equal(t, fixture, rep, "expected the fixture to be returned")
}

func TestOrganizationList(t *testing.T) {
	fixture := &api.OrganizationPage{
		Organizations: []*api.Organization{
			{
				ID:        "001",
				Name:      "Events R Us",
				Domain:    "events.io",
				Created:   "2023-02-06T13:59:16-06:00",
				LastLogin: "2023-02-07T13:59:16-06:00",
			},
			{
				ID:        "002",
				Name:      "Rotational Labs",
				Domain:    "rotational.io",
				Created:   "2023-02-06T13:59:16-06:00",
				LastLogin: "2023-02-07T13:59:16-06:00",
			},
		},
	}

	// Create a test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/v1/organization", r.URL.Path)

		params := r.URL.Query()
		require.Equal(t, "2", params.Get("page_size"))
		require.Equal(t, "12", params.Get("next_page_token"))

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(fixture)
	}))
	defer ts.Close()

	// Create a client to execute tests against the test server
	client, err := api.New(ts.URL)
	require.NoError(t, err, "could not create api client")

	req := &api.PageQuery{
		PageSize:      2,
		NextPageToken: "12",
	}

	out, err := client.OrganizationList(context.Background(), req)
	require.NoError(t, err, "could not execute api request")
	require.Equal(t, fixture, out, "unexpected response error")
}

func TestOrganization(t *testing.T) {
	fixture := &api.Organization{
		ID:       "001",
		Name:     "Events R Us",
		Domain:   "events.io",
		Created:  "2023-02-06T13:59:16-06:00",
		Modified: "2023-02-07T13:59:16-06:00",
	}

	// Create a test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/v1/organization/001", r.URL.Path)

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(fixture)
	}))
	defer ts.Close()

	// Create a client to execute tests against the test server
	client, err := api.New(ts.URL)
	require.NoError(t, err, "could not create client")

	// Execute the request
	out, err := client.OrganizationDetail(context.Background(), "001")
	require.NoError(t, err, "could not execute organization request")
	require.Equal(t, fixture, out, "expected the fixture to be returned")
}

func TestTenantList(t *testing.T) {
	fixture := &api.TenantPage{
		Tenants: []*api.Tenant{
			{
				ID:              "001",
				Name:            "tenant01",
				EnvironmentType: "Dev",
			},
			{
				ID:              "002",
				Name:            "tenant02",
				EnvironmentType: "Prod",
			},
			{
				ID:              "003",
				Name:            "tenant03",
				EnvironmentType: "Stage",
			},
		},
		NextPageToken: "23",
	}

	// Creates a test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/v1/tenant", r.URL.Path)

		params := r.URL.Query()
		require.Equal(t, "2", params.Get("page_size"))
		require.Equal(t, "12", params.Get("next_page_token"))

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(fixture)
	}))
	defer ts.Close()

	// Creates a client to execute tests against the server
	client, err := api.New(ts.URL)
	require.NoError(t, err, "could not create api client")

	req := &api.PageQuery{
		PageSize:      2,
		NextPageToken: "12",
	}

	out, err := client.TenantList(context.Background(), req)
	require.NoError(t, err, "could not execute api request")
	require.Equal(t, fixture, out, "unexpected response error")
}

func TestTenantCreate(t *testing.T) {
	fixture := &api.Tenant{
		ID:              "1234",
		Name:            "feist",
		EnvironmentType: "Dev",
	}

	// Creates a test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/v1/tenant", r.URL.Path)

		in := &api.Tenant{}
		err := json.NewDecoder(r.Body).Decode(in)
		require.NoError(t, err, "could not decode request")

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(fixture)
	}))
	defer ts.Close()

	// Creates a client to execute tests against the test server
	client, err := api.New(ts.URL)
	require.NoError(t, err, "could not create api client")

	out, err := client.TenantCreate(context.Background(), &api.Tenant{})
	require.NoError(t, err, "could not execute api request")
	require.Equal(t, fixture, out, "unexpected response error")
}

func TestTenantDetail(t *testing.T) {
	fixture := &api.Tenant{
		ID:              "001",
		Name:            "tenant01",
		EnvironmentType: "Dev",
	}

	// Creates a test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/v1/tenant/tenant01", r.URL.Path)

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(fixture)
	}))
	defer ts.Close()

	// Creates a client to execute tests against the test server
	client, err := api.New(ts.URL)
	require.NoError(t, err, "could not create api client")

	out, err := client.TenantDetail(context.Background(), "tenant01")
	require.NoError(t, err, "could not execute api request")
	require.Equal(t, fixture, out, "unexpected response error")
}

func TestTenantUpdate(t *testing.T) {
	fixture := &api.Tenant{
		ID:              "001",
		Name:            "tenant01",
		EnvironmentType: "Dev",
	}

	// Creates a test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPut, r.Method)
		require.Equal(t, "/v1/tenant/001", r.URL.Path)

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(fixture)
	}))
	defer ts.Close()

	// Creates a client to execute tests against the server
	client, err := api.New(ts.URL)
	require.NoError(t, err, "could not execute api request")

	req := &api.Tenant{
		ID:              "001",
		Name:            "tenant02",
		EnvironmentType: "Prod",
	}

	rep, err := client.TenantUpdate(context.Background(), req)
	require.NoError(t, err, "could not execute api request")
	require.Equal(t, fixture, rep, "unexpected response error")
}

func TestTenantDelete(t *testing.T) {
	fixture := &api.Reply{}

	// Creates a test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodDelete, r.Method)
		require.Equal(t, "/v1/tenant/tenant01", r.URL.Path)

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(fixture)
	}))
	defer ts.Close()

	// Creates a client to execute tests against the test server
	client, err := api.New(ts.URL)
	require.NoError(t, err, "could not create api client")

	err = client.TenantDelete(context.Background(), "tenant01")
	require.NoError(t, err, "could not execute api request")
}

func TestTenantStats(t *testing.T) {
	fixture := []*api.StatValue{
		{
			Name:  "projects",
			Value: 10,
		},
		{
			Name:  "topics",
			Value: 5,
		},
		{
			Name:  "keys",
			Value: 3,
		},
	}

	// Creates a test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/v1/tenant/002/stats", r.URL.Path)

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(fixture)
	}))
	defer ts.Close()

	// Create a client to execute tests against the test server
	client, err := api.New(ts.URL)
	require.NoError(t, err, "could not create api client")

	out, err := client.TenantStats(context.Background(), "002")
	require.NoError(t, err, "could not execute api request")
	require.Equal(t, fixture, out, "unexpected response body")
}

func TestMemberList(t *testing.T) {
	fixture := &api.MemberPage{
		Members: []*api.Member{
			{
				ID:   "002",
				Name: "Ryan Moore",
				Role: "Admin",
			},
		},
		NextPageToken: "23",
	}
	// Creates a test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/v1/members", r.URL.Path)

		params := r.URL.Query()
		require.Equal(t, "2", params.Get("page_size"))
		require.Equal(t, "12", params.Get("next_page_token"))

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(fixture)
	}))
	defer ts.Close()

	// Creates a client to execute tests against the test server
	client, err := api.New(ts.URL)
	require.NoError(t, err, "could not create api client")

	req := &api.PageQuery{
		PageSize:      2,
		NextPageToken: "12",
	}

	out, err := client.MemberList(context.Background(), req)
	require.NoError(t, err, "could not execute api request")
	require.Equal(t, fixture, out, "unexpected response error")
}

func TestMemberCreate(t *testing.T) {
	fixture := &api.Member{
		ID:   "002",
		Name: "Ryan Moore",
		Role: "Admin",
	}

	// Creates a test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/v1/members", r.URL.Path)

		in := &api.Member{}
		err := json.NewDecoder(r.Body).Decode(in)
		require.NoError(t, err, "could not decode request")

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(fixture)
	}))
	defer ts.Close()

	// Create a client to execute tests against the test server
	client, err := api.New(ts.URL)
	require.NoError(t, err, "could not create api client")

	req := &api.Member{}

	out, err := client.MemberCreate(context.Background(), req)
	require.NoError(t, err, "could not execute api request")
	require.Equal(t, fixture, out, "unexpected response error")
}

func TestMemberDetail(t *testing.T) {
	fixture := &api.Member{
		ID:   "001",
		Name: "Luke Hamilton",
		Role: "Admin",
	}

	// Creates a test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/v1/members/member001", r.URL.Path)

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(fixture)
	}))
	defer ts.Close()

	// Creates a client to execute tests against the test server
	client, err := api.New(ts.URL)
	require.NoError(t, err, "could not create api client")

	out, err := client.MemberDetail(context.Background(), "member001")
	require.NoError(t, err, "could not execute api request")
	require.Equal(t, fixture, out, "unexpected response error")
}

func TestMemberUpdate(t *testing.T) {
	fixture := &api.Member{
		ID:   "001",
		Name: "member01",
		Role: "Admin",
	}

	// Creates a test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPut, r.Method)
		require.Equal(t, "/v1/members/001", r.URL.Path)

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(fixture)
	}))
	defer ts.Close()

	// Creates a client to execute tests against the server
	client, err := api.New(ts.URL)
	require.NoError(t, err, "could not execute api request")

	req := &api.Member{
		ID:   "001",
		Name: "member02",
		Role: "Admin",
	}

	rep, err := client.MemberUpdate(context.Background(), req)
	require.NoError(t, err, "could not execute api request")
	require.Equal(t, fixture, rep, "unexpected response error")
}

func TestMemberRoleUpdate(t *testing.T) {
	fixture := &api.Member{
		Role: permissions.RoleObserver,
	}

	// Create a test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/v1/members/member001", r.URL.Path)

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(fixture)
	}))
	defer ts.Close()

	// Create a client to execute test against the test server
	client, err := api.New(ts.URL)
	require.NoError(t, err, "could not create api client")

	req := &api.UpdateRoleParams{
		Role: permissions.RoleAdmin,
	}

	out, err := client.MemberRoleUpdate(context.Background(), "member001", req)
	require.NoError(t, err, "could not execute api request")
	require.Equal(t, fixture, out, "expected member role to update")
}

func TestMemberDelete(t *testing.T) {
	fixture := &api.MemberDeleteReply{
		APIKeys: []string{"key001", "key002"},
		Token:   "token001",
	}

	// Creates a test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodDelete, r.Method)
		require.Equal(t, "/v1/members/member001", r.URL.Path)

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(fixture)
	}))
	defer ts.Close()

	// Creates a client to execute tests against the test server
	client, err := api.New(ts.URL)
	require.NoError(t, err, "could not create api client")

	out, err := client.MemberDelete(context.Background(), "member001")
	require.NoError(t, err, "could not execute api request")
	require.Equal(t, fixture, out, "response did not match fixture")
}

func TestProfileDetail(t *testing.T) {
	fixture := &api.Member{
		ID:    "001",
		Name:  "Leopold Wentzel",
		Email: "leopold.wentzel@gmail.com",
	}

	// Create a test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/v1/profile", r.URL.Path)
		require.Equal(t, http.MethodGet, r.Method)

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(fixture)
	}))
	defer ts.Close()

	// Create a client to execute tests against the test server
	client, err := api.New(ts.URL)
	require.NoError(t, err, "could not create api client")

	out, err := client.ProfileDetail(context.Background())
	require.NoError(t, err, "could not execute api request")
	require.Equal(t, fixture, out, "unexpected response error")
}

func TestProfileUpdate(t *testing.T) {
	fixture := &api.Member{
		ID:    "001",
		Name:  "Leopold A Wentzel",
		Email: "leopold.wentzel@gmail.com",
	}

	// Create a test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/v1/profile", r.URL.Path)
		require.Equal(t, http.MethodPut, r.Method)

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(fixture)
	}))
	defer ts.Close()

	// Create a client to execute tests against the server
	client, err := api.New(ts.URL)
	require.NoError(t, err, "could not execute api request")

	req := &api.Member{
		Name:  "Leopold A Wentzel",
		Email: "leopold.wentzel@gmail.com",
	}
	out, err := client.ProfileUpdate(context.Background(), req)
	require.NoError(t, err, "could not execute api request")
	require.Equal(t, fixture, out, "unexpected response error")
}

func TestTenantProjectList(t *testing.T) {
	fixture := &api.TenantProjectPage{
		TenantID: "01",
		TenantProjects: []*api.Project{
			{
				ID:          "001",
				Name:        "project01",
				Description: "This is my first project",
				Owner: api.Member{
					Name: "Luke Hamilton",
				},
				Status:       "Active",
				ActiveTopics: 12,
			},
		},
		NextPageToken: "23",
	}

	// Creates a test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/v1/tenant/tenant01/projects", r.URL.Path)

		params := r.URL.Query()
		require.Equal(t, "2", params.Get("page_size"))
		require.Equal(t, "12", params.Get("next_page_token"))

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(fixture)
	}))
	defer ts.Close()

	// Creates a client to executes tests against the test server
	client, err := api.New(ts.URL)
	require.NoError(t, err, "could not create api client")

	req := &api.PageQuery{
		PageSize:      2,
		NextPageToken: "12",
	}

	out, err := client.TenantProjectList(context.Background(), "tenant01", req)
	require.NoError(t, err, "could not execute api request")
	require.Equal(t, fixture, out, "unexpected response error")
}

func TestTenantProjectCreate(t *testing.T) {
	fixture := &api.Project{
		ID:          "001",
		Name:        "project01",
		Description: "My first project",
	}

	// Creates a test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/v1/tenant/tenant01/projects", r.URL.Path)

		in := &api.Project{}
		err := json.NewDecoder(r.Body).Decode(in)
		require.NoError(t, err, "could not decode request")

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(fixture)
	}))
	defer ts.Close()

	// Creates a client to execute tests against the test server
	client, err := api.New(ts.URL)
	require.NoError(t, err, "could not create api client")

	out, err := client.TenantProjectCreate(context.Background(), "tenant01", &api.Project{})
	require.NoError(t, err, "could not execute api request")
	require.Equal(t, fixture, out, "unexpected response error")
}

func TestTenantProjectPatch(t *testing.T) {
	fixture := &api.Project{
		ID:          "tenant001",
		Name:        "Some project",
		Description: "Updated description",
	}

	// Creates a test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPatch, r.Method)
		require.Equal(t, "/v1/tenant/tenant001/projects/project001", r.URL.Path)

		in := &api.Project{}
		err := json.NewDecoder(r.Body).Decode(in)
		require.NoError(t, err, "could not decode request")

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(fixture)
	}))
	defer ts.Close()

	// Creates a client to execute tests against the test server
	client, err := api.New(ts.URL)
	require.NoError(t, err, "could not create api client")

	in := &api.Project{
		Description: "Updated description",
	}
	out, err := client.TenantProjectPatch(context.Background(), "tenant001", "project001", in)
	require.NoError(t, err, "could not execute api request")
	require.Equal(t, fixture, out, "response did not match fixture")
}

func TestTenantProjectStats(t *testing.T) {
	fixture := []*api.StatValue{
		{
			Name:  "projects",
			Value: 3,
		},
		{
			Name:  "topics",
			Value: 2,
		},
		{
			Name:  "keys",
			Value: 3,
		},
	}

	// Create a test server.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/v1/tenant/001/projects/stats", r.URL.Path)

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(fixture)
	}))
	defer ts.Close()

	// Create a client to execute tests against the server.
	client, err := api.New(ts.URL)
	require.NoError(t, err, "could not create api client")

	out, err := client.TenantProjectStats(context.Background(), "001")
	require.NoError(t, err, "could not execute api request")
	require.Equal(t, fixture, out, "unexpected response error")
}

func TestProjectList(t *testing.T) {
	fixture := &api.ProjectPage{
		Projects: []*api.Project{
			{
				ID:   "001",
				Name: "project01",
			},
		},
		NextPageToken: "23",
	}

	// Creates a test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/v1/projects", r.URL.Path)

		params := r.URL.Query()
		require.Equal(t, "2", params.Get("page_size"))
		require.Equal(t, "12", params.Get("next_page_token"))

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(fixture)
	}))
	defer ts.Close()

	// Creates a client to execute tests against the test server
	client, err := api.New(ts.URL)
	require.NoError(t, err, "could not create api client")

	req := &api.PageQuery{
		PageSize:      2,
		NextPageToken: "12",
	}

	out, err := client.ProjectList(context.Background(), req)
	require.NoError(t, err, "could not execute api request")
	require.Equal(t, fixture, out, "unexpected response error")
}

func TestProjectCreate(t *testing.T) {
	fixture := &api.Project{
		ID:   "001",
		Name: "project01",
	}

	// Creates a test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/v1/projects", r.URL.Path)

		in := &api.Project{}
		err := json.NewDecoder(r.Body).Decode(in)
		require.NoError(t, err, "could not decode request")

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(fixture)
	}))
	defer ts.Close()

	// Creates a client to execute tests against the test server
	client, err := api.New(ts.URL)
	require.NoError(t, err, "could not create api client")

	out, err := client.ProjectCreate(context.Background(), &api.Project{})
	require.NoError(t, err, "could not execute api request")
	require.Equal(t, fixture, out, "unexpected response error")
}

func TestProjectDetail(t *testing.T) {
	fixture := &api.Project{
		ID:   "001",
		Name: "project01",
	}

	// Creates a test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/v1/projects/project001", r.URL.Path)

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(fixture)
	}))
	defer ts.Close()

	// Creates a client to execute tests against the test server
	client, err := api.New(ts.URL)
	require.NoError(t, err, "could not create api client")

	out, err := client.ProjectDetail(context.Background(), "project001")
	require.NoError(t, err, "could not execute api request")
	require.Equal(t, fixture, out, "unexpected response error")
}

func TestProjectUpdate(t *testing.T) {
	fixture := &api.Project{
		ID:   "001",
		Name: "project01",
	}

	// Creates a test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPut, r.Method)
		require.Equal(t, "/v1/projects/001", r.URL.Path)

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(fixture)
	}))
	defer ts.Close()

	// Creates a client to execute tests against the server
	client, err := api.New(ts.URL)
	require.NoError(t, err, "could not execute api request")

	req := &api.Project{
		ID:   "001",
		Name: "project02",
	}

	rep, err := client.ProjectUpdate(context.Background(), req)
	require.NoError(t, err, "could not execute api request")
	require.Equal(t, fixture, rep, "unexpected response error")
}

func TestProjectPatch(t *testing.T) {
	fixture := &api.Project{
		ID:          "001",
		Name:        "Some project",
		Description: "Updated description",
	}

	// Create a test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPatch, r.Method)
		require.Equal(t, "/v1/projects/001", r.URL.Path)

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(fixture)
	}))
	defer ts.Close()

	// Create a client to request the server
	client, err := api.New(ts.URL)
	require.NoError(t, err, "could not execute api request")

	in := &api.Project{
		Description: "Updated description",
	}

	rep, err := client.ProjectPatch(context.Background(), "001", in)
	require.NoError(t, err, "could not execute api request")
	require.Equal(t, fixture, rep, "response did not match fixture")
}

func TestProjectDelete(t *testing.T) {
	fixture := &api.Reply{}

	// Creates a test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodDelete, r.Method)
		require.Equal(t, "/v1/projects/project001", r.URL.Path)

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(fixture)
	}))
	defer ts.Close()

	// Creates a client to execute tests against the test server
	client, err := api.New(ts.URL)
	require.NoError(t, err, "could not create api client")

	err = client.ProjectDelete(context.Background(), "project001")
	require.NoError(t, err, "could not execute api request")
}

func TestProjectTopicList(t *testing.T) {
	fixture := &api.ProjectTopicPage{
		ProjectID: "001",
		Topics: []*api.Topic{
			{
				ID:   "005",
				Name: "topic002",
			},
		},
		NextPageToken: "23",
	}

	// Creates a test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/v1/projects/project001/topics", r.URL.Path)

		params := r.URL.Query()
		require.Equal(t, "12", params.Get("next_page_token"))
		require.Equal(t, "2", params.Get("page_size"))

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(fixture)
	}))
	defer ts.Close()

	// Create a client to execute tests against the test server
	client, err := api.New(ts.URL)
	require.NoError(t, err, "could not create api client")

	req := &api.PageQuery{
		PageSize:      2,
		NextPageToken: "12",
	}

	out, err := client.ProjectTopicList(context.Background(), "project001", req)
	require.NoError(t, err, "could not execute api request")
	require.Equal(t, fixture, out, "unexpected response returned")
}

func TestProjectTopicCreate(t *testing.T) {
	fixture := &api.Topic{
		ID:   "001",
		Name: "topic01",
	}
	// Create a test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/v1/projects/project001/topics", r.URL.Path)

		in := &api.Topic{}
		err := json.NewDecoder(r.Body).Decode(in)
		require.NoError(t, err, "could not decode request")

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(fixture)
	}))
	defer ts.Close()

	// Create a client to execute tests against the test server
	client, err := api.New(ts.URL)
	require.NoError(t, err, "could not create api client")

	out, err := client.ProjectTopicCreate(context.Background(), "project001", &api.Topic{})
	require.NoError(t, err, "could not execute api request")
	require.Equal(t, fixture, out, "unexpected response returned")
}

func TestProjectQuery(t *testing.T) {
	fixture := &api.ProjectQueryResponse{
		Results: []*api.QueryResult{
			{
				Metadata: map[string]string{
					"key": "value",
				},
				Mimetype: "application/json",
				Version:  "Document v1.0.1",
				Data:     "{\"foo\": \"bar\"}",
				Created:  time.Now().UTC().Format(time.RFC3339),
			},
		},
	}

	// Create a test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/v1/projects/project001/query", r.URL.Path)

		in := &api.ProjectQueryRequest{}
		err := json.NewDecoder(r.Body).Decode(in)
		require.NoError(t, err, "could not decode request")

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(fixture)
	}))
	defer ts.Close()

	// Create a client to execute requests against the test server
	client, err := api.New(ts.URL)
	require.NoError(t, err, "could not create api client")

	req := &api.ProjectQueryRequest{
		ProjectID: "project001",
		Query:     "SELECT * FROM topic001 LIMIT 10",
	}
	out, err := client.ProjectQuery(context.Background(), req)
	require.NoError(t, err, "could not execute api request")
	require.Equal(t, fixture, out, "unexpected response returned")
}

func TestTopicList(t *testing.T) {
	fixture := &api.TopicPage{
		Topics: []*api.Topic{
			{
				ID:   "005",
				Name: "topic01",
			},
		},
		NextPageToken: "23",
	}

	// Creates a test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/v1/topics", r.URL.Path)

		params := r.URL.Query()
		require.Equal(t, "12", params.Get("next_page_token"))
		require.Equal(t, "2", params.Get("page_size"))

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(fixture)
	}))
	defer ts.Close()

	// Create a client to execute tests against the test server
	client, err := api.New(ts.URL)
	require.NoError(t, err, "could not create api client")

	req := &api.PageQuery{
		PageSize:      2,
		NextPageToken: "12",
	}

	out, err := client.TopicList(context.Background(), req)
	require.NoError(t, err, "could not execute api request")
	require.Equal(t, fixture, out, "unexpected response returned")
}

func TestTopicDetail(t *testing.T) {
	fixture := &api.Topic{
		ID:   "001",
		Name: "topic01",
	}

	// Creates a test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/v1/topics/topic001", r.URL.Path)

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(fixture)
	}))
	defer ts.Close()

	// Creates a client to execute tests against the test server
	client, err := api.New(ts.URL)
	require.NoError(t, err, "could not create api client")

	out, err := client.TopicDetail(context.Background(), "topic001")
	require.NoError(t, err, "could not execute api request")
	require.Equal(t, fixture, out, "unexpected response error")
}

func TestTopicEvents(t *testing.T) {
	fixture := []*api.EventTypeInfo{
		{
			Type:     "Document",
			Version:  "1.0.0",
			Mimetype: "application/json",
			Events: &api.StatValue{
				Value:   12345678,
				Percent: 96.0,
			},
			Storage: &api.StatValue{
				Value:   512,
				Units:   "MB",
				Percent: 98.5,
			},
		},
		{
			Type:     "Feed Item",
			Version:  "0.8.1",
			Mimetype: "application/rss",
			Events: &api.StatValue{
				Value:   98765,
				Percent: 4.0,
			},
			Storage: &api.StatValue{
				Value:   4.3,
				Units:   "KB",
				Percent: 1.5,
			},
		},
	}

	// Create a test server to return the fixture
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/v1/topics/topic001/events", r.URL.Path)

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(fixture)
	}))
	defer ts.Close()

	// Create a client to call the test server
	client, err := api.New(ts.URL)
	require.NoError(t, err, "could not create api client")

	out, err := client.TopicEvents(context.Background(), "topic001")
	require.NoError(t, err, "could not execute api request")
	require.Equal(t, fixture, out, "unexpected response returned")
}

func TestTopicStats(t *testing.T) {
	fixture := []*api.StatValue{
		{
			Name:  "publishers",
			Value: 2,
		},
		{
			Name:  "subscribers",
			Value: 3,
		},
		{
			Name:  "total_events",
			Value: 100,
		},
		{
			Name:  "storage",
			Value: 256,
			Units: "MB",
		},
	}

	// Create a test server to return the fixture
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/v1/topics/topic001/stats", r.URL.Path)

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(fixture)
	}))

	// Create a client to call the test server
	client, err := api.New(ts.URL)
	require.NoError(t, err, "could not create api client")

	out, err := client.TopicStats(context.Background(), "topic001")
	require.NoError(t, err, "could not execute api request")
	require.Equal(t, fixture, out, "unexpected response returned")
}

func TestTopicUpdate(t *testing.T) {
	fixture := &api.Topic{
		ID:   "001",
		Name: "topic01",
	}

	// Creates a test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPut, r.Method)
		require.Equal(t, "/v1/topics/001", r.URL.Path)

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(fixture)
	}))
	defer ts.Close()

	// Creates a client to execute tests against the server
	client, err := api.New(ts.URL)
	require.NoError(t, err, "could not execute api request")

	req := &api.Topic{
		ID:   "001",
		Name: "topic02",
	}

	rep, err := client.TopicUpdate(context.Background(), req)
	require.NoError(t, err, "could not execute api request")
	require.Equal(t, fixture, rep, "unexpected response error")
}

func TestTopicDelete(t *testing.T) {
	fixture := &api.Confirmation{
		ID:    "topic001",
		Token: "token",
	}

	// Creates a test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodDelete, r.Method)
		require.Equal(t, "/v1/topics/topic001", r.URL.Path)

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(fixture)
	}))
	defer ts.Close()

	// Creates a client to execute tests against the test server
	client, err := api.New(ts.URL)
	require.NoError(t, err, "could not create api client")

	req := &api.Confirmation{
		ID: "topic001",
	}
	out, err := client.TopicDelete(context.TODO(), req)
	require.NoError(t, err, "could not execute api request")
	require.Equal(t, fixture, out, "unexpected response error")
}

func TestProjectAPIKeyList(t *testing.T) {
	fixture := &api.ProjectAPIKeyPage{
		ProjectID: "001",
		APIKeys: []*api.APIKeyPreview{
			{
				ID:          "001",
				ClientID:    "client001",
				Name:        "myapikey",
				Permissions: "Full",
				Status:      "Active",
				LastUsed:    time.Now().Format(time.RFC3339Nano),
				Created:     time.Now().Format(time.RFC3339Nano),
				Modified:    time.Now().Format(time.RFC3339Nano),
			},
		},
		PrevPageToken: "21",
		NextPageToken: "23",
	}

	// Create a test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/v1/projects/project001/apikeys", r.URL.Path)
		params := r.URL.Query()
		require.Equal(t, "2", params.Get("page_size"))
		require.Equal(t, "12", params.Get("next_page_token"))

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(fixture)
	}))
	defer ts.Close()

	// Creates a client to executes tests against the test server
	client, err := api.New(ts.URL)
	require.NoError(t, err, "could not create api client")

	req := &api.PageQuery{
		PageSize:      2,
		NextPageToken: "12",
	}

	out, err := client.ProjectAPIKeyList(context.Background(), "project001", req)
	require.NoError(t, err, "could not execute api request")
	require.Equal(t, fixture, out, "unexpected response error")
}

func TestProjectAPIKeyCreate(t *testing.T) {
	fixture := &api.APIKey{
		ID:           "001",
		ClientID:     "client001",
		ClientSecret: "segredo",
		Name:         "myapikey",
		Owner:        "Ryan Moore",
		Permissions:  []string{"Read", "Write", "Delete"},
		Created:      time.Now().Format(time.RFC3339Nano),
	}

	//Creates a test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/v1/projects/project001/apikeys", r.URL.Path)

		in := &api.APIKey{}
		err := json.NewDecoder(r.Body).Decode(in)
		require.NoError(t, err, "could not decode request")

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(fixture)
	}))
	defer ts.Close()

	// Creates a client to execute tests against the test server
	client, err := api.New(ts.URL)
	require.NoError(t, err, "could not create api client")

	req := &api.APIKey{
		ID: "001",
	}
	out, err := client.ProjectAPIKeyCreate(context.Background(), "project001", req)
	require.NoError(t, err, "could not execute api request")
	require.Equal(t, fixture, out, "unexpected response error")
}

func TestAPIKeyCreate(t *testing.T) {
	fixture := &api.APIKey{
		ID:           "001",
		ClientID:     "client001",
		ClientSecret: "segredo",
		Name:         "myapikey",
		Owner:        "Ryan Moore",
		Permissions:  []string{"Read", "Write", "Delete"},
		Created:      time.Now().Format(time.RFC3339Nano),
	}

	// Create a test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/v1/apikeys", r.URL.Path)

		in := &api.APIKey{}
		err := json.NewDecoder(r.Body).Decode(in)
		require.NoError(t, err, "could not decode request")

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(fixture)
	}))
	defer ts.Close()

	// Creates a client to execute tests against the test server
	client, err := api.New(ts.URL)
	require.NoError(t, err, "could not create api client")

	out, err := client.APIKeyCreate(context.Background(), &api.APIKey{})
	require.NoError(t, err, "could not execute api request")
	require.Equal(t, fixture, out, "response did not match")
}

func TestAPIKeyList(t *testing.T) {
	fixture := &api.APIKeyPage{
		APIKeys: []*api.APIKey{
			{
				ID:           "001",
				ClientID:     "client001",
				ClientSecret: "segredo",
				Name:         "myapikey",
				Owner:        "Ryan Moore",
				Permissions:  []string{"Read", "Write", "Delete"},
				Created:      time.Now().Format(time.RFC3339Nano),
			},
		},
		PrevPageToken: "21",
		NextPageToken: "23",
	}

	// Create a test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/v1/apikeys", r.URL.Path)

		params := r.URL.Query()
		require.Equal(t, "2", params.Get("page_size"))
		require.Equal(t, "12", params.Get("next_page_token"))

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(fixture)
	}))
	defer ts.Close()

	// Creates a client to execute tests against the test server
	client, err := api.New(ts.URL)
	require.NoError(t, err, "could not create api client")

	req := &api.PageQuery{
		PageSize:      2,
		NextPageToken: "12",
	}

	out, err := client.APIKeyList(context.Background(), req)
	require.NoError(t, err, "could not execute api request")
	require.Equal(t, fixture, out, "unexpected response error")
}

func TestAPIKeyDetail(t *testing.T) {
	fixture := &api.APIKey{
		ID:           "001",
		ClientID:     "client001",
		ClientSecret: "segredo",
		Name:         "myapikey",
		Owner:        "Ryan Moore",
		Permissions:  []string{"Read", "Write", "Delete"},
		Created:      time.Now().Format(time.RFC3339Nano),
	}

	// Creates a test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/v1/apikeys/apikey001", r.URL.Path)

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(fixture)
	}))
	defer ts.Close()

	// Creates a client to execute tests against the test server
	client, err := api.New(ts.URL)
	require.NoError(t, err, "could not create api client")

	out, err := client.APIKeyDetail(context.Background(), "apikey001")
	require.NoError(t, err, "could not execute api request")
	require.Equal(t, fixture, out, "unexpected response error")
}

func TestAPIKeyUpdate(t *testing.T) {
	fixture := &api.APIKey{
		ID:   "101",
		Name: "apikey01",
	}

	// Creates a test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPut, r.Method)
		require.Equal(t, "/v1/apikeys/101", r.URL.Path)

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(fixture)
	}))
	defer ts.Close()

	// Creates a client to execute tests against the server
	client, err := api.New(ts.URL)
	require.NoError(t, err, "could not execute api request")

	req := &api.APIKey{
		ID:   "101",
		Name: "apikey02",
	}

	rep, err := client.APIKeyUpdate(context.Background(), req)
	require.NoError(t, err, "could not execute api request")
	require.Equal(t, fixture, rep, "unexpected response error")
}

func TestAPIKeyDelete(t *testing.T) {
	fixture := &api.Reply{}

	// Creates a test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodDelete, r.Method)
		require.Equal(t, "/v1/apikeys/apikey001", r.URL.Path)

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(fixture)
	}))
	defer ts.Close()

	// Creates a client to execute tests against the test server
	client, err := api.New(ts.URL)
	require.NoError(t, err, "could not create api client")

	err = client.APIKeyDelete(context.Background(), "apikey001")
	require.NoError(t, err, "could not execute api request")
}

func TestAPIKeyPermissions(t *testing.T) {
	fixture := []string{"topics:read", "topics:destroy"}

	// Creates a test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/v1/apikeys/permissions", r.URL.Path)

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(fixture)
	}))
	defer ts.Close()

	// Creates a client to execute tests against the test server
	client, err := api.New(ts.URL)
	require.NoError(t, err, "could not create api client")

	out, err := client.APIKeyPermissions(context.Background())
	require.NoError(t, err, "could not execute api request")
	require.Equal(t, fixture, out, "unexpected response error")
}

func TestCSRFProtect(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := r.Cookie("csrf_token")
		if err != nil {
			t.Log("no csrf_token cookie")
			http.Error(w, "no csrf_token cookie", http.StatusBadRequest)
			return
		}

		ref, err := r.Cookie("csrf_reference_token")
		if err != nil {
			t.Log("no csrf_reference_token cookie")
			http.Error(w, "no csrf_reference_token cookie", http.StatusBadRequest)
			return
		}

		if token.Value != ref.Value {
			t.Log("csrf_token does not match reference")
			http.Error(w, "csrf_token does not match reference", http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	// Creates a client to execute tests against the test server
	client, err := api.New(ts.URL)
	require.NoError(t, err, "could not create api client")

	apiv1, ok := client.(*api.APIv1)
	require.True(t, ok, "client was not an APIv1 client")

	err = apiv1.SetCSRFProtect(true)
	require.NoError(t, err, "could not set csrf protection")

	req, err := apiv1.NewRequest(context.Background(), http.MethodPost, "/", nil, nil)
	require.NoError(t, err, "could not create POST request")

	_, err = apiv1.Do(req, nil, true)
	require.NoError(t, err, "csrf protect failed")
}

func TestAuthCookies(t *testing.T) {
	// NOTE: JWT token secret is: supersecretsquirrel

	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessToken := &http.Cookie{
			Name:     middleware.AccessTokenCookie,
			Value:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyLCJhdWQiOlsiaHR0cHM6Ly9sb2NhbGhvc3QvYXV0aCJdLCJpc3MiOiJodHRwczovL2xvY2FsaG9zdCJ9.raJFGdWDm-OtxTnbYaU2-emtVaJlFZPJ_98cLH5hz5Y",
			Path:     "/",
			Expires:  time.Now().Add(10 * time.Minute),
			Secure:   true,
			HttpOnly: true,
		}
		http.SetCookie(w, accessToken)

		refreshToken := &http.Cookie{
			Name:     middleware.RefreshTokenCookie,
			Value:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyLCJhdWQiOlsiaHR0cHM6Ly9sb2NhbGhvc3QvcmVmcmVzaCJdLCJpc3MiOiJodHRwczovL2xvY2FsaG9zdCJ9.0e8MJCKY_-1itkI4MKqGSdkUVqbk5aqDcK6Lui8LKrI",
			Path:     "/",
			Expires:  time.Now().Add(10 * time.Minute),
			Secure:   true,
			HttpOnly: true,
		}
		http.SetCookie(w, refreshToken)

		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	client, err := api.New(ts.URL, api.WithClient(ts.Client()))
	require.NoError(t, err, "could not create https client")

	apiv1, ok := client.(*api.APIv1)
	require.True(t, ok, "client is not an APIv1")

	req, err := apiv1.NewRequest(context.Background(), http.MethodGet, "/", nil, nil)
	require.NoError(t, err, "could not create api request")

	t.Run("AccessToken", func(t *testing.T) {
		_, err := apiv1.Do(req, nil, true)
		require.NoError(t, err, "could not execute request")

		accessToken, err := apiv1.AccessToken()
		require.NoError(t, err, "could not retrieve access token")
		require.Equal(t, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyLCJhdWQiOlsiaHR0cHM6Ly9sb2NhbGhvc3QvYXV0aCJdLCJpc3MiOiJodHRwczovL2xvY2FsaG9zdCJ9.raJFGdWDm-OtxTnbYaU2-emtVaJlFZPJ_98cLH5hz5Y", accessToken)
	})

	t.Run("RefreshToken", func(t *testing.T) {
		_, err := apiv1.Do(req, nil, true)
		require.NoError(t, err, "could not execute request")

		refreshToken, err := apiv1.RefreshToken()
		require.NoError(t, err, "could not retrieve refresh token")
		require.Equal(t, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyLCJhdWQiOlsiaHR0cHM6Ly9sb2NhbGhvc3QvcmVmcmVzaCJdLCJpc3MiOiJodHRwczovL2xvY2FsaG9zdCJ9.0e8MJCKY_-1itkI4MKqGSdkUVqbk5aqDcK6Lui8LKrI", refreshToken)

	})
}
