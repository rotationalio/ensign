package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"

	"github.com/google/go-querystring/query"
	"github.com/oklog/ulid/v2"
)

// New creates a new API v1 client that implements the Tenant Client interface.
func New(endpoint string, opts ...ClientOption) (_ TenantClient, err error) {
	c := &APIv1{}
	if c.endpoint, err = url.Parse(endpoint); err != nil {
		return nil, fmt.Errorf("could not parse endpoint: %s", err)
	}

	// Applies our options
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

		// Creates cookie jar for CSRF
		if c.client.Jar, err = cookiejar.New(nil); err != nil {
			return nil, fmt.Errorf("could not create cookiejar: %s", err)
		}
	}

	return c, nil
}

// APIv1 implements the TenantClient interface
type APIv1 struct {
	endpoint *url.URL
	client   *http.Client
	creds    string
}

// Ensures the APIv1 implements the TenantClient interface
var _ TenantClient = &APIv1{}

//===========================================================================
// Client Methods
//===========================================================================

func (s *APIv1) Status(ctx context.Context) (out *StatusReply, err error) {
	// Make the HTTP request
	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodGet, "/v1/status", nil, nil); err != nil {
		return nil, err
	}

	// NOTE: We cannot use s.Do because we want to parse 503 Unavailable errors
	var rep *http.Response
	if rep, err = s.client.Do(req); err != nil {
		return nil, err
	}

	defer rep.Body.Close()

	// Detects other erros
	if rep.StatusCode != http.StatusOK && rep.StatusCode != http.StatusServiceUnavailable {
		return nil, fmt.Errorf("%s", rep.Status)
	}

	// Deserializes JSON data from the response
	out = &StatusReply{}
	if err = json.NewDecoder(rep.Body).Decode(out); err != nil {
		return nil, fmt.Errorf("could not deserialize status reply: %s", err)
	}

	return out, nil
}

func (s *APIv1) SignUp(ctx context.Context, in *ContactInfo) (err error) {
	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodPost, "/v1/notifications/signup", in, nil); err != nil {
		return err
	}

	if _, err = s.Do(req, nil, true); err != nil {
		return err
	}
	return nil
}

func (s *APIv1) Register(ctx context.Context, in *RegisterRequest) (err error) {
	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodPost, "/v1/register", in, nil); err != nil {
		return nil
	}

	if _, err = s.Do(req, nil, true); err != nil {
		return err
	}
	return nil
}

func (s *APIv1) Login(ctx context.Context, in *LoginRequest) (out *AuthReply, err error) {
	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodPost, "/v1/login", in, nil); err != nil {
		return nil, err
	}

	out = &AuthReply{}
	if _, err = s.Do(req, out, true); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *APIv1) Refresh(ctx context.Context, in *RefreshRequest) (out *AuthReply, err error) {
	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodPost, "/v1/refresh", in, nil); err != nil {
		return nil, err
	}

	out = &AuthReply{}
	if _, err = s.Do(req, out, true); err != nil {
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

func (s *APIv1) OrganizationDetail(ctx context.Context, id string) (out *Organization, err error) {
	if id == "" {
		return nil, ErrOrganizationIDRequired
	}

	path := fmt.Sprintf("/v1/organization/%s", id)
	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodGet, path, nil, nil); err != nil {
		return nil, err
	}

	out = &Organization{}
	if _, err = s.Do(req, out, true); err != nil {
		return nil, err
	}

	return out, nil
}

func (s *APIv1) TenantList(ctx context.Context, in *PageQuery) (out *TenantPage, err error) {
	var params url.Values
	if params, err = query.Values(in); err != nil {
		return nil, fmt.Errorf("could not encode query params: %w", err)
	}

	// Make the HTTP request
	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodGet, "/v1/tenant", nil, &params); err != nil {
		return nil, err
	}

	out = &TenantPage{}
	if _, err = s.Do(req, out, true); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *APIv1) TenantCreate(ctx context.Context, in *Tenant) (out *Tenant, err error) {
	// Make the HTTP Request
	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodPost, "/v1/tenant", in, nil); err != nil {
		return nil, err
	}

	// Make the HTTP response
	out = &Tenant{}
	var rep *http.Response
	if rep, err = s.Do(req, out, true); err != nil {
		return nil, err
	}

	if rep.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("expected status created, received %s", rep.Status)
	}
	return out, nil
}

