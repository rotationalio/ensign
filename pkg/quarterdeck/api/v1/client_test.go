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
	ulids "github.com/rotationalio/ensign/pkg/utils/ulid"
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
	fixture := &api.RegisterReply{}

	// Create a test server
	ts := httptest.NewServer(testhandler(fixture, http.MethodPost, "/v1/register"))
	defer ts.Close()

	// Create a client and execute endpoint request
	client, err := api.New(ts.URL)
	require.NoError(t, err, "could not create api client")

	req := &api.RegisterRequest{}

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

	rep, err := client.Refresh(context.TODO())
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

	req := &api.PageQuery{}
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
