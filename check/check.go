// Package check contains functions that can check an IP address.
package check

import (
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

// Default contains checks recommended for most people.
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
	ThreadCrowd,
	Tls,
}

// na returns "n/a" if s is empty.
func na(s string) string {
	if s == "" {
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
