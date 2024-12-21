// Checkip is a command-line tool that provides information on IP addresses.
package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"sync"

	"github.com/jreisinger/checkip/check"
	"github.com/jreisinger/checkip/cli"
)

func init() {
	log.SetFlags(0)
	log.SetPrefix(filepath.Base(os.Args[0]) + ": ")
}

var j = flag.Bool("j", false, "detailed output in JSON")
var p = flag.Int("p", 5, "check `n` IP addresses in parallel")

type Result struct {
	IP     net.IP
	Checks cli.Checks
}

func main() {
	flag.Parse()

	ipaddrs := make(chan net.IP)
	results := make(chan Result)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		cli.GetIpAddrs(flag.Args(), ipaddrs)
		wg.Done()
	}()

	for i := 0; i < *p; i++ {
		wg.Add(1)
		go func() {
			for ipaddr := range ipaddrs {
				checks, errors := cli.Run(check.Funcs, ipaddr)
				for _, e := range errors {
					log.Print(e)
				}
				results <- Result{ipaddr, checks}
			}
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		if *j {
			result.Checks.PrintJSON(result.IP)
		} else {
			fmt.Printf("--- %s ---\n", result.IP.String())
			result.Checks.SortByName()
			result.Checks.PrintSummary()
			result.Checks.PrintMalicious()
		}
	}
}
