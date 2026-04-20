package check

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (fn roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}

type errReadCloser struct {
	data []byte
	err  error
	read bool
}

func (r *errReadCloser) Read(p []byte) (int, error) {
	if r.read {
		return 0, r.err
	}

	r.read = true
	n := copy(p, r.data)
	return n, r.err
}

func (r *errReadCloser) Close() error {
	return nil
}

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
		oldClient := downloadHTTPClient
		downloadHTTPClient = &http.Client{
			Timeout: time.Second,
			Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, "https://example.com/file.txt", req.URL.String())
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader("just a simple file")),
					Header:     make(http.Header),
					Request:    req,
				}, nil
			}),
		}
		t.Cleanup(func() {
			downloadHTTPClient = oldClient
		})

		rc, err := downloadFile("https://example.com/file.txt")
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

func TestUpdateFileKeepsExistingFileWhenRefreshFails(t *testing.T) {
	cacheFile := filepath.Join(t.TempDir(), "cache.txt")
	require.NoError(t, os.WriteFile(cacheFile, []byte("cached content"), 0600))

	oldTime := time.Now().Add(-8 * 24 * time.Hour)
	require.NoError(t, os.Chtimes(cacheFile, oldTime, oldTime))

	oldClient := downloadHTTPClient
	downloadHTTPClient = &http.Client{
		Timeout: time.Second,
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			assert.Equal(t, "https://example.com/file.txt", req.URL.String())
			return &http.Response{
				StatusCode: http.StatusOK,
				Body: &errReadCloser{
					data: []byte("new content"),
					err:  io.ErrUnexpectedEOF,
				},
				Header:  make(http.Header),
				Request: req,
			}, nil
		}),
	}
	t.Cleanup(func() {
		downloadHTTPClient = oldClient
	})

	err := updateFile(cacheFile, "https://example.com/file.txt", "")
	require.Error(t, err)

	got, readErr := os.ReadFile(cacheFile)
	require.NoError(t, readErr)
	assert.Equal(t, "cached content", string(got))

	entries, readDirErr := os.ReadDir(filepath.Dir(cacheFile))
	require.NoError(t, readDirErr)
	require.Len(t, entries, 1)
	assert.Equal(t, filepath.Base(cacheFile), entries[0].Name())
}

func TestUpdateFileReplacesExistingFileAfterSuccessfulRefresh(t *testing.T) {
	cacheFile := filepath.Join(t.TempDir(), "cache.txt")
	require.NoError(t, os.WriteFile(cacheFile, []byte("cached content"), 0600))

	oldTime := time.Now().Add(-8 * 24 * time.Hour)
	require.NoError(t, os.Chtimes(cacheFile, oldTime, oldTime))

	oldClient := downloadHTTPClient
	downloadHTTPClient = &http.Client{
		Timeout: time.Second,
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			assert.Equal(t, "https://example.com/file.txt", req.URL.String())
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader("new content")),
				Header:     make(http.Header),
				Request:    req,
			}, nil
		}),
	}
	t.Cleanup(func() {
		downloadHTTPClient = oldClient
	})

	err := updateFile(cacheFile, "https://example.com/file.txt", "")
	require.NoError(t, err)

	got, readErr := os.ReadFile(cacheFile)
	require.NoError(t, readErr)
	assert.Equal(t, "new content", string(got))

	entries, readDirErr := os.ReadDir(filepath.Dir(cacheFile))
	require.NoError(t, readDirErr)
	require.Len(t, entries, 1)
	assert.Equal(t, filepath.Base(cacheFile), entries[0].Name())
}
