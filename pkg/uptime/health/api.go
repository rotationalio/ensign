package health

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"hash/fnv"
	"mime"
	"net/http"
	"time"

	"github.com/vmihailenco/msgpack/v5"
)

// APIMonitor extends the HTTPMonitor to read JSON responses from the request to detect
// version and uptime information as well as to read the service state directly from
// the response of the server. Standard Rotational services generally provide more
// details at their status endpoints that can be parsed into this struct.
//
// The API monitor will parse the response body from any status code so long as the
// response contains a
type APIMonitor struct {
	HTTPMonitor
}

var _ Monitor = &APIMonitor{}

func NewAPIMonitor(endpoint string, opts ...MonitorOption) (mon *APIMonitor, err error) {
	var client *HTTPMonitor
	if client, err = NewHTTPMonitor(endpoint, opts...); err != nil {
		return nil, err
	}

	mon = &APIMonitor{
		HTTPMonitor: *client,
	}
	return mon, nil
}

// Executes an HTTP request to the wrapped endpoint and creates an APIServiceStatus
// with the status code and parsed API response body or a client error. This method does
// not return request errors (e.g. on timeout or could not connect) since this is a
// state signal. If an error is returned it is because the status check could not be
// executed at all. If the response body cannot be parsed then the status code is used
// to determine the service status.
func (h *APIMonitor) Status(ctx context.Context) (_ ServiceStatus, err error) {
	// Create a request to the specified endpoint with the specified method
	var req *http.Request
	if req, err = h.NewRequest(ctx); err != nil {
		return nil, err
	}

	// Set content-based headers on the request
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Accept-Language", "en-US,en")

	// Create the service state to return with the timestamp of the request
	state := &APIServiceStatus{}
	state.Timestamp = time.Now()

	// Execute the request and check for errors
	var rep *http.Response
	if rep, err = h.Do(req); err != nil {
		// Save the error on the state status
		state.Error = err.Error()
		state.ErrorType = h.CheckError(err)
		return state, nil
	}
	defer rep.Body.Close()

	// If the content type is JSON and the status code is not 204 attempt to parse the
	// body response from JSON into a generic map in the service status.
	if rep.StatusCode != http.StatusNoContent {
		media, _, _ := mime.ParseMediaType(rep.Header.Get("Content-Type"))
		if media == "application/json" {
			state.Content = make(map[string]interface{})
			if err = json.NewDecoder(rep.Body).Decode(&state.Content); err != nil {
				// Do not return errors, simply set the error as part of the content.
				state.Content["parse_error"] = err.Error()
			}
		}
	}

	state.StatusCode = rep.StatusCode
	return state, nil
}

type APIServiceStatus struct {
	HTTPServiceStatus
	Content map[string]interface{} `msgpack:"content"`
}

// Unmarshal from msgpack binary data.
func (h *APIServiceStatus) Unmarshal(data []byte) error {
	return msgpack.Unmarshal(data, h)
}

// Marshal to msgpack binary data for storage.
func (h *APIServiceStatus) Marshal() ([]byte, error) {
	return msgpack.Marshal(h)
}

// Hashes the status code, parsed status, version, error, and error type for
// comparison purposes.
func (h *APIServiceStatus) Hash() []byte {
	if h.hash == nil {
		sig := fnv.New128()

		// Write the Status code
		code := make([]byte, 2)
		binary.LittleEndian.PutUint16(code, uint16(h.StatusCode))
		sig.Write(code)

		// Write the parsed status
		state, _ := h.ParseStatus()
		binstate := make([]byte, 2)
		binary.LittleEndian.PutUint16(binstate, uint16(state))
		sig.Write(binstate)

		// Write the version
		sig.Write([]byte(h.Version()))

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

func (h *APIServiceStatus) Status() Status {
	// If the status can be parsed from the content then return that status.
	if status, err := h.ParseStatus(); err == nil {
		return status
	}

	// Fallback to using the HTTP status or error.
	return h.HTTPServiceStatus.Status()
}

// ParseStatus from the API response, returning an error if the response did not contain
// a status value or if the status in the response could not be parsed.
func (h *APIServiceStatus) ParseStatus() (Status, error) {
	if len(h.Content) == 0 {
		return Unknown, ErrNoContent
	}

	if status, ok := h.Content["status"]; ok {
		if s, isstr := status.(string); isstr {
			return ParseStatus(s)
		}
		return Unknown, ErrUnparsableStatus
	}
	return Unknown, ErrNoStatusResponse
}

// Version attempts to get version information from the API response.
func (h *APIServiceStatus) Version() string {
	if len(h.Content) > 0 {
		if version, ok := h.Content["version"]; ok {
			if s, isstr := version.(string); isstr {
				return s
			}
		}
	}
	return ""
}
