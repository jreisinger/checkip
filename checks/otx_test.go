package checks

import (
	"net"
	"net/http"
	"testing"

	"github.com/jreisinger/checkip/check"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCheckOTX(t *testing.T) {
	t.Run("given valid response then result and no error is returned", func(t *testing.T) {
		handlerFn := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.WriteHeader(http.StatusOK)
			rw.Write(LoadResponse(t, "otx_response.json"))
		})

		testUrl := SetMockHttpClient(t, handlerFn)
		setOTXUrl(t, testUrl)

		result, err := CheckOTX(net.ParseIP("118.25.6.39"))
		require.NoError(t, err)
		assert.Equal(t, "otx.alienvault.com", result.Name)
		assert.Equal(t, check.TypeSec, result.Type)
		assert.Equal(t, true, result.Malicious)
		assert.Equal(t, check.EmptyInfo{}, result.Info)
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
