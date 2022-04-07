package checks

import (
	"fmt"
	"net"
	"regexp"
	"strconv"

	"github.com/jreisinger/checkip/check"
)

// BlockList searches the ipaddr in http://api.blocklist.de.
func BlockList(ipddr net.IP) (check.Result, error) {
	result := check.Result{
		Name: "blocklist.de",
		Type: check.TypeSec,
		Info: check.EmptyInfo{},
	}

	url := fmt.Sprintf("http://api.blocklist.de/api.php?ip=%s&start=1", ipddr)

	resp, err := check.DefaultHttpClient.Get(url, map[string]string{}, map[string]string{})
	if err != nil {
		return result, check.NewError(err)
	}

	number := regexp.MustCompile(`\d+`)
	numbers := number.FindAll(resp, 2)

	attacks, err := strconv.Atoi(string(numbers[0]))
	if err != nil {
		return result, check.NewError(err)
	}
	reports, err := strconv.Atoi(string(numbers[1]))
	if err != nil {
		return result, check.NewError(err)
	}

	result.Malicious = attacks > 0 && reports > 0

	return result, nil
}
