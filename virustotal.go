package checkip

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"
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

// Check fills in data for a given IP address from virustotal API. It returns
// false if the IP address is considered malicious or suspicious by some
// analysis. See https://developers.virustotal.com/v3.0/reference#ip-object
func (vt *VirusTotal) Check(ipaddr net.IP) (bool, error) {
	apiKey, err := getConfigValue("VIRUSTOTAL_API_KEY")
	if err != nil {
		return true, fmt.Errorf("can't call API: %w", err)
	}

	// curl --header "x-apikey:$VIRUSTOTAL_API_KEY" https://www.virustotal.com/api/v3/ip_addresses/1.1.1.1

	baseURL, err := url.Parse("https://www.virustotal.com/api/v3/ip_addresses/" + ipaddr.String())
	if err != nil {
		return true, err
	}

	req, err := http.NewRequest("GET", baseURL.String(), nil)
	if err != nil {
		return true, err
	}

	// Set request headers.
	req.Header.Set("x-apikey", apiKey)

	client := newHTTPClient(5 * time.Second)
	resp, err := client.Do(req)
	if err != nil {
		return true, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return true, fmt.Errorf("%s", resp.Status)
	}

	if err := json.NewDecoder(resp.Body).Decode(vt); err != nil {
		return true, err
	}

	return vt.isOK(), nil
}

func (vt *VirusTotal) isOK() bool {
	// https://developers.virustotal.com/reference#ip-object
	return vt.Data.Attributes.Reputation >= 0
}

// String returns the result of the check.
func (vt *VirusTotal) String() string {
	what := "malicious"
	if vt.Data.Attributes.Reputation >= 0 {
		what = "harmless"
	}
	return fmt.Sprintf("%s with reputation of %d (the higher the absolute number, the more trustworthy)", what, vt.Data.Attributes.Reputation)
}
