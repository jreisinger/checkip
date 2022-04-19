package check

import (
	"strings"
	"testing"
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
