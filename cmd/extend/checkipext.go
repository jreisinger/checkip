// Checkip is a command-line tool that provides information on IP addresses.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/jreisinger/checkip/check"
	"github.com/jreisinger/checkip/cli"
)

var funcRegistry = make(map[string]interface{})

func init() {
	log.SetFlags(0)
	log.SetPrefix(os.Args[0] + ": ")
	for _, v := range check.All {
		name := runtime.FuncForPC(reflect.ValueOf(v).Pointer()).Name()
		funcRegistry[name] = v
	}
}

var j = flag.Bool("j", false, "detailed output in JSON")
var m = flag.Bool("m", false, "MISP event output in JSON")
var p = flag.Int("p", 5, "check `n` IP addresses in parallel")
var t = flag.String("t", "", "list of checks")
var d = flag.Bool("d", false, "debug")

type IpAndResults struct {
	IP      net.IP
	Results cli.Checks
	Rating  cli.Rating
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, " %s [-flag] IP [IP liste]\n", os.Args[0])

		flag.PrintDefaults()

		fmt.Fprintf(os.Stderr, "\n  Available Checks :\n  ")
		for v := range funcRegistry {
			v = strings.Replace(v, "github.com/jreisinger/checkip/check.", "", -1)
			fmt.Fprintf(os.Stderr, "%s, ", v)
		}
		fmt.Fprintln(os.Stderr, "")
	}

	flag.Parse()

	check.Debug = *d

	tcheck := ""
	if *t == "" {
		tcheck, _ = check.GetConfigValue("CHECKS")
	} else {
		tcheck = *t
	}
	tcheck = strings.Replace(tcheck, " ", "", -1)
	split := strings.Split(tcheck, ",")
	for _, s := range split {
		fname := "github.com/jreisinger/checkip/check." + s
		if funcRegistry[fname] != nil {
			check.AddUse(funcRegistry[fname])
		}
	}

	checks := check.Funcs
	if len(check.Use) > 0 {
		if *j == false && *m == false {
			fmt.Println("Checks: " + tcheck)
		}
		checks = check.Use
	}

	ipaddrsCh := make(chan net.IP)
	resultsCh := make(chan IpAndResults)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		getIpAddrs(flag.Args(), ipaddrsCh)
		wg.Done()
	}()

	var ra cli.Rating
	for i := 0; i < *p; i++ {
		wg.Add(1)
		go func() {
			for ipaddr := range ipaddrsCh {
				r, errors := cli.Run(checks, ipaddr)
				for _, e := range errors {
					if *d == true {
						log.Print(e)
					}
				}
				ra = r.Cotation()
				resultsCh <- IpAndResults{ipaddr, r, ra}
			}
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(resultsCh)
	}()

    currentTime := time.Now()
    date := fmt.Sprintf("%s", currentTime.Format("2006-01-02"))
	if *m {
		fmt.Printf("{\"Event\":{\"date\":\"%s\"",date)
		// distribution: 5 inherit frome event, 2 connected community, 0 community
		fmt.Printf(`,"threat_level_id":"1","info":"testevent","published":false,"analysis":"0","distribution":"2","Attribute":[`)
	}

	var res []string
	for c := range resultsCh {
		ra = c.Rating
		rating := fmt.Sprintf("%s1 - %s", ra.Type, ra.TypeDesc)
		if *j {
			c.Results.PrintExtJSON(c.IP, ra)
		} else if *m {
			c.Results.PrintMISPJSON(c.IP, ra)
			fmt.Printf(",")
		} else {
			fmt.Printf("--- %s ---\n", c.IP.String())
			c.Results.SortByName()
			r := c.Results.ExtPrintSummary()
			if ra.Type != "" {
				res = append(res, fmt.Sprintf("%s [%s]", r, rating))
				fmt.Printf("%-15s %s\n", "Cotation", rating)
			}
			c.Results.PrintMalicious()
		}
	}

    if *m {
		fmt.Printf("{}]}}")
	}

	if len(res) > 0 {
		fmt.Printf("\nIOC: %s\n", strings.Join(res, ", "))
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
