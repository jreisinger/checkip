package check

import (
	"fmt"
	"net"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAbuseIPDB(t *testing.T) {
	t.Run("given valid response then result and no error is returned", func(t *testing.T) {
		handlerFn := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.WriteHeader(http.StatusOK)
			rw.Write(loadResponse(t, "abuseipdb_response.json"))
		})
		setAbuseIPDBMockConfig(t)
		testUrl := setMockHttpClient(t, handlerFn)
		setAbuseIPDBMockUrl(t, testUrl)

		result, err := AbuseIPDB(net.ParseIP("118.25.6.39"))
		require.NoError(t, err)
		assert.Equal(t, "abuseipdb.com", result.Description)
		assert.Equal(t, InfoAndIsMalicious, result.Type)
		assert.Equal(t, true, result.IpAddrIsMalicious)
		assert.Equal(t, "domain: tencent.com, usage type: Data Center/Web Hosting/Transit", result.IpAddrInfo.Summary())
	})

	t.Run("given non 2xx response then error is returned", func(t *testing.T) {
		handlerFn := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.WriteHeader(http.StatusInternalServerError)
		})
		setAbuseIPDBMockConfig(t)
		testUrl := setMockHttpClient(t, handlerFn)
		setAbuseIPDBMockUrl(t, testUrl)

		_, err := AbuseIPDB(net.ParseIP("118.25.6.39"))
		require.Error(t, err)
	})
}

// --- test helpers ---

// setAbuseIPDBMockUrl temporarily sets abuseIPDBUrl variable to testUrl.
func setAbuseIPDBMockUrl(t *testing.T, testUrl string) {
	origUrl := abuseIPDBUrl
	abuseIPDBUrl = testUrl
	t.Cleanup(func() {
		abuseIPDBUrl = origUrl
	})
}

// setAbuseIPDBMockConfig sets ABUSEIPDB_API_KEY to a dummy value.
func setAbuseIPDBMockConfig(t *testing.T) {
	setMockConfig(t, func(key string) (string, error) {
		if key == "ABUSEIPDB_API_KEY" {
			return "123-secret-789", nil
		}
		return "", fmt.Errorf("unexpected key %s received", key)
	})
}
