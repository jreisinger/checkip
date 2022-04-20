// Checkip is a command-line tool that checks an IP address.
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

var j = flag.Bool("j", false, "output all results in JSON")
var c = flag.Int("c", 5, "number of concurrent checks")

func main() {
	flag.Parse()
	ipaddrs := parseArgs(flag.Args())

	// tokens is a counting semaphore used to
	// enforce a limit on concurrent checks.
	var tokens = make(chan struct{}, *c)

	var wg sync.WaitGroup
	for _, ipaddr := range ipaddrs {
		wg.Add(1)
		go func(ipaddr net.IP) {
			defer wg.Done()
			tokens <- struct{}{} // acquire a token

			results, errors := cli.Run(check.Default, ipaddr)
			for _, e := range errors {
				log.Print(e)
			}
			if *j {
				results.PrintJSON(ipaddr)
			} else {
				if len(ipaddrs) > 1 {
					fmt.Printf("--- %s ---\n", ipaddr.String())
				}
				results.SortByName()
				results.PrintSummary()
				results.PrintMalicious()
			}

			<-tokens // release the token
		}(ipaddr)
	}
	wg.Wait()
}

func parseArgs(args []string) []net.IP {
	var ipaddrs []net.IP

	for _, arg := range args {
		ipaddrs = append(ipaddrs, parseIP(arg))
	}

	// Get IP addresses from stdin.
	if len(ipaddrs) == 0 {
		input := bufio.NewScanner(os.Stdin)
		for input.Scan() {
			ipaddrs = append(ipaddrs, parseIP(input.Text()))
		}
		if err := input.Err(); err != nil {
			log.Print(err)
		}
	}

	return ipaddrs
}

func parseIP(ip string) net.IP {
	ipaddr := net.ParseIP(ip)
	if ipaddr == nil {
		log.Printf("wrong IP address: %s", ip)
	}
	return ipaddr
}
