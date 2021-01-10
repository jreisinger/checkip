package check

import (
	"net"
	"testing"
)

func TestASNCheck(t *testing.T) {
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
		a := &AS{}
		ip := net.ParseIP(tp.ip)
		_, err := a.Check(ip)
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
