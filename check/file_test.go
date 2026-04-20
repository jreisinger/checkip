package check

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
		got, err := getCachePath(tc.filename)
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

func TestDownloadFile(t *testing.T) {
	t.Run("given valid response then result and no error is returned", func(t *testing.T) {
		handlerFn := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.WriteHeader(http.StatusOK)
			rw.Write(loadResponse(t, "file.txt"))
		})

		testUrl := setMockHttpClient(t, handlerFn)
		rc, err := downloadFile(testUrl)
		require.NoError(t, err)
		defer rc.Close()

		b, err := io.ReadAll(rc)
		require.NoError(t, err)
		assert.Equal(t, "just a simple file", string(b))
	})
}

func TestStoreFileReturnsCreateError(t *testing.T) {
	outFile := filepath.Join(t.TempDir(), "missing", "file.txt")

	err := storeFile(outFile, io.NopCloser(strings.NewReader("content")))

	require.Error(t, err)
}

func TestExtractGzFileReturnsCreateError(t *testing.T) {
	var compressed bytes.Buffer
	zw := gzip.NewWriter(&compressed)
	_, err := zw.Write([]byte("content"))
	require.NoError(t, err)
	require.NoError(t, zw.Close())

	outFile := filepath.Join(t.TempDir(), "missing", "file.txt")

	err = extractGzFile(outFile, io.NopCloser(bytes.NewReader(compressed.Bytes())))

	require.Error(t, err)
}
