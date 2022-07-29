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
var c = flag.Int("c", 5, "IP addresses being checked concurrently")

type IpAndResults struct {
	IP      net.IP
	Results cli.Results
}

func main() {
	flag.Parse()

	checks := check.Default
	if *a {
		checks = check.All
	}

	ipaddrsCh := make(chan net.IP)
	resultsCh := make(chan IpAndResults)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		getIpAddrs(flag.Args(), ipaddrsCh)
		wg.Done()
	}()

	for i := 0; i < *c; i++ {
		wg.Add(1)
		go func() {
			for ipaddr := range ipaddrsCh {
				r, errors := cli.Run(checks, ipaddr)
				for _, e := range errors {
					log.Print(e)
				}
				resultsCh <- IpAndResults{ipaddr, r}
			}
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(resultsCh)
	}()

	for c := range resultsCh {
		if *j {
			c.Results.PrintJSON(c.IP)
		} else {
			fmt.Printf("--- %s ---\n", c.IP.String())
			c.Results.SortByName()
			c.Results.PrintSummary()
			c.Results.PrintMalicious()
		}
	}
}

// getIpAddrs parses IP addresses supplied as command line arguments or as
// STDIN. It sends the received IP addresses down the ipaddrsCh.
func getIpAddrs(args []string, ipaddrsCh chan net.IP) {
	defer close(ipaddrsCh)

	if len(args) == 0 { // get IP addresses from stdin.
		input := bufio.NewScanner(os.Stdin)
		for input.Scan() {
			ipaddr := net.ParseIP(input.Text())
			if ipaddr == nil {
				log.Printf("wrong IP address: %s", input.Text())
				continue
			}
			ipaddrsCh <- ipaddr
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
			ipaddrsCh <- ipaddr
		}
	}
}
