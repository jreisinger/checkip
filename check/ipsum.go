package check

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/jreisinger/checkip"
)

// IPSum checks how many blacklists the ipaddr is found on. It uses
// https://raw.githubusercontent.com/stamparm/ipsum/master/ipsum.txt.
func IPSum(ipaddr net.IP) (checkip.Result, error) {
	result := checkip.Result{
		Name: "github.com/stamparm/ipsum",
		Type: checkip.TypeSec,
	}

	// file := "/var/tmp/ipsum.txt"
	file, err := getDbFilesPath("ipsum.txt")
	if err != nil {
		return result, err
	}

	url := "https://raw.githubusercontent.com/stamparm/ipsum/master/ipsum.txt"

	if err := updateFile(file, url, ""); err != nil {
		return result, newCheckError(err)
	}

	blackLists, err := searchIPSumBlacklists(ipaddr, file)
	if err != nil {
		return result, newCheckError(fmt.Errorf("searching %s in %s: %v", ipaddr, file, err))
	}
	result.Malicious = blackLists > 0

	return result, nil
}

// searchIPSumBlacklists searches the ippadrr in tsvFile for number of blacklists
func searchIPSumBlacklists(ipaddr net.IP, tsvFile string) (int, error) {
	tsv, err := os.Open(tsvFile)
	if err != nil {
		return 0, err
	}

	s := bufio.NewScanner(tsv)
	for s.Scan() {
		line := s.Text()
		if strings.HasPrefix(line, "#") { // skip comments
			continue
		}
		fields := strings.Fields(line)
		if ipaddr.Equal(net.ParseIP(fields[0])) { // IP address found
			numOfBlacklists, err := strconv.Atoi(fields[1])
			if err != nil {
				return 0, err
			}
			return numOfBlacklists, nil
		}
	}
	if s.Err() != nil {
		return 0, err
	}
	return 0, nil
}
