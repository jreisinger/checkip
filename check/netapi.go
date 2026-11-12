package check

import (
	"bufio"
	"net"
	"os"
	"strings"
)

// NetAPI searches the ipaddr in https://netapi.com/zf-api/f/compromised/ip/live/,
// the IPv4 addresses that attacked NetAPI's own honeypot sensors within the
// last 48 hours (one address per line). See https://netapi.com/compromised-ips/.
func NetAPI(ipaddr net.IP) (Check, error) {
	result := Check{
		Description: "netapi.com",
		Type:        IsMalicious,
	}

	file, err := getNetAPIFile()
	if err != nil {
		return result, err
	}
	defer file.Close()

	input := bufio.NewScanner(file)
	for input.Scan() {
		line := strings.TrimSpace(input.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if net.ParseIP(line).Equal(ipaddr) {
			result.IpAddrIsMalicious = true
			break
		}
	}
	if err := input.Err(); err != nil {
		return result, newCheckError(err)
	}

	return result, nil
}

// getNetAPIFile downloads (if outdated) and returns open file containing the
// NetAPI list of honeypot attackers.
var getNetAPIFile = func() (*os.File, error) {
	file, err := getCachePath("netapi.com_ip_live.txt")
	if err != nil {
		return nil, newCheckError(err)
	}

	u := "https://netapi.com/zf-api/f/compromised/ip/live/"
	if err := updateFile(file, u, ""); err != nil {
		return nil, newCheckError(err)
	}

	return os.Open(file)
}
