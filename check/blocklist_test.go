package check

import (
	"net"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBlockList(t *testing.T) {
	t.Run("given IP present in blocklist.de file then result and no error is returned", func(t *testing.T) {
		setBlockLisFileMock(t, getBlockListFileMock)
		result, err := BlockList(net.ParseIP("66.249.70.34"))
		require.NoError(t, err)
		assert.Equal(t, "blocklist.de", result.Description)
		assert.Equal(t, IsMalicious, result.Type)
		assert.Equal(t, true, result.IpAddrIsMalicious)
	})
}

func setBlockLisFileMock(t *testing.T, fn func() (*os.File, error)) {
	orig := getBlockListFile
	getBlockListFile = fn
	t.Cleanup(func() {
		getBlockListFile = orig
	})
}

func getBlockListFileMock() (*os.File, error) {
	file, err := os.Open("testdata/blocklist.de_all.list")
	if err != nil {
		return nil, err
	}
	return file, err
}
