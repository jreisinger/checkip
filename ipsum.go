package checkip

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

// IPsum holds information from github.com/stamparm/ipsum.
type IPsum struct {
	NumOfBlacklists int
}

func (ip *IPsum) Name() string { return "github.com/stamparm/ipsum" }

// Check checks how many blackists the IP address is found on.
func (ip *IPsum) Check(ipaddr net.IP) error {
	file := "/var/tmp/ipsum.txt"
	url := "https://raw.githubusercontent.com/stamparm/ipsum/master/ipsum.txt"

	if err := updateFile(file, url, ""); err != nil {
		return fmt.Errorf("can't update %s from %s: %v", file, url, err)
	}

	if err := ip.search(ipaddr, file); err != nil {
		return fmt.Errorf("searching %s in %s: %v", ipaddr, file, err)
	}

	return nil
}

func (ip *IPsum) IsMalicious() bool {
	return ip.NumOfBlacklists > 0
}

// search searches the ippadrr in tsvFile and if found fills in IPsum data.
func (ip *IPsum) search(ipaddr net.IP, tsvFile string) error {
	tsv, err := os.Open(tsvFile)
	if err != nil {
		return err
	}

	s := bufio.NewScanner(tsv)
	for s.Scan() {
		line := s.Text()
		if strings.HasPrefix(line, "#") { // skip comments
			continue
		}
		fields := strings.Fields(line)
		if ipaddr.Equal(net.ParseIP(fields[0])) { // IP address found
			ip.NumOfBlacklists, err = strconv.Atoi(fields[1])
			if err != nil {
				return err
			}
			break
		}
	}
	if s.Err() != nil {
		return err
	}

	return nil
}
