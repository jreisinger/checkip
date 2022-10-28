// Package check contains functions that can check an IP address.
package check

import (
	"regexp"

	"github.com/jreisinger/checkip"
)

// All contains all available checks.
var All = []checkip.Check{
	AbuseIPDB,
	BlockList,
	CinsScore,
	DBip,
	DnsMX,
	DnsName,
	Firehol,
	IPSum,
	IPtoASN,
	MaxMind,
	OTX,
	PhishStats,
	Ping,
	Shodan,
	ThreadCrowd,
	Tls,
	UrlScan,
	VirusTotal,
}

// Default contains subset of all available checks recommended for most people.
var Default = []checkip.Check{
	BlockList,
	CinsScore,
	DBip,
	DnsName,
	Firehol,
	IPSum,
	IPtoASN,
	OTX,
	Ping,
	Shodan,
	Tls,
}

// na returns n/a if s is empty or contains only whitespace.
func na(s string) string {
	ws := regexp.MustCompile(`^\s+$`)
	if s == "" || ws.MatchString(s) {
		return "n/a"
	}
	return s
}

// nonEmpty returns strings that are not empty.
func nonEmpty(strings ...string) []string {
	var ss []string
	for _, s := range strings {
		if s != "" {
			ss = append(ss, s)
		}
	}
	return ss
}
