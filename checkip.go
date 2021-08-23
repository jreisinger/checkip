// Package checkip checks an IP address using various public services.
package checkip

import (
	"fmt"
	"net"
	"net/http"
	"sort"
	"strconv"
	"time"
)

// Checker checks an IP address. It returns false if it considers the IP address
// to be suspicious. You can print the Checker to see what it has found about
// the IP address.
type Checker interface {
	Check(ip net.IP) (ok bool, err error)
	fmt.Stringer
}

// Checkers is for giving names to Checkers.
type Checkers map[string]Checker

// CheckAndPrint checks an IP address running Checkers concurrently. It prints
// the name and the info the Checker has found about the IP address. If the
// Checker considers the IP address suspicious (not ok) the info is highlighted.
// If the Checker returns en error the info is lowlighted. Resuls are
// alphabetically sorted by the Checker name.
func (checkers Checkers) CheckAndPrint(ip net.IP) {
	ch := make(chan string)
	longest := longestName(checkers)
	format := "%-" + strconv.Itoa(longest) + "s %s"
	for name, checker := range checkers {
		go checkAndPrint(ip, format, name, checker, ch)
	}
	var results []string
	for range checkers {
		results = append(results, <-ch)
	}
	sort.Strings(results)
	for _, res := range results {
		fmt.Println(res)
	}
}

func checkAndPrint(ip net.IP, format, name string, checker Checker, ch chan string) {
	ok, err := checker.Check(ip)
	switch {
	case err != nil:
		ch <- fmt.Sprintf(format, name, lowlight(err.Error()))
	case !ok:
		ch <- fmt.Sprintf(format, name, highlight(checker.String()))
	default:
		ch <- fmt.Sprintf(format, name, checker)
	}
}

func longestName(checkers Checkers) int {
	var longest int
	for name := range checkers {
		if len(name) > longest {
			longest = len(name)
		}
	}
	return longest
}

// newHTTPClient creates an HTTP client. Clients and Transports are safe for
// concurrent use by multiple goroutines and for efficiency should only be
// created once and re-used. See https://golang.org/pkg/net/http/ for more.
func newHTTPClient(timeout time.Duration) *http.Client {
	return &http.Client{Timeout: timeout}
}
