package checks

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"

	"github.com/jreisinger/checkip/check"
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
	return fmt.Sprintf("network: %s, SAN: %s", check.Na(v.Data.Attributes.Network), check.Na(strings.Join(v.Data.Attributes.LastHTTPScert.Extensions.SAN, ", ")))
}

func (v virusTotal) JsonString() (string, error) {
	b, err := json.Marshal(v)
	return string(b), err
}

// VirusTotal gets generic information and security reputation about the ippaddr
// from https://www.virustotal.com/api.
func VirusTotal(ipaddr net.IP) (check.Result, error) {
	apiKey, err := check.GetConfigValue("VIRUSTOTAL_API_KEY")
	if err != nil {
		return check.Result{}, check.NewError(err)
	}
	if apiKey == "" {
		return check.Result{}, nil
	}

	// curl --header "x-apikey:$VIRUSTOTAL_API_KEY" https://www.virustotal.com/api/v3/ip_addresses/1.1.1.1
	headers := map[string]string{"x-apikey": apiKey}
	apiUrl := "https://www.virustotal.com/api/v3/ip_addresses/" + ipaddr.String()
	var virusTotal virusTotal
	if err := check.DefaultHttpClient.GetJson(apiUrl, headers, map[string]string{}, &virusTotal); err != nil {
		return check.Result{}, check.NewError(err)
	}

	return check.Result{
		Name:      "virustotal.com",
		Type:      check.TypeInfoSec,
		Info:      virusTotal,
		Malicious: virusTotal.Data.Attributes.Reputation < 0,
	}, nil
}
