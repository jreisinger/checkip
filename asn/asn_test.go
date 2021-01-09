package asn

import (
	"net"
	"testing"
)

func TestForIP(t *testing.T) {
	type testpair struct {
		ip          string
		countryCode string
		description string
	}
	testpairs := []testpair{
		{"1.1.1.1", "US", "CLOUDFLARENET - Cloudflare, Inc."},
		{"8.8.8.8", "US", "GOOGLE - Google LLC"},
	}
	for _, tp := range testpairs {
		a := New()
		ip := net.ParseIP(tp.ip)
		err := a.ForIP(ip)
		if a.CountryCode != tp.countryCode || err != nil {
			t.Errorf("country code for %s was expected to be '%s' but is '%s' with %v",
				tp.ip, tp.countryCode, a.CountryCode, err)
		}
		if a.Description != tp.description || err != nil {
			t.Errorf("description for %s was expected to be '%s' but is '%s' with %v",
				tp.ip, tp.description, a.Description, err)
		}
	}
}
