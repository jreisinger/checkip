// Demo demonstrates how you can use checkip as a library.
package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/jreisinger/checkip"
	"github.com/jreisinger/checkip/check"
)

func main() {
	ipaddr := net.ParseIP(os.Args[1])
	if ipaddr == nil {
		log.Fatalf("wrong IP address: %s\n", flag.Args()[0])
	}

	chks := []checkip.Check{
		check.AbuseIPDB,
	}

	for _, c := range chks {
		r, err := c(ipaddr)
		if err != nil {
			log.Print(err)
		}
		fmt.Printf("%-10s %s\n", r.Name, r.Info.Summary())
	}
}
