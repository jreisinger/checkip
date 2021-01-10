package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/jreisinger/checkip/check"
)

// Version is the default version od checkip.
var Version = "dev"

func main() {
	log.SetFlags(0) // no timestamp in error messages
	handleFlags()

	ipaddr := net.ParseIP(os.Args[1])
	if ipaddr == nil {
		log.Fatalf("invalid IP address: %v\n", os.Args[1])
	}

	ch := make(chan string)
	checkers := []check.Checker{
		&check.DNS{},
		&check.ThreatCrowd{},
		&check.VirusTotal{},
		&check.Geo{},
		&check.AS{},
		&check.AbuseIPDB{},
	}
	for _, chk := range checkers {
		go check.Run(chk, ipaddr, ch)
	}
	for range checkers {
		fmt.Print(<-ch)
	}
}
