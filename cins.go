package checkip

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

// CINSArmy says whether the IP address was found among CountIPs "bad guys" IP
// addresses according to https://cinsscore.com/#list
type CINSArmy struct {
	BadGuyIP bool
	CountIPs int
}

// Check fills in the CINSArmy data.
// https://logz.io/blog/open-source-threat-intelligence-feeds/
func (c *CINSArmy) Check(ipaddr net.IP) (bool, error) {
	file := "/var/tmp/cins.txt"
	url := "http://cinsscore.com/list/ci-badguys.txt"

	if err := update(file, url, ""); err != nil {
		return true, fmt.Errorf("can't update %s from %s: %v", file, url, err)
	}

	if err := c.search(ipaddr, file); err != nil {
		return true, fmt.Errorf("searching %s in %s: %v", ipaddr, file, err)
	}

	return c.isOK(), nil
}

func (c *CINSArmy) isOK() bool {
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

// String returns the result of the check.
func (c *CINSArmy) String() string {
	s := fmt.Sprintf("found among %d \"bad guy\" IP addresses", c.CountIPs)
	if !c.BadGuyIP {
		s = "not " + s
	}
	return s
}
