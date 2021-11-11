package checker

import "github.com/jreisinger/checkip/check"

var DefaultCheckers = []check.Check{
	CheckAs,
	CheckAbuseIPDB,
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
