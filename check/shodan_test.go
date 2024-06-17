package check

import (
	"net"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShodan(t *testing.T) {
	apiKey, err := getConfigValue("SHODAN_API_KEY")
	if err != nil || apiKey == "" {
		return
	}

	t.Run("given valid response then result and no error is returned", func(t *testing.T) {
		handlerFn := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.WriteHeader(http.StatusOK)
			rw.Write(loadResponse(t, "shodan_response.json"))
		})

		testUrl := setMockHttpClient(t, handlerFn)
		setShodanUrl(t, testUrl)

		result, err := Shodan(net.ParseIP("118.25.6.39"))
		require.NoError(t, err)
		assert.Equal(t, "shodan.io", result.Description)
		assert.Equal(t, InfoAndIsMalicious, result.Type)
		assert.Equal(t, true, result.IpAddrIsMalicious)
	})

	t.Run("given non 2xx response then error is returned", func(t *testing.T) {
		handlerFn := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.WriteHeader(http.StatusInternalServerError)
		})

		testUrl := setMockHttpClient(t, handlerFn)
		setShodanUrl(t, testUrl)

		_, err := Shodan(net.ParseIP("118.25.6.39"))
		require.Error(t, err)
	})
}

// --- test helpers ---

func setShodanUrl(t *testing.T, testUrl string) {
	url := shodanUrl
	shodanUrl = testUrl
	t.Cleanup(func() {
		shodanUrl = url
	})
}
