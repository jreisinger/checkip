package check

import (
	"errors"
	"io"
	"net"
	"os"
	"path/filepath"
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

const subnet4 = "50.16.16.211"

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
		{
			strings.NewReader(subnet4),
			"50.16.16.211",
			true,
			false,
		},
		{
			strings.NewReader(subnet4),
			"1.1.1.1",
			false,
			false,
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
		if err == nil && tc.err {
			t.Fatal("expected error, got nil")
		}
		if err != nil && tc.err && !strings.Contains(err.Error(), "parse FireHOL entry") {
			t.Fatalf("error %q doesn't include FireHOL context", err)
		}
		if found != tc.found {
			t.Fatalf("%s found %t but expected %t", tc.ipaddr, found, tc.found)
		}
	}
}

func TestFireholPrefixesErrorsWithCheckName(t *testing.T) {
	origGetCachePath := fireholGetCachePath
	origUpdateFile := fireholUpdateFile
	origSearch := fireholSearch
	fireholGetCachePath = func(string) (string, error) {
		return filepath.Join(t.TempDir(), "firehol_level1.netset"), nil
	}
	fireholUpdateFile = func(string, string, string) error {
		return nil
	}
	fireholSearch = func(io.Reader, net.IP) (bool, error) {
		return false, errors.New("boom")
	}
	t.Cleanup(func() {
		fireholGetCachePath = origGetCachePath
		fireholUpdateFile = origUpdateFile
		fireholSearch = origSearch
	})

	file := filepath.Join(t.TempDir(), "firehol_level1.netset")
	if err := os.WriteFile(file, []byte(""), 0600); err != nil {
		t.Fatal(err)
	}
	fireholGetCachePath = func(string) (string, error) {
		return file, nil
	}

	_, err := Firehol(net.ParseIP("1.1.1.1"))
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "Firehol: boom") {
		t.Fatalf("error = %q, want Firehol prefix", err)
	}
}
