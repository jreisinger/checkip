package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/oschwald/geoip2-golang"
)

const geoDB = "GeoLite2-City.mmdb"

func main() {
	log.SetFlags(0) // no timestamp

	if len(os.Args) != 2 {
		log.Fatalf("usage: %v %s\n", os.Args[0], "IPADDR")
	}

	db, err := geoip2.Open(geoDB)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	ip := net.ParseIP(os.Args[1])
	if ip == nil {
		log.Fatalf("invalid IP address: %v\n", os.Args[1])
	}

	record, err := db.City(ip)
	if err != nil {
		log.Fatal(err)
	}

	city := record.City.Names["en"]
	country := record.Country.Names["en"]
	isoCode := record.Country.IsoCode

	if city != "" || country != "" || isoCode != "" {
		fmt.Printf("%v, %v (%v)\n", city, country, isoCode)
	}
}
