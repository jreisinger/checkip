// Package cli contains functions for running checks from command-line.
package cli

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net"
	"os"
	"sort"
	"sync"

	"github.com/jreisinger/checkip/check"
)

// Run runs checks concurrently against the ippaddr.
func Run(checkFuncs []check.Func, ipaddr net.IP) (Checks, []error) {
	var checksMu struct {
		sync.Mutex
		checks []check.Check
	}
	var errorsMu struct {
		sync.Mutex
		errors []error
	}

	var wg sync.WaitGroup
	for _, cf := range checkFuncs {
		wg.Add(1)
		go func(cf check.Func) {
			defer wg.Done()
			c, err := cf(ipaddr)
			if err != nil {
				errorsMu.Lock()
				errorsMu.errors = append(errorsMu.errors, err)
				errorsMu.Unlock()
				return
			}
			checksMu.Lock()
			checksMu.checks = append(checksMu.checks, c)
			checksMu.Unlock()
		}(cf)
	}
	wg.Wait()

	return checksMu.checks, errorsMu.errors
}

type Checks []check.Check

// PrintJSON prints detailed results in JSON format.
func (checks Checks) PrintJSON(ipaddr net.IP) {
	// if len(rs) == 0 {
	// 	return
	// }

	_, _, prob := checks.maliciousStats()

	out := struct {
		IpAddr        net.IP `json:"ipAddr"`
		MaliciousProb string `json:"maliciousProb"`
		Checks        Checks `json:"checks"`
	}{
		ipaddr,
		fmt.Sprintf("%.2f", prob),
		checks,
	}

	enc := json.NewEncoder(os.Stdout)
	if err := enc.Encode(out); err != nil {
		log.Fatal(err)
	}
}

// SortByName sorts Results by name.
func (checks Checks) SortByName() {
	sort.Slice(checks, func(i, j int) bool {
		return checks[i].Description < checks[j].Description
	})
}

// PrintSummary prints summary results from Info and InfoSec checks.
func (checks Checks) PrintSummary() {
	for _, r := range checks {
		// To avoid "invalid memory address or nil pointer dereference"
		// runtime error and printing empty summary info.
		if r.IpAddrInfo == nil || r.IpAddrInfo.Summary() == "" {
			continue
		}

		if r.Type == check.Info || r.Type == check.InfoAndIsMalicious {
			fmt.Printf("%-15s %s\n", r.Description, r.IpAddrInfo.Summary())
		}
	}
}

// PrintMalicious prints how many of the InfoSec and Sec checks consider the IP
// address to be malicious.
func (checks Checks) PrintMalicious() {
	total, malicious, prob := checks.maliciousStats()
	msg := fmt.Sprintf("%-15s %.0f%% (%d/%d) ",
		"malicious prob.", math.Round(prob*100), malicious, total)
	switch {
	case prob >= 0.50:
		msg += `🚫`
	case prob >= 0.15:
		msg += `🤏`
	default:
		msg += `✅`
	}
	fmt.Println(msg)
}

func (checks Checks) maliciousStats() (total, malicious int, prob float64) {
	for _, r := range checks {
		// if r.Info == nil {
		// 	continue
		// }
		if r.Type == check.IsMalicious || r.Type == check.InfoAndIsMalicious {
			total++
			if r.IpAddrIsMalicious {
				malicious++
			}
		}
	}
	prob = float64(malicious) / float64(total)
	return total, malicious, prob
}

// GetIpAddrs parses IP addresses supplied as command line arguments or on
// STDIN. It keeps sending the received IP addresses down the ipaddrs channel.
// When there's no more input it closes the channel and returns.
func GetIpAddrs(args []string, ipaddrs chan<- net.IP) {
	defer close(ipaddrs)

	if len(args) == 0 { // get IP addresses from stdin.
		input := bufio.NewScanner(os.Stdin)
		for input.Scan() {
			ipaddr := net.ParseIP(input.Text())
			if ipaddr == nil {
				log.Printf("wrong IP address: %s", input.Text())
				continue
			}
			ipaddrs <- ipaddr
		}
		if err := input.Err(); err != nil {
			log.Print(err)
		}
	} else {
		for _, arg := range args {
			ipaddr := net.ParseIP(arg)
			if ipaddr == nil {
				log.Printf("wrong IP address: %s", arg)
				continue
			}
			ipaddrs <- ipaddr
		}
	}
}
