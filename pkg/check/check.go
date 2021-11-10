// Package check checks an IP address using various public services. An IP address is
// checked by running one or more Checkers. There are two kinds of Checkers. An
// InfoChecker just gathers some useful information about the IP address. A
// SecChecker says whether the IP address is considered malicious or not.
package check

import (
	"net"
	"sync"
)

const (
	TypeInfoSec Type = "InfoSec"
	TypeInfo    Type = "Info"
	TypeSec     Type = "Sec"
)

type Type string

type Check func(ipaddr net.IP) Result

// Run runs checkers concurrently checking the ipaddr.
func Run(checkers []Check, ipaddr net.IP) Results {
	var res []Result

	var wg sync.WaitGroup
	for _, chk := range checkers {
		wg.Add(1)
		go func(c Check) {
			defer wg.Done()
			res = append(res, c(ipaddr))
		}(chk)
	}
	wg.Wait()
	return res
}
