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
		PrevPageToken: "21",
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
		TenantName:      "tenant01",
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
		TenantName:      "tenant01",
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
		TenantName:      "tenant02",
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

func TestTenantMemberList(t *testing.T) {
	fixture := &api.TenantMemberPage{
		TenantID: "002",
		TenantMembers: []*api.Member{
			{
				ID:   "002",
				Name: "Luke Hamilton",
				Role: "Admin",
			},
		},
		PrevPageToken: "1212",
		NextPageToken: "1214",
	}

	// Creates a test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/v1/tenant/tenant002/members", r.URL.Path)

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

	out, err := client.TenantMemberList(context.Background(), "tenant002", req)
	require.NoError(t, err, "could not execute api request")
	require.Equal(t, fixture, out, "unexpected response error")
}

func TestTenantMemberCreate(t *testing.T) {
	fixture := &api.Member{
		ID:   "02",
		Name: "Luke Hamilton",
		Role: "Admin",
	}

	// Creates a test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/v1/tenant/tenant02/members", r.URL.Path)

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

	out, err := client.TenantMemberCreate(context.Background(), "tenant02", &api.Member{})
	require.NoError(t, err, "could not execute api request")
	require.Equal(t, fixture, out, "unexpected response error")
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
		PrevPageToken: "21",
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

func TestMemberDelete(t *testing.T) {
	fixture := &api.Reply{}

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

	err = client.MemberDelete(context.Background(), "member001")
	require.NoError(t, err, "could not execute api request")
}

func TestTenantProjectList(t *testing.T) {
	fixture := &api.TenantProjectPage{
		TenantID: "01",
		TenantProjects: []*api.Project{
			{
				ID:   "001",
				Name: "project01",
			},
		},
		PrevPageToken: "21",
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
		ID:   "001",
		Name: "project01",
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

func TestProjectList(t *testing.T) {
	fixture := &api.ProjectPage{
		Projects: []*api.Project{
			{
				ID:   "001",
				Name: "project01",
			},
		},
		PrevPageToken: "21",
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
		TenantTopics: []*api.Topic{
			{
				ID:   "005",
				Name: "topic002",
			},
		},
		PrevPageToken: "21",
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

func TestTopicList(t *testing.T) {
	fixture := &api.TopicPage{
		Topics: []*api.Topic{
			{
				ID:   "005",
				Name: "topic01",
			},
		},
		PrevPageToken: "21",
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

func TestTopicCreate(t *testing.T) {
	fixture := &api.Topic{
		ID:   "001",
		Name: "topic01",
	}
	// Create a test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/v1/topics", r.URL.Path)

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

	out, err := client.TopicCreate(context.Background(), &api.Topic{})
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
	fixture := &api.Reply{}

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

	err = client.TopicDelete(context.TODO(), "topic001")
	require.NoError(t, err, "could not execute api request")
}

func TestProjectAPIKeyList(t *testing.T) {
	fixture := &api.ProjectAPIKeyPage{
		ProjectID: "001",
		APIKeys: []*api.APIKey{
			{
				ID:           001,
				ClientID:     "client001",
				ClientSecret: "segredo",
				Name:         "myapikey",
				Owner:        "Ryan Moore",
				Permissions:  []string{"Read", "Write", "Delete"},
				Created:      time.Now().Format(time.RFC3339Nano),
				Modified:     time.Now().Format(time.RFC3339Nano),
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
		ID:           001,
		ClientID:     "client001",
		ClientSecret: "segredo",
		Name:         "myapikey",
		Owner:        "Ryan Moore",
		Permissions:  []string{"Read", "Write", "Delete"},
		Created:      time.Now().Format(time.RFC3339Nano),
		Modified:     time.Now().Format(time.RFC3339Nano),
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

	out, err := client.ProjectAPIKeyCreate(context.Background(), "project001", &api.APIKey{})
	require.NoError(t, err, "could not execute api request")
	require.Equal(t, fixture, out, "unexpected response error")
}

func TestAPIKeyList(t *testing.T) {
	fixture := &api.APIKeyPage{
		APIKeys: []*api.APIKey{
			{
				ID:           001,
				ClientID:     "client001",
				ClientSecret: "segredo",
				Name:         "myapikey",
				Owner:        "Ryan Moore",
				Permissions:  []string{"Read", "Write", "Delete"},
				Created:      time.Now().Format(time.RFC3339Nano),
				Modified:     time.Now().Format(time.RFC3339Nano),
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

func TestAPIKeyCreate(t *testing.T) {
	fixture := &api.APIKey{
		ID:           001,
		ClientID:     "client001",
		ClientSecret: "segredo",
		Name:         "myapikey",
		Owner:        "Ryan Moore",
		Permissions:  []string{"Read", "Write", "Delete"},
		Created:      time.Now().Format(time.RFC3339Nano),
		Modified:     time.Now().Format(time.RFC3339Nano),
	}

	//Create a test server
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
	require.Equal(t, fixture, out, "unexpected response error")
}

func TestAPIKeyDetail(t *testing.T) {
	fixture := &api.APIKey{
		ID:           001,
		ClientID:     "client001",
		ClientSecret: "segredo",
		Name:         "myapikey",
		Owner:        "Ryan Moore",
		Permissions:  []string{"Read", "Write", "Delete"},
		Created:      time.Now().Format(time.RFC3339Nano),
		Modified:     time.Now().Format(time.RFC3339Nano),
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
		ID:   101,
		Name: "apikey01",
	}

	// Creates a test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPut, r.Method)
		require.Equal(t, "/v1/apikey/101", r.URL.Path)

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(fixture)
	}))
	defer ts.Close()

	// Creates a client to execute tests against the server
	client, err := api.New(ts.URL)
	require.NoError(t, err, "could not execute api request")

	req := &api.APIKey{
		ID:   101,
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
