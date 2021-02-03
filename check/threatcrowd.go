package check

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"

	. "github.com/logrusorgru/aurora"
)

// https://github.com/AlienVault-OTX/ApiV2#votes
var votesMeaning = map[int]string{
	-1: fmt.Sprintf("voted %s by most users", Magenta("malicious")),
	0:  "voted malicious/harmless by equal number of users",
	1:  "voted harmless by most users",
}

// ThreatCrowd holds information about an IP address from
// https://www.threatcrowd.org voting.
type ThreatCrowd struct {
	Votes int `json:"votes"`
}

// Do retrieves information about an IP address from the ThreatCrowd API:
// https://www.threatcrowd.org/searchApi/v2/ip/report. If the IP address is
// voted malicious it returns false.
func (t *ThreatCrowd) Do(ipaddr net.IP) (bool, error) {
	baseURL, err := url.Parse("https://www.threatcrowd.org/searchApi/v2/ip/report")
	if err != nil {
		return false, err
	}

	params := url.Values{}
	params.Add("ip", ipaddr.String())
	baseURL.RawQuery = params.Encode()

	req, err := http.NewRequest("GET", baseURL.String(), nil)
	if err != nil {
		return false, err
	}

	client := NewHTTPClient(5 * time.Second)
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("search threatcrowd failed: %s", resp.Status)
	}

	if err := json.NewDecoder(resp.Body).Decode(t); err != nil {
		return false, err
	}

	if t.isNotOK() {
		return false, nil
	}

	return true, nil
}

func (t *ThreatCrowd) isNotOK() bool {
	// https://github.com/AlienVault-OTX/ApiV2#votes
	return t.Votes < 0
}

// Name returns the name of the check.
func (t *ThreatCrowd) Name() string {
	return fmt.Sprint("ThreatCrowd")
}

// String returns the result of the check.
func (t *ThreatCrowd) String() string {

	return fmt.Sprintf("%s", votesMeaning[t.Votes])
}
