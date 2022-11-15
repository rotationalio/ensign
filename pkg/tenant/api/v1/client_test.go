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

func TestTenantList(t *testing.T) {
	fixture := &api.TenantPage{
		Tenants: []*api.Tenant{
			{
				ID:              "001",
				TenantName:      "tenant01",
				EnvironmentType: "Dev",
			},
			{
				ID:              "002",
				TenantName:      "tenant02",
				EnvironmentType: "Prod",
			},
			{
				ID:              "003",
				TenantName:      "tenant03",
				EnvironmentType: "Stage",
			},
		},
		PrevPageToken: "2121",
		NextPageToken: "4040",
	}

	// Creates a test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/v1/tenant", r.URL.Path)

		rURL, _ := url.Parse("/v1/tenant?next_page_token=1212&page_size=2")

		var params url.Values = rURL.Query()

		require.Equal(t, "1212", params.Get("next_page_token"))
		require.Equal(t, "2", params.Get("page_size"))

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(fixture)
	}))
	defer ts.Close()

	// Creates a client to execute tests against the server
	client, err := api.New(ts.URL)
	require.NoError(t, err, "could not create api client")

	req := &api.PageQuery{}

	out, err := client.TenantList(context.TODO(), req)
	require.NoError(t, err, "could not execute api request")
	require.Equal(t, fixture, out, "unexpected response returned")
}

func TestTenantCreate(t *testing.T) {
	fixture := &api.Tenant{
		ID:              "1234",
		TenantName:      "feist",
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
		w.WriteHeader(http.StatusNoContent)
		json.NewEncoder(w).Encode(fixture)
	}))
	defer ts.Close()

	// Creates a client to execute tests against the test server
	client, err := api.New(ts.URL)
	require.NoError(t, err, "could not create api client")

	req := &api.Tenant{
		ID:              "1234",
		TenantName:      "feist",
		EnvironmentType: "Dev",
	}

	err = client.TenantCreate(context.TODO(), req)
	require.NoError(t, err, "could not execute api request")
}

func TestTenantDetail(t *testing.T) {
	fixture := &api.Tenant{
		ID:         "001",
		TenantName: "tenant01",
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

	out, err := client.TenantDetail(context.TODO(), "tenant01")
	require.NoError(t, err, "could not execute api request")
	require.Equal(t, fixture, out, "unexpected result occurred")
}

func TestTenantUpdate(t *testing.T) {
	fixture := &api.Tenant{
		ID:              "001",
		TenantName:      "tenant01",
		EnvironmentType: "Dev",
	}

	// Creates a new test server
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
		TenantName:      "tenant02",
		EnvironmentType: "Prod",
	}

	rep, err := client.TenantUpdate(context.TODO(), req)
	require.NoError(t, err, "could not execute api request")
	require.Equal(t, fixture, rep, "unexpected response returned")
}
func TestTenantDelete(t *testing.T) {
	fixture := &api.Reply{}

	// Creates a new test server
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

	err = client.TenantDelete(context.TODO(), "tenant01")
	require.NoError(t, err, "could not execute api request")
}

func TestTenantMemberList(t *testing.T) {
	fixture := &api.TenantMemberPage{
		TenantMembers: []*api.TenantMember{
			{
				TenantID:   "01",
				MemberID:   "002",
				MemberName: "Luke Hamilton",
				MemberRole: "Admin",
			},
		},
		PrevPageToken: "1212",
		NextPageToken: "1214",
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/v1/tenant/tenant01/members", r.URL.Path)

		params := url.Values{}
		params.Set("next_page_token", "1212")
		params.Add("page_size", "2")

		require.Equal(t, "1212", params.Get("next_page_token"))
		require.Equal(t, "2", params.Get("page_size"))

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(fixture)
	}))
	defer ts.Close()

	client, err := api.New(ts.URL)
	require.NoError(t, err, "could not create api client")

	req := &api.PageQuery{}

	out, err := client.TenantMemberList(context.TODO(), "tenant01", req)
	require.NoError(t, err, "could not execute api request")
	require.Equal(t, fixture, out, "unexpected result occurred")
}

func TestTenantMemberCreate(t *testing.T) {
	fixture := &api.TenantMember{
		TenantID:   "01",
		MemberID:   "02",
		MemberName: "Luke Hamilton",
		MemberRole: "Admin",
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/v1/tenant/tenant01/members", r.URL.Path)

		in := &api.TenantMember{}
		err := json.NewDecoder(r.Body).Decode(in)
		require.NoError(t, err, "could not decode request")

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusNoContent)
		json.NewEncoder(w).Encode(fixture)
	}))
	defer ts.Close()

	client, err := api.New(ts.URL)
	require.NoError(t, err, "could not create api client")

	req := &api.TenantMember{
		TenantID:   "01",
		MemberID:   "02",
		MemberName: "Luke Hamilton",
		MemberRole: "Admin",
	}

	err = client.TenantMemberCreate(context.TODO(), "tenant01", req)
	require.NoError(t, err, "could not execute api request")
}

func TestMemberList(t *testing.T) {
	fixture := &api.MemberPage{
		Members: []*api.Member{
			{
				ID:         "002",
				MemberName: "Ryan Moore",
				MemberRole: "Admin",
			},
		},
		PrevPageToken: "2121",
		NextPageToken: "2123",
	}
	// Creates a test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/v1/members", r.URL.Path)

		params := r.URL.Query()

		require.Equal(t, "1212", params.Get("next_page_token"))
		require.Equal(t, "2", params.Get("page_size"))

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(fixture)
	}))
	defer ts.Close()

	client, err := api.New(ts.URL)
	require.NoError(t, err, "could not create api client")

	req := &api.PageQuery{}

	out, err := client.MemberList(context.TODO(), req)
	require.NoError(t, err, "could not execute api request")
	require.Equal(t, fixture, out, "unexpected response error")
}

func TestMemberCreate(t *testing.T) {
	fixture := &api.Member{}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/v1/members", r.URL.Path)

		in := &api.Member{}
		err := json.NewDecoder(r.Body).Decode(in)
		require.NoError(t, err, "could not decode request")

		w.Header().Add("Content-Type", "application/json; charset=utft-8")
		w.WriteHeader(http.StatusNoContent)
		json.NewEncoder(w).Encode(fixture)
	}))
	defer ts.Close()

	client, err := api.New(ts.URL)
	require.NoError(t, err, "could not create api client")

	req := &api.Member{}

	err = client.MemberCreate(context.TODO(), req)
	require.NoError(t, err, "could not execute api request")
}

func TestSignUp(t *testing.T) {
	// Creates a Test Server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/v1/notifications/signup", r.URL.Path)

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	// Creates a client to execute tests against the test server
	client, err := api.New(ts.URL)
	require.NoError(t, err)

	contact := &api.ContactInfo{
		FirstName:    "Jane",
		LastName:     "Eyere",
		Email:        "jane@example.com",
		Country:      "SG",
		Title:        "Director",
		Organization: "Simple, PTE",
	}

	err = client.SignUp(context.Background(), contact)
	require.NoError(t, err, "could not execute signup request")
}
