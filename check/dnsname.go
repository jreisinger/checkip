package check

import (
	"encoding/json"
	"net"
	"strings"

	"github.com/jreisinger/checkip"
)

// Names are the DNS names of the given IP address.
type Names []string

func (n Names) Summary() string {
	return strings.Join(n, ", ")
}

func (n Names) Json() ([]byte, error) {
	return json.Marshal(n)
}

// DnsName does a reverse lookup for a given IP address to get its names.
func DnsName(ipaddr net.IP) (checkip.Result, error) {
	result := checkip.Result{
		Name: "dns name",
		Type: checkip.TypeInfo,
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
	result.Info = Names(names)

	return result, nil
}
