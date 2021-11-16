// Package check defines how to check an IP address. It provides types and
// functions that are useful for writing checks.
package check

import (
	"net"
)

const (
	TypeInfo    Type = iota // provides generic information about the IP address
	TypeSec                 // says whether the IP address is considered malicious
	TypeInfoSec             // provides both generic and security information about the IP address
)

// Type is the type of a check.
type Type int

// String returns the name of the check type: Info, Sec or InfoSec.
func (t Type) String() string {
	switch t {
	case TypeInfo:
		return "Info"
	case TypeSec:
		return "Sec"
	case TypeInfoSec:
		return "InfoSec"
	default:
		return "Unknown check type"
	}
}

// Check checks an IP address providing generic and/or security information.
type Check func(ipaddr net.IP) (Result, error)

// Result is the results of a check.
type Result struct {
	Name      string `json:"name"`      // check name
	Type      Type   `json:"type"`      // check type
	Info      Info   `json:"info"`      // provided by TypeInfo check
	Malicious bool   `json:"malicious"` // provided by TypeSec check
}

// Info is some generic information provided by an Info check.
type Info interface {
	Summary() string
	JsonString() (string, error) // all data in JSON format
}

// EmptyInfo is returned by checks that don't provide generic information about
// an IP address.
type EmptyInfo struct {
}

func (EmptyInfo) Summary() string {
	return Na("")
}

func (EmptyInfo) JsonString() (string, error) {
	return "{}", nil
}

// Na returns "n/a" if s is empty.
func Na(s string) string {
	if s == "" {
		return "n/a"
	}
	return s
}

// NonEmpty returns strings that are not empty.
func NonEmpty(strings ...string) []string {
	var ss []string
	for _, s := range strings {
		if s != "" {
			ss = append(ss, s)
		}
	}
	return ss
}
