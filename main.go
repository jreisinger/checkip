package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/jreisinger/checkip/asn"
	"github.com/jreisinger/checkip/geodb"
)

var outputPrefix = map[string]string{
	"geo": "Geo (maxmind.com)",
	"asn": "ASN (iptoasn.com)",
}

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
	if err := g.ForIP(ip); err != nil {
		fmt.Printf("%s: %v\n", outputPrefix["geo"], err)
	} else {
		fmt.Printf("%s: %v\n", outputPrefix["geo"], strings.Join(g.Location, ", "))
	}

	a, err := asn.ForIP(ip)
	if err != nil {
		fmt.Printf("%s: %v\n", outputPrefix["asn"], err)
	} else {
		fmt.Printf("%s: %d, %s - %s, %s, %s\n", outputPrefix["asn"], a.Number, a.FirsIP, a.LastIP, a.Description, a.CountryCode)
	}
}
