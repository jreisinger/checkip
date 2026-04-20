package check

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

type cins struct {
	BadGuyIP bool
	CountIPs int
}

var cinsScoreURL = "https://cinsscore.com/list/ci-badguys.txt"

// CinsScore searches ipaddr in https://cinsscore.com/list/ci-badguys.txt.
func CinsScore(ipaddr net.IP) (Check, error) {
	result := Check{
		Description: "cinsscore.com",
		Type:        IsMalicious,
	}

	// file := "/var/tmp/cins.txt"
	file, err := getCachePath("cins.txt")
	if err != nil {
		return result, err
	}

	if err := updateFile(file, cinsScoreURL, ""); err != nil {
		return result, newCheckError(err)
	}

	cins, err := cinsSearch(ipaddr, file)
	if err != nil {
		return result, newCheckError(fmt.Errorf("searching %s in %s: %v", ipaddr, file, err))
	}

	result.IpAddrIsMalicious = cins.BadGuyIP

	return result, nil
}

// cinsSearch searches the ippadrr in filename fills in ET data.
func cinsSearch(ipaddr net.IP, filename string) (cins, error) {
	file, err := os.Open(filename)
	if err != nil {
		return cins{}, err
	}
	defer file.Close()

	var cinsArmy cins
	s := bufio.NewScanner(file)
	for s.Scan() {
		line := s.Text()
		cinsArmy.CountIPs++
		if line == ipaddr.String() {
			cinsArmy.BadGuyIP = true
		}
	}
	if s.Err() != nil {
		return cins{}, s.Err()
	}
	return cinsArmy, nil
}
