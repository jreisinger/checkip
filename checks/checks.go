// Package checks contains functions that check an IP address. Checks provide
// generic and/or security information about the IP address.
package checks

import "github.com/jreisinger/checkip/check"

// Default checks you should use.
var Default = []check.Check{
	CheckAbuseIPDB,
	CheckAS,
	CheckBlockList,
	CheckCins,
	CheckDNS,
	CheckGeo,
	CheckIPSum,
	CheckMX,
	CheckOTX,
	CheckPing,
	CheckShodan,
	CheckThreadCrowd,
	CheckUrlscan,
	CheckVirusTotal,
}
