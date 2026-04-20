package check

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

var fireholGetCachePath = getCachePath
var fireholUpdateFile = updateFile
var fireholSearch = ipFound

// Firehol checks whether the ipaddr is found on blacklist
// https://iplists.firehol.org/?ipset=firehol_level1.
func Firehol(ipaddr net.IP) (Check, error) {
	result := Check{
		Description: "firehol.org",
		Type:        IsMalicious,
	}

	file, err := fireholGetCachePath("firehol_level1.netset")
	if err != nil {
		return result, newCheckError(err)
	}

	url := "https://iplists.firehol.org/files/firehol_level1.netset"

	if err := fireholUpdateFile(file, url, ""); err != nil {
		return result, newCheckError(err)
	}

	f, err := os.Open(file)
	if err != nil {
		return result, newCheckError(err)
	}
	defer f.Close()

	found, err := fireholSearch(f, ipaddr)
	if err != nil {
		return result, newCheckError(err)
	}
	result.IpAddrIsMalicious = found

	return result, nil
}

// ipFound says whether ippaddr was found in subnets. Subnets contains subnets
// in CIDR notation, one per line. Empty lines and comment lines are ignored.
func ipFound(subnets io.Reader, ipaddr net.IP) (bool, error) {
	lines, err := io.ReadAll(subnets)
	if err != nil {
		return false, err
	}

	for _, line := range strings.Split(string(lines), "\n") {
		line = strings.TrimSpace(line)

		// skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		ipnet, err := parseSubnet(line)
		if err != nil {
			return false, fmt.Errorf("parse FireHOL entry %q: %w", line, err)
		}

		if ipnet.Contains(ipaddr) {
			return true, nil
		}
	}

	return false, nil
}

func parseSubnet(line string) (*net.IPNet, error) {
	if ip := net.ParseIP(line); ip != nil {
		if ipv4 := ip.To4(); ipv4 != nil {
			return &net.IPNet{IP: ipv4, Mask: net.CIDRMask(32, 32)}, nil
		}
		return &net.IPNet{IP: ip, Mask: net.CIDRMask(128, 128)}, nil
	}

	_, ipnet, err := net.ParseCIDR(line)
	if err != nil {
		return nil, err
	}
	return ipnet, nil
}
