package checker

import "github.com/jreisinger/checkip/check"

var DefaultCheckers = []check.Check{
	CheckAbuseIPDB,
	CheckAs,
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
