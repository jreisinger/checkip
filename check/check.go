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

// Funcs contains all available functions for checking IP addresses.
var Funcs = []Func{
	AbuseIPDB,
	BlockList,
	CinsScore,
	DBip,
	DnsMX,
	DnsName,
	Firehol,
	IPSum,
	IPtoASN,
	IsOnAWS,
	MaxMind,
	OTX,
	Ping,
	SansISC,
	Shodan,
	Spur,
	Tls,
	UrlScan,
	VirusTotal,
	GreyNoise,
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
