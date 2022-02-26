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

// CheckName does a reverse lookup for a given IP address to get its names.
func CheckName(ipaddr net.IP) (check.Result, error) {
	names, err := net.LookupAddr(ipaddr.String())
	if err != nil {
		return check.Result{}, check.NewError(err)
	}

	for i := range names {
		names[i] = strings.TrimSuffix(names[i], ".")
	}

	return check.Result{
		Name: "dns name",
		Type: check.TypeInfo,
		Info: Names(names),
	}, nil
}
