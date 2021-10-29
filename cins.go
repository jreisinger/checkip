package checkip

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

// CINSArmy holds information about an IP address from
// https://cinsscore.com/#list. I found CINSArmy mentioned at
// https://logz.io/blog/open-source-threat-intelligence-feeds/.
type CINSArmy struct {
	BadGuyIP bool
	CountIPs int
}

// Check fills in the CINSArmy data.
func (c *CINSArmy) Check(ipaddr net.IP) error {
	file := "/var/tmp/cins.txt"
	url := "http://cinsscore.com/list/ci-badguys.txt"

	if err := updateFile(file, url, ""); err != nil {
		return fmt.Errorf("can't update %s from %s: %v", file, url, err)
	}

	if err := c.search(ipaddr, file); err != nil {
		return fmt.Errorf("searching %s in %s: %v", ipaddr, file, err)
	}

	return nil
}

// IsOK returns true if the IP address is not considered suspicious.
func (c *CINSArmy) IsOK() bool {
	return !c.BadGuyIP
}

// search searches the ippadrr in filename fills in ET data.
func (c *CINSArmy) search(ipaddr net.IP, filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}

	s := bufio.NewScanner(file)
	for s.Scan() {
		line := s.Text()
		c.CountIPs++
		if line == ipaddr.String() {
			c.BadGuyIP = true
		}
	}
	if s.Err() != nil {
		return err
	}

	return nil
}
