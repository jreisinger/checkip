package checkip

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"
)

// ThreatCrowd holds information about an IP address from threatcrowd.org
// voting.
type ThreatCrowd struct {
	Votes int `json:"votes"`
}

// Check retrieves information about an IP address from the ThreatCrowd API:
// https://www.threatcrowd.org/searchApi/v2/ip/report. It returns false if the
// IP address is voted malicious by most users.
func (t *ThreatCrowd) Check(ipaddr net.IP) (bool, error) {
	baseURL, err := url.Parse("https://www.threatcrowd.org/searchApi/v2/ip/report")
	if err != nil {
		return true, err
	}

	params := url.Values{}
	params.Add("ip", ipaddr.String())
	baseURL.RawQuery = params.Encode()

	req, err := http.NewRequest("GET", baseURL.String(), nil)
	if err != nil {
		return true, err
	}

	client := newHTTPClient(5 * time.Second)
	resp, err := client.Do(req)
	if err != nil {
		return true, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return true, fmt.Errorf("search threatcrowd failed: %s", resp.Status)
	}

	if err := json.NewDecoder(resp.Body).Decode(t); err != nil {
		return true, err
	}

	return t.isOK(), nil
}

func (t *ThreatCrowd) isOK() bool {
	// https://github.com/AlienVault-OTX/ApiV2#votes
	return t.Votes >= 0
}

// String returns the result of the check.
func (t *ThreatCrowd) String() string {
	// https://github.com/AlienVault-OTX/ApiV2#votes
	votesMeaning := map[int]string{
		-1: "voted malicious by most users",
		0:  "voted malicious/harmless by equal number of users",
		1:  "voted harmless by most users",
	}

	return votesMeaning[t.Votes]
}