func (s *APIv1) TenantDetail(ctx context.Context, id string) (out *Tenant, err error) {
	if id == "" {
		return nil, ErrTenantIDRequired
	}

	path := fmt.Sprintf("/v1/tenant/%s", id)

	// Make the HTTP request
	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodGet, path, nil, nil); err != nil {
		return nil, err
	}

	if _, err = s.Do(req, &out, true); err != nil {
		return nil, err
	}

	return out, nil
}

func (s *APIv1) TenantUpdate(ctx context.Context, in *Tenant) (out *Tenant, err error) {
	if in.ID == "" {
		return nil, ErrTenantIDRequired
	}

	path := fmt.Sprintf("/v1/tenant/%s", in.ID)

	// Make the HTTP request
	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodPut, path, in, nil); err != nil {
		return nil, err
	}

	if _, err = s.Do(req, &out, true); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *APIv1) TenantDelete(ctx context.Context, id string) (err error) {
	if id == "" {
		return ErrTenantIDRequired
	}

	path := fmt.Sprintf("/v1/tenant/%s", id)

	// Make the HTTP request
	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodDelete, path, nil, nil); err != nil {
		return err
	}
	if _, err = s.Do(req, nil, true); err != nil {
		return err
	}
	return nil
}

func (s *APIv1) TenantStats(ctx context.Context, id string) (out []*StatValue, err error) {
	if id == "" {
		return nil, ErrTenantIDRequired
	}

	path := fmt.Sprintf("/v1/tenant/%s/stats", id)

	// Make the HTTP request
	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodGet, path, nil, nil); err != nil {
		return nil, err
	}

	if _, err = s.Do(req, &out, true); err != nil {
		return nil, err
	}

	return out, nil
}

func (s *APIv1) MemberList(ctx context.Context, in *PageQuery) (out *MemberPage, err error) {
	var params url.Values
	if params, err = query.Values(in); err != nil {
		return nil, fmt.Errorf("could not encode query params: %w", err)
	}

	// Make the HTTP request
	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodGet, "/v1/members", nil, &params); err != nil {
		return nil, err
	}

	out = &MemberPage{}
	if _, err = s.Do(req, out, true); err != nil {
		return nil, err
	}

	return out, nil
}

func (s *APIv1) MemberCreate(ctx context.Context, in *Member) (out *Member, err error) {
	// Make the HTTP request
	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodPost, "v1/members", in, nil); err != nil {
		return nil, err
	}

	// Make the HTTP response
	out = &Member{}
	var rep *http.Response
	if rep, err = s.Do(req, out, true); err != nil {
		return nil, err
	}

	if rep.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("expected status created, received %s", rep.Status)
	}
	return out, nil
}

func (s *APIv1) MemberDetail(ctx context.Context, id string) (out *Member, err error) {
	if id == "" {
		return nil, ErrMemberIDRequired
	}

	path := fmt.Sprintf("/v1/members/%s", id)

	// Make the HTTP request
	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodGet, path, nil, nil); err != nil {
		return nil, err
	}

	if _, err = s.Do(req, &out, true); err != nil {
		return nil, err
	}

	return out, nil
}

func (s *APIv1) MemberUpdate(ctx context.Context, in *Member) (out *Member, err error) {
	if in.ID == "" {
		return nil, ErrMemberIDRequired
	}

	path := fmt.Sprintf("/v1/members/%s", in.ID)

	// Make the HTTP request
	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodPut, path, in, nil); err != nil {
		return nil, err
	}

	if _, err = s.Do(req, &out, true); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *APIv1) MemberDelete(ctx context.Context, id string) (err error) {
	if id == "" {
		return ErrMemberIDRequired
	}

	path := fmt.Sprintf("/v1/members/%s", id)

	// Make the HTTP request
	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodDelete, path, nil, nil); err != nil {
		return err
	}
	if _, err = s.Do(req, nil, true); err != nil {
		return err
	}
	return nil
}

func (s *APIv1) TenantProjectList(ctx context.Context, id string, in *PageQuery) (out *TenantProjectPage, err error) {
	if id == "" {
		return nil, ErrTenantIDRequired
	}

	path := fmt.Sprintf("v1/tenant/%s/projects", id)

	var params url.Values
	if params, err = query.Values(in); err != nil {
		return nil, fmt.Errorf("could not encode query params: %w", err)
	}

	// Make the HTTP request
	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodGet, path, nil, &params); err != nil {
		return nil, err
	}

	out = &TenantProjectPage{}
	if _, err = s.Do(req, out, true); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *APIv1) TenantProjectCreate(ctx context.Context, id string, in *Project) (out *Project, err error) {
	if id == "" {
		return nil, ErrTenantIDRequired
	}

	path := fmt.Sprintf("v1/tenant/%s/projects", id)

	// Make the HTTP request
	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodPost, path, in, nil); err != nil {
		return nil, err
	}

	out = &Project{}

	// Make the HTTP response
	var rep *http.Response
	if rep, err = s.Do(req, out, true); err != nil {
		return nil, err
	}

	if rep.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("expected status created, received %s", rep.Status)
	}

	return out, nil
}

