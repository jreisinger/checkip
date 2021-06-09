package checkip

import (
	"net"
	"testing"
)

func TestAS(t *testing.T) {
	testIPs := []string{"1.1.1.1", "8.8.8.8"}
	for _, ip := range testIPs {
		a := &AS{}
		ip := net.ParseIP(ip)
		_, err := a.Check(ip)
		if err != nil {
			t.Errorf("checking %s: %v", ip, err)
		}
		if a.CountryCode == "" {
			t.Errorf("country code for %s is empty", ip)
		}
		if a.Description == "" {
			t.Errorf("description for %s is empty", ip)
		}
	}
}
