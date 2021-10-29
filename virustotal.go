package checkip

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
)

// VirusTotal holds information about an IP address from virustotal.com.
type VirusTotal struct {
	Data struct {
		Attributes struct {
			Reputation        int `json:"reputation"`
			LastAnalysisStats struct {
				Harmless   int `json:"harmless"`
				Malicious  int `json:"malicious"`
				Suspicious int `json:"suspicious"`
				Timeout    int `json:"timeout"`
				Undetected int `json:"undetected"`
			} `json:"last_analysis_stats"`
		} `json:"attributes"`
	} `json:"data"`
}

// Check fills in data about ippaddr from https://www.virustotal.com/api
func (vt *VirusTotal) Check(ipaddr net.IP) error {
	apiKey, err := getConfigValue("VIRUSTOTAL_API_KEY")
	if err != nil {
		return fmt.Errorf("can't call API: %w", err)
	}

	// curl --header "x-apikey:$VIRUSTOTAL_API_KEY" https://www.virustotal.com/api/v3/ip_addresses/1.1.1.1

	headers := map[string]string{
		"x-apikey": apiKey,
	}
	resp, err := makeAPIcall("https://www.virustotal.com/api/v3/ip_addresses/"+ipaddr.String(), headers, map[string]string{})
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%s", resp.Status)
	}

	if err := json.NewDecoder(resp.Body).Decode(vt); err != nil {
		return err
	}

	return nil
}

// IsOK returns true if the IP address is not considered suspicious.
func (vt *VirusTotal) IsOK() bool {
	// https://developers.virustotal.com/reference#ip-object
	return vt.Data.Attributes.Reputation >= 0
}
