package asn

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/jreisinger/checkip/util"
)

// AS holds information about an Autonomous System.
type AS struct {
	FirstIP     net.IP
	LastIP      net.IP
	CountryCode string
	Number      int
	Description string
}

// New creates AS.
func New() *AS {
	return &AS{}
}

// ForIP fills in AS data for a given IP address.
func (a *AS) ForIP(ipaddr net.IP) error {
	file := "/var/tmp/ip2asn-combined.tsv"
	url := "https://iptoasn.com/data/ip2asn-combined.tsv.gz"
	if err := update(file, url); err != nil {
		return fmt.Errorf("can't update %s from %s: %v", file, url, err)
	}

	if err := a.search(ipaddr, file); err != nil {
		return fmt.Errorf("searching %s in %s: %v", ipaddr, file, err)
	}

	return nil
}

// search searches the ippadrr in tsvFile and if found fills in AS data.
func (a *AS) search(ipaddr net.IP, tsvFile string) error {
	tsv, err := os.Open(tsvFile)
	if err != nil {
		return err
	}

	s := bufio.NewScanner(tsv)
	for s.Scan() {
		line := s.Text()
		fields := strings.Split(line, "\t")
		a.FirstIP = net.ParseIP(fields[0])
		a.LastIP = net.ParseIP(fields[1])
		if isBetween(ipaddr, a.FirstIP, a.LastIP) {
			a.Number, err = strconv.Atoi(fields[2])
			if err != nil {
				return fmt.Errorf("converting string to int: %v", err)
			}
			a.CountryCode = fields[3]
			a.Description = fields[4]
			return nil
		}
	}
	if s.Err() != nil {
		return err
	}

	return nil
}

func isBetween(ipAddr, firstIPAddr, lastIPAddr net.IP) bool {
	if bytes.Compare(ipAddr, firstIPAddr) >= 0 && bytes.Compare(ipAddr, lastIPAddr) <= 0 {
		return true
	}
	return false
}

// update downloads and creates file from url if not present, updates if file is
// older than a week.
func update(filepath, url string) error {
	file, err := os.Stat(filepath)

	if os.IsNotExist(err) {
		r, err := util.DownloadFile(url)
		if err != nil {
			return err
		}
		if err := util.ExtractGzFile(filepath, r); err != nil {
			return err
		}

		return nil // don't check ModTime if file does not exist
	}

	if util.IsOlderThanOneWeek(file.ModTime()) {
		r, err := util.DownloadFile(url)
		if err != nil {
			return err
		}
		if err := util.ExtractGzFile(filepath, r); err != nil {
			return err
		}
	}

	return nil
}
