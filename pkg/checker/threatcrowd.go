package checker

import (
	"github.com/jreisinger/checkip/pkg/check"
	"net"
)

// ThreatCrowd holds information about an IP address from threatcrowd.org.
type ThreatCrowd struct {
	Votes int `json:"votes"`
}

// CheckThreadCrowd retrieves information from
// https://www.threatcrowd.org/searchApi/v2/ip/report.
func CheckThreadCrowd(ipaddr net.IP) check.Result {
	queryParams := map[string]string{
		"ip": ipaddr.String(),
	}

	// https://github.com/AlienVault-OTX/ApiV2#votes
	// -1 	voted malicious by most users
	// 0 	voted malicious/harmless by equal number of users
	// 1:  	voted harmless by most users
	var threadCrowd ThreatCrowd
	if err := check.DefaultHttpClient.GetJson("https://www.threatcrowd.org/searchApi/v2/ip/report", map[string]string{}, queryParams, &threadCrowd); err != nil {
		return check.Result{Error: check.NewResultError(err)}
	}

	return check.Result{
		Name:        "threatcrowd.org",
		Type:        check.TypeSec,
		Data:        check.EmptyData{},
		IsMalicious: threadCrowd.Votes < 0,
	}
}
