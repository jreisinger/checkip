package checker

import (
	"bufio"
	"fmt"
	"github.com/jreisinger/checkip/pkg/check"
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

// CheckCins fills in the CINSArmy data.
func CheckCins(ipaddr net.IP) check.Result {
	file := "/var/tmp/cins.txt"
	url := "http://cinsscore.com/list/ci-badguys.txt"

	if err := check.UpdateFile(file, url, ""); err != nil {
		return check.Result{ResultError: check.NewResultError(err)}
	}

	cins, err := cinsSearch(ipaddr, file)
	if err != nil {
		return check.Result{ResultError: check.NewResultError(fmt.Errorf("searching %s in %s: %v", ipaddr, file, err))}
	}

	return check.Result{
		Name:        "cinsscore.com",
		Type:        check.TypeSec,
		Data:        check.EmptyData{},
		IsMalicious: cins.BadGuyIP,
	}
}

// search searches the ippadrr in filename fills in ET data.
func cinsSearch(ipaddr net.IP, filename string) (CINSArmy, error) {
	file, err := os.Open(filename)
	if err != nil {
		return CINSArmy{}, err
	}

	var cinsArmy CINSArmy
	s := bufio.NewScanner(file)
	for s.Scan() {
		line := s.Text()
		cinsArmy.CountIPs++
		if line == ipaddr.String() {
			cinsArmy.BadGuyIP = true
		}
	}
	if s.Err() != nil {
		return CINSArmy{}, err
	}
	return cinsArmy, nil
}