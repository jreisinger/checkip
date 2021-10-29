// Checkip finds information on an IP address using various public services.
package checkip

import (
	"net"
)

// Checker runs a check of an IP address.
type Checker interface {
	Check(ip net.IP) error
}

// InfoChecker finds information about an IP address.
type InfoChecker interface {
	Info() string
	Checker
}

// SecChecker checks an IP address is ok from the security point of view.
type SecChecker interface {
	IsOK() bool
	Checker
}
