package checks

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/jreisinger/checkip/check"
)

// CheckIPSum checks how many blacklists the ipaddr is found on.
func CheckIPSum(ipaddr net.IP) (check.Result, error) {
	file := "/var/tmp/ipsum.txt"
	url := "https://raw.githubusercontent.com/stamparm/ipsum/master/ipsum.txt"

	if err := check.UpdateFile(file, url, ""); err != nil {
		return check.Result{}, check.NewError(err)
	}

	blackLists, err := searchIPSumBlacklists(ipaddr, file)
	if err != nil {
		return check.Result{}, check.NewError(fmt.Errorf("searching %s in %s: %v", ipaddr, file, err))
	}

	return check.Result{
		Name:            "github.com/stamparm/ipsum",
		Type:            check.TypeSec,
		Info:            check.EmptyInfo{},
		IPaddrMalicious: blackLists > 0,
	}, nil
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
