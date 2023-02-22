package health

import (
	"fmt"
	"net/http"
)

// MonitorOption are used to configure monitors to make correct requests. MonitorOptions
// are used across different monitors and not all options are applicable. See the
// docstring of the monitor to determine which options are available for it.
type MonitorOption func(*Options) error

// Options are configurations provided to Monitors.
type Options struct {
	HTTPClient *http.Client // Used to specify an HTTP client to the HTTP monitor.
	HTTPMethod string       // Used to specify the method to make status requests with.
}

func NewOptions(opts ...MonitorOption) (*Options, error) {
	options := &Options{}
	for _, opt := range opts {
		if err := opt(options); err != nil {
			return nil, err
		}
	}
	return options, nil
}

// Use the specified http client instead of the default http client.
func WithHTTPClient(client *http.Client) MonitorOption {
	return func(o *Options) error {
		o.HTTPClient = client
		return nil
	}
}

// Use the specified http method instead of GET
func WithHTTPMethod(method string) MonitorOption {
	return func(o *Options) error {
		allowed := map[string]struct{}{
			http.MethodGet:     {},
			http.MethodHead:    {},
			http.MethodPost:    {},
			http.MethodPut:     {},
			http.MethodPatch:   {},
			http.MethodConnect: {},
			http.MethodTrace:   {},
		}

		if _, ok := allowed[method]; !ok {
			return fmt.Errorf("method %q is not an allowed HTTP method", method)
		}

		o.HTTPMethod = method
		return nil
	}
}
