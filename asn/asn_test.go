package asn

import (
	"net"
	"testing"
)

func TestForIP(t *testing.T) {
	type testpair struct {
		ip          string
		countryCode string
	}
	testpairs := []testpair{
		{"1.1.1.1", "US"},
		{"8.8.8.8", "US"},
	}
	for _, tp := range testpairs {
		a := New()
		ip := net.ParseIP(tp.ip)
		a.ForIP(ip)
		if a.CountryCode != tp.countryCode {
			t.Errorf("country code for %s is expected to be %s but is %s", tp.ip, tp.countryCode, a.CountryCode)
		}
	}
}
