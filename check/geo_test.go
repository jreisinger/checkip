package check

import (
	"net"
	"testing"

	"github.com/jreisinger/checkip/util"
)

func TestGeoCheck(t *testing.T) {
	// This is needed for the tests not to fail on travis-ci.org.
	if v, _ := util.GetConfigValue("GEOIP_LICENSE_KEY"); v == "" {
		t.Skip("skipping test; GEOIP_LICENSE_KEY is not set")
	}

	type testpair struct {
		ip      string
		country string
	}
	testpairs := []testpair{
		{"1.1.1.1", "Australia"},
		{"8.8.8.8", "United States"},
	}
	for _, tp := range testpairs {
		g := &Geo{}
		ip := net.ParseIP(tp.ip)
		g.Do(ip)
		if g.Country != tp.country {
			t.Errorf("%s is expected to be in %s but is in %s", tp.ip, tp.country, g.Country)
		}
	}

}
