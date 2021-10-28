// Checkip checks an IP address using all checkers from checkip package.
package main

import (
	"flag"
	"fmt"
	"net"
	"os"

	"github.com/jreisinger/checkip"
	"github.com/logrusorgru/aurora"
)

func main() {
	flag.Parse()

	if len(flag.Args()) != 1 {
		fmt.Printf("Usage: %s <ipaddr>\n", os.Args[0])
		os.Exit(1)
	}

	ipaddr := net.ParseIP(flag.Arg(0))
	if ipaddr == nil {
		fmt.Fprintf(os.Stderr, "%s: wrong IP address: %s\n", os.Args[0], flag.Arg(0))
		os.Exit(1)
	}

	// secCheckers can tell you wether the IP address is suspicious.
	secCheckers := []checkip.Checker{
		&checkip.AbuseIPDB{},
		&checkip.CINSArmy{},
		&checkip.ET{},
		&checkip.OTX{},
		&checkip.IPsum{},
		&checkip.ThreatCrowd{},
		&checkip.VirusTotal{},
	}

	// infoCheckers just give you information about an IP address. They
	// always return ok == true.
	infoCheckers := []checkip.Checker{
		&checkip.AS{},
		&checkip.DNS{},
		&checkip.Geo{},
		&checkip.IP{},
		&checkip.Shodan{},
	}

	n := checkip.Run(secCheckers, ipaddr)
	perc := float64(n) / float64(len(secCheckers)) * 100
	var msg string
	switch {
	case perc < 15:
		msg = fmt.Sprint(aurora.Green("Malicious"))
	case perc < 50:
		msg = fmt.Sprint(aurora.Yellow("Malicious"))
	default:
		msg = fmt.Sprint(aurora.Red("Malicious"))
	}
	fmt.Printf("%s\t%.0f%% (%d out of %d checkers)\n", msg, perc, n, len(secCheckers))

	checkip.Run(infoCheckers, ipaddr)
	for _, c := range infoCheckers {
		fmt.Println(c)
	}
}
