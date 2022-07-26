// Checkip is a command-line tool that provides information on IP addresses.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"sync"

	"github.com/jreisinger/checkip/check"
	"github.com/jreisinger/checkip/cli"
)

func init() {
	log.SetFlags(0)
	log.SetPrefix(os.Args[0] + ": ")
}

var a = flag.Bool("a", false, "run all available checks")
var j = flag.Bool("j", false, "output all results in JSON")
var c = flag.Int("c", 5, "number of concurrent checks")

func main() {
	flag.Parse()
	ipaddrs := parseArgs(flag.Args())

	checks := check.Default
	if *a {
		checks = check.All
	}

	// tokens is a counting semaphore used to
	// enforce a limit on concurrent checks.
	var tokens = make(chan struct{}, *c)

	resultsPerIP := make(map[string]cli.Results)

	var wg sync.WaitGroup
	for _, ipaddr := range ipaddrs {
		wg.Add(1)
		go func(ipaddr net.IP) {
			defer wg.Done()
			tokens <- struct{}{} // acquire a token

			r, errors := cli.Run(checks, ipaddr)
			resultsPerIP[ipaddr.String()] = r
			for _, e := range errors {
				log.Print(e)
			}

			<-tokens // release the token
		}(ipaddr)
	}
	wg.Wait()

	for ip, results := range resultsPerIP {
		if *j {
			results.PrintJSON(net.ParseIP(ip))
		} else {
			if len(ipaddrs) > 1 {
				fmt.Printf("--- %s ---\n", ip)
			}
			results.SortByName()
			results.PrintSummary()
			results.PrintMalicious()
		}
	}
}

func parseArgs(args []string) []net.IP {
	var ipaddrs []net.IP

	for _, arg := range args {
		ipaddr := net.ParseIP(arg)
		if ipaddr == nil {
			log.Printf("wrong IP address: %s", arg)
			continue
		}
		ipaddrs = append(ipaddrs, ipaddr)
	}

	// Get IP addresses from stdin.
	if len(ipaddrs) == 0 {
		input := bufio.NewScanner(os.Stdin)
		for input.Scan() {
			ipaddr := net.ParseIP(input.Text())
			if ipaddr == nil {
				log.Printf("wrong IP address: %s", input.Text())
				continue
			}
			ipaddrs = append(ipaddrs, ipaddr)
		}
		if err := input.Err(); err != nil {
			log.Print(err)
		}
	}

	return ipaddrs
}
