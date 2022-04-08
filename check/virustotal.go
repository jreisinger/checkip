package check

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"

	"github.com/jreisinger/checkip"
)

type virusTotal struct {
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

func (v virusTotal) Summary() string {
	return fmt.Sprintf("network: %s, SAN: %s", na(v.Data.Attributes.Network), na(strings.Join(v.Data.Attributes.LastHTTPScert.Extensions.SAN, ", ")))
}

func (v virusTotal) JsonString() (string, error) {
	b, err := json.Marshal(v)
	return string(b), err
}

// VirusTotal gets generic information and security reputation about the ippaddr
// from https://www.virustotal.com/api.
func VirusTotal(ipaddr net.IP) (checkip.Result, error) {
	result := checkip.Result{
		Name: "virustotal.com",
		Type: checkip.TypeInfoSec,
	}

	apiKey, err := getConfigValue("VIRUSTOTAL_API_KEY")
	if err != nil {
		return result, newCheckError(err)
	}
	if apiKey == "" {
		return result, nil
	}

	// curl --header "x-apikey:$VIRUSTOTAL_API_KEY" https://www.virustotal.com/api/v3/ip_addresses/1.1.1.1
	headers := map[string]string{"x-apikey": apiKey}
	apiUrl := "https://www.virustotal.com/api/v3/ip_addresses/" + ipaddr.String()
	var virusTotal virusTotal
	if err := defaultHttpClient.GetJson(apiUrl, headers, map[string]string{}, &virusTotal); err != nil {
		return result, newCheckError(err)
	}

	result.Info = virusTotal
	result.Malicious = virusTotal.Data.Attributes.Reputation < 0

	return result, nil
}
