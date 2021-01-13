package check

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/jreisinger/checkip/util"
)

// AbuseIPDB holds information about an IP address from abuseipdb.com database.
type AbuseIPDB struct {
	Data struct {
		IPAddress            string        `json:"ipAddress"`
		IsPublic             bool          `json:"isPublic"`
		IPVersion            int           `json:"ipVersion"`
		IsWhitelisted        bool          `json:"isWhitelisted"`
		AbuseConfidenceScore int           `json:"abuseConfidenceScore"`
		CountryCode          string        `json:"countryCode"`
		UsageType            string        `json:"usageType"`
		Isp                  string        `json:"isp"`
		Domain               string        `json:"domain"`
		Hostnames            []interface{} `json:"hostnames"`
		CountryName          string        `json:"countryName"`
		TotalReports         int           `json:"totalReports"`
		NumDistinctUsers     int           `json:"numDistinctUsers"`
		LastReportedAt       time.Time     `json:"lastReportedAt"`
		Reports              []struct {
			ReportedAt          time.Time `json:"reportedAt"`
			Comment             string    `json:"comment"`
			Categories          []int     `json:"categories"`
			ReporterID          int       `json:"reporterId"`
			ReporterCountryCode string    `json:"reporterCountryCode"`
			ReporterCountryName string    `json:"reporterCountryName"`
		} `json:"reports"`
	} `json:"data"`
}

// Do fills in AbuseIPDB data for a given IP address. See the AbuseIPDB API
// documentation for more https://docs.abuseipdb.com/?shell#check-endpoint
func (a *AbuseIPDB) Do(ipaddr net.IP) (bool, error) {
	apiKey, err := util.GetConfigValue("ABUSEIPDB_API_KEY")
	if err != nil {
		return false, fmt.Errorf("can't call API: %w", err)
	}

	baseURL, err := url.Parse("https://api.abuseipdb.com/api/v2/check")
	if err != nil {
		return false, err
	}

	// Add GET paramaters.
	params := url.Values{}
	params.Add("ipAddress", ipaddr.String())
	baseURL.RawQuery = params.Encode()

	req, err := http.NewRequest("GET", baseURL.String(), nil)
	if err != nil {
		return false, err
	}

	// Set request headers.
	req.Header.Set("Key", apiKey)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("calling API: %s", resp.Status)
	}

	if err := json.NewDecoder(resp.Body).Decode(a); err != nil {
		return false, err
	}
	if a.Data.AbuseConfidenceScore > 0 {
		return false, nil
	}

	return true, nil
}

// Name returns the name of the check.
func (a *AbuseIPDB) Name() string {
	return fmt.Sprint("AbuseIPDB")
}

// String returns the result of the check.
func (a *AbuseIPDB) String() string {
	return fmt.Sprintf("malicious with %d%% confidence | %v", a.Data.AbuseConfidenceScore, a.Data.Domain)
}
