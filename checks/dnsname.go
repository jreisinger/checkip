package checks

import (
	"encoding/json"
	"net"
	"strings"

	"github.com/jreisinger/checkip/check"
)

// Names are the DNS names of the given IP address.
type Names []string

func (n Names) Summary() string {
	return check.Na(strings.Join(n, ", "))
}

func (n Names) JsonString() (string, error) {
	b, err := json.Marshal(n)
	return string(b), err
}

// DnsName does a reverse lookup for a given IP address to get its names.
func DnsName(ipaddr net.IP) (check.Result, error) {
	result := check.Result{
		Name: "dns name",
		Type: check.TypeInfo,
	}

	names, _ := net.LookupAddr(ipaddr.String())
	for i := range names {
		names[i] = strings.TrimSuffix(names[i], ".")
	}
	result.Info = Names(names)

	return result, nil
}
