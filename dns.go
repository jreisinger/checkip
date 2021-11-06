package checkip

import (
	"net"
	"strings"
)

// DNS holds the DNS names from net.LookupAddr.
type DNS struct {
	Names []string
}

func (d *DNS) String() string { return "net.LookupAddr" }

// Check does a reverse lookup for a given IP address.
func (d *DNS) Check(ipaddr net.IP) error {
	// NOTE: We are ignoring error. It says: "nodename nor servname
	// provided, or not known" if there is no DNS name for the IP address.
	names, _ := net.LookupAddr(ipaddr.String())
	d.Names = names
	return nil
}

// Info returns interesting information from the check.
func (d *DNS) Info() string {
	return strings.Join(d.Names, ", ")
}
