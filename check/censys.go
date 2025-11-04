package check

import (
	"encoding/json"
	"fmt"
	"net"
	"sort"
	"strings"
)

type censys struct {
	Result result `json:"result"`
}

type result struct {
	Ressource ressource `json:"resource"`
}

type ressource struct {
	Data censysData      `json:"services"`
	OS   operatingSystem `json:"operating_system"`
}

type operatingSystem struct {
	Product string `json:"product"`
	Vendor  string `json:"vendor"`
	Version string `json:"version"`
	Edition string `json:"edition"`
}

type censysData []struct {
	Port        int    `json:"port"`
	Transport   string `json:"transport_protocol"` // tcp, udp
	ServiceName string `json:"protocol"`
}

var censysUrl = "https://api.platform.censys.io/v3/global/asset/host"

// Censys gets generic information from search.censys.io.
func Censys(ipaddr net.IP) (Check, error) {
	result := Check{
		Description: "censys.io",
		Type:        InfoAndIsMalicious,
	}

	headers := map[string]string{
		"Accept":       "application/vnd.censys.api.v3.host.v1+json",
		"Content-Type": "application/x-www-form-urlencoded;charset=UTF-8",
	}

	// mandatory CENSYS_KEY (token in v3)
	apiKey, err := getConfigValue("CENSYS_KEY")
	if err != nil {
		return result, newCheckError(err)
	}
	if apiKey == "" {
		result.MissingCredentials = "CENSYS_KEY"
		return result, nil
	}
	headers["Authorization"] = "Bearer " + apiKey

	// optional CENSYS_ORG_ID for starter and entreprise plans
	apiOrgID, _ := getConfigValue("CENSYS_ORG_ID")
	if apiOrgID != "" {
		headers["X-Organization-ID"] = apiOrgID
	}

	var censys censys
	apiURL := fmt.Sprintf("%s/%s", censysUrl, ipaddr)
	if err := defaultHttpClient.GetJson(apiURL, headers, map[string]string{}, &censys); err != nil {
		return result, newCheckError(err)
	}
	result.IpAddrInfo = censys

	for _, d := range censys.Result.Ressource.Data {
		port := d.Port
		if port != 80 && port != 443 && port != 53 { // undecidable ports
			result.IpAddrIsMalicious = true
		}
	}

	return result, nil
}

type byPortC censysData

func (x byPortC) Len() int           { return len(x) }
func (x byPortC) Less(i, j int) bool { return x[i].Port < x[j].Port }
func (x byPortC) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

// Summary returns interesting information from the check.
func (c censys) Summary() string {
	var portInfo []string
	sort.Sort(byPortC(c.Result.Ressource.Data))
	for _, d := range c.Result.Ressource.Data {
		service := make(map[string]int)
		sport := fmt.Sprintf("%s/%d", strings.ToLower(d.Transport), d.Port)
		service[sport]++

		if service[sport] > 1 {
			continue
		}

		portInfo = append(portInfo, fmt.Sprintf("%s (%s)", sport, strings.ToLower(d.ServiceName)))
	}

	s := c.Result.Ressource
	return fmt.Sprintf("OS: %s %s, open: %s", na(s.OS.Vendor), na(s.OS.Product), strings.Join(portInfo, ", "))
}

func (c censys) Json() ([]byte, error) {
	return json.Marshal(c)
}
