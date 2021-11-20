package checks

import (
	"net"
	"net/url"
	"path"

	"github.com/jreisinger/checkip/check"
)

var otxUrl = "https://otx.alienvault.com/api/v1/indicators/IPv4"

type otx struct {
	PulseInfo struct {
		Count int `json:"count"`
	} `json:"pulse_info"`
}

// CheckOTX counts pulses on otx.alienvault.com to find out whether the ipaddr
// is malicious.
func CheckOTX(ipaddr net.IP) (check.Result, error) {
	u, err := url.Parse(otxUrl)
	if err != nil {
		return check.Result{}, check.NewError(err)
	}
	u.Path = path.Join(u.Path, ipaddr.String())

	var otx otx
	if err := check.DefaultHttpClient.GetJson(u.String(), map[string]string{}, map[string]string{}, &otx); err != nil {
		return check.Result{}, check.NewError(err)
	}

	return check.Result{
		Name:      "otx.alienvault.com",
		Type:      check.TypeSec,
		Info:      check.EmptyInfo{},
		Malicious: otx.PulseInfo.Count > 10,
	}, nil
}
