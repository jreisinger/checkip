package asn

import (
	"net"
	"testing"
)

func TestForIP(t *testing.T) {
	a := New()
	ip := net.ParseIP("1.1.1.1")
	err := a.ForIP(ip)
	if err != nil {
		t.Errorf("error getting asn info: %v", err)
	}
	if a.CountryCode != "CN" {
		t.Errorf("expected country code waf CN, got %s", a.CountryCode)
	}
}
