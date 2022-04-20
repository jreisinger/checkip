// Checkip is a command-line tool that checks an IP address.
package main

import (
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

	if len(flag.Args()) < 1 {
		log.Fatal("supply at least one IP address")
	}

	var ipaddrs []net.IP

	for _, arg := range flag.Args() {
		ipaddr := net.ParseIP(arg)
		if ipaddr == nil {
			log.Printf("wrong IP address: %s", arg)
			continue
		}
		ipaddrs = append(ipaddrs, ipaddr)
	}

	// tokens is a counting semaphore used to
	// enforce a limit on concurrent requests.
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
