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
	wellKnown := []net.IP{
		net.ParseIP("1.1.1.1"),
		net.ParseIP("4.4.4.4"),
		net.ParseIP("8.8.8.8"),
	}
	var known WellKnown
	for _, wk := range wellKnown {
		if string(ipaddr) == string(wk) {
			known = true
		}
	}
	return checkip.Result{Name: "well known", Info: known}, nil
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

func ExampleCheck() {
	ipaddr := net.ParseIP("2.2.2.2")
	result, _ := IsWellKnown(ipaddr)
	fmt.Println(result)
	// Output: {well known Info false false}
}
