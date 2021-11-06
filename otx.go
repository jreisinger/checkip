package checkip

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
)

// OTX holds information from otx.alienvault.com.
type OTX struct {
	PulseInfo struct {
		Count int `json:"count"`
	} `json:"pulse_info"`
}

func (otx *OTX) String() string { return "otx.alienvault.com" }

// Check gets data from https://otx.alienvault.com/api.
func (otx *OTX) Check(ipaddr net.IP) error {
	apiurl := fmt.Sprintf("https://otx.alienvault.com/api/v1/indicators/IPv4/%s/", ipaddr.String())

	resp, err := makeAPIcall(apiurl, map[string]string{}, map[string]string{})
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

func (otx *OTX) IsMalicious() bool {
	return otx.PulseInfo.Count > 10
}
