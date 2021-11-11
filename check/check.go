// Package check checks an IP address using various public services. An IP
// address is checked by running one or more Checkers.
package check

import (
	"net"
	"sync"
)

const (
	TypeInfo    Type = "Info" // provides some useful information about the IP address
	TypeSec     Type = "Sec"  // says whether the IP address is considered malicious
	TypeInfoSec Type = "InfoSec"
)

type Type string

type Check func(ipaddr net.IP) Result

// Run runs checkers concurrently checking the ipaddr.
func Run(checks []Check, ipaddr net.IP) Results {
	var res []Result

	var wg sync.WaitGroup
	for _, chk := range checks {
		wg.Add(1)
		go func(c Check) {
			defer wg.Done()
			res = append(res, c(ipaddr))
		}(chk)
	}
	wg.Wait()
	return res
}
