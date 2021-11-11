// Package checks contains functions that can check an IP address.
package checks

import "github.com/jreisinger/checkip/check"

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
