package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/jreisinger/checkip/geodb"
)

func main() {
	log.SetFlags(0) // no timestamp

	if len(os.Args) != 2 {
		log.Fatalf("usage: %v %s\n", os.Args[0], "IPADDR")
	}

	licenseKey := os.Getenv("GEOIP_LICENSE_KEY")
	if licenseKey == "" {
		log.Fatalf("environment variable GEOIP_LICENSE_KEY not defined")
	}

	geoDBFilepath := "/var/tmp/GeoLite2-City.mmdb"
	geoDBUrl := "https://download.maxmind.com/app/geoip_download?edition_id=GeoLite2-City&license_key=" + licenseKey + "&suffix=tar.gz"

	g := &geodb.GeoDB{Filepath: geoDBFilepath, URL: geoDBUrl}

	if err := g.Update(geoDBUrl); err != nil {
		log.Fatalf("can't update geo DB: %v\n", err)
	}

	if err := g.Load(); err != nil {
		log.Fatalf("can't load geo DB: %v\n", err)
	}

	defer g.Close()

	ip := net.ParseIP(os.Args[1])
	if ip == nil {
		log.Fatalf("invalid IP address: %v\n", os.Args[1])
	}

	record, err := g.DB.City(ip)
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
