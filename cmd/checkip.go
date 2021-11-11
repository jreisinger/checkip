package cmd

import (
	"log"
	"net"
	"os"

	"github.com/jreisinger/checkip/pkg/check"
	"github.com/jreisinger/checkip/pkg/checker"
)

func Exec() {
	log.SetFlags(0)
	log.SetPrefix(os.Args[0] + ": ")

	if len(os.Args[1:]) != 1 {
		log.Fatal("supply an IP address")
	}

	ipaddr := net.ParseIP(os.Args[1])
	if ipaddr == nil {
		log.Fatalf("wrong IP address: %s\n", os.Args[1])
	}

	results := check.Run(checker.DefaultCheckers, ipaddr)
	results.SortByName()
	results.Print()
}
