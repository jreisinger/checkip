package check

import (
	"net"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNetAPI(t *testing.T) {
	t.Run("given IP present in the NetAPI list then it is malicious", func(t *testing.T) {
		setNetAPIFileMock(t)
		result, err := NetAPI(net.ParseIP("198.51.100.23"))
		require.NoError(t, err)
		assert.Equal(t, "netapi.com", result.Description)
		assert.Equal(t, IsMalicious, result.Type)
		assert.Equal(t, true, result.IpAddrIsMalicious)
	})

	t.Run("given IP not in the NetAPI list then it is not malicious", func(t *testing.T) {
		setNetAPIFileMock(t)
		result, err := NetAPI(net.ParseIP("192.0.2.1"))
		require.NoError(t, err)
		assert.Equal(t, false, result.IpAddrIsMalicious)
	})
}

func setNetAPIFileMock(t *testing.T) {
	orig := getNetAPIFile
	getNetAPIFile = func() (*os.File, error) {
		return os.Open("testdata/netapi_ip_live.txt")
	}
	t.Cleanup(func() {
		getNetAPIFile = orig
	})
}
