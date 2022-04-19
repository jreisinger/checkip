package check

import (
	"io"
	"net"
	"strings"
	"testing"
)

const subnet = `
# 192.168.1.0/24
192.168.2.0/24
192.168.3.0/24
`

const subnet2 = "27.124.64.0/20"

const subnet3 = "a"

func TestIpFound(t *testing.T) {
	testcases := []struct {
		subnets io.Reader
		ipaddr  string
		found   bool
		err     bool
	}{
		{
			strings.NewReader(subnet),
			"192.168.1.1",
			false,
			false,
		},
		{
			strings.NewReader(subnet),
			"192.168.2.1",
			true,
			false,
		},
		{
			strings.NewReader(subnet),
			"192.168.3.1",
			true,
			false,
		},

		{
			strings.NewReader(subnet2),
			"27.124.64.0", // Network
			true,
			false,
		},
		{
			strings.NewReader(subnet2),
			"27.124.64.1", // HostMin
			true,
			false,
		},
		{
			strings.NewReader(subnet2),
			"27.124.79.254", // HostMax
			true,
			false,
		},
		{
			strings.NewReader(subnet2),
			"27.124.79.255", // Broadcast
			true,
			false,
		},
		{
			strings.NewReader(subnet2),
			"1.1.1.1",
			false,
			false,
		},

		{
			strings.NewReader(subnet3),
			"1.1.1.1",
			false,
			true,
		},
	}
	for _, tc := range testcases {
		ipaddr := net.ParseIP(tc.ipaddr)
		if ipaddr == nil {
			t.Fatalf("can't parse IP address: %s", tc.ipaddr)
		}
		found, err := ipFound(tc.subnets, ipaddr)
		if err != nil && !tc.err {
			t.Fatal(err)
		}
		if found != tc.found {
			t.Fatalf("%s found %t but expected %t", tc.ipaddr, found, tc.found)
		}
	}
}
