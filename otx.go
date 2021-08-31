package checkip

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"
)

// OTX holds IP address reputation data from otx.alienvault.com.
type OTX struct {
	PulseInfo struct {
		Count int `json:"count"`
	} `json:"pulse_info"`
}

// Check gets data from https://otx.alienvault.com/api. It returns false if
// there are more than 10 pulses about the IP address.
func (otx *OTX) Check(ipaddr net.IP) (bool, error) {
	otxurl := fmt.Sprintf("https://otx.alienvault.com/api/v1/indicators/IPv4/%s/", ipaddr.String())
	baseURL, err := url.Parse(otxurl)
	if err != nil {
		return true, err
	}

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
		return true, fmt.Errorf("%s", resp.Status)
	}

	if err := json.NewDecoder(resp.Body).Decode(otx); err != nil {
		return true, err
	}

	return otx.isOK(), nil
}

func (otx *OTX) isOK() bool {
	return otx.PulseInfo.Count <= 10
}

// String returns the result of the check.
func (otx *OTX) String() string {
	return fmt.Sprintf("%d pulses", otx.PulseInfo.Count)
}
