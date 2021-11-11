// Package check checks an IP address using various public services.
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

type EmptyInfo struct {
}

func (EmptyInfo) String() string {
	return Na("")
}

func (EmptyInfo) JsonString() (string, error) {
	return "{}", nil
}

func Na(s string) string {
	if s == "" {
		return "n/a"
	}
	return s
}

func NonEmpty(strings ...string) []string {
	var ss []string
	for _, s := range strings {
		if s != "" {
			ss = append(ss, s)
		}
	}
	return ss
}
