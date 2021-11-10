package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/jreisinger/checkip/api"
	"github.com/jreisinger/checkip/cmd"
	checkip "github.com/jreisinger/checkip/pkg"
)

var s = flag.Bool("s", false, "serve JSON API")

func main() {
	flag.Parse()

	checkers := []checkip.Checker{
		&checkip.AS{},
		&checkip.AbuseIPDB{},
		&checkip.Blocklist{},
		&checkip.CINSArmy{},
		&checkip.DNS{},
		&checkip.Geo{},
		&checkip.IPsum{},
		&checkip.OTX{},
		&checkip.Shodan{},
		&checkip.ThreatCrowd{},
		&checkip.VirusTotal{},
	}

	if *s {
		c := api.Checkers(checkers)
		http.HandleFunc("/api/v1/", c.Handler)
		log.Fatal(http.ListenAndServe(":8000", nil))
	} else {

		log.SetFlags(0)
		log.SetPrefix(os.Args[0] + ": ")

		if len(flag.Args()) != 1 {
			log.Fatal("missing IP address")
		}

		ipaddr := net.ParseIP(flag.Arg(0))
		if ipaddr == nil {
			log.Fatalf("wrong IP address: %s\n", flag.Arg(0))
		}

		results := checkip.Run(checkers, ipaddr)

		cmd.Print(results)
	}
}
