package check

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"
)

type ipapi struct {
	IP           string `json:"ip"`
	Rir          string `json:"rir"`
	IsBogon      bool   `json:"is_bogon"`
	IsMobile     bool   `json:"is_mobile"`
	IsSatellite  bool   `json:"is_satellite"`
	IsCrawler    bool   `json:"is_crawler"`
	IsDatacenter bool   `json:"is_datacenter"`
	IsTor        bool   `json:"is_tor"`
	IsProxy      bool   `json:"is_proxy"`
	IsVpn        bool   `json:"is_vpn"`
	IsAbuser     bool   `json:"is_abuser"`
	Vpn          struct {
		IsVpn   bool   `json:"is_vpn"`
		Service string `json:"service"`
		URL     string `json:"url"`
	} `json:"vpn"`
}

var ipapiUrl = "https://api.ipapi.is"

// IpAPI gets generic information from api.ipapi.is.
func IpAPI(ipaddr net.IP) (Check, error) {
	result := Check{
		Description: "ipapi.is",
		Type:        InfoAndIsMalicious,
	}

	headers := map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	}

	// optional API_KEY
	apiKey, _ := getConfigValue("IP_API_KEY")
	if apiKey != "" {
		headers["Token"] = apiKey
	}

	var ipapi ipapi
	apiURL := fmt.Sprintf("%s?q=%s", ipapiUrl, ipaddr)
	if apiKey != "" {
		apiURL = fmt.Sprintf("%s?k=%s&q=%s", ipapiUrl, apiKey, ipaddr)
	}
	if err := defaultHttpClient.GetJson(apiURL, headers, map[string]string{}, &ipapi); err != nil {
		return result, newCheckError(err)
	}
	result.IpAddrInfo = ipapi
	if ipapi.IsVpn == true || ipapi.IsTor == true || ipapi.IsAbuser == true {
		result.IpAddrIsMalicious = true
	}

	return result, nil
}

// Info returns interesting information from the check.
func (o ipapi) Summary() string {
	var stype []string

	if o.IsVpn == true {
		if o.Vpn.Service != "" {
			stype = append(stype, fmt.Sprintf("VPN (%s)", o.Vpn.Service))
		} else {
			stype = append(stype, "VPN")
		}

	}

	return fmt.Sprintf("%s", strings.Join(stype, ", "))
}

func (o ipapi) Json() ([]byte, error) {
	return json.Marshal(o)
}
