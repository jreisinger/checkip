package abuseipdb

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"
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

// New creates AS.
func New() *AbuseIPDB {
	return &AbuseIPDB{}
}

// ForIP fills in AbuseIPDB data for a given IP address. See the AbuseIPDB API
// documentation for more https://docs.abuseipdb.com/?shell#check-endpoint
func (a *AbuseIPDB) ForIP(ipaddr net.IP) error {
	apiKey := os.Getenv("ABUSEIPDB_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("can't call API: environment variable ABUSEIPDB_API_KEY is not set")
	}

	baseURL, err := url.Parse("https://api.abuseipdb.com/api/v2/check")
	if err != nil {
		return err
	}

	// Add GET paramaters.
	params := url.Values{}
	params.Add("ipAddress", ipaddr.String())
	baseURL.RawQuery = params.Encode()

	req, err := http.NewRequest("GET", baseURL.String(), nil)
	if err != nil {
		return err
	}

	// Set request headers.
	req.Header.Set("Key", apiKey)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
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
