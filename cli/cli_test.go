package cli

import (
	"testing"

	"github.com/jreisinger/checkip/check"
)

func TestMaliciousStatsSkipsChecksWithMissingCredentials(t *testing.T) {
	checks := Checks{
		{
			Description:        "skipped",
			Type:               check.InfoAndIsMalicious,
			MissingCredentials: "API_KEY",
		},
		{
			Description:       "signal",
			Type:              check.IsMalicious,
			IpAddrIsMalicious: true,
		},
		{
			Description: "info",
			Type:        check.Info,
		},
	}

	total, malicious, prob := checks.maliciousStats()

	if total != 1 {
		t.Fatalf("total = %d, want 1", total)
	}
	if malicious != 1 {
		t.Fatalf("malicious = %d, want 1", malicious)
	}
	if prob != 1 {
		t.Fatalf("prob = %v, want 1", prob)
	}
}

func TestMaliciousStatsReturnsZeroWhenNoChecksRan(t *testing.T) {
	checks := Checks{
		{
			Description:        "skipped",
			Type:               check.IsMalicious,
			MissingCredentials: "API_KEY",
		},
	}

	total, malicious, prob := checks.maliciousStats()

	if total != 0 {
		t.Fatalf("total = %d, want 0", total)
	}
	if malicious != 0 {
		t.Fatalf("malicious = %d, want 0", malicious)
	}
	if prob != 0 {
		t.Fatalf("prob = %v, want 0", prob)
	}
}
