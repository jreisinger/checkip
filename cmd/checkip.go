package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"sync"

	"github.com/jreisinger/checkip"
	"github.com/logrusorgru/aurora"
)

func init() {
	log.SetFlags(0)
	log.SetPrefix(os.Args[0] + ": ")
}

func main() {
	if len(os.Args[1:]) != 1 {
		fmt.Printf("Usage: %s <ipaddr>\n", os.Args[0])
		os.Exit(1)
	}

	ipaddr := net.ParseIP(os.Args[1])
	if ipaddr == nil {
		fmt.Fprintf(os.Stderr, "%s: wrong IP address: %s\n", os.Args[0], flag.Arg(0))
		os.Exit(1)
	}

	checkers := []checkip.Checker{
		&checkip.AS{},
		&checkip.AbuseIPDB{},
		&checkip.CINSArmy{},
		&checkip.DNS{},
		&checkip.ET{},
		&checkip.Geo{},
		&checkip.IP{},
		&checkip.IPsum{},
		&checkip.OTX{},
		&checkip.Shodan{},
		&checkip.ThreatCrowd{},
		&checkip.VirusTotal{},
	}

	var wg sync.WaitGroup
	for _, c := range checkers {
		wg.Add(1)
		go func(c checkip.Checker) {
			defer wg.Done()
			if err := c.Check(ipaddr); err != nil {
				log.Print(err)
			}
		}(c)
	}
	wg.Wait()

	var total, malicious int
	for _, c := range checkers {
		switch ip := c.(type) {
		case checkip.InfoChecker:
			fmt.Println(ip.Info())
		case checkip.SecChecker:
			total++
			if !ip.IsOK() {
				malicious++
			}
		}

	}
	perc := float64(malicious) / float64(total) * 100
	var msg string
	switch {
	case perc < 15:
		msg = fmt.Sprint(aurora.Green("Malicious"))
	case perc < 50:
		msg = fmt.Sprint(aurora.Yellow("Malicious"))
	default:
		msg = fmt.Sprint(aurora.Red("Malicious"))
	}
	fmt.Printf("%s\t%.0f%% (%d out of %d checkers)\n", msg, perc, malicious, total)
}
