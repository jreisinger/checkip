package dns

import (
	"net"
	"testing"
)

func TestForIP(t *testing.T) {
	d := New()
	ip := net.ParseIP("1.1.1.1")
	err := d.ForIP(ip)
	if err != nil {
		t.Errorf("dns.ForIP doesn't work: %v", err)
	}
	if d.Names[0] != "one.one.one.one." {
		t.Errorf("1.1.1.1 didn't resolve to one.one.one.one. but to %s", d.Names[0])
	}
}