func (s *APIv1) ProjectList(ctx context.Context, in *PageQuery) (out *ProjectPage, err error) {

	var params url.Values
	if params, err = query.Values(in); err != nil {
		return nil, fmt.Errorf("could not encode query params: %w", err)
	}

	// Make the HTTP request
	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodGet, "/v1/projects", nil, &params); err != nil {
		return nil, err
	}

	out = &ProjectPage{}
	if _, err = s.Do(req, out, true); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *APIv1) ProjectCreate(ctx context.Context, in *Project) (out *Project, err error) {
	// Make the HTTP request
	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodPost, "/v1/projects", in, nil); err != nil {
		return nil, err
	}

	out = &Project{}

	// Make the HTTP response
	var rep *http.Response
	if rep, err = s.Do(req, out, true); err != nil {
		return nil, err
	}

	if rep.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("expected status created, received %s", rep.Status)
	}
	return out, nil
}

func (s *APIv1) ProjectDetail(ctx context.Context, id string) (out *Project, err error) {
	if id == "" {
		return nil, ErrProjectIDRequired
	}

	path := fmt.Sprintf("/v1/projects/%s", id)

	// Make the HTTP request
	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodGet, path, nil, nil); err != nil {
		return nil, err
	}

	if _, err = s.Do(req, &out, true); err != nil {
		return nil, err
	}

	return out, nil
}

func (s *APIv1) ProjectUpdate(ctx context.Context, in *Project) (out *Project, err error) {
	if in.ID == "" {
		return nil, ErrProjectIDRequired
	}

	path := fmt.Sprintf("/v1/projects/%s", in.ID)

	// Make the HTTP request
	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodPut, path, in, nil); err != nil {
		return nil, err
	}

	if _, err = s.Do(req, &out, true); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *APIv1) ProjectDelete(ctx context.Context, id string) (err error) {
	path := fmt.Sprintf("/v1/projects/%s", id)

	// Make the HTTP request
	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodDelete, path, nil, nil); err != nil {
		return err
	}
	if _, err = s.Do(req, nil, true); err != nil {
		return err
	}
	return nil
}

func (s *APIv1) ProjectTopicList(ctx context.Context, id string, in *PageQuery) (out *ProjectTopicPage, err error) {
	if id == "" {
		return nil, ErrProjectIDRequired
	}

	path := fmt.Sprintf("/v1/projects/%s/topics", id)

	var params url.Values
	if params, err = query.Values(in); err != nil {
		return nil, err
	}
	// Make the HTTP request
	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodGet, path, nil, &params); err != nil {
		return nil, err
	}

	out = &ProjectTopicPage{}
	if _, err = s.Do(req, out, true); err != nil {
		return nil, err
	}

	return out, nil
}

func (s *APIv1) ProjectTopicCreate(ctx context.Context, id string, in *Topic) (out *Topic, err error) {
	if id == "" {
		return nil, ErrProjectIDRequired
	}

	path := fmt.Sprintf("/v1/projects/%s/topics", id)

	// Make the HTTP request
	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodPost, path, in, nil); err != nil {
		return nil, err
	}

	// Make the HTTP response
	out = &Topic{}
	var rep *http.Response
	if rep, err = s.Do(req, out, true); err != nil {
		return nil, err
	}

	if rep.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("expected status created, received %s", rep.Status)
	}

	return out, err
}

func (s *APIv1) TopicList(ctx context.Context, in *PageQuery) (out *TopicPage, err error) {
	var params url.Values
	if params, err = query.Values(in); err != nil {
		return nil, err
	}

	// Make the HTTP request
	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodGet, "/v1/topics", nil, &params); err != nil {
		return nil, err
	}

	out = &TopicPage{}
	if _, err = s.Do(req, out, true); err != nil {
		return nil, err
	}

	return out, nil
}

