package checkip

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

// ET (Emerging Threats) says whether the IP address was found among CountIPs
// compromised IP addresses according to rules.emergingthreats.net
type ET struct {
	CompromisedIP bool
	CountIPs      int
}

// Check checks whether the ippaddr is not among compromised IP addresses from
// The Emerging Threats Intelligence feed (ET).
// https://logz.io/blog/open-source-threat-intelligence-feeds/
func (e *ET) Check(ipaddr net.IP) (bool, error) {
	file := "/var/tmp/et.txt"
	url := "https://rules.emergingthreats.net/blockrules/compromised-ips.txt"

	if err := update(file, url, ""); err != nil {
		return true, fmt.Errorf("can't update %s from %s: %v", file, url, err)
	}

	if err := e.search(ipaddr, file); err != nil {
		return true, fmt.Errorf("searching %s in %s: %v", ipaddr, file, err)
	}

	return true, nil
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

// String returns the result of the check.
func (e *ET) String() string {
	s := fmt.Sprintf("found among %d compromised IP addresses", e.CountIPs)
	if !e.CompromisedIP {
		s = "not " + s
	}
	return s
}
