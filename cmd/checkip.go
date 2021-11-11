package cmd

import (
	"encoding/json"
	"flag"
	"log"
	"net"
	"os"

	"github.com/jreisinger/checkip/check"
	"github.com/jreisinger/checkip/checker"
)

// var v = flag.Bool("v", false, "be verbose")
var j = flag.Bool("j", false, "output all data in JSON")

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
		enc := json.NewEncoder(os.Stdout)
		if err := enc.Encode(results); err != nil {
			log.Fatal(err)
		}
	} else {
		results.Print()
	}
}