func (s *APIv1) TopicDetail(ctx context.Context, id string) (out *Topic, err error) {
	path := fmt.Sprintf("/v1/topics/%s", id)

	// Make the HTTP request
	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodGet, path, nil, nil); err != nil {
		return nil, err
	}

	if _, err = s.Do(req, &out, true); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *APIv1) TopicUpdate(ctx context.Context, in *Topic) (out *Topic, err error) {
	if in.ID == "" {
		return nil, ErrTopicIDRequired
	}

	path := fmt.Sprintf("/v1/topics/%s", in.ID)

	// Make the HTTP request
	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodPut, path, in, nil); err != nil {
		return nil, err
	}

	if _, err = s.Do(req, &out, true); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *APIv1) TopicDelete(ctx context.Context, in *Confirmation) (out *Confirmation, err error) {
	if in.ID == "" {
		return nil, ErrTopicIDRequired
	}

	path := fmt.Sprintf("/v1/topics/%s", in.ID)

	// Make the HTTP request
	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodDelete, path, in, nil); err != nil {
		return nil, err
	}
	if _, err = s.Do(req, &out, true); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *APIv1) ProjectAPIKeyList(ctx context.Context, id string, in *PageQuery) (out *ProjectAPIKeyPage, err error) {
	if id == "" {
		return nil, ErrProjectIDRequired
	}

	path := fmt.Sprintf("/v1/projects/%s/apikeys", id)

	var params url.Values
	if params, err = query.Values(in); err != nil {
		return nil, fmt.Errorf("could not encode query params: %w", err)
	}

	// Make the HTTP request
	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodGet, path, nil, &params); err != nil {
		return nil, err
	}

	out = &ProjectAPIKeyPage{}
	if _, err = s.Do(req, out, true); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *APIv1) ProjectAPIKeyCreate(ctx context.Context, id string, in *APIKey) (out *APIKey, err error) {
	if id == "" {
		return nil, ErrProjectIDRequired
	}

	path := fmt.Sprintf("/v1/projects/%s/apikeys", id)

	// Make the HTTP Request
	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodPost, path, in, nil); err != nil {
		return nil, err
	}

	// Make the HTTP response
	out = &APIKey{}
	var rep *http.Response
	if rep, err = s.Do(req, out, true); err != nil {
		return nil, err
	}

	if rep.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("expected status created, received %s", rep.Status)
	}
	return out, nil
}

func (s *APIv1) APIKeyCreate(ctx context.Context, in *APIKey) (out *APIKey, err error) {
	// Make the HTTP Request
	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodPost, "/v1/apikeys", in, nil); err != nil {
		return nil, err
	}

	// Make the HTTP response
	out = &APIKey{}
	var rep *http.Response
	if rep, err = s.Do(req, out, true); err != nil {
		return nil, err
	}

	if rep.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("expected status created, received %s", rep.Status)
	}
	return out, nil
}

func (s *APIv1) APIKeyList(ctx context.Context, in *PageQuery) (out *APIKeyPage, err error) {
	var params url.Values
	if params, err = query.Values(in); err != nil {
		return nil, fmt.Errorf("could not encode query params: %w", err)
	}

	// Make the HTTP request
	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodGet, "/v1/apikeys", nil, &params); err != nil {
		return nil, err
	}

	out = &APIKeyPage{}
	if _, err = s.Do(req, out, true); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *APIv1) APIKeyDetail(ctx context.Context, id string) (out *APIKey, err error) {
	if id == "" {
		return nil, ErrAPIKeyIDRequired
	}

	path := fmt.Sprintf("/v1/apikeys/%s", id)

	// Make the HTTP request
	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodGet, path, nil, nil); err != nil {
		return nil, err
	}

	if _, err = s.Do(req, &out, true); err != nil {
		return nil, err
	}

	return out, nil
}

func (s *APIv1) APIKeyUpdate(ctx context.Context, in *APIKey) (out *APIKey, err error) {
	if in.ID == "" {
		return nil, ErrAPIKeyIDRequired
	}

	path := fmt.Sprintf("/v1/apikeys/%s", in.ID)

	// Make the HTTP request
	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodPut, path, in, nil); err != nil {
		return nil, err
	}

	if _, err = s.Do(req, &out, true); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *APIv1) APIKeyDelete(ctx context.Context, id string) (err error) {
	if id == "" {
		return ErrAPIKeyIDRequired
	}

	path := fmt.Sprintf("/v1/apikeys/%s", id)

	// Make the HTTP request
	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodDelete, path, nil, nil); err != nil {
		return err
	}
	if _, err = s.Do(req, nil, true); err != nil {
		return err
	}
	return nil
}

