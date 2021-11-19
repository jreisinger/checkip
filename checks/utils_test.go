package checks

import (
	"github.com/jreisinger/checkip/check"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
)

// SetMockConfig helper to replace GetConfigValue function
func SetMockConfig(t *testing.T, fn func(key string) (string, error)) {
	defaultConfig := check.GetConfigValue
	check.GetConfigValue = fn
	t.Cleanup(func() {
		check.GetConfigValue = defaultConfig
	})
}

// SetSetMockHttpClient sets check.DefaultHttpClient to httptest handler and returns test url
func SetMockHttpClient(t *testing.T, handlerFn http.HandlerFunc) string {
	server := httptest.NewServer(handlerFn)
	defaultHttpClient := check.DefaultHttpClient
	check.DefaultHttpClient = check.NewHttpClient(server.Client())
	t.Cleanup(func() {
		server.Close()
		check.DefaultHttpClient = defaultHttpClient
	})
	return server.URL
}

// LoadResponse loads named file form testdata directory
func LoadResponse(t *testing.T, name string) []byte {
	b, err := ioutil.ReadFile(filepath.Join("testdata", name))
	require.NoError(t, err)
	return b
}
