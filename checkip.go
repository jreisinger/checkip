// Checkip checks an IP address using various public services.
package checkip

import (
	"fmt"
	"net"
	"sync"
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
	var mu sync.Mutex
	var suspicious int
	var wg sync.WaitGroup
	for _, checker := range checkers {
		wg.Add(1)
		go func(checker Checker) {
			ok, err := checker.Check(ipaddr)
			if err == nil && !ok {
				mu.Lock()
				suspicious++
				mu.Unlock()
			}
			wg.Done()
		}(checker)
	}
	wg.Wait()
	return suspicious
}
