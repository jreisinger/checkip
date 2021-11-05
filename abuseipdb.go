package checkip

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"
)

// Only return reports within the last x amount of days. Default is 30.
var maxAgeInDays = "90"

// AbuseIPDB holds information about an IP address from abuseipdb.com database.
type AbuseIPDB struct {
	Data Data `json:"data"`
}
type Data struct {
	IPAddress            string        `json:"ipAddress"`
	IsPublic             bool          `json:"isPublic"`
	IPVersion            int           `json:"ipVersion"`
	IsWhitelisted        bool          `json:"isWhitelisted"`
	AbuseConfidenceScore int           `json:"abuseConfidenceScore"`
	CountryCode          string        `json:"countryCode"`
	CountryName          string        `json:"countryName"`
	UsageType            string        `json:"usageType"`
	Isp                  string        `json:"isp"`
	Domain               string        `json:"domain"`
	Hostnames            []interface{} `json:"hostnames"`
	TotalReports         int           `json:"totalReports"`
	NumDistinctUsers     int           `json:"numDistinctUsers"`
	LastReportedAt       time.Time     `json:"lastReportedAt"`
}

func (a *AbuseIPDB) Name() string { return "abuseipdb.com" }

// Check fills in AbuseIPDB data for a given IP address. It gets the data from
// api.abuseipdb.com/api/v2/check (docs.abuseipdb.com/#check-endpoint).
func (a *AbuseIPDB) Check(ipaddr net.IP) error {
	apiKey, err := getConfigValue("ABUSEIPDB_API_KEY")
	if err != nil {
		return fmt.Errorf("can't call API: %w", err)
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
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("calling API: %s", resp.Status)
	}

	if err := json.NewDecoder(resp.Body).Decode(a); err != nil {
		return err
	}

	return nil
}

func (a *AbuseIPDB) IsMalicious() bool {
	return a.Data.TotalReports > 0 && !a.Data.IsWhitelisted && a.Data.AbuseConfidenceScore > 25
}
