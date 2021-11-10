package checkip

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
)

// VirusTotal holds information about an IP address from virustotal.com.
type VirusTotal struct {
	Data struct {
		Attributes struct {
			Reputation        int    `json:"reputation"`
			Network           string `json:"network"`
			ASowner           string `json:"as_owner"`
			LastAnalysisStats struct {
				Harmless   int `json:"harmless"`
				Malicious  int `json:"malicious"`
				Suspicious int `json:"suspicious"`
				Timeout    int `json:"timeout"`
				Undetected int `json:"undetected"`
			} `json:"last_analysis_stats"`
			LastHTTPScert struct {
				Extensions struct {
					SAN []string `json:"subject_alternative_name"`
				} `json:"extensions"`
			} `json:"last_https_certificate"`
		} `json:"attributes"`
	} `json:"data"`
}

func (vt *VirusTotal) String() string { return "virustotal.com" }

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
	apiurl := "https://www.virustotal.com/api/v3/ip_addresses/" + ipaddr.String()
	resp, err := makeAPIcall(apiurl, headers, map[string]string{})
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("calling %s: %s", apiurl, resp.Status)
	}

	if err := json.NewDecoder(resp.Body).Decode(vt); err != nil {
		return err
	}

	return nil
}

func (vt *VirusTotal) IsMalicious() bool {
	// https://developers.virustotal.com/reference#ip-object
	return vt.Data.Attributes.Reputation < 0
}

func (vt *VirusTotal) Info() string {
	return fmt.Sprintf("AS onwer: %s, network: %s, SAN: %s", na(vt.Data.Attributes.ASowner), na(vt.Data.Attributes.Network), na(strings.Join(vt.Data.Attributes.LastHTTPScert.Extensions.SAN, ", ")))
}
