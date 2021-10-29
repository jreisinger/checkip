package checkip

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"
)

// OTX holds information from otx.alienvault.com.
type OTX struct {
	PulseInfo struct {
		Count int `json:"count"`
	} `json:"pulse_info"`
}

// Check gets data from https://otx.alienvault.com/api.
func (otx *OTX) Check(ipaddr net.IP) error {
	otxurl := fmt.Sprintf("https://otx.alienvault.com/api/v1/indicators/IPv4/%s/", ipaddr.String())
	baseURL, err := url.Parse(otxurl)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("GET", baseURL.String(), nil)
	if err != nil {
		return err
	}

	client := newHTTPClient(5 * time.Second)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%s", resp.Status)
	}

	if err := json.NewDecoder(resp.Body).Decode(otx); err != nil {
		return err
	}

	return nil
}

// IsOK returns true if the IP address is not considered suspicious.
func (otx *OTX) IsOK() bool {
	return otx.PulseInfo.Count <= 10
}
