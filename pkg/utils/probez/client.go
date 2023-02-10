package probez

import (
	"context"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

type Probe interface {
	Live(ctx context.Context) (bool, int, error)
	Ready(ctx context.Context) (bool, int, error)
	Healthy(ctx context.Context) (bool, int, error)
}

func NewProbe(endpoint string) (_ Probe, err error) {
	p := &probe{
		client: &http.Client{
			Transport:     nil,
			CheckRedirect: nil,
			Timeout:       5 * time.Second,
		},
	}

	if p.client.Jar, err = cookiejar.New(nil); err != nil {
		return nil, err
	}

	if p.baseURL, err = url.Parse(endpoint); err != nil {
		return nil, err
	}

	return p, nil
}

type probe struct {
	client  *http.Client
	baseURL *url.URL
}

func (p *probe) Live(ctx context.Context) (bool, int, error) {
	return p.do(ctx, "/livez")
}

func (p *probe) Ready(ctx context.Context) (bool, int, error) {
	return p.do(ctx, "/readyz")
}

func (p *probe) Healthy(ctx context.Context) (bool, int, error) {
	return p.do(ctx, "/healthz")
}

const (
	userAgent    = "Probe/v1"
	accept       = "text/plain"
	acceptEncode = "gzip, deflate, br"
)

func (p *probe) do(ctx context.Context, path string) (z bool, status int, err error) {
	u := p.baseURL.ResolveReference(&url.URL{Path: path})

	var req *http.Request
	if req, err = http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil); err != nil {
		return false, 0, err
	}

	req.Header.Add("User-Agent", userAgent)
	req.Header.Add("Accept", accept)
	req.Header.Add("Accept-Encoding", acceptEncode)

	var rep *http.Response
	if rep, err = p.client.Do(req); err != nil {
		return false, 0, err
	}
	defer rep.Body.Close()

	z = rep.StatusCode >= 200 && rep.StatusCode < 300
	return z, rep.StatusCode, nil
}
