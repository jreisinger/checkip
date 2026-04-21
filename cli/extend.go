// Package cli contains functions for running checks from command-line.
package cli

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"regexp"
	"strings"

    "github.com/logrusorgru/aurora"
	"github.com/jreisinger/checkip/check"
)

type Rating struct {
	Type     string
	TypeDesc string
	ASNn     string
	ASNd     string
	Country  string
}

func (rs Checks) Cotation() Rating {
	var res Rating
	for _, r := range rs {
		// To avoid "invalid memory address or nil pointer dereference"
		// runtime error and printing empty summary info.
		if r.IpAddrInfo == nil || r.IpAddrInfo.Summary() == "" {
			continue
		}

		summary := r.IpAddrInfo.Summary()
		if r.Description == "IOCLoc" {
			re := regexp.MustCompile(`\(.*([A-Z]{2}?).*\).*(AS\d+) - (.*)$`)
			submatch := re.FindStringSubmatch(summary)
			if len(submatch) > 0 {
				res.Country = submatch[1]
				res.ASNn = submatch[2]
				res.ASNd = submatch[3]
			}
			continue
		}
		desc := strings.ToLower(summary)
		switch {
		case res.Type == "" && (strings.Contains(desc, "data center") || strings.Contains(desc, "network")):
			res.Type = "A"
			res.TypeDesc = "server"
		case res.Type == "" && strings.Contains(r.IpAddrInfo.Summary(), "open:"):
			res.Type = "A"
			res.TypeDesc = "server"
		case strings.Contains(desc, "vpn") || strings.Contains(desc, "avast"):
			res.Type = "B"
			res.TypeDesc = "vpn"
		case strings.Contains(desc, "mikrotik") || strings.Contains(desc, "fixed line"):
			res.Type = "C"
			res.TypeDesc = "botnet"
		case strings.Contains(desc, "mobile"):
			res.Type = "D"
			res.TypeDesc = "mobile"
		case strings.Contains(desc, "akamai") || strings.Contains(desc, "amazon") || strings.Contains(desc, "content delivery"):
			res.Type = "E"
			res.TypeDesc = "cdn"
		}

	}
	return res
}

// PrintExtJSON prints detailed results in JSON format.
func (checks Checks) PrintExtJSON(ipaddr net.IP, rating Rating) {
	// if len(rs) == 0 {
	// 	return
	// }

	_, _, prob := checks.maliciousStats()

	out := struct {
		IpAddr        net.IP `json:"ipAddr"`
		MaliciousProb string `json:"maliciousProb"`
		Checks        Checks `json:"checks"`
		Cotation      Rating `json:"cotation"`
	}{
		ipaddr,
		fmt.Sprintf("%.2f", prob),
		checks,
		rating,
	}

	enc := json.NewEncoder(os.Stdout)
	if err := enc.Encode(out); err != nil {
		log.Fatal(err)
	}
}

// PrintMISPJSON prints detailed results in JSON format.
func (checks Checks) PrintMISPJSON(ipaddr net.IP, rating Rating) {
	// if len(rs) == 0 {
	// 	return
	// }

	type tag struct {
		Name string `json:"name"`
	}

	type Attribute struct {
		Type         string `json:"type"`
		Category     string `json:"category"`
		ToIds        bool   `json:"to_ids"`
		Distribution string `json:"distribution"`
		Comment      string `json:"comment"`
		Value        net.IP `json:"value"`
		Tag          []tag  `json:"Tag"`
	}

	attribute := Attribute{Type: "ip-src", Category: "Network activity", Distribution: "5", Value: ipaddr}
	attribute.Comment = fmt.Sprintf("[%s1 - %s] (%s) %s - %s", rating.Type, rating.TypeDesc, rating.Country, rating.ASNn, rating.ASNd)
	attribute.Tag = append(attribute.Tag, tag{Name: fmt.Sprintf("quoting-scale:source-stability=\"%s\"", rating.Type)})
	attribute.Tag = append(attribute.Tag, tag{Name: "quoting-scale:information-precision=\"1\""})

	enc := json.NewEncoder(os.Stdout)
	if err := enc.Encode(attribute); err != nil {
		log.Fatal(err)
	}
}

// ExtPrintSummary add IpAddrInfo.Summary for IOCLoc check
func (rs Checks) ExtPrintSummary() string {
	au := aurora.NewAurora(true)
	res := ""
	for _, r := range rs {
		// To avoid "invalid memory address or nil pointer dereference"
		// runtime error and printing empty summary info.
		if r.IpAddrInfo == nil || r.IpAddrInfo.Summary() == "" {
			continue
		}

		if r.Type == check.Info || r.Type == check.InfoAndIsMalicious {
			if r.Description == "MyDB" || r.Description == "Misp" {
				r.Description = fmt.Sprintf("<%s>\t\t",au.Yellow(r.Description))
			}
			fmt.Printf("%-15s %s\n", r.Description, r.IpAddrInfo.Summary())
			if r.Description == "IOCLoc" {
				res = r.IpAddrInfo.Summary()
			}

		}

	}
	return res
}
