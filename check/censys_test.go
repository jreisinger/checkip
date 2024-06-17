package check

import (
	"net"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCensys(t *testing.T) {

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
			rw.Write(loadResponse(t, "censys_response.json"))
		})

		testUrl := setMockHttpClient(t, handlerFn)
		setCensysUrl(t, testUrl)

		result, err := Censys(net.ParseIP("118.25.6.39"))
		require.NoError(t, err)
		assert.Equal(t, "censys.io", result.Description)
		assert.Equal(t, InfoAndIsMalicious, result.Type)
		assert.Equal(t, true, result.IpAddrIsMalicious)
		assert.Equal(t, "MikroTik, RB760iGS, udp/161 (snmp), tcp/2000 (mikrotik_bw), tcp/51922 (ssh)", result.IpAddrInfo.Summary())
	})

	t.Run("given non 2xx response then error is returned", func(t *testing.T) {
		handlerFn := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.WriteHeader(http.StatusInternalServerError)
		})

		testUrl := setMockHttpClient(t, handlerFn)
		setCensysUrl(t, testUrl)

		_, err := Censys(net.ParseIP("118.25.6.39"))
		require.Error(t, err)
	})
}

// --- test helpers ---

func setCensysUrl(t *testing.T, testUrl string) {
	url := censysUrl
	censysUrl = testUrl
	t.Cleanup(func() {
		censysUrl = url
	})
}
