package checkip

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
)

// Shodan holds information about an IP address from shodan.io scan data.
type Shodan struct {
	Org   string `json:"org"`
	Data  data   `json:"data"`
	Os    string `json:"os"`
	Ports []int  `json:"ports"`
}

type data []struct {
	Product string `json:"product"`
	Version string `json:"version"`
	Port    int    `json:"port"`
}

// Check fills in Shodan data for a given IP address. Its get the data from
// https://api.shodan.io
func (s *Shodan) Check(ipaddr net.IP) (bool, error) {
	apiKey, err := getConfigValue("SHODAN_API_KEY")
	if err != nil {
		return false, fmt.Errorf("can't call API: %w", err)
	}

	resp, err := http.Get(fmt.Sprintf("https://api.shodan.io/shodan/host/%s?key=%s", ipaddr, apiKey))
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&s); err != nil {
		return false, err
	}

	return s.isOK(), nil
}

func (s *Shodan) isOK() bool {
	return !s.gotServiceVersion()
}

func (s *Shodan) gotServiceVersion() bool {
	for _, d := range s.Data {
		if d.Version != "" {
			return true
		}
	}
	return false
}

// String returns the result of the check.
func (s *Shodan) String() string {
	os := "OS unknown"
	if s.Os != "" {
		os = s.Os
	}

	var portInfo []string
	for _, d := range s.Data {
		product := "service unknown"
		if d.Product != "" {
			product = d.Product
		}

		version := "version unknown"
		if d.Version != "" {
			version = d.Version
		}

		portInfo = append(portInfo, fmt.Sprintf("%d (%s, %s)", d.Port, product, version))
	}

	portStr := "port"
	if len(portInfo) != 1 {
		portStr += "s"
	}
	if len(portInfo) > 0 {
		portStr += ":"
	}

	return fmt.Sprintf("%s, %d open %s %s", os, len(portInfo), portStr, strings.Join(portInfo, ", "))
}

func joinPortData(ds data) string {
	var ss []string
	for _, d := range ds {
		s := fmt.Sprintf("%d (%s %s)", d.Port, d.Product, d.Version)
		ss = append(ss, s)
	}
	return strings.Join(ss, ", ")
}

func joinInts(ints []int) string {
	var ss []string
	for _, i := range ints {
		a := strconv.Itoa(i)
		ss = append(ss, a)
	}
	return strings.Join(ss, ", ")
}
