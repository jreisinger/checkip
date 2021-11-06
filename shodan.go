package checkip

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"sort"
	"strings"
)

// Shodan holds information about an IP address from shodan.io scan data.
type Shodan struct {
	Org   string `json:"org"`
	Data  data   `json:"data"`
	OS    string `json:"os"`
	Ports []int  `json:"ports"`
}

type data []struct {
	Product   string `json:"product"`
	Version   string `json:"version"`
	Port      int    `json:"port"`
	Transport string `json:"transport"` // tcp, udp
}

func (s *Shodan) String() string { return "shodan.io" }

// Check fills in Shodan data for a given IP address. Its get the data from
// https://api.shodan.io.
func (s *Shodan) Check(ipaddr net.IP) error {
	apiKey, err := getConfigValue("SHODAN_API_KEY")
	if err != nil {
		return fmt.Errorf("can't call API: %w", err)
	}

	apiURL := fmt.Sprintf("https://api.shodan.io/shodan/host/%s?key=%s", ipaddr, apiKey)
	resp, err := makeAPIcall(apiURL, map[string]string{}, map[string]string{})
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// StatusNotFound is returned when shodan doesn't know the IP address.
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNotFound {
		return fmt.Errorf("calling %s: %s", apiURL, resp.Status)
	}

	if err := json.NewDecoder(resp.Body).Decode(&s); err != nil {
		return err
	}

	return nil
}

type byPort data

func (x byPort) Len() int           { return len(x) }
func (x byPort) Less(i, j int) bool { return x[i].Port < x[j].Port }
func (x byPort) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

// Info returns interesting information from the check.
func (s *Shodan) Info() string {
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
			portInfo = append(portInfo, fmt.Sprintf("%s/%d (%s, %s)", d.Transport, d.Port, product, version))
		}
	}

	portStr := "port"
	if len(portInfo) != 1 {
		portStr += "s"
	}
	if len(portInfo) > 0 {
		portStr += ":"
	}

	return fmt.Sprintf("OS: %s, %d open %s %s", na(s.OS), len(portInfo), portStr, strings.Join(portInfo, ", "))
}
