// Checkip is a command-line tool that provides information on IP addresses.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/jreisinger/checkip/check"
	"github.com/jreisinger/checkip/cli"
)

func init() {
	log.SetFlags(0)
	log.SetPrefix(os.Args[0] + ": ")
	// add extended definitions
	for _, d := range check.ExtDefinitions {
		check.Definitions = append(check.Definitions, d)
	}
}

var j = flag.Bool("j", false, "detailed output in JSON")
var noCache = flag.Bool("no-cache", false, "disable in-memory and persistent result cache")
var noActive = flag.Bool("no-active", false, "disable active checks that contact the target IP directly")
var m = flag.Bool("m", false, "MISP event output in JSON")
var p = flag.Int("p", 5, "check `n` IP addresses in parallel")
var t = flag.String("t", "", "list of checks")
var a = flag.String("a", "", "append to list of checks")
var d = flag.Bool("d", false, "debug")

// IpAndResults extend Results with Rating
type IpAndResults struct {
	IP      net.IP
	Results cli.Checks
	Rating  cli.Rating
}

func validateParallelism(parallelism int) error {
	if parallelism < 1 {
		return fmt.Errorf("invalid -p value %d: must be > 0", parallelism)
	}
	return nil
}

func selectedDefinitions(list []string, disableActive bool) []check.Definition {
	if !disableActive {
		return check.InList(list, check.Definitions)
	}
	return check.InList(list, check.WithoutActive(check.Definitions))
}

func getListOfChecks(cmd string, cmdappend string) ([]string, error) {
	txtlist, _ := check.GetConfigValue("CHECKS")
	if txtlist == "" {
		fmt.Fprintf(os.Stderr, "You can set a list of checks in CHECKS entry from config file\n")
		for _, d := range check.Definitions {
			txtlist = txtlist + ", " + d.Name
		}
	}

	if cmd != "" {
		txtlist = cmd
	}

	if cmdappend != "" {
		txtlist = txtlist + ", " + cmdappend
	}

	txtlist = strings.Replace(txtlist, " ", "", -1)
	split := strings.Split(txtlist, ",")

	rxAlpha := regexp.MustCompile("^[ a-zA-Z./-]+$")
	var list []string
	for _, s := range split {
		if rxAlpha.MatchString(s) == false {
			continue
		}
		exist := false
		for _, d := range check.Definitions {
			name := strings.Replace(d.Name, " ", "", -1)
			if s == name {
				exist = true
			}
		}
		if exist == true {
			list = append(list, s)
		} else {
			return list, fmt.Errorf("invalid check «%s»", s)
		}
	}

	return list, nil
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, " %s [-flag] IP [IP liste]\n", os.Args[0])

		flag.PrintDefaults()

		fmt.Fprintf(os.Stderr, "\n  Available Checks :\n  ")
		for _, d := range check.Definitions {
			fmt.Fprintf(os.Stderr, "%s, ", d.Name)
		}
		fmt.Fprintln(os.Stderr, "")
	}

	flag.Parse()
	if err := validateParallelism(*p); err != nil {
		log.Fatal(err)
	}

	check.Debug = *d

	list, err := getListOfChecks(*t, *a)
	if err != nil {
		log.Fatal(err)
	}

	if len(list) > 0 {
		if *j == false && *m == false {
			fmt.Println("Checks: " + strings.Join(list, ", "))
		}
	}

	runner := cli.NewRunnerWithOptions(selectedDefinitions(list, *noActive), cli.RunnerOptions{
		DisableCache: *noCache,
	})

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
				r, errors := runner.Run(ipaddr)
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

	// Experimental MISP event output
	currentTime := time.Now()
	date := fmt.Sprintf("%s", currentTime.Format("2006-01-02"))
	if *m {
		fmt.Printf("{\"Event\":{\"date\":\"%s\"", date)
		// distribution: 5 inherit from event, 2 connected community, 0 community
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
