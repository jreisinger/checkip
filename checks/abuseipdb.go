package checks

import (
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/jreisinger/checkip/check"
)

// Only return reports within the last x amount of days. Default is 30.
const abuseIPDBMaxAgeInDays = "90"

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
	return fmt.Sprintf("domain: %s, usage type: %s", check.Na(a.Domain), check.Na(a.UsageType))
}

func (d abuseIPDB) JsonString() (string, error) {
	b, err := json.Marshal(d)
	return string(b), err
}

// CheckAbuseIPDB uses api.abuseipdb.com to get generic information about ipaddr
// and see if the ipaddr has been reported as malicious.
func CheckAbuseIPDB(ipaddr net.IP) (check.Result, error) {
	apiKey, err := check.GetConfigValue("ABUSEIPDB_API_KEY")
	if err != nil {
		return check.Result{}, check.NewError(err)
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
	if err := check.DefaultHttpClient.GetJson("https://api.abuseipdb.com/api/v2/check", headers, queryParams, &data); err != nil {
		return check.Result{}, check.NewError(err)
	}

	return check.Result{
		Name:      "abuseipdb.com",
		Type:      check.TypeInfoSec,
		Info:      data.AbuseIPDB,
		Malicious: data.AbuseIPDB.TotalReports > 0 && !data.AbuseIPDB.IsWhitelisted && data.AbuseIPDB.AbuseConfidenceScore > 25,
	}, nil
}
