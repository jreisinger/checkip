package check

import (
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/jreisinger/checkip"
)

// Only return reports within the last x amount of days. Default is 30.
const abuseIPDBMaxAgeInDays = "90"

var abuseIPDBUrl = "https://api.abuseipdb.com/api/v2/check"

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

func (a abuseIPDB) JsonString() (string, error) {
	b, err := json.Marshal(a)
	return string(b), err
}

// AbuseIPDB uses https://api.abuseipdb.com/ to get generic information about
// ipaddr and see if the ipaddr has been reported as malicious.
func AbuseIPDB(ipaddr net.IP) (checkip.Result, error) {
	result := checkip.Result{Name: "abuseipdb.com", Type: checkip.TypeInfoSec}

	apiKey, err := getConfigValue("ABUSEIPDB_API_KEY")
	if err != nil {
		return result, newCheckError(err)
	}
	if apiKey == "" {
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

	var data struct {
		AbuseIPDB abuseIPDB `json:"data"`
	}
	// docs.abuseipdb.com/#check-endpoint
	if err := defaultHttpClient.GetJson(abuseIPDBUrl, headers, queryParams, &data); err != nil {
		return result, newCheckError(err)
	}

	result.Info = data.AbuseIPDB
	result.Malicious = data.AbuseIPDB.TotalReports > 0 &&
		!data.AbuseIPDB.IsWhitelisted &&
		data.AbuseIPDB.AbuseConfidenceScore > 25

	return result, nil

}
