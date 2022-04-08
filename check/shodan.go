package check

import (
	"encoding/json"
	"fmt"
	"net"
	"sort"
	"strings"

	"github.com/jreisinger/checkip"
)

type shodan struct {
	Org   string     `json:"org"`
	Data  shodanData `json:"data"`
	OS    string     `json:"os"`
	Ports []int      `json:"ports"`
}

type shodanData []struct {
	Product   string `json:"product"`
	Version   string `json:"version"`
	Port      int    `json:"port"`
	Transport string `json:"transport"` // tcp, udp
}

// Shodan gets generic information from https://api.shodan.io.
func Shodan(ipaddr net.IP) (checkip.Result, error) {
	result := checkip.Result{
		Name: "shodan.io",
		Type: checkip.TypeInfo,
	}

	apiKey, err := checkip.GetConfigValue("SHODAN_API_KEY")
	if err != nil {
		return result, checkip.NewError(err)
	}
	if apiKey == "" {
		return result, nil
	}

	var shodan shodan
	apiURL := fmt.Sprintf("https://api.shodan.io/shodan/host/%s?key=%s", ipaddr, apiKey)
	if err := checkip.DefaultHttpClient.GetJson(apiURL, map[string]string{}, map[string]string{}, &shodan); err != nil {
		return result, checkip.NewError(err)
	}
	result.Info = shodan

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
			ss := checkip.NonEmpty(product, version)
			portInfo = append(portInfo, fmt.Sprintf("%s/%d (%s)", d.Transport, d.Port, strings.Join(ss, ", ")))
		}
	}

	portStr := "port"
	if len(portInfo) != 1 {
		portStr += "s"
	}
	if len(portInfo) > 0 {
		portStr += ":"
	}
	return fmt.Sprintf("OS: %s, %d open %s %s", checkip.Na(s.OS), len(portInfo), portStr, strings.Join(portInfo, ", "))
}

func (s shodan) JsonString() (string, error) {
	b, err := json.Marshal(s)
	return string(b), err
}
