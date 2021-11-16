// Package cli contains functions for checking IP addresses from CLI.
package cli

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"sort"
	"sync"

	"github.com/jreisinger/checkip/check"
)

func Run(checks []check.Check, ipaddr net.IP) (Results, []error) {
	var results Results
	var errors []error

	var wg sync.WaitGroup
	for _, chk := range checks {
		wg.Add(1)
		go func(c check.Check) {
			defer wg.Done()
			r, err := c(ipaddr)
			if err != nil {
				errors = append(errors, err)
				// log.Printf("check failed: %v", err)
				return
			}
			results = append(results, r)
		}(chk)
	}
	wg.Wait()
	return results, errors
}

type Results []check.Result

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

// PrintInfo prints results from Info and InfoSec checkers.
func (rs Results) PrintInfo() {
	for _, r := range rs {
		if r.Type == check.TypeInfo || r.Type == check.TypeInfoSec {
			fmt.Printf("%-15s %s\n", r.Name, r.Info.Summary())
		}
	}
}

// PrintProbabilityMalicious prints the probability the IP address is malicious.
func (rs Results) PrintProbabilityMalicious() {
	msg := fmt.Sprintf("%s\t%.0f%% ", "Malicious", rs.probabilityMalicious()*100)
	switch {
	case rs.probabilityMalicious() >= 0.50:
		msg += `🚫`
	case rs.probabilityMalicious() >= 0.10:
		msg += `🤏`
	default:
		msg += `✅`
	}
	fmt.Println(msg)
}

func (rs Results) probabilityMalicious() float64 {
	var malicious, totalSec float64
	for _, r := range rs {
		if r.Type == check.TypeSec || r.Type == check.TypeInfoSec {
			totalSec++
			if r.Malicious {
				malicious++
			}
		}
	}
	return malicious / totalSec
}
