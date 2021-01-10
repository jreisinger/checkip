// Package check allows you to run various IP address checks.
package check

import (
	"fmt"
	"net"

	. "github.com/logrusorgru/aurora"
)

// Checker represents an IP address checker.
type Checker interface {
	Check(addr net.IP) (bool, error)
	Name() string
	String() string // output
}

// Run runs a Checker against and IP address and returns the result over a
// channel.
func Run(chkr Checker, ipaddr net.IP, ch chan string) {
	format := "%-11s %s\n"
	ok, err := chkr.Check(ipaddr)
	if err != nil {
		ch <- fmt.Sprintf(format, Gray(11, chkr.Name()), Gray(11, err))
		return
	}
	if ok {
		ch <- fmt.Sprintf(format, chkr.Name(), chkr)
	} else {
		ch <- fmt.Sprintf(format, Magenta(chkr.Name()), Magenta(chkr))
	}
}
