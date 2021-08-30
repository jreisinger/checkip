// Package checkip checks an IP address using various public services.
package checkip

import (
	"fmt"
	"net"
	"net/http"
	"time"
)

// Checker checks an IP address. It returns false if it considers the IP address
// to be suspicious. You can print the Checker to see what it has found about
// the IP address.
type Checker interface {
	Check(ip net.IP) (ok bool, err error)
	fmt.Stringer
}

// newHTTPClient creates an HTTP client. Clients and Transports are safe for
// concurrent use by multiple goroutines and for efficiency should only be
// created once and re-used. See https://golang.org/pkg/net/http/ for more.
func newHTTPClient(timeout time.Duration) *http.Client {
	return &http.Client{Timeout: timeout}
}
