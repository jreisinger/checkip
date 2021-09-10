package checkip

import (
	"net/http"
	"time"
)

// newHTTPClient creates an HTTP client. Clients and Transports are safe for
// concurrent use by multiple goroutines and for efficiency should only be
// created once and re-used. See https://golang.org/pkg/net/http/ for more.
func newHTTPClient(timeout time.Duration) *http.Client {
	return &http.Client{Timeout: timeout}
}
