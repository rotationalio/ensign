package health

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"hash/fnv"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"

	"github.com/vmihailenco/msgpack/v5"
)

// HTTPMonitor is a simple status check that uses the status code to determine the
// state that the endpoint is in. All that is compared is the status code from the
// response -- no data in the response is parsed. This check can be used to embed more
// detailed status checks, however.
type HTTPMonitor struct {
	endpoint *url.URL
	method   string
	client   *http.Client
}

var _ Monitor = &HTTPMonitor{}

func NewHTTPMonitor(endpoint string, opts ...MonitorOption) (mon *HTTPMonitor, err error) {
	var conf *Options
	if conf, err = NewOptions(opts...); err != nil {
		return nil, err
	}

	mon = &HTTPMonitor{
		client: conf.HTTPClient,
		method: conf.HTTPMethod,
	}

	// Parse the URL endpoint
	if mon.endpoint, err = url.Parse(endpoint); err != nil {
		return nil, err
	}

	// Create the default client
	if mon.client == nil {
		mon.client = &http.Client{
			Transport:     nil,
			CheckRedirect: mon.CheckRedirect,
			Timeout:       5 * time.Second,
		}

		if mon.client.Jar, err = cookiejar.New(nil); err != nil {
			return nil, fmt.Errorf("could not create cookeijar: %w", err)
		}
	}

	// Use the default method if none is specified
	if mon.method == "" {
		mon.method = http.MethodGet
	}

	return mon, nil
}

// Executes an HTTP request to the wrapped endpoint and creates an HTTPServiceStatus
// with the resulting status code or error. This method does not return request errors,
// e.g. if the client request fails or times out (since this is a state signal). If an
// error is returned from this method it is because the status check could not be made
// at all, e.g. the request could not be created to even execute a request.
func (h *HTTPMonitor) Status(ctx context.Context) (_ ServiceStatus, err error) {
	// Create a request to the specified endpoint with the specified method
	var req *http.Request
	if req, err = h.NewRequest(ctx); err != nil {
		return nil, err
	}

	state := &HTTPServiceStatus{
		Timestamp: time.Now(),
	}

	var rep *http.Response
	if rep, err = h.Do(req); err != nil {
		// Save the error on the state status
		state.Error = err.Error()
		state.ErrorType = h.CheckError(err)
		return state, nil
	}
	rep.Body.Close()

	state.StatusCode = rep.StatusCode
	return state, nil
}

// Create a new HTTP request with the specified headers
func (h *HTTPMonitor) NewRequest(ctx context.Context) (req *http.Request, err error) {
	if req, err = http.NewRequestWithContext(ctx, h.method, h.endpoint.String(), nil); err != nil {
		return nil, fmt.Errorf("could not create request: %w", err)
	}

	// Set the headers on the request
	// TODO: set any additional headers on the request as necessary
	req.Header.Add("User-Agent", "Rotational Uptime Monitor/v1")
	return req, nil
}

// Execute an HTTP request with the monitor client - useful for embeddings.
func (h *HTTPMonitor) Do(req *http.Request) (*http.Response, error) {
	return h.client.Do(req)
}

// If the server returns a 301 or 308 then update the monitor's endpoint.
func (h *HTTPMonitor) CheckRedirect(req *http.Request, via []*http.Request) error {
	// Determine the status code of the last request
	if len(via) > 1 {
		code := via[len(via)-1].Response.StatusCode
		if code == http.StatusMovedPermanently || code == http.StatusPermanentRedirect {
			// Update the endpoint to prevent redirects in the next status request.
			h.endpoint = req.URL
		}
	}

	// Default redirect behavior: only allow 10 requests before returning an error.
	if len(via) >= 10 {
		return ErrTooManyRedirects
	}
	return nil
}

// CheckError attempts to determine the status of a service based on the error returned
// by executing the request without a successful response.
func (h *HTTPMonitor) CheckError(err error) Status {
	switch {
	case errors.Is(err, ErrTooManyRedirects):
		return Unhealthy
	case errors.Is(err, context.DeadlineExceeded):
		return Degraded
	}
	return Offline
}

// HTTPServiceStatus determines the health of a service by only its status code or by
// the presence of an error if it is returned by the client (e.g. could not connect).
type HTTPServiceStatus struct {
	StatusCode int       `msgpack:"status_code"`
	Error      string    `msgpack:"error"`
	ErrorType  Status    `msgpack:"error_type"`
	Timestamp  time.Time `msgpack:"timestamp"`
	hash       []byte
}

var _ ServiceStatus = &HTTPServiceStatus{}

// Unmarshal from msgpack binary data.
func (h *HTTPServiceStatus) Unmarshal(data []byte) error {
	return msgpack.Unmarshal(data, h)
}

// Marshal to msgpack binary data for storage.
func (h *HTTPServiceStatus) Marshal() ([]byte, error) {
	return msgpack.Marshal(h)
}

// Hashes the status code, error, and error type for comparison purposes.
func (h *HTTPServiceStatus) Hash() []byte {
	if h.hash == nil {
		sig := fnv.New128()

		// Write the Status code
		code := make([]byte, 2)
		binary.LittleEndian.PutUint16(code, uint16(h.StatusCode))
		sig.Write(code)

		// Write the error
		sig.Write([]byte(h.Error))

		// Write the error type
		etype := make([]byte, 2)
		binary.LittleEndian.PutUint16(etype, uint16(h.ErrorType))
		sig.Write(etype)

		buf := make([]byte, 0, 16)
		h.hash = sig.Sum(buf)
	}
	return h.hash
}

func (h *HTTPServiceStatus) Status() Status {
	// If we have neither a status code nor an error then return unknown.
	if h.StatusCode == 0 && h.Error == "" {
		return Unknown
	}

	// If we have an error then determine the error from the error type otherwise simply
	// report offline since we were unable to connect and make a request to the server.
	if h.Error != "" {
		if h.ErrorType > 0 {
			return h.ErrorType
		}
		return Offline
	}

	// Determine the status of the server from the status code returned.
	switch {
	case h.StatusCode >= 200 && h.StatusCode < 300:
		return Online
	case h.StatusCode == http.StatusTooManyRequests || h.StatusCode == http.StatusRequestTimeout:
		return Degraded
	case h.StatusCode == http.StatusNotFound:
		// This can happen if a proxy isn't configured correctly for the service
		return Unhealthy
	case h.StatusCode == http.StatusServiceUnavailable:
		return Maintenance
	case h.StatusCode == http.StatusBadGateway || h.StatusCode == http.StatusGatewayTimeout:
		return Offline
	default:
		return Unhealthy
	}
}

func (h *HTTPServiceStatus) CheckedAt() time.Time {
	return h.Timestamp
}
