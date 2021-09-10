package checkip

import (
	"fmt"
	"net"
)

// IP holds info from net.IP.
type IP struct {
	Private     bool
	DefaultMask net.IPMask
}

// Check fills in data to IP.
func (i *IP) Check(ipaddr net.IP) (bool, error) {
	i.Private = ipaddr.IsPrivate()
	i.DefaultMask = ipaddr.DefaultMask()
	return true, nil
}

// String returns the result of the check.
func (i *IP) String() string {
	return fmt.Sprintf("RFC 1918 private: %v, default mask: %d", i.Private, i.DefaultMask)
}
