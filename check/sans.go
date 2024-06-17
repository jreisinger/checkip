package check

import (
	"encoding/json"
	"fmt"
	"net"
	"net/url"
	"path"
)

const sansUrl = "https://isc.sans.edu/api/ip/"

type sans struct {
	Ip struct {
		Count          int    // (also reports or records) total number of packets blocked from this IP
		Attacks        int    // (also targets) number of unique destination IP addresses for these packets
		AsAbuseContact string `json:"asabusecontact"`
	}
}

func (s sans) Summary() string {
	return fmt.Sprintf("attacks: %d, abuse contact: %s", s.Ip.Attacks, s.Ip.AsAbuseContact)
}

func (s sans) Json() ([]byte, error) {
	return json.Marshal(s)
}

// SansISC gets info from SANS Internet Storm Center API.
// curl "https://isc.sans.edu/api/ip/${IPADDR}?json"
func SansISC(ipaddr net.IP) (Check, error) {
	result := Check{
		Description: "isc.sans.edu",
		Type:        InfoAndIsMalicious,
	}

	u, err := url.Parse(sansUrl)
	if err != nil {
		return result, newCheckError(err)
	}

	u.Path = path.Join(u.Path, ipaddr.String())

	var sans sans
	if err := defaultHttpClient.GetJson(u.String(), map[string]string{}, map[string]string{"json": ""}, &sans); err != nil {
		return result, newCheckError(err)
	}

	result.IpAddrInfo = sans
	result.IpAddrIsMalicious = sans.Ip.Attacks > 0

	return result, nil
}
