package checkip

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

// IP holds information from net.IP.
type IP struct {
	Private     bool
	DefaultMask net.IPMask
}

// Check fills in IP data.
func (i *IP) Check(ipaddr net.IP) error {
	i.Private = ipaddr.IsPrivate()
	i.DefaultMask = ipaddr.DefaultMask()
	return nil
}

// Info returns interesting information from the check.
func (i *IP) Info() string {
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
