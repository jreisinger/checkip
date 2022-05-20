package check

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// loadResponse loads named file form testdata directory.
func loadResponse(t *testing.T, name string) []byte {
	b, err := ioutil.ReadFile(filepath.Join("testdata", name))
	require.NoError(t, err)
	return b
}

// setMockHttpClient temporarily sets defaultHttpClient to handlerFn and returns
// test URL.
func setMockHttpClient(t *testing.T, handlerFn http.HandlerFunc) string {
	server := httptest.NewServer(handlerFn)
	dHC := defaultHttpClient
	defaultHttpClient = newHttpClient(server.Client())
	t.Cleanup(func() {
		server.Close()
		defaultHttpClient = dHC
	})
	return server.URL
}

// setMockConfig temporarily replaces getConfigValue function with fn.
func setMockConfig(t *testing.T, fn func(key string) (string, error)) {
	origGetConfigValue := getConfigValue
	getConfigValue = fn
	t.Cleanup(func() {
		getConfigValue = origGetConfigValue
	})
}
