package check

import (
	"encoding/json"
	"fmt"
	"net"
	"time"
)

// Only return reports within the last x amount of days. Default is 30.
const abuseIPDBMaxAgeInDays = "90"

var abuseIPDBUrl = "https://api.abuseipdb.com/api/v2/check"

// abuseIPDB represents JSON data returned by the AbuseIPDB API.
type abuseIPDB struct {
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

func (a abuseIPDB) Summary() string {
	return fmt.Sprintf("domain: %s, usage type: %s", na(a.Domain), na(a.UsageType))
}

func (a abuseIPDB) Json() ([]byte, error) {
	return json.Marshal(a)
}

// AbuseIPDB uses api.abuseipdb.com to get generic information about ipaddr and
// to see if the ipaddr has been reported as malicious.
func AbuseIPDB(ipaddr net.IP) (Check, error) {
	result := Check{Description: "abuseipdb.com", Type: InfoAndIsMalicious}

	apiKey, err := getConfigValue("ABUSEIPDB_API_KEY")
	if err != nil {
		return result, newCheckError(err)
	}
	if apiKey == "" { // we don't consider missing to be an error
		result.MissingCredentials = "ABUSEIPDB_API_KEY"
		return result, nil
	}

	headers := map[string]string{
		"Key":          apiKey,
		"Accept":       "application/json",
		"Content-Type": "application/x-www-form-urlencoded",
	}

	queryParams := map[string]string{
		"ipAddress":    ipaddr.String(),
		"maxAgeInDays": abuseIPDBMaxAgeInDays,
	}

	var response struct {
		Data abuseIPDB `json:"data"`
	}
	// docs.abuseipdb.com/#check-endpoint
	if err := defaultHttpClient.GetJson(abuseIPDBUrl, headers, queryParams, &response); err != nil {
		return result, newCheckError(err)
	}

	result.IpAddrInfo = response.Data
	result.IpAddrIsMalicious = response.Data.TotalReports > 0 &&
		!response.Data.IsWhitelisted &&
		response.Data.AbuseConfidenceScore > 25

	return result, nil
}
