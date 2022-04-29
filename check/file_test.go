package check

import (
	"strings"
	"testing"
	"time"
)

func TestGetDbFilesPath(t *testing.T) {
	testcases := []struct {
		filename string
		suffix   string
	}{
		{"cins.txt", "cins.txt"},
		{"dbip-city-lite.mmdb", "dbip-city-lite.mmdb"},
		{"GeoLite2-City.mmdb", "GeoLite2-City.mmdb"},
	}
	for _, tc := range testcases {
		got, err := getDbFilesPath(tc.filename)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.HasSuffix(got, tc.suffix) {
			t.Fatalf("path %s doesn't have %s suffix", got, tc.suffix)
		}
	}
}

func TestIsOlderThanOneWeek(t *testing.T) {
	testcases := []struct {
		t                time.Time
		olderThanOneWeek bool
	}{
		{time.Now(), false},
		{time.Now().Add(-time.Hour * 24 * 6), false},
		{time.Now().Add(-time.Hour * 24 * 8), true},
	}
	for i, tc := range testcases {
		got := isOlderThanOneWeek(tc.t)
		if got != tc.olderThanOneWeek {
			t.Fatalf("test case %d: got %t expected %t", i+1, got, tc.olderThanOneWeek)
		}
	}
}
