package checker

import (
	"fmt"
	"net"
	"regexp"
	"strconv"

	"github.com/jreisinger/checkip/pkg/check"
)

// CheckBlockList fills in BlockList data for a given IP address. It gets the data from
// http://api.blocklist.de
func CheckBlockList(ipddr net.IP) check.Result {
	url := fmt.Sprintf("http://api.blocklist.de/api.php?ip=%s&start=1", ipddr)

	resp, err := check.DefaultHttpClient.Get(url, map[string]string{}, map[string]string{})
	if err != nil {
		return check.Result{Error: check.NewResultError(err)}
	}

	number := regexp.MustCompile(`\d+`)
	numbers := number.FindAll(resp, 2)

	attacks, err := strconv.Atoi(string(numbers[0]))
	if err != nil {
		return check.Result{Error: check.NewResultError(err)}
	}
	reports, err := strconv.Atoi(string(numbers[1]))
	if err != nil {
		return check.Result{Error: check.NewResultError(err)}
	}

	return check.Result{
		CheckName:         "blocklist.de",
		CheckType:         check.TypeInfoSec,
		Data:              check.EmptyData{},
		IsIPaddrMalicious: attacks > 0 && reports > 0,
	}
}
