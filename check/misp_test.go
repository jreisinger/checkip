package check

import (
	"fmt"
	"net"
	"net/http"
	"testing"

	"github.com/logrusorgru/aurora"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMisp(t *testing.T) {
	au := aurora.NewAurora(true)
	t.Run("given valid response then result and no error is returned", func(t *testing.T) {
		handlerFn := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.WriteHeader(http.StatusOK)
			rw.Write(loadResponse(t, "misp_response.json"))
		})
		setMispMockConfig(t)
		testUrl := setMockHttpClient(t, handlerFn)
		setMispMockUrl(t, testUrl)

		result, err := Misp(net.ParseIP("154.16.192.230"))
		require.NoError(t, err)
		assert.Equal(t, fmt.Sprintf("<%s>\t\t", au.Yellow("Misp")), result.Description)
		assert.Equal(t, InfoAndIsMalicious, result.Type)
		assert.Equal(t, true, result.IpAddrIsMalicious)
		assert.Equal(t, "* [13 orgc: 2] Something bad happened: [B1 - vpn] express, [B2 - vpn] express", result.IpAddrInfo.Summary())
	})

	t.Run("given non 2xx response then error is returned", func(t *testing.T) {
		handlerFn := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.WriteHeader(http.StatusInternalServerError)
		})
		setMispMockConfig(t)
		testUrl := setMockHttpClient(t, handlerFn)
		setMispMockUrl(t, testUrl)

		_, err := Misp(net.ParseIP("118.25.6.39"))
		require.Error(t, err)
	})
}

// --- test helpers ---

// setMispMockUrl temporarily sets mispURL variable to testUrl.
func setMispMockUrl(t *testing.T, testUrl string) {
	origUrl := mispURL
	mispURL = testUrl
	t.Cleanup(func() {
		mispURL = origUrl
	})
}

// setMispMockConfig sets MISP_API_KEY to a dummy value.
func setMispMockConfig(t *testing.T) {
	setMockConfig(t, func(key string) (string, error) {
		if key == "MISP_KEY" {
			return "123-secret-789", nil
		}
		if key == "MISP_URL" {
			return "http://localhost", nil
		}
		return "", fmt.Errorf("unexpected key %s received", key)
	})
}
