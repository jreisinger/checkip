package checkip

import (
	"net/http"
	"net/url"
	"time"
)

const timeout = 5 * time.Second

// newHTTPClient creates an HTTP client. Clients and Transports are safe for
// concurrent use by multiple goroutines and for efficiency should only be
// created once and re-used. See https://golang.org/pkg/net/http/ for more.
func newHTTPClient(timeout time.Duration) *http.Client {
	return &http.Client{Timeout: timeout}
}

func makeAPIcall(apiurl string, headers map[string]string, queryParams map[string]string) (*http.Response, error) {
	var r *http.Response

	apiURL, err := url.Parse(apiurl)
	if err != nil {
		return r, err
	}

	// Set query parameters.
	vals := url.Values{}
	for k, v := range queryParams {
		vals.Add(k, v)
	}
	apiURL.RawQuery = vals.Encode()

	req, err := http.NewRequest("GET", apiURL.String(), nil)
	if err != nil {
		return r, err
	}

	// Set request headers.
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	client := newHTTPClient(timeout)
	r, err = client.Do(req)
	return r, err
}
