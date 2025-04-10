package check

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"
)

type spur struct {
	As struct {
		Number       int    `json:"number"`
		Organization string `json:"organization"`
	} `json:"as"`
	Client struct {
		Behaviors []string `json:"behaviors"`
		Count     int      `json:"count"`
		Proxies   []string `json:"proxies"`
		Types     []string `json:"types"`
	} `json:"client"`
	Infrastructure string `json:"infrastructure"`
	IP             string `json:"ip"`
	Location       struct {
		City    string `json:"city"`
		Country string `json:"country"`
		State   string `json:"state"`
	} `json:"location"`
	Organization string   `json:"organization"`
	Risks        []string `json:"risks"`
	Services     []string `json:"services"`
	Tunnels      []struct {
		Anonymous bool     `json:"anonymous"`
		Entries   []string `json:"entries"`
		Operator  string   `json:"operator"`
		Type      string   `json:"type"`
	} `json:"tunnels"`
}

var spurUrl = "https://api.spur.us/v2/context"

// Spur gets generic information from api.spur.io.
func Spur(ipaddr net.IP) (Check, error) {
	result := Check{
		Description: "spur.io",
		Type:        InfoAndIsMalicious,
	}

	apiKey, err := getConfigValue("SPUR_API_KEY")
	if err != nil {
		return result, newCheckError(err)
	}
	if apiKey == "" {
		result.MissingCredentials = "SPUR_API_KEY"
		return result, nil
	}

	headers := map[string]string{
		"Token":        apiKey,
		"Accept":       "application/json",
		"Content-Type": "application/json",
	}

	var spur spur
	apiURL := fmt.Sprintf("%s/%s", spurUrl, ipaddr)
	if err := defaultHttpClient.GetJson(apiURL, headers, map[string]string{}, &spur); err != nil {
		return result, newCheckError(err)
	}
	result.IpAddrInfo = spur
	for _, t := range spur.Tunnels {
		if t.Anonymous == true {
			result.IpAddrIsMalicious = true
		}
	}

	return result, nil
}

// Info returns interesting information from the check.
func (s spur) Summary() string {
	var operators []string
	var stype []string
	for _, t := range s.Tunnels {
		if t.Anonymous == true {
			stype = append(stype, t.Type)
		}
		if t.Operator != "" {
			operators = append(operators, t.Operator)
		}
	}
	if len(s.Tunnels) == 0 {
		stype = append(stype, "Residential")
		operators = s.Risks
	}

	return fmt.Sprintf("%s: %s", strings.Join(stype, ", "), strings.Join(operators, ", "))
}

func (s spur) Json() ([]byte, error) {
	return json.Marshal(s)
}
