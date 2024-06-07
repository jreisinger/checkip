// Package checkip defines how to Check an IP address.
package checkip

import (
	"encoding/json"
	"net"
)

// Existing Check types.
const (
	TypeInfo    Type = iota // generic information about the IP address
	TypeSec                 // whether the IP address is considered malicious
	TypeInfoSec             // both of the above
)

// Type is the type of a Check.
type Type int

// String returns the name of the Check type.
func (t Type) String() string {
	return [...]string{"Info", "Sec", "InfoSec"}[t]
}

func (t Type) MarshalJSON() ([]byte, error) {
	var s string
	switch t {
	case TypeInfo:
		s = "info"
	case TypeSec:
		s = "security"
	case TypeInfoSec:
		s = "infoAndSecurity"
	}
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
