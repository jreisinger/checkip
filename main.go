package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/jreisinger/checkip/abuseipdb"
	"github.com/jreisinger/checkip/asn"
	"github.com/jreisinger/checkip/dns"
	"github.com/jreisinger/checkip/geo"
	"github.com/jreisinger/checkip/threatcrowd"
)

var checkOutputPrefix = map[string]string{
	"asn":         "ASN         ",
	"dns":         "DNS         ",
	"geo":         "GEO         ",
	"abuseipdb":   "AbuseIPDB   ",
	"threatcrowd": "ThreatCrowd ",
}

// Version is the default version od checkip.
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

	ch := make(chan string)

	go func(ch chan string) {
		d := dns.New()
		if err := d.ForIP(ip); err != nil {
			ch <- fmt.Sprintf("%s %v\n", checkOutputPrefix["dns"], err)
		} else {
			ch <- fmt.Sprintf("%s %v\n", checkOutputPrefix["dns"], strings.Join(d.Names, " | "))
		}
	}(ch)

	go func(ch chan string) {
		a := asn.New()
		if err := a.ForIP(ip); err != nil {
			ch <- fmt.Sprintf("%s %v\n", checkOutputPrefix["asn"], err)
		} else {
			ch <- fmt.Sprintf("%s %d | %s - %s | %s | %s\n", checkOutputPrefix["asn"], a.Number, a.FirstIP, a.LastIP, a.Description, a.CountryCode)
		}
	}(ch)

	go func(ch chan string) {
		g := geo.New()
		if err := g.ForIP(ip); err != nil {
			ch <- fmt.Sprintf("%s %v\n", checkOutputPrefix["geo"], err)
		} else {
			ch <- fmt.Sprintf("%s %v\n", checkOutputPrefix["geo"], strings.Join(g.Location, " | "))
		}
	}(ch)

	go func(ch chan string) {
		a := abuseipdb.New()
		if err := a.ForIP(ip); err != nil {
			ch <- fmt.Sprintf("%s %v\n", checkOutputPrefix["abuseipdb"], err)
		} else {
			abuseConfidenceScore := a.Data.AbuseConfidenceScore
			domain := a.Data.Domain
			ch <- fmt.Sprintf("%s malicious with %d%% confidence | %v\n", checkOutputPrefix["abuseipdb"], abuseConfidenceScore, domain)
		}
	}(ch)

	go func(ch chan string) {
		// https://github.com/AlienVault-OTX/ApiV2#votes
		votesMeaning := map[int]string{
			-1: "most users have voted this malicious",
			0:  "equal number of users have voted this malicious and not malicious",
			1:  "most users have voted this not malicious",
		}

		t := threatcrowd.New()
		if err := t.ForIP(ip); err != nil {
			ch <- fmt.Sprintf("%s %v\n", checkOutputPrefix["threatcrowd"], err)
		} else {
			ch <- fmt.Sprintf("%s %v\n", checkOutputPrefix["threatcrowd"], votesMeaning[t.Votes])
		}
	}(ch)

	for i := 0; i < len(checkOutputPrefix); i++ {
		fmt.Printf("%s", <-ch)
	}
}
