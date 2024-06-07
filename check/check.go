// Package check contains functions that can check an IP address.
package check

import (
	"encoding/json"
	"fmt"
	"net"
	"regexp"
)

// Existing Check types.
const (
	TypeInfo    Type = iota // generic information about the IP address
	TypeSec                 // whether the IP address is considered malicious
	TypeInfoSec             // both of the above
)

// Checks contains all available checks.
var Checks = []Check{
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
	PhishStats,
	Ping,
	SansISC,
	Shodan,
	Tls,
	UrlScan,
	VirusTotal,
}

// Type is the type of a Check.
type Type int

// String returns the name of the Check type.
func (t Type) String() string {
	return [...]string{"info", "sec", "infosec"}[t]
}

func (t Type) MarshalJSON() ([]byte, error) {
	s := fmt.Sprint(t)
	return json.Marshal(s)
}

// Check provides generic and/or security information about an IP address.
type Check func(ipaddr net.IP) (Result, error)

// Result is the information provided by a Check.
type Result struct {
	Name               string `json:"name"` // check name, max 15 chars
	Type               Type   `json:"type"` // check type
	MissingCredentials string `json:"missing_credentials,omitempty"`
	Malicious          bool   `json:"malicious"` // provided by TypeSec and TypeInfoSec check type
	Info               Info   `json:"info"`
}

// Info is generic information provided by a TypeInfo or TypeInfoSec Check.
type Info interface {
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
