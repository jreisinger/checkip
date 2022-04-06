package checks

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/jreisinger/checkip/check"
)

const days = 30 // limit search to last 30 days

// UrlScan gets data from urlscan.io. When a URL is submitted to urlscan.io, an
// automated process will browse to the URL like a regular user and record the
// activity that this page navigation creates.
func UrlScan(ipaddr net.IP) (check.Result, error) {
	apiKey, err := check.GetConfigValue("URLSCAN_API_KEY")
	if err != nil {
		return check.Result{}, check.NewError(err)
	}
	if apiKey == "" {
		return check.Result{}, nil
	}

	url := "https://urlscan.io/api/v1/search"
	headers := map[string]string{
		"API-Key":      apiKey,
		"Content-Type": "application/json",
	}
	queryParams := map[string]string{
		"q": fmt.Sprintf("page.ip:%s AND date:>now-%dd", ipaddr, days),
	}

	var u urlscan

	if err := check.DefaultHttpClient.GetJson(url, headers, queryParams, &u); err != nil {
		return check.Result{}, check.NewError(err)
	}

	var maliciousVerdicts int

	for _, r := range u.Results {
		var ur urlscanResult
		err := check.DefaultHttpClient.GetJson(r.Result, headers, map[string]string{}, &ur)
		if err != nil {
			return check.Result{}, check.NewError(err)
		}
		if ur.Verdicts.Overall.Malicious {
			maliciousVerdicts++
		}
		// time.Sleep(time.Millisecond * 100)
	}

	return check.Result{
		Name:      "urlscan.io",
		Type:      check.TypeInfoSec,
		Info:      u,
		Malicious: float64(maliciousVerdicts)/float64(len(u.Results)) > 0.1,
	}, nil
}

type urlscan struct {
	Results []struct {
		IndexedAt time.Time `json:"indexedAt"`
		Page      struct {
			IP       string `json:"ip"`
			MimeType string `json:"mimeType"`
			URL      string `json:"url"`
			Status   string `json:"status"`
		} `json:"page"`
		Result     string `json:"result"`
		Screenshot string `json:"screenshot"`
	} `json:"results"`
}

type urlscanResult struct {
	Verdicts struct {
		Overall struct {
			Malicious bool `json:"malicious"`
		} `json:"overall"`
	} `json:"verdicts"`
}

// Strings tells how many scanned URLs are associated with the IP address.
func (u urlscan) Summary() string {
	urlCnt := make(map[string]int)
	for _, r := range u.Results {
		urlCnt[r.Page.URL]++
	}

	var urls []string
	for url := range urlCnt {
		urls = append(urls, url)
	}

	switch len(urls) {
	case 0:
		return "0 related URLs"
	case 1:
		return fmt.Sprintf("1 related URL: %s", urls[0])
	default:
		return fmt.Sprintf("%d related URLs: %s", len(urls), strings.Join(urls, ", "))
	}
}

func (u urlscan) JsonString() (string, error) {
	b, err := json.Marshal(u)
	return string(b), err
}
