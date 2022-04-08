package check

import (
	"fmt"
	"net"
	"regexp"
	"strconv"

	"github.com/jreisinger/checkip"
)

// BlockList searches the ipaddr in http://api.blocklist.de.
func BlockList(ipddr net.IP) (checkip.Result, error) {
	result := checkip.Result{
		Name: "blocklist.de",
		Type: checkip.TypeSec,
		Info: checkip.EmptyInfo{},
	}

	url := fmt.Sprintf("http://api.blocklist.de/api.php?ip=%s&start=1", ipddr)

	resp, err := defaultHttpClient.Get(url, map[string]string{}, map[string]string{})
	if err != nil {
		return result, newCheckError(err)
	}

	number := regexp.MustCompile(`\d+`)
	numbers := number.FindAll(resp, 2)

	attacks, err := strconv.Atoi(string(numbers[0]))
	if err != nil {
		return result, newCheckError(err)
	}
	reports, err := strconv.Atoi(string(numbers[1]))
	if err != nil {
		return result, newCheckError(err)
	}

	result.Malicious = attacks > 0 && reports > 0

	return result, nil
}
