package checker

import (
	"encoding/json"
	"fmt"
	"net"
	"sort"
	"strings"

	"github.com/jreisinger/checkip/check"
)

// Shodan holds information about an IP address from shodan.io scan data.
type Shodan struct {
	Org   string     `json:"org"`
	Data  ShodanData `json:"data"`
	OS    string     `json:"os"`
	Ports []int      `json:"ports"`
}

type ShodanData []struct {
	Product   string `json:"product"`
	Version   string `json:"version"`
	Port      int    `json:"port"`
	Transport string `json:"transport"` // tcp, udp
}

// CheckShodan gets data from https://api.shodan.io.
func CheckShodan(ipaddr net.IP) (check.Result, error) {
	apiKey, err := check.GetConfigValue("SHODAN_API_KEY")
	if err != nil {
		return check.Result{}, check.NewError(err)
	}

	var shodan Shodan
	apiURL := fmt.Sprintf("https://api.shodan.io/shodan/host/%s?key=%s", ipaddr, apiKey)
	if err := check.DefaultHttpClient.GetJson(apiURL, map[string]string{}, map[string]string{}, &shodan); err != nil {
		return check.Result{}, check.NewError(err)
	}

	return check.Result{
		Name: "shodan.io",
		Type: check.TypeInfo,
		Info: shodan,
	}, nil
}

type byPort ShodanData

func (x byPort) Len() int           { return len(x) }
func (x byPort) Less(i, j int) bool { return x[i].Port < x[j].Port }
func (x byPort) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

// Info returns interesting information from the check.
func (s Shodan) String() string {
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
			ss := check.NonEmpty(product, version)
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
	return fmt.Sprintf("OS: %s, %d open %s %s", check.Na(s.OS), len(portInfo), portStr, strings.Join(portInfo, ", "))
}

func (s Shodan) JsonString() (string, error) {
	b, err := json.Marshal(s)
	return string(b), err
}
