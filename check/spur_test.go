package check

import (
	"net"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSpur(t *testing.T) {
	apiKey, err := getConfigValue("SPUR_API_KEY")
	if err != nil || apiKey == "" {
		return
	}

	t.Run("given valid response then result and no error is returned", func(t *testing.T) {
		handlerFn := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.WriteHeader(http.StatusOK)
			rw.Write(loadResponse(t, "spur_response.json"))
		})

		testUrl := setMockHttpClient(t, handlerFn)
		setSpurUrl(t, testUrl)

		result, err := Spur(net.ParseIP("148.72.164.177"))
		require.NoError(t, err)
		assert.Equal(t, "spur.io", result.Description)
		assert.Equal(t, InfoAndIsMalicious, result.Type)
		assert.Equal(t, true, result.IpAddrIsMalicious)
		assert.Equal(t, "VPN: NORD_VPN", result.IpAddrInfo.Summary())
	})

	t.Run("given non 2xx response then error is returned", func(t *testing.T) {
		handlerFn := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.WriteHeader(http.StatusInternalServerError)
		})

		testUrl := setMockHttpClient(t, handlerFn)
		setSpurUrl(t, testUrl)

		_, err := Spur(net.ParseIP("148.72.164.177"))
		require.Error(t, err)
	})
}

// --- test helpers ---

func setSpurUrl(t *testing.T, testUrl string) {
	url := spurUrl
	spurUrl = testUrl
	t.Cleanup(func() {
		spurUrl = url
	})
}
