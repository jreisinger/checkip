// Checkip checks an IP address using various public services.
package checkip

import (
	"fmt"
	"net"
	"sync"

	"github.com/logrusorgru/aurora"
)

// Checker checks an IP address. ok is false if it considers the IP address to
// be suspicious. If the check fails (err != nil), ok must be true - presumption
// of innocence. Checker can be printed to see what it has found about the IP
// address.
type Checker interface {
	Check(ip net.IP) (ok bool, err error)
	fmt.Stringer
}

// Run runs checkers concurrently and returns the number of checkers that
// consider the IP address to be suspicious.
func Run(checkers []Checker, ipaddr net.IP) int {
	var suspicious int
	var wg sync.WaitGroup
	for _, checker := range checkers {
		wg.Add(1)
		go func(checker Checker) {
			ok, err := checker.Check(ipaddr)
			if err == nil && !ok {
				suspicious++
			}
			wg.Done()
		}(checker)
	}
	return suspicious
}

// RunAndPrint runs checkers concurrently and print the results. checkers maps
// names to checkers. Format defines how to print the name and checker results
// (e.g. "%-25s %s").
func RunAndPrint(checkers map[string]Checker, ipaddr net.IP, format string) {
	var wg sync.WaitGroup
	format += "\n"
	for name, checker := range checkers {
		wg.Add(1)
		go func(checker Checker, name string) {
			ok, err := checker.Check(ipaddr)
			switch {
			case err != nil:
				fmt.Printf(format, name, aurora.Gray(11, err.Error()))
			case !ok:
				fmt.Printf(format, name, aurora.Magenta(checker.String()))
			default:
				fmt.Printf(format, name, checker)
			}
			wg.Done()
		}(checker, name)
	}
	wg.Wait()
}
