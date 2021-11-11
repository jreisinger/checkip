package checks

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"

	"github.com/jreisinger/checkip/check"
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
func CheckDNS(ipaddr net.IP) (check.Result, error) {
	// NOTE: We are ignoring error. It says: "nodename nor servname
	// provided, or not known" if there is no DNS name for the IP address.
	names, _ := net.LookupAddr(ipaddr.String())
	// if err != nil {
	// 	return check.Result{}, check.NewError(err)
	// }

	return check.Result{
		Name: "net.LookupAddr",
		Type: check.TypeInfo,
		Info: DNS{Names: names},
	}, nil
}
