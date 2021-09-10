// Checkip checks an IP address using all checkers from checkip package.
package main

import (
	"flag"
	"fmt"
	"net"
	"os"

	"github.com/jreisinger/checkip"
)

var s = flag.Bool("s", false, "only print how many checkers consider the ipaddr suspicious")

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

	// checkers can tell you wether the IP address is suspicious.
	checkers := map[string]checkip.Checker{
		"abuseipdb.com":             &checkip.AbuseIPDB{},
		"cinsscore.com":             &checkip.CINSArmy{},
		"emergingthreats.net":       &checkip.ET{},
		"otx.alienvault.com":        &checkip.OTX{},
		"github.com/stamparm/ipsum": &checkip.IPsum{},
		"shodan.io":                 &checkip.Shodan{},
		"threatcrowd.org":           &checkip.ThreatCrowd{},
		"virustotal.com":            &checkip.VirusTotal{},
	}

	// infoCheckers just give you information about an IP address. They
	// always return ok == true.
	infoCheckers := map[string]checkip.Checker{
		"iptoasn.com":          &checkip.AS{},
		"net.LookupAddr":       &checkip.DNS{},
		"maxmind.com GeoLite2": &checkip.Geo{},
		"net.IP":               &checkip.IP{},
	}

	if *s {
		var checkerSlice []checkip.Checker
		for _, v := range checkers {
			checkerSlice = append(checkerSlice, v)
		}
		n := checkip.Run(checkerSlice, ipaddr)
		perc := float64(n) / float64(len(checkers)) * 100.0
		fmt.Printf("%02.0f%% (%d/%d)\n", perc, n, len(checkers))
	} else {
		for k, v := range infoCheckers {
			checkers[k] = v
		}
		checkip.RunAndPrint(checkers, ipaddr, "%-25s %s")
	}

}
