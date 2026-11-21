package check

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strings"
)

var vpnDetectionUrl = "https://api.vpndetection.io"

// vpnDetection is what VPNDetection knows about an IP address. The proxy
// flags come only with a paid plan's key; a free key leaves them false.
type vpnDetection struct {
	IP         string `json:"ip"`
	IsVPN      bool   `json:"is_vpn"`
	IsTor      bool   `json:"is_tor"`
	IsRelay    bool   `json:"is_relay"`
	IsResProxy bool   `json:"is_resproxy"`
	IsDcProxy  bool   `json:"is_dcproxy"`
	IsMobProxy bool   `json:"is_mobproxy"`
	IsPubProxy bool   `json:"is_pubproxy"`
	VPN        struct {
		ProviderName string `json:"provider_name"`
	} `json:"vpn"`
	Relay struct {
		ProviderName string `json:"provider_name"`
	} `json:"relay"`
	ASN struct {
		ASN  int    `json:"asn"`
		Name string `json:"name"`
	} `json:"asn"`
}

// VPNDetection gets VPN, Tor exit, relay and proxy flags from
// api.vpndetection.io. A free key is enough for the VPN, Tor and relay flags.
func VPNDetection(ipaddr net.IP) (Check, error) {
	result := Check{
		Description: "vpndetection.io",
		Type:        InfoAndIsMalicious,
	}

	apiKey, err := getConfigValue("VPNDETECTION_API_KEY")
	if err != nil {
		return result, newCheckError(err)
	}
	if apiKey == "" {
		result.MissingCredentials = "VPNDETECTION_API_KEY"
		return result, nil
	}

	// Only public addresses have anything to flag.
	if !ipaddr.IsGlobalUnicast() || ipaddr.IsPrivate() {
		return result, nil
	}

	headers := map[string]string{
		"Authorization": "Bearer " + apiKey,
		"Accept":        "application/json",
	}

	var vd vpnDetection
	apiURL := fmt.Sprintf("%s/%s", vpnDetectionUrl, ipaddr)
	if err := defaultHttpClient.GetJson(apiURL, headers, map[string]string{}, &vd); err != nil {
		return result, newCheckError(err)
	}
	if vd.IP == "" {
		return result, newCheckError(errors.New("unexpected response: missing ip"))
	}
	result.IpAddrInfo = vd
	// A relay like iCloud Private Relay hides ordinary users, so it doesn't count.
	result.IpAddrIsMalicious = vd.IsVPN || vd.IsTor || vd.IsResProxy || vd.IsDcProxy || vd.IsMobProxy || vd.IsPubProxy

	return result, nil
}

// Summary returns what the address was flagged as, or that it wasn't.
func (v vpnDetection) Summary() string {
	var flags []string
	if v.IsVPN {
		flags = append(flags, withName("VPN", v.VPN.ProviderName))
	}
	if v.IsTor {
		flags = append(flags, "Tor exit")
	}
	if v.IsRelay {
		flags = append(flags, withName("relay", v.Relay.ProviderName))
	}
	for _, p := range []struct {
		on   bool
		name string
	}{
		{v.IsResProxy, "residential proxy"},
		{v.IsDcProxy, "datacenter proxy"},
		{v.IsMobProxy, "mobile proxy"},
		{v.IsPubProxy, "public proxy"},
	} {
		if p.on {
			flags = append(flags, p.name)
		}
	}
	if len(flags) == 0 {
		return "not flagged"
	}
	return strings.Join(flags, ", ")
}

func (v vpnDetection) Json() ([]byte, error) {
	return json.Marshal(v)
}

func withName(flag, name string) string {
	if name == "" {
		return flag
	}
	return fmt.Sprintf("%s: %s", flag, name)
}
