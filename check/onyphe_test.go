package check

import (
	"net"
	"net/http"
	"testing"

	"github.com/jreisinger/checkip"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOnyphe(t *testing.T) {

	apiKey, err := getConfigValue("CENSYS_KEY")
	if err != nil || apiKey == "" {
		return
	}
	apiSec, err := getConfigValue("CENSYS_SEC")
	if err != nil || apiSec == "" {
		return
	}

	t.Run("given valid response then result and no error is returned", func(t *testing.T) {
		handlerFn := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.WriteHeader(http.StatusOK)
			rw.Write(loadResponse(t, "onyphe_response.json"))
		})

		testUrl := setMockHttpClient(t, handlerFn)
		setOnypheUrl(t, testUrl)

		result, err := Onyphe(net.ParseIP("118.25.6.39"))
		require.NoError(t, err)
		assert.Equal(t, "onyphe.io", result.Name)
		assert.Equal(t, checkip.TypeInfoSec, result.Type)
		assert.Equal(t, true, result.Malicious)
		assert.Equal(t, "Open: snmp udp/161 (RouterOS, Mikrotik), winbox tcp/8291 (Linux Kernel, Linux)", result.Info.Summary())
	})


	t.Run("given non 2xx response then error is returned", func(t *testing.T) {
		handlerFn := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.WriteHeader(http.StatusInternalServerError)
		})

		testUrl := setMockHttpClient(t, handlerFn)
		setOnypheUrl(t, testUrl)

		_, err := Onyphe(net.ParseIP("118.25.6.39"))
		require.Error(t, err)
	})
}

// --- test helpers ---

func setOnypheUrl(t *testing.T, testUrl string) {
	url := onypheUrl
	onypheUrl = testUrl
	t.Cleanup(func() {
		onypheUrl = url
	})
}
