package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/jreisinger/checkip/asn"
	"github.com/jreisinger/checkip/dns"
	"github.com/jreisinger/checkip/geodb"
)

var outputPrefix = map[string]string{
	"geo": "Geo",
	"asn": "ASN",
	"dns": "DNS",
}

var Version = "dev"

func main() {
	log.SetFlags(0) // no timestamp

	if len(os.Args) != 2 {
		log.Fatalf("usage: %v %s\n", os.Args[0], "<IPADDR|version>")
	}

	if os.Args[1] == "version" {
		fmt.Println(Version)
		os.Exit(0)
	}

	ip := net.ParseIP(os.Args[1])
	if ip == nil {
		log.Fatalf("invalid IP address: %v\n", os.Args[1])
	}

	d := dns.New()
	if err := d.ForIP(ip); err != nil {
		fmt.Printf("%s: %v\n", outputPrefix["dns"], err)
	} else {
		fmt.Printf("%s: %v\n", outputPrefix["dns"], strings.Join(d.Names, ", "))
	}

	a := asn.New()
	if err := a.ForIP(ip); err != nil {
		fmt.Printf("%s: %v\n", outputPrefix["asn"], err)
	} else {
		fmt.Printf("%s: %d, %s - %s, %s, %s\n", outputPrefix["asn"], a.Number, a.FirstIP, a.LastIP, a.Description, a.CountryCode)
	}

	g := geodb.New()
	if err := g.ForIP(ip); err != nil {
		fmt.Printf("%s: %v\n", outputPrefix["geo"], err)
	} else {
		fmt.Printf("%s: %v\n", outputPrefix["geo"], strings.Join(g.Location, ", "))
	}

}
