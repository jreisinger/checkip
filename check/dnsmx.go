package check

import (
	"bytes"
	"encoding/json"
	"net"
	"strings"
)

// mx maps mx records to domain names.
type mx struct {
	Records map[string][]string `json:"records"` // domain => MX records
}

func (mx mx) Summary() string {
	var s string
	for domain, mxRecords := range mx.Records {
		if domain == "" && len(mxRecords) == 0 {
			continue
		}
		for i := range mxRecords {
			mxRecords[i] = strings.TrimSuffix(mxRecords[i], ".")
		}
		s += domain + ": " + strings.Join(mxRecords, ", ")
	}
	return na(s)
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

	names, _ := net.LookupAddr(ipaddr.String()) // NOTE: ignoring error

	// Enrich names with a name with 'www.' removed.
	// [www.csh.ac.at.] => [www.csh.ac.at. csh.ac.at.]
	for _, n := range names {
		t := strings.TrimPrefix(n, "www.")
		names = append(names, t)
	}

	// Enrich names with a domain name from AbuseIPDB.
	// [www.csh.ac.at. csh.ac.at.] = > [www.csh.ac.at. csh.ac.at. aco.net]
	r, err := AbuseIPDB(ipaddr)
	if err != nil {
		return result, newCheckError(err)
	}
	if r.IpAddrInfo == nil {
		return result, nil
	}
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
	names = append(names, a.Domain)

	var mx mx
	for _, n := range names {
		var mxRecords2 []string
		mxRecords, _ := net.LookupMX(n) // NOTE: ingoring error
		for _, r := range mxRecords {
			mxRecords2 = append(mxRecords2, r.Host)
		}
		if _, ok := mx.Records[n]; !ok {
			mx.Records = make(map[string][]string)
		}
		mx.Records[n] = mxRecords2
	}
	result.IpAddrInfo = mx

	return result, nil
}
