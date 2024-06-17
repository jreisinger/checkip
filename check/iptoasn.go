package check

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

type autonomousSystem struct {
	Number      int    `json:"-"`
	FirstIP     net.IP `json:"-"`
	LastIP      net.IP `json:"-"`
	Description string `json:"description"`
	CountryCode string `json:"-"`
}

func (a autonomousSystem) Summary() string {
	return a.Description
}

func (a autonomousSystem) Json() ([]byte, error) {
	return json.Marshal(a)
}

// IPtoASN gets info about autonomous system of the ipaddr. The data is taken
// from https://iptoasn.com/data/ip2asn-combined.tsv.gz.
func IPtoASN(ipaddr net.IP) (Check, error) {
	result := Check{
		Description: "iptoasn.com",
		Type:        Info,
	}

	// file := "/var/tmp/ip2asn-combined.tsv"
	file, err := getCachePath("ip2asn-combined.tsv")
	if err != nil {
		return result, err
	}

	url := "https://iptoasn.com/data/ip2asn-combined.tsv.gz"

	if err := updateFile(file, url, "gz"); err != nil {
		return result, newCheckError(err)
	}

	as, err := asSearch(ipaddr, file)
	if err != nil {
		return result, newCheckError(fmt.Errorf("searching %s in %s: %v", ipaddr, file, err))
	}
	result.IpAddrInfo = as

	return result, nil
}

// search the ippadrr in tsvFile and if found fills in AS data.
func asSearch(ipaddr net.IP, tsvFile string) (autonomousSystem, error) {
	tsv, err := os.Open(tsvFile)
	if err != nil {
		return autonomousSystem{}, err
	}

	as := autonomousSystem{}
	s := bufio.NewScanner(tsv)
	for s.Scan() {
		line := s.Text()
		fields := strings.Split(line, "\t")
		as.FirstIP = net.ParseIP(fields[0])
		as.LastIP = net.ParseIP(fields[1])
		if ipIsBetween(ipaddr, as.FirstIP, as.LastIP) {
			as.Number, err = strconv.Atoi(fields[2])
			if err != nil {
				return autonomousSystem{}, fmt.Errorf("converting string to int: %v", err)
			}
			as.CountryCode = fields[3]
			as.Description = fields[4]
			return as, nil
		}
	}
	if s.Err() != nil {
		return autonomousSystem{}, err
	}
	return as, nil
}

func ipIsBetween(ipAddr, firstIPAddr, lastIPAddr net.IP) bool {
	if bytes.Compare(ipAddr, firstIPAddr) >= 0 && bytes.Compare(ipAddr, lastIPAddr) <= 0 {
		return true
	}
	return false
}
