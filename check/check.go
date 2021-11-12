// Package check defines how to check an IP address.
package check

import (
	"net"
)

type Type string

const (
	TypeInfo    Type = "Info" // provides some useful information about the IP address
	TypeSec     Type = "Sec"  // says whether the IP address is considered malicious
	TypeInfoSec Type = "InfoSec"
)

// Check checks an IP address providing generic and/or security information.
type Check func(ipaddr net.IP) (Result, error)

type Result struct {
	Name            string
	Type            Type
	IPaddrMalicious bool
	Info            Info
}

type Info interface {
	String() string
	JsonString() (string, error)
}

// EmptyInfo is returned by checks that don't provide generic information about
// an IP address.
type EmptyInfo struct {
}

func (EmptyInfo) String() string {
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
