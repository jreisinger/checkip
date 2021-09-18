package checkip

import (
	"net"
	"testing"
)

func TestRun(t *testing.T) {
	// infoCheckers just give you information about an IP address. They
	// always return ok == true, i.e. they consider no IP suspicious.
	infoCheckers := []Checker{&AS{}, &DNS{}, &Geo{}, &IP{}}

	tests := []struct {
		ip         string
		suspicious int
	}{
		{"1.1.1.1", 0},
		{"1.1.1.2", 0},
		{"8.8.8.8", 0},
		{"8.8.4.4", 0},
	}

	for _, test := range tests {
		ipaddr := net.ParseIP(test.ip)
		got := Run(infoCheckers, ipaddr)
		if got != test.suspicious {
			t.Fatalf("%s considered suspiscious by %d info checkers", test.ip, got)
		}
	}
}
