// Package cli contains functions for checking IP addresses from CLI tools.
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
	"github.com/logrusorgru/aurora"
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
	enc := json.NewEncoder(os.Stdout)
	if err := enc.Encode(rs); err != nil {
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
		if r.Type == "Info" || r.Type == "InfoSec" {
			fmt.Printf("%-15s %s\n", r.Name, r.Info.String())
		}
	}
}

// PrintProbabilityMalicious prints the probability the IP address is malicious.
func (rs Results) PrintProbabilityMalicious() {
	var msg string
	switch {
	case rs.probabilityMalicious() <= 0.15:
		msg = fmt.Sprint(aurora.Green("Malicious"))
	case rs.probabilityMalicious() <= 0.50:
		msg = fmt.Sprint(aurora.Yellow("Malicious"))
	default:
		msg = fmt.Sprint(aurora.Red("Malicious"))
	}

	fmt.Printf("%s\t%.0f%%\n", msg, rs.probabilityMalicious()*100)
}

func (rs Results) probabilityMalicious() float64 {
	var malicious, totalSec float64
	for _, r := range rs {
		if r.Type == "Sec" || r.Type == "InfoSec" {
			totalSec++
			if r.IPaddrMalicious {
				malicious++
			}
		}
	}
	return malicious / totalSec
}
