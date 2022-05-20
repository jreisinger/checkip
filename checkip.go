// Package checkip defines how to Check an IP address.
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

// Check provides generic and/or security information about an IP address.
type Check func(ipaddr net.IP) (Result, error)

// Result is the information provided by a Check.
type Result struct {
	Name      string `json:"name"`      // check name, max 15 chars
	Type      Type   `json:"type"`      // check type
	Malicious bool   `json:"malicious"` // provided by TypeSec check type
	Info      Info   `json:"info"`
}

// Info is generic information provided by a TypeInfo or TypeInfoSec Check.
type Info interface {
	Summary() string       // summary info
	Json() ([]byte, error) // all info in JSON format
}
