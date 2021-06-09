package main

import (
	"fmt"
	"net"
	"os"

	"github.com/jreisinger/checkip"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s <ipaddr>\n", os.Args[0])
		os.Exit(0)
	}

	ipaddr := net.ParseIP(os.Args[1])
	if ipaddr == nil {
		fmt.Fprintf(os.Stderr, "%s: wrong IP address: %s\n", os.Args[0], os.Args[1])
		os.Exit(1)
	}

	checkers := checkip.Checkers{
		"abuseipdb.com":             &checkip.AbuseIPDB{},
		"iptoasn.com":               &checkip.AS{},
		"net.LookupAddr":            &checkip.DNS{},
		"maxmind.com geolocation":   &checkip.Geo{},
		"github.com/stamparm/ipsum": &checkip.IPsum{},
		"otx.alienvault.com":        &checkip.OTX{},
		"shodan.io":                 &checkip.Shodan{},
		"threatcrowd.org":           &checkip.ThreatCrowd{},
		"virustotal.com":            &checkip.VirusTotal{},
	}
	checkers.Run(ipaddr)
}
