package main

import (
	"flag"
	"fmt"
	"net"
	"os"

	"github.com/jreisinger/checkip"
)

var c = flag.Bool("c", false, "use only checkers that can tell whether IP address is suspicious")

func main() {
	flag.Parse()

	if len(flag.Args()) != 1 {
		fmt.Printf("Usage: %s <ipaddr>\n", os.Args[0])
		os.Exit(0)
	}

	ipaddr := net.ParseIP(flag.Arg(0))
	if ipaddr == nil {
		fmt.Fprintf(os.Stderr, "%s: wrong IP address: %s\n", os.Args[0], flag.Arg(0))
		os.Exit(1)
	}

	// checkers can tell you wether the IP address is suspicious.
	checkers := checkip.Checkers{
		"abuseipdb.com":             &checkip.AbuseIPDB{},
		"otx.alienvault.com":        &checkip.OTX{},
		"github.com/stamparm/ipsum": &checkip.IPsum{},
		"shodan.io":                 &checkip.Shodan{},
		"threatcrowd.org":           &checkip.ThreatCrowd{},
		"virustotal.com":            &checkip.VirusTotal{},
	}

	// infoCheckers just give you information about an IP address.
	infoCheckers := checkip.Checkers{
		"iptoasn.com":          &checkip.AS{},
		"net.LookupAddr":       &checkip.DNS{},
		"maxmind.com GeoLite2": &checkip.Geo{},
	}

	if !*c {
		for k, v := range infoCheckers {
			checkers[k] = v
		}
	}

	checkers.CheckAndPrint(ipaddr)
}
