// Checkip is a command-line tool that checks an IP address.
package main

import (
	"flag"
	"log"
	"net"
	"os"

	"github.com/jreisinger/checkip/check"
	"github.com/jreisinger/checkip/cli"
)

func init() {
	log.SetFlags(0)
	log.SetPrefix(os.Args[0] + ": ")
}

var j = flag.Bool("j", false, "output all results in JSON")

func main() {
	flag.Parse()

	if len(flag.Args()) != 1 {
		log.Fatal("supply an IP address")
	}

	ipaddr := net.ParseIP(flag.Args()[0])
	if ipaddr == nil {
		log.Fatalf("wrong IP address: %s\n", flag.Args()[0])
	}

	results, errors := cli.Run(check.Default, ipaddr)
	for _, e := range errors {
		log.Print(e)
	}
	results.SortByName()
	if *j {
		results.PrintJSON()
	} else {
		results.PrintSummary()
		results.PrintMalicious()
	}
}
