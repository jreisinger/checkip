package checker

import (
	"encoding/json"
	"fmt"
	"github.com/jreisinger/checkip/pkg/check"
	"net"
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

func (v VirusTotal) String() string {
	return fmt.Sprintf("AS onwer: %s, network: %s, SAN: %s", check.Na(v.Data.Attributes.ASowner), check.Na(v.Data.Attributes.Network), check.Na(strings.Join(v.Data.Attributes.LastHTTPScert.Extensions.SAN, ", ")))
}

func (v VirusTotal) JsonString() (string, error) {
	b, err := json.Marshal(v)
	return string(b), err
}

// CheckVirusTotal fills in data about ippaddr from https://www.virustotal.com/api
func CheckVirusTotal(ipaddr net.IP) check.Result {
	apiKey, err := check.GetConfigValue("VIRUSTOTAL_API_KEY")
	if err != nil {
		return check.Result{ResultError: check.NewResultError(err)}
	}

	// curl --header "x-apikey:$VIRUSTOTAL_API_KEY" https://www.virustotal.com/api/v3/ip_addresses/1.1.1.1
	headers := map[string]string{"x-apikey": apiKey}
	apiUrl := "https://www.virustotal.com/api/v3/ip_addresses/" + ipaddr.String()
	var virusTotal VirusTotal
	if err := check.DefaultHttpClient.GetJson(apiUrl, headers, map[string]string{}, &virusTotal); err != nil {
		return check.Result{ResultError: check.NewResultError(err)}
	}

	return check.Result{
		Name:        "virustotal.com",
		Type:        check.TypeInfoSec,
		Data:        virusTotal,
		IsMalicious: virusTotal.Data.Attributes.Reputation < 0,
	}
}
