package check

import (
	"bufio"
	"net"
	"os"
	"strings"
)

// BlockList searches the ipaddr in lists.blocklist.de/lists/dnsbl/all.list.
func BlockList(ipaddr net.IP) (Check, error) {
	result := Check{
		Description: "blocklist.de",
		Type:        IsMalicious,
	}

	file, err := getBlockListFile()
	if err != nil {
		return result, err
	}
	defer file.Close()

	input := bufio.NewScanner(file)
	for input.Scan() {
		fields := strings.Split(input.Text(), ":")
		if net.ParseIP(fields[0]).Equal(ipaddr) {
			result.IpAddrIsMalicious = true
			break
		}
	}
	if err := input.Err(); err != nil {
		return result, err
	}

	return result, nil
}

// getBlockListFile downloads (if outdated) and returns open file containing
// blocklist.de database.
var getBlockListFile = func() (*os.File, error) {
	file, err := getCachePath("blocklist.de_all.list")
	if err != nil {
		return nil, err
	}

	u := "https://lists.blocklist.de/lists/dnsbl/all.list"
	if err := updateFile(file, u, ""); err != nil {
		return nil, newCheckError(err)
	}

	return os.Open(file)
}
