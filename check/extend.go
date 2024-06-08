// Package check contains functions that can check an IP address.
package check

// Debug set by main flag
var Debug bool

// GetConfigValue export getConfigValue function for main
func GetConfigValue(key string) (string, error) {
        return getConfigValue(key)
}

// Use : function list to be used in checks
var Use = []Func{}

// AddUse : methode to add function in checks
func AddUse(s interface{}) {
        Use = append(Use,s.(Func))
}

// All contains all available checks.
var All = []Func{
	IOCLoc,
	Spur,
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
	IOCLoc,
	IsOnAWS,
	MaxMind,
	MyDB,
	Onyphe,
	OTX,
	PhishStats,
	Ping,
	SansISC,
	Shodan,
	Spur,
	Tls,
	UrlScan,
	VirusTotal,
}

