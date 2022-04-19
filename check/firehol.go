package check

import (
	"io"
	"net"
	"os"
	"strings"

	"github.com/jreisinger/checkip"
)

type firehol struct {
	blacklisted bool
}

// Firehol checks whether the ipaddr is found on a blacklist from
// https://iplists.firehol.org/?ipset=firehol_level1.
func Firehol(ipaddr net.IP) (checkip.Result, error) {
	result := checkip.Result{
		Name: "firehol.org",
		Type: checkip.TypeSec,
	}

	file, err := getDbFilesPath("firehol_level1.netset")
	if err != nil {
		return result, err
	}

	url := "https://iplists.firehol.org/files/firehol_level1.netset"

	if err := updateFile(file, url, ""); err != nil {
		return result, newCheckError(err)
	}

	f, err := os.Open(file)
	if err != nil {
		return result, err
	}
	defer f.Close()

	found, err := ipFound(f, ipaddr)
	if err != nil {
		return result, err
	}
	result.Malicious = found

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
		// skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		_, ipnet, err := net.ParseCIDR(line)
		if err != nil {
			return false, err
		}

		if ipnet.Contains(ipaddr) {
			return true, nil
		}
	}

	return false, nil
}
