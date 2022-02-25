package main

import (
	"flag"
	"log"
	"net"
	"os"

	"github.com/jreisinger/checkip/checks"
	"github.com/jreisinger/checkip/cli"
)

func init() {
	log.SetFlags(0)
	log.SetPrefix(os.Args[0] + ": ")
}

var j = flag.Bool("j", false, "output all results in JSON")
var a = flag.Bool("a", false, "run also active checks that interact with ipaddr")

func main() {
	flag.Parse()

	if len(flag.Args()) != 1 {
		log.Fatal("supply an IP address")
	}

	ipaddr := net.ParseIP(flag.Args()[0])
	if ipaddr == nil {
		log.Fatalf("wrong IP address: %s\n", flag.Args()[0])
	}

	var chks = checks.Passive
	if *a {
		chks = append(chks, checks.Active...)
	}

	results, errors := cli.Run(chks, ipaddr)
	for _, e := range errors {
		log.Print(e)
	}
	results.SortByName()
	if *j {
		results.PrintJSON()
	} else {
		results.PrintInfo()
		results.PrintProbabilityMalicious()
	}
}
