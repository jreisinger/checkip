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
			LastAnalysisStats struct {
				Harmless   int `json:"harmless"`
				Malicious  int `json:"malicious"`
				Suspicious int `json:"suspicious"`
				Timeout    int `json:"timeout"`
				Undetected int `json:"undetected"`
			} `json:"last_analysis_stats"`
			TotalVotes struct {
				Harmless  int
				Malicious int
			} `json:"total_votes"`
		} `json:"attributes"`
	} `json:"data"`
}

// Check fills in data for a given IP address from virustotal API. It returns
// false if the IP address is considered malicious or suspicious by some
// analysis. See https://developers.virustotal.com/v3.0/reference#ip-object
func (vt *VirusTotal) Check(ipaddr net.IP) (bool, error) {
	apiKey, err := getConfigValue("VIRUSTOTAL_API_KEY")
	if err != nil {
		return false, fmt.Errorf("can't call API: %w", err)
	}

	// curl --header "x-apikey:$VIRUSTOTAL_API_KEY" https://www.virustotal.com/api/v3/ip_addresses/1.1.1.1

	baseURL, err := url.Parse("https://www.virustotal.com/api/v3/ip_addresses/" + ipaddr.String())
	if err != nil {
		return false, err
	}

	req, err := http.NewRequest("GET", baseURL.String(), nil)
	if err != nil {
		return false, err
	}

	// Set request headers.
	req.Header.Set("x-apikey", apiKey)

	client := newHTTPClient(5 * time.Second)
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("%s", resp.Status)
	}

	if err := json.NewDecoder(resp.Body).Decode(vt); err != nil {
		return false, err
	}

	return vt.isOK(), nil
}

func (vt *VirusTotal) isOK() bool {
	return vt.Data.Attributes.LastAnalysisStats.Malicious <= 0
}

// String returns the result of the check.
func (vt *VirusTotal) String() string {
	return fmt.Sprintf("%d harmless, %d suspicious, %d malicious analysis results",
		vt.Data.Attributes.LastAnalysisStats.Harmless,
		vt.Data.Attributes.LastAnalysisStats.Suspicious,
		vt.Data.Attributes.LastAnalysisStats.Malicious,
	)
}
