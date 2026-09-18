package check

import (
	"net"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGreyNoise(t *testing.T) {
	t.Run("given valid response then result and no error is returned", func(t *testing.T) {
		handlerFn := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.WriteHeader(http.StatusOK)
			rw.Write(loadResponse(t, "greynoise_response.json"))
		})

		testUrl := setMockHttpClient(t, handlerFn)
		setGreynoiseUrl(t, testUrl+"/")

		result, err := GreyNoise(net.ParseIP("1.2.3.4"))
		require.NoError(t, err)
		assert.Equal(t, "greynoise.io", result.Description)
		assert.Equal(t, InfoAndIsMalicious, result.Type)
		assert.Equal(t, true, result.IpAddrIsMalicious)
		assert.Equal(t, "Success", result.IpAddrInfo.Summary())
	})

	t.Run("given 404 response then IP is reported as not observed with no error", func(t *testing.T) {
		handlerFn := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.WriteHeader(http.StatusNotFound)
		})

		testUrl := setMockHttpClient(t, handlerFn)
		setGreynoiseUrl(t, testUrl+"/")

		result, err := GreyNoise(net.ParseIP("1.2.3.4"))
		require.NoError(t, err)
		assert.Equal(t, false, result.IpAddrIsMalicious)
		info, ok := result.IpAddrInfo.(grey)
		require.True(t, ok)
		assert.Equal(t, "IP not observed scanning the internet or contained in RIOT data set.", info.Message)
		assert.Equal(t, "", info.Classification)
	})

	t.Run("given non 2xx, non 404 response then error is returned", func(t *testing.T) {
		handlerFn := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.WriteHeader(http.StatusInternalServerError)
		})

		testUrl := setMockHttpClient(t, handlerFn)
		setGreynoiseUrl(t, testUrl+"/")

		_, err := GreyNoise(net.ParseIP("1.2.3.4"))
		require.Error(t, err)
	})
}

// --- test helpers ---

func setGreynoiseUrl(t *testing.T, testUrl string) {
	url := greynoiseurl
	greynoiseurl = testUrl
	t.Cleanup(func() {
		greynoiseurl = url
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
		assert.Equal(t, "Test message", summary)
	})

	t.Run("Summary method handles empty message", func(t *testing.T) {
		g := grey{
			IP:   "1.2.3.4",
			Riot: false,
		}
		summary := g.Summary()
		assert.Equal(t, "n/a", summary)
	})
}
