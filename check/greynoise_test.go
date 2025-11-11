package check

import (
	"net"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGreyNoise(t *testing.T) {
	t.Run("given malicious IP response then malicious result is returned", func(t *testing.T) {
		handlerFn := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.WriteHeader(http.StatusOK)
			rw.Write(loadResponse(t, "greynoise_malicious_response.json"))
		})
		testUrl := setMockHttpClient(t, handlerFn)
		setGreyNoiseMockUrl(t, testUrl)

		result, err := GreyNoise(net.ParseIP("1.2.3.4"))
		require.NoError(t, err)
		assert.Equal(t, "greynoise.io", result.Description)
		assert.Equal(t, InfoAndIsMalicious, result.Type)
		assert.Equal(t, true, result.IpAddrIsMalicious)
		assert.Contains(t, result.IpAddrInfo.Summary(), "IP 1.2.3.4")
		assert.Contains(t, result.IpAddrInfo.Summary(), "riot: false")
		assert.Contains(t, result.IpAddrInfo.Summary(), "This IP has been observed scanning")
	})

	t.Run("given benign IP response then non-malicious result is returned", func(t *testing.T) {
		handlerFn := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.WriteHeader(http.StatusOK)
			rw.Write(loadResponse(t, "greynoise_benign_response.json"))
		})
		testUrl := setMockHttpClient(t, handlerFn)
		setGreyNoiseMockUrl(t, testUrl)

		result, err := GreyNoise(net.ParseIP("8.8.8.8"))
		require.NoError(t, err)
		assert.Equal(t, "greynoise.io", result.Description)
		assert.Equal(t, InfoAndIsMalicious, result.Type)
		assert.Equal(t, false, result.IpAddrIsMalicious)
		assert.Contains(t, result.IpAddrInfo.Summary(), "IP 8.8.8.8")
		assert.Contains(t, result.IpAddrInfo.Summary(), "riot: true")
		assert.Contains(t, result.IpAddrInfo.Summary(), "known service provider")
	})

	t.Run("given 404 response then no error is returned with default message", func(t *testing.T) {
		handlerFn := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.WriteHeader(http.StatusNotFound)
			rw.Write(loadResponse(t, "greynoise_404_response.json"))
		})
		testUrl := setMockHttpClient(t, handlerFn)
		setGreyNoiseMockUrl(t, testUrl)

		result, err := GreyNoise(net.ParseIP("192.168.1.1"))
		require.NoError(t, err)
		assert.Equal(t, "greynoise.io", result.Description)
		assert.Equal(t, InfoAndIsMalicious, result.Type)
		assert.Equal(t, false, result.IpAddrIsMalicious)
		assert.Contains(t, result.IpAddrInfo.Summary(), "IP 192.168.1.1")
		assert.Contains(t, result.IpAddrInfo.Summary(), "riot: false")
		assert.Contains(t, result.IpAddrInfo.Summary(), "IP not observed scanning")
	})

	t.Run("given non-404 error response then error is returned", func(t *testing.T) {
		handlerFn := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.WriteHeader(http.StatusInternalServerError)
		})
		testUrl := setMockHttpClient(t, handlerFn)
		setGreyNoiseMockUrl(t, testUrl)

		_, err := GreyNoise(net.ParseIP("1.2.3.4"))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "500 Internal Server Error")
	})

	t.Run("given network error then error is returned", func(t *testing.T) {
		setGreyNoiseMockUrl(t, "http://invalid-url-that-does-not-exist.local")

		_, err := GreyNoise(net.ParseIP("1.2.3.4"))
		require.Error(t, err)
	})
}

// Test helper functions

// setGreyNoiseMockUrl temporarily sets greynoiseurl variable to testUrl.
func setGreyNoiseMockUrl(t *testing.T, testUrl string) {
	origUrl := greynoiseurl
	greynoiseurl = testUrl + "/"
	t.Cleanup(func() {
		greynoiseurl = origUrl
	})
}

// Test the grey struct methods
func TestGreyStruct(t *testing.T) {
	t.Run("Summary method formats correctly", func(t *testing.T) {
		g := grey{
			IP:      "1.2.3.4",
			Riot:    true,
			Message: "Test message",
		}
		summary := g.Summary()
		assert.Equal(t, "IP 1.2.3.4, riot: true, message: Test message", summary)
	})

	t.Run("Summary method handles empty message", func(t *testing.T) {
		g := grey{
			IP:   "1.2.3.4",
			Riot: false,
		}
		summary := g.Summary()
		assert.Equal(t, "IP 1.2.3.4, riot: false, message: n/a", summary)
	})

	t.Run("Json method marshals correctly", func(t *testing.T) {
		g := grey{
			IP:             "1.2.3.4",
			Noise:          true,
			Riot:           false,
			Classification: "malicious",
		}
		jsonBytes, err := g.Json()
		require.NoError(t, err)
		assert.Contains(t, string(jsonBytes), `"ip":"1.2.3.4"`)
		assert.Contains(t, string(jsonBytes), `"noise":true`)
		assert.Contains(t, string(jsonBytes), `"riot":false`)
		assert.Contains(t, string(jsonBytes), `"classification":"malicious"`)
	})
}
