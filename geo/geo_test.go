package geo

import (
	"net"
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	g := New()
	if g.Filepath != "/var/tmp/GeoLite2-City.mmdb" {
		t.Errorf("default geodb path is wrong: %s", g.Filepath)
	}
}

func TestForIP(t *testing.T) {
	if os.Getenv("GEOIP_LICENSE_KEY") == "" {
		t.Skip("skipping test; GEOIP_LICENSE_KEY is not set")
	}
	type testpair struct {
		ip    string
		state string
	}
	testpairs := []testpair{
		{"1.1.1.1", "Australia"},
		{"8.8.8.8", "United States"},
	}
	for _, tp := range testpairs {
		g := New()
		ip := net.ParseIP(tp.ip)
		g.ForIP(ip)
		if g.Location[1] != tp.state {
			t.Errorf("%s is expected to be in %s but is in %s", tp.ip, tp.state, g.Location[1])
		}
	}

}
