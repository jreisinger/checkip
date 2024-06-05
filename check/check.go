// Package check contains functions that can check an IP address.
package check

import (
	"regexp"

	"github.com/jreisinger/checkip"
)

// Use : function list to be used in checks
var Use = []checkip.Check{}

// AddUse : methode to add function in checks
func AddUse(s interface{}) {
	Use = append(Use,s.(checkip.Check))
}

// All contains all available checks.
var All = []checkip.Check{
	AbuseIPDB,
	BlockList,
	CinsScore,
	Censys,
	DBip,
	DnsMX,
	DnsName,
	Firehol,
	IPSum,
	IPtoASN,
	IsOnAWS,
	MaxMind,
	OTX,
	Onyphe,
	PhishStats,
	Ping,
	SansISC,
	Shodan,
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
	IsOnAWS,
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
