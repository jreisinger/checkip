package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/jreisinger/geoip/geodb"
)

func main() {
	log.SetFlags(0) // no timestamp

	if len(os.Args) != 2 {
		log.Fatalf("usage: %v %s\n", os.Args[0], "IPADDR")
	}

	ip := net.ParseIP(os.Args[1])
	if ip == nil {
		log.Fatalf("invalid IP address: %v\n", os.Args[1])
	}

	g := geodb.New()

	if err := g.Update(); err != nil {
		log.Fatalf("can't update geo DB: %v\n", err)
	}

	if err := g.Open(); err != nil {
		log.Fatalf("can't load geo DB: %v\n", err)
	}
	defer g.Close()

	if err := g.GetLocation(ip); err != nil {
		log.Fatalf("can't get location: %v\n", err)
	}
	fmt.Printf("%v\n", strings.Join(g.Location, ", "))
}
