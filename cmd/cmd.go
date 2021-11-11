package cmd

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"sort"
	"sync"

	"github.com/jreisinger/checkip/check"
	"github.com/jreisinger/checkip/checker"
	"github.com/logrusorgru/aurora"
)

func init() {
	log.SetFlags(0)
	log.SetPrefix(os.Args[0] + ": ")
}

var j = flag.Bool("j", false, "output all results in JSON")

func Exec() {
	flag.Parse()

	if len(flag.Args()) != 1 {
		log.Fatal("supply an IP address")
	}

	ipaddr := net.ParseIP(flag.Args()[0])
	if ipaddr == nil {
		log.Fatalf("wrong IP address: %s\n", flag.Args()[0])
	}

	results, errors := run(checker.DefaultCheckers, ipaddr)
	for _, e := range errors {
		log.Print(e)
	}
	results.SortByName()
	if *j {
		results.printJSON()
	} else {
		results.printInfo()
		results.printProbabilityMalicious()
	}
}

func run(checks []check.Check, ipaddr net.IP) (results, []error) {
	var results results
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

type results []check.Result

func (rs results) printJSON() {
	if len(rs) == 0 {
		return
	}
	enc := json.NewEncoder(os.Stdout)
	if err := enc.Encode(rs); err != nil {
		log.Fatal(err)
	}
}

func (rs results) SortByName() {
	sort.Slice(rs, func(i, j int) bool {
		return rs[i].Name < rs[j].Name
	})
}

// printInfo prints results from Info and InfoSec checkers.
func (rs results) printInfo() {
	for _, r := range rs {
		if r.Type == "Info" || r.Type == "InfoSec" {
			fmt.Printf("%-15s %s\n", r.Name, r.Info.String())
		}
	}
}

// printProbabilityMalicious prints the probability the IP address is malicious.
func (rs results) printProbabilityMalicious() {
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

func (rs results) probabilityMalicious() float64 {
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
