package check_test

import (
	"net"
	"testing"

	"github.com/jreisinger/checkip/check"
)

func TestIsOnAWS(t *testing.T) {
	tests := []struct {
		ipAddr  string
		summary string
	}{
		{
			ipAddr:  "54.74.0.27",
			summary: "true, prefix: 54.74.0.0/15, region: eu-west-1, sevices: [AMAZON EC2]",
		},
		{
			ipAddr:  "1.1.1.1",
			summary: "false",
		},
	}

	for i, test := range tests {
		ipaddr := net.ParseIP(test.ipAddr)
		r, err := check.IsOnAWS(ipaddr)
		if err != nil {
			t.Fatalf("check.IsOnAWS(%s) failed: %v", ipaddr, err)
		}
		if r.Info.Summary() != test.summary {
			t.Errorf("test %d\ngot\t%q\nwant\t%q", i, r.Info.Summary(), test.summary)
		}
	}
}
