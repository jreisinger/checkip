// Package checkip defines Check type, which provides information about an IP
// address.
package checkip

import (
	"net"
)

// Existing Check types.
const (
	TypeInfo    Type = iota // generic information about the IP address
	TypeSec                 // whether the IP address is considered malicious
	TypeInfoSec             // both of the above
)

// Type is the type of a Check.
type Type int32

// String returns the name of the Check type: Info, Sec or InfoSec.
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

// Result is the information returned by a Check.
type Result struct {
	Name      string `json:"name"`      // check name
	Type      Type   `json:"type"`      // check type
	Info      Info   `json:"info"`      // provided by TypeInfo check
	Malicious bool   `json:"malicious"` // provided by TypeSec check
}

// Info is some generic information provided by an Info or InfoSec Check.
type Info interface {
	Summary() string             // summary info
	JsonString() (string, error) // all info in JSON format
}

// EmptyInfo is returned by TypeSec Checks that don't provide generic
// information about an IP address.
type EmptyInfo struct {
}

// Summary returns empty string.
func (EmptyInfo) Summary() string {
	return ""
}

// JsonString returns empty JSON string.
func (EmptyInfo) JsonString() (string, error) {
	return "{}", nil
}
