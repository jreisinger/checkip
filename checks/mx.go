package checks

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"

	"github.com/jreisinger/checkip/check"
)

type MX struct {
	Servers []string `json:names`
}

func (m MX) Summary() string {
	msg := "mail server"
	if len(m.Servers) > 1 {
		msg += "s"
	}
	return fmt.Sprintf("%s: %s", msg, check.Na(strings.Join(m.Servers, ", ")))
}

func (m MX) JsonString() (string, error) {
	b, err := json.Marshal(m)
	return string(b), err
}

func CheckMX(ipaddr net.IP) (check.Result, error) {
	names, err := net.LookupAddr(ipaddr.String())
	if err != nil {
		return check.Result{}, check.NewError(err)
	}

	// Enrinch names with a name with 'www.' removed.
	// [www.csh.ac.at.] => [www.csh.ac.at. csh.ac.at.]
	for _, n := range names {
		t := strings.TrimPrefix(n, "www.")
		names = append(names, t)
	}

	var mx MX

	for _, n := range names {
		mxRecords, _ := net.LookupMX(n) // NOTE: ingoring error
		for _, r := range mxRecords {
			mx.Servers = append(mx.Servers, r.Host)
		}
	}

	return check.Result{
		Name: "mx",
		Type: check.TypeInfo,
		Info: mx,
	}, nil
}
