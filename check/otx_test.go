package check

import (
	"net"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOTX(t *testing.T) {
	t.Run("given valid response then result and no error is returned", func(t *testing.T) {
		handlerFn := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.WriteHeader(http.StatusOK)
			rw.Write(loadResponse(t, "otx_response.json"))
		})

		testUrl := setMockHttpClient(t, handlerFn)
		setOTXUrl(t, testUrl)

		result, err := OTX(net.ParseIP("118.25.6.39"))
		require.NoError(t, err)
		assert.Equal(t, "otx.alienvault.com", result.Description)
		assert.Equal(t, IsMalicious, result.Type)
		assert.Equal(t, true, result.IpAddrIsMalicious)
	})

	t.Run("given non 2xx response then error is returned", func(t *testing.T) {
		handlerFn := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.WriteHeader(http.StatusInternalServerError)
		})

		testUrl := setMockHttpClient(t, handlerFn)
		setOTXUrl(t, testUrl)

		_, err := OTX(net.ParseIP("118.25.6.39"))
		require.Error(t, err)
	})
}

// --- test helpers ---

func setOTXUrl(t *testing.T, testUrl string) {
	url := otxUrl
	otxUrl = testUrl
	t.Cleanup(func() {
		otxUrl = url
	})
}
