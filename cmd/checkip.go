// Checkip checks an IP address using all checkers from checkip package.
package main

import (
	"flag"
	"fmt"
	"net"
	"os"

	"github.com/jreisinger/checkip"
)

var i = flag.Bool("i", false, "run only information checkers")

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

	if *i {
		checkers = infoCheckers
	} else {
		for k, v := range infoCheckers {
			checkers[k] = v
		}
	}
	checkip.RunAndPrint(checkers, ipaddr, "%-25s %s")
}
