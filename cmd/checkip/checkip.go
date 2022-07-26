// Checkip is a command-line tool that provides information on IP addresses.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/jreisinger/checkip"
	"github.com/jreisinger/checkip/check"
	"github.com/jreisinger/checkip/cli"
)

func init() {
	log.SetFlags(0)
	log.SetPrefix(os.Args[0] + ": ")
}

var a = flag.Bool("a", false, "run all available checks")
var j = flag.Bool("j", false, "output all results in JSON")
var c = flag.Int("c", 5, "number of concurrent checkers")

type IpAndResults struct {
	IP      net.IP
	Results cli.Results
}

func main() {
	flag.Parse()
	ipaddrs := parseArgs(flag.Args())

	checks := check.Default
	if *a {
		checks = check.All
	}

	ipaddrsCh := make(chan net.IP)
	ipAndResultsCh := make(chan IpAndResults)

	for i := 0; i < *c; i++ {
		go checker(checks, ipaddrsCh, ipAndResultsCh)
	}

	go func() {
		for _, ipaddr := range ipaddrs {
			ipaddrsCh <- ipaddr
		}
	}()

	for range ipaddrs {
		c := <-ipAndResultsCh
		if *j {
			c.Results.PrintJSON(c.IP)
		} else {
			if len(ipaddrs) > 1 {
				fmt.Printf("--- %s ---\n", c.IP.String())
			}
			c.Results.SortByName()
			c.Results.PrintSummary()
			c.Results.PrintMalicious()
		}
	}
}

// checker runs checks against IP addresses coming from ipaddrs and sends back
// the IP address and checks results.
func checker(checks []checkip.Check, ipaddrs chan net.IP, ipAndResults chan IpAndResults) {
	for ipaddr := range ipaddrs {
		r, errors := cli.Run(checks, ipaddr)
		for _, e := range errors {
			log.Print(e)
		}
		ipAndResults <- IpAndResults{ipaddr, r}
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
