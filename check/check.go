// Package check allows you to run various IP address checks.
package check

import (
	"fmt"
	"net"
	"net/http"
	"sort"
	"time"

	. "github.com/logrusorgru/aurora"
)

// Check represents an IP address checker.
type Check interface {
	Do(addr net.IP) (bool, error)
	Name() string
	String() string // result of the check
}

// RunAndFormat runs a check of an IP address and returns formated result over a
// channel. countNotOK holds the number of checkers that think the IP address is
// not OK.
func RunAndFormat(chk Check, ipaddr net.IP, ch chan string, countNotOK *int) {
	format := "%-11s %s"
	ok, err := chk.Do(ipaddr)
	if err != nil {
		ch <- fmt.Sprintf(format, chk.Name(), Gray(11, err))
		return
	}
	if ok {
		ch <- fmt.Sprintf(format, chk.Name(), chk)
	} else {
		*countNotOK++
		ch <- fmt.Sprintf(format, chk.Name(), Magenta(chk))
	}
}

// RunAndPrint runs concurrent checks of an IP address and prints sorted
// results. countNotOK holds the number of checkers that think the IP address is
// not OK.
func RunAndPrint(checks []Check, ipaddr net.IP, countNotOK *int) {
	var results []string

	chn := make(chan string)
	for _, chk := range checks {
		go RunAndFormat(chk, ipaddr, chn, countNotOK)
	}
	for range checks {
		results = append(results, <-chn)
	}

	sort.Strings(results)
	for _, result := range results {
		fmt.Println(result)
	}
}

// GetAvailable returns all available checks.
func GetAvailable() []Check {
	availableChecks := []Check{
		&AbuseIPDB{},
		&AS{},
		&DNS{},
		&Geo{},
		&IPsum{},
		&OTX{},
		&ThreatCrowd{},
		&VirusTotal{},
	}
	return availableChecks
}

// NewHTTPClient creates an HTTP client. Clients and Transports are safe for
// concurrent use by multiple goroutines and for efficiency should only be
// created once and re-used. See https://golang.org/pkg/net/http/ for more.
func NewHTTPClient(timeout time.Duration) *http.Client {
	return &http.Client{Timeout: timeout}
}
