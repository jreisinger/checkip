// Package checkip checks an IP address using various public services.
package checkip

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/logrusorgru/aurora"
)

// Checker checks an IP address. It returns false if it considers the IP address
// to be suspicious. You can print the Checker to see what it has found about
// the IP address.
type Checker interface {
	Check(ip net.IP) (ok bool, err error)
	fmt.Stringer
}

// Run runs checkers concurrently and returns the number of checkers that
// consider the IP address suspicious.
func Run(checkers []Checker, ipaddr net.IP) int {
	ch := make(chan bool)
	for _, checker := range checkers {
		go func(checker Checker) {
			ok, err := checker.Check(ipaddr)
			if err == nil {
				ch <- ok
			}
		}(checker)
	}
	var suspicious int
	for range checkers {
		ok := <-ch
		if !ok {
			suspicious++
		}
	}
	return suspicious
}

// RunAndPrint runs checkers concurrently and print the results. checkers maps a
// name to a checker. Format defines how to print the name and checker results
// (e.g. "%-25s %s").
func RunAndPrint(checkers map[string]Checker, ipaddr net.IP, format string) {
	ch := make(chan string)
	for name, checker := range checkers {
		go func(checker Checker, name string) {
			ok, err := checker.Check(ipaddr)
			switch {
			case err != nil:
				ch <- fmt.Sprintf(format, name, aurora.Gray(11, err.Error()))
			case !ok:
				ch <- fmt.Sprintf(format, name, aurora.Magenta(checker.String()))
			default:
				ch <- fmt.Sprintf(format, name, checker)
			}
		}(checker, name)
	}
	for range checkers {
		fmt.Println(<-ch)
	}
}

// newHTTPClient creates an HTTP client. Clients and Transports are safe for
// concurrent use by multiple goroutines and for efficiency should only be
// created once and re-used. See https://golang.org/pkg/net/http/ for more.
func newHTTPClient(timeout time.Duration) *http.Client {
	return &http.Client{Timeout: timeout}
}
