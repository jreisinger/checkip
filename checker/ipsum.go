package checker

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/jreisinger/checkip/check"
)

// CheckIPSum checks how many blacklists the IP address is found on.
func CheckIPSum(ipaddr net.IP) check.Result {
	file := "/var/tmp/ipsum.txt"
	url := "https://raw.githubusercontent.com/stamparm/ipsum/master/ipsum.txt"

	if err := check.UpdateFile(file, url, ""); err != nil {
		return check.Result{Error: check.NewResultError(err)}
	}

	blackLists, err := searchIPSumBlacklists(ipaddr, file)
	if err != nil {
		return check.Result{Error: check.NewResultError(fmt.Errorf("searching %s in %s: %v", ipaddr, file, err))}
	}

	return check.Result{
		CheckName:         "github.com/stamparm/ipsum",
		CheckType:         check.TypeSec,
		Data:              check.EmptyData{},
		IsIPaddrMalicious: blackLists > 0,
	}
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
