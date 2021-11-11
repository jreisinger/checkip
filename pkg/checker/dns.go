package checker

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"

	"github.com/jreisinger/checkip/pkg/check"
)

// DNS holds the DNS names from net.LookupAddr.
type DNS struct {
	Names []string `json:"names"`
}

func (d DNS) String() string {
	msg := "DNS name"
	if len(d.Names) > 1 {
		msg += "s"
	}
	return fmt.Sprintf("%s: %s", msg, check.Na(strings.Join(d.Names, ", ")))
}

func (d DNS) JsonString() (string, error) {
	b, err := json.Marshal(d)
	return string(b), err
}

// CheckDNS does a reverse lookup for a given IP address.
func CheckDNS(ipaddr net.IP) check.Result {
	// NOTE: We are ignoring error. It says: "nodename nor servname
	// provided, or not known" if there is no DNS name for the IP address.
	names, _ := net.LookupAddr(ipaddr.String())

	return check.Result{
		CheckName: "net.LookupAddr",
		CheckType: check.TypeInfo,
		Data:      DNS{Names: names},
	}
}
