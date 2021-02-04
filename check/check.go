// Package check allows you to run various IP address checks.
package check

import (
	"encoding/json"
	"fmt"
	"log"
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

type JSONOutput struct {
	IPaddr  net.IP
	Results []Result
}

type Result struct {
	Name  string
	Msg   string
	NotOK bool
	Err   error
}

type byName []Result

func (r byName) Len() int           { return len(r) }
func (r byName) Swap(i, j int)      { r[j], r[i] = r[i], r[j] }
func (r byName) Less(i, j int) bool { return r[i].Name < r[j].Name }

func run(chk Check, ipaddr net.IP, ch chan Result) {
	var result Result
	result.Name = chk.Name()
	ok, err := chk.Do(ipaddr)
	result.Msg = chk.String()
	if err != nil {
		result.Err = err
	}
	if !ok {
		result.NotOK = true
	}
	ch <- result
}

// CountNotOK is the numbers of checks that say that an IP address is not ok.
var CountNotOK int

// RunAndPrint runs concurrent checks of an IP address and prints sorted
// results. It updates CountNotOK when a check says the IP address is not OK.
func RunAndPrint(checks []Check, ipaddr net.IP, ch chan string) {
	var results []Result

	chn := make(chan Result)
	for _, chk := range checks {
		go run(chk, ipaddr, chn)
	}
	for range checks {
		results = append(results, <-chn)
	}

	s := fmt.Sprintf("----------- %15s ----------\n", ipaddr)

	sort.Sort(byName(results))
	for _, r := range results {
		format := "%s %s\n"
		if r.Err != nil {
			s += fmt.Sprintf(format, util.Lowlight(fmt.Sprintf("%-11s", r.Name)), util.Lowlight(fmt.Sprintf("%s", r.Err)))
		} else if r.NotOK {
			s += fmt.Sprintf(format, util.Highlight(fmt.Sprintf("%-11s", r.Name)), r.Msg)
			CountNotOK++
		} else {
			s += fmt.Sprintf(format, fmt.Sprintf("%-11s", r.Name), r.Msg)
		}
	}

	ch <- s
}

func RunAndJSON(checks []Check, ipaddr net.IP, ch chan string) {
	var j JSONOutput
	var results []Result

	chn := make(chan Result)
	for _, chk := range checks {
		go run(chk, ipaddr, chn)
	}
	for range checks {
		results = append(results, <-chn)
	}

	j.IPaddr = ipaddr
	j.Results = results

	js, err := json.Marshal(j)
	if err != nil {
		log.Fatal(err)
	}

	ch <- fmt.Sprintf("%s", js)
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

// NewHTTPClient creates an HTTP client. Clients and Transports are safe for
// concurrent use by multiple goroutines and for efficiency should only be
// created once and re-used. See https://golang.org/pkg/net/http/ for more.
func NewHTTPClient(timeout time.Duration) *http.Client {
	return &http.Client{Timeout: timeout}
}
