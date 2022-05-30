package checkip_test

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/jreisinger/checkip"
	"github.com/jreisinger/checkip/cli"
)

// IsWellKnown implements checkip.Check.
func IsWellKnown(ipaddr net.IP) (checkip.Result, error) {
	res := checkip.Result{Name: "well known"}

	wellKnown := []net.IP{
		net.ParseIP("1.1.1.1"),
		net.ParseIP("4.4.4.4"),
		net.ParseIP("8.8.8.8"),
	}

	for _, wk := range wellKnown {
		if string(ipaddr) == string(wk) {
			res.Info = WellKnown(true)
		}
	}

	return res, nil
}

// WellKnown implements checkip.Info.
type WellKnown bool

func (wk WellKnown) Json() ([]byte, error) {
	return json.Marshal(wk)
}

func (wk WellKnown) Summary() string {
	return fmt.Sprintf("%v", wk)
}

func Example() {
	ipaddr := net.ParseIP("1.1.1.1")
	results, _ := cli.Run([]checkip.Check{IsWellKnown}, ipaddr)
	results.PrintSummary()
	// Output: well known      true
}
