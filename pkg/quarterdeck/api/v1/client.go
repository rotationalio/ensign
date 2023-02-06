package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"

	"github.com/google/go-querystring/query"
)

// New creates a new API v1 client that implements the Quarterdeck Client interface.
func New(endpoint string, opts ...ClientOption) (_ QuarterdeckClient, err error) {
	c := &APIv1{}
	if c.endpoint, err = url.Parse(endpoint); err != nil {
		return nil, fmt.Errorf("could not parse endpoint: %s", err)
	}

	// Apply our options
	for _, opt := range opts {
		if err = opt(c); err != nil {
			return nil, err
		}
	}

	// If an http client isn't specified, create a default client.
	if c.client == nil {
		c.client = &http.Client{
			Transport:     nil,
			CheckRedirect: nil,
			Timeout:       30 * time.Second,
		}

		// Create cookie jar for CSRF
		if c.client.Jar, err = cookiejar.New(nil); err != nil {
			return nil, fmt.Errorf("could not create cookiejar: %s", err)
		}
	}

	return c, nil
}

// APIv1 implements the QuarterdeckClient interface
type APIv1 struct {
	endpoint *url.URL     // the base url for all requests
	client   *http.Client // used to make http requests to the server
	creds    Credentials  // default credentials used to authorize requests
}

// Ensure the APIv1 implements the QuarterdeckClient interface
var _ QuarterdeckClient = &APIv1{}

//===========================================================================
// Client Methods
//===========================================================================

func (s *APIv1) Status(ctx context.Context) (out *StatusReply, err error) {
	// Make the HTTP request
	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodGet, "/v1/status", nil, nil); err != nil {
		return nil, err
	}

	// NOTE: we cannot use s.Do because we want to parse 503 Unavailable errors
	var rep *http.Response
	if rep, err = s.client.Do(req); err != nil {
		return nil, err
	}
	defer rep.Body.Close()

	// Detect other errors
	if rep.StatusCode != http.StatusOK && rep.StatusCode != http.StatusServiceUnavailable {
		return nil, fmt.Errorf("%s", rep.Status)
	}

	// Deserialize the JSON data from the response
	out = &StatusReply{}
	if err = json.NewDecoder(rep.Body).Decode(out); err != nil {
		return nil, fmt.Errorf("could not deserialize status reply: %s", err)
	}
	return out, nil
}

func (s *APIv1) Register(ctx context.Context, in *RegisterRequest) (out *RegisterReply, err error) {
	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodPost, "/v1/register", in, nil); err != nil {
		return nil, err
	}

	if _, err = s.Do(req, &out, true); err != nil {
		return nil, err
	}

	return out, nil
}

func (s *APIv1) Login(ctx context.Context, in *LoginRequest) (out *LoginReply, err error) {
	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodPost, "/v1/login", in, nil); err != nil {
		return nil, err
	}

	if _, err = s.Do(req, &out, true); err != nil {
		return nil, err
	}

	return out, nil
}

func (s *APIv1) Authenticate(ctx context.Context, in *APIAuthentication) (out *LoginReply, err error) {
	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodPost, "/v1/authenticate", in, nil); err != nil {
		return nil, err
	}

	if _, err = s.Do(req, &out, true); err != nil {
		return nil, err
	}

	return out, nil
}

func (s *APIv1) Refresh(ctx context.Context, in *RefreshRequest) (out *LoginReply, err error) {
	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodPost, "/v1/refresh", in, nil); err != nil {
		return nil, err
	}

	if _, err = s.Do(req, &out, true); err != nil {
		return nil, err
	}

	return out, nil
}

//===========================================================================
// API Keys Resource
//===========================================================================

func (s *APIv1) APIKeyList(ctx context.Context, in *APIPageQuery) (out *APIKeyList, err error) {
	var params url.Values
	if params, err = query.Values(in); err != nil {
		return nil, fmt.Errorf("could not encode query params: %s", err)
	}

	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodGet, "/v1/apikeys", nil, &params); err != nil {
		return nil, err
	}

	if _, err = s.Do(req, &out, true); err != nil {
		return nil, err
	}

	return out, nil
}

func (s *APIv1) APIKeyCreate(ctx context.Context, in *APIKey) (out *APIKey, err error) {
	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodPost, "/v1/apikeys", in, nil); err != nil {
		return nil, err
	}

	if _, err = s.Do(req, &out, true); err != nil {
		return nil, err
	}

	return out, nil
}

func (s *APIv1) APIKeyDetail(ctx context.Context, id string) (out *APIKey, err error) {
	endpoint := fmt.Sprintf("/v1/apikeys/%s", id)

	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodGet, endpoint, nil, nil); err != nil {
		return nil, err
	}

	if _, err = s.Do(req, &out, true); err != nil {
		return nil, err
	}

	return out, nil
}

func (s *APIv1) APIKeyUpdate(ctx context.Context, in *APIKey) (out *APIKey, err error) {
	endpoint := fmt.Sprintf("/v1/apikeys/%s", in.ID.String())

	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodPut, endpoint, in, nil); err != nil {
		return nil, err
	}

	if _, err = s.Do(req, &out, true); err != nil {
		return nil, err
	}

	return out, nil
}

