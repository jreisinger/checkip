// Package cli contains functions for running checkip.Checks from command-line.
package cli

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net"
	"os"
	"sort"
	"sync"

	"github.com/jreisinger/checkip"
)

// Run runs Checks concurrently against the ippaddr.
func Run(checks []checkip.Check, ipaddr net.IP) (Results, []error) {
	var results Results
	var errors []error

	var wg sync.WaitGroup
	for _, chk := range checks {
		wg.Add(1)
		go func(c checkip.Check) {
			defer wg.Done()
			r, err := c(ipaddr)
			if err != nil {
				errors = append(errors, err)
				return
			}
			results = append(results, r)
		}(chk)
	}
	wg.Wait()
	return results, errors
}

// Results are generic or security information provided by of a Check.
type Results []checkip.Result

// PrintJSON prints all Results in JSON.
func (rs Results) PrintJSON() {
	if len(rs) == 0 {
		return
	}
	out := struct {
		Check Results `json:"checks"`
	}{
		rs,
	}
	enc := json.NewEncoder(os.Stdout)
	if err := enc.Encode(out); err != nil {
		log.Fatal(err)
	}
}

func (rs Results) SortByName() {
	sort.Slice(rs, func(i, j int) bool {
		return rs[i].Name < rs[j].Name
	})
}

// PrintInfo prints summary results from Info and InfoSec checkers.
func (rs Results) PrintInfo() {
	for _, r := range rs {
		if r.Info == nil {
			continue
		}
		if r.Type == checkip.TypeInfo || r.Type == checkip.TypeInfoSec {
			fmt.Printf("%-14s --> %s\n", r.Name, r.Info.Summary())
		}
	}
}

// PrintMalicious prints how many of the InfoSec and Sec checkers consider the
// IP address to be malicious.
func (rs Results) PrintMalicious() {
	total, malicious, prob := rs.maliciousStats()
	msg := fmt.Sprintf("%-14s --> %.0f%% (%d/%d) ",
		"Malicious", math.Round(prob*100), malicious, total)
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

func (rs Results) maliciousStats() (total, malicious int, prob float64) {
	for _, r := range rs {
		if r.Info == nil {
			continue
		}
		if r.Type == checkip.TypeSec || r.Type == checkip.TypeInfoSec {
			total++
			if r.Malicious {
				malicious++
			}
		}
	}
	prob = float64(malicious) / float64(total)
	return total, malicious, prob
}
