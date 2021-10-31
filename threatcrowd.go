package checkip

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
)

// ThreatCrowd holds information about an IP address from threatcrowd.org.
type ThreatCrowd struct {
	Votes int `json:"votes"`
}

// Check retrieves information from
// https://www.threatcrowd.org/searchApi/v2/ip/report.
func (t *ThreatCrowd) Check(ipaddr net.IP) error {
	queryParams := map[string]string{
		"ip": ipaddr.String(),
	}
	resp, err := makeAPIcall("https://www.threatcrowd.org/searchApi/v2/ip/report", map[string]string{}, queryParams)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("search threatcrowd failed: %s", resp.Status)
	}

	if err := json.NewDecoder(resp.Body).Decode(t); err != nil {
		return err
	}

	return nil
}

func (t *ThreatCrowd) IsMalicious() bool {
	// https://github.com/AlienVault-OTX/ApiV2#votes
	// -1 	voted malicious by most users
	// 0 	voted malicious/harmless by equal number of users
	// 1:  	voted harmless by most users
	return t.Votes < 0
}
