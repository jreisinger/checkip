package checkip

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"
)

// AbuseIPDB holds information about an IP address from abuseipdb.com database.
type AbuseIPDB struct {
	Data struct {
		AbuseConfidenceScore int `json:"abuseConfidenceScore"`
		TotalReports         int `json:"totalReports"`
	} `json:"data"`
}

// Check fills in AbuseIPDB data for a given IP address. Its get the data from
// api.abuseipdb.com/api/v2/check (docs.abuseipdb.com/#check-endpoint).
func (a *AbuseIPDB) Check(ipaddr net.IP) (bool, error) {
	apiKey, err := getConfigValue("ABUSEIPDB_API_KEY")
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

	client := newHTTPClient(5 * time.Second)
	resp, err := client.Do(req)
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

	return a.isOK(), nil
}

func (a *AbuseIPDB) isOK() bool {
	return a.Data.AbuseConfidenceScore <= 25
}

// String returns the result of the check.
func (a *AbuseIPDB) String() string {
	return fmt.Sprintf("reported abusive %d times with %d%% confidence",
		a.Data.TotalReports,
		a.Data.AbuseConfidenceScore,
	)
}
