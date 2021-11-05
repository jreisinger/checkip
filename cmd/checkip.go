// Checkip quickly finds information about an IP address from a CLI.
package main

import (
	"flag"
	"log"
	"net"
	"os"

	"github.com/jreisinger/checkip"
)

func init() {
	log.SetFlags(0)
	log.SetPrefix(os.Args[0] + ": ")
}

var j = flag.Bool("j", false, "output JSON")

func main() {
	flag.Parse()

	if len(flag.Args()) != 1 {
		log.Fatal("missing IP address")
	}

	ipaddr := net.ParseIP(flag.Arg(0))
	if ipaddr == nil {
		log.Fatalf("wrong IP address: %s\n", flag.Arg(0))
	}

	checkers := []checkip.Checker{
		&checkip.AS{},
		&checkip.AbuseIPDB{},
		&checkip.Blocklist{},
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

	results := checkip.Run(checkers, ipaddr)
	if *j {
		checkip.PrintJSON(results)
	} else {
		checkip.Print(results)
	}
}
