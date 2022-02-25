// Package checks contains functions that check an IP address. Checks provide
// generic and/or security information about the IP address.
package checks

import "github.com/jreisinger/checkip/check"

// Passive checks don't interact directly with the target IP address.
var Passive = []check.Check{
	CheckAbuseIPDB,
	CheckAS,
	CheckBlockList,
	CheckCins,
	CheckGeo,
	CheckIPSum,
	CheckMX,
	CheckName,
	CheckOTX,
	CheckShodan,
	CheckThreadCrowd,
	CheckUrlscan,
	CheckVirusTotal,
}

// Active checks interact with the target IP address.
var Active = []check.Check{
	CheckPing,
	CheckTcpPorts,
}
