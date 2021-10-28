package checkip

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
)

// Only return reports within the last x amount of days. Default is 30.
var maxAgeInDays = "90"

// AbuseIPDB holds information about an IP address from abuseipdb.com database.
type AbuseIPDB struct {
	Data struct {
		AbuseConfidenceScore int  `json:"abuseConfidenceScore"`
		TotalReports         int  `json:"totalReports"`
		IsWhitelisted        bool `json:"isWhitelisted"`
	} `json:"data"`
}

// Check fills in AbuseIPDB data for a given IP address. Its get the data from
// api.abuseipdb.com/api/v2/check (docs.abuseipdb.com/#check-endpoint). It
// returns false if the IP address is not whitelisted and AbuseConfidenceScore >
// 25.
func (a *AbuseIPDB) Check(ipaddr net.IP) (bool, error) {
	apiKey, err := getConfigValue("ABUSEIPDB_API_KEY")
	if err != nil {
		return true, fmt.Errorf("can't call API: %w", err)
	}

	headers := map[string]string{
		"Key":          apiKey,
		"Accept":       "application/json",
		"Content-Type": "application/x-www-form-urlencoded",
	}

	queryParams := map[string]string{
		"ipAddress":    ipaddr.String(),
		"maxAgeInDays": maxAgeInDays,
	}

	resp, err := makeAPIcall("https://api.abuseipdb.com/api/v2/check", headers, queryParams)
	if err != nil {
		return true, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return true, fmt.Errorf("calling API: %s", resp.Status)
	}

	if err := json.NewDecoder(resp.Body).Decode(a); err != nil {
		return true, err
	}

	return a.isOK(), nil
}

func (a *AbuseIPDB) isOK() bool {
	return a.Data.TotalReports == 0 || a.Data.IsWhitelisted || a.Data.AbuseConfidenceScore <= 25
}

// String returns the result of the check.
func (a *AbuseIPDB) String() string {
	return fmt.Sprintf("%d reports in last %s days (abuse confidence score: %d%%, whitelisted: %v)",
		a.Data.TotalReports,
		maxAgeInDays,
		a.Data.AbuseConfidenceScore,
		a.Data.IsWhitelisted,
	)
}
