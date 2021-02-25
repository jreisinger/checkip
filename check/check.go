// Package check allows you to run various IP address checks.
package check

import (
	"fmt"
	"net"
	"net/http"
	"sort"
	"time"

	"github.com/jreisinger/checkip/util"
)

// Check represents an IP address checker.
type Check interface {
	Do(addr net.IP) (bool, error)
	Name() string
	String() string // result of the check
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
		&Shodan{},
		&ThreatCrowd{},
		&VirusTotal{},
	}
	return availableChecks
}

type checkResult struct {
	name  string
	msg   string
	notOK bool
	err   error
}

type byName []checkResult

func (r byName) Len() int           { return len(r) }
func (r byName) Swap(i, j int)      { r[j], r[i] = r[i], r[j] }
func (r byName) Less(i, j int) bool { return r[i].name < r[j].name }

func run(chk Check, ipaddr net.IP, ch chan checkResult) {
	var result checkResult
	result.name = chk.Name()
	ok, err := chk.Do(ipaddr)
	result.msg = chk.String()
	if err != nil {
		result.err = err
	}
	if !ok {
		result.notOK = true
	}
	ch <- result
}

// RunAndPrint runs concurrent checks of an IP address and prints sorted
// results. It returns the number of checks that say the IP address is not OK.
func RunAndPrint(checks []Check, ipaddr net.IP) (countNotOK int) {
	var results []checkResult

	chn := make(chan checkResult)
	for _, chk := range checks {
		go run(chk, ipaddr, chn)
	}
	for range checks {
		results = append(results, <-chn)
	}

	sort.Sort(byName(results))
	for _, r := range results {
		format := "%s %s"
		s := fmt.Sprintf(format, fmt.Sprintf("%-11s", r.name), r.msg)
		if r.err != nil {
			s = fmt.Sprintf(format, util.Lowlight(fmt.Sprintf("%-11s", r.name)), util.Lowlight(fmt.Sprintf("%s", r.err)))
		} else if r.notOK {
			s = fmt.Sprintf(format, util.Highlight(fmt.Sprintf("%-11s", r.name)), r.msg)
			countNotOK++
		}
		fmt.Println(s)
	}

	return countNotOK
}

// NewHTTPClient creates an HTTP client. Clients and Transports are safe for
// concurrent use by multiple goroutines and for efficiency should only be
// created once and re-used. See https://golang.org/pkg/net/http/ for more.
func NewHTTPClient(timeout time.Duration) *http.Client {
	return &http.Client{Timeout: timeout}
}
