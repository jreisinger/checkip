package check

import (
	"encoding/json"
	"fmt"
	"net"
	"sort"
	"strings"
)

type shodan struct {
	Org   string     `json:"org"`
	Data  shodanData `json:"data"`
	OS    string     `json:"os"`
	Ports []int      `json:"ports"`
	Vulns []string   `json:"vulns"`
}

type shodanData []struct {
	Product   string `json:"product"`
	Version   string `json:"version"`
	Port      int    `json:"port"`
	Transport string `json:"transport"` // tcp, udp
}

var shodanUrl = "https://api.shodan.io"

// Shodan gets generic information from api.shodan.io.
func Shodan(ipaddr net.IP) (Check, error) {
	result := Check{
		Description: "shodan.io",
		Type:        InfoAndIsMalicious,
	}

	apiKey, err := getConfigValue("SHODAN_API_KEY")
	if err != nil {
		return result, newCheckError(err)
	}
	if apiKey == "" {
		result.MissingCredentials = "SHODAN_API_KEY"
		return result, nil
	}

	var shodan shodan
	apiURL := fmt.Sprintf("%s/shodan/host/%s?key=%s", shodanUrl, ipaddr, apiKey)
	if err := defaultHttpClient.GetJson(apiURL, map[string]string{}, map[string]string{}, &shodan); err != nil {
		return result, newCheckError(err)
	}
	if len(shodan.Vulns) == 0 {
		shodan.Vulns = make([]string, 0)
	}
	result.IpAddrInfo = shodan

	if len(shodan.Vulns) > 0 {
		result.IpAddrIsMalicious = true
	}

	return result, nil
}

type byPort shodanData

func (x byPort) Len() int           { return len(x) }
func (x byPort) Less(i, j int) bool { return x[i].Port < x[j].Port }
func (x byPort) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

// Info returns interesting information from the check.
func (s shodan) Summary() string {
	var portInfo []string
	sort.Sort(byPort(s.Data))
	for _, d := range s.Data {
		var product string
		if d.Product != "" {
			product = d.Product
		}

		var version string
		if d.Version != "" {
			version = d.Version
		}

		if product == "" && version == "" {
			portInfo = append(portInfo, fmt.Sprintf("%s/%d", d.Transport, d.Port))
		} else {
			ss := nonEmpty(product, version)
			portInfo = append(portInfo, fmt.Sprintf("%s/%d (%s)", d.Transport, d.Port, strings.Join(ss, ", ")))
		}
	}

	return fmt.Sprintf("OS: %s, open: %s, vulns: %s", na(s.OS), strings.Join(portInfo, ", "), na(strings.Join(s.Vulns, ", ")))
}

func (s shodan) Json() ([]byte, error) {
	return json.Marshal(s)
}