func (s *APIv1) APIKeyPermissions(ctx context.Context) (out []string, err error) {
	// Make the HTTP request
	var req *http.Request
	if req, err = s.NewRequest(ctx, http.MethodGet, "/v1/apikeys/permissions", nil, nil); err != nil {
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
	userAgent    = "Tenant API Client/v1"
	accept       = "application/json"
	acceptLang   = "en-US, en"
	acceptEncode = "gzip, deflate, br"
	contentType  = "application/json; charset=utf-8"
)

func (s *APIv1) NewRequest(ctx context.Context, method, path string, data interface{}, params *url.Values) (req *http.Request, err error) {
	// Resolves the URL reference from the path
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

	// Creates the http request
	if req, err = http.NewRequestWithContext(ctx, method, url.String(), body); err != nil {
		return nil, fmt.Errorf("could not create request: %s", err)
	}

	// Sets the headers on the request
	req.Header.Add("User-Agent", userAgent)
	req.Header.Add("Accept", accept)
	req.Header.Add("Accept-Language", acceptLang)
	req.Header.Add("Accept-Encoding", acceptEncode)
	req.Header.Add("Content-Type", contentType)
	req.Header.Add("X-Request-ID", ulid.Make().String())

	// Add authentication if it is available
	if s.creds != "" {
		req.Header.Add("Authorization", "Bearer "+s.creds)
	}

	// Adds CSRF protection if it is available
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
// deserializes response data into the specified struct.
func (s *APIv1) Do(req *http.Request, data interface{}, checkStatus bool) (rep *http.Response, err error) {
	if rep, err = s.client.Do(req); err != nil {
		return rep, fmt.Errorf("could not execute request: %s", err)
	}
	defer rep.Body.Close()

	// Detects http status errors if they've occurred
	if checkStatus {
		if rep.StatusCode < 200 || rep.StatusCode >= 300 {
			// Attempt to read the error response from JSON, if available
			var reply Reply
			if err = json.NewDecoder(rep.Body).Decode(&reply); err == nil {
				if reply.Error != "" {
					return rep, fmt.Errorf("[%d] %s", rep.StatusCode, reply.Error)
				}
			}
			return rep, errors.New(rep.Status)
		}
	}

	// Deserializes the JSON data from the body
	if data != nil && rep.StatusCode >= 200 && rep.StatusCode < 300 && rep.StatusCode != http.StatusNoContent {
		// Checks the content type to ensure data deserialization is possible
		if ct := rep.Header.Get("Content-Type"); ct != contentType {
			return rep, fmt.Errorf("unexpected content type: %q", ct)
		}

		if err = json.NewDecoder(rep.Body).Decode(data); err != nil {
			return nil, fmt.Errorf("could not deserialize response data: %s", err)
		}
	}

	return rep, nil
}

// SetCredentials is a helper function for external users to override credentials at
// runtime by directly passing in the token, which is useful for testing.
// TODO: Pass in a credentials interface instead of the token string.
func (c *APIv1) SetCredentials(token string) {
	c.creds = token
}

// SetCSRFProtect is a helper function to set CSRF cookies on the client. This is not
// possible in a browser because of the HttpOnly flag. This method should only be used
// for testing purposes and an error is returned if the URL is not localhost. For live
// clients - the server should set these cookies. If protect is false, then the cookies
// are removed from the client by setting the cookies to an empty slice.
func (c *APIv1) SetCSRFProtect(protect bool) error {
	if c.client.Jar == nil {
		return errors.New("client does not have a cookie jar, cannot set cookies")
	}

	if c.endpoint.Hostname() != "127.0.0.1" && c.endpoint.Hostname() != "localhost" {
		return fmt.Errorf("csrf protect is for local testing only, cannot set cookies for %s", c.endpoint.Hostname())
	}

	// The URL for the cookies
	u := c.endpoint.ResolveReference(&url.URL{Path: "/"})

	var cookies []*http.Cookie
	if protect {
		cookies = []*http.Cookie{
			{
				Name:     "csrf_token",
				Value:    "testingcsrftoken",
				Expires:  time.Now().Add(10 * time.Minute),
				HttpOnly: false,
			},
			{
				Name:     "csrf_reference_token",
				Value:    "testingcsrftoken",
				Expires:  time.Now().Add(10 * time.Minute),
				HttpOnly: true,
			},
		}
	} else {
		cookies = c.client.Jar.Cookies(u)
		for _, cookie := range cookies {
			cookie.MaxAge = -1
		}
	}

	c.client.Jar.SetCookies(u, cookies)
	return nil
}
