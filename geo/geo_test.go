package geo

import (
	"testing"
)

func TestNew(t *testing.T) {
	g := New()
	if g.Filepath != "/var/tmp/GeoLite2-City.mmdb" {
		t.Errorf("default geodb path is wrong: %s", g.Filepath)
	}
}
