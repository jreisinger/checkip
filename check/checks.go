// Package check contains functions that check an IP address.
package check

import (
	"github.com/jreisinger/checkip"
)

// Passive checks don't interact directly with the target IP address.
var Passive = []checkip.Check{
	AbuseIPDB,
	BlockList,
	CinsScore,
	DBip,
	DnsMX,
	DnsName,
	IPSum,
	IPtoASN,
	MaxMind,
	OTX,
	Shodan,
	ThreadCrowd,
	UrlScan,
	VirusTotal,
}

// Active checks interact with the target IP address.
var Active = []checkip.Check{
	Ping,
	TcpPorts,
}
