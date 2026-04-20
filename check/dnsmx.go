package check

import (
	"bytes"
	"encoding/json"
	"net"
	"sort"
	"strings"
)

var dnsLookupAddr = net.LookupAddr
var dnsLookupMX = net.LookupMX
var abuseIPDBCheck = AbuseIPDB

// mx maps mx records to domain names.
type mx struct {
	Records map[string][]string `json:"records"` // domain => MX records
}

func (mx mx) Summary() string {
	var records []string
	var domains []string
	for domain := range mx.Records {
		domains = append(domains, domain)
	}
	sort.Strings(domains)

	for _, domain := range domains {
		mxRecords := mx.Records[domain]
		if domain == "" || len(mxRecords) == 0 {
			continue
		}
		for i := range mxRecords {
			mxRecords[i] = trimTrailingDot(mxRecords[i])
		}
		records = append(records, domain+": "+strings.Join(mxRecords, ", "))
	}
	return strings.Join(records, ", ")
}

func (mx mx) Json() ([]byte, error) {
	return json.Marshal(mx)
}

// DnsMX performs reverse lookup and consults AbuseIPDB to get domain names fo
// the ipaddr. Then it looks up MX records (mail servers) for the given domain
// names.
func DnsMX(ipaddr net.IP) (Check, error) {
	result := Check{
		Description: "dns mx",
		Type:        Info,
	}

	names, _ := dnsLookupAddr(ipaddr.String()) // NOTE: ignoring error
	for i := range names {
		names[i] = trimTrailingDot(names[i])
	}

	// Enrich names with a name with 'www.' removed.
	// [www.csh.ac.at] => [www.csh.ac.at csh.ac.at]
	for _, n := range names {
		t := strings.TrimPrefix(n, "www.")
		names = append(names, t)
	}

	// Enrich names with a domain name from AbuseIPDB.
	// [www.csh.ac.at. csh.ac.at.] = > [www.csh.ac.at. csh.ac.at. aco.net]
	r, err := abuseIPDBCheck(ipaddr)
	if err != nil {
		return result, newCheckError(err)
	}
	if r.IpAddrInfo != nil {
		j, err := r.IpAddrInfo.Json()
		if err != nil {
			return result, newCheckError(err)
		}
		b := bytes.NewReader(j)
		decoder := json.NewDecoder(b)
		var a abuseIPDB
		if err := decoder.Decode(&a); err != nil {
			return result, newCheckError(err)
		}
		names = append(names, trimTrailingDot(a.Domain))
	}

	var mx mx
	for _, n := range names {
		if n == "" {
			continue
		}
		var mxRecords2 []string
		mxRecords, _ := dnsLookupMX(n) // NOTE: ingoring error
		for _, r := range mxRecords {
			mxRecords2 = append(mxRecords2, trimTrailingDot(r.Host))
		}
		if len(mxRecords2) == 0 {
			continue
		}
		if mx.Records == nil {
			mx.Records = make(map[string][]string)
		}
		mx.Records[n] = mxRecords2
	}
	result.IpAddrInfo = mx

	return result, nil
}
