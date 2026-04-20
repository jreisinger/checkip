// Package check contains types and functions for getting information on IP addresses.
package check

import (
	"encoding/json"
	"fmt"
	"net"
	"regexp"
)

// Type of a check.
const (
	Info               Type = iota // some information about the IP address
	IsMalicious                    // whether the IP address is considered malicious
	InfoAndIsMalicious             // both of the above
)

// CachePolicy says whether a check result can be reused within one process.
type CachePolicy int

const (
	// CacheProcess reuses the check result for repeated IP addresses within a
	// single process.
	CacheProcess CachePolicy = iota
	// CacheNone always runs the check live.
	CacheNone
)

// Definition describes a check and how it should be executed.
type Definition struct {
	// Name should be unique across all registered checks.
	Name string
	Run  Func
	// Cache defaults to CacheProcess.
	Cache CachePolicy
}

// Definitions contains all available checks and their execution policy.
var Definitions = []Definition{
	{Name: "abuseipdb.com", Run: AbuseIPDB},
	{Name: "blocklist.de", Run: BlockList},
	{Name: "cinsscore.com", Run: CinsScore},
	{Name: "db-ip.com", Run: DBip},
	{Name: "dns MX", Run: DnsMX},
	{Name: "dns name", Run: DnsName},
	{Name: "firehol.org", Run: Firehol},
	{Name: "ipsum.app", Run: IPSum},
	{Name: "iptoasn.com", Run: IPtoASN},
	{Name: "is on AWS", Run: IsOnAWS},
	{Name: "maxmind.com", Run: MaxMind},
	{Name: "otx.alienvault.com", Run: OTX},
	{Name: "ping", Run: Ping},
	{Name: "isc.sans.edu", Run: SansISC},
	{Name: "shodan.io", Run: Shodan},
	{Name: "spur.us", Run: Spur},
	{Name: "tls", Run: Tls},
	{Name: "urlscan.io", Run: UrlScan},
	{Name: "virustotal.com", Run: VirusTotal},
}

// Funcs contains all available check functions, derived from Definitions for
// backward compatibility.
var Funcs = funcs(Definitions)

func funcs(definitions []Definition) []Func {
	funcs := make([]Func, 0, len(definitions))
	for _, definition := range definitions {
		funcs = append(funcs, definition.Run)
	}
	return funcs
}

// Type is the type of a Check.
type Type int

// String returns the name of the Check type.
func (t Type) String() string {
	return [...]string{"Info", "IsMalicious", "InfoAndIsMalicious"}[t]
}

func (t Type) MarshalJSON() ([]byte, error) {
	s := fmt.Sprint(t)
	return json.Marshal(s)
}

// Func gathers generic and/or security information about an IP address.
type Func func(ipaddr net.IP) (Check, error)

// Check contains information on the check itself and
// the obtained information about an IP address
type Check struct {
	Description        string `json:"description"` // max 15 chars
	Type               Type   `json:"type"`
	MissingCredentials string `json:"missingCredentials,omitempty"`
	IpAddrIsMalicious  bool   `json:"ipAddrIsMalicious"`
	IpAddrInfo         IpInfo `json:"ipAddrInfo"`
}

// IpInfo is generic information on an IP address.
type IpInfo interface {
	Summary() string       // summary info
	Json() ([]byte, error) // all info in JSON format
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
