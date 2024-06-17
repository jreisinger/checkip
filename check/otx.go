package check

import (
	"net"
	"net/url"
	"path"
)

var otxUrl = "https://otx.alienvault.com/api/v1/indicators/IPv4"

type otx struct {
	PulseInfo struct {
		Count int `json:"count"`
	} `json:"pulse_info"`
}

// OTX counts pulses to find out whether the ipaddr is malicious. Is uses
// https://otx.alienvault.com/api/v1/indicators/IPv4.
func OTX(ipaddr net.IP) (Check, error) {
	result := Check{
		Description: "otx.alienvault.com",
		Type:        IsMalicious,
	}

	u, err := url.Parse(otxUrl)
	if err != nil {
		return result, newCheckError(err)
	}
	u.Path = path.Join(u.Path, ipaddr.String())

	var otx otx
	if err := defaultHttpClient.GetJson(u.String(), map[string]string{}, map[string]string{}, &otx); err != nil {
		return result, newCheckError(err)
	}
	result.IpAddrIsMalicious = otx.PulseInfo.Count > 10

	return result, nil
}
