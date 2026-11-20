package check

import (
	"net"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVPNDetection(t *testing.T) {
	t.Run("given VPN response then it is malicious and the key is sent", func(t *testing.T) {
		handlerFn := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			assert.Equal(t, "Bearer 123-secret-789", req.Header.Get("Authorization"))
			assert.Equal(t, "/185.159.157.20", req.URL.Path)
			rw.WriteHeader(http.StatusOK)
			rw.Write(loadResponse(t, "vpndetection_vpn_response.json"))
		})
		setVPNDetectionMockConfig(t, "123-secret-789")
		testUrl := setMockHttpClient(t, handlerFn)
		setVPNDetectionMockUrl(t, testUrl)

		result, err := VPNDetection(net.ParseIP("185.159.157.20"))
		require.NoError(t, err)
		assert.Equal(t, "vpndetection.io", result.Description)
		assert.Equal(t, InfoAndIsMalicious, result.Type)
		assert.Equal(t, true, result.IpAddrIsMalicious)
		assert.Equal(t, "VPN", result.IpAddrInfo.Summary())
	})

	t.Run("given relay response then it is not malicious", func(t *testing.T) {
		handlerFn := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.WriteHeader(http.StatusOK)
			rw.Write(loadResponse(t, "vpndetection_relay_response.json"))
		})
		setVPNDetectionMockConfig(t, "123-secret-789")
		testUrl := setMockHttpClient(t, handlerFn)
		setVPNDetectionMockUrl(t, testUrl)

		result, err := VPNDetection(net.ParseIP("172.224.226.6"))
		require.NoError(t, err)
		assert.Equal(t, false, result.IpAddrIsMalicious)
		assert.Equal(t, "relay: Apple iCloud Private Relay", result.IpAddrInfo.Summary())
	})

	t.Run("given no API key then the check is not run", func(t *testing.T) {
		handlerFn := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			t.Error("no request expected without an API key")
		})
		setVPNDetectionMockConfig(t, "")
		testUrl := setMockHttpClient(t, handlerFn)
		setVPNDetectionMockUrl(t, testUrl)

		result, err := VPNDetection(net.ParseIP("185.159.157.20"))
		require.NoError(t, err)
		assert.Equal(t, "VPNDETECTION_API_KEY", result.MissingCredentials)
		assert.Nil(t, result.IpAddrInfo)
	})

	t.Run("given non 2xx response then error is returned", func(t *testing.T) {
		handlerFn := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.WriteHeader(http.StatusInternalServerError)
		})
		setVPNDetectionMockConfig(t, "123-secret-789")
		testUrl := setMockHttpClient(t, handlerFn)
		setVPNDetectionMockUrl(t, testUrl)

		_, err := VPNDetection(net.ParseIP("185.159.157.20"))
		require.Error(t, err)
	})

	t.Run("given response without ip then error is returned", func(t *testing.T) {
		handlerFn := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.WriteHeader(http.StatusOK)
			rw.Write([]byte(`{"error":"bad request"}`))
		})
		setVPNDetectionMockConfig(t, "123-secret-789")
		testUrl := setMockHttpClient(t, handlerFn)
		setVPNDetectionMockUrl(t, testUrl)

		_, err := VPNDetection(net.ParseIP("185.159.157.20"))
		require.Error(t, err)
	})
}

// --- test helpers ---

func setVPNDetectionMockUrl(t *testing.T, testUrl string) {
	origUrl := vpnDetectionUrl
	vpnDetectionUrl = testUrl
	t.Cleanup(func() {
		vpnDetectionUrl = origUrl
	})
}

func setVPNDetectionMockConfig(t *testing.T, apiKey string) {
	setMockConfig(t, func(key string) (string, error) {
		if key == "VPNDETECTION_API_KEY" {
			return apiKey, nil
		}
		return "", nil
	})
}
