package check

import (
	"encoding/json"
	"net"
	"strings"
)

// dnsNames are the DNS names of the given IP address.
type dnsNames []string

func (n dnsNames) Summary() string {
	return strings.Join(n, ", ")
}

func (n dnsNames) Json() ([]byte, error) {
	return json.Marshal(n)
}

// DnsName does a reverse lookup for a given IP address to get its names.
func DnsName(ipaddr net.IP) (Check, error) {
	result := Check{
		Description: "dns name",
		Type:        Info,
	}

	names, err := net.LookupAddr(ipaddr.String())
	if err != nil {
		if len(names) == 0 {
			// IP address does not resolve to any names,
			// ignore this error.
		} else {
			return result, err
		}
	}
	for i := range names {
		names[i] = strings.TrimSuffix(names[i], ".")
	}
	result.IpAddrInfo = dnsNames(names)

	return result, nil
}
