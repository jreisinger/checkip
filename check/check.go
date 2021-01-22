// Package check allows you to run various IP address checks.
package check

import (
	"fmt"
	"net"

	. "github.com/logrusorgru/aurora"
)

// Check represents an IP address checker.
type Check interface {
	Do(addr net.IP) (bool, error)
	Name() string
	String() string // result of the check
}

// Run runs a check of an IP address and returns the result over a channel.
func Run(chk Check, ipaddr net.IP, ch chan string) {
	format := "%-11s %s\n"
	ok, err := chk.Do(ipaddr)
	if err != nil {
		ch <- fmt.Sprintf(format, Gray(11, chk.Name()), err)
		return
	}
	if ok {
		ch <- fmt.Sprintf(format, chk.Name(), chk)
	} else {
		ch <- fmt.Sprintf(format, Magenta(chk.Name()), chk)
	}
}

// GetAvailable returns all available checks.
func GetAvailable() []Check {
	availableChecks := []Check{
		&AbuseIPDB{},
		&AS{},
		&DNS{},
		&Geo{},
		&IPsum{},
		&OTX{},
		&ThreatCrowd{},
		&VirusTotal{},
	}
	return availableChecks
}
