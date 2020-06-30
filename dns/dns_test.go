package dns

import (
	"net"
	"testing"
)

func TestForIP(t *testing.T) {
	type testpair struct {
		ip   string
		name string
	}
	testpairs := []testpair{
		{"1.1.1.1", "one.one.one.one."},
		{"8.8.8.8", "dns.google."},
	}
	for _, tp := range testpairs {
		d := New()
		ip := net.ParseIP(tp.ip)
		d.ForIP(ip)
		if d.Names[0] != tp.name {
			t.Errorf("%s didn't resolve to %s but to %s", tp.ip, tp.name, d.Names[0])
		}
	}
}
