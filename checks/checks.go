// Package checks contains functions that check an IP address providing generic
// and/or security information about it.
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
	CheckOTX,
	CheckShodan,
	CheckThreadCrowd,
	CheckVirusTotal,
}
