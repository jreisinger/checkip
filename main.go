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
	"github.com/jreisinger/checkip/virustotal"

	. "github.com/logrusorgru/aurora"
)

// Standard print format.
var format = map[string]string{
	"asn":         "ASN         ",
	"dns":         "DNS         ",
	"geo":         "GEO         ",
	"abuseipdb":   "AbuseIPDB   ",
	"threatcrowd": "ThreatCrowd ",
	"virustotal":  "VirusTotal  ",
}

// Error print format.
var formatErr = map[string]string{
	"asn":         fmt.Sprint(Gray(11, format["asn"])),
	"dns":         fmt.Sprint(Gray(11, format["dns"])),
	"geo":         fmt.Sprint(Gray(11, format["geo"])),
	"abuseipdb":   fmt.Sprint(Gray(11, format["abuseipdb"])),
	"threatcrowd": fmt.Sprint(Gray(11, format["threatcrowd"])),
	"virustotal":  fmt.Sprint(Gray(11, format["virustotal"])),
}

// Problem print format.
var formatProb = map[string]string{
	"asn":         fmt.Sprint(Magenta(format["asn"])),
	"dns":         fmt.Sprint(Magenta(format["dns"])),
	"geo":         fmt.Sprint(Magenta(format["geo"])),
	"abuseipdb":   fmt.Sprint(Magenta(format["abuseipdb"])),
	"threatcrowd": fmt.Sprint(Magenta(format["threatcrowd"])),
	"virustotal":  fmt.Sprint(Magenta(format["virustotal"])),
}

// Version is the default version od checkip.
var Version = "dev"

func main() {
	log.SetFlags(0) // no timestamp in error messages
	handleFlags()

	ip := net.ParseIP(os.Args[1])
	if ip == nil {
		log.Fatalf("invalid IP address: %v\n", os.Args[1])
	}

	ch := make(chan string)

	go func(ch chan string) {
		d := dns.New()
		if err := d.ForIP(ip); err != nil {
			ch <- fmt.Sprintf("%s %v\n", formatErr["dns"], err)
		} else {
			ch <- fmt.Sprintf("%s %v\n", format["dns"], strings.Join(d.Names, " | "))
		}
	}(ch)

	go func(ch chan string) {
		a := asn.New()
		if err := a.ForIP(ip); err != nil {
			ch <- fmt.Sprintf("%s %v\n", formatErr["asn"], err)
		} else {
			ch <- fmt.Sprintf("%s %d | %s - %s | %s | %s\n", format["asn"], a.Number, a.FirstIP, a.LastIP, a.Description, a.CountryCode)
		}
	}(ch)

	go func(ch chan string) {
		g := geo.New()
		if err := g.ForIP(ip); err != nil {
			ch <- fmt.Sprintf("%s %v\n", formatErr["geo"], err)
		} else {
			ch <- fmt.Sprintf("%s %v\n", format["geo"], strings.Join(g.Location, " | "))
		}
	}(ch)

	go func(ch chan string) {
		a := abuseipdb.New()
		if err := a.ForIP(ip); err != nil {
			ch <- fmt.Sprintf("%s %v\n", formatErr["abuseipdb"], err)
		} else {
			abuseConfidenceScore := a.Data.AbuseConfidenceScore
			domain := a.Data.Domain
			f := format["abuseipdb"]
			if abuseConfidenceScore > 0 {
				f = formatProb["abuseipdb"]
			}
			ch <- fmt.Sprintf("%s malicious with %d%% confidence | %v\n", f, abuseConfidenceScore, domain)
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
			ch <- fmt.Sprintf("%s %v\n", formatErr["threatcrowd"], err)
		} else {
			f := format["threatcrowd"]
			if t.Votes < 0 {
				f = formatProb["threatcrowd"]
			}
			ch <- fmt.Sprintf("%s %v\n", f, votesMeaning[t.Votes])
		}
	}(ch)

	go func(ch chan string) {
		v := virustotal.New()
		if err := v.ForIP(ip); err != nil {
			ch <- fmt.Sprintf("%s %v\n", formatErr["virustotal"], err)
		} else {
			f := format["virustotal"]
			if v.Data.Attributes.LastAnalysisStats.Malicious > 0 || v.Data.Attributes.LastAnalysisStats.Suspicious > 0 {
				f = formatProb["virustotal"]
			}
			ch <- fmt.Sprintf("%s scannners results: %d malicious, %d suspicious, %d harmless\n", f, v.Data.Attributes.LastAnalysisStats.Malicious, v.Data.Attributes.LastAnalysisStats.Suspicious, v.Data.Attributes.LastAnalysisStats.Harmless)
		}
	}(ch)

	for i := 0; i < len(format); i++ {
		fmt.Printf("%s", <-ch)
	}
}
