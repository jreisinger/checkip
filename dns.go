package checkip

import (
	"fmt"
	"net"
	"strings"
)

// DNS holds the DNS names from net.LookupAddr.
type DNS struct {
	Names []string
}

// Check does a reverse lookup for a given IP address.
func (d *DNS) Check(ipaddr net.IP) (bool, error) {
	// We are ignoring error. It says: nodename nor servname provided, or
	// not known
	names, _ := net.LookupAddr(ipaddr.String())
	d.Names = names
	return true, nil
}

// String returns the result of the check.
func (d *DNS) String() string {
	return fmt.Sprintf("DNS names\t%s", strings.Join(d.Names, ", "))
}
