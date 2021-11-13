package checks

import (
	"fmt"
	"net"

	"github.com/jreisinger/checkip/check"
)

type otx struct {
	PulseInfo struct {
		Count int `json:"count"`
	} `json:"pulse_info"`
}

// CheckOTX counts pulses on otx.alienvault.com to find out whether the ipaddr
// is malicious.
func CheckOTX(ipaddr net.IP) (check.Result, *check.Error) {
	apiurl := fmt.Sprintf("https://otx.alienvault.com/api/v1/indicators/IPv4/%s/", ipaddr.String())

	var otx otx
	if err := check.DefaultHttpClient.GetJson(apiurl, map[string]string{}, map[string]string{}, &otx); err != nil {
		return check.Result{}, check.NewError(err)
	}

	return check.Result{
		Name:            "otx.alienvault.com",
		Type:            check.TypeSec,
		Info:            check.EmptyInfo{},
		IPaddrMalicious: otx.PulseInfo.Count > 10,
	}, nil
}