func (s *APIv1) APIKeyDelete(ctx context.Context, id string) (err error) {
	endpoint := fmt.Sprintf("/v1/apikeys/%s", id)

	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodDelete, endpoint, nil, nil); err != nil {
		return err
	}

	if _, err = s.Do(req, nil, true); err != nil {
		return err
	}

	return nil
}

//===========================================================================
// Project Resource
//===========================================================================

func (s *APIv1) ProjectCreate(ctx context.Context, in *Project) (out *Project, err error) {
	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodPost, "/v1/projects", in, nil); err != nil {
		return nil, err
	}

	if _, err = s.Do(req, &out, true); err != nil {
		return nil, err
	}

	return out, nil
}

//===========================================================================
// Users Resource
//===========================================================================

func (s *APIv1) UserUpdate(ctx context.Context, in *User) (out *User, err error) {
	endpoint := fmt.Sprintf("/v1/users/%s", in.UserID.String())

	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodPut, endpoint, in, nil); err != nil {
		return nil, err
	}

	if _, err = s.Do(req, &out, true); err != nil {
		return nil, err
	}

	return out, nil
}

func (s *APIv1) UserList(ctx context.Context, in *UserPageQuery) (out *UserList, err error) {
	var params url.Values
	if params, err = query.Values(in); err != nil {
		return nil, fmt.Errorf("could not encode query params: %s", err)
	}

	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodGet, "/v1/users", nil, &params); err != nil {
		return nil, err
	}

	if _, err = s.Do(req, &out, true); err != nil {
		return nil, err
	}

	return out, nil
}

//===========================================================================
// Helper Methods
//===========================================================================

const (
	userAgent    = "Quarterdeck API Client/v1"
	accept       = "application/json"
	acceptLang   = "en-US,en"
	acceptEncode = "gzip, deflate, br"
	contentType  = "application/json; charset=utf-8"
)

func (s *APIv1) NewRequest(ctx context.Context, method, path string, data interface{}, params *url.Values) (req *http.Request, err error) {
	// Resolve the URL reference from the path
	url := s.endpoint.ResolveReference(&url.URL{Path: path})
	if params != nil && len(*params) > 0 {
		url.RawQuery = params.Encode()
	}

	var body io.ReadWriter
	switch {
	case data == nil:
		body = nil
	default:
		body = &bytes.Buffer{}
		if err = json.NewEncoder(body).Encode(data); err != nil {
			return nil, fmt.Errorf("could not serialize request data as json: %s", err)
		}
	}

	// Create the http request
	if req, err = http.NewRequestWithContext(ctx, method, url.String(), body); err != nil {
		return nil, fmt.Errorf("could not create request: %s", err)
	}

	// Set the headers on the request
	req.Header.Add("User-Agent", userAgent)
	req.Header.Add("Accept", accept)
	req.Header.Add("Accept-Language", acceptLang)
	req.Header.Add("Accept-Encoding", acceptEncode)
	req.Header.Add("Content-Type", contentType)

	// Use credentials from the client object unless they are available in the context
	var (
		ok    bool
		creds Credentials
	)
	if creds, ok = CredsFromContext(ctx); !ok {
		creds = s.creds
	}

	// Add authentication if it's available (add Authorization header)
	if creds != nil {
		var token string
		if token, err = creds.AccessToken(); err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", "Bearer "+token)
	}

	// Add CSRF protection if its available
	if s.client.Jar != nil {
		cookies := s.client.Jar.Cookies(url)
		for _, cookie := range cookies {
			if cookie.Name == "csrf_token" {
				req.Header.Add("X-CSRF-TOKEN", cookie.Value)
			}
		}
	}

	return req, nil
}

// Do executes an http request against the server, performs error checking, and
// deserializes the response data into the specified struct.
func (s *APIv1) Do(req *http.Request, data interface{}, checkStatus bool) (rep *http.Response, err error) {
	if rep, err = s.client.Do(req); err != nil {
		return rep, fmt.Errorf("could not execute request: %s", err)
	}
	defer rep.Body.Close()

	// Detect http status errors if they've occurred
	if checkStatus {
		if rep.StatusCode < 200 || rep.StatusCode >= 300 {
			// Attempt to read the error response from JSON, if available
			serr := &StatusError{
				StatusCode: rep.StatusCode,
			}

			if err = json.NewDecoder(rep.Body).Decode(&serr.Reply); err == nil {
				return rep, serr
			}

			serr.Reply = unsuccessful
			return rep, serr
		}
	}

	// Deserialize the JSON data from the body
	if data != nil && rep.StatusCode >= 200 && rep.StatusCode < 300 && rep.StatusCode != http.StatusNoContent {
		// Check the content type to ensure data deserialization is possible
		if ct := rep.Header.Get("Content-Type"); ct != contentType {
			return rep, fmt.Errorf("unexpected content type: %q", ct)
		}

		if err = json.NewDecoder(rep.Body).Decode(data); err != nil {
			return nil, fmt.Errorf("could not deserialize response data: %s", err)
		}
	}

	return rep, nil
}
