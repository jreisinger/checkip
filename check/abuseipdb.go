package check

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/jreisinger/checkip/util"
)

// AbuseIPDB holds information about an IP address from abuseipdb.com database.
type AbuseIPDB struct {
	Data struct {
		AbuseConfidenceScore int    `json:"abuseConfidenceScore"`
		UsageType            string `json:"usageType"`
		Domain               string `json:"domain"`
		TotalReports         int    `json:"totalReports"`
	} `json:"data"`
}

// Do fills in AbuseIPDB data for a given IP address. Its get the data from
// https://api.abuseipdb.com/api/v2/check
// (https://docs.abuseipdb.com/#check-endpoint).
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

	client := NewHTTPClient(5 * time.Second)
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

	if a.isNotOK() {
		return false, nil
	}

	return true, nil
}

func (a *AbuseIPDB) isNotOK() bool {
	return a.Data.AbuseConfidenceScore > 25
}

// Name returns the name of the check.
func (a *AbuseIPDB) Name() string {
	return fmt.Sprint("AbuseIPDB")
}

// String returns the result of the check.
func (a *AbuseIPDB) String() string {
	confidence := strconv.Itoa(a.Data.AbuseConfidenceScore)
	if a.isNotOK() {
		confidence = fmt.Sprintf("%s", util.Highlight(confidence+"%"))
	} else {
		confidence = fmt.Sprintf("%s", (confidence + "%"))
	}
	return fmt.Sprintf("reported abusive %d times with %s confidence (%s)",
		a.Data.TotalReports,
		confidence,
		a.Data.Domain,
	)
}
