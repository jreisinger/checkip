package checkip

import (
	"fmt"
	"net"
	"strconv"
	"strings"
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
	private := "RFC 1918 private"
	if !i.Private {
		private = "not " + private
	}
	var mask []string
	for _, b := range i.DefaultMask {
		mask = append(mask, strconv.Itoa(int(b)))
	}
	return fmt.Sprintf("IP address\t%s, default mask %s", private, strings.Join(mask, "."))
}
