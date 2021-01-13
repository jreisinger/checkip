package check

import (
	"fmt"
	"net"
	"strings"
)

// DNS holds the DNS names of the checked IP address.
type DNS struct {
	Names []string
}

// Do looks up the DNS names for the given IP address.
func (d *DNS) Do(ipaddr net.IP) (bool, error) {
	names, err := net.LookupAddr(ipaddr.String())
	if err != nil {
		return false, err
	}
	d.Names = names
	return true, nil
}

// Name returns the name of the check.
func (d *DNS) Name() string {
	return fmt.Sprint("DNS")
}

// String returns the result of the check.
func (d *DNS) String() string {
	return fmt.Sprintf("%s", strings.Join(d.Names, " | "))
}
