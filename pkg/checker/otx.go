package checker

import (
	"fmt"
	"github.com/jreisinger/checkip/pkg/check"
	"net"
)

// OTX holds information from otx.alienvault.com.
type OTX struct {
	PulseInfo struct {
		Count int `json:"count"`
	} `json:"pulse_info"`
}

// CheckOTX gets data from https://otx.alienvault.com/api.
func CheckOTX(ipaddr net.IP) check.Result {
	apiurl := fmt.Sprintf("https://otx.alienvault.com/api/v1/indicators/IPv4/%s/", ipaddr.String())

	var otx OTX
	if err := check.DefaultHttpClient.GetJson(apiurl, map[string]string{}, map[string]string{}, &otx); err != nil {
		return check.Result{Error: check.NewResultError(err)}
	}

	return check.Result{
		Name:        "otx.alienvault.com",
		Type:        check.TypeSec,
		Data:        check.EmptyData{},
		IsMalicious: otx.PulseInfo.Count > 10,
	}
}
