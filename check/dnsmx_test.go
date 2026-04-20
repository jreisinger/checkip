package check

import (
	"net"
	"reflect"
	"testing"
)

func TestDnsMXUsesReverseLookupWithoutAbuseIPDBInfo(t *testing.T) {
	origLookupAddr := dnsLookupAddr
	origLookupMX := dnsLookupMX
	origAbuseIPDBCheck := abuseIPDBCheck
	dnsLookupAddr = func(string) ([]string, error) {
		return []string{"www.example.com."}, nil
	}
	dnsLookupMX = func(name string) ([]*net.MX, error) {
		switch name {
		case "www.example.com":
			return []*net.MX{{Host: "mx-www.example.com."}}, nil
		case "example.com":
			return []*net.MX{{Host: "mx.example.com."}}, nil
		default:
			return nil, nil
		}
	}
	abuseIPDBCheck = func(net.IP) (Check, error) {
		return Check{}, nil
	}
	t.Cleanup(func() {
		dnsLookupAddr = origLookupAddr
		dnsLookupMX = origLookupMX
		abuseIPDBCheck = origAbuseIPDBCheck
	})

	result, err := DnsMX(net.ParseIP("1.2.3.4"))
	if err != nil {
		t.Fatalf("DnsMX returned error: %v", err)
	}

	info, ok := result.IpAddrInfo.(mx)
	if !ok {
		t.Fatalf("IpAddrInfo type = %T, want mx", result.IpAddrInfo)
	}

	want := map[string][]string{
		"www.example.com": {"mx-www.example.com"},
		"example.com":     {"mx.example.com"},
	}
	if !reflect.DeepEqual(info.Records, want) {
		t.Fatalf("records = %#v, want %#v", info.Records, want)
	}
}

func TestMxSummaryReturnsEmptyWhenNoMXRecordsExist(t *testing.T) {
	info := mx{Records: map[string][]string{
		"one.one.one.one": {},
	}}

	if got := info.Summary(); got != "" {
		t.Fatalf("Summary() = %q, want empty string", got)
	}
}
