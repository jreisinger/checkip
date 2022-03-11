package checks

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/jreisinger/checkip/check"
)

type AS struct {
	Number      int    `json:"-"`
	FirstIP     net.IP `json:"-"`
	LastIP      net.IP `json:"-"`
	Description string `json:"description"`
	CountryCode string `json:"-"`
}

func (a AS) Summary() string {
	return fmt.Sprintf("AS description: %s", check.Na(a.Description))
}

func (a AS) JsonString() (string, error) {
	b, err := json.Marshal(a)
	return string(b), err
}

// CheckAS gets info about autonomous system (AS) of the ipaddr. The data is
// taken from a TSV file ip2asn-combined downloaded from iptoasn.com.
func CheckAS(ipaddr net.IP) (check.Result, error) {
	file := "/var/tmp/ip2asn-combined.tsv"
	url := "https://iptoasn.com/data/ip2asn-combined.tsv.gz"

	if err := check.UpdateFile(file, url, "gz"); err != nil {
		return check.Result{}, check.NewError(err)
	}

	as, err := asSearch(ipaddr, file)
	if err != nil {
		return check.Result{}, check.NewError(fmt.Errorf("searching %s in %s: %v", ipaddr, file, err))
	}

	return check.Result{
		Name: "iptoasn.com",
		Type: check.TypeInfo,
		Info: as,
	}, nil
}

// search the ippadrr in tsvFile and if found fills in AS data.
func asSearch(ipaddr net.IP, tsvFile string) (AS, error) {
	tsv, err := os.Open(tsvFile)
	if err != nil {
		return AS{}, err
	}

	as := AS{}
	s := bufio.NewScanner(tsv)
	for s.Scan() {
		line := s.Text()
		fields := strings.Split(line, "\t")
		as.FirstIP = net.ParseIP(fields[0])
		as.LastIP = net.ParseIP(fields[1])
		if ipIsBetween(ipaddr, as.FirstIP, as.LastIP) {
			as.Number, err = strconv.Atoi(fields[2])
			if err != nil {
				return AS{}, fmt.Errorf("converting string to int: %v", err)
			}
			as.CountryCode = fields[3]
			as.Description = fields[4]
			return as, nil
		}
	}
	if s.Err() != nil {
		return AS{}, err
	}
	return as, nil
}

func ipIsBetween(ipAddr, firstIPAddr, lastIPAddr net.IP) bool {
	if bytes.Compare(ipAddr, firstIPAddr) >= 0 && bytes.Compare(ipAddr, lastIPAddr) <= 0 {
		return true
	}
	return false
}
