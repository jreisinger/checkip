package main

import (
	"flag"
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
	"github.com/jreisinger/checkip/virustotal"
)

var checkOutputPrefix = map[string]string{
	"asn":         "ASN         ",
	"dns":         "DNS         ",
	"geo":         "GEO         ",
	"abuseipdb":   "AbuseIPDB   ",
	"threatcrowd": "ThreatCrowd ",
	"virustotal":  "VirusTotal  ",
}

// Version is the default version od checkip.
var Version = "dev"

func main() {
	log.SetFlags(0) // no timestamp

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "%s [flags] <ipaddr>\n", os.Args[0])
		flag.PrintDefaults()
	}

	version := flag.Bool("version", false, "version")

	flag.Parse()

	if *version {
		fmt.Println(Version)
		os.Exit(0)
	}

	if len(flag.Args()) != 1 {
		log.Fatalf("missing IP address to check")
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

	go func(ch chan string) {
		v := virustotal.New()
		if err := v.ForIP(ip); err != nil {
			ch <- fmt.Sprintf("%s %v\n", checkOutputPrefix["virustotal"], err)
		} else {
			ch <- fmt.Sprintf("%s scannners results: %d malicious, %d suspicious, %d harmless\n", checkOutputPrefix["virustotal"], v.Data.Attributes.LastAnalysisStats.Malicious, v.Data.Attributes.LastAnalysisStats.Suspicious, v.Data.Attributes.LastAnalysisStats.Harmless)
		}
	}(ch)

	for i := 0; i < len(checkOutputPrefix); i++ {
		fmt.Printf("%s", <-ch)
	}
}
