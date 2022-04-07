package checks

import (
	"net"

	"github.com/jreisinger/checkip/check"
)

type threatCrowd struct {
	Votes int `json:"votes"`
}

// ThreadCrowd threatcrowd.org to find out whether the ipaddr was voted
// malicious.
func ThreadCrowd(ipaddr net.IP) (check.Result, error) {
	result := check.Result{
		Name: "threatcrowd.org",
		Type: check.TypeSec,
		Info: check.EmptyInfo{},
	}

	queryParams := map[string]string{
		"ip": ipaddr.String(),
	}

	// https://github.com/AlienVault-OTX/ApiV2#votes
	// -1 	voted malicious by most users
	// 0 	voted malicious/harmless by equal number of users
	// 1:  	voted harmless by most users
	var threadCrowd threatCrowd
	if err := check.DefaultHttpClient.GetJson("https://www.threatcrowd.org/searchApi/v2/ip/report", map[string]string{}, queryParams, &threadCrowd); err != nil {
		return result, check.NewError(err)
	}
	result.Malicious = threadCrowd.Votes < 0

	return result, nil
}
