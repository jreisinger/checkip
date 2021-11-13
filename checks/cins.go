package checks

import (
	"bufio"
	"fmt"
	"net"
	"os"

	"github.com/jreisinger/checkip/check"
)

type cins struct {
	BadGuyIP bool
	CountIPs int
}

// CheckCins searches ipaddr in https://cinsscore.com/list/ci-badguys.txt.
func CheckCins(ipaddr net.IP) (check.Result, *check.Error) {
	file := "/var/tmp/cins.txt"
	url := "http://cinsscore.com/list/ci-badguys.txt"

	if err := check.UpdateFile(file, url, ""); err != nil {
		return check.Result{}, check.NewError(err)
	}

	cins, err := cinsSearch(ipaddr, file)
	if err != nil {
		return check.Result{}, check.NewError(fmt.Errorf("searching %s in %s: %v", ipaddr, file, err))
	}

	return check.Result{
		Name:            "cinsscore.com",
		Type:            check.TypeSec,
		Info:            check.EmptyInfo{},
		IPaddrMalicious: cins.BadGuyIP,
	}, nil
}

// search searches the ippadrr in filename fills in ET data.
func cinsSearch(ipaddr net.IP, filename string) (cins, error) {
	file, err := os.Open(filename)
	if err != nil {
		return cins{}, err
	}

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
		return cins{}, err
	}
	return cinsArmy, nil
}
