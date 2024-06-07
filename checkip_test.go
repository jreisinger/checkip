package checkip_test

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/jreisinger/checkip/check"
	"github.com/jreisinger/checkip/cli"
)

// Example shows how to run two IP address checks
//
//   - CheckIsWellKnown that we define here
//   - check.DnsName
func Example() {
	ipaddr := net.ParseIP("1.1.1.1")
	results, _ := cli.Run([]check.Func{CheckIsWellKnown, check.DnsName}, ipaddr)
	results.PrintSummary()
	// Output: well known      true
	// dns name        one.one.one.one
}

func CheckIsWellKnown(ipaddr net.IP) (check.Check, error) {
	res := check.Check{Description: "well known"}

	wellKnown := []net.IP{
		net.ParseIP("1.1.1.1"),
		net.ParseIP("4.4.4.4"),
		net.ParseIP("8.8.8.8"),
	}

	for _, wk := range wellKnown {
		if string(ipaddr) == string(wk) {
			res.IpAddrInfo = IsWellKnown(true)
			break
		}
	}

	return res, nil
}

type IsWellKnown bool

func (iwk IsWellKnown) Json() ([]byte, error) {
	return json.Marshal(iwk)
}

func (iwk IsWellKnown) Summary() string {
	return fmt.Sprintf("%v", iwk)
}
