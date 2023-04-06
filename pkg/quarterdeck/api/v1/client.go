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

	"github.com/cenkalti/backoff/v4"
	"github.com/google/go-querystring/query"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
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
			return nil, fmt.Errorf("could not create cookiejar: %w", err)
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

func (s *APIv1) VerifyEmail(ctx context.Context, in *VerifyRequest) (err error) {
	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodPost, "/v1/verify", in, nil); err != nil {
		return err
	}

	if _, err = s.Do(req, nil, true); err != nil {
		return err
	}

	return nil
}

//===========================================================================
// Organization Resource
//===========================================================================

func (s *APIv1) OrganizationDetail(ctx context.Context, id string) (out *Organization, err error) {
	endpoint := fmt.Sprintf("/v1/organizations/%s", id)

	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodGet, endpoint, nil, nil); err != nil {
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

func (s *APIv1) APIKeyPermissions(ctx context.Context) (out []string, err error) {
	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodGet, "/v1/apikeys/permissions", nil, nil); err != nil {
		return nil, err
	}

	if _, err = s.Do(req, &out, true); err != nil {
		return nil, err
	}

	return out, nil
}

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

func (s *APIv1) ProjectAccess(ctx context.Context, in *Project) (out *LoginReply, err error) {
	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodPost, "/v1/projects/access", in, nil); err != nil {
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

func (s *APIv1) UserDetail(ctx context.Context, id string) (out *User, err error) {
	endpoint := fmt.Sprintf("/v1/users/%s", id)

	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodGet, endpoint, nil, nil); err != nil {
		return nil, err
	}

	if _, err = s.Do(req, &out, true); err != nil {
		return nil, err
	}

	return out, nil
}

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

func (s *APIv1) UserDelete(ctx context.Context, id string) (err error) {
	endpoint := fmt.Sprintf("/v1/users/%s", id)

	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodDelete, endpoint, nil, nil); err != nil {
		return err
	}

	if _, err = s.Do(req, nil, true); err != nil {
		return err
	}

	return nil
}

func (s *APIv1) InvitePreview(ctx context.Context, token string) (out *UserInvitePreview, err error) {
	endpoint := fmt.Sprintf("/v1/invites/%s", token)

	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodGet, endpoint, nil, nil); err != nil {
		return nil, err
	}

	if _, err = s.Do(req, &out, true); err != nil {
		return nil, err
	}

	return out, nil
}

func (s *APIv1) InviteCreate(ctx context.Context, in *UserInviteRequest) (out *UserInviteReply, err error) {
	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodPost, "/v1/invites", in, nil); err != nil {
		return nil, err
	}

	if _, err = s.Do(req, &out, true); err != nil {
		return nil, err
	}

	return out, nil
}

//===========================================================================
// Accounts Resource
// This endpoint is intended for users to manage their own account
//===========================================================================

func (s *APIv1) AccountUpdate(ctx context.Context, in *User) (out *User, err error) {
	endpoint := fmt.Sprintf("/v1/accounts/%s", in.UserID.String())

	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodPut, endpoint, in, nil); err != nil {
		return nil, err
	}

	if _, err = s.Do(req, &out, true); err != nil {
		return nil, err
	}

	return out, nil
}

//===========================================================================
// Client Utility Methods
//===========================================================================

// Wait for ready polls the Quarterdeck status endpoint until it responds with an 200
// response, retrying with exponential backoff or until the context deadline is expired.
// If the user does not supply a context with a deadline, then a default deadline of
// 5 minutes is used so that this method does not block indefinitely. If the Quarterdeck
// service is ready (e.g. responds to a status request) then no error is returned,
// otherwise an error is returned if Quarterdeck never responds.
//
// NOTE: if Quarterdeck returns a 503 Service Unavailable because it is in maintenance
// mode, this method will continue to wait until the deadline for Quarterdeck to exit
// from maintenance mode and be ready again.
func (s *APIv1) WaitForReady(ctx context.Context) (err error) {
	// If context does not have a deadline, create a context with a default deadline.
	if _, hasDeadline := ctx.Deadline(); !hasDeadline {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 5*time.Minute)
		defer cancel()
	}

	// Create the status request to send until ready
	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodGet, "/v1/status", nil, nil); err != nil {
		return err
	}

	// Create a closure to repeatedly call the Quarterdeck status endpoint
	checkReady := func() (err error) {
		var rep *http.Response
		if rep, err = s.client.Do(req); err != nil {
			return err
		}
		defer rep.Body.Close()

		if rep.StatusCode < 200 || rep.StatusCode >= 300 {
			return &StatusError{StatusCode: rep.StatusCode, Reply: Reply{Success: false, Error: http.StatusText(rep.StatusCode)}}
		}
		return nil
	}

	// Create exponential backoff ticker for retries
	ticker := backoff.NewExponentialBackOff()

	// Keep checking if Quarterdeck is ready until it is ready or until the context expires.
	for {
		// Execute the status request
		if err = checkReady(); err == nil {
			// Success - Quarterdeck is ready for requests!
			return nil
		}

		// Log the error warning that we're still waiting to connect to quarterdeck
		log.Warn().Err(err).Msg("waiting to connect to quarterdeck")
		wait := time.After(ticker.NextBackOff())

		// Wait for the context to be done or for the ticker to move to the next backoff.
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-wait:
		}
	}
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

	// If there is a request ID on the context, set it on the request, otherwise generate one
	var requestID string
	if requestID, _ = RequestIDFromContext(ctx); requestID == "" {
		requestID = ulid.Make().String()
	}
	req.Header.Add("X-Request-ID", requestID)

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
