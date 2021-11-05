package checkip

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"regexp"
	"strconv"
)

// Blocklist holds information about an IP address from blocklist.de database.
type Blocklist struct {
	Attacks int
	Reports int
}

func (b *Blocklist) Name() string { return "blocklist.de" }

// Check fills in Bloclist data for a given IP address. It gets the data from
// http://api.blocklist.de
func (b *Blocklist) Check(ipddr net.IP) error {
	apiurl := fmt.Sprintf("http://api.blocklist.de/api.php?ip=%s&start=1", ipddr)

	resp, err := makeAPIcall(apiurl, map[string]string{}, map[string]string{})
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	number := regexp.MustCompile(`\d+`)
	numbers := number.FindAll(body, 2)
	if b.Attacks, err = strconv.Atoi(string(numbers[0])); err != nil {
		return err
	}
	if b.Reports, err = strconv.Atoi(string(numbers[1])); err != nil {
		return err
	}

	return nil
}

func (b *Blocklist) IsMalicious() bool {
	return b.Attacks > 0 && b.Reports > 0
}
