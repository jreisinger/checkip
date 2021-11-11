package cmd

import (
	"flag"
	"github.com/jreisinger/checkip/pkg/check"
	"github.com/jreisinger/checkip/pkg/checker"
	"log"
	"net"
	"os"
)

func init() {
	log.SetFlags(0)
	log.SetPrefix(os.Args[0] + ": ")
}

var j = flag.Bool("j", false, "print all data in JSON")

func CheckIP() {
	flag.Parse()

	if len(flag.Args()) != 1 {
		log.Fatal("missing IP address")
	}

	ipaddr := net.ParseIP(flag.Arg(0))
	if ipaddr == nil {
		log.Fatalf("wrong IP address: %s\n", flag.Arg(0))
	}

	results := check.Run(checker.DefaultCheckers, ipaddr)
	results.SortByName()
	results.Print()
}
