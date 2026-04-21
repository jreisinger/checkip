package check

import (
	"fmt"
	"net"
	"net/http"
	"testing"


	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMyDB(t *testing.T) {
	t.Run("given valid response then result and no error is returned", func(t *testing.T) {
		handlerFn := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.WriteHeader(http.StatusOK)
			rw.Write(loadResponse(t, "mydb_response.json"))
		})
		setMyDBMockConfig(t)
		testUrl := setMockHttpClient(t, handlerFn)
		setMyDBMockUrl(t, testUrl)

		result, err := MyDB(net.ParseIP("154.16.192.230"))
		require.NoError(t, err)
		assert.Equal(t, "MyDB", result.Description)
		assert.Equal(t, InfoAndIsMalicious, result.Type)
		assert.Equal(t, true, result.IpAddrIsMalicious)
		assert.Equal(t, "Registred on 2025-01-20T19:55:02+01:00 ** Survey: 154.16.192.0/24 -- vpn PIA", result.IpAddrInfo.Summary())
	})

	t.Run("given non 2xx response then error is returned", func(t *testing.T) {
		handlerFn := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.WriteHeader(http.StatusInternalServerError)
		})
		setMyDBMockConfig(t)
		testUrl := setMockHttpClient(t, handlerFn)
		setMyDBMockUrl(t, testUrl)

		_, err := MyDB(net.ParseIP("118.25.6.39"))
		require.Error(t, err)
	})
}

// --- test helpers ---

// setMyDBMockUrl temporarily sets myDBUrl variable to testUrl.
func setMyDBMockUrl(t *testing.T, testUrl string) {
	origUrl := myDBUrl
	myDBUrl = testUrl
	t.Cleanup(func() {
		myDBUrl = origUrl
	})
}

// setMyDBMockConfig sets MYDB_API_KEY to a dummy value.
func setMyDBMockConfig(t *testing.T) {
	setMockConfig(t, func(key string) (string, error) {
		if key == "MYDB_API_KEY" {
			return "123-secret-789", nil
		}
		if key == "MYDB_URL" {
			return "http://localhost", nil
		}
		return "", fmt.Errorf("unexpected key %s received", key)
	})
}
