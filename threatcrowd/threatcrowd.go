package threatcrowd

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
)

// ThreatCrowd holds information about an IP address from https://www.threatcrowd.org voting.
type ThreatCrowd struct {
	ResponseCode string `json:"response_code"`
	Resolutions  []struct {
		LastResolved string `json:"last_resolved"`
		Domain       string `json:"domain"`
	} `json:"resolutions"`
	Hashes     []string      `json:"hashes"`
	References []interface{} `json:"references"`
	Votes      int           `json:"votes"`
	Permalink  string        `json:"permalink"`
}

// New creates ThreatCrowd.
func New() *ThreatCrowd {
	return &ThreatCrowd{}
}

// ForIP fills in ThreatCrowd data for a given IP address.
func (t *ThreatCrowd) ForIP(ipaddr net.IP) error {
	// curl https://www.threatcrowd.org/searchApi/v2/ip/report/?ip=188.40.75.132

	baseURL, err := url.Parse("https://www.threatcrowd.org/searchApi/v2/ip/report")
	if err != nil {
		return err
	}

	params := url.Values{}
	params.Add("ip", ipaddr.String())
	baseURL.RawQuery = params.Encode()

	req, err := http.NewRequest("GET", baseURL.String(), nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("search threatcrowd failed: %s", resp.Status)
	}

	if err := json.NewDecoder(resp.Body).Decode(t); err != nil {
		return err
	}

	return nil
}
