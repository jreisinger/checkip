package check

import (
	"net"
	"net/url"
	"path"

	"github.com/jreisinger/checkip"
)

var otxUrl = "https://otx.alienvault.com/api/v1/indicators/IPv4"

type otx struct {
	PulseInfo struct {
		Count int `json:"count"`
	} `json:"pulse_info"`
}

// OTX counts pulses on otx.alienvault.com to find out whether the ipaddr is
// malicious.
func OTX(ipaddr net.IP) (checkip.Result, error) {
	result := checkip.Result{
		Name: "otx.alienvault.com",
		Type: checkip.TypeSec,
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
	result.Malicious = otx.PulseInfo.Count > 10

	return result, nil
}
