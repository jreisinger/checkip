package checkip

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
)

// WellKnown implements checkip.Info.
type WellKnown bool

func (wk WellKnown) Json() ([]byte, error) {
	return json.Marshal(wk)
}

func (wk WellKnown) Summary() string {
	return fmt.Sprintf("%v", wk)
}

// IsWellKnown implements checkip.Check.
func IsWellKnown(ipaddr net.IP) (Result, error) {
	res := Result{Name: "well known"}

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

func Example() {
	ipaddr := net.ParseIP("1.1.1.1")
	result, err := IsWellKnown(ipaddr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result)
	// Output: {well known Info false true}
}
