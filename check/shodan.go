package check

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/jreisinger/checkip/util"
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

// Do fills in Shodan data for a given IP address. Its get the data from
// https://api.shodan.io
func (s *Shodan) Do(ipaddr net.IP) (bool, error) {
	apiKey, err := util.GetConfigValue("SHODAN_API_KEY")
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

	return true, nil
}

// Name returns the name of the check.
func (s *Shodan) Name() string {
	return fmt.Sprint("Shodan")
}

// String returns the result of the check.
func (s *Shodan) String() string {
	if s.Os == "" {
		s.Os = "n/a"
	}
	return fmt.Sprintf("OS: %s, open ports: %s", s.Os, joinData(s.Data))
}

func joinData(ds data) string {
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
