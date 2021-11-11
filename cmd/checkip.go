package cmd

import (
	"flag"
	"log"
	"net"
	"os"

	"github.com/jreisinger/checkip/check"
	"github.com/jreisinger/checkip/checker"
)

var j = flag.Bool("j", false, "output all results in JSON")

func Exec() {
	flag.Parse()

	log.SetFlags(0)
	log.SetPrefix(os.Args[0] + ": ")

	if len(flag.Args()) != 1 {
		log.Fatal("supply an IP address")
	}

	ipaddr := net.ParseIP(flag.Args()[0])
	if ipaddr == nil {
		log.Fatalf("wrong IP address: %s\n", flag.Args()[0])
	}

	results := check.Run(checker.DefaultCheckers, ipaddr)
	results.SortByName()
	if *j {
		results.PrintJSON()
	} else {
		results.PrintInfo()
		results.PrintProbabilityMalicious()
	}
}
