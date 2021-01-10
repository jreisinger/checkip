package check

import (
	"net"
	"testing"
)

func TestGeoCheck(t *testing.T) {
	type testpair struct {
		ip    string
		state string
	}
	testpairs := []testpair{
		{"1.1.1.1", "Australia"},
		{"8.8.8.8", "United States"},
	}
	for _, tp := range testpairs {
		g := &Geo{}
		ip := net.ParseIP(tp.ip)
		g.Check(ip)
		if g.Location[1] != tp.state {
			t.Errorf("%s is expected to be in %s but is in %s", tp.ip, tp.state, g.Location[1])
		}
	}

}
