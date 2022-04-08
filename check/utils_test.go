package check

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/jreisinger/checkip"
	"github.com/stretchr/testify/require"
)

// SetMockConfig helper to replace GetConfigValue function
func SetMockConfig(t *testing.T, fn func(key string) (string, error)) {
	defaultConfig := checkip.GetConfigValue
	checkip.GetConfigValue = fn
	t.Cleanup(func() {
		checkip.GetConfigValue = defaultConfig
	})
}

// SetSetMockHttpClient sets checkip.DefaultHttpClient to httptest handler and returns test url
func SetMockHttpClient(t *testing.T, handlerFn http.HandlerFunc) string {
	server := httptest.NewServer(handlerFn)
	defaultHttpClient := checkip.DefaultHttpClient
	checkip.DefaultHttpClient = checkip.NewHttpClient(server.Client())
	t.Cleanup(func() {
		server.Close()
		checkip.DefaultHttpClient = defaultHttpClient
	})
	return server.URL
}

// LoadResponse loads named file form testdata directory
func LoadResponse(t *testing.T, name string) []byte {
	b, err := ioutil.ReadFile(filepath.Join("testdata", name))
	require.NoError(t, err)
	return b
}
