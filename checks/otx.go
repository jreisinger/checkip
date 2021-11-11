package checks

import (
	"fmt"
	"net"

	"github.com/jreisinger/checkip/check"
)

// OTX holds information from otx.alienvault.com.
type OTX struct {
	PulseInfo struct {
		Count int `json:"count"`
	} `json:"pulse_info"`
}

func CheckOTX(ipaddr net.IP) (check.Result, error) {
	apiurl := fmt.Sprintf("https://otx.alienvault.com/api/v1/indicators/IPv4/%s/", ipaddr.String())

	var otx OTX
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
