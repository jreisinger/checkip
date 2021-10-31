package checkip

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

// ET (Emerging Threats) holds information about an IP address from
// rules.emergingthreats.net
type ET struct {
	CompromisedIP bool
	CountIPs      int
}

// Check checks whether the ippaddr is not among compromised IP addresses from
// The Emerging Threats Intelligence feed (ET). I found ET mentioned at
// https://logz.io/blog/open-source-threat-intelligence-feeds/
func (e *ET) Check(ipaddr net.IP) error {
	file := "/var/tmp/et.txt"
	url := "https://rules.emergingthreats.net/blockrules/compromised-ips.txt"

	if err := updateFile(file, url, ""); err != nil {
		return fmt.Errorf("can't update %s from %s: %v", file, url, err)
	}

	if err := e.search(ipaddr, file); err != nil {
		return fmt.Errorf("searching %s in %s: %v", ipaddr, file, err)
	}

	return nil
}

func (e *ET) IsMalicious() bool {
	return e.CompromisedIP
}

// search searches the ippadrr in filename fills in ET data.
func (e *ET) search(ipaddr net.IP, filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}

	s := bufio.NewScanner(file)
	for s.Scan() {
		line := s.Text()
		e.CountIPs++
		if line == ipaddr.String() {
			e.CompromisedIP = true
		}
	}
	if s.Err() != nil {
		return err
	}

	return nil
}
